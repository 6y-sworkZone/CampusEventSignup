<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 900px; margin: 0 auto;">
      <h2 class="page-title">{{ isEdit ? '编辑活动' : '创建活动' }}</h2>
      
      <el-card v-loading="loading">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
          <el-form-item label="活动标题" prop="title">
            <el-input v-model="form.title" placeholder="请输入活动标题（2-50字）" maxlength="50" />
          </el-form-item>

          <el-form-item label="活动封面">
            <el-upload
              class="cover-uploader"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleCoverSuccess"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <img v-if="form.cover_image" :src="form.cover_image" class="cover-image" />
              <div v-else class="upload-placeholder">
                <el-icon size="40"><Plus /></el-icon>
                <div>点击上传封面</div>
              </div>
            </el-upload>
          </el-form-item>

          <el-form-item label="活动描述" prop="description">
            <RichEditor v-model="form.description" />
          </el-form-item>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="开始时间" prop="start_time">
                <el-date-picker
                  v-model="form.start_time"
                  type="datetime"
                  placeholder="选择开始时间"
                  style="width: 100%"
                  value-format="YYYY-MM-DDTHH:mm:ss"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="结束时间" prop="end_time">
                <el-date-picker
                  v-model="form.end_time"
                  type="datetime"
                  placeholder="选择结束时间"
                  style="width: 100%"
                  value-format="YYYY-MM-DDTHH:mm:ss"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="报名截止时间" prop="registration_end">
            <el-date-picker
              v-model="form.registration_end"
              type="datetime"
              placeholder="选择报名截止时间"
              style="width: 100%"
              value-format="YYYY-MM-DDTHH:mm:ss"
            />
          </el-form-item>

          <el-form-item label="最大参与人数" prop="max_participants">
            <el-input-number v-model="form.max_participants" :min="1" :max="10000" />
          </el-form-item>

          <el-form-item label="活动标签">
            <el-select v-model="tags" multiple placeholder="选择活动标签">
              <el-option label="讲座" value="讲座" />
              <el-option label="比赛" value="比赛" />
              <el-option label="文艺" value="文艺" />
              <el-option label="体育" value="体育" />
              <el-option label="学术" value="学术" />
              <el-option label="志愿" value="志愿" />
              <el-option label="其他" value="其他" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="isEdit" label="活动状态">
            <el-tag :type="getStatusType(form.status)">{{ getStatusText(form.status) }}</el-tag>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
              {{ isEdit ? '保存修改' : '创建活动' }}
            </el-button>
            <el-button @click="$router.back()">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/utils/request'
import { ElMessage } from 'element-plus'
import RichEditor from '@/components/RichEditor.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isEdit = computed(() => !!route.params.id)

const loading = ref(false)
const submitLoading = ref(false)
const formRef = ref()
const tags = ref([])

const form = ref({
  title: '',
  description: '',
  cover_image: '',
  start_time: '',
  end_time: '',
  registration_end: '',
  max_participants: 100,
  status: 'open'
})

const rules = {
  title: [
    { required: true, message: '请输入活动标题', trigger: 'blur' },
    { min: 2, max: 50, message: '标题长度为2-50字', trigger: 'blur' }
  ],
  start_time: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  registration_end: [{ required: true, message: '请选择报名截止时间', trigger: 'change' }],
  max_participants: [{ required: true, message: '请输入最大参与人数', trigger: 'blur' }]
}

const uploadUrl = computed(() => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || ''
  return `${baseUrl}/api/upload/cover`
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

const fetchEvent = async () => {
  loading.value = true
  try {
    const res = await api.get(`/events/${route.params.id}`)
    form.value = {
      title: res.data.title,
      description: res.data.description,
      cover_image: res.data.cover_image,
      start_time: res.data.start_time,
      end_time: res.data.end_time,
      registration_end: res.data.registration_end,
      max_participants: res.data.max_participants,
      status: res.data.status
    }
    tags.value = res.data.tags?.map((t) => t.name) || []
  } finally {
    loading.value = false
  }
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
  }
  return isImage && isLt5M
}

const handleCoverSuccess = (response) => {
  form.value.cover_image = response.data.url
  ElMessage.success('封面上传成功')
}

const handleSubmit = async () => {
  await formRef.value.validate()
  
  const startTime = new Date(form.value.start_time)
  const endTime = new Date(form.value.end_time)
  const registrationEnd = new Date(form.value.registration_end)

  if (endTime <= startTime) {
    ElMessage.error('结束时间必须晚于开始时间')
    return
  }
  if (registrationEnd >= startTime) {
    ElMessage.error('报名截止时间必须早于活动开始时间')
    return
  }

  submitLoading.value = true
  try {
    const data = {
      ...form.value,
      tags: tags.value
    }

    if (isEdit.value) {
      await api.put(`/events/${route.params.id}`, data)
      ElMessage.success('修改成功')
    } else {
      await api.post('/events', data)
      ElMessage.success('创建成功')
    }
    router.push('/admin/events')
  } finally {
    submitLoading.value = false
  }
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
  if (isEdit.value) {
    fetchEvent()
  }
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #303133;
}

.cover-uploader {
  display: block;
}

.cover-image {
  width: 100%;
  max-width: 400px;
  height: 200px;
  object-fit: cover;
  border-radius: 8px;
}

.upload-placeholder {
  width: 400px;
  height: 200px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-placeholder:hover {
  border-color: #409eff;
}
</style>
