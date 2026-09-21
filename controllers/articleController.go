package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"QuickGin/forms"
	"QuickGin/services"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	svc *services.ArticleService
}

func NewArticleController() *ArticleController {
	return &ArticleController{svc: services.NewArticleService()}
}

// ListArticles godoc
// @Summary      List published articles
// @Tags         articles
// @Produce      json
// @Param        limit   query     int  false  "Max results (default 20, max 50)"
// @Param        offset  query     int  false  "Pagination offset (default 0)"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]string
// @Router       /articles [get]
func (ctrl *ArticleController) ListArticles(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20, 50)
	offset := parseIntQuery(c, "offset", 0, 0)

	articles, err := ctrl.svc.ListArticles(limit, offset)
	if err != nil {
		ctrl.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

// GetArticle godoc
// @Summary      Get a single article
// @Tags         articles
// @Produce      json
// @Param        id   path      string  true  "Article ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /articles/{id} [get]
func (ctrl *ArticleController) GetArticle(c *gin.Context) {
	article, err := ctrl.svc.GetArticle(c.Param("id"))
	if err != nil {
		ctrl.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// CreateArticle godoc
// @Summary      Create an article
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        body  body      forms.CreateArticleForm  true  "Article details"
// @Success      201   {object}  map[string]interface{}
// @Failure      406   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /articles [post]
func (ctrl *ArticleController) CreateArticle(c *gin.Context) {
	var form forms.CreateArticleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		message := forms.Translate(err, forms.CreateArticleFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	article, err := ctrl.svc.CreateArticle(form)
	if err != nil {
		ctrl.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": article})
}

// UpdateArticle godoc
// @Summary      Update an article
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Article ID"
// @Param        body  body      forms.UpdateArticleForm  true  "Fields to update"
// @Success      200   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]string
// @Failure      406   {object}  map[string]string
// @Router       /articles/{id} [put]
func (ctrl *ArticleController) UpdateArticle(c *gin.Context) {
	var form forms.UpdateArticleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		message := forms.Translate(err, forms.UpdateArticleFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}
	form.ID = c.Param("id")

	article, err := ctrl.svc.UpdateArticle(form)
	if err != nil {
		ctrl.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// DeleteArticle godoc
// @Summary      Delete an article
// @Tags         articles
// @Produce      json
// @Param        id   path      string  true  "Article ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /articles/{id} [delete]
func (ctrl *ArticleController) DeleteArticle(c *gin.Context) {
	if err := ctrl.svc.DeleteArticle(c.Param("id")); err != nil {
		ctrl.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (ctrl *ArticleController) handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrArticleNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "System operation error encountered: " + err.Error()})
}

func parseIntQuery(c *gin.Context, key string, def, max int) int {
	raw := c.Query(key)
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return def
	}
	if max > 0 && n > max {
		return max
	}
	return n
}
