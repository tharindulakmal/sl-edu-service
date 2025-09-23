package repository

import (
	"database/sql"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type SubjectRepository struct {
	DB *sql.DB
}

func NewSubjectRepository(db *sql.DB) *SubjectRepository {
	return &SubjectRepository{DB: db}
}

func (r *SubjectRepository) GetSubjectsByGradeID(gradeID int) ([]models.Subject, error) {
	rows, err := r.DB.Query("SELECT id, name FROM subjects WHERE grade_id = ?", gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	return subjects, nil
}
