<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <el-icon size="40" color="#409eff"><Calendar /></el-icon>
        <h2>创建账户</h2>
        <p>加入校园活动报名系统</p>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="form.role">
            <el-radio label="student">学生</el-radio>
            <el-radio label="admin">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item v-if="form.role === 'student'" label="学号" prop="studentId">
          <el-input v-model="form.studentId" placeholder="请输入学号" />
        </el-form-item>
        <el-form-item v-if="form.role === 'admin'" label="工号" prop="employeeId">
          <el-input v-model="form.employeeId" placeholder="请输入工号" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入11位手机号" maxlength="11" />
        </el-form-item>
        <el-form-item label="学院" prop="college">
          <el-select v-model="form.college" placeholder="请选择学院" class="w-full">
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
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码（至少6位）" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="w-full" :loading="loading" @click="handleRegister">注 册</el-button>
        </el-form-item>
      </el-form>
      <div class="register-footer">
        <span>已有账户？</span>
        <router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const loading = ref(false)
const form = ref({
  role: 'student',
  username: '',
  studentId: '',
  employeeId: '',
  phone: '',
  college: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.value.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const validatePhone = (rule, value, callback) => {
  if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的手机号'))
  } else {
    callback()
  }
}

const rules = {
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  username: [{ required: true, min: 2, max: 50, message: '用户名长度为2-50个字符', trigger: 'blur' }],
  studentId: [{ required: true, message: '请输入学号', trigger: 'blur' }],
  employeeId: [{ required: true, message: '请输入工号', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { validator: validatePhone, trigger: 'blur' }
  ],
  college: [{ required: true, message: '请选择学院', trigger: 'change' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  await formRef.value.validate()
  loading.value = true
  try {
    await userStore.register({
      username: form.value.username,
      student_id: form.value.studentId,
      employee_id: form.value.employeeId,
      phone: form.value.phone,
      college: form.value.college,
      password: form.value.password,
      role: form.value.role
    })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  margin: 16px 0 8px;
  color: #303133;
}

.register-header p {
  color: #909399;
  margin: 0;
}

.register-footer {
  text-align: center;
  margin-top: 20px;
  color: #909399;
}

.register-footer a {
  color: #409eff;
  text-decoration: none;
  margin-left: 8px;
}
</style>
