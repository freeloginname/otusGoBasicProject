package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/freeloginname/otusGoBasicProject/internal/repository/transaction"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type handler struct {
	DB        *pgxpool.Pool
	SecretKey []byte
}

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool, secretKey []byte) {
	h := &handler{
		DB:        db,
		SecretKey: secretKey,
	}

	routes := r.Group("/users")
	routes.POST("/", h.CreateUser)
	routes.POST("/login", h.LoginUser)
	routes.GET("/", h.CheckAuth, h.GetUserProfile)
}

func (h handler) CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	var tokenString string

	if authHeader == "" {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.Set("error", "Authorization header is missing or user session is not established")
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing or user session is not established"})
			c.Abort()
			// c.Next()
			// c.AbortWithStatus(http.StatusUnauthorized)
			c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Authorization header is missing or user session is not established"})
			return
		}
		tokenString = cookie
		// c.String(http.StatusOK, "Cookie value: %s", cookie)

		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		// c.AbortWithStatus(http.StatusUnauthorized)
		// return
	} else {
		tokenString = authHeader
	}

	// authToken := strings.Split(authHeader, " ")
	// if len(authToken) != 2 || authToken[0] != "Bearer" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// tokenString := authToken[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.SecretKey, nil
	})
	if err != nil || !token.Valid {
		c.Set("error", "Invalid or expired token")
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		// c.Next()
		// c.AbortWithStatus(http.StatusUnauthorized)
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Set("error", "Invalid token")
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		// c.Next()
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "Invalid token"})
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.Set("error", "token expired")
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.Abort()
		// c.Next()
		// c.AbortWithStatus(http.StatusUnauthorized)
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "token expired"})
		return
	}

	userFound, err := transaction.GetUser(c, h.DB, claims["name"].(string))
	if err != nil {
		c.Abort()
		// c.Next()
		// c.AbortWithStatus(http.StatusUnauthorized)
		c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"error": "User does not exists"})
		return
	}

	c.Set("currentUser", userFound.Name)
	c.Next()
}
