package http

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestDeploymentHandler_CreateDockerDeployment(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	appRepo := mocks.NewMockAppsRepository(controller)
	session := mocks.NewMockStore(controller)
	deploymentRepo := mocks.NewMockDeploymentRepository(controller)

	testDockerUrl := "foo/bar:latest"
	appName := "testApp"
	testHeader := "testHeader"
	appRepo.EXPECT().GetApp(appName).Return(mockApp, nil)
	session.EXPECT().Get(testHeader).Return(mockAccount, nil)
	deploymentRepo.EXPECT().CreateDockerDeployment(mockApp, testDockerUrl).Return(&types.DeploymentResult{
		Address: "https://new.app",
		Version: "v1",
	}, nil)
	type payload struct{
		AppName   string `json:"app_name"`
		DockerUrl string `json:"docker_url"`
	}
	p := &payload{AppName: appName, DockerUrl: testDockerUrl}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/apps/deploy", buf)
	r.Header.Add(tokenHeaderName, testHeader)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	handler := NewDeploymentHandler(deploymentRepo, appRepo, session)
	handler.CreateDockerDeployment(ctx)

	assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("expected http response code %d;" +
		" got %d instead", http.StatusOK, w.Code))
}

func TestDeploymentHandler_CreateDockerDeployment_Access_Control(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	appRepo := mocks.NewMockAppsRepository(controller)
	session := mocks.NewMockStore(controller)
	deploymentRepo := mocks.NewMockDeploymentRepository(controller)

	testDockerUrl := "foo/bar:latest"
	appName := "testApp"
	testHeader := "testHeader"
	// use a different account,
	// simulate a situation whereby a different user try to deploy to an app
	// that does not belong to them
	mAccount := &types.Account{
		Model:          gorm.Model{ID: 2},
		Name:           "Foo Bar",
		Email:          "foo@bar.co",
		AccountToken:   "",
	}
	appRepo.EXPECT().GetApp(appName).Return(mockApp, nil)
	session.EXPECT().Get(testHeader).Return(mAccount, nil)
	type payload struct{
		AppName   string `json:"app_name"`
		DockerUrl string `json:"docker_url"`
	}
	p := &payload{AppName: appName, DockerUrl: testDockerUrl}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/apps/deploy", buf)
	r.Header.Add(tokenHeaderName, testHeader)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	handler := NewDeploymentHandler(deploymentRepo, appRepo, session)
	handler.CreateDockerDeployment(ctx)

	assert.Equal(t, http.StatusForbidden, w.Code, fmt.Sprintf("expected http response code %d;" +
		" got %d instead", http.StatusForbidden, w.Code))
}