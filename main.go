package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	// Create a context that is canceled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		// Handle the interrupt signal
		cleanup()
		cancel()
		os.Exit(1)
	}()

	initializers.LoadConfig("config.yml")
	rootCmd.AddCommand(cmd.CreateProjectCommand)

	// Execute the root command with context
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// cleanup function to handle any necessary cleanup before exiting
func cleanup() {
	// Perform any necessary cleanup here
	// For example, you can print a message or close any open resources
	fmt.Println("\nReceived interrupt signal, exiting...")
}
