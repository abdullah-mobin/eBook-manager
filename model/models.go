package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Book struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		Name   string             `bson:"name"`
		Author string             `bson:"author"`
		Type   string             `bson:"type"`
		PDF    string             `bson:"pdf"`
	}
	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty"`
		Name     string             `bson:"name"`
		Email    string             `bson:"email"`
		Password string             `bson:"password"`
		UserType string             `bson:"usertype"`
	}
	GUser struct {
		Username  string `json:"name"`
		Useremail string `json:"email"`
	}
)
