package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
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
			assert.NotNil(t, note.UpdatedAt)

			// Verify directly in DB
			var (
				dbTitle, dbDescription   string
				dbCreatedAt, dbUpdatedAt time.Time
			)
			err = conn.QueryRow(ctx, "SELECT title, description, created_at, updated_at FROM core.notes WHERE id=$1", note.ID.String()).
				Scan(&dbTitle, &dbDescription, &dbCreatedAt, &dbUpdatedAt)
			assert.NoError(t, err)
			assert.Equal(t, dto.Title, dbTitle)
			assert.Equal(t, dto.Description, dbDescription)
			assert.True(t, note.CreatedAt.Equal(dbCreatedAt))
			assert.True(t, note.UpdatedAt.Equal(dbUpdatedAt))
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

	t.Run("FetchNoteByID", func(t *testing.T) {
		t.Run("should fetch an existing note given its ID", func(t *testing.T) {
			t.Parallel()

			dto := repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			}

			note, err := repo.CreateNote(ctx, dto)
			assert.NoError(t, err)

			fetchedNote, err := repo.FetchNoteByID(ctx, note.ID)
			assert.NoError(t, err)
			assert.Equal(t, note.ID, fetchedNote.ID)
			assert.Equal(t, dto.Title, fetchedNote.Title)
			assert.Equal(t, dto.Description, fetchedNote.Description)
			assert.True(t, note.CreatedAt.Equal(fetchedNote.CreatedAt))
			assert.NotNil(t, fetchedNote.UpdatedAt)
			assert.True(t, note.UpdatedAt.Equal(*fetchedNote.UpdatedAt))
		})

		t.Run("should return an error when fetching a non-existing note", func(t *testing.T) {
			t.Parallel()

			_, err := repo.FetchNoteByID(ctx, uuid.New())
			assert.Error(t, err)
			assert.Equal(t, sql.ErrNoRows, err)
		})
	})

	t.Run("UpdateNote", func(t *testing.T) {
		t.Run("should update a note's title, description and updated_at", func(t *testing.T) {
			t.Parallel()

			// Create the original note
			note, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Update the note's title and description
			newTitle := gofakeit.Sentence(3)
			newDescription := gofakeit.Sentence(10)
			updatedNote, err := repo.UpdateNote(ctx, note.ID, repository.UpdateNoteDTO{
				Title:       &newTitle,
				Description: &newDescription,
			})
			assert.NoError(t, err)
			assert.Equal(t, newTitle, updatedNote.Title)
			assert.Equal(t, newDescription, updatedNote.Description)
			assert.NotNil(t, updatedNote.UpdatedAt)
			assert.True(t, updatedNote.UpdatedAt.After(*note.UpdatedAt))

			// DB validation
			var (
				dbTitle, dbDescription string
				dbUpdatedAt            time.Time
			)
			err = conn.QueryRow(
				ctx, "SELECT title, description, updated_at FROM core.notes WHERE id=$1", note.ID.String(),
			).Scan(&dbTitle, &dbDescription, &dbUpdatedAt)
			assert.NoError(t, err)
			assert.Equal(t, newTitle, dbTitle)
			assert.Equal(t, newDescription, dbDescription)
			assert.True(t, updatedNote.UpdatedAt.Equal(dbUpdatedAt))
		})

		t.Run("should fail if updating title to an existing title", func(t *testing.T) {
			t.Parallel()

			// Create note A
			noteA, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Create note B
			noteB, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Attempt to rename B to A’s title
			_, err = repo.UpdateNote(ctx, noteB.ID, repository.UpdateNoteDTO{
				Title: &noteA.Title,
			})
			assert.Error(t, err)
		})
	})

	t.Run("DeleteNote", func(t *testing.T) {
		t.Run("should delete a note given the note ID", func(t *testing.T) {
			t.Parallel()

			// Create a note
			note, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Delete the note
			err = repo.DeleteNote(ctx, note.ID)
			assert.NoError(t, err)

			// FetchByID should return no rows found
			_, err = repo.FetchNoteByID(ctx, note.ID)
			assert.Error(t, err)

			// Validate DB
			var count int
			err = conn.QueryRow(ctx,
				"SELECT COUNT(*) FROM core.notes WHERE id=$1", note.ID.String(),
			).Scan(&count)
			assert.NoError(t, err)
			assert.Zero(t, count)
		})
	})

	t.Run("FetchNotes", func(t *testing.T) {
		t.Run("should fetch all notes", func(t *testing.T) {
			// Setup a separate postgres instance
			_, conn, cleanupFunc, err := tests.SetupPostgresDB(ctx)
			assert.NoError(t, err)

			defer func() {
				err := cleanupFunc()
				assert.NoError(t, err)
			}()

			repo := repository.NewRepository(conn)

			// Insert 3 random notes
			for range 3 {
				repo.CreateNote(ctx, repository.CreateNoteDTO{
					Title:       gofakeit.Sentence(3),
					Description: gofakeit.Sentence(10),
				})
			}

			notes, err := repo.FetchNotes(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(notes))
		})
	})
}
