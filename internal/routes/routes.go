package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/tharindulakmal/sl-edu-service/internal/handlers"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api/v1")

	// Grades
	gradeRepo := repository.NewGradeRepository(db)
	gradeHandler := handlers.NewGradeHandler(gradeRepo)
	grades := api.Group("/grades")
	{
		grades.GET("", gradeHandler.GetGrades)
	}

	subjectRepo := repository.NewSubjectRepository(db)
	subjectHandler := handlers.NewSubjectHandler(subjectRepo)
	subjects := api.Group("/subjects")
	{
		subjects.GET("", subjectHandler.GetSubjectsByGrade)
	}

	lessonRepo := repository.NewLessonRepository(db)
	lessonHandler := handlers.NewLessonHandler(lessonRepo)
	lessons := api.Group("/lessons")

	lessons.GET("", lessonHandler.GetLessons)

	topicRepo := repository.NewTopicRepository(db)
	topicHandler := handlers.NewTopicHandler(topicRepo)

	api.GET("/tutor/topics", topicHandler.GetTopics)

	smartNoteRepo := repository.NewSmartNoteRepository(db)
	smartNoteHandler := handlers.NewSmartNoteHandler(smartNoteRepo)

	note := api.Group("/note")

	note.GET("/smartnote/:lessonId/:topicId/:subId", smartNoteHandler.GetSmartNote)

	questionRepo := repository.NewQuestionRepository(db)
	questionHandler := handlers.NewQuestionHandler(questionRepo)

	question := api.Group("/mcq")
	{
		question.GET("/questions/:id", questionHandler.GetQuestionByID)
		question.GET("/questions", questionHandler.GetQuestions)

		question.POST("/questions", questionHandler.CreateQuestion)
		question.PUT("/questions/:id", questionHandler.UpdateQuestion)
		question.DELETE("/questions/:id", questionHandler.DeleteQuestion)
	}

}
