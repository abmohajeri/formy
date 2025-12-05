package routes

import (
	"core/config"
	"core/controllers/api"
	"core/controllers/telegram"
	"core/services"
	"os"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(config.RateLimitMiddleware(time.Minute, 30))

	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "views/assets")

	recaptcha.Init(os.Getenv("RECAPTCHA_SECRET_KEY"))
	router.POST("/:data", api.CreateFormData)

	services.InitTelegram()
	router.POST("/"+services.Token, telegram.TelegramWebhookHandler)

	router.GET("/captcha", services.AltchaHandler)

	return router
}
