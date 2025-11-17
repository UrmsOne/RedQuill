<template>
  <div class="llm-models-page">
      <a-card title="LLM模型管理" class="content-card">
        <template #extra>
          <a-button type="primary" @click="showCreateModal">
            <PlusOutlined />
            添加模型
          </a-button>
        </template>
        
        <a-table
          :columns="columns"
          :data-source="models"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            
            <template v-if="column.key === 'capabilities'">
              <a-tag v-for="cap in record.capabilities" :key="cap" color="blue">
                {{ cap }}
              </a-tag>
            </template>
            
            <template v-if="column.key === 'actions'">
              <a-space>
                <a @click="testModel(record)">测试</a>
                <a @click="editModel(record)">编辑</a>
                <a @click="deleteModel(record.id, record.display_name || record.name)">删除</a>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
      
      <!-- 创建模型模态框 -->
      <a-modal
        v-model:open="createModalVisible"
        title="添加LLM模型"
        @ok="handleCreate"
        :confirm-loading="createLoading"
        width="600px"
      >
        <a-form
          :model="createForm"
          :rules="createRules"
          layout="vertical"
        >
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="模型名称" name="name">
                <a-input v-model:value="createForm.name" placeholder="如：deepseek-chat" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="显示名称" name="display_name">
                <a-input v-model:value="createForm.display_name" placeholder="前端显示名称" />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-form-item label="描述" name="description">
            <a-textarea v-model:value="createForm.description" placeholder="模型用途描述" />
          </a-form-item>
          
          <a-form-item label="能力" name="capabilities">
            <a-select v-model:value="createForm.capabilities" mode="multiple" placeholder="选择模型能力">
              <a-select-option value="chat">对话</a-select-option>
              <a-select-option value="generation">文本生成</a-select-option>
              <a-select-option value="analysis">分析</a-select-option>
              <a-select-option value="optimization">优化</a-select-option>
              <a-select-option value="creative_writing">创意写作</a-select-option>
              <a-select-option value="code_generation">代码生成</a-select-option>
              <a-select-option value="chinese_writing">中文写作</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-form-item label="状态" name="status">
            <a-select v-model:value="createForm.status" placeholder="选择状态">
              <a-select-option value="active">激活</a-select-option>
              <a-select-option value="inactive">未激活</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-divider>模型配置</a-divider>
          
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="提供商" name="provider">
                <a-select v-model:value="createForm.config.provider" placeholder="选择提供商">
                  <a-select-option value="openai">OpenAI</a-select-option>
                  <a-select-option value="azure">Azure OpenAI</a-select-option>
                  <a-select-option value="deepseek">DeepSeek</a-select-option>
                  <a-select-option value="doubao">豆包</a-select-option>
                  <a-select-option value="qwen">千问</a-select-option>
                  <a-select-option value="wenxin">文心一言</a-select-option>
                  <a-select-option value="ollama">Ollama</a-select-option>
                  <a-select-option value="mock">Mock</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="模型名称" name="model_name">
                <a-input v-model:value="createForm.config.model_name" placeholder="如：gpt-3.5-turbo" />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-form-item label="API密钥" name="api_key">
            <a-input-password v-model:value="createForm.config.api_key" placeholder="输入API密钥" />
          </a-form-item>
          
          <a-form-item label="基础URL" name="base_url">
            <a-input v-model:value="createForm.config.base_url" placeholder="如：https://api.openai.com" />
          </a-form-item>
          
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="温度" name="temperature">
                <a-input-number v-model:value="createForm.config.temperature" :min="0" :max="2" :step="0.1" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="最大Token" name="max_tokens">
                <a-input-number v-model:value="createForm.config.max_tokens" :min="1" :max="32000" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="超时时间(秒)" name="timeout">
                <a-input-number v-model:value="createForm.config.timeout" :min="10" :max="600" />
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </a-modal>
      
      <!-- 测试模型模态框 -->
      <a-modal
        v-model:open="testModalVisible"
        :title="`测试模型 - ${currentTestModel?.display_name || currentTestModel?.name || ''}`"
        @ok="handleTest"
        :confirm-loading="testLoading"
        :ok-button-props="{ disabled: testSuccess }"
        width="600px"
      >
        <div v-if="currentTestModel" class="model-info">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="模型名称">{{ currentTestModel.name }}</a-descriptions-item>
            <a-descriptions-item label="显示名称">{{ currentTestModel.display_name }}</a-descriptions-item>
            <a-descriptions-item label="提供商">{{ currentTestModel.config?.provider }}</a-descriptions-item>
            <a-descriptions-item label="模型">{{ currentTestModel.config?.model_name }}</a-descriptions-item>
          </a-descriptions>
          <a-divider />
        </div>
        
        <a-form
          :model="testForm"
          :rules="testRules"
          layout="vertical"
        >
          <a-form-item label="测试消息" name="message">
            <a-textarea 
              v-model:value="testForm.message" 
              placeholder="输入测试消息..."
              :rows="4"
              :disabled="testSuccess"
            />
          </a-form-item>
          
          <a-form-item label="流式响应" name="stream">
            <a-switch v-model:checked="testForm.stream" :disabled="testSuccess" />
          </a-form-item>
        </a-form>
        
        <div v-if="testResult || streamingResult || isStreaming" class="test-result">
          <a-divider>测试结果</a-divider>
          <div v-if="testSuccess" class="success-indicator">
            <a-alert 
              message="测试成功" 
              type="success" 
              show-icon 
              :closable="false"
            />
          </div>
          <div class="result-content">
            <div v-if="isStreaming" class="streaming-indicator">
              <a-spin size="small" />
              <span style="margin-left: 8px;">正在生成响应...</span>
            </div>
            <div v-if="testForm.stream" class="streaming-content">
              <pre>{{ streamingResult || (isStreaming ? '等待内容...' : '') }}</pre>
            </div>
            <div v-else class="normal-content">
              <pre>{{ testResult }}</pre>
            </div>
          </div>
        </div>
      </a-modal>
      
      <!-- 编辑模型模态框 -->
      <a-modal
        v-model:open="editModalVisible"
        title="编辑LLM模型"
        @ok="handleEdit"
        :confirm-loading="editLoading"
        width="600px"
      >
        <a-form
          :model="editForm"
          :rules="editRules"
          layout="vertical"
        >
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="模型名称" name="name">
                <a-input v-model:value="editForm.name" placeholder="请输入模型名称" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="显示名称" name="display_name">
                <a-input v-model:value="editForm.display_name" placeholder="请输入显示名称" />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-form-item label="模型ID" name="model_id">
            <a-input v-model:value="editForm.model_id" placeholder="请输入模型ID" />
          </a-form-item>
          
          <a-form-item label="描述" name="description">
            <a-textarea v-model:value="editForm.description" placeholder="请输入模型描述" :rows="3" />
          </a-form-item>
          
          <a-form-item label="模型能力" name="capabilities">
            <a-select
              v-model:value="editForm.capabilities"
              mode="multiple"
              placeholder="请选择模型能力"
              style="width: 100%"
            >
              <a-select-option value="chat">对话</a-select-option>
              <a-select-option value="completion">文本补全</a-select-option>
              <a-select-option value="embedding">嵌入</a-select-option>
              <a-select-option value="image">图像生成</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="温度范围">
                <a-slider
                  v-model:value="editForm.temperature_range"
                  range
                  :min="0"
                  :max="2"
                  :step="0.1"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="每Token成本" name="cost_per_token">
                <a-input-number
                  v-model:value="editForm.cost_per_token"
                  :min="0"
                  :step="0.0001"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-form-item label="状态" name="status">
            <a-select v-model:value="editForm.status" placeholder="请选择状态">
              <a-select-option value="active">激活</a-select-option>
              <a-select-option value="inactive">未激活</a-select-option>
              <a-select-option value="maintenance">维护中</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-divider>配置信息</a-divider>
          
          <a-form-item label="提供商" name="config.provider">
            <a-select v-model:value="editForm.config.provider" placeholder="请选择提供商">
              <a-select-option value="openai">OpenAI</a-select-option>
              <a-select-option value="azure">Azure OpenAI</a-select-option>
              <a-select-option value="ollama">Ollama</a-select-option>
              <a-select-option value="deepseek">DeepSeek</a-select-option>
              <a-select-option value="doubao">豆包</a-select-option>
              <a-select-option value="qwen">通义千问</a-select-option>
              <a-select-option value="wenxin">文心一言</a-select-option>
              <a-select-option value="mock">Mock</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-form-item label="API密钥" name="config.api_key">
            <a-input-password v-model:value="editForm.config.api_key" placeholder="请输入API密钥" />
          </a-form-item>
          
          <a-form-item label="基础URL" name="config.base_url">
            <a-input v-model:value="editForm.config.base_url" placeholder="请输入基础URL" />
          </a-form-item>
          
          <a-form-item label="模型名称" name="config.model_name">
            <a-input v-model:value="editForm.config.model_name" placeholder="请输入模型名称" />
          </a-form-item>
          
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="温度">
                <a-input-number
                  v-model:value="editForm.config.temperature"
                  :min="0"
                  :max="2"
                  :step="0.1"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="最大Token数">
                <a-input-number
                  v-model:value="editForm.config.max_tokens"
                  :min="1"
                  :max="32000"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="超时时间(秒)">
                <a-input-number
                  v-model:value="editForm.config.timeout"
                  :min="1"
                  :max="300"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </a-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { api, streamGenerate } from '@/utils/api'

