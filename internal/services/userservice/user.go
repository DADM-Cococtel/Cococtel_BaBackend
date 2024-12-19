package userservice

import (
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/dtos"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/entities"
	"github.com/Cococtel/Cococtel_BaBackend/internal/repository/userrepository"
	"github.com/Cococtel/Cococtel_BaBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type (
	IUser interface {
		VerifyUser(*gin.Context) utils.ApiError
		RegisterUser(*gin.Context, dtos.Register) (entities.User, utils.ApiError)
		LoginUser(*gin.Context, dtos.Login) (entities.SuccessfulLogin, utils.ApiError)
		ValidateLogin(*gin.Context, dtos.TwoFactorAuth) (entities.SuccessfulLogin, utils.ApiError)
		GetQRDoubleAuth(*gin.Context, dtos.GenerateQR) (string, utils.ApiError)
		NotifyQRRead(*gin.Context, string) utils.ApiError
		//UpdateProfile(*gin.Context, dtos.Profile) utils.ApiError
		//GetUser(*gin.Context, string) (entities.User, utils.ApiError)
		//UpdatePassword(*gin.Context, string) utils.ApiError
		//SendEmailRecoveryPassword(*gin.Context, string) utils.ApiError
		//UpdateUserType(*gin.Context, string, dtos.AccountType) utils.ApiError
	}
	user struct {
		userRepository userrepository.IUser
	}
)

func NewUser(userRepository userrepository.IUser) IUser {
	return &user{
		userRepository: userRepository,
	}
}

func (us *user) VerifyUser(ctx *gin.Context) utils.ApiError {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(defines.ErrInvalidToken, http.StatusUnauthorized)
	}
	_, err = us.userRepository.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(defines.ErrInvalidToken, http.StatusUnauthorized)
	}
	return nil
}

func (us *user) RegisterUser(ctx *gin.Context, register dtos.Register) (entities.User, utils.ApiError) {
	login, usr, err := getUserFromRegister(register)
	if err != nil {
		log.Println(err)
		return entities.User{}, utils.NewApiError(defines.ErrAlreadyExists, http.StatusConflict)
	}
	if us.userRepository.Exists(login.LoginType, &login.Username, &usr.Email, &usr.Phone) {
		return entities.User{}, utils.NewApiError(defines.ErrAlreadyExists, http.StatusConflict)
	}

	err = us.userRepository.SaveLogin(ctx, login)
	if err != nil {
		log.Println(err)
		return entities.User{}, utils.NewApiError(err, http.StatusInternalServerError)
	}

	err = us.userRepository.SaveUser(ctx, usr)
	if err != nil {
		log.Println(err)
		return entities.User{}, utils.NewApiError(err, http.StatusInternalServerError)
	}
	/* TODO: Conectar microservicio de envio de correo, notificaciones, etc.
	err = utils.SendEmail(usr.Email, defines.SuccessfulRegisterTitle, defines.SuccessfulRegisterDescription, defines.LinkToVitalit, defines.SuccessfulRegisterSubject)
	if err != nil {
		log.Println(err)
		log.Print("Error enviando el correo: ", err)
	}
	*/
	return usr, nil
}

func (us *user) LoginUser(ctx *gin.Context, userToLogin dtos.Login) (entities.SuccessfulLogin, utils.ApiError) {
	var usr entities.Login
	var err error
	//TODO: AGREGAR AUTENTICACION CON GOOGLE
	/*
		if userToLogin.Type == "google" {
			usrrr := loginWithGoogle(ctx, userToLogin)
		}
	*/
	switch getDataType(*userToLogin.User) {
	case defines.EmailLoginUser:
		usr, err = us.userRepository.GetLoginByEmail(ctx, *userToLogin.User)
	case defines.PhoneLoginUser:
		usr, err = us.userRepository.GetLoginByPhone(ctx, *userToLogin.User)
	default:
		usr, err = us.userRepository.GetLoginByUsername(ctx, *userToLogin.User)
	}
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, utils.NewApiError(err, http.StatusNotFound)
	}
	err = utils.CompareHashAndPassword(usr.Password, *userToLogin.Password)
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, utils.NewApiError(defines.ErrIncorrectPassword, http.StatusConflict)
	}
	if usr.DoubleAuth {
		return entities.SuccessfulLogin{
			UserID:     usr.UserID,
			DoubleAuth: usr.DoubleAuth,
		}, nil
	}
	logged, apiErr := us.ValidateLogin(ctx, dtos.TwoFactorAuth{User: usr.UserID, Code: "000000"})
	if apiErr != nil {
		return entities.SuccessfulLogin{}, apiErr
	}
	return logged, nil
}

