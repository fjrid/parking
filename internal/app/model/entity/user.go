package entity

import "time"

type (
	User struct {
		ID         int64
		Email      string
		Password   string
		Role       string
		CreatedAt  time.Time
		CreatedBy  string
		ModifiedAt time.Time
		ModifiedBy string
		DeletedAt  *time.Time
		DeletedBy  *string
	}
)

var (
	UserTableName = "users"
	UserTable     = struct {
		ID         string
		Email      string
		Password   string
		Role       string
		CreatedAt  string
		CreatedBy  string
		ModifiedAt string
		ModifiedBy string
		DeletedAt  string
		DeletedBy  string
	}{
		ID:         "id",
		Email:      "email",
		Password:   "password",
		Role:       "role",
		CreatedAt:  "created_at",
		CreatedBy:  "created_by",
		ModifiedAt: "modified_at",
		ModifiedBy: "modified_by",
		DeletedAt:  "deleted_at",
		DeletedBy:  "deleted_by",
	}
)
