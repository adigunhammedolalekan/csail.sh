package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saas/hostgolang/pkg/repository"
	"github.com/saas/hostgolang/pkg/session"
	"github.com/saas/hostgolang/pkg/types"
	"log"
	"net/http"
)

type DeploymentHandler struct {
	repo repository.DeploymentRepository
	appRepo repository.AppsRepository
	store session.Store
}

func NewDeploymentHandler(repo repository.DeploymentRepository, appRepo repository.AppsRepository, s session.Store) *DeploymentHandler {
	return &DeploymentHandler{repo:repo, store: s, appRepo: appRepo}
}

func (handler *DeploymentHandler) CreateDeploymentHandler(ctx *gin.Context) {
	account := handler.ensureAccount(ctx)
	if account == nil {
		return
	}
	file, err := ctx.FormFile("bin")
	appName := ctx.PostForm("app_name")
	if appName == "" {
		BadRequestResponse(ctx, "bad request: app name is missing")
		return
	}
	if err != nil {
		BadRequestResponse(ctx, "deployment binary file is missing")
		return
	}
	app, err := handler.appRepo.GetApp(appName)
	if err != nil {
		BadRequestResponse(ctx, fmt.Sprintf("application not found: %s", appName))
		return
	}
	if app.AccountId != account.ID {
		ForbiddenRequestResponse(ctx, "this app does not belong to you")
		return
	}
	fi, err := file.Open()
	if err != nil {
		log.Println("failed to open the attached bin: ", err)
		BadRequestResponse(ctx, "deployment binary file is missing")
		return
	}
	r, err := handler.repo.CreateDeployment(app, fi)
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "deployment created", Data: r})
}

func (handler *DeploymentHandler) UpdateEnvironmentVars(ctx *gin.Context) {
	account := handler.ensureAccount(ctx)
	if account == nil {
		return
	}
	appName := ctx.Param("appName")
	if appName == "" {
		appName = ctx.Query("appName")
	}
	if appName == "" {
		BadRequestResponse(ctx, "application name is missing")
		return
	}
	var request struct{
		Key string `json:"key"`
		Value string `json:"value"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequestResponse(ctx, "invalid or malformed request body")
		return
	}
	err := handler.appRepo.UpdateEnvironmentVars(appName, map[string]string{request.Key: request.Value})
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "environment config vars updated"})
}

func (handler *DeploymentHandler) GetEnvironmentVars(ctx *gin.Context) {
	account := handler.ensureAccount(ctx)
	if account == nil {
		return
	}
	appName := ctx.Param("appName")
	if appName == "" {
		appName = ctx.Query("appName")
	}
	data, err := handler.appRepo.GetEnvironmentVars(appName)
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: data})
}

func (handler *DeploymentHandler) GetApplicationLogsHandler(ctx *gin.Context) {
	account := handler.ensureAccount(ctx)
	if account == nil {
		return
	}
	appName := ctx.Param("appName")
	if appName == "" {
		BadRequestResponse(ctx, "application name is missing")
		return
	}
	app, err := handler.appRepo.GetApp(appName)
	if err != nil {
		BadRequestResponse(ctx, err.Error())
		return
	}
	if app.AccountId != account.ID {
		ForbiddenRequestResponse(ctx, "This app does not belong to you")
		return
	}
	data, err := handler.repo.GetApplicationLogs(appName)
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	type logResponse struct {
		Logs string `json:"logs"`
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: &logResponse{Logs: data}})
}

func (handler *DeploymentHandler) ensureAccount(ctx *gin.Context) *types.Account {
	token := ctx.GetHeader(tokenHeaderName)
	if token == "" {
		ForbiddenRequestResponse(ctx, "authenticate token is missing from request")
		return nil
	}
	account, err := handler.store.Get(token)
	if err != nil {
		ForbiddenRequestResponse(ctx, "authentication token not found or it has expired. Please re-authenticate")
		return nil
	}
	return account
}
