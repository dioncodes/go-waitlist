package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/dioncodes/go-waitlist/internal/handler/waitlist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"Content-Type", "Origin", "Authorization"}

	if os.Getenv("ENV") == "dev" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		corsConfig.AllowOrigins = strings.Split(corsOrigins, ",")
	}

	r.Use(cors.New(corsConfig))

	// r.LoadHTMLGlob(os.Getenv("BASE_DIR") + "/templates/web/*.html")

	v1 := r.Group("/v1")

	waitlist.RegisterRoutes(v1)

	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"details": "system up and running",
		})
	})
}
