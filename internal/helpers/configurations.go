package helpers

import (
	"fmt"

	"github.com/TejasGhatte/go-sail/internal/models"
	"github.com/TejasGhatte/go-sail/internal/initializers"
)

// Provider interface defines methods for generating database connection and migration code
type Provider interface {
	GetImports() []string
	GetConnectionCode() string
	GetMigrationCode() string
	GetDBVariable() string
}

type CombinationProvider struct {
	Database    models.DatabaseConfig
	ORM         models.ORMConfig
	Combination models.CombinationConfig
	MigrationCode string
}

func (cp *CombinationProvider) GetImports() []string {
	imports := []string{
		fmt.Sprintf("%q", cp.Database.DriverPkg),
		fmt.Sprintf("%q", cp.ORM.ImportPath),
	}
	for _, additionalImport := range cp.Combination.AdditionalImports {
		imports = append(imports, fmt.Sprintf("%q", additionalImport))
	}

	return imports
}

func (cp *CombinationProvider) GetConnectionCode() string {
	return  fmt.Sprintf(`
	var err error
	dsn := fmt.Sprintf(%q, "your_username", "your_password", "your_database")
	DB, err = %s
	if err != nil {
		fmt.Println("failed to connect to database")
	}
	fmt.Println("Connect to database")
	`, cp.Combination.DSNTemplate, cp.Combination.InitFunc)
}

func (cp *CombinationProvider) GetMigrationCode() string {
	return cp.MigrationCode
}

func (cp *CombinationProvider) GetDBVariable() string {
	return fmt.Sprintf("*%s.DB", cp.ORM.Name)
}

// ProviderFactory creates a provider for a specific database and ORM combination
func ProviderFactory(database, orm string) (Provider, error) {
	dbConfig := initializers.Config.Databases[database]

	ormConfig := initializers.Config.ORMs[orm]

	combinationConfig := initializers.Config.Combinations[database][orm]

	return &CombinationProvider{
		Database:    dbConfig,
		ORM:         ormConfig,
		Combination: combinationConfig,
		MigrationCode: initializers.Config.MigrationCode[orm],
	}, nil
}