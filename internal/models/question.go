package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Custom type for []string stored as JSON in DB
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	if err := json.Unmarshal(bytes, a); err != nil {
		// fallback: return empty slice instead of nil on bad data
		*a = []string{}
		return nil
	}
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type Question struct {
	ID            int         `json:"id" db:"id"`
	GradeID       int         `json:"gradeId" db:"grade_id"`
	LessonID      int         `json:"lessonId" db:"lesson_id"`
	TopicID       *int        `json:"topicId,omitempty" db:"topic_id"`
	SubtopicID    *int        `json:"subtopicId,omitempty" db:"subtopic_id"`
	TutorID       *int        `json:"tutorId,omitempty" db:"tutor_id"`
	TuteID        *int        `json:"tuteId,omitempty" db:"tute_id"`
	Question      string      `json:"question" db:"question"`
	QuestionImg   *string     `json:"questionImgUrl,omitempty" db:"question_img_url"`
	CorrectAnswer string      `json:"correctAnswer" db:"correct_answer"`
	Theory        *string     `json:"theory,omitempty" db:"theory"`
	Solution      *string     `json:"solution,omitempty" db:"solution"`
	OtherAnswers  StringArray `json:"otherAnswers" db:"other_answers"` // stored as JSON
	CreatedAt     string      `json:"createdAt" db:"created_at"`
}
