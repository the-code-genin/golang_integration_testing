//go:generate mockgen -source=interface.go -destination interface_mock.go -package repository . Repository

package repository

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateNote(ctx context.Context, dto CreateNoteDTO) (*Note, error)
	UpdateNote(ctx context.Context, id uuid.UUID, dto UpdateNoteDTO) (*Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error

	FetchNotes(ctx context.Context) ([]Note, error)
	FetchNoteByID(ctx context.Context, id uuid.UUID) (*Note, error)
}

type CreateNoteDTO struct {
	Title, Description string
}

type UpdateNoteDTO struct {
	Title, Description *string
}
