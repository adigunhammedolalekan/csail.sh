package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/mocks"
	"github.com/saas/hostgolang/pkg/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockAccount = &types.Account{
	Model: gorm.Model{
		ID: 1,
	},
	Name:           "test",
	Email:          "test@gmail.com",
	Password:       "password",
	GithubId:       "",
	CompanyName:    "company",
	CompanyWebsite: "company.com",
	AccountToken:   "token",
}
var mockApp = &types.App{
	Model: gorm.Model{
		ID: 1,
	},
	AccountId: 1,
	AppName:   "testApp",
	AccessUrl: "",
}

func TestApiHandler_CreateAccountHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockOpt := &types.NewAccountOpts{
		Name:           "test",
		Email:          "test@gmail.com",
		Password:       "password",
		CompanyName:    "company",
		CompanyWebsite: "company.com",
	}

	accountRepo := mocks.NewMockAccountRepository(controller)
	sessionStore := mocks.NewMockStore(controller)
	handler := NewApiHandler(nil, accountRepo, sessionStore)
	accountRepo.EXPECT().CreateAccount(mockOpt).Return(mockAccount, nil)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(mockOpt); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/account/new", buf)
	handler.CreateAccountHandler(ctx)

	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("expected http code %d; got %d", want, got)
	}
	assert.Equal(t, w.Code, http.StatusOK, "expected code does not match")
}

func TestCreateAccountHandlerBadInput(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountRepo := mocks.NewMockAccountRepository(controller)
	handler := NewApiHandler(nil, accountRepo, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/account/new", nil)
	handler.CreateAccountHandler(ctx)

	if got, want := w.Code, http.StatusBadRequest; got != want {
		t.Fatalf("expected http code %d; got %d", want, got)
	}
}

func TestApiHandler_AuthenticateAccountHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	opt := &types.AuthenticateAccountOpts{
		Email:    "test@gmail.com",
		Password: "password",
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(opt); err != nil {
		t.Fatal(err)
	}
	accountRepo := mocks.NewMockAccountRepository(controller)
	accountRepo.EXPECT().AuthenticateAccount(opt).Return(mockAccount, nil)

	handler := NewApiHandler(nil, accountRepo, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/account/authenticate", buf)

	handler.AuthenticateAccountHandler(ctx)

	if want, got := http.StatusOK, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}

func TestApiHandler_AuthenticateAccountHandler2BadInput(t *testing.T) {
	handler := NewApiHandler(nil, nil, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/account/authenticate", nil)

	handler.AuthenticateAccountHandler(ctx)

	if want, got := http.StatusBadRequest, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}

func TestApiHandler_CreateAppHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "testApp"
	appRepo := mocks.NewMockAppsRepository(controller)
	sessionStore := mocks.NewMockStore(controller)
	sessionStore.EXPECT().Get("testUser").Return(mockAccount, nil)
	appRepo.EXPECT().CreateApp(name, "SR", mockAccount.ID).Return(mockApp, nil)

	handler := NewApiHandler(appRepo, nil, sessionStore)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/me/apps", nil)
	req.Header.Set(tokenHeaderName, "testUser")
	req.URL.RawQuery += "&name=" + name + "&plan=" + "SR"
	ctx.Request = req

	handler.CreateAppHandler(ctx)
	if want, got := http.StatusOK, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}

func TestApiHandler_CreateAppHandlerAuthError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	name := "testApp"
	appRepo := mocks.NewMockAppsRepository(controller)
	sessionStore := mocks.NewMockStore(controller)
	sessionStore.EXPECT().Get("testUser").Return(nil, errors.New("account data not found"))

	handler := NewApiHandler(appRepo, nil, sessionStore)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/me/apps", nil)
	req.Header.Set(tokenHeaderName, "testUser")
	req.URL.RawQuery += "&name=" + name
	ctx.Request = req

	handler.CreateAppHandler(ctx)
	if want, got := http.StatusForbidden, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}

func TestApiHandler_GetAccountApps(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	testUser := "testUser"
	appRepo := mocks.NewMockAppsRepository(controller)
	sessionStore := mocks.NewMockStore(controller)

	sessionStore.EXPECT().Get(testUser).Return(mockAccount, nil)
	appRepo.EXPECT().GetAccountApps(mockAccount.ID).Return([]types.App{*mockApp}, nil)

	handler := NewApiHandler(appRepo, nil, sessionStore)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/me/apps", nil)
	req.Header.Set(tokenHeaderName, testUser)
	ctx.Request = req

	handler.GetAccountApps(ctx)
	if want, got := http.StatusOK, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}

func TestApiHandler_GetAccountApps_Forbidden(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	testUser := "testUser"
	appRepo := mocks.NewMockAppsRepository(controller)
	sessionStore := mocks.NewMockStore(controller)

	sessionStore.EXPECT().Get(testUser).Return(nil, errors.New("account not found"))
	handler := NewApiHandler(appRepo, nil, sessionStore)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/me/apps", nil)
	req.Header.Set(tokenHeaderName, testUser)
	ctx.Request = req

	handler.GetAccountApps(ctx)
	if want, got := http.StatusForbidden, w.Code; got != want {
		t.Fatalf("expected httpCode %d; got %d instead", want, got)
	}
}