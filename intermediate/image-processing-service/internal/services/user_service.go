package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserService struct {
	db     *database.MongoDB
	logger *zap.Logger
}

func NewUserService(db *database.MongoDB, logger *zap.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := s.db.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		s.logger.Error("Failed to insert user", logger.Error(err))
		return fmt.Errorf("failed to insert user: %w", err)
	}

	s.logger.Info("User created successfully", logger.String("id", user.ID.Hex()))
	return nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.logger.Error("Failed to get user by email", logger.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	collection := s.db.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.logger.Error("Failed to get user", logger.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id string, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid user ID", logger.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user.UpdatedAt = time.Now()
	collection := s.db.Collection("users")

	update := bson.M{
		"$set": bson.M{
			"name": user.Name,
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		s.logger.Error("Failed to update user", logger.Error(err))
		return err
	}

	s.logger.Info("User updated successfully", logger.String("user_id", id))
	return nil
}

// EmailExists checks if the email already exists in the database
func (s *UserService) EmailExists(email string) (bool, error) {
	collection := s.db.Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.M{"email": email})
	if err != nil {
		s.logger.Error("Failed to count users", logger.Error(err))
		return false, err
	}

	return count > 0, nil
}
