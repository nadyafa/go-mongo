package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nadyafa/go-mongo/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{c}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Id = primitive.NewObjectID()
	collection := uc.client.Database("mongo-golang").Collection("users")

	if _, err := collection.InsertOne(context.TODO(), u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

func (uc UserController) GetUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var user models.User
	collection := uc.client.Database("mongo-golang").Collection("users")

	if err := collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

func (uc UserController) UpdateUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var updateUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")

	update := bson.M{"$set": updateUser}
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": oid}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Update user with id: %s\n", oid.Hex())
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(r.Context(), bson.M{"_id": oid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user with idL %s\n", oid.Hex())
}
