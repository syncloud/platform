import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Activate from '../../src/views/Activate.vue'
import { ElButton, ElStep, ElSteps } from 'element-plus'

jest.setTimeout(30000)

test('Activate free domain', async () => {
  let redirectEmail = ''
  let redirectPassword = ''
  let domain = ''
  let deviceUsername = ''
  let devicePassword = ''
  let availabilityDomain = ''
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  delete window.location
  window.location = ''
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(function (config) {
    const request = JSON.parse(config.data)
    redirectEmail = request.redirect_email
    redirectPassword = request.redirect_password
    domain = request.domain
    deviceUsername = request.device_username
    devicePassword = request.device_password
    return [200, { success: true }]
  })
  mock.onPost('/rest/redirect/domain/availability').reply(function (config) {
    const request = JSON.parse(config.data)
    availabilityDomain = request.domain
    redirectEmail = request.email
    redirectPassword = request.password
    return [200, { success: true }]
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

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
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
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
  await wrapper.find('#domain_input').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#device_password_confirm').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(availabilityDomain).toBe('domain.test.com')
  expect(redirectEmail).toBe('r email')
  expect(redirectPassword).toBe('r password')
  expect(domain).toBe('domain.test.com')
  expect(deviceUsername).toBe('user')
  expect(devicePassword).toBe('password')
  expect(window.location).toMatch(new RegExp('^/\\?t=.*'))

  wrapper.unmount()
})

test('Activate free domain error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(500, {
    message: 'not ok'
  })

  mock.onPost('/rest/redirect/domain/availability').reply(200, {
    message: 'ok'
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

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
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
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
  await wrapper.find('#domain_input').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#device_password_confirm').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(error).toBe('not ok')

  wrapper.unmount()
})

test('Activate free domain availability error', async () => {
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(500, {
    message: 'not ok'
  })

  mock.onPost('/rest/redirect/domain/availability').reply(400, {
    success: false,
    parameters_messages: [
      { parameter: 'domain', messages: ['domain is already taken'] }
    ]
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

  const wrapper = mount(Activate,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: true,
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
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
  await wrapper.find('#domain_input').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')

  await flushPromises()

  expect(wrapper.find('#domain_alert').text()).toBe('domain is already taken')

  wrapper.unmount()
})

test('Activate free domain availability generic error', async () => {
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(500, {
    message: 'not ok'
  })

  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }

  mock.onPost('/rest/redirect/domain/availability').reply(400, {
    success: false,
    message: 'authentication failed'
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

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
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
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
  await wrapper.find('#domain_input').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')

  await flushPromises()

  expect(error).toBe('authentication failed')

  wrapper.unmount()
})

test('Activate free domain availability error email', async () => {
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(500, {
    message: 'not ok'
  })

  mock.onPost('/rest/redirect/domain/availability').reply(400, {
    success: false,
    parameters_messages: [
      { parameter: 'email', messages: ['wrong email'] }
    ]
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

  const wrapper = mount(Activate,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: true,
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
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
  await wrapper.find('#domain_input').setValue('domain')
  await wrapper.find('#btn_next').trigger('click')

  await flushPromises()

  expect(wrapper.find('#email_alert').text()).toBe('wrong email')

  wrapper.unmount()
})

test('Activate premium domain', async () => {
  let redirectEmail = ''
  let redirectPassword = ''
  let domain = ''
  let deviceUsername = ''
  let devicePassword = ''
  let availabilityDomain = ''
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  delete window.location
  window.location = ''
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(function (config) {
    const request = JSON.parse(config.data)
    redirectEmail = request.redirect_email
    redirectPassword = request.redirect_password
    domain = request.domain
    deviceUsername = request.device_username
    devicePassword = request.device_password
    return [200, { success: true }]
  })
  mock.onPost('/rest/redirect/domain/availability').reply(function (config) {
    const request = JSON.parse(config.data)
    availabilityDomain = request.domain
    redirectEmail = request.email
    redirectPassword = request.password
    return [200, { success: true }]
  })
  mock.onGet('/rest/redirect_info').reply(200, { success: true, data: { domain: 'test.com' } })

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
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_premium_domain').trigger('click')
  await wrapper.find('#email').setValue('r email')
  await wrapper.find('#redirect_password').setValue('r password')
  await wrapper.find('#domain_premium').setValue('example.com')
  await wrapper.find('#btn_next').trigger('click')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#device_password_confirm').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(redirectEmail).toBe('r email')
  expect(redirectPassword).toBe('r password')
  expect(domain).toBe('example.com')
  expect(deviceUsername).toBe('user')
  expect(devicePassword).toBe('password')
  expect(window.location).toMatch(new RegExp('^/\\?t=.*'))
  expect(availabilityDomain).toBe('example.com')

  wrapper.unmount()
})

test('Activated while page is open (mostly in local dev)', async () => {
  let redirectEmail = ''
  let redirectPassword = ''
  let domain = ''
  let deviceUsername = ''
  let devicePassword = ''
  let availabilityDomain = ''
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }
  delete window.location
  window.location = ''
  const mock = new MockAdapter(axios)
  mock.onPost('/rest/activate/managed').reply(function (config) {
    const request = JSON.parse(config.data)
    redirectEmail = request.redirect_email
    redirectPassword = request.redirect_password
    domain = request.domain
    deviceUsername = request.device_username
    devicePassword = request.device_password
    return [200, { success: true }]
  })
  mock.onPost('/rest/redirect/domain/availability').reply(function (config) {
    const request = JSON.parse(config.data)
    availabilityDomain = request.domain
    redirectEmail = request.email
    redirectPassword = request.password
    return [200, { success: true }]
  })
  mock.onGet('/rest/redirect_info').reply(502, { message: 'Device is already activated' })

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
          Dialog: true,
          'el-button': ElButton,
          'el-steps': ElSteps,
          'el-step': ElStep
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_premium_domain').trigger('click')
  await wrapper.find('#email').setValue('r email')
  await wrapper.find('#redirect_password').setValue('r password')
  await wrapper.find('#domain_premium').setValue('example.com')
  await wrapper.find('#btn_next').trigger('click')
  await wrapper.find('#device_username').setValue('user')
  await wrapper.find('#device_password').setValue('password')
  await wrapper.find('#device_password_confirm').setValue('password')
  await wrapper.find('#btn_activate').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(redirectEmail).toBe('r email')
  expect(redirectPassword).toBe('r password')
  expect(domain).toBe('example.com')
  expect(deviceUsername).toBe('user')
  expect(devicePassword).toBe('password')
  expect(window.location).toMatch(new RegExp('^/\\?t=.*'))
  expect(availabilityDomain).toBe('example.com')

  wrapper.unmount()
})
