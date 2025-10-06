package handlers

import (
	"net/http"
	"strconv"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	repo repository.QuestionRepository
}

func NewQuestionHandler(repo repository.QuestionRepository) *QuestionHandler {
	return &QuestionHandler{repo: repo}
}

// GET /api/v1/tutor/questions/:id
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	question, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}
	c.JSON(http.StatusOK, question)
}

// GET /api/v1/tutor/questions?lessonId=1&page=1&pageSize=10
func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	filters := map[string]interface{}{}
	if v := c.Query("lessonId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["lessonId"] = id
		}
	}
	if v := c.Query("subjectId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["subjectId"] = id
		}
	}
	if v := c.Query("gradeId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["gradeId"] = id
		}
	}
	if v := c.Query("topicId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["topicId"] = id
		}
	}
	if v := c.Query("subtopicId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["subtopicId"] = id
		}
	}
	if v := c.Query("tutorId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["tutorId"] = id
		}
	}
	if v := c.Query("tuteId"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters["tuteId"] = id
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// get list
	questions, err := h.repo.GetList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get total count
	totalCount, err := h.repo.Count(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return with metadata
	c.JSON(http.StatusOK, gin.H{
		"data":       questions,
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": totalCount,
	})
}

// POST /api/v1/tutor/questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var q models.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if q.GradeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeId is required"})
		return
	}
	id, err := h.repo.Create(&q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	q.ID = int(id)
	c.JSON(http.StatusCreated, q)
}

// PUT /api/v1/tutor/questions/:id
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	var q models.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	q.ID = id
	if q.GradeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeId is required"})
		return
	}
	if err := h.repo.Update(&q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

// DELETE /api/v1/tutor/questions/:id
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
