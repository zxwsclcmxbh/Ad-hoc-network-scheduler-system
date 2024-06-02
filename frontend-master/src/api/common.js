import axios from 'axios'
export function getCount() {
  return axios.get('/iai/api/platform/getBaseCount')
}
