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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	testdataDir := filepath.Join("testdata", "add")
	initialPath := filepath.Join(testdataDir, "initial")
	expectPath := filepath.Join(testdataDir, "expected")
	licensePath := filepath.Join(testdataDir, "test_license")

	testDir, err := ioutil.TempDir("/tmp", "addlicense_test_")
	assert.NoError(t, err, "Failed to create a testing directory")
	assert.NoError(t, copyDir(initialPath, testDir), "Failed to copy initial files to test dir")

	license, err := ioutil.ReadFile(licensePath)
	assert.NoError(t, err, "Failed to read license for testing")
	assert.NoError(t, RemoveLicense(testDir, license, []string{"ignore"}), "Failed to add license")

	assert.NoError(t, compDir(testDir, expectPath), "")

	os.RemoveAll(testDir)
}

func TestAdd(t *testing.T) {
	testdataDir := filepath.Join("testdata", "add")
	initialPath := filepath.Join(testdataDir, "initial")
	expectPath := filepath.Join(testdataDir, "expected")
	licensePath := filepath.Join(testdataDir, "test_license")

	testDir, err := ioutil.TempDir("/tmp", "addlicense_test_")
	assert.NoError(t, err, "Failed to create a testing directory")
	assert.NoError(t, copyDir(initialPath, testDir), "Failed to copy initial files to test dir")

	license, err := ioutil.ReadFile(licensePath)
	assert.NoError(t, err, "Failed to read license for testing")
	assert.NoError(t, AddLicense(testDir, license, []string{"ignore"}), "Failed to add license")

	assert.NoError(t, compDir(testDir, expectPath), "")

	os.RemoveAll(testDir)
}

func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fname := filepath.Base(path)
		testFile := filepath.Join(dest, fname)

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(testFile, input, 0644); err != nil {
			return err
		}

		return nil
	})
}

func compDir(result string, expected string) error {
	return filepath.Walk(result, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fname := filepath.Base(path)
		result, err := ioutil.ReadFile(filepath.Join(result, fname))
		if err != nil {
			return errors.Wrap(err, "Failed to read the result file")
		}

		expect, err := ioutil.ReadFile(filepath.Join(expected, fname))
		if err != nil {
			return errors.Wrap(err, "Failed to read the expect file")
		}

		if bytes.Compare(result, expect) != 0 {
			return errors.Wrap(err, "The result file is different from expected")
		}

		return nil
	})
}
