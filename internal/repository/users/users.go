package users

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"Assignment5/pkg/modules"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(filter modules.UserFilter) (modules.PaginatedResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 5
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	allowedOrderBy := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"gender":     "gender",
		"birth_date": "birth_date",
	}

	orderBy := "id"
	if val, ok := allowedOrderBy[filter.OrderBy]; ok {
		orderBy = val
	}

	where := []string{"1=1"}
	args := []interface{}{}
	argID := 1

	if filter.ID != nil {
		where = append(where, fmt.Sprintf("id = $%d", argID))
		args = append(args, *filter.ID)
		argID++
	}

	if filter.Name != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+filter.Name+"%")
		argID++
	}

	if filter.Email != "" {
		where = append(where, fmt.Sprintf("email ILIKE $%d", argID))
		args = append(args, "%"+filter.Email+"%")
		argID++
	}

	if filter.Gender != "" {
		where = append(where, fmt.Sprintf("gender = $%d", argID))
		args = append(args, filter.Gender)
		argID++
	}

	if filter.BirthDate != "" {
		where = append(where, fmt.Sprintf("birth_date = $%d", argID))
		args = append(args, filter.BirthDate)
		argID++
	}

	whereSQL := strings.Join(where, " AND ")

	var totalCount int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE %s`, whereSQL)
	err := r.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return modules.PaginatedResponse{}, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date, created_at, updated_at
		FROM users
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argID, argID+1)

	args = append(args, filter.Limit, filter.Offset)

	var users []modules.User
	err = r.db.Select(&users, query, args...)
	if err != nil {
		return modules.PaginatedResponse{}, err
	}

	return modules.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}, nil
}

func (r *Repository) GetCommonFriends(u1, u2 int) ([]modules.User, error) {

	query := `
	SELECT u.*
	FROM users u
	JOIN user_friends f1 ON u.id=f1.friend_id
	JOIN user_friends f2 ON u.id=f2.friend_id
	WHERE f1.user_id=$1 AND f2.user_id=$2
	`

	var users []modules.User
	err := r.db.Select(&users, query, u1, u2)

	return users, err
}
