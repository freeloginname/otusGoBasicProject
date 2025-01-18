package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) CreateUserForm(c *gin.Context) {

	c.HTML(http.StatusOK, "create_user_form.tmpl", gin.H{"title": "Создать пользователя"})
}
