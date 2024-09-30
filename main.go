package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/TejasGhatte/go-sail/cmd"
	"github.com/TejasGhatte/go-sail/internal/initializers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-sail",
	Short: "A CLI for generating project templates for Go backend frameworks",
	Long:  `go-sail is a CLI tool that generates project templates for Go backend frameworks like Fiber, Echo, and Gin, with pre-configured logging and caching, helping developers quickly set up and initialize projects. Users can choose their own database and ORM configurations, and go-sail generates the necessary files for the project.`,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handling ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		cleanup(cmd.ProjectName)
		cancel()
		os.Exit(1)
	}()

	initializers.LoadConfig("config.yml")
	rootCmd.AddCommand(cmd.CreateProjectCommand)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cleanup(projectName string) {
	fmt.Println("\nReceived interrupt signal, exiting...")
	if projectName != "" {
		currentDir, _ := os.Getwd()
		projectDir := filepath.Join(currentDir, projectName)
		if err := os.RemoveAll(projectDir); err != nil {
			fmt.Printf("Failed to remove project directory %s: %v\n", projectDir, err)
		} else {
			fmt.Printf("Successfully removed project directory %s\n", projectDir)
		}
	}
}
