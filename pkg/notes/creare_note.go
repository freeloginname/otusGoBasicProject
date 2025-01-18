package notes

import (
	"net/http"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/notes"
	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
)

// type Note struct {
// 	ID     pgtype.UUID `db:"id" json:"id"`
// 	UserID pgtype.UUID `db:"user_id" json:"user_id"`
// 	Name   string      `db:"name" json:"name"`
// 	Text   string      `db:"text" json:"text"`
// }

type CreateNoteRequestBody struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	UserName string `json:"user_name"`
}

func (h handler) CreateNote(c *gin.Context) {
	body := CreateNoteRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		// c.HTML(http.StatusBadRequest, "unauthorized.tmpl", gin.H{"error": "Произошла ошибка"})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userName := body.UserName
	if userName == "" {
		anyTypeUserName, _ := c.Get("currentUser")
		userName = anyTypeUserName.(string)
	}

	user, _ := c.Get("currentUser")
	if user != userName {
		// c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Доступ запрещен"})
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Доступ запрещен"})
		return
	}

	var note notes.Note

	note.Name = body.Name
	note.Text = body.Text
	// TODO добавить авторизацию

	_, err := transaction.CreateNote(c, h.DB, userName, note.Name, note.Text)
	if err != nil {
		// c.HTML(http.StatusBadRequest, "unauthorized.tmpl", gin.H{"error": "Произошла ошибка"})
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	c.JSON(http.StatusCreated, map[string]string{"success": "Note Created"})
}
