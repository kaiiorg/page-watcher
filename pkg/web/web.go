package web

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kaiiorg/page-watcher/pkg/config"
	"github.com/kaiiorg/page-watcher/pkg/repositories/page_repository"
	"github.com/rs/zerolog/log"
	"html/template"
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

	demoContent string
}

func NewWebDiffPreviewer(c *config.Web, pageRepository page_repository.PageRepository) (*WebDiffPreviewer, error) {
	gin.SetMode(gin.ReleaseMode)
	w := &WebDiffPreviewer{
		config:         c,
		pageRepository: pageRepository,
		gin:            gin.New(),
		demoContent:    uuid.NewString(),
	}

	t := template.New("page-watcher")
	t.Funcs(template.FuncMap{
		"replaceEndLines": w.replaceEndLines,
		"isInsert":        w.isInsert,
		"isDel":           w.isDel,
		"isEqual":         w.isEqual,
	})

	templ, err := t.ParseFS(templatesFS, "templates/*.gohtml")
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
	w.gin.GET("/demo", w.demo)
	w.gin.POST("/demo/change", w.demoChange)

	w.gin.NoRoute(w.noRoute)
	w.gin.Use(static.Serve("/", static.EmbedFolder(staticFS, "static")))

	go func() {
		err := w.gin.Run(fmt.Sprintf(":%d", w.config.Port))
		log.Warn().Err(err).Msg("API has stopped")
	}()
}
