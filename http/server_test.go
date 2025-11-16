package http_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	h "github.com/the-code-genin/golang_integration_testing/http"
	"github.com/the-code-genin/golang_integration_testing/repository"
	"github.com/the-code-genin/golang_integration_testing/service"
	"github.com/the-code-genin/golang_integration_testing/tests"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	// Setup postgres database
	_, conn, cleanupFunc, err := tests.SetupPostgresDB(ctx)
	assert.NoError(t, err)

	defer func() {
		err := cleanupFunc()
		assert.NoError(t, err)
	}()

	// Setup repository and service layers
	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)

	// Setup the server
	server := httptest.NewServer(h.NewServer(svc).Handler())
	defer server.Close()

	// Setup HTTP client
	httpClient := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("Endpoint to create notes", func(t *testing.T) {
		t.Run("should return a created status code given a valid title and description", func(t *testing.T) {
			t.Parallel()

			title := gofakeit.Sentence(3)
			description := gofakeit.Sentence(10)

			resp := httpClient.POST("/v1/notes").
				WithHeader("Content-Type", "application/json").
				WithJSON(map[string]string{
					"title":       title,
					"description": description,
				}).Expect().
				Status(http.StatusCreated).
				JSON().Object().
				ContainsSubset(map[string]any{
					"title":       title,
					"description": description,
				}).ContainsKey("id").
				ContainsKey("created_at").
				ContainsKey("updated_at")

			noteID, err := uuid.Parse(resp.Value("id").String().Raw())
			assert.NoError(t, err)

			// Confirm that the note is created in the DB
			note, err := repo.FetchNoteByID(ctx, noteID)
			assert.NoError(t, err)
			assert.Equal(t, title, note.Title)
			assert.Equal(t, description, note.Description)

			resp.Value("created_at").String().AsDateTime().IsEqual(note.CreatedAt)
		})

		t.Run("should return a bad request if the input parameters are invalid", func(t *testing.T) {
			t.Parallel()

			testcases := []struct {
				name    string
				payload map[string]string
			}{
				{
					name: "missing title",
					payload: map[string]string{
						"description": gofakeit.Sentence(10),
					},
				},
				{
					name: "missing description",
					payload: map[string]string{
						"title": gofakeit.Sentence(3),
					},
				},
				{
					name:    "both missing",
					payload: map[string]string{},
				},
			}

			for _, tc := range testcases {
				t.Run(tc.name, func(t *testing.T) {
					httpClient.POST("/v1/notes").
						WithHeader("Content-Type", "application/json").
						WithJSON(tc.payload).
						Expect().
						Status(http.StatusBadRequest).
						JSON().Object().
						ContainsSubset(map[string]any{
							"message": "bad request",
						})
				})
			}
		})

		t.Run("should return an internal error if an existing note's title is specified", func(t *testing.T) {
			t.Parallel()

			title := gofakeit.Sentence(3)

			_, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       title,
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Attempt to create duplicate
			httpClient.POST("/v1/notes").
				WithHeader("Content-Type", "application/json").
				WithJSON(map[string]string{
					"title":       title,
					"description": gofakeit.Sentence(10),
				}).Expect().
				Status(http.StatusInternalServerError).
				JSON().Object().
				ContainsSubset(map[string]any{
					"message": service.ErrInternal.Error(),
				})
		})
	})

	t.Run("FetchNoteByID", func(t *testing.T) {
		t.Run("should return an ok status code when given a valid ID", func(t *testing.T) {
			t.Parallel()

			// First, create a note in the DB
			note, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			// Fetch the note via HTTP
			resp := httpClient.GET("/v1/notes/{id}", note.ID.String()).
				Expect().
				Status(http.StatusOK).
				JSON().Object().
				ContainsSubset(map[string]any{
					"id":          note.ID.String(),
					"title":       note.Title,
					"description": note.Description,
				})

			resp.Value("created_at").String().AsDateTime().IsEqual(note.CreatedAt)
			resp.Value("updated_at").String().AsDateTime().IsEqual(*note.UpdatedAt)
		})

		t.Run("should return a not found error if a non-existent ID is provided", func(t *testing.T) {
			t.Parallel()

			// Generate a random UUID that doesn't exist
			nonExistentID := uuid.New()

			// Attempt to fetch
			httpClient.GET("/v1/notes/{id}", nonExistentID.String()).
				Expect().
				Status(http.StatusNotFound).
				JSON().Object().
				ContainsSubset(map[string]any{
					"message": service.ErrNoteNotFound.Error(),
				})
		})
	})

	t.Run("UpdateNote", func(t *testing.T) {
		t.Run("should return a 200 ok response given a valid ID and update payload", func(t *testing.T) {
			t.Parallel()

			// Create a note first
			note, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			newTitle := gofakeit.Sentence(3)
			newDescription := gofakeit.Sentence(10)

			resp := httpClient.PATCH("/v1/notes/{id}", note.ID.String()).
				WithHeader("Content-Type", "application/json").
				WithJSON(map[string]string{
					"title":       newTitle,
					"description": newDescription,
				}).Expect().
				Status(http.StatusOK).
				JSON().Object().
				ContainsSubset(map[string]any{
					"id":          note.ID.String(),
					"title":       newTitle,
					"description": newDescription,
				})

			updatedNote, err := repo.FetchNoteByID(ctx, note.ID)
			assert.NoError(t, err)
			assert.Equal(t, newTitle, updatedNote.Title)
			assert.Equal(t, newDescription, updatedNote.Description)

			resp.Value("updated_at").String().AsDateTime().IsEqual(*updatedNote.UpdatedAt)
		})
	})

	t.Run("DeleteNote", func(t *testing.T) {
		t.Run("should delete a note given a valid ID", func(t *testing.T) {
			t.Parallel()

			note, err := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			assert.NoError(t, err)

			httpClient.DELETE("/v1/notes/{id}", note.ID.String()).
				Expect().
				Status(http.StatusNoContent)

			_, err = repo.FetchNoteByID(ctx, note.ID)
			assert.Error(t, err)
			assert.Equal(t, sql.ErrNoRows, err)
		})
	})

	t.Run("FetchNotes", func(t *testing.T) {
		t.Run("should fetch all notes", func(t *testing.T) {
			t.Parallel()

			// Create multiple notes
			note1, _ := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})
			note2, _ := repo.CreateNote(ctx, repository.CreateNoteDTO{
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Sentence(10),
			})

			resp := httpClient.GET("/v1/notes").
				Expect().
				Status(http.StatusOK).
				JSON().Array()

			resp.Find(func(_ int, value *httpexpect.Value) bool {
				return value.Object().Value("id").String().Raw() == note1.ID.String()
			}).Object().ContainsSubset(map[string]any{
				"title":       note1.Title,
				"description": note1.Description,
			})

			resp.Find(func(_ int, value *httpexpect.Value) bool {
				return value.Object().Value("id").String().Raw() == note2.ID.String()
			}).Object().ContainsSubset(map[string]any{
				"title":       note2.Title,
				"description": note2.Description,
			})
		})
	})
}
