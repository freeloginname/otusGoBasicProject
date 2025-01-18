package ui

import (
	"net/http"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
)

// Deprecated
func (h handler) GetNotes(c *gin.Context) {

	//http.Get()
	user, _ := c.Get("currentUser")

	userName := c.Param("user_name")
	if user != userName {
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Доступ запрещен"})
		return
	}
	notesList, err := transaction.GetAllUserNotes(c, h.DB, userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	c.HTML(http.StatusOK, "notes_get_notes.tmpl", gin.H{"Notes": notesList, "title": "Заметки пользователя", "user_name": userName})
}
