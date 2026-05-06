import { mount, RouterLinkStub } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import AppCenter from '../../src/views/AppCenter.vue'

jest.setTimeout(30000)

async function setFilter (wrapper, value) {
  wrapper.vm.filter = value
  await wrapper.vm.$nextTick()
}

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

test('Filter searches id, name, and description', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('rest/apps/available').reply(200,
    {
      data: [
        { id: 'files', name: 'File browser', description: 'Browse files', icon: '/images/files.png' },
        { id: 'nextcloud', name: 'Nextcloud file sharing', description: 'Sync and share', icon: '/images/nc.png' },
        { id: 'photoprism', name: 'PhotoPrism', description: 'Photos and videos', icon: '/images/pp.png' }
      ]
    }
  )

  const mockRouter = { push: jest.fn() }
  const wrapper = mount(AppCenter,
    {
      attachTo: document.body,
      global: {
        components: { RouterLink: RouterLinkStub },
        stubs: { Error: { template: '<span/>', methods: { showAxios: () => {} } } },
        mocks: { $route: { path: '/appcenter' }, $router: mockRouter }
      }
    }
  )
  await flushPromises()

  // matches by id (snap name)
  await setFilter(wrapper,'files')
  expect(wrapper.text()).toContain('File browser')
  expect(wrapper.text()).not.toContain('PhotoPrism')

  // matches by display name
  await setFilter(wrapper,'browser')
  expect(wrapper.text()).toContain('File browser')
  expect(wrapper.text()).not.toContain('Nextcloud')
  expect(wrapper.text()).not.toContain('PhotoPrism')

  // matches by description
  await setFilter(wrapper,'photos')
  expect(wrapper.text()).toContain('PhotoPrism')
  expect(wrapper.text()).not.toContain('File browser')

  // case-insensitive
  await setFilter(wrapper,'NEXTCLOUD')
  expect(wrapper.text()).toContain('Nextcloud file sharing')

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
