
    







## **Fully functional ECOMMERCE API USING GIN FRAMEWORK AND MONGODB** ##

![GOPHER BUY RELAX](https://github.com/5olitude/ecommerce/blob/main/rest.png){:height="36px" width="36px"}

-----**Initial  Release v2.30** ⚠️Not tested the efficiency

***project  structure***


* **Ecommerce** 📁
  - controllers 📁
    - controllers.go📝
  - database📁
    - database.go📝
  - middleware📁
    - middleware.go📝
  - models📁
    - models.go📝
  - routes📁
    - routes.go📝
  - tokens📁
    - tokengen.go📝
  - go.sum 📝
  - main.go📝


  ####   ``` Using Mongodb for small scale ecommerce industry is not a good idea instead use redis or mysql```

  
  ## API FUNCTIONALITY CURRENTLY ADDED?
     - Signup 🔒
     - Login  🔒
     - Product listing General View 👀
     - Adding the products to DB    
     - Sorting the products from DB using regex 👀
     - Adding  the products to cart 🛒
     - Removing the Product from cart🛒
     - Viewing the items in cart with total price🛒💰
     - Adding  the Home Address 🏠
     - Adding  the Work Address 🏢
     - Editing the Address ✂️
     - Deleting the Adress 🗑️
     - Checkout the Items from Cart
     - Buy Now products💰
     #### future implementations?

     - Pagination 1>>2>>3
     - Admin Part
     - etc***

    #### Packages Required?
    (github.com/dgrijalva/jwt-go) [moved to a new repo due to security implementations]

	(github.com/gin-gonic/gin v1.7.4)[Framework in golang -gin used for rest api]implementation )

	(github.com/go-playground/validator/v10 v10.9.0)[Validating the users struct such as email , phone etc]

	(go.mongodb.org/mongo-driver v1.7.2)[go driver for database connection with mongodb]

  (go.mongodb.org/mongo-driver/mongo)

  (go.mongodb.org/mongo-driver/mongo/options)

	(golang.org/x/crypto) [Encrypting the users pssword]

    go version 1.16

    ## Code At  glance in Database(database.go)

    first we need to define the database , make sure the mongodb installed in your system.
    The important thing we have to remember is we need to create  a database as well as two collections  , Two collections that is one for storing the user informations and the other for storing the product informations.We must have to configure the url and port where mongo is running

        //code example of url and port in databasetup.go
             mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")) --port and url
        //code example of defining collection name
            var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
    like to know more about configurations visit := 
    
    https://www.mongodb.com/blog/post/quick-start-golang-mongodb-starting-and-setup

    ## Code At glance in models.go
    This part defines us how should our database looks like, I think its not about programming skills but if you have creativity and basic syntax ideas and have some ability to defines struct from creativity 90% of work completed.Before Jumping into other codes we should have a rough idea about our plan so its better to lookup into models .

 - We have to define a product first , having a unique id , name and price

         type Product struct {
	     Product_ID   primitive.ObjectID `bson:"_id"`
	     Product_Name *string            `json:"product_name"`
	     Price        *uint64            `json:"price"`
	     Rating       *uint8             `json:"rating"`
	     Image        *string            `json:"image"`
       } 

- We have to define an slice of array products where a user  can store individual products

        type ProductUser struct {
	    Product_ID   primitive.ObjectID `bson:"_id"`
	    Product_Name *string            `json:"product_name" bson:"product_name"`
	    Price        int                `json:"price"  bson:"price"`
	    Rating       *uint              `json:"rating" bson:"rating"`
	    Image        *string            `json:"image"  bson:"image"`
      }
        
- The next struct we have to define the Address 

         type Address struct {
	    Address_id primitive.ObjectID `bson:"_id"`
	    House      *string            `json:"house_name" bson:"house_name"`
	    Street     *string            `json:"street_name" bson:"street_name"`
	    City       *string            `json:"city_name" bson:"city_name"`
	    Pincode    *string            `json:"pin_code" bson:"pin_code"`
        }
       
