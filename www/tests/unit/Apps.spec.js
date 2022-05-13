import { mount, RouterLinkStub } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Apps from '@/views/Apps'

jest.setTimeout(30000)

test('Show apps', async () => {
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/apps/installed').reply(function (_) {
    return [200, {
      data: [
        {
          id: 'wordpress',
          name: 'WordPress',
          icon: '/images/wordpress-128.png'
        },
        {
          id: 'diaspora',
          name: 'Diaspora',
          icon: '/images/penguin.png'
        }
      ]
    }]
  })

  const wrapper = mount(Apps,
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
          $route: { path: '/apps' },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.text()).toContain('WordPress')
  expect(wrapper.text()).toContain('Diaspora')

  wrapper.unmount()
})

test('Show error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/apps/installed').reply(function (_) {
    return [500, {
      message: 'not ok'
    }]
  })

  const wrapper = mount(Apps,
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
          $route: { path: '/apps' },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(error).toBe('not ok')

  wrapper.unmount()
})
