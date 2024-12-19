package userservice

import (
	"fmt"
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/dtos"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/entities"
	"github.com/Cococtel/Cococtel_BaBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/xlzd/gotp"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func getUserFromRegister(register dtos.Register) (entities.Login, entities.User, error) {
	var password string
	var err error
	if *register.Type == defines.NormalLoginType && register.Password != nil {
		password, err = utils.GenerateHash(*register.Password)
		if err != nil {
			log.Println(err)
			return entities.Login{}, entities.User{}, err
		}
	}

	user := entities.User{
		UserID:   utils.GenerateUUID(),
		Name:     *register.Name,
		Lastname: *register.Lastname,
		Email:    *register.Email,
		Country:  *register.Country,
		Phone:    *register.Phone,
		Image:    *register.Image,
	}

	login := entities.Login{
		LoginID:   utils.GenerateUUID(),
		Username:  *register.Username,
		Password:  password,
		Secret:    gotp.RandomSecret(16),
		LoginType: *register.Type,
		UserID:    user.UserID,
	}

	return login, user, nil
}

func getDataType(data string) string {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	phoneRegex := regexp.MustCompile(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`)

	if emailRegex.MatchString(data) {
		return "email"
	} else if phoneRegex.MatchString(data) {
		return "phone"
	}
	return "username"
}

func loginWithGoogle(ctx *gin.Context, login dtos.Login) {
	key := *login.User   // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("our-google-client-id", "our-google-client-secret", "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}

func verifyOTP(secret string, code string) bool {
	totp := gotp.NewDefaultTOTP(secret)
	if totp.Verify(code, time.Now().Unix()) {
		return true
	} else {
		return false
	}
}
