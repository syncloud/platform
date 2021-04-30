import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Activate from '@/views/Activate'

jest.setTimeout(30000)

test('Activate free domain', async () => {
  let redirectEmail = ''
  let redirectPassword = ''
  let domain = ''
  let deviceUsername = ''
  let devicePassword = ''
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  let reloaded = false
  delete window.location
  window.location = {
    reload (resetCache) {
      reloaded = true
    }
  }
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate').reply(function (config) {
    const request = JSON.parse(config.data)
    redirectEmail = request.redirect_email
    redirectPassword = request.redirect_password
    domain = request.user_domain
    deviceUsername = request.device_username
    devicePassword = request.device_password
    return [200, { success: true }]
  })

  const wrapper = mount(Activate,
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
          Dialog: true
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_free_domain').trigger('click')
  await wrapper.find('#email').setValue('r email')
  await wrapper.find('#redirect_password').setValue('r password')
  await wrapper.find('#user_domain').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(redirectEmail).toBe('r email')
  expect(redirectPassword).toBe('r password')
  expect(domain).toBe('domain')
  expect(deviceUsername).toBe('user')
  expect(devicePassword).toBe('password')
  expect(reloaded).toBe(true)

  wrapper.unmount()
})

test('Activate free domain error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate').reply(500, {
    message: 'not ok'
  })

  const wrapper = mount(Activate,
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
          Dialog: true
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#email').setValue('r email')
  await wrapper.find('#redirect_password').setValue('r password')
  await wrapper.find('#user_domain').setValue('domain')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(error).toBe('not ok')

  wrapper.unmount()
})

test('Activate custom domain', async () => {
  let domain = ''
  let deviceUsername = ''
  let devicePassword = ''
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  let reloaded = false
  delete window.location
  window.location = {
    reload (resetCache) {
      reloaded = true
    }
  }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate_custom_domain').reply(function (config) {
    const request = JSON.parse(config.data)
    domain = request.full_domain
    deviceUsername = request.device_username
    devicePassword = request.device_password
    return [200, { success: true }]
  })

  const wrapper = mount(Activate,
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
          Dialog: true
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#domain_type_custom').trigger('click')
  await wrapper.find('#full_domain').setValue('domain')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(domain).toBe('domain')
  expect(deviceUsername).toBe('user')
  expect(devicePassword).toBe('password')
  expect(reloaded).toBe(true)

  wrapper.unmount()
})
