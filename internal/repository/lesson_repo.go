package repository

import (
	"database/sql"
	"fmt"
)

type Lesson struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type LessonRepositoryInterface interface {
	GetLessonsBySubject(subjectId int) ([]Lesson, error)
}

type LessonRepository struct {
	DB *sql.DB
}

func NewLessonRepository(db *sql.DB) LessonRepositoryInterface {
	return &LessonRepository{DB: db}
}

func (r *LessonRepository) GetLessonsBySubject(subjectId int) ([]Lesson, error) {
	rows, err := r.DB.Query("SELECT id, name, image_url FROM lessons WHERE subject_id = ?", subjectId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch lessons: %w", err)
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.Name, &l.ImageUrl); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}
