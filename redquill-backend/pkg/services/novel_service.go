// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: novel_service.go
/@Description: Novel service implementation
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

// NovelService 小说服务
type NovelService struct {
	client *mongo.Client
	dbName string
}

// NewNovelService 创建小说服务
func NewNovelService(client *mongo.Client, dbName string) *NovelService {
	return &NovelService{
		client: client,
		dbName: dbName,
	}
}

// PostNovels 创建小说
func (s *NovelService) PostNovels(ctx context.Context, title, authorID, status, currentPhase string, blueprint models.ProjectBlueprint, aiContext models.AIContext) (models.Novel, error) {
	coll := s.client.Database(s.dbName).Collection("novels")

	now := time.Now()
	novel := models.Novel{
		Title:            title,
		AuthorID:         authorID,
		Status:           status,
		CurrentPhase:     currentPhase,
		Ctime:            now.Unix(),
		Mtime:            now.Unix(),
		ProjectBlueprint: blueprint,
		AIContext:        aiContext,
		ExtraInfo:        make(map[string]interface{}), // 初始化ExtraInfo
	}

	res, err := coll.InsertOne(ctx, novel)
	if err != nil {
		return models.Novel{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		novel.ID = oid.Hex()
	}

	return novel, nil
}

// GetNovels 获取小说详情
func (s *NovelService) GetNovels(ctx context.Context, id string) (models.Novel, error) {
	coll := s.client.Database(s.dbName).Collection("novels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Novel{}, errors.New("invalid id")
	}

	var novel models.Novel
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&novel); err != nil {
		return models.Novel{}, err
	}

	return novel, nil
}

// PutNovels 更新小说
func (s *NovelService) PutNovels(ctx context.Context, id string, title *string, status *string, currentPhase *string, blueprint *models.ProjectBlueprint, aiContext *models.AIContext, extraInfo *map[string]interface{}) (models.Novel, error) {
	coll := s.client.Database(s.dbName).Collection("novels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Novel{}, errors.New("invalid id")
	}

	update := bson.M{"mtime": time.Now().Unix()}
	set := bson.M{}

	if title != nil {
		set["title"] = *title
	}
	if status != nil {
		set["status"] = *status
	}
	if currentPhase != nil {
		set["current_phase"] = *currentPhase
	}
	if blueprint != nil {
		set["project_blueprint"] = *blueprint
	}
	if aiContext != nil {
		set["ai_context"] = *aiContext
	}
	if extraInfo != nil {
		set["extra_info"] = *extraInfo
	}

	for k, v := range set {
		update[k] = v
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out models.Novel
	if err := coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&out); err != nil {
		return models.Novel{}, err
	}

	return out, nil
}

