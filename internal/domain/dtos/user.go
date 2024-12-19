package dtos

import "github.com/golang-jwt/jwt/v5"

type (
	Register struct {
		Name     *string `json:"name"`
		Lastname *string `json:"lastname,omitempty"`
		Phone    *string `json:"phone,omitempty"`
		Email    *string `json:"email,omitempty"`
		Country  *string `json:"country,omitempty"`
		Image    *string `json:"image,omitempty"`
		Username *string `json:"username,omitempty"`
		Password *string `json:"password,omitempty"`
		Type     *string `json:"type,omitempty"`
	}

	Login struct {
		User     *string `json:"user,omitempty"`
		Password *string `json:"password,omitempty"`
		Type     *string `json:"type,omitempty"`
		//AdditionalData     *string `json:"user,omitempty"`
	}

	TwoFactorAuth struct {
		User string `json:"user"`
		Code string `json:"code"`
	}

	JwtCustomClaims struct {
		User   string `json:"user"`
		Secret string `json:"secret"`
		jwt.RegisteredClaims
	}

	GenerateQR struct {
		User string `json:"user"`
	}
)
