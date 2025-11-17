import { ref, computed, watch } from 'vue'
import { useNovelStore } from '@/stores/novel'

/**
 * 小说数据自动获取和填充的Composable
 * 用于AI内容生成时自动从小说数据中获取相关信息
 */
export function useNovelData(novelId: string) {
  const novelStore = useNovelStore()
  
  // 响应式数据
  const novel = ref(null)
  const storyCores = ref([])
  const worldviews = ref([])
  const characters = ref([])
  const chapters = ref([])
  const outlines = ref([])
  
  // 加载状态
  const loading = ref(false)
  const storyCoresLoading = ref(false)
  const worldviewsLoading = ref(false)
  const charactersLoading = ref(false)
  const chaptersLoading = ref(false)
  const outlinesLoading = ref(false)
  
  // 计算属性 - 自动提取的小说信息
  const novelInfo = computed(() => {
    if (!novel.value) return null
    
    return {
      // 基本信息
      title: novel.value.title,
      genre: novel.value.project_blueprint?.genre,
      sub_genre: novel.value.project_blueprint?.sub_genre,
      target_audience: novel.value.project_blueprint?.target_audience,
      core_conflict: novel.value.project_blueprint?.core_conflict,
      
      // AI上下文
      recent_summary: novel.value.ai_context?.recent_summary,
      current_focus: novel.value.ai_context?.current_focus,
      style_guideline: novel.value.ai_context?.style_guideline,
      emotional_tone: novel.value.ai_context?.emotional_tone,
      
      // 状态信息
      status: novel.value.status,
      current_phase: novel.value.current_phase,
      
      // 时间信息
      ctime: novel.value.ctime,
      mtime: novel.value.mtime
    }
  })
  
  // 计算属性 - 故事核心信息
  const storyCoreInfo = computed(() => {
    if (!storyCores.value.length) return null
    
    return {
      available_cores: storyCores.value.map(core => ({
        id: core.id,
        title: core.title,
        core_conflict: core.core_conflict,
        theme: core.theme,
        innovation: core.innovation
      })),
      latest_core: storyCores.value[0] // 假设最新的在第一位
    }
  })
  
  // 计算属性 - 世界观信息
  const worldviewInfo = computed(() => {
    if (!worldviews.value.length) return null
    
    return {
      available_worldviews: worldviews.value.map(wv => ({
        id: wv.id,
        power_system: wv.power_system,
        society_structure: wv.society_structure,
        geography: wv.geography
      })),
      latest_worldview: worldviews.value[0]
    }
  })
  
  // 计算属性 - 角色信息
  const characterInfo = computed(() => {
    if (!characters.value.length) return null
    
    return {
      available_characters: characters.value.map(char => ({
        id: char.id,
        name: char.name,
        role: char.role,
        personality: char.soul_profile?.personality,
        background: char.soul_profile?.background
      })),
      main_characters: characters.value.filter(char => char.role === 'protagonist'),
      supporting_characters: characters.value.filter(char => char.role === 'supporting')
    }
  })
  
  // 计算属性 - 章节信息
  const chapterInfo = computed(() => {
    if (!chapters.value.length) return null
    
    return {
      available_chapters: chapters.value.map(chapter => ({
        id: chapter.id,
        title: chapter.title,
        summary: chapter.summary,
        status: chapter.status
      })),
      latest_chapter: chapters.value[0],
      completed_chapters: chapters.value.filter(chapter => chapter.status === 'completed')
    }
  })
  
  // 计算属性 - 大纲信息
  const outlineInfo = computed(() => {
    if (!outlines.value.length) return null
    
    return {
      available_outlines: outlines.value.map(outline => ({
        id: outline.id,
        title: outline.title,
        summary: outline.summary,
        chapter_count: outline.chapters?.length || 0,
        arc_count: outline.story_arcs?.length || 0
      })),
      latest_outline: outlines.value[0]
    }
  })
  
  // 计算属性 - ExtraInfo信息
  const extraInfo = computed(() => {
    if (!novel.value?.extra_info) return null
    
    return {
      story_core: novel.value.extra_info.story_core,
      worldview: novel.value.extra_info.worldview,
      character: novel.value.extra_info.character,
      outline: novel.value.extra_info.outline,
      chapter: novel.value.extra_info.chapter,
      all_phases: novel.value.extra_info
    }
  })
  
  // 获取小说基本信息
  const fetchNovel = async () => {
    try {
      loading.value = true
      console.log('开始获取小说信息...', novelId)
      await novelStore.fetchNovel(novelId)
      novel.value = novelStore.currentNovel
      console.log('小说信息获取完成:', novel.value)
    } catch (error) {
      console.error('获取小说信息失败:', error)
    } finally {
      loading.value = false
    }
  }
  
  // 获取故事核心
  const fetchStoryCores = async () => {
    try {
      storyCoresLoading.value = true
      await novelStore.fetchStoryCores(novelId)
      storyCores.value = novelStore.storyCores
    } catch (error) {
      console.error('获取故事核心失败:', error)
    } finally {
      storyCoresLoading.value = false
    }
  }
  
  // 获取世界观
  const fetchWorldviews = async () => {
    try {
      worldviewsLoading.value = true
      await novelStore.fetchWorldview(novelId)
      worldviews.value = novelStore.worldviews
    } catch (error) {
      console.error('获取世界观失败:', error)
    } finally {
      worldviewsLoading.value = false
    }
  }
  
  // 获取角色
  const fetchCharacters = async () => {
    try {
      charactersLoading.value = true
      await novelStore.fetchCharacters(novelId)
      characters.value = novelStore.characters
    } catch (error) {
      console.error('获取角色失败:', error)
    } finally {
      charactersLoading.value = false
    }
  }
  
  // 获取章节
  const fetchChapters = async () => {
    try {
      chaptersLoading.value = true
      await novelStore.fetchChapters(novelId)
      chapters.value = novelStore.chapters
    } catch (error) {
      console.error('获取章节失败:', error)
    } finally {
      chaptersLoading.value = false
    }
  }
  
  // 获取大纲
  const fetchOutlines = async () => {
    try {
      outlinesLoading.value = true
      await novelStore.fetchOutlines(novelId)
      outlines.value = novelStore.outlines
    } catch (error) {
      console.error('获取大纲失败:', error)
    } finally {
      outlinesLoading.value = false
    }
  }
  
  // 获取所有相关数据
  const fetchAllData = async () => {
    await Promise.all([
      fetchNovel(),
      fetchStoryCores(),
      fetchWorldviews(),
      fetchCharacters(),
      fetchOutlines(),
      fetchChapters()
    ])
  }
  
  // 根据当前阶段自动获取相关数据
  const fetchRelevantData = async () => {
    if (!novel.value) {
      await fetchNovel()
    }
    
    const phase = novel.value?.current_phase
    const promises = [fetchNovel()]
    
    switch (phase) {
      case 'story_core':
        promises.push(fetchStoryCores())
        break
      case 'worldview':
        promises.push(fetchStoryCores(), fetchWorldviews())
        break
      case 'characters':
        promises.push(fetchStoryCores(), fetchWorldviews(), fetchCharacters())
        break
      case 'outlining':
        promises.push(fetchStoryCores(), fetchWorldviews(), fetchCharacters(), fetchOutlines())
        break
      case 'writing':
        promises.push(fetchStoryCores(), fetchWorldviews(), fetchCharacters(), fetchOutlines(), fetchChapters())
        break
      default:
        promises.push(fetchStoryCores(), fetchWorldviews(), fetchCharacters(), fetchOutlines(), fetchChapters())
    }
    
    await Promise.all(promises)
  }
  
  // 自动填充表单数据
  const getFormDefaults = (type: string) => {
    const defaults: Record<string, any> = {}
    
    if (!novelInfo.value) return defaults
    
    switch (type) {
      case 'story-core':
        defaults.genre = novelInfo.value.genre || ''
        defaults.sub_genre = novelInfo.value.sub_genre || ''
        defaults.target_audience = novelInfo.value.target_audience || ''
        break
        
      case 'worldview':
        defaults.genre = novelInfo.value.genre || ''
        if (storyCoreInfo.value?.available_cores.length) {
          defaults.selected_concept = storyCoreInfo.value.available_cores[0].title
        }
        break
        
      case 'character':
        if (storyCoreInfo.value?.latest_core) {
          defaults.story_core = storyCoreInfo.value.latest_core.core_conflict
        }
        if (worldviewInfo.value?.latest_worldview) {
          defaults.worldview = '修真界世界观' // 这里可以根据实际世界观数据设置
        }
        break
        
      case 'outline':
        defaults.genre = novelInfo.value.genre || ''
        defaults.target_audience = novelInfo.value.target_audience || ''
        defaults.total_chapters = 50 // 默认50章
        break
        
      case 'chapter':
        if (storyCoreInfo.value?.latest_core) {
          defaults.novel_context = {
            story_core: storyCoreInfo.value.latest_core.core_conflict,
            worldview: '修真界世界观',
            current_arc: '当前故事线'
          }
        }
        if (characterInfo.value?.main_characters.length) {
          defaults.characters_involved = characterInfo.value.main_characters.map(char => char.name)
        }
        break
    }
    
    return defaults
  }
  
  // 监听小说变化，自动更新相关数据
  watch(() => novel.value?.current_phase, (newPhase, oldPhase) => {
    if (newPhase !== oldPhase) {
      fetchRelevantData()
    }
  })
  
  return {
    // 数据
    novel,
    storyCores,
    worldviews,
    characters,
    chapters,
    outlines,
    
    // 加载状态
    loading,
    storyCoresLoading,
    worldviewsLoading,
    charactersLoading,
    chaptersLoading,
    outlinesLoading,
    
    // 计算属性
    novelInfo,
    storyCoreInfo,
    worldviewInfo,
    characterInfo,
    chapterInfo,
    outlineInfo,
    extraInfo,
    
    // 方法
    fetchNovel,
    fetchStoryCores,
    fetchWorldviews,
    fetchCharacters,
    fetchChapters,
    fetchOutlines,
    fetchAllData,
    fetchRelevantData,
    getFormDefaults
  }
}
