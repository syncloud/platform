import { defineStore } from 'pinia'
import axios from 'axios'

export const useSpaceStore = defineStore('space', {
  state: () => ({
    low: false,
    mounts: []
  }),
  getters: {
    lowest (state) {
      const low = state.mounts.filter(mount => mount.low)
      if (low.length === 0) {
        return undefined
      }
      return low.reduce((a, b) => (a.free_kb <= b.free_kb ? a : b))
    }
  },
  actions: {
    load () {
      return axios.get('/rest/storage/space')
        .then(response => {
          const data = response.data.data
          this.low = !!(data && data.low)
          this.mounts = (data && data.mounts) || []
        })
        .catch(() => {
          this.low = false
          this.mounts = []
        })
    }
  }
})
