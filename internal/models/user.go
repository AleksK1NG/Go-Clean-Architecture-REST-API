package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User full model
type User struct {
	ID          uuid.UUID  `json:"user_id" db:"user_id" redis:"user_id" validate:"omitempty,uuid"`
	FirstName   string     `json:"first_name" db:"first_name" redis:"first_name" validate:"required,lte=30"`
	LastName    string     `json:"last_name" db:"last_name" redis:"last_name" validate:"required,lte=30"`
	Email       string     `json:"email" db:"email" redis:"email" validate:"omitempty,lte=60,email"`
	Password    string     `json:"password,omitempty" db:"password" redis:"password" validate:"required,gte=6"`
	Role        *string    `json:"role,omitempty" db:"role" redis:"role" validate:"omitempty,lte=10"`
	About       *string    `json:"about,omitempty" db:"about" redis:"about" validate:"omitempty,lte=1024"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar" redis:"avatar" validate:"omitempty,lte=512,url"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" redis:"phone_number" validate:"omitempty,lte=20"`
	Address     *string    `json:"address,omitempty" db:"address" redis:"address" validate:"omitempty,lte=250"`
	City        *string    `json:"city,omitempty" db:"city" redis:"city" validate:"omitempty,lte=24"`
	Country     *string    `json:"country,omitempty" db:"country" redis:"country" validate:"omitempty,lte=24"`
	Gender      *string    `json:"gender,omitempty" db:"gender" redis:"gender" validate:"omitempty,lte=10"`
	Postcode    *int       `json:"postcode,omitempty" db:"postcode" redis:"postcode" validate:"omitempty"`
	Birthday    *time.Time `json:"birthday,omitempty" db:"birthday" redis:"birthday" validate:"omitempty,lte=10"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" redis:"updated_at"`
	LoginDate   time.Time  `json:"login_date" db:"login_date" redis:"login_date"`
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

// User full model
type UserUpdate struct {
	ID          uuid.UUID  `json:"user_id" db:"user_id" validate:"required,omitempty"`
	FirstName   string     `json:"first_name" db:"first_name" validate:"lte=30"`
	LastName    string     `json:"last_name" db:"last_name" validate:"lte=30"`
	Email       string     `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Role        *string    `json:"role,omitempty" db:"role" validate:"omitempty,lte=10"`
	About       *string    `json:"about,omitempty" db:"about" validate:"omitempty,lte=1024"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar" validate:"omitempty,lte=512,url"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" validate:"omitempty,lte=20"`
	Address     *string    `json:"address,omitempty" db:"address" validate:"omitempty,lte=250"`
	City        *string    `json:"city,omitempty" db:"city" validate:"omitempty,lte=24"`
	Country     *string    `json:"country,omitempty" db:"country" validate:"omitempty,lte=24"`
	Gender      *string    `json:"gender,omitempty" db:"gender" validate:"omitempty,lte=10"`
	Postcode    *int       `json:"postcode,omitempty" db:"postcode" validate:"omitempty"`
	Balance     float64    `json:"balance" db:"balance"`
	Birthday    *time.Time `json:"birthday,omitempty" db:"birthday" validate:"omitempty,lte=10"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LoginDate   time.Time  `json:"login_date" db:"login_date"`
}

// Prepare user for register
func (u *UserUpdate) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}

// All Users response
type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}
