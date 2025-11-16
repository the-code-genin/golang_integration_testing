package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/the-code-genin/golang_integration_testing/repository"
	"go.uber.org/mock/gomock"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	// Setup mocked repo
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepository(ctrl)

	service := NewService(mockRepo)

	t.Run("CreateNote", func(t *testing.T) {
		t.Run("should create a note given a title and description", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			expectedNote := &repository.Note{
				ID:          uuid.New(),
				Title:       dto.Title,
				Description: dto.Description,
			}

			// Expect the repository CreateNote to be called once
			mockRepo.EXPECT().
				CreateNote(gomock.Any(), dto).
				Return(expectedNote, nil)

			note, err := service.CreateNote(ctx, dto)
			assert.NoError(t, err)
			assert.Equal(t, expectedNote, note)
		})

		t.Run("should return ErrNoteTitleTaken if a duplicate key error was returned from the DBAL", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			// Simulate repository returning an error for duplicate title
			mockRepo.EXPECT().
				CreateNote(gomock.Any(), dto).
				Return(nil, errors.New("duplicate key error")).
				Times(1)

			note, err := service.CreateNote(ctx, dto)
			assert.Error(t, err)
			assert.Equal(t, ErrNoteTitleTaken, err)
			assert.Nil(t, note)
		})

		t.Run("should return ErrInternal if an unknown error occurred while creating the note", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			mockRepo.EXPECT().
				CreateNote(gomock.Any(), dto).
				Return(nil, assert.AnError).
				Times(1)

			note, err := service.CreateNote(ctx, dto)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
			assert.Nil(t, note)
		})
	})

	t.Run("FetchNoteByID", func(t *testing.T) {
		id := uuid.New()

		t.Run("should fetch a note given a valid ID", func(t *testing.T) {
			expectedNote := &repository.Note{
				ID:          id,
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			mockRepo.EXPECT().FetchNoteByID(gomock.Any(), id).Return(expectedNote, nil)

			note, err := service.FetchNoteByID(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, expectedNote, note)
		})

		t.Run("should return ErrNoteNotFound if repository returns sql.ErrNoRows", func(t *testing.T) {
			mockRepo.EXPECT().FetchNoteByID(gomock.Any(), id).Return(nil, sql.ErrNoRows)

			note, err := service.FetchNoteByID(ctx, id)
			assert.Error(t, err)
			assert.Equal(t, ErrNoteNotFound, err)
			assert.Nil(t, note)
		})

		t.Run("should return ErrInternal for unknown errors", func(t *testing.T) {
			mockRepo.EXPECT().FetchNoteByID(gomock.Any(), id).Return(nil, assert.AnError)

			note, err := service.FetchNoteByID(ctx, id)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
			assert.Nil(t, note)
		})
	})

	t.Run("UpdateNote", func(t *testing.T) {
		id := uuid.New()
		title := gofakeit.Sentence(3)
		desc := gofakeit.Sentence(10)
		dto := repository.UpdateNoteDTO{Title: &title, Description: &desc}
		expectedNote := &repository.Note{ID: id, Title: title, Description: desc}

		t.Run("should update a note successfully", func(t *testing.T) {
			mockRepo.EXPECT().UpdateNote(gomock.Any(), id, dto).Return(expectedNote, nil)

			note, err := service.UpdateNote(ctx, id, dto)
			assert.NoError(t, err)
			assert.Equal(t, expectedNote, note)
		})

		t.Run("should return ErrNoteTitleTaken for duplicate key", func(t *testing.T) {
			mockRepo.EXPECT().UpdateNote(gomock.Any(), id, dto).Return(nil, errors.New("duplicate key error"))

			note, err := service.UpdateNote(ctx, id, dto)
			assert.Error(t, err)
			assert.Equal(t, ErrNoteTitleTaken, err)
			assert.Nil(t, note)
		})

		t.Run("should return ErrInternal for unknown errors", func(t *testing.T) {
			mockRepo.EXPECT().UpdateNote(gomock.Any(), id, dto).Return(nil, assert.AnError)

			note, err := service.UpdateNote(ctx, id, dto)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
			assert.Nil(t, note)
		})
	})

	t.Run("DeleteNote", func(t *testing.T) {
		id := uuid.New()

		t.Run("should delete a note successfully", func(t *testing.T) {
			mockRepo.EXPECT().DeleteNote(gomock.Any(), id).Return(nil)

			err := service.DeleteNote(ctx, id)
			assert.NoError(t, err)
		})

		t.Run("should return ErrInternal for repository errors", func(t *testing.T) {
			mockRepo.EXPECT().DeleteNote(gomock.Any(), id).Return(assert.AnError)

			err := service.DeleteNote(ctx, id)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
		})
	})

	t.Run("FetchNotes", func(t *testing.T) {
		expectedNotes := []repository.Note{
			{ID: uuid.New(), Title: gofakeit.Sentence(3), Description: gofakeit.Sentence(10)},
			{ID: uuid.New(), Title: gofakeit.Sentence(3), Description: gofakeit.Sentence(10)},
		}

		t.Run("should fetch all notes successfully", func(t *testing.T) {
			mockRepo.EXPECT().FetchNotes(gomock.Any()).Return(expectedNotes, nil)

			notes, err := service.FetchNotes(ctx)
			assert.NoError(t, err)
			assert.Equal(t, expectedNotes, notes)
		})

		t.Run("should return ErrInternal for repository errors", func(t *testing.T) {
			mockRepo.EXPECT().FetchNotes(gomock.Any()).Return(nil, assert.AnError)

			notes, err := service.FetchNotes(ctx)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
			assert.Nil(t, notes)
		})
	})
}
