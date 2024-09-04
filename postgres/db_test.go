package postgres_test

import (
	"fmt"
	"os"
	"os/exec"
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
	MustCloseBD(t, db)
}

func MustOpenDB(tb testing.TB) *postgres.DB {
	tb.Helper()

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
	if err := dockerCmd.Run(); err != nil {
		tb.Fatalf("failed to start db container server: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbName)
	db := postgres.NewDB(dsn)
	var err error

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
	fmt.Println(mgCmd.String())
	mgCmdOuput, err := mgCmd.CombinedOutput()
	if err != nil {
		tb.Fatalf("failed to run migrations: %v, output: %s", err, string(mgCmdOuput))
	}

	return db
}

func MustCloseBD(tb testing.TB, db *postgres.DB) {
	tb.Helper()

	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}

	stopCmd := exec.Command("docker", "stop", containerName)
	if err := stopCmd.Run(); err != nil {
		tb.Fatalf("failed to stop db conatiner server: %v", err)
	}

	rmCmd := exec.Command("docker", "rm", containerName)
	if err := rmCmd.Run(); err != nil {
		tb.Fatalf("failed to remove db container: %v", err)
	}
}
