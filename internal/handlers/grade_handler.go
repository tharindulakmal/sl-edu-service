package handlers

import (
	"net/http"

	"github.com/tharindulakmal/sl-edu-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	Repo repository.GradeRepositoryInterface
}

func NewGradeHandler(repo repository.GradeRepositoryInterface) *GradeHandler {
	return &GradeHandler{Repo: repo}
}

func (h *GradeHandler) GetGrades(c *gin.Context) {
	grades, err := h.Repo.GetAllGrades()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grades)
}
