package domain

import (
	"time"
)

type User struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	Name         string       `json:"name" bson:"name"`
	Email        string       `json:"email" bson:"email"`
	Phone        string       `json:"phone" bson:"phone"`
	Password     string       `json:"password" bson:"password"`
	RegisteredAt time.Time    `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time    `json:"lastVisitAt" bson:"lastVisitAt"`
	Verification Verification `json:"verification" bson:"verification"`
	Trainers     []string     `json:"trainers" bson:"trainers"`
}
