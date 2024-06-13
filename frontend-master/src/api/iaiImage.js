import axios from 'axios'

export function getIaiImageList(isTotal) {
  return axios.get('/iai/api/brainController/getIaiImageList', { params: { isTotal: isTotal }})
}

export function editIaiImageRecord(data) {
  return axios.post('/iai/api/brainController/editIaiImageRecord', data)
}
