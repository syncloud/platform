<template>
  <div class="wrapper">
    <div class="content">

      <div class="wd12" id="block1">
        <div CLASS="block1" style="padding: 50px 0 0 0;">
          <h1>Backup</h1>
        </div>
        <div class="row-no-gutters settingsblock">

          <div id="backupGrid" style="width: 100%; height: 300px" class="ag-theme-balham"></div>

        </div>

      </div>

    </div>
  </div>

  <div id="backup_action_confirmation" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span>
          </button>
          <h4 class="modal-title"><span id="confirm_caption"></span></h4>
        </div>
        <div class="modal-body">
          <input type="hidden" id="backup_file" v-model="file"/>
          <input type="hidden" id="backup_action" v-model="action"/>
          <div class="bodymod">
            <div class="btext">
              <span id="confirm_question"></span>
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close
            </button>
            <button type="button" id="btn_confirm" class="btn buttonlight bwidth smbutton"
                    data-dismiss="modal" @click="submit">OK
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
import 'bootstrap'
import Error from '../components/Error.vue'
import toastr from 'toastr'
import { Grid } from 'ag-grid-community'
import axios from 'axios'
import * as Common from '../js/common.js'

export default {
  name: 'Backup',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      file: '',
      action: '',
      grid: undefined,
      gridOptions: undefined
    }
  },
  components: {
    Error
  },
  mounted () {
    this.gridOptions = {
      defaultColDef: {
        cellStyle: { 'text-align': 'left' }
      },
      columnDefs: [
        {
          headerName: 'File',
          field: 'file',
          resizable: true,
          sortable: true,
          suppressMovable: true,
          filter: 'agTextColumnFilter',
          floatingFilter: true
        },
        {
          headerName: 'Actions',
          width: 100,
          resizable: false,
          suppressMovable: true,
          cellRenderer: (params) => {
            const div = document.createElement('div')
            div.innerHTML = `
                <i class='fa fa-undo fa-2x' style='padding-left: 20px;  cursor:pointer;'></i>
                <i class='fa fa-trash fa-2x' style='padding-left: 20px;  cursor:pointer;'></i>
             `
            const buttons = div.querySelectorAll('i')
            buttons[0].addEventListener('click', () => {
              this.file = params.data.file
              this.action = 'restore'
              $('#confirm_caption').html('Restore')
              $('#confirm_question').html('Do you want to restore: ' + params.data.file + '?')
              $('#backup_action_confirmation').modal('show')
            })
            buttons[1].addEventListener('click', () => {
              this.file = params.data.file
              this.action = 'remove'
              $('#confirm_caption').html('Remove')
              $('#confirm_question').html('Do you want to remove: ' + params.data.file + '?')
              $('#backup_action_confirmation').modal('show')
            })
            return div
          }
        }
      ],
      suppressDragLeaveHidesColumns: true
    }

    const eGridDiv = document.querySelector('#backupGrid')
    eGridDiv.innerHTML = ''
    this.grid = new Grid(eGridDiv, this.gridOptions)
    this.reload()
  },
  methods: {
    submit () {
      switch (this.action) {
        case 'restore':
          this.restore()
          break
        case 'remove':
          this.remove()
          break
      }
    },
    remove () {
      axios.post('/rest/backup/remove', { file: this.file })
        .then(_ => {
          this.reload()
        })
        .catch(err => this.$refs.error.showToast(err))
    },
    restore () {
      const that = this
      axios
        .post('/rest/backup/restore', { file: this.file })
        .then(_ => {
          toastr.info('Restoring an app from a backup')

          Common.runAfterJobIsComplete(
            setTimeout,
            () => {
              toastr.info('Backup restore has finished')
              this.reload()
            },
            err => that.$refs.error.showToast(err),
            Common.JOB_STATUS_URL,
            Common.JOB_STATUS_PREDICATE)
        })
        .catch(err => {
          this.$refs.error.showToast(err)
        })
    },
    reload () {
      axios.get('/rest/backup/list')
        .then((response) => {
          this.gridOptions.api.setRowData(response.data.data)
          this.gridOptions.api.sizeColumnsToFit()
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
@import 'toastr/build/toastr.css';
@import "ag-grid-community/dist/styles/ag-grid.css";
@import "ag-grid-community/dist/styles/ag-theme-balham.css";
</style>
