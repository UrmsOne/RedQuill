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
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"redquill-backend/pkg/models"
	"strings"
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

// GenerateCharactersFromOutline 根据大纲批量生成角色
func (s *NovelGenerationService) GenerateCharactersFromOutline(ctx context.Context, novelID, llmModelID, outlineID string, userRequirements string) ([]models.Character, error) {
	// 1. 获取大纲数据
	novelService := NewNovelService(s.client, s.dbName)
	outline, err := novelService.GetOutline(ctx, outlineID)
	if err != nil {
		return nil, err
	}

	// 2. 获取故事核心和世界观
	storyCores, err := novelService.GetStoryCores(ctx, novelID)
	if err != nil || len(storyCores) == 0 {
		return nil, errors.New("no story core found")
	}

	worldview, err := novelService.GetWorldviews(ctx, novelID)
	if err != nil {
		return nil, errors.New("no worldview found")
	}

	// 3. 构建输入数据
	generationReq := models.GenerationRequest{
		NovelID:    novelID,
		LLMModelID: llmModelID,
		InputData: map[string]interface{}{
			"outline_content":   s.buildOutlineContent(outline),
			"story_core":        s.buildStoryCoreContent(storyCores[0]),
			"worldview":         s.buildWorldviewContent(worldview),
			"user_requirements": userRequirements,
		},
		TemplateType: "batch_character",
		Stream:       false,
	}

	// 4. 调用LLM生成
	templateService := NewPromptTemplateService(s.client, s.dbName)
	response, err := templateService.GenerateWithLLM(ctx, generationReq)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Error)
	}

	// 5. 解析响应数据
	charactersData, ok := response.Data["characters"].([]interface{})
	if !ok {
		return nil, errors.New("invalid characters data format")
	}

	// 6. 解析并保存角色
	var characters []models.Character
	for _, charData := range charactersData {
		charMap, ok := charData.(map[string]interface{})
		if !ok {
			continue
		}

		character := s.parseCharacterFromMap(charMap)
		character.NovelID = novelID

		// 保存到数据库
		savedChar, err := novelService.PostCharacters(
			ctx,
			novelID,
			character.Name,
			character.Type,
			character.CoreAttributes,
			character.SoulProfile,
		)
		if err != nil {
			// 记录错误但继续处理其他角色
			continue
		}
		characters = append(characters, savedChar)
	}

	// 7. 保存ExtraInfo
	extraInfo := map[string]interface{}{
		"generation_time": response.Data["generation_time"],
		"token_count":     response.TokenCount,
		"usage_count":     response.UsageCount,
		"raw_response":    response.Data,
		"outline_id":      outlineID,
		"character_count": len(characters),
	}
	if err := novelService.UpdateNovelExtraInfo(ctx, novelID, "batch_character", extraInfo); err != nil {
		// 记录错误但不影响主流程
	}

	return characters, nil
}

// buildOutlineContent 构建大纲内容文本
func (s *NovelGenerationService) buildOutlineContent(outline models.Outline) string {
	content := fmt.Sprintf("标题：%s\n概要：%s\n", outline.Title, outline.Summary)

	if len(outline.KeyThemes) > 0 {
		content += fmt.Sprintf("关键主题：%s\n", strings.Join(outline.KeyThemes, "、"))
	}

	if len(outline.StoryArcs) > 0 {
		content += "故事弧线：\n"
		for _, arc := range outline.StoryArcs {
			content += fmt.Sprintf("- %s（第%d-%d章）：%s\n", arc.Name, arc.StartChapter, arc.EndChapter, arc.Description)
		}
	}

	if len(outline.Chapters) > 0 {
		content += "章节信息：\n"
		for _, chapter := range outline.Chapters {
			content += fmt.Sprintf("第%d章 %s：%s\n", chapter.ChapterNumber, chapter.Title, chapter.Summary)
			if len(chapter.Characters) > 0 {
				content += fmt.Sprintf("  涉及角色：%s\n", strings.Join(chapter.Characters, "、"))
			}
		}
	}

	return content
}

