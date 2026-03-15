package modules

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Gender    string    `db:"gender" json:"gender"`
	BirthDate time.Time `db:"birth_date" json:"birth_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserFilter struct {
	ID        *int
	Name      string
	Email     string
	Gender    string
	BirthDate string
	OrderBy   string
	Limit     int
	Offset    int
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type CommonFriendsResponse struct {
	User1         int    `json:"user1"`
	User2         int    `json:"user2"`
	CommonFriends []User `json:"commonFriends"`
	Count         int    `json:"count"`
}
