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
	flagNameCopyDatastore      = "copy-datastore"
)

type copyArgs struct {
	pgversion          string
	dbHost             string
	dbPassword         string
	dbUsername         string
	dbPort             string
	dbName             string
	dbDumpPath         string
	mgnlWorkspaces     cli.StringSlice
	normalizeTableName bool
	copyDatastore      bool
}

type connectionArgs struct {
	pgversion  string
	dbHost     string
	dbPassword string
	dbUsername string
	dbPort     string
	dbName     string
}

func (a copyArgs) validate() error {
	if a.pgversion == "" {
		return fmt.Errorf("pgversion cannot be empty")
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

var version string

func main() {
	var copyArgs copyArgs
	var connectionArgs connectionArgs

	app := &cli.App{
		Name:    "mgnl-workspace-copy",
		Usage:   "Copy Magnolia workspace data from Postgresql dump file",
		Version: version,
	}

	copyFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagNamePgVersion,
			Aliases:     nil,
			Value:       "12",
			Usage:       "Postgres version of DB dump file. Support version: 11, 12",
			Destination: &copyArgs.pgversion,
			EnvVars:     []string{"PGVERSION"},
		},
		&cli.StringFlag{
			Name:        flagNameDbHost,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB host. Example: dev-magnolia-helm-author-db-0.dev-magnolia-helm-author-db.dev",
			Destination: &copyArgs.dbHost,
			EnvVars:     []string{"DB_HOST"},
		},
		&cli.StringFlag{
			Name:        flagNameDbUsername,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB username",
			Destination: &copyArgs.dbUsername,
			EnvVars:     []string{"DB_USERNAME"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPassword,
			Aliases:     nil,
			Value:       "postgres",
			Usage:       "Postgres DB user password",
			Destination: &copyArgs.dbPassword,
			EnvVars:     []string{"DB_PASSWORD"},
		},
		&cli.StringFlag{
			Name:        flagNameDbPort,
			Aliases:     nil,
			Value:       "5432",
			Usage:       "Postgres DB port",
			Destination: &copyArgs.dbPort,
			EnvVars:     []string{"DB_PORT"},
		},
		&cli.StringFlag{
			Name:        flagNameDbName,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB name. Example: author, public",
			Destination: &copyArgs.dbName,
			EnvVars:     []string{"DB_NAME"},
		},
		&cli.StringFlag{
			Name:        flagNameDbDumpFile,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres DB dump file path",
			Destination: &copyArgs.dbDumpPath,
			EnvVars:     []string{"DB_DUMP_PATH"},
		},
		&cli.StringSliceFlag{
			Name:        flagNameMgnlWorkspaces,
			Aliases:     nil,
			Value:       nil,
			Usage:       "Magnolia workspaces will be copyd",
			Destination: &copyArgs.mgnlWorkspaces,
			EnvVars:     []string{"MGNL_WORKSPACES"},
		},
		&cli.BoolFlag{
			Name:        flagNameNormalizeTableName,
			Aliases:     nil,
			Value:       true,
			Usage:       "Auto normalize table name. Convert exported workspace name to db table name.",
			Destination: &copyArgs.normalizeTableName,
			EnvVars:     []string{"NORMALIZE_TABLE_NAME"},
		},
		&cli.BoolFlag{
			Name:        flagNameCopyDatastore,
			Aliases:     nil,
			Value:       true,
			Usage:       "Copy ds_datastore table",
			Destination: &copyArgs.copyDatastore,
			EnvVars:     []string{"COPY_DATASTORE"},
		},
	}

	connectionFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        flagNamePgVersion,
			Aliases:     nil,
			Value:       "",
			Usage:       "Postgres version to use for workspace data copy. Support version: 11, 12. If not set the default Postgres version 15 will be used",
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
			Name:  "copy",
			Usage: "run db workspace copy",
			Flags: copyFlags,
			Action: func(ctx *cli.Context) error {
				if err := copyArgs.validate(); err != nil {
					return fmt.Errorf("invalid argument: %w", err)
				}
				return copyWorkspaces(ctx.Context, copyArgs)
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
