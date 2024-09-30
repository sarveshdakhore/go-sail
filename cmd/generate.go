package cmd

import (
	"fmt"

	"github.com/TejasGhatte/go-sail/internal/errors"
	"github.com/TejasGhatte/go-sail/internal/scripts"
	"github.com/spf13/cobra"
)

var CreateProjectCommand *cobra.Command

func init() {
	CreateProjectCommand = &cobra.Command{
		Use:   "create [project-name]",
		Short: "Creates a new go project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			ctx := cmd.Context()
			if err := scripts.CreateProject(ctx, projectName); err != nil {
				if err == errors.ErrInterrupt {
					fmt.Println("Program Exited: interrupt")
				} else {
					fmt.Printf("Program Exited: %v\n", err)
				}
			}
		},
	}
}
