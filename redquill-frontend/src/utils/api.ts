import axios from 'axios'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      message.error('登录已过期，请重新登录')
    } else if (error.response?.status === 403) {
      message.error('没有权限访问')
    } else if (error.response?.status >= 500) {
      message.error('服务器错误，请稍后重试')
    }
    return Promise.reject(error)
  }
)

// 流式请求工具
export const createStreamRequest = (url: string, data: any, onMessage: (content: string) => void, onDone: () => void, onError: (error: string) => void) => {
  const authStore = useAuthStore()
  const token = authStore.token

  const eventSource = new EventSource(`${url}?${new URLSearchParams({
    ...data,
    token: token || ''
  })}`)

  eventSource.onmessage = (event) => {
    try {
      const parsed = JSON.parse(event.data)
      if (parsed.content) {
        onMessage(parsed.content)
      }
      if (parsed.done) {
        onDone()
        eventSource.close()
      }
    } catch (error) {
      console.error('解析流式数据失败:', error)
    }
  }

  eventSource.onerror = (error) => {
    onError('连接错误')
    eventSource.close()
  }

  return eventSource
}

// 流式生成请求
export const streamGenerate = async (
  endpoint: string,
  data: any,
  onMessage: (content: string) => void,
  onDone: () => void,
  onError: (error: string) => void
) => {
  try {
    const authStore = useAuthStore()
    const response = await fetch(`/api/v1${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`,
        'Accept': 'text/event-stream',
        'Cache-Control': 'no-cache'
      },
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const reader = response.body?.getReader()
    if (!reader) {
      throw new Error('无法读取响应流')
    }

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (let i = 0; i < lines.length; i++) {
        const line = lines[i]
        console.log('处理SSE行:', line)
        
        // 跳过空行
        if (line.trim() === '') {
          continue
        }
        
        // 处理SSE事件
        if (line.startsWith('event:')) {
          const eventType = line.slice(6).trim()
          console.log('事件类型:', eventType)
          
          // 查找对应的data行（在当前行之后）
          for (let j = i + 1; j < lines.length; j++) {
            const nextLine = lines[j]
            console.log('查找data行:', nextLine)
            
            if (nextLine.startsWith('data:')) {
              const eventData = nextLine.slice(5).trim()
              
              if (eventType === 'data') {
                try {
                  const parsed = JSON.parse(eventData)
                  console.log('解析到数据事件:', parsed)
                  if (parsed.content !== undefined) {
                    console.log('发送内容到回调:', parsed.content)
                    onMessage(parsed.content)
                  }
                  if (parsed.done === true) {
                    console.log('流式响应完成')
                    onDone()
                    return
                  }
                } catch (error) {
                  console.error('解析数据事件失败:', error)
                  // 如果不是JSON，直接显示文本
                  if (eventData && eventData !== '') {
                    console.log('发送原始数据到回调:', eventData)
                    onMessage(eventData)
                  }
                }
              } else if (eventType === 'error') {
                try {
                  const parsed = JSON.parse(eventData)
                  onError(parsed.error || '未知错误')
                  return
                } catch (error) {
                  onError(eventData || '解析错误信息失败')
                  return
                }
              }
              break
            }
          }
        } else if (line.startsWith('data: ')) {
          // 处理没有事件类型的data行
          const eventData = line.slice(6).trim()
          console.log('处理data行:', eventData)
          
          if (eventData === '[DONE]') {
            onDone()
            return
          }

          try {
            const parsed = JSON.parse(eventData)
            console.log('解析data行JSON:', parsed)
            if (parsed.content !== undefined) {
              onMessage(parsed.content)
            }
            if (parsed.done === true) {
              onDone()
              return
            }
            if (parsed.error) {
              onError(parsed.error)
              return
            }
          } catch (error) {
            console.error('解析事件数据失败:', error)
            // 如果不是JSON，直接显示文本
            if (eventData && eventData !== '') {
              onMessage(eventData)
            }
          }
        }
      }
    }
  } catch (error: any) {
    onError(error.message || '请求失败')
  }
}

export { api }
