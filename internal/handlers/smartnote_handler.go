package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

type SmartNoteHandler struct {
	Repo repository.SmartNoteRepositoryInterface
}

func NewSmartNoteHandler(repo repository.SmartNoteRepositoryInterface) *SmartNoteHandler {
	return &SmartNoteHandler{Repo: repo}
}

func (h *SmartNoteHandler) GetSmartNote(c *gin.Context) {
	lessonIdStr := c.Param("lessonId")
	topicIdStr := c.Param("topicId")
	subIdStr := c.Param("subId")

	lessonId, err1 := strconv.Atoi(lessonIdStr)
	topicId, err2 := strconv.Atoi(topicIdStr)
	subId, err3 := strconv.Atoi(subIdStr)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	sn, err := h.Repo.GetSmartNote(lessonId, topicId, subId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sn)
}
