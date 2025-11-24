package server

import (
	"Project/internal/config"
	"Project/internal/handlers"
	"Project/internal/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Router struct {
	mux *http.ServeMux
}

func New(repo repository.Repository, config config.Config) *Router {
	mux := http.NewServeMux()

	qh := handlers.NewQuestionsHandler(repo)
	ah := handlers.NewAnswersHandler(repo)

	mux.HandleFunc("/questions/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/questions")
		path = strings.Trim(path, "/")

		if strings.HasSuffix(r.URL.Path, "/answers/") {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			parts := strings.Split(path, "/")
			if len(parts) < 2 {
				http.Error(w, "invalid path", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(parts[0])
			if err != nil {
				http.Error(w, "invalid question id", http.StatusBadRequest)
				return
			}

			ah.CreateAnswer(w, r, uint(id))
			return
		}

		switch r.Method {
		case http.MethodGet:
			qh.HandleQuestions(w, r)

		case http.MethodPost:
			if path != "" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			qh.CreateQuestion(w, r)

		case http.MethodDelete:
			qh.HandleQuestions(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/answers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ah.GetAnswer(w, r)
		case http.MethodDelete:
			ah.DeleteAnswer(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return &Router{mux: mux}
}

func (r *Router) RunRouter(config config.Config) {
	log.Printf("Start listening on %s", config.Port)
	err := http.ListenAndServe(config.Port, r.mux)
	if err != nil {
		log.Fatal(err)
	}
}
