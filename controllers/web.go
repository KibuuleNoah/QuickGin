package controllers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/KibuuleNoah/QuickGin/middleware"
	"github.com/gin-gonic/gin"
)

type NavLink struct {
	Path  string
	Label string
	Icon  template.HTML // template.HTML tells Go NOT to escape the SVG string
}

var (
	iconHome = template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
		<path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
	</svg>`)

	iconSettings = template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
		<path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
		<path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
	</svg>`)

	iconUpload = template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
		<path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
	</svg>`)

	iconWallet = template.HTML(`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
		<path stroke-linecap="round" stroke-linejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/>
	</svg>`)
)

// dashboardNavLinks is the single source of truth for what appears in the
// sidebar, top nav, and mobile nav. Add / remove entries here only.
var dashboardNavLinks = []NavLink{
	{Path: "/", Label: "Home", Icon: iconHome},
	{Path: "/settings", Label: "Settings", Icon: iconSettings},
	{Path: "/submissions", Label: "Submissions", Icon: iconUpload},
	{Path: "/wallet", Label: "Wallet", Icon: iconWallet},
}

type WebController struct{}

func NewWebController(router *gin.Engine) *WebController {
	// Group static routes and apply NoCache middleware
	static := router.Group("/dist")
	static.Use(middleware.NoCacheMiddleware())
	{
		static.Static("", "web/dist/")
	}
	// ── HTML templates ──
	// Gin parses all .html files in ./templates/ and makes them available
	// to c.HTML() by filename (e.g. "auth.html", "dashboard.html").
	//
	// We use a custom FuncMap so templates can call {{ safeHTML .Icon }}
	// to render SVG strings without Go escaping the angle brackets.
	router.SetFuncMap(template.FuncMap{
		"safeHTML": func(s template.HTML) template.HTML { return s },
	})
	router.LoadHTMLGlob("web/templates/*")

	// ── Page routes ──
	router.GET("/auth", AuthPage)
	router.GET("/user", UserDashboardPage)

	// ── 404 fallback ──
	router.NoRoute(func(c *gin.Context) {
		// API calls get JSON; page requests get a redirect to auth
		if strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.Redirect(http.StatusFound, "/auth")
	})

	return &WebController{}
}

/*func servePage(template string, extra gin.H) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := gin.H{}
		for k, v := range extra {
			data[k] = v
		}
		c.HTML(http.StatusOK, template, data)
	}
}*/

func AuthPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{})
}

func UserDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"NavLinks": dashboardNavLinks,
	})
}
