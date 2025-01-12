package users

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_template/db"
	"golang_template/helper"
	"golang_template/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"message": "I'm good!"})
}

func List(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}
	cursor, err := db.CollUsers().Find(context.Background(), filter)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error fetching users")
		return
	}
	defer cursor.Close(context.Background())

	users := []bson.M{}
	if err := cursor.All(context.Background(), &users); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error parsing users")
		return
	}

	helper.JSONDataResponse(w, http.StatusOK, users)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := helper.ConvertObjectId(w, vars["id"])

	filter := bson.M{"_id": objectID}
	fmt.Println(filter)

	var result models.Users
	err := db.CollUsers().FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		helper.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	helper.JSONDataResponse(w, http.StatusOK, result)
}

func Create(w http.ResponseWriter, r *http.Request) {

	var newUser models.Users
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Username == "" || newUser.Password == "" {
		helper.ErrorResponse(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	hashedPassword, err := helper.EncryptPassword(newUser.Password)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error encrypting password")
		return
	}
	newUser.Password = string(hashedPassword)

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	collection := db.CollUsers()
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	helper.JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"data":    newUser,
	})
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := helper.ConvertObjectId(w, vars["id"])

	var updatedUser models.Users
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateFields := bson.M{}

	if updatedUser.Name != "" {
		updateFields["name"] = updatedUser.Name
	}

	if updatedUser.Email != "" {
		updateFields["email"] = updatedUser.Email
	}

	if updatedUser.Username != "" {
		updateFields["username"] = updatedUser.Username
	}

	if updatedUser.Password != "" {

		encryptedPassword, err := helper.EncryptPassword(updatedUser.Password)
		if err != nil {
			helper.ErrorResponse(w, http.StatusInternalServerError, "Error encrypting password")
			return
		}
		updateFields["password"] = encryptedPassword
	}

	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 0 {
		helper.ErrorResponse(w, http.StatusBadRequest, "No fields to update")
		return
	}

	collection := db.CollUsers()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	if result.MatchedCount == 0 {
		helper.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	var updatedUserFromDb models.Users
	err = collection.FindOne(context.Background(), filter).Decode(&updatedUserFromDb)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error retrieving updated user")
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"data":    updatedUserFromDb,
	})
}
