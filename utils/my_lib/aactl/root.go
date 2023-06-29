package aactl

import (
	"fmt"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/aactl/aac"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/aactl/templates"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var programName string

func NewAACCommand() *cobra.Command {
	programName := GetProgramName(filepath.Base(os.Args[0]))
	// used in cobra templates to display either `kubectl virt` or `virtctl`
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return programName
		},
	)

	// used to enable replacement of `ProgramName` placeholder for cobra.Example, which has no template support
	cobra.AddTemplateFunc(
		"prepare", func(s string) string {
			// order matters!
			result := strings.Replace(s, "kubectl", "kubectl aac", -1)
			result = strings.Replace(result, "{{ProgramName}}", programName, -1)
			return result
		},
	)

	rootCmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         filepath.Base(os.Args[0]) + " controls alertManager alerts related operations on your alerts_arms_center cluster.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(cmd.UsageString())
		},
	}

	optionsCmd := &cobra.Command{
		Use:    "options",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(cmd.UsageString())
		},
	}
	optionsCmd.SetUsageTemplate(templates.OptionsUsageTemplate())
	rootCmd.SetUsageTemplate(templates.MainUsageTemplate())
	rootCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(
		aac.NewAACAlertCommand(),
	)
	return rootCmd
}

func GetProgramName(binary string) string {
	if strings.HasSuffix(binary, "-aac") {
		return fmt.Sprintf("%s aac", strings.TrimSuffix(binary, "-aac"))
	}
	return "aactl"
}

func Execute() {
	cmd := NewAACCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(cmd.Root().ErrOrStderr(), strings.TrimSpace(err.Error()))
		os.Exit(1)
	}
}
