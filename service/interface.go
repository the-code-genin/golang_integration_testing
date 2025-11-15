//go:generate mockgen -source=interface.go -destination interface_mock.go -package service . Service

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/the-code-genin/golang_integration_testing/repository"
)

type Service interface {
	CreateNote(ctx context.Context, dto repository.CreateNoteDTO) (*repository.Note, error)
	UpdateNote(ctx context.Context, id uuid.UUID, dto repository.UpdateNoteDTO) (*repository.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) error

	FetchNotes(ctx context.Context) ([]repository.Note, error)
	FetchNoteByID(ctx context.Context, id uuid.UUID) (*repository.Note, error)
}
