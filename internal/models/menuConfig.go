package models

import menuconfigmodels "github.com/tharindulakmal/sl-edu-service/internal/models/menuconfig"

type (
	MenuConfigPagedResponse[T any] = menuconfigmodels.PagedResponse[T]
	MenuConfigGrade                = menuconfigmodels.Grade
	MenuConfigGradeUpsert          = menuconfigmodels.GradeUpsert
	MenuConfigSubject              = menuconfigmodels.Subject
	MenuConfigSubjectUpsert        = menuconfigmodels.SubjectUpsert
	MenuConfigLesson               = menuconfigmodels.Lesson
	MenuConfigLessonUpsert         = menuconfigmodels.LessonUpsert
	MenuConfigTopic                = menuconfigmodels.Topic
	MenuConfigTopicUpsert          = menuconfigmodels.TopicUpsert
	MenuConfigSubtopic             = menuconfigmodels.Subtopic
	MenuConfigSubtopicUpsert       = menuconfigmodels.SubtopicUpsert
	MenuConfigTutor                = menuconfigmodels.Tutor
	MenuConfigTutorUpsert          = menuconfigmodels.TutorUpsert
	MenuConfigYear                 = menuconfigmodels.Year
	MenuConfigYearUpsert           = menuconfigmodels.YearUpsert
	MenuConfigTutorial             = menuconfigmodels.Tutorial
	MenuConfigTutorialUpsert       = menuconfigmodels.TutorialUpsert
)
