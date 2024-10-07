package pgsql

import (
	"context"
	"cybertask/config"
	"cybertask/internal/logger"
	"cybertask/model"
	"database/sql"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type postgresDB struct {
	l   *logger.Logger
	bun *bun.DB

	shutdown sync.Once
}

// Postgres database constructor.
func New(l *logger.Logger, cfg config.Database) (*postgresDB, error) {

	dsn := assembleDSN(cfg)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	bundb := bun.NewDB(sqldb, pgdialect.New())
	pgdb := &postgresDB{
		l:        l,
		bun:      bundb,
		shutdown: sync.Once{},
	}

	return pgdb, nil

}

// CreateTask inserts single task into underlying database.
func (db *postgresDB) CreateTask(ctx context.Context, task model.Task) error {
	_, err := db.bun.
		NewInsert().
		Model(task).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// CreateTask inserts single task into underlying database.
func (db *postgresDB) CreateTasks(ctx context.Context, tasks []model.Task) error {
	_, err := db.bun.
		NewInsert().
		Model(&tasks).
		Column("id", "header", "description", "created_at", "status").
		Table("tasks").
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// DeleteTask deletes single task from underlying database.
func (db *postgresDB) DeleteTask(ctx context.Context, id model.TaskID) error {
	_, err := db.bun.
		NewDelete().
		Table("tasks").
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// GetTask searches and returns single unique task by id.
func (db *postgresDB) GetTask(ctx context.Context, id model.TaskID) (model.Task, error) {

	task := model.Task{}

	err := db.bun.
		NewSelect().
		Model(&task).
		Column("id", "header", "description", "created_at", "status").
		Table("tasks").
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// UpdateTask modifies state of the task in underlying database.
func (db *postgresDB) UpdateTask(ctx context.Context, task model.Task) error {
	_, err := db.bun.
		NewUpdate().
		Model(task).
		Column("id", "header", "description", "created_at", "status").
		Table("tasks").
		Where("id = ?", task.ID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

// Multi-call safe database shutdown.
func (db *postgresDB) Shutdown() error {
	var err error
	db.shutdown.Do(
		func() {
			err = db.bun.Close()
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func assembleDSN(cfg config.Database) string {

	user := cfg.User
	port := cfg.Port
	dbname := cfg.Dbname
	address := cfg.Address

	password := viper.GetString("PG_PASSWORD")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user,
		password,
		address,
		port,
		dbname)
}
