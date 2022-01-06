package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	httpServerAddress string
}

// Configure Swagger Services
func (n *NSXAPI) Configure(ginEngine *gin.Engine) *gin.Engine {
	logrus.Info("Configuring services")

	n.ConfigureSwaggo(ginEngine, n.httpServerAddress)
	n.configureMVC()

	return ginEngine
}

func (n *NSXAPI) configureMVC() {
}

// ConfigureSwaggo All configurations for swagger
func (n *NSXAPI) ConfigureSwaggo(ginEngine *gin.Engine, httpServerAddress string) {
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
