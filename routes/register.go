package routes

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID must be a number"})
		return
	}
	event, err := models.GetEventById(eventId)
	if errors.Is(err, models.ErrEventNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event"})
		return
	}

	err = event.RegisterUser(userId)
	if errors.Is(err, models.ErrAlreadyRegistered) {
		context.JSON(http.StatusConflict, gin.H{"message": "You are already registered for this event"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register for the event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Successfully registered for the event"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID must be a number"})
		return
	}
	event, err := models.GetEventById(eventId)
	if errors.Is(err, models.ErrEventNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event"})
		return
	}

	err = event.CancelUserRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to cancel registration"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Registration cancelled successfully"})
}
