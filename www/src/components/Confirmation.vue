<template>
  <div ref="confirmation" :id="id" class="modal fade bs-are-use-sure" tabindex="-1"
       role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content" >
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">
            <slot name="title"></slot>
          </h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              <slot name="text"></slot>
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" id="btn_partition_cancel" class="btn buttonlight bwidth smbutton"
                    data-dismiss="modal" @click="close">Close
            </button>
            <button type="button" id="btn_partition_action" class="btn buttonlight bwidth smbutton"
                    data-dismiss="modal" @click="yes">Yes
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
<script>
import 'bootstrap'
import $ from 'jquery'

export default {
  name: 'Confirmation',
  emits: ['confirm', 'cancel'],
  props: {
    id: String
  },
  data () {
    return {
      active: Boolean
    }
  },
  mounted () {
    this.active = false
  },
  methods: {
    show () {
      console.debug('show')
      this.active = true
      console.debug(this.active)
      $(this.$refs.confirmation).modal('show')
    },
    yes () {
      this.active = false
      this.$emit('confirm')
    },
    close () {
      this.active = false
      this.$emit('cancel')
    }
  }
}
</script>