// buildStoryCoreContent 构建故事核心内容文本
func (s *NovelGenerationService) buildStoryCoreContent(storyCore models.StoryCore) string {
	return fmt.Sprintf(`标题：%s
核心冲突：%s
主题：%s
创新点：%s
商业潜力：%s
目标受众：%s`,
		storyCore.Title,
		storyCore.CoreConflict,
		storyCore.Theme,
		storyCore.Innovation,
		storyCore.CommercialPotential,
		storyCore.TargetAudience,
	)
}

// buildWorldviewContent 构建世界观内容文本
func (s *NovelGenerationService) buildWorldviewContent(worldview models.Worldview) string {
	content := ""

	if worldview.PowerSystem.Name != "" {
		content += fmt.Sprintf("力量体系：%s\n", worldview.PowerSystem.Name)
		content += fmt.Sprintf("修炼方式：%s\n", worldview.PowerSystem.CultivationMethod)
	}

	if worldview.SocietyStructure.Hierarchy != "" {
		content += fmt.Sprintf("社会结构：%s\n", worldview.SocietyStructure.Hierarchy)
		content += fmt.Sprintf("经济体系：%s\n", worldview.SocietyStructure.EconomicSystem)
	}

	if len(worldview.Geography.MajorRegions) > 0 {
		content += fmt.Sprintf("主要地域：%s\n", strings.Join(worldview.Geography.MajorRegions, "、"))
	}

	if len(worldview.SpecialRules) > 0 {
		content += fmt.Sprintf("特殊规则：%s\n", strings.Join(worldview.SpecialRules, "、"))
	}

	return content
}

// parseCharacterFromMap 从map解析角色数据
func (s *NovelGenerationService) parseCharacterFromMap(charMap map[string]interface{}) models.Character {
	// 解析灵魂档案
	soulProfileData := charMap["soul_profile"].(map[string]interface{})
	personalityData := soulProfileData["personality"].(map[string]interface{})
	backgroundData := soulProfileData["background"].(map[string]interface{})
	motivationsData := soulProfileData["motivations"].(map[string]interface{})

	soulProfile := models.SoulProfile{
		Personality: models.Personality{
			CoreTraits:        s.getStringArray(personalityData, "core_traits"),
			MoralCompass:      s.getString(personalityData, "moral_compass"),
			InternalConflicts: s.getStringArray(personalityData, "internal_conflicts"),
			Fears:             s.getStringArray(personalityData, "fears"),
			Desires:           s.getStringArray(personalityData, "desires"),
		},
		Background: models.Background{
			Origin:         s.getString(backgroundData, "origin"),
			DefiningEvents: s.getStringArray(backgroundData, "defining_events"),
			HiddenSecrets:  s.getStringArray(backgroundData, "hidden_secrets"),
		},
		Motivations: models.Motivations{
			ImmediateGoal: s.getString(motivationsData, "immediate_goal"),
			LongTermGoal:  s.getString(motivationsData, "long_term_goal"),
			CoreDrive:     s.getString(motivationsData, "core_drive"),
		},
	}

	// 解析核心属性
	coreAttributesData := charMap["core_attributes"].(map[string]interface{})
	relationshipsData := coreAttributesData["relationships"].(map[string]interface{})

	coreAttributes := models.CoreAttributes{
		CultivationLevel: s.getString(coreAttributesData, "cultivation_level"),
		CurrentItems:     s.getStringArray(coreAttributesData, "current_items"),
		Abilities:        s.getStringArray(coreAttributesData, "abilities"),
		Relationships: map[string][]string{
			"enemies": s.getStringArray(relationshipsData, "enemies"),
			"allies":  s.getStringArray(relationshipsData, "allies"),
			"mentors": s.getStringArray(relationshipsData, "mentors"),
		},
	}

	return models.Character{
		Name:           s.getString(charMap, "name"),
		Type:           s.getString(charMap, "type"),
		SoulProfile:    soulProfile,
		CoreAttributes: coreAttributes,
	}
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
		Title:     s.getString(outlineData, "title"),
		Summary:   s.getString(outlineData, "summary"),
		Chapters:  chapters,
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
