package mongodb

import (
	"context"
	"fmt"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type redirectURLRepository struct {
	collection *mongo.Collection
}

func NewRedirectURLRepository(mongodb *mongo.Database, collectionName string) interfaces.RedirectURLRepository {
	return &redirectURLRepository{
		collection: mongodb.Collection(collectionName),
	}
}

func (r *redirectURLRepository) Store(ctx context.Context, in entities.RedirectURL) (*entities.RedirectURL, error) {
	_, err := r.collection.InsertOne(ctx, in)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *redirectURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*entities.RedirectURL, error) {
	filter := map[string]interface{}{"short_code": shortCode}
	var result entities.RedirectURL
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to find document: %v", err)
	}
	return &result, nil
}

func (r *redirectURLRepository) GetByUUID(ctx context.Context, uuid string) (*entities.RedirectURL, error) {
	filter := map[string]interface{}{"uuid": uuid}
	var result entities.RedirectURL
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to find document: %v", err)
	}
	return &result, nil
}

func (r *redirectURLRepository) Update(ctx context.Context, uuid string, in entities.RedirectURL) (*entities.RedirectURL, error) {
	filter := map[string]interface{}{"uuid": uuid}
	update := map[string]interface{}{
		"$set": in,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return &in, nil
}

func (r *redirectURLRepository) Delete(ctx context.Context, uuid string) error {
	filter := map[string]interface{}{"uuid": uuid}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entities.ErrURLNotFound
		}
		return fmt.Errorf("failed to delete document: %v", err)
	}
	return nil
}
