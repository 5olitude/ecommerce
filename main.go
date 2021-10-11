package main

import (
	"ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// This is good and follows the advice of: https://12factor.net/config.
	// But you should do this for all config: mongodb (credentials, database, collections), SECRET_KEY in tokengen.go.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	// The authentication middleware is applied to all routes, including the /users/signup route. So nobody can actually use the application.
	router.Use(middleware.Authentication())
	// Your routes are inconsistent starting with and without '/'.
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", controllers.RemoveItem())
	router.GET("listcart", controllers.GetItemFromCart())
	router.POST("addaddress", controllers.AddAddress())
	router.PUT("edithomeaddress", controllers.EditHomeAddress())
	router.PUT("editworkaddress", controllers.EditWorkAddress())
	router.GET("deleteaddresses", controllers.DeleteAddress())
	router.GET("cartcheckout", controllers.BuyFromCart())
	router.GET("instantbuy", controllers.InstantBuy())
	//router.GET("logout", controllers.Logout())
	//break :)

	// Log the error that the router can possibly return.
	log.Fatal(router.Run(":" + port))
}
