package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yihuaf/addlicense/libaddlicense"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "addlicense [flags] path...",
		Short:        "CLI used to add license to source files",
		SilenceUsage: true,
	}

	rootCmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		logrus.SetReportCaller(true)
		return nil
	}

	rootCmd.PersistentFlags().String("license", "", "Path to license file")
	rootCmd.Args = cobra.MinimumNArgs(1)
	rootCmd.RunE = func(c *cobra.Command, args []string) error {
		licensePath, err := c.Flags().GetString("license")
		if err != nil {
			return err
		}

		if licensePath == "" {
			return errors.New("License file path can't be empty")
		}

		license, err := ioutil.ReadFile(licensePath)
		if err != nil {
			return err
		}

		for _, dir := range args {
			if err := libaddlicense.AddLicense(dir, license); err != nil {
				return errors.Wrapf(err, "Failed to add license to: %s", dir)
			}
		}

		return nil
	}

	rootCmd.Execute()
}
