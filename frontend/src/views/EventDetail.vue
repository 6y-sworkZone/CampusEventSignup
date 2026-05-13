<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 1000px; margin: 0 auto;" v-loading="loading">
      <el-card v-if="event" class="detail-card">
        <div class="cover-section">
          <img :src="event.cover_image || 'https://picsum.photos/1000/400'" class="detail-cover" alt="活动封面" />
          <div class="cover-overlay">
            <div class="event-title">{{ event.title }}</div>
            <el-tag :type="event.status === 'open' ? 'success' : 'info'" size="large">
              {{ getStatusText(event.status) }}
            </el-tag>
          </div>
        </div>

        <div class="detail-content">
          <div class="info-grid">
            <div class="info-item">
              <el-icon color="#409eff"><Calendar /></el-icon>
              <div>
                <div class="label">开始时间</div>
                <div class="value">{{ formatDateTime(event.start_time) }}</div>
              </div>
            </div>
            <div class="info-item">
              <el-icon color="#409eff"><Clock /></el-icon>
              <div>
                <div class="label">结束时间</div>
                <div class="value">{{ formatDateTime(event.end_time) }}</div>
              </div>
            </div>
            <div class="info-item">
              <el-icon color="#409eff"><Calendar /></el-icon>
              <div>
                <div class="label">报名截止</div>
                <div class="value">{{ formatDateTime(event.registration_end) }}</div>
              </div>
            </div>
            <div class="info-item">
              <el-icon color="#409eff"><User /></el-icon>
              <div>
                <div class="label">报名人数</div>
                <div class="value">{{ event.current_count }}/{{ event.max_participants }}人</div>
              </div>
            </div>
          </div>

          <div class="tags-section mb-20">
            <el-tag v-for="tag in event.tags" :key="tag.id" size="large">{{ tag.name }}</el-tag>
          </div>

          <div class="description-section mb-20">
            <h3>活动详情</h3>
            <div v-html="event.description || '暂无详细描述'"></div>
          </div>

          <div class="action-section" v-if="userStore.isLoggedIn">
            <template v-if="event.status === 'open'">
              <el-button type="primary" size="large" :loading="registerLoading" @click="handleRegister">
                立即报名
              </el-button>
              <el-button size="large" @click="handleSignin" v-if="canSignin">
                签到
              </el-button>
            </template>
            <template v-else-if="event.status === 'ongoing'">
              <el-button size="large" @click="handleSignin">
                签到
              </el-button>
            </template>
          </div>
        </div>
      </el-card>
    </div>
  </div>

  <el-dialog v-model="signinDialogVisible" title="活动签到" width="400px">
    <el-form :model="signinForm" label-width="80px">
      <el-form-item label="签到码">
        <el-input v-model="signinForm.code" placeholder="请输入4位签到码" maxlength="4" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="signinDialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="signinLoading" @click="submitSignin">确认签到</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const registerLoading = ref(false)
const signinLoading = ref(false)
const event = ref(null)
const signinDialogVisible = ref(false)
const signinForm = ref({ code: '' })

const canSignin = computed(() => {
  if (!event.value) return false
  const startTime = new Date(event.value.start_time)
  const now = new Date()
  return startTime.getTime() - now.getTime() <= 30 * 60 * 1000
})

const fetchEvent = async () => {
  loading.value = true
  try {
    const res = await api.get(`/events/${route.params.id}`)
    event.value = res.data
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  try {
    await ElMessageBox.confirm('确认要报名这个活动吗？', '报名确认', {
      confirmButtonText: '确认报名',
      cancelButtonText: '取消',
      type: 'info'
    })
    registerLoading.value = true
    const res = await api.post(`/events/${event.value.id}/register`)
    ElMessage.success(res.data.message || '报名成功')
    fetchEvent()
  } finally {
    registerLoading.value = false
  }
}

const handleSignin = () => {
  signinDialogVisible.value = true
}

const submitSignin = async () => {
  if (!signinForm.value.code || signinForm.value.code.length !== 4) {
    ElMessage.error('请输入4位签到码')
    return
  }
  signinLoading.value = true
  try {
    await api.post(`/events/${event.value.id}/signin`, { sign_code: signinForm.value.code })
    ElMessage.success('签到成功')
    signinDialogVisible.value = false
  } finally {
    signinLoading.value = false
  }
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

onMounted(() => {
  fetchEvent()
})
</script>

<style scoped>
.detail-card {
  overflow: hidden;
}

.cover-section {
  position: relative;
  margin: -20px -20px 20px;
}

.detail-cover {
  width: 100%;
  height: 300px;
  object-fit: cover;
  display: block;
}

.cover-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 30px 20px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: white;
}

.event-title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 12px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.info-item .label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
}

.info-item .value {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.tags-section {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.description-section h3 {
  font-size: 18px;
  margin-bottom: 12px;
  color: #303133;
}

.action-section {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .event-title {
    font-size: 22px;
  }
}
</style>
