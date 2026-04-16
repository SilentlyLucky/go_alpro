package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/modules/user/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/query"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/validation"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/utils"
	"github.com/Caknoooo/go-pagination"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type (
	UserController interface {
		Me(ctx *gin.Context)
		GetAllUser(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GetUserById(ctx *gin.Context)
		GetAllUsers(ctx *gin.Context)
	}

	userController struct {
		userService    service.UserService
		userValidation *validation.UserValidation
		db             *gorm.DB
	}
)

func NewUserController(injector *do.Injector, us service.UserService) UserController {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	userValidation := validation.NewUserValidation()
	return &userController{
		userService:    us,
		userValidation: userValidation,
		db:             db,
	}
}

func (c *userController) GetAllUser(ctx *gin.Context) {
	var filter = &query.UserFilter{}
	filter.BindPagination(ctx)

	ctx.ShouldBindQuery(filter)

	users, total, err := pagination.PaginatedQueryWithIncludable[query.User](c.db, filter)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	paginationResponse := pagination.CalculatePagination(filter.Pagination, total)
	response := pagination.NewPaginatedResponse(http.StatusOK, dto.MESSAGE_SUCCESS_GET_LIST_USER, users, paginationResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Me(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.userService.GetUserById(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (ctrl *userController) GetUserById(c *gin.Context) {
	userId := c.Param("id")

	ctx := c.Request.Context()

	user, err := ctrl.userService.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengambil data", "data": user})
}

func (ctrl *userController) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := ctrl.userService.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengambil semua data", "data": users})
}

func (c *userController) Update(ctx *gin.Context) {
	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.userValidation.ValidateUserUpdateRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.MustGet("user_id").(string)
	result, err := c.userService.Update(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Delete(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	if err := c.userService.Delete(ctx.Request.Context(), userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}
