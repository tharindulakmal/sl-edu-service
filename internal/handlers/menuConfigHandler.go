package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	menuconfigmodels "github.com/tharindulakmal/sl-edu-service/internal/models/menuconfig"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
	"github.com/tharindulakmal/sl-edu-service/internal/validator"
)

type (
	Grade          = menuconfigmodels.Grade
	GradeUpsert    = menuconfigmodels.GradeUpsert
	Subject        = menuconfigmodels.Subject
	SubjectUpsert  = menuconfigmodels.SubjectUpsert
	Lesson         = menuconfigmodels.Lesson
	LessonUpsert   = menuconfigmodels.LessonUpsert
	Topic          = menuconfigmodels.Topic
	TopicUpsert    = menuconfigmodels.TopicUpsert
	Subtopic       = menuconfigmodels.Subtopic
	SubtopicUpsert = menuconfigmodels.SubtopicUpsert
	Tutor          = menuconfigmodels.Tutor
	TutorUpsert    = menuconfigmodels.TutorUpsert
	Year           = menuconfigmodels.Year
	YearUpsert     = menuconfigmodels.YearUpsert
	Tutorial       = menuconfigmodels.Tutorial
	TutorialUpsert = menuconfigmodels.TutorialUpsert
)

type PagedResponse[T any] = menuconfigmodels.PagedResponse[T]

type Handler struct {
	repo *repository.MenuConfigRepository
}

func NewHandler(repo *repository.MenuConfigRepository) *Handler {
	return &Handler{repo: repo}
}

func RegisterAdminMenuConfigRoutes(group *gin.RouterGroup, db *sql.DB) {
	repo := repository.NewMenuConfigRepository(db)
	handler := NewHandler(repo)

	group.GET("/grades", handler.listGrades)
	group.POST("/grades", handler.createGrade)
	group.GET("/grades/:id", handler.getGrade)
	group.PUT("/grades/:id", handler.updateGrade)
	group.DELETE("/grades/:id", handler.deleteGrade)

	group.GET("/subjects", handler.listSubjects)
	group.POST("/subjects", handler.createSubject)
	group.GET("/subjects/:id", handler.getSubject)
	group.PUT("/subjects/:id", handler.updateSubject)
	group.DELETE("/subjects/:id", handler.deleteSubject)

	group.GET("/lessons", handler.listLessons)
	group.POST("/lessons", handler.createLesson)
	group.GET("/lessons/:id", handler.getLesson)
	group.PUT("/lessons/:id", handler.updateLesson)
	group.DELETE("/lessons/:id", handler.deleteLesson)

	group.GET("/topics", handler.listTopics)
	group.POST("/topics", handler.createTopic)
	group.GET("/topics/:id", handler.getTopic)
	group.PUT("/topics/:id", handler.updateTopic)
	group.DELETE("/topics/:id", handler.deleteTopic)

	group.GET("/subtopics", handler.listSubtopics)
	group.POST("/subtopics", handler.createSubtopic)
	group.GET("/subtopics/:id", handler.getSubtopic)
	group.PUT("/subtopics/:id", handler.updateSubtopic)
	group.DELETE("/subtopics/:id", handler.deleteSubtopic)

	group.GET("/tutors", handler.listTutors)
	group.POST("/tutors", handler.createTutor)
	group.GET("/tutors/:id", handler.getTutor)
	group.PUT("/tutors/:id", handler.updateTutor)
	group.DELETE("/tutors/:id", handler.deleteTutor)

	group.GET("/years", handler.listYears)
	group.POST("/years", handler.createYear)
	group.GET("/years/:id", handler.getYear)
	group.PUT("/years/:id", handler.updateYear)
	group.DELETE("/years/:id", handler.deleteYear)

	group.GET("/tutorials", handler.listTutorials)
	group.POST("/tutorials", handler.createTutorial)
	group.GET("/tutorials/:id", handler.getTutorial)
	group.PUT("/tutorials/:id", handler.updateTutorial)
	group.DELETE("/tutorials/:id", handler.deleteTutorial)
}

func (h *Handler) listGrades(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	grades, total, err := h.repo.ListGrades(c.Request.Context(), search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Grade]{Data: grades, TotalCount: total})
}

func (h *Handler) getGrade(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.repo.GetGrade(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, grade)
}

func (h *Handler) createGrade(c *gin.Context) {
	var input GradeUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateGradeUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.repo.CreateGrade(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, grade)
}

func (h *Handler) updateGrade(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input GradeUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateGradeUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.repo.UpdateGrade(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, grade)
}

