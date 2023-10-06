package controller

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Umur      int64      `json:"umur"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type UserReq struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Umur  int    `form:"umur"`
}

func GetUsersController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User

		query := "SELECT id,name,email,umur,created_at,updated_at FROM users"

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

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": users,
		})
	}
}

func GetUserByIdController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		query := "SELECT id, name, email, umur,created_at,updated_at FROM users WHERE id = $1"
		var UpdatedAt sql.NullTime
		var user User
		err := db.QueryRowx(query, id).Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Umur,
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
			"data": user,
		})
	}
}

func AddUserController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserReq
		var user User

		err := c.Bind(&req)
		if err != nil {
			return err
		}

		query := `
		INSERT INTO users (name, email, umur, created_at)
		VALUES ($1, $2, $3, now())  
		RETURNING id, name, email, umur, created_at
		`
		row := db.QueryRowx(query, req.Name, req.Email, req.Umur)

		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Berhasil mengedit data",
			"data":    user,
		})
	}
}

func UpdateUserController(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserReq
		var user User
		id := c.Param("id")

		err := c.Bind(&req)
		if err != nil {
			return err
		}

		query := `
		UPDATE users
		SET name = $1, email = $2, umur = $3, updated_at = now() WHERE id = $4 
		RETURNING id, name, email, umur, created_at, updated_at
		`
		row := db.QueryRowx(query, req.Name, req.Email, req.Umur, id)

		err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Umur, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}

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
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "data berhasil di hapus",
		})
	}
}
