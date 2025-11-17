<template>
  <div class="outline-generate">
    <!-- 数据状态显示 -->
    <a-alert
      v-if="dataLoading"
      message="正在获取故事核心和世界观数据..."
      type="info"
      show-icon
      style="margin-bottom: 16px"
    />
    
    <a-alert
      v-else-if="!storyCore || !worldview"
      message="缺少必要数据"
      :description="`故事核心: ${storyCore ? '已获取' : '未获取'}, 世界观: ${worldview ? '已获取' : '未获取'}`"
      type="warning"
      show-icon
      style="margin-bottom: 16px"
    />
    
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
      <a-row :gutter="24">
        <a-col :span="12">
          <a-form-item label="总章节数" name="total_chapters">
            <a-input-number 
              v-model:value="form.total_chapters" 
              placeholder="请输入总章节数"
              :min="1"
              :max="1000"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        
        <a-col :span="12">
          <a-form-item label="小说类型" name="genre">
            <a-input v-model:value="form.genre" placeholder="如：玄幻、都市、历史等" />
          </a-form-item>
        </a-col>
      </a-row>
      
      <a-form-item label="目标受众" name="target_audience">
        <a-input v-model:value="form.target_audience" placeholder="目标读者群体" />
      </a-form-item>
      
      <a-form-item label="用户指导" name="user_ideas">
        <a-textarea 
          v-model:value="form.user_ideas" 
          placeholder="请输入您对大纲生成的具体要求，如：希望重点突出某个角色、特定的情节发展、特殊的叙事结构等"
          :rows="4"
        />
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
        <a-tabs v-model:activeKey="resultTabActive" size="small">
          <a-tab-pane key="overview" tab="概览">
            <div class="result-content">
              <h3>{{ result.title }}</h3>
              <p><strong>概要：</strong>{{ result.summary }}</p>
              <div v-if="result.key_themes && result.key_themes.length > 0">
                <p><strong>关键主题：</strong></p>
                <a-tag v-for="theme in result.key_themes" :key="theme" color="blue">
                  {{ theme }}
                </a-tag>
              </div>
            </div>
          </a-tab-pane>
          
          <a-tab-pane key="story-arcs" tab="故事弧线">
            <div v-if="result.story_arcs && result.story_arcs.length > 0">
              <a-list
                :data-source="result.story_arcs"
                item-layout="vertical"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="item.name"
                      :description="item.description"
                    />
                    <template #extra>
                      <a-tag color="green">第{{ item.start_chapter }}-{{ item.end_chapter }}章</a-tag>
                    </template>
                    <div class="arc-theme">
                      <strong>主题：</strong>{{ item.theme }}
                    </div>
                  </a-list-item>
                </template>
              </a-list>
            </div>
            <a-empty v-else description="暂无故事弧线" />
          </a-tab-pane>
          
          <a-tab-pane key="chapters" tab="章节列表">
            <div v-if="result.chapters && result.chapters.length > 0">
              <a-list
                :data-source="result.chapters"
                item-layout="vertical"
                :pagination="{ pageSize: 10, showSizeChanger: true }"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="`第${item.chapter_number}章 ${item.title}`"
                      :description="item.summary"
                    />
                    <template #extra>
                      <a-space direction="vertical" size="small">
                        <a-tag color="blue">{{ item.word_count }}字</a-tag>
                        <a-tag color="purple">{{ item.pov }}</a-tag>
                        <a-tag color="orange">{{ item.location }}</a-tag>
                      </a-space>
                    </template>
                    
                    <div class="chapter-details">
                      <div v-if="item.key_events && item.key_events.length > 0">
                        <strong>关键事件：</strong>
                        <ul>
                          <li v-for="event in item.key_events" :key="event">{{ event }}</li>
                        </ul>
                      </div>
                      
                      <div v-if="item.characters && item.characters.length > 0">
                        <strong>出场角色：</strong>
                        <a-tag v-for="character in item.characters" :key="character" color="cyan">
                          {{ character }}
                        </a-tag>
                      </div>
                      
                      <div v-if="item.outline && item.outline.goal">
                        <strong>章节目标：</strong>{{ item.outline.goal }}
                      </div>
                    </div>
                  </a-list-item>
                </template>
              </a-list>
            </div>
            <a-empty v-else description="暂无章节信息" />
          </a-tab-pane>
        </a-tabs>
        
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
      </a-card>
    </div>
    
    <!-- 流式生成结果 -->
    <div v-if="streamContent" class="stream-result">
      <a-card title="流式生成结果">
        <a-tabs v-model:activeKey="streamTabActive" size="small">
          <a-tab-pane key="stream" tab="流式内容">
            <div class="stream-content">
              <pre>{{ streamContent }}</pre>
              <span v-if="streaming" class="stream-cursor">|</span>
            </div>
          </a-tab-pane>
          <a-tab-pane key="parsed" tab="解析结果" v-if="result">
            <div class="result-content">
              <h3>{{ result.title }}</h3>
              <p><strong>概要：</strong>{{ result.summary }}</p>
              <div v-if="result.key_themes && result.key_themes.length > 0">
                <p><strong>关键主题：</strong></p>
                <a-tag v-for="theme in result.key_themes" :key="theme" color="blue">
                  {{ theme }}
                </a-tag>
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
    
    <!-- 原始数据模态框 -->
    <a-modal
      v-model:open="rawDataModalVisible"
      title="原始数据"
      width="80%"
      :footer="null"
    >
      <pre class="raw-data-content">{{ JSON.stringify(rawData, null, 2) }}</pre>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import { api, streamGenerate } from '@/utils/api'
