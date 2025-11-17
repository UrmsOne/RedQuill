<template>
  <div class="dashboard">
      <a-row :gutter="24">
        <a-col :span="24">
          <a-card title="欢迎使用RedQuill" class="welcome-card">
            <p>AI小说生成平台，让创作更简单</p>
            <a-button type="primary" @click="$router.push('/app/novels')">
              开始创作
            </a-button>
          </a-card>
        </a-col>
      </a-row>
      
      <a-row :gutter="24" style="margin-top: 24px">
        <a-col :span="8">
          <a-card title="我的小说" class="stat-card">
            <a-statistic
              title="总数量"
              :value="novelStats.total"
              :loading="loading"
            />
            <template #extra>
              <a-button type="link" @click="$router.push('/app/novels')">
                查看全部
              </a-button>
            </template>
          </a-card>
        </a-col>
        
        <a-col :span="8">
          <a-card title="LLM模型" class="stat-card">
            <a-statistic
              title="可用模型"
              :value="llmStats.total"
              :loading="loading"
            />
            <template #extra>
              <a-button type="link" @click="$router.push('/app/llm-models')">
                管理模型
              </a-button>
            </template>
          </a-card>
        </a-col>
        
        <a-col :span="8">
          <a-card title="Prompt模板" class="stat-card">
            <a-statistic
              title="模板数量"
              :value="promptStats.total"
              :loading="loading"
            />
            <template #extra>
              <a-button type="link" @click="$router.push('/app/prompts')">
                管理模板
              </a-button>
            </template>
          </a-card>
        </a-col>
      </a-row>
      
      <a-row :gutter="24" style="margin-top: 24px">
        <a-col :span="12">
          <a-card title="最近的小说" class="recent-card">
            <a-list
              :data-source="recentNovels"
              :loading="loading"
              item-layout="horizontal"
            >
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta
                    :title="item.title"
                    :description="`状态: ${item.status} | 阶段: ${item.current_phase}`"
                  />
                  <template #actions>
                    <a @click="$router.push(`/app/novel/${item.id}`)">查看</a>
                  </template>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
        
        <a-col :span="12">
          <a-card title="快速操作" class="quick-actions">
            <a-space direction="vertical" size="middle" style="width: 100%">
              <a-button type="primary" block @click="$router.push('/app/novels')">
                创建新小说
              </a-button>
              <a-button block @click="$router.push('/app/llm-models')">
                配置LLM模型
              </a-button>
              <a-button block @click="$router.push('/app/prompts')">
                管理Prompt模板
              </a-button>
            </a-space>
          </a-card>
        </a-col>
      </a-row>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNovelStore } from '@/stores/novel'
import { api } from '@/utils/api'

const novelStore = useNovelStore()

const loading = ref(false)
const novelStats = ref({ total: 0 })
const llmStats = ref({ total: 0 })
const promptStats = ref({ total: 0 })
const recentNovels = ref([])

const fetchStats = async () => {
  try {
    loading.value = true
    
    // 获取小说统计
    const novelsResponse = await api.get('/novels', { params: { pageSize: 1 } })
    novelStats.value.total = novelsResponse.data.pagination?.total || 0
    recentNovels.value = novelsResponse.data.items?.slice(0, 5) || []
    
    // 获取LLM模型统计
    const llmResponse = await api.get('/llm-models', { params: { pageSize: 1 } })
    llmStats.value.total = llmResponse.data.pagination?.total || 0
    
    // 获取Prompt统计
    const promptResponse = await api.get('/prompts', { params: { pageSize: 1 } })
    promptStats.value.total = promptResponse.data.pagination?.total || 0
    
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.dashboard {
  padding: 24px;
}

.welcome-card {
  text-align: center;
  margin-bottom: 24px;
}

.welcome-card p {
  font-size: 16px;
  color: #666;
  margin: 16px 0;
}

.stat-card {
  text-align: center;
}

.recent-card,
.quick-actions {
  height: 300px;
}

.quick-actions .ant-space {
  height: 100%;
  justify-content: center;
}
</style>
