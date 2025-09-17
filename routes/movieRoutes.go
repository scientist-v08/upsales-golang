package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/favmovies/controller"
	"github.com/scientist-v08/favmovies/middleware"
)

func RegisterMovieRoutes(r *gin.Engine) {
	r.POST("/create/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.PostMovies)
	r.GET("/get/allMovies", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.GetMoviesPaginated)
	r.POST("/update/movie", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.ApproveOrEditMovies)
	r.POST("/update/moviereject", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.DeleteOrRejectMovie)
	r.GET("/get/moviesToBeApproved", middleware.RequireAnyRole("ROLE_ADMIN"), controller.GetMoviesToBeApproved)

	// Grouping movie-related routes
	// movies := r.Group("/movies")
	// {
	// 	movies.POST("/create", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.PostMovies)
	// 	movies.GET("/all", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.GetMoviesPaginated)
	// 	movies.POST("/update", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.ApproveOrEditMovies)
	// 	movies.POST("/reject", middleware.RequireAnyRole("ROLE_USER", "ROLE_ADMIN"), controller.DeleteOrRejectMovie)
	// 	movies.GET("/toBeApproved", middleware.RequireAnyRole("ROLE_ADMIN"), controller.GetMoviesToBeApproved)
	// }
}