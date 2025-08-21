package controllers

import (
	"net/http"

	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService}
}

// @Summary Check if the user is authenticated
// @Description Checks if the provided JWT token is valid. This endpoint requires an authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/check [get]
func (c *AuthController) CheckAuth(ctx *gin.Context) {
	// This endpoint should be protected by a middleware that validates the JWT.
	// If the request reaches this handler, it means the token is valid.
	response, _ := c.authService.CheckAuth()
	ctx.JSON(http.StatusOK, response)
}

// @Summary Confirm password reset OTP
// @Description Confirms the OTP and sets a new password.
// @Tags auth
// @Accept json
// @Produce json
// @Param confirmOTP body models.ConfirmOtpRequest true "Confirm OTP request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/confirm-password-reset-otp [post]
func (c *AuthController) ConfirmPasswordResetOtp(ctx *gin.Context) {
	var request models.ConfirmOtpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	response, err := c.authService.ConfirmPasswordResetOtp(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Login a user
// @Description Authenticates a user and returns a JWT access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginData body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var request models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	response, err := c.authService.Login(request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Request password reset
// @Description Sends a password reset OTP to the user's registered phone or email.
// @Tags auth
// @Accept json
// @Produce json
// @Param passwordReset body models.PasswordResetRequest true "Password reset request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/password-reset [post]
func (c *AuthController) RequestPasswordReset(ctx *gin.Context) {
	var request models.PasswordResetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	response, err := c.authService.RequestPasswordReset(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Reset password via email
// @Description Sends a password reset OTP to the user's registered email.
// @Tags auth
// @Accept json
// @Produce json
// @Param resetPassword body models.ResetPasswordRequest true "Reset password via email request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/reset-password/email [post]
func (c *AuthController) ResetPasswordViaEmail(ctx *gin.Context) {
	var request models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	response, err := c.authService.ResetPasswordViaEmail(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Update password
// @Description Allows an authenticated user to change their own password.
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param updatePassword body models.UpdatePasswordRequest true "Update password request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/update-password [post]
func (c *AuthController) UpdatePassword(ctx *gin.Context) {
	// This endpoint should also be protected by middleware.
	// The user ID should be extracted from the JWT claims, not a request body/param.
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "invalid token"})
		return
	}

	var request models.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := c.authService.UpdatePassword(userID.(string), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Reset password with token
// @Description Allows a user to reset their password using a provided token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Password reset token"
// @Param newPassword body models.ResetPasswordWithTokenRequest true "New password"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
func (c *AuthController) ResetPasswordWithToken(ctx *gin.Context) {
	// 1. Read "token" from query parameters
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "reset token is required"})
		return
	}

	// 2. Bind the new password from the JSON body
	var req models.ResetPasswordWithTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 3. Call the service
	resp, err := c.authService.ResetPasswordWithToken(token, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
