/* eslint-disable */
<template>
  <el-card class="box-card">
    <div slot="header" class="clearfix">
      <span>编排流程图</span>
      <!-- <el-button style="float: right; padding: 3px 0" type="text" @click="openImport">导入</el-button>
    <el-button style="float: right; padding: 3px 0" type="text" @click="exportJson" >导出</el-button> -->
    </div>
    <el-row>
      <div class="super-flow-base-demo">
        <!-- <div class="node-container">
          <span
              class="node-item"
              v-for="item in nodeItemList"
              :key="item.label"
              @mousedown="evt => nodeItemMouseDown(evt, item.value)">
            {{ item.label }}
          </span>
        </div> -->
        <div ref="flowContainer" class="flow-container">
          <super-flow
            ref="superFlow"
            :has-mark-line="false"
            :node-list="nodeList"
            :link-list="linkList"
            :origin="origin"
            :link-addable="false"
            :link-editabel="false"
            :graph-menu="graphMenuList"
            :node-menu="nodeMenuList"
            :link-menu="linkMenuList"
          >
            <template v-slot:node="{ meta, node }">
              <div
                :class="`flow-node flow-node-${nodesettings[meta.id].prop}`"
                @dblclick="openConfig(node)" @mouseover="testhover(node)" @mouseout="clearhover"
              >
                <header>
                  {{ nodesettings[meta.id].name }}
                </header>
                <div class="desc">
                  {{ nodesettings[meta.id].desc }}
                </div>

                <div class="status">
                  <el-col :span="12">

                    <div v-if="! result[node.id]">
                      加载运行状态...
                    </div>
                    <div v-else>
                      <span v-if="result[node.id].ready">✅工作正常</span>
                      <span v-else>{{ result[node.id].state }} {{ result[node.id] }}</span>
                    </div>

                    <!-- <Promised :promised="getstatus(node)">
                      <template v-slot:combined="{ isPending, isDelayElapsed, data, error }">
                        <span>✅工作正常 {{error}} {{data}}</span>
                      </template>
                    </Promised> -->
                  </el-col>
                  <el-col :span="12">
                    <span style="float: right">节点:{{ nodeData[node.id].node }}</span>
                  </el-col>
                </div>
              </div>
            </template>
          </super-flow>
        </div>
        <el-drawer
          title="查看节点属性"
          size="100%"
          :visible.sync="drawer"
          direction="rtl"
          :before-close="handleClose"
          :destroy-on-close="true"
        >
          <div style="width: 100%; padding: 5px">
            <el-row>
              <!--            <el-col :span="8">-->
              <!--              <el-card v-if="Object.keys(tmpschema).length != 0">-->
              <!--                <el-descriptions-->
              <!--                  style="margin-top: 0px"-->
              <!--                  title="配置"-->
              <!--                  :column="1"-->
              <!--                  border-->
              <!--                >-->
              <!--                  <el-descriptions-item-->
              <!--                    v-for="(v, k) in tmpschema"-->
              <!--                    :key="k"-->
              <!--                    :label="v.name"-->
              <!--                  >-->
              <!--                    {{ tmpdata[k] }}-->
              <!--                  </el-descriptions-item>-->
              <!--                </el-descriptions>-->
              <!--              </el-card>-->
              <!--            </el-col>-->
              <el-col :span="8">
                <el-card v-if="mapper_needed">
                  <el-descriptions
                    style="margin-top: 0px"
                    title="映射"
                    :column="1"
                    border
                  >
                    <el-descriptions-item
                      v-for="item in tmp_mapper_info"
                      :key="item.key"
                      :label="item.name"
                    >
                      <div
                        v-if="
                        tmpdata.mapper[item.key].type == 'default' ||
                          tmpdata.mapper[item.key].type == 'custom'
                      "
                      >
                        {{ tmpdata.mapper[item.key].type }}({{
                          tmpdata.mapper[item.key].value
                        }})
                      </div>
                      <div v-else>{{ tmpdata.mapper[item.key].type }}</div>
                    </el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-col>
              <el-col :span="8">
                <el-card style="height: 150px">
                  <el-descriptions
                    style="margin-top: 0"
                    title="运行节点"
                    :column="1"
                    border
                  >
                    <template slot="extra">
                      <el-button type="primary" size="small" @click="migrateConfig" v-if="isStart"> 迁移节点</el-button>
                      <el-button type="primary" size="small" @click="nodeStrategyConfig" v-if="!isStart"> 更换编排策略</el-button>
                    </template>
                    <el-descriptions-item label="节点">
                      {{ tmpdata.node }}
                    </el-descriptions-item>
                  </el-descriptions>
                  <div>
                    <el-drawer
                      title="节点编排"
                      :append-to-body="true"
                      :before-close="innerHandleClose"
                      :visible.sync="innerDrawer"
                    >
                      <el-form :model="tmpdata" label-width="100px">
                        <el-form-item label="选择端:">
                          <el-radio v-model="tmpdata.nodeStrategy.nodeType" label="cloud">云端</el-radio>
                          <el-radio v-model="tmpdata.nodeStrategy.nodeType" label="edge">边缘端</el-radio>
                          <el-radio v-model="tmpdata.nodeStrategy.nodeType" label="random">随机</el-radio>
                        </el-form-item>
                        <el-form-item v-if="tmpdata.nodeStrategy.nodeType === 'cloud'" label="GPU:">
                          <el-switch
                            v-model="tmpdata.nodeStrategy.isGPU"
                            active-text="使用"
                            inactive-text="不使用"
                          />
                        </el-form-item>
                        <el-form-item v-else-if="tmpdata.nodeStrategy.nodeType === 'edge'" label="厂区:">
                          <el-select v-model="tmpdata.nodeStrategy.nodeArea" placeholder="请选择">
                            <el-option
                              v-for="item in area"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                        </el-form-item>
                        <div
                          v-if="tmpdata.nodeStrategy.nodeType !=='random'&& tmpdata.nodeStrategy.nodeArea!=='random'"
                        >
                          <el-form-item label="运行节点:">
                            <el-switch
                              v-model="tmpdata.nodeStrategy.isNodeRandom"
                              active-text="随机"
                              inactive-text="自选"
                            />
                          </el-form-item>
                          <el-form-item v-if="!tmpdata.nodeStrategy.isNodeRandom" label="运行节点:">
                            <el-select v-model="tmpdata.nodeStrategy.nodePrefer" :filterable="true" clearable>
                              <el-option
                                v-for="node in nodeFilter(node_info,tmpdata.nodeStrategy.nodeType,tmpdata.nodeStrategy.nodeArea,tmpdata.nodeStrategy.isGPU)"
                                :key="node.key" :value="node.key"
                              >
                                <span>{{ node.key }}</span>
                                <el-tag v-for="tag in node.tag" :key="tag" size="mini">{{ tag }}</el-tag>
                                <i :class="'iconfont icon-a-xinhaowangluo-'+node.status"/>
                              </el-option>
                            </el-select>
                          </el-form-item>
                        </div>
                        <el-form-item style="margin-left: 100px">
                          <el-button @click="saveNodeStrategy">确认</el-button>
                        </el-form-item>
                      </el-form>
                    </el-drawer>
                  </div>
                </el-card>
              </el-col>
              <el-col :span="8">
                <el-card style="height: 150px">
                  <el-descriptions
                    style="margin-top: 0"
                    title="使用镜像"
                    :column="1"
                    border
                  >
                    <template slot="extra">
                      <el-button v-if="isUpdateImage" type="primary" size="small" @click="updateImageConfig"> 更换镜像
                      </el-button>
                    </template>
                    <el-descriptions-item label="镜像">
                      {{ podImage }}
                    </el-descriptions-item>
                  </el-descriptions>
                </el-card>
              </el-col>
            </el-row>

            <el-row>
              <el-tabs v-model="activeName" type="border-card">
                <el-tab-pane :lazy="true" label="实时监控" name="monitor">
                  <el-row>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=9&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=10&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=11&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                  </el-row>
                  <el-row>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=12&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=13&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="8">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=14&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                  </el-row>
                </el-tab-pane>
                <el-tab-pane :lazy="true" label="历史数据" name="history">
                  <el-row>
                    <el-col :span="12">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=2&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="12">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=3&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                  </el-row>
                  <el-row>
                    <el-col :span="12">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=4&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                    <el-col :span="12">
                      <iframe
                        :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=5&var-pod=${podname}`"
                        width="100%" height="100%" frameborder="0"
                      />
                    </el-col>
                  </el-row>
                </el-tab-pane>
                <el-tab-pane :lazy="true" label="日志查看" name="log">
                  <iframe
                    :src="`${grafana_base}d-solo/linjDM44k/test?orgId=1&from=now-1h&to=now&theme=light&panelId=20&var-taskid=${taskid}&var-podname=${pod}`"
                    width="800" height="200" frameborder="0"
                  />
                </el-tab-pane>
              </el-tabs>
            </el-row>
          </div>
        </el-drawer>

        <el-dialog :visible.sync="nodeDialog" title="选择迁移节点">
          <el-select v-model="new_node">
            <el-option v-for="(nodeprop,node) in node_info" v-if="node!==tmpdata.node" :key="node" :value="node">
              <span>{{ node }}</span>
              <el-tag v-for="tag in nodeprop.tag" :key="tag" size="mini">{{ tag }}</el-tag>
              <i :class="'iconfont icon-a-xinhaowangluo-'+nodeprop.status"/>
            </el-option>
          </el-select>
          <el-button @click="migrate">迁移</el-button>
        </el-dialog>

        <el-dialog :visible.sync="imageDialog" title="搜索镜像">
          <el-col :span="14">
            <el-input v-model="newPodImage" placeholder="请搜索镜像" :disabled="true"/>
          </el-col>
          <el-col :span="6"/>
          <el-col :span="4">
            <el-popover v-model="searchvisible" placement="bottom" title="搜索镜像仓库" width="900" trigger="click">
              <vue-simple-suggest
                v-model="searchchoice"
                :list="searchFunc"
                mode="select"
                @select="onSuggestSelect"
              >
                <input v-model="searchtext">
                <div slot="suggestion-item" slot-scope="{ suggestion }" class="custom">{{ suggestion.title }}
                </div>
              </vue-simple-suggest>
              <el-button slot="reference">搜索</el-button>
            </el-popover>
          </el-col>
          <el-button @click="updateImage">更换</el-button>
        </el-dialog>
      </div>
    </el-row>
    <el-row>
      <Topology :task="taskfortopology" :highlight="highlight"></Topology>
    </el-row>
  </el-card>
</template>

<script>
import {
  getBlockSettings,
  getDeviceMapper,
  getNodes,
  getTask,
  getPodName,
  getPodStatus,
  podMigration,
  searchImage,
  getPodHostName,
  updatePodImage,
  NodeAffinityMigration
} from '@/api/task'
import Topology from '@/components/Topology/Task'
import VueSimpleSuggest from 'vue-simple-suggest'
import 'vue-simple-suggest/dist/styles.css'
import { nodeAffinityParse } from '@/utils/nodeAffinity' // Using a css-loader

export default {
  components: { Topology, VueSimpleSuggest },
  data() {
    return {
      area: [{
        value: 'huadu',
        label: '花都工厂'
      }, {
        value: 'huangpu',
        label: '黄埔创新中心'
      }, {
        value: 'random',
        label: '随机'
      }],
      isStart: false,
      innerDrawer: false,
      isUpdateImage: false,
      searchvisible: false,
      searchtext: '',
      searchchoice: '',
      podImage: '',
      newPodImage: '',
      taskfortopology: {},
      highlight: '',
      dragConf: {
        isDown: false,
        isMove: false,
        offsetTop: 0,
        offsetLeft: 0,
        clientX: 0,
        clientY: 0,
        ele: null,
        info: null
      },
      node_info: {},
      result: {},
      interval: {},
      taskid: '',
      podname: 'model-5c3b97cfa75b4cc0be5d679f8112aa78-86c834e59fa644daa73dql87',
      grafana_base: 'http://10.112.134.196:30505/',
      pod: '',
      nodesettings: {},
      podid: '',
      mapper_needed: false,
      activeName: 'monitor',
      mapper_type: 'http',
      origin: [681, 465],
      nodeList: [],
      linkList: [],
      nodeData: {},
      drawer: false,
      tmpschema: {},
      tmpdata: {
        node: '',
        nodeStrategy: {
          isGPU: false,
          nodeArea: '',
          isNodeRandom: false,
          nodeType: 'random',
          nodePrefer: ''
        },
        nodeAffinity: ''
      },
      new_node: '',
      tmpnodename: '',
      nodeDialog: false,
      imageDialog: false,
      importtext: '',
      nodeItemList: [],
      graphMenuList: [],
      mapper_info: {},
      tmp_mapper_info: [],
      nodeMenuList: [],
      linkMenuList: []
    }
  },
  mounted() {
  },
  async created() {
    this.taskid = this.$route.query.taskid
    let nodesettings = await getBlockSettings()
    const devicemapper = await getDeviceMapper()

    nodesettings = nodesettings.data
    this.nodesettings = nodesettings
    this.mapper_info = devicemapper.data
    this.node_info = await getNodes()
    // this.importtext = JSON.stringify(require("./test.json"));
    await this.importJson()
    this.taskReady = 1
  },
  beforeDestroy() {
    this.nodeList.map(this.cancelInterval)
  },
  methods: {
    async searchFunc() {
      console.log(this.searchtext)
      const result = await searchImage(this.searchtext)
      console.log(result)
      return result.data.data
    },
    onSuggestSelect() {
      console.log(this.searchchoice)
      this.newPodImage = this.searchchoice.image
      console.log(JSON.stringify(this.tmpdata))
      this.searchvisible = false
    },
    testhover(e) {
      console.log(e)
      this.highlight = e.id
    },
    clearhover() {
      this.highlight = ''
    },
    migrateConfig() {
      this.nodeDialog = true
    },
    nodeStrategyConfig() {
      this.innerDrawer = true
      console.log(this.tmpdata)
    },
    nodeFilter(node_info, nodeType, nodeArea, isGPU) {
      // this.tmpdata.nodeStrategy.nodePrefer = ''
      if (nodeType === 'edge') {
        node_info = this.nodeTypeFilter(node_info, '边缘端')
        node_info = this.nodeAreaFilter(node_info, nodeArea)
      } else if (nodeType === 'cloud') {
        node_info = this.nodeTypeFilter(node_info, '云端')
        node_info = this.nodeGPUFilter(node_info, isGPU)
      }
      return node_info
    },
    nodeTypeFilter(node_info, nodeType) {
      const arr = []
      Object.keys(node_info).forEach(function(key) {
        let allObject = {}
        allObject = Object.assign({ key: key }, node_info[key])
        arr.push(allObject)
      })
      const newArr = []
      for (let i = 0; i < arr.length; i++) {
        if (arr[i].tag.length > 0) {
          for (let j = 0; j < arr[i].tag.length; j++) {
            if (arr[i].tag[j] === nodeType) {
              newArr.push(arr[i])
            }
          }
        }
      }
      return newArr
    },
    nodeAreaFilter(node_info, nodeArea) {
      const arr = []
      Object.keys(node_info).forEach(function(key) {
        let allObject = {}
        allObject = Object.assign({ key: key }, node_info[key])
        arr.push(allObject)
      })
      const newArr = []
      for (let i = 0; i < arr.length; i++) {
        if (arr[i].tag.length > 0) {
          if (arr[i].area === nodeArea) {
            newArr.push(arr[i])
          }
        }
      }
      return newArr
    },
    nodeGPUFilter(node_info, isGPU) {
      let flag = ''
      isGPU ? flag = 'GPU' : flag = 'CPU'
      const arr = []
      Object.keys(node_info).forEach(function(key) {
        let allObject = {}
        allObject = Object.assign({ key: key }, node_info[key])
        arr.push(allObject)
      })
      const newArr = []
      for (let i = 0; i < arr.length; i++) {
        if (arr[i].tag.length > 0) {
          for (let j = 0; j < arr[i].tag.length; j++) {
            if (arr[i].tag[j] === flag) {
              newArr.push(arr[i])
            }
          }
        }
      }
      return newArr
    },
    innerHandleClose(done) {
      this.$confirm('确认关闭？')
        .then(_ => {
          done()
        })
        .catch(_ => {
        })
    },
    updateImageConfig() {
      this.imageDialog = true
    },
    migrate() {
      this.$confirm('确认要迁移运行节点？').then(() => {
        this.nodeData[this.podid].node = this.new_node
        let loading = this.$loading({ fullscreen: true })
        podMigration({
          definition: JSON.stringify(this.exportJson()),
          task_id: this.taskid,
          pod_id: this.podid,
          target_node: this.new_node
        }).then((res) => {
          if (res.data.code === 0) {
            console.log(res.data)
            this.$message({ message: '迁移成功', type: 'success' })
            this.importJson()
          } else {
            this.$message.error(res.data.msg)
          }
          console.log(res.data)
          this.nodeDialog = false
          this.drawer = false
          console.log(res.data)
          loading.close()
        }).catch(() => {
          this.nodeDialog = false
          this.$message.error('失败')
          loading.close()
        })
        // this.exportJson()
      }).catch(() => {
        this.nodeDialog = false
      })
    },
    updateImage() {
      console.log(this.newPodImage)
      console.log(this.nodeData[this.podid])
      this.$confirm('确认要更新镜像吗？').then(() => {
        this.podImage = this.newPodImage
        let loading = this.$loading({ fullscreen: true })
        updatePodImage({
          task_id: this.taskid,
          pod_id: this.podid,
          new_image: this.newPodImage,
          definition: JSON.stringify(this.exportJson())
        }).then((res) => {
          if (res.data.code === 0) {
            console.log(res.data)
            this.$message({ message: '更新成功', type: 'success' })
            this.importJson()
          } else {
            this.$message.error(res.data.msg)
          }
          console.log(res.data)
          this.imageDialog = false
          this.nodeDialog = false
          this.drawer = false
          console.log(res.data)
          loading.close()
        }).catch(() => {
          this.nodeDialog = false
          this.$message.error('失败')
          loading.close()
        })
        this.exportJson()
      }).catch(() => {
        this.nodeDialog = false
      })
    },
    openImport() {
      this.nodeDialog = true
    },
    async importJson() {
      const resp = await getTask(this.taskid)
      this.taskfortopology = resp.data.data
      console.log(this.taskfortopology)
      const input = JSON.parse(resp.data.data.definition)
      console.log(input.data)
      for (const i in input.data) {
        getPodHostName({ 'podid': i }).then((res) => {
          console.log(res.data.data)
          input.data[i]['node'] = res.data.data
        })
      }
      console.log(input.data)
      this.nodeList = input.graph.nodeList
      this.origin = input.graph.origin
      this.linkList = input.graph.linkList
      this.nodeData = input.data
      this.nodeDialog = false
      this.importtext = ''
      this.nodeList.map(this.creatSetInterval)
      this.nodeList.map(this.getstatus)
    },
    handleClose(done) {
      this.isUpdateImage = false
      done()
    },
    creatSetInterval(node) {
      this.interval[node.id] = setInterval(() => {
        this.getstatus(node)
      }, 300 * 1000)
    },
    cancelInterval(node) {
      clearInterval(this.interval[node.id])
    },
    async getpodname(podid, taskid) {
      const resp = await getPodName({ podid: podid, taskid: taskid })
      return resp.data.Data
    },
    getstatus(node) {
      this.$set(this.result, node.id, undefined)
      this.result[node.id] = undefined
      this.getpodname(node.id, this.taskid).then((result) => {
        return getPodStatus({ podname: result })
      }).then((resp) => {
        console.log(resp.data)
        // this.$set(this.result,node.id,resp.data.Data)
        this.$set(this.result, node.id, resp.data.Data)
        console.log(this.result)
        // this.$forceUpdate()
      })
    },
    async openConfig(node) {
      this.podid = node.id
      this.activeName = 'monitor'
      this.podname = await this.getpodname(node.id, this.taskid)
      this.podImage = this.result[node.id].image
      this.pod = node.meta.id
      console.log(this.pod)
      if (this.pod === 'model') {
        this.isUpdateImage = true
      }
      if (this.pod === 'intake-http-form' || this.pod === 'intake-http-json') {
        this.isStart = true
      }
      var tmpdata = {
        node: ''
      }
      this.tmpnodename = node.id
      var tmpnodemeta = this.nodesettings[node.meta.id]
      this.tmpschema = this.nodesettings[node.meta.id].schema
      // if (tmpnodemeta.name.includes('Kafka')) {
      //   this.mapper_type = 'kafka'
      // } else {
      //   this.mapper_type = 'http'
      // }
      console.log(tmpnodemeta)
      if (tmpnodemeta.prop === 'start') {
        this.mapper_needed = true
        this.tmp_mapper_info = this.mapper_info[node.meta.id]
        tmpdata.mapper = {}
        const mapper = tmpdata.mapper
        for (const k in this.tmp_mapper_info) {
          mapper[this.tmp_mapper_info[k].key] = { type: '', value: '' }
        }
      } else {
        this.mapper_needed = false
      }
      if (node.id in this.nodeData) {
        for (const k in this.nodeData[node.id]) {
          tmpdata[k] = this.nodeData[node.id][k]
          if (k in this.tmpschema) {
            if (this.tmpschema[k].type === 'choice') {
              for (const index in this.tmpschema[k].choices) {
                console.log(
                  tmpdata[k],
                  JSON.stringify(this.tmpschema[k].choices[index])
                )
                if (tmpdata[k] === this.tmpschema[k].choices[index].val) {
                  tmpdata[k] = this.tmpschema[k].choices[index].label
                }
              }
            }
          }
        }
      } else {
        for (const k in this.tmpschema) {
          if (this.tmpschema[k].type === 'input') {
            tmpdata[k] = ''
          } else {
            tmpdata[k] = this.tmpschema[k].value
          }
        }
      }
      this.tmpdata = tmpdata
      console.log(this.tmpdata)
      this.drawer = true
    },
    save() {
      this.nodeData[this.tmpnodename] = {}
      console.log(JSON.stringify(this.tmpdata))
      for (var k in this.tmpdata) {
        this.nodeData[this.tmpnodename][k] = this.tmpdata[k]
      }
      this.tmpschema = {}
      this.mapper_needed = false
      this.tmpdata = {}
      this.tmpnodename = ''
      this.drawer = false
    },
    saveNodeStrategy() {
      console.log(this.tmpdata)
      this.tmpdata['nodeAffinity'] = nodeAffinityParse(this.tmpdata.nodeStrategy, '')
      console.log(this.tmpdata)
      this.$confirm('确认要更改节点策略吗？').then(() => {
        this.nodeData[this.podid].nodeAffinity = this.tmpdata['nodeAffinity']
        let loading = this.$loading({ fullscreen: true })
        NodeAffinityMigration({
          task_id: this.taskid,
          pod_id: this.podid,
          new_nodeAffinity: this.tmpdata['nodeAffinity'],
          definition: JSON.stringify(this.exportJson())
        }).then((res) => {
          if (res.data.code === 0) {
            console.log(res.data)
            this.$message({ message: '更换成功', type: 'success' })
            this.importJson()
          } else {
            this.$message.error(res.data.msg)
          }
          console.log(res.data)
          this.imageDialog = false
          this.nodeDialog = false
          this.drawer = false
          console.log(res.data)
          loading.close()
        }).catch(() => {
          this.nodeDialog = false
          this.$message.error('失败')
          loading.close()
        })
        this.exportJson()
      }).catch(() => {
        this.nodeDialog = false
      })
      this.innerDrawer = false
    },
    exportJson() {
      var result = {
        graph: this.$refs.superFlow.toJSON(),
        data: this.nodeData
      }
      console.log(JSON.stringify(result))
      return result
    }
  }
}
</script>

<style lang="less">
.node-item {
  @node-item-height: 30px;

  font-size: 14px;
  display: inline-block;
  height: @node-item-height;
  width: 120px;
  margin-top: 20px;
  background-color: #ffffff;
  line-height: 15px;
  box-shadow: 1px 1px 4px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  user-select: none;
  text-align: center;
  z-index: 6;

  &:hover {
    box-shadow: 1px 1px 8px rgba(0, 0, 0, 0.4);
  }
}

.status {
  font-size: 10px;
}

.super-flow-base-demo {
  width: 100%;
  height: 550px;
  margin: 0 auto;
  background-color: #f5f5f5;
  @list-width: 200px;

  > .node-container {
    width: @list-width;
    float: left;
    height: 100%;
    text-align: center;
    background-color: #ffffff;
  }

  > .flow-container {
    width: 100%;
    float: left;
    height: 100%;
    overflow: hidden;
  }

  .super-flow__node {
    display: flex;

    .flow-node {
      flex-direction: column;
      height: 100%;
      width: 100%;

      > header {
        font-size: 12px;
        height: 16px;
        line-height: 16px;
        padding: 0 12px;
        color: #ffffff;
        text-align: center;
      }

      .desc {
        text-align: center;
        line-height: 12px;
        font-size: 10px;
        overflow: hidden;
        padding: 6px 12px;
        word-break: break-all;
        height: 32px;
      }

      &.flow-node-start {
        > header {
          background-color: #55abfc;
        }
      }

      &.flow-node-process {
        > header {
          background-color: #bc1d16;
        }
      }

      &.flow-node-model {
        > header {
          background-color: rgba(188, 181, 58, 0.76);
        }
      }

      &.flow-node-output {
        > header {
          background-color: #30b95c;
        }
      }
    }
  }
}
</style>
