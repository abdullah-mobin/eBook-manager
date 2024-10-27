package handler

import (
	"context"

	"github.com/abdullah-mobin/ebook-manager/config"
	"github.com/abdullah-mobin/ebook-manager/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func WelcomeMsg(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to ebook manager",
	})
}

func GetAllBooks(c *fiber.Ctx) error {

	var books []model.Book

	cursor, err := config.MongoCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return CustomError(c, "Failed to query books", err, 500)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var b model.Book
		err := cursor.Decode(&b)
		if err != nil {
			return CustomError(c, "Error decoding collection", err, 500)
		}
		books = append(books, b)
	}

	err = cursor.Err()
	if err != nil {
		return CustomError(c, "Error iterating collection", err, 500)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"books": books,
	})
}

func GetBooksById(c *fiber.Ctx) error {
	var book model.Book

	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return CustomError(c, "incorrect id formate", err, 400)
	}

	err = config.MongoCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return CustomError(c, "Book not found", err, 404)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"book": book,
	})
}

func GetBooksByName(c *fiber.Ctx) error {
	var book model.Book
	name := c.Query("name")

	err := config.MongoCollection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return CustomError(c, "Book not found", err, 404)
		}
		return CustomError(c, "Error occurred while searching", err, 500)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"book": book,
	})
}

func CreateNewBook(c *fiber.Ctx) error {
	var book model.Book
	err := c.BodyParser(&book)
	if err != nil {
		return CustomError(c, "Invalid input formate", err, 400)
	}
	if book.Name == "" || book.Author == "" || book.Type == "" || book.PDF == "" {
		return CustomError(c, "All fields required", err, 400)
	}
	result, er := config.MongoCollection.InsertOne(context.TODO(), book)
	if er != nil {
		return CustomError(c, "Failed to create new book", er, 500)
	}

	book.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ID":      book.ID,
		"message": "book created successfully",
	})
}

func UpdateBook(c *fiber.Ctx) error {
	var book model.Book
	id := c.Params("id")

	obID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return CustomError(c, "Invalid id formate", err, 400)
	}
	err = c.BodyParser(&book)
	if err != nil {
		return CustomError(c, "Invalid input formate", err, 400)
	}

	updates := bson.M{}
	if book.Name != "" {
		updates["name"] = book.Name
	}
	if book.Author != "" {
		updates["author"] = book.Author
	}
	if book.Type != "" {
		updates["type"] = book.Type
	}
	if book.PDF != "" {
		updates["pdf"] = book.PDF
	}

	result, er := config.MongoCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": obID},
		bson.M{"$set": updates},
	)
	if er != nil {
		return CustomError(c, "Failed to update data", er, 500)
	}

	if result.MatchedCount == 0 {
		return CustomError(c, "Book not found to update", nil, 404)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"updated items count": result.ModifiedCount,
		"message":             "Book updated successfully",
	})
}

func DeleteBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book model.Book

	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return CustomError(c, "Invalid id formate", err, 400)
	}

	err = config.MongoCollection.FindOne(context.TODO(), bson.M{"_id": obId}).Decode(&book)
	if err != nil {
		return CustomError(c, "Book not found", err, 404)
	}

	result, er := config.MongoCollection.DeleteOne(context.TODO(), bson.M{"_id": obId})
	if er != nil {
		return CustomError(c, "Failed to delete book", er, 500)
	}
	if result.DeletedCount == 0 {
		return CustomError(c, "Book Not Found", fiber.ErrNotFound, 404)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"book deleted": result.DeletedCount,
		"deleted book": book,
		"message":      "Book Deleted Successfully",
	})
}

func CustomError(c *fiber.Ctx, msg string, err error, code int) error {
	return c.Status(code).JSON(fiber.Map{
		"error":   err,
		"message": msg,
	})
}
