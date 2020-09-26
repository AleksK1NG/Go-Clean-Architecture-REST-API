package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User full model
type User struct {
	ID          uuid.UUID  `json:"user_id" db:"user_id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password"`
	Role        *string    `json:"role,omitempty" db:"role"`
	About       *string    `json:"about,omitempty" db:"about"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number"`
	Address     *string    `json:"address,omitempty" db:"address"`
	City        *string    `json:"city,omitempty" db:"city"`
	Country     *string    `json:"country,omitempty" db:"country"`
	Gender      *string    `json:"gender,omitempty" db:"gender"`
	Postcode    *int       `json:"postcode,omitempty" db:"postcode"`
	Balance     float64    `json:"balance" db:"balance"`
	Birthday    *time.Time `json:"birthday,omitempty" db:"birthday"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LoginDate   time.Time  `json:"login_date" db:"login_date"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}

	return nil
}
