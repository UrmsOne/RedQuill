import { defineStore } from 'pinia'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { api } from '@/utils/api'

export interface Novel {
  id: string
  title: string
  author_id: string
  status: string
  current_phase: string
  ctime: number
  mtime: number
  project_blueprint: ProjectBlueprint
  ai_context: AIContext
  extra_info?: Record<string, any> // 存放各个阶段AI生成返回的信息
}

export interface ProjectBlueprint {
  genre: string
  sub_genre: string
  total_chapters: number
  core_conflict: string
  target_audience: string
  commercial_focus: string
}

export interface AIContext {
  recent_summary: string
  current_focus: string
  style_guideline: string
  emotional_tone: string
}

export interface StoryCore {
  id: string
  novel_id: string
  title: string
  core_conflict: string
  theme: string
  innovation: string
  commercial_potential: string
  target_audience: string
}

export interface Worldview {
  id: string
  novel_id: string
  power_system: PowerSystem
  society_structure: SocietyStructure
  geography: Geography
  special_rules: string[]
}

export interface PowerSystem {
  name: string
  levels: string[]
  cultivation_method: string
  limitations: string
}

export interface SocietyStructure {
  hierarchy: string
  major_factions: Faction[]
  economic_system: string
}

export interface Faction {
  name: string
  type: string
  influence: string
}

export interface Geography {
  major_regions: string[]
  special_locations: string[]
}

export interface Character {
  id: string
  novel_id: string
  name: string
  type: string
  core_attributes: CoreAttributes
  soul_profile: SoulProfile
}

export interface CoreAttributes {
  cultivation_level: string
  current_items: string[]
  abilities: string[]
  relationships: Relationships
}

export interface Relationships {
  enemies: string[]
  allies: string[]
  mentors: string[]
}

export interface SoulProfile {
  personality: Personality
  background: Background
  motivations: Motivations
}

export interface Personality {
  core_traits: string[]
  moral_compass: string
  internal_conflicts: string[]
  fears: string[]
  desires: string[]
}

export interface Background {
  origin: string
  defining_events: string[]
  hidden_secrets: string[]
}

export interface Motivations {
  immediate_goal: string
  long_term_goal: string
  core_drive: string
}

export interface Chapter {
  id: string
  novel_id: string
  chapter_number: number
  title: string
  content: string
  summary: string
  word_count: number
  outline: ChapterOutline
  quality_metrics: QualityMetrics
  character_development: Record<string, string>
}

export interface ChapterOutline {
  goal: string
  key_events: string[]
  dramatic_points: number
}

export interface QualityMetrics {
  score: number
  strengths: string[]
  improvement_areas: string[]
}

export interface Outline {
  id: string
  novel_id: string
  title: string
  summary: string
  chapters: ChapterInfo[]
  story_arcs: StoryArc[]
  key_themes: string[]
  ctime: number
  mtime: number
}

export interface ChapterInfo {
  chapter_number: number
  title: string
  summary: string
  key_events: string[]
  characters: string[]
  location: string
  pov: string
  word_count: number
  outline: ChapterOutline
}

export interface StoryArc {
  name: string
  description: string
  start_chapter: number
  end_chapter: number
  theme: string
}

