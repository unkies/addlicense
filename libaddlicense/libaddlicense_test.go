// MIT License
//
// Copyright (c) 2020 yihuaf
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package libaddlicense

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	initialPath := filepath.Join("testdata", "initial")
	expectPath := filepath.Join("testdata", "expected")
	licensePath := filepath.Join("testdata", "test_license")

	testDir, err := ioutil.TempDir("/tmp", "addlicense_test_")
	assert.NoError(t, err, "Failed to create a testing directory")

	if err := filepath.Walk(initialPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fname := filepath.Base(path)
		testFile := filepath.Join(testDir, fname)

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(testFile, input, 0644); err != nil {
			return err
		}

		return nil
	}); err != nil {
		assert.NoError(t, err, "Failed to walk")
	}

	license, err := ioutil.ReadFile(licensePath)
	assert.NoError(t, err, "Failed to read license for testing")
	assert.NoError(t, AddLicenseWithIgnore(testDir, license, []string{"ignore"}), "Failed to add license")

	if err := filepath.Walk(testDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fname := filepath.Base(path)
		result, err := ioutil.ReadFile(filepath.Join(testDir, fname))
		assert.NoError(t, err, "Failed to read the result file")

		expect, err := ioutil.ReadFile(filepath.Join(expectPath, fname))
		assert.NoError(t, err, "Failed to read the expect file")

		ret := bytes.Compare(result, expect)
		assert.Equal(t, 0, ret, "The result file is different from expected")

		return nil
	}); err != nil {
		assert.FailNow(t, "Failed to walk")
	}

	os.RemoveAll(testDir)
}