const loading = ref(false)
const models = ref([])
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
const editModalVisible = ref(false)
const editLoading = ref(false)
const currentEditModel = ref<any>(null)
const createForm = reactive({
  name: '',
  model_id: '',
  display_name: '',
  description: '',
  capabilities: [],
  temperature_range: [0, 2],
  cost_per_token: 0,
  status: 'active',
  config: {
    provider: '',
    api_key: '',
    base_url: '',
    model_name: '',
    temperature: 0.7,
    max_tokens: 2048,
    timeout: 60
  }
})

const createRules = {
  name: [{ required: true, message: '请输入模型名称' }],
  display_name: [{ required: true, message: '请输入显示名称' }],
  capabilities: [{ required: true, message: '请选择模型能力' }],
  status: [{ required: true, message: '请选择状态' }],
  'config.provider': [{ required: true, message: '请选择提供商' }],
  'config.api_key': [{ required: true, message: '请输入API密钥' }],
  'config.base_url': [{ required: true, message: '请输入基础URL' }],
  'config.model_name': [{ required: true, message: '请输入模型名称' }]
}

const editForm = reactive({
  name: '',
  model_id: '',
  display_name: '',
  description: '',
  capabilities: [],
  temperature_range: [0, 2],
  cost_per_token: 0,
  status: 'active',
  config: {
    provider: '',
    api_key: '',
    base_url: '',
    model_name: '',
    temperature: 0.7,
    max_tokens: 2048,
    timeout: 60
  }
})

