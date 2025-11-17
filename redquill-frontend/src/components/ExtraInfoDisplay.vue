<template>
  <div class="extra-info-display">
    <a-card title="AI生成信息" size="small" v-if="extraInfo">
      <a-tabs v-model:activeKey="activeTab" size="small">
        <a-tab-pane key="story_core" tab="故事核心" v-if="extraInfo.story_core">
          <div class="info-section">
            <a-descriptions size="small" :column="2">
              <a-descriptions-item label="生成时间">
                {{ formatTime(extraInfo.story_core.generation_time) }}
              </a-descriptions-item>
              <a-descriptions-item label="Token消耗">
                {{ extraInfo.story_core.token_count }}
              </a-descriptions-item>
              <a-descriptions-item label="使用次数">
                {{ extraInfo.story_core.usage_count }}
              </a-descriptions-item>
            </a-descriptions>
            <a-divider />
            <a-collapse size="small">
              <a-collapse-panel key="raw" header="原始响应">
                <pre class="raw-response">{{ JSON.stringify(extraInfo.story_core.raw_response, null, 2) }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </a-tab-pane>
        
        <a-tab-pane key="worldview" tab="世界观" v-if="extraInfo.worldview">
          <div class="info-section">
            <a-descriptions size="small" :column="2">
              <a-descriptions-item label="生成时间">
                {{ formatTime(extraInfo.worldview.generation_time) }}
              </a-descriptions-item>
              <a-descriptions-item label="Token消耗">
                {{ extraInfo.worldview.token_count }}
              </a-descriptions-item>
              <a-descriptions-item label="使用次数">
                {{ extraInfo.worldview.usage_count }}
              </a-descriptions-item>
            </a-descriptions>
            <a-divider />
            <a-collapse size="small">
              <a-collapse-panel key="raw" header="原始响应">
                <pre class="raw-response">{{ JSON.stringify(extraInfo.worldview.raw_response, null, 2) }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </a-tab-pane>
        
        <a-tab-pane key="character" tab="角色" v-if="extraInfo.character">
          <div class="info-section">
            <a-descriptions size="small" :column="2">
              <a-descriptions-item label="生成时间">
                {{ formatTime(extraInfo.character.generation_time) }}
              </a-descriptions-item>
              <a-descriptions-item label="Token消耗">
                {{ extraInfo.character.token_count }}
              </a-descriptions-item>
              <a-descriptions-item label="使用次数">
                {{ extraInfo.character.usage_count }}
              </a-descriptions-item>
            </a-descriptions>
            <a-divider />
            <a-collapse size="small">
              <a-collapse-panel key="raw" header="原始响应">
                <pre class="raw-response">{{ JSON.stringify(extraInfo.character.raw_response, null, 2) }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </a-tab-pane>
        
        <a-tab-pane key="chapter" tab="章节" v-if="extraInfo.chapter">
          <div class="info-section">
            <a-descriptions size="small" :column="2">
              <a-descriptions-item label="生成时间">
                {{ formatTime(extraInfo.chapter.generation_time) }}
              </a-descriptions-item>
              <a-descriptions-item label="Token消耗">
                {{ extraInfo.chapter.token_count }}
              </a-descriptions-item>
              <a-descriptions-item label="使用次数">
                {{ extraInfo.chapter.usage_count }}
              </a-descriptions-item>
            </a-descriptions>
            <a-divider />
            <a-collapse size="small">
              <a-collapse-panel key="raw" header="原始响应">
                <pre class="raw-response">{{ JSON.stringify(extraInfo.chapter.raw_response, null, 2) }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </a-tab-pane>
        
        <a-tab-pane key="all" tab="全部信息" v-if="extraInfo.all_phases">
          <div class="info-section">
            <a-collapse size="small">
              <a-collapse-panel 
                v-for="(phase, key) in extraInfo.all_phases" 
                :key="key" 
                :header="getPhaseTitle(key)"
              >
                <a-descriptions size="small" :column="2">
                  <a-descriptions-item label="生成时间">
                    {{ formatTime(phase.generation_time) }}
                  </a-descriptions-item>
                  <a-descriptions-item label="Token消耗">
                    {{ phase.token_count }}
                  </a-descriptions-item>
                  <a-descriptions-item label="使用次数">
                    {{ phase.usage_count }}
                  </a-descriptions-item>
                </a-descriptions>
                <a-divider />
                <pre class="raw-response">{{ JSON.stringify(phase.raw_response, null, 2) }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    
    <a-empty v-else description="暂无AI生成信息" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  extraInfo: any
}>()

const activeTab = computed(() => {
  if (!props.extraInfo) return 'all'
  
  // 根据可用的阶段设置默认标签页
  if (props.extraInfo.story_core) return 'story_core'
  if (props.extraInfo.worldview) return 'worldview'
  if (props.extraInfo.character) return 'character'
  if (props.extraInfo.chapter) return 'chapter'
  return 'all'
})

const formatTime = (timestamp: any) => {
  if (!timestamp) return '未知'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

const getPhaseTitle = (key: string) => {
  const titles: Record<string, string> = {
    story_core: '故事核心',
    worldview: '世界观',
    character: '角色',
    chapter: '章节'
  }
  return titles[key] || key
}
</script>

<style scoped>
.extra-info-display {
  margin-top: 16px;
}

.info-section {
  padding: 8px 0;
}

.raw-response {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 300px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
