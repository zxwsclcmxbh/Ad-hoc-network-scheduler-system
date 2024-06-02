<template>
  <div class="dashboard-editor-container">

    <panel-group @handleSetLineChartData="handleSetLineChartData" />
    <el-tabs type="border-card">
      <el-tab-pane label="资源使用情况">
        <el-row style="padding:16px 16px 0;margin-bottom:32px;">
          <el-tabs type="border-card">
            <el-tab-pane label="百分比">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=9`" width="100%" height="400" frameborder="0" />
            </el-tab-pane>
            <el-tab-pane label="用量">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=2`" width="100%" height="400" frameborder="0" />
            </el-tab-pane>
          </el-tabs>
        </el-row>
        <el-row style="padding:16px 16px 0;margin-bottom:32px;">
          <el-tabs type="border-card">
            <el-tab-pane label="百分比">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=8`" width="100%" height="400" frameborder="0" />
            </el-tab-pane>
            <el-tab-pane label="用量">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=4`" width="100%" height="400" frameborder="0" />
            </el-tab-pane>
          </el-tabs>
        </el-row>
        <!-- <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <div style="">
        <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=2`" width="100%" height="400" frameborder="0" />
      </div>
    </el-row> -->
        <!-- <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <div style="">
        <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=4`" width="100%" height="400" frameborder="0" />
      </div>
    </el-row> -->

        <el-row :gutter="32">
          <el-col :xs="24" :sm="24" :lg="12">
            <div class="chart-wrapper">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=6`" width="100%" height="300" frameborder="0" />
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :lg="12">
            <div class="chart-wrapper">
              <iframe :src="`${grafana_base}d-solo/1RS3pE4Vk/zhu-ye?orgId=1&refresh=5s&theme=light&panelId=7`" width="100%" height="300" frameborder="0" />
            </div>
          </el-col>
        </el-row>
      </el-tab-pane>
      <el-tab-pane label="系统拓扑">
        <el-row style="padding:16px 16px 0;margin-bottom:32px;"><topology /></el-row>

      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
import Topology from '@/components/Topology'

const lineChartData = {
  newVisitis: {
    expectedData: [100, 120, 161, 134, 105, 160, 165],
    actualData: [120, 82, 91, 154, 162, 140, 145]
  },
  messages: {
    expectedData: [200, 192, 120, 144, 160, 130, 140],
    actualData: [180, 160, 151, 106, 145, 150, 130]
  },
  purchases: {
    expectedData: [80, 100, 121, 104, 105, 90, 100],
    actualData: [120, 90, 100, 138, 142, 130, 130]
  },
  shoppings: {
    expectedData: [130, 140, 141, 142, 145, 150, 160],
    actualData: [120, 82, 91, 154, 162, 140, 130]
  }
}

export default {
  name: 'DashboardAdmin',
  components: {
    Topology,
    PanelGroup

  },
  data() {
    return {
      lineChartData: lineChartData.newVisitis,
      grafana_base: 'http://10.112.134.196:30505/'
    }
  },
  methods: {
    handleSetLineChartData(type) {
      this.lineChartData = lineChartData[type]
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width:1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
