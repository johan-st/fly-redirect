package repo

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type TursoRepo struct {
	db *sql.DB
}

func NewTursoRepo() (*TursoRepo, error) {
	dbURL := os.Getenv("DB_TURSO_URL")
	dbToken := os.Getenv("DB_TURSO_TOKEN")

	if dbURL == "" || dbToken == "" {
		log.Fatal("Missing required database environment variables")
	}

	db, err := sql.Open("libsql", dbURL+"?authToken="+dbToken)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &TursoRepo{
		db: db,
	}

	if err := repo.RunMigrations(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (t *TursoRepo) Close() error {
	return t.db.Close()
}

// HealthCheck returns an error if the database is not reachable
func (t *TursoRepo) HealthCheck() error {
	return t.db.Ping()
}

func (t *TursoRepo) GetRedirectCount() (int, error) {
	var count int
	err := t.db.QueryRow("SELECT count FROM redirects WHERE id = 1").Scan(&count)
	return count, err
}

func (t *TursoRepo) IncrementRedirectCount() error {
	result, err := t.db.Exec("UPDATE redirects SET count = count + 1 WHERE id = 1")
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

// MIGRATIONS

func (t *TursoRepo) RunMigrations() error {
	_, err := t.db.Exec(`
		CREATE TABLE IF NOT EXISTS redirects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			count INTEGER DEFAULT 0
		);
		
		INSERT OR IGNORE INTO redirects (id, count) VALUES (1, 0);
	`)
	return err
}

/**
 * Schema_v1
 *
 * CREATE TABLE redirects (
 * 	id INTEGER PRIMARY KEY AUTOINCREMENT,
 * 	count INTEGER DEFAULT 0
 * );
 */
