package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/repositories"
	"github.com/umwaribenie/final_user_management/utils"
)

var ctx = context.Background()

type AuthService interface {
	CheckAuth() (models.SuccessResponse, error)
	ConfirmPasswordResetOtp(request models.ConfirmOtpRequest) (models.SuccessResponse, error)
	Login(request models.LoginRequest) (models.LoginResponse, error)
	RequestPasswordReset(request models.PasswordResetRequest) (models.SuccessResponse, error)
	ResetPasswordViaEmail(request models.ResetPasswordRequest) (models.SuccessResponse, error)
	UpdatePassword(userID string, request models.UpdatePasswordRequest) (models.SuccessResponse, error)
	ResetPasswordWithToken(tokenString string, newPassword string) (models.SuccessResponse, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	redisClient *redis.Client
}

// NewAuthService constructor
func NewAuthService(userRepo repositories.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *authService) ResetPasswordViaEmail(request models.ResetPasswordRequest) (models.SuccessResponse, error) {
	// 1. Find user
	var (
		user *models.User
		err  error
	)
	if request.Username != "" {
		user, err = s.userRepo.FindByUsername(request.Username)
	} else if request.ClientID != "" {
		user, err = s.userRepo.FindByID(request.ClientID)
	} else {
		return models.SuccessResponse{}, errors.New("username or clientID is required")
	}
	if err != nil {
		log.Printf("Password-reset attempt for non-existent user: %v", err)
		return models.SuccessResponse{Message: "If an account with that identifier exists, an OTP has been sent."}, nil
	}
	if user.Email == "" {
		return models.SuccessResponse{}, errors.New("user has no email")
	}

	// 2. Generate OTP
	otp := utils.GenerateOTP()

	// 3. Store OTP two ways
	redisOTPKey := "otp:" + otp
	redisEmailKey := "email:" + user.Email
	pipe := s.redisClient.TxPipeline()
	pipe.Set(ctx, redisOTPKey, user.Email, 5*time.Minute)
	pipe.Set(ctx, redisEmailKey, otp, 5*time.Minute)
	if _, err = pipe.Exec(ctx); err != nil {
		log.Printf("Redis error: %v", err)
		return models.SuccessResponse{}, errors.New("failed to store OTP")
	}

	// 4. Send email
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	emailCfg := utils.EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"), SMTPPort: smtpPort,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
	}
	body := fmt.Sprintf("Your OTP for password reset is: <strong>%s</strong><br>It expires in 5 minutes.", otp)
	if err := utils.SendPasswordResetEmail(emailCfg, user.Email, body); err != nil {
		log.Printf("Email send error: %v", err)
		return models.SuccessResponse{}, errors.New("failed to send OTP email")
	}

	log.Printf("OTP %s emailed to %s", otp, user.Email)
	return models.SuccessResponse{Message: "Password reset OTP sent via email"}, nil
}

func (s *authService) CheckAuth() (models.SuccessResponse, error) {

	return models.SuccessResponse{Message: "User is authenticated"}, nil
}

func (s *authService) ConfirmPasswordResetOtp(req models.ConfirmOtpRequest) (models.SuccessResponse, error) {
	// 1. Lookup email by OTP
	email, err := s.redisClient.Get(ctx, "otp:"+req.Otp).Result()
	if err == redis.Nil {
		return models.SuccessResponse{}, errors.New("invalid or expired OTP")
	} else if err != nil {
		return models.SuccessResponse{}, err
	}

	// 2. Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return models.SuccessResponse{}, errors.New("user not found")
	}

	// 3. Hash and update password
	hashed, _ := utils.HashPassword(req.Password)
	if err := s.userRepo.UpdatePassword(user.ID, hashed); err != nil {
		return models.SuccessResponse{}, err
	}

	// 4. Cleanup Redis keys
	s.redisClient.Del(ctx, "otp:"+req.Otp, "email:"+email)

	return models.SuccessResponse{Message: "Password reset successful"}, nil
}

// services/auth_service.go

