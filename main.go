//Recipes API
//
//This is a sample recipes API. You can find out more about the API at
//
//Schemes: http
//Host: localhost:8080
//BasePath: /
//Version: 1.0.0
//Contact: Me
//
//Consumes:
//- application/json
//
//Produces:
//- application/json
//swagger:meta

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipes-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	fmt.Println(status)
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

//swagger:operation POST /recipes recipe newRecipe
//Create new recipe
//---
//produces:
//- application/json
//responses:
//'200':
//description: Successfull create recipe
//'404':
//description: Error parse object
// func NewRecipeHandler(c *gin.Context) {
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	recipe.ID = primitive.NewObjectID()
// 	recipe.PublishedAt = time.Now()
// 	_, err = collection.InsertOne(ctx, recipe)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, recipe)
// }

//swagger:operation GET /recipes recipes listRecipes
//Returns list of recipes
//---
//produces:
//- application/json
//responses:
//'200':
//description: Successful operation

// func ListRecipesHandler(c *gin.Context) {
// 	cur, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer cur.Close(ctx)
// 	recipes := make([]Recipe, 0)
// 	for cur.Next(ctx) {
// 		var recipe Recipe
// 		cur.Decode(&recipe)
// 		recipes = append(recipes, recipe)
// 	}
// 	c.JSON(http.StatusOK, recipes)
// }

//swagger:operation PUT /recipes/{id} recipes updateRecipe
//Update an existing recipe
//---
//parameters:
//- name: id
//in: path
//description: ID of the recipe
//required: true
//types: string
//produces:
//- application/json
//responses:
//'200':
//description: Successful operation
//'400':
//description: Invalid input
//'404':
//description: Invalid recipe ID
// func UpdateRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.D{{"$set", bson.D{
// 		{"name", recipe.Name},
// 		{"instructions", recipe.Instructions},
// 		{"ingredients", recipe.Ingredients},
// 		{"tags", recipe.Tags},
// 	}}})

// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
// }

//swagger:operation DELETE /recipes/{id} recipes deleteRecipe
//Delete an existing recipe
//---
//parameters:
//- name: id
//in: path
//description: ID of the recipe
//required: true
//types: string
//produces:
//- application/json
//responses:
//'200':
//description: Recipe has been deleted
//'404':
//description: Recipe not found

// func DeleteRecipeHandler(c *gin.Context) {
// 	// id := c.Param("id")
// 	index := -1
// 	// for i := 0; i < len(recipes); i++ {
// 	// 	if recipes[i].ID == id {
// 	// 		index = i
// 	// 	}
// 	// }

// 	if index == -1 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "Recipe not found",
// 		})
// 		return
// 	}

// 	recipes = append(recipes[:index], recipes[index+1:]...)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Recipe has been deleted",
// 	})
// }

//swagger:operation GET /recipes/search recipes searchRecipe
//Search Recipes exist in recipes
//---
//produces:
//- application/json
//responses:
//'200':
//description: list recipes

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)
// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}

// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}
// 	c.JSON(http.StatusOK, listOfRecipes)
// }

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.Run()
}
