package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole defines the type for user roles which can be user or admin.
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// UserStatus defines the type for user statuses (active,inactive and deleted).
type UserStatus string

const (
	ActiveStatus   UserStatus = "active"
	InactiveStatus UserStatus = "inactive"
	DeletedStatus  UserStatus = "deleted"
)

type User struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClientID       string         `gorm:"unique" json:"clientId"`
	Email          string         `gorm:"uniqueIndex" json:"email"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	NationalID     *string        `gorm:"unique" json:"nationalId,omitempty"`
	PassportNumber *string        `gorm:"unique" json:"passportNumber,omitempty"`
	Password       string         `json:"-"`
	Phone          string         `gorm:"unique" json:"phone"`
	ProfilePicture *string        `json:"profilePicture,omitempty"`
	Username       string         `gorm:"uniqueIndex" json:"username"`
	Slug           string         `gorm:"uniqueIndex" json:"slug"`
	Role           UserRole       `gorm:"type:varchar(50);default:'user'" json:"role"`
	Status         UserStatus     `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
