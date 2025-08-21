package controllers

import (
	"net/http"

	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService}
}

// @Summary Get all users
// @Description Retrieves a list of all users, with support for pagination, searching, and filtering.
// @Tags users
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of users per page" default(10)
// @Param from query string false "Start date for user creation (YYYY-MM-DD)"
// @Param to query string false "End date for user creation (YYYY-MM-DD)"
// @Param search query string false "Search term for user details (first name, last name, email, username)"
// @Param role query string false "Filter by user role" Enums(user, admin)
// @Param status query string false "Filter by user status" Enums(active, inactive, deleted)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	var request models.GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := c.userService.GetAllUsers(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Register a new user
// @Description Creates a new user account with a default 'user' role.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/register [post]
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var request models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := c.userService.RegisterUser(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// @Summary Register a new user by admin
// @Description Creates a new user account with a specified role.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserByAdminRequest true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/registerusersbyadmin [post]
func (c *UserController) RegisterUserByAdmin(ctx *gin.Context) {
	var request models.CreateUserByAdminRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := c.userService.RegisterUserByAdmin(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// @Summary Find a user by slug
// @Description Retrieves a user's details using their URL-friendly slug.
// @Tags users
// @Accept json
// @Produce json
// @Param slug path string true "User Slug"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /users/slug/{slug} [get]
func (c *UserController) GetUserBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	user, err := c.userService.GetUserBySlug(slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Update password by admin
// @Description Allows an admin to update a user's password.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param updatePassword body models.UpdatePasswordRequest true "Update password request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /users/{id}/update-password/admin [post]
func (c *UserController) UpdatePasswordByAdmin(ctx *gin.Context) {
	userID := ctx.Param("id")
	var request models.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	response, err := c.userService.UpdatePasswordByAdmin(userID, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Get a user by ID
// @Description Retrieves a user's details using their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Delete a user
// @Description Deletes a user by their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	response, err := c.userService.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Update a user
// @Description Updates a user's details by their unique ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User data to update"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [patch]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var request models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := c.userService.UpdateUser(id, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
