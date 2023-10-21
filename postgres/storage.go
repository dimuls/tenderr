package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"tednerr/entity"
)

//go:embed migrations/*.sql
var fs embed.FS

type Storage struct {
	DB *sqlx.DB
}

func (s *Storage) Migrate() error {
	srcDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("new iofs: %w", err)
	}

	driver, err := postgres.WithInstance(s.DB.DB, &postgres.Config{
		MigrationsTable: "migration",
	})
	if err != nil {
		return fmt.Errorf("new postgres migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", srcDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("new migration: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (s *Storage) Classes() (cs []entity.Class, err error) {
	return cs, s.DB.Select(&cs, `
		select * from class where rules is not null;
	`)
}

func (s *Storage) SetClasses(cs []entity.Class) (err error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("begig tx: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = s.DB.Exec(`
		delete from class;
	`)
	if err != nil {
		return fmt.Errorf("delete from class: %w", err)
	}

	for _, c := range cs {
		_, err = s.DB.Exec(`
			insert into class (id, name, rules) values ($1, $2, $3)
		`, c.ID, c.Name, pq.Array(c.Rules))
		if err != nil {
			return fmt.Errorf("insert class: %w", err)
		}
	}

	return err
}

func (s *Storage) RemoveClass(classID uuid.UUID) error {
	_, err := s.DB.Exec(`
		delete from class where id = $1
	`, classID)
	return err
}