const editRules = {
  name: [{ required: true, message: '请输入模型名称' }],
  display_name: [{ required: true, message: '请输入显示名称' }],
  capabilities: [{ required: true, message: '请选择模型能力' }],
  status: [{ required: true, message: '请选择状态' }],
  'config.provider': [{ required: true, message: '请选择提供商' }],
  'config.api_key': [{ required: true, message: '请输入API密钥' }],
  'config.base_url': [{ required: true, message: '请输入基础URL' }],
  'config.model_name': [{ required: true, message: '请输入模型名称' }]
}

const testModalVisible = ref(false)
const testLoading = ref(false)
const currentTestModel = ref<any>(null)
const testForm = reactive({
  message: '你好，请介绍一下自己',
  stream: false
})
const testRules = {
  message: [{ required: true, message: '请输入测试消息' }]
}
const testResult = ref('')
const streamingResult = ref('')
const isStreaming = ref(false)
const testSuccess = ref(false)

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '显示名称',
    dataIndex: 'display_name',
    key: 'display_name',
    width: 150
  },
  {
    title: '提供商',
    dataIndex: ['config', 'provider'],
    key: 'provider',
    width: 100
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: '能力',
    dataIndex: 'capabilities',
    key: 'capabilities',
    width: 200
  },
  {
    title: '使用次数',
    dataIndex: 'usage_count',
    key: 'usage_count',
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

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'green',
    inactive: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    active: '激活',
    inactive: '未激活'
  }
  return texts[status] || status
}

