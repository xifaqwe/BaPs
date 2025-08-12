package pack

import (
	"github.com/gin-gonic/gin"
)

func ShopList(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"shop_list": []interface{}{},      
			"category_list": []interface{}{}, 
			"refresh_time": 0,                 
		},
	})
}

func RegisterShopRoutes(router *gin.Engine) {
	router.GET("/shop/list", ShopList)
	router.POST("/shop/list", ShopList)
}