func (h *Handler) deleteGrade(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteGrade(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listSubjects(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	gradeIDParam := strings.TrimSpace(c.Query("gradeId"))
	var gradeID *int64
	if gradeIDParam != "" {
		id, convErr := strconv.ParseInt(gradeIDParam, 10, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gradeId"})
			return
		}
		gradeID = &id
	}

	subjects, total, err := h.repo.ListSubjects(c.Request.Context(), gradeID, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Subject]{Data: subjects, TotalCount: total})
}

func (h *Handler) getSubject(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.repo.GetSubject(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, subject)
}

func (h *Handler) createSubject(c *gin.Context) {
	var input SubjectUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateSubjectUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "grades", input.GradeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grade not found"})
		return
	}

	subject, err := h.repo.CreateSubject(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

func (h *Handler) updateSubject(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input SubjectUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateSubjectUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "grades", input.GradeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grade not found"})
		return
	}

	subject, err := h.repo.UpdateSubject(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, subject)
}

func (h *Handler) deleteSubject(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteSubject(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listLessons(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	subjectParam := strings.TrimSpace(c.Query("subjectId"))
	var subjectID *int64
	if subjectParam != "" {
		id, convErr := strconv.ParseInt(subjectParam, 10, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subjectId"})
			return
		}
		subjectID = &id
	}

	lessons, total, err := h.repo.ListLessons(c.Request.Context(), subjectID, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Lesson]{Data: lessons, TotalCount: total})
}

func (h *Handler) getLesson(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson, err := h.repo.GetLesson(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *Handler) createLesson(c *gin.Context) {
	var input LessonUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateLessonUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "subjects", input.SubjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject not found"})
		return
	}

	lesson, err := h.repo.CreateLesson(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

func (h *Handler) updateLesson(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input LessonUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateLessonUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "subjects", input.SubjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject not found"})
		return
	}

	lesson, err := h.repo.UpdateLesson(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *Handler) deleteLesson(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteLesson(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listTopics(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	lessonParam := strings.TrimSpace(c.Query("lessonId"))
	var lessonID *int64
	if lessonParam != "" {
		id, convErr := strconv.ParseInt(lessonParam, 10, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lessonId"})
			return
		}
		lessonID = &id
	}

	topics, total, err := h.repo.ListTopics(c.Request.Context(), lessonID, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Topic]{Data: topics, TotalCount: total})
}

func (h *Handler) getTopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topic, err := h.repo.GetTopic(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *Handler) createTopic(c *gin.Context) {
	var input TopicUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTopicUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "lessons", input.LessonID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson not found"})
		return
	}

	topic, err := h.repo.CreateTopic(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, topic)
}

func (h *Handler) updateTopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input TopicUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTopicUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "lessons", input.LessonID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson not found"})
		return
	}

	topic, err := h.repo.UpdateTopic(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (h *Handler) deleteTopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteTopic(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listSubtopics(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	topicParam := strings.TrimSpace(c.Query("topicId"))
	var topicID *int64
	if topicParam != "" {
		id, convErr := strconv.ParseInt(topicParam, 10, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topicId"})
			return
		}
		topicID = &id
	}

	subtopics, total, err := h.repo.ListSubtopics(c.Request.Context(), topicID, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Subtopic]{Data: subtopics, TotalCount: total})
}

func (h *Handler) getSubtopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtopic, err := h.repo.GetSubtopic(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, subtopic)
}

func (h *Handler) createSubtopic(c *gin.Context) {
	var input SubtopicUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateSubtopicUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "topics", input.TopicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic not found"})
		return
	}

	subtopic, err := h.repo.CreateSubtopic(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subtopic)
}

func (h *Handler) updateSubtopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input SubtopicUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateSubtopicUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := h.repo.CheckParentExists(c.Request.Context(), "topics", input.TopicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic not found"})
		return
	}

	subtopic, err := h.repo.UpdateSubtopic(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, subtopic)
}

func (h *Handler) deleteSubtopic(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteSubtopic(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listTutors(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	tutors, total, err := h.repo.ListTutors(c.Request.Context(), search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Tutor]{Data: tutors, TotalCount: total})
}

func (h *Handler) getTutor(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutor, err := h.repo.GetTutor(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, tutor)
}

func (h *Handler) createTutor(c *gin.Context) {
	var input TutorUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTutorUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutor, err := h.repo.CreateTutor(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tutor)
}

func (h *Handler) updateTutor(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input TutorUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTutorUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutor, err := h.repo.UpdateTutor(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, tutor)
}

func (h *Handler) deleteTutor(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteTutor(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listYears(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	years, total, err := h.repo.ListYears(c.Request.Context(), search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Year]{Data: years, TotalCount: total})
}

func (h *Handler) getYear(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	year, err := h.repo.GetYear(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, year)
}

func (h *Handler) createYear(c *gin.Context) {
	var input YearUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateYearUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	year, err := h.repo.CreateYear(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, year)
}

func (h *Handler) updateYear(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input YearUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateYearUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	year, err := h.repo.UpdateYear(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, year)
}

func (h *Handler) deleteYear(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteYear(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) listTutorials(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	tutorials, total, err := h.repo.ListTutorials(c.Request.Context(), search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[Tutorial]{Data: tutorials, TotalCount: total})
}

func (h *Handler) getTutorial(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutorial, err := h.repo.GetTutorial(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, tutorial)
}

func (h *Handler) createTutorial(c *gin.Context) {
	var input TutorialUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTutorialUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutorial, err := h.repo.CreateTutorial(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tutorial)
}

func (h *Handler) updateTutorial(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input TutorialUpsert
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateTutorialUpsert(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tutorial, err := h.repo.UpdateTutorial(c.Request.Context(), id, input)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, tutorial)
}

func (h *Handler) deleteTutorial(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteTutorial(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func parsePagination(c *gin.Context) (int, int, error) {
	pageStr := strings.TrimSpace(c.DefaultQuery("page", "1"))
	pageSizeStr := strings.TrimSpace(c.DefaultQuery("pageSize", "10"))

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid page parameter")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pageSize parameter")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize, nil
}

func parseIDParam(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

func handleRepoError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, repository.ErrMenuConfigNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
