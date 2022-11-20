import { mount } from '@vue/test-utils'
import Backup from '../../src/views/Backup.vue'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import { ElOption, ElSelect, ElTable, ElInput, ElTableColumn, ElButton } from 'element-plus'

test('show list of backups', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/backup/list').reply(200,
    {
      success: true,
      data: [
        { path: '/data/platform/backup', file: 'files-2019-0515-123506.tar.gz' },
        { path: '/data/platform/backup', file: 'nextcloud-2019-0515-123506.tar.gz' }
      ]
    }
  )
  mock.onGet('/rest/backup/auto').reply(200,
    {
      success: true,
      data: { auto: 'no', auto_day: 0, auto_time: 0}
    }
  )

  const wrapper = mount(Backup,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-option': ElOption,
          'el-select': ElSelect,
          'el-table': ElTable,
          'el-table-column': ElTableColumn,
          'el-input': ElInput,
          'el-button': ElButton,
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          },
        }
      }
    }
  )

  await flushPromises()

  expect(wrapper.text()).toContain('files-2019-0515-123506.tar.gz')
  expect(wrapper.find('#auto').attributes().disabled).not.toBeDefined()
  expect(wrapper.find('#auto-day').attributes().disabled).toBeDefined()
  expect(wrapper.find('#auto-time').attributes().disabled).toBeDefined()

  wrapper.unmount()
})

test('save auto mode', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/backup/list').reply(200,
    {
      success: true,
      data: []
    }
  )
  let auto = 'backup'
  let autoDay = 1
  let autoTime = 2
  mock.onGet('/rest/backup/auto').reply(200,
    {
      success: true,
      data: { auto: auto, day: autoDay, time: autoTime}
    }
  )
  let saved = false
  mock.onPost('/rest/backup/auto').reply(function (config) {
    let request = JSON.parse(config.data)
    auto = request.auto
    autoDay = request.day
    autoTime = request.time
    saved = true
    return [200, { success: true }]
  })

  const wrapper = mount(Backup,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-option': ElOption,
          'el-select': ElSelect,
          'el-table': ElTable,
          'el-table-column': ElTableColumn,
          'el-input': ElInput,
          'el-button': ElButton,
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          },
        }
      }
    }
  )

  await flushPromises()

  expect(wrapper.find('#auto').attributes().disabled).not.toBeDefined()
  expect(wrapper.find('#auto-day').attributes().disabled).not.toBeDefined()
  expect(wrapper.find('#auto-time').attributes().disabled).not.toBeDefined()
  
  // await wrapper.find('#auto').trigger('click')
  // await flushPromises()
  // await wrapper.find('#auto-backup').trigger('click')
  
  // await wrapper.find('#auto-day').trigger('click')
  // await wrapper.find('#auto-day-monday').trigger('click')
  
  await wrapper.find('#save').trigger('click')
  
  await flushPromises()

  expect(auto).toBe('backup')
  expect(autoDay).toBe(1)
  expect(autoTime).toBe(2)
  expect(saved).toBe(true)
  
  wrapper.unmount()
})
