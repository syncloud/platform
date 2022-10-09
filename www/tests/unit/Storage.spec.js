import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Storage from '../../src/views/Storage.vue'
import { ElSwitch, ElRadio, ElRadioGroup, ElCheckbox, ElCheckboxGroup } from 'element-plus'

jest.setTimeout(30000)

test('Activate single partition', async () => {
  let deviceAction = ''
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/disk/activate').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [200, { success: true }]
  })

  const wrapper = mount(Storage,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          'el-switch': ElSwitch,
          'el-radio': ElRadio,
          'el-radio-group': ElRadioGroup,
          'el-checkbox': ElCheckbox,
          'el-checkbox-group': ElCheckboxGroup,
          Confirmation: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#partition_confirmation').trigger('confirm')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

test('Activate single partition error', async () => {
  let deviceAction = ''
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: true,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/disk/activate').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [500, { message: 'not ok' }]
  })

  const wrapper = mount(Storage,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          'el-switch': ElSwitch,
          'el-radio': ElRadio,
          'el-radio-group': ElRadioGroup,
          'el-checkbox': ElCheckbox,
          'el-checkbox-group': ElCheckboxGroup,
          Confirmation: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()
  
  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#partition_confirmation').trigger('confirm')
  
  await flushPromises()

  expect(error).toBe('not ok')
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

test('Activate single partition service error', async () => {
  let deviceAction = ''
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: true,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/disk/activate').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [200, { success: false, message: 'not ok' }]
  })

  const wrapper = mount(Storage,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          'el-switch': ElSwitch,
          'el-radio': ElRadio,
          'el-radio-group': ElRadioGroup,
          'el-checkbox': ElCheckbox,
          'el-checkbox-group': ElCheckboxGroup,
          Confirmation: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()
  
  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#partition_confirmation').trigger('confirm')
  
  await flushPromises()

  expect(error).toBe('not ok')
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

