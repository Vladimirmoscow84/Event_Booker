package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) CreateBookingHandler(c *ginext.Context) {
	idStr := c.Param("id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request struct {
		UserID int `json:"user_id" binding:"required"`
	}
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingID, err := r.bookingCreator.CreateBooking(c.Request.Context(), eventID, request.UserID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"booking_id": bookingID})
}
