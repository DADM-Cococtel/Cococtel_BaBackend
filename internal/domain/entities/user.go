package entities

type (
	User struct {
		UserID   string `json:"user_id"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Email    string `json:"email"`
		Country  string `json:"country"`
		Phone    string `json:"phone"`
		Image    string `json:"image"`
	}

	Login struct {
		LoginID    string `json:"id_login"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Secret     string `json:"secret"`
		DoubleAuth bool   `json:"double_auth"`
		LoginType  string `json:"login_type"`
		UserID     string `json:"user_id"`
	}

	SuccessfulLogin struct {
		UserID      string `json:"id"`
		Name        string `json:"name"`
		DoubleAuth  bool   `json:"double_auth"`
		Expiration  string `json:"expiration"`
		Token       string `json:"token"`
		AccountType string `json:"account_type"`
	}

	AccountType struct {
		UserID     string `json:"user_id"`
		Expiration string `json:"expiration"`
		NewType    string `json:"new_type"`
	}
)
