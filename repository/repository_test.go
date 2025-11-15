package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/the-code-genin/golang_integration_testing/repository"
	"github.com/the-code-genin/golang_integration_testing/tests"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()

	// Setup postgres database
	_, conn, cleanupFunc, err := tests.SetupPostgresDB(ctx)
	assert.NoError(t, err)

	defer func() {
		err := cleanupFunc()
		assert.NoError(t, err)
	}()

	repo := repository.NewRepository(conn)

	t.Run("CreateNote", func(t *testing.T) {
		t.Run("should create a note given a title and description", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			note, err := repo.CreateNote(ctx, dto)
			assert.NoError(t, err)
			assert.Equal(t, dto.Title, note.Title)
			assert.Equal(t, dto.Description, note.Description)
			assert.NotNil(t, note.ID)
			assert.False(t, note.CreatedAt.IsZero())

			// Verify directly in DB
			var (
				dbTitle, dbDescription string
				dbCreatedAt            time.Time
			)
			err = conn.QueryRow(ctx, "SELECT title, description, created_at FROM core.notes WHERE id=$1", note.ID.String()).
				Scan(&dbTitle, &dbDescription, &dbCreatedAt)
			assert.NoError(t, err)
			assert.Equal(t, dto.Title, dbTitle)
			assert.Equal(t, dto.Description, dbDescription)
			assert.True(t, note.CreatedAt.Equal(dbCreatedAt))
		})

		t.Run("should fail if an existing note's title is specified", func(t *testing.T) {
			t.Parallel()

			title := gofakeit.Sentence(3)

			_, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       title,
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Attempt to create duplicate
			_, err = repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       title,
				Description: gofakeit.Sentence(10),
			})
			assert.Error(t, err)

			// Confirm only one exists in DB
			var count int
			err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM core.notes WHERE title=$1", title).
				Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
		})
	})
}
