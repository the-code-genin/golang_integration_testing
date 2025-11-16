package http_test

import (
	"context"
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
		t.Run("should create a note given a title and description", func(t *testing.T) {
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

		t.Run("should fail if title or description is empty", func(t *testing.T) {
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
						ContainsKey("message")
				})
			}
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
}
