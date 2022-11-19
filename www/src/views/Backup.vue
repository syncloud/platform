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

  <el-table :data="data" style="width: 100%" table-layout="fixed">
    <el-table-column label="File" prop="file" width="250px"/>
    <el-table-column align="right">
      <template #header>
        <el-input v-model="search" size="small" placeholder="Type to search" />
      </template>
      <template #default="scope">
        <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
          >Edit</el-button
        >
        <el-button
          size="small"
          type="danger"
          @click="handleDelete(scope.$index, scope.row)"
          >Delete</el-button
        >
      </template>
    </el-table-column>
  </el-table>

      </div>

    </div>
  </div>

  <Confirmation :visible="confirmationVisible" id="confirmation" @confirm="submit"
                @cancel="confirmationVisible = false">
    <template v-slot:title>
      <span v-if="action === 'restore'">Restore</span>
      <span v-if="this.action === 'remove'">Remove</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          <span v-if="action === 'restore'">Do you want to restore: {{ file }}?</span>
          <span v-if="action === 'remove'">Do you want to remove: {{ file }}?</span>
        </div>
      </div>
    </template>
  </Confirmation>

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
import Confirmation from '../components/Confirmation.vue'

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
      gridOptions: undefined,
      confirmationVisible: false,
      data: [],
    }
  },
  components: {
    Error,
    Confirmation
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
              this.confirmationVisible = true
            })
            buttons[1].addEventListener('click', () => {
              this.file = params.data.file
              this.action = 'remove'
              this.confirmationVisible = true
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
      this.confirmationVisible = false
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
          this.data = response.data.data
          this.gridOptions.api.setRowData(response.data.data)
          this.gridOptions.api.sizeColumnsToFit()
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
@import 'toastr/build/toastr.css';
@import "ag-grid-community/dist/styles/ag-grid.css";
@import "ag-grid-community/dist/styles/ag-theme-balham.css";
</style>
