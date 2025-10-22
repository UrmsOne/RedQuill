<template>
  <div class="worldview-generate" v-if="novel">
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
      <a-form-item label="选择故事核心" name="selected_concept">
        <a-select v-model:value="form.selected_concept" placeholder="请选择故事核心" :loading="storyCoresLoading">
          <a-select-option v-for="core in storyCores" :key="core.id" :value="core.title">
            {{ core.title }}
          </a-select-option>
        </a-select>
      </a-form-item>
      
      <a-form-item label="类型" name="genre">
        <a-input v-model:value="form.genre" placeholder="如：玄幻、都市、历史等" />
      </a-form-item>
      
      <a-form-item label="个人建议" name="user_ideas">
        <a-textarea 
          v-model:value="form.user_ideas" 
          placeholder="请输入您对世界观的个人建议和想法，例如：希望有魔法体系、科技与魔法结合、特定的社会结构等"
          :rows="3"
          :maxlength="500"
          show-count
        />
      </a-form-item>
      
      <!-- 故事核心信息显示 -->
      <div v-if="form.story_core.title" class="story-core-info">
        <a-card title="故事核心信息" size="small">
          <div class="story-core-content">
            <p><strong>标题：</strong>{{ form.story_core.title }}</p>
            <p><strong>核心冲突：</strong>{{ form.story_core.core_conflict }}</p>
            <p><strong>主题：</strong>{{ form.story_core.theme }}</p>
            <p><strong>创新点：</strong>{{ form.story_core.innovation }}</p>
            <p><strong>商业潜力：</strong>{{ form.story_core.commercial_potential }}</p>
            <p><strong>目标受众：</strong>{{ form.story_core.target_audience }}</p>
          </div>
        </a-card>
      </div>
      
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
          <h4>修炼体系</h4>
          <p><strong>名称：</strong>{{ result.power_system?.name }}</p>
          <p><strong>等级：</strong>{{ result.power_system?.levels?.join(' → ') }}</p>
          <p><strong>修炼方法：</strong>{{ result.power_system?.cultivation_method }}</p>
          <p><strong>限制：</strong>{{ result.power_system?.limitations }}</p>
          
          <h4>社会结构</h4>
          <p><strong>等级制度：</strong>{{ result.society_structure?.hierarchy }}</p>
          <p><strong>经济体系：</strong>{{ result.society_structure?.economic_system }}</p>
          
          <h4>地理环境</h4>
          <p><strong>主要区域：</strong>{{ result.geography?.major_regions?.join(', ') }}</p>
          <p><strong>特殊地点：</strong>{{ result.geography?.special_locations?.join(', ') }}</p>
          
          <h4>特殊规则</h4>
          <ul>
            <li v-for="rule in result.special_rules" :key="rule">{{ rule }}</li>
          </ul>
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
        <a-tabs v-model:activeKey="streamTabActive">
          <a-tab-pane key="stream" tab="流式内容">
            <div class="stream-content">
              <pre>{{ streamContent }}</pre>
              <span v-if="streaming" class="stream-cursor">|</span>
            </div>
          </a-tab-pane>
          
          <a-tab-pane key="parsed" tab="解析结果" v-if="result">
            <div class="result-content">
              <h4>修炼体系</h4>
              <p><strong>名称：</strong>{{ result.power_system?.name }}</p>
              <p><strong>等级：</strong>{{ result.power_system?.levels?.join(' → ') }}</p>
              <p><strong>修炼方法：</strong>{{ result.power_system?.cultivation_method }}</p>
              <p><strong>限制：</strong>{{ result.power_system?.limitations }}</p>
              
              <h4>社会结构</h4>
              <p><strong>等级制度：</strong>{{ result.society_structure?.hierarchy }}</p>
              <p><strong>经济体系：</strong>{{ result.society_structure?.economic_system }}</p>
              
              <h4>地理环境</h4>
              <p><strong>主要区域：</strong>{{ result.geography?.major_regions?.join(', ') }}</p>
              <p><strong>特殊地点：</strong>{{ result.geography?.special_locations?.join(', ') }}</p>
              
              <h4>特殊规则</h4>
              <ul>
                <li v-for="rule in result.special_rules" :key="rule">{{ rule }}</li>
              </ul>
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
          </a-tab-pane>
        </a-tabs>
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
  selected_concept: '',
  genre: '',
  llm_model_id: '',
  user_ideas: '', // 用户个人建议
  // 故事核心相关字段
  story_core: {
    title: '',
    core_conflict: '',
    theme: '',
    innovation: '',
    commercial_potential: '',
    target_audience: ''
  }
})

