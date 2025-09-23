package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

type SubjectHandler struct {
	Repo *repository.SubjectRepository
}

func NewSubjectHandler(repo *repository.SubjectRepository) *SubjectHandler {
	return &SubjectHandler{Repo: repo}
}

func (h *SubjectHandler) GetSubjectsByGrade(c *gin.Context) {
	gradeIDStr := c.Query("gradeId")
	gradeID, err := strconv.Atoi(gradeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gradeId"})
		return
	}

	subjects, err := h.Repo.GetSubjectsByGradeID(gradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}
