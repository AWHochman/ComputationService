package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/compute", compute)
	router.Run(":8080")
}

func compute(c *gin.Context) {

}