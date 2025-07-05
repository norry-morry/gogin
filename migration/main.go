// Package main はデータベースのマイグレーション処理を行うエントリーポイントです。
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"
)

var migrationFilePath = "file://./migrations/"

func main() {
	log.Println("start migration:", time.Now())
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Host := os.Getenv("DB_HOST")
	log.Println("DB_HOST:", Host)
	flag.Parse()
	command := flag.Arg(0)
	migrationFileName := flag.Arg(1)

	if command == "" {
		showUsage()
		os.Exit(1)
	}

	m := newMigrate()
	//version, dirty, _ := m.Version()
	version, dirty, versionErr := m.Version()
	if versionErr != nil {
		log.Fatalf("failed to get migration version: %v", versionErr)
	}
	force := flag.Bool("f", false, "force execute fixed sql")
	if dirty && !*force {
		log.Println("force=true: force execute current version sql")
		//m.Force(int(version))
		if err := m.Force(int(version)); err != nil {
			log.Fatal(err) // 必要に応じて他のハンドリングでもOK
		}
	}

	switch command {
	case "new":
		newMigration(migrationFileName)
	case "up":
		up(m)
	case "down":
		down(m)
	case "drop":
		drop(m)
	case "version":
		showVersionInfo(m.Version())
	default:
		log.Println("\nerror: invalid command '", command, "'")
		showUsage()
		os.Exit(0)
	}
}

func generateDsn() string {
	var dsn string

	User := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASS")
	Host := os.Getenv("DB_HOST")
	Port := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", User, Password, Host, Port, DBName)

	return dsn
}

func newMigrate() *migrate.Migrate {
	dsn := generateDsn()
	db, openErr := sql.Open("mysql", dsn)
	if openErr != nil {
		log.Println(errors.Wrap(openErr, "error occurred. sql.Open()"))
		os.Exit(1)
	}

	driver, instanceErr := mysql.WithInstance(db, &mysql.Config{})
	if instanceErr != nil {
		log.Println(errors.Wrap(instanceErr, "error occurred. mysql.WithInstance()"))
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationFilePath,
		"mysql",
		driver,
	)

	if err != nil {
		log.Println(errors.Wrap(err, "error occurred. migrate.NewWithDatabaseInstance()"))
		os.Exit(1)
	}
	return m
}

func showUsage() {
	log.Println(`
-------------------------------------
Usage:
  go run migration/main.go <command>
Commands:
  new FILENAME	Create new up & down migration files
  up		Apply up migrations
  down		Apply down migrations
  drop		Drop everything
  version	Check current migrate version
-------------------------------------`)
}

func newMigration(name string) {
	if name == "" {
		log.Println("\nerror: migration file name must be supplied as an argument")
		os.Exit(1)
	}
	base := fmt.Sprintf("./migrations/%s_%s", time.Now().Format("20060102030405"), name)
	ext := ".sql"
	createFile(base + ".up" + ext)
	createFile(base + ".down" + ext)
}

func createFile(fname string) {
	if _, err := os.Create(fname); err != nil {
		panic(err)
	}
}

func up(m *migrate.Migrate) {
	log.Println("Before:")
	showVersionInfo(m.Version())
	err := m.Up()
	if err != nil {
		if err.Error() != "no change" {
			panic(err)
		}
		log.Println("\nno change")
	} else {
		log.Println("\nUpdated:")
		version, dirty, err := m.Version()
		showVersionInfo(version, dirty, err)
	}
}

func down(m *migrate.Migrate) {
	log.Println("Before:")
	showVersionInfo(m.Version())
	err := m.Steps(-1)
	if err != nil {
		panic(err)
	} else {
		log.Println("\nUpdated:")
		showVersionInfo(m.Version())
	}
}

func drop(m *migrate.Migrate) {
	err := m.Drop()
	if err != nil {
		panic(err)
	} else {
		log.Println("Dropped all migrations")
		return
	}
}

func showVersionInfo(version uint, dirty bool, err error) {
	log.Println("-------------------")
	log.Println("version : ", version)
	log.Println("dirty   : ", dirty)
	log.Println("error   : ", err)
	log.Println("-------------------")
}
