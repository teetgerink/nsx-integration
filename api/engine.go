package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	//DatabaseOpenFunction give the default to postgresOpen
	DatabaseOpenFunction = postgres.Open
)

const (
	//SWAGGERPATH used for swaggo
	SWAGGERPATH = "/swagger/*any"
	//SWAGGERHTMLPATH path for html file
	SWAGGERHTMLPATH = "/swagger/index.html"
	//SWAGGERJSONPATH json path
	SWAGGERJSONPATH = "/swagger/doc.json"
)

//NSXAPI struct contains httpServerAddress
type NSXAPI struct {
	httpServerAddress    string
	ConnectionString     string
	DatabaseOpenFunction func(dsn string) gorm.Dialector
	Database             *gorm.DB
}

// Configure Swagger Services
func (n *NSXAPI) Configure(ginEngine *gin.Engine) *gin.Engine {
	logrus.Info("Configuring services")
	n.configureDatabase(&gorm.Config{})
	n.configureSwaggo(ginEngine, n.httpServerAddress)
	n.configureMVC()

	return ginEngine
}

// configureMVC running the services controllers and listeners
func (n *NSXAPI) configureMVC() {
}

func (n *NSXAPI) configureDatabase(config *gorm.Config, configureModels ...interface{}) {
	logrus.Info("Connecting and automatically create database models")

	// By default the ConnectionString is loaded via an environment variable but this allows us to override it
	// in the scenario of running parallel tests
	if n.ConnectionString == "" {
		n.ConnectionString = os.Getenv("DATABASE_CONNECTION_URL")
	}

	// By default the DatabaseOpenFunction is implemented with gorm's postgres.Open but this allows us to override it
	if n.DatabaseOpenFunction == nil {
		logrus.Info("Setting Database Open Function with default postgres.Open")
		n.DatabaseOpenFunction = DatabaseOpenFunction
	}

	db, err := gorm.Open(n.DatabaseOpenFunction(n.ConnectionString), config)
	n.Database = db

	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(configureModels...); err != nil {
		panic(err)
	}
}

// configureSwaggo All configurations for swagger
func (n *NSXAPI) configureSwaggo(ginEngine *gin.Engine, httpServerAddress string) {
	swaggerURL := fmt.Sprintf("%v%v", httpServerAddress, SWAGGERJSONPATH)

	logrus.Infof("Configuring swagger on url %v and redirect on /", swaggerURL)
	url := ginSwagger.URL(swaggerURL)

	// Main route redirects to swagger
	ginEngine.GET("", func(c *gin.Context) {
		c.Redirect(301, SWAGGERHTMLPATH)
	})

	// Serve swagger files
	ginEngine.GET(SWAGGERPATH, ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
