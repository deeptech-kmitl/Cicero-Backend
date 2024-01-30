package usersRepositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersPattern"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.User, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
	GetProfile(userId string) (*users.User, error)
	DeleteOauth(oauthId string) error
	UpdateProfile(req *users.UserUpdate) error
}

type usersRepository struct {
	db *sqlx.DB
}

func UsersRepositoryHandler(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.User, error) {
	result := usersPattern.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"fname",
		"lname",
		"phone",
		"role_id"
	FROM "User"
	WHERE "email" = $1;`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Oauth" (
		"user_id",
		"access_token"
	)
	VALUES ($1, $2)
	RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"fname",
		"lname",
		"phone",
		"role_id"
	FROM "User"
	WHERE "id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	DELETE FROM "Oauth"
	WHERE "id" = $1;`

	if _, err := r.db.ExecContext(ctx, query, oauthId); err != nil {
		return fmt.Errorf("oauth not found")
	}
	return nil
}

func (r *usersRepository) UpdateProfile(req *users.UserUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	query := `
	UPDATE "User" SET`

	if req.FirstName != "" {
		values = append(values, req.FirstName)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"fname" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.LastName != "" {
		values = append(values, req.LastName)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"lname" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.Email != "" {
		values = append(values, req.Email)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"email" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.Phone != "" {
		values = append(values, req.Phone)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"phone" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.Avatar != "" {
		values = append(values, req.Avatar)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"avatar" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("update profile user failed: %v", err)
	}
	return nil
}
