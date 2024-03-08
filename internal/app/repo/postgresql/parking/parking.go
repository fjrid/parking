package parking

import (
	"context"
	"database/sql"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/pkg/sqkit"
	"go.uber.org/dig"
)

type (
	ParkingRepositoryImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}

	ParkingRepository interface {
		CreateParking(ctx context.Context, data dto.CreateParkingRequest) (res entity.Parking, err error)
		CreateParkingSlot(ctx context.Context, data dto.CreateParkingSlotRequest) (res entity.ParkingSlot, err error)
		FindParkingByName(ctx context.Context, name string) (results []entity.Parking, err error)
		FindSlotByID(ctx context.Context, id int64) (results []entity.ParkingSlot, err error)
		CreateParkingBook(ctx context.Context, data dto.CreateParkingBookRequest) (res entity.ParkingBook, err error)
		UpdateParkingSlot(ctx context.Context, data entity.ParkingSlot) (res entity.ParkingSlot, err error)
		FindBookByUserIDAndStatus(ctx context.Context, userID int64, status string) (results []entity.ParkingBook, err error)
		FindBookByIDAndUserIDAndStatus(ctx context.Context, id int64, userID int64, status string) (results []entity.ParkingBook, err error)
		UpdateParkingBook(ctx context.Context, data entity.ParkingBook) (res entity.ParkingBook, err error)
		FindParkingSlot(ctx context.Context, viewPagination dto.ViewPagination, opts ...sqkit.SelectOption) (results []entity.SummaryParkingSlot, vp dto.ViewPagination, err error)
		ParkingBookSummary(ctx context.Context, startDate string, endDate string) (result dto.SummaryParkingBook, err error)
	}
)
