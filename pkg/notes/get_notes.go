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

func (h handler) GetNotes(c *gin.Context) {
	user, _ := c.Get("currentUser")
	userName := user.(string)

	// userName := c.Param("user_name")
	// if user != userName {
	if userName == "" {
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Доступ запрещен"})
		return
	}
	notesList, err := transaction.GetAllUserNotes(c, h.DB, userName)
	if err != nil {
		c.HTML(http.StatusBadRequest, "unauthorized.tmpl", gin.H{"error": "Не удалось получить данные о заметках"})
		// c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	// Для передачи в api:
	// c.JSON(http.StatusCreated, &notesList)
	c.HTML(http.StatusOK, "notes_get_notes.tmpl",
		gin.H{"Notes": notesList, "title": "Заметки пользователя", "userName": userName})
}
