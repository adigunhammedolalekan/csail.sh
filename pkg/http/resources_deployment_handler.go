package http

import (
	"github.com/gin-gonic/gin"
	"github.com/saas/hostgolang/pkg/repository"
	"github.com/saas/hostgolang/pkg/session"
	"github.com/saas/hostgolang/pkg/types"
	"net/http"
)

type ResourcesDeploymentHandler struct {
	repo repository.ResourcesDeployment
	appRepo repository.AppsRepository
	store session.Store
}

func NewResourcesDeploymentHandler(repo repository.ResourcesDeployment, appRepo repository.AppsRepository, store session.Store) *ResourcesDeploymentHandler {
	return &ResourcesDeploymentHandler{repo: repo, store: store, appRepo: appRepo}
}

func (handler *ResourcesDeploymentHandler) CreateResourceHandler(ctx *gin.Context) {
	account := handler.ensureAccount(ctx)
	if account == nil {
		return
	}
	appName := ctx.Param("appName")
	app, err := handler.appRepo.GetApp(appName)
	if err != nil {
		BadRequestResponse(ctx, err.Error())
		return
	}
	if app.AccountId != account.ID {
		ForbiddenRequestResponse(ctx, "forbidden")
		return
	}
	resourceName := ctx.Query("name")
	if resourceName == "" {
		BadRequestResponse(ctx, "resource name is missing")
		return
	}
	opt := &types.DeployResourcesOpt{
		AppName:     appName,
		Name:        resourceName,
		Memory:      0.2,
		Cpu:         0.1,
		StorageSize: 2, // in GB
	}
	result, err := handler.repo.DeployResource(opt)
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: result})
}


func (handler *ResourcesDeploymentHandler) ensureAccount(ctx *gin.Context) *types.Account {
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
