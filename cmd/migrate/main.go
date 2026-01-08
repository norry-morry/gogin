// Package main はデータベースのマイグレーション処理を行うエントリーポイントです。
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"gorm.io/gorm/logger"

	"resume/internal/bootstrap"
	"resume/internal/config"
	infradb "resume/internal/infra/db"
)

// var migrationFilePath = "file://../../migrations"
// 追加: グローバルフラグは Parse 前に定義
var migDirFlag = flag.String("dir", "", "path to migrations dir")
var force = flag.Bool("f", false, "force execute current version sql (when dirty)")
var truncateAllow = []string{
	"user_profiles",
	// "users",               // ← ユーザー自体も消すならここに
	// 他、テストデータ用テーブルを列挙
}
var seedDirFlag = flag.String("seeddir", "", "path to seed dir")

func main() {
	log.Println("start migration:", time.Now())
	bootstrap.LoadDotEnv()

	flag.Parse()
	command := flag.Arg(0)
	migrationFileName := flag.Arg(1)
	seedFileName := flag.Arg(1)

	if command == "" {
		showUsage()
		os.Exit(1)
	}

	// まず new だけは db 接続不要。先に分岐して終了する
	if command == "new" {
		newMigration(migrationFileName)
		return
	}

	if command == "newfixture" {
		newFixture(migrationFileName)
		return
	}

	if command == "newseed" {
		newSeedFile(seedFileName)
		return
	}

	// ここからは migrate を初期化
	m := newMigrate()

	version, dirty, versionErr := m.Version()
	if versionErr != nil {
		if errors.Is(versionErr, migrate.ErrNilVersion) {
			log.Println("version : 0 (nil version)")
		} else {
			log.Fatalf("failed to get migration version: %v", versionErr)
		}
	} else {
		log.Printf("version : %d, dirty : %v\n", version, dirty)
	}

	if dirty {
		if *force {
			log.Println("force=true: force execute current version sql")
			if err := m.Force(int(version)); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("database is dirty at version %d (use -f to force)", version)
		}
	}

	switch command {
	case "up":
		up(m)
	case "down":
		down(m)
	case "drop":
		drop(m)
	case "version":
		showVersionInfo(m.Version())
	case "seedup":
		seedUp()
	case "seedreset":
		seedReset()
	case "seeddrop":
		seedDrop()
	default:
		log.Println("\nerror: invalid command '", command, "'")
		showUsage()
		os.Exit(0)
	}
}

// ここを "file://../../migrations" 固定ではなく、解決関数を使って生成
func newMigrate() *migrate.Migrate {
	cfg := config.Load()

	gdb, err := infradb.NewGorm(cfg, logger.Default.LogMode(logger.Info))
	if err != nil {
		log.Fatal(errors.Wrap(err, "infradb.NewGorm"))
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatal(errors.Wrap(err, "gorm.db().db()"))
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "mysql.WithInstance"))
	}

	dir := resolveMigrationsDir()
	sourceURL := "file://" + dir

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "mysql", driver)
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrate.NewWithDatabaseInstance"))
	}
	return m
}

func showUsage() {
	log.Println(`
-------------------------------------
Usage:
  go run migration/main.go <command>
Commands:
  new FILENAME		Create new up & down migration files
  up				Apply up migrations
  down				Apply down migrations
  drop				Drop everything
  version			Check current migrate version
  newseed FILENAME	Create new seed SQL file (timestamped)
  seedup			Execute all seed SQL files in order
  seedreset			Truncate registered seed tables, then reapply seeds
  seeddrop			Truncate registered seed tables only (delete all seed data)
-------------------------------------`)
}

// 既存の newMigration はディレクトリ解決だけ共通化
func newMigration(name string) {
	if name == "" {
		log.Fatal("missing migration file name (e.g., `go run ./cmd/migrate new countries`)")
	}
	dir := resolveMigrationsDir()

	ts := time.Now().In(time.FixedZone("JST", 9*60*60)).Format("20060102150405")
	upPath := filepath.Join(dir, fmt.Sprintf("%s_schema_%s.up.sql", ts, name))
	downPath := filepath.Join(dir, fmt.Sprintf("%s_schema_%s.down.sql", ts, name))

	createFile(upPath)
	createFile(downPath)
	log.Printf("created:\n  %s\n  %s\n", upPath, downPath)
}

// newFixture は マスタ系テーブルの挿入用データを作成する
func newFixture(name string) {
	if name == "" {
		log.Fatal("missing migration file name (e.g., `go run ./cmd/migrate newfixture countries`)")
	}
	dir := resolveMigrationsDir()

	ts := time.Now().In(time.FixedZone("JST", 9*60*60)).Format("20060102150405")
	upPath := filepath.Join(dir, fmt.Sprintf("%s_fixture_%s.up.sql", ts, name))
	downPath := filepath.Join(dir, fmt.Sprintf("%s_fixture_%s.down.sql", ts, name))

	createFile(upPath)
	createFile(downPath)
	log.Printf("created:\n  %s\n  %s\n", upPath, downPath)
}

