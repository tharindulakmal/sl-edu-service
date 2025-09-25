package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

type LessonHandler struct {
	Repo repository.LessonRepositoryInterface
}

func NewLessonHandler(repo repository.LessonRepositoryInterface) *LessonHandler {
	return &LessonHandler{Repo: repo}
}

func (h *LessonHandler) GetLessons(c *gin.Context) {
	subjectIdStr := c.Query("subject")
	subjectId, err := strconv.Atoi(subjectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subjectId"})
		return
	}

	lessons, err := h.Repo.GetLessonsBySubject(subjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lessons)
}
