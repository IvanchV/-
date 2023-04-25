package main

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

type movie struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Cast   string `json:"cast"`
	Year   int    `json:"year"`
}

var ListOfMovies []movie

func GetMovies(ctx *fiber.Ctx) error {
	return ctx.JSON(ListOfMovies)
}

func PostMovies(ctx *fiber.Ctx) error {
	var newMovie movie

	if err := ctx.BodyParser(&newMovie); err != nil {
		err := ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Ошибка добавления"})
		if err != nil {
			return err
		}
		return err
	}

	ListOfMovies = append(ListOfMovies, newMovie)
	return ctx.Status(http.StatusCreated).JSON(newMovie)
}

func GetMoviesByNumber(ctx *fiber.Ctx) error {
	number := ctx.Params("number")

	for _, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			return ctx.JSON(a)
		}
	}
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "movies not found"})
}

func Sleep(ctx *fiber.Ctx) error {
	time.Sleep(10 * time.Second)
	return ctx.JSON(fiber.Map{"message": "Задержка в 10 секунд завершена"})
}

func DeleteByNumber(ctx *fiber.Ctx) error {
	number := ctx.Params("number")

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			ListOfMovies = append(ListOfMovies[:i], ListOfMovies[i+1:]...)
			return ctx.Status(http.StatusNoContent).JSON(a)
		}
	}
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "movie not delete"})
}

func UpdateFilm(ctx *fiber.Ctx) error {
	number := ctx.Params("number")

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			err := ctx.BodyParser(&a)
			if err != nil {
				return err
			}
			ListOfMovies[i] = a
			return ctx.JSON(a)
		}
	}
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": "movies not found"})
}

func GetGif(ctx *fiber.Ctx) error {
	fileName := "gif.gif"
	return ctx.SendFile(fileName)
}

func PostGif(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("gif")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to upload gif",
		})
	}

	err = ctx.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save gif",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Gif uploaded successfully",
	})
}

func main() {
	server := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})
	server.Get("/movies", GetMovies)
	server.Post("/new_movie", PostMovies)
	server.Get("/movie/:number", GetMoviesByNumber)
	server.Delete("/delete/:number", DeleteByNumber)
	server.Put("/update/:number", UpdateFilm)
	server.Get("/sleep", Sleep)
	server.Get("/gif", GetGif)
	server.Post("/new_gif", PostGif)
	err := server.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
