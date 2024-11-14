package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"note_app/types"
	"os"
)

const userColl = "users"

type UserStore interface {
	GetUser(context.Context, string) (*types.User, error)
	GetUserByObjectID(context.Context, primitive.ObjectID) (*types.User, error)

	UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error
	DeleteUser(context.Context, string) error
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	AddFriend(context.Context, Map, string, string) error
	RemoveFriend(context.Context, Map, string, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}
	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetUser(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *MongoUserStore) GetUserByObjectID(ctx context.Context, id primitive.ObjectID) (*types.User, error) {

	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *MongoUserStore) AddFriend(ctx context.Context, filter Map, id string, userID string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter["_id"] = oid

	update := bson.M{"$addToSet": bson.M{"friends": uid}}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoUserStore) RemoveFriend(ctx context.Context, filter Map, id string, userID string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter["_id"] = oid

	update := bson.M{"$pull": bson.M{"friends": uid}}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}
