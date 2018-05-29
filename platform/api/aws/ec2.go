// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/coreos/mantle/util"
)

func (a *API) AddKey(name, key string) error {
	_, err := a.ec2.ImportKeyPair(&ec2.ImportKeyPairInput{
		KeyName:           &name,
		PublicKeyMaterial: []byte(key),
	})

	return err
}

func (a *API) DeleteKey(name string) error {
	_, err := a.ec2.DeleteKeyPair(&ec2.DeleteKeyPairInput{
		KeyName: &name,
	})

	return err
}

// CreateInstances creates EC2 instances with a given name tag, optional ssh key name, user data. The image ID, instance type, and security group set in the API will be used. CreateInstances will block until all instances are running and have an IP address.
func (a *API) CreateInstances(name, keyname, userdata string, count uint64) ([]*ec2.Instance, error) {
	cnt := int64(count)

	var ud *string
	if len(userdata) > 0 {
		tud := base64.StdEncoding.EncodeToString([]byte(userdata))
		ud = &tud
	}

	err := a.ensureInstanceProfile(a.opts.IAMInstanceProfile)
	if err != nil {
		return nil, fmt.Errorf("error verifying IAM instance profile: %v", err)
	}

	sgId, err := a.getSecurityGroupID(a.opts.SecurityGroup)
	if err != nil {
		return nil, fmt.Errorf("error resolving security group: %v", err)
	}

	vpcId, err := a.getVPCID(sgId)
	if err != nil {
		return nil, fmt.Errorf("error resolving vpc: %v", err)
	}

	subnetId, err := a.getSubnetID(vpcId)
	if err != nil {
		return nil, fmt.Errorf("error resolving subnet: %v", err)
	}

	key := &keyname
	if keyname == "" {
		key = nil
	}
	inst := ec2.RunInstancesInput{
		ImageId:          &a.opts.AMI,
		MinCount:         &cnt,
		MaxCount:         &cnt,
		KeyName:          key,
		InstanceType:     &a.opts.InstanceType,
		SecurityGroupIds: []*string{&sgId},
		SubnetId:         &subnetId,
		UserData:         ud,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: &a.opts.IAMInstanceProfile,
		},
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String(ec2.ResourceTypeInstance),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String(name),
					},
					&ec2.Tag{
						Key:   aws.String("CreatedBy"),
						Value: aws.String("mantle"),
					},
				},
			},
		},
	}

	reservations, err := a.ec2.RunInstances(&inst)
	if err != nil {
		return nil, fmt.Errorf("error running instances: %v", err)
	}

	ids := make([]string, len(reservations.Instances))
	for i, inst := range reservations.Instances {
		ids[i] = *inst.InstanceId
	}

	// loop until all machines are online
	var insts []*ec2.Instance

	// 10 minutes is a pretty reasonable timeframe for AWS instances to work.
	timeout := 10 * time.Minute
	// don't make api calls too quickly, or we will hit the rate limit
	delay := 10 * time.Second
	err = util.WaitUntilReady(timeout, delay, func() (bool, error) {
		desc, err := a.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice(ids),
		})
		if err != nil {
			return false, err
		}
		insts = desc.Reservations[0].Instances

		for _, i := range insts {
			if *i.State.Name != ec2.InstanceStateNameRunning || i.PublicIpAddress == nil {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		a.TerminateInstances(ids)
		return nil, fmt.Errorf("waiting for instances to run: %v", err)
	}

	return insts, nil
}

// gcEC2 will terminate ec2 instances older than gracePeriod.
// It will only operate on ec2 instances tagged with 'mantle' to avoid stomping
// on other resources in the account.
func (a *API) gcEC2(gracePeriod time.Duration) error {
	durationAgo := time.Now().Add(-1 * gracePeriod)

	instances, err := a.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:CreatedBy"),
				Values: aws.StringSlice([]string{"mantle"}),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error describing instances: %v", err)
	}

	toTerminate := []string{}

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			if instance.LaunchTime.After(durationAgo) {
				plog.Debugf("ec2: skipping instance %s due to being too new", *instance.InstanceId)
				// Skip, still too new
				continue
			}

			if instance.State != nil {
				switch *instance.State.Name {
				case ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning, ec2.InstanceStateNameStopped:
					toTerminate = append(toTerminate, *instance.InstanceId)
				case ec2.InstanceStateNameTerminated, ec2.InstanceStateNameShuttingDown:
				default:
					plog.Infof("ec2: skipping instance in state %s", *instance.State.Name)
				}
			} else {
				plog.Warningf("ec2 instance had no state: %s", *instance.InstanceId)
			}
		}
	}

	return a.TerminateInstances(toTerminate)
}

// TerminateInstances schedules EC2 instances to be terminated.
func (a *API) TerminateInstances(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	input := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(ids),
	}

	if _, err := a.ec2.TerminateInstances(input); err != nil {
		return err
	}

	return nil
}

