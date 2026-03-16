package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"url-shortening-service/internal/services/url_shortening"
	"url-shortening-service/models"
	"url-shortening-service/storage/mongodb"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UP and Running...")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code := url_shortening.CodeGenerator()

	entry := models.URLEntry{
		ID:          primitive.NewObjectID(),
		URL:         req.URL,
		ShortCode:   code,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessCount: 0,
	}

	collection := mongodb.GetCollection()
	_, err = collection.InsertOne(context.Background(), entry)

	if err != nil {
		http.Error(w, "failed to create short url", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	var entry models.URLEntry

	collection := mongodb.GetCollection()
	filter := bson.M{"shortcode": code}
	err := collection.FindOne(context.Background(), filter).Decode(&entry)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get url", http.StatusInternalServerError)
		return
	}

	update := bson.M{"$inc": bson.M{"accesscount": 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		http.Error(w, "failed to update url access count", http.StatusInternalServerError)
		return
	}

	entry.AccessCount += 1

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(entry)
}

func UpdateURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	var req models.ShortenRequest
	var entry models.URLEntry
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil || req.URL == "" {
		http.Error(w, "bad request or url is empty", http.StatusBadRequest)
		return
	}

	collection := mongodb.GetCollection()
	filter := bson.M{"shortcode": code}
	err = collection.FindOne(context.Background(), filter).Decode(&entry)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	update := bson.M{"$set": bson.M{"url": req.URL, "updatedAt": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "failed to update url", http.StatusInternalServerError)
		return
	}

	err = collection.FindOne(context.Background(), filter).Decode(&entry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(entry)
}

func DeleteURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	collection := mongodb.GetCollection()
	filter := bson.M{"shortcode": code}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "failed to delete url", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(w, "failed to delete url", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	var entry models.URLEntry

	collection := mongodb.GetCollection()
	filter := bson.M{"shortcode": code}
	err := collection.FindOne(context.Background(), filter).Decode(&entry)
	if err != nil {
		http.Error(w, "failed to get url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(entry)
}
