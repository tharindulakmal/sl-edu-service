package repository

import (
	"database/sql"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type SmartNoteRepositoryInterface interface {
	GetSmartNote(lessonID, topicID, subID int) (models.SmartNote, error)
}

type SmartNoteRepository struct {
	DB *sql.DB
}

func NewSmartNoteRepository(db *sql.DB) SmartNoteRepositoryInterface {
	return &SmartNoteRepository{DB: db}
}

func (r *SmartNoteRepository) GetSmartNote(lessonID, topicID, subID int) (models.SmartNote, error) {
	var sn models.SmartNote
	err := r.DB.QueryRow(`
        SELECT sub_topic_name, image_def_url, definition, theory, image_theory_url, example, image_example_url
        FROM smart_notes
        WHERE lesson_id = ? AND topic_id = ? AND subtopic_id = ? LIMIT 1
    `, lessonID, topicID, subID).Scan(
		&sn.SubTopicName, &sn.ImageDefUrl, &sn.Definition,
		&sn.Theory, &sn.ImageTheoryUrl, &sn.Example, &sn.ImageExampleUrl,
	)
	return sn, err
}
