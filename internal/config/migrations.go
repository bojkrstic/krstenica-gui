package config

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Direction - points to migration direction. Usually you should use either Up or ToVersion. All other directions should
// be used with caution
type Direction string

const (
	// Up - set it when you want to migrate to the latest version, should usually be default
	Up Direction = "up"
	// Down is not supported, use StepBack or ToVersion instead. We don't want to remove whole DB with migration
	Down Direction = "down"
	// ToVersion - set it when you want to migrate to specific version, can be higher or lower
	ToVersion Direction = "to_version"
	// StepBack - set it when you want to rollback to previous version
	// IMPORTANT! be careful using it with multiple pods in k8s, because each pod will revert already reverted migration
	StepBack Direction = "step_back"

	stepBack = -1
)

type MigrationConfig struct {
	Direction Direction `mapstructure:"direction"`
	Version   int       `mapstructure:"version"`
}

// PostgresMigrate runs postgres migrations. Returns version was migrated to.
// Full Down migration is not supported, use StepBack or ToVersion instead
func PostgresMigrate(connectionStr string, cnf MigrationConfig, migrationFiles embed.FS) (int, error) {
	if cnf.Direction == Down {
		return 0, fmt.Errorf("down migration not supported, you can use step_back instead")
	}

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return 0, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, fmt.Errorf("failed to init postgres driver: %w", err)
	}

	// use preloaded migration files
	d, err := iofs.New(migrationFiles, "postgres")
	if err != nil {
		return 0, fmt.Errorf("failed to create postgres migrator source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return 0, fmt.Errorf("failed to create postgres migrator instance: %w", err)
	}

	switch cnf.Direction {
	case Up: // usual UP migration, until the latest version
		err = m.Up()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				return 0, fmt.Errorf("failed to run postgres migrations: %w", err)
			}
		}
	case ToVersion: // migrate to specific version
		if cnf.Version == 0 {
			return 0, fmt.Errorf("invalid migration config version: %d", cnf.Version)
		}

		v, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) { // if version is 0, that is fine if we go up
				v = 0
			} else {
				return 0, fmt.Errorf("failed to get postgres migration version: %w", err)
			}
		}

		steps, err := getSteps(int(v), cnf.Version)
		if err != nil {
			return 0, fmt.Errorf("failed to get postgres migration steps: %w", err)
		}

		if dirty {
			return 0, fmt.Errorf("postgres migration is dirty, please run down migration first")
		}

		if steps != 0 {
			err = m.Steps(steps)
			if err != nil {
				return 0, fmt.Errorf("failed to run postgres migrations: %w", err)
			}
		}
	case StepBack: // migrate back to previous version
		err = m.Steps(stepBack)
		if err != nil {
			return 0, fmt.Errorf("failed to run postgres migrations: %w", err)
		}
	default:
		return 0, fmt.Errorf("invalid migration config diretion: %s", cnf.Direction)
	}

	v, _, err := m.Version()
	if err != nil {
		return 0, fmt.Errorf(`migration went well, but version fetch failed.
								If it was down to 0 migration then this error is fine, err: %w`, err)
	}

	return int(v), nil
}

func getSteps(current, needed int) (int, error) {
	if needed < 1 {
		return 0, fmt.Errorf("invalid migration config version: %d", needed)
	}

	// in this case we have valid math because steps can be negative
	return needed - current, nil
}
