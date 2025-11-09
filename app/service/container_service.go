package service

import (
	"context"
	"strconv"
	"time"
	"yard-planning/app/model"
	"yard-planning/app/repository"
	"yard-planning/app/web"
	"yard-planning/response"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ContainerService interface {
	SuggestPosition(ctx context.Context, request *web.ContainerRequest) (*web.PositionResponse, *response.CustomError)

	PlaceContainer(ctx context.Context, request *web.PlacementRequest) (*web.PositionResponse, *response.CustomError)

	PickupContainer(ctx context.Context, request *web.PickupRequest) (*web.GeneralResponse, *response.CustomError)
}

type ContainerServiceImpl struct {
	YardRepository              repository.YardRepository
	YardPlanRepository          repository.YardPlanRepository
	ContainerPositionRepository repository.ContainerPositionRepository
	DB                          *gorm.DB
	Validate                    *validator.Validate
}

func NewContainerService(
	yardRepo repository.YardRepository,
	planRepo repository.YardPlanRepository,
	containerRepo repository.ContainerPositionRepository,
	DB *gorm.DB,
	validate *validator.Validate,
) ContainerService {
	return &ContainerServiceImpl{
		YardRepository:              yardRepo,
		YardPlanRepository:          planRepo,
		ContainerPositionRepository: containerRepo,
		DB:                          DB,
		Validate:                    validate,
	}
}

func (s *ContainerServiceImpl) SuggestPosition(ctx context.Context, request *web.ContainerRequest) (*web.PositionResponse, *response.CustomError) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	if request.Size == "" || request.Height == "" || request.Type == "" {
		return nil, response.BadRequestError("Container specifications (Size, Height, Type) are required for suggestion.")
	}

	// check container if exist
	var existingPosition model.ContainerPosition
	err := s.ContainerPositionRepository.FindByContainerNumber(s.DB, &existingPosition, request.ContainerNumber)

	if err == nil && existingPosition.ID > 0 {
		blockName := ""
		if existingPosition.BlockID > 0 {
			var block model.Block
			s.YardRepository.FindBlockByID(s.DB, &block, existingPosition.BlockID)
			blockName = block.Name
		}

		return nil, response.GeneralError(
			"Container " + request.ContainerNumber + " is already placed at position: " +
				blockName + " S" + strconv.Itoa(existingPosition.SlotNumber) +
				" R" + strconv.Itoa(existingPosition.RowNumber) + " T" + strconv.Itoa(existingPosition.TierNumber),
		)
	}

	// Find yard
	var yard model.Yard
	if err := s.YardRepository.FindYardByName(s.DB, &yard, request.YardName); err != nil {
		return nil, response.NotFoundError("Yard not found.")
	}

	var blocks []model.Block

	// Find blocks
	if request.BlockName != "" {
		var block model.Block
		err = s.YardRepository.FindBlockByNameAndYardID(s.DB, &block, request.BlockName, yard.ID)
		if err != nil {
			return nil, response.NotFoundError("Block not found in this Yard.")
		}
		blocks = append(blocks, block)
	} else {
		err = s.YardRepository.FindBlocksByYardID(s.DB, &blocks, yard.ID)
		if err != nil {
			return nil, response.GeneralError("Failed to fetch blocks for the yard: " + err.Error())
		}
	}

	if len(blocks) == 0 {
		return nil, response.GeneralError("No blocks found in the yard to search.")
	}

	// Iterate blocks
	for _, block := range blocks {

		var activePlans []model.YardPlan
		if err := s.YardPlanRepository.FindActivePlansByBlock(s.DB, &activePlans, block.ID); err != nil {
			continue
		}

		if len(activePlans) == 0 {
			continue
		}

		// Iterate plans
		for _, plan := range activePlans {
			// check container match
			if plan.ContainerSize != request.Size || plan.ContainerHeight != request.Height || plan.ContainerType != request.Type {
				continue
			}

			//Iterate tiers
			for t := 1; t <= block.Tiers; t++ {
				//Iterate rows
				for r := plan.RowStart; r <= plan.RowEnd; r++ {
					//Iterate slots
					for slotNum := plan.SlotStart; slotNum <= plan.SlotEnd; slotNum++ {

						slotNumbersToCheck := []int{slotNum}

						if request.Size == "40ft" {
							if (slotNum-plan.SlotStart+1)%2 != 1 || (slotNum+1) > plan.SlotEnd {
								continue
							}
							slotNumbersToCheck = append(slotNumbersToCheck, slotNum+1)
						}

						// Check available container, 0 = empty
						count, err := s.ContainerPositionRepository.CheckPositionAvailability(
							s.DB,
							block.ID,
							r,
							t,
							slotNumbersToCheck,
						)

						if err != nil {
							return nil, response.GeneralError("Database check failed.")
						}

						if count == 0 {
							return &web.PositionResponse{
								Block:      block.Name,
								Slot:       slotNum,
								Row:        r,
								Tier:       t,
								YardPlanID: &plan.ID,
							}, nil
						}
					}
				}
			}
		}
	}

	return nil, response.GeneralError("No empty position found matching the active yard plans criteria.")
}

