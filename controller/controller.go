package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"ngetes/helpers"
	"ngetes/model"
)

type User struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Umur        int64      `json:"umur"`
	Address     *string    `json:"address"`
	PhoneNumber *string    `json:"phone_number"`
	Gender      *string    `json:"gender"`
	Status      *string    `json:"status"`
	City        *string    `json:"city"`
	Province    *string    `json:"province"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UserReq struct {
	Name        string `form:"name" validate:"required"`
	Email       string `form:"email" validate:"required,email"`
	Umur        int    `form:"umur" validate:"required,numeric"`
	Password    string `form:"password" validate:"required"`
	Address     string `form:"address" validate:"required"`
	PhoneNumber string `form:"phone_number" validate:"required"`
	Gender      string `form:"gender"`
	Status      string `form:"status"`
	City        string `form:"city" validate:"required"`
	Province    string `form:"province" validate:"required"`
}

type UserDel struct {
	Id []int64 `json:"id"`
}

func GetUsersController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User

		query := `
		SELECT users.id, users.name, users.email, users.umur, detail_users.address, detail_users.phone_number,detail_users.gender,detail_users.status,detail_users.city,detail_users.province,users.created_at,users.updated_at
		FROM users
		LEFT JOIN detail_users
		ON users.id = detail_users.id_user`

		rows, err := db.Query(query)
		if err != nil {
			return err
		}

		for rows.Next() {
			var user User
			var updatedAt sql.NullTime
			err = rows.Scan(
				&user.Id,
				&user.Name,
				&user.Email,
				&user.Umur,
				&user.Address,
				&user.PhoneNumber,
				&user.Status,
				&user.Gender,
				&user.City,
				&user.Province,
				&user.CreatedAt,
				&updatedAt,
			)
			if updatedAt.Valid {
				user.UpdatedAt = &updatedAt.Time
			}
			if err != nil {
				return err
			}
			users = append(users, user)
		}

		user := c.Get("jwt-res").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["id"].(string)
		fmt.Println(name)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": users,
		})

	}

}

func GetUserByIdController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		query := `SELECT users.id, users.name, users.email, users.umur, detail_users.address, detail_users.phone_number,detail_users.gender,detail_users.status,detail_users.city,detail_users.province,users.created_at,users.updated_at
		FROM users
		LEFT JOIN detail_users
		ON users.id = detail_users.id_user
		WHERE users.id = $1`

		var UpdatedAt sql.NullTime
		var user User
		err := db.QueryRowx(query, id).Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Umur,
			&user.Address,
			&user.PhoneNumber,
			&user.Gender,
			&user.Status,
			&user.City,
			&user.Province,
			&user.CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"message": "Data pengguna tidak ditemukan",
				})
			}

		}
		if UpdatedAt.Valid {
			user.UpdatedAt = &UpdatedAt.Time
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    user,
			"message": "Data pengguna ditemukan",
		})
	}
}

func AddUserController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserReq
		var user User
		validate := validator.New()

		err := c.Bind(&req)
		if err != nil {
			return err
		}
		err = validate.Struct(req)
		if err != nil {
			var errorMessage []string
			validationErrors := err.(validator.ValidationErrors)
			for _, err := range validationErrors {
				errorMessage = append(errorMessage, err.Error())
			}
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": errorMessage,
			})
		}

		password, err := helpers.HashPassword(req.Password)
		if err != nil {
			return err
		}

		query := `
		INSERT INTO users (name, email, umur, created_at, password)
		VALUES ($1, $2, $3, now(), $4)  
		RETURNING id, name, email, umur, created_at
		`
		row := db.QueryRowx(query, req.Name, req.Email, req.Umur, password)
		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt)
		if err != nil {
			return err
		}

		query2 := `
		INSERT INTO detail_users (id_user, address, phone_number, gender, status, city, province, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,now())
		RETURNING address, phone_number, gender, status, city, province
		`
		row2 := db.QueryRowx(query2, user.Id, req.Address, req.PhoneNumber, req.Gender, req.Status, req.City, req.Province)
		err = row2.Scan(&user.Address, &user.PhoneNumber, &user.Gender, &user.Status, &user.City, &user.Province)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Berhasil menambahkan data",
			"data":    user,
		})
	}
}

func UpdateUserController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserReq
		var user User
		validate := validator.New()
		id := c.Param("id")

		err := c.Bind(&req)
		if err != nil {
			return err
		}
		err = validate.Struct(req)
		if err != nil {
			var errorMessage []string
			validationErrors := err.(validator.ValidationErrors)
			for _, err := range validationErrors {
				errorMessage = append(errorMessage, err.Error())
			}
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": errorMessage,
			})
		}

		password, err := helpers.HashPassword(req.Password)
		if err != nil {
			return err
		}

		query := `
		UPDATE users
		SET name = $1, email = $2, umur = $3, updated_at = now(),  password = $5 WHERE id = $4
		RETURNING id, name, email, umur, created_at, updated_at
		`
		row := db.QueryRowx(query, req.Name, req.Email, req.Umur, id, password)

		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}
		query2 := `UPDATE detail_users
		SET address = $1, phone_number = $2, gender = $3, status = $4, city = $5, province = $6, updated_at = now() WHERE id_user = $7
		RETURNING address, phone_number, gender, status, city, province
		`

		row2 := db.QueryRowx(query2, req.Address, req.PhoneNumber, req.Gender, req.Status, req.City, req.Province, user.Id)
		err = row2.Scan(&user.Address, &user.PhoneNumber, &user.Gender, &user.Status, &user.City, &user.Province)
		if err != nil {
			return err
		}
		fmt.Println(user)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Berhasil mengedit data",
			"data":    user,
		})
	}
}

func DeleteUsersController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		query := "DELETE FROM users WHERE id = $1"

		_, err := db.Exec(query, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"message": "Data pengguna tidak ditemukan",
				})
			}
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "data berhasil di hapus",
		})
	}
}

func RegisterController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.RegisRequest
		var user model.RegisReponse
		validate := validator.New()

		err := c.Bind(&req)
		if err != nil {
			return err
		}
		err = validate.Struct(req)
		if err != nil {
			var errorMessage []string
			validationErrors := err.(validator.ValidationErrors)
			for _, err := range validationErrors {
				errorMessage = append(errorMessage, err.Error())
			}
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": errorMessage,
			})
		}

		password, err := helpers.HashPassword(req.Password)
		if err != nil {
			return err
		}

		query := `
		INSERT INTO users (name, email, umur, created_at, password)
		VALUES ($1, $2, $3, now(), $4)  
		RETURNING id, name, email, umur, created_at
		`
		row := db.QueryRowx(query, req.Name, req.Email, req.Umur, password)
		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Berhasil menambahkan user",
			"data":    user,
		})
	}
}

func LoginController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.LoginRequest
		var user model.UserModel
		var err error
		validate := validator.New()

		err = c.Bind(&req)
		if err != nil {
			return err
		}
		err = validate.Struct(req)
		if err != nil {
			var errorMessage []string
			validationErrors := err.(validator.ValidationErrors)
			for _, err := range validationErrors {
				errorMessage = append(errorMessage, err.Error())
			}
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": errorMessage,
			})
		}
		query := `SELECT id, name, email, umur, created_at, updated_at, password FROM users WHERE email = $1`
		row := db.QueryRowx(query, req.Email)

		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt, &user.UpdatedAt, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Email tidak terdaftar",
				})
			}
			return err
		}

		match, err := helpers.ComparePassword(user.Password, req.Password)
		if err != nil {
			if !match {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Kata sandi tidak cocok",
				})
			}
			return err
		}

		claims := &model.MyClaims{
			Id: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		}

		sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
		token, err := sign.SignedString([]byte("secret"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token":   token,
			"message": "Login berhasil",
			"data":    user,
		})
	}
}
