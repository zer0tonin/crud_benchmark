package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var listTemplate *template.Template
var getTemplate *template.Template

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to load config")
	}

	listTemplate = template.Must(
		template.ParseFiles("./templates/list.html"),
	)
	getTemplate = template.Must(
		template.ParseFiles("./templates/get.html"),
	)
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(viper.GetString("mongo.uri")))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged mongo.")

	coll := client.Database("myDb").Collection("posts")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/" {
				List(w, r, coll)
			} else {
				Get(w, r, coll)
			}
		case http.MethodPost:
			Create(w, r, coll)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Listening on " + viper.GetString("host"))
	http.ListenAndServe(viper.GetString("host"), nil)
}
