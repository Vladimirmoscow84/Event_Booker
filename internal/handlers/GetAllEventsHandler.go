package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) GetAllEventsHandler(c *ginext.Context) {
	events, err := r.eventsGetter.ListEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
