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
            
            <div v-if="!storyCores || storyCores.length === 0" class="empty-state">
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
            
            <div v-if="!worldviews || worldviews.length === 0" class="empty-state">
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
              <div style="margin-top: 16px;">
                <a-button type="link" @click="viewWorldview">查看详情</a-button>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
      
      <a-row :gutter="24" style="margin-top: 24px">
        <a-col :span="8">
          <a-card title="大纲列表" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateOutline">
                生成大纲
              </a-button>
            </template>
            
            <div v-if="!outlines || outlines.length === 0" class="empty-state">
              <a-empty description="暂无大纲" />
            </div>
            
            <div v-else>
              <a-list
                :data-source="outlines"
                item-layout="vertical"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta
                      :title="item.title"
                      :description="item.summary"
                    />
                    <template #actions>
                      <a @click="viewOutline(item)">查看详情</a>
                    </template>
                    <div class="outline-meta">
                      <a-tag color="blue">{{ item.chapters?.length || 0 }}章</a-tag>
                      <a-tag color="green">{{ item.story_arcs?.length || 0 }}个弧线</a-tag>
                    </div>
                  </a-list-item>
                </template>
              </a-list>
            </div>
          </a-card>
        </a-col>
        
        <a-col :span="8">
          <a-card title="角色列表" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateCharacter">
                生成角色
              </a-button>
            </template>
            
            <div v-if="!characters || characters.length === 0" class="empty-state">
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
        
        <a-col :span="8">
          <a-card title="章节列表" class="content-card">
            <template #extra>
              <a-button type="primary" @click="generateChapter">
                生成章节
              </a-button>
            </template>
            
            <div v-if="!chapters || chapters.length === 0" class="empty-state">
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
      
      <!-- 故事核心详情模态框 -->
      <a-modal
        v-model:open="storyCoreModalVisible"
        :title="selectedStoryCore?.title || '故事核心详情'"
        width="800px"
        :footer="null"
      >
        <div v-if="selectedStoryCore" class="detail-content">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="标题">{{ selectedStoryCore.title }}</a-descriptions-item>
            <a-descriptions-item label="核心冲突">{{ selectedStoryCore.core_conflict }}</a-descriptions-item>
            <a-descriptions-item label="主题">{{ selectedStoryCore.theme }}</a-descriptions-item>
            <a-descriptions-item label="创新点">{{ selectedStoryCore.innovation }}</a-descriptions-item>
            <a-descriptions-item label="商业潜力">{{ selectedStoryCore.commercial_potential }}</a-descriptions-item>
            <a-descriptions-item label="目标受众">{{ selectedStoryCore.target_audience }}</a-descriptions-item>
          </a-descriptions>
        </div>
      </a-modal>
      
      <!-- 世界观详情模态框 -->
      <a-modal
        v-model:open="worldviewModalVisible"
        title="世界观详情"
        width="900px"
        :footer="null"
      >
        <div v-if="selectedWorldview" class="detail-content">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="修炼体系">
              <div>
                <p><strong>名称：</strong>{{ selectedWorldview.power_system?.name }}</p>
                <p><strong>等级：</strong>{{ selectedWorldview.power_system?.levels?.join(' → ') }}</p>
                <p><strong>修炼方法：</strong>{{ selectedWorldview.power_system?.cultivation_method }}</p>
                <p><strong>限制：</strong>{{ selectedWorldview.power_system?.limitations }}</p>
              </div>
            </a-descriptions-item>
            <a-descriptions-item label="社会结构">
              <div>
                <p><strong>等级制度：</strong>{{ selectedWorldview.society_structure?.hierarchy }}</p>
                <p><strong>经济体系：</strong>{{ selectedWorldview.society_structure?.economic_system }}</p>
              </div>
            </a-descriptions-item>
            <a-descriptions-item label="地理环境">
              <div>
                <p><strong>主要区域：</strong>{{ selectedWorldview.geography?.major_regions?.join(', ') }}</p>
                <p><strong>特殊地点：</strong>{{ selectedWorldview.geography?.special_locations?.join(', ') }}</p>
              </div>
            </a-descriptions-item>
            <a-descriptions-item label="特殊规则">
              <ul>
                <li v-for="rule in selectedWorldview.special_rules" :key="rule">{{ rule }}</li>
              </ul>
            </a-descriptions-item>
          </a-descriptions>
        </div>
      </a-modal>
      
      <!-- 角色详情模态框 -->
      <a-modal
        v-model:open="characterModalVisible"
        :title="selectedCharacter?.name || '角色详情'"
        width="900px"
        :footer="null"
      >
        <div v-if="selectedCharacter" class="detail-content">
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="名称">{{ selectedCharacter.name }}</a-descriptions-item>
            <a-descriptions-item label="类型">{{ selectedCharacter.type }}</a-descriptions-item>
            <a-descriptions-item label="修炼境界">{{ selectedCharacter.core_attributes?.cultivation_level }}</a-descriptions-item>
            <a-descriptions-item label="能力">
              <a-tag v-for="ability in selectedCharacter.core_attributes?.abilities" :key="ability" color="blue" style="margin-right: 8px;">
                {{ ability }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="核心特质">
              <a-tag v-for="trait in selectedCharacter.soul_profile?.personality?.core_traits" :key="trait" color="purple" style="margin-right: 8px;">
                {{ trait }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="道德观">{{ selectedCharacter.soul_profile?.personality?.moral_compass }}</a-descriptions-item>
            <a-descriptions-item label="核心驱动力">{{ selectedCharacter.soul_profile?.motivations?.core_drive }}</a-descriptions-item>
            <a-descriptions-item label="近期目标">{{ selectedCharacter.soul_profile?.motivations?.immediate_goal }}</a-descriptions-item>
            <a-descriptions-item label="长期目标">{{ selectedCharacter.soul_profile?.motivations?.long_term_goal }}</a-descriptions-item>
          </a-descriptions>
        </div>
      </a-modal>
      
      <!-- 章节详情模态框 -->
      <a-modal
        v-model:open="chapterModalVisible"
        :title="selectedChapter ? `第${selectedChapter.chapter_number}章 ${selectedChapter.title}` : '章节详情'"
        width="1000px"
        :footer="null"
      >
        <div v-if="selectedChapter" class="detail-content">
          <a-descriptions :column="2" bordered style="margin-bottom: 24px;">
            <a-descriptions-item label="章节号">{{ selectedChapter.chapter_number }}</a-descriptions-item>
            <a-descriptions-item label="字数">{{ selectedChapter.word_count }}字</a-descriptions-item>
            <a-descriptions-item label="摘要" :span="2">{{ selectedChapter.summary }}</a-descriptions-item>
            <a-descriptions-item label="章节目标" :span="2" v-if="selectedChapter.outline?.goal">
              {{ selectedChapter.outline.goal }}
            </a-descriptions-item>
            <a-descriptions-item label="关键事件" :span="2" v-if="selectedChapter.outline?.key_events?.length">
              <a-tag v-for="event in selectedChapter.outline.key_events" :key="event" color="blue" style="margin-right: 8px;">
                {{ event }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="角色发展" :span="2" v-if="selectedChapter.character_development && Object.keys(selectedChapter.character_development).length > 0">
              <div v-for="(development, character) in selectedChapter.character_development" :key="character">
                <strong>{{ character }}：</strong>{{ development }}
              </div>
            </a-descriptions-item>
            <a-descriptions-item label="质量评分" v-if="selectedChapter.quality_metrics?.score">
              {{ selectedChapter.quality_metrics.score }}/10
            </a-descriptions-item>
            <a-descriptions-item label="优点" v-if="selectedChapter.quality_metrics?.strengths?.length">
              <a-tag v-for="strength in selectedChapter.quality_metrics.strengths" :key="strength" color="green" style="margin-right: 8px;">
                {{ strength }}
              </a-tag>
            </a-descriptions-item>
          </a-descriptions>
          
          <div v-if="selectedChapter.content" class="chapter-content-section">
            <h4>章节内容</h4>
            <div class="chapter-content-text">{{ selectedChapter.content }}</div>
          </div>
        </div>
      </a-modal>
      
      <!-- 大纲详情模态框 -->
      <a-modal
        v-model:open="outlineModalVisible"
        :title="selectedOutline?.title || '大纲详情'"
        width="1200px"
        :footer="null"
      >
        <OutlineDetail 
          v-if="selectedOutline" 
          :outline="selectedOutline"
          @updated="handleOutlineUpdated"
          @deleted="handleOutlineDeleted"
        />
      </a-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useNovelStore } from '@/stores/novel'
import type { Novel, StoryCore, Worldview, Character, Chapter, Outline } from '@/stores/novel'
import OutlineDetail from '@/components/OutlineDetail.vue'

const route = useRoute()
const router = useRouter()
const novelStore = useNovelStore()

const novel = ref<Novel | null>(null)
const storyCores = ref<StoryCore[]>([])
const worldviews = ref<Worldview[]>([])
const characters = ref<Character[]>([])
const chapters = ref<Chapter[]>([])
const outlines = ref<Outline[]>([])

// 模态框状态
const storyCoreModalVisible = ref(false)
const worldviewModalVisible = ref(false)
const characterModalVisible = ref(false)
const chapterModalVisible = ref(false)
const outlineModalVisible = ref(false)

// 选中的详情数据
const selectedStoryCore = ref<StoryCore | null>(null)
const selectedWorldview = ref<Worldview | null>(null)
const selectedCharacter = ref<Character | null>(null)
const selectedChapter = ref<Chapter | null>(null)
const selectedOutline = ref<Outline | null>(null)

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
    
    await novelStore.fetchOutlines(novelId)
    outlines.value = novelStore.outlines
    
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

const generateOutline = () => {
  router.push(`/app/novel/${novelId}/generate?type=outline`)
}

const generateCharacter = () => {
  router.push(`/app/novel/${novelId}/generate?type=character`)
}

const generateChapter = () => {
  router.push(`/app/novel/${novelId}/generate?type=chapter`)
}

const viewStoryCore = (storyCore: StoryCore) => {
  selectedStoryCore.value = storyCore
  storyCoreModalVisible.value = true
}

const viewOutline = (outline: Outline) => {
  selectedOutline.value = outline
  outlineModalVisible.value = true
}

const viewCharacter = (character: Character) => {
  selectedCharacter.value = character
  characterModalVisible.value = true
}

const viewChapter = (chapter: Chapter) => {
  selectedChapter.value = chapter
  chapterModalVisible.value = true
}

// 处理世界观查看（世界观只有一个）
const viewWorldview = () => {
  if (worldviews.value && worldviews.value.length > 0) {
    selectedWorldview.value = worldviews.value[0]
    worldviewModalVisible.value = true
  } else {
    message.warning('暂无世界观数据')
  }
}

// 处理大纲更新
const handleOutlineUpdated = (updatedOutline: Outline) => {
  const index = outlines.value.findIndex(o => o.id === updatedOutline.id)
  if (index !== -1) {
    outlines.value[index] = updatedOutline
    selectedOutline.value = updatedOutline
  }
  message.success('大纲已更新')
}

// 处理大纲删除
const handleOutlineDeleted = (deletedId: string) => {
  outlines.value = outlines.value.filter(o => o.id !== deletedId)
  outlineModalVisible.value = false
  selectedOutline.value = null
  message.success('大纲已删除')
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

.outline-meta {
  margin-top: 8px;
}

.outline-meta .ant-tag {
  margin-right: 8px;
}

.detail-content {
  max-height: 70vh;
  overflow-y: auto;
}

.chapter-content-section {
  margin-top: 24px;
}

.chapter-content-section h4 {
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 600;
}

.chapter-content-text {
  padding: 16px;
  background: #f5f5f5;
  border-radius: 4px;
  white-space: pre-wrap;
  word-wrap: break-word;
  line-height: 1.8;
  max-height: 500px;
  overflow-y: auto;
}
</style>
