package usersRepositories

import (
	"context"
	"encoding/json"
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
	AddWishlist(userId, prodId string) error
	RemoveWishlist(userId, prodId string) error
	CheckWishlist(userId, prodId string) (bool, error)
	FindWishlist(userId string) (*users.WishlistRes, error)
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

func (r *usersRepository) AddWishlist(userId, prodId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Wishlist" (
		"user_id",
		"product_id"
	)
	VALUES ($1, $2);
	`

	if _, err := r.db.ExecContext(ctx, query, userId, prodId); err != nil {
		switch err.Error() {
		case "ERROR: insert or update on table \"Wishlist\" violates foreign key constraint \"Wishlist_product_id_fkey\" (SQLSTATE 23503)":
			return fmt.Errorf("product not found")
		default:
			return fmt.Errorf("add wishlist failed: %v", err)
		}

	}
	return nil

}

func (r *usersRepository) RemoveWishlist(userId, prodId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	DELETE FROM "Wishlist"
	WHERE "user_id" = $1
	AND "product_id" = $2
	`

	if _, err := r.db.ExecContext(ctx, query, userId, prodId); err != nil {
		return fmt.Errorf("remove wishlist failed: %v", err)
	}
	return nil
}

func (r *usersRepository) CheckWishlist(userId, prodId string) (bool, error) {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "Wishlist"
	WHERE "user_id" = $1
	AND "product_id" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, prodId); err != nil {
		return false, fmt.Errorf("check wishlist failed: %v", err)
	}
	return check, nil
}

func (r *usersRepository) FindWishlist(userId string) (*users.WishlistRes, error) {
	query := `
	SELECT
		COALESCE(array_to_json(array_agg("wishlist")), '[]'::json)
	FROM (
		SELECT
			"p"."product_title",
			"p"."product_price",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "Image" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM
			"Wishlist" wl
		JOIN "Product" p ON wl."product_id" = p."id"
		WHERE wl."user_id" = $1
	) AS "wishlist";
	`
	WishlistBytes := make([]byte, 0)
	if err := r.db.Get(&WishlistBytes, query, userId); err != nil {
		return nil, fmt.Errorf("get wishlist failed: %v", err)
	}

	wishList := make(users.WishlistRes, 0)
	if err := json.Unmarshal(WishlistBytes, &wishList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wishlist: %v", err)
	}

	return &wishList, nil

}
