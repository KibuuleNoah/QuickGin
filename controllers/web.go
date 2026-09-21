package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pajo/db"

	"github.com/gin-gonic/gin"
)

type ShareItemMetaData struct {
	Title       string `db:"name"`          // Maps from items.name
	Description string `db:"description"`   // Maps from item_details.description
	Image       string `db:"thumbnail_url"` // Maps from items.thumbnail_url
	URL         string `db:"-"`             // Ignored by the DB; set dynamically in Go code
}

type WebController struct{}

func NewWebController(router *gin.Engine) *WebController {
	router.SetFuncMap(template.FuncMap{
		"safeHTML": func(s template.HTML) template.HTML { return s },
	})
	router.LoadHTMLGlob("web/templates/*")

	// ── Page routes ──
	router.GET("/share/item/:itemID", ShareItem)

	return &WebController{}
}

func ShareItem(c *gin.Context) {
	itemID := c.Param("itemID")

	var meta ShareItemMetaData

	query := `
		SELECT 
			i.name, 
			d.description, 
			i.thumbnail_url 
		FROM items i
		INNER JOIN item_details d ON i.id = d.item_id 
		WHERE i.id = $1
	`

	err := db.AppDB().Get(&meta, query, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		log.Println("ShareItem: failed fetching item metadata:", err)
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	meta.URL = fmt.Sprintf("%s/items/%s", "https://pajo.ug", itemID)

	c.HTML(http.StatusOK, "shareItem.html", meta)
}
