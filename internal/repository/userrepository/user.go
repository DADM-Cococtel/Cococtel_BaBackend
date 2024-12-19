package userrepository

import (
	"database/sql"
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"log"
)

type (
	IUser interface {
		SaveUser(*gin.Context, entities.User) error
		SaveLogin(*gin.Context, entities.Login) error
		Exists(string, *string, *string, *string) bool
		GetLoginByUsername(*gin.Context, string) (entities.Login, error)
		GetLoginByEmail(*gin.Context, string) (entities.Login, error)
		GetLoginByPhone(*gin.Context, string) (entities.Login, error)
		GetSuccessfulLogin(*gin.Context, string) (entities.SuccessfulLogin, error)
		UpdatePassword(*gin.Context, string, string) error
		UpdateDoubleAuth(*gin.Context, string) error
		UpdateAccountType(*gin.Context, entities.AccountType) error
		GetUser(*gin.Context, string) (entities.User, error)
		UpdateUser(*gin.Context, entities.User) error
		GetSecret(*gin.Context, string) (*string, error)
	}
	userRepository struct {
		db *sql.DB
	}
)

func NewUserRepository(db *sql.DB) IUser {
	return &userRepository{db: db}
}

func (ur *userRepository) SaveUser(ctx *gin.Context, user entities.User) error {
	stmt, err := ur.db.Prepare(defines.SaveUser)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(
		&user.UserID, &user.Name, &user.Lastname, &user.Email, &user.Country, &user.Phone, &user.Image)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) SaveLogin(ctx *gin.Context, login entities.Login) error {
	stmt, err := ur.db.Prepare(defines.SaveLogin)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(&login.LoginID, &login.Username, &login.Password, &login.Secret, &login.LoginType, &login.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) Exists(loginType string, username, email, phone *string) bool {
	var exists int

	err := ur.db.QueryRow(defines.CheckDuplicity, loginType, email, phone, username).Scan(&exists)
	if err != nil {
		log.Println(err)
		log.Println("Error checking duplicity")
	}

	return exists == 1

}

func (ur *userRepository) GetLoginByUsername(ctx *gin.Context, username string) (entities.Login, error) {
	row := ur.db.QueryRow(defines.GetLoginByUsername, username)
	login := entities.Login{}
	err := row.Scan(&login.LoginID, &login.Username, &login.Password, &login.Secret, &login.UserID, &login.DoubleAuth)
	if err != nil {
		log.Println(err)
		return entities.Login{}, err
	}
	return login, nil
}

func (ur *userRepository) GetLoginByEmail(ctx *gin.Context, email string) (entities.Login, error) {
	row := ur.db.QueryRow(defines.GetLoginByEmail, email)
	var login entities.Login
	err := row.Scan(&login.LoginID, &login.Username, &login.Password, &login.Secret, &login.UserID, &login.DoubleAuth)
	if err != nil {
		log.Println(err)
		return entities.Login{}, err
	}
	return login, nil
}

func (ur *userRepository) GetLoginByPhone(ctx *gin.Context, phone string) (entities.Login, error) {
	row := ur.db.QueryRow(defines.GetLoginByPhone, phone)
	var login entities.Login
	err := row.Scan(&login.LoginID, &login.Username, &login.Password, &login.Secret, &login.UserID, &login.DoubleAuth)
	if err != nil {
		log.Println(err)
		return entities.Login{}, err
	}
	return login, nil
}

func (ur *userRepository) GetSuccessfulLogin(ctx *gin.Context, userID string) (entities.SuccessfulLogin, error) {
	row := ur.db.QueryRow(defines.GetSuccessfulLogin, userID)
	var user entities.SuccessfulLogin
	err := row.Scan(&user.UserID, &user.Name, &user.Expiration, &user.AccountType, &user.DoubleAuth)
	if err != nil {
		log.Println(err)
		return entities.SuccessfulLogin{}, err
	}
	return user, nil
}

func (ur *userRepository) UpdatePassword(ctx *gin.Context, loginID, password string) error {
	stmt, err := ur.db.Prepare(defines.UpdatePassword)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(&password, &loginID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) UpdateDoubleAuth(ctx *gin.Context, loginID string) error {
	stmt, err := ur.db.Prepare(defines.UpdateDoubleAuth)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(&loginID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) UpdateAccountType(ctx *gin.Context, accountType entities.AccountType) error {
	stmt, err := ur.db.Prepare(defines.UpdateUserType)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(&accountType.NewType, &accountType.Expiration, &accountType.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) GetUser(ctx *gin.Context, userID string) (entities.User, error) {
	row := ur.db.QueryRow(defines.GetUser, userID)
	user := entities.User{UserID: userID}
	err := row.Scan(&user.Name, &user.Lastname, &user.Email, &user.Phone, &user.Email, &user.Country, &user.Image)
	if err != nil {
		log.Println(err)
		return entities.User{}, err
	}
	return user, nil
}

func (ur *userRepository) UpdateUser(ctx *gin.Context, user entities.User) error {
	stmt, err := ur.db.Prepare(defines.UpdateUser)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(&user.Name, &user.Lastname, &user.Phone, &user.Email, &user.Country, &user.Image, &user.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ur *userRepository) GetSecret(ctx *gin.Context, userID string) (*string, error) {
	row := ur.db.QueryRow(defines.GetSecret, userID)
	var secret string
	err := row.Scan(&secret)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &secret, nil
}
