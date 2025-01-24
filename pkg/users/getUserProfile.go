package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetUserProfile(c *gin.Context) {
	user, _ := c.Get("currentUser")

	// c.JSON(200, gin.H{
	// 	"user": user,
	// })
	// previousError, _ := c.Get("error")
	// if c.IsAborted() {
	// 	c.HTML(http.StatusUnauthorized, "unauthorized.tmpl",
	// 	 gin.H{"title": "Информация о пользователе", "error": previousError})
	// }
	c.HTML(http.StatusOK, "get_user_profile_form.tmpl", gin.H{"title": "Информация о пользователе", "user": user})
}
