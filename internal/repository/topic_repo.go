package repository

import (
	"database/sql"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type TopicRepositoryInterface interface {
	GetTopicsByLesson(lessonID int) ([]models.Topic, error)
	GetDefaultSmartNote(lessonID int) (models.SmartNote, error)
}

type TopicRepository struct {
	DB *sql.DB
}

func NewTopicRepository(db *sql.DB) TopicRepositoryInterface {
	return &TopicRepository{DB: db}
}

func (r *TopicRepository) GetTopicsByLesson(lessonID int) ([]models.Topic, error) {
	// fetch topics -- schema now stores topic label in `name`
	rows, err := r.DB.Query("SELECT id, name FROM topics WHERE lesson_id = ? ORDER BY created_at", lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.TopicID, &t.TopicName); err != nil {
			return nil, err
		}

		// fetch subtopics for this topic
		subRows, err := r.DB.Query("SELECT id, topic_id, name FROM subtopics WHERE topic_id = ? ORDER BY created_at", t.TopicID)
		if err != nil {
			return nil, err
		}

		var subs []models.SubTopic
		for subRows.Next() {
			var s models.SubTopic
			if err := subRows.Scan(&s.SubTopicID, &s.TopicID, &s.SubTopicName); err != nil {
				return nil, err
			}
			subs = append(subs, s)
		}
		subRows.Close()

		t.SubTopicList = subs
		topics = append(topics, t)
	}

	return topics, nil
}

func (r *TopicRepository) GetDefaultSmartNote(lessonID int) (models.SmartNote, error) {
	var sn models.SmartNote
	err := r.DB.QueryRow(`
        SELECT sub_topic_name, image_def_url, definition, theory, image_theory_url, example, image_example_url
        FROM smart_notes
        WHERE lesson_id = ? AND is_default = TRUE LIMIT 1
    `, lessonID).Scan(
		&sn.SubTopicName, &sn.ImageDefUrl, &sn.Definition, &sn.Theory,
		&sn.ImageTheoryUrl, &sn.Example, &sn.ImageExampleUrl,
	)
	return sn, err
}
