package parser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "post"
)

type AutoGenerated struct {
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Links struct {
				Previous interface{} `json:"previous"`
				Current  string      `json:"current"`
				Next     string      `json:"next"`
			} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID     int    `json:"id"`
		UserID int    `json:"user_id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	} `json:"data"`
}

func Parser(w http.ResponseWriter, r *http.Request) {
	//parse json file
	resp, err := http.Get("https://gorest.co.in/public/v1/posts")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//parse json to struct
	var autoGenerated AutoGenerated
	err = json.Unmarshal(body, &autoGenerated)
	if err != nil {
		log.Fatal(err)
	}
	//connect to postgresql
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created connection to database")

	//insert unmarshalled data into database
	for i := 0; i < len(autoGenerated.Data); i++ {
		_, err = db.Exec("INSERT INTO posts (id, user_id, title, body) VALUES ($1, $2, $3, $4)",
			autoGenerated.Data[i].ID, autoGenerated.Data[i].UserID, autoGenerated.Data[i].Title, autoGenerated.Data[i].Body)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Successfully inserted data into database")

}
