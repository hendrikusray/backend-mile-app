package main

import (
	"fmt"
	"net/http"
	"os"

	cnf "mile-app-test/config"
	"mile-app-test/middleware"

	taskHTTP "mile-app-test/internal/task/delivery/http"
	taskRepo "mile-app-test/internal/task/repository"
	taskUsecase "mile-app-test/internal/task/usecase"

	userHTTP "mile-app-test/internal/user/delivery/http"
	user "mile-app-test/internal/user/repository"
	userUsecase "mile-app-test/internal/user/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Init application...")
}

func main() {
	fmt.Println("Setup application...")
	if err := cnf.StartConfig(); err != nil {
		panic(err)
	}
	setupRouter()
}

func setupRouter() {
	gin.SetMode(gin.ReleaseMode)

	client := cnf.Connect()
	cnf.EnsureMongoIndexes(client)

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://<your-vercel-project>.vercel.app",
			"https://<custom-domain-kamu>",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/v1")

	// USER
	userRepo := user.NewUserRepository(client)
	userUseCase := userUsecase.NewUserUseCase(userRepo, 1000)
	userHTTP.NewUserHandler(v1, userUseCase)

	// TASK (dengan middleware auth)
	tr := taskRepo.NewTaskRepository(client)
	tu := taskUsecase.NewTaskUseCase(tr)
	taskHTTP.RegisterTaskRoutes(v1, tu, middleware.Auth())

	// Health
	addHealth(v1)
	v2 := r.Group("/v2")
	addHealth(v2)

	// Listen ke PORT dari env (Railway/Render mengisi ini)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}

func addHealth(rg *gin.RouterGroup) {
	h := rg.Group("/health")
	h.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	h.GET("/db-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
}