func (us *user) ValidateLogin(ctx *gin.Context, twoFactorAuth dtos.TwoFactorAuth) (entities.SuccessfulLogin, utils.ApiError) {
	usr, err := us.userRepository.GetSuccessfulLogin(ctx, twoFactorAuth.User)
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, utils.NewApiError(err, http.StatusNotFound)
	}
	secret, err := us.userRepository.GetSecret(ctx, usr.UserID)
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, utils.NewApiError(err, http.StatusInternalServerError)
	}
	if usr.DoubleAuth {
		if !verifyOTP(*secret, twoFactorAuth.Code) {
			return entities.SuccessfulLogin{}, utils.NewApiError(defines.ErrInvalidOTP, http.StatusBadRequest)
		}
	}
	usr.Token, err = utils.GenerateJWTToken(*secret, usr.UserID, 72)
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, utils.NewApiError(err, http.StatusInternalServerError)
	}
	parsedDate, err := time.Parse(usr.Expiration, time.Layout)
	if parsedDate.Before(time.Now()) {
		accountType := entities.AccountType{
			UserID:     usr.UserID,
			Expiration: "",
			NewType:    defines.FreeType,
		}
		_ = us.userRepository.UpdateAccountType(ctx, accountType)
	}
	return usr, nil
}

func (us *user) GetQRDoubleAuth(ctx *gin.Context, user dtos.GenerateQR) (string, utils.ApiError) {
	login, err := us.userRepository.GetSuccessfulLogin(ctx, user.User)
	if err != nil {
		log.Println(err)
		return "", utils.NewApiError(err, http.StatusNotFound)
	}
	if login.DoubleAuth {
		return "", utils.NewApiError(defines.ErrAuthAlreadyEnabled, http.StatusConflict)
	}
	secret, err := us.userRepository.GetSecret(ctx, login.UserID)
	if err != nil {
		log.Println(err)
		return "", utils.NewApiError(err, http.StatusInternalServerError)
	}
	uri := utils.GenerateTOTPWithSecret(login.UserID, *secret)
	return uri, nil
}

