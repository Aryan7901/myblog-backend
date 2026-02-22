package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func BlogRoutes(router *gin.Engine) {
	router.GET("/blogs/all", controllers.GetAllBlogs)
	router.GET("/blogs/blog/:bid", controllers.GetBlogById)

	authorized:=router.Group("")

	authorized.Use(middleware.CheckAuth())
	authorized.POST("/blogs/comment/:bid", controllers.MakeComment)
	authorized.PATCH("/blogs/comment/:cid", controllers.UpdateComment)
	authorized.DELETE("/blogs/comment/:cid", controllers.DeleteComment)
}
