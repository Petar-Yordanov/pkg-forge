package engine

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Petar-Yordanov/pkg-forge/manifest/parser"
	_ "modernc.org/sqlite"
)

type StateRow struct {
	ManifestName   string
	Ordinal        int
	EntryKind      string
	PackageName    string
	Version        string
	PackageManager string
	StepHash       string
}

type StateStore struct {
	db *sql.DB
}

func OpenStateStore(path string) (*StateStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &StateStore{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *StateStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *StateStore) init() error {
	schema := `
CREATE TABLE IF NOT EXISTS manifest_state (
	manifest_name TEXT NOT NULL,
	ordinal INTEGER NOT NULL,
	entry_kind TEXT NOT NULL,
	package_name TEXT NOT NULL,
	version TEXT NOT NULL,
	package_manager TEXT NOT NULL,
	step_hash TEXT NOT NULL,
	PRIMARY KEY (manifest_name, ordinal)
);

CREATE INDEX IF NOT EXISTS idx_manifest_state_manifest_name
ON manifest_state(manifest_name);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *StateStore) LoadManifest(manifestName string) ([]StateRow, error) {
	rows, err := s.db.Query(`
SELECT manifest_name, ordinal, entry_kind, package_name, version, package_manager, step_hash
FROM manifest_state
WHERE manifest_name = ?
ORDER BY ordinal ASC
`, manifestName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StateRow
	for rows.Next() {
		var row StateRow
		if err := rows.Scan(
			&row.ManifestName,
			&row.Ordinal,
			&row.EntryKind,
			&row.PackageName,
			&row.Version,
			&row.PackageManager,
			&row.StepHash,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *StateStore) ReplaceManifest(manifestName string, rows []StateRow) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM manifest_state WHERE manifest_name = ?`, manifestName); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
INSERT INTO manifest_state (
	manifest_name,
	ordinal,
	entry_kind,
	package_name,
	version,
	package_manager,
	step_hash
) VALUES (?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range rows {
		if _, err := stmt.Exec(
			row.ManifestName,
			row.Ordinal,
			row.EntryKind,
			row.PackageName,
			row.Version,
			row.PackageManager,
			row.StepHash,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func BuildStateRow(manifestName string, ordinal int, e parser.Entry) (StateRow, error) {
	hash, err := HashEntry(e)
	if err != nil {
		return StateRow{}, err
	}

	return StateRow{
		ManifestName:   manifestName,
		Ordinal:        ordinal,
		EntryKind:      e.Kind,
		PackageName:    e.Name,
		Version:        e.Version,
		PackageManager: e.PackageManager,
		StepHash:       hash,
	}, nil
}

func HashEntry(e parser.Entry) (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("marshal entry for hashing: %w", err)
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
