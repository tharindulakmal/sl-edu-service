package repository

import "github.com/tharindulakmal/sl-edu-service/internal/models"

type MockGradeRepository struct {
	Grades []models.Grade
	Err    error
}

func (m *MockGradeRepository) GetAllGrades() ([]models.Grade, error) {
	return m.Grades, m.Err
}
