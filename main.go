package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func main() {
	// Crie um novo servidor Gin
	server := gin.New()

	// Configure o roteamento para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	// Conecte-se ao MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://leonanfreitas:qtHEEIuQ84ePeavH@cluster0.2inszqq.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Coleção de usuários no MongoDB
	collection := client.Database("users_db").Collection("users")

	server.POST("/processar", func(c *gin.Context) {
		// Obtenha o valor do campo "nome" do formulário
		nome := c.PostForm("nome")

		// Obtenha o valor do campo "idade" do formulário
		idade := c.PostForm("idade")

		// Crie um novo documento BSON com base nos valores do formulário
		user := User{
			Name: nome,
			Age:  idade,
		}

		// Insira o documento no MongoDB
		insertResult, err := collection.InsertOne(ctx, user)
		if err != nil {
			panic(err)
		}

		// Imprima o ID do documento inserido
		fmt.Println("ID do documento inserido:", insertResult.InsertedID)
	})

	server.Run(":8080") // Altere a porta, se necessário
}