func (a *API) CreateTags(resources []string, tags map[string]string) error {
	tagObjs := make([]*ec2.Tag, 0, len(tags))
	for key, value := range tags {
		tagObjs = append(tagObjs, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}
	_, err := a.ec2.CreateTags(&ec2.CreateTagsInput{
		Resources: aws.StringSlice(resources),
		Tags:      tagObjs,
	})
	if err != nil {
		return fmt.Errorf("error creating tags: %v", err)
	}
	return err
}

// GetConsoleOutput returns the console output. Returns "", nil if no logs
// are available.
func (a *API) GetConsoleOutput(instanceID string) (string, error) {
	res, err := a.ec2.GetConsoleOutput(&ec2.GetConsoleOutputInput{
		InstanceId: aws.String(instanceID),
	})
	if err != nil {
		return "", fmt.Errorf("couldn't get console output of %v: %v", instanceID, err)
	}

	if res.Output == nil {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(*res.Output)
	if err != nil {
		return "", fmt.Errorf("couldn't decode console output of %v: %v", instanceID, err)
	}

	return string(decoded), nil
}

// getSecurityGroupID gets a security group matching the given name.
// If the security group does not exist, it's created.
func (a *API) getSecurityGroupID(name string) (string, error) {
	// using a Filter on group-name rather than the explicit GroupNames parameter
	// disentangles this call from checking only inside of the default VPC
	sgIds, err := a.ec2.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{&name},
			},
		},
	})

	if len(sgIds.SecurityGroups) == 0 {
		return a.createSecurityGroup(name)
	}

	if err != nil {
		return "", fmt.Errorf("unable to get security group named %v: %v", name, err)
	}

	return *sgIds.SecurityGroups[0].GroupId, nil
}

// getVPCID gets a VPC for the given security group
func (a *API) getVPCID(sgId string) (string, error) {
	sgs, err := a.ec2.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{&sgId},
	})
	if err != nil {
		return "", fmt.Errorf("listing vpc's: %v", err)
	}
	for _, sg := range sgs.SecurityGroups {
		if sg.VpcId != nil {
			return *sg.VpcId, nil
		}
	}
	return "", fmt.Errorf("no vpc found for security group %v", sgId)
}

// getSubnetID gets a subnet for the given VPC.
func (a *API) getSubnetID(vpc string) (string, error) {
	subIds, err := a.ec2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{&vpc},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("unable to get subnets for vpc %v: %v", vpc, err)
	}
	for _, id := range subIds.Subnets {
		if id.SubnetId != nil {
			return *id.SubnetId, nil
		}
	}
	return "", fmt.Errorf("no subnets found for vpc %v", vpc)
}

// creates an InternetGateway and attaches it to the given VPC
func (a *API) createInternetGateway(vpcId *string) (string, error) {
	igw, err := a.ec2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return "", err
	}
	if igw.InternetGateway == nil || igw.InternetGateway.InternetGatewayId == nil {
		return "", fmt.Errorf("internet gateway was nil")
	}
	err = a.tagCreatedByMantle([]string{*igw.InternetGateway.InternetGatewayId})
	if err != nil {
		return "", err
	}
	_, err = a.ec2.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: igw.InternetGateway.InternetGatewayId,
		VpcId:             vpcId,
	})
	if err != nil {
		return "", fmt.Errorf("attaching internet gateway to vpc: %v", err)
	}
	return *igw.InternetGateway.InternetGatewayId, nil
}

// createSubnets creates a subnet in each availability zone for the region
// that is associated with the given VPC associated with the given RouteTable
func (a *API) createSubnets(vpcId, routeTableId *string) error {
	azs, err := a.ec2.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return fmt.Errorf("retrieving availability zones: %v", err)
	}

	for i, az := range azs.AvailabilityZones {
		// 15 is the maximum amount of zones possible when giving them a /20
		// CIDR range inside of a /16 network.
		if i > 15 {
			return nil
		}

		if az.ZoneName == nil {
			continue
		}

		name := *az.ZoneName
		sub, err := a.ec2.CreateSubnet(&ec2.CreateSubnetInput{
			AvailabilityZone: aws.String(name),
			VpcId:            vpcId,
			// Increment the CIDR block by 16 every time
			CidrBlock: aws.String(fmt.Sprintf("172.31.%d.0/20", i*16)),
		})
		if err != nil {
			// Some availability zones get returned but cannot have subnets
			// created inside of them
			if awsErr, ok := (err).(awserr.Error); ok {
				if awsErr.Code() == "InvalidParameterValue" {
					continue
				}
			}
			return fmt.Errorf("creating subnet: %v", err)
		}
		if sub.Subnet == nil || sub.Subnet.SubnetId == nil {
			return fmt.Errorf("subnet was nil after creation")
		}
		err = a.tagCreatedByMantle([]string{*sub.Subnet.SubnetId})
		if err != nil {
			return err
		}
		_, err = a.ec2.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
			SubnetId: sub.Subnet.SubnetId,
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}

		_, err = a.ec2.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: routeTableId,
			SubnetId:     sub.Subnet.SubnetId,
		})
		if err != nil {
			return fmt.Errorf("associating subnet with route table: %v", err)
		}
	}

	return nil
}

