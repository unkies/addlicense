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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// AddLicenseWithIgnore will add license to all files under the root, except
// path that matches the ignore patterns. The pattern used are regex patterns
// and match is implemented using golang regexp package.
func AddLicenseWithIgnore(root string, license []byte, ignore []string) error {
	var wg errgroup.Group

	if err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Failed to walk the path")
		}

		if f.IsDir() {
			return nil
		}

		for _, pattern := range ignore {
			matched, err := regexp.MatchString(pattern, path)
			if err != nil {
				return err
			}

			// Ignore path that matches the ignore pattern.
			if matched {
				return nil
			}
		}

		wg.Go(func() error {
			return AddLicenseSingle(path, license)
		})
		return nil
	}); err != nil {
		return errors.Wrap(err, "Failed to add license")
	}

	if err := wg.Wait(); err != nil {
		return nil
	}

	return nil
}

// AddLicense adds license to all files under the root. It will try to parse the
// file extension and add the right header accordingly. It will also handle
// shebang lines correctly.
func AddLicense(root string, license []byte) error {
	return AddLicenseWithIgnore(root, license, []string{})
}

// AddLicenseSingle add the license to a single file in path.
func AddLicenseSingle(path string, license []byte) error {
	header, err := licenseHeader(path, license)
	if err != nil || header == nil {
		return errors.Wrap(err, "Failed to create a header from the license")
	}

	f, err := ioutil.ReadFile(path)
	if err != nil || hasLicenseHeader(f, header) {
		return err
	}

	// If the file has shebang the requires to be at the beginning of the file,
	// we need to set the license after the shebang.
	line := hashBang(f)
	if len(line) > 0 {
		f = f[len(line):]
		if line[len(line)-1] != '\n' {
			line = append(line, '\n')
		}

		header = append(line, header...)
	}

	f = append(header, f...)

	// When we write the file, we need to preserve the file mode.
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, f, info.Mode()); err != nil {
		return err
	}

	return nil
}

// licenseHeader process the license into header based on the file extension. We
// assume the license passed in is read-only.
func licenseHeader(path string, license []byte) ([]byte, error) {
	var out []byte
	var err error

	switch fileExtension(path) {
	default:
		return nil, nil
	case ".c", ".h":
		out, err = prefix(license, "/*", " * ", " */")
	case ".js", ".mjs", ".cjs", ".jsx", ".tsx", ".css", ".tf", ".ts":
		out, err = prefix(license, "/**", " * ", " */")
	case ".cc", ".cpp", ".cs", ".go", ".hh", ".hpp", ".java", ".m", ".mm", ".proto", ".rs", ".scala", ".swift", ".dart", ".groovy", ".kt", ".kts":
		out, err = prefix(license, "", "// ", "")
	case ".py", ".sh", ".yaml", ".yml", ".dockerfile", "dockerfile", ".rb", "gemfile":
		out, err = prefix(license, "", "# ", "")
	case ".el", ".lisp":
		out, err = prefix(license, "", ";; ", "")
	case ".erl":
		out, err = prefix(license, "", "% ", "")
	case ".hs", ".sql":
		out, err = prefix(license, "", "-- ", "")
	case ".html", ".xml", ".vue":
		out, err = prefix(license, "<!--", " ", "-->")
	case ".php":
		out, err = prefix(license, "", "// ", "")
	case ".ml", ".mli", ".mll", ".mly":
		out, err = prefix(license, "(**", "   ", "*)")
	}

	return out, err
}

func prefix(license []byte, top, mid, bot string) ([]byte, error) {
	buf := bytes.NewBuffer(license)

	var out bytes.Buffer
	if top != "" {
		fmt.Fprintln(&out, top)
	}

	s := bufio.NewScanner(buf)
	for s.Scan() {
		fmt.Fprintln(&out, strings.TrimRightFunc(mid+s.Text(), unicode.IsSpace))
	}

	if bot != "" {
		fmt.Fprintln(&out, bot)
	}

	fmt.Fprintln(&out)
	return out.Bytes(), nil
}

func fileExtension(name string) string {
	if v := filepath.Ext(name); v != "" {
		return strings.ToLower(v)
	}
	return strings.ToLower(filepath.Base(name))
}

var head = []string{
	"#!",                       // shell script
	"<?xml",                    // XML declaratioon
	"<!doctype",                // HTML doctype
	"# encoding:",              // Ruby encoding
	"# frozen_string_literal:", // Ruby interpreter instruction
	"<?php",                    // PHP opening tag
}

func hashBang(b []byte) []byte {
	var line []byte
	for _, c := range b {
		line = append(line, c)
		if c == '\n' {
			break
		}
	}

	first := strings.ToLower(string(line))
	for _, h := range head {
		if strings.HasPrefix(first, h) {
			return line
		}
	}

	return nil
}

func hasLicenseHeader(file []byte, header []byte) bool {
	return bytes.Contains(file, header)
}

func stripHeader(file []byte, license []byte) []byte {
	return bytes.ReplaceAll(file, license, []byte{})
}
