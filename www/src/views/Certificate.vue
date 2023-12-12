<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Certificate</h1>
        <div class="row-no-gutters settingsblock">

          <div class="col2">
            <div class="setline">
              <span class="span">Valid: </span>
              <i v-if="valid" class="material-icons icon-good" id="valid_good">check_circle</i>
              <i v-if="!valid" class="material-icons icon-bad" id="valid_bad">error</i>
            </div>
            <div class="setline">
              <span class="span">Valid days: </span>
              <span class="span" id="valid_days">{{ validDays }}</span>
            </div>
            <div class="setline">
              <span class="span">Real: </span>
              <i v-if="real" class="material-icons icon-good" id="real_good">check_circle</i>
              <i v-if="!real" class="material-icons icon-bad" id="real_bad">error</i>
            </div>

            <div class="setline">
              <span class="span">You can see more details</span>
              <div class="spandiv">
                <router-link to="/certificate/log" class="apps hlink">
                  <button class="buttonblue bwidth smbutton">Log</button>
                </router-link>
              </div>
            </div>

          </div>

        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Certificate',
  components: {
    Error
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      valid: false,
      real: false,
      validDays: 0,
      loading: undefined
    }
  },
  mounted () {
    this.progressShow()
    axios.get('/rest/certificate')
      .then((resp) => {
        const data = resp.data.data
        this.valid = data.is_valid
        this.real = data.is_real
        this.validDays = data.valid_for_days
        this.progressHide()
      })
      .catch(err => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      })
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
.icon-good {
  font-size: 20px;
  vertical-align: -15%;
  color: green;
}
.icon-bad {
  font-size: 20px;
  vertical-align: -15%;
  color: tomato;
}
</style>
