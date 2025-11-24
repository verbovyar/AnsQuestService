package repository

import "Project/internal/models"

type Repository interface {
	ListQuestions() ([]models.Question, error)
	CreateQuestion(text string) (*models.Question, error)
	GetQuestionWithAnswers(id uint) (*models.Question, error)
	DeleteQuestion(id uint) error
	CreateAnswer(questionID uint, userID, text string) (*models.Answer, error)
	GetAnswer(id uint) (*models.Answer, error)
	DeleteAnswer(id uint) error
	QuestionExists(id uint) (bool, error)
}
