package usercontroller

import (
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/dtos"
	"github.com/Cococtel/Cococtel_BaBackend/internal/services/userservice"
	"github.com/Cococtel/Cococtel_BaBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
)

type (
	IUser interface {
		VerifyUser() gin.HandlerFunc
		RegisterUser() gin.HandlerFunc
		LoginUser() gin.HandlerFunc
		ValidateLogin() gin.HandlerFunc
		GetQRDoubleAuth() gin.HandlerFunc
		NotifyQRRead() gin.HandlerFunc
		//EditProfile() gin.HandlerFunc
		//GetUser() gin.HandlerFunc
		//UpdatePassword() gin.HandlerFunc
		//SendEmailToRecoveryPassword() gin.HandlerFunc
		//UpdateUserType() gin.HandlerFunc
	}
	user struct {
		userService userservice.IUser
	}
)

func NewUser(userService userservice.IUser) IUser {
	return &user{
		userService: userService,
	}
}

func (u *user) RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userToRegister dtos.Register
		er := ctx.ShouldBindBodyWith(&userToRegister, binding.JSON)
		if er != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		usr, err := u.userService.RegisterUser(ctx, userToRegister)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusCreated, usr)
	}
}

func (u *user) VerifyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := u.userService.VerifyUser(ctx)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, "ok")
	}
}

func (u *user) LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userToLogin dtos.Login
		if err := ctx.ShouldBindBodyWith(&userToLogin, binding.JSON); err != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		usrToken, err := u.userService.LoginUser(ctx, userToLogin)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, usrToken)
	}
}

func (u *user) ValidateLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var twoFactorAuth dtos.TwoFactorAuth
		if err := ctx.ShouldBindBodyWith(&twoFactorAuth, binding.JSON); err != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		usr, err := u.userService.ValidateLogin(ctx, twoFactorAuth)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, usr)
	}
}

func (u *user) GetQRDoubleAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var generateQR dtos.GenerateQR
		if err := ctx.ShouldBindBodyWith(&generateQR, binding.JSON); err != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		qr, err := u.userService.GetQRDoubleAuth(ctx, generateQR)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, qr)
	}
}

func (u *user) NotifyQRRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := u.userService.NotifyQRRead(ctx, id)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, "ok")
	}
}

/*
func (u *user) EditProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var profile dtos.Profile
		if err := ctx.ShouldBindBodyWith(&profile, binding.JSON); err != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		err := u.userService.UpdateProfile(ctx, profile)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, defines.Ok)
	}
}

func (u *user) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		usr, err := u.userService.GetUser(ctx, id)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, usr)
	}
}

func (u *user) UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var recoveryPassword dtos.RecoveryPassword
		er := ctx.ShouldBindBodyWith(&recoveryPassword, binding.JSON)
		if er != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		err := u.userService.UpdatePassword(ctx, recoveryPassword.Password)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusCreated, defines.Ok)
	}
}

func (u *user) SendEmailToRecoveryPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var recoveryPassword dtos.LoginPasswordRecovery
		er := ctx.ShouldBindBodyWith(&recoveryPassword, binding.JSON)
		if er != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		err := u.userService.SendEmailRecoveryPassword(ctx, recoveryPassword.Login)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusCreated, defines.Ok)
	}
}

func (u *user) UpdateUserType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var accountConfig dtos.AccountType
		er := ctx.ShouldBindBodyWith(&accountConfig, binding.JSON)
		if er != nil {
			utils.Error(ctx, http.StatusBadRequest, defines.ErrShouldBindJSON.Error())
			return
		}
		err := u.userService.UpdateUserType(ctx, id, accountConfig)
		if err != nil {
			log.Println(err)
			utils.Error(ctx, err.Status(), err.Message().Error())
			return
		}
		utils.Success(ctx, http.StatusOK, defines.Ok)
	}
}
*/
