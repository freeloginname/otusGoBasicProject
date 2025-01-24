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

func (h handler) GetNote(c *gin.Context) {
	user, _ := c.Get("currentUser")
	userName := user.(string)

	// userName := c.Param("user_name")
	// if user != userName {
	if userName == "" {
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Доступ запрещен"})
		return
	}

	name := c.Param("name")

	note, err := transaction.GetUserNoteByName(c, h.DB, userName, name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "unauthorized.tmpl", gin.H{"error": "Не удалось получить данные о заметке"})
		// c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }
	// Для передачи в api:
	/// c.JSON(http.StatusOK, &note)
	c.HTML(http.StatusOK, "notes_get_note.tmpl",
		gin.H{"Note": note, "title": "Заметки пользователя", "user_name": userName, "note_name": name})
}
