package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tz/models"

	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "post"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	//connect to postgresql
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//check if the connection is established
	err = db.Ping()
	if err != nil {
		panic(err)

	}
	fmt.Println("Successfully created connection to database")
	return db

}

// GetPost will return a single post by its id
func GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the post id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getPost function with post id to retrieve a single post
	post, err := getPost(int64(id))

	if err != nil {
		log.Fatalf("Unable to get post. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(post)
}

func GetAllPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the posts in the db
	posts, err := getAllPosts()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the posts as response
	json.NewEncoder(w).Encode(posts)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty post of type models.User
	var post models.Data

	// decode the json request to post
	err = json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update user to update the post
	updatedRows := updatePost(int64(id), post)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deletePost, convert the int to int64
	deletedRows := deletePost(int64(id))

	// format the message string
	msg := fmt.Sprintf("Post updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func getPost(id int64) (models.Data, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a post of models.Data type
	var post models.Data

	// create the select sql query
	sqlStatement := `SELECT * FROM posts WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Body)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return post, nil
	case nil:
		return post, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty post on error
	return post, err
}

// get one post from the DB by its id
func getAllPosts() ([]models.Data, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var posts []models.Data

	// create the select sql query
	sqlStatement := "SELECT * FROM posts"

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		// create a post of models.Data type
		var post models.Data

		// unmarshal the row object to post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)

		}

		// append the post to the posts slice
		posts = append(posts, post)

	}
	// return empty post on error
	return posts, err
}

// update post in the DB
func updatePost(id int64, post models.Data) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE posts SET user_id=$2, title=$3, body=$4 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, post.UserID, post.Title, post.Body)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete post in the DB
func deletePost(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE  FROM posts WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
