package controller

import (
	"net/http"
	"yard-planning/app/service"
	"yard-planning/app/web"
	"yard-planning/response"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (c *UserControllerImpl) Register(ctx *gin.Context) {
	userCreateRequest := new(web.Register)
	if err := ctx.ShouldBindJSON(userCreateRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	userResponse, customErr := c.UserService.Register(ctx.Request.Context(), userCreateRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Status:  true,
		Message: "Register successfully!",
		Data:    userResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *UserControllerImpl) Login(ctx *gin.Context) {
	loginRequest := new(web.LoginUserRequest)
	if err := ctx.ShouldBindJSON(loginRequest); err != nil {
		customErr := response.BadRequestError("Invalid request body")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	loginResponse, customErr := c.UserService.Login(ctx.Request.Context(), loginRequest)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	webResponse := response.WebResponse{
		Status:  true,
		Message: "Login successfully!",
		Data:    loginResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
