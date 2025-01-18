package notes

import (
	"net/http"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
)

// type Note struct {
// 	ID     pgtype.UUID `db:"id" json:"id"`
// 	UserID pgtype.UUID `db:"user_id" json:"user_id"`
// 	Name   string      `db:"name" json:"name"`
// 	Text   string      `db:"text" json:"text"`
// }

func (h handler) DeleteNote(c *gin.Context) {

	user, _ := c.Get("currentUser")
	userName := user.(string)
	if userName == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Доступ запрещен"})
		// c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Доступ запрещен"})
		return
	}

	noteName := c.Param("name")

	err := transaction.DeleteUserNoteByName(c, h.DB, userName, noteName)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	c.JSON(http.StatusCreated, map[string]string{"success": "Note Deleted or does not exist"})
}
