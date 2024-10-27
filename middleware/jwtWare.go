package middleware

import (
	"time"

	"github.com/abdullah-mobin/ebook-manager/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProtectedRoute() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.Secrete),
		ErrorHandler: jwtErr,
	})
}

func NewJwt(id primitive.ObjectID) (string, error) {
	claim := jwt.MapClaims{
		"aud":    "ebook-manager.idk",
		"iss":    "ebook-manager.idk",
		"iat":    time.Now().Unix(),
		"user":   id,
		"expire": time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	sign, err := token.SignedString([]byte(config.Secrete))
	if err != nil {
		return "", err
	}

	return sign, nil
}

func jwtErr(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized request",
			"message": "user not authorized to proced",
		})
	}
	return nil
}
