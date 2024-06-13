<template>
  <div>
    <div style="display: flex;justify-content:space-between">
      <el-radio-group v-model="isPublic" style="margin-top: 25px;margin-left: 30px" @change="changeIsPublic">
        <el-radio-button :label="true">公有镜像</el-radio-button>
        <el-radio-button :label="false">我的镜像</el-radio-button>
      </el-radio-group>
      <el-input
        v-model="search"
        style="margin-right: 30px"
        placeholder="请输入关键词搜索"
        clearable
        class="search"
      />
    </div>
    <el-divider />
    <el-row v-model="isPublic" style="margin-left: 40px;margin-right: 40px">
      <el-col
        v-for="item in imageData.filter(data => !search || data.title.toLowerCase().includes(search.toLowerCase()))"
        :key="item.id"
        :span="5"
        style="margin:20px 30px 30px 10px"
      >
        <el-card class="imageCard">
          <div slot="header" class="clearfix">
            <span>{{ item.title }}</span>
            <el-button v-if="!isPublic" style="float: right; padding: 3px 0" type="text" @click="edit(item)">编辑</el-button>
          </div>
          <div v-if="!isPublic && item.isPublic === '1'" class="text item " style="color: #2ac06d">
            公有
          </div>
          <div v-if="!isPublic && item.isPublic === '0'" class="text item " style="color: #dd1100">
            私有
          </div>
          <div class="text item">
            描述：{{ item.desc }}
          </div>
          <div class="text item ">
            镜像：{{ item.image }}
          </div>
        </el-card>
      </el-col>

      <el-drawer
        ref="drawer"
        title="镜像编辑"
        :before-close="handleClose"
        :visible.sync="dialog"
        direction="rtl"
        custom-class="demo-drawer"
      >
        <div class="demo-drawer__content">
          <el-form :model="form">
            <el-form-item label="镜像名称" label-width="80px">
              <el-input v-model="form.title" autocomplete="off" :value="form.title" />
            </el-form-item>
            <el-form-item label="可见范围" label-width="80px">
              <el-switch
                v-model="form.isPublic"
                active-text="公有"
                inactive-text="私有"
              />
            </el-form-item>
            <el-form-item label="镜像描述" label-width="80px">
              <el-input v-model="form.desc" autocomplete="off" :value="form.desc" />
            </el-form-item>
          </el-form>
          <div style="display: flex;justify-content:center;margin-top: 80px">
            <el-button style="margin-right: 20px" @click="cancelForm">取 消</el-button>
            <el-button type="primary" :loading="loading" @click="$refs.drawer.closeDrawer()">
              {{ loading ? '提交中 ...' : '确 定' }}
            </el-button>
          </div>
        </div>
      </el-drawer>
    </el-row>
  </div>
</template>

<script>
import { getIaiImageList, editIaiImageRecord } from '@/api/iaiImage'
import moment from 'moment'

export default {
  data() {
    return {
      imageData: [],
      search: '',
      isPublic: true,
      dialog: false,
      loading: false,
      form: {
        id: '',
        userId: '',
        images: '',
        title: '',
        isPublic: '1',
        desc: ''
      },
      timer: null
    }
  },
  mounted() {
    getIaiImageList(true).then(res => {
      console.log(res)
      if (res.data.code === 0) {
        this.imageData = res.data.data
      }
    })
  },
  methods: {
    changeIsPublic(e) {
      console.log(e)
      getIaiImageList(e).then(res => {
        console.log(res)
        if (res.data.code === 0) {
          this.imageData = res.data.data
        }
      })
    },
    edit(item) {
      this.form.title = item.title
      this.form.desc = item.desc
      this.form.id = item.id
      item.isPublic === '1' ? this.form.isPublic = true : this.form.isPublic = false
      this.dialog = true
    },
    // 日期UTC格式化
    formatTime(value) {
      if (value === '') {
        return ''
      }
      return moment(value).format('YYYY-MM-DD HH:mm:ss')
    },

    handleClose(done) {
      if (this.loading) {
        return
      }
      this.$confirm('确定要更新镜像信息吗？')
        .then(_ => {
          this.loading = true
          console.log(this.form)
          const data = {
            'title': this.form.title,
            'desc': this.form.desc,
            'id': this.form.id,
            'isPublic': this.form.isPublic
          }
          console.log(data)
          editIaiImageRecord(data).then(res => {
            console.log(res)
          })
          this.timer = setTimeout(() => {
            done()
            // 动画关闭需要一定的时间
            setTimeout(() => {
              this.loading = false
            }, 400)
          }, 2000)
          location.reload()
        })
        .catch(_ => {
        })
    },
    cancelForm() {
      this.loading = false
      this.dialog = false
      clearTimeout(this.timer)
    }
  }
}
</script>

<style>
.text {
  font-size: 13px;
  color: #999;
}

.item {
  margin-bottom: 18px;
}

.search {
  margin-left: 20px;
  margin-top: 20px;
  display: flex;
  width: 400px;
}

.imageCard {
  border-radius: 15px;
}
</style>
