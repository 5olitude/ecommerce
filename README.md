## **Fully functional ECOMMERCE API USING GIN FRAMEWORK AND MONGODB** ##

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



