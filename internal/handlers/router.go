package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Event_Booker/internal/middleware"
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

type authService interface {
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type Router struct {
	Router         *ginext.Engine
	eventCreator   eventCreator
	bookingCreator bookingCreator
	bookingPayer   bookingPayer
	eventsGetter   eventsGetter
	userCreator    userCreator
	authService    authService
}

func New(router *ginext.Engine, eCreator eventCreator, bCreator bookingCreator, bPayer bookingPayer, eGetter eventsGetter, uCreator userCreator, aService authService) *Router {
	return &Router{
		Router:         router,
		eventCreator:   eCreator,
		bookingCreator: bCreator,
		bookingPayer:   bPayer,
		eventsGetter:   eGetter,
		userCreator:    uCreator,
		authService:    aService,
	}
}

func (r *Router) Routes(jwtSecret string) {

	authMiddleware := middleware.JWTMiddleware(jwtSecret)
	adminMiddleware := middleware.AdminOnly()

	//Открытые маршруты
	r.Router.POST("/auth/register", r.RegisterHandler)
	r.Router.POST("/auth/login", r.LoginHandler)
	r.Router.GET("/events/:id", r.GetEventHandler)
	r.Router.GET("/events", r.GetAllEventsHandler)

	// Юзерские маршруты
	r.Router.POST("/events/:id/book", authMiddleware, r.CreateBookingHandler)
	r.Router.POST("/events/:id/confirm", authMiddleware, r.ConfirmBookingHandler)

	//Админские маршруты
	r.Router.POST("/users", authMiddleware, adminMiddleware, r.CreateUserHandler)
	r.Router.POST("/events", authMiddleware, adminMiddleware, r.CreateEventHandler)

	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.GET("/user", func(c *gin.Context) { c.File("./web/user.html") })
	r.Router.GET("/admin", func(c *gin.Context) { c.File("./web/admin.html") })
	r.Router.Static("/static", "./web")
}
