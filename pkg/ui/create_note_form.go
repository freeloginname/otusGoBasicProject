package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) CreateNoteForm(c *gin.Context) {

	c.HTML(http.StatusOK, "create_note_form.tmpl", gin.H{"title": "Создать заметку"})
}
