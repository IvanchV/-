package main

import (
	"github.com/gin-gonic/gin"
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

func GetMovies(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, ListOfMovies)
}

func PostMovies(ctx *gin.Context) {
	var newMovie movie

	if err := ctx.BindJSON(&newMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка добавления"})
		return
	}

	ListOfMovies = append(ListOfMovies, newMovie)
	ctx.IndentedJSON(http.StatusCreated, newMovie)
}

func GetMoviesByNumber(ctx *gin.Context) {
	number := ctx.Param("number")

	for _, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "movies not found"})
}

func Sleep(ctx *gin.Context) {
	time.Sleep(10 * time.Second)
	ctx.JSON(http.StatusOK, gin.H{"message": "Задержка в 10 секунд завершена"})
}

func DeleteByNumber(ctx *gin.Context) {
	number := ctx.Param("number")

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			ListOfMovies = append(ListOfMovies[:i], ListOfMovies[i+1:]...)
			ctx.IndentedJSON(http.StatusNoContent, a)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not delete"})
}

func UpdateFilm(ctx *gin.Context) {
	number := ctx.Param("number")

	for i, a := range ListOfMovies {
		if strconv.Itoa(a.Number) == number {
			err := ctx.BindJSON(&a)
			if err != nil {
				return
			}
			ListOfMovies[i] = a
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "movies not found"})
}

func GetGif(ctx *gin.Context) {
	fileName := "gif.gif"
	ctx.File(fileName)
}

/*func PostGif(ctx *gin.Context) {
	imageBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile("example.gif", imageBytes, 0644)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
	})
}*/

func PostGif1(c *gin.Context) {
	file, err := c.FormFile("gif")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload gif"})
		return
	}

	err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save gif"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gif uploaded successfully"})
}

func main() {
	server := gin.Default()
	server.GET("/movies", GetMovies)
	server.POST("/new_movie", PostMovies)
	server.GET("/movie/:number", GetMoviesByNumber)
	server.DELETE("/delete/:number", DeleteByNumber)
	server.PUT("update/:number", UpdateFilm)
	server.GET("/sleep", Sleep)
	server.GET("/gif", GetGif)
	//server.POST("/new_gif", PostGif)
	server.POST("/new_gif", PostGif1)
	err := server.Run(":8080")
	if err != nil {
		return
	}

}
