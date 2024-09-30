package scripts

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	if err := runGoModTidy(name); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v", err)
	}
	if err := scanAndDownloadImports(name); err != nil {
		return fmt.Errorf("failed to download required libraries: %v", err)
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
func runGoModTidy(projectName string) error {
	currentDir, _ := os.Getwd()
	projectDir := filepath.Join(currentDir, projectName)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy command failed: %v", err)
	}

	return nil
}

func scanAndDownloadImports(projectName string) error {
	currentDir, _ := os.Getwd()
	projectDir := filepath.Join(currentDir, projectName)

	// extract imports first
	imports, err := extractImports(projectDir)
	if err != nil {
		return fmt.Errorf("failed to extract imports: %v", err)
	}

	// download required imports
	for _, imp := range imports {
		if err := runGoGet(projectDir, imp); err != nil {
			return fmt.Errorf("failed to run go get for import %s: %v", imp, err)
		}
	}

	if err := runGoModTidy(projectName); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v", err)
	}

	// remove unused imports using goimports by formating files
	if err := runGoImports(projectDir); err != nil {
		return fmt.Errorf("failed to run goimports: %v", err)
	}

	return nil
}

func extractImports(projectDir string) ([]string, error) {
	var imports []string
	uniqueImports := make(map[string]struct{})

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fileImports, err := parseImports(path)
			if err != nil {
				return err
			}
			for _, imp := range fileImports {
				if _, exists := uniqueImports[imp]; !exists {
					uniqueImports[imp] = struct{}{}
					imports = append(imports, imp)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imports, nil
}

func parseImports(filePath string) ([]string, error) {
	var imports []string

	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, importPath)
	}

	return imports, nil
}
func runGoGet(projectDir, importPath string) error {
	cmd := exec.Command("go", "get", importPath)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get command failed for import %s: %v", importPath, err)
	}

	return nil
}

func runGoImports(projectDir string) error {
	// Find all Go files in the project directory
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Run goimports on each Go file
			cmd := exec.Command("goimports", "-w", path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("goimports command failed for file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to run goimports on project: %v", err)
	}
	return nil
}
