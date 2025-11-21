package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) ConfirmBookingHandler(c *ginext.Context) {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := r.bookingPayer.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	event, err := r.eventsGetter.GetEvent(c.Request.Context(), booking.EventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	if !event.RequiresPayment {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event does not require payment"})
		return
	}

	if err := r.bookingPayer.ConfirmBooking(c.Request.Context(), bookingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "confirmed"})
}
