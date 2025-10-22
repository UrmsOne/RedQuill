<template>
  <div class="chapter-generate" v-if="novel">
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleGenerate"
    >
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
          placeholder="请描述前面的剧情发展..."
          :rows="3"
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
        <div class="result-content">
          <h3>第{{ result.chapter_number }}章 {{ result.title }}</h3>
          <p><strong>字数：</strong>{{ result.word_count }}字</p>
          <p><strong>摘要：</strong>{{ result.summary }}</p>
          
          <h4>章节大纲</h4>
          <p><strong>目标：</strong>{{ result.outline?.goal }}</p>
          <p><strong>关键事件：</strong>{{ result.outline?.key_events?.join(' → ') }}</p>
          <p><strong>戏剧冲突点：</strong>{{ result.outline?.dramatic_points }}个</p>
          
          <h4>角色发展</h4>
          <div v-for="(development, character) in result.character_development" :key="character">
            <p><strong>{{ character }}：</strong>{{ development }}</p>
          </div>
          
          <h4>质量评估</h4>
          <p><strong>评分：</strong>{{ result.quality_metrics?.score }}/10</p>
          <p><strong>优点：</strong>{{ result.quality_metrics?.strengths?.join(', ') }}</p>
          <p><strong>改进点：</strong>{{ result.quality_metrics?.improvement_areas?.join(', ') }}</p>
          
          <h4>章节内容</h4>
          <div class="chapter-content">
            {{ result.content }}
          </div>
        </div>
      </a-card>
    </div>
    
    <!-- 流式生成结果 -->
    <div v-if="streamContent" class="stream-result">
      <a-card title="流式生成结果">
        <div class="stream-content">
          <pre>{{ streamContent }}</pre>
          <span v-if="streaming" class="stream-cursor">|</span>
        </div>
      </a-card>
    </div>
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
  chapter_goal: '',
  previous_summary: '',
  characters_involved: [],
  plot_templates: [],
  llm_model_id: ''
})

const rules = {
  chapter_goal: [{ required: true, message: '请输入章节目标' }],
  llm_model_id: [{ required: true, message: '请选择LLM模型' }]
}

const characters = ref([])
const charactersLoading = ref(false)
const models = ref([])
const modelsLoading = ref(false)
const generating = ref(false)
const streaming = ref(false)
const result = ref(null)
const streamContent = ref('')

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
    const response = await api.post('/generate/chapter', {
      novel_id: props.novel.id,
      llm_model_id: form.llm_model_id,
      input_data: {
        novel_context: {
          story_core: props.novel.project_blueprint?.core_conflict || '',
          worldview: '修真界世界观',
          current_arc: '当前故事线'
        },
        chapter_goal: form.chapter_goal,
        characters_involved: form.characters_involved.map(name => {
          const character = characters.value.find(c => c.name === name)
          return {
            soul_profile: character?.soul_profile || {},
            core_attributes: character?.core_attributes || {},
            current_state: '当前状态'
          }
        }),
        previous_summary: form.previous_summary,
        plot_templates: form.plot_templates
      },
      stream: false
    })
    
    result.value = response.data
    message.success('生成成功')
    emit('generated', 'chapter', response.data)
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
      '/generate/chapter',
      {
        novel_id: props.novel.id,
        llm_model_id: form.llm_model_id,
        input_data: {
          novel_context: {
            story_core: props.novel.project_blueprint?.core_conflict || '',
            worldview: '修真界世界观',
            current_arc: '当前故事线'
          },
          chapter_goal: form.chapter_goal,
          characters_involved: form.characters_involved.map(name => {
            const character = characters.value.find(c => c.name === name)
            return {
              soul_profile: character?.soul_profile || {},
              core_attributes: character?.core_attributes || {},
              current_state: '当前状态'
            }
          }),
          previous_summary: form.previous_summary,
          plot_templates: form.plot_templates
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
          result.value = parsed
          emit('generated', 'chapter', parsed)
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
  max-height: 400px;
  overflow-y: auto;
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
