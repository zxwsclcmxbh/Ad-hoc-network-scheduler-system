/* eslint-disable */
<template>
  <div>
    <el-steps :active="active" finish-status="success">
      <el-step title="填写基本信息" />
      <el-step title="绘制编排流程图" />
      <el-step title="编排结果" />
    </el-steps>
    <div v-if="active===1">
      <el-form v-model="task_desc" label-width="100px">
        <el-form-item label="项目名称">
          <el-input v-model="task_desc.name" />
        </el-form-item>
        <el-form-item label="项目描述">
          <el-input v-model="task_desc.desc" type="textarea" :rows="4" placeholder="请输入内容" />
        </el-form-item>
        <el-form-item>
          <el-button @click="next_step">下一步</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card v-else-if="active===2">
      <div slot="header" class="clearfix">
        <span>编排流程图</span>
        <el-button style="float:right" @click="exportJson">提交编排</el-button>
        <el-button style="float:right" @click="prev_step">上一步</el-button>
      </div>
      <div class="super-flow-base-demo">
        <div class="node-container">
          <span
            v-for="item in nodeItemList"
            :key="item.label"
            class="node-item"
            @mousedown="evt => nodeItemMouseDown(evt, item.value)"
          >
            {{ item.label }}
          </span>
        </div>
        <div ref="flowContainer" class="flow">
          <super-flow
            ref="superFlow"
            :has-mark-line="false"
            :node-list="nodeList"
            :link-list="linkList"
            :origin="origin"
            :graph-menu="graphMenuList"
            :node-menu="nodeMenuList"
            :link-menu="linkMenuList"
            :enter-intercept="enterIntercept"
            :output-intercept="outputIntercept"
          >
            <template v-slot:node="{meta,node}">
              <div :class="`flow-node flow-node-${nodesettings[meta.id].prop}`" @dblclick="openConfig(node)">
                <header>
                  {{ nodesettings[meta.id].name }}
                </header>
                <section>
                  {{ nodesettings[meta.id].desc }}
                </section>
              </div>
            </template>
          </super-flow>
        </div>
        <el-drawer title="编辑节点属性" size="40%" :visible.sync="drawer" direction="rtl" :before-close="handleClose">
          <div style="width:100%; padding:5px">
            <el-form :model="tmpdata" label-width="100px">
              <el-form-item v-for="(v,k) in tmpschema" :key="k" :label="v.name">
                <div v-if="v.type==='input'">
                  <el-input v-model="tmpdata[k]" :placeholder="v.placeholder" />
                </div>
                <div v-else-if="v.type==='model'">
                  <el-col :span="14">
                    <el-input v-model="tmpdata[k]" :placeholder="v.placeholder" :disabled="true" />
                  </el-col>
                  <el-col :span="6" />
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
                </div>
                <div v-else>
                  <el-select v-model="tmpdata[k]" :placeholder="v.placeholder">
                    <el-option v-for="item in v.choices" :key="item.label" :label="item.label" :value="item.val" />
                  </el-select>
                </div>
              </el-form-item>
              <div v-if="mapper_needed">
                <el-form-item v-for="item in tmp_mapper_info" :key="item.key" :label="item.name">
                  <el-col :span="10">
                    <el-select v-model="tmpdata.mapper[item.key].type">
                      <el-option v-for="i in item.options" :key="i.value" :label="i.label" :value="i.value" />
                    </el-select>
                  </el-col>
                  <el-col :span="2" />
                  <el-col
                    v-if="tmpdata.mapper[item.key].type==='default' || tmpdata.mapper[item.key].type==='custom'"
                    :span="12"
                  >
                    <el-input v-model="tmpdata.mapper[item.key].value" />
                  </el-col>
                </el-form-item>
              </div>
              <el-form-item v-if="!isForceSelectEdgeNode" label="节点编排">
                <el-button type="primary" plain @click="nodeCombination">点击进行设置</el-button>
                <el-drawer
                  title="节点编排"
                  :append-to-body="true"
                  :before-close="innerHandleClose"
                  :visible.sync="innerDrawer"
                >
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
                  <div v-if="tmpdata.nodeStrategy.nodeType !=='random'&& tmpdata.nodeStrategy.nodeArea!=='random'">
                    <el-form-item label="运行节点:">
                      <el-switch
                        v-model="tmpdata.nodeStrategy.isNodeRandom"
                        active-text="随机"
                        inactive-text="自选"
                      />
                    </el-form-item>
                    <el-form-item v-if="!tmpdata.nodeStrategy.isNodeRandom" label="运行节点:">
                      <el-select v-model="tmpdata.nodeStrategy.nodePrefer" :filterable="true" clearable>
                        <el-option v-for="node in nodeFilter(node_info,tmpdata.nodeStrategy.nodeType,tmpdata.nodeStrategy.nodeArea,tmpdata.nodeStrategy.isGPU)" :key="node.key" :value="node.key">
                          <span>{{ node.key }}</span>
                          <el-tag v-for="tag in node.tag" :key="tag" size="mini">{{ tag }}</el-tag>
                          <i :class="'iconfont icon-a-xinhaowangluo-'+node.status" />
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </div>
                  <el-form-item style="margin-left: 100px">
                    <el-button @click="saveNodeCombination">确认</el-button>
                  </el-form-item>
                </el-drawer>
              </el-form-item>
              <el-form-item v-if="isForceSelectEdgeNode" label="运行节点">
                <el-select v-model="tmpdata.node" :filterable="true" clearable>
                  <el-option v-for="node in nodeTypeFilter(node_info,'边缘端')" :key="node.key" :value="node.key">
                    <span>{{ node.key }}</span>
                    <el-tag v-for="tag in node.tag" :key="tag" size="mini">{{ tag }}</el-tag>
                    <i :class="'iconfont icon-a-xinhaowangluo-'+node.status" />
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button @click="save">确认</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-drawer>
      </div>
    </el-card>
    <div v-else>
      <el-result icon="success" title="部署成功" sub-title="稍后可查看部署状态">
        <template slot="extra">
          <el-button type="primary" size="medium">返回</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script>
