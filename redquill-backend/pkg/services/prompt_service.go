// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_service.go
/@Description: Prompt service implementation
/*/

package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"redquill-backend/pkg/common"
	"redquill-backend/pkg/models"
)

// PromptService Prompt服务
type PromptService struct {
	client *mongo.Client
	dbName string
}

// NewPromptService 创建Prompt服务
func NewPromptService(client *mongo.Client, dbName string) *PromptService {
	return &PromptService{
		client: client,
		dbName: dbName,
	}
}

// PostPrompts 创建Prompt
func (s *PromptService) PostPrompts(ctx context.Context, name, promptType, category, description, content string, variables, tags []string, public bool, creatorID, username string) (models.Prompt, error) {
	coll := s.client.Database(s.dbName).Collection("prompts")
	
	// 检查Prompt名称是否已存在
	var existing models.Prompt
	if err := coll.FindOne(ctx, bson.M{"name": name, "creator_id": creatorID}).Decode(&existing); err == nil {
		return models.Prompt{}, errors.New("prompt name already exists for this user")
	}
	
	now := time.Now()
	prompt := models.Prompt{
		Name:        name,
		Type:        promptType,
		Category:    category,
		Description: description,
		Content:     content,
		Variables:   variables,
		UsageCount:  0,
		Tags:        tags,
		Public:      public,
		CreatorID:   creatorID,
		Username:    username,
		Ctime:       now.Unix(),
		Mtime:       now.Unix(),
	}
	
	res, err := coll.InsertOne(ctx, prompt)
	if err != nil {
		return models.Prompt{}, err
	}
	
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		prompt.ID = oid.Hex()
	}
	
	return prompt, nil
}

// GetPrompts 获取Prompt详情
func (s *PromptService) GetPrompts(ctx context.Context, id string) (models.Prompt, error) {
	coll := s.client.Database(s.dbName).Collection("prompts")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Prompt{}, errors.New("invalid id")
	}
	
	var prompt models.Prompt
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&prompt); err != nil {
		return models.Prompt{}, err
	}
	
	return prompt, nil
}

// PutPrompts 更新Prompt
func (s *PromptService) PutPrompts(ctx context.Context, id string, name *string, promptType *string, category *string, description *string, content *string, variables *[]string, tags *[]string, public *bool) (models.Prompt, error) {
	coll := s.client.Database(s.dbName).Collection("prompts")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Prompt{}, errors.New("invalid id")
	}
	
	update := bson.M{"mtime": time.Now().Unix()}
	set := bson.M{}
	
	if name != nil {
		set["name"] = *name
	}
	if promptType != nil {
		set["type"] = *promptType
	}
	if category != nil {
		set["category"] = *category
	}
	if description != nil {
		set["description"] = *description
	}
	if content != nil {
		set["content"] = *content
	}
	if variables != nil {
		set["variables"] = *variables
	}
	if tags != nil {
		set["tags"] = *tags
	}
	if public != nil {
		set["public"] = *public
	}
	
	for k, v := range set {
		update[k] = v
	}
	
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out models.Prompt
	if err := coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&out); err != nil {
		return models.Prompt{}, err
	}
	
	return out, nil
}

// DeletePrompts 删除Prompt
func (s *PromptService) DeletePrompts(ctx context.Context, id string) error {
	coll := s.client.Database(s.dbName).Collection("prompts")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	
	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ListPrompts 分页查询Prompt列表
func (s *PromptService) ListPrompts(ctx context.Context, page, pageSize int64, sortExpr, keyword string) (PagedPrompts, error) {
	coll := s.client.Database(s.dbName).Collection("prompts")
	
	// 构建过滤条件
	kwFilter := common.BuildKeywordFilter(keyword, []string{"name", "description", "content", "category"})
	filter := common.MergeFilters(bson.M{}, kwFilter)
	
	// 构建选项
	sort := common.BuildSort(sortExpr)
	opts := common.BuildFindOptions(page, pageSize, sort, bson.M{})
	
	items, total, err := common.FindWithPagination[models.Prompt](ctx, coll, filter, opts)
	if err != nil {
		return PagedPrompts{}, err
	}
	
	totalPages := total / common.NormalizePageSize(pageSize)
	if total%common.NormalizePageSize(pageSize) != 0 {
		totalPages++
	}
	
	return PagedPrompts{
		Items: items,
		Pagination: common.Pagination{
			Page:      common.NormalizePage(page),
			PageSize:  common.NormalizePageSize(pageSize),
			Total:     total,
			TotalPage: totalPages,
		},
	}, nil
}

// PagedPrompts 分页Prompt结果
type PagedPrompts struct {
	Items      []models.Prompt `json:"items"`
	Pagination common.Pagination `json:"pagination"`
}
