package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bektosh03/monolith/api/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	res, err := h.db.Exec(
		`INSERT INTO users VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT users_email_key DO NOTHING`,
		user.Name, user.Email, user.Password,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.String(http.StatusBadRequest, "already exists")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
	pageValue := c.Query("page")
	if pageValue == "" {
		pageValue = "1"
	}
	limitValue := c.Query("limit")
	if limitValue == "" {
		limitValue = "10"
	}
	page, err := strconv.Atoi(pageValue)
	if err != nil {
		c.String(http.StatusBadRequest, "page value should be numeric")
		return
	}
	limit, err := strconv.Atoi(limitValue)
	if err != nil {
		c.String(http.StatusBadRequest, "page value should be numeric")
		return
	}
	var (
		users  = []models.User{}
		userName sql.NullString
		offset = (page - 1) * limit
	)
	rows, err := h.db.Query(
		`SELECT name, email, password FROM users OFFSET $1 LIMIT $2`,
		offset, limit,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&userName,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		user.Name = userName.String
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	var user models.User
	row := h.db.QueryRow(`SELECT * FROM users WHERE email = $1`, email)
	err := row.Scan(
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		c.String(http.StatusNotFound, "user with this email does not exist")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	email := c.Param("email")
	stmt, err := h.db.Prepare(`DELETE FROM users WHERE email = $1`)
	if err != nil {
		println(err)
		return
	}
	fmt.Println(stmt)
	res, err := stmt.Exec(email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n < 1 {
		c.String(http.StatusNotFound, "not exists")
		return
	}
	c.JSON(
		http.StatusOK,
		struct {
			Status string `json:"status"`
		}{
			Status: "deleted",
		},
	)
}
