package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang_template/db"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"golang_template/helper"
	"golang_template/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds models.Users
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	collection := db.CollUsers()
	var user models.Users
	err := collection.FindOne(context.Background(), bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
	}

	err = storeSession(session)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error storing session")
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "SignIn successful",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

func generateJWT(userID primitive.ObjectID) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func storeSession(session models.Session) error {

	collection := db.CollSessions()

	filter := bson.M{"user_id": session.UserID}
	var existingSession models.Session
	err := collection.FindOne(context.Background(), filter).Decode(&existingSession)

	if err == nil {

		update := bson.M{
			"$set": bson.M{
				"token":      session.Token,
				"created_at": session.CreatedAt,
			},
		}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		return err
	} else if err != mongo.ErrNoDocuments {

		return err
	}

	_, err = collection.InsertOne(context.Background(), session)
	return err
}

func SignOut(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		helper.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		helper.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token format")
		return
	}

	err := removeSession(tokenString)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error signing out")
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "SignOut successful",
	})
}

func removeSession(token string) error {
	collection := db.CollSessions()

	_, err := collection.DeleteOne(context.Background(), bson.M{"token": token})
	return err
}

func CleanupExpiredSessions(w http.ResponseWriter, r *http.Request) {
	collection := db.CollSessions()

	sessionExpiration := time.Minute

	_, err := collection.DeleteMany(
		context.Background(),
		bson.M{
			"created_at": bson.M{"$lt": time.Now().Add(-sessionExpiration)},
		},
	)
	if err != nil {
		log.Println("Error cleaning up expired sessions:", err)
	}
}
