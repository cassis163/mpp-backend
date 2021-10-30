package routing

import (
	"database/sql"
	"fmt"
	"mpp/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB) {
	adress := "localhost:8090"
	r := gin.Default()
	r.Use(cors.Default())
	setup(r, db)
	r.Run(adress)
	fmt.Println("Server started at " + adress)
}

func setup(r *gin.Engine, db *sql.DB) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/movies", func(c *gin.Context) { controllers.Index(c, db) })
	r.GET("/movies/:id", func(c *gin.Context) { controllers.Get(c, db) })
	r.POST("/movies", func(c *gin.Context) { controllers.Create(c, db) })
	r.GET("/summaries", func (c *gin.Context) { controllers.GetSummaries(c, db) })
}
