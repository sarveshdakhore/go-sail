package main

import (
	"github.com/TejasGhatte/go-sail/cmd"
	"github.com/TejasGhatte/go-sail/internal/initializers"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "go-sail",
	Short: "A cli for generating project templates for Go backend frameworks",
	Long: `go-sail is a CLI tool that generates project templates for Go backend frameworks like Fiber, Echo, and Gin, with pre-configured logging and caching, helping developers quickly set up and initialize projects. User can choose their own database and orm configurations, and go-sail generates the necessary files for the project.`,
}

func main() {
	initializers.LoadConfig("config.yml")
	rootCmd.AddCommand(cmd.CreateProjectCommand)
	cobra.CheckErr(rootCmd.Execute())
}
