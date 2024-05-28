package main

import (
	"Postgres-GORM-BookStore/models"
	"Postgres-GORM-BookStore/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)

}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := Book{}

	err := ctx.BodyParser(&book)

	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil

}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	book := models.Books{}

	id := ctx.Params("id")

	if id != "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(book, id)

	if err.Error != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil

}

func (r *Repository) GetBookByID(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	book := models.Books{}

	err := r.DB.First(&book, id).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find book",
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book found",
		"data":    book,
	})
	return nil

}

func (r *Repository) GetBooks(ctx *fiber.Ctx) error {

	books := []models.Books{}

	err := r.DB.Find(&books).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find books",
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books found",
		"data":    books,
	})
	return nil

}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Could not load the database")
	}
	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("Could not migrate database")
	}

	repository := Repository{DB: db}

	app := fiber.New()
	repository.SetupRoutes(app)
	app.Listen(":8080")

}
