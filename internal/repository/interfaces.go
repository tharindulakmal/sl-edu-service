package repository

import "github.com/tharindulakmal/sl-edu-service/internal/models"

type GradeRepositoryInterface interface {
	GetAllGrades() ([]models.Grade, error)
}
