package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	return &Note{}, nil
}

func (r *repository) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *repository) FetchNotes(ctx context.Context) ([]Note, error) {
	return []Note{}, nil
}

func (r *repository) FetchNoteByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	return &Note{}, nil
}
