package controller

import (
	"net/http"
	"yard-planning/app/service"
	"yard-planning/app/web"
	"yard-planning/response"

	"github.com/gin-gonic/gin"
)

type ContainerController interface {
	SuggestPosition(ctx *gin.Context)
	PlaceContainer(ctx *gin.Context)
	PickupContainer(ctx *gin.Context)
}

type ContainerControllerImpl struct {
	ContainerService service.ContainerService
}

func NewContainerController(containerService service.ContainerService) ContainerController {
	return &ContainerControllerImpl{
		ContainerService: containerService,
	}
}

func (c *ContainerControllerImpl) SuggestPosition(ctx *gin.Context) {
	request := new(web.ContainerRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		customErr := response.BadRequestError("Invalid request body or missing required fields.")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	positionResponse, customErr := c.ContainerService.SuggestPosition(ctx.Request.Context(), request)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	finalResponse := web.SuggestedPositionResponse{
		SuggestedPosition: web.PositionResponse{
			Block: positionResponse.Block,
			Slot:  positionResponse.Slot,
			Row:   positionResponse.Row,
			Tier:  positionResponse.Tier,
		},
	}

	webResponse := response.WebResponse{
		Status:  true,
		Message: "Suggested position successfully retrieved.",
		Data:    finalResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (c *ContainerControllerImpl) PlaceContainer(ctx *gin.Context) {
	request := new(web.PlacementRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		customErr := response.BadRequestError("Invalid request body or missing required fields.")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	_, customErr := c.ContainerService.PlaceContainer(ctx.Request.Context(), request)

	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}
	webResponse := response.WebResponse{
		Status:  true,
		Message: "Success",
	}

	ctx.JSON(http.StatusCreated, webResponse)
}

func (c *ContainerControllerImpl) PickupContainer(ctx *gin.Context) {
	request := new(web.PickupRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		customErr := response.BadRequestError("Invalid request body or missing required fields.")
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	_, customErr := c.ContainerService.PickupContainer(ctx.Request.Context(), request)
	if customErr != nil {
		ctx.JSON(customErr.StatusCode, customErr)
		return
	}

	webResponse := response.WebResponse{
		Status:  true,
		Message: "Success",
	}

	ctx.JSON(http.StatusOK, webResponse)
}
