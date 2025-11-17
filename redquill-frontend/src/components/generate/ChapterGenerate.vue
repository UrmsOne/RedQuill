<template>
  <div class="chapter-generate" v-if="novel">
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
      <a-form-item label="章节号" name="chapter_number">
        <a-select 
          v-model:value="form.chapter_number" 
          placeholder="请选择或输入章节号"
          :loading="chaptersLoading || outlinesLoading"
          show-search
          allow-clear
          @search="handleChapterNumberSearch"
        >
          <a-select-option 
            v-for="option in chapterNumberOptions" 
            :key="option.value" 
            :value="option.value"
          >
            {{ option.label }}
          </a-select-option>
        </a-select>
        <div v-if="selectedChapterInfo" class="chapter-info-preview" style="margin-top: 8px; padding: 8px; background: #f5f5f5; border-radius: 4px;">
          <p style="margin: 0;"><strong>章节标题：</strong>{{ selectedChapterInfo.title }}</p>
          <p style="margin: 4px 0 0 0;"><strong>章节概要：</strong>{{ selectedChapterInfo.summary }}</p>
          <p v-if="selectedChapterInfo.characters && selectedChapterInfo.characters.length > 0" style="margin: 4px 0 0 0;">
            <strong>涉及角色：</strong>{{ selectedChapterInfo.characters.join('、') }}
          </p>
        </div>
      </a-form-item>

      <a-form-item label="大纲" name="outline_id">
        <a-select 
          v-model:value="form.outline_id" 
          placeholder="请选择大纲（可选，选择后将从大纲获取章节信息）"
          :loading="outlinesLoading"
          allow-clear
          @change="handleOutlineChange"
        >
          <a-select-option 
            v-for="outline in outlines" 
            :key="outline.id" 
            :value="outline.id"
          >
            {{ outline.title }} ({{ outline.chapters?.length || 0 }}章)
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="章节目标" name="chapter_goal">
        <a-textarea 
          v-model:value="form.chapter_goal" 
          placeholder="请描述本章节的目标和要达成的效果..."
          :rows="3"
        />
      </a-form-item>
      
      <a-form-item label="前情提要" name="previous_summary">
        <a-textarea 
          v-model:value="form.previous_summary" 
          placeholder="将自动填充上一个章节的摘要和正文内容..."
          :rows="8"
          :auto-size="{ minRows: 4, maxRows: 12 }"
        />
      </a-form-item>
      
      <a-form-item label="参与角色" name="characters_involved">
        <a-select 
          v-model:value="form.characters_involved" 
          mode="multiple"
          placeholder="请选择参与的角色"
          :loading="charactersLoading"
        >
          <a-select-option v-for="character in characters" :key="character.id" :value="character.name">
            {{ character.name }} ({{ character.type }})
          </a-select-option>
        </a-select>
      </a-form-item>
      
      <a-form-item label="情节模板" name="plot_templates">
        <a-select 
          v-model:value="form.plot_templates" 
          mode="multiple"
          placeholder="请选择情节模板"
        >
          <a-select-option value="冲突升级">冲突升级</a-select-option>
          <a-select-option value="角色成长">角色成长</a-select-option>
          <a-select-option value="情感转折">情感转折</a-select-option>
          <a-select-option value="悬念设置">悬念设置</a-select-option>
          <a-select-option value="高潮爆发">高潮爆发</a-select-option>
        </a-select>
      </a-form-item>
      
      <a-form-item label="LLM模型" name="llm_model_id">
        <a-select v-model:value="form.llm_model_id" placeholder="请选择LLM模型" :loading="modelsLoading">
          <a-select-option v-for="model in models" :key="model.id" :value="model.id">
            {{ model.display_name || model.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      
      <a-form-item>
        <a-space>
          <a-button 
            type="primary" 
            html-type="submit" 
            :loading="generating"
            :disabled="!form.llm_model_id"
          >
            {{ generating ? '生成中...' : '开始生成' }}
          </a-button>
          
          <a-button 
            @click="handleStreamGenerate"
            :loading="streaming"
            :disabled="!form.llm_model_id"
          >
            {{ streaming ? '流式生成中...' : '流式生成' }}
          </a-button>
        </a-space>
      </a-form-item>
    </a-form>
    
    <!-- 生成结果 -->
    <div v-if="result" class="result-section">
      <a-card title="生成结果">
        <template #extra>
          <a-space>
            <a-button 
              type="primary" 
              @click="handleSaveToDatabase"
              :loading="saving"
            >
              {{ saving ? '保存中...' : '保存到数据库' }}
            </a-button>
            <a-button @click="handleViewRawData">
              查看原始数据
            </a-button>
          </a-space>
        </template>
        <div class="result-content">
          <h3>第{{ result.chapter_number || form.chapter_number || '?' }}章 {{ result.title || '未命名章节' }}</h3>
          <p v-if="result.word_count"><strong>字数：</strong>{{ result.word_count }}字</p>
          <p v-if="result.summary"><strong>摘要：</strong>{{ result.summary }}</p>
          
          <h4 v-if="result.outline">章节大纲</h4>
          <p v-if="result.outline?.goal"><strong>目标：</strong>{{ result.outline.goal }}</p>
          <p v-if="result.outline?.key_events && result.outline.key_events.length > 0">
            <strong>关键事件：</strong>{{ result.outline.key_events.join(' → ') }}
          </p>
          <p v-if="result.outline?.dramatic_points">
            <strong>戏剧冲突点：</strong>{{ result.outline.dramatic_points }}个
          </p>
          
          <h4 v-if="result.character_development && Object.keys(result.character_development).length > 0">角色发展</h4>
          <div v-for="(development, character) in result.character_development" :key="character">
            <p><strong>{{ character }}：</strong>{{ development }}</p>
          </div>
          
          <h4 v-if="result.quality_metrics">质量评估</h4>
          <p v-if="result.quality_metrics?.score">
            <strong>评分：</strong>{{ result.quality_metrics.score }}/10
          </p>
          <p v-if="result.quality_metrics?.strengths && result.quality_metrics.strengths.length > 0">
            <strong>优点：</strong>{{ result.quality_metrics.strengths.join(', ') }}
          </p>
          <p v-if="result.quality_metrics?.improvement_areas && result.quality_metrics.improvement_areas.length > 0">
            <strong>改进点：</strong>{{ result.quality_metrics.improvement_areas.join(', ') }}
          </p>
          
          <h4>章节内容</h4>
          <div class="chapter-content" v-if="result.content">
            {{ result.content }}
          </div>
          <div v-else class="chapter-content-empty">
            <a-empty description="暂无内容" />
          </div>
        </div>
      </a-card>
    </div>
    
    <!-- 流式生成结果 -->
    <div v-if="streamContent" class="stream-result">
      <a-card title="流式生成结果">
        <a-tabs v-model:activeKey="streamTabActive">
          <a-tab-pane key="stream" tab="流式内容">
            <div class="stream-content">
              <pre>{{ streamContent }}</pre>
              <span v-if="streaming" class="stream-cursor">|</span>
            </div>
          </a-tab-pane>
          <a-tab-pane key="parsed" tab="解析结果" v-if="result">
            <div class="result-content">
              <h3>第{{ result.chapter_number || form.chapter_number || '?' }}章 {{ result.title || '未命名章节' }}</h3>
              <p v-if="result.word_count"><strong>字数：</strong>{{ result.word_count }}字</p>
              <p v-if="result.summary"><strong>摘要：</strong>{{ result.summary }}</p>
              
              <h4 v-if="result.outline">章节大纲</h4>
              <p v-if="result.outline?.goal"><strong>目标：</strong>{{ result.outline.goal }}</p>
              <p v-if="result.outline?.key_events && result.outline.key_events.length > 0">
                <strong>关键事件：</strong>{{ result.outline.key_events.join(' → ') }}
              </p>
              <p v-if="result.outline?.dramatic_points">
                <strong>戏剧冲突点：</strong>{{ result.outline.dramatic_points }}个
              </p>
              
              <h4 v-if="result.character_development && Object.keys(result.character_development).length > 0">角色发展</h4>
              <div v-for="(development, character) in result.character_development" :key="character">
                <p><strong>{{ character }}：</strong>{{ development }}</p>
              </div>
              
              <h4 v-if="result.quality_metrics">质量评估</h4>
              <p v-if="result.quality_metrics?.score">
                <strong>评分：</strong>{{ result.quality_metrics.score }}/10
              </p>
              <p v-if="result.quality_metrics?.strengths && result.quality_metrics.strengths.length > 0">
                <strong>优点：</strong>{{ result.quality_metrics.strengths.join(', ') }}
              </p>
              <p v-if="result.quality_metrics?.improvement_areas && result.quality_metrics.improvement_areas.length > 0">
                <strong>改进点：</strong>{{ result.quality_metrics.improvement_areas.join(', ') }}
              </p>
              
              <h4>章节内容</h4>
              <div class="chapter-content" v-if="result.content">
                {{ result.content }}
              </div>
              <div v-else class="chapter-content-empty">
                <a-empty description="暂无内容" />
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
        <template #extra v-if="!streaming && result">
          <a-space>
            <a-button 
              type="primary" 
              @click="handleSaveToDatabase"
              :loading="saving"
            >
              {{ saving ? '保存中...' : '保存到数据库' }}
            </a-button>
            <a-button @click="handleViewRawData">
              查看原始数据
            </a-button>
          </a-space>
        </template>
      </a-card>
    </div>
    
    <!-- 原始数据查看模态框 -->
    <a-modal
      v-model:open="rawDataModalVisible"
      title="原始数据"
      width="80%"
      :footer="null"
    >
      <pre style="background: #f5f5f5; padding: 16px; border-radius: 6px; max-height: 500px; overflow-y: auto;">{{ JSON.stringify(rawData, null, 2) }}</pre>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import { api, streamGenerate } from '@/utils/api'

const props = defineProps<{
  novel: any
}>()

const emit = defineEmits<{
  generated: [type: string, data: any]
}>()

const form = reactive({
  chapter_number: undefined as number | undefined,
  outline_id: undefined as string | undefined,
  chapter_goal: '',
  previous_summary: '',
  characters_involved: [],
  plot_templates: [],
  llm_model_id: ''
})

const rules = {
  chapter_number: [{ required: true, message: '请选择或输入章节号' }],
  chapter_goal: [{ required: true, message: '请输入章节目标' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const characters = ref([])
const charactersLoading = ref(false)
const chapters = ref([])
const chaptersLoading = ref(false)
const outlines = ref([])
const outlinesLoading = ref(false)
const models = ref([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const result = ref(null)
const streamContent = ref('')
const streamTabActive = ref('stream')
const saving = ref(false)
const rawData = ref(null)
const rawDataModalVisible = ref(false)
const selectedChapterInfo = ref<any>(null)
const chapterNumberOptions = ref<Array<{ value: number; label: string }>>([])

const fetchCharacters = async () => {
  try {
    charactersLoading.value = true
    const response = await api.get(`/characters/${props.novel.id}`)
    characters.value = response.data
  } catch (error) {
    message.error('获取角色列表失败')
  } finally {
    charactersLoading.value = false
  }
}

// 填充前情提要（上一个章节的摘要和正文内容）
const fillPreviousSummary = () => {
  if (!form.chapter_number) {
    return
  }
  
  const previousChapterNumber = form.chapter_number - 1
  if (previousChapterNumber > 0) {
    // 从已存在的章节中查找上一个章节
    const previousChapter = chapters.value.find((c: any) => c.chapter_number === previousChapterNumber)
    if (previousChapter) {
      // 构建前情提要：包含摘要和正文内容
      const parts: string[] = []
      
      // 添加摘要
      if (previousChapter.summary) {
        parts.push(`【章节摘要】\n${previousChapter.summary}`)
      }
      
      // 添加正文内容
      if (previousChapter.content) {
        parts.push(`【章节正文】\n${previousChapter.content}`)
      }
      
      // 如果两部分都有内容，用换行分隔；如果只有一部分，直接使用
      if (parts.length > 0) {
        form.previous_summary = parts.join('\n\n')
      } else {
        // 如果既没有摘要也没有正文，清空前情提要
        form.previous_summary = ''
      }
    } else {
      // 如果确实没有上一个章节，清空前情提要
      form.previous_summary = ''
    }
  } else {
    // 如果是第1章，没有前情提要
    form.previous_summary = ''
  }
}

const fetchChapters = async () => {
  try {
    chaptersLoading.value = true
    const response = await api.get(`/chapters/${props.novel.id}`)
    chapters.value = response.data || []
    updateChapterNumberOptions()
    
    // 章节列表加载完成后，如果已选择章节号，自动填充前情提要
    if (form.chapter_number) {
      fillPreviousSummary()
    }
  } catch (error) {
    message.error('获取章节列表失败')
  } finally {
    chaptersLoading.value = false
  }
}

const fetchOutlines = async () => {
  try {
    outlinesLoading.value = true
    const response = await api.get(`/outlines/${props.novel.id}`)
    outlines.value = response.data || []
    updateChapterNumberOptions()
  } catch (error) {
    message.error('获取大纲列表失败')
  } finally {
    outlinesLoading.value = false
  }
}

// 更新章节号选项
const updateChapterNumberOptions = () => {
  const options: Array<{ value: number; label: string }> = []
  
  // 从已存在的章节中获取章节号
  const existingChapterNumbers = new Set<number>()
  chapters.value.forEach((chapter: any) => {
    if (chapter.chapter_number) {
      existingChapterNumbers.add(chapter.chapter_number)
      options.push({
        value: chapter.chapter_number,
        label: `第${chapter.chapter_number}章 ${chapter.title || ''} (已存在)`
      })
    }
  })
  
  // 从大纲中获取章节号
  outlines.value.forEach((outline: any) => {
    if (outline.chapters && Array.isArray(outline.chapters)) {
      outline.chapters.forEach((chapterInfo: any) => {
        if (chapterInfo.chapter_number && !existingChapterNumbers.has(chapterInfo.chapter_number)) {
          existingChapterNumbers.add(chapterInfo.chapter_number)
          options.push({
            value: chapterInfo.chapter_number,
            label: `第${chapterInfo.chapter_number}章 ${chapterInfo.title || ''} (来自大纲)`
          })
        }
      })
    }
  })
  
  // 计算下一个章节号
  let nextChapterNumber = 1
  if (existingChapterNumbers.size > 0) {
    const maxChapterNumber = Math.max(...Array.from(existingChapterNumbers))
    nextChapterNumber = maxChapterNumber + 1
  }
  
  // 添加"新建章节"选项
  options.push({
    value: nextChapterNumber,
    label: `第${nextChapterNumber}章 (新建)`
  })
  
  // 按章节号排序
  options.sort((a, b) => a.value - b.value)
  
  chapterNumberOptions.value = options
  
  // 如果没有选择章节号，默认选择下一个章节号
  if (!form.chapter_number && nextChapterNumber > 0) {
    form.chapter_number = nextChapterNumber
  }
}

// 处理大纲变化
const handleOutlineChange = async (outlineId: string | undefined) => {
  if (!outlineId) {
    selectedChapterInfo.value = null
    updateChapterNumberOptions()
    return
  }
  
  // 先从本地查找，如果没有完整信息则从API获取
  let outline = outlines.value.find((o: any) => o.id === outlineId)
  if (!outline || !outline.chapters || outline.chapters.length === 0) {
    try {
      const response = await api.get(`/outline/${outlineId}`)
      outline = response.data
      
      // 更新本地缓存
      const index = outlines.value.findIndex((o: any) => o.id === outlineId)
      if (index !== -1) {
        outlines.value[index] = outline
      } else {
        outlines.value.push(outline)
      }
    } catch (error) {
      console.error('获取大纲详情失败:', error)
      message.error('获取大纲详情失败')
      return
    }
  }
  
  if (outline && outline.chapters) {
    updateChapterNumberOptions()
    
    // 如果已选择章节号，尝试从大纲中找到对应的章节信息
    if (form.chapter_number) {
      const chapterInfo = outline.chapters.find((c: any) => c.chapter_number === form.chapter_number)
      if (chapterInfo) {
        selectedChapterInfo.value = chapterInfo
        // 自动填充章节目标（动态更新）
        if (chapterInfo.outline?.goal) {
          form.chapter_goal = chapterInfo.outline.goal
        } else {
          // 如果大纲中没有章节目标，清空章节目标
          form.chapter_goal = ''
        }
        // 自动选择角色
        if (chapterInfo.characters && chapterInfo.characters.length > 0) {
          form.characters_involved = chapterInfo.characters
        }
      } else {
        selectedChapterInfo.value = null
        // 如果大纲中没有找到对应章节，清空章节目标
        if (form.chapter_number) {
          form.chapter_goal = ''
        }
      }
    }
  }
}

// 处理章节号搜索
const handleChapterNumberSearch = (value: string) => {
  if (value && /^\d+$/.test(value)) {
    const chapterNumber = parseInt(value)
    if (chapterNumber > 0) {
      // 检查是否已存在该选项
      const exists = chapterNumberOptions.value.some(opt => opt.value === chapterNumber)
      if (!exists) {
        chapterNumberOptions.value.push({
          value: chapterNumber,
          label: `第${chapterNumber}章 (手动输入)`
        })
        // 按章节号排序
        chapterNumberOptions.value.sort((a, b) => a.value - b.value)
      }
    }
  }
}

// 监听章节号变化
watch(() => form.chapter_number, async (newChapterNumber) => {
  if (!newChapterNumber) {
    selectedChapterInfo.value = null
    return
  }
  
  // 自动填充上一个章节的摘要到前情提要
  fillPreviousSummary()
  
  // 如果选择了大纲，从大纲中查找章节信息
  if (form.outline_id) {
    // 先从本地查找
    let outline = outlines.value.find((o: any) => o.id === form.outline_id)
    
    // 如果本地没有完整信息，则从API获取
    if (!outline || !outline.chapters || outline.chapters.length === 0) {
      try {
        const response = await api.get(`/outline/${form.outline_id}`)
        outline = response.data
        
        // 更新本地缓存
        const index = outlines.value.findIndex((o: any) => o.id === form.outline_id)
        if (index !== -1) {
          outlines.value[index] = outline
        } else {
          outlines.value.push(outline)
        }
      } catch (error) {
        console.error('获取大纲详情失败:', error)
        selectedChapterInfo.value = null
        return
      }
    }
    
    if (outline && outline.chapters) {
      const chapterInfo = outline.chapters.find((c: any) => c.chapter_number === newChapterNumber)
      if (chapterInfo) {
        selectedChapterInfo.value = chapterInfo
        // 自动填充章节目标（动态更新，不管之前是否有值）
        if (chapterInfo.outline?.goal) {
          form.chapter_goal = chapterInfo.outline.goal
        } else {
          // 如果大纲中没有章节目标，清空章节目标
          form.chapter_goal = ''
        }
        // 自动选择角色
        if (chapterInfo.characters && chapterInfo.characters.length > 0) {
          form.characters_involved = chapterInfo.characters
        }
      } else {
        selectedChapterInfo.value = null
        // 如果大纲中没有找到对应章节，清空章节目标
        form.chapter_goal = ''
      }
    }
  } else {
    // 从已存在的章节中查找
    const existingChapter = chapters.value.find((c: any) => c.chapter_number === newChapterNumber)
    if (existingChapter) {
      selectedChapterInfo.value = {
        title: existingChapter.title,
        summary: existingChapter.summary,
        characters: []
      }
      // 如果已存在的章节有outline信息，也填充章节目标
      if (existingChapter.outline?.goal) {
        form.chapter_goal = existingChapter.outline.goal
      } else {
        // 如果已存在章节没有章节目标，清空
        form.chapter_goal = ''
      }
    } else {
      selectedChapterInfo.value = null
      // 如果没有找到已存在的章节，清空章节目标
      form.chapter_goal = ''
    }
  }
})

const fetchModels = async () => {
  try {
    modelsLoading.value = true
    const response = await api.get('/llm-models')
    models.value = response.data.items || response.data
  } catch (error) {
    message.error('获取模型列表失败')
  } finally {
    modelsLoading.value = false
  }
}

// 获取章节大纲信息
const getChapterOutline = async () => {
  if (!form.outline_id || !form.chapter_number) {
    return null
  }
  
  // 先从本地缓存的大纲列表中查找
  let outline = outlines.value.find((o: any) => o.id === form.outline_id)
  
  // 如果本地没有找到，或者没有完整的章节信息，则从API获取
  if (!outline || !outline.chapters || outline.chapters.length === 0) {
    try {
      const response = await api.get(`/outline/${form.outline_id}`)
      outline = response.data
      
      // 更新本地缓存
      const index = outlines.value.findIndex((o: any) => o.id === form.outline_id)
      if (index !== -1) {
        outlines.value[index] = outline
      } else {
        outlines.value.push(outline)
      }
    } catch (error) {
      console.error('获取大纲详情失败:', error)
      return null
    }
  }
  
  if (!outline || !outline.chapters) {
    return null
  }
  
  // 从大纲中找到对应章节号的章节信息
  const chapterInfo = outline.chapters.find((c: any) => c.chapter_number === form.chapter_number)
  return chapterInfo || null
}

const handleGenerate = async () => {
  try {
    generating.value = true
    
    // 获取章节大纲信息
    const chapterOutline = await getChapterOutline()
    
    // 构建输入数据
    const inputData: any = {
      chapter_number: form.chapter_number,
      chapter_goal: form.chapter_goal,
      previous_summary: form.previous_summary,
      characters_involved: form.characters_involved.map((name: string) => {
        const character = characters.value.find((c: any) => c.name === name)
        return {
          name: name,
          soul_profile: character?.soul_profile || {},
          core_attributes: character?.core_attributes || {},
          current_state: '当前状态'
        }
      }),
      plot_templates: form.plot_templates
    }
    
    // 如果选择了大纲，添加大纲ID和章节大纲信息
    if (form.outline_id) {
      inputData.outline_id = form.outline_id
      
      // 如果找到了章节大纲信息，添加到输入数据中
      if (chapterOutline) {
        inputData.characters_outline = chapterOutline
      }
    }
    
    // 添加用户输入（可选）
    inputData.user_input = {
      chapter_goal: form.chapter_goal,
      plot_templates: form.plot_templates
    }
    
    const response = await api.post('/generate/chapter', {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: inputData,
      stream: false
    })
    
    // 确保章节号存在
    if (!response.data.chapter_number && form.chapter_number) {
      response.data.chapter_number = form.chapter_number
    }
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'chapter', response.data)
    
    // 刷新章节列表
    await fetchChapters()
  } catch (error: any) {
    message.error(error.response?.data?.error || '生成失败')
  } finally {
    generating.value = false
  }
}

// 清理和解析JSON内容
const cleanAndParseJSON = (content: string) => {
  // 移除Markdown代码块标记
  let cleaned = content.replace(/```json\n?/g, '').replace(/```\n?/g, '')
  
  // 移除可能的其他Markdown标记
  cleaned = cleaned.replace(/^```.*$/gm, '')
  
  // 移除前后的空白字符
  cleaned = cleaned.trim()
  
  // 查找JSON部分的结束位置（最后一个完整的JSON对象）
  let jsonEnd = -1
  let braceCount = 0
  let inString = false
  let escapeNext = false
  
  for (let i = 0; i < cleaned.length; i++) {
    const char = cleaned[i]
    
    if (escapeNext) {
      escapeNext = false
      continue
    }
    
    if (char === '\\') {
      escapeNext = true
      continue
    }
    
    if (char === '"' && !escapeNext) {
      inString = !inString
      continue
    }
    
    if (!inString) {
      if (char === '{') {
        if (braceCount === 0) {
          // 找到第一个 {，开始计数
        }
        braceCount++
      } else if (char === '}') {
        braceCount--
        if (braceCount === 0) {
          jsonEnd = i
          break
        }
      }
    }
  }
  
  // 提取JSON部分和内容部分
  let jsonPart = ''
  let contentPart = ''
  
  if (jsonEnd !== -1) {
    jsonPart = cleaned.substring(0, jsonEnd + 1)
    contentPart = cleaned.substring(jsonEnd + 1).trim()
  } else {
    // 如果找不到完整的JSON，尝试使用正则表达式
    const jsonMatch = cleaned.match(/\{[\s\S]*?\}/)
    if (jsonMatch) {
      jsonPart = jsonMatch[0]
      const matchIndex = cleaned.indexOf(jsonMatch[0])
      contentPart = cleaned.substring(matchIndex + jsonMatch[0].length).trim()
    } else {
      jsonPart = cleaned
    }
  }
  
  // 解析JSON
  try {
    const parsed = JSON.parse(jsonPart)
    
    // 如果JSON中没有content字段，但找到了内容部分，则添加进去
    if (!parsed.content && contentPart) {
      // 移除可能的markdown标记
      const cleanedContent = contentPart
        .replace(/```json\n?/g, '')
        .replace(/```\n?/g, '')
        .replace(/【正文开始】/g, '')
        .replace(/【正文结束】/g, '')
        .trim()
      
      if (cleanedContent && cleanedContent.length > 10) {
        parsed.content = cleanedContent
      }
    }
    
    return parsed
  } catch (error) {
    // 如果解析失败，尝试更宽松的方式
    const firstBrace = cleaned.indexOf('{')
    const lastBrace = cleaned.lastIndexOf('}')
    
    if (firstBrace !== -1 && lastBrace !== -1 && lastBrace > firstBrace) {
      try {
        const parsed = JSON.parse(cleaned.substring(firstBrace, lastBrace + 1))
        // 同样尝试提取内容部分
        if (!parsed.content) {
          const afterJson = cleaned.substring(lastBrace + 1).trim()
          if (afterJson && afterJson.length > 10) {
            parsed.content = afterJson
              .replace(/```json\n?/g, '')
              .replace(/```\n?/g, '')
              .replace(/【正文开始】/g, '')
              .trim()
          }
        }
        return parsed
      } catch (e) {
        console.error('提取JSON部分后解析失败:', e)
      }
    }
    
    throw error
  }
}

const handleStreamGenerate = async () => {
  try {
    streaming.value = true
    streamContent.value = ''
    
    // 获取章节大纲信息
    const chapterOutline = await getChapterOutline()
    
    // 构建输入数据
    const inputData: any = {
      chapter_number: form.chapter_number,
      chapter_goal: form.chapter_goal,
      previous_summary: form.previous_summary,
      characters_involved: form.characters_involved.map((name: string) => {
        const character = characters.value.find((c: any) => c.name === name)
        return {
          name: name,
          soul_profile: character?.soul_profile || {},
          core_attributes: character?.core_attributes || {},
          current_state: '当前状态'
        }
      }),
      plot_templates: form.plot_templates
    }
    
    // 如果选择了大纲，添加大纲ID和章节大纲信息
    if (form.outline_id) {
      inputData.outline_id = form.outline_id
      
      // 如果找到了章节大纲信息，添加到输入数据中
      if (chapterOutline) {
        inputData.characters_outline = chapterOutline
      }
    }
    
    // 添加用户输入（可选）
    inputData.user_input = {
      chapter_goal: form.chapter_goal,
      plot_templates: form.plot_templates
    }
    
    await streamGenerate(
      '/generate/chapter',
      {
        novel_id: props.novel.id,
        llm_model_id: form.llm_model_id,
        input_data: inputData,
        stream: true
      },
      (content: string) => {
        streamContent.value += content
      },
      () => {
        streaming.value = false
        message.success('流式生成完成')
        try {
          const parsed = cleanAndParseJSON(streamContent.value)
          
          // 确保章节号存在
          if (!parsed.chapter_number && form.chapter_number) {
            parsed.chapter_number = form.chapter_number
          }
          
          result.value = parsed
          streamTabActive.value = 'parsed' // 自动切换到解析结果标签页
          emit('generated', 'chapter', parsed)
          // 刷新章节列表
          fetchChapters()
        } catch (error) {
          console.error('解析生成内容失败:', error)
          console.error('原始内容:', streamContent.value)
          message.error('解析生成内容失败，请查看流式内容')
        }
      },
      (error: string) => {
        streaming.value = false
        message.error('流式生成失败: ' + error)
      }
    )
  } catch (error: any) {
    streaming.value = false
    message.error('流式生成失败')
  }
}

// 保存到数据库
const handleSaveToDatabase = async () => {
  if (!result.value || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }
  
  // 验证必需字段
  const chapterNumber = result.value.chapter_number || form.chapter_number
  if (!chapterNumber) {
    message.error('章节号不能为空')
    return
  }
  
  const title = result.value.title || '未命名章节'
  const content = result.value.content || ''
  if (!content) {
    message.error('章节内容不能为空')
    return
  }
  
  try {
    saving.value = true
    
    // 构建章节数据，确保所有必需字段都有值
    const chapterData: any = {
      novel_id: props.novel.id,
      chapter_number: chapterNumber,
      title: title,
      content: content,
      summary: result.value.summary || '',
      outline: {
        goal: result.value.outline?.goal || '',
        key_events: result.value.outline?.key_events || [],
        dramatic_points: result.value.outline?.dramatic_points || 0
      },
      quality_metrics: {
        score: result.value.quality_metrics?.score || 0,
        strengths: result.value.quality_metrics?.strengths || [],
        improvement_areas: result.value.quality_metrics?.improvement_areas || []
      },
      character_development: result.value.character_development || {}
    }
    
    // 调用章节创建接口
    const response = await api.post('/chapter', chapterData)
    
    message.success('已保存到数据库')
    
    // 触发父组件更新
    emit('generated', 'chapter', response.data)
    
    // 刷新章节列表
    await fetchChapters()
  } catch (error: any) {
    message.error('保存失败: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

// 查看原始数据
const handleViewRawData = () => {
  rawData.value = result.value
  rawDataModalVisible.value = true
}

// 监听novel变化，自动填充表单
watch(() => props.novel, (newNovel, oldNovel) => {
  console.log('ChapterGenerate: novel数据变化', { newNovel, oldNovel })
  if (newNovel) {
    try {
      autoFillForm()
    } catch (error) {
      console.error('ChapterGenerate: 自动填充表单失败', error)
    }
  }
}, { immediate: true })

onMounted(() => {
  if (props.novel) {
    fetchCharacters()
    fetchChapters()
    fetchOutlines()
    fetchModels()
    
    // 自动填充表单数据
    autoFillForm()
  }
})

// 自动填充表单数据
const autoFillForm = () => {
  if (!props.novel) return
  
  // 从小说基本信息填充章节目标
  if (props.novel.project_blueprint?.core_conflict) {
    form.chapter_goal = `推进核心冲突"${props.novel.project_blueprint.core_conflict}"的发展`
  }
  
  // 从AI上下文填充章节目标
  if (props.novel.ai_context?.current_focus) {
    form.chapter_goal += `\n当前焦点：${props.novel.ai_context.current_focus}`
  }
  
  // 从风格指南填充章节目标
  if (props.novel.ai_context?.style_guideline) {
    form.chapter_goal += `\n风格要求：${props.novel.ai_context.style_guideline}`
  }
  
  // 从情感基调填充章节目标
  if (props.novel.ai_context?.emotional_tone) {
    form.chapter_goal += `\n情感基调：${props.novel.ai_context.emotional_tone}`
  }
  
  // 自动选择主要角色
  if (characters.value.length > 0) {
    const mainCharacters = characters.value.filter(char => char.role === 'protagonist')
    if (mainCharacters.length > 0) {
      form.characters_involved = mainCharacters.map(char => char.name)
    }
  }
}
</script>

<style scoped>
.chapter-generate {
  max-width: 800px;
}

.result-section {
  margin-top: 24px;
}

.result-content {
  line-height: 1.6;
}

.result-content h3 {
  color: #1890ff;
  margin-bottom: 16px;
}

.result-content h4 {
  color: #1890ff;
  margin: 16px 0 8px 0;
}

.result-content p {
  margin-bottom: 8px;
}

.chapter-content {
  background: #f9f9f9;
  padding: 16px;
  border-radius: 6px;
  margin-top: 16px;
  white-space: pre-wrap;
  line-height: 1.8;
  max-height: 600px;
  overflow-y: auto;
  min-height: 100px;
}

.chapter-content-empty {
  margin-top: 16px;
  padding: 40px;
  text-align: center;
  background: #fafafa;
  border-radius: 6px;
  border: 1px dashed #d9d9d9;
}

.stream-result {
  margin-top: 24px;
}

.stream-content {
  position: relative;
  background: #f5f5f5;
  padding: 16px;
  border-radius: 6px;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 400px;
  overflow-y: auto;
}

.stream-cursor {
  animation: blink 1s infinite;
  color: #1890ff;
  font-weight: bold;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
