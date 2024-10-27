package handler

import (
	"context"
	"encoding/json"

	"github.com/abdullah-mobin/ebook-manager/config"
	"github.com/abdullah-mobin/ebook-manager/middleware"
	"github.com/abdullah-mobin/ebook-manager/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(c *fiber.Ctx) error {

	var input model.User
	err := c.BodyParser(&input)
	if err != nil {
		return UserCustomErr(c, "Invalid input formate", fiber.ErrBadRequest)
	}

	if input.Name == "" {
		return UserCustomErr(c, "Name required", fiber.ErrBadRequest)
	}
	if input.Email == "" {
		return UserCustomErr(c, "Email required", fiber.ErrBadRequest)
	}
	if input.Password == "" {
		return UserCustomErr(c, "password required", fiber.ErrBadRequest)
	}
	if input.UserType == "" {
		return UserCustomErr(c, "User type required", fiber.ErrBadRequest)
	}
	if !IsValidEmail(input.Email) {
		return UserCustomErr(c, "invalid email", fiber.ErrBadRequest)
	}

	result, er := config.UserCollection.InsertOne(context.TODO(), input)
	if er != nil {
		return UserCustomErr(c, "Failed to register user, retry later", fiber.ErrInternalServerError)
	}

	input.ID = result.InsertedID.(primitive.ObjectID)

	token, er := middleware.NewJwt(input.ID)
	if er != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    fiber.ErrInternalServerError,
			"redirect": "/auth/login",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"$user_id":   input.ID,
		"$user_type": input.UserType,
		"token":      token,
		"message":    "Signup successfull",
		"redirect":   "/home",
	})
}

func Login(c *fiber.Ctx) error {

	var input, dbsUser model.User
	err := c.BodyParser(&input)
	if err != nil {
		return UserCustomErr(c, "invalid input", fiber.ErrBadRequest)
	}
	if input.Email == "" {
		return UserCustomErr(c, "Email required", fiber.ErrBadRequest)
	}
	if !IsValidEmail(input.Email) {
		return UserCustomErr(c, "invalid email", fiber.ErrBadRequest)

	}
	if input.Password == "" {
		return UserCustomErr(c, "Password required", fiber.ErrBadRequest)
	}

	email := input.Email
	err = config.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&dbsUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserCustomErr(c, "Create New Account", fiber.ErrNotFound)
		}
		return UserCustomErr(c, "Server offline, retry later", fiber.ErrInternalServerError)
	}

	if input.Password != dbsUser.Password {
		return UserCustomErr(c, "Wrong Password, Input correct password or reset password", fiber.ErrUnauthorized)
	}
	token, er := middleware.NewJwt(input.ID)
	if er != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    fiber.ErrInternalServerError,
			"redirect": "/auth/login",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":    token,
		"$user_id": dbsUser.ID,
		"message":  "Login successfull",
		"redirect": "/ebook-manager/home",
	})
}

func Googlelogin(c *fiber.Ctx) error {
	oauth2 := config.Oauth2()
	url := oauth2.AuthCodeURL("randomstate")
	return c.Redirect(url)
}

func GoogleCalback(c *fiber.Ctx) error {

	state := c.Query("state")
	if state != "randomstate" {
		return c.SendString("States don't Match!!")
	}

	oauth2 := config.Oauth2()
	code := c.FormValue("code")
	token, err := oauth2.Exchange(context.Background(), code)
	if err != nil {
		return UserCustomErr(c, "failed exchange token", err)
	}

	client := oauth2.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return UserCustomErr(c, "cant get user info", err)
	}
	defer response.Body.Close()

	var user model.GUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return UserCustomErr(c, "failed to decode user info", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "signup successful",
		"info":    user,
	})
}
