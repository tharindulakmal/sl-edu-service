package handlers

import (
	"net/http"
	"strconv"

	"github.com/tharindulakmal/sl-edu-service/internal/models"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	Repo repository.TopicRepositoryInterface
}

func NewTopicHandler(repo repository.TopicRepositoryInterface) *TopicHandler {
	return &TopicHandler{Repo: repo}
}

func (h *TopicHandler) GetTopics(c *gin.Context) {
	lessonIdStr := c.Query("lessonId")
	lessonId, err := strconv.Atoi(lessonIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lessonId"})
		return
	}

	topics, err := h.Repo.GetTopicsByLesson(lessonId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defaultSN, err := h.Repo.GetDefaultSmartNote(lessonId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.TopicsResponse{
		Topics:           topics,
		DefaultSmartNote: defaultSN,
	}

	c.JSON(http.StatusOK, response)
}
