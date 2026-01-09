package proto

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdAPI generates an API proto file.
var CmdAPI = &cobra.Command{
	Use:   "api [name] [version] [protoName]",
	Short: "Generate an API proto file",
	Long:  "Generate a standard Kratos API proto file in the current project. Example: ke api demo v1 user",
	RunE:  runAPI,
}

func runAPI(_ *cobra.Command, args []string) error {
	var (
		name      string
		version   = "v1"
		protoName string
	)

	switch len(args) {
	case 0:
		// interactive input for name
		prompt := &survey.Input{
			Message: "Enter the API name:",
			Help:    "It will generate api/<name>/v1/<name>.proto",
		}
		if err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		protoName = name
	case 1:
		name = args[0]
		protoName = name
	case 2:
		name = args[0]
		version = args[1]
		protoName = name
	case 3:
		name = args[0]
		version = args[1]
		protoName = args[2]
	default:
		return fmt.Errorf("too many arguments")
	}

	p := newProto(name, version, protoName)

	if err := p.Generate(); err != nil {
		return err
	}

	fmt.Printf("✔ Created %s\n", p.FilePath())
	return nil
}
