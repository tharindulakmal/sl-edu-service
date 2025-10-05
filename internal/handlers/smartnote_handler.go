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
	gradeParam := c.Query("gradeId")
	if gradeParam == "" {
		gradeParam = c.Query("GradeId")
	}
	subjectParam := c.Query("subjectId")
	if subjectParam == "" {
		subjectParam = c.Query("SubjectId")
	}
	lessonParam := c.Query("lessonId")
	if lessonParam == "" {
		lessonParam = c.Query("LessonId")
	}

	gradeID, err := strconv.Atoi(gradeParam)
	if err != nil || gradeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gradeId"})
		return
	}

	subjectID, err := strconv.Atoi(subjectParam)
	if err != nil || subjectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subjectId"})
		return
	}

	lessonID, err := strconv.Atoi(lessonParam)
	if err != nil || lessonID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lessonId"})
		return
	}

	var topicIDPtr *int
	if topicParam := c.Query("topicId"); topicParam != "" {
		if topicVal, convErr := strconv.Atoi(topicParam); convErr == nil {
			topicIDPtr = &topicVal
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topicId"})
			return
		}
	}

	var subIDPtr *int
	if subParam := c.Query("subtopicId"); subParam != "" {
		if subVal, convErr := strconv.Atoi(subParam); convErr == nil {
			subIDPtr = &subVal
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subtopicId"})
			return
		}
	}

	sn, err := h.Repo.GetSmartNote(gradeID, subjectID, lessonID, topicIDPtr, subIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sn)
}
