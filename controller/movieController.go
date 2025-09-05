package controller

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/favmovies/initializers"
	"github.com/scientist-v08/favmovies/model"
)

func requiredFeild(feild string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"Error": feild + "is a required feild"})
}

func Contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func PostMovies(c *gin.Context) {
	// Parse multipart form (max 2MB for form data)
		if err := c.Request.ParseMultipartForm(2 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
			return
		}

		// Define required fields
		requiredFields := []string{"Title", "Type", "Director", "Budget", "Location", "Duration", "Year of Release"}
		values := make(map[string]string)

		// Validate each field
		for _, field := range requiredFields {
			if valuesSlice, ok := c.Request.MultipartForm.Value[field]; !ok || len(valuesSlice) == 0 || strings.TrimSpace(valuesSlice[0]) == "" {
				requiredFeild(field, c)
				return
			} else {
				values[field] = strings.TrimSpace(valuesSlice[0])
			}
		}

		// Get the uploaded file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image: " + err.Error()})
			return
		}

		// Validate file size (e.g., max 1MB)
		if file.Size > 1*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image size exceeds 5MB"})
			return
		}

		// Validate file type (e.g., JPEG or PNG)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG and PNG images are allowed"})
			return
		}

		// Open and read the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image: " + err.Error()})
			return
		}
		defer src.Close()

		imageBytes, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to read image: " + err.Error()})
			return
		}

		roles, _ := c.Get("roles")
		userRoles := roles.([]string)
		isAdmin := false
		if Contains(userRoles, "ROLE_ADMIN") {
        	isAdmin = true
    	}

		var mimeType string

		switch ext {
			case ".jpg", ".jpeg":
				mimeType = "image/jpeg"
			case ".png":
				mimeType = "image/png"
			default:
				mimeType = "application/octet-stream" // fallback
		}

		movie := model.Movies{
			Title:         	 values["Title"],
			Type:          	 values["Type"],
			Director:      	 values["Director"],
			Budget:        	 values["Budget"],
			Location:      	 values["Location"],
			Duration:      	 values["Duration"],
			YearOfRelease: 	 values["Year of Release"],
			Image:           imageBytes,
			IsAdminApproved: isAdmin,
			MimeType:		 mimeType,
		}

		if err := initializers.DB.Create(&movie).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Movie Post created successfully"})
}

func GetMoviesPaginated(c *gin.Context) {
	// Pagination
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	offset := (pageNumber - 1) * pageSize

	// Sorting
	sortBy := c.DefaultQuery("sortBy", "created_at") // default sort column
	order := c.DefaultQuery("order", "desc")         // asc or desc

	// Search filters
	searchTitle := c.Query("title")
	searchDirector := c.Query("director")

	// Base query (only approved movies)
	query := initializers.DB.Model(&model.Movies{}).Where("is_admin_approved = ?", true)

	// Apply search filters
	if searchTitle != "" {
		query = query.Where("title ILIKE ?", "%"+searchTitle+"%")
	}
	if searchDirector != "" {
		query = query.Where("director ILIKE ?", "%"+searchDirector+"%")
	}

	// Count total (with filters)
	var count int64
	query.Count(&count)

	// Get paginated + sorted results
	var movies []model.Movies
	if err := query.
		Order(fmt.Sprintf("%s %s", sortBy, order)).
		Limit(pageSize).
		Offset(offset).
		Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"movies":        movies,
		"totalElements": count,
		"pageNumber":    pageNumber,
		"pageSize":      pageSize,
	})
}

func GetMoviesToBeApproved(c *gin.Context) {
	var movies []model.Movies
	query := initializers.DB.Model(&model.Movies{}).Where("is_admin_approved = ?", false)
	if err := query.
		Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Response
	c.JSON(http.StatusOK, gin.H{
		"movies":        movies,
	})
}

func DeleteOrRejectMovie(c *gin.Context) {
	// Now update that particular row in the DB
	var obtainedMovie struct {
		ID              int64 `json:"id"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		Director        string `json:"director"`
		Budget          string `json:"budget"`
		Location        string `json:"location"`
		Duration        string `json:"duration"`
		YearOfRelease   string `json:"year_of_release"`
		IsAdminApproved bool   `json:"is_admin_approved"`
	}
	c.Bind(&obtainedMovie)

	// Obtain the id
	id := obtainedMovie.ID

	// Delete the post
	result := initializers.DB.Where("id = ?", id).Delete(&model.Movies{})

	// Check if it was deleted
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Unable to delete",
		})
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"Success": "Movie deleted/rejected",
	})
}

func ApproveOrEditMovies(c *gin.Context) {
	// Now update that particular row in the DB
	var obtainedMovie struct {
		ID              int64 `json:"id"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		Director        string `json:"director"`
		Budget          string `json:"budget"`
		Location        string `json:"location"`
		Duration        string `json:"duration"`
		YearOfRelease   string `json:"year_of_release"`
		IsAdminApproved bool   `json:"is_admin_approved"`
	}
	c.Bind(&obtainedMovie)

	// Obtain the id
	id := obtainedMovie.ID

	// Obtain the data of that particular ID from the DB
	var movie model.Movies
	result := initializers.DB.First(&movie, id)

	// Now check if it exists in the DB
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "This post doesn't exist in the DB",
		})
		return
	}

	movie.Title = obtainedMovie.Title
	movie.Type = obtainedMovie.Type
	movie.Director = obtainedMovie.Director
	movie.Budget = obtainedMovie.Budget
	movie.Location = obtainedMovie.Location
	movie.Duration = obtainedMovie.Duration
	movie.YearOfRelease = obtainedMovie.YearOfRelease
	movie.IsAdminApproved = obtainedMovie.IsAdminApproved

	// Save changes
	if err := initializers.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	// Respond with success message
	c.JSON(200, gin.H{
		"Success": "Updated/Approved",
	})
}
