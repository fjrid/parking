package parking

import (
	"context"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/internal/app/repo/postgresql/parking"
	"go.uber.org/dig"
)

type (
	ParkingServiceImpl struct {
		dig.In
		ParkingRepo parking.ParkingRepository
	}

	ParkingService interface {
		CreateParking(ctx context.Context, request dto.CreateParkingRequest) (err error)
		BookingPark(ctx context.Context, request dto.BookingParkingRequest) (parkingBook entity.ParkingBook, err error)
		FinishBook(ctx context.Context, request dto.FinishBookRequest) (parkingBook entity.ParkingBook, err error)
		ListSlotPark(ctx context.Context, viewPagination dto.ViewPagination, request dto.ListSlotParkRequest) (results []entity.SummaryParkingSlot, vp dto.ViewPagination, err error)
		ChangeMaintenance(ctx context.Context, request dto.ChangeMaintenanceRequest) (err error)
		ParkingBookSummary(ctx context.Context, request dto.ParkingBookSummaryRequest) (result dto.SummaryParkingBook, err error)
	}
)
