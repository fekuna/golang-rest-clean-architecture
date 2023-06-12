package models

import (
	"time"

	"github.com/google/uuid"
)

type Avatar struct {
	AvatarID  uuid.UUID `json:"avatar_id" db:"avatar_id" redis:"avatar_id" validate:"omitempty"`
	Bucket    string    `json:"bucket" db:"bucket" redis:"bucket" validate:"required"`
	FilePath  string    `json:"file_path" db:"file_path" redis:"file_path" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
}
