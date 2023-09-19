package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Produtos struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type Config struct {
	MongoDBURI     string `json:"MONGODB_URI"`
	DatabaseName   string `json:"DATABASE_NAME"`
	CollectionName string `json:"COLLECTION_NAME"`
}

func main() {
	// Crie um novo servidor Gin
	server := gin.New()
	// Carregue os modelos HTML
	server.LoadHTMLGlob("static/*")

	if os.Getenv("MONGODB_URI") == "" {
		// Abra o arquivo config.json e em caso de erro, imprime o erro e encerra a função
		configFile, err := os.Open("config.json")
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo de configuração:", err)
			return
		}

		//Garante que o arquivo seja fechado no fim da execução
		defer configFile.Close()

		//Criada váriavel do tipo Config (struct no início) para armazenar os dados do arquivo .json
		var config Config

		// É criado um novo decodificador para ler os dados do arquivo configFile
		decoder := json.NewDecoder(configFile)

		//O código a seguir é que vai, efetivamente, decodificar e converter o JSON em uma estrutura de dados do tipo Config
		//Para entender a criação de err como váriavel de bloco, verificar documentação da função createServer()
		//Em caso de erro, a função é imediatamente encerrada e o erro é impresso
		if err := decoder.Decode(&config); err != nil {
			fmt.Println("Erro ao decodificar o arquivo de configuração:", err)
			return
		}

		//Seta as variáveis de ambiente com os valores do arquivo de configuração (propriedades do struct Config)
		os.Setenv("MONGODB_URI", config.MongoDBURI)
		os.Setenv("DATABASE_NAME", config.DatabaseName)
		os.Setenv("COLLECTION_NAME", config.CollectionName)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	DB_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")

	// Conecte-se ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(context.Background())

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
		insertResult, err := collection.InsertOne(context.Background(), produto)
		if err != nil {
			panic(err)
		}

		// Imprima o ID do documento inserido
		fmt.Println("ID do documento inserido:", insertResult.InsertedID)
	})

	server.Run(":8080") // Altere a porta, se necessário
}
