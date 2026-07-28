package reclaim

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient points a Client at a stub server, capturing the request body it
// sends so the encoding can be asserted on.
func newTestClient(t *testing.T, response string, captured *[]byte) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		*captured = body

		w.Header().Set("Content-Type", "application/json")
		_, err = io.WriteString(w, response)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := New()
	client.baseUrl = server.URL
	return client
}

// Titles used to be interpolated into a JSON string literal with fmt.Sprintf, so
// any quote, backslash or newline produced a malformed body and an opaque API
// error. An agent generating titles hits this routinely.
func TestCreateTask_encodesAwkwardTitles(t *testing.T) {
	titles := map[string]string{
		"double quotes": `Review the "urgent" dashboard`,
		"backslash":     `check C:\Users\ptanton\notes`,
		"newline":       "first line\nsecond line",
		"tab":           "before\tafter",
		"json fragment": `{"title": "nested", "priority": "P1"}`,
		"unicode":       "review café résumé — 日本語",
	}

	for name, title := range titles {
		t.Run(name, func(t *testing.T) {
			var captured []byte
			client := newTestClient(t, `{"id": 1}`, &captured)

			_, err := client.CreateTask(&CreateTaskRequest{
				Title:              title,
				MinChunkSize:       2,
				MaxChunkSize:       16,
				TimeChunksRequired: 2,
				Priority:           P3,
			})
			require.NoError(t, err)

			var decoded CreateTaskRequest
			require.NoError(t, json.Unmarshal(captured, &decoded), "request body was not valid JSON: %s", captured)
			assert.Equal(t, title, decoded.Title)
		})
	}
}

func TestCreateTask_defaultsStatusAndCategory(t *testing.T) {
	var captured []byte
	client := newTestClient(t, `{"id": 1}`, &captured)

	_, err := client.CreateTask(&CreateTaskRequest{Title: "a task", Priority: P3})
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(captured, &decoded))
	assert.Equal(t, "NEW", decoded["status"])
	assert.Equal(t, "WORK", decoded["eventCategory"])
	assert.NotContains(t, decoded, "due", "an unset due must be omitted, not sent as the zero time")
	assert.NotContains(t, decoded, "notes")
}

func TestCreateTask_sendsDueAndNotes(t *testing.T) {
	var captured []byte
	client := newTestClient(t, `{"id": 1}`, &captured)

	due := time.Date(2026, 8, 14, 17, 30, 0, 0, time.UTC)
	_, err := client.CreateTask(&CreateTaskRequest{Title: "a task", Priority: P3, Due: &due, Notes: "some notes"})
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(captured, &decoded))
	assert.Equal(t, "2026-08-14T17:30:00Z", decoded["due"])
	assert.Equal(t, "some notes", decoded["notes"])
}

func TestCreateTask_readsBackTheCreatedTask(t *testing.T) {
	var captured []byte
	client := newTestClient(t, `{"id": 13354821, "title": "a task", "due": "2026-08-14T18:30:00+01:00", "snoozeUntil": null}`, &captured)

	task, err := client.CreateTask(&CreateTaskRequest{Title: "a task", Priority: P3})
	require.NoError(t, err)

	assert.Equal(t, 13354821, task.Id)
	require.NotNil(t, task.Due)
	assert.Equal(t, time.Date(2026, 8, 14, 17, 30, 0, 0, time.UTC), task.Due.UTC())
	assert.Nil(t, task.SnoozeUntil)
}

func TestSnoozeTask_returnsTheUpdatedTask(t *testing.T) {
	var captured []byte
	client := newTestClient(t, `{"id": 42, "snoozeUntil": "2026-08-05T09:00:00Z"}`, &captured)

	snoozeUntil := time.Date(2026, 8, 5, 9, 0, 0, 0, time.UTC)
	task, err := client.SnoozeTask(42, snoozeUntil)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(captured, &decoded))
	assert.Equal(t, "2026-08-05T09:00:00Z", decoded["snoozeUntil"])

	assert.Equal(t, 42, task.Id)
	require.NotNil(t, task.SnoozeUntil)
	assert.Equal(t, snoozeUntil, task.SnoozeUntil.UTC())
}

func TestCreateTask_errorIncludesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"message":"minChunkSize must be positive"}`)
	}))
	t.Cleanup(server.Close)

	client := New()
	client.baseUrl = server.URL

	_, err := client.CreateTask(&CreateTaskRequest{Title: "a task", Priority: P3})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "minChunkSize must be positive")
}
