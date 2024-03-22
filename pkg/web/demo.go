package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type demoChangeFormData struct {
	Content string `form:"content"`
}

func (w *WebDiffPreviewer) demoChange(c *gin.Context) {
	data := &demoChangeFormData{}
	err := c.ShouldBind(data)
	if err != nil {
		w.error(c, err)
		return
	}

	w.demoContent = data.Content

	c.Redirect(http.StatusFound, "/demo")
}

func (w *WebDiffPreviewer) demo(c *gin.Context) {
	c.HTML(http.StatusOK, "demo.gohtml", gin.H{"content": w.demoContent})
}
