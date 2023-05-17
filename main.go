package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	movie struct {
		ID          int     `json:"id"`
		ImdbID      string  `json:"imdbID"`
		Title       string  `json:"title"`
		Year        int     `json:"year"`
		Rating      float64 `json:"rating"`
		IsSuperHero bool    `json:"isSuperHero"`
	}
)

var (
	movies = map[int]*movie{}
	seq    = 1
	lock   = sync.Mutex{}
)

//----------
// Handlers
//----------

func createMovie(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	m := &movie{
		ID: seq,
	}
	if err := c.Bind(m); err != nil {
		return err
	}
	movies[m.ID] = m
	seq++
	return c.JSON(http.StatusCreated, m)
}

func getMovie(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, movies[id])
}

func updateMovie(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	m := new(movie)
	if err := c.Bind(m); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	movies[id].Rating = m.Rating
	return c.JSON(http.StatusOK, movies[id])
}

func deleteMovie(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(movies, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllMovies(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, movies)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/go-health-check", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server running smooth.")
	})

	// Routes
	e.GET("/movies", getAllMovies)
	e.POST("/movies", createMovie)
	e.GET("/movies/:id", getMovie)
	e.PUT("/movies/:id", updateMovie)
	e.DELETE("/movies/:id", deleteMovie)

	e.Logger.Fatal(e.Start(":3001"))
}
