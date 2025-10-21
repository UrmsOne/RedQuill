// Package models
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_model.go
/@Description: Prompt model data structure
/*/

package models

// Prompt Prompt模板
type Prompt struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"`
	Category    string   `json:"category" bson:"category"`
	Description string   `json:"description" bson:"description"`
	Content     string   `json:"content" bson:"content"`
	Variables   []string `json:"variables" bson:"variables"`
	UsageCount  int64    `json:"usage_count" bson:"usage_count"`
	Tags        []string `json:"tags" bson:"tags"`
	Public      bool     `json:"public" bson:"public"`
	CreatorID   string   `json:"creator_id" bson:"creator_id"`
	Username    string   `json:"username" bson:"username"`
	Ctime       int64    `json:"ctime" bson:"ctime"`
	Mtime       int64    `json:"mtime" bson:"mtime"`
}
