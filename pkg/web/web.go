package web

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/kaiiorg/page-watcher/pkg/config"
	"github.com/kaiiorg/page-watcher/pkg/repositories/page_repository"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed templates/*
	templatesFS embed.FS
	//go:embed static/*
	staticFS embed.FS
)

type WebDiffPreviewer struct {
	config         *config.Web
	pageRepository page_repository.PageRepository

	gin *gin.Engine
}

func NewWebDiffPreviewer(c *config.Web, pageRepository page_repository.PageRepository) (*WebDiffPreviewer, error) {
	gin.SetMode(gin.ReleaseMode)
	w := &WebDiffPreviewer{
		config:         c,
		pageRepository: pageRepository,
		gin:            gin.New(),
	}

	templ, err := template.ParseFS(templatesFS, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	w.gin.SetHTMLTemplate(templ)

	return w, nil
}

func (w *WebDiffPreviewer) Run() {
	w.gin.Use(gin.LoggerWithFormatter(w.log))
	w.gin.Use(gin.Recovery())
	w.gin.SetTrustedProxies(nil)

	w.gin.GET("/", w.index)
	w.gin.GET("/changes/:base64Name", w.latestChange)
	w.gin.GET("/changes/:base64Name/:pageChange", w.specificChange)

	w.gin.NoRoute(w.noRoute)
	w.gin.Use(static.Serve("/", static.EmbedFolder(staticFS, "static")))

	go func() {
		err := w.gin.Run(fmt.Sprintf(":%d", w.config.Port))
		log.Warn().Err(err).Msg("API has stopped")
	}()
}

func (w *WebDiffPreviewer) log(loggerParams gin.LogFormatterParams) string {
	log.Info().
		Str("clientIP", loggerParams.ClientIP).
		Str("method", loggerParams.Method).
		Str("path", loggerParams.Path).
		Int("status", loggerParams.StatusCode).
		Str("latency", loggerParams.Latency.String()).
		Msg(loggerParams.ErrorMessage)
	return ""
}

func (w *WebDiffPreviewer) noRoute(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}

func (w *WebDiffPreviewer) index(c *gin.Context) {
	pages, err := w.pageRepository.GetDistinctPages()
	if err != nil {
		c.Error(err)
		return
	}

	pageMap := map[string]string{}
	for _, page := range pages {
		pageMap[base64.StdEncoding.EncodeToString([]byte(page))] = page
	}

	c.HTML(http.StatusOK, "index.gohtml", gin.H{"pages": pageMap})
}

func (w *WebDiffPreviewer) latestChange(c *gin.Context) {
	base64Name := c.Param("base64Name")
	name, err := base64.StdEncoding.DecodeString(base64Name)
	if err != nil {
		c.Error(err)
		return
	}

	latestChange, err := w.pageRepository.GetLatestChange(string(name))
	if err != nil {
		c.Error(err)
		return
	}

	diff, err := latestChange.DecodeDiff()
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(
		http.StatusOK,
		"change.gohtml",
		gin.H{
			"page":  latestChange,
			"lines": strings.Split(latestChange.Text, "\n"),
			"diff":  diff,
		},
	)
}

func (w *WebDiffPreviewer) specificChange(c *gin.Context) {

}
