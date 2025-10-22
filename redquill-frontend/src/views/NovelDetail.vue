<template>
  <div class="novel-detail">
      <a-breadcrumb style="margin-bottom: 16px">
        <a-breadcrumb-item>
          <router-link to="/app/novels">小说管理</router-link>
        </a-breadcrumb-item>
        <a-breadcrumb-item>{{ novel?.title || '小说详情' }}</a-breadcrumb-item>
      </a-breadcrumb>
      
      <a-card v-if="novel" title="小说信息" class="content-card">
        <template #extra>
          <a-space>
            <a-button @click="editNovel">编辑</a-button>
            <a-button type="primary" @click="startGenerate">开始生成</a-button>
          </a-space>
        </template>
        
        <a-descriptions :column="2">
          <a-descriptions-item label="标题">{{ novel.title }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(novel.status)">
              {{ getStatusText(novel.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="当前阶段">
            <a-tag :color="getPhaseColor(novel.current_phase)">
              {{ getPhaseText(novel.current_phase) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="类型">{{ novel.project_blueprint.genre }}</a-descriptions-item>
          <a-descriptions-item label="子类型">{{ novel.project_blueprint.sub_genre }}</a-descriptions-item>
          <a-descriptions-item label="目标受众">{{ novel.project_blueprint.target_audience }}</a-descriptions-item>
          <a-descriptions-item label="核心冲突" :span="2">
            {{ novel.project_blueprint.core_conflict }}
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
      
      <a-row :gutter="24" style="margin-top: 24px">
        <a-col :span="12">
          <a-card title="故事核心" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateStoryCore">
                生成故事核心
              </a-button>
            </template>
            
            <div v-if="storyCores.length === 0" class="empty-state">
              <a-empty description="暂无故事核心" />
            </div>
            
            <div v-else>
              <a-list
                :data-source="storyCores"
                item-layout="vertical"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="item.title"
                      :description="item.core_conflict"
                    />
                    <template #actions>
                      <a @click="viewStoryCore(item)">查看详情</a>
                    </template>
                  </a-list-item>
                </template>
              </a-list>
            </div>
          </a-card>
        </a-col>
        
        <a-col :span="12">
          <a-card title="世界观" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateWorldview">
                生成世界观
              </a-button>
            </template>
            
            <div v-if="worldviews.length === 0" class="empty-state">
              <a-empty description="暂无世界观" />
            </div>
            
            <div v-else>
              <a-descriptions v-for="worldview in worldviews" :key="worldview.id" :column="1">
                <a-descriptions-item label="修炼体系">
                  {{ worldview.power_system.name }}
                </a-descriptions-item>
                <a-descriptions-item label="社会结构">
                  {{ worldview.society_structure.hierarchy }}
                </a-descriptions-item>
                <a-descriptions-item label="地理环境">
                  {{ worldview.geography.major_regions.join(', ') }}
                </a-descriptions-item>
              </a-descriptions>
            </div>
          </a-card>
        </a-col>
      </a-row>
      
      <a-row :gutter="24" style="margin-top: 24px">
        <a-col :span="12">
          <a-card title="角色列表" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateCharacter">
                生成角色
              </a-button>
            </template>
            
            <div v-if="characters.length === 0" class="empty-state">
              <a-empty description="暂无角色" />
            </div>
            
            <div v-else>
              <a-list
                :data-source="characters"
                item-layout="horizontal"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="item.name"
                      :description="`${item.type} | ${item.core_attributes.cultivation_level}`"
                    />
                    <template #actions>
                      <a @click="viewCharacter(item)">查看详情</a>
                    </template>
                  </a-list-item>
                </template>
              </a-list>
            </div>
          </a-card>
        </a-col>
        
        <a-col :span="12">
          <a-card title="章节列表" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateChapter">
                生成章节
              </a-button>
            </template>
            
            <div v-if="chapters.length === 0" class="empty-state">
              <a-empty description="暂无章节" />
            </div>
            
            <div v-else>
              <a-list
                :data-source="chapters"
                item-layout="horizontal"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="`第${item.chapter_number}章 ${item.title}`"
                      :description="`${item.word_count}字 | ${item.summary}`"
                    />
                    <template #actions>
                      <a @click="viewChapter(item)">查看内容</a>
                    </template>
                  </a-list-item>
                </template>
              </a-list>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useNovelStore } from '@/stores/novel'
import type { Novel, StoryCore, Worldview, Character, Chapter } from '@/stores/novel'

const route = useRoute()
const router = useRouter()
const novelStore = useNovelStore()

const novel = ref<Novel | null>(null)
const storyCores = ref<StoryCore[]>([])
const worldviews = ref<Worldview[]>([])
const characters = ref<Character[]>([])
const chapters = ref<Chapter[]>([])

const novelId = route.params.id as string

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    drafting: 'blue',
    writing: 'green',
    completed: 'success',
    paused: 'orange'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    drafting: '草稿',
    writing: '写作中',
    completed: '已完成',
    paused: '暂停'
  }
  return texts[status] || status
}

const getPhaseColor = (phase: string) => {
  const colors: Record<string, string> = {
    story_core: 'purple',
    worldview: 'cyan',
    characters: 'magenta',
    outlining: 'blue',
    writing: 'green'
  }
  return colors[phase] || 'default'
}

const getPhaseText = (phase: string) => {
  const texts: Record<string, string> = {
    story_core: '故事核心',
    worldview: '世界观',
    characters: '角色',
    outlining: '大纲',
    writing: '写作'
  }
  return texts[phase] || phase
}

const fetchNovelData = async () => {
  try {
    // 获取小说详情
    await novelStore.fetchNovel(novelId)
    novel.value = novelStore.currentNovel
    
    // 获取相关数据
    await novelStore.fetchStoryCores(novelId)
    storyCores.value = novelStore.storyCores
    
    await novelStore.fetchWorldview(novelId)
    worldviews.value = novelStore.worldviews
    
    await novelStore.fetchCharacters(novelId)
    characters.value = novelStore.characters
    
    await novelStore.fetchChapters(novelId)
    chapters.value = novelStore.chapters
    
  } catch (error) {
    message.error('获取小说数据失败')
  }
}

const editNovel = () => {
  message.info('编辑功能开发中')
}

const startGenerate = () => {
  router.push(`/app/novel/${novelId}/generate`)
}

const generateStoryCore = () => {
  router.push(`/app/novel/${novelId}/generate?type=story-core`)
}

const generateWorldview = () => {
  router.push(`/app/novel/${novelId}/generate?type=worldview`)
}

const generateCharacter = () => {
  router.push(`/app/novel/${novelId}/generate?type=character`)
}

const generateChapter = () => {
  router.push(`/app/novel/${novelId}/generate?type=chapter`)
}

const viewStoryCore = (storyCore: StoryCore) => {
  message.info('查看故事核心详情')
}

const viewCharacter = (character: Character) => {
  message.info('查看角色详情')
}

const viewChapter = (chapter: Chapter) => {
  message.info('查看章节内容')
}

onMounted(() => {
  fetchNovelData()
})
</script>

<style scoped>
.novel-detail {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}
</style>
