package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	menuconfigmodels "github.com/tharindulakmal/sl-edu-service/internal/models/menuconfig"
)

type (
	Grade              = menuconfigmodels.Grade
	GradeUpsert        = menuconfigmodels.GradeUpsert
	Subject            = menuconfigmodels.Subject
	SubjectUpsert      = menuconfigmodels.SubjectUpsert
	CatalogResponse    = menuconfigmodels.CatalogResponse
	GradeSubject       = menuconfigmodels.GradeSubject
	GradeSubjectLesson = menuconfigmodels.GradeSubjectLesson
	Topic              = menuconfigmodels.Topic
	TopicUpsert        = menuconfigmodels.TopicUpsert
	Subtopic           = menuconfigmodels.Subtopic
	SubtopicUpsert     = menuconfigmodels.SubtopicUpsert
	Tutor              = menuconfigmodels.Tutor
	TutorUpsert        = menuconfigmodels.TutorUpsert
	Year               = menuconfigmodels.Year
	YearUpsert         = menuconfigmodels.YearUpsert
	Tutorial           = menuconfigmodels.Tutorial
	TutorialUpsert     = menuconfigmodels.TutorialUpsert
)

var (
	ErrMenuConfigNotFound              = errors.New("menuconfig: not found")
	ErrGradeSubjectAlreadyLinked       = errors.New("menuconfig: grade subject link exists")
	ErrGradeSubjectLessonAlreadyLinked = errors.New("menuconfig: grade subject lesson link exists")
)

type MenuConfigRepository struct {
	db *sql.DB
}

func NewMenuConfigRepository(db *sql.DB) *MenuConfigRepository {
	return &MenuConfigRepository{db: db}
}

