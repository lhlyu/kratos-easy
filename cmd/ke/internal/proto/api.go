package proto

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdAPI generates an API proto file.
var CmdAPI = &cobra.Command{
	Use:   "api [name]",
	Short: "Generate an API proto file",
	Long:  "Generate a standard Kratos API proto file in the current project. Example: ke api demo",
	RunE:  runAPI,
}

func runAPI(_ *cobra.Command, args []string) error {
	var name string

	// Allow `ke api demo` or interactive input
	if len(args) > 0 {
		name = args[0]
	} else {
		prompt := &survey.Input{
			Message: "Enter the API name:",
			Help:    "It will generate api/<name>/v1/<name>.proto",
		}
		if err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	p := newProto(name)

	if err := p.Generate(); err != nil {
		return err
	}

	fmt.Printf("✔ Created %s\n", p.FilePath())
	return nil
}
