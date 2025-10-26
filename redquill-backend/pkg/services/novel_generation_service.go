// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: novel_generation_service.go
/@Description: Novel generation service implementation
/*/

package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"redquill-backend/pkg/models"
)

// NovelGenerationService 小说生成服务
type NovelGenerationService struct {
	client *mongo.Client
	dbName string
}

// NewNovelGenerationService 创建小说生成服务
func NewNovelGenerationService(client *mongo.Client, dbName string) *NovelGenerationService {
	return &NovelGenerationService{
		client: client,
		dbName: dbName,
	}
}

// GenerateStoryCore 生成故事核心
func (s *NovelGenerationService) GenerateStoryCore(ctx context.Context, novelID, llmModelID string, inputData map[string]interface{}) (models.StoryCore, error) {
	// 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:      novelID,
		LLMModelID:   llmModelID,
		InputData:    inputData,
		TemplateType: "story_core",
		Stream:       false,
	}

	// 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return models.StoryCore{}, err
	}

	if !response.Success {
		return models.StoryCore{}, errors.New(response.Error)
	}

	// 解析响应数据
	concepts, ok := response.Data["concepts"].([]interface{})
	if !ok || len(concepts) == 0 {
		return models.StoryCore{}, errors.New("no concepts generated")
	}

	// 取第一个概念作为故事核心
	concept := concepts[0].(map[string]interface{})

	storyCore := models.StoryCore{
		NovelID:             novelID,
		Title:               s.getString(concept, "title"),
		CoreConflict:        s.getString(concept, "core_conflict"),
		Theme:               s.getString(concept, "theme"),
		Innovation:          s.getString(concept, "innovation"),
		CommercialPotential: s.getString(concept, "commercial_potential"),
		TargetAudience:      s.getString(concept, "target_audience"),
	}

	// 保存到数据库
	novelService := NewNovelService(s.client, s.dbName)
	storyCore, err = novelService.PostStoryCores(ctx, storyCore.NovelID, storyCore.Title, storyCore.CoreConflict, storyCore.Theme, storyCore.Innovation, storyCore.CommercialPotential, storyCore.TargetAudience)
	if err != nil {
		return models.StoryCore{}, err
	}

	// 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "story_core", extraInfo); err != nil {
		// 记录错误但不影响主流程
		// log.Printf("Failed to update extra info: %v", err)
	}

	return storyCore, nil
}

// GenerateWorldview 生成世界观
func (s *NovelGenerationService) GenerateWorldview(ctx context.Context, novelID, llmModelID string, inputData map[string]interface{}) (models.Worldview, error) {
	// 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:      novelID,
		LLMModelID:   llmModelID,
		InputData:    inputData,
		TemplateType: "worldview",
		Stream:       false,
	}

	// 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return models.Worldview{}, err
	}

	if !response.Success {
		return models.Worldview{}, errors.New(response.Error)
	}

	// 解析响应数据
	worldviewData := response.Data

	powerSystem := models.PowerSystem{
		Name:              s.getString(worldviewData, "power_system.name"),
		Levels:            s.getStringArray(worldviewData, "power_system.levels"),
		CultivationMethod: s.getString(worldviewData, "power_system.cultivation_method"),
		Limitations:       s.getString(worldviewData, "power_system.limitations"),
	}

	societyStructure := models.SocietyStructure{
		Hierarchy:      s.getString(worldviewData, "society_structure.hierarchy"),
		EconomicSystem: s.getString(worldviewData, "society_structure.economic_system"),
		MajorFactions:  s.parseFactions(worldviewData),
	}

	geography := models.Geography{
		MajorRegions:     s.getStringArray(worldviewData, "geography.major_regions"),
		SpecialLocations: s.getStringArray(worldviewData, "geography.special_locations"),
	}

	specialRules := s.getStringArray(worldviewData, "special_rules")

	// 保存到数据库
	novelService := NewNovelService(s.client, s.dbName)
	worldview, err := novelService.PostWorldviews(ctx, novelID, powerSystem, societyStructure, geography, specialRules)
	if err != nil {
		return models.Worldview{}, err
	}

	// 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "worldview", extraInfo); err != nil {
		// 记录错误但不影响主流程
		// log.Printf("Failed to update extra info: %v", err)
	}

	return worldview, nil
}

// GenerateCharacter 生成角色
func (s *NovelGenerationService) GenerateCharacter(ctx context.Context, novelID, llmModelID string, inputData map[string]interface{}) (models.Character, error) {
	// 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:      novelID,
		LLMModelID:   llmModelID,
		InputData:    inputData,
		TemplateType: "character",
		Stream:       false,
	}

	// 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return models.Character{}, err
	}

	if !response.Success {
		return models.Character{}, errors.New(response.Error)
	}

	// 解析响应数据
	characterData := response.Data

	// 解析灵魂档案
	soulProfile := models.SoulProfile{
		Personality: models.Personality{
			CoreTraits:        s.getStringArray(characterData, "soul_profile.personality.core_traits"),
			MoralCompass:      s.getString(characterData, "soul_profile.personality.moral_compass"),
			InternalConflicts: s.getStringArray(characterData, "soul_profile.personality.internal_conflicts"),
			Fears:             s.getStringArray(characterData, "soul_profile.personality.fears"),
			Desires:           s.getStringArray(characterData, "soul_profile.personality.desires"),
		},
		Background: models.Background{
			Origin:         s.getString(characterData, "soul_profile.background.origin"),
			DefiningEvents: s.getStringArray(characterData, "soul_profile.background.defining_events"),
			HiddenSecrets:  s.getStringArray(characterData, "soul_profile.background.hidden_secrets"),
		},
		Motivations: models.Motivations{
			ImmediateGoal: s.getString(characterData, "soul_profile.motivations.immediate_goal"),
			LongTermGoal:  s.getString(characterData, "soul_profile.motivations.long_term_goal"),
			CoreDrive:     s.getString(characterData, "soul_profile.motivations.core_drive"),
		},
	}

	// 解析核心属性
	coreAttributes := models.CoreAttributes{
		CultivationLevel: s.getString(characterData, "core_attributes.cultivation_level"),
		CurrentItems:     s.getStringArray(characterData, "core_attributes.current_items"),
		Abilities:        s.getStringArray(characterData, "core_attributes.abilities"),
		Relationships:    s.parseRelationships(characterData),
	}

	// 保存到数据库
	novelService := NewNovelService(s.client, s.dbName)
	character, err := novelService.PostCharacters(ctx, novelID, s.getString(characterData, "name"), s.getString(inputData, "character_type"), coreAttributes, soulProfile)
	if err != nil {
		return models.Character{}, err
	}

	// 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "character", extraInfo); err != nil {
		// 记录错误但不影响主流程
		// log.Printf("Failed to update extra info: %v", err)
	}

	return character, nil
}

// GenerateChapter 生成章节
func (s *NovelGenerationService) GenerateChapter(ctx context.Context, novelID, llmModelID string, inputData map[string]interface{}) (models.Chapter, error) {
	// 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:      novelID,
		LLMModelID:   llmModelID,
		InputData:    inputData,
		TemplateType: "chapter",
		Stream:       false,
	}

	// 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return models.Chapter{}, err
	}

	if !response.Success {
		return models.Chapter{}, errors.New(response.Error)
	}

	// 解析响应数据
	chapterData := response.Data

	// 解析章节大纲
	outline := models.ChapterOutline{
		Goal:           s.getString(chapterData, "outline.goal"),
		KeyEvents:      s.getStringArray(chapterData, "outline.key_events"),
		DramaticPoints: s.getInt(chapterData, "outline.dramatic_points"),
	}

	// 解析质量指标
	qualityMetrics := models.QualityMetrics{
		Score:            s.getInt(chapterData, "quality_metrics.score"),
		Strengths:        s.getStringArray(chapterData, "quality_metrics.strengths"),
		ImprovementAreas: s.getStringArray(chapterData, "quality_metrics.improvement_areas"),
	}

	// 解析角色发展
	characterDevelopment := s.parseCharacterDevelopment(chapterData)

	// 获取章节内容
	content := s.getString(chapterData, "content")

	// 保存到数据库
	novelService := NewNovelService(s.client, s.dbName)
	chapter, err := novelService.PostChapters(ctx, novelID, s.getInt(inputData, "chapter_number"), s.getString(chapterData, "title"), content, s.getString(chapterData, "summary"), outline, qualityMetrics, characterDevelopment)
	if err != nil {
		return models.Chapter{}, err
	}

	// 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "chapter", extraInfo); err != nil {
		// 记录错误但不影响主流程
		// log.Printf("Failed to update extra info: %v", err)
	}

	return chapter, nil
}

// 辅助方法
func (s *NovelGenerationService) getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (s *NovelGenerationService) getStringArray(data map[string]interface{}, key string) []string {
	if val, ok := data[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, len(arr))
			for i, v := range arr {
				if str, ok := v.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return []string{}
}

func (s *NovelGenerationService) getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

func (s *NovelGenerationService) parseFactions(data map[string]interface{}) []models.Faction {
	factions := []models.Faction{}
	if factionsData, ok := data["society_structure.major_factions"]; ok {
		if arr, ok := factionsData.([]interface{}); ok {
			for _, v := range arr {
				if factionMap, ok := v.(map[string]interface{}); ok {
					faction := models.Faction{
						Name:      s.getString(factionMap, "name"),
						Type:      s.getString(factionMap, "type"),
						Influence: s.getString(factionMap, "influence"),
					}
					factions = append(factions, faction)
				}
			}
		}
	}
	return factions
}

func (s *NovelGenerationService) parseRelationships(data map[string]interface{}) map[string][]string {
	relationships := make(map[string][]string)
	if relData, ok := data["core_attributes.relationships"]; ok {
		if relMap, ok := relData.(map[string]interface{}); ok {
			for key, val := range relMap {
				if arr, ok := val.([]interface{}); ok {
					result := make([]string, len(arr))
					for i, v := range arr {
						if str, ok := v.(string); ok {
							result[i] = str
						}
					}
					relationships[key] = result
				}
			}
		}
	}
	return relationships
}

func (s *NovelGenerationService) parseCharacterDevelopment(data map[string]interface{}) map[string]string {
	development := make(map[string]string)
	if devData, ok := data["character_development"]; ok {
		if devMap, ok := devData.(map[string]interface{}); ok {
			for key, val := range devMap {
				if str, ok := val.(string); ok {
					development[key] = str
				}
			}
		}
	}
	return development
}

// GenerateOutline 生成大纲
func (s *NovelGenerationService) GenerateOutline(ctx context.Context, novelID, llmModelID string, inputData map[string]interface{}) (models.Outline, error) {
	// 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:      novelID,
		LLMModelID:   llmModelID,
		InputData:    inputData,
		TemplateType: "outline",
		Stream:       false,
	}

	// 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return models.Outline{}, err
	}

	if !response.Success {
		return models.Outline{}, errors.New(response.Error)
	}

	// 解析响应数据
	outlineData := response.Data

	// 解析章节信息
	chapters := s.parseChapters(outlineData)
	
	// 解析故事弧线
	storyArcs := s.parseStoryArcs(outlineData)
	
	// 解析关键主题
	keyThemes := s.getStringArray(outlineData, "key_themes")

	// 构建大纲对象
	outline := models.Outline{
		NovelID:   novelID,
		Title:    s.getString(outlineData, "title"),
		Summary:  s.getString(outlineData, "summary"),
		Chapters: chapters,
		StoryArcs: storyArcs,
		KeyThemes: keyThemes,
	}

	// 保存到数据库
	novelService := NewNovelService(s.client, s.dbName)
	outline, err = novelService.PostOutlines(ctx, outline)
	if err != nil {
		return models.Outline{}, err
	}

	// 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "outline", extraInfo); err != nil {
		// 记录错误但不影响主流程
		// log.Printf("Failed to update extra info: %v", err)
	}

	return outline, nil
}

// parseChapters 解析章节信息
func (s *NovelGenerationService) parseChapters(data map[string]interface{}) []models.ChapterInfo {
	chapters := []models.ChapterInfo{}
	if chaptersData, ok := data["chapters"]; ok {
		if arr, ok := chaptersData.([]interface{}); ok {
			for _, v := range arr {
				if chapterMap, ok := v.(map[string]interface{}); ok {
					chapter := models.ChapterInfo{
						ChapterNumber: s.getInt(chapterMap, "chapter_number"),
						Title:         s.getString(chapterMap, "title"),
						Summary:       s.getString(chapterMap, "summary"),
						KeyEvents:     s.getStringArray(chapterMap, "key_events"),
						Characters:    s.getStringArray(chapterMap, "characters"),
						Location:      s.getString(chapterMap, "location"),
						POV:           s.getString(chapterMap, "pov"),
						WordCount:     s.getInt(chapterMap, "word_count"),
						Outline: models.ChapterOutline{
							Goal:           s.getString(chapterMap, "outline.goal"),
							KeyEvents:      s.getStringArray(chapterMap, "outline.key_events"),
							DramaticPoints: s.getInt(chapterMap, "outline.dramatic_points"),
						},
					}
					chapters = append(chapters, chapter)
				}
			}
		}
	}
	return chapters
}

// parseStoryArcs 解析故事弧线
func (s *NovelGenerationService) parseStoryArcs(data map[string]interface{}) []models.StoryArc {
	arcs := []models.StoryArc{}
	if arcsData, ok := data["story_arcs"]; ok {
		if arr, ok := arcsData.([]interface{}); ok {
			for _, v := range arr {
				if arcMap, ok := v.(map[string]interface{}); ok {
					arc := models.StoryArc{
						Name:         s.getString(arcMap, "name"),
						Description:  s.getString(arcMap, "description"),
						StartChapter: s.getInt(arcMap, "start_chapter"),
						EndChapter:   s.getInt(arcMap, "end_chapter"),
						Theme:        s.getString(arcMap, "theme"),
					}
					arcs = append(arcs, arc)
				}
			}
		}
	}
	return arcs
}
