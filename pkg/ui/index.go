package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) Index(c *gin.Context) {

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "Главная страница"})
}
