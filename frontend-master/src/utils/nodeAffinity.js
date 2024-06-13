export function nodeAffinityParse(nodeStrategy, nodeValue) {
  let typeValue, displayValue, nodePrefer, nodeAffinityJSON
  // eslint-disable-next-line prefer-const
  typeValue = nodeStrategy.nodeType
  nodePrefer = nodeStrategy.nodePrefer
  if (nodeStrategy.nodeType === 'cloud' && nodeStrategy.isGPU && nodeStrategy.isNodeRandom) {
    displayValue = 'GPU_cloud'
    nodePrefer = 'random'
  } else if (nodeStrategy.nodeType === 'cloud' && nodeStrategy.isGPU && !nodeStrategy.isNodeRandom) {
    displayValue = 'GPU_cloud'
  } else if (nodeStrategy.nodeType === 'cloud' && !nodeStrategy.isGPU && !nodeStrategy.isNodeRandom) {
    displayValue = 'CPU_cloud'
  } else if (nodeStrategy.nodeType === 'cloud' && !nodeStrategy.isGPU && nodeStrategy.isNodeRandom) {
    displayValue = 'CPU_cloud'
    nodePrefer = 'random'
  } else if (nodeStrategy.nodeType === 'edge' && nodeStrategy.isNodeRandom) {
    displayValue = nodeStrategy.nodeArea
    nodeStrategy.nodePrefer = 'random'
  } else if (nodeStrategy.nodeType === 'edge') {
    displayValue = nodeStrategy.nodeArea
  }
  console.log('nodeValue:' + nodeValue)
  console.log('typeValue:' + typeValue)
  console.log('displayValue:' + displayValue)
  console.log('nodePrefer:' + nodePrefer)
  // 指定节点
  if (nodeValue !== '') {
    console.log('指定节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'kubernetes.io/hostname',
                operator: 'In',
                values: [nodeValue]
              }
            ]
          }
        ]
      }
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 随机选端
  if (typeValue === 'random') {
    console.log('随机选端')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              }
            ]
          }
        ]
      }
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在边缘端随机选display
  if (typeValue === 'edge' && displayValue === 'random') {
    console.log('在边缘端随机选display')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              }
            ]
          }
        ]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [
        {
          weight: 50,
          preference: {
            matchExpressions: [
              {
                key: 'type',
                operator: 'In',
                values: [
                  'edge'
                ]
              }
            ]
          }
        }
      ]
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在边缘端指定display,但是随机节点
  if (typeValue === 'edge' && (nodePrefer === 'random' || nodePrefer === '') && displayValue !== 'random' && displayValue !== '') {
    console.log('在边缘端指定display,但是随机节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              }
            ]
          }
        ]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [
        {
          weight: 50,
          preference: {
            matchExpressions: [
              {
                key: 'area',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        },
        {
          weight: 1,
          preference: {
            matchExpressions: [
              {
                key: 'type',
                operator: 'In',
                values: [typeValue]
              }
            ]
          }
        }
      ]
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在边缘端指定display,指定节点
  if (typeValue === 'edge' && nodePrefer !== 'random' && nodePrefer !== '' && displayValue !== 'random' && displayValue !== '') {
    console.log('在边缘端指定display,指定节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              }
            ]
          }
        ]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [
        {
          weight: 50,
          preference: {
            matchExpressions: [
              {
                key: 'kubernetes.io/hostname',
                operator: 'In',
                values: [nodePrefer]
              }
            ]
          }
        },
        {
          weight: 20,
          preference: {
            matchExpressions: [
              {
                key: 'area',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        },
        {
          weight: 1,
          preference: {
            matchExpressions: [
              {
                key: 'type',
                operator: 'In',
                values: [typeValue]
              }
            ]
          }
        }
      ]
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在云端使用GPU并且随机选节点
  if (typeValue === 'cloud' && displayValue === 'GPU_cloud' && (nodePrefer === 'random' || nodePrefer === '')) {
    console.log('在云端使用GPU并且随机选节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              },
              {
                key: 'display',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        ]
      }
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在云端使用GPU并且指定节点
  if (typeValue === 'cloud' && displayValue === 'GPU_cloud' && nodePrefer !== 'random' && nodePrefer !== '') {
    console.log('在云端使用GPU并且指定节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              },
              {
                key: 'display',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        ]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [
        {
          weight: 50,
          preference: {
            matchExpressions: [
              {
                key: 'kubernetes.io/hostname',
                operator: 'In',
                values: [nodePrefer]
              }
            ]
          }
        }
      ]
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在云端使用CPU并且随机节点
  if (typeValue === 'cloud' && displayValue === 'CPU_cloud' && (nodePrefer === 'random' || nodePrefer === '')) {
    console.log('在云端使用CPU并且随机节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              },
              {
                key: 'display',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        ]
      }
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  // 在云端使用CPU并且指定节点
  if (typeValue === 'cloud' && displayValue === 'CPU_cloud' && nodePrefer !== 'random' && nodePrefer !== '') {
    console.log('在云端使用CPU并且指定节点')
    nodeAffinityJSON = {
      requiredDuringSchedulingIgnoredDuringExecution: {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'type',
                operator: 'NotIn',
                values: [
                  'master'
                ]
              },
              {
                key: 'display',
                operator: 'In',
                values: [displayValue]
              }
            ]
          }
        ]
      },
      preferredDuringSchedulingIgnoredDuringExecution: [
        {
          weight: 50,
          preference: {
            matchExpressions: [
              {
                key: 'kubernetes.io/hostname',
                operator: 'In',
                values: [nodePrefer]
              }
            ]
          }
        }
      ]
    }
    return JSON.stringify(nodeAffinityJSON)
  }
  return null
}
