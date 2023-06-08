package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID    uuid.UUID `json:"session_id" db:"session_id" redis:"session_id" validate:"omitempty"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token" redis:"refresh_token" validate:"required"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at" redis:"expires_at" validate:"required"`
	UserID       uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"required"`
}

// func (s *Session) HashRefreshToken() error {
// 	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(s.RefreshToken), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	s.RefreshToken = string(hashedRefreshToken)
// 	return nil
// }
