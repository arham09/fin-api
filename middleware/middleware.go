package middleware

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var signingKey = []byte("aqOeh4ck3R")

type Middleware struct {
}

// InitMiddleware intialize the middleware
func InitMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")

		if auth == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Unauthorized",
			})
		}

		splitToken := strings.Split(auth, "Bearer")
		validToken := strings.TrimSpace(splitToken[1])

		token, err := jwt.Parse(validToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return signingKey, nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": err.Error(),
			})
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "JWT is invalid",
			})
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		c.Set("userId", int(claims["userid"].(float64)))

		return next(c)
	}
}

// CORS will handle the CORS middleware
func (m *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}
