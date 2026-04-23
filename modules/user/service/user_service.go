package service

import (
    "github.com/Mobilizes/materi-be-alpro/database/entities"
    "github.com/Mobilizes/materi-be-alpro/modules/user/dto"
    "github.com/Mobilizes/materi-be-alpro/modules/user/repository"
    "github.com/Mobilizes/materi-be-alpro/pkg/helpers"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    hashedPassword, err := helpers.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    user := &entities.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    err = s.repo.Create(user)
    if err != nil {
        return nil, err
    }

    return mapUserResponse(user), nil
}

func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    return mapUserResponse(user), nil
}

func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
    users, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    responses := make([]dto.UserResponse, 0, len(users))
    for i := range users {
        responses = append(responses, *mapUserResponse(&users[i]))
    }

    return responses, nil
}

func mapUserResponse(user *entities.User) *dto.UserResponse {
    return &dto.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    }
}
