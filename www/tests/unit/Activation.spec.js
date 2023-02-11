import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Activation from '../../src/views/Activation.vue'
import { ElButton } from 'element-plus'

jest.setTimeout(30000)

test('Activation url', async () => {
  
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  delete window.location
  window.location = ''
  const mock = new MockAdapter(axios)
  let availabilityDomain = ''
  let redirectEmail = ''
  let redirectPassword = ''
  mock.onPost('/rest/redirect/domain/availability').reply(function (config) {
    const request = JSON.parse(config.data)
    availabilityDomain = request.domain
    redirectEmail = request.email
    redirectPassword = request.password
    return [200, { success: true }]
  })
  mock.onGet('/rest/device/url').reply(200, { success: true, data: 'test.com' })

  const wrapper = mount(Activation,
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
          Dialog: true,
          'el-button': ElButton,
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_device_domain').text()).toBe('test.com')
  
  wrapper.unmount()
})

test('Activation reactivate', async () => {
  
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  delete window.location
  window.location = ''
  let deactivated = false
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/deactivate').reply(function (_) {
    deactivated = true
    return [200, { success: true }]
  })
  mock.onGet('/rest/device/url').reply(200, { success: true, data: 'test.com' })

  const wrapper = mount(Activation,
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
          Dialog: true,
          'el-button': ElButton,
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#txt_device_domain').text()).toBe('test.com')
  wrapper.find('#btn_reactivate').trigger('click')

  await flushPromises()

  expect(deactivated).toBe(true)

  wrapper.unmount()
})

