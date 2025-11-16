package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/golang_integration_testing/service"
)

type Server struct {
	service service.Service
	router  *gin.Engine
}

func (s *Server) sendNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": message,
	})
}

func (s *Server) sendBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": message,
	})
}

func (s *Server) sendInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}

func (s *Server) sendCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func (s *Server) sendOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func (s *Server) sendNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func NewServer(svc service.Service) *Server {
	router := gin.Default()

	server := &Server{
		service: svc,
		router:  router,
	}

	g := router.Group("/v1/notes")
	{
		g.POST("", server.createNoteHandler)
		g.GET("", server.fetchNotesHandler)
		g.GET("/:id", server.fetchNoteByIDHandler)
		g.PATCH("/:id", server.updateNoteHandler)
		g.DELETE("/:id", server.deleteNoteHandler)
	}

	router.NoRoute(func(c *gin.Context) {
		server.sendNotFound(c, "The route you requested for was not found on this server")
	})

	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) Handler() http.Handler {
	return s.router.Handler()
}
