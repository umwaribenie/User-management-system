package models

// GetAllUsersRequest maps to the query parameters for fetching all users.
type GetAllUsersRequest struct {
	PageNumber int    `form:"pageNumber"`
	PageSize   int    `form:"pageSize"`
	From       string `form:"from"`
	To         string `form:"to"`
	Search     string `form:"search"`
	Role       string `form:"role"`
	Status     string `form:"status"`
}

// CreateUserRequest is the model for self-registration.
type CreateUserRequest struct {
	ClientID       string  `json:"clientId" binding:"required"`
	Email          string  `json:"email" binding:"required,email"`
	FirstName      string  `json:"firstName" binding:"required"`
	LastName       string  `json:"lastName" binding:"required"`
	NationalID     *string `json:"nationalId"`
	PassportNumber *string `json:"passportNumber"`
	Password       string  `json:"password" binding:"required,min=6"`
	Phone          string  `json:"phone" binding:"required"`
	ProfilePicture *string `json:"profilePicture"`
	Username       string  `json:"username" binding:"required"`
}

// CreateUserByAdminRequest is the model for admin-driven user creation.
type CreateUserByAdminRequest struct {
	Email          string   `json:"email" binding:"required,email"`
	FirstName      string   `json:"firstName" binding:"required"`
	LastName       string   `json:"lastName" binding:"required"`
	NationalID     *string  `json:"nationalId"`
	PassportNumber *string  `json:"passportNumber"`
	Password       string   `json:"password" binding:"required,min=6"`
	Phone          string   `json:"phone" binding:"required"`
	ProfilePicture *string  `json:"profilePicture"`
	Role           UserRole `json:"role" binding:"required,oneof=user admin"` // Updated validation
	Username       string   `json:"username" binding:"required"`
}

// UpdatePasswordRequest is used for all password update scenarios.
type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
	OldPassword string `json:"oldPassword" binding:"required"`
}

// UpdateUserRequest is for patching a user's details.
type UpdateUserRequest struct {
	Email          *string   `json:"email,omitempty"`
	FirstName      *string   `json:"firstName,omitempty"`
	LastName       *string   `json:"lastName,omitempty"`
	NationalID     *string   `json:"nationalId,omitempty"`
	PassportNumber *string   `json:"passportNumber,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	ProfilePicture *string   `json:"profilePicture,omitempty"`
	Role           *UserRole `json:"role,omitempty" binding:"omitempty,oneof=user admin"`
	Username       *string   `json:"username,omitempty"`
}

// Auth-related request models
type ConfirmOtpRequest struct {
	Otp      string `json:"otp" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	ClientID string `json:"clientId"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username"`
}

type PasswordResetRequest struct {
	ClientID string `json:"clientId"`
	Username string `json:"username"`
}

type ResetPasswordRequest struct {
	ClientID string `json:"clientId"`
	Username string `json:"username"`
}

type ResetPasswordWithTokenRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
