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

	err = PopulateDirectory(options)
	if err != nil {
		return err
	}

	return nil
}

func PopulateDirectory(ctx *models.Options) error {
	if err := GitClone(ctx.ProjectName, ctx.Framework, initializers.Config.Repositories[ctx.Framework]); err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	currentDir, _ := os.Getwd()
	folder := filepath.Join(currentDir, ctx.ProjectName, "initializers")

	if ctx.Database != "" && ctx.ORM != "" {
		provider, err := helpers.ProviderFactory(ctx.Database, ctx.ORM)
		if err != nil {
			return fmt.Errorf("error creating database provider: %v", err)
		}

		err = helpers.GenerateDatabaseFile(folder, provider)
		if err != nil {
			return fmt.Errorf("error generating database file: %v", err)
		}

		err = helpers.GenerateMigrationFile(folder, provider)
		if err != nil {
			return fmt.Errorf("error generating migration file: %v", err)
		}

	}
	return nil
}
