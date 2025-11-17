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
    
    <!-- 批量生成区域 -->
    <a-divider>批量生成</a-divider>
    
    <a-card title="根据大纲生成所有角色" size="small" style="margin-bottom: 16px;">
      <template #extra>
        <a-button 
          type="primary" 
          @click="handleBatchGenerate"
          :loading="batchGenerating"
          :disabled="!availableOutlines.length"
        >
          批量生成角色 ({{ availableOutlines.length }})
        </a-button>
      </template>
      
      <a-form layout="vertical">
        <a-form-item label="选择大纲">
          <a-select 
            v-model:value="batchForm.outline_id" 
            placeholder="请选择用于生成角色的大纲"
            :loading="outlinesLoading"
          >
            <a-select-option 
              v-for="outline in availableOutlines" 
              :key="outline.id" 
              :value="outline.id"
            >
              {{ outline.title }} ({{ outline.chapters?.length || 0 }}章)
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="用户要求">
          <a-textarea 
            v-model:value="batchForm.user_requirements" 
            placeholder="请输入对角色生成的具体要求，如：希望重点突出某个角色、特定的角色关系等"
            :rows="3"
          />
        </a-form-item>
        
        <a-form-item label="LLM模型">
          <a-select v-model:value="batchForm.llm_model_id" placeholder="请选择LLM模型">
            <a-select-option v-for="model in models" :key="model.id" :value="model.id">
              {{ model.display_name || model.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-card>
    
    <!-- 流式生成内容 -->
    <div v-if="batchGenerating || batchStreamContent" class="stream-content">
      <a-card title="生成中..." size="small">
        <template #extra>
          <a-tag v-if="batchGenerating" color="processing">生成中</a-tag>
          <a-tag v-else color="success">生成完成</a-tag>
        </template>
        
        <div class="stream-text">
          <pre>{{ batchStreamContent }}</pre>
        </div>
      </a-card>
    </div>
    
    <!-- 批量生成结果 -->
    <div v-if="batchResult.length > 0" class="batch-result">
      <a-card title="批量生成结果" size="small">
        <template #extra>
          <a-space>
            <a-tag color="green">{{ batchResult.length }} 个角色</a-tag>
            <a-button 
              type="primary" 
              @click="handleBatchSaveToDatabase"
              :loading="batchSaving"
            >
              {{ batchSaving ? '保存中...' : '保存到数据库' }}
            </a-button>
            <a-button @click="handleViewBatchRawData">
              查看原始数据
            </a-button>
          </a-space>
        </template>
        
        <a-list :data-source="batchResult" item-layout="vertical">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta
                :title="item.name"
                :description="`类型: ${getCharacterTypeText(item.type)}`"
              />
              <template #actions>
                <a @click="viewCharacter(item)">查看详情</a>
              </template>
              <div class="character-summary">
                <p><strong>核心驱动力：</strong>{{ item.soul_profile?.motivations?.core_drive || '未定义' }}</p>
                <p><strong>修炼境界：</strong>{{ item.core_attributes?.cultivation_level || '未定义' }}</p>
              </div>
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
    
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

// 批量生成相关状态
const batchGenerating = ref(false)
const batchResult = ref([])
const batchStreamContent = ref('')
const batchSaving = ref(false)
const availableOutlines = ref([])
const outlinesLoading = ref(false)

const batchForm = reactive({
  outline_id: '',
  user_requirements: '',
  llm_model_id: ''
})

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

// 获取可用大纲
const fetchAvailableOutlines = async () => {
  if (!props.novel?.id) {
    console.log('没有novel.id，跳过获取大纲列表')
    return
  }
  
  try {
    console.log('开始获取大纲列表，novel.id:', props.novel.id)
    outlinesLoading.value = true
    const response = await api.get(`/outlines/${props.novel.id}`)
    availableOutlines.value = response.data || []
    console.log('获取到大纲列表:', availableOutlines.value)
  } catch (error) {
    console.warn('获取大纲列表失败:', error)
    availableOutlines.value = []
  } finally {
    outlinesLoading.value = false
  }
}

// 批量生成角色（流式）
const handleBatchGenerate = async () => {
  console.log('开始批量生成角色', {
    outline_id: batchForm.outline_id,
    llm_model_id: batchForm.llm_model_id,
    novel_id: props.novel?.id
  })
  
  if (!batchForm.outline_id || !batchForm.llm_model_id) {
    message.error('请选择大纲和LLM模型')
    return
  }
  
  try {
    batchGenerating.value = true
    batchResult.value = []
    batchStreamContent.value = ''
    
    console.log('开始获取数据...')
    
    // 获取大纲、故事核心、世界观数据
    const [outlineData, storyCoreData, worldviewData] = await Promise.all([
      fetchOutlineData(batchForm.outline_id),
      fetchStoryCoreData(),
      fetchWorldviewData()
    ])
    
    console.log('数据获取完成，开始流式生成...', {
      outlineData: outlineData.substring(0, 100) + '...',
      storyCoreData: storyCoreData.substring(0, 100) + '...',
      worldviewData: worldviewData.substring(0, 100) + '...'
    })
    
    await streamGenerate(
      '/generate/characters-from-outline',
      {
        novel_id: props.novel.id,
        llm_model_id: batchForm.llm_model_id,
        outline_content: outlineData,
        story_core: storyCoreData,
        worldview: worldviewData,
        user_requirements: batchForm.user_requirements
      },
      (content: string) => {
        console.log('收到流式内容:', content)
        // 累积流式内容
        batchStreamContent.value += content
      },
      () => {
        console.log('流式生成完成')
        batchGenerating.value = false
        message.success('批量生成完成')
        // 解析生成的内容
        try {
          const parsed = cleanAndParseJSON(batchStreamContent.value)
          if (parsed && parsed.characters) {
            batchResult.value = parsed.characters
            message.success(`成功生成 ${parsed.characters.length} 个角色`)
            emit('generated', 'character', parsed.characters)
          } else {
            message.warning('生成的内容格式不正确')
          }
        } catch (error) {
          console.error('解析生成内容失败:', error)
          message.error('解析生成内容失败，请查看流式内容')
        }
      },
      (error: string) => {
        console.error('流式生成错误:', error)
        batchGenerating.value = false
        message.error('批量生成失败: ' + error)
      }
    )
  } catch (error: any) {
    console.error('批量生成异常:', error)
    batchGenerating.value = false
    message.error(error.response?.data?.error || error.message || '批量生成失败')
  }
}

// 获取大纲数据
const fetchOutlineData = async (outlineId: string) => {
  try {
    console.log('获取大纲数据，outlineId:', outlineId)
    const response = await api.get(`/outline/${outlineId}`)
    const outline = response.data
    
    console.log('获取到大纲数据:', outline)
    
    if (!outline) {
      throw new Error('大纲数据为空')
    }
    
    let content = `标题：${outline.title || '未命名'}\n概要：${outline.summary || '无概要'}\n`
    
    if (outline.key_themes && outline.key_themes.length > 0) {
      content += `关键主题：${outline.key_themes.join('、')}\n`
    }
    
    if (outline.story_arcs && outline.story_arcs.length > 0) {
      content += '故事弧线：\n'
      outline.story_arcs.forEach((arc: any) => {
        content += `- ${arc.name}：${arc.description} (第${arc.start_chapter}-${arc.end_chapter}章)\n`
      })
    }
    
    if (outline.chapters && outline.chapters.length > 0) {
      content += '章节概览：\n'
      outline.chapters.forEach((chapter: any) => {
        content += `第${chapter.chapter_number}章：${chapter.title} - ${chapter.summary}\n`
      })
    }
    
    return content
  } catch (error) {
    console.error('获取大纲数据失败:', error)
    throw new Error('获取大纲数据失败: ' + (error as Error).message)
  }
}

// 获取故事核心数据
const fetchStoryCoreData = async () => {
  try {
    const response = await api.get(`/story-cores/${props.novel.id}`)
    const storyCores = response.data
    
    if (!storyCores || storyCores.length === 0) {
      throw new Error('没有找到故事核心数据')
    }
    
    const storyCore = storyCores[0]
    let content = `核心冲突：${storyCore.core_conflict}\n`
    content += `主题：${storyCore.theme}\n`
    content += `创新点：${storyCore.innovation}\n`
    content += `商业潜力：${storyCore.commercial_potential}\n`
    content += `目标受众：${storyCore.target_audience}\n`
    
    return content
  } catch (error) {
    console.error('获取故事核心数据失败:', error)
    throw new Error('获取故事核心数据失败')
  }
}

// 获取世界观数据
const fetchWorldviewData = async () => {
  try {
    const response = await api.get(`/worldview/${props.novel.id}`)
    const worldview = response.data
    
    let content = '力量体系：\n'
    content += `- 名称：${worldview.power_system?.name || ''}\n`
    content += `- 修炼方法：${worldview.power_system?.cultivation_method || ''}\n`
    content += `- 限制：${worldview.power_system?.limitations || ''}\n`
    
    if (worldview.power_system?.levels && worldview.power_system.levels.length > 0) {
      content += `- 等级：${worldview.power_system.levels.join('、')}\n`
    }
    
    content += '社会结构：\n'
    content += `- 等级制度：${worldview.society_structure?.hierarchy || ''}\n`
    content += `- 经济体系：${worldview.society_structure?.economic_system || ''}\n`
    
    if (worldview.society_structure?.major_factions && worldview.society_structure.major_factions.length > 0) {
      content += '- 主要势力：\n'
      worldview.society_structure.major_factions.forEach((faction: any) => {
        content += `  * ${faction.name} (${faction.type}) - ${faction.influence}\n`
      })
    }
    
    content += '地理环境：\n'
    if (worldview.geography?.major_regions && worldview.geography.major_regions.length > 0) {
      content += `- 主要区域：${worldview.geography.major_regions.join('、')}\n`
    }
    if (worldview.geography?.special_locations && worldview.geography.special_locations.length > 0) {
      content += `- 特殊地点：${worldview.geography.special_locations.join('、')}\n`
    }
    
    if (worldview.special_rules && worldview.special_rules.length > 0) {
      content += `特殊规则：${worldview.special_rules.join('、')}\n`
    }
    
    return content
  } catch (error) {
    console.error('获取世界观数据失败:', error)
    throw new Error('获取世界观数据失败')
  }
}

// 清理Markdown格式并提取JSON
const cleanAndParseJSON = (content: string) => {
  try {
    // 移除Markdown代码块标记
    let cleaned = content.replace(/```json\s*/g, '').replace(/```\s*/g, '')
    
    // 尝试找到JSON对象
    const jsonMatch = cleaned.match(/\{[\s\S]*\}/)
    if (jsonMatch) {
      cleaned = jsonMatch[0]
    }
    
    // 解析JSON
    return JSON.parse(cleaned)
  } catch (error) {
    console.error('清理和解析JSON失败:', error)
    throw error
  }
}

// 查看角色详情
const viewCharacter = (character: any) => {
  message.info(`查看角色: ${character.name}`)
  // 这里可以打开角色详情模态框或跳转到详情页面
}

// 批量保存到数据库
const handleBatchSaveToDatabase = async () => {
  if (!batchResult.value || batchResult.value.length === 0 || !props.novel?.id) {
    message.error('没有可保存的数据')
    return
  }
  
  try {
    batchSaving.value = true
    
    // 批量保存角色
    const savePromises = batchResult.value.map(async (character: any) => {
      const characterData = {
        novel_id: props.novel.id,
        name: character.name,
        type: character.type,
        core_attributes: character.core_attributes || {},
        soul_profile: character.soul_profile || {}
      }
      
      return api.post('/character', characterData)
    })
    
    const responses = await Promise.all(savePromises)
    
    message.success(`成功保存 ${responses.length} 个角色到数据库`)
    
    // 触发父组件更新
    emit('generated', 'character', responses.map(r => r.data))
    
    // 清空批量结果
    batchResult.value = []
    batchStreamContent.value = ''
    
  } catch (error: any) {
    console.error('批量保存角色失败:', error)
    message.error('保存失败: ' + (error.response?.data?.error || error.message))
  } finally {
    batchSaving.value = false
  }
}

// 查看批量生成原始数据
const handleViewBatchRawData = () => {
  rawData.value = batchResult.value
  rawDataModalVisible.value = true
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
    fetchAvailableOutlines()
    
    // 自动填充表单数据
    autoFillForm()
  }
})

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

.stream-content {
  margin-bottom: 16px;
}

.stream-text {
  max-height: 300px;
  overflow-y: auto;
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #d9d9d9;
}

.stream-text pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.batch-result {
  margin-top: 16px;
}

.character-summary {
  margin-top: 8px;
}

.character-summary p {
  margin: 4px 0;
  font-size: 12px;
  color: #666;
}
</style>
