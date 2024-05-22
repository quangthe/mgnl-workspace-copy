package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func checkConnection(ctx context.Context, args connectionArgs) error {
	logrus.Info("check db connection: args=", args)

	cmdPath := resolvePath(binPathPattern, args.pgversion, "psql")

	if args.pgversion == "" {
		logrus.Info("pgversion is not set, use default psql")
		cmdPath = "/usr/local/bin/psql"
	}

	cmd := exec.Command(
		cmdPath,
		"-h", args.dbHost,
		"-p", args.dbPort,
		"-U", args.dbUsername,
		"-d", args.dbName,
		"-c", "select 1",
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", args.dbPassword))

	if err := runCommand(cmd); err != nil {
		return err
	}

	logrus.Info("db connection is ok")

	return nil
}
