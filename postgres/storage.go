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

func (s *Storage) Elements() (es []entity.Element, err error) {
	err = s.DB.Select(&es, `
		select * from element;
	`)
	return
}

func (s *Storage) AddUserError(ue entity.UserError) error {
	_, err := s.DB.Exec(`
		insert into user_error (id, element_id, message, created_at, contact)
		values ($1, $2, $3, $4, $5)
	`, ue.ID, ue.ElementID, ue.Message, ue.CreatedAt, ue.Contact)
	return err
}

func (s *Storage) UserErrors() (ues []entity.UserError, err error) {
	err = s.DB.Select(&ues, `
		select * from user_error;
	`)
	return
}

func (s *Storage) AddErrorNotification(n entity.ErrorNotification) error {
	_, err := s.DB.Exec(`
		insert into error_notification (
		    id,
			element_id,
			message,
			resolved,
			created_at,
			resolved_at
		)
		values ($1, $2, $3, $4, $5, $6)
	`, n.ID, n.ElementID, n.Message, n.Resolved, n.CreatedAt, n.ResolvedAt)
	return err
}

func (s *Storage) ResolveErrorNotification(id uuid.UUID, message string) error {
	_, err := s.DB.Exec(`
		update error_notification set resolved = true, resolve_message = $1, resolved_at = now() where id = $2
	`, message, id)
	return err
}

func (s *Storage) ErrorNotifications() (ns []entity.ErrorNotification, err error) {
	err = s.DB.Select(&ns, `
		select * from error_notification order by created_at desc
	`)
	return ns, err
}

func (s *Storage) AddErrorResolveWaiter(w entity.ErrorResolveWaiter) error {
	_, err := s.DB.Exec(`
		insert into error_resolve_waiter (
		    id,
			error_notification_id,
			contact
		)
		values ($1, $2, $3)
	`, w.ID, w.ErrorNotificationID, w.Contact)
	return err
}

func (s *Storage) ErrorResolveWaiterStats() (ss []entity.ErrorResolveWaiterStats, err error) {
	err = s.DB.Select(&ss, `
		select error_notification_id, count(*) as waiters_count from error_resolve_waiter
			group by error_notification_id order by waiters_count desc 
	`)
	return
}

func (s *Storage) RemoveErrorResolveWaiters(enID uuid.UUID) (ws []entity.ErrorResolveWaiter, err error) {
	err = s.DB.Select(&ws, `
		delete from error_resolve_waiter where error_notification_id=$1 returning *
	`, enID)
	return
}
