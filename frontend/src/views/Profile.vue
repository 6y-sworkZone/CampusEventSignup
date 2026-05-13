<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 800px; margin: 0 auto;">
      <h2 class="page-title">个人中心</h2>
      
      <el-card v-loading="loading" class="profile-card">
        <div class="avatar-section">
          <el-avatar :size="120" :src="user?.avatar">
            <el-icon size="60"><User /></el-icon>
          </el-avatar>
          <el-upload
            class="avatar-uploader"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
            accept="image/*"
          >
            <el-button type="primary" size="small">更换头像</el-button>
          </el-upload>
        </div>

        <el-divider />

        <el-form :model="form" label-width="100px">
          <el-form-item label="用户名">
            <el-input v-model="form.username" disabled />
          </el-form-item>
          <el-form-item label="角色">
            <el-tag>{{ user?.role === 'admin' ? '管理员' : '学生' }}</el-tag>
          </el-form-item>
          <el-form-item v-if="user?.student_id" label="学号">
            <el-input v-model="form.student_id" disabled />
          </el-form-item>
          <el-form-item v-if="user?.employee_id" label="工号">
            <el-input v-model="form.employee_id" disabled />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="form.phone" />
          </el-form-item>
          <el-form-item label="学院">
            <el-select v-model="form.college" class="w-full">
              <el-option label="计算机学院" value="计算机学院" />
              <el-option label="电子信息学院" value="电子信息学院" />
              <el-option label="商学院" value="商学院" />
              <el-option label="文学院" value="文学院" />
              <el-option label="理学院" value="理学院" />
              <el-option label="工学院" value="工学院" />
              <el-option label="法学院" value="法学院" />
              <el-option label="医学院" value="医学院" />
            </el-select>
          </el-form-item>
        </el-form>

        <el-divider />

        <h3>修改密码</h3>
        <el-form :model="passwordForm" label-width="100px">
          <el-form-item label="原密码">
            <el-input v-model="passwordForm.old_password" type="password" show-password />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="passwordForm.new_password" type="password" show-password />
          </el-form-item>
          <el-form-item label="确认密码">
            <el-input v-model="passwordForm.confirm_password" type="password" show-password />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="passwordLoading" @click="changePassword">修改密码</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/request'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const loading = ref(false)
const passwordLoading = ref(false)
const user = ref(null)
const form = ref({
  username: '',
  student_id: '',
  employee_id: '',
  phone: '',
  college: ''
})
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const uploadUrl = computed(() => {
  return `${import.meta.env.VITE_API_BASE_URL || ''}/api/upload/avatar`
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

const fetchUser = async () => {
  loading.value = true
  try {
    const res = await api.get('/user/me')
    user.value = res.data
    form.value = {
      username: res.data.username,
      student_id: res.data.student_id || '',
      employee_id: res.data.employee_id || '',
      phone: res.data.phone || '',
      college: res.data.college || ''
    }
  } finally {
    loading.value = false
  }
}

const beforeAvatarUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
  }
  return isImage && isLt2M
}

const handleAvatarSuccess = (response) => {
  ElMessage.success('头像上传成功')
  fetchUser()
}

const changePassword = async () => {
  if (!passwordForm.value.old_password) {
    ElMessage.error('请输入原密码')
    return
  }
  if (!passwordForm.value.new_password || passwordForm.value.new_password.length < 6) {
    ElMessage.error('新密码至少6位')
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.error('两次输入密码不一致')
    return
  }

  passwordLoading.value = true
  try {
    await api.post('/user/change-password', {
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    })
    ElMessage.success('密码修改成功')
    passwordForm.value = {
      old_password: '',
      new_password: '',
      confirm_password: ''
    }
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  fetchUser()
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #303133;
}

.profile-card {
  padding: 30px;
}

.avatar-section {
  text-align: center;
}

.avatar-uploader {
  margin-top: 16px;
}
</style>
