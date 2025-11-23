package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) CreateBookingHandler(c *gin.Context) {
	// 1. Получаем event id из path
	idStr := c.Param("id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	// 2. Получаем user_id из контекста, который положил JWT middleware
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	// 3. Делегируем в сервис (или storage), который создаёт бронь
	bookingID, err := r.bookingCreator.CreateBooking(c.Request.Context(), eventID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking_id": bookingID})
}