func (s *ContainerServiceImpl) PlaceContainer(ctx context.Context, request *web.PlacementRequest) (*web.PositionResponse, *response.CustomError) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	// Check if container is exist
	var existingPosition model.ContainerPosition
	err := s.ContainerPositionRepository.FindByContainerNumber(s.DB, &existingPosition, request.ContainerNumber)

	if err == nil && existingPosition.ID > 0 {
		return nil, response.GeneralError("Container number " + request.ContainerNumber + " is already placed in the Yard.")
	}

	// Check yard 3. Cari Yard, Block, dan Cek Batasan Block
	var yard model.Yard
	if err := s.YardRepository.FindYardByName(s.DB, &yard, request.YardName); err != nil {
		return nil, response.NotFoundError("Yard not found.")
	}

	//check block
	var block model.Block
	if err := s.YardRepository.FindBlockByNameAndYardID(s.DB, &block, request.BlockName, yard.ID); err != nil {
		return nil, response.NotFoundError("Block not found in the specified Yard.")
	}

	slot := request.Slot
	row := request.Row
	tier := request.Tier
	size := request.Size
	is40ft := size == "40ft"

	if slot > block.Slots || row > block.Rows || tier > block.Tiers || slot < 1 || row < 1 || tier < 1 {
		return nil, response.BadRequestError("Placement position is outside the Block dimensions.")
	}

	// placement and availability check
	slotNumbersToCheck := []int{slot}

	if is40ft {
		//must be placed in odd numbered slot
		if slot%2 == 0 {
			return nil, response.BadRequestError("40ft containers must start at an odd Slot number (Slot N).")
		}

		//needs 2 slots
		nextSlot := slot + 1

		if nextSlot > block.Slots {
			return nil, response.BadRequestError("Not enough space for 40ft container (requires Slot " + strconv.Itoa(nextSlot) + ").")
		}
		slotNumbersToCheck = append(slotNumbersToCheck, nextSlot)
	}

	// check availabiltiy in database
	count, err := s.ContainerPositionRepository.CheckPositionAvailability(
		s.DB,
		block.ID,
		row,
		tier,
		slotNumbersToCheck,
	)

	if err != nil {
		return nil, response.GeneralError("Database check failed: " + err.Error())
	}

	if count > 0 {
		return nil, response.GeneralError("Position already occupied by other container(s).")
	}

	// Check and get yard_plan
	var yardPlanID *int = nil
	yardPlan, err := s.YardPlanRepository.FindApplicablePlan(s.DB, block.ID, slot, row, size, request.Height, request.Type)
	if err == nil && yardPlan != nil {
		yardPlanID = &yardPlan.ID
	}

	var newPosition model.ContainerPosition

	txErr := s.DB.Transaction(func(tx *gorm.DB) error {
		newPosition = model.ContainerPosition{
			ContainerNumber: request.ContainerNumber,
			BlockID:         block.ID,
			SlotNumber:      slot,
			RowNumber:       row,
			TierNumber:      tier,

			ContainerSize:   size,
			ContainerHeight: request.Height,
			ContainerType:   request.Type,
			ContainerStatus: "STORAGE",

			ArrivalDate: time.Now(),
			YardPlanID:  yardPlanID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if saveErr := s.ContainerPositionRepository.Save(tx, &newPosition); saveErr != nil {
			return saveErr
		}
		return nil
	})

	if txErr != nil {
		return nil, response.RepositoryError("Failed to place container: " + txErr.Error())
	}

	return &web.PositionResponse{
		Block: block.Name,
		Slot:  newPosition.SlotNumber,
		Row:   newPosition.RowNumber,
		Tier:  newPosition.TierNumber,
	}, nil
}

func (s *ContainerServiceImpl) PickupContainer(ctx context.Context, request *web.PickupRequest) (*web.GeneralResponse, *response.CustomError) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	// check container
	var container model.ContainerPosition
	err := s.ContainerPositionRepository.FindByContainerNumber(s.DB, &container, request.ContainerNumber)
	if err != nil {
		return nil, response.NotFoundError("Container not found at any position or already picked up.")
	}

	// check stacking
	isStacked, err := s.ContainerPositionRepository.IsStackedAbove(
		s.DB,
		container.BlockID,
		container.SlotNumber,
		container.RowNumber,
		container.TierNumber,
	)

	if err != nil {
		return nil, response.GeneralError("Database check failed: " + err.Error())
	}

	if isStacked {
		return nil, response.GeneralError("Conflict: Cannot perform pickup. Another container is stacked on top.")
	}

	txErr := s.DB.Transaction(func(tx *gorm.DB) error {
		if deleteErr := s.ContainerPositionRepository.Delete(tx, container.ID); deleteErr != nil {
			return deleteErr
		}

		return nil
	})

	if txErr != nil {
		return nil, response.RepositoryError("Failed to perform container pickup: " + txErr.Error())
	}

	return &web.GeneralResponse{
		Message: "Success: Container picked up successfully.",
	}, nil
}
