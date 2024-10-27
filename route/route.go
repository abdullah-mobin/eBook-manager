package route

import (
	"github.com/abdullah-mobin/ebook-manager/handler"
	"github.com/abdullah-mobin/ebook-manager/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetRoute(app *fiber.App) {
	app.Get("/home", handler.WelcomeMsg)
	SetAuthRoute(app)
	SetEBookRoute(app)
	SetUserRoute(app)
}

func SetAuthRoute(auth *fiber.App) {
	auth.Post("/signup", handler.SignUp).Name("Sign Up Route")
	auth.Post("/login", handler.Login).Name("Login Route")
	auth.Get("/auth-google", handler.Googlelogin).Name("Sign Up / login with google")
	auth.Get("/google_callback", handler.GoogleCalback).Name("google callback for oauth")
}

func SetEBookRoute(app *fiber.App) {
	home := app.Group("/book")

	home.Get("/all-books", handler.GetAllBooks).Name("All books")
	home.Get("/id/:id", handler.GetBooksById).Name("Single book by id")
	home.Get("/name", handler.GetBooksByName).Name("Single book by name as query")

	protected := app.Group("/admin", middleware.ProtectedRoute())

	protected.Post("/create-book", handler.CreateNewBook).Name("Register new book")
	protected.Put("/update-book/:id", handler.UpdateBook).Name("Update existing book")
	protected.Delete("/delete-book/:id", handler.DeleteBookByID).Name("Remove existing book by id")
}

func SetUserRoute(app *fiber.App) {
	user := app.Group("/user", middleware.ProtectedRoute())

	user.Get("/profile", handler.GetUserByEmail).Name("User profile")
	user.Get("/users", handler.GetAllUsers).Name("get all users")
	user.Get("/user-by-email", handler.GetUserByEmail).Name("get user by email")
	user.Post("/new-user", handler.CreateUser).Name("register new user")
	user.Put("/update-user/:email", handler.UpdateUser).Name("Update existing user")
	user.Delete("/delete-user/:email", handler.DeleteUser).Name("Delete existing user")
}
