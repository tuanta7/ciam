package domain

type Account struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FullName  string `json:"full_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
