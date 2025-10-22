<template>
  <div class="generate-page">
      <a-breadcrumb style="margin-bottom: 16px">
        <a-breadcrumb-item>
          <router-link to="/novels">小说管理</router-link>
        </a-breadcrumb-item>
        <a-breadcrumb-item>
          <router-link :to="`/novel/${novelId}`">{{ novel?.title || '小说' }}</router-link>
        </a-breadcrumb-item>
        <a-breadcrumb-item>AI生成</a-breadcrumb-item>
      </a-breadcrumb>
      
      <a-card title="AI内容生成" class="content-card">
        <a-tabs v-model:activeKey="activeTab" @change="handleTabChange">
          <a-tab-pane key="story-core" tab="故事核心">
            <StoryCoreGenerate 
              :novel="novel"
              @generated="handleGenerated"
            />
          </a-tab-pane>
          
          <a-tab-pane key="worldview" tab="世界观">
            <WorldviewGenerate 
              :novel="novel"
              @generated="handleGenerated"
            />
          </a-tab-pane>
          
          <a-tab-pane key="character" tab="角色">
            <CharacterGenerate 
              :novel="novel"
              @generated="handleGenerated"
            />
          </a-tab-pane>
          
          <a-tab-pane key="chapter" tab="章节">
            <ChapterGenerate 
              :novel="novel"
              @generated="handleGenerated"
            />
          </a-tab-pane>
        </a-tabs>
      </a-card>
      
      <!-- AI生成信息 -->
      <a-card title="AI生成信息" class="content-card" style="margin-top: 16px;">
        <ExtraInfoDisplay :extraInfo="extraInfo" />
      </a-card>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useNovelData } from '@/composables/useNovelData'
import StoryCoreGenerate from '@/components/generate/StoryCoreGenerate.vue'
import WorldviewGenerate from '@/components/generate/WorldviewGenerate.vue'
import CharacterGenerate from '@/components/generate/CharacterGenerate.vue'
import ChapterGenerate from '@/components/generate/ChapterGenerate.vue'
import ExtraInfoDisplay from '@/components/ExtraInfoDisplay.vue'

const route = useRoute()
const novelId = route.params.id as string
const activeTab = ref('story-core')

// 使用小说数据composable
const {
  novel,
  storyCores,
  worldviews,
  characters,
  chapters,
  loading,
  storyCoresLoading,
  worldviewsLoading,
  charactersLoading,
  chaptersLoading,
  novelInfo,
  storyCoreInfo,
  worldviewInfo,
  characterInfo,
  chapterInfo,
  fetchNovel,
  fetchStoryCores,
  fetchWorldviews,
  fetchCharacters,
  fetchChapters,
  fetchAllData,
  fetchRelevantData,
  getFormDefaults
} = useNovelData(novelId)

const handleTabChange = (key: string) => {
  activeTab.value = key
  console.log('切换标签页:', key, '当前小说数据:', novel.value)
  // 切换标签页时自动获取相关数据
  fetchRelevantData()
}

const handleGenerated = async (type: string, data: any) => {
  console.log(`${type} 生成完成:`, data)
  
  // 根据生成类型刷新相关数据
  const refreshPromises = []
  
  switch (type) {
    case 'story-core':
      refreshPromises.push(fetchStoryCores())
      break
    case 'worldview':
      refreshPromises.push(fetchWorldviews())
      break
    case 'character':
      refreshPromises.push(fetchCharacters())
      break
    case 'chapter':
      refreshPromises.push(fetchChapters())
      break
  }
  
  // 同时刷新小说基本信息
  refreshPromises.push(fetchNovel())
  
  await Promise.all(refreshPromises)
}

onMounted(async () => {
  // 初始加载时获取所有相关数据
  await fetchRelevantData()
  
  // 根据URL参数设置默认标签页
  const type = route.query.type as string
  if (type) {
    activeTab.value = type
  }
})
</script>

<style scoped>
.generate-page {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}
</style>
