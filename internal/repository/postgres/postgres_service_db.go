package postgres

import (
	"Project/internal/models"
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListQuestions() ([]models.Question, error) {
	var qs []models.Question

	if err := r.db.Order("created_at DESC").Find(&qs).Error; err != nil {
		return nil, err
	}

	return qs, nil
}

func (r *PostgresRepository) CreateQuestion(text string) (*models.Question, error) {
	q := &models.Question{Text: text}

	if err := r.db.Create(q).Error; err != nil {
		return nil, err
	}

	return q, nil
}

func (r *PostgresRepository) GetQuestionWithAnswers(id uint) (*models.Question, error) {
	var q models.Question

	if err := r.db.Preload("Answers").First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &q, nil
}

func (r *PostgresRepository) DeleteQuestion(id uint) error {
	return r.db.Delete(&models.Question{}, id).Error
}

func (r *PostgresRepository) CreateAnswer(questionID uint, userID, text string) (*models.Answer, error) {
	a := &models.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (r *PostgresRepository) GetAnswer(id uint) (*models.Answer, error) {
	var a models.Answer

	if err := r.db.First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

func (r *PostgresRepository) DeleteAnswer(id uint) error {
	return r.db.Delete(&models.Answer{}, id).Error
}

func (r *PostgresRepository) QuestionExists(id uint) (bool, error) {
	var count int64

	if err := r.db.Model(&models.Question{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
