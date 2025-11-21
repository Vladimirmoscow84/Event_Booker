package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type eventCreator interface {
	CreateEvent(ctx context.Context, event *model.Event) (int, error)
}

type bookingCreator interface {
	CreateBooking(ctx context.Context, eventID, userID int) (int, error)
}
type bookingPayer interface {
	ConfirmBooking(ctx context.Context, bookingID int) error
	GetBooking(ctx context.Context, id int) (*model.Booking, error)
}

type eventsGetter interface {
	GetEvent(ctx context.Context, id int) (*model.Event, error)
	ListEvents(ctx context.Context) ([]*model.Event, error)
}

type userCreator interface {
	CreateUser(ctx context.Context, email string) (int, error)
}

type Router struct {
	Router         *ginext.Engine
	eventCreator   eventCreator
	bookingCreator bookingCreator
	bookingPayer   bookingPayer
	eventsGetter   eventsGetter
	userCreator    userCreator
}

func New(router *ginext.Engine, eCreator eventCreator, bCreator bookingCreator, bPayer bookingPayer, eGetter eventsGetter, uCreator userCreator) *Router {
	return &Router{
		Router:         router,
		eventCreator:   eCreator,
		bookingCreator: bCreator,
		bookingPayer:   bPayer,
		eventsGetter:   eGetter,
		userCreator:    uCreator,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/events", r.CreateEventHandler)
	r.Router.POST("/events/:id/book", r.CreateBookingHandler)
	r.Router.POST("/events/:id/confirm", r.ConfirmBookingHandler)
	r.Router.GET("/events/:id", r.GetEventHandler)
	r.Router.GET("/events", r.GetAllEventsHandler)
	r.Router.GET("/events", r.GetAllEventsHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}
