<template>
  <div class="page-container" style="padding-top: 20px;">
    <div style="max-width: 1400px; margin: 0 auto;">
      <h2 class="page-title">数据统计</h2>

      <el-row :gutter="20" class="mb-20">
        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-icon" style="background: #409eff;">
              <el-icon size="28"><Calendar /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_events || 0 }}</div>
              <div class="stat-label">总活动数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-icon" style="background: #67c23a;">
              <el-icon size="28"><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_users || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-icon" style="background: #e6a23c;">
              <el-icon size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_registrations || 0 }}</div>
              <div class="stat-label">总报名数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-icon" style="background: #f56c6c;">
              <el-icon size="28"><School /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.college_stats?.length || 0 }}</div>
              <div class="stat-label">参与学院</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :xs="24" :lg="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>学院分布</span>
                <el-button type="primary" size="small" @click="exportData">导出数据</el-button>
              </div>
            </template>
            <div ref="collegeChart" class="chart-container"></div>
          </el-card>
        </el-col>
        <el-col :xs="24" :lg="12">
          <el-card class="chart-card">
            <template #header>
              <span>7天报名趋势</span>
            </template>
            <div ref="trendChart" class="chart-container"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '@/utils/request'
import * as echarts from 'echarts'

const collegeChart = ref(null)
const trendChart = ref(null)
const stats = ref({})
let collegeChartInstance = null
let trendChartInstance = null

const fetchStats = async () => {
  try {
    const res = await api.get('/stats/dashboard')
    stats.value = res.data
    renderCharts()
  } catch (e) {
    console.error('获取统计数据失败', e)
  }
}

const renderCharts = () => {
  if (collegeChart.value && stats.value.college_stats) {
    collegeChartInstance = echarts.init(collegeChart.value)
    collegeChartInstance.setOption({
      tooltip: {
        trigger: 'item'
      },
      legend: {
        orient: 'vertical',
        right: '5%',
        top: 'center'
      },
      series: [
        {
          name: '学院分布',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['40%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: stats.value.college_stats.map((item) => ({
            value: item.count,
            name: item.college
          }))
        }
      ]
    })
  }

  if (trendChart.value && stats.value.daily_stats) {
    trendChartInstance = echarts.init(trendChart.value)
    trendChartInstance.setOption({
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: stats.value.daily_stats.map((item) => {
          const date = new Date(item.date)
          return `${date.getMonth() + 1}/${date.getDate()}`
        })
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '报名人数',
          type: 'line',
          stack: 'Total',
          data: stats.value.daily_stats.map((item) => item.count),
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.5)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
            ])
          },
          smooth: true
        }
      ]
    })
  }
}

const exportData = () => {
  window.open('/api/stats/export', '_blank')
}

const handleResize = () => {
  collegeChartInstance?.resize()
  trendChartInstance?.resize()
}

onMounted(() => {
  fetchStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  collegeChartInstance?.dispose()
  trendChartInstance?.dispose()
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #303133;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.chart-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  width: 100%;
  height: 350px;
}
</style>
