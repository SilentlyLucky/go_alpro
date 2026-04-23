package validation

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/Mobilizes/materi-be-alpro/modules/user/dto"
)

func ValidateCreateUser(c *gin.Context) (*dto.CreateUserRequest, error) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        return nil, err
    }
    return &req, nil
}

func ValidateGetUserByID(c *gin.Context) (uint, error) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return 0, err
    }
    return uint(id), nil
}
