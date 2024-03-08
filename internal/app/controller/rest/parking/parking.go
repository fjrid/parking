package parking

import (
	"net/http"

	"github.com/fjrid/parking/internal/app/infra"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/service/parking"
	"github.com/fjrid/parking/pkg/echokit"
	"github.com/fjrid/parking/pkg/utinterface"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book repo
	ParkingRestCtrl struct {
		dig.In
		ParkingSvc parking.ParkingService
	}
)

var _ echokit.Router = (*ParkingRestCtrl)(nil)

// SetRoute to define API Route
func (c *ParkingRestCtrl) SetRoute(e echokit.Server) {
	r := e.Group("/parking")
	r.Use(infra.JWTMiddleware())

	r.POST("", c.Create, infra.AdminMiddleware)
	r.POST("/:parkingSlotID/book", c.Book)
	r.POST("/:bookingID/finish", c.Finish)
	r.GET("/slot", c.FindSlot)
	r.POST("/:parkingSlotID/maintenance", c.ChangeMaintenance)
	r.GET("/book-summary", c.BookingSummary)
}

func (c *ParkingRestCtrl) Create(ec echo.Context) (err error) {
	var request dto.CreateParkingRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	ctx := ec.Request().Context()

	err = c.ParkingSvc.CreateParking(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
	})
}

func (c *ParkingRestCtrl) Book(ec echo.Context) (err error) {
	var request dto.BookingParkingRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	request.ParkingSlotID = utinterface.ToInt(ec.Param("parkingSlotID"), 0)

	ctx := ec.Request().Context()

	result, err := c.ParkingSvc.BookingPark(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
		Data:    result,
	})
}

func (c *ParkingRestCtrl) Finish(ec echo.Context) (err error) {
	var request dto.FinishBookRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	request.BookID = utinterface.ToInt(ec.Param("bookingID"), 0)

	ctx := ec.Request().Context()

	result, err := c.ParkingSvc.FinishBook(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
		Data:    result,
	})
}

func (c *ParkingRestCtrl) FindSlot(ec echo.Context) (err error) {
	var (
		request = dto.ListSlotParkRequest{
			Status:     ec.QueryParam("status"),
			CarNumber:  ec.QueryParam("car_number"),
			ParkName:   ec.QueryParam("park_name"),
			ParkNumber: ec.QueryParam("park_number"),
		}
		viewPagination = dto.ViewPagination{
			Limit:  utinterface.ToInt(ec.QueryParam("limit"), 10),
			Offset: utinterface.ToInt(ec.QueryParam("offset"), 0),
		}
	)

	ctx := ec.Request().Context()

	result, vp, err := c.ParkingSvc.ListSlotPark(ctx, viewPagination, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
		Data:    result,
		Meta: map[string]interface{}{
			"pagination": vp,
		},
	})
}

func (c *ParkingRestCtrl) ChangeMaintenance(ec echo.Context) (err error) {
	var request dto.ChangeMaintenanceRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	request.ParkingSlotID = utinterface.ToInt(ec.Param("parkingSlotID"), 0)

	ctx := ec.Request().Context()

	err = c.ParkingSvc.ChangeMaintenance(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
	})
}

func (c *ParkingRestCtrl) BookingSummary(ec echo.Context) (err error) {
	var (
		request = dto.ParkingBookSummaryRequest{
			StartDate: ec.QueryParam("start_date"),
			EndDate:   ec.QueryParam("end_date"),
		}
	)

	ctx := ec.Request().Context()

	result, err := c.ParkingSvc.ParkingBookSummary(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
		Data:    result,
	})
}
