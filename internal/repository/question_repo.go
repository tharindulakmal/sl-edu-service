package repository

import (
	"database/sql"
	"fmt"
	"strings"

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
	query := `SELECT id, grade_id, lesson_id, topic_id, subtopic_id, tutor_id, tute_id,
					 question, question_img_url, correct_answer, theory, solution,
					 other_answers, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') as created_at
			  FROM questions WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var q models.Question
	if err := row.Scan(&q.ID, &q.GradeID, &q.LessonID, &q.TopicID, &q.SubtopicID, &q.TutorID, &q.TuteID,
		&q.Question, &q.QuestionImg, &q.CorrectAnswer, &q.Theory, &q.Solution,
		&q.OtherAnswers, &q.CreatedAt); err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *questionRepository) GetList(filters map[string]interface{}, page, pageSize int) ([]models.Question, error) {
	whereParts := []string{"1=1"}
	args := []interface{}{}

	gradeID, hasGrade := filters["gradeId"]
	subjectID, hasSubject := filters["subjectId"]

	switch {
	case hasSubject && !hasGrade:
		whereParts = append(whereParts, "l.subject_id = ?")
		args = append(args, subjectID)
	case hasGrade && !hasSubject:
		whereParts = append(whereParts, "EXISTS (SELECT 1 FROM grade_subjects gs WHERE gs.grade_id = ? AND gs.subject_id = l.subject_id)")
		args = append(args, gradeID)
	case hasGrade && hasSubject:
		whereParts = append(whereParts, "l.subject_id = ?")
		args = append(args, subjectID)
		whereParts = append(whereParts, "EXISTS (SELECT 1 FROM grade_subjects gs WHERE gs.grade_id = ? AND gs.subject_id = l.subject_id)")
		args = append(args, gradeID)
	}

	if lessonID, ok := filters["lessonId"]; ok {
		whereParts = append(whereParts, "q.lesson_id = ?")
		args = append(args, lessonID)
	}
	if topicID, ok := filters["topicId"]; ok {
		whereParts = append(whereParts, "q.topic_id = ?")
		args = append(args, topicID)
	}
	if subtopicID, ok := filters["subtopicId"]; ok {
		whereParts = append(whereParts, "q.subtopic_id = ?")
		args = append(args, subtopicID)
	}
	if tutorID, ok := filters["tutorId"]; ok {
		whereParts = append(whereParts, "q.tutor_id = ?")
		args = append(args, tutorID)
	}
	if tuteID, ok := filters["tuteId"]; ok {
		whereParts = append(whereParts, "q.tute_id = ?")
		args = append(args, tuteID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
                SELECT q.id, q.grade_id, q.lesson_id, q.topic_id, q.subtopic_id, q.tutor_id, q.tute_id,
                       q.question, q.question_img_url, q.correct_answer, q.theory, q.solution,
                       q.other_answers,
                       DATE_FORMAT(q.created_at, '%%Y-%%m-%%dT%%H:%%i:%%sZ') as created_at
                FROM questions q
                INNER JOIN lessons l ON q.lesson_id = l.id
                %s
                ORDER BY q.id DESC
                LIMIT ? OFFSET ?`, whereClause)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.GradeID, &q.LessonID, &q.TopicID, &q.SubtopicID, &q.TutorID, &q.TuteID,
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
			lesson_id, grade_id, topic_id, subtopic_id, tutor_id, tute_id,
			question, question_img_url, correct_answer, theory, solution, other_answers
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query,
		q.LessonID, q.GradeID, q.TopicID, q.SubtopicID, q.TutorID, q.TuteID,
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
		SET lesson_id=?, grade_id=?, topic_id=?, subtopic_id=?, tutor_id=?, tute_id=?,
		    question=?, question_img_url=?, correct_answer=?, theory=?, solution=?, other_answers=?
		WHERE id=?`
	_, err := r.db.Exec(query,
		q.LessonID, q.GradeID, q.TopicID, q.SubtopicID, q.TutorID, q.TuteID,
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
	whereParts := []string{"1=1"}
	args := []interface{}{}

	gradeID, hasGrade := filters["gradeId"]
	subjectID, hasSubject := filters["subjectId"]

	switch {
	case hasSubject && !hasGrade:
		whereParts = append(whereParts, "l.subject_id = ?")
		args = append(args, subjectID)
	case hasGrade && !hasSubject:
		whereParts = append(whereParts, "EXISTS (SELECT 1 FROM grade_subjects gs WHERE gs.grade_id = ? AND gs.subject_id = l.subject_id)")
		args = append(args, gradeID)
	case hasGrade && hasSubject:
		whereParts = append(whereParts, "l.subject_id = ?")
		args = append(args, subjectID)
		whereParts = append(whereParts, "EXISTS (SELECT 1 FROM grade_subjects gs WHERE gs.grade_id = ? AND gs.subject_id = l.subject_id)")
		args = append(args, gradeID)
	}

	if lessonID, ok := filters["lessonId"]; ok {
		whereParts = append(whereParts, "q.lesson_id = ?")
		args = append(args, lessonID)
	}
	if topicID, ok := filters["topicId"]; ok {
		whereParts = append(whereParts, "q.topic_id = ?")
		args = append(args, topicID)
	}
	if subtopicID, ok := filters["subtopicId"]; ok {
		whereParts = append(whereParts, "q.subtopic_id = ?")
		args = append(args, subtopicID)
	}
	if tutorID, ok := filters["tutorId"]; ok {
		whereParts = append(whereParts, "q.tutor_id = ?")
		args = append(args, tutorID)
	}
	if tuteID, ok := filters["tuteId"]; ok {
		whereParts = append(whereParts, "q.tute_id = ?")
		args = append(args, tuteID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM questions q INNER JOIN lessons l ON q.lesson_id = l.id %s", whereClause)

	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
