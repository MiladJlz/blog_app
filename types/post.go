package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	minContentLen = 10
)

type PathParameter struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"66db2c856699531daa9abc16"`
}
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"66db2c856699531daa9abc16"`
	Content   string             `bson:"content" json:"content" example:"This is example."`
	Author    primitive.ObjectID `bson:"author" json:"author" example:"66db21cdb5d96466fa5f3c3c"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at" example:"2024-09-06T16:23:33.648Z"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at" example:"2024-09-06T16:23:33.648Z"`
}

type CreatePostParams struct {
	Content   string    `json:"content" example:"This is example."`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author" example:"66db2c856699531daa9abc16"`
}
type UpdatePostParams struct {
	Content string `json:"content" example:"This is example."`
}

func (p UpdatePostParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.Content) > 0 {
		m["content"] = p.Content
	}
	m["updated_at"] = time.Now()
	return m
}
func (params CreatePostParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.Content) < minContentLen {
		errors["content"] = fmt.Sprintf("content length should be at least %d characters", minContentLen)
	}

	return errors
}
func NewPostFromParams(params CreatePostParams) *Post {
	oid, _ := primitive.ObjectIDFromHex(params.Author)

	return &Post{
		Content:   params.Content,
		Author:    oid,
		CreatedAt: time.Now(),
	}
}
