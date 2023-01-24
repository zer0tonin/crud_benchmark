package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title     string
	Body      string
	CreatedAt string
}

func (p Post) StringID() string {
	return p.ID.Hex()
}

func List(w http.ResponseWriter, r *http.Request, coll *mongo.Collection) {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"createdat", -1}}).SetLimit(10)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var posts []Post
	if err := cursor.All(context.TODO(), &posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = listTemplate.Execute(
		w,
		map[string]interface{}{
			"posts": posts,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Get(w http.ResponseWriter, r *http.Request, coll *mongo.Collection) {
	id := path.Base(r.URL.Path)
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := coll.FindOne(context.TODO(), bson.M{"_id": bsonID})
	if result == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var post Post
	if err := result.Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = getTemplate.Execute(
		w,
		post,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PostPayload struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string
}

func Create(w http.ResponseWriter, r *http.Request, coll *mongo.Collection) {
	var payload *PostPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	post := PostPayload{
		Title:     payload.Title,
		Body:      payload.Body,
		CreatedAt: fmt.Sprint(now.Unix()),
	}
	coll.InsertOne(context.TODO(), post)

	// redirect to list
	http.Redirect(w, r, "/", http.StatusCreated)
}
