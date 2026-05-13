<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 1400px; margin: 0 auto;">
      <div class="flex-between mb-20">
        <h2 class="page-title">活动管理</h2>
        <el-button type="primary" @click="$router.push('/admin/events/create')">创建活动</el-button>
      </div>

      <el-card v-loading="loading">
        <el-table :data="events" stripe>
          <el-table-column prop="title" label="活动名称" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="时间" width="280">
            <template #default="{ row }">
              <div>{{ formatDateTime(row.start_time) }}</div>
              <div style="font-size: 12px; color: #909399;">至 {{ formatDateTime(row.end_time) }}</div>
            </template>
          </el-table-column>
          <el-table-column label="报名人数" width="120">
            <template #default="{ row }">
              {{ row.current_count }}/{{ row.max_participants }}
            </template>
          </el-table-column>
          <el-table-column label="签到码" width="100">
            <template #default="{ row }">
              <el-tag type="warning">{{ row.sign_code }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="$router.push(`/admin/events/${row.id}/edit`)">编辑</el-button>
              <el-button link type="success" size="small" @click="viewDetail(row)">详情</el-button>
              <el-button v-if="row.status === 'open'" link type="warning" size="small" @click="endEvent(row)">结束</el-button>
              <el-button link type="danger" size="small" @click="deleteEvent(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination mt-20">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="fetchEvents"
            @size-change="fetchEvents"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const events = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const fetchEvents = async () => {
  loading.value = true
  try {
    const res = await api.get('/events', {
      params: { page: page.value, page_size: pageSize.value }
    })
    events.value = res.data.list
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

const viewDetail = (row) => {
  router.push(`/events/${row.id}`)
}

const endEvent = async (row) => {
  try {
    await ElMessageBox.confirm('确认要结束这个活动吗？结束后将无法继续报名', '结束确认', {
      confirmButtonText: '确认结束',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await api.post(`/events/${row.id}/end`)
    ElMessage.success('活动已结束')
    fetchEvents()
  } catch {}
}

const deleteEvent = async (row) => {
  try {
    await ElMessageBox.confirm('确认要删除这个活动吗？此操作不可恢复', '删除确认', {
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await api.delete(`/events/${row.id}`)
    ElMessage.success('删除成功')
    fetchEvents()
  } catch {}
}

const formatDateTime = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const getStatusText = (status) => {
  const map = {
    draft: '草稿',
    open: '报名中',
    ongoing: '进行中',
    ended: '已结束',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getStatusType = (status) => {
  const map = {
    draft: 'info',
    open: 'success',
    ongoing: 'warning',
    ended: 'info',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

onMounted(() => {
  fetchEvents()
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin: 0;
  color: #303133;
}

.pagination {
  text-align: center;
}
</style>
