import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Updates from '../../src/views/Updates.vue'

jest.setTimeout(30000)

test('Update platform', async () => {
  let upgraded = false
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: {
          id: 'platform',
          name: 'Platform',
          required: true,
          ui: false,
          url: 'http://platform.odroid-c2.syncloud.it'
        },
        current_version: '1',
        installed_version: '2'
      },
      success: true
    }]
  })

  mock.onGet('/rest/installer/version').reply(200,
    {
      data: {
        store_version: '3',
        installed_version: '3'
      },
      success: true
    }
  )
  mock.onPost('/rest/app/upgrade').reply(function (_) {
    upgraded = true
    return [200, { success: true }]
  })
  
  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, { success: true, data: { is_running: false } }]
  })
  
  const wrapper = mount(Updates,
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
          Switch: true,
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  expect(wrapper.find('#btn_platform_upgrade').exists()).toBe(true)
  expect(wrapper.find('#btn_installer_upgrade').exists()).toBe(false)

  await wrapper.find('#btn_platform_upgrade').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(upgraded).toBe(true)
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Update installer', async () => {
  let upgraded = false
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: {
          id: 'platform',
          name: 'Platform',
          required: true,
          ui: false,
          url: 'http://platform.odroid-c2.syncloud.it'
        },
        current_version: '1',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onGet('/rest/installer/version').reply(200,
    {
      data: {
        store_version: '2',
        installed_version: '3'
      },
      success: true
    }
  )
  mock.onPost('/rest/installer/upgrade').reply(function (_) {
    upgraded = true
    return [200, { success: true }]
  })
  let statusCalled = false
  mock.onGet('/rest/job/status').reply(function (_) {
    statusCalled = true
    return [200, {success: true, data:{status: "Idle"}}]
  })
  
  const wrapper = mount(Updates,
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
          Switch: true,
          Dialog: true
        }
      }
    }
  )
  await flushPromises()

  expect(wrapper.find('#btn_platform_upgrade').exists()).toBe(false)
  expect(wrapper.find('#btn_installer_upgrade').exists()).toBe(true)

  await wrapper.find('#btn_installer_upgrade').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(upgraded).toBe(true)
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Update installer error', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: {
          id: 'platform',
          name: 'Platform',
          required: true,
          ui: false,
          url: 'http://platform.odroid-c2.syncloud.it'
        },
        current_version: '1',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onGet('/rest/installer/version').reply(200,
    {
      data: {
        store_version: '2',
        installed_version: '3'
      },
      success: true
    }
  )
  mock.onPost('/rest/installer/upgrade').reply(function (_) {
    return [500, { success: false }]
  })

  const wrapper = mount(Updates,
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
          Switch: true,
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_installer_upgrade').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})
