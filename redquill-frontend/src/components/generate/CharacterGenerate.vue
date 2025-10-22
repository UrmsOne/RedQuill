<template>
  <div class="character-generate" v-if="novel">
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
      <a-form-item label="角色类型" name="character_type">
        <a-select v-model:value="form.character_type" placeholder="请选择角色类型">
          <a-select-option value="protagonist">主角</a-select-option>
          <a-select-option value="antagonist">反派</a-select-option>
          <a-select-option value="supporting">配角</a-select-option>
          <a-select-option value="mentor">导师</a-select-option>
        </a-select>
      </a-form-item>
      
      <a-form-item label="角色要求" name="role_requirements">
        <a-textarea 
          v-model:value="form.role_requirements" 
          placeholder="请描述角色的具体要求..."
          :rows="3"
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
        <div class="result-content">
          <h3>{{ result.name }}</h3>
          <p><strong>类型：</strong>{{ getCharacterTypeText(result.type) }}</p>
          
          <h4>核心属性</h4>
          <p><strong>修炼等级：</strong>{{ result.core_attributes?.cultivation_level }}</p>
          <p><strong>当前物品：</strong>{{ result.core_attributes?.current_items?.join(', ') }}</p>
          <p><strong>能力：</strong>{{ result.core_attributes?.abilities?.join(', ') }}</p>
          
          <h4>性格特征</h4>
          <p><strong>核心特质：</strong>{{ result.soul_profile?.personality?.core_traits?.join(', ') }}</p>
          <p><strong>道德准则：</strong>{{ result.soul_profile?.personality?.moral_compass }}</p>
          <p><strong>内心冲突：</strong>{{ result.soul_profile?.personality?.internal_conflicts?.join(', ') }}</p>
          <p><strong>恐惧：</strong>{{ result.soul_profile?.personality?.fears?.join(', ') }}</p>
          <p><strong>欲望：</strong>{{ result.soul_profile?.personality?.desires?.join(', ') }}</p>
          
          <h4>背景故事</h4>
          <p><strong>出身：</strong>{{ result.soul_profile?.background?.origin }}</p>
          <p><strong>关键事件：</strong>{{ result.soul_profile?.background?.defining_events?.join(', ') }}</p>
          <p><strong>隐藏秘密：</strong>{{ result.soul_profile?.background?.hidden_secrets?.join(', ') }}</p>
          
          <h4>动机目标</h4>
          <p><strong>近期目标：</strong>{{ result.soul_profile?.motivations?.immediate_goal }}</p>
          <p><strong>长期目标：</strong>{{ result.soul_profile?.motivations?.long_term_goal }}</p>
          <p><strong>核心驱动力：</strong>{{ result.soul_profile?.motivations?.core_drive }}</p>
        </div>
        <template #extra>
          <a-space>
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
          <a-tab-pane key="parsed" tab="解析结果" v-if="parsedResult">
            <div class="result-content">
              <h3>{{ parsedResult?.name || '角色名称' }}</h3>
              <p><strong>类型：</strong>{{ getCharacterTypeText(parsedResult?.type) }}</p>
              
              <h4>核心属性</h4>
              <p><strong>修炼等级：</strong>{{ parsedResult?.core_attributes?.cultivation_level }}</p>
              <p><strong>当前物品：</strong>{{ parsedResult?.core_attributes?.current_items?.join(', ') }}</p>
              <p><strong>能力：</strong>{{ parsedResult?.core_attributes?.abilities?.join(', ') }}</p>
              
              <h4>性格特征</h4>
              <p><strong>核心特质：</strong>{{ parsedResult?.soul_profile?.personality?.core_traits?.join(', ') }}</p>
              <p><strong>道德准则：</strong>{{ parsedResult?.soul_profile?.personality?.moral_compass }}</p>
              <p><strong>内心冲突：</strong>{{ parsedResult?.soul_profile?.personality?.internal_conflicts?.join(', ') }}</p>
              <p><strong>恐惧：</strong>{{ parsedResult?.soul_profile?.personality?.fears?.join(', ') }}</p>
              <p><strong>欲望：</strong>{{ parsedResult?.soul_profile?.personality?.desires?.join(', ') }}</p>
              
              <h4>背景故事</h4>
              <p><strong>出身：</strong>{{ parsedResult?.soul_profile?.background?.origin }}</p>
              <p><strong>关键事件：</strong>{{ parsedResult?.soul_profile?.background?.defining_events?.join(', ') }}</p>
              <p><strong>隐藏秘密：</strong>{{ parsedResult?.soul_profile?.background?.hidden_secrets?.join(', ') }}</p>
              
              <h4>动机目标</h4>
              <p><strong>近期目标：</strong>{{ parsedResult?.soul_profile?.motivations?.immediate_goal }}</p>
              <p><strong>长期目标：</strong>{{ parsedResult?.soul_profile?.motivations?.long_term_goal }}</p>
              <p><strong>核心驱动力：</strong>{{ parsedResult?.soul_profile?.motivations?.core_drive }}</p>
            </div>
          </a-tab-pane>
        </a-tabs>
        <template #extra v-if="!streaming && parsedResult">
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
  character_type: 'protagonist',
  role_requirements: '',
  llm_model_id: ''
})

