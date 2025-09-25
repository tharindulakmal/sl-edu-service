package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tharindulakmal/sl-edu-service/internal/models"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

func TestGetAllGrades(t *testing.T) {
	// Arrange
	mockRepo := &repository.MockGradeRepository{
		Grades: []models.Grade{
			{ID: 1, Grade: "Grade 1"},
			{ID: 2, Grade: "Grade 2"},
		},
		Err: nil,
	}

	handler := NewGradeHandler(mockRepo)

	router := gin.Default()
	router.GET("/grades", handler.GetGrades)

	req, _ := http.NewRequest("GET", "/grades", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Grade 1")
	assert.Contains(t, w.Body.String(), "Grade 2")
}
