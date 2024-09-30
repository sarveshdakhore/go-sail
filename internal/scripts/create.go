package scripts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/TejasGhatte/go-sail/internal/errors"
	"github.com/TejasGhatte/go-sail/internal/helpers"
	"github.com/TejasGhatte/go-sail/internal/initializers"
	"github.com/TejasGhatte/go-sail/internal/models"
	"github.com/TejasGhatte/go-sail/internal/prompts"
	"github.com/briandowns/spinner"
)

func CreateProject(ctx context.Context, name string) error {
	framework, err := prompts.SelectFramework(ctx)
	if err != nil {
		if err == errors.ErrInterrupt {
			return err
		}
		return err
	}

	database, err := prompts.SelectDatabase(ctx)
	if err != nil {
		if err == errors.ErrInterrupt {
			return err
		}
		return err
	}

	var orm string
	if database != "" {
		orm, err = prompts.SelectORM(ctx)
		if err != nil {
			if err == errors.ErrInterrupt {
				return err
			}
			return err
		}
	}

	fmt.Println("Generating project with the following options:")
	fmt.Printf("Framework: %s, Database: %s, ORM: %s\n", framework, database, orm)

	options := &models.Options{
		ProjectName: name,
		Framework:   framework,
		Database:    database,
		ORM:         orm,
	}

	// Spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	err = PopulateDirectory(ctx, options)
	if err != nil {
		return err
	}

	return nil
}

func PopulateDirectory(ctx context.Context, options *models.Options) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := GitClone(ctx, options.ProjectName, options.Framework, initializers.Config.Repositories[options.Framework]); err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	currentDir, _ := os.Getwd()
	folder := filepath.Join(currentDir, options.ProjectName, "initializers")

	if options.Database != "" && options.ORM != "" {
		provider, err := helpers.ProviderFactory(options.Database, options.ORM)
		if err != nil {
			return fmt.Errorf("error creating database provider: %v", err)
		}

		err = helpers.GenerateDatabaseFile(ctx, folder, provider)
		if err != nil {
			return fmt.Errorf("error generating database file: %v", err)
		}

		err = helpers.GenerateMigrationFile(ctx, folder, provider)
		if err != nil {
			return fmt.Errorf("error generating migration file: %v", err)
		}
	}
	return nil
}
