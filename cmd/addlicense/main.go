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

package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unkies/addlicense/libaddlicense"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "addlicense [flags] path...",
		Short: "CLI used to add license to source files",
	}

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		logrus.SetReportCaller(true)
		debug, err := c.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}

	rootCmd.PersistentFlags().String("license", "", "Path to license file")
	rootCmd.PersistentFlags().StringArray("ignore", []string{}, "Patterns to ignore. Follow the shell pattern")
	rootCmd.Args = cobra.MinimumNArgs(1)
	rootCmd.RunE = func(c *cobra.Command, args []string) error {
		licensePath, err := c.Flags().GetString("license")
		if err != nil {
			return err
		}

		if licensePath == "" {
			return errors.New("License file path can't be empty")
		}

		ignore, err := c.Flags().GetStringArray("ignore")
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"root":         args,
			"license file": licensePath,
			"ignore":       ignore,
		}).Debug("Adding license")

		if err := run(args, licensePath, ignore); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"license path": licensePath,
				"directories":  args,
			}).Fatal("Failed to add license")
		}

		return nil
	}

	rootCmd.Execute()
}

func run(dirs []string, licensePath string, ignore []string) error {
	license, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := libaddlicense.AddLicenseWithIgnore(dir, license, ignore); err != nil {
			return errors.Wrapf(err, "Failed to add license to: %s", dir)
		}
	}

	return nil
}
