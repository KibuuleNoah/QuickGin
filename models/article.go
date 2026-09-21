package models

// Article ...
type Article struct {
	ID        string `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Slug      string `db:"slug" json:"slug"`
	Body      string `db:"body" json:"body"`
	AuthorID  string `db:"author_id" json:"authorId"`
	Published bool   `db:"published" json:"published"`
	CreatedAt int64  `db:"created_at" json:"createdAt"`
	UpdatedAt int64  `db:"updated_at" json:"updatedAt"`
}
