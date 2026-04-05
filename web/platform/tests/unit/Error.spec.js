import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Error from '../../src/components/Error.vue'
import { ElButton, ElDialog } from 'element-plus'

jest.setTimeout(30000)

test('Send', async () => {
  const mockRouter = { push: jest.fn() }
  let sent = false
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/logs/send').reply(function (_) {
    sent = true
    return [200, { success: true }]
  })

  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        },
        mocks: {
          $route: { path: '/' },
          $router: mockRouter
        }
      }
    }
  )

  await wrapper.vm.showAxios({ response: { status: 400 } })

  await wrapper.find('#btn_error_send_logs').trigger('click')
  await flushPromises()

  expect(sent).toBe(true)
  wrapper.unmount()
})

test('Unauthorised', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        },
        mocks: {
          $route: { path: '/app' },
          $router: mockRouter
        }
      }
    }
  )

  await wrapper.vm.showAxios({ response: { status: 401 } })
  expect(mockRouter.push).toHaveBeenCalledWith('/login')

  wrapper.unmount()
})

test('Not activated', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        },
        mocks: {
          $route: { path: '/' },
          $router: mockRouter
        }
      }
    }
  )

  await wrapper.vm.showAxios({ response: { status: 501 } })
  expect(mockRouter.push).toHaveBeenCalledWith('/activate')

  wrapper.unmount()
})

test('No data', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        }
      }
    }
  )

  await wrapper.vm.showAxios({ response: { status: 500 } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('No response', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        }
      }
    })

  await wrapper.vm.showAxios({ })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('No message', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        }
      }
    })

  await wrapper.vm.showAxios({ response: { status: 500, data: {} } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('Message', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        }
      }
    })

  await wrapper.vm.showAxios({ response: { status: 500, data: { message: 'test error' } } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('test error')

  wrapper.unmount()
})

test('Show send logs by default', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-dialog': ElDialog,
          'el-button': ElButton
        }
      }
    })

  await wrapper.vm.showAxios({ response: { status: 500, data: { message: 'test error' } } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('test error')
  expect(wrapper.find('#btn_error_send_logs').exists()).toBe(true)

  wrapper.unmount()
})

test('Disable send logs', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error, {
    attachTo: document.body,
    global: {
      stubs: {
        'el-dialog': ElDialog,
        'el-button': ElButton
      }
    },
    propsData: {
      enableLogs: false
    }
  })

  await wrapper.vm.showAxios({ response: { status: 500, data: { message: 'test error' } } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('test error')
  expect(wrapper.find('#btn_error_send_logs').exists()).toBe(false)

  wrapper.unmount()
})
