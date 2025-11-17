<template>
  <div class="story-core-generate">
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
      <a-row :gutter="24">
        <a-col :span="12">
          <a-form-item label="类型" name="genre">
            <a-input v-model:value="form.genre" placeholder="如：玄幻、都市、历史等" />
          </a-form-item>
        </a-col>
        
        <a-col :span="12">
          <a-form-item label="子类型" name="sub_genre">
            <a-input v-model:value="form.sub_genre" placeholder="如：重生、系统、修仙等" />
          </a-form-item>
        </a-col>
      </a-row>
      
      <a-form-item label="用户想法" name="user_ideas">
        <a-textarea 
          v-model:value="form.user_ideas" 
          placeholder="请描述您对故事的想法和期望..."
          :rows="4"
        />
      </a-form-item>
      
      <a-form-item label="目标受众" name="target_audience">
        <a-input v-model:value="form.target_audience" placeholder="目标读者群体" />
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
        <div class="result-content">
          <h3>{{ result.title }}</h3>
          <p><strong>核心冲突：</strong>{{ result.core_conflict }}</p>
          <p><strong>主题：</strong>{{ result.theme }}</p>
          <p><strong>创新点：</strong>{{ result.innovation }}</p>
          <p><strong>商业潜力：</strong>{{ result.commercial_potential }}</p>
          <p><strong>目标受众：</strong>{{ result.target_audience }}</p>
        </div>
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
              <!-- 多个concepts选择 -->
              <div v-if="result.concepts && result.concepts.length > 1" class="concepts-selection">
                <h3>请选择故事核心概念：</h3>
                <a-radio-group v-model:value="selectedConceptIndex" @change="handleConceptChange">
                  <div v-for="(concept, index) in result.concepts" :key="index" class="concept-option">
                    <a-radio :value="index" class="concept-radio">
                      <div class="concept-card">
                        <h4>{{ concept.title }}</h4>
                        <p><strong>核心冲突：</strong>{{ concept.core_conflict }}</p>
                        <p><strong>主题：</strong>{{ concept.theme }}</p>
                        <p><strong>创新点：</strong>{{ concept.innovation }}</p>
                        <p><strong>商业潜力：</strong>{{ concept.commercial_potential }}</p>
                        <p><strong>目标受众：</strong>{{ concept.target_audience }}</p>
                      </div>
                    </a-radio>
                  </div>
                </a-radio-group>
              </div>
              
              <!-- 单个concept或已选择的概念 -->
              <div v-else-if="result.concepts && result.concepts.length === 1" class="single-concept">
                <h3>{{ result.concepts[0].title }}</h3>
                <p><strong>核心冲突：</strong>{{ result.concepts[0].core_conflict }}</p>
                <p><strong>主题：</strong>{{ result.concepts[0].theme }}</p>
                <p><strong>创新点：</strong>{{ result.concepts[0].innovation }}</p>
                <p><strong>商业潜力：</strong>{{ result.concepts[0].commercial_potential }}</p>
                <p><strong>目标受众：</strong>{{ result.concepts[0].target_audience }}</p>
              </div>
              
              <!-- 直接显示结果（兼容旧格式） -->
              <div v-else class="direct-result">
                <h3>{{ result.title }}</h3>
                <p><strong>核心冲突：</strong>{{ result.core_conflict }}</p>
                <p><strong>主题：</strong>{{ result.theme }}</p>
                <p><strong>创新点：</strong>{{ result.innovation }}</p>
                <p><strong>商业潜力：</strong>{{ result.commercial_potential }}</p>
                <p><strong>目标受众：</strong>{{ result.target_audience }}</p>
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

const props = defineProps<{
  novel: any
}>()

const emit = defineEmits<{
  generated: [type: string, data: any]
}>()

const form = reactive({
  genre: '',
  sub_genre: '',
  user_ideas: '',
  target_audience: '',
  llm_model_id: ''
})

