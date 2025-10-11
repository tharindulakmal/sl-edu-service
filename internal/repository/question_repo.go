package repository

import (
	"database/sql"
	"fmt"
	"strconv"
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
	joinClause, whereClause, args := buildQuestionQueryParts(filters)

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
                SELECT q.id, q.grade_id, q.lesson_id, q.topic_id, q.subtopic_id, q.tutor_id, q.tute_id,
                       q.question, q.question_img_url, q.correct_answer, q.theory, q.solution,
                       q.other_answers,
                       DATE_FORMAT(q.created_at, '%%Y-%%m-%%dT%%H:%%i:%%sZ') as created_at
                FROM questions q
                %s
                %s
                ORDER BY q.id DESC
                LIMIT ? OFFSET ?`, joinClause, whereClause)

	args = append(args, pageSize, offset)

	if gradeID, ok := getFilterInt64(filters, "gradeId"); ok {
		if _, hasSubject := getFilterInt64(filters, "subjectId"); !hasSubject {
			if _, hasLesson := getFilterInt64(filters, "lessonId"); !hasLesson {
				fmt.Printf("[questions:list] grade-only SQL (gradeId=%d): %s | args=%v\n", gradeID, query, args)
			}
		}
	}

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
	joinClause, whereClause, args := buildQuestionQueryParts(filters)

	query := fmt.Sprintf("SELECT COUNT(*) FROM questions q %s %s", joinClause, whereClause)

	if gradeID, ok := getFilterInt64(filters, "gradeId"); ok {
		if _, hasSubject := getFilterInt64(filters, "subjectId"); !hasSubject {
			if _, hasLesson := getFilterInt64(filters, "lessonId"); !hasLesson {
				fmt.Printf("[questions:count] grade-only SQL (gradeId=%d): %s | args=%v\n", gradeID, query, args)
			}
		}
	}

	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func buildQuestionQueryParts(filters map[string]interface{}) (string, string, []interface{}) {
	joinParts := make([]string, 0)
	whereParts := []string{"1=1"}
	args := make([]interface{}, 0)

	gradeID, hasGrade := getFilterInt64(filters, "gradeId")
	subjectID, hasSubject := getFilterInt64(filters, "subjectId")
	lessonID, hasLesson := getFilterInt64(filters, "lessonId")

	gradeOnly := hasGrade && !hasSubject && !hasLesson

	if hasLesson {
		whereParts = append(whereParts, "q.lesson_id = ?")
		args = append(args, lessonID)
	} else {
		if (hasGrade && hasSubject) || (!gradeOnly && hasSubject) {
			joinParts = append(joinParts, "LEFT JOIN lessons l ON q.lesson_id IS NOT NULL AND q.lesson_id <> 0 AND q.lesson_id = l.id")
		}
		switch {
		case hasGrade && hasSubject:
			whereParts = append(whereParts,
				"(((q.lesson_id IS NULL OR q.lesson_id = 0) AND q.grade_id = ?) OR (q.lesson_id IS NOT NULL AND q.lesson_id <> 0 AND EXISTS (SELECT 1 FROM grade_subject_lessons gsl WHERE gsl.grade_id = ? AND gsl.subject_id = ? AND gsl.lesson_id = q.lesson_id)))",
			)
			args = append(args, gradeID, gradeID, subjectID)
		case hasGrade:
			whereParts = append(whereParts,
				"(((q.lesson_id IS NULL OR q.lesson_id = 0) AND q.grade_id = ?) OR (q.lesson_id IS NOT NULL AND q.lesson_id <> 0 AND EXISTS (SELECT 1 FROM grade_subject_lessons gsl WHERE gsl.grade_id = ? AND gsl.lesson_id = q.lesson_id))))",
			)
			args = append(args, gradeID, gradeID)
		case hasSubject:
			whereParts = append(whereParts, "(q.lesson_id IS NOT NULL AND q.lesson_id <> 0 AND l.subject_id = ?)")
			args = append(args, subjectID)
		}
	}

	if topicID, ok := getFilterInt64(filters, "topicId"); ok {
		whereParts = append(whereParts, "q.topic_id = ?")
		args = append(args, topicID)
	}
	if subtopicID, ok := getFilterInt64(filters, "subtopicId"); ok {
		whereParts = append(whereParts, "q.subtopic_id = ?")
		args = append(args, subtopicID)
	}
	if tutorID, ok := getFilterInt64(filters, "tutorId"); ok {
		whereParts = append(whereParts, "q.tutor_id = ?")
		args = append(args, tutorID)
	}
	if tuteID, ok := getFilterInt64(filters, "tuteId"); ok {
		whereParts = append(whereParts, "q.tute_id = ?")
		args = append(args, tuteID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	joinClause := ""
	if len(joinParts) > 0 {
		joinClause = strings.Join(joinParts, "\n")
	}

	return joinClause, whereClause, args
}

func getFilterInt64(filters map[string]interface{}, key string) (int64, bool) {
	val, ok := filters[key]
	if !ok || val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		if v == 0 {
			return 0, false
		}
		return int64(v), true
	case int64:
		if v == 0 {
			return 0, false
		}
		return v, true
	case int32:
		if v == 0 {
			return 0, false
		}
		return int64(v), true
	case float64:
		if v == 0 {
			return 0, false
		}
		return int64(v), true
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" || trimmed == "0" {
			return 0, false
		}
		parsed, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil || parsed == 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
