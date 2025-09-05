package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/favmovies/controller"
	"github.com/scientist-v08/favmovies/initializers"
	"github.com/scientist-v08/favmovies/middleware"
)

var router *gin.Engine

func init() {
	// Load env + DB (just like in main.go)
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()

	// Setup Gin router
	router = gin.Default()
	router.Use(gin.Recovery())

	// CORS
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Routes
	router.POST("/create/user", controller.SignUp)
	router.POST("/create/admin", controller.AdminSignUp)
	router.POST("/login/user", controller.Login)
	router.POST("/create/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.PostMovies)
	router.GET("/get/allMovies", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.GetMoviesPaginated)
	router.POST("/update/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.ApproveOrEditMovies)
	router.POST("/update/moviereject", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.DeleteOrRejectMovie)
	router.GET("/get/moviesToBeApproved", middleware.RequireAnyRole("ROLE_ADMIN"), controller.GetMoviesToBeApproved)
}

// Vercel entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