import VueSimpleSuggest from 'vue-simple-suggest'
import 'vue-simple-suggest/dist/styles.css' // Using a css-loader
import { getBlockSettings, getDeviceMapper, getNodes, searchImage, submitTask } from '@/api/task'
import { nodeAffinityParse } from '@/utils/nodeAffinity'
export default {
  components: {
    VueSimpleSuggest
  },
  data() {
    return {
      // isGPU: false,
      // nodeArea: '',
      // isNodeRandom: false,
      // nodeType: 'random',
      // nodePrefer: '',
      tempNodeAffinity: {},
      active: 1,
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
      searchvisible: false,
      searchtext: '',
      searchchoice: '',
      nodesettings: {},
      node_info: {},
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
      search: false,
      mapper_needed: false,
      // mapper_type: 'http',
      origin: [681, 465],
      nodeList: [],
      linkList: [],
      nodeData: {},
      task_desc: {
        name: '',
        desc: ''
      },
      drawer: false,
      innerDrawer: false,
      isForceSelectEdgeNode: false,
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
      tmpnodename: '',
      dialog: false,
      importtext: '',
      nodeItemList: [],
      graphMenuList: [],
      mapper_info: {},
      tmp_mapper_info: [],
      nodeMenuList: [
        [
          {
            label: '编辑',
            selected: (node, coordinate) => {
              this.openConfig(node)
            }
          }
        ],
        [
          {
            label: '删除',
            disable: false,

            selected(node, coordinate) {
              node.remove()
            }
          }
        ]
      ],
      linkMenuList: [
        [
          {
            label: '删除',
            disable: false,
            selected: (link, coordinate) => {
              link.remove()
            }
          }
        ]
      ]
    }
  },
  mounted() {
    document.addEventListener('mousemove', this.docMousemove)
    document.addEventListener('mouseup', this.docMouseup)
    this.$once('hook:beforeDestroy', () => {
      document.removeEventListener('mousemove', this.docMousemove)
      document.removeEventListener('mouseup', this.docMouseup)
    })
  },
  async created() {
    const nodesettings = await getBlockSettings()
    // let nodesettings= require("./data.json")

    const devicemapper = await getDeviceMapper()
    // console.log(await getNodes())
    // nodesettings = nodesettings
    this.nodesettings = nodesettings.data
    this.mapper_info = devicemapper.data
    // console.log(this.nodesettings)
    // console.log(JSON.stringify(this.nodesettings),require("./data.json"))
    this.node_info = await getNodes()
    console.log(this.node_info)
    for (const id in this.nodesettings) {
      var i = this.nodesettings[id]
      const meta = {
        id: id
      }
      const value = {
        width: 120,
        height: 50,
        meta: meta
      }
      var temp = {}
      if (i.prop === 'start') {
        temp = {
          meta: meta,
          label: i.name,
          // disable(graph) {
          //   return !!graph.nodeList.find(node => nodesettings[node.meta.id].prop === 'start')
          // },
          selected: (graph, coordinate) => {
            // const start = graph.nodeList.find(node => nodesettings[node.meta.id].prop === 'start')
            // if (!start) {
            graph.addNode({
              ...value,
              coordinate: coordinate
            })
            // }
          }
        }
      } else if (i.prop === 'output') {
        temp = {
          label: i.name,
          meta: meta,
          disable(graph) {
            return !!graph.nodeList.find(node => nodesettings[node.meta.id].prop === 'output')
          },
          selected: (graph, coordinate) => {
            const start = graph.nodeList.find(node => nodesettings[node.meta.id].prop === 'output')
            if (!start) {
              graph.addNode({
                ...value,
                coordinate: coordinate
              })
            }
          }
        }
      } else {
        temp = {
          label: i.name,
          meta: meta,
          disable: false,
          selected: (graph, coordinate) => {
            graph.addNode({
              ...value,
              coordinate: coordinate
            })
          }
        }
      }
      this.graphMenuList.push([temp])
      this.nodeItemList.push({
        label: i.name,
        value: () => {
          return value
        }
      })
    }
    // console.log()
    // console.log(JSON.stringify(this.nodeItemList))
    // this.graphMenuList=[this.graphMenuList]
  },
  methods: {
    next_step() {
      this.active += 1
    },
    prev_step() {
      this.active -= 1
    },
    onSuggestSelect() {
      console.log(this.searchchoice)
      this.tmpdata.model = this.searchchoice.image
      console.log(JSON.stringify(this.tmpdata))
      this.searchvisible = false
    },
    async searchFunc() {
      console.log(this.searchtext)
      const result = await searchImage(this.searchtext)
      return result.data.data
    },
    openImport() {
      this.dialog = true
    },
    importJson() {
      const input = JSON.parse(this.importtext)
      this.nodeList = input.graph.nodeList
      this.origin = input.graph.origin
      this.linkList = input.graph.linkList
      this.nodeData = input.data
      this.dialog = false
      this.importtext = ''
    },
    handleClose(done) {
      this.$confirm('确认关闭？')
        .then(_ => {
          this.isForceSelectEdgeNode = false
          this.mapper_needed = false
          this.tmpschema = {}
          this.tmpdata = {
            node: '',
            nodeStrategy: {
              isGPU: false,
              nodeArea: '',
              isNodeRandom: false,
              nodeType: 'random',
              nodePrefer: ''
            },
            nodeAffinity: ''
          }
          this.tmpnodename = ''

          done()
        })
        .catch(_ => {
        })
    },
    innerHandleClose(done) {
      this.$confirm('确认关闭？')
        .then(_ => {
          done()
        })
        .catch(_ => {
        })
    },
    nodeCombination() {
      this.innerDrawer = true
      console.log(this.node_info)
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
    openConfig(node) {
      const tmpdata = {
        node: '',
        nodeStrategy: {
          isGPU: false,
          nodeArea: '',
          isNodeRandom: false,
          nodeType: 'random',
          nodePrefer: ''
        },
        nodeAffinity: ''
      }
      this.tmpnodename = node.id
      const tmpnodemeta = this.nodesettings[node.meta.id]
      this.tmpschema = this.nodesettings[node.meta.id].schema
      // if (tmpnodemeta.name.includes('Kafka')) {
      //   this.mapper_type = 'kafka'
      // } else {
      //   this.mapper_type = 'http'
      // }
      if (tmpnodemeta.prop === 'start') {
        this.isForceSelectEdgeNode = true
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
      this.drawer = true
    },
    save() {
      this.tmpdata['nodeAffinity'] = nodeAffinityParse(this.tmpdata.nodeStrategy, this.tmpdata.node)
      if (this.tmpdata['nodeAffinity'] === null) {
        this.$message.error('节点编排无效！')
      }
      this.nodeData[this.tmpnodename] = {}
      console.log(JSON.stringify(this.tmpdata))
      for (const k in this.tmpdata) {
        this.nodeData[this.tmpnodename][k] = this.tmpdata[k]
      }
      this.isForceSelectEdgeNode = false
      this.tmpschema = {}
      this.mapper_needed = false
      this.tmpdata = {
        node: '',
        nodeStrategy: {
          isGPU: false,
          nodeArea: '',
          isNodeRandom: false,
          nodeType: 'random',
          nodePrefer: ''
        },
        nodeAffinity: ''
      }
      console.log(this.tmpdata)
      this.tmpnodename = ''
      this.drawer = false
    },
    saveNodeCombination() {
      this.innerDrawer = false
    },
    enterIntercept(fromNode, toNode, graph) {
      console.log()
      const from = graph.linkList.filter(link => link.start.id === fromNode.id)
      const to = graph.linkList.filter(link => link._end.id === toNode.id)
      const duplicate = graph.linkList.filter(link => link._end.id === toNode.id && link.start.id === fromNode.id)
      console.log('出发节点', from.length, '接收节点', to.length, '重复节点', duplicate.length, '最大输出', this.nodesettings[fromNode.meta.id].max_output, '最大输入', this.nodesettings[toNode.meta.id].max_input)
      if (duplicate.length >= 1) return false
      if (this.nodesettings[fromNode.meta.id].max_output === -1) {
        return true
      } else {
        if (from.length >= this.nodesettings[fromNode.meta.id].max_output) return false
      }
      if (this.nodesettings[toNode.meta.id].max_input === -1) {
        return true
      } else {
        if (to.length >= this.nodesettings[toNode.meta.id].max_input) return false
      }

      const fromType = this.nodesettings[fromNode.meta.id].prop
      switch (this.nodesettings[toNode.meta.id].prop) {
        case 'start':
          return false
        case 'process':
          return [
            'start',
            'model',
            'process',
            'fork'
          ].includes(fromType)
        case 'model':
          return [
            'start',
            'process',
            'model',
            'fork'
          ].includes(fromType)
        case 'output':
          return [
            'start',
            'model',
            'process',
            'fork'
          ].includes(fromType)
        case 'fork':
          return [
            'start',
            'model',
            'process',
            'fork'
          ].includes(fromType)
        default:
          return true
      }
    },
    outputIntercept(node, graph) {
      return !(this.nodesettings[node.meta.id].prop === 'output')
    },
    // 提交编排结果
    exportJson() {
      var result = {
        graph: this.$refs.superFlow.toJSON(),
        data: this.nodeData
      }
      console.log(result)
      // this.$alert(JSON.stringify(result), '复制内容', {
      //     confirmButtonText: '确定',
      //   });
      this.$confirm('确认提交部署？')
        .then(_ => {
          var loading = this.$loading({ fullscreen: true })
          submitTask({
            'definition': JSON.stringify(result),
            'name': this.task_desc.name,
            'description': this.task_desc.desc
          }).then((resp) => {
            loading.close()
            if (resp.data.code === 200) {
              this.$message({ message: '部署成功', type: 'success' })
              this.active += 1
            } else {
              this.$message({ messgae: resp.data.msg, type: 'error' })
            }
          }
          ).catch((err) => {
            this.$message({ messgae: err, type: 'error' })
          })
        })
        .catch(_ => {
        })
    },
    docMousemove({ clientX, clientY }) {
      const conf = this.dragConf

      if (conf.isMove) {
        conf.ele.style.top = clientY - conf.offsetTop + 'px'
        conf.ele.style.left = clientX - conf.offsetLeft + 'px'
      } else if (conf.isDown) {
        // 鼠标移动量大于 5 时 移动状态生效
        conf.isMove =
          Math.abs(clientX - conf.clientX) > 5 ||
          Math.abs(clientY - conf.clientY) > 5
      }
    },
    docMouseup({ clientX, clientY }) {
      const conf = this.dragConf
      conf.isDown = false

      if (conf.isMove) {
        const {
          top,
          right,
          bottom,
          left
        } = this.$refs.flowContainer.getBoundingClientRect()

        // 判断鼠标是否进入 flow container
        if (
          clientX > left &&
          clientX < right &&
          clientY > top &&
          clientY < bottom
        ) {
          // 获取拖动元素左上角相对 super flow 区域原点坐标
          const coordinate = this.$refs.superFlow.getMouseCoordinate(
            clientX - conf.offsetLeft,
            clientY - conf.offsetTop
          )

          // 添加节点
          this.$refs.superFlow.addNode({
            coordinate,
            ...conf.info
          })
        }

        conf.isMove = false
      }

      if (conf.ele) {
        conf.ele.remove()
        conf.ele = null
      }
    },
    nodeItemMouseDown(evt, infoFun) {
      const {
        clientX,
        clientY,
        currentTarget
      } = evt

      const {
        top,
        left
      } = evt.currentTarget.getBoundingClientRect()

      const conf = this.dragConf
      const ele = currentTarget.cloneNode(true)
      const info = infoFun()
      console.log(info)
      console.log(this.$refs.superFlow.graph.nodeList.find(node => this.nodesettings[node.meta.id].prop === 'start'))
      if ((this.nodesettings[info.meta.id].prop === 'output' && this.$refs.superFlow.graph.nodeList.find(node => this.nodesettings[node.meta.id].prop === 'output'))) {
        return
      }
      Object.assign(this.dragConf, {
        offsetLeft: clientX - left,
        offsetTop: clientY - top,
        clientX: clientX,
        clientY: clientY,
        info: info,
        ele,
        isDown: true
      })

      ele.style.position = 'fixed'
      ele.style.margin = '0'
      ele.style.top = clientY - conf.offsetTop + 'px'
      ele.style.left = clientX - conf.offsetLeft + 'px'

      this.$el.appendChild(this.dragConf.ele)
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
  background-color: #FFFFFF;
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

.super-flow-base-demo {
  width: 100%;
  height: 500px;
  margin: 0 auto;
  background-color: #f5f5f5;
  @list-width: 200px;

  > .node-container {
    width: @list-width;
    float: left;
    height: 100%;
    text-align: center;
    background-color: #FFFFFF;
  }

  > .flow {
    width: calc(100% - @list-width);
    float: left;
    height: 100%;
    overflow: hidden;
  }

  .super-flow__node {
    .flow-node {
      > header {
        font-size: 12px;
        height: 16px;
        line-height: 16px;
        padding: 0 12px;
        color: #ffffff;
      }

      > section {
        text-align: center;
        line-height: 12px;
        font-size: 10px;
        overflow: hidden;
        padding: 6px 12px;
        word-break: break-all;
      }

      &.flow-node-start {
        > header {
          background-color: #55abfc;
        }
      }

      &.flow-node-process {
        > header {
          background-color: #BC1D16;
        }
      }

      &.flow-node-fork {
        >header {
          background-color: #BC1D16;
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