// createRouteTable creates a RouteTable with a local target for destination
// 172.31.0.0/16 as well as an InternetGateway for destination 0.0.0.0/0
func (a *API) createRouteTable(vpcId *string) (string, error) {
	rt, err := a.ec2.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: vpcId,
	})
	if err != nil {
		return "", err
	}
	if rt.RouteTable == nil || rt.RouteTable.RouteTableId == nil {
		return "", fmt.Errorf("route table was nil after creation")
	}

	igw, err := a.createInternetGateway(vpcId)
	if err != nil {
		return "", fmt.Errorf("creating internet gateway: %v", err)
	}

	_, err = a.ec2.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igw),
		RouteTableId:         rt.RouteTable.RouteTableId,
	})
	if err != nil {
		return "", fmt.Errorf("creating remote route: %v", err)
	}

	return *rt.RouteTable.RouteTableId, nil
}

// createVPC creates a VPC with an IPV4 CidrBlock of 172.31.0.0/16
func (a *API) createVPC() (string, error) {
	vpc, err := a.ec2.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("172.31.0.0/16"),
	})
	if err != nil {
		return "", fmt.Errorf("creating VPC: %v", err)
	}
	if vpc.Vpc == nil || vpc.Vpc.VpcId == nil {
		return "", fmt.Errorf("vpc was nil after creation")
	}
	err = a.tagCreatedByMantle([]string{*vpc.Vpc.VpcId})
	if err != nil {
		return "", err
	}

	_, err = a.ec2.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: vpc.Vpc.VpcId,
	})
	if err != nil {
		return "", fmt.Errorf("modifying VPC attributes: %v", err)
	}
	_, err = a.ec2.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsSupport: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: vpc.Vpc.VpcId,
	})
	if err != nil {
		return "", fmt.Errorf("modifying VPC attributes: %v", err)
	}

	routeTable, err := a.createRouteTable(vpc.Vpc.VpcId)
	if err != nil {
		return "", fmt.Errorf("creating RouteTable: %v", err)
	}
	err = a.tagCreatedByMantle([]string{routeTable})
	if err != nil {
		return "", err
	}

	err = a.createSubnets(vpc.Vpc.VpcId, &routeTable)
	if err != nil {
		return "", fmt.Errorf("creating subnets: %v", err)
	}

	return *vpc.Vpc.VpcId, nil
}

// createSecurityGroup creates a security group with tcp/22 access allowed from the
// internet.
func (a *API) createSecurityGroup(name string) (string, error) {
	vpcId, err := a.createVPC()
	if err != nil {
		return "", err
	}
	sg, err := a.ec2.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(name),
		Description: aws.String("mantle security group for testing"),
		VpcId:       aws.String(vpcId),
	})
	if err != nil {
		return "", err
	}
	plog.Debugf("created security group %v", *sg.GroupId)

	allowedIngresses := []ec2.AuthorizeSecurityGroupIngressInput{
		{
			// SSH access from the public internet
			// Full access from inside the same security group
			GroupId: sg.GroupId,
			IpPermissions: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
					FromPort: aws.Int64(22),
					ToPort:   aws.Int64(22),
				},
				{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int64(1),
					ToPort:     aws.Int64(65535),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: sg.GroupId,
							VpcId:   &vpcId,
						},
					},
				},
				{
					IpProtocol: aws.String("udp"),
					FromPort:   aws.Int64(1),
					ToPort:     aws.Int64(65535),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: sg.GroupId,
							VpcId:   &vpcId,
						},
					},
				},
				{
					IpProtocol: aws.String("icmp"),
					FromPort:   aws.Int64(-1),
					ToPort:     aws.Int64(-1),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{
						{
							GroupId: sg.GroupId,
							VpcId:   &vpcId,
						},
					},
				},
			},
		},
	}

	for _, input := range allowedIngresses {
		_, err := a.ec2.AuthorizeSecurityGroupIngress(&input)

		if err != nil {
			// We created the SG but can't add all the needed rules, let's try to
			// bail gracefully
			_, delErr := a.ec2.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			})
			if delErr != nil {
				return "", fmt.Errorf("created sg %v (%v) but couldn't authorize it. Manual deletion may be required: %v", *sg.GroupId, name, err)
			}
			return "", fmt.Errorf("created sg %v (%v), but couldn't authorize it and thus deleted it: %v", *sg.GroupId, name, err)
		}
	}
	return *sg.GroupId, err
}

func (a *API) tagCreatedByMantle(resources []string) error {
	return a.CreateTags(resources, map[string]string{
		"CreatedBy": "mantle",
	})
}