// func (s *authService) RequestPasswordReset(request models.PasswordResetRequest) (models.SuccessResponse, error) {
// 	// 1. Find the user by username or clientId
// 	var user *models.User
// 	var err error

// 	if request.Username != "" {
// 		user, err = s.userRepo.FindByUsername(request.Username)
// 	} else if request.ClientID != "" {
// 		user, err = s.userRepo.FindByID(request.ClientID)
// 	} else {
// 		return models.SuccessResponse{}, errors.New("username or clientID is required")
// 	}

// 	if err != nil {
// 		log.Printf("Password reset attempt for non-existent user (Username: %s, ClientID: %s). Error: %v", request.Username, request.ClientID, err)
// 		return models.SuccessResponse{Message: "If an account with that identifier exists, an OTP has been sent."}, nil
// 	}

// 	// 2. Generate OTP
// 	otp := utils.GenerateOTP()
// 	if otp == "" {
// 		return models.SuccessResponse{}, errors.New("failed to generate OTP")
// 	}

// 	// 3. Store OTP in Redis with a 5-minute expiration (keyed by email if present, else phone)
// 	var redisKey string
// 	if user.Email != "" {
// 		redisKey = "otp:" + user.Email
// 	} else if user.Phone != "" {
// 		redisKey = "otp:" + user.Phone
// 	} else {
// 		return models.SuccessResponse{}, errors.New("user has neither email nor phone for OTP delivery")
// 	}
// 	err = s.redisClient.Set(ctx, redisKey, otp, 5*time.Minute).Err()
// 	if err != nil {
// 		log.Printf("Failed to store OTP for user %s in Redis: %v", user.Email, err)
// 		return models.SuccessResponse{}, errors.New("failed to set up password reset process")
// 	}

// 	// 4. Send OTP via Email (if user has email)
// 	emailSent := false
// 	if user.Email != "" {
// 		emailConfig := utils.EmailConfig{
// 			SMTPHost:     os.Getenv("SMTP_HOST"),
// 			SMTPPort:     587, // Or use os.Getenv("SMTP_PORT") and convert to int
// 			SMTPUsername: os.Getenv("SMTP_USERNAME"),
// 			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
// 			FromEmail:    os.Getenv("FROM_EMAIL"),
// 		}
// 		emailBody := fmt.Sprintf("Your One-Time Password (OTP) for password reset is: <strong>%s</strong><br>It will expire in 5 minutes.", otp)
// 		if err := utils.SendPasswordResetEmail(emailConfig, user.Email, emailBody); err != nil {
// 			log.Printf("Failed to send OTP email to %s: %v", user.Email, err)
// 		} else {
// 			emailSent = true
// 			log.Printf("OTP sent via email to %s", user.Email)
// 		}
// 	}

// 	// 5. Send OTP via SMS (if user has phone)
// 	smsSent := false
// 	if user.Phone != "" {
// 		if err := utils.SendPasswordResetOTPSMS(user.Phone, otp); err != nil {
// 			log.Printf("Failed to send OTP SMS to %s: %v", user.Phone, err)
// 		} else {
// 			smsSent = true
// 			log.Printf("OTP sent via SMS to %s", user.Phone)
// 		}
// 	}

// 	// 6. Return a generic message
// 	if emailSent && smsSent {
// 		return models.SuccessResponse{Message: "Password reset OTP sent via email and SMS"}, nil
// 	} else if emailSent {
// 		return models.SuccessResponse{Message: "Password reset OTP sent via email"}, nil
// 	} else if smsSent {
// 		return models.SuccessResponse{Message: "Password reset OTP sent via SMS"}, nil
// 	} else {
// 		return models.SuccessResponse{}, errors.New("failed to send OTP via both email and SMS")
// 	}
// }
// services/auth_service.go

// services/auth_service.go

