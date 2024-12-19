package defines

const (
	// Registro
	SaveUser = `INSERT INTO users (
                  	id_user, name_user, lastname_user, email_user, country_user, phone_user, image_user)
                  	VALUES (?, ?, ?, ?, ?, ?, ?);`
	SaveLogin      = `INSERT INTO login (id_login, username_login, password_login, secret_login, type_login, id_user) VALUES (?, ?, ?, ?, ?, ?);`
	CheckDuplicity = `SELECT CASE WHEN EXISTS (
            SELECT 1 
            FROM login l
            INNER JOIN users u ON l.id_user = u.id_user
            WHERE l.type_login = ?
              AND (u.email_user = ?
                   OR u.phone_user = ?
                   OR l.username_login = ?)
        	) THEN 1
        	ELSE 0
    		END;`
	// Login
	GetLoginByUsername = `SELECT id_login, username_login, password_login, secret_login, id_user, 
							IF(double_login="0", false, true) 
							FROM login 
							WHERE username_login = ?;`
	GetLoginByEmail = `SELECT login.id_login, login.username_login, login.password_login, login.secret_login, users.id_user, 
						IF(login.double_login="0", false, true) 				
       					FROM login JOIN users ON login.id_user = users.id_user 
						WHERE users.email_user = ?;`
	GetLoginByPhone = `SELECT login.id_login, login.username_login, login.password_login, login.secret_login, users.id_user, 
						IF(login.double_login="0", false, true) 				
       					FROM login JOIN users ON login.id_user = users.id_user 
						WHERE users.phone_user = ?;`
	GetSuccessfulLogin = `SELECT users.id_user, users.name_user, login.expiration_login, login.account_login,
       						IF(login.double_login="0", false, true) 
							FROM users JOIN login ON login.id_user = users.id_user 
							WHERE login.id_user = ?;`
	GetSecret = `SELECT login.secret_login FROM login WHERE login.id_user = ?;`
	// CURL
	UpdatePassword   = `UPDATE login SET password_login = ? WHERE id_user = ?;`
	UpdateDoubleAuth = `UPDATE login SET double_login = CASE 
                     		WHEN double_login = 0 THEN 1
                    		WHEN double_login = 1 THEN 0
                   		END
						WHERE id_user = ?;`
	UpdateUserType = `UPDATE login 
					SET account_login  = ?, expiration_login = ?
					WHERE id_user = ?;`
	GetUser = `SELECT name_user, lastname_user, email_user, phone_user, email_user, country_user, image_user 
				FROM users 
				WHERE id_user = ?`
	UpdateUser = `UPDATE users 
					SET name_user = ?, lastname_user = ?, phone_user = ?, email_user = ?, country_user = ?, image_user = ? 
					WHERE id_user = ?;`
)
