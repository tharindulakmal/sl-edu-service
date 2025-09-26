package handlers

import (
	"net/http"
	"strconv"

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

	questions, err := h.repo.GetList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
