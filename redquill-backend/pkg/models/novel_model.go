// Package models
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: novel_model.go
/@Description: Novel model data structure
/*/

package models

// Novel 作品
type Novel struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Title        string `json:"title" bson:"title"`
	AuthorID     string `json:"author_id" bson:"author_id"`
	Status       string `json:"status" bson:"status"`               // drafting|writing|completed|paused
	CurrentPhase string `json:"current_phase" bson:"current_phase"` // story_core|worldview|characters|outlining|writing
	Ctime        int64  `json:"ctime" bson:"ctime"`
	Mtime        int64  `json:"mtime" bson:"mtime"`

	ProjectBlueprint ProjectBlueprint     `json:"project_blueprint" bson:"project_blueprint"`
	AIContext        AIContext            `json:"ai_context" bson:"ai_context"`
	ExtraInfo        map[string]interface{} `json:"extra_info" bson:"extra_info"` // 存放各个阶段AI生成返回的信息
}

// ProjectBlueprint 项目蓝图
type ProjectBlueprint struct {
	Genre           string `json:"genre" bson:"genre"`
	SubGenre        string `json:"sub_genre" bson:"sub_genre"`
	TotalChapters   int    `json:"total_chapters" bson:"total_chapters"`
	CoreConflict    string `json:"core_conflict" bson:"core_conflict"`
	TargetAudience  string `json:"target_audience" bson:"target_audience"`
	CommercialFocus string `json:"commercial_focus" bson:"commercial_focus"`
}

// AIContext AI上下文
type AIContext struct {
	RecentSummary  string `json:"recent_summary" bson:"recent_summary"`
	CurrentFocus   string `json:"current_focus" bson:"current_focus"`
	StyleGuideline string `json:"style_guideline" bson:"style_guideline"`
	EmotionalTone  string `json:"emotional_tone" bson:"emotional_tone"`
}

// StoryCore 故事核心
type StoryCore struct {
	ID                  string `json:"id" bson:"_id,omitempty"`
	NovelID             string `json:"novel_id" bson:"novel_id"`
	Title               string `json:"title" bson:"title"`
	CoreConflict        string `json:"core_conflict" bson:"core_conflict"`
	Theme               string `json:"theme" bson:"theme"`
	Innovation          string `json:"innovation" bson:"innovation"`
	CommercialPotential string `json:"commercial_potential" bson:"commercial_potential"`
	TargetAudience      string `json:"target_audience" bson:"target_audience"`
	Ctime               int64  `json:"ctime" bson:"ctime"`
}

// Worldview 世界观
type Worldview struct {
	ID               string           `json:"id" bson:"_id,omitempty"`
	NovelID          string           `json:"novel_id" bson:"novel_id"`
	PowerSystem      PowerSystem      `json:"power_system" bson:"power_system"`
	SocietyStructure SocietyStructure `json:"society_structure" bson:"society_structure"`
	Geography        Geography        `json:"geography" bson:"geography"`
	SpecialRules     []string         `json:"special_rules" bson:"special_rules"`
	Ctime            int64            `json:"ctime" bson:"ctime"`
}

// PowerSystem 力量体系
type PowerSystem struct {
	Name              string   `json:"name" bson:"name"`
	Levels            []string `json:"levels" bson:"levels"`
	CultivationMethod string   `json:"cultivation_method" bson:"cultivation_method"`
	Limitations       string   `json:"limitations" bson:"limitations"`
}

// SocietyStructure 社会结构
type SocietyStructure struct {
	Hierarchy      string    `json:"hierarchy" bson:"hierarchy"`
	MajorFactions  []Faction `json:"major_factions" bson:"major_factions"`
	EconomicSystem string    `json:"economic_system" bson:"economic_system"`
}

// Faction 势力
type Faction struct {
	Name      string `json:"name" bson:"name"`
	Type      string `json:"type" bson:"type"`
	Influence string `json:"influence" bson:"influence"`
}

// Geography 地理
type Geography struct {
	MajorRegions     []string `json:"major_regions" bson:"major_regions"`
	SpecialLocations []string `json:"special_locations" bson:"special_locations"`
}

// Character 角色
type Character struct {
	ID             string         `json:"id" bson:"_id,omitempty"`
	NovelID        string         `json:"novel_id" bson:"novel_id"`
	Name           string         `json:"name" bson:"name"`
	Type           string         `json:"type" bson:"type"` // protagonist|antagonist|supporting|love_interest
	CoreAttributes CoreAttributes `json:"core_attributes" bson:"core_attributes"`
	SoulProfile    SoulProfile    `json:"soul_profile" bson:"soul_profile"`
	GrowthTrack    []GrowthEvent  `json:"growth_track" bson:"growth_track"`
	Ctime          int64          `json:"updated_at" bson:"updated_at"`
}

// CoreAttributes 核心属性
type CoreAttributes struct {
	CultivationLevel string              `json:"cultivation_level" bson:"cultivation_level"`
	CurrentItems     []string            `json:"current_items" bson:"current_items"`
	Abilities        []string            `json:"abilities" bson:"abilities"`
	Relationships    map[string][]string `json:"relationships" bson:"relationships"`
}

// SoulProfile 灵魂档案
type SoulProfile struct {
	Personality Personality `json:"personality" bson:"personality"`
	Background  Background  `json:"background" bson:"background"`
	Motivations Motivations `json:"motivations" bson:"motivations"`
}

// Personality 性格
type Personality struct {
	CoreTraits        []string `json:"core_traits" bson:"core_traits"`
	MoralCompass      string   `json:"moral_compass" bson:"moral_compass"`
	InternalConflicts []string `json:"internal_conflicts" bson:"internal_conflicts"`
	Fears             []string `json:"fears" bson:"fears"`
	Desires           []string `json:"desires" bson:"desires"`
}

// Background 背景
type Background struct {
	Origin         string   `json:"origin" bson:"origin"`
	DefiningEvents []string `json:"defining_events" bson:"defining_events"`
	HiddenSecrets  []string `json:"hidden_secrets" bson:"hidden_secrets"`
}

// Motivations 动机
type Motivations struct {
	ImmediateGoal string `json:"immediate_goal" bson:"immediate_goal"`
	LongTermGoal  string `json:"long_term_goal" bson:"long_term_goal"`
	CoreDrive     string `json:"core_drive" bson:"core_drive"`
}

// GrowthEvent 成长事件
type GrowthEvent struct {
	Chapter int    `json:"chapter" bson:"chapter"`
	Event   string `json:"event" bson:"event"`
	Change  string `json:"change" bson:"change"`
}

// Chapter 章节
type Chapter struct {
	ID                   string            `json:"id" bson:"_id,omitempty"`
	NovelID              string            `json:"novel_id" bson:"novel_id"`
	ChapterNumber        int               `json:"chapter_number" bson:"chapter_number"`
	Title                string            `json:"title" bson:"title"`
	Content              string            `json:"content" bson:"content"`
	WordCount            int               `json:"word_count" bson:"word_count"`
	Summary              string            `json:"summary" bson:"summary"`
	Outline              ChapterOutline    `json:"outline" bson:"outline"`
	QualityMetrics       QualityMetrics    `json:"quality_metrics" bson:"quality_metrics"`
	CharacterDevelopment map[string]string `json:"character_development" bson:"character_development"`
	Ctime                int64             `json:"ctime" bson:"ctime"`
}

// ChapterOutline 章节大纲
type ChapterOutline struct {
	Goal           string   `json:"goal" bson:"goal"`
	KeyEvents      []string `json:"key_events" bson:"key_events"`
	DramaticPoints int      `json:"dramatic_points" bson:"dramatic_points"`
}

// QualityMetrics 质量指标
type QualityMetrics struct {
	Score            int      `json:"score" bson:"score"`
	Strengths        []string `json:"strengths" bson:"strengths"`
	ImprovementAreas []string `json:"improvement_areas" bson:"improvement_areas"`
}

// WritingSession 创作会话
type WritingSession struct {
	ID             string         `json:"id" bson:"_id,omitempty"`
	NovelID        string         `json:"novel_id" bson:"novel_id"`
	CurrentChapter int            `json:"current_chapter" bson:"current_chapter"`
	SessionContext SessionContext `json:"session_context" bson:"session_context"`
	LastUpdated    int64          `json:"last_updated" bson:"last_updated"`
}

// SessionContext 会话上下文
type SessionContext struct {
	RecentEvents     string   `json:"recent_events" bson:"recent_events"`
	CurrentArc       string   `json:"current_arc" bson:"current_arc"`
	PendingConflicts []string `json:"pending_conflicts" bson:"pending_conflicts"`
}

// Outline 大纲
type Outline struct {
	ID           string        `json:"id" bson:"_id,omitempty"`
	NovelID      string        `json:"novel_id" bson:"novel_id"`
	Title        string        `json:"title" bson:"title"`
	Summary      string        `json:"summary" bson:"summary"`
	Chapters     []ChapterInfo `json:"chapters" bson:"chapters"`
	StoryArcs    []StoryArc    `json:"story_arcs" bson:"story_arcs"`
	KeyThemes    []string      `json:"key_themes" bson:"key_themes"`
	Ctime        int64         `json:"ctime" bson:"ctime"`
	Mtime        int64         `json:"mtime" bson:"mtime"`
}

// ChapterInfo 章节信息
type ChapterInfo struct {
	ChapterNumber int               `json:"chapter_number" bson:"chapter_number"`
	Title         string            `json:"title" bson:"title"`
	Summary       string            `json:"summary" bson:"summary"`
	KeyEvents     []string          `json:"key_events" bson:"key_events"`
	Characters    []string          `json:"characters" bson:"characters"`
	Location      string            `json:"location" bson:"location"`
	POV           string            `json:"pov" bson:"pov"` // 视角角色
	WordCount     int               `json:"word_count" bson:"word_count"`
	Outline       ChapterOutline    `json:"outline" bson:"outline"`
}

// StoryArc 故事弧线
type StoryArc struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	StartChapter int   `json:"start_chapter" bson:"start_chapter"`
	EndChapter   int   `json:"end_chapter" bson:"end_chapter"`
	Theme       string `json:"theme" bson:"theme"`
}