const fetchModels = async (page = 1, pageSize = 10) => {
  try {
    loading.value = true
    const response = await api.get('/llm-models', { params: { page, pageSize } })
    models.value = response.data.items || response.data
    pagination.value.total = response.data.pagination?.total || models.value.length
  } catch (error) {
    message.error('获取模型列表失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchModels(pag.current, pag.pageSize)
}

const showCreateModal = () => {
  createModalVisible.value = true
}

const handleCreate = async () => {
  try {
    createLoading.value = true
    await api.post('/llm-model', createForm)
    createModalVisible.value = false
    fetchModels()
    message.success('创建成功')
    
    // 重置表单
    Object.assign(createForm, {
      name: '',
      model_id: '',
      display_name: '',
      description: '',
      capabilities: [],
      temperature_range: [0, 2],
      cost_per_token: 0,
      status: 'active',
      config: {
        provider: '',
        api_key: '',
        base_url: '',
        model_name: '',
        temperature: 0.7,
        max_tokens: 2048,
        timeout: 60
      }
    })
  } catch (error: any) {
    message.error(error.response?.data?.error || '创建失败')
  } finally {
    createLoading.value = false
  }
}

const handleEdit = async () => {
  if (!currentEditModel.value) {
    message.error('未选择要编辑的模型')
    return
  }
  
  try {
    editLoading.value = true
    
    const response = await api.put(`/llm-model/${currentEditModel.value.id}`, editForm)
    
    message.success('模型更新成功')
    editModalVisible.value = false
    await fetchModels()
  } catch (error: any) {
    message.error(error.response?.data?.error || '更新失败')
  } finally {
    editLoading.value = false
  }
}

const testModel = (model: any) => {
  currentTestModel.value = model
  testModalVisible.value = true
  testForm.message = '你好，请介绍一下自己'
  testForm.stream = false
  testResult.value = ''
  streamingResult.value = ''
  isStreaming.value = false
  testSuccess.value = false
}

const handleTest = async () => {
  if (!currentTestModel.value) {
    message.error('未选择测试模型')
    return
  }
  
  try {
    testLoading.value = true
    testResult.value = ''
    streamingResult.value = ''
    
    if (testForm.stream) {
      // 流式响应处理
      isStreaming.value = true
      await handleStreamTest()
    } else {
      isStreaming.value = false
      // 普通响应处理
      const response = await api.post(`/llm-model/${currentTestModel.value.id}/test`, {
        messages: [{ role: 'user', content: testForm.message }],
        stream: false
      })
      
      testResult.value = response.data.data || response.data.message
      testSuccess.value = true
      message.success('测试完成')
    }
  } catch (error: any) {
    message.error(error.response?.data?.error || '测试失败')
  } finally {
    testLoading.value = false
    // 注意：isStreaming在流式测试的回调中已经设置为false，这里不需要重复设置
    if (!testForm.stream) {
      isStreaming.value = false
    }
  }
}

const handleStreamTest = async () => {
  try {
    // isStreaming已经在handleTest中设置为true
    streamingResult.value = ''
    
    // 使用流式API
    await streamGenerate(
      `/llm-model/${currentTestModel.value.id}/test`,
      {
        messages: [{ role: 'user', content: testForm.message }],
        stream: true
      },
      (content: string) => {
        // 接收到新内容时，追加到结果中
        console.log('接收到流式内容:', content)
        streamingResult.value += content
        console.log('当前streamingResult:', streamingResult.value)
      },
      () => {
        // 流式响应完成
        isStreaming.value = false
        testSuccess.value = true
        message.success('流式测试完成')
      },
      (error: string) => {
        // 流式响应错误
        isStreaming.value = false
        message.error('流式测试失败: ' + error)
      }
    )
  } catch (error: any) {
    console.error('流式测试错误:', error)
    isStreaming.value = false
    message.error('流式测试失败: ' + (error.response?.data?.error || error.message))
  }
}

const editModel = (model: any) => {
  currentEditModel.value = model
  editModalVisible.value = true
  
  // 填充表单数据
  editForm.name = model.name || ''
  editForm.model_id = model.model_id || ''
  editForm.display_name = model.display_name || ''
  editForm.description = model.description || ''
  editForm.capabilities = model.capabilities || []
  editForm.temperature_range = model.temperature_range || [0, 2]
  editForm.cost_per_token = model.cost_per_token || 0
  editForm.status = model.status || 'active'
  
  // 填充配置数据
  if (model.config) {
    editForm.config.provider = model.config.provider || ''
    editForm.config.api_key = model.config.api_key || ''
    editForm.config.base_url = model.config.base_url || ''
    editForm.config.model_name = model.config.model_name || ''
    editForm.config.temperature = model.config.temperature || 0.7
    editForm.config.max_tokens = model.config.max_tokens || 2048
    editForm.config.timeout = model.config.timeout || 60
  }
}

const deleteModel = async (id: string, modelName: string) => {
  // 使用ES6 import的Modal组件
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除模型 "${modelName}" 吗？此操作不可撤销。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await api.delete(`/llm-model/${id}`)
        fetchModels()
        message.success('删除成功')
      } catch (error: any) {
        message.error(error.response?.data?.error || '删除失败')
      }
    },
    onCancel: () => {
      // 取消操作，不显示任何提示
    }
  })
}

onMounted(() => {
  fetchModels()
})
</script>

<style scoped>
.llm-models-page {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}

.test-result {
  margin-top: 16px;
}

.result-content {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 6px;
  max-height: 300px;
  overflow-y: auto;
}

.result-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.model-info {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 6px;
  margin-bottom: 16px;
}

.success-indicator {
  margin-bottom: 12px;
}

.streaming-indicator {
  display: flex;
  align-items: center;
  padding: 8px 0;
  color: #1890ff;
  font-size: 14px;
}

.streaming-content {
  margin-top: 8px;
}

.streaming-content pre {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #d9d9d9;
  min-height: 100px;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.normal-content pre {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #d9d9d9;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
}
</style>
