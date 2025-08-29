package repository

import (
	"database/sql"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type GradeRepository struct {
	DB *sql.DB
}

func NewGradeRepository(db *sql.DB) *GradeRepository {
	return &GradeRepository{DB: db}
}

func (r *GradeRepository) GetAllGrades() ([]models.Grade, error) {
	rows, err := r.DB.Query("SELECT id, grade FROM grades")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []models.Grade
	for rows.Next() {
		var g models.Grade
		if err := rows.Scan(&g.ID, &g.Grade); err != nil {
			return nil, err
		}
		grades = append(grades, g)
	}

	return grades, nil
}
