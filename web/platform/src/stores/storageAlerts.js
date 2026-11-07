import { defineStore } from 'pinia'
import axios from 'axios'

export const useStorageAlertsStore = defineStore('storageAlerts', {
  state: () => ({
    problems: [],
    fresh: false
  }),
  actions: {
    load () {
      return axios.get('/rest/storage/alerts')
        .then(response => {
          const data = response.data.data
          this.problems = (data && data.problems) || []
          this.fresh = !!(data && data.new)
        })
        .catch(() => {
          this.problems = []
          this.fresh = false
        })
    },
    seen () {
      return axios.post('/rest/storage/alerts/seen')
        .then(() => this.load())
    }
  }
})
