package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/favmovies/controller"
	"github.com/scientist-v08/favmovies/initializers"
	"github.com/scientist-v08/favmovies/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{allowedOrigin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/create/user", controller.SignUp)
	r.POST("/create/admin", controller.AdminSignUp)
	r.POST("/login/user", controller.Login)
	r.POST("/create/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.PostMovies)
	r.GET("/get/allMovies", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.GetMoviesPaginated)
	r.POST("/update/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.ApproveOrEditMovies)
	r.POST("/update/moviereject", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.DeleteOrRejectMovie)
	r.GET("/get/moviesToBeApproved", middleware.RequireAnyRole("ROLE_ADMIN"), controller.GetMoviesToBeApproved)
	r.Run()
}
