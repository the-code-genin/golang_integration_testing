package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) Repository {
	return &repository{conn}
}

func (r *repository) CreateNote(ctx context.Context, dto CreateNoteDTO) (*Note, error) {
	// Generate an ID and timestamp for the note
	id := uuid.New()
	createdAt := time.Now()

	// Insert the note into the database
	_, err := r.conn.Exec(
		ctx,
		`INSERT INTO core.notes (id, title, description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $4)`,
		id.String(), dto.Title, dto.Description, createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Return the created note
	return &Note{
		ID:          id,
		Title:       dto.Title,
		Description: dto.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   &createdAt,
	}, nil
}

func (r *repository) UpdateNote(ctx context.Context, id uuid.UUID, dto UpdateNoteDTO) (*Note, error) {
	args := []any{id.String()}
	setClauses := []string{}

	// Always update updated_at
	args = append(args, time.Now())
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(args)))

	if dto.Title != nil {
		args = append(args, *dto.Title)
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", len(args)))
	}

	if dto.Description != nil {
		args = append(args, *dto.Description)
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
	}

	_, err := r.conn.Exec(
		ctx,
		fmt.Sprintf(`UPDATE core.notes SET %s WHERE id = $1`, strings.Join(setClauses, ", ")),
		args...,
	)
	if err != nil {
		return nil, err
	}

	return r.FetchNoteByID(ctx, id)
}

func (r *repository) DeleteNote(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `DELETE FROM core.notes WHERE id = $1`, id.String())
	return err
}

func (r *repository) FetchNotes(ctx context.Context) ([]Note, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, title, description, created_at, updated_at
		FROM core.notes
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var updatedAt time.Time
		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		note.UpdatedAt = &updatedAt
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *repository) FetchNoteByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	var note Note
	var updatedAt time.Time

	err := r.conn.QueryRow(ctx, `
		SELECT id, title, description, created_at, updated_at
		FROM core.notes
		WHERE id = $1
	`, id.String()).Scan(&note.ID, &note.Title, &note.Description, &note.CreatedAt, &updatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	note.UpdatedAt = &updatedAt
	return &note, nil
}
