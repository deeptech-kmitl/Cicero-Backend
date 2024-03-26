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
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserRegisterRes, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
	GetProfile(userId string) (*users.User, error)
	DeleteOauth(oauthId string) error
	UpdateProfile(req *users.UserUpdate) error
	AddWishlist(userId, prodId string) error
	RemoveWishlist(userId, prodId string) error
	CheckWishlist(userId, prodId string) (bool, error)
	FindWishlist(userId string) (*users.WishlistRes, error)
	CheckCart(userId, prodId, size string) (bool, error)
	AddCart(req *users.AddCartReq) error
	AddCartAgain(req *users.AddCartReq) error
	RemoveCart(userId, cartId string) error
	FindCart(userId string) ([]*users.Cart, error)
	DecreaseQtyCart(userId, cartId string) (int, error)
	IncreaseQtyCart(userId, cartId string) (int, error)
	UpdateSizeCart(req *users.UpdateSizeReq) (string, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserRegisterRes, error) {
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
		"dob",
		"avatar",
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
		"role_id",
		"avatar",
		"dob"
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

	if req.Dob != "" {
		values = append(values, req.Dob)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"dob" = $%d?`, lastIndex))

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
			"p"."id",
			"p"."product_title",
			"p"."product_price",
			"p"."product_size",
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

func (r *usersRepository) AddCart(req *users.AddCartReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Cart" (
		"user_id",
		"product_id",
		"size"
	)
	VALUES ($1, $2, $3);
	`

	if _, err := r.db.ExecContext(ctx, query, req.UserId, req.ProductId, req.Size); err != nil {
		switch err.Error() {
		case "ERROR: insert or update on table \"Cart\" violates foreign key constraint \"Cart_product_id_fkey\" (SQLSTATE 23503)":
			return fmt.Errorf("product not found")
		default:
			return fmt.Errorf("add cart failed: %v", err)
		}

	}
	return nil
}

func (r *usersRepository) RemoveCart(userId, cartId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//check user id is same in db
	queryCheck := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "Cart"
	WHERE "id" = $1
	AND "user_id" = $2;
	`

	var check bool
	if err := r.db.Get(&check, queryCheck, cartId, userId); err != nil {
		return fmt.Errorf("check cart failed: %v", err)
	}

	if !check {
		return fmt.Errorf("no permission to remove cart")
	}

	query := `
	DELETE FROM "Cart"
	WHERE "id" = $1;
	`

	if _, err := r.db.ExecContext(ctx, query, cartId); err != nil {
		return fmt.Errorf("remove cart failed: %v", err)
	}
	return nil
}

func (r *usersRepository) CheckCart(userId, prodId, size string) (bool, error) {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "Cart"
	WHERE "user_id" = $1
	AND "product_id" = $2
	AND "size" = $3;`

	var check bool
	if err := r.db.Get(&check, query, userId, prodId, size); err != nil {
		return false, fmt.Errorf("check cart failed: %v", err)
	}
	return check, nil
}

func (r *usersRepository) AddCartAgain(req *users.AddCartReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	UPDATE "Cart"
	SET "qty" = "qty" + 1
	WHERE "user_id" = $1
	AND "product_id" = $2
	AND "size" = $3;
	`
	if _, err := r.db.ExecContext(ctx, query, req.UserId, req.ProductId, req.Size); err != nil {
		return fmt.Errorf("add cart again failed: %v", err)
	}

	return nil
}

func (r *usersRepository) FindCart(userId string) ([]*users.Cart, error) {
	query := `
	SELECT
		COALESCE(array_to_json(array_agg("cart")), '[]'::json)
	FROM (
		SELECT
			"c"."id" AS "cart_id",
			"p"."id",
			"p"."product_title",
			"p"."product_price",
			"p"."product_desc",
			"c"."qty",
			"c"."size",
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
			"Cart" c
		JOIN "Product" p ON c."product_id" = p."id"
		WHERE c."user_id" = $1
	) AS "cart";
	`
	CartBytes := make([]byte, 0)
	if err := r.db.Get(&CartBytes, query, userId); err != nil {
		return nil, fmt.Errorf("get cart failed: %v", err)
	}

	cart := make([]*users.Cart, 0)
	if err := json.Unmarshal(CartBytes, &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %v", err)
	}

	return cart, nil

}

// func (r *usersRepository) GetQtyCart(userId, prodId string) (int, error) {
// 	query := `
// 	SELECT
// 		qty
// 	FROM "Cart"
// 	WHERE "user_id" = $1
// 	AND "product_id" = $2;
// 	`
// 	var qty int
// 	if err := r.db.Get(qty, query, userId, prodId); err != nil {
// 		return 0, fmt.Errorf("cart not found")
// 	}
// 	return qty, nil

// }

func (r *usersRepository) DecreaseQtyCart(userId, cartId string) (int, error) {
	query := `
	UPDATE "Cart"
	SET "qty" = "qty" - 1
	WHERE "user_id" = $1
	AND "id" = $2
	RETURNING "qty";
	`

	fmt.Println("cartId", cartId)

	var qty int
	if err := r.db.Get(&qty, query, userId, cartId); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return 0, fmt.Errorf("no permission to decrease qty cart")
		default:
			return 0, fmt.Errorf("decrease qty cart failed: %v", err)
		}
	}
	return qty, nil

}

func (r *usersRepository) IncreaseQtyCart(userId, cartId string) (int, error) {
	query := `
	UPDATE "Cart"
	SET "qty" = "qty" + 1
	WHERE "user_id" = $1
	AND "id" = $2
	RETURNING "qty";
	`

	var qty int
	if err := r.db.Get(&qty, query, userId, cartId); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return 0, fmt.Errorf("no permission to increase qty cart")
		default:
			return 0, fmt.Errorf("increase qty cart failed: %v", err)
		}
	}
	return qty, nil

}

func (r *usersRepository) UpdateSizeCart(req *users.UpdateSizeReq) (string, error) {
	query := `
	UPDATE "Cart"
	SET "size" = $1
	WHERE "user_id" = $2
	AND "id" = $3
	RETURNING "size";
	`

	var size string
	if err := r.db.Get(&size, query, req.Size, req.UserId, req.CartId); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return "", fmt.Errorf("no permission to update size")
		default:
			return "", fmt.Errorf("update size cart failed: %v", err)
		}
	}
	return size, nil
}
