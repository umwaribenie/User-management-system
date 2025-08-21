package services

import (
	"errors"
	"math"

	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/repositories"
	"github.com/umwaribenie/final_user_management/utils"
)

type UserService interface {
	GetAllUsers(params models.GetAllUsersRequest) (models.PaginatedResponse, error)
	RegisterUser(request models.CreateUserRequest) (*models.User, error)
	RegisterUserByAdmin(request models.CreateUserByAdminRequest) (*models.User, error)
	GetUserBySlug(slug string) (*models.User, error)
	UpdatePasswordByAdmin(id string, request models.UpdatePasswordRequest) (models.SuccessResponse, error)
	GetUserByID(id string) (*models.User, error)
	DeleteUser(id string) (models.SuccessResponse, error)
	UpdateUser(id string, request models.UpdateUserRequest) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) GetAllUsers(params models.GetAllUsersRequest) (models.PaginatedResponse, error) {
	if params.PageNumber == 0 {
		params.PageNumber = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}

	users, total, err := s.userRepo.FindAll(params)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(params.PageSize)))
	if lastPage == 0 && total > 0 {
		lastPage = 1
	}

	var nextPage, previousPage *int
	if params.PageNumber < lastPage {
		next := params.PageNumber + 1
		nextPage = &next
	}
	if params.PageNumber > 1 {
		prev := params.PageNumber - 1
		previousPage = &prev
	}

	return models.PaginatedResponse{
		CurrentPage:  params.PageNumber,
		LastPage:     lastPage,
		List:         users,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		Status:       "success",
		Total:        total,
	}, nil
}

func (s *userService) RegisterUser(request models.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ClientID:       request.ClientID,
		Email:          request.Email,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		NationalID:     request.NationalID,
		PassportNumber: request.PassportNumber,
		Password:       hashedPassword,
		Phone:          request.Phone,
		ProfilePicture: request.ProfilePicture,
		Username:       request.Username,
		Role:           models.RoleUser,
		Status:         models.ActiveStatus,
		Slug:           utils.GenerateSlug(request.FirstName + " " + request.LastName),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) RegisterUserByAdmin(request models.CreateUserByAdminRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          request.Email,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		NationalID:     request.NationalID,
		PassportNumber: request.PassportNumber,
		Password:       hashedPassword,
		Phone:          request.Phone,
		ProfilePicture: request.ProfilePicture,
		Username:       request.Username,
		Role:           request.Role,
		Status:         models.ActiveStatus,
		Slug:           utils.GenerateSlug(request.FirstName + " " + request.LastName),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserBySlug(slug string) (*models.User, error) {
	return s.userRepo.FindBySlug(slug)
}

func (s *userService) UpdatePasswordByAdmin(id string, request models.UpdatePasswordRequest) (models.SuccessResponse, error) {
	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return models.SuccessResponse{}, err
	}
	if err := s.userRepo.UpdatePassword(id, hashedPassword); err != nil {
		return models.SuccessResponse{}, err
	}
	return models.SuccessResponse{Message: "Password updated successfully"}, nil
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) DeleteUser(id string) (models.SuccessResponse, error) {
	if err := s.userRepo.Delete(id); err != nil {
		return models.SuccessResponse{}, err
	}
	return models.SuccessResponse{Message: "User deleted successfully"}, nil
}

func (s *userService) UpdateUser(id string, request models.UpdateUserRequest) (*models.User, error) {
	// 1. Retrieve the existing user from the database.
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Apply the updates from the request to the existing user object.

	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.LastName = *request.LastName
	}
	if request.NationalID != nil {
		user.NationalID = request.NationalID
	}
	if request.PassportNumber != nil {
		user.PassportNumber = request.PassportNumber
	}
	if request.Phone != nil {
		user.Phone = *request.Phone
	}
	if request.ProfilePicture != nil {
		user.ProfilePicture = request.ProfilePicture
	}
	if request.Role != nil {
		user.Role = *request.Role
	}
	if request.Username != nil {
		user.Username = *request.Username
	}

	// 3. Update the user in the database.
	if err := s.userRepo.Update(id, user); err != nil {
		return nil, err
	}

	// 4. Return the updated user object.
	return user, nil
}
