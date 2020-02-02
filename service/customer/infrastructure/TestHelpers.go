// +build test

package infrastructure

import (
	"database/sql"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/eventstore/postgres/database"
)

func SetUpDIContainer() (*cmd.DIContainer, error) {
	config, err := cmd.NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	diContainer, err := cmd.NewDIContainer(
		db,
		events.UnmarshalCustomerEvent,
	)
	if err != nil {
		return nil, err
	}

	return diContainer, nil
}
