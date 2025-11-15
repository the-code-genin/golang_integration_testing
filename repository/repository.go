package repository

import (
	"context"

	"github.com/google/uuid"
)

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateNote(ctx context.Context, dto CreateNoteDTO) (*Note, error) {
	return &Note{}, nil
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