func (r *MenuConfigRepository) ListGrades(ctx context.Context, search string, page, pageSize int) ([]Grade, int, error) {
	baseQuery := "SELECT id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM grades"
	countQuery := "SELECT COUNT(*) FROM grades"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)
	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	grades := make([]Grade, 0)
	for rows.Next() {
		var g Grade
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			return nil, 0, err
		}
		grades = append(grades, g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

func (r *MenuConfigRepository) GetGrade(ctx context.Context, id int64) (*Grade, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM grades WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var g Grade
	if err := stmt.QueryRowContext(ctx, id).Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &g, nil
}

func (r *MenuConfigRepository) CreateGrade(ctx context.Context, input GradeUpsert) (*Grade, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO grades (name) VALUES (?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetGrade(ctx, id)
}

func (r *MenuConfigRepository) UpdateGrade(ctx context.Context, id int64, input GradeUpsert) (*Grade, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE grades SET name = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name), id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetGrade(ctx, id)
}

func (r *MenuConfigRepository) DeleteGrade(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM grades WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListSubjects(ctx context.Context, gradeID *int64, search string, page, pageSize int) ([]Subject, int, error) {
	baseQuery := "SELECT id, grade_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subjects"
	countQuery := "SELECT COUNT(*) FROM subjects"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if gradeID != nil {
		filters = append(filters, "grade_id = ?")
		args = append(args, *gradeID)
		countArgs = append(countArgs, *gradeID)
	}
	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	subjects := make([]Subject, 0)
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.GradeID, &s.Name, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		subjects = append(subjects, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *MenuConfigRepository) GetSubject(ctx context.Context, id int64) (*Subject, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, grade_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subjects WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var s Subject
	if err := stmt.QueryRowContext(ctx, id).Scan(&s.ID, &s.GradeID, &s.Name, &s.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *MenuConfigRepository) CreateSubject(ctx context.Context, input SubjectUpsert) (*Subject, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO subjects (grade_id, name) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.GradeID, strings.TrimSpace(input.Name))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetSubject(ctx, id)
}

func (r *MenuConfigRepository) UpdateSubject(ctx context.Context, id int64, input SubjectUpsert) (*Subject, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE subjects SET grade_id = ?, name = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.GradeID, strings.TrimSpace(input.Name), id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetSubject(ctx, id)
}

func (r *MenuConfigRepository) DeleteSubject(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM subjects WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) LinkGradeSubject(ctx context.Context, gradeID, subjectID int64) (*GradeSubject, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := ensureGradeExists(ctx, tx, gradeID); err != nil {
		return nil, err
	}
	if err := ensureSubjectExists(ctx, tx, subjectID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO grade_subjects (grade_id, subject_id) VALUES (?, ?)", gradeID, subjectID); err != nil {
		if isDuplicateEntry(err) {
			return nil, ErrGradeSubjectAlreadyLinked
		}
		return nil, err
	}

	var gs GradeSubject
	row := tx.QueryRowContext(ctx, "SELECT grade_id, subject_id, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') FROM grade_subjects WHERE grade_id = ? AND subject_id = ?", gradeID, subjectID)
	if err := row.Scan(&gs.GradeID, &gs.SubjectID, &gs.CreatedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &gs, nil
}

func (r *MenuConfigRepository) UnlinkGradeSubject(ctx context.Context, gradeID, subjectID int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM grade_subjects WHERE grade_id = ? AND subject_id = ?", gradeID, subjectID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) LinkGradeSubjectLesson(ctx context.Context, gradeID, subjectID, lessonID int64) (*GradeSubjectLesson, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := ensureGradeExists(ctx, tx, gradeID); err != nil {
		return nil, err
	}
	if err := ensureSubjectExists(ctx, tx, subjectID); err != nil {
		return nil, err
	}
	if err := ensureLessonBelongsToSubject(ctx, tx, lessonID, subjectID); err != nil {
		return nil, err
	}
	if err := ensureGradeSubjectLinkExists(ctx, tx, gradeID, subjectID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO grade_subject_lessons (grade_id, subject_id, lesson_id) VALUES (?, ?, ?)", gradeID, subjectID, lessonID); err != nil {
		if isDuplicateEntry(err) {
			return nil, ErrGradeSubjectLessonAlreadyLinked
		}
		return nil, err
	}

	var gsl GradeSubjectLesson
	row := tx.QueryRowContext(ctx, "SELECT grade_id, subject_id, lesson_id, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') FROM grade_subject_lessons WHERE grade_id = ? AND subject_id = ? AND lesson_id = ?", gradeID, subjectID, lessonID)
	if err := row.Scan(&gsl.GradeID, &gsl.SubjectID, &gsl.LessonID, &gsl.CreatedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &gsl, nil
}

func (r *MenuConfigRepository) UnlinkGradeSubjectLesson(ctx context.Context, gradeID, subjectID, lessonID int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM grade_subject_lessons WHERE grade_id = ? AND subject_id = ? AND lesson_id = ?", gradeID, subjectID, lessonID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListLessons(ctx context.Context, gradeID, subjectID *int64, search string, page, pageSize int) ([]menuconfigmodels.Lesson, int, error) {
	baseQuery := "SELECT l.id, l.subject_id, l.name, DATE_FORMAT(l.created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM lessons l"
	countQuery := "SELECT COUNT(*) FROM lessons l"
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if gradeID != nil {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM grade_subjects gs WHERE gs.subject_id = l.subject_id AND gs.grade_id = ?)")
		args = append(args, *gradeID)
		countArgs = append(countArgs, *gradeID)
	}
	if subjectID != nil {
		conditions = append(conditions, "l.subject_id = ?")
		args = append(args, *subjectID)
		countArgs = append(countArgs, *subjectID)
	}
	if search != "" {
		conditions = append(conditions, "l.name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}

	if len(conditions) > 0 {
		where := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += where
		countQuery += where
	}

	baseQuery += " ORDER BY l.id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	lessons := make([]menuconfigmodels.Lesson, 0)
	for rows.Next() {
		var l menuconfigmodels.Lesson
		if err := rows.Scan(&l.ID, &l.SubjectID, &l.Name, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		lessons = append(lessons, l)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *MenuConfigRepository) GetLesson(ctx context.Context, id int64) (*menuconfigmodels.Lesson, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, subject_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM lessons WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var l menuconfigmodels.Lesson
	if err := stmt.QueryRowContext(ctx, id).Scan(&l.ID, &l.SubjectID, &l.Name, &l.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *MenuConfigRepository) CreateLesson(ctx context.Context, input menuconfigmodels.LessonUpsert) (*menuconfigmodels.Lesson, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO lessons (subject_id, name) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.SubjectID, strings.TrimSpace(input.Name))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetLesson(ctx, id)
}

func (r *MenuConfigRepository) UpdateLesson(ctx context.Context, id int64, input menuconfigmodels.LessonUpsert) (*menuconfigmodels.Lesson, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE lessons SET subject_id = ?, name = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.SubjectID, strings.TrimSpace(input.Name), id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetLesson(ctx, id)
}

func (r *MenuConfigRepository) DeleteLesson(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM lessons WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListTopics(ctx context.Context, lessonID *int64, search string, page, pageSize int) ([]Topic, int, error) {
	baseQuery := "SELECT id, lesson_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM topics"
	countQuery := "SELECT COUNT(*) FROM topics"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if lessonID != nil {
		filters = append(filters, "lesson_id = ?")
		args = append(args, *lessonID)
		countArgs = append(countArgs, *lessonID)
	}
	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	topics := make([]Topic, 0)
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.LessonID, &t.Name, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		topics = append(topics, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return topics, total, nil
}

func (r *MenuConfigRepository) GetTopic(ctx context.Context, id int64) (*Topic, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, lesson_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM topics WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t Topic
	if err := stmt.QueryRowContext(ctx, id).Scan(&t.ID, &t.LessonID, &t.Name, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *MenuConfigRepository) CreateTopic(ctx context.Context, input TopicUpsert) (*Topic, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO topics (lesson_id, name) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.LessonID, strings.TrimSpace(input.Name))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetTopic(ctx, id)
}

func (r *MenuConfigRepository) UpdateTopic(ctx context.Context, id int64, input TopicUpsert) (*Topic, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE topics SET lesson_id = ?, name = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.LessonID, strings.TrimSpace(input.Name), id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetTopic(ctx, id)
}

func (r *MenuConfigRepository) DeleteTopic(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM topics WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListSubtopics(ctx context.Context, topicID *int64, search string, page, pageSize int) ([]Subtopic, int, error) {
	baseQuery := "SELECT id, topic_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subtopics"
	countQuery := "SELECT COUNT(*) FROM subtopics"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if topicID != nil {
		filters = append(filters, "topic_id = ?")
		args = append(args, *topicID)
		countArgs = append(countArgs, *topicID)
	}
	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	subtopics := make([]Subtopic, 0)
	for rows.Next() {
		var s Subtopic
		if err := rows.Scan(&s.ID, &s.TopicID, &s.Name, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		subtopics = append(subtopics, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return subtopics, total, nil
}

func (r *MenuConfigRepository) GetSubtopic(ctx context.Context, id int64) (*Subtopic, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, topic_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subtopics WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var s Subtopic
	if err := stmt.QueryRowContext(ctx, id).Scan(&s.ID, &s.TopicID, &s.Name, &s.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *MenuConfigRepository) CreateSubtopic(ctx context.Context, input SubtopicUpsert) (*Subtopic, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO subtopics (topic_id, name) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.TopicID, strings.TrimSpace(input.Name))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetSubtopic(ctx, id)
}

func (r *MenuConfigRepository) UpdateSubtopic(ctx context.Context, id int64, input SubtopicUpsert) (*Subtopic, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE subtopics SET topic_id = ?, name = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.TopicID, strings.TrimSpace(input.Name), id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetSubtopic(ctx, id)
}

func (r *MenuConfigRepository) DeleteSubtopic(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM subtopics WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListTutors(ctx context.Context, search string, page, pageSize int) ([]Tutor, int, error) {
	baseQuery := "SELECT id, name, email, phone, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM tutors"
	countQuery := "SELECT COUNT(*) FROM tutors"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tutors := make([]Tutor, 0)
	for rows.Next() {
		var t Tutor
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Phone, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		tutors = append(tutors, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return tutors, total, nil
}

func (r *MenuConfigRepository) GetTutor(ctx context.Context, id int64) (*Tutor, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, email, phone, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM tutors WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t Tutor
	if err := stmt.QueryRowContext(ctx, id).Scan(&t.ID, &t.Name, &t.Email, &t.Phone, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *MenuConfigRepository) CreateTutor(ctx context.Context, input TutorUpsert) (*Tutor, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO tutors (name, email, phone) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name), input.Email, input.Phone)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetTutor(ctx, id)
}

func (r *MenuConfigRepository) UpdateTutor(ctx context.Context, id int64, input TutorUpsert) (*Tutor, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE tutors SET name = ?, email = ?, phone = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name), input.Email, input.Phone, id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetTutor(ctx, id)
}

func (r *MenuConfigRepository) DeleteTutor(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM tutors WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListYears(ctx context.Context, search string, page, pageSize int) ([]Year, int, error) {
	baseQuery := "SELECT id, value, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM years"
	countQuery := "SELECT COUNT(*) FROM years"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if search != "" {
		filters = append(filters, "CAST(value AS CHAR) LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	years := make([]Year, 0)
	for rows.Next() {
		var y Year
		if err := rows.Scan(&y.ID, &y.Value, &y.CreatedAt); err != nil {
			return nil, 0, err
		}
		years = append(years, y)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return years, total, nil
}

func (r *MenuConfigRepository) GetYear(ctx context.Context, id int64) (*Year, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, value, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM years WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var y Year
	if err := stmt.QueryRowContext(ctx, id).Scan(&y.ID, &y.Value, &y.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &y, nil
}

func (r *MenuConfigRepository) CreateYear(ctx context.Context, input YearUpsert) (*Year, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO years (value) VALUES (?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.Value)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetYear(ctx, id)
}

func (r *MenuConfigRepository) UpdateYear(ctx context.Context, id int64, input YearUpsert) (*Year, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE years SET value = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, input.Value, id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetYear(ctx, id)
}

func (r *MenuConfigRepository) DeleteYear(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM years WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) ListTutorials(ctx context.Context, search string, page, pageSize int) ([]Tutorial, int, error) {
	baseQuery := "SELECT id, name, url, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM tutorials"
	countQuery := "SELECT COUNT(*) FROM tutorials"
	filters := make([]string, 0)
	args := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if search != "" {
		filters = append(filters, "name LIKE CONCAT('%', ?, '%')")
		args = append(args, search)
		countArgs = append(countArgs, search)
	}
	if len(filters) > 0 {
		where := " WHERE " + strings.Join(filters, " AND ")
		baseQuery += where
		countQuery += where
	}
	baseQuery += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offsetFromPage(page, pageSize))

	stmt, err := r.db.PrepareContext(ctx, baseQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tutorials := make([]Tutorial, 0)
	for rows.Next() {
		var t Tutorial
		if err := rows.Scan(&t.ID, &t.Name, &t.URL, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		tutorials = append(tutorials, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.count(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return tutorials, total, nil
}

func (r *MenuConfigRepository) GetTutorial(ctx context.Context, id int64) (*Tutorial, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, url, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM tutorials WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var t Tutorial
	if err := stmt.QueryRowContext(ctx, id).Scan(&t.ID, &t.Name, &t.URL, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMenuConfigNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *MenuConfigRepository) CreateTutorial(ctx context.Context, input TutorialUpsert) (*Tutorial, error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO tutorials (name, url) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name), input.URL)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetTutorial(ctx, id)
}

func (r *MenuConfigRepository) UpdateTutorial(ctx context.Context, id int64, input TutorialUpsert) (*Tutorial, error) {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE tutorials SET name = ?, url = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, strings.TrimSpace(input.Name), input.URL, id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrMenuConfigNotFound
	}
	return r.GetTutorial(ctx, id)
}

func (r *MenuConfigRepository) DeleteTutorial(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM tutorials WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMenuConfigNotFound
	}
	return nil
}

func (r *MenuConfigRepository) FetchCatalog(ctx context.Context) (CatalogResponse, error) {
	resp := CatalogResponse{
		Grades:              make([]Grade, 0),
		Subjects:            make([]Subject, 0),
		GradeSubjects:       make([]GradeSubject, 0),
		GradeSubjectLessons: make([]GradeSubjectLesson, 0),
		Lessons:             make([]menuconfigmodels.Lesson, 0),
		Topics:              make([]Topic, 0),
		Subtopics:           make([]Subtopic, 0),
	}

	gradesRows, err := r.db.QueryContext(ctx, "SELECT id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM grades ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer gradesRows.Close()

	for gradesRows.Next() {
		var grade Grade
		if err := gradesRows.Scan(&grade.ID, &grade.Name, &grade.CreatedAt); err != nil {
			return resp, err
		}
		resp.Grades = append(resp.Grades, grade)
	}
	if err := gradesRows.Err(); err != nil {
		return resp, err
	}

	subjectsRows, err := r.db.QueryContext(ctx, "SELECT id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subjects ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer subjectsRows.Close()

	for subjectsRows.Next() {
		var subject Subject
		if err := subjectsRows.Scan(&subject.ID, &subject.Name, &subject.CreatedAt); err != nil {
			return resp, err
		}
		resp.Subjects = append(resp.Subjects, subject)
	}
	if err := subjectsRows.Err(); err != nil {
		return resp, err
	}

	gradeSubjectsRows, err := r.db.QueryContext(ctx, "SELECT grade_id, subject_id, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM grade_subjects ORDER BY grade_id ASC, subject_id ASC")
	if err != nil {
		return resp, err
	}
	defer gradeSubjectsRows.Close()

	for gradeSubjectsRows.Next() {
		var gs GradeSubject
		if err := gradeSubjectsRows.Scan(&gs.GradeID, &gs.SubjectID, &gs.CreatedAt); err != nil {
			return resp, err
		}
		resp.GradeSubjects = append(resp.GradeSubjects, gs)
	}
	if err := gradeSubjectsRows.Err(); err != nil {
		return resp, err
	}

	gradeSubjectLessonsRows, err := r.db.QueryContext(ctx, "SELECT grade_id, subject_id, lesson_id, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM grade_subject_lessons ORDER BY grade_id ASC, subject_id ASC, lesson_id ASC")
	if err != nil {
		return resp, err
	}
	defer gradeSubjectLessonsRows.Close()

	for gradeSubjectLessonsRows.Next() {
		var gsl GradeSubjectLesson
		if err := gradeSubjectLessonsRows.Scan(&gsl.GradeID, &gsl.SubjectID, &gsl.LessonID, &gsl.CreatedAt); err != nil {
			return resp, err
		}
		resp.GradeSubjectLessons = append(resp.GradeSubjectLessons, gsl)
	}
	if err := gradeSubjectLessonsRows.Err(); err != nil {
		return resp, err
	}

	lessonsRows, err := r.db.QueryContext(ctx, "SELECT id, subject_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM lessons ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer lessonsRows.Close()

	for lessonsRows.Next() {
		var lesson menuconfigmodels.Lesson
		if err := lessonsRows.Scan(&lesson.ID, &lesson.SubjectID, &lesson.Name, &lesson.CreatedAt); err != nil {
			return resp, err
		}
		resp.Lessons = append(resp.Lessons, lesson)
	}
	if err := lessonsRows.Err(); err != nil {
		return resp, err
	}

	topicsRows, err := r.db.QueryContext(ctx, "SELECT id, lesson_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM topics ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer topicsRows.Close()

	for topicsRows.Next() {
		var topic Topic
		if err := topicsRows.Scan(&topic.ID, &topic.LessonID, &topic.Name, &topic.CreatedAt); err != nil {
			return resp, err
		}
		resp.Topics = append(resp.Topics, topic)
	}
	if err := topicsRows.Err(); err != nil {
		return resp, err
	}

	subtopicsRows, err := r.db.QueryContext(ctx, "SELECT id, topic_id, name, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') AS created_at FROM subtopics ORDER BY id ASC")
	if err != nil {
		return resp, err
	}
	defer subtopicsRows.Close()

	for subtopicsRows.Next() {
		var subtopic Subtopic
		if err := subtopicsRows.Scan(&subtopic.ID, &subtopic.TopicID, &subtopic.Name, &subtopic.CreatedAt); err != nil {
			return resp, err
		}
		resp.Subtopics = append(resp.Subtopics, subtopic)
	}
	if err := subtopicsRows.Err(); err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *MenuConfigRepository) count(ctx context.Context, query string, args ...interface{}) (int, error) {
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var total int
	if err := stmt.QueryRowContext(ctx, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func offsetFromPage(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return (page - 1) * pageSize
}

func (r *MenuConfigRepository) CheckParentExists(ctx context.Context, table string, id int64) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = ? LIMIT 1", table)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists int
	if err := stmt.QueryRowContext(ctx, id).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ensureGradeExists(ctx context.Context, tx *sql.Tx, gradeID int64) error {
	var exists int
	if err := tx.QueryRowContext(ctx, "SELECT 1 FROM grades WHERE id = ?", gradeID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("grade not found: %w", ErrMenuConfigNotFound)
		}
		return err
	}
	return nil
}

func ensureSubjectExists(ctx context.Context, tx *sql.Tx, subjectID int64) error {
	var exists int
	if err := tx.QueryRowContext(ctx, "SELECT 1 FROM subjects WHERE id = ?", subjectID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("subject not found: %w", ErrMenuConfigNotFound)
		}
		return err
	}
	return nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

func ensureLessonBelongsToSubject(ctx context.Context, tx *sql.Tx, lessonID, subjectID int64) error {
	var storedSubjectID int64
	if err := tx.QueryRowContext(ctx, "SELECT subject_id FROM lessons WHERE id = ?", lessonID).Scan(&storedSubjectID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("lesson not found: %w", ErrMenuConfigNotFound)
		}
		return err
	}
	if storedSubjectID != subjectID {
		return fmt.Errorf("lesson does not belong to subject: %w", ErrMenuConfigNotFound)
	}
	return nil
}

func ensureGradeSubjectLinkExists(ctx context.Context, tx *sql.Tx, gradeID, subjectID int64) error {
	var exists int
	if err := tx.QueryRowContext(ctx, "SELECT 1 FROM grade_subjects WHERE grade_id = ? AND subject_id = ?", gradeID, subjectID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("grade subject link not found: %w", ErrMenuConfigNotFound)
		}
		return err
	}
	return nil
}
