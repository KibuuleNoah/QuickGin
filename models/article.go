package models

import (
	"errors"
	"fmt"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/forms"
)

// Response represents the top-level structure
type AllArticleResponse struct {
	Results []Result `json:"results"`
}

type OneArticleResponse struct {
	Data ArticleResponse `json:"data"`
}

// Result represents each result item
type Result struct {
	Data []ArticleResponse `json:"data"`
	Meta Meta              `json:"meta"`
}

// Article represents an individual article
type ArticleResponse struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	UpdatedAt int64        `json:"updatedAt"`
	CreatedAt int64        `json:"createdAt"`
	User      UserResponse `json:"user"`
}

// User represents the article's author
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Meta represents metadata for pagination or total count
type Meta struct {
	Total int `json:"total"`
}

// Article ...
type Article struct {
	ID        int64    `db:"id" json:"id"`
	UserID    int64    `db:"user_id" json:"-"`
	Title     string   `db:"title" json:"title"`
	Content   string   `db:"content" json:"content"`
	UpdatedAt int64    `db:"updated_at" json:"updatedAt"`
	CreatedAt int64    `db:"created_at" json:"createdAt"`
	User      *JSONRaw `db:"user" json:"user"`
}

// ArticleModel ...
type ArticleModel struct{}

// Create ...
func (m ArticleModel) Create(userID string, form forms.CreateArticleForm) (articleID int64, err error) {
	err = db.AppDB().QueryRow("INSERT INTO public.article(user_id, title, content) VALUES($1, $2, $3) RETURNING id", userID, form.Title, form.Content).Scan(&articleID)
	return articleID, err
}

// One ...
func (m ArticleModel) One(userID, id string) (article Article, err error) {
	err = db.AppDB().Get(&article, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.article a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 AND a.id=$2 LIMIT 1", userID, id)
	return article, err
}

// All ...
func (m ArticleModel) All(userID string) (articles []DataList, err error) {
	// err = db.AppDB().Select(&articles, "SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') AS data, (SELECT row_to_json(n) FROM ( SELECT count(a.id) AS total FROM public.article AS a WHERE a.user_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.article a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 ORDER BY a.id DESC) d", userID)

	MockManyDataLists := func(count int) []DataList {
		list := make([]DataList, count)
		for i := 0; i < count; i++ {
			list[i] = DataList{
				Data: JSONRaw(fmt.Sprintf(`{"id": %d, "name": "Item %d"}`, i+1, i+1)),
				Meta: JSONRaw(fmt.Sprintf(`{"index": %d, "total": %d}`, i, count)),
			}
		}
		return list
	}

	// Usage
	articles = MockManyDataLists(5)
	return articles, err
}

// Update ...
func (m ArticleModel) Update(userID string, id string, form forms.CreateArticleForm) (err error) {
	//METHOD 1
	//Check the article by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	operation, err := db.AppDB().Exec("UPDATE public.article SET title=$3, content=$4 WHERE id=$1 AND user_id=$2", id, userID, form.Title, form.Content)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("updated 0 records")
	}

	return err
}

// Delete ...
func (m ArticleModel) Delete(userID, id string) (err error) {

	operation, err := db.AppDB().Exec("DELETE FROM public.article WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
