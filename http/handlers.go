package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/the-code-genin/golang_integration_testing/repository"
)

func (s *Server) createNoteHandler(c *gin.Context) {
	var req repository.CreateNoteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request body: %v", err)
		s.sendBadRequest(c, err.Error())
		return
	}

	note, err := s.service.CreateNote(c.Request.Context(), req)
	if err != nil {
		log.Printf("unable to create note: %v", err)
		s.sendInternalError(c, err.Error())
		return
	}

	s.sendCreated(c, *note)
}
