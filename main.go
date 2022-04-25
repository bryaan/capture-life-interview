package main

import (
	"capture-life/controllers"
	"capture-life/models"
	"github.com/gin-gonic/gin"
)
import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	// Blog Post Routes
	r.GET("/blog-posts", controllers.GetAllBlogPosts)
	r.GET("/blog-post/:id", controllers.GetBlogPost)
	r.POST("/blog-post", controllers.CreateBlogPost)
	r.PUT("/blog-post/:id", controllers.UpdateBlogPost)
	r.DELETE("/blog-post/:id", controllers.DeleteBlogPost)

	// Comment Routes
	r.GET("/comments/:blog-post-id", controllers.GetCommentsForBlogPost)
	r.POST("/comments", controllers.CreateComment)
	r.PUT("/comments/:id", controllers.UpdateComment)
	r.DELETE("/comments/:id", controllers.DeleteComment)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