func (us *user) NotifyQRRead(ctx *gin.Context, userID string) utils.ApiError {
	err := us.userRepository.UpdateDoubleAuth(ctx, userID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	return nil
}

/*
func (us *user) UpdateProfile(ctx *gin.Context, update dtos.Profile) utils.ApiError {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusUnauthorized)
	}
	userToUpdate, err := us.userRepository.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusNotFound)
	}
	loginToUpdate, err := us.userRepository.GetLoginByEmail(ctx, userToUpdate.Email)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	sameEmail, finalUser := updateUser(userToUpdate, update)
	if !sameEmail {
		_, err = us.userRepository.GetLoginByEmail(ctx, finalUser.Email)
		if !errors.Is(err, sql.ErrNoRows) {
			return utils.NewApiError(defines.ErrAlreadyExistsEmail, http.StatusBadRequest)
		}
	}
	if loginToUpdate.Username != update.Username {
		_, err = us.userRepository.GetLoginByUsername(ctx, update.Username)
		if !errors.Is(err, sql.ErrNoRows) {
			return utils.NewApiError(defines.ErrAlreadyExistsUsername, http.StatusBadRequest)
		}
		loginToUpdate.Username = update.Username
	}
	err = us.userRepository.UpdateUser(ctx, finalUser)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	loginToUpdate.FtLogin = update.FtLogin == 1
	err = us.userRepository.UpdateLogin(ctx, loginToUpdate)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	return nil
}

func (us *user) GetUser(ctx *gin.Context, userID string) (entities.User, utils.ApiError) {
	userId, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		log.Println(err)
		return entities.User{}, utils.NewApiError(err, http.StatusUnauthorized)
	}
	if userID != userId {
		return entities.User{}, utils.NewApiError(defines.ErrInvalidToken, http.StatusUnauthorized)
	}
	usr, err := us.userRepository.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return entities.User{}, utils.NewApiError(err, http.StatusNotFound)
	}
	return usr, nil
}

func (us *user) UpdatePassword(ctx *gin.Context, password string) utils.ApiError {
	userID, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusUnauthorized)
	}
	userToUpdate, err := us.userRepository.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusNotFound)
	}
	loginToUpdate, err := us.userRepository.GetLoginByEmail(ctx, userToUpdate.Email)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	passwordHashed, err := utils.GenerateHash(password)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	loginToUpdate.Password = passwordHashed
	err = us.userRepository.UpdateLogin(ctx, loginToUpdate)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	return nil
}

func (us *user) SendEmailRecoveryPassword(ctx *gin.Context, loginRecovery string) utils.ApiError {
	var login entities.Login
	var err error
	if getDataType(loginRecovery) == defines.EmailField {
		login, err = us.userRepository.GetLoginByEmail(ctx, loginRecovery)
	} else {
		login, err = us.userRepository.GetLoginByUsername(ctx, loginRecovery)
	}
	usr, err := us.userRepository.GetUserByLoginID(ctx, login.LoginID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusNotFound)
	}
	token, err := utils.GenerateJWTToken(login.Secret, usr.UserID, defines.RecoveryExpirationTime)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	err = utils.SendEmail(usr.Email, defines.RecoveryPasswordTitle, defines.RecoveryPasswordDescription, defines.RecoveryPasswordPath+token, defines.RecoveryPasswordSubject)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusInternalServerError)
	}
	return nil
}

func (us *user) UpdateUserType(ctx *gin.Context, userID string, accountConfig dtos.AccountType) utils.ApiError {
	timeNow := int(time.Now().Unix())
	userId, err := utils.GetUserIDFromToken(ctx)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusUnauthorized)
	}
	if userID != userId {
		return utils.NewApiError(defines.ErrInvalidToken, http.StatusUnauthorized)
	}
	userToUpdate, err := us.userRepository.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(err, http.StatusNotFound)
	}
	userUpdate := entities.AccountType{
		LoginID:  userToUpdate.Login,
		Duration: defines.DurationDay + timeNow,
		NewType:  defines.PremiumType,
	}
	if accountConfig.Duration < 1 {
		err = us.userRepository.UpdateAccountType(ctx, userUpdate)
		if err != nil {
			log.Println(err)
			return utils.NewApiError(defines.ErrCantUpdateAccountType, http.StatusInternalServerError)
		}
		return nil
	}

	if accountConfig.ReferralCode != "null" && accountConfig.ReferralCode != "" {
		referralUser, err := us.userRepository.GetLoginByUsername(ctx, accountConfig.ReferralCode)
		if err != nil {
			log.Println(err)
		}
		referralUserUpdate := entities.AccountType{
			LoginID:  referralUser.LoginID,
			Duration: defines.DurationMonth + timeNow,
			NewType:  defines.PremiumType,
		}
		err = us.userRepository.UpdateAccountType(ctx, referralUserUpdate)
		if err != nil {
			log.Println(err)
		}
	}
	userUpdate.Duration = timeNow + defines.DurationMonth*accountConfig.Duration

	err = us.userRepository.UpdateAccountType(ctx, userUpdate)
	if err != nil {
		log.Println(err)
		return utils.NewApiError(defines.ErrCantUpdateAccountType, http.StatusInternalServerError)
	}
	return nil
}
*/
