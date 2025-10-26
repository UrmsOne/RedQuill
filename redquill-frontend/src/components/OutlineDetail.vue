<template>
  <div class="outline-detail">
    <a-card :title="outline?.title || '大纲详情'" class="outline-card">
      <template #extra>
        <a-space>
          <a-button @click="handleEdit">编辑</a-button>
          <a-button type="primary" @click="handleExport">导出</a-button>
          <a-button danger @click="handleDelete">删除</a-button>
        </a-space>
      </template>
      
      <a-tabs v-model:activeKey="activeTab" @change="handleTabChange">
        <a-tab-pane key="overview" tab="概览">
          <div class="overview-content">
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="标题" :span="2">
                {{ outline?.title || '未命名' }}
              </a-descriptions-item>
              <a-descriptions-item label="概要" :span="2">
                {{ outline?.summary || '暂无概要' }}
              </a-descriptions-item>
              <a-descriptions-item label="总章节数">
                {{ outline?.chapters?.length || 0 }} 章
              </a-descriptions-item>
              <a-descriptions-item label="故事弧线数">
                {{ outline?.story_arcs?.length || 0 }} 个
              </a-descriptions-item>
              <a-descriptions-item label="关键主题" :span="2">
                <a-tag v-for="theme in outline?.key_themes" :key="theme" color="blue">
                  {{ theme }}
                </a-tag>
                <span v-if="!outline?.key_themes?.length">暂无主题</span>
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ formatTime(outline?.ctime) }}
              </a-descriptions-item>
              <a-descriptions-item label="更新时间">
                {{ formatTime(outline?.mtime) }}
              </a-descriptions-item>
            </a-descriptions>
          </div>
        </a-tab-pane>
        
        <a-tab-pane key="story-arcs" tab="故事弧线">
          <div v-if="outline?.story_arcs?.length" class="story-arcs-content">
            <a-timeline>
              <a-timeline-item 
                v-for="arc in outline.story_arcs" 
                :key="arc.name"
                :color="getArcColor(arc)"
              >
                <template #dot>
                  <a-badge :count="`${arc.start_chapter}-${arc.end_chapter}`" />
                </template>
                <a-card size="small" :title="arc.name">
                  <p><strong>描述：</strong>{{ arc.description }}</p>
                  <p><strong>主题：</strong>{{ arc.theme }}</p>
                  <p><strong>章节范围：</strong>第{{ arc.start_chapter }}章 - 第{{ arc.end_chapter }}章</p>
                </a-card>
              </a-timeline-item>
            </a-timeline>
          </div>
          <a-empty v-else description="暂无故事弧线" />
        </a-tab-pane>
        
        <a-tab-pane key="chapters" tab="章节列表">
          <div v-if="outline?.chapters?.length" class="chapters-content">
            <a-list
              :data-source="outline.chapters"
              item-layout="vertical"
              :pagination="{ 
                pageSize: 10, 
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
              }"
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
                    <a-row :gutter="16">
                      <a-col :span="12">
                        <div v-if="item.key_events?.length">
                          <strong>关键事件：</strong>
                          <ul>
                            <li v-for="event in item.key_events" :key="event">{{ event }}</li>
                          </ul>
                        </div>
                      </a-col>
                      <a-col :span="12">
                        <div v-if="item.characters?.length">
                          <strong>出场角色：</strong>
                          <div style="margin-top: 4px;">
                            <a-tag v-for="character in item.characters" :key="character" color="cyan" size="small">
                              {{ character }}
                            </a-tag>
                          </div>
                        </div>
                      </a-col>
                    </a-row>
                    
                    <div v-if="item.outline?.goal" class="chapter-goal">
                      <strong>章节目标：</strong>{{ item.outline.goal }}
                    </div>
                  </div>
                </a-list-item>
              </template>
            </a-list>
          </div>
          <a-empty v-else description="暂无章节信息" />
        </a-tab-pane>
        
        <a-tab-pane key="timeline" tab="时间线">
          <div v-if="outline?.chapters?.length" class="timeline-content">
            <a-timeline mode="left">
              <a-timeline-item 
                v-for="chapter in outline.chapters.slice(0, 20)" 
                :key="chapter.chapter_number"
                :color="getChapterColor(chapter)"
              >
                <template #dot>
                  <a-badge :count="chapter.chapter_number" />
                </template>
                <a-card size="small" :title="chapter.title">
                  <p>{{ chapter.summary }}</p>
                  <div class="chapter-meta">
                    <a-tag color="blue" size="small">{{ chapter.word_count }}字</a-tag>
                    <a-tag color="purple" size="small">{{ chapter.pov }}</a-tag>
                    <a-tag color="orange" size="small">{{ chapter.location }}</a-tag>
                  </div>
                </a-card>
              </a-timeline-item>
            </a-timeline>
            <div v-if="outline.chapters.length > 20" class="timeline-more">
              <a-alert 
                message="显示前20章，完整时间线请查看章节列表" 
                type="info" 
                show-icon 
              />
            </div>
          </div>
          <a-empty v-else description="暂无章节信息" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
    
    <!-- 编辑模态框 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑大纲"
      width="800px"
      @ok="handleSaveEdit"
      @cancel="handleCancelEdit"
    >
      <a-form
        :model="editForm"
        layout="vertical"
        ref="editFormRef"
      >
        <a-form-item label="标题" name="title">
          <a-input v-model:value="editForm.title" placeholder="请输入大纲标题" />
        </a-form-item>
        
        <a-form-item label="概要" name="summary">
          <a-textarea 
            v-model:value="editForm.summary" 
            placeholder="请输入大纲概要"
            :rows="4"
          />
        </a-form-item>
        
        <a-form-item label="关键主题" name="key_themes">
          <a-select
            v-model:value="editForm.key_themes"
            mode="tags"
            placeholder="请输入关键主题"
            style="width: 100%"
          >
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useNovelStore } from '@/stores/novel'
import type { Outline } from '@/stores/novel'

