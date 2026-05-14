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
      <div class="hot-tags-section" v-if="hotTags.length > 0">
        <div class="section-header flex-between">
          <h2>热门标签</h2>
        </div>
        <div class="tags-container">
          <el-tag
            v-for="tag in hotTags"
            :key="tag.name"
            :type="selectedTag === tag.name ? 'primary' : 'info'"
            size="large"
            class="hot-tag"
            @click="toggleTag(tag.name)"
            effect="plain"
          >
            {{ tag.name }} ({{ tag.count }})
          </el-tag>
        </div>
      </div>

      <div class="section-header flex-between" style="margin-top: 30px;">
        <h2>最新活动</h2>
        <div class="header-actions">
          <el-button v-if="selectedTag" type="text" @click="clearTag">
            清除筛选
          </el-button>
          <el-button link @click="$router.push('/events')">查看更多</el-button>
        </div>
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
        <p>{{ selectedTag ? '该标签下暂无活动' : '暂无活动' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/request'

const loading = ref(false)
const events = ref([])
const hotTags = ref([])
const selectedTag = ref('')

const fetchHotTags = async () => {
  try {
    const res = await api.get('/events/tags/hot')
    hotTags.value = res.data
  } catch (e) {
    console.error('获取热门标签失败', e)
  }
}

const fetchEvents = async () => {
  loading.value = true
  try {
    const params = { page: 1, page_size: 8 }
    if (selectedTag.value) {
      params.tag = selectedTag.value
    }
    const res = await api.get('/events', { params })
    events.value = res.data.list
  } finally {
    loading.value = false
  }
}

const toggleTag = (tagName) => {
  if (selectedTag.value === tagName) {
    selectedTag.value = ''
  } else {
    selectedTag.value = tagName
  }
  fetchEvents()
}

const clearTag = () => {
  selectedTag.value = ''
  fetchEvents()
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
  fetchHotTags()
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-header h2 {
  font-size: 20px;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hot-tags-section {
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.hot-tag {
  cursor: pointer;
  font-size: 14px;
  padding: 8px 16px;
  transition: all 0.3s;
}

.hot-tag:hover {
  transform: translateY(-2px);
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

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .banner-content h1 {
    font-size: 32px;
  }
}
</style>
