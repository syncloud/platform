import { mount, RouterLinkStub } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import AppCenter from '@/views/AppCenter'

jest.setTimeout(30000)

test('Show apps', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('rest/apps/available').reply(200,
    {
      data: [
        { id: 'app1', name: 'App1', icon: '/images/1.png' },
        { id: 'app2', name: 'App2', icon: '/images/2.png' }
      ]
    }
  )

  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const wrapper = mount(AppCenter,
    {
      attachTo: document.body,
      global: {
        components: {
          RouterLink: RouterLinkStub
        },
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/appcenter' },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(error).toBe('')
  expect(wrapper.text()).toContain('App1')
  expect(wrapper.text()).toContain('App2')

  wrapper.unmount()
})

test('Show error', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('rest/apps/available').reply(500,
    {
      message: 'not ok'
    }
  )

  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const wrapper = mount(AppCenter,
    {
      attachTo: document.body,
      global: {
        components: {
          RouterLink: RouterLinkStub
        },
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/appcenter' },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(error).toBe('not ok')

  wrapper.unmount()
})
