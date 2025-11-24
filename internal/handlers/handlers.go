package handlers

import (
	"Project/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type QuestionsHandler struct {
	repo repository.Repository
}

func NewQuestionsHandler(repo repository.Repository) *QuestionsHandler {
	return &QuestionsHandler{repo: repo}
}

func (h *QuestionsHandler) HandleQuestions(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/questions")
	path = strings.Trim(path, "/")

	if r.Method == http.MethodGet && path == "" {
		h.ListQuestions(w, r)
		return
	}

	if r.Method == http.MethodGet && path != "" {
		if strings.Contains(path, "/answers") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		h.GetQuestion(w, r, uint(id))
		return
	}

	if r.Method == http.MethodDelete && path != "" {
		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		h.DeleteQuestion(w, r, uint(id))
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (h *QuestionsHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	qs, err := h.repo.ListQuestions()
	if err != nil {
		http.Error(w, "failed to fetch questions", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, qs)
}

type createQuestionRequest struct {
	Text string `json:"text"`
}

func (h *QuestionsHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	q, err := h.repo.CreateQuestion(req.Text)
	if err != nil {
		http.Error(w, "failed to create question", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, q)
}

func (h *QuestionsHandler) GetQuestion(w http.ResponseWriter, r *http.Request, id uint) {
	q, err := h.repo.GetQuestionWithAnswers(id)
	if err != nil {
		http.Error(w, "failed to fetch question", http.StatusInternalServerError)
		return
	}
	if q == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, q)
}

func (h *QuestionsHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request, id uint) {
	if err := h.repo.DeleteQuestion(id); err != nil {
		http.Error(w, "failed to delete question", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type AnswersHandler struct {
	repo repository.Repository
}

func NewAnswersHandler(repo repository.Repository) *AnswersHandler {
	return &AnswersHandler{repo: repo}
}

var req struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

func (h *AnswersHandler) CreateAnswer(w http.ResponseWriter, r *http.Request, questionID uint) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Text) == "" {
		http.Error(w, "user_id and text are required", http.StatusBadRequest)
		return
	}

	exists, err := h.repo.QuestionExists(questionID)
	if err != nil {
		http.Error(w, "failed to check question", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "question not found", http.StatusBadRequest)
		return
	}

	a, err := h.repo.CreateAnswer(questionID, req.UserID, req.Text)
	if err != nil {
		http.Error(w, "failed to create answer", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

func (h *AnswersHandler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/answers")
	path = strings.Trim(path, "/")
	if path == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	a, err := h.repo.GetAnswer(uint(id))
	if err != nil {
		http.Error(w, "failed to fetch answer", http.StatusInternalServerError)
		return
	}

	if a == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, a)
}

func (h *AnswersHandler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/answers")
	path = strings.Trim(path, "/")
	if path == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteAnswer(uint(id)); err != nil {
		http.Error(w, "failed to delete answer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
