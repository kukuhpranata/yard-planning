package web

type ContainerRequest struct {
	YardName        string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`

	Size   string `json:"container_size" validate:"required,oneof=20ft 40ft"`
	Height string `json:"container_height" validate:"required,oneof=8.6ft 9.6ft"`
	Type   string `json:"container_type" validate:"required"`

	// Optional
	BlockName string `json:"block"`
	Slot      int    `json:"slot"`
	Row       int    `json:"row"`
	Tier      int    `json:"tier"`
}

type PlacementRequest struct {
	YardName        string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`

	BlockName string `json:"block" validate:"required"`
	Slot      int    `json:"slot" validate:"required,min=1"`
	Row       int    `json:"row" validate:"required,min=1"`
	Tier      int    `json:"tier" validate:"required,min=1"`

	// needs for yard_plan check
	Size   string `json:"container_size" validate:"required,oneof=20ft 40ft"`
	Height string `json:"container_height" validate:"required,oneof=8.6ft 9.6ft"`
	Type   string `json:"container_type" validate:"required"`
}

type PickupRequest struct {
	YardName        string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`
}

type PositionRequest struct {
	BlockID    int  `json:"-"`
	Slot       int  `json:"slot" validate:"required,min=1"`
	Row        int  `json:"row" validate:"required,min=1"`
	Tier       int  `json:"tier" validate:"required,min=1"`
	YardPlanID *int `json:"-"`
}

type PositionResponse struct {
	Block string `json:"block"`
	Slot  int    `json:"slot"`
	Row   int    `json:"row"`
	Tier  int    `json:"tier"`

	BlockID    int  `json:"-"`
	YardPlanID *int `json:"-"`
}

type SuggestedPositionResponse struct {
	SuggestedPosition PositionResponse `json:"suggested_position"`
}

type GeneralResponse struct {
	Message string `json:"message"`
}
