package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	datastoreTable = "ds_datastore"
	binPathPattern = "/usr/lib/postgresql/%s/bin"
)

func copyWorkspaces(ctx context.Context, args copyArgs) error {
	logrus.Info("run db copy: args=", args)

	for _, tableName := range args.mgnlWorkspaces.Value() {
		tn := tableName
		if args.normalizeTableName {
			tn = normalizeTableName(tableName)
		}
		logrus.Info("copy table=", tn)

		if err := checkTableExists(tn, args); err != nil {
			logrus.Warn("table=", tn, " does not exists: ", err)
		} else {
			// if table exists, clean up data in the table
			if err := emptyTable(tn, args); err != nil {
				return err
			}
		}

		// this assumes the table exists in the dump file
		// TODO: Check if table exists in the dump file before doing restore!
		logrus.Info("restore data from dump file: table=", tn)
		cmd := exec.Command(
			resolvePath(binPathPattern, args.pgversion, "pg_restore"),
			"--verbose",
			"--data-only",
			"-h", args.dbHost,
			"-p", args.dbPort,
			"-U", args.dbUsername,
			"-d", args.dbName,
			"-t", tn,
			"--jobs", "10",
			args.dbDumpPath,
		)
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", args.dbPassword))

		if err := runCommand(cmd); err != nil {
			return err
		}
	}

	if args.copyDatastore {
		logrus.Info("copy datastore")

		if err := emptyTable(datastoreTable, args); err != nil {
			return err
		}

		cmd := exec.Command(
			resolvePath(binPathPattern, args.pgversion, "pg_restore"),
			"--verbose",
			"--data-only",
			"-h", args.dbHost,
			"-p", args.dbPort,
			"-U", args.dbUsername,
			"-d", args.dbName,
			"-t", datastoreTable,
			"--jobs", "10",
			args.dbDumpPath,
		)
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", args.dbPassword))

		if err := runCommand(cmd); err != nil {
			return err
		}
	}

	if args.showVersionWarning {
		if !slices.Contains(args.mgnlWorkspaces.Value(), "magnolia-mgnlversion") {
			logrus.Warn("pm_mgnlversion_bundle workspace may be missing")
		}
		if !slices.Contains(args.mgnlWorkspaces.Value(), "magnolia_conf_sec-mgnlVersion") {
			logrus.Warn("version_bundle workspace may be missing")
		}
	}

	return nil
}

func resolvePath(pathPattern string, postgresVersion string, tool string) string {
	return filepath.Join(fmt.Sprintf(pathPattern, postgresVersion), tool)
}

func checkTableExists(tableName string, args copyArgs) error {
	logrus.Info("check if table exists: table=", tableName)
	cmd := exec.Command(
		resolvePath(binPathPattern, args.pgversion, "psql"),
		"-h", args.dbHost,
		"-p", args.dbPort,
		"-U", args.dbUsername,
		"-d", args.dbName,
		"-c", fmt.Sprintf("SELECT * from %s LIMIT 1", tableName),
	)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", args.dbPassword))

	if err := runCommand(cmd); err != nil {
		return err
	}

	return nil
}

func emptyTable(tableName string, args copyArgs) error {
	logrus.Info("empty data: table=", tableName)
	cmd := exec.Command(
		resolvePath(binPathPattern, args.pgversion, "psql"),
		"-h", args.dbHost,
		"-p", args.dbPort,
		"-U", args.dbUsername,
		"-d", args.dbName,
		"-c", fmt.Sprintf("TRUNCATE TABLE %s;", tableName),
	)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", args.dbPassword))

	if err := runCommand(cmd); err != nil {
		return err
	}
	return nil
}

func normalizeTableName(tableName string) string {
	logrus.Info("normalize table=", tableName)
	if strings.EqualFold(tableName, "magnolia-mgnlversion") {
		return "pm_mgnlversion_bundle"
	}
	if strings.EqualFold(tableName, "magnolia_conf_sec-mgnlVersion") {
		return "version_bundle"
	}

	tn := strings.Replace(tableName, "-", "_x002d_", -1)

	return fmt.Sprintf("pm_%s_bundle", tn)
}

func runCommand(cmd *exec.Cmd) error {
	if cmd == nil {
		return fmt.Errorf("command pointer cannot be nil")
	}
	logrus.Info("command is ", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/tmp"
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command %s: %w", cmd.String(), err)
	}
	logrus.Info("successfully ran command: ", cmd.String())
	return nil
}