export const useNovelStore = defineStore('novel', () => {
  const novels = ref<Novel[]>([])
  const currentNovel = ref<Novel | null>(null)
  const storyCores = ref<StoryCore[]>([])
  const worldviews = ref<Worldview[]>([])
  const characters = ref<Character[]>([])
  const chapters = ref<Chapter[]>([])
  const outlines = ref<Outline[]>([])
  const loading = ref(false)

  // 获取小说列表
  const fetchNovels = async (params?: any) => {
    try {
      loading.value = true
      console.log('Store: 开始获取小说列表...', params)
      const response = await api.get('/novels', { params })
      console.log('Store: API响应:', response.data)
      // 确保正确处理空响应和null值
      const items = response.data?.items
      novels.value = Array.isArray(items) ? items : []
      console.log('Store: 设置小说列表:', novels.value.length, '条记录')
    } catch (error: any) {
      console.error('Store: 获取小说列表失败:', error)
      message.error(error.response?.data?.error || '获取小说列表失败')
      novels.value = []
    } finally {
      loading.value = false
      console.log('Store: loading状态已重置:', loading.value)
    }
  }

  // 获取小说详情
  const fetchNovel = async (id: string) => {
    try {
      loading.value = true
      const response = await api.get(`/novel/${id}`)
      currentNovel.value = response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取小说详情失败')
    } finally {
      loading.value = false
    }
  }

  // 创建小说
  const createNovel = async (novelData: Partial<Novel>) => {
    try {
      loading.value = true
      const response = await api.post('/novel', novelData)
      novels.value.unshift(response.data)
      message.success('创建成功')
      return response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '创建失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 更新小说
  const updateNovel = async (id: string, novelData: Partial<Novel>) => {
    try {
      loading.value = true
      const response = await api.put(`/novel/${id}`, novelData)
      const index = novels.value.findIndex(n => n.id === id)
      if (index !== -1) {
        novels.value[index] = response.data
      }
      if (currentNovel.value?.id === id) {
        currentNovel.value = response.data
      }
      message.success('更新成功')
      return response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '更新失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 删除小说
  const deleteNovel = async (id: string) => {
    try {
      loading.value = true
      await api.delete(`/novel/${id}`)
      novels.value = novels.value.filter(n => n.id !== id)
      if (currentNovel.value?.id === id) {
        currentNovel.value = null
      }
      message.success('删除成功')
      return true
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取故事核心
  const fetchStoryCores = async (novelId: string) => {
    try {
      const response = await api.get(`/story-cores/${novelId}`)
      storyCores.value = response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取故事核心失败')
    }
  }

  // 获取世界观
  const fetchWorldview = async (novelId: string) => {
    try {
      const response = await api.get(`/worldview/${novelId}`)
      worldviews.value = response.data ? [response.data] : []
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取世界观失败')
    }
  }

  // 获取角色列表
  const fetchCharacters = async (novelId: string) => {
    try {
      const response = await api.get(`/characters/${novelId}`)
      characters.value = response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取角色列表失败')
    }
  }

  // 获取章节列表
  const fetchChapters = async (novelId: string) => {
    try {
      const response = await api.get(`/chapters/${novelId}`)
      chapters.value = response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取章节列表失败')
    }
  }

  // 获取大纲列表
  const fetchOutlines = async (novelId: string) => {
    try {
      const response = await api.get(`/outlines/${novelId}`)
      outlines.value = response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取大纲列表失败')
    }
  }

  // 获取单个大纲
  const fetchOutline = async (id: string) => {
    try {
      const response = await api.get(`/outline/${id}`)
      return response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '获取大纲详情失败')
      return null
    }
  }

  // 创建大纲
  const createOutline = async (outlineData: Partial<Outline>) => {
    try {
      loading.value = true
      const response = await api.post('/outline', outlineData)
      outlines.value.unshift(response.data)
      message.success('创建成功')
      return response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '创建失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 更新大纲
  const updateOutline = async (id: string, outlineData: Partial<Outline>) => {
    try {
      loading.value = true
      const response = await api.put(`/outline/${id}`, outlineData)
      const index = outlines.value.findIndex(o => o.id === id)
      if (index !== -1) {
        outlines.value[index] = response.data
      }
      message.success('更新成功')
      return response.data
    } catch (error: any) {
      message.error(error.response?.data?.error || '更新失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 删除大纲
  const deleteOutline = async (id: string) => {
    try {
      loading.value = true
      await api.delete(`/outline/${id}`)
      outlines.value = outlines.value.filter(o => o.id !== id)
      message.success('删除成功')
      return true
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    novels,
    currentNovel,
    storyCores,
    worldviews,
    characters,
    chapters,
    outlines,
    loading,
    fetchNovels,
    fetchNovel,
    createNovel,
    updateNovel,
    deleteNovel,
    fetchStoryCores,
    fetchWorldview,
    fetchCharacters,
    fetchChapters,
    fetchOutlines,
    fetchOutline,
    createOutline,
    updateOutline,
    deleteOutline
  }
})
