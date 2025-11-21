package handlers

import (
	"net/http"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) CreateEventHandler(c *ginext.Context) {
	var event model.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := r.eventCreator.CreateEvent(c.Request.Context(), &event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
