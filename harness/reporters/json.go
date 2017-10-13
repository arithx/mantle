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
        "encoding/json"
	"fmt"
	"os"
        "time"
)

type jsonReporter struct {
	Tests []jsonTest `json:"tests"`
	Result string `json:"result"`
	file *os.File
	filename string

	// Context variables
	Platform string `json:"platform"`
	Version string `json:"version"`
}

func NewJSONReporter(filename, platform, version string) *jsonReporter {
	return &jsonReporter{
		Platform: platform,
		Version: version,
		filename: filename,
	}
}

func (r *jsonReporter) OpenFile(filepath string) error {
	f, err := os.Create(filepath)
	r.file = f
	return err
}

func (r *jsonReporter) ReportTest(name, result string, duration time.Duration, b []byte) {
	r.Tests = append(r.Tests, jsonTest{
		Name: name,
		Result: result,
		Duration: duration,
		Output: fmt.Sprintf("%s", string(b)),
	})
}

func (r *jsonReporter) Cleanup() error {
	return r.file.Close()
}

func (r *jsonReporter) Output() error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(r.file, "%s", b); err != nil {
		return err
	}
	return nil
}

func (r *jsonReporter) SetResult(result string) {
	r.Result = result
}

func (r *jsonReporter) Filename() string {
	return r.filename
}

type jsonTest struct {
	Name string `json:"name"`
	Result string `json:"result"`
	Duration time.Duration `json:"duration"`
	Output string `json:"output"`
}