const rules = {
  character_type: [{ required: true, message: '请选择角色类型' }],
  role_requirements: [{ required: true, message: '请输入角色要求' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const models = ref([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const saving = ref(false)
const result = ref(null)
const streamContent = ref('')
const streamTabActive = ref('stream')
const parsedResult = ref(null)
const rawData = ref(null)
const rawDataModalVisible = ref(false)

const getCharacterTypeText = (type: string) => {
  const types: Record<string, string> = {
    protagonist: '主角',
    antagonist: '反派',
    supporting: '配角',
    mentor: '导师'
  }
  return types[type] || type
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
    const response = await api.post('/generate/character', {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        story_core: props.novel.project_blueprint?.core_conflict || '',
        worldview: '修真界世界观',
        character_type: form.character_type,
        role_requirements: form.role_requirements
      },
      stream: false
    })
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'character', response.data)
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
      '/generate/character',
      {
        novel_id: props.novel.id,
        llm_model_id: form.llm_model_id,
        input_data: {
          story_core: props.novel.project_blueprint?.core_conflict || '',
          worldview: '修真界世界观',
          character_type: form.character_type,
          role_requirements: form.role_requirements
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
          const parsed = JSON.parse(streamContent.value)
          parsedResult.value = parsed
          // 自动切换到解析结果tab
          streamTabActive.value = 'parsed'
          // 不设置result.value，避免与普通生成冲突
          // 不触发emit，避免重复触发
        } catch (error) {
          console.error('解析生成内容失败:', error)
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

// 监听novel变化，自动填充表单
watch(() => props.novel, (newNovel, oldNovel) => {
  console.log('CharacterGenerate: novel数据变化', { newNovel, oldNovel })
  if (newNovel) {
    try {
      autoFillForm()
    } catch (error) {
      console.error('CharacterGenerate: 自动填充表单失败', error)
    }
  }
}, { immediate: true })

onMounted(() => {
  if (props.novel) {
    fetchModels()
    
    // 自动填充表单数据
    autoFillForm()
  }
})

// 自动填充表单数据
const autoFillForm = () => {
  if (!props.novel) return
  
  // 从小说基本信息填充角色要求
  if (props.novel.project_blueprint?.core_conflict) {
    form.role_requirements = `基于核心冲突"${props.novel.project_blueprint.core_conflict}"的角色需求`
  }
  
  // 从AI上下文填充角色要求
  if (props.novel.ai_context?.current_focus) {
    form.role_requirements += `\n当前焦点：${props.novel.ai_context.current_focus}`
  }
  
  // 从风格指南填充角色要求
  if (props.novel.ai_context?.style_guideline) {
    form.role_requirements += `\n风格要求：${props.novel.ai_context.style_guideline}`
  }
}

// 保存到数据库
const handleSaveToDatabase = async () => {
  if (!parsedResult.value || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }

  try {
    saving.value = true
    
    // 构建保存请求
    const saveData = {
      novel_id: props.novel.id,
      name: (parsedResult.value as any)?.name || '未命名角色',
      type: form.character_type,
      core_attributes: (parsedResult.value as any)?.core_attributes || {},
      soul_profile: (parsedResult.value as any)?.soul_profile || {}
    }

    // 调用保存接口
    const response = await api.post('/character', saveData)
    
    message.success('角色保存成功')
    // 不触发emit，避免重复触发
    // emit('generated', 'character', response.data)
    
    // 清空流式内容
    streamContent.value = ''
    parsedResult.value = null
    
  } catch (error: any) {
    console.error('保存角色失败:', error)
    message.error('保存失败: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

// 查看原始数据
const handleViewRawData = () => {
  rawData.value = parsedResult.value
  rawDataModalVisible.value = true
}
</script>

<style scoped>
.character-generate {
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

.stream-actions {
  margin-top: 16px;
  text-align: right;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
