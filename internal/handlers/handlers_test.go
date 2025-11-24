package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Project/internal/handlers"
	"Project/internal/models"
	"Project/internal/repository"
	"Project/internal/repository/postgres"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestRepo(t *testing.T) repository.Repository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&models.Question{}, &models.Answer{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return postgres.New(db)
}

func TestCreateQuestion(t *testing.T) {
	repo := newTestRepo(t)
	h := handlers.NewQuestionsHandler(repo)

	body := bytes.NewBufferString(`{"text":"test question"}`)
	req := httptest.NewRequest(http.MethodPost, "/questions/", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.CreateQuestion(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var got struct {
		ID   uint   `json:"id"`
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Text != "test question" {
		t.Fatalf("expected text %q, got %q", "test question", got.Text)
	}

	if got.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}
}
