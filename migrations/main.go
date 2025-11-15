package migrations

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
)

//go:embed *.up.sql
var upMigrationsFS embed.FS

//go:embed *.down.sql
var downMigrationsFS embed.FS

// Read file names in the directory fsys
func readFileNames(fsys embed.FS) ([]string, error) {
	entries, err := fsys.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	fileNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames, nil
}

// Squash the migration files in the directory fsys into a Buffer,
// with the option to squash the files in ascending or descending order.
func squashMigrations(fsys embed.FS, asc bool) (*bytes.Buffer, error) {
	// Read the file names
	fileNames, err := readFileNames(fsys)
	if err != nil {
		return nil, fmt.Errorf("could not list migration files: %w", err)
	}

	// Sort the file names
	if asc {
		sort.Strings(fileNames)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))
	}

	// Squash the files
	buf := bytes.Buffer{}
	for _, fileName := range fileNames {
		content, err := fsys.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		if _, err := buf.WriteString("\n-- " + fileName + "\n"); err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}

		if _, err := buf.Write(content); err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}

	return &buf, nil
}

// Squash the up and down migrations in the migrations folder into Buffers.
func SquashMigrations() (upMigrations, downMigrations *bytes.Buffer, err error) {
	upMigrations, err = squashMigrations(upMigrationsFS, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to squash up migrations: %w", err)
	}

	downMigrations, err = squashMigrations(downMigrationsFS, false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to squash down migrations: %w", err)
	}

	return upMigrations, downMigrations, nil
}
