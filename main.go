package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Produtos struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	// Crie um novo servidor Gin
	server := gin.New()
	// Carregue os modelos HTML
	server.LoadHTMLGlob("static/*")

	MONGODB_URI := os.Getenv("MONGODB_URI")
	DB_NAME := os.Getenv("DB_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")

	// Conecte-se ao MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// Coleção de produtos no MongoDB
	collection := client.Database(DB_NAME).Collection(COLLECTION_NAME)

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

		var produtos []Produtos

		// Itere pelos documentos no cursor.
		for cursor.Next(context.Background()) {
			var produto Produtos
			if err := cursor.Decode(&produto); err != nil {
				c.AbortWithStatus(500)
				return
			}
			produtos = append(produtos, produto)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Produtos": produtos,
		})
	})

	server.POST("/processar", func(c *gin.Context) {
		// Obtenha o valor do campo "nome" do formulário
		nome := c.PostForm("nome")

		// Obtenha o valor do campo "idade" do formulário
		valor := c.PostForm("valor")

		// Crie um novo documento BSON com base nos valores do formulário
		produto := Produtos{
			Name:  nome,
			Value: valor,
		}

		// Insira o documento no MongoDB
		insertResult, err := collection.InsertOne(ctx, produto)
		if err != nil {
			panic(err)
		}

		// Imprima o ID do documento inserido
		fmt.Println("ID do documento inserido:", insertResult.InsertedID)
	})

	server.Run(":8080") // Altere a porta, se necessário
}
