package menuconfig

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PagedResponse is a generic response wrapper for paginated data.
type PagedResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"totalCount"`
}

func parseFlexibleInt64(raw json.RawMessage, field string) (int64, error) {
	if len(raw) == 0 {
		return 0, nil
	}

	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return 0, nil
	}

	var asInt int64
	if err := json.Unmarshal(raw, &asInt); err == nil {
		return asInt, nil
	}

	var asStr string
	if err := json.Unmarshal(raw, &asStr); err == nil {
		asStr = strings.TrimSpace(asStr)
		if asStr == "" {
			return 0, nil
		}

		parsed, err := strconv.ParseInt(asStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid integer", field)
		}
		return parsed, nil
	}

	return 0, fmt.Errorf("%s must be a number or numeric string", field)
}

type Grade struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type GradeUpsert struct {
	Name string `json:"name"`
}

type Subject struct {
	ID        int64  `json:"id"`
	GradeID   int64  `json:"gradeId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type SubjectUpsert struct {
	GradeID int64  `json:"gradeId"`
	Name    string `json:"name"`
}

func (s *SubjectUpsert) UnmarshalJSON(data []byte) error {
	type subjectUpsertJSON struct {
		GradeID json.RawMessage `json:"gradeId"`
		Name    string          `json:"name"`
	}

	var aux subjectUpsertJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.Name = aux.Name
	if len(aux.GradeID) == 0 {
		s.GradeID = 0
		return nil
	}

	var idInt int64
	if err := json.Unmarshal(aux.GradeID, &idInt); err == nil {
		s.GradeID = idInt
		return nil
	}

	var idStr string
	if err := json.Unmarshal(aux.GradeID, &idStr); err == nil {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			s.GradeID = 0
			return nil
		}

		parsed, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return fmt.Errorf("gradeId must be a valid integer")
		}
		s.GradeID = parsed
		return nil
	}

	return fmt.Errorf("gradeId must be a number or numeric string")
}

type Lesson struct {
	ID        int64  `json:"id"`
	SubjectID int64  `json:"subjectId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type LessonUpsert struct {
	SubjectID int64  `json:"subjectId"`
	Name      string `json:"name"`
}

func (l *LessonUpsert) UnmarshalJSON(data []byte) error {
	type lessonUpsertJSON struct {
		SubjectID json.RawMessage `json:"subjectId"`
		Name      string          `json:"name"`
	}

	var aux lessonUpsertJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	l.Name = aux.Name
	subjectID, err := parseFlexibleInt64(aux.SubjectID, "subjectId")
	if err != nil {
		return err
	}
	l.SubjectID = subjectID
	return nil
}

type Topic struct {
	ID        int64  `json:"id"`
	LessonID  int64  `json:"lessonId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type TopicUpsert struct {
	LessonID int64  `json:"lessonId"`
	Name     string `json:"name"`
}

func (t *TopicUpsert) UnmarshalJSON(data []byte) error {
	type topicUpsertJSON struct {
		LessonID json.RawMessage `json:"lessonId"`
		Name     string          `json:"name"`
	}

	var aux topicUpsertJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Name = aux.Name
	lessonID, err := parseFlexibleInt64(aux.LessonID, "lessonId")
	if err != nil {
		return err
	}
	t.LessonID = lessonID
	return nil
}

type Subtopic struct {
	ID        int64  `json:"id"`
	TopicID   int64  `json:"topicId"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type SubtopicUpsert struct {
	TopicID int64  `json:"topicId"`
	Name    string `json:"name"`
}

func (s *SubtopicUpsert) UnmarshalJSON(data []byte) error {
	type subtopicUpsertJSON struct {
		TopicID json.RawMessage `json:"topicId"`
		Name    string          `json:"name"`
	}

	var aux subtopicUpsertJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.Name = aux.Name
	topicID, err := parseFlexibleInt64(aux.TopicID, "topicId")
	if err != nil {
		return err
	}
	s.TopicID = topicID
	return nil
}

type Tutor struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"createdAt"`
}

type TutorUpsert struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type Year struct {
	ID        int64  `json:"id"`
	Value     int    `json:"value"`
	CreatedAt string `json:"createdAt"`
}

type YearUpsert struct {
	Value int `json:"value"`
}

type Tutorial struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	URL       *string `json:"url"`
	CreatedAt string  `json:"createdAt"`
}

type TutorialUpsert struct {
	Name string  `json:"name"`
	URL  *string `json:"url"`
}
