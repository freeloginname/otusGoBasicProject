package users

import (
	"net/http"
	"time"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/notes"
	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequestBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h Handler) LoginUser(c *gin.Context) {
	body := LoginUserRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var user notes.User
	user.Name = body.Name
	user.Password = body.Password

	userFound, err := transaction.GetUser(c, h.DB, user.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist or invalid password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userFound.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString(h.SecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}
	// var cookiesPath strings.Builder
	// cookiesPath.WriteString("/notes/")
	// cookiesPath.WriteString(userFound.Name)

	// поменять в продакшене параметры для куки
	// c.SetCookie("token", token, 86400, cookiesPath.String(), "localhost", false, false)
	c.SetCookie("token", token, 86400, "/", "localhost:8080", false, false)

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}
