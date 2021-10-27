package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mpp/config"
	"mpp/models"
	"mpp/util"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context, db *sql.DB) {
	movies := getMovies(db)
	c.JSON(200, movies)
}

func Get(c *gin.Context, db *sql.DB) {
	idParam := c.Param("id")
	row := db.QueryRow("SELECT * FROM movies WHERE id = ?;", idParam)
	movie := new(models.Movie)
	err := row.Scan(&movie.Id, &movie.Name, &movie.Year, &movie.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{
				"message": "Record was not found",
			})
		} else {
			util.CheckErr(err)
		}
	} else {
		c.JSON(200, movie)
	}
}

func Create(c *gin.Context, db *sql.DB) {
	movie := new(models.Movie)
	if err := c.BindJSON(&movie); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, movie)

	insertMovieSQL := `INSERT INTO movies(id, name, year, score) VALUES (?, ?, ?, ?);`
	statement, err := db.Prepare(insertMovieSQL) // Prepare statement.
	// This is good to avoid SQL injections
	util.CheckErr(err)
	_, err = statement.Exec(movie.Id, movie.Name, movie.Year, movie.Score)
	util.CheckErr(err)
	c.JSON(200, gin.H{
		"message": "Succesfully inserted the movie",
	})
}

func GetSummaries(c *gin.Context, db *sql.DB) {
	movies := getMovies(db)
	var wg sync.WaitGroup

	for _, movie := range(movies) {
		// TODO: Add worker pool
		wg.Add(1)
		go updateSummary(movie, db, &wg)
	}

	wg.Wait()
	c.JSON(200, gin.H{
		"message": "Successfully generated summaries",
	})
}

func updateSummary(movie models.Movie, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	summary := getSummary(movie.Id, movie.Name)
	setSummary(movie.Id, movie.Name, *summary, db)
}

func getSummary(id string, title string) *string {
	fmt.Println("Getting summary for " + title)

	res, err := http.Get("http://www.omdbapi.com/?apikey=" + config.API_KEY + "&i=" + id + "&plot=full&r=json")
	util.CheckErr(err)

	fmt.Println("Received summary for " + title)

	body, err := ioutil.ReadAll(res.Body)
	util.CheckErr(err)
	defer res.Body.Close()

	omdbItem := models.OmdbItem{}
	if res.StatusCode == http.StatusOK {
		err := json.Unmarshal(body, &omdbItem)
		util.CheckErr(err)
	}

	return &omdbItem.Plot
}

func setSummary(id string, name string, summary string, db *sql.DB) {
	_, err := db.Exec("UPDATE movies SET summary = ? WHERE id = ?;", summary, id)
	util.CheckErr(err)

	fmt.Println("Written summary for " + name)
}

func getMovies(db *sql.DB) []models.Movie {
	rows, err := db.Query("SELECT * FROM movies;")
	movies := []models.Movie{}
	util.CheckErr(err)
	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(&movie.Id, &movie.Name, &movie.Year, &movie.Score, &movie.Summary)
		util.CheckErr(err)
		movies = append(movies, movie)
	}
	rows.Close()

	return movies
}
