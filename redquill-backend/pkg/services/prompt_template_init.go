// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_template_init.go
/@Description: Initialize prompt templates based on novels.md
/*/

package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"redquill-backend/pkg/models"
)

// InitializePromptTemplates 初始化Prompt模板
func InitializePromptTemplates(client *mongo.Client, dbName string) error {
	ctx := context.Background()
	coll := client.Database(dbName).Collection("prompt_templates")

	// 检查是否已经初始化
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 已经初始化过了
	}

	templates := []models.PromptTemplate{
		{
			Name:        "故事核心生成",
			Type:        "story_core",
			Phase:       "story_core",
			Description: "基于用户选择的题材和初步想法，生成具有爆款潜力的故事核心方案",
			Content: `【角色】
你是一位资深网文编辑，精通各类题材的爆款规律。

【任务】
基于用户选择的题材和初步想法，生成3个具有爆款潜力的故事核心方案。

【输入数据】
{
  "genre": "{genre}",
  "sub_genre": "{sub_genre}", 
  "user_ideas": "{user_ideas}",
  "target_audience": "{target_audience}"
}

【输出要求】
请严格按照以下JSON格式输出3个完整方案：
{
  "concepts": [
    {
      "title": "故事标题",
      "core_conflict": "核心矛盾（一句话概括主角面临的主要冲突）",
      "theme": "故事主题",
      "innovation": "创新亮点",
      "commercial_potential": "商业潜力分析",
      "target_audience": "目标读者群体"
    }
  ]
}`,
			Variables:  []string{"genre", "sub_genre", "user_ideas", "target_audience"},
			UsageCount: 0,
			CreatorID:  "system",
			Creator:    "system",
			Ctime:      time.Now().Unix(),
			Mtime:      time.Now().Unix(),
		},
		{
			Name:        "世界观构建",
			Type:        "worldview",
			Phase:       "worldview",
			Description: "为选定的故事核心构建完整的世界观体系",
			Content: `【角色】
你是世界架构师，擅长构建自洽且富有魅力的故事舞台。

【任务】
为选定的故事核心构建完整的世界观体系。

【输入数据】
{
  "genre": "{genre}",
  "title" : "{title}",
  "core_conflict" : "{core_conflict}",
  "theme" : "{theme}",
  "innovation" : "{innovation}",
  "commercial_potential" : "{commercial_potential}",
  "target_audience" : "{target_audience}",
  "user_ideas": "{user_ideas}",
}

【输出要求】
请严格按照以下JSON格式输出完整世界观：
{
  "power_system": {
    "name": "力量体系名称",
    "levels": ["等级1", "等级2", ...],
    "cultivation_method": "修炼方式描述",
    "limitations": "力量限制条件"
  },
  "society_structure": {
    "hierarchy": "社会阶层描述",
    "major_factions": [
      {
        "name": "势力名称",
        "type": "正派|反派|中立",
        "influence": "势力影响力描述"
      }
    ],
    "economic_system": "经济体系描述"
  },
  "geography": {
    "major_regions": ["主要地域1", "主要地域2"],
    "special_locations": ["特殊地点1", "特殊地点2"]
  },
  "special_rules": ["特殊规则1", "特殊规则2"]
}`,
			Variables:  []string{"title", "core_conflict", "theme", "innovation", "commercial_potential", "target_audience", "genre"},
			UsageCount: 0,
			CreatorID:  "system",
			Creator:    "system",
			Ctime:      time.Now().Unix(),
			Mtime:      time.Now().Unix(),
		},
		{
			Name:        "角色灵魂塑造",
			Type:        "character",
			Phase:       "characters",
			Description: "基于故事核心和世界观，深度塑造主要角色的内在灵魂",
			Content: `【角色】
你是顶尖的角色设计师，擅长创造有血有肉、让读者印象深刻的人物。

【任务】
基于故事核心和世界观，深度塑造主要角色的内在灵魂。

【输入数据】
{
  "story_core": "{story_core}",
  "worldview": "{worldview}",
  "character_type": "{character_type}",
  "role_requirements": "{role_requirements}"
}

【输出要求】
请严格按照以下JSON格式输出角色灵魂档案：
{
  "name": "角色名称",
  "soul_profile": {
    "personality": {
      "core_traits": ["特质1", "特质2"],
      "moral_compass": "道德观描述",
      "internal_conflicts": ["内心矛盾1", "内心矛盾2"],
      "fears": ["恐惧1", "恐惧2"],
      "desires": ["欲望1", "欲望2"]
    },
    "background": {
      "origin": "出身背景",
      "defining_events": ["关键事件1", "关键事件2"],
      "hidden_secrets": ["隐藏秘密1", "隐藏秘密2"]
    },
    "motivations": {
      "immediate_goal": "近期目标",
      "long_term_goal": "长期目标", 
      "core_drive": "核心驱动力"
    }
  },
  "core_attributes": {
    "cultivation_level": "初始修炼境界",
    "current_items": ["初始物品1", "初始物品2"],
    "abilities": ["初始能力1", "初始能力2"],
    "relationships": {
      "enemies": ["敌人1", "敌人2"],
      "allies": ["盟友1", "盟友2"],
      "mentors": ["导师1", "导师2"]
    }
  }
}`,
			Variables:  []string{"story_core", "worldview", "character_type", "role_requirements"},
			UsageCount: 0,
			CreatorID:  "system",
			Creator:    "system",
			Ctime:      time.Now().Unix(),
			Mtime:      time.Now().Unix(),
		},
		{
			Name:        "章节内容生成",
			Type:        "chapter",
			Phase:       "writing",
			Description: "根据章节大纲和目标，生成具体章节内容",
			Content: `【角色】
你是{novel_title}的御用写手，完全沉浸在故事的世界中。

【任务】
根据章节大纲和目标，生成具体章节内容。

【输入数据】
{
  "novel_context": {
    "story_core": "{story_core}",
    "worldview": "{worldview}",
    "current_arc": "{current_arc}"
  },
  "chapter_goal": "{chapter_goal}",
  "characters_involved": [
    {
      "soul_profile": "{character_soul_profile}",
      "core_attributes": "{character_core_attributes}",
      "current_state": "{character_current_state}"
    }
  ],
  "previous_summary": "{previous_summary}",
  "plot_templates": ["{plot_template1}", "{plot_template2}"]
}

【输出要求】
请先输出章节元数据JSON，然后输出正文内容：
{
  "title": "章节标题",
  "summary": "本章内容摘要",
  "plot_advancements": ["剧情推进点1", "剧情推进点2"],
  "character_development": {
    "角色名": "在本章中的成长变化"
  },
  "next_chapter_hook": "为下一章埋下的钩子",
  "outline": {
    "goal": "本章核心目标",
    "key_events": ["关键事件1", "关键事件2"],
    "dramatic_points": 3
  },
  "quality_metrics": {
    "score": 8,
    "strengths": ["优点1", "优点2"],
    "improvement_areas": ["待改进领域1", "待改进领域2"]
  }
}

【正文开始】
（此处生成2000字左右的章节正文内容）`,
			Variables:  []string{"novel_title", "story_core", "worldview", "current_arc", "chapter_goal", "characters_involved", "previous_summary", "plot_templates"},
			UsageCount: 0,
			CreatorID:  "system",
			Creator:    "system",
			Ctime:      time.Now().Unix(),
			Mtime:      time.Now().Unix(),
		},
		{
			Name:        "内容质量审核",
			Type:        "quality_review",
			Phase:       "writing",
			Description: "对生成的章节内容进行全面的质量评估",
			Content: `【角色】
你是苛刻的网文读者和专业的文学评论家。

【任务】
对生成的章节内容进行全面的质量评估。

【输入数据】
{
  "chapter_content": "{chapter_content}",
  "chapter_metadata": "{chapter_metadata}",
  "novel_context": {
    "story_core": "{story_core}",
    "worldview": "{worldview}"
  },
  "quality_standards": {
    "role_consistency": "角色一致性要求",
    "plot_advancement": "剧情推进要求", 
    "emotional_impact": "情感冲击要求"
  }
}

【输出要求】
请严格按照以下JSON格式输出质量报告：
{
  "overall_score": "总体评分1-10",
  "strengths": ["优点1", "优点2"],
  "issues": [
    {
      "type": "问题类型：role_inconsistency|pacing_issue|logic_error",
      "location": "问题位置描述",
      "description": "问题详细描述", 
      "suggestion": "修改建议"
    }
  ],
  "optimization_suggestions": ["优化建议1", "优化建议2"]
}`,
			Variables:  []string{"chapter_content", "chapter_metadata", "story_core", "worldview", "quality_standards"},
			UsageCount: 0,
			CreatorID:  "system",
			Creator:    "system",
			Ctime:      time.Now().Unix(),
			Mtime:      time.Now().Unix(),
		},
	}

	// 插入模板
	for _, template := range templates {
		_, err := coll.InsertOne(ctx, template)
		if err != nil {
			return err
		}
	}

	return nil
}
