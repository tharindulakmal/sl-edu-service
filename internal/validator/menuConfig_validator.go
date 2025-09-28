package validator

import (
	"errors"
	"strings"

	menuconfigmodels "github.com/tharindulakmal/sl-edu-service/internal/models/menuconfig"
)

var (
	errNameRequired = errors.New("name is required")
	errNameTooLong  = errors.New("name must be at most 120 characters")
)

func validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errNameRequired
	}
	if len(trimmed) > 120 {
		return errNameTooLong
	}
	return nil
}

func ValidateGradeUpsert(input menuconfigmodels.GradeUpsert) error {
	return validateName(input.Name)
}

func ValidateSubjectUpsert(input menuconfigmodels.SubjectUpsert) error {
	if input.GradeID == 0 {
		return errors.New("gradeId is required")
	}
	return validateName(input.Name)
}

func ValidateLessonUpsert(input menuconfigmodels.LessonUpsert) error {
	if input.SubjectID == 0 {
		return errors.New("subjectId is required")
	}
	return validateName(input.Name)
}

func ValidateTopicUpsert(input menuconfigmodels.TopicUpsert) error {
	if input.LessonID == 0 {
		return errors.New("lessonId is required")
	}
	return validateName(input.Name)
}

func ValidateSubtopicUpsert(input menuconfigmodels.SubtopicUpsert) error {
	if input.TopicID == 0 {
		return errors.New("topicId is required")
	}
	return validateName(input.Name)
}

func ValidateTutorUpsert(input menuconfigmodels.TutorUpsert) error {
	return validateName(input.Name)
}

func ValidateYearUpsert(input menuconfigmodels.YearUpsert) error {
	if input.Value < 1900 || input.Value > 2100 {
		return errors.New("value must be between 1900 and 2100")
	}
	return nil
}

func ValidateTutorialUpsert(input menuconfigmodels.TutorialUpsert) error {
	return validateName(input.Name)
}