import type { Outline } from '@/stores/novel'

const props = defineProps<{
  novel: any
}>()

const emit = defineEmits<{
  generated: [type: string, data: any]
}>()

const form = reactive({
  total_chapters: 50,
  genre: '',
  target_audience: '',
  user_ideas: '',
  llm_model_id: ''
})

const rules = {
  total_chapters: [{ required: true, message: '请输入总章节数' }],
  genre: [{ required: true, message: '请输入小说类型' }],
  target_audience: [{ required: true, message: '请输入目标受众' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const models = ref<any[]>([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const result = ref<any>(null)
const streamContent = ref('')
const saving = ref(false)
const rawData = ref<any>(null)
const rawDataModalVisible = ref(false)
const resultTabActive = ref('overview')
const streamTabActive = ref('stream')

// 故事核心和世界观数据
const storyCore = ref<any>(null)
const worldview = ref<any>(null)
const dataLoading = ref(false)

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

// 获取故事核心数据
const fetchStoryCore = async () => {
  if (!props.novel?.id) return
  
  try {
    dataLoading.value = true
    const response = await api.get(`/story-cores/${props.novel.id}`)
    console.log('获取故事核心响应:', response.data)
    
    // 后端返回的是数组，取第一个作为当前使用的故事核心
    if (Array.isArray(response.data) && response.data.length > 0) {
      storyCore.value = response.data[0]
      console.log('设置storyCore.value为:', storyCore.value)
    } else {
      storyCore.value = null
      console.log('没有找到故事核心数据')
    }
  } catch (error) {
    console.warn('获取故事核心失败:', error)
    storyCore.value = null
  } finally {
    dataLoading.value = false
  }
}

// 获取世界观数据
const fetchWorldview = async () => {
  if (!props.novel?.id) return
  
  try {
    dataLoading.value = true
    const response = await api.get(`/worldview/${props.novel.id}`)
    console.log('获取世界观响应:', response.data)
    
    // 后端返回的是单个对象
    worldview.value = response.data
    console.log('设置worldview.value为:', worldview.value)
  } catch (error) {
    console.warn('获取世界观失败:', error)
    worldview.value = null
  } finally {
    dataLoading.value = false
  }
}

// 获取所有必要数据
const fetchAllData = async () => {
  await Promise.all([
    fetchModels(),
    fetchStoryCore(),
    fetchWorldview()
  ])
}

const handleGenerate = async () => {
  try {
    generating.value = true
    
    const storyCoreText = getStoryCoreText()
    const worldviewText = getWorldviewText()
    
    const requestData = {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        story_core: storyCoreText,
        worldview: worldviewText,
        total_chapters: form.total_chapters,
        genre: form.genre,
        target_audience: form.target_audience,
        user_ideas: form.user_ideas
      },
      stream: false
    }
    
    console.log('发送生成请求:', requestData)
    console.log('story_core内容:', storyCoreText)
    console.log('worldview内容:', worldviewText)
    
    const response = await api.post('/generate/outline', requestData)
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'outline', response.data)
  } catch (error: any) {
    message.error(error.response?.data?.error || '生成失败')
  } finally {
    generating.value = false
  }
}

const handleStreamGenerate = async () => {
  try {
    streaming.value = true
    streamContent.value = ''
    
    const storyCoreText = getStoryCoreText()
    const worldviewText = getWorldviewText()
    
    const requestData = {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        story_core: storyCoreText,
        worldview: worldviewText,
        total_chapters: form.total_chapters,
        genre: form.genre,
        target_audience: form.target_audience,
        user_ideas: form.user_ideas
      },
      stream: true
    }
    
    console.log('发送流式生成请求:', requestData)
    console.log('story_core内容:', storyCoreText)
    console.log('worldview内容:', worldviewText)
    
    await streamGenerate(
      '/generate/outline',
      requestData,
      (content: string) => {
        streamContent.value += content
      },
      () => {
        streaming.value = false
        message.success('流式生成完成')
        // 解析生成的内容
        try {
          const parsed = cleanAndParseJSON(streamContent.value)
          result.value = parsed
          resultTabActive.value = 'overview'
          streamTabActive.value = 'parsed' // 自动切换到解析结果标签页
          emit('generated', 'outline', parsed)
        } catch (error) {
          console.error('解析生成内容失败:', error)
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

// 获取故事核心文本
const getStoryCoreText = () => {
  console.log('getStoryCoreText: storyCore.value =', storyCore.value)
  console.log('getStoryCoreText: props.novel?.extra_info?.story_core =', props.novel?.extra_info?.story_core)
  
  // 优先使用从后端获取的数据
  if (storyCore.value) {
    console.log('使用从后端获取的storyCore数据')
    if (typeof storyCore.value === 'string') {
      return storyCore.value
    }
    
    // 如果是对象，提取关键信息
    const result = `标题：${storyCore.value.title || '未命名'}
核心冲突：${storyCore.value.core_conflict || '未定义'}
主题：${storyCore.value.theme || '未定义'}
创新点：${storyCore.value.innovation || '未定义'}
商业潜力：${storyCore.value.commercial_potential || '未定义'}
目标受众：${storyCore.value.target_audience || '未定义'}`
    console.log('生成的storyCore文本:', result)
    return result
  }
  
  // 回退到从novel.extra_info获取
  if (props.novel?.extra_info?.story_core) {
    console.log('使用从novel.extra_info获取的storyCore数据')
    const storyCoreData = props.novel.extra_info.story_core
    if (typeof storyCoreData === 'string') {
      return storyCoreData
    }
    
    const result = `标题：${storyCoreData.title || '未命名'}
核心冲突：${storyCoreData.core_conflict || '未定义'}
主题：${storyCoreData.theme || '未定义'}
创新点：${storyCoreData.innovation || '未定义'}
商业潜力：${storyCoreData.commercial_potential || '未定义'}
目标受众：${storyCoreData.target_audience || '未定义'}`
    console.log('从extra_info生成的storyCore文本:', result)
    return result
  }
  
  console.log('没有找到storyCore数据，返回默认信息')
  return '暂无故事核心信息'
}

// 获取世界观文本
const getWorldviewText = () => {
  // 优先使用从后端获取的数据
  if (worldview.value) {
    if (typeof worldview.value === 'string') {
      return worldview.value
    }
    
    // 如果是对象，提取关键信息
    let text = ''
    if (worldview.value.power_system) {
      text += `力量体系：${worldview.value.power_system.name || '未定义'}\n`
      text += `修炼方式：${worldview.value.power_system.cultivation_method || '未定义'}\n`
    }
    if (worldview.value.society_structure) {
      text += `社会结构：${worldview.value.society_structure.hierarchy || '未定义'}\n`
      text += `经济体系：${worldview.value.society_structure.economic_system || '未定义'}\n`
    }
    if (worldview.value.geography) {
      text += `主要地域：${(worldview.value.geography.major_regions || []).join(', ')}\n`
      text += `特殊地点：${(worldview.value.geography.special_locations || []).join(', ')}\n`
    }
    if (worldview.value.special_rules) {
      text += `特殊规则：${(worldview.value.special_rules || []).join(', ')}\n`
    }
    
    return text || '世界观信息不完整'
  }
  
  // 回退到从novel.extra_info获取
  if (props.novel?.extra_info?.worldview) {
    const worldviewData = props.novel.extra_info.worldview
    if (typeof worldviewData === 'string') {
      return worldviewData
    }
    
    let text = ''
    if (worldviewData.power_system) {
      text += `力量体系：${worldviewData.power_system.name || '未定义'}\n`
      text += `修炼方式：${worldviewData.power_system.cultivation_method || '未定义'}\n`
    }
    if (worldviewData.society_structure) {
      text += `社会结构：${worldviewData.society_structure.hierarchy || '未定义'}\n`
      text += `经济体系：${worldviewData.society_structure.economic_system || '未定义'}\n`
    }
    if (worldviewData.geography) {
      text += `主要地域：${(worldviewData.geography.major_regions || []).join(', ')}\n`
      text += `特殊地点：${(worldviewData.geography.special_locations || []).join(', ')}\n`
    }
    if (worldviewData.special_rules) {
      text += `特殊规则：${(worldviewData.special_rules || []).join(', ')}\n`
    }
    
    return text || '世界观信息不完整'
  }
  
  return '暂无世界观信息'
}

// 保存到数据库
const handleSaveToDatabase = async () => {
  if (!result.value || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }

  try {
    saving.value = true
    
    // 构建大纲数据
    const outlineData = {
      novel_id: props.novel.id,
      title: result.value.title || '未命名大纲',
      summary: result.value.summary || '',
      chapters: result.value.chapters || [],
      story_arcs: result.value.story_arcs || [],
      key_themes: result.value.key_themes || []
    }

    // 调用大纲创建接口
    const response = await api.post('/outline', outlineData)
    
    message.success('已保存到数据库')
    
    // 触发父组件更新
    emit('generated', 'outline', response.data)
  } catch (error: any) {
    message.error('保存失败: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

// 查看原始数据
const handleViewRawData = () => {
  if (!result.value) {
    message.error('没有可查看的数据')
    return
  }
  
  rawData.value = result.value
  rawDataModalVisible.value = true
}

// 清理Markdown格式并提取JSON
const cleanAndParseJSON = (content: string) => {
  try {
    // 移除Markdown代码块标记
    let cleaned = content.replace(/```json\s*/g, '').replace(/```\s*/g, '')
    
    // 移除可能的其他Markdown标记
    cleaned = cleaned.replace(/```\s*/g, '')
    
    // 查找JSON对象的开始和结束
    const jsonStart = cleaned.indexOf('{')
    const jsonEnd = cleaned.lastIndexOf('}')
    
    if (jsonStart !== -1 && jsonEnd !== -1 && jsonEnd > jsonStart) {
      cleaned = cleaned.substring(jsonStart, jsonEnd + 1)
    }
    
    // 尝试解析JSON
    return JSON.parse(cleaned)
  } catch (error) {
    console.error('清理和解析JSON失败:', error)
    throw error
  }
}

// 自动填充表单数据
const autoFillForm = () => {
  console.log('OutlineGenerate: 开始自动填充表单', props.novel)
  if (!props.novel) {
    console.log('OutlineGenerate: 没有小说数据，跳过填充')
    return
  }
  
  // 从小说基本信息填充
  if (props.novel.project_blueprint) {
    form.genre = props.novel.project_blueprint.genre || ''
    form.target_audience = props.novel.project_blueprint.target_audience || ''
    form.total_chapters = props.novel.project_blueprint.total_chapters || 50
    console.log('OutlineGenerate: 填充基本信息', {
      genre: form.genre,
      target_audience: form.target_audience,
      total_chapters: form.total_chapters
    })
  }
  
  // 初始化user_ideas为空字符串
  form.user_ideas = ''
  
  console.log('OutlineGenerate: 表单填充完成', form)
}

// 监听novel变化，自动填充表单
watch(() => props.novel, (newNovel, oldNovel) => {
  console.log('OutlineGenerate: novel数据变化', { newNovel, oldNovel })
  if (newNovel) {
    console.log('OutlineGenerate: 开始自动填充表单')
    autoFillForm()
  }
}, { immediate: true })

onMounted(() => {
  fetchAllData()
  
  // 自动填充表单数据
  autoFillForm()
})
</script>

<style scoped>
.outline-generate {
  max-width: 1000px;
}

.raw-data-content {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 500px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.stream-content {
  max-height: 400px;
  overflow-y: auto;
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
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

.result-content p {
  margin-bottom: 12px;
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

.chapter-details {
  margin-top: 12px;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
}

.chapter-details > div {
  margin-bottom: 8px;
}

.chapter-details ul {
  margin: 4px 0 0 20px;
  padding: 0;
}

.chapter-details li {
  margin-bottom: 4px;
}

.arc-theme {
  margin-top: 8px;
  padding: 8px;
  background: #f0f8ff;
  border-radius: 4px;
  border-left: 3px solid #1890ff;
}
</style>