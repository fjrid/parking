package entity

import "time"

type (
	Parking struct {
		ID         int64
		Name       string
		CreatedAt  time.Time
		CreatedBy  string
		ModifiedAt time.Time
		ModifiedBy string
		DeletedAt  *time.Time
		DeletedBy  *string
	}

	ParkingSlot struct {
		ID          int64
		ParkingID   int64
		ParkingName string
		Number      int64
		Status      string
		CreatedAt   time.Time
		CreatedBy   string
		ModifiedAt  time.Time
		ModifiedBy  string
		DeletedAt   *time.Time
		DeletedBy   *string
	}

	ParkingBook struct {
		ID            int64
		UserID        int64
		ParkingSlotID int64
		StartTime     time.Time
		EndTime       *time.Time
		Fee           float64
		Status        string
		CarNumber     string
		CreatedAt     time.Time
		CreatedBy     string
		ModifiedAt    time.Time
		ModifiedBy    string
		DeletedAt     *time.Time
		DeletedBy     *string
	}

	SummaryParkingSlot struct {
		ID            int64
		ParkingName   string
		ParkingNumber int64
		Status        string
		CarNumber     *string
	}
)

var (
	ParkingTableName = "parking"
	ParkingTable     = struct {
		ID         string
		Name       string
		CreatedAt  string
		CreatedBy  string
		ModifiedAt string
		ModifiedBy string
		DeletedAt  string
		DeletedBy  string
	}{
		ID:         "id",
		Name:       "name",
		CreatedAt:  "created_at",
		CreatedBy:  "created_by",
		ModifiedAt: "modified_at",
		ModifiedBy: "modified_by",
		DeletedAt:  "deleted_at",
		DeletedBy:  "deleted_by",
	}

	ParkingSlotTableName = "parking_slot"
	ParkingSlotTable     = struct {
		ID         string
		ParkingID  string
		Number     string
		Status     string
		CreatedAt  string
		CreatedBy  string
		ModifiedAt string
		ModifiedBy string
		DeletedAt  string
		DeletedBy  string
	}{
		ID:         "id",
		ParkingID:  "parking_id",
		Number:     "number",
		Status:     "status",
		CreatedAt:  "created_at",
		CreatedBy:  "created_by",
		ModifiedAt: "modified_at",
		ModifiedBy: "modified_by",
		DeletedAt:  "deleted_at",
		DeletedBy:  "deleted_by",
	}

	ParkingBookTableName = "parking_book"
	ParkingBookTable     = struct {
		ID            string
		UserID        string
		ParkingSlotID string
		StartTime     string
		EndTime       string
		Fee           string
		Status        string
		CarNumber     string
		CreatedAt     string
		CreatedBy     string
		ModifiedAt    string
		ModifiedBy    string
		DeletedAt     string
		DeletedBy     string
	}{
		ID:            "id",
		UserID:        "user_id",
		ParkingSlotID: "parking_slot_id",
		StartTime:     "start_time",
		EndTime:       "end_time",
		Fee:           "fee",
		Status:        "status",
		CarNumber:     "car_number",
		CreatedAt:     "created_at",
		CreatedBy:     "created_by",
		ModifiedAt:    "modified_at",
		ModifiedBy:    "modified_by",
		DeletedAt:     "deleted_at",
		DeletedBy:     "deleted_by",
	}
)
