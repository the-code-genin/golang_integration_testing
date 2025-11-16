package http

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/the-code-genin/golang_integration_testing/repository"
	"github.com/the-code-genin/golang_integration_testing/service"
)

func (s *Server) createNoteHandler(c *gin.Context) {
	var req repository.CreateNoteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request body: %v", err)
		s.sendBadRequest(c, "bad request")
		return
	}

	note, err := s.service.CreateNote(c, req)
	if err != nil {
		log.Printf("unable to create note: %v", err)

		switch {
		case errors.Is(err, service.ErrNoteTitleTaken):
			s.sendConflict(c, err.Error())
		default:
			s.sendInternalError(c, err.Error())
		}
		return
	}

	s.sendCreated(c, *note)
}

func (s *Server) fetchNoteByIDHandler(c *gin.Context) {
	id, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Printf("invalid UUID: %v", err)
		s.sendBadRequest(c, "invalid note ID")
		return
	}

	note, err := s.service.FetchNoteByID(c, id)
	if err != nil {
		log.Printf("unable to fetch note: %v", err)

		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			s.sendNotFound(c, err.Error())
		default:
			s.sendInternalError(c, err.Error())
		}
		return
	}

	s.sendOk(c, *note)
}

func (s *Server) fetchNotesHandler(c *gin.Context) {
	notes, err := s.service.FetchNotes(c)
	if err != nil {
		log.Printf("unable to fetch notes: %v", err)
		s.sendInternalError(c, err.Error())
		return
	}
	s.sendOk(c, notes)
}

func (s *Server) updateNoteHandler(c *gin.Context) {
	id, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Printf("invalid UUID: %v", err)
		s.sendBadRequest(c, "invalid note ID")
		return
	}

	var req repository.UpdateNoteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request body: %v", err)
		s.sendBadRequest(c, "bad request")
		return
	}

	note, err := s.service.UpdateNote(c, id, req)
	if err != nil {
		log.Printf("unable to update note: %v", err)

		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			s.sendNotFound(c, err.Error())
		default:
			s.sendInternalError(c, err.Error())
		}
		return
	}

	s.sendOk(c, *note)
}

func (s *Server) deleteNoteHandler(c *gin.Context) {
	id, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Printf("invalid UUID: %v", err)
		s.sendBadRequest(c, "invalid note ID")
		return
	}

	err = s.service.DeleteNote(c, id)
	if err != nil {
		log.Printf("unable to delete note: %v", err)

		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			s.sendNotFound(c, err.Error())
		default:
			s.sendInternalError(c, err.Error())
		}
		return
	}

	s.sendNoContent(c)
}