const rules = {
  genre: [{ required: true, message: '请输入类型' }],
  sub_genre: [{ required: true, message: '请输入子类型' }],
  user_ideas: [{ required: true, message: '请输入用户想法' }],
  target_audience: [{ required: true, message: '请输入目标受众' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const models = ref([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const result = ref(null)
const streamContent = ref('')
const saving = ref(false)
const rawData = ref(null)
const rawDataModalVisible = ref(false)
const streamTabActive = ref('stream')
const selectedConceptIndex = ref(0)
const selectedConcept = ref(null)

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

const handleGenerate = async () => {
  try {
    generating.value = true
    const response = await api.post('/generate/story-core', {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        genre: form.genre,
        sub_genre: form.sub_genre,
        user_ideas: form.user_ideas,
        target_audience: form.target_audience
      },
      stream: false
    })
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'story-core', response.data)
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
    
    await streamGenerate(
      '/generate/story-core',
      {
        novel_id: props.novel.id,
        llm_model_id: form.llm_model_id,
        input_data: {
          genre: form.genre,
          sub_genre: form.sub_genre,
          user_ideas: form.user_ideas,
          target_audience: form.target_audience
        },
        stream: true
      },
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
          
          // 初始化选择状态
          selectedConceptIndex.value = 0
          if (parsed.concepts && parsed.concepts.length > 0) {
            selectedConcept.value = parsed.concepts[0]
          }
          
          streamTabActive.value = 'parsed' // 自动切换到解析结果标签页
          emit('generated', 'story-core', parsed)
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

// 保存到数据库
const handleSaveToDatabase = async () => {
  if (!result.value || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }

  try {
    saving.value = true
    
    // 确定要保存的数据
    let dataToSave = result.value
    
    // 如果有多个concepts且用户已选择，保存选中的概念
    if (result.value.concepts && result.value.concepts.length > 1 && selectedConcept.value) {
      dataToSave = selectedConcept.value
    } else if (result.value.concepts && result.value.concepts.length === 1) {
      // 如果只有一个concept，直接使用
      dataToSave = result.value.concepts[0]
    }
    
    // 构建故事核心数据
    const storyCoreData = {
      novel_id: props.novel.id,
      title: dataToSave.title || '未命名故事核心',
      core_conflict: dataToSave.core_conflict || '',
      theme: dataToSave.theme || '',
      innovation: dataToSave.innovation || '',
      commercial_potential: dataToSave.commercial_potential || '',
      target_audience: dataToSave.target_audience || ''
    }

    // 调用故事核心创建接口
    const response = await api.post('/story-core', storyCoreData)
    
    message.success('已保存到数据库')
    
    // 触发父组件更新
    emit('generated', 'story-core', response.data)
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

// 处理概念选择变化
const handleConceptChange = (e: any) => {
  const index = e.target.value
  selectedConceptIndex.value = index
  
  if (result.value && result.value.concepts && result.value.concepts[index]) {
    selectedConcept.value = result.value.concepts[index]
  }
}

// 监听novel变化，自动填充表单
watch(() => props.novel, (newNovel, oldNovel) => {
  console.log('StoryCoreGenerate: novel数据变化', { newNovel, oldNovel })
  if (newNovel) {
    console.log('StoryCoreGenerate: 开始自动填充表单')
    autoFillForm()
  }
}, { immediate: true })

onMounted(() => {
  fetchModels()
  
  // 自动填充表单数据
  autoFillForm()
})

// 自动填充表单数据
const autoFillForm = () => {
  console.log('StoryCoreGenerate: 开始自动填充表单', props.novel)
  if (!props.novel) {
    console.log('StoryCoreGenerate: 没有小说数据，跳过填充')
    return
  }
  
  // 从小说基本信息填充
  if (props.novel.project_blueprint) {
    form.genre = props.novel.project_blueprint.genre || ''
    form.sub_genre = props.novel.project_blueprint.sub_genre || ''
    form.target_audience = props.novel.project_blueprint.target_audience || ''
    console.log('StoryCoreGenerate: 填充基本信息', {
      genre: form.genre,
      sub_genre: form.sub_genre,
      target_audience: form.target_audience
    })
  }
  
  // 从AI上下文填充用户想法
  if (props.novel.ai_context?.recent_summary) {
    form.user_ideas = props.novel.ai_context.recent_summary
    console.log('StoryCoreGenerate: 填充用户想法', form.user_ideas)
  }
  
  // 如果有核心冲突，可以作为用户想法的补充
  if (props.novel.project_blueprint?.core_conflict && !form.user_ideas) {
    form.user_ideas = `核心冲突：${props.novel.project_blueprint.core_conflict}`
    console.log('StoryCoreGenerate: 填充核心冲突', form.user_ideas)
  }
  
  console.log('StoryCoreGenerate: 表单填充完成', form)
}
</script>

<style scoped>
.story-core-generate {
  max-width: 800px;
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

.concepts-selection {
  margin-bottom: 16px;
}

.concept-option {
  margin-bottom: 16px;
}

.concept-radio {
  width: 100%;
}

.concept-card {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 16px;
  margin-top: 8px;
  background: #fafafa;
  transition: all 0.3s;
}

.concept-radio:checked + .concept-card {
  border-color: #1890ff;
  background: #e6f7ff;
}

.concept-card h4 {
  margin: 0 0 12px 0;
  color: #1890ff;
  font-size: 16px;
}

.concept-card p {
  margin: 8px 0;
  line-height: 1.5;
}

.single-concept,
.direct-result {
  padding: 16px;
  background: #f9f9f9;
  border-radius: 6px;
}

.single-concept h3,
.direct-result h3 {
  margin: 0 0 16px 0;
  color: #1890ff;
  font-size: 18px;
}

.single-concept p,
.direct-result p {
  margin: 12px 0;
  line-height: 1.6;
}
</style>