// UpdateNovelExtraInfo 更新小说ExtraInfo
func (s *NovelService) UpdateNovelExtraInfo(ctx context.Context, id string, phase string, data interface{}) error {
	coll := s.client.Database(s.dbName).Collection("novels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	// 更新ExtraInfo字段
	update := bson.M{
		"$set": bson.M{
			"extra_info." + phase: data,
			"mtime":               time.Now().Unix(),
		},
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

// GetNovelExtraInfo 获取小说ExtraInfo
func (s *NovelService) GetNovelExtraInfo(ctx context.Context, id string, phase string) (interface{}, error) {
	coll := s.client.Database(s.dbName).Collection("novels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var novel models.Novel
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&novel); err != nil {
		return nil, err
	}

	if novel.ExtraInfo == nil {
		return nil, nil
	}

	return novel.ExtraInfo[phase], nil
}

// DeleteNovels 删除小说
func (s *NovelService) DeleteNovels(ctx context.Context, id string) error {
	coll := s.client.Database(s.dbName).Collection("novels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ListNovels 分页查询小说列表
func (s *NovelService) ListNovels(ctx context.Context, page, pageSize int64, sortExpr, keyword string) (PagedNovels, error) {
	coll := s.client.Database(s.dbName).Collection("novels")

	// 构建过滤条件
	kwFilter := common.BuildKeywordFilter(keyword, []string{"title", "project_blueprint.genre", "project_blueprint.core_conflict"})
	filter := common.MergeFilters(bson.M{}, kwFilter)

	// 构建选项
	sort := common.BuildSort(sortExpr)
	opts := common.BuildFindOptions(page, pageSize, sort, bson.M{})

	items, total, err := common.FindWithPagination[models.Novel](ctx, coll, filter, opts)
	if err != nil {
		return PagedNovels{}, err
	}

	totalPages := total / common.NormalizePageSize(pageSize)
	if total%common.NormalizePageSize(pageSize) != 0 {
		totalPages++
	}

	return PagedNovels{
		Items: items,
		Pagination: common.Pagination{
			Page:      common.NormalizePage(page),
			PageSize:  common.NormalizePageSize(pageSize),
			Total:     total,
			TotalPage: totalPages,
		},
	}, nil
}

// PagedNovels 分页小说结果
type PagedNovels struct {
	Items      []models.Novel    `json:"items"`
	Pagination common.Pagination `json:"pagination"`
}

// PostStoryCores 创建故事核心
func (s *NovelService) PostStoryCores(ctx context.Context, novelID, title, coreConflict, theme, innovation, commercialPotential, targetAudience string) (models.StoryCore, error) {
	coll := s.client.Database(s.dbName).Collection("story_cores")

	now := time.Now()
	storyCore := models.StoryCore{
		NovelID:             novelID,
		Title:               title,
		CoreConflict:        coreConflict,
		Theme:               theme,
		Innovation:          innovation,
		CommercialPotential: commercialPotential,
		TargetAudience:      targetAudience,
		Ctime:               now.Unix(),
	}

	res, err := coll.InsertOne(ctx, storyCore)
	if err != nil {
		return models.StoryCore{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		storyCore.ID = oid.Hex()
	}

	return storyCore, nil
}

// GetStoryCores 获取故事核心
func (s *NovelService) GetStoryCores(ctx context.Context, novelID string) ([]models.StoryCore, error) {
	coll := s.client.Database(s.dbName).Collection("story_cores")

	cursor, err := coll.Find(ctx, bson.M{"novel_id": novelID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var storyCores []models.StoryCore
	if err := cursor.All(ctx, &storyCores); err != nil {
		return nil, err
	}

	return storyCores, nil
}

// PostWorldviews 创建世界观
func (s *NovelService) PostWorldviews(ctx context.Context, novelID string, powerSystem models.PowerSystem, societyStructure models.SocietyStructure, geography models.Geography, specialRules []string) (models.Worldview, error) {
	coll := s.client.Database(s.dbName).Collection("worldviews")

	now := time.Now()
	worldview := models.Worldview{
		NovelID:          novelID,
		PowerSystem:      powerSystem,
		SocietyStructure: societyStructure,
		Geography:        geography,
		SpecialRules:     specialRules,
		Ctime:            now.Unix(),
	}

	res, err := coll.InsertOne(ctx, worldview)
	if err != nil {
		return models.Worldview{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		worldview.ID = oid.Hex()
	}

	return worldview, nil
}

// GetWorldviews 获取世界观
func (s *NovelService) GetWorldviews(ctx context.Context, novelID string) (models.Worldview, error) {
	coll := s.client.Database(s.dbName).Collection("worldviews")

	var worldview models.Worldview
	if err := coll.FindOne(ctx, bson.M{"novel_id": novelID}).Decode(&worldview); err != nil {
		return models.Worldview{}, err
	}

	return worldview, nil
}

// PostCharacters 创建角色
func (s *NovelService) PostCharacters(ctx context.Context, novelID, name, charType string, coreAttributes models.CoreAttributes, soulProfile models.SoulProfile) (models.Character, error) {
	coll := s.client.Database(s.dbName).Collection("characters")

	now := time.Now()
	character := models.Character{
		NovelID:        novelID,
		Name:           name,
		Type:           charType,
		CoreAttributes: coreAttributes,
		SoulProfile:    soulProfile,
		GrowthTrack:    []models.GrowthEvent{},
		Ctime:          now.Unix(),
	}

	res, err := coll.InsertOne(ctx, character)
	if err != nil {
		return models.Character{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		character.ID = oid.Hex()
	}

	return character, nil
}

// GetCharacters 获取角色列表
func (s *NovelService) GetCharacters(ctx context.Context, novelID string) ([]models.Character, error) {
	coll := s.client.Database(s.dbName).Collection("characters")

	cursor, err := coll.Find(ctx, bson.M{"novel_id": novelID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var characters []models.Character
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, err
	}

	return characters, nil
}

// PostChapters 创建章节
func (s *NovelService) PostChapters(ctx context.Context, novelID string, chapterNumber int, title, content, summary string, outline models.ChapterOutline, qualityMetrics models.QualityMetrics, characterDevelopment map[string]string) (models.Chapter, error) {
	coll := s.client.Database(s.dbName).Collection("chapters")

	now := time.Now()
	chapter := models.Chapter{
		NovelID:              novelID,
		ChapterNumber:        chapterNumber,
		Title:                title,
		Content:              content,
		WordCount:            len([]rune(content)),
		Summary:              summary,
		Outline:              outline,
		QualityMetrics:       qualityMetrics,
		CharacterDevelopment: characterDevelopment,
		Ctime:                now.Unix(),
	}

	res, err := coll.InsertOne(ctx, chapter)
	if err != nil {
		return models.Chapter{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		chapter.ID = oid.Hex()
	}

	return chapter, nil
}

// GetChapters 获取章节列表
func (s *NovelService) GetChapters(ctx context.Context, novelID string) ([]models.Chapter, error) {
	coll := s.client.Database(s.dbName).Collection("chapters")

	cursor, err := coll.Find(ctx, bson.M{"novel_id": novelID}, options.Find().SetSort(bson.M{"chapter_number": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chapters []models.Chapter
	if err := cursor.All(ctx, &chapters); err != nil {
		return nil, err
	}

	return chapters, nil
}

// GetChapter 获取单个章节
func (s *NovelService) GetChapter(ctx context.Context, id string) (models.Chapter, error) {
	coll := s.client.Database(s.dbName).Collection("chapters")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Chapter{}, errors.New("invalid id")
	}

	var chapter models.Chapter
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&chapter); err != nil {
		return models.Chapter{}, err
	}

	return chapter, nil
}

// PostWritingSessions 创建创作会话
func (s *NovelService) PostWritingSessions(ctx context.Context, novelID string, currentChapter int, sessionContext models.SessionContext) (models.WritingSession, error) {
	coll := s.client.Database(s.dbName).Collection("writing_sessions")

	now := time.Now()
	session := models.WritingSession{
		NovelID:        novelID,
		CurrentChapter: currentChapter,
		SessionContext: sessionContext,
		LastUpdated:    now.Unix(),
	}

	res, err := coll.InsertOne(ctx, session)
	if err != nil {
		return models.WritingSession{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		session.ID = oid.Hex()
	}

	return session, nil
}

// GetWritingSessions 获取创作会话
func (s *NovelService) GetWritingSessions(ctx context.Context, novelID string) (models.WritingSession, error) {
	coll := s.client.Database(s.dbName).Collection("writing_sessions")

	var session models.WritingSession
	if err := coll.FindOne(ctx, bson.M{"novel_id": novelID}).Decode(&session); err != nil {
		return models.WritingSession{}, err
	}

	return session, nil
}
