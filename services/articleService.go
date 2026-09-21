package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"QuickGin/db"
	"QuickGin/forms"
	"QuickGin/models"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
)

type ArticleServiceConfig struct {
	DB *sqlx.DB
}

type ArticleService struct {
	cfg ArticleServiceConfig
}

func NewArticleService() *ArticleService {
	return &ArticleService{cfg: ArticleServiceConfig{DB: db.AppDB()}}
}

var ErrArticleNotFound = errors.New("article record not found")

func (s *ArticleService) ListArticles(limit, offset int) ([]models.Article, error) {
	var articles []models.Article
	err := s.cfg.DB.Select(&articles, `
		SELECT id, title, slug, body, author_id, published, created_at, updated_at
		FROM articles
		WHERE published = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed fetching articles: %w", err)
	}
	return articles, nil
}

func (s *ArticleService) GetArticle(id string) (*models.Article, error) {
	var article models.Article
	err := s.cfg.DB.Get(&article, `
		SELECT id, title, slug, body, author_id, published, created_at, updated_at
		FROM articles WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("failed fetching article: %w", err)
	}
	return &article, nil
}

func (s *ArticleService) CreateArticle(form forms.CreateArticleForm) (*models.Article, error) {
	articleSlug := slug.Make(form.Title)

	var article models.Article
	err := s.cfg.DB.QueryRowx(`
		INSERT INTO articles (title, slug, body, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, slug, body, author_id, published, created_at, updated_at`,
		form.Title, articleSlug, form.Body, form.AuthorID,
	).StructScan(&article)
	if err != nil {
		return nil, fmt.Errorf("failed inserting article record: %w", err)
	}

	return &article, nil
}

func (s *ArticleService) UpdateArticle(form forms.UpdateArticleForm) (*models.Article, error) {
	sets := []string{}
	args := []interface{}{}
	i := 1

	if form.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", i))
		args = append(args, *form.Title)
		i++
		sets = append(sets, fmt.Sprintf("slug = $%d", i))
		args = append(args, slug.Make(*form.Title))
		i++
	}
	if form.Body != nil {
		sets = append(sets, fmt.Sprintf("body = $%d", i))
		args = append(args, *form.Body)
		i++
	}
	if form.Published != nil {
		sets = append(sets, fmt.Sprintf("published = $%d", i))
		args = append(args, *form.Published)
		i++
	}

	if len(sets) == 0 {
		return s.GetArticle(form.ID)
	}

	sets = append(sets, "updated_at = now()")
	args = append(args, form.ID)

	query := fmt.Sprintf(`
		UPDATE articles SET %s
		WHERE id = $%d
		RETURNING id, title, slug, body, author_id, published, created_at, updated_at`,
		strings.Join(sets, ", "), i)

	var article models.Article
	err := s.cfg.DB.QueryRowx(query, args...).StructScan(&article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("failed updating article record: %w", err)
	}

	return &article, nil
}

func (s *ArticleService) DeleteArticle(id string) error {
	res, err := s.cfg.DB.Exec("DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed deleting article record: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed analyzing deletion statement response: %w", err)
	}

	if rows == 0 {
		return ErrArticleNotFound
	}

	return nil
}