- If the user has ordered something the struct look like this  having an embedded struct inside a struct , here we define the ProductUser as a slice(A person can buy more than one product right?) and a payement struct to define Cash on delivery or digital payement

         type Order struct {
	        Order_ID       primitive.ObjectID `bson:"_id"`
	        Order_Cart     []ProductUser      `json:"order_list"  bson:"order_list"`
	        Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
        	  Price          int                `json:"total_price" bson:"total_price"`
         	  Discount       *int               `json:"discount"    bson:"discount"`
        	  Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
           }
         

The Payement struct is something look like this 

              type Payment struct {
	          Digital bool `json:"digital" bson:"digital"`
	          COD     bool `json:"cod"     bson:"cod"`
        }
         
***we can define those structs as the simple fields or blocks (genrally known as documents and subdocuments in mongodb)in the user databse or the user struct***

- Finally in the user struct we are going to embedd the simple structs .The new fields in struct are ID, Name ,Email, Token etc  

        type User struct {
	     ID              primitive.ObjectID `json:"_id" bson:"_id"`
	    First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	    Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30"`
    	Password        *string            `json:"password"   validate:"required,min=6"`
	    Email           *string            `json:"email"      validate:"email,required"`
	    Phone           *string            `json:"phone"      validate:"required"`
	    Token           *string            `json:"token"`
	    Refresh_Token   *string            `josn:"refresh_token"`
	    Created_At      time.Time          `json:"created_at"`
	    Updated_At      time.Time          `json:"updtaed_at"`
	    User_ID         string             `json:"user_id"`
	    UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	    Address_Details []Address          `json:"address" bson:"address"`
	    Order_Status    []Order            `json:"orders" bson:"orders"`
                }



## Code At  Glance in controllers.go

   This file mainly describes about the token authentication process . We have used the JWT authentication from dgrijalwa but now the repository has changed . I have used the same implemtaion for signup and login from 
    https://dev.to/joojodontoh/build-user-authentication-in-golang-with-jwt-and-mongodb-2igd , this blog is clear and precise about jwt auhentication rather than my explanation here.

  ***There is an important think we have to remember when defining  array struct  the mongodb converts the array to a nil in document field***

  So  to overcome this problem we make an empty array in signup function like this,whever a user calls the signup function it initialise the documents to empty array

         user.UserCart  =   make([]models.ProductUser, 0)
		 user.Address_Details = make([]models.Address, 0)
		 user.Order_Status = make([]models.Order, 0)


- **SIGNUP FUNCTION API CALL (POST REQUEST)**
      
  http://localhost:8000/users/signup

         {
          "first_name":"Joseph",
          "last_name":"Hermis",
          "email":"josephhermis@protonmail.com",
          "password":"unlucky",
          "phone":"+1558426655"
        }
        Response :"Successfully Signed Up!!"

  
- **LOGIN FUNCTION API CALL (POST REQUEST)**

    http://localhost:8000/users/login
              
              {
                "email":"josephhermis@protonmail.com",
                "password":"unlucky"
              }

    response will be like this 

             {
          "_id": "***********************",
          "first_name": "joseph",
          "last_name": "hermis",
          "password": "$2a$14$UIYjkTfnFnhg4qhIfhtYnuK9qsBQifPKgu/WPZAYBaaN17j0eTQZa",
          "email": "josephhermis@protonomail.com",
          "phone": "+1558921455",
          "token": "eyJc0Bwcm90b25vbWFpbC5jb20iLCJGaXJzdF9OYW1lIjoiam9zZXBoIiwiTGFzdF9OYW1lIjoiaGVybWlzIiwiVWlkIjoiNjE2MTRmNTM5ZjI5YmU5NDJiZDlkZjhlIiwiZXhwIjoxNjMzODUzNjUxfQ.NbcpVtPLJJqRF44OLwoanynoejsjdJb5_v2qB41SmB8",
          "Refresh_Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnLCJVaWQiOiIiLCJleHAiOjE2MzQzNzIwNTF9.ocpU8-0gCJsejmCeeEiL8DXhFcZsW7Z3OCN34HgIf2c",
          "created_at": "2021-10-09T08:14:11Z",
          "updtaed_at": "2021-10-09T08:14:11Z",
          "user_id": "61614f539f29be942bd9df8e",
          "usercart": [],
          "address": [],
          "orders": []
            } 
          Login Function call create an outlayer for our collection
