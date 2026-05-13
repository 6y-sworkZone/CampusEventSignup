<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 1200px; margin: 0 auto;">
      <h2 class="page-title">我的报名</h2>
      
      <el-tabs v-model="activeTab" class="mb-20">
        <el-tab-pane label="全部" name="" />
        <el-tab-pane label="已报名" name="registered" />
        <el-tab-pane label="已签到" name="signed" />
        <el-tab-pane label="候补" name="waitlist" />
        <el-tab-pane label="已取消" name="cancelled" />
      </el-tabs>

      <div v-loading="loading">
        <div v-if="registrations.length === 0" class="empty-state">
          <el-icon class="empty-icon"><Document /></el-icon>
          <p>暂无报名记录</p>
        </div>

        <el-card v-for="reg in registrations" :key="reg.id" class="reg-card mb-20" shadow="hover">
          <div class="reg-content">
            <div class="reg-info">
              <h3 class="event-title" @click="$router.push(`/events/${reg.event_id}`)">{{ reg.event.title }}</h3>
              <div class="event-meta">
                <span><el-icon size="14"><Calendar /></el-icon> {{ formatDateTime(reg.event.start_time) }}</span>
                <span><el-icon size="14"><User /></el-icon> {{ reg.event.current_count }}/{{ reg.event.max_participants }}人</span>
              </div>
              <div class="event-tags">
                <el-tag size="small" v-for="tag in reg.event.tags" :key="tag.id">{{ tag.name }}</el-tag>
              </div>
            </div>
            <div class="reg-status">
              <el-tag :type="getStatusType(reg)" size="large">
                {{ getStatusText(reg) }}
              </el-tag>
              <div class="reg-time">报名时间: {{ formatDateTime(reg.created_at) }}</div>
            </div>
            <div class="reg-actions">
              <el-button 
                v-if="canCancel(reg)" 
                type="danger" 
                size="small" 
                @click="handleCancel(reg)"
              >
                取消报名
              </el-button>
              <el-button 
                v-if="canSignin(reg)" 
                type="primary" 
                size="small"
                @click="handleSignin(reg)"
              >
                签到
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
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
import { ref, onMounted, watch } from 'vue'
import api from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const signinLoading = ref(false)
const registrations = ref([])
const activeTab = ref('')
const signinDialogVisible = ref(false)
const signinForm = ref({ code: '' })
const currentReg = ref(null)

const fetchRegistrations = async () => {
  loading.value = true
  try {
    const res = await api.get('/registrations/my', { params: { status: activeTab.value } })
    registrations.value = res.data
  } finally {
    loading.value = false
  }
}

watch(activeTab, () => {
  fetchRegistrations()
})

const canCancel = (reg) => {
  if (reg.status === 'cancelled' || reg.signed) return false
  const startTime = new Date(reg.event.start_time)
  const now = new Date()
  return startTime.getTime() - now.getTime() > 2 * 60 * 60 * 1000
}

const canSignin = (reg) => {
  if (reg.status === 'cancelled' || reg.signed || reg.is_waitlist) return false
  const startTime = new Date(reg.event.start_time)
  const now = new Date()
  return startTime.getTime() - now.getTime() <= 30 * 60 * 1000
}

const handleCancel = async (reg) => {
  try {
    await ElMessageBox.confirm('确认要取消报名吗？活动开始前2小时内不可取消', '取消确认', {
      confirmButtonText: '确认取消',
      cancelButtonText: '保留',
      type: 'warning'
    })
    await api.post(`/events/${reg.event_id}/cancel`)
    ElMessage.success('取消成功')
    fetchRegistrations()
  } catch {}
}

const handleSignin = (reg) => {
  currentReg.value = reg
  signinDialogVisible.value = true
}

const submitSignin = async () => {
  if (!signinForm.value.code || signinForm.value.code.length !== 4) {
    ElMessage.error('请输入4位签到码')
    return
  }
  signinLoading.value = true
  try {
    await api.post(`/events/${currentReg.value.event_id}/signin`, { sign_code: signinForm.value.code })
    ElMessage.success('签到成功')
    signinDialogVisible.value = false
    fetchRegistrations()
  } finally {
    signinLoading.value = false
  }
}

const formatDateTime = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const getStatusText = (reg) => {
  if (reg.signed) return '已签到'
  if (reg.is_waitlist) return '候补排队'
  const map = {
    registered: '已报名',
    cancelled: '已取消',
    signed: '已签到'
  }
  return map[reg.status] || reg.status
}

const getStatusType = (reg) => {
  if (reg.signed) return 'success'
  if (reg.is_waitlist) return 'warning'
  if (reg.status === 'cancelled') return 'danger'
  return 'info'
}

onMounted(() => {
  fetchRegistrations()
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #303133;
}

.reg-card {
  cursor: default;
}

.reg-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.reg-info {
  flex: 1;
}

.event-title {
  font-size: 18px;
  font-weight: 500;
  margin: 0 0 10px;
  cursor: pointer;
  color: #303133;
}

.event-title:hover {
  color: #409eff;
}

.event-meta {
  display: flex;
  gap: 20px;
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
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

.reg-status {
  text-align: center;
}

.reg-time {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.reg-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

@media (max-width: 768px) {
  .reg-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .reg-actions {
    flex-direction: row;
  }
}
</style>