// newSeedFile は 新たにseeder用のsqlファイルを生成
func newSeedFile(name string) {
	if name == "" {
		log.Fatal("missing seed file name (e.g., `go run ./cmd/migrate newseed countries`)")
	}
	dir := resolveSeedDir()

	ts := time.Now().In(time.FixedZone("JST", 9*60*60)).Format("20060102_150405")
	fullpath := filepath.Join(dir, fmt.Sprintf("%s_%s.sql", ts, name))

	createFile(fullpath)
	log.Printf("created:\n  %s\n", fullpath)
}

func createFile(p string) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(f)
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

// dirty は **golang-migrate が管理している「マイグレーションの整合性フラグ」**です。
func showVersionInfo(version uint, dirty bool, err error) {
	log.Println("-------------------")
	log.Println("version : ", version)
	log.Println("dirty   : ", dirty)
	log.Println("error   : ", err)
	log.Println("-------------------")
}

// 追加: マイグレーションディレクトリ解決（-dir > 環境変数 > 自動検出）
func resolveMigrationsDir() string {
	// 1) フラグ・環境変数 優先
	if *migDirFlag != "" {
		return *migDirFlag
	}
	if v := os.Getenv("MIGRATIONS_DIR"); v != "" {
		return v
	}

	// 2) このソース(main.go)の物理パス基準で推定
	if _, filename, _, ok := runtime.Caller(0); ok {
		base := filepath.Dir(filename)
		candidates := []string{
			filepath.Clean(filepath.Join(base, "../../migrations")), // /go/src/migrations
			filepath.Clean(filepath.Join(base, "../migrations")),    // 念のため
		}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && st.IsDir() {
				return c
			}
		}
	}

	// 3) 最後のフォールバック
	fallbacks := []string{
		"./migrations",
		"/go/src/migrations",
	}
	for _, c := range fallbacks {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}

	log.Fatal("migrations dir not found. set -dir or MIGRATIONS_DIR, or ensure ./migrations exists relative to repo root")
	return ""
}

func resolveSeedDir() string {
	// 1) フラグ・環境変数 優先
	if *seedDirFlag != "" {
		return *seedDirFlag
	}
	if v := os.Getenv("SEED_DIR"); v != "" {
		return v
	}

	// 2) このソース(main.go)の物理パス基準で推定
	if _, filename, _, ok := runtime.Caller(0); ok {
		base := filepath.Dir(filename)
		candidates := []string{
			filepath.Clean(filepath.Join(base, "../../seeders")), // /go/src/migrations
			filepath.Clean(filepath.Join(base, "../seeders")),    // 念のため
		}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && st.IsDir() {
				return c
			}
		}
	}

	// 3) 最後のフォールバック
	fallbacks := []string{
		"./seeders",
		"/go/src/seeders",
	}
	for _, c := range fallbacks {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}

	log.Fatal("seed dir not found. set -dir or SEED_DIR, or ensure ./seeders exists relative to repo root")
	return ""
}

func openSQL() *sql.DB {
	cfg := config.Load()
	gdb, err := infradb.NewGorm(cfg, logger.Default.LogMode(logger.Warn))
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

func seedUp() {
	dir := resolveSeedDir() // 既に実装済（good）
	db := openSQL()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// *.sql を昇順ソート（タイムスタンプ接頭語を想定）
	var paths []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".sql") {
			paths = append(paths, filepath.Join(dir, f.Name()))
		}
	}
	sort.Strings(paths)

	for _, p := range paths {
		log.Printf(">> applying seed: %s\n", filepath.Base(p))
		b, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}

		// 1ファイル＝1トランザクション
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		if err := execSQLBatch(tx, string(b)); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				log.Printf("failed to rollback: %v", rbErr)
			}
			log.Fatalf("seed failed on %s: %v", filepath.Base(p), err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("seedup done.")
}

func execSQLBatch(tx *sql.Tx, src string) error {
	// 超シンプル：; 区切りで分割（ストアド等は非対応）
	stmts := strings.Split(src, ";")
	for _, s := range stmts {
		q := strings.TrimSpace(s)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func seedDrop() {
	db := openSQL()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()
	for _, t := range truncateAllow {
		log.Printf("TRUNCATE %s", t)
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec("TRUNCATE TABLE " + t); err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("seeddrop done.")
}

func seedReset() {
	seedDrop()
	seedUp()
}

// note これはwikiにまとめた方がいいかな
//		🧭 dirty の意味
//		状態	意味	どういう時になるか
//		false	正常終了	マイグレーションが最後まで成功した
//		true	途中失敗・中断状態	実行中にエラー・強制停止があり、完全には反映されていない
