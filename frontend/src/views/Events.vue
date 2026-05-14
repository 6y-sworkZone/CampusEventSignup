<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 1400px; margin: 0 auto;">
      <div class="filter-section card mb-20">
        <el-form :inline="true" :model="filters">
          <el-form-item label="状态">
            <el-select v-model="filters.status" placeholder="全部状态" clearable>
              <el-option label="报名中" value="open" />
              <el-option label="进行中" value="ongoing" />
              <el-option label="已结束" value="ended" />
            </el-select>
          </el-form-item>
          <el-form-item label="搜索">
            <el-input v-model="filters.keyword" placeholder="搜索活动" clearable />
          </el-form-item>
          <el-form-item label="标签">
            <el-select v-model="filters.tag" placeholder="选择标签" clearable>
              <el-option label="讲座" value="讲座" />
              <el-option label="比赛" value="比赛" />
              <el-option label="文艺" value="文艺" />
              <el-option label="体育" value="体育" />
              <el-option label="学术" value="学术" />
              <el-option label="志愿" value="志愿" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="resetAndFetch" :loading="loading">搜索</el-button>
            <el-button @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div v-loading="loading && events.length === 0">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="event in events" :key="event.id">
            <el-card class="event-card" shadow="hover" @click="$router.push(`/events/${event.id}`)">
              <template #cover>
                <div class="event-cover-wrapper">
                  <img :src="event.cover_image || 'https://picsum.photos/400/200?random=' + event.id" class="event-cover" alt="活动封面" />
                  <div class="event-status" :class="event.status">
                    {{ getStatusText(event.status) }}
                  </div>
                </div>
              </template>
              <div class="event-info">
                <h3 class="event-title">{{ event.title }}</h3>
                <div class="event-meta">
                  <span><el-icon size="14"><Calendar /></el-icon> {{ formatDate(event.start_time) }}</span>
                  <span><el-icon size="14"><User /></el-icon> {{ event.current_count }}/{{ event.max_participants }}</span>
                </div>
                <div class="event-tags">
                  <el-tag size="small" v-for="tag in event.tags" :key="tag.id">{{ tag.name }}</el-tag>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <div v-if="events.length === 0 && !loading" class="empty-state">
          <el-icon class="empty-icon"><Search /></el-icon>
          <p>没有找到符合条件的活动</p>
        </div>

        <div v-if="loading && events.length > 0" class="loading-more">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>

        <div v-if="!hasMore && events.length > 0" class="no-more">
          没有更多了
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '@/utils/request'

const loading = ref(false)
const events = ref([])
const total = ref(0)
const hasMore = ref(true)

const filters = ref({
  page: 1,
  page_size: 16,
  status: '',
  keyword: '',
  tag: ''
})

const fetchEvents = async (isLoadMore = false) => {
  if (loading.value) return
  
  loading.value = true
  try {
    const res = await api.get('/events', { params: filters.value })
    if (isLoadMore) {
      events.value = [...events.value, ...res.data.list]
    } else {
      events.value = res.data.list
    }
    total.value = res.data.total
    hasMore.value = events.value.length < total.value
  } finally {
    loading.value = false
  }
}

const resetAndFetch = () => {
  filters.value.page = 1
  hasMore.value = true
  events.value = []
  fetchEvents(false)
}

const resetFilters = () => {
  filters.value = {
    page: 1,
    page_size: 16,
    status: '',
    keyword: '',
    tag: ''
  }
  resetAndFetch()
}

const handleScroll = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
  const scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight
  const clientHeight = document.documentElement.clientHeight || window.innerHeight

  if (scrollTop + clientHeight >= scrollHeight - 100) {
    if (!loading.value && hasMore.value) {
      filters.value.page++
      fetchEvents(true)
    }
  }
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
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

onMounted(() => {
  fetchEvents(false)
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.filter-section {
  padding: 20px;
}

.event-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.3s;
}

.event-card:hover {
  transform: translateY(-5px);
}

.event-cover-wrapper {
  position: relative;
}

.event-status {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  color: white;
}

.event-status.open {
  background: #67c23a;
}

.event-status.ongoing {
  background: #e6a23c;
}

.event-status.ended {
  background: #909399;
}

.event-status.cancelled {
  background: #f56c6c;
}

.event-info h3 {
  margin: 0 0 12px;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event-meta {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 13px;
  margin-bottom: 12px;
}

.event-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.event-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #909399;
  gap: 8px;
}

.no-more {
  text-align: center;
  padding: 30px;
  color: #c0c4cc;
  font-size: 14px;
}
</style>
