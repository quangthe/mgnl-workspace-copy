package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	flagNamePgVersion          = "pgversion"
	flagNameDbHost             = "host"
	flagNameDbUsername         = "username"
	flagNameDbPassword         = "password"
	flagNameDbPort             = "port"
	flagNameDbName             = "dbname"
	flagNameDbDumpFile         = "dump-file"
	flagNameMgnlWorkspaces     = "mgnl-workspaces"
	flagNameNormalizeTableName = "normalize-table-name"
	flagNameMigrateDatastore   = "migrate-datastore"
)

type migratorArgs struct {
	pgversion          string
	dbHost             string
	dbPassword         string
	dbUsername         string
	dbPort             string
	dbName             string
	dbDumpPath         string
	mgnlWorkspaces     cli.StringSlice
	normalizeTableName bool
	migrateDatastore   bool
}

type connectionArgs struct {
	pgversion  string
	dbHost     string
	dbPassword string
	dbUsername string
	dbPort     string
	dbName     string
}

func (a migratorArgs) validate() error {
	if a.pgversion == "" {
		return fmt.Errorf("pgversion cannot be empty")
	}
	if a.pgversion != "11" && a.pgversion != "12" {
		return fmt.Errorf("supported postgresql version: 11, 12")
	}
	if a.dbHost == "" {
		return fmt.Errorf("db host cannot be empty")
	}
	if a.dbPort == "" {
		return fmt.Errorf("db port cannot be empty")
	}
	if a.dbUsername == "" {
		return fmt.Errorf("db username cannot be empty")
	}
	if a.dbName == "" {
		return fmt.Errorf("db name cannot be empty")
	}
	if a.dbDumpPath == "" {
		return fmt.Errorf("db dump file path cannot be empty")
	}
	if len(a.mgnlWorkspaces.Value()) == 0 {
		return fmt.Errorf("mgnl workspaces list cannot be empty")
	}
	return nil
}

func (a connectionArgs) validate() error {
	if a.pgversion != "11" && a.pgversion != "12" && a.pgversion != "" {
		return fmt.Errorf("supported postgresql version: 11, 12")
	}
	if a.dbHost == "" {
		return fmt.Errorf("db host cannot be empty")
	}
	if a.dbPort == "" {
		return fmt.Errorf("db port cannot be empty")
	}
	if a.dbUsername == "" {
		return fmt.Errorf("db username cannot be empty")
	}
	if a.dbName == "" {
		return fmt.Errorf("db name cannot be empty")
	}
	return nil
}

var version = "dev"

func main() {
	var migratorArgs migratorArgs
	var connectionArgs connectionArgs

	app := &cli.App{}
	app.Name = "mgnl-workspace-copy"
	app.Usage = "Copy Magnolia workspace data from Postgresql dump file"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.VersionFlag,
	}

	migrateFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagNamePgVersion,
			Aliases:     nil,
			Value:       "12",
			Usage:       "Postgres version to use for migration. Support version: 11, 12",
			Destination: &migratorArgs.pgversion,
			EnvVars:     []string{"PGVERSION"},
		},
		&cli.StringFlag{
			Name:        flagNameDbHost,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB host. Example: dev-magnolia-helm-author-db-0.dev-magnolia-helm-author-db.dev",
			Destination: &migratorArgs.dbHost,
			EnvVars:     []string{"DB_HOST"},
		},
		&cli.StringFlag{
			Name:        flagNameDbUsername,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB username",
			Destination: &migratorArgs.dbUsername,
			EnvVars:     []string{"DB_USERNAME"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPassword,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB user password",
			Destination: &migratorArgs.dbPassword,
			EnvVars:     []string{"DB_PASSWORD"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPort,
			Aliases:     nil,
			Value:       "5432",
			Usage:       "Postgres DB port",
			Destination: &migratorArgs.dbPort,
			EnvVars:     []string{"DB_PORT"},
		},
		&cli.StringFlag{
			Name:        flagNameDbName,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB name. Example: author, public",
			Destination: &migratorArgs.dbName,
			EnvVars:     []string{"DB_NAME"},
		},
		&cli.StringFlag{
			Name:        flagNameDbDumpFile,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB dump file path",
			Destination: &migratorArgs.dbDumpPath,
			EnvVars:     []string{"DB_DUMP_PATH"},
		},
		&cli.StringSliceFlag{
			Name:        flagNameMgnlWorkspaces,
			Aliases:     nil,
			Value:       nil,
			Usage:       "Magnolia workspaces will be migrated",
			Destination: &migratorArgs.mgnlWorkspaces,
			EnvVars:     []string{"MGNL_WORKSPACES"},
		},
		&cli.BoolFlag{
			Name:        flagNameNormalizeTableName,
			Aliases:     nil,
			Value:       true,
			Usage:       "Auto normalize table name. Convert exported workspace name to db table name.",
			Destination: &migratorArgs.normalizeTableName,
			EnvVars:     []string{"NORMALIZE_TABLE_NAME"},
		},
		&cli.BoolFlag{
			Name:        flagNameMigrateDatastore,
			Aliases:     nil,
			Value:       true,
			Usage:       "Migrate ds_datastore table",
			Destination: &migratorArgs.migrateDatastore,
			EnvVars:     []string{"MIGRATE_DATASTORE"},
		},
	}

	connectionFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagNamePgVersion,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres version to use for migration. Support version: 11, 12. If not set the default Postgres version 15 will be used",
			Destination: &connectionArgs.pgversion,
			EnvVars:     []string{"PGVERSION"},
		},
		&cli.StringFlag{
			Name:        flagNameDbHost,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB host. Example: dev-magnolia-helm-author-db-0.dev-magnolia-helm-author-db.dev",
			Destination: &connectionArgs.dbHost,
			EnvVars:     []string{"DB_HOST"},
		},
		&cli.StringFlag{
			Name:        flagNameDbUsername,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB username",
			Destination: &connectionArgs.dbUsername,
			EnvVars:     []string{"DB_USERNAME"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPassword,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB user password",
			Destination: &connectionArgs.dbPassword,
			EnvVars:     []string{"DB_PASSWORD"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPort,
			Aliases:     nil,
			Value:       "5432",
			Usage:       "Postgres DB port",
			Destination: &connectionArgs.dbPort,
			EnvVars:     []string{"DB_PORT"},
		},
		&cli.StringFlag{
			Name:        flagNameDbName,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB name. Example: author, public",
			Destination: &connectionArgs.dbName,
			EnvVars:     []string{"DB_NAME"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "migrate",
			Usage: "run db migration",
			Flags: migrateFlags,
			Action: func(ctx *cli.Context) error {
				if err := migratorArgs.validate(); err != nil {
					return fmt.Errorf("invalid argument: %w", err)
				}
				return migrator(ctx.Context, migratorArgs)
			},
		},
		{
			Name:  "check",
			Usage: "check db connection",
			Flags: connectionFlags,
			Action: func(ctx *cli.Context) error {
				if err := connectionArgs.validate(); err != nil {
					return fmt.Errorf("invalid argument: %w", err)
				}
				return checkConnection(ctx.Context, connectionArgs)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