const props = defineProps<{
  outline: Outline | null
}>()

const emit = defineEmits<{
  updated: [outline: Outline]
  deleted: [id: string]
}>()

const novelStore = useNovelStore()

const activeTab = ref('overview')
const editModalVisible = ref(false)
const editFormRef = ref()

const editForm = reactive({
  title: '',
  summary: '',
  key_themes: [] as string[]
})

// 格式化时间
const formatTime = (timestamp?: number) => {
  if (!timestamp) return '未知'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 获取故事弧线颜色
const getArcColor = (arc: any) => {
  const colors = ['blue', 'green', 'red', 'orange', 'purple', 'cyan']
  return colors[arc.start_chapter % colors.length]
}

// 获取章节颜色
const getChapterColor = (chapter: any) => {
  const colors = ['blue', 'green', 'red', 'orange', 'purple', 'cyan']
  return colors[chapter.chapter_number % colors.length]
}

// 处理标签页切换
const handleTabChange = (key: string) => {
  activeTab.value = key
}

// 处理编辑
const handleEdit = () => {
  if (!props.outline) return
  
  editForm.title = props.outline.title || ''
  editForm.summary = props.outline.summary || ''
  editForm.key_themes = [...(props.outline.key_themes || [])]
  editModalVisible.value = true
}

// 处理保存编辑
const handleSaveEdit = async () => {
  if (!props.outline) return
  
  try {
    const updatedOutline = await novelStore.updateOutline(props.outline.id, {
      title: editForm.title,
      summary: editForm.summary,
      key_themes: editForm.key_themes
    })
    
    if (updatedOutline) {
      emit('updated', updatedOutline)
      editModalVisible.value = false
    }
  } catch (error) {
    console.error('更新大纲失败:', error)
  }
}

// 处理取消编辑
const handleCancelEdit = () => {
  editModalVisible.value = false
}

// 处理删除
const handleDelete = () => {
  if (!props.outline) return
  
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这个大纲吗？此操作不可恢复。',
    okText: '确认删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        const success = await novelStore.deleteOutline(props.outline!.id)
        if (success) {
          emit('deleted', props.outline!.id)
        }
      } catch (error) {
        console.error('删除大纲失败:', error)
      }
    }
  })
}

// 处理导出
const handleExport = () => {
  if (!props.outline) return
  
  const exportData = {
    title: props.outline.title,
    summary: props.outline.summary,
    key_themes: props.outline.key_themes,
    story_arcs: props.outline.story_arcs,
    chapters: props.outline.chapters,
    export_time: new Date().toISOString()
  }
  
  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.outline.title || '大纲'}_${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  message.success('导出成功')
}

// 监听outline变化
watch(() => props.outline, (newOutline) => {
  if (newOutline) {
    activeTab.value = 'overview'
  }
}, { immediate: true })
</script>

<style scoped>
.outline-detail {
  max-width: 1200px;
}

.outline-card {
  margin-bottom: 24px;
}

.overview-content {
  padding: 16px 0;
}

.story-arcs-content {
  padding: 16px 0;
}

.chapters-content {
  padding: 16px 0;
}

.timeline-content {
  padding: 16px 0;
}

.timeline-more {
  margin-top: 16px;
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

.chapter-goal {
  margin-top: 8px;
  padding: 8px;
  background: #f0f8ff;
  border-radius: 4px;
  border-left: 3px solid #1890ff;
}

.chapter-meta {
  margin-top: 8px;
}

.chapter-meta .ant-tag {
  margin-right: 4px;
}
</style>