func (s *authService) RequestPasswordReset(request models.PasswordResetRequest) (models.SuccessResponse, error) {
	// 1. Find user
	var (
		user *models.User
		err  error
	)
	if request.Username != "" {
		user, err = s.userRepo.FindByUsername(request.Username)
	} else if request.ClientID != "" {
		user, err = s.userRepo.FindByID(request.ClientID)
	} else {
		return models.SuccessResponse{}, errors.New("username or clientID is required")
	}
	if err != nil {
		log.Printf("Password-reset attempt for non-existent user: %v", err)
		return models.SuccessResponse{Message: "If an account with that identifier exists, an OTP has been sent."}, nil
	}
	if user.Email == "" {
		return models.SuccessResponse{}, errors.New("user has no email")
	}

	// 2. Generate OTP
	otp := utils.GenerateOTP()

	// 3. Store OTP both ways in Redis (5-minute TTL)
	redisOTPKey := "otp:" + otp
	redisEmailKey := "email:" + user.Email
	pipe := s.redisClient.TxPipeline()
	pipe.Set(ctx, redisOTPKey, user.Email, 5*time.Minute)
	pipe.Set(ctx, redisEmailKey, otp, 5*time.Minute)
	if _, err = pipe.Exec(ctx); err != nil {
		log.Printf("Redis error: %v", err)
		return models.SuccessResponse{}, errors.New("failed to store OTP")
	}

	// 4. Email the OTP
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	emailCfg := utils.EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"), SMTPPort: smtpPort,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
	}
	body := fmt.Sprintf("Your OTP for password reset is: <strong>%s</strong><br>It expires in 5 minutes.", otp)
	if err := utils.SendPasswordResetEmail(emailCfg, user.Email, body); err != nil {
		log.Printf("Email send error: %v", err)
		return models.SuccessResponse{}, errors.New("failed to send OTP email")
	}

	log.Printf("OTP %s emailed to %s", otp, user.Email)
	return models.SuccessResponse{Message: "Password reset OTP sent via email"}, nil
}

func (s *authService) UpdatePassword(userID string, request models.UpdatePasswordRequest) (models.SuccessResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return models.SuccessResponse{}, err
	}

	if !utils.CheckPasswordHash(request.OldPassword, user.Password) {
		return models.SuccessResponse{}, errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return models.SuccessResponse{}, err
	}

	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return models.SuccessResponse{}, err
	}

	return models.SuccessResponse{Message: "Password updated successfully"}, nil
}

func (s *authService) ResetPasswordWithToken(tokenString string, newPassword string) (models.SuccessResponse, error) {
	// 1. Verify token and extract email
	email, err := utils.VerifyPasswordResetToken(tokenString)
	if err != nil {
		return models.SuccessResponse{}, errors.New("invalid or expired reset token")
	}

	// 2. Find the user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return models.SuccessResponse{}, errors.New("user not found")
	}

	// 3. Hash the new password
	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		return models.SuccessResponse{}, err
	}

	// 4. Update the user’s password in the DB
	if err := s.userRepo.UpdatePassword(user.ID, hashedPwd); err != nil {
		return models.SuccessResponse{}, err
	}

	return models.SuccessResponse{Message: "Password has been reset successfully"}, nil
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999-100000) + 100000)
}
func (s *authService) Login(request models.LoginRequest) (models.LoginResponse, error) {
	var user *models.User
	var err error

	// The `Username` field in the request can be used for either username or email login
	if request.Username != "" {
		// Try to find by username first
		user, err = s.userRepo.FindByUsername(request.Username)
		if err != nil || user == nil {
			// If not found by username, try by email
			user, err = s.userRepo.FindByEmail(request.Username)
		}
	} else if request.ClientID != "" {
		// Support login by ClientID if needed
		user, err = s.userRepo.FindByID(request.ClientID)
	} else {
		return models.LoginResponse{}, errors.New("username or clientID is required for login")
	}

	// If user is still not found after checking all methods, credentials are invalid
	if err != nil || user == nil {
		return models.LoginResponse{}, errors.New("invalid credentials")
	}

	// Verify the password
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return models.LoginResponse{}, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		return models.LoginResponse{}, errors.New("failed to generate token")
	}

	return models.LoginResponse{AccessToken: token}, nil
}
