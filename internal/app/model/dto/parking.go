package dto

import "time"

type (
	CreateParkingRequest struct {
		Name      string `json:"name" valid:"required"`
		TotalSlot int64  `json:"total_slot" valid:"required"`

		CreatedBy string `json:"-"`
	}

	CreateParkingSlotRequest struct {
		ParkingID int64
		Number    int64
		Status    string

		CreatedBy string
	}

	BookingParkingRequest struct {
		ParkingSlotID int64  `json:"parking_slot_id" valid:"required"`
		CarNumber     string `json:"car_number" valid:"required"`
	}

	CreateParkingBookRequest struct {
		ParkingSlotID int64
		UserID        int64
		StartTime     time.Time
		Status        string
		CarNumber     string

		CreatedBy string
	}

	FinishBookRequest struct {
		BookID int64 `json:"book_id"`
	}

	ListSlotParkRequest struct {
		Status     string
		CarNumber  string
		ParkName   string
		ParkNumber string
	}

	ChangeMaintenanceRequest struct {
		ParkingSlotID int64 `json:"parking_slot_id" valid:"required"`
	}

	SummaryParkingBook struct {
		TotalTime    string  `json:"total_time"`
		TotalBooking int64   `json:"total_booking"`
		TotalFee     float64 `json:"total_fee"`
	}

	ParkingBookSummaryRequest struct {
		StartDate string `json:"start_date" valid:"required"`
		EndDate   string `json:"end_date" valid:"required"`
	}
)
