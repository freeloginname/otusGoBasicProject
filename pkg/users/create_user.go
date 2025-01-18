package users

import (
	"net/http"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/notes"
	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequestBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h handler) CreateUser(c *gin.Context) {
	body := CreateUserRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var user notes.User

	user.Name = body.Name
	user.Password = body.Password

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err = transaction.CreateUser(c, h.DB, user.Name, string(passwordHash))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// if result := h.DB.Create(&user); result.Error != nil {
	// 	c.AbortWithError(http.StatusNotFound, result.Error)
	// 	return
	// }
	c.JSON(http.StatusCreated, map[string]string{"success": "User Created"})

	// c.HTML(http.StatusOK, "users_create_user.tmpl", gin.H{"title": "Создать пользователя"})
}
