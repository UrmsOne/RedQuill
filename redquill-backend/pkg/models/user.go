// Package models
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: user.go
/@Description:
/*/

package models

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"` // hashed
	Ctime    int64  `json:"ctime" bson:"ctime"`
	Mtime    int64  `json:"mtime" bson:"mtime"`
}
