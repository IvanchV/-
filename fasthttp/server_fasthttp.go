package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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

func GetMovies(ctx *fasthttp.RequestCtx) {
	response, err := json.Marshal(ListOfMovies)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(response)
}

func GetMoviesByNumber(ctx *fasthttp.RequestCtx) {
	number := ctx.UserValue("number").(string)
	for _, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			response, err := json.Marshal(a)
			if err != nil {
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			_, _ = ctx.Write(response)
			return
		}
	}
	ctx.Error("Movies Not Found", fasthttp.StatusNotFound)
}

func PostMovie(ctx *fasthttp.RequestCtx) {
	var newMovie movie
	if err := json.Unmarshal(ctx.PostBody(), &newMovie); err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = ctx.Write([]byte(`{"error": "Ошибка добавления"}`))
		return
	}

	ListOfMovies = append(ListOfMovies, newMovie)

	response, err := json.Marshal(newMovie)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = ctx.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	_, _ = ctx.Write(response)
}

func PutMovie(ctx *fasthttp.RequestCtx) {
	number := ctx.UserValue("number").(string)
	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			if err := json.Unmarshal(ctx.PostBody(), &a); err != nil {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				_, _ = ctx.Write([]byte(`{"error": "Ошибка обновления"}`))
				return
			}
			ListOfMovies[i] = a
			response, err := json.Marshal(a)
			if err != nil {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				_, _ = ctx.Write([]byte(`{"error": "Internal Server Error"}`))
				return
			}
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			_, _ = ctx.Write(response)
			return
		}
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	_, _ = ctx.Write([]byte(`{"message": "movies not found"}`))
}

func DeleteMovie(ctx *fasthttp.RequestCtx) {
	number := ctx.UserValue("number").(string)
	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			ListOfMovies = append(ListOfMovies[:i], ListOfMovies[i+1:]...)
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	_, _ = ctx.Write([]byte(`{"message": "movie not delete"}`))
}

func GetGif(ctx *fasthttp.RequestCtx) {
	file := "gif.gif"
	ctx.SendFile(file)
}

func PostGif(ctx *fasthttp.RequestCtx) {
	file, err := ctx.FormFile("gif")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBody([]byte("Failed to upload gif"))
		return
	}
	err = fasthttp.SaveMultipartFile(file, "./uploads/"+file.Filename)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody([]byte(`{"message": "Failed to save gif"}`))
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody([]byte(`{"message": "Gif uploaded successfully"}`))
}

func Sleep(ctx *fasthttp.RequestCtx) {
	time.Sleep(10 * time.Second)
	ctx.SetContentType("application/json")
	_, _ = ctx.Write([]byte(`{"message": "Задержка в 10 секунд завершена"}`))
}

func main() {
	server := fasthttprouter.New()
	server.GET("/movies", GetMovies)
	server.POST("/new_movie", PostMovie)
	server.GET("/movie/:number", GetMoviesByNumber)
	server.PUT("/put/:number", PutMovie)
	server.DELETE("/delete/:number", DeleteMovie)
	server.GET("/gif", GetGif)
	server.POST("/new_gif", PostGif)
	server.GET("/sleep", Sleep)

	err := fasthttp.ListenAndServe(":8080", server.Handler)
	if err != nil {
		return
	}
}