const rules = {
  selected_concept: [{ required: true, message: '请选择故事核心' }],
  genre: [{ required: true, message: '请输入类型' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const storyCores = ref([])
const storyCoresLoading = ref(false)
const models = ref([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const result = ref(null)
const streamContent = ref('')
const streamTabActive = ref('stream')

const fetchStoryCores = async () => {
  try {
    storyCoresLoading.value = true
    const response = await api.get(`/story-cores/${props.novel.id}`)
    storyCores.value = response.data
  } catch (error) {
    message.error('获取故事核心失败')
  } finally {
    storyCoresLoading.value = false
  }
}

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
    const response = await api.post('/generate/worldview', {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        selected_concept: form.selected_concept,
        genre: form.genre,
        user_ideas: form.user_ideas,
        // 将story_core的数据往外提一层
        title: form.story_core.title,
        core_conflict: form.story_core.core_conflict,
        theme: form.story_core.theme,
        innovation: form.story_core.innovation,
        commercial_potential: form.story_core.commercial_potential,
        target_audience: form.story_core.target_audience
      },
      stream: false
    })
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'worldview', response.data)
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
  
  // 尝试解析JSON
  try {
    return JSON.parse(cleaned)
  } catch (error) {
    // 如果解析失败，尝试提取JSON部分
    const jsonMatch = cleaned.match(/\{[\s\S]*\}/)
    if (jsonMatch) {
      return JSON.parse(jsonMatch[0])
    }
    throw error
  }
}

const handleStreamGenerate = async () => {
  try {
    streaming.value = true
    streamContent.value = ''
    
    await streamGenerate(
      '/generate/worldview',
      {
        novel_id: props.novel.id,
        llm_model_id: form.llm_model_id,
        input_data: {
          selected_concept: form.selected_concept,
          genre: form.genre,
          user_ideas: form.user_ideas,
          // 将story_core的数据往外提一层
          title: form.story_core.title,
          core_conflict: form.story_core.core_conflict,
          theme: form.story_core.theme,
          innovation: form.story_core.innovation,
          commercial_potential: form.story_core.commercial_potential,
          target_audience: form.story_core.target_audience
        },
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
          result.value = parsed
          streamTabActive.value = 'parsed' // 自动切换到解析结果标签页
          emit('generated', 'worldview', parsed)
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
const saving = ref(false)
const handleSaveToDatabase = async () => {
  if (!result.value || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }
  
  try {
    saving.value = true
    
    // 构建世界观数据
    const worldviewData = {
      novel_id: props.novel.id,
      power_system: result.value.power_system || {},
      society_structure: result.value.society_structure || {},
      geography: result.value.geography || {},
      special_rules: result.value.special_rules || []
    }
    
    // 调用世界观创建接口
    const response = await api.post('/worldview', worldviewData)
    
    message.success('已保存到数据库')
    
    // 触发父组件更新
    emit('generated', 'worldview', response.data)
  } catch (error: any) {
    message.error('保存失败: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

// 查看原始数据
const rawData = ref(null)
const rawDataModalVisible = ref(false)
const handleViewRawData = () => {
  rawData.value = result.value
  rawDataModalVisible.value = true
}

// 监听故事核心选择变化，更新故事核心信息
watch(() => form.selected_concept, (newConcept) => {
  if (newConcept && storyCores.value.length > 0) {
    const selectedCore = storyCores.value.find(core => core.title === newConcept)
    if (selectedCore) {
      form.story_core = {
        title: selectedCore.title || '',
        core_conflict: selectedCore.core_conflict || '',
        theme: selectedCore.theme || '',
        innovation: selectedCore.innovation || '',
        commercial_potential: selectedCore.commercial_potential || '',
        target_audience: selectedCore.target_audience || ''
      }
    }
  }
})

// 监听novel变化，自动填充表单
watch(() => props.novel, (newNovel, oldNovel) => {
  console.log('WorldviewGenerate: novel数据变化', { newNovel, oldNovel })
  if (newNovel) {
    try {
      autoFillForm()
    } catch (error) {
      console.error('WorldviewGenerate: 自动填充表单失败', error)
    }
  }
}, { immediate: true })

onMounted(() => {
  if (props.novel) {
    fetchStoryCores()
    fetchModels()
    
    // 自动填充表单数据
    autoFillForm()
  }
})

// 自动填充表单数据
const autoFillForm = () => {
  if (!props.novel) return
  
  try {
    console.log('WorldviewGenerate: 自动填充表单', { novel: props.novel, storyCores: storyCores.value })
    
    // 从小说基本信息填充类型
    if (props.novel.project_blueprint) {
      form.genre = props.novel.project_blueprint.genre || ''
    }
    
    // 如果有故事核心，自动选择最新的
    if (storyCores.value.length > 0) {
      form.selected_concept = storyCores.value[0].title
    }
  } catch (error) {
    console.error('WorldviewGenerate: 自动填充表单失败', error)
  }
}
</script>

<style scoped>
.worldview-generate {
  max-width: 800px;
}

.result-section {
  margin-top: 24px;
}

.result-content {
  line-height: 1.6;
}

.result-content h4 {
  color: #1890ff;
  margin: 16px 0 8px 0;
}

.result-content p {
  margin-bottom: 8px;
}

.result-content ul {
  margin-left: 20px;
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

.story-core-info {
  margin: 16px 0;
}

.story-core-content p {
  margin: 8px 0;
  line-height: 1.6;
}

.story-core-content strong {
  color: #1890ff;
  font-weight: 600;
}
</style>
