<template>
  <div>
    <h1>任务管理应用</h1>

    <!-- 添加任务表单 -->
    <div class="task-form">
      <input
        v-model="newTask"
        class="task-input"
        placeholder="添加新任务..."
        @keyup.enter="addTask"
      />
      <button class="add-button" @click="addTask">添加</button>
    </div>

    <!-- 任务列表 -->
    <div class="task-list">
      <div v-if="loading">加载中...</div>
      <div v-else-if="tasks.length === 0">暂无任务</div>
      <div
        v-else
        v-for="task in tasks"
        :key="task.id"
        class="task-item"
        :class="{ completed: task.completed }"
      >
        <input
          type="checkbox"
          :checked="task.completed"
          @change="toggleTask(task)"
        />
        <span class="task-title">{{ task.title }}</span>
        <button @click="deleteTask(task.id)">删除</button>
      </div>
    </div>

    <!-- 状态消息 -->
    <div v-if="message" class="message">{{ message }}</div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'

export default {
  name: 'App',
  setup() {
    const API_URL = 'http://localhost:8080/api'
    const tasks = ref([])
    const newTask = ref('')
    const loading = ref(false)
    const message = ref('')

    // 显示临时消息
    const showMessage = (msg) => {
      message.value = msg
      setTimeout(() => {
        message.value = ''
      }, 3000)
    }

    // 获取所有任务
    const fetchTasks = async () => {
      loading.value = true
      console.log('开始获取任务列表')
      console.log('请求URL:', `${API_URL}/tasks`)

      try {
        const response = await axios.get(`${API_URL}/tasks`)
        console.log('获取任务成功:', response)
        tasks.value = response.data
        console.log('当前任务列表:', tasks.value)
      } catch (error) {
        console.error('获取任务失败:', error)
        console.error('错误详情:', error.response ? error.response.data : '无响应数据')
        showMessage('获取任务失败')
      } finally {
        loading.value = false
      }
    }

    // 添加新任务
    const addTask = async () => {
      if (!newTask.value.trim()) return

      console.log('准备添加任务:', newTask.value.trim())

      try {
        const taskData = {
          title: newTask.value.trim(),
          completed: false
        }
        console.log('发送的数据:', taskData)
        console.log('请求URL:', `${API_URL}/tasks`)

        const response = await axios.post(`${API_URL}/tasks`, taskData, {
          headers: {
            'Content-Type': 'application/json'
          }
        })

        console.log('服务器响应:', response)
        tasks.value.push(response.data)
        newTask.value = ''
        showMessage('任务添加成功')
      } catch (error) {
        console.error('添加任务失败:', error)
        console.error('错误详情:', error.response ? error.response.data : '无响应数据')
        showMessage('添加任务失败')
      }
    }

    // 切换任务状态
    const toggleTask = async (task) => {
      try {
        const updatedTask = { ...task, completed: !task.completed }
        await axios.put(`${API_URL}/tasks/${task.id}`, updatedTask)

        // 更新本地状态
        const index = tasks.value.findIndex(t => t.id === task.id)
        if (index !== -1) {
          tasks.value[index].completed = !tasks.value[index].completed
        }
      } catch (error) {
        console.error('更新任务失败:', error)
        showMessage('更新任务失败')
      }
    }

    // 删除任务
    const deleteTask = async (id) => {
      try {
        await axios.delete(`${API_URL}/tasks/${id}`)
        tasks.value = tasks.value.filter(task => task.id !== id)
        showMessage('任务已删除')
      } catch (error) {
        console.error('删除任务失败:', error)
        showMessage('删除任务失败')
      }
    }

    // 组件挂载时获取任务
    onMounted(fetchTasks)

    return {
      tasks,
      newTask,
      loading,
      message,
      addTask,
      toggleTask,
      deleteTask
    }
  }
}
</script>

<style scoped>
.message {
  margin-top: 2em;
  padding: 0.5em;
  background-color: #4caf50;
  color: white;
  border-radius: 4px;
}
</style>
