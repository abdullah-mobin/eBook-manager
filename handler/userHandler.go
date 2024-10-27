package handler

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/abdullah-mobin/ebook-manager/config"
	"github.com/abdullah-mobin/ebook-manager/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserCustomErr(c *fiber.Ctx, message string, err error) error {
	return c.JSON(fiber.Map{
		"error":      err,
		"suggestion": message,
	})
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CreateUser(c *fiber.Ctx) error {

	var input model.User

	err := c.BodyParser(&input)
	if err != nil {
		return UserCustomErr(c, "Invalid input formate", err)
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
		fmt.Println("see i'm ok")
		return UserCustomErr(c, "User type required", fiber.ErrBadRequest)
	}
	if !IsValidEmail(input.Email) {
		return UserCustomErr(c, "invalid email", fiber.ErrBadRequest)
	}

	result, er := config.UserCollection.InsertOne(context.TODO(), input)
	if er != nil {
		return UserCustomErr(c, "Failed to register user, retry later", er)
	}

	input.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":   input.ID,
		"user_type": input.UserType,
		"message":   "user created succesfully",
	})
}

func GetUserByEmail(c *fiber.Ctx) error {
	var user model.User
	email := c.Query("email")

	err := config.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserCustomErr(c, "User not found", fiber.ErrNotFound)
		}
		return c.JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error finding account",
		})
	}

	return c.JSON(user)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []model.User

	cursor, err := config.UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return UserCustomErr(c, "Failed to query users", fiber.ErrInternalServerError)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var u model.User
		err := cursor.Decode(&u)
		if err != nil {
			return UserCustomErr(c, "Error decoding collection", fiber.ErrInternalServerError)
		}
		users = append(users, u)
	}

	err = cursor.Err()
	if err != nil {
		return UserCustomErr(c, "Error iterating collection", fiber.ErrInternalServerError)
	}

	return c.JSON(users)
}

func UpdateUser(c *fiber.Ctx) error {
	var user model.User
	email := c.Params("email")

	err := c.BodyParser(&user)
	if err != nil {
		return UserCustomErr(c, "Invalid input formate", fiber.ErrBadRequest)
	}

	updates := bson.M{}
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.UserType != "" {
		updates["usertype"] = user.UserType
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}

	result, er := config.UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": email},
		bson.M{"$set": updates},
	)
	if er != nil {
		return UserCustomErr(c, "Failed to update data", fiber.ErrInternalServerError)
	}

	if result.MatchedCount == 0 {
		return UserCustomErr(c, "User not found to update", fiber.ErrNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"updated items count": result.ModifiedCount,
		"message":             "User updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	var user model.User
	email := c.Params("email")

	err := config.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return UserCustomErr(c, "User Not Found", err)
	}

	res, er := config.UserCollection.DeleteOne(context.TODO(), bson.M{"email": email})
	if er != nil {
		return UserCustomErr(c, "failed to delete user", fiber.ErrInternalServerError)
	}
	if res.DeletedCount == 0 {
		return UserCustomErr(c, "User Not Found", fiber.ErrNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user deleted": res.DeletedCount,
		"deleted user": user,
		"message":      "User account deleted successfully",
	})
}
