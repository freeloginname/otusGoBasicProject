package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) LoginUserForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login_user_form.tmpl", gin.H{"title": "Вход"})
}
