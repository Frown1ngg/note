package main

import (
	"notes_project/server"
)

func init() {
	server.InitServer()
}
func main() {
	server.StartServer()
	// router := gin.Default()
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "lol",
	// 	})
	// })
	// router.Run(":8080")
}
