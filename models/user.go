package models

// Role defines user permissions within an organization
type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
)

// User ...
type User struct {
	ID          string `db:"id" json:"id"`
	PhoneNo     string `db:"phone_no" json:"phoneNo"`
	Verified    bool   `db:"verified" json:"verified"`
	Username    string `db:"username" json:"username"`
	AvatarUrl   string `db:"avatar_url" json:"avatarUrl"`
	Role        *Role  `db:"role" json:"role"`
	UpdatedAt   int64  `db:"updated_at" json:"-"`
	CreatedAt   int64  `db:"created_at" json:"-"`
	LastLoginAt int64  `db:"last_login_at" json:"lastLoginAt,omitempty"`
}

// Fetch One User of the currect struct by id
// same as  err := DB.Get(&user, "SELECT id, phone_no, name, updated_at, created_at FROM public.user WHERE id=$1 LIMIT 1", form.UserID)
// func (u *User) () (err error) {
// 	err = db.AppDB().Get(u, "SELECT id, phone_no, name, updated_at, created_at FROM public.user WHERE id=$1 LIMIT 1", u.ID)
// d	return err
// }