- **Admin add Product Function  POST REQUEST**
   
   **note this function is not seperated from normal user fixed soon for admin**
   
   http://localhost:8000/admin/addproduct

        {
        "product_name":"laptop",
        "price":300,
        "rating":10,
        "image":"1.jpg"
      }

       Response : "Successfully added our Product Admin!!"

- **View all the Products in db GET REQUEST**
    
    pagination added soon in next release

    http://localhost:8000/users/productview

              Response 
              [
        [
        {
            "Product_ID": "6153ff8edef2c3c0a02ae39a",
            "product_name": "notepad",
            "price": 50,
            "rating": 10,
            "image": "penc.jpg"
        },
        {
            "Product_ID": "616152679f29be942bd9df8f",
            "product_name": "laptop",
            "price": 300,
            "rating": 10,
            "image": "1.jpg"
        },
        {
            "Product_ID": "616152ee9f29be942bd9df90",
            "product_name": "top",
            "price": 300,
            "rating": 10,
            "image": "1.jpg"
        },
        {
            "Product_ID": "616152fa9f29be942bd9df91",
            "product_name": "table",
            "price": 300,
            "rating": 10,
            "image": "1.jpg"
        },
        {
            "Product_ID": "616153039f29be942bd9df92",
            "product_name": "apple",
            "price": 300,
            "rating": 10,
            "image": "1.jpg"
        }
      ]


-  **Search Product by regex function (GET REQUEST)**

defines the word search sorting 
     http://localhost:8000/users/search?name=le

         response
         [
            {
                "Product_ID": "616152fa9f29be942bd9df91",
                "product_name": "table",
                "price": 300,
                "rating": 10,
                "image": "1.jpg"
            },
            {
                "Product_ID": "616153039f29be942bd9df92",
                "product_name": "apple",
                "price": 300,
                "rating": 10,
                "image": "1.jpg"
            }
        ]
      
The corresponding Query to mongodb is **ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})**


- **Adding the Products to the Cart (GET REQUEST)**

    http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx

    Corresponding mongodb  query 

          filter := bson.D{primitive.E{Key: "_id", Value: id}}
		  update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
		  _, err = UserCollection.UpdateOne(ctx, filter, update)

- **Removing Item From the Cart (GET REQUEST)**

    http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx

    Corresponding mongodb  query

           filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
	  	 update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": removed_id}}}
	       _, err = UserCollection.UpdateMany(ctx, filter, update)

-  **Listing the item in the users cart (GET REQUEST) and total price**

    http://localhost:8000/listcart?id=xxxxxxuser_idxxxxxxxxxx

      
      Corresponding Mongodb Query (WE are using the aggrgate operation to find sum)

        filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

-  **Addding the Address (POST REQUEST)**

     http://localhost:8000/addadress?id=user_id*************

     The Address array is limited to two values home and work address more than two address is not acceptable

        {
          "house_name":"jupyterlab",
          "street_name":"notebook",
          "city_name":"mars",
          "pin_code":"685607"
        }

-  **Editing the Home Address(PUT REQUEST)**

     http://localhost:8000/edithomeaddress?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx

     
-  **Editing the Work Address(PUT REQUEST)**
  
     http://localhost:8000/editworkaddress?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx


-  **Delete Addresses(GET REQUEST)**

      http://localhost:8000/deleteaddresses?id=xxxxxxxxxuser_idxxxxxxxxxxxxx

      delete both addresses

-  **Cart Checkout Function and placing the order(GET REQUEST)**
 
     After placing the order the items have to be deleted from cart functonality added

     http://localhost:8000?id=xxuser_idxxx

-  **Instantly Buying the  Products(GET REQUEST)**
      
      http://localhost:8000?pid=xxproduct_idxxx&id=xxxxuser_idxxxx


##   Code At Glance in main.go

All the routes defined here requires the api authentication key 


