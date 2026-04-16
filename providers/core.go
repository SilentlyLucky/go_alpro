package providers

import (
	"github.com/Caknoooo/go-gin-clean-starter/config"
	authController "github.com/Caknoooo/go-gin-clean-starter/modules/auth/controller"
	authRepo "github.com/Caknoooo/go-gin-clean-starter/modules/auth/repository"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	userController "github.com/Caknoooo/go-gin-clean-starter/modules/user/controller"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/repository"
	userService "github.com/Caknoooo/go-gin-clean-starter/modules/user/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (authService.JWTService, error) {
		return authService.NewJWTService(), nil
	})

	db, errDB := do.InvokeNamed[*gorm.DB](injector, constants.DB)
	if errDB != nil {
		panic("GAGAL LOAD DATABASE: " + errDB.Error())
	}

	jwtService, errJWT := do.InvokeNamed[authService.JWTService](injector, constants.JWTService)
	if errJWT != nil {
		panic("GAGAL LOAD JWT: " + errJWT.Error())
	}

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := authRepo.NewRefreshTokenRepository(db)

	userService := userService.NewUserService(userRepository, db)
	authService := authService.NewAuthService(userRepository, refreshTokenRepository, jwtService, db)

	do.Provide(
		injector, func(i *do.Injector) (userController.UserController, error) {
			return userController.NewUserController(i, userService), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (authController.AuthController, error) {
			return authController.NewAuthController(i, authService), nil
		},
	)
}
