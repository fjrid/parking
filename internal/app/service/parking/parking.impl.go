package parking

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/pkg/dbtx"
	"github.com/fjrid/parking/pkg/sqkit"
)

// @ctor
func NewParkingSvc(impl ParkingServiceImpl) ParkingService {
	return &impl
}

func (s *ParkingServiceImpl) CreateParking(ctx context.Context, request dto.CreateParkingRequest) (err error) {

	isValid, err := govalidator.ValidateStruct(request)
	if err != nil {
		return
	}

	if !isValid {
		err = errors.New("data request invalid")
		return
	}

	userInterface := ctx.Value(constant.CONTEXT_USER_KEY)
	if userInterface != nil {
		user := userInterface.(dto.JWTClaims)

		request.CreatedBy = user.Email
	}

	tx := dbtx.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.Rollback(err)
		}
	}()

	availableParking, err := s.ParkingRepo.FindParkingByName(ctx, request.Name)
	if err != nil {
		return
	}

	if len(availableParking) > 0 {
		err = errors.New("parking name already used")
		return
	}

	parking, err := s.ParkingRepo.CreateParking(ctx, request)
	if err != nil {
		return
	}

	parkingSlots := make([]dto.CreateParkingSlotRequest, 0)
	for i := 1; i <= int(request.TotalSlot); i++ {
		parkingSlots = append(parkingSlots, dto.CreateParkingSlotRequest{
			ParkingID: parking.ID,
			Number:    int64(i),
			Status:    constant.PARKING_STATUS_AVAILABLE,
			CreatedBy: request.CreatedBy,
		})
	}

	for _, parkingSlot := range parkingSlots {
		_, err = s.ParkingRepo.CreateParkingSlot(ctx, parkingSlot)
		if err != nil {
			return
		}
	}

	tx.Commit()

	return
}

func (s *ParkingServiceImpl) BookingPark(ctx context.Context, request dto.BookingParkingRequest) (parkingBook entity.ParkingBook, err error) {
	isValid, err := govalidator.ValidateStruct(request)
	if err != nil {
		return
	}

	if !isValid {
		err = errors.New("data request invalid")
		return
	}

	tx := dbtx.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.Rollback(err)
		}
	}()

	parkingSlots, err := s.ParkingRepo.FindSlotByID(ctx, request.ParkingSlotID)
	if err != nil {
		return
	}

	if len(parkingSlots) == 0 {
		err = errors.New("parking slot not found")
		return
	}

	parkingSlot := parkingSlots[0]

	if parkingSlot.Status != constant.PARKING_STATUS_AVAILABLE {
		err = errors.New("parking slot unavailable")
		return
	}

	var user dto.JWTClaims
	userInterface := ctx.Value(constant.CONTEXT_USER_KEY)
	if userInterface != nil {
		user = userInterface.(dto.JWTClaims)
	}

	bookings, err := s.ParkingRepo.FindBookByUserIDAndStatus(ctx, user.ID, constant.BOOK_STATUS_ON_GOING)
	if err != nil {
		return
	}

	if len(bookings) > 0 {
		err = errors.New("you have active booking")
		return
	}

	parkingBook, err = s.ParkingRepo.CreateParkingBook(ctx, dto.CreateParkingBookRequest{
		UserID:        user.ID,
		ParkingSlotID: parkingSlot.ID,
		StartTime:     time.Now(),
		Status:        constant.BOOK_STATUS_ON_GOING,
		CreatedBy:     user.Email,
		CarNumber:     request.CarNumber,
	})
	if err != nil {
		return
	}

	parkingSlot.Status = constant.PARKING_STATUS_BOOKED
	parkingSlot.ModifiedAt = time.Now()
	parkingSlot.ModifiedBy = user.Email
	_, err = s.ParkingRepo.UpdateParkingSlot(ctx, parkingSlot)
	if err != nil {
		return
	}

	tx.Commit()

	return
}

