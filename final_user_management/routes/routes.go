package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/umwaribenie/final_user_management/controllers"
)

// SetupRouter connects all the user and auth endpoints.
func SetupRouter(
	router *gin.Engine,
	userController *controllers.UserController,
	authController *controllers.AuthController,
) {
	// User routes
	u := router.Group("/users")
	{
		u.GET("/", userController.GetAllUsers)
		u.POST("/register", userController.RegisterUser)
		u.POST("/registerusersbyadmin", userController.RegisterUserByAdmin)
		u.GET("/slug/:slug", userController.GetUserBySlug)
		u.POST("/:id/update-password/admin", userController.UpdatePasswordByAdmin)
		u.GET("/:id", userController.GetUserByID)
		u.DELETE("/:id", userController.DeleteUser)
		u.PATCH("/:id", userController.UpdateUser)
	}

	// Auth routes
	a := router.Group("/auth")
	{
		a.POST("/password-reset", authController.RequestPasswordReset)
		a.POST("/confirm-password-reset-otp", authController.ConfirmPasswordResetOtp)
		a.POST("/login", authController.Login)
		a.POST("/reset-password/email", authController.ResetPasswordViaEmail)
		a.POST("/update-password", authController.UpdatePassword)
		a.POST("/reset-password", authController.ResetPasswordWithToken)
		a.GET("/check", authController.CheckAuth)

	}
}
