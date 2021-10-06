package main

import (
	"ecommerce/controllers"
	"ecommerce/middleware"
	"ecommerce/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", controllers.AddToCart())
	router.GET("/removeitem", controllers.RemoveItem())
	router.GET("listcart", controllers.GetItemFromCart())
	router.Run(":" + port)
}
