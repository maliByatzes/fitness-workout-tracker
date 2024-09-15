package postgres_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/maliByatzes/fwt/postgres"
)

const (
	containerName = "test-postgres"
	port          = "5433"
	image         = "postgres:16-alpine"
	user          = "testuser"
	password      = "testpassword"
	dbName        = "testdb"
)

func TestDB(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

func MustOpenDB(tb testing.TB) *postgres.DB {
	tb.Helper()

	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbName)
	db := postgres.NewDB(dsn)

	containerCmd := exec.Command("docker", "container", "inspect", "-f", `'{{.State.Running}}'`, containerName)
	containerCmdOutput, err := containerCmd.CombinedOutput()
	cntCmdOutputStr := string(containerCmdOutput)
	if err != nil {
		if strings.Contains(cntCmdOutputStr, "Error response from daemon: No such container:") {
			dockerCmd := exec.Command(
				"docker",
				"run",
				"--name",
				containerName,
				"-p",
				fmt.Sprintf("%s:5432", port),
				"-e",
				fmt.Sprintf("POSTGRES_USER=%s", user),
				"-e",
				fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
				"-e",
				fmt.Sprintf("POSTGRES_DB=%s", dbName),
				"-d",
				image)
			dockerCmdOutput, err := dockerCmd.CombinedOutput()
			if err != nil {
				tb.Fatalf("failed to start db container server: %v, output: %s", err, string(dockerCmdOutput))
			}

			for i := 0; i < 10; i++ {
				err = db.Open()
				if err == nil {
					err = db.DB.Ping()
				}

				if err == nil {
					break
				}

				time.Sleep(2 * time.Second)
			}
			if err != nil {
				tb.Fatalf("failed to connect to db: %v", err)
			}

			wd, _ := os.Getwd()
			mgCmd := exec.Command("migrate", "-database", fmt.Sprintf(`%s`, dsn), "-path", wd+"/migrations", "up")
			mgCmdOuput, err := mgCmd.CombinedOutput()
			if err != nil {
				tb.Fatalf("failed to run migrations: %v, output: %s", err, string(mgCmdOuput))
			}

			return db
		} else {
			tb.Fatalf("failed to run container check cmd: %v", err)
		}
	}

	newCntStr := strings.Trim(cntCmdOutputStr, "\n")
	if newCntStr == `'false'` {
		startCmd := exec.Command("docker", "start", containerName)
		startCmdOutput, err := startCmd.CombinedOutput()
		if err != nil {
			tb.Fatalf("failed to start docker container: %d, output: %s", err, startCmdOutput)
		}
	}

	for i := 0; i < 10; i++ {
		err = db.Open()
		if err == nil {
			err = db.DB.Ping()
		}

		if err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}
	if err != nil {
		tb.Fatalf("failed to connect to db: %v", err)
	}

	wd, _ := os.Getwd()
	mgCmd := exec.Command("migrate", "-database", fmt.Sprintf(`%s`, dsn), "-path", wd+"/migrations", "up")
	mgCmdOuput, err := mgCmd.CombinedOutput()
	if err != nil {
		tb.Fatalf("failed to run up migrations: %v, output: %s", err, string(mgCmdOuput))
	}

	return db
}

func MustCloseDB(tb testing.TB, db *postgres.DB) {
	tb.Helper()

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbName)
	wd, _ := os.Getwd()
	mgCmd := exec.Command("migrate", "-database", fmt.Sprintf(`%s`, dsn), "-path", wd+"/migrations", "down", "-all")
	mgCmdOuput, err := mgCmd.CombinedOutput()
	if err != nil {
		tb.Fatalf("failed to run down migrations: %v, output: %s", err, string(mgCmdOuput))
	}

	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}

	stopCmd := exec.Command("docker", "stop", containerName)
	if err := stopCmd.Run(); err != nil {
		tb.Fatalf("failed to stop db conatiner server: %v", err)
	}
}
