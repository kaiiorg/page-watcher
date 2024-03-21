package web

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (w *WebDiffPreviewer) noRoute(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}

func (w *WebDiffPreviewer) error(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "error.gohtml", gin.H{"error": err.Error()})
}

func (w *WebDiffPreviewer) index(c *gin.Context) {
	pages, err := w.pageRepository.GetDistinctPages()
	if err != nil {
		w.error(c, err)
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
		w.error(c, err)
		return
	}

	latestChange, err := w.pageRepository.GetLatestChange(string(name))
	if err != nil {
		w.error(c, err)
		return
	}

	diff, err := latestChange.DecodeDiff()
	if err != nil {
		w.error(c, err)
		return
	}

	c.HTML(
		http.StatusOK,
		"change.gohtml",
		gin.H{
			"page":  latestChange,
			"lines": strings.Split(latestChange.NormalizedText, "\n"),
			"diff":  diff,
		},
	)
}

func (w *WebDiffPreviewer) specificChange(c *gin.Context) {

}
