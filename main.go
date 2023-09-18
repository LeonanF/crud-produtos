package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	// Carregue os modelos HTML
	server.LoadHTMLGlob("static/*")

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

	// Configure o roteamento para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		// Crie uma consulta para recuperar todos os usuários.
		filter := bson.D{}

		// Execute a consulta.
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		var users []User

		// Itere pelos documentos no cursor.
		for cursor.Next(context.Background()) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				c.AbortWithStatus(500)
				return
			}
			users = append(users, user)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Users": users,
		})
	})

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
