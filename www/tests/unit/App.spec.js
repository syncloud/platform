import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import App from '../../src/views/App.vue'

jest.setTimeout(30000)

test('Install', async () => {
  const showError = jest.fn()
  let app
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: null
      },
      success: true
    }]
  })

  mock.onPost('/rest/app/install').reply(function (config) {
    app = JSON.parse(config.data).app_id
    return [200, { success: true }]
  })

  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, { success: true, data: { is_running: false } }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#btn_upgrade').exists()).toBe(false)
  expect(wrapper.find('#btn_remove').exists()).toBe(false)

  await wrapper.find('#btn_install').trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(app).toBe('files')
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Upgrade', async () => {
  const showError = jest.fn()
  let app
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/app/upgrade').reply(function (config) {
    app = JSON.parse(config.data).app_id
    return [200, { success: true }]
  })
  
  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, { success: true, data: { is_running: false } }]
  })
  
  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#btn_install').exists()).toBe(false)
  expect(wrapper.find('#btn_remove').exists()).toBe(true)

  await wrapper.find('#btn_upgrade').trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(app).toBe('files')
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Remove', async () => {
  const showError = jest.fn()
  let app
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/app/remove').reply(function (config) {
    app = JSON.parse(config.data).app_id
    return [200, { success: true }]
  })
  
  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, { success: true, data: { is_running: false } }]
  })
  
  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#btn_install').exists()).toBe(false)
  expect(wrapper.find('#btn_upgrade').exists()).toBe(true)

  await wrapper.find('#btn_remove').trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(app).toBe('files')
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Action error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/app/remove').reply(function (_) {
    return [500, { message: 'not ok' }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_remove').trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')

  await flushPromises()

  expect(error).toBe('not ok')
  wrapper.unmount()
})

test('Show error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [500, {
      success: false,
      message: 'not ok'
    }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(error).toBe('not ok')
  wrapper.unmount()
})

test('Backup', async () => {
  const showError = jest.fn()
  let app
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/backup/create').reply(function (config) {
    app = JSON.parse(config.data).app
    return [200, { success: true }]
  })
  
  let statusCalled = false
  mock.onGet('/rest/job/status').reply(function (_) {
    statusCalled = true
    return [200, {success: true, data:{status: "Idle"}}]
  })
  
  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)

  await wrapper.find('#btn_backup').trigger('click')
  await wrapper.find('#btn_backup_confirm').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(app).toBe('files')
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Backup error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/backup/create').reply(function (_) {
    return [500, {
      message: 'not ok'
    }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_backup').trigger('click')
  await wrapper.find('#btn_backup_confirm').trigger('click')

  await flushPromises()

  expect(error).toBe('not ok')
  wrapper.unmount()
})

test('Backup service error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/app').reply(function (_) {
    return [200, {
      data: {
        app: { id: 'files', name: 'Files', required: false, ui: false, url: 'http://files.odroid-c2.syncloud.it' },
        current_version: '2',
        installed_version: '1'
      },
      success: true
    }]
  })

  mock.onPost('/rest/backup/create').reply(function (_) {
    return [200, { success: false, message: 'not ok' }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#btn_backup').trigger('click')
  await wrapper.find('#btn_backup_confirm').trigger('click')

  await flushPromises()

  expect(error).toBe('not ok')
  wrapper.unmount()
})
