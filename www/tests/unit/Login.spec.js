import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Login from '../../src/views/Login.vue'
import { ElButton } from 'element-plus'

jest.setTimeout(30000)

test('Login', async () => {
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }

  let username
  let password

  const mock = new MockAdapter(axios)
  mock.onPost('/rest/login').reply(function (config) {
    const request = JSON.parse(config.data)
    username = request.username
    password = request.password
    return [200, { success: true }]
  })

  const wrapper = mount(Login,
    {
      attachTo: document.body,
      props: {
        checkUserSession: jest.fn(),
        activated: true
      },
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          'el-button': ElButton
        },
        mocks: {
          $route: { path: '/login' },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#username').setValue('username')
  await wrapper.find('#password').setValue('password')
  await wrapper.find('#btn_login').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(username).toBe('username')
  expect(password).toBe('password')
  expect(mockRouter.push).toHaveBeenCalledWith('/')

  wrapper.unmount()
})
