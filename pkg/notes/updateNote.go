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

type UpdateNoteRequestBody struct {
	Text string `json:"text"`
}

func (h handler) UpdateNote(c *gin.Context) {
	body := UpdateNoteRequestBody{}
	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, _ := c.Get("currentUser")
	userName := user.(string)

	if userName == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Доступ запрещен"})
		return
	}

	noteName := c.Param("name")
	// text := c.Param("text")

	err := transaction.UpdateNote(c, h.DB, userName, noteName, body.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }

	c.JSON(http.StatusCreated, map[string]string{"success": "Note Updated"})
}
