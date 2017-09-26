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

package types

import (
	"testing"
)

func strChecker(t *testing.T, name, expected, received string) {
	if expected != received {
		t.Fatalf("%s is wrong:\n\texpected: %s\n\treceived: %s", name, expected, received)
	}
}

func pStrChecker(t *testing.T, name string, expected, received *string) {
	if expected == nil && received == nil {
		return
	}
	if (expected == nil && received != nil) || (expected != nil && received == nil) {
		t.Fatalf("%s is wrong")
	} else if *expected != *received {
		t.Fatalf("%s is wrong:\n\texpected: %s\n\treceived: %s", name, *expected, *received)
	}
}

func boolChecker(t *testing.T, name string, expected, received bool) {
	if expected != received {
		t.Fatalf("%s is wrong:\n\texpected: %t\n\treceived: %t", name, expected, received)
	}
}

func pBoolChecker(t *testing.T, name string, expected, received *bool) {
	if expected == nil && received == nil {
		return
	}
	if (expected == nil && received != nil) || (expected != nil && received == nil) {
		t.Fatalf("%s is wrong")
	} else if *expected != *received {
		t.Fatalf("%s is wrong:\n\texpected: %t\n\treceived: %t", name, *expected, *received)
	}
}

func intChecker(t *testing.T, name string, expected, received int) {
	if expected != received {
		t.Fatalf("%s is wrong:\n\texpected: %d\n\treceived: %d", name, expected, received)
	}
}

func pIntChecker(t *testing.T, name string, expected, received *int) {
	if expected == nil && received == nil {
		return
	}
	if (expected == nil && received != nil) || (expected != nil && received == nil) {
		t.Fatalf("%s is wrong")
	} else if *expected != *received {
		t.Fatalf("%s is wrong:\n\texpected: %d\n\treceived: %d", name, *expected, *received)
	}
}
