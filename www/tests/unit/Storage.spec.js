import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Storage from '../../src/views/Storage.vue'
import Confirmation from '../../src/components/Confirmation.vue'
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
  mock.onPost('/rest/storage/disk/activate/partition').reply(function (config) {
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
            template: '<span :id="id"><slot name="text"></slot></span>',
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
  await expect(wrapper.find('#format').isVisible()).toBe(true)
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
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
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
  mock.onPost('/rest/storage/disk/activate/partition').reply(function (config) {
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
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
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
  mock.onPost('/rest/storage/disk/activate/partition').reply(function (config) {
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


test('Activate multiple disks', async () => {
  let devices = ''
  let error = ''
  const showError = (err) => {
    error = err
  }
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
  mock.onPost('/rest/storage/disk/activate/disk').reply(function (config) {
    devices = JSON.parse(config.data).devices
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
            template: '<span :id="id"><slot name="text"></slot></span>',
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
  
  await wrapper.find('#multi').trigger('click')
  await wrapper.find('#disk_0').trigger('click')
  await wrapper.find('#disk_1').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  await expect(wrapper.find('#format').isVisible()).toBe(true)
  await wrapper.find('#partition_confirmation').trigger('confirm')
  
  await flushPromises()
  
  expect(error).toBe('')
  expect(devices).toEqual(['/dev/sdb', '/dev/sdc'])
  wrapper.unmount()
})


test('Show single partition', async () => {
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
            { active: true, device: '/dev/sdb1', size: '931.5G' }
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
  
  await expect(wrapper.find('#partition_0_0').classes()).toContain('is-checked')
  
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Show single partition none', async () => {
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
  
  await expect(wrapper.find('#none').classes()).toContain('is-checked')
  
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Show multi disk', async () => {
  const showError = jest.fn()
  
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
  
  await expect(wrapper.find('#multi').attributes('aria-checked')).toBe("true")
  await expect(wrapper.find('#disk_0').element.parentElement.getAttribute('class')).toContain('is-checked')
  await expect(wrapper.find('#disk_1').element.parentElement.getAttribute('class')).not.toContain('is-checked')
  
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})
