package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/the-code-genin/golang_integration_testing/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo}
}

func (s *service) CreateNote(
	ctx context.Context, dto repository.CreateNoteDTO,
) (*repository.Note, error) {
	note, err := s.repo.CreateNote(ctx, dto)
	if err != nil {
		log.Printf("an error occurred while creating a note: %v", err)
		return nil, ErrInternal
	}
	return note, nil
}

func (s *service) UpdateNote(
	ctx context.Context, id uuid.UUID, dto repository.UpdateNoteDTO,
) (*repository.Note, error) {
	return &repository.Note{}, nil
}

func (s *service) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *service) FetchNotes(ctx context.Context) ([]repository.Note, error) {
	return []repository.Note{}, nil
}

func (s *service) FetchNoteByID(
	ctx context.Context, id uuid.UUID,
) (*repository.Note, error) {
	return &repository.Note{}, nil
}
