package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/repository"
	"github.com/saas/hostgolang/pkg/session"
	"github.com/saas/hostgolang/pkg/types"
	"net/http"
)

const tokenHeaderName = "X-Auth-Token"

type ApiHandler struct {
	appRepo     repository.AppsRepository
	accountRepo repository.AccountRepository
	store       session.Store
}

func NewApiHandler(app repository.AppsRepository, account repository.AccountRepository, store session.Store) *ApiHandler {
	return &ApiHandler{
		appRepo:     app,
		accountRepo: account,
		store:       store,
	}
}

func (handler *ApiHandler) CreateAccountHandler(ctx *gin.Context) {
	opt := &types.NewAccountOpts{}
	if err := ctx.BindJSON(opt); err != nil {
		BadRequestResponse(ctx, "invalid or malformed request body")
		return
	}
	account, err := handler.accountRepo.CreateAccount(opt)
	if err != nil {
		BadRequestResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "account created", Data: account})
}

func (handler *ApiHandler) AuthenticateAccountHandler(ctx *gin.Context) {
	opt := &types.AuthenticateAccountOpts{}
	if err := ctx.BindJSON(opt); err != nil {
		BadRequestResponse(ctx, "invalid or malformed request body")
		return
	}
	account, err := handler.accountRepo.AuthenticateAccount(opt)
	if err != nil {
		ForbiddenRequestResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: account})
}

func (handler *ApiHandler) CreateAppHandler(ctx *gin.Context) {
	token := ctx.GetHeader(tokenHeaderName)
	if token == "" {
		ForbiddenRequestResponse(ctx, "authentication token is missing from request")
		return
	}
	account, err := handler.store.Get(token)
	if err != nil {
		ForbiddenRequestResponse(ctx, "failed to retrieve authentication data. Please reauthenticate")
		return
	}
	var request struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		request.Name = ctx.Query("name")
	}
	app, err := handler.appRepo.CreateApp(request.Name, account.ID)
	if err != nil {
		InternalServerErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: app})
}

func (handler *ApiHandler) GetAccountApps(ctx *gin.Context) {
	token := ctx.GetHeader(tokenHeaderName)
	if token == "" {
		ForbiddenRequestResponse(ctx, "authentication token is missing from request")
		return
	}
	account, err := handler.store.Get(token)
	if err != nil {
		ForbiddenRequestResponse(ctx, "failed to retrieve authentication data. Please re-authenticate")
		return
	}

	data, err := handler.appRepo.GetAccountApps(account.ID)
	switch err {
	case gorm.ErrRecordNotFound:
		ctx.JSON(http.StatusNotFound, &SuccessResponse{Error: true, Message: "no app has been created for this account", Data: nil})
	case nil:
		ctx.JSON(http.StatusOK, &SuccessResponse{Error: false, Message: "success", Data: data})
	default:
		InternalServerErrorResponse(ctx, "failed to retrieve app list at this time. Please retry later")
	}
}
