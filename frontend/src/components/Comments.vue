<template>
  <div class="comments-section">
    <h3>评论区 ({{ comments.length }})</h3>
    
    <div v-if="userStore.isLoggedIn" class="comment-input-section">
      <el-avatar :size="40" :src="userStore.user?.avatar">
        <el-icon><User /></el-icon>
      </el-avatar>
      <div class="input-wrapper">
        <el-input
          v-model="newComment"
          type="textarea"
          :rows="3"
          placeholder="写下你的评论..."
          maxlength="500"
          show-word-limit
        />
        <div class="submit-btn-wrapper">
          <el-button type="primary" :loading="submitting" @click="submitComment">
            发表评论
          </el-button>
        </div>
      </div>
    </div>
    
    <div v-else class="login-prompt">
      <el-button type="text" @click="$router.push('/login')">登录后即可发表评论</el-button>
    </div>

    <div class="comments-list" v-loading="loading">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <el-avatar :size="40" :src="comment.user?.avatar">
          <el-icon><User /></el-icon>
        </el-avatar>
        <div class="comment-content">
          <div class="comment-header">
            <span class="username">{{ comment.user?.username }}</span>
            <span class="time">{{ formatDate(comment.created_at) }}</span>
            <el-button
              v-if="userStore.isAdmin || (userStore.isLoggedIn && userStore.user?.id === comment.user_id)"
              type="text"
              size="small"
              class="delete-btn"
              @click="deleteComment(comment.id)"
            >
              删除
            </el-button>
          </div>
          <div class="comment-text">{{ comment.content }}</div>
          
          <div class="replies-section">
            <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
              <el-avatar :size="32" :src="reply.user?.avatar">
                <el-icon><User /></el-icon>
              </el-avatar>
              <div class="reply-content">
                <div class="reply-header">
                  <span class="username">{{ reply.user?.username }}</span>
                  <span class="time">{{ formatDate(reply.created_at) }}</span>
                  <el-button
                    v-if="userStore.isAdmin || (userStore.isLoggedIn && userStore.user?.id === reply.user_id)"
                    type="text"
                    size="small"
                    class="delete-btn"
                    @click="deleteComment(reply.id)"
                  >
                    删除
                  </el-button>
                </div>
                <div class="reply-text">{{ reply.content }}</div>
              </div>
            </div>
          </div>
          
          <el-button
            v-if="userStore.isLoggedIn"
            type="text"
            size="small"
            class="reply-btn"
            @click="showReplyInput(comment)"
          >
            回复
          </el-button>
          
          <div v-if="replyingTo === comment.id" class="reply-input-section">
            <el-input
              v-model="replyContent"
              type="textarea"
              :rows="2"
              placeholder="写下你的回复..."
              maxlength="500"
            />
            <div class="reply-actions">
              <el-button size="small" @click="cancelReply">取消</el-button>
              <el-button size="small" type="primary" :loading="submitting" @click="submitReply(comment.id)">
                回复
              </el-button>
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="comments.length === 0 && !loading" class="empty-comments">
        <el-icon size="48" color="#909399"><ChatDotRound /></el-icon>
        <p>暂无评论，快来发表第一条评论吧！</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  eventId: {
    type: [Number, String],
    required: true
  }
})

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const comments = ref([])
const newComment = ref('')
const replyContent = ref('')
const replyingTo = ref(null)

const fetchComments = async () => {
  loading.value = true
  try {
    const res = await api.get(`/events/${props.eventId}/comments`)
    comments.value = res.data
  } finally {
    loading.value = false
  }
}

const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  
  submitting.value = true
  try {
    await api.post(`/events/${props.eventId}/comments`, {
      content: newComment.value.trim()
    })
    ElMessage.success('评论成功')
    newComment.value = ''
    fetchComments()
  } finally {
    submitting.value = false
  }
}

const showReplyInput = (comment) => {
  replyingTo.value = comment.id
  replyContent.value = ''
}

const cancelReply = () => {
  replyingTo.value = null
  replyContent.value = ''
}

const submitReply = async (parentId) => {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  
  submitting.value = true
  try {
    await api.post(`/events/${props.eventId}/comments`, {
      content: replyContent.value.trim(),
      parent_id: parentId
    })
    ElMessage.success('回复成功')
    replyingTo.value = null
    replyContent.value = ''
    fetchComments()
  } finally {
    submitting.value = false
  }
}

const deleteComment = async (commentId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await api.delete(`/comments/${commentId}`)
    ElMessage.success('删除成功')
    fetchComments()
  } catch {}
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

onMounted(() => {
  fetchComments()
})
</script>

<style scoped>
.comments-section {
  margin-top: 30px;
}

.comments-section h3 {
  margin: 0 0 20px;
  font-size: 18px;
  color: #303133;
}

.comment-input-section {
  display: flex;
  gap: 12px;
  margin-bottom: 30px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.input-wrapper {
  flex: 1;
}

.submit-btn-wrapper {
  margin-top: 12px;
  text-align: right;
}

.login-prompt {
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
  text-align: center;
  margin-bottom: 30px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  display: flex;
  gap: 12px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.username {
  font-weight: 500;
  color: #303133;
}

.time {
  font-size: 12px;
  color: #909399;
}

.delete-btn {
  margin-left: auto;
  color: #f56c6c;
}

.comment-text {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 8px;
}

.replies-section {
  margin-top: 12px;
  margin-left: 20px;
  padding-left: 20px;
  border-left: 2px solid #e4e7ed;
}

.reply-item {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 4px;
}

.reply-text {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}

.reply-btn {
  color: #409eff;
  padding: 0;
}

.reply-input-section {
  margin-top: 12px;
  margin-left: 20px;
}

.reply-actions {
  margin-top: 8px;
  text-align: right;
}

.empty-comments {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-comments p {
  margin: 16px 0 0;
}
</style>
