<template>
  <div class="prompts-page">
      <a-card title="Prompt模板管理" class="content-card">
        <template #extra>
          <a-button type="primary" @click="showCreateModal">
            <PlusOutlined />
            添加模板
          </a-button>
        </template>
        
        <a-table
          :columns="columns"
          :data-source="prompts"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'type'">
              <a-tag :color="getTypeColor(record.type)">
                {{ getTypeText(record.type) }}
              </a-tag>
            </template>
            
            <template v-if="column.key === 'public'">
              <a-tag :color="record.public ? 'green' : 'orange'">
                {{ record.public ? '公开' : '私有' }}
              </a-tag>
            </template>
            
            <template v-if="column.key === 'actions'">
              <a-space>
                <a @click="viewPrompt(record)">查看</a>
                <a @click="editPrompt(record)">编辑</a>
                <a @click="deletePrompt(record.id)">删除</a>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
      
      <!-- 创建模板模态框 -->
      <a-modal
        v-model:open="createModalVisible"
        title="添加Prompt模板"
        @ok="handleCreate"
        :confirm-loading="createLoading"
        width="800px"
      >
        <a-form
          :model="createForm"
          :rules="createRules"
          layout="vertical"
        >
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="模板名称" name="name">
                <a-input v-model:value="createForm.name" placeholder="如：故事核心生成" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="模板类型" name="type">
                <a-select v-model:value="createForm.type" placeholder="选择模板类型">
                  <a-select-option value="story_core">故事核心</a-select-option>
                  <a-select-option value="worldview">世界观</a-select-option>
                  <a-select-option value="character">角色</a-select-option>
                  <a-select-option value="chapter">章节</a-select-option>
                  <a-select-option value="quality_review">质量评估</a-select-option>
                  <a-select-option value="custom">自定义</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="创作阶段" name="phase">
                <a-select v-model:value="createForm.phase" placeholder="选择创作阶段">
                  <a-select-option value="story_core">故事核心</a-select-option>
                  <a-select-option value="worldview">世界观</a-select-option>
                  <a-select-option value="characters">角色</a-select-option>
                  <a-select-option value="outlining">大纲</a-select-option>
                  <a-select-option value="writing">写作</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="分类" name="category">
                <a-input v-model:value="createForm.category" placeholder="如：玄幻、都市等" />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-form-item label="模板内容" name="content">
            <a-textarea 
              v-model:value="createForm.content" 
              placeholder="请输入Prompt模板内容，使用{{变量名}}表示变量..."
              :rows="8"
            />
          </a-form-item>
          
          <a-form-item label="变量列表" name="variables">
            <a-select 
              v-model:value="createForm.variables" 
              mode="tags" 
              placeholder="输入变量名，按回车添加"
            />
          </a-form-item>
          
          <a-form-item label="描述" name="description">
            <a-textarea v-model:value="createForm.description" placeholder="模板用途描述" />
          </a-form-item>
          
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="标签" name="tags">
                <a-select 
                  v-model:value="createForm.tags" 
                  mode="tags" 
                  placeholder="输入标签，按回车添加"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="是否公开" name="public">
                <a-switch v-model:checked="createForm.public" />
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </a-modal>
      
      <!-- 查看模板模态框 -->
      <a-modal
        v-model:open="viewModalVisible"
        title="查看Prompt模板"
        :footer="null"
        width="800px"
      >
        <div v-if="currentPrompt" class="prompt-detail">
          <a-descriptions :column="2">
            <a-descriptions-item label="名称">{{ currentPrompt.name }}</a-descriptions-item>
            <a-descriptions-item label="类型">
              <a-tag :color="getTypeColor(currentPrompt.type)">
                {{ getTypeText(currentPrompt.type) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="创作阶段">{{ currentPrompt.phase }}</a-descriptions-item>
            <a-descriptions-item label="分类">{{ currentPrompt.category }}</a-descriptions-item>
            <a-descriptions-item label="使用次数">{{ currentPrompt.usage_count }}</a-descriptions-item>
            <a-descriptions-item label="是否公开">
              <a-tag :color="currentPrompt.public ? 'green' : 'orange'">
                {{ currentPrompt.public ? '公开' : '私有' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="变量" :span="2">
              <a-tag v-for="variable in currentPrompt.variables" :key="variable" color="blue">
                {{ variable }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="标签" :span="2">
              <a-tag v-for="tag in currentPrompt.tags" :key="tag" color="purple">
                {{ tag }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="描述" :span="2">{{ currentPrompt.description }}</a-descriptions-item>
          </a-descriptions>
          
          <a-divider>模板内容</a-divider>
          <div class="prompt-content">
            <pre>{{ currentPrompt.content }}</pre>
          </div>
        </div>
      </a-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { api } from '@/utils/api'

const loading = ref(false)
const prompts = ref([])
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`
})

const createModalVisible = ref(false)
const createLoading = ref(false)
const createForm = reactive({
  name: '',
  type: '',
  phase: '',
  category: '',
  content: '',
  variables: [],
  description: '',
  tags: [],
  public: false
})

const createRules = {
  name: [{ required: true, message: '请输入模板名称' }],
  type: [{ required: true, message: '请选择模板类型' }],
  phase: [{ required: true, message: '请选择创作阶段' }],
  content: [{ required: true, message: '请输入模板内容' }]
}

const viewModalVisible = ref(false)
const currentPrompt = ref(null)

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    width: 100
  },
  {
    title: '创作阶段',
    dataIndex: 'phase',
    key: 'phase',
    width: 100
  },
  {
    title: '分类',
    dataIndex: 'category',
    key: 'category',
    width: 100
  },
  {
    title: '使用次数',
    dataIndex: 'usage_count',
    key: 'usage_count',
    width: 100
  },
  {
    title: '是否公开',
    dataIndex: 'public',
    key: 'public',
    width: 100
  },
  {
    title: '创建时间',
    dataIndex: 'ctime',
    key: 'ctime',
    width: 180,
    customRender: ({ text }: any) => new Date(text * 1000).toLocaleString()
  },
  {
    title: '操作',
    key: 'actions',
    width: 150
  }
]

const getTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    story_core: 'purple',
    worldview: 'cyan',
    character: 'magenta',
    chapter: 'blue',
    quality_review: 'green',
    custom: 'orange'
  }
  return colors[type] || 'default'
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    story_core: '故事核心',
    worldview: '世界观',
    character: '角色',
    chapter: '章节',
    quality_review: '质量评估',
    custom: '自定义'
  }
  return texts[type] || type
}

const fetchPrompts = async (page = 1, pageSize = 10) => {
  try {
    loading.value = true
    const response = await api.get('/prompts', { params: { page, pageSize } })
    prompts.value = response.data.items || response.data
    pagination.value.total = response.data.pagination?.total || prompts.value.length
  } catch (error) {
    message.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchPrompts(pag.current, pag.pageSize)
}

const showCreateModal = () => {
  createModalVisible.value = true
}

const handleCreate = async () => {
  try {
    createLoading.value = true
    await api.post('/prompt', createForm)
    createModalVisible.value = false
    fetchPrompts()
    message.success('创建成功')
    
    // 重置表单
    Object.assign(createForm, {
      name: '',
      type: '',
      phase: '',
      category: '',
      content: '',
      variables: [],
      description: '',
      tags: [],
      public: false
    })
  } catch (error: any) {
    message.error(error.response?.data?.error || '创建失败')
  } finally {
    createLoading.value = false
  }
}

const viewPrompt = (prompt: any) => {
  currentPrompt.value = prompt
  viewModalVisible.value = true
}

const editPrompt = (prompt: any) => {
  message.info('编辑功能开发中')
}

const deletePrompt = async (id: string) => {
  try {
    await api.delete(`/prompt/${id}`)
    fetchPrompts()
    message.success('删除成功')
  } catch (error: any) {
    message.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  fetchPrompts()
})
</script>

<style scoped>
.prompts-page {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}

.prompt-detail {
  max-height: 600px;
  overflow-y: auto;
}

.prompt-content {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 6px;
  margin-top: 16px;
}

.prompt-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
}
</style>
