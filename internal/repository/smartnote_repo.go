package repository

import (
	"database/sql"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type SmartNoteRepositoryInterface interface {
	GetSmartNote(gradeID, subjectID, lessonID int, topicID, subID *int) (models.SmartNote, error)
}

type SmartNoteRepository struct {
	DB *sql.DB
}

func NewSmartNoteRepository(db *sql.DB) SmartNoteRepositoryInterface {
	return &SmartNoteRepository{DB: db}
}

func (r *SmartNoteRepository) GetSmartNote(gradeID, subjectID, lessonID int, topicID, subID *int) (models.SmartNote, error) {
	var sn models.SmartNote
	query := `
		SELECT sn.sub_topic_name, sn.image_def_url, sn.definition, sn.theory,
		       sn.image_theory_url, sn.example, sn.image_example_url
		FROM smart_notes sn
			INNER JOIN lessons l ON sn.lesson_id = l.id
			INNER JOIN subjects s ON l.subject_id = s.id
		WHERE sn.lesson_id = ?
		  AND l.subject_id = ?
		  AND s.grade_id = ?`
	args := []interface{}{lessonID, subjectID, gradeID}
	if topicID != nil {
		query += " AND sn.topic_id = ?"
		args = append(args, *topicID)
	}
	if subID != nil {
		query += " AND sn.subtopic_id = ?"
		args = append(args, *subID)
	}
	query += " ORDER BY sn.id LIMIT 1"

	err := r.DB.QueryRow(query, args...).Scan(
		&sn.SubTopicName,
		&sn.ImageDefUrl,
		&sn.Definition,
		&sn.Theory,
		&sn.ImageTheoryUrl,
		&sn.Example,
		&sn.ImageExampleUrl,
	)
	return sn, err
}
