package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HtmlRenderer struct {}
func NewHtmlRenderer() *HtmlRenderer {
	return &HtmlRenderer{}
}

func (r *HtmlRenderer) RenderLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func (r *HtmlRenderer) SignUp(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

func (r *HtmlRenderer) ForgotPassword(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "reset.html", nil)
}
