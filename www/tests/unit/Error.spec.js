import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Error from '@/components/Error.vue'

jest.setTimeout(30000)

test('Send', async () => {
  let sent = false
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/send_log').reply(function (_) {
    sent = true
    return [200, { success: true }]
  })

  const wrapper = mount(Error,
    {
      attachTo: document.body
    }
  )

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
  const wrapper = mount(Error, { attachTo: document.body })

  await wrapper.vm.showAxios({ response: { status: 500 } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('No response', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error, { attachTo: document.body })

  await wrapper.vm.showAxios({ })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('No message', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error, { attachTo: document.body })

  await wrapper.vm.showAxios({ response: { status: 500, data: {} } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('Server Error')

  wrapper.unmount()
})

test('Message', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error, { attachTo: document.body })

  await wrapper.vm.showAxios({ response: { status: 500, data: { message: 'test error' } } })

  expect(mockRouter.push).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_error').text()).toBe('test error')

  wrapper.unmount()
})

test('Show send logs by default', async () => {
  const mockRouter = { push: jest.fn() }
  const wrapper = mount(Error, { attachTo: document.body })

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

test('Parameters message', async () => {
  const mockRouter = { push: jest.fn() }

  const wrapper = mount(Error, {
    attachTo: document.body,
    props: {
      testing: true
    }
  })

  await wrapper.vm.showAxios({
    response: {
      status: 500,
      data: { parameters_messages: [{ parameter: 'test_parameter1', messages: ['1', '2'] }] }
    }
  })
  await flushPromises()

  expect(mockRouter.push).toHaveBeenCalledTimes(0)

  expect(wrapper.find('#test_parameter1').exists()).toBe(true)

  expect(wrapper.find('#test_parameter1_alert').text()).toBe('1\n2')

  wrapper.unmount()
})

test('Parameters message clean', async () => {
  const mockRouter = { push: jest.fn() }

  const wrapper = mount(Error, {
    attachTo: document.body,
    props: {
      testing: true
    }
  })

  await wrapper.vm.showAxios({
    response: {
      status: 500,
      data: { parameters_messages: [{ parameter: 'test_parameter1', messages: ['1', '2'] }] }
    }
  })
  await wrapper.vm.showAxios({
    response: {
      status: 500,
      data: { parameters_messages: [{ parameter: 'test_parameter1', messages: ['3', '4'] }] }
    }
  })
  await flushPromises()

  expect(mockRouter.push).toHaveBeenCalledTimes(0)

  expect(wrapper.find('#test_parameter1').exists()).toBe(true)

  const alerts = wrapper.findAll('#test_parameter1_alert')
  expect(alerts.length).toBe(1)
  expect(alerts[0].text()).toBe('3\n4')

  wrapper.unmount()
})
