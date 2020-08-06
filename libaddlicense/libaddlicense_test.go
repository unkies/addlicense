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
	licensePath := filepath.Join("testdata", "license.txt")

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
	assert.NoError(t, AddLicense(testDir, license), "Failed to add license")

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
