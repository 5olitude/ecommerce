package controllers

import (
	"context"
	"ecommerce/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddToCart adds products to the cart of the user.
// GET request
// http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx
func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productcart []models.ProductUser
		productqueryid := c.Query("id")
		userid := c.Query("normal")
		productid, _ := primitive.ObjectIDFromHex(productqueryid)
		if productqueryid == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Product Query not matched"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// You have sometimes multiple defer cancel() in your functions.
		defer cancel()
		searchfromdb, err := ProductCollection.Find(ctx, bson.M{"_id": productid})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Invalid ID refer")
			return
		}
		err = searchfromdb.All(ctx, &productcart)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		id, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			// Again giving the user the error that the database returns
			// is not good idea. Gives normal user useless information
			// and bad users too much information.
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Added to the cart")
	}
}

//function to remove item from cart
//GET Request
//http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx
func RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		remove_id := c.Query("id")
		user_id := c.Query("normal")
		if remove_id == "" || user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Query"})
			c.Abort()
			return
		}
		removed_id, _ := primitive.ObjectIDFromHex(remove_id)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			// Use log instead of fmt for logging
			log.Println(err)
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Again two cancels in a single function.
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": removed_id}}}
		_, err = UserCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Server Error")
			return
		}
		c.IndentedJSON(200, "Successfully removed from cart")
	}
}

//function to get all items in the cart and total price
//GET request
//http://localhost:8000/listcart?id=xxxxxxuser_idxxxxxxxxxx
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			// Why would you panic here?
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var getcartitems models.User
		var ordercart models.Order
		ordercart.Order_ID = primitive.NewObjectID()
		ordercart.Orderered_At = time.Now()
		ordercart.Order_Cart = make([]models.ProductUser, 0)
		ordercart.Payment_Method.COD = true
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		currentresults, err := UserCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

		ctx.Done()
		if err != nil {
			panic(err)
		}
		var getusercart []bson.M
		if err = currentresults.All(ctx, &getusercart); err != nil {
			panic(err)
		}
		var total_price int32
		for _, user_item := range getusercart {
			price := user_item["total"]
			total_price = price.(int32)
		}
		ordercart.Price = int(total_price)
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
		_, err = UserCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "something went wrong")
		}
		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&getcartitems)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "something went wrong")
		}
		var ktx, kancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer kancel()
		filter2 := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
		_, err = UserCollection.UpdateOne(ctx, filter2, update2)
		if err != nil {
			c.IndentedJSON(500, "something went wrong")
		}
		usercart_empty := make([]models.ProductUser, 0)
		filtered := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
		_, err = UserCollection.UpdateOne(ctx, filtered, updated)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Errror")
		}

		ktx.Done()
		c.IndentedJSON(200, "Successfully Placed the order")

	}
}

func InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		item_id := c.Query("pid")
		user_id := c.Query("id")
		if item_id == "" || user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Code"})
			c.Abort()
			return
		}
		itemt_id, err := primitive.ObjectIDFromHex(item_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal Server Erroe")
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal Server Error")
		}
		var product_details models.ProductUser
		var orders_detail models.Order
		orders_detail.Order_ID = primitive.NewObjectID()
		orders_detail.Orderered_At = time.Now()
		orders_detail.Order_Cart = make([]models.ProductUser, 0)
		orders_detail.Payment_Method.COD = true
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Always use the defer cancel() as close to the creation point as possible.
		defer cancel()
		err = ProductCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: itemt_id}}).Decode(&product_details)
		if err != nil {
			log.Println(err)
			// This kind of error is already a lot better for the user.
			c.IndentedJSON(400, "Something Wrong happened")
		}
		orders_detail.Price = product_details.Price
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		filter2 := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
		_, err = UserCollection.UpdateOne(ctx, filter2, update2)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "something wrong happened")
		}
		c.IndentedJSON(200, "Successully placed the order ")
	}
}
