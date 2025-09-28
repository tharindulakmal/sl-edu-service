package repository

import (
	"database/sql"
	"fmt"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
)

type QuestionRepository interface {
	GetByID(id int) (*models.Question, error)
	GetList(filters map[string]interface{}, page, pageSize int) ([]models.Question, error)
	Create(q *models.Question) (int64, error)
	Update(q *models.Question) error
	Delete(id int) error
	Count(filters map[string]interface{}) (int, error)
}

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) GetByID(id int) (*models.Question, error) {
	query := `SELECT id, lesson_id, topic_id, subtopic_id, tutor_id, tute_id,
					 question, question_img_url, correct_answer, theory, solution,
					 other_answers, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') as created_at
			  FROM questions WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var q models.Question
	if err := row.Scan(&q.ID, &q.LessonID, &q.TopicID, &q.SubtopicID, &q.TutorID, &q.TuteID,
		&q.Question, &q.QuestionImg, &q.CorrectAnswer, &q.Theory, &q.Solution,
		&q.OtherAnswers, &q.CreatedAt); err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *questionRepository) GetList(filters map[string]interface{}, page, pageSize int) ([]models.Question, error) {
	where := "1=1"
	args := []interface{}{}

	if lessonId, ok := filters["lessonId"]; ok {
		where += " AND lesson_id = ?"
		args = append(args, lessonId)
	}
	if topicId, ok := filters["topicId"]; ok {
		where += " AND topic_id = ?"
		args = append(args, topicId)
	}
	if subtopicId, ok := filters["subtopicId"]; ok {
		where += " AND subtopic_id = ?"
		args = append(args, subtopicId)
	}
	if tutorId, ok := filters["tutorId"]; ok {
		where += " AND tutor_id = ?"
		args = append(args, tutorId)
	}
	if tuteId, ok := filters["tuteId"]; ok {
		where += " AND tute_id = ?"
		args = append(args, tuteId)
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, lesson_id, topic_id, subtopic_id, tutor_id, tute_id,
		       question, question_img_url, correct_answer, theory, solution,
		       other_answers,
		       DATE_FORMAT(created_at, '%%Y-%%m-%%dT%%H:%%i:%%sZ') as created_at
		FROM questions
		WHERE %s
		ORDER BY id ASC
		LIMIT ? OFFSET ?`, where)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.LessonID, &q.TopicID, &q.SubtopicID, &q.TutorID, &q.TuteID,
			&q.Question, &q.QuestionImg, &q.CorrectAnswer, &q.Theory, &q.Solution,
			&q.OtherAnswers, &q.CreatedAt); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *questionRepository) Create(q *models.Question) (int64, error) {
	query := `
		INSERT INTO questions (
			lesson_id, topic_id, subtopic_id, tutor_id, tute_id,
			question, question_img_url, correct_answer, theory, solution, other_answers
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query,
		q.LessonID, q.TopicID, q.SubtopicID, q.TutorID, q.TuteID,
		q.Question, q.QuestionImg, q.CorrectAnswer, q.Theory, q.Solution, q.OtherAnswers,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *questionRepository) Update(q *models.Question) error {
	query := `
		UPDATE questions
		SET lesson_id=?, topic_id=?, subtopic_id=?, tutor_id=?, tute_id=?,
		    question=?, question_img_url=?, correct_answer=?, theory=?, solution=?, other_answers=?
		WHERE id=?`
	_, err := r.db.Exec(query,
		q.LessonID, q.TopicID, q.SubtopicID, q.TutorID, q.TuteID,
		q.Question, q.QuestionImg, q.CorrectAnswer, q.Theory, q.Solution, q.OtherAnswers,
		q.ID,
	)
	return err
}

func (r *questionRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM questions WHERE id = ?", id)
	return err
}

func (r *questionRepository) Count(filters map[string]interface{}) (int, error) {
	where := "1=1"
	args := []interface{}{}

	if lessonId, ok := filters["lessonId"]; ok {
		where += " AND lesson_id = ?"
		args = append(args, lessonId)
	}
	if topicId, ok := filters["topicId"]; ok {
		where += " AND topic_id = ?"
		args = append(args, topicId)
	}
	if subtopicId, ok := filters["subtopicId"]; ok {
		where += " AND subtopic_id = ?"
		args = append(args, subtopicId)
	}
	if tutorId, ok := filters["tutorId"]; ok {
		where += " AND tutor_id = ?"
		args = append(args, tutorId)
	}
	if tuteId, ok := filters["tuteId"]; ok {
		where += " AND tute_id = ?"
		args = append(args, tuteId)
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM questions WHERE %s", where)

	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
