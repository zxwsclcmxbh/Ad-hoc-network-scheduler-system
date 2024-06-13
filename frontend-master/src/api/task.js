import axios from 'axios'
export async function getNodesT() {
  const resp = await axios.get('/iai/api/brainController/getNodeMetrics')
  const data = resp.data
  var result = {}
  for (const i in data.data) {
    const item = data.data[i]
    var label
    if (item.metadata.labels.display) {
      label = item.metadata.labels.display.split('_').map((item) => {
        switch (item) {
          case 'GPU':
          case 'CPU':
            return item
          case 'mingluo':
            return '明珞'
          case 'wusuo':
            return '五所'
          case 'cloud':
            return '云端'
          case 'edge':
            return '边缘端'
        }
      })
    } else {
      label = []
    }
    const cpu_capacity = item.capacity.cpu
    let mem_capacity = item.capacity.memory.replace('Mi', '').replace('Ki', '')
    if (mem_capacity >= 1024 * 1024) {
      mem_capacity = mem_capacity / 1024
    }
    if (item.usage) {
      const cpu_usage = item.usage.cpu.replace('m', '') / 1000

      const mem_usage = item.usage.memory.replace('Ki', '') / 1024

      const cpu_percent = cpu_usage / cpu_capacity * 100
      const mem_percent = mem_usage / mem_capacity * 100
      result[item.metadata.name] = {
        type: item.metadata.labels.type && item.metadata.labels.type == 'master' ? 'master' : 'node',
        label,
        cpu_percent: cpu_percent,
        cpu_capacity: cpu_capacity,
        mem_percent,
        mem_capacity: mem_capacity,
        status: item.status
      }
    } else {
      result[item.metadata.name] = {
        type: item.metadata.labels.type && item.metadata.labels.type == 'master' ? 'master' : 'node',
        label,
        cpu_percent: 'unkonwn',
        cpu_capacity: cpu_capacity,
        mem_percent: 'unkonwn',
        mem_capacity: mem_capacity,
        status: item.status
      }
    }
  }
  return result
}
export async function getNodes() {
  const resp = await axios.get('/iai/api/brainController/getNodeMetrics')
  const data = resp.data
  var result = {}
  for (const i in data.data) {
    const item = data.data[i]
    var label
    if (item.metadata.labels.display) {
      label = item.metadata.labels.display.split('_').map((item) => {
        switch (item) {
          case 'GPU':
          case 'CPU':
            return item
          case 'mingluo':
            return '明珞'
          case 'wusuo':
            return '五所'
          case 'cloud':
            return '云端'
          case 'edge':
            return '边缘端'
        }
      })

    } else {
      label = []
    }
    const cpu_capacity = item.capacity.cpu
    let mem_capacity = item.capacity.memory.replace('Mi', '').replace('Ki', '')
    if (mem_capacity >= 1024 * 1024) {
      mem_capacity = mem_capacity / 1024
    }
    if (item.usage) {
      const cpu_usage = item.usage.cpu.replace('m', '') / 1000

      const mem_usage = item.usage.memory.replace('Ki', '') / 1024

      const cpu_percent = cpu_usage / cpu_capacity * 100
      const mem_percent = mem_usage / mem_capacity * 100
      let status = '01'
      if (cpu_percent < 25 && mem_percent < 25) {
        status = '04'
      } else if ((cpu_percent <= 50 && mem_percent <= 50) || (cpu_percent <= 20 && mem_percent <= 70)) {
        status = '03'
      } else if ((cpu_percent <= 75 && mem_percent <= 75) || (cpu_percent <= 50 && mem_percent <= 85)) {
        status = '02'
      } else {
        status = '01'
      }
      result[item.metadata.name] = {
        tag: label,
        area: item.metadata.labels.area,
        cpu_usage: cpu_usage,
        cpu_capacity: cpu_capacity,
        mem_usage: mem_usage,
        mem_capacity: mem_capacity,
        cpu: cpu_percent,
        mem: mem_percent,
        status: status
      }
    } else {
      result[item.metadata.name] = {
        area: item.metadata.labels.area,
        tag: label,
        cpu_usage: 0,
        cpu_capacity: cpu_capacity,
        mem_usage: 0,
        mem_capacity: mem_capacity,
        cpu: 0,
        mem: 0,
        status: '01'
      }
    }
  }
  return result
}

export function getPodStatus(data) {
  return axios.post('/iai/api/brainController/getIaiPodStatus', data)
}
export function getPodName(data) {
  return axios.post('/iai/api/brainController/getIaiPodName', data)
}
export function getPodHostName(data) {
  return axios.post('/iai/api/brainController/getIaiPodHostName', data)
}
export function getBlockSettings() {
  return axios.get('/iai/api/platform/getSchemaConfig')
}

export function getDeviceMapper() {
  return axios.get('/iai/api/platform/getDataMapperConfig')
}

export function searchImage(name) {
  return axios.get('/iai/api/platform/searchImageWarehouse', { params: { search: name }})
}

export function submitTask(data) {
  return axios.post('/iai/api/brainController/createTaskNew', data)
}

export function getTask(taskid) {
  return axios.get('/iai/api/brainController/getTask', { params: { taskid: taskid }})
}

export function getTaskList() {
  return axios.get('/iai/api/brainController/getTaskList')
}

export function deleteTask(data) {
  return axios.post('/iai/api/brainController/deleteTask', data)
}

export function podMigration(data) {
  return axios.post('/iai/api/brainController/podMigration', data)
}
export function updatePodImage(data) {
  return axios.post('/iai/api/brainController/podImagesUpdate', data)
}

export function NodeAffinityMigration(data) {
  return axios.post('/iai/api/brainController/nodeAffinityMigration',data)
}

export function getTasks() {
  return axios.get('/iai/api/brainController/getTasks')
}
