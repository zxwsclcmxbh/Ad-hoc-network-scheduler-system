<template>
  <el-table
    v-loading="loading"
    :data="taskData.filter(data => !search || data.name.toLowerCase().includes(search.toLowerCase()))"
    :default-sort="{prop: 'CreatedAt', order: 'descending'}"
  >
    <el-table-column
      label="创建时间"
      prop="CreatedAt"
      sortable
    />
    <el-table-column
      label="任务名称"
      prop="name"
      align="center"
    />
    <el-table-column label="任务状态" align="center">
      <template v-slot="{row}">
        <el-tag :type="row.TaskStatus | statusFilter">
          {{ row.TaskStatus }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column
      label="任务描述"
      prop="description"
    />
    <el-table-column
      align="right"
    >
      <template slot="header" slot-scope="scope">
        <el-input
          v-model="search"
          size="mini"
          placeholder="输入关键字搜索"
          @input="change($event,scope.$index)"
        />
      </template>
      <template v-slot="scope">
        <el-button
          slot="reference"
          @click="taskDetail(scope.$index, scope.row)"
        >详情
        </el-button>
        <el-button
          type="danger"
          @click="deleteTask(scope.$index, scope.row)"
        >删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import { getTaskList, deleteTask } from '@/api/task'
import moment from 'moment'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        UP: 'success',
        creating: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      podData: [],
      taskData: [],
      search: '',
      loading: false,
      interval: undefined
    }
  },
  mounted() {
    this.update()
    this.interval = setInterval(this.update, 5000)
  },
  beforeDestroy() {
    clearInterval(this.interval)
  },
  methods: {
    update() {
      this.loading = true
      getTaskList().then(res => {
        console.log(res)
        if (res.data.code === 200) {
          this.taskData = res.data.data
          console.log(this.taskData)
          for (let i = 0; i < this.taskData.length; i++) {
            this.taskData[i].CreatedAt = this.formatTime(this.taskData[i].CreatedAt)
          }
        }
        this.loading = false
      })
    },
    change(e) {
      this.$forceUpdate()
      console.log(e)
    },
    // 日期UTC格式化
    formatTime(value) {
      if (value === '') {
        return ''
      }
      return moment(value).format('YYYY-MM-DD HH:mm:ss')
    },
    taskDetail(index, row) {
      console.log(index, row)
      this.$router.push({ path: '/task/view', query: { taskid: row.TaskId }})
    },
    deleteTask(index, row) {
      this.$confirm('是否删除任务，删除后不可恢复。', '删除任务?', { confirmButtonText: '删除', cancelButtonText: '取消' }).then(() => {
        deleteTask({ task_id: row.TaskId }).then(() => { this.update() })
      }).catch(() => {})

      // console.log(index, row)
    }
  }
}
</script>
