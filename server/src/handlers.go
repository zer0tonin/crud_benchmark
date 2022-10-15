package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

type Post struct {
	ID int
	Title string
	Body string
}

var posts = []Post{
	{
		ID: 1,
		Title: "first post",
		Body: "hello world",
	},
	{
		ID: 2,
		Title: "second post",
		Body: "hello again",
	},
}

func List(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFiles("./templates/list.html"),
	)
	err := t.Execute(
		w,
		map[string]interface{}{
			"posts": posts,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var post *Post
	for _, p := range posts {
		if p.ID == id {
			post = &p
			break
		}
	}
	if post == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	t := template.Must(
		template.ParseFiles("./templates/get.html"),
	)
	err = t.Execute(
		w,
		post,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PostPayload struct {
	Title string `json:"title"`
	Body string `json:"body"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	var payload *PostPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	latest := posts[len(posts) - 1]
	post := Post{
		ID: latest.ID + 1,
		Title: payload.Title,
		Body: payload.Body,
	}
	posts = append(posts, post)

	// redirect to list
	http.Redirect(w, r, "/", http.StatusCreated)
}
