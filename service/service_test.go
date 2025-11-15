package service

import (
	"context"
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

		t.Run("should return ErrInternal if an error occurred while creating the note", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			// Simulate repository returning an error for duplicate title
			mockRepo.EXPECT().
				CreateNote(ctx, dto).
				Return(nil, assert.AnError).
				Times(1)

			note, err := service.CreateNote(ctx, dto)
			assert.Error(t, err)
			assert.Equal(t, ErrInternal, err)
			assert.Nil(t, note)
		})
	})
}
