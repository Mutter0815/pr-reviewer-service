package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/Health", func(c *gin.Context) { c.JSON(200, "OK") })
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("err")
	}
}
