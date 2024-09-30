package model

import "go.mongodb.org/mongo-driver/bson/primitive"

/*  NOTE
	struct tags should be separated by one space.
*/

type Model interface{}

type IdObj struct {
	Id primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
}

type PublicUser struct {
	Name        string   `json:"name" bson:"name,omitempty"`
	Description string   `json:"description" bson:"description,omitempty"`
	Skills      []string `json:"skills" bson:"skills,omitempty"`
}

type User struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Email       string             `json:"email" bson:"email,omitempty"`
	Password    string             `json:"password" bson:"password,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Skills      []string           `json:"skills" bson:"skills,omitempty"`
}

type ProjectCheck struct {
	OwnerId primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
}

type Project struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	OwnerId  	primitive.ObjectID   `json:"ownerId" bson:"ownerId,omitempty"`
	Name        string   			 `json:"name" bson:"name,omitempty"`
	Description string   			 `json:"description" bson:"description,omitempty"`
	Category    string   			 `json:"category" bson:"category,omitempty"`
	Tags        []string 			 `json:"tags" bson:"tags,omitempty"`
	Seeking     []string 			 `json:"seeking" bson:"seeking,omitempty"`
	// DateCreated string
	// Creator primitive.ObjectID
	// Admins []string
	//  []string
	// Location    string   `json:"location"`
}

type ProjectRequest struct {
	Id        primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID   `json:"userId" bson:"userId,omitempty"`
	ProjectId primitive.ObjectID   `json:"projectId" bson:"projectId,omitempty"`
	Message   string 			   `json:"message" bson:"message,omitempty"`
	// NumVotes int
	// VotesNeeded == |Project.admins|
	// Kind: “application” || “invite”
	// Rejected: bool
	}