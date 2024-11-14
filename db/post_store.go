package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"note_app/types"
	"os"
)

const postColl = "posts"

type Map map[string]any

type PostStore interface {
	InsertPost(context.Context, *types.Post) (*types.Post, error)
	UpdatePost(ctx context.Context, filter Map, params types.UpdatePostParams) error
	DeletePost(context.Context, string) error
	GetPosts(context.Context) ([]*types.Post, error)
	GetPostByID(context.Context, string) (*types.Post, error)
	GetPostsByUserID(context.Context, string) ([]*types.Post, error)
}
type MongoPostStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoPostStore(client *mongo.Client) *MongoPostStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoPostStore{
		client: client,
		coll:   client.Database(dbname).Collection(postColl),
	}
}

func (s *MongoPostStore) UpdatePost(ctx context.Context, filter Map, params types.UpdatePostParams) error {
	oid, _ := primitive.ObjectIDFromHex(filter["_id"].(string))

	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoPostStore) DeletePost(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)

	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoPostStore) InsertPost(ctx context.Context, post *types.Post) (*types.Post, error) {
	res, err := s.coll.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	post.ID = res.InsertedID.(primitive.ObjectID)
	return post, nil
}

func (s *MongoPostStore) GetPosts(ctx context.Context) ([]*types.Post, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var posts []*types.Post
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *MongoPostStore) GetPostByID(ctx context.Context, id string) (*types.Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var post types.Post
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&post); err != nil {
		return nil, err
	}
	return &post, nil
}
func (s *MongoPostStore) GetPostsByUserID(ctx context.Context, id string) ([]*types.Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cur, err := s.coll.Find(ctx, bson.M{"author": oid})
	if err != nil {
		return nil, err
	}
	var posts []*types.Post
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}
