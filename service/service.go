package service

import (
	"context"
	"database/sql"
	"errors"
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
	note, err := s.repo.UpdateNote(ctx, id, dto)
	if err != nil {
		log.Printf("an error occurred while updating note with id %s: %v", id.String(), err)
		return nil, ErrInternal
	}
	return note, nil
}

func (s *service) DeleteNote(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteNote(ctx, id)
	if err != nil {
		log.Printf("an error occurred while deleting note with id %s: %v", id.String(), err)
		return ErrInternal
	}
	return nil
}

func (s *service) FetchNotes(ctx context.Context) ([]repository.Note, error) {
	notes, err := s.repo.FetchNotes(ctx)
	if err != nil {
		log.Printf("an error occurred while fetching notes: %v", err)
		return nil, ErrInternal
	}
	return notes, nil
}

func (s *service) FetchNoteByID(
	ctx context.Context, id uuid.UUID,
) (*repository.Note, error) {
	note, err := s.repo.FetchNoteByID(ctx, id)
	if err != nil {
		log.Printf("an error occurred while fetching note with id %s: %v", id.String(), err)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoteNotFound
		default:
			return nil, ErrInternal
		}
	}
	return note, nil
}
