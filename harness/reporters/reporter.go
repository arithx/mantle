// Copyright 2017 CoreOS, Inc.
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

package reporters

import (
	"fmt"
	"time"

	"github.com/coreos/pkg/multierror"
)

type Reporters []Reporter

func (reps *Reporters) ReportTest(name, result string, duration time.Duration, b []byte) {
	for _, r := range *reps {
		r.ReportTest(name, result, duration, b)
	}
}

func (reps *Reporters) Cleanup() error {
	var err multierror.Error

	for _, r := range *reps {
		e := r.Cleanup()
		if e != nil {
			fmt.Println(e)
			err = append(err, e)
		}
	}

	return err.AsError()
}

func (reps *Reporters) Output() error {
	for _, r := range *reps {
		err := r.Output()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (reps *Reporters) SetResult(s string) {
	for _, r := range *reps {
		r.SetResult(s)
	}
}

func (reps *Reporters) OpenFile(pathFunc func(string) string) error {
	for _, r := range *reps {
		err := r.OpenFile(pathFunc(r.Filename()))
		if err != nil {
			return err
		}
	}
	return nil
}

type Reporter interface {
	OpenFile(string) error
	ReportTest(string, string, time.Duration, []byte)
	Cleanup() error
	Output() error
	SetResult(string)
	Filename() string
}
