package routes

import (
	"swift-api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.GET("/swift-codes/:swift-code", controller.GetSwiftCode)
		v1.GET("/swift-codes/country/:countryISO2", controller.GetByCountry)
		v1.POST("/swift-codes", controller.AddSwiftCode)
		v1.DELETE("/swift-codes/:swift-code", controller.DeleteSwiftCode)
		v1.POST("/import-swift", controller.ImportSWIFTData)
	}
}
