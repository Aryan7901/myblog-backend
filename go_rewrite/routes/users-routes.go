package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/user/signup", controllers.Signup)
	router.POST("/user/login", controllers.Login)
	
	authorized:=router.Group("")
	authorized.Use(middleware.CheckAuth())

	authorized.GET("/user/list", controllers.GetUserBlogs)
	authorized.POST("/user/new-blog", controllers.CreateBlog)
	authorized.PATCH("/user/:bid", controllers.UpdateBlog)
	authorized.DELETE("/user/:bid", controllers.DeleteBlog)
}
