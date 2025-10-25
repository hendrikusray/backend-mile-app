package main

import (
	"fmt"
	cnf "mile-app-test/config"
	"mile-app-test/middleware"
	"net/http"

	taskHTTP "mile-app-test/internal/task/delivery/http"
	taskRepo "mile-app-test/internal/task/repository"
	taskUsecase "mile-app-test/internal/task/usecase"
	userHTTP "mile-app-test/internal/user/delivery/http"
	user "mile-app-test/internal/user/repository"
	userUsecase "mile-app-test/internal/user/usecase"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Init application...")

}

func main() {

	fmt.Println("Setup application...")
	err := cnf.StartConfig()
	if err != nil {
		panic(err)
	}
	setupRouter()
}

func setupRouter() {

	client := cnf.Connect()
	r := gin.Default()
	v1 := r.Group("/v1")

	userRepo := user.NewUserRepository(client)
	userUseCase := userUsecase.NewUserUseCase(userRepo, 1000)
	tr := taskRepo.NewTaskRepository(client)
	tu := taskUsecase.NewTaskUseCase(tr)

	taskHTTP.RegisterTaskRoutes(v1, tu, middleware.Auth())
	userHTTP.NewUserHandler(v1, userUseCase)

	addHealth(v1)

	v2 := r.Group("/v2")
	addHealth(v2)
	err := r.Run()

	panic(err)
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
