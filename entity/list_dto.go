package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type BatchRequestPayload struct {
	Ids []primitive.ObjectID `form:"ids" json:"ids"`
}
