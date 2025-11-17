<template>
  <div class="users-page">
      <a-card title="用户管理" class="content-card">
        <a-table
          :columns="columns"
          :data-source="users"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'actions'">
              <a-space>
                <a @click="editUser(record)">编辑</a>
                <a @click="deleteUser(record.id)">删除</a>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>
      
      <!-- 编辑用户模态框 -->
      <a-modal
        v-model:open="editModalVisible"
        title="编辑用户"
        @ok="handleEdit"
        :confirm-loading="editLoading"
      >
        <a-form
          :model="editForm"
          :rules="editRules"
          layout="vertical"
        >
          <a-form-item label="用户名" name="name">
            <a-input v-model:value="editForm.name" placeholder="请输入用户名" />
          </a-form-item>
          
          <a-form-item label="邮箱" name="email">
            <a-input v-model:value="editForm.email" placeholder="请输入邮箱" />
          </a-form-item>
          
          <a-form-item label="新密码" name="password">
            <a-input-password v-model:value="editForm.password" placeholder="留空表示不修改密码" />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { api } from '@/utils/api'

const loading = ref(false)
const users = ref([])
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`
})

const editModalVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive({
  id: '',
  name: '',
  email: '',
  password: ''
})

const editRules = {
  name: [{ required: true, message: '请输入用户名' }],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' }
  ]
}

const columns = [
  {
    title: '用户名',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
    width: 200
  },
  {
    title: '创建时间',
    dataIndex: 'ctime',
    key: 'ctime',
    width: 180,
    customRender: ({ text }: any) => new Date(text * 1000).toLocaleString()
  },
  {
    title: '修改时间',
    dataIndex: 'mtime',
    key: 'mtime',
    width: 180,
    customRender: ({ text }: any) => new Date(text * 1000).toLocaleString()
  },
  {
    title: '操作',
    key: 'actions',
    width: 150
  }
]

const fetchUsers = async (page = 1, pageSize = 10) => {
  try {
    loading.value = true
    const response = await api.get('/users', { params: { page, pageSize } })
    users.value = response.data.items || response.data
    pagination.value.total = response.data.pagination?.total || users.value.length
  } catch (error) {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchUsers(pag.current, pag.pageSize)
}

const editUser = (user: any) => {
  editForm.id = user.id
  editForm.name = user.name
  editForm.email = user.email
  editForm.password = ''
  editModalVisible.value = true
}

const handleEdit = async () => {
  try {
    editLoading.value = true
    const updateData: any = {
      name: editForm.name,
      email: editForm.email
    }
    
    if (editForm.password) {
      updateData.password = editForm.password
    }
    
    await api.put(`/user/${editForm.id}`, updateData)
    editModalVisible.value = false
    fetchUsers()
    message.success('更新成功')
  } catch (error: any) {
    message.error(error.response?.data?.error || '更新失败')
  } finally {
    editLoading.value = false
  }
}

const deleteUser = async (id: string) => {
  try {
    await api.delete(`/user/${id}`)
    fetchUsers()
    message.success('删除成功')
  } catch (error: any) {
    message.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-page {
  padding: 24px;
}

.content-card {
  margin-bottom: 24px;
}
</style>
