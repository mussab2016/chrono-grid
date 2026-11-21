package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DBEngine struct {
	postgres *embeddedpostgres.EmbeddedPostgres
	Client   *sql.DB
}

func InitDB(port uint32) (*DBEngine, error) {
	// Define local data paths within the project directory
	dataPath := filepath.Join(".", "data")
	binPath := filepath.Join(".", "bin", "embedded-postgres")
	// psql -h localhost -p 5433 -U unisched_user -d unisched_db  psql postgresql://unisched_user:unisched_pass@localhost:5433/unisched_db
	user := "unisched_user"
	pass := "unisched_pass"
	dbName := "unisched_db"
	// Set up embedded PostgreSQL engine
	config := embeddedpostgres.DefaultConfig().
		Username(user).
		Password(pass).
		Database(dbName).
		Port(port). // Use a non-default port to prevent conflicts with any previously installed Postgres
		DataPath(dataPath).
		RuntimePath(binPath)

	pgEngine := embeddedpostgres.NewDatabase(config)

	// Start the database engine
	log.Println("Cleaning port and lock files...")
	killPortWindows(int(port))
	time.Sleep(200 * time.Millisecond)

	log.Println("Starting embedded PostgreSQL engine...")
	if err := pgEngine.Start(); err != nil {
		return nil, fmt.Errorf("failed to start embedded PostgreSQL: %v", err)
	}

	// Connect to the database using the local connection string

	connStr := fmt.Sprintf(
		"host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		port, user, pass, dbName,
	)

	dbClient, err := sql.Open("postgres", connStr)
	if err != nil {
		pgEngine.Stop()
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := dbClient.Ping(); err != nil {
		pgEngine.Stop()
		return nil, fmt.Errorf("database doesn't respond: %v", err)
	}

	log.Println("connecting successfully to the database...")

	dbEngine := &DBEngine{
		postgres: pgEngine,
		Client:   dbClient,
	}

	// Execute pending migration files automatically on startup
	if err := dbEngine.runMigrations(dbName); err != nil {
		dbEngine.Close()
		return nil, fmt.Errorf("migration execution failed: %w", err)
	}

	log.Println("Database initialized and migrations applied successfully!")
	return dbEngine, nil
}

func killPortWindows(port int) {
	psScript := fmt.Sprintf(
		"Get-NetTCPConnection -LocalPort %d -ErrorAction SilentlyContinue | Where-Object { $_.OwningProcess -gt 0 } | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue }",
		port,
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	_ = cmd.Run()
}

// runMigrations scans db/migrations directory and applies up.sql scripts
func (d *DBEngine) runMigrations(db string) error {
	driver, err := postgres.WithInstance(d.Client, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Path to local migration files
	absPath, err := filepath.Abs("server/db/migrations")
	if err != nil {
		return fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}
	migrationPath := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))

	m, err := migrate.NewWithDatabaseInstance(migrationPath, db, driver)

	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	// Apply all new migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// handle the case of being dirty
		version, dirty, vErr := m.Version()
		if vErr == nil && dirty {
			log.Printf("Database is in dirty state at version %d. Forcing version cleanup...\n", version)
			// unlock and install the current version to try again
			_ = m.Force(int(version))
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	//if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	//	return fmt.Errorf("failed to apply migrations: %w", err)
	//}

	log.Println("Database schema is up to date.")
	return nil
}

// Close gracefully stops DB connection and postgres process
func (d *DBEngine) Close() {
	if d.Client != nil {
		_ = d.Client.Close()
	}
	if d.postgres != nil {
		log.Println("Stopping embedded PostgreSQL server...")
		_ = d.postgres.Stop()
	}
}

func (d *DBEngine) GenerateData() {
	seedBytes, err := os.ReadFile("./server/db/seed/seed.sql")
	if err != nil {
		log.Fatalf("Failed to read seed.sql: %v", err)
	}
	if _, err := d.Client.Exec(string(seedBytes)); err != nil {
		log.Fatalf("Failed to execute seed.sql: %v", err)
	}
	log.Println("Database seeded successfully!")
}
