package main

import (
	"log"

	"github.com/lhlyu/kratos-easy/cmd/ke/internal/proto"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ke",
	Short:   "ke: An elegant toolkit for Go microservices.",
	Long:    `ke: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(proto.CmdAPI)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
