package helper

import (
	"context"
	"encoding/json"
	"golang_template/db"
	"golang_template/models"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ConvertObjectId(w http.ResponseWriter, id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal error: Invalid ID format")
		return primitive.NilObjectID
	}
	return objectID
}

func JSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error to response")
		return
	}
}

func JSONDataResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Status:  true,
		Message: "success",
		Data:    data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Error to response")
		return
	}
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"error": message})
}

func ValidateRegex(pattern, value string) bool {
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorResponse(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token format")
			return
		}

		valid, err := isValidSession(tokenString)
		if err != nil || !valid {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidSession(token string) (bool, error) {
	collection := db.CollSessions()

	var session models.Session
	err := collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return false, err
	}

	sessionExpiration := 24 * time.Hour

	if time.Since(session.CreatedAt) > sessionExpiration {
		removeSession(token)
		return false, nil
	}

	return true, nil
}

func removeSession(token string) error {
	collection := db.CollSessions()

	_, err := collection.DeleteOne(context.Background(), bson.M{"token": token})
	return err
}

func EncryptPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error encrypting password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Error comparing passwords:", err)
		return false
	}
	return true
}
