package domain

import (
	"context"
	"time"
)

type Role int32

const (
	RoleUnspecified     Role = 0
	RoleSuperadmin      Role = 1
	RoleAdmin           Role = 2
	RoleDoctor          Role = 3
	RoleNurse           Role = 4
	RolePharmacist      Role = 5
	RoleAdmission       Role = 6
	RoleCashier         Role = 7
	RoleLabTechnician   Role = 8
	RoleRadiographer    Role = 9
	RoleNutritionist    Role = 10
	RolePhysiotherapist Role = 11
	RoleOther           Role = 12
)

type User struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	FullName       string     `json:"full_name"`
	Role           Role       `json:"role"`
	PractitionerID string     `json:"practitioner_id"`
	IsActive       bool       `json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UserFilter struct {
	Query    string
	Role     Role
	IsActive *bool
	Limit    int
	Offset   int
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, filter UserFilter) ([]User, int64, error)
	Update(ctx context.Context, u *User) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error
}
