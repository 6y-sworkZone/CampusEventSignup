<template>
  <div class="page-container">
    <div class="home-banner">
      <div class="banner-content">
        <h1>校园活动报名系统</h1>
        <p>发现精彩活动，丰富校园生活</p>
        <el-button type="primary" size="large" @click="$router.push('/events')">浏览活动</el-button>
      </div>
    </div>
    
    <div class="home-content" style="max-width: 1400px; margin: 0 auto;">
      <div class="section-header flex-between">
        <h2>最新活动</h2>
        <el-button link @click="$router.push('/events')">查看更多</el-button>
      </div>
      
      <el-row :gutter="20" v-loading="loading">
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
        <el-icon class="empty-icon"><Document /></el-icon>
        <p>暂无活动</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/request'

const loading = ref(false)
const events = ref([])

const fetchEvents = async () => {
  loading.value = true
  try {
    const res = await api.get('/events', { params: { page: 1, page_size: 8 } })
    events.value = res.data.list
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日`
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
  fetchEvents()
})
</script>

<style scoped>
.home-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 80px 20px;
  text-align: center;
  color: white;
  margin-bottom: 30px;
}

.banner-content h1 {
  font-size: 48px;
  margin-bottom: 16px;
}

.banner-content p {
  font-size: 18px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.section-header {
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 24px;
  color: #303133;
  margin: 0;
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

@media (max-width: 768px) {
  .banner-content h1 {
    font-size: 32px;
  }
}
</style>
