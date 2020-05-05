package server

import (
	"github.com/adigunhammedolalekan/namegenerator"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v6"
	"github.com/saas/hostgolang/pkg/config"
	"github.com/saas/hostgolang/pkg/http"
	proxy "github.com/saas/hostgolang/pkg/proxyclient"
	"github.com/saas/hostgolang/pkg/repository"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/session"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	addr   string
	router *gin.Engine
}

func NewServer(addr string) (*Server, error) {
	if err := godotenv.Load("creds.env"); err != nil {
		return nil, err
	}
	db, err := services.CreateDatabaseConnection(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}
	cfg := &config.Config{
		ProxyServerAddress: os.Getenv("PROXY_SERVER_URL"),
		ProxySecret:        os.Getenv("PROXY_SECRET"),
		ServerUrl: os.Getenv("SERVER_URL"),
		Registry: config.RegistryConfig{
			Url:      "registry.hostgolang.com",
			Username: "MASTER",
			Password: "manman",
		},
	}

	proxyClient, err := proxy.NewProxyClient(cfg)
	if err != nil {
		return nil, err
	}

	k8sConfig, err := createK8sClient()
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	dy, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	k8sService := services.NewK8sService(k8sClient, dy, cfg, k8sConfig)
	dockerService, err := createDockerService(cfg)
	if err != nil {
		return nil, err
	}
	minioClient, err := minio.New(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		false)
	if err != nil {
		return nil, err
	}
	storageClient, err := services.NewMinioStorageClient(minioClient)
	if err != nil {
		return nil, err
	}
	resourseK8sClient := services.NewResourcesService(k8sClient, k8sService)
	sessionStore := session.NewRedisSessionStore(redisClient)
	accountRepo := repository.NewAccountRepository(db, sessionStore)
	appRepo := repository.NewAppsRepository(db, namegenerator.NewNameGenerator(time.Now().UnixNano()), k8sService)
	deploymentRepo := repository.NewDeploymentRepository(db, dockerService, k8sService,
		proxyClient, appRepo, storageClient, cfg)
	resourceRepo := repository.NewResourcesDeploymentRepository(db, appRepo, accountRepo, resourseK8sClient, dockerService)

	apiHandler := http.NewApiHandler(appRepo, accountRepo, sessionStore)
	deploymentHandler := http.NewDeploymentHandler(deploymentRepo, appRepo, sessionStore)
	resourcesDeploymentHandler := http.NewResourcesDeploymentHandler(resourceRepo, appRepo, sessionStore)
	router := gin.Default()
	apiRouter := router.Group("/api")
	apiRouter.POST("/account/new", apiHandler.CreateAccountHandler)
	apiRouter.POST("/account/authenticate", apiHandler.AuthenticateAccountHandler)
	apiRouter.POST("/me/apps", apiHandler.CreateAppHandler)
	apiRouter.GET("/me/apps", apiHandler.GetAccountApps)
	apiRouter.POST("/apps/docker/deploy", deploymentHandler.CreateDockerDeployment)
	apiRouter.GET("/apps/configs/:appName", deploymentHandler.GetEnvironmentVars)
	apiRouter.GET("/apps/logs/:appName", deploymentHandler.GetApplicationLogsHandler)
	apiRouter.POST("/apps/configs/:appName", deploymentHandler.UpdateEnvironmentVars)
	apiRouter.DELETE("/apps/configs/unset/:appName", deploymentHandler.DeleteEnvironmentVars)
	apiRouter.GET("/apps/scale/:appName", deploymentHandler.ScaleAppHandler)
	apiRouter.GET("/apps/ps/:appName", deploymentHandler.ListRunningInstances)
	apiRouter.PUT("/apps/rollback/:appName", deploymentHandler.RollbackDeploymentHandler)
	apiRouter.POST("/apps/resource/new/:appName", resourcesDeploymentHandler.CreateResourceHandler)
	apiRouter.DELETE("/apps/resource/remove/:appName", resourcesDeploymentHandler.DeleteResourceHandler)
	apiRouter.GET("/apps/releases/:appName", deploymentHandler.GetReleasesHandler)
	apiRouter.POST("/apps/domain/new", deploymentHandler.AddDomainHandler)
	apiRouter.GET("/apps/resource/dump/:appName", resourcesDeploymentHandler.DumpDatabaseHandler)
	apiRouter.GET("/apps/resource/restore/:appName", resourcesDeploymentHandler.RestoreDatabaseHandler)
	apiRouter.DELETE("/apps/domain/remove", deploymentHandler.RemoveDomainHandler)
	apiRouter.POST("/apps/use", apiHandler.UseAppHandler)
	apiRouter.GET("/status", apiHandler.StatusHandler)
	return &Server{
		addr:   addr,
		router: router,
	}, nil
}

func createK8sClient() (*rest.Config, error) {
	k8sConfigPath := ""
	if k8sConfigPath = os.Getenv("K8S_CONFIG_DIR"); k8sConfigPath == "" {
		k8sConfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	c, err := clientcmd.BuildConfigFromFlags("", k8sConfigPath)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func createDockerService(cfg *config.Config) (services.DockerService, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.37"), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return services.NewDockerService(docker, cfg), nil
}

func (s *Server) Run() error {
	if err := s.router.Run(s.addr); err != nil {
		return err
	}
	return nil
}
