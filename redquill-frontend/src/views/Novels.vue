<template>
  <div class="novels-page">
      <a-card title="小说管理" class="content-card">
        <template #extra>
          <a-button type="primary" @click="showCreateModal">
            <PlusOutlined />
            创建小说
          </a-button>
        </template>
        
        <a-table
          :columns="columns"
          :data-source="novels"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column?.key === 'status' && record">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            
            <template v-if="column?.key === 'current_phase' && record">
              <a-tag :color="getPhaseColor(record.current_phase)">
                {{ getPhaseText(record.current_phase) }}
              </a-tag>
            </template>
            
            <template v-if="column?.key === 'actions' && record">
              <a-space>
                <a @click="viewNovel(record)">查看</a>
                <a @click="editNovel(record)">编辑</a>
                <a @click="generateContent(record)">生成</a>
                <a-popconfirm
                  title="确定要删除这本小说吗？"
                  @confirm="deleteNovel(record.id)"
                >
                  <a style="color: #ff4d4f">删除</a>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
      
      <!-- 创建小说模态框 -->
      <a-modal
        v-model:open="createModalVisible"
        title="创建小说"
        @ok="handleCreate"
        :confirm-loading="createLoading"
      >
        <a-form
          :model="createForm"
          :rules="createRules"
          layout="vertical"
        >
          <a-form-item label="小说标题" name="title">
            <a-input v-model:value="createForm.title" placeholder="请输入小说标题" />
          </a-form-item>
          
          <a-form-item label="小说状态" name="status">
            <a-select v-model:value="createForm.status" placeholder="请选择状态">
              <a-select-option value="drafting">草稿</a-select-option>
              <a-select-option value="writing">写作中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="paused">暂停</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-form-item label="当前阶段" name="current_phase">
            <a-select v-model:value="createForm.current_phase" placeholder="请选择阶段">
              <a-select-option value="story_core">故事核心</a-select-option>
              <a-select-option value="worldview">世界观</a-select-option>
              <a-select-option value="characters">角色</a-select-option>
              <a-select-option value="outlining">大纲</a-select-option>
              <a-select-option value="writing">写作</a-select-option>
            </a-select>
          </a-form-item>
          
          <a-form-item label="类型" name="genre">
            <a-input v-model:value="createForm.project_blueprint.genre" placeholder="如：玄幻、都市、历史等" />
          </a-form-item>
          
          <a-form-item label="子类型" name="sub_genre">
            <a-input v-model:value="createForm.project_blueprint.sub_genre" placeholder="如：重生、系统、修仙等" />
          </a-form-item>
          
          <a-form-item label="目标受众" name="target_audience">
            <a-input v-model:value="createForm.project_blueprint.target_audience" placeholder="目标读者群体" />
          </a-form-item>
          
          <a-form-item label="核心冲突" name="core_conflict">
            <a-textarea v-model:value="createForm.project_blueprint.core_conflict" placeholder="故事的核心冲突" />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useNovelStore } from '@/stores/novel'
import type { Novel } from '@/stores/novel'

const router = useRouter()
const novelStore = useNovelStore()

const loading = ref(false)
const novels = ref<Novel[]>([])
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
  title: '',
  status: 'drafting',
  current_phase: 'story_core',
  project_blueprint: {
    genre: '',
    sub_genre: '',
    total_chapters: 100,
    core_conflict: '',
    target_audience: '',
    commercial_focus: ''
  },
  ai_context: {
    recent_summary: '',
    current_focus: '',
    style_guideline: '',
    emotional_tone: ''
  }
})

const createRules = {
  title: [{ required: true, message: '请输入小说标题' }],
  status: [{ required: true, message: '请选择小说状态' }],
  current_phase: [{ required: true, message: '请选择当前阶段' }]
}

const columns = [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    width: 200
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: '阶段',
    dataIndex: 'current_phase',
    key: 'current_phase',
    width: 120
  },
  {
    title: '类型',
    dataIndex: ['project_blueprint', 'genre'],
    key: 'genre',
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
    width: 200
  }
]

const getStatusColor = (status: string | null | undefined) => {
  if (!status) return 'default'
  const colors: Record<string, string> = {
    drafting: 'blue',
    writing: 'green',
    completed: 'success',
    paused: 'orange'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string | null | undefined) => {
  if (!status) return '未知'
  const texts: Record<string, string> = {
    drafting: '草稿',
    writing: '写作中',
    completed: '已完成',
    paused: '暂停'
  }
  return texts[status] || status
}

const getPhaseColor = (phase: string | null | undefined) => {
  if (!phase) return 'default'
  const colors: Record<string, string> = {
    story_core: 'purple',
    worldview: 'cyan',
    characters: 'magenta',
    outlining: 'blue',
    writing: 'green'
  }
  return colors[phase] || 'default'
}

const getPhaseText = (phase: string | null | undefined) => {
  if (!phase) return '未知'
  const texts: Record<string, string> = {
    story_core: '故事核心',
    worldview: '世界观',
    characters: '角色',
    outlining: '大纲',
    writing: '写作'
  }
  return texts[phase] || phase
}

const fetchNovels = async (page = 1, pageSize = 10) => {
  try {
    loading.value = true
    console.log('开始获取小说列表...')
    await novelStore.fetchNovels({ page, pageSize })
    novels.value = novelStore.novels || []
    // 修复分页总数计算
    pagination.value.total = novelStore.novels?.length || 0
    console.log('小说列表获取完成:', novels.value.length, '条记录')
  } catch (error) {
    console.error('获取小说列表失败:', error)
    message.error('获取小说列表失败')
    novels.value = []
    pagination.value.total = 0
  } finally {
    loading.value = false
    console.log('loading状态已重置:', loading.value)
  }
}

const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchNovels(pag.current, pag.pageSize)
}

const showCreateModal = () => {
  createModalVisible.value = true
}

const handleCreate = async () => {
  try {
    createLoading.value = true
    const novel = await novelStore.createNovel(createForm)
    if (novel) {
      createModalVisible.value = false
      fetchNovels()
      // 重置表单
      Object.assign(createForm, {
        title: '',
        status: 'drafting',
        current_phase: 'story_core',
        project_blueprint: {
          genre: '',
          sub_genre: '',
          total_chapters: 100,
          core_conflict: '',
          target_audience: '',
          commercial_focus: ''
        },
        ai_context: {
          recent_summary: '',
          current_focus: '',
          style_guideline: '',
          emotional_tone: ''
        }
      })
    }
  } catch (error) {
    message.error('创建失败')
  } finally {
    createLoading.value = false
  }
}

const viewNovel = (novel: Novel) => {
  router.push(`/app/novel/${novel.id}`)
}

const editNovel = (novel: Novel) => {
  // TODO: 实现编辑功能
  message.info('编辑功能开发中')
}

const generateContent = (novel: Novel) => {
  router.push(`/app/novel/${novel.id}/generate`)
}

const deleteNovel = async (id: string) => {
  try {
    await novelStore.deleteNovel(id)
    fetchNovels()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(() => {
  fetchNovels()
})
</script>

<style scoped>
.novels-page {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}
</style>
