package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	conn *sql.DB
)

type User struct {
	Name string
	Age  int
}

func main() {
	// Crie um novo servidor Gin
	server := gin.New()

	// Configure o roteamento para servir o arquivo "index.html" como página principal
	server.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	// Rota para processar o envio do formulário
	server.POST("/processar", func(c *gin.Context) {
		// Recupere os dados do formulário
		nome := c.PostForm("nome")
		idade := c.PostForm("idade")

		idadeInt, err := strconv.Atoi(idade)
		if err != nil {
			// Lida com erros de conversão, se necessário
			c.String(http.StatusBadRequest, "Erro: O valor não é um número inteiro válido.")
			return
		}

		user := User{
			Name: nome,
			Age:  idadeInt,
		}

		DatabaseConnection()

		if insertError := DatabaseInsertion(user); insertError != nil {
			c.String(http.StatusInternalServerError, "Erro ao inserir no banco de dados: "+err.Error())
			return
		}

		c.String(http.StatusOK, "Dados do formulário recebidos e inseridos no banco de dados com sucesso!")

	})

	// Inicia o servidor
	server.Run()
}

func DatabaseConnection() {

	var err error

	// Acesse as variáveis de ambiente para obter informações sensíveis
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	stringConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conecte-se ao banco de dados
	conn, err = sql.Open("mysql", stringConnection)
	if err != nil {
		panic(err)
	}

}

func DatabaseInsertion(user User) error {
	// Use um comando preparado com placeholders
	_, err := conn.Exec("INSERT INTO user_data (nome, idade) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		fmt.Println("Erro ao inserir no banco de dados:", err.Error())
		return err
	}

	return nil
}