func (s *ParkingServiceImpl) FinishBook(ctx context.Context, request dto.FinishBookRequest) (parkingBook entity.ParkingBook, err error) {
	tx := dbtx.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.Rollback(err)
		}
	}()

	var user dto.JWTClaims
	userInterface := ctx.Value(constant.CONTEXT_USER_KEY)
	if userInterface != nil {
		user = userInterface.(dto.JWTClaims)
	}

	parkingBooks, err := s.ParkingRepo.FindBookByIDAndUserIDAndStatus(ctx, request.BookID, user.ID, constant.BOOK_STATUS_ON_GOING)
	if err != nil {
		return
	}

	if len(parkingBooks) == 0 {
		err = errors.New("booking not found")
		return
	}

	parkingBook = parkingBooks[0]
	now := time.Now()
	startTime := time.Date(parkingBook.StartTime.Year(), parkingBook.StartTime.Month(), parkingBook.StartTime.Day(), parkingBook.StartTime.Hour(), parkingBook.StartTime.Minute(), parkingBook.StartTime.Second(), parkingBook.StartTime.Nanosecond(), time.UTC)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)

	diffTime := endTime.Sub(startTime)
	fee := math.Ceil(diffTime.Hours()) * 10

	parkingBook.EndTime = &endTime
	parkingBook.Status = constant.BOOK_STATUS_FINISHED
	parkingBook.Fee = fee
	parkingBook.ModifiedAt = time.Now()
	parkingBook.ModifiedBy = user.Email

	_, err = s.ParkingRepo.UpdateParkingBook(ctx, parkingBook)
	if err != nil {
		return
	}

	parkingSlots, err := s.ParkingRepo.FindSlotByID(ctx, parkingBook.ParkingSlotID)
	if err != nil {
		return
	}

	if len(parkingSlots) == 0 {
		err = errors.New("parking slot not found")
		return
	}

	parkingSlot := parkingSlots[0]
	parkingSlot.Status = constant.PARKING_STATUS_AVAILABLE
	parkingSlot.ModifiedAt = time.Now()
	parkingSlot.ModifiedBy = user.Email

	_, err = s.ParkingRepo.UpdateParkingSlot(ctx, parkingSlot)
	if err != nil {
		return
	}

	tx.Commit()

	return
}

func (s *ParkingServiceImpl) ListSlotPark(ctx context.Context, viewPagination dto.ViewPagination, request dto.ListSlotParkRequest) (results []entity.SummaryParkingSlot, vp dto.ViewPagination, err error) {
	condition := []sqkit.SelectOption{}

	if request.Status != "" {
		condition = append(condition, sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.Status): request.Status})
	}

	if request.CarNumber != "" {
		condition = append(condition, sqkit.Where{
			fmt.Sprintf("%s.%s LIKE ('%s')", entity.ParkingBookTableName, entity.ParkingBookTable.CarNumber, "%"+request.CarNumber+"%"),
		})
	}

	if request.ParkName != "" {
		condition = append(condition, sqkit.Where{
			fmt.Sprintf("%s.%s LIKE ('%s')", entity.ParkingTableName, entity.ParkingTable.Name, "%"+request.ParkName+"%"),
		})
	}

	if request.ParkNumber != "" {
		condition = append(condition, sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.Number): request.ParkNumber})
	}

	return s.ParkingRepo.FindParkingSlot(ctx, viewPagination, condition...)
}

func (s *ParkingServiceImpl) ChangeMaintenance(ctx context.Context, request dto.ChangeMaintenanceRequest) (err error) {
	isValid, err := govalidator.ValidateStruct(request)
	if err != nil {
		return
	}

	if !isValid {
		err = errors.New("data request invalid")
		return
	}

	tx := dbtx.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.Rollback(err)
		}
	}()

	parkingSlots, err := s.ParkingRepo.FindSlotByID(ctx, request.ParkingSlotID)
	if err != nil {
		return
	}

	if len(parkingSlots) == 0 {
		err = errors.New("parking slot not found")
		return
	}

	parkingSlot := parkingSlots[0]

	if parkingSlot.Status == constant.PARKING_STATUS_BOOKED {
		err = errors.New("parking already booked")
		return
	}

	var user dto.JWTClaims
	userInterface := ctx.Value(constant.CONTEXT_USER_KEY)
	if userInterface != nil {
		user = userInterface.(dto.JWTClaims)
	}

	if parkingSlot.Status == constant.PARKING_STATUS_AVAILABLE {
		parkingSlot.Status = constant.PARKING_STATUS_MAINTENANCE
	} else {
		parkingSlot.Status = constant.PARKING_STATUS_AVAILABLE
	}

	parkingSlot.ModifiedAt = time.Now()
	parkingSlot.ModifiedBy = user.Email
	_, err = s.ParkingRepo.UpdateParkingSlot(ctx, parkingSlot)
	if err != nil {
		return
	}

	tx.Commit()

	return
}

func (s *ParkingServiceImpl) ParkingBookSummary(ctx context.Context, request dto.ParkingBookSummaryRequest) (result dto.SummaryParkingBook, err error) {
	isValid, err := govalidator.ValidateStruct(request)
	if err != nil {
		return
	}

	if !isValid {
		err = errors.New("data request invalid")
		return
	}

	return s.ParkingRepo.ParkingBookSummary(ctx, request.StartDate, request.EndDate)
}
