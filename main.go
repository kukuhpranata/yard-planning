package main

import (
	"log"
	"strings"
	"yard-planning/app/controller"
	"yard-planning/app/repository"
	"yard-planning/app/service"
	"yard-planning/database"
	"yard-planning/helper/token"
	"yard-planning/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {

	db, err := database.NewPostgresClient()
	if err != nil {
		panic(err)
	}

	validate := validator.New()

	// Initialize repositories
	userRepository := repository.NewUserRepository()
	yardRepository := repository.NewYardRepository()
	yardPlanRepository := repository.NewYardPlanRepository()
	containerPositionRepository := repository.NewContainerPositionRepository()

	// Initialize services
	userService := service.NewUserService(userRepository, db, validate)
	containerService := service.NewContainerService(yardRepository, yardPlanRepository, containerPositionRepository, db, validate)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	containerController := controller.NewContainerController(containerService)

	router := gin.Default()

	api := router.Group("/api")
	{

		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		api.POST("/suggestion", containerController.SuggestPosition)
		api.POST("/placement", containerController.PlaceContainer)
		api.POST("/pickup", containerController.PickupContainer)

		auth := api.Group("/auth")
		auth.Use(CheckAuth())
		{
		}
	}

	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedError("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateJwtToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedError(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Next()
	}
}
