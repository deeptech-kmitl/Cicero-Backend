package usersPattern

import (
	"context"
	"fmt"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*users.UserRegisterRes, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser {
	if !isAdmin {
		return newCustomer(db, req)
	}
	return newAdmin(db, req)
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}

}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "User" (
		email,
		password,
		fname,
		lname,
		phone,
		dob,
		role_id
		)
	VALUES ($1, $2, $3, $4, $5, $6, 1)
	RETURNING "id";
	`
	if err := f.db.QueryRowContext(ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.FirstName,
		f.req.LastName,
		f.req.Phone,
		f.req.Dob,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"User_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		case "ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("phone number has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil
}

func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "User" (
		email,
		password,
		fname,
		lname,
		phone,
		dob,
		role_id
		)
	VALUES ($1, $2, $3, $4, $5, $6, 2)
	RETURNING "id";
	`
	if err := f.db.QueryRowContext(ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.FirstName,
		f.req.LastName,
		f.req.Phone,
		f.req.Dob,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"User_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		case "ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("phone number has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil
}

func (f *userReq) Result() (*users.UserRegisterRes, error) {
	query := `
	SELECT
		"u"."id",
		"u"."email",
		"u"."fname",
		"u"."lname",
		"u"."phone",
		"u"."dob",
		"u"."role_id"
	FROM "User" "u"
	WHERE "u"."id" = $1;`

	user := new(users.UserRegisterRes)
	if err := f.db.Get(user, query, f.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}
