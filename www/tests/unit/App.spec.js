import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import App from '../../src/views/App.vue'
import { ElButton, ElCol, ElProgress, ElRow } from 'element-plus'

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
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#app_confirmation').trigger('confirm')
  await flushPromises()

  expect(wrapper.find('#progress').exists()).toBe(false)

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(app).toBe('files')
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Install of the same app is already in progress on open', async () => {
  const showError = jest.fn()
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

  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, {
      success: true,
      data: {
        is_running: true,
        progress: {
          files: {
            app: 'files',
            summary: 'Downloading'
          }
        }
      }
    }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: true
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
  expect(wrapper.find('#btn_install').exists()).toBe(false)
  expect(wrapper.find('#app_name').text()).toBe('Files')
  expect(wrapper.find('#progress_summary').text()).toBe('Downloading')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(statusCalled).toBeTruthy()
  wrapper.unmount()
})

test('Install of different app is already in progress on open', async () => {
  const showError = jest.fn()
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

  let statusCalled = false
  mock.onGet('/rest/installer/status').reply(function (_) {
    statusCalled = true
    return [200, {
      success: true,
      data: {
        is_running: true,
        progress: {
          another_app: {
            app: 'another_app',
            summary: 'Downloading'
          }
        }
      }
    }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: true
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
  expect(wrapper.find('#btn_install').exists()).toBe(true)

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
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
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#app_confirmation').trigger('confirm')
  await flushPromises()

  expect(wrapper.find('#progress').exists()).toBe(false)

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
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#app_confirmation').trigger('confirm')
  await flushPromises()

  expect(wrapper.find('#progress').exists()).toBe(false)

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

  mock.onGet('/rest/installer/status').reply(function (_) {
    return [200, { success: true, data: { is_running: false } }]
  })

  mock.onPost('/rest/app/remove').reply(function (_) {
    return [500, { message: 'not ok' }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#app_confirmation').trigger('confirm')

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
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: true
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
  mock.onGet('/rest/installer/status').reply(function (_) {
    return [200, { success: true, data: { is_running: false } }]
  })
  mock.onPost('/rest/backup/create').reply(function (config) {
    app = JSON.parse(config.data).app
    return [200, { success: true }]
  })

  let statusCalled = false
  mock.onGet('/rest/job/status').reply(function (_) {
    statusCalled = true
    return [200, { success: true, data: { status: 'Idle' } }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#backup_confirmation').trigger('confirm')
  await flushPromises()

  expect(wrapper.find('#progress').exists()).toBe(false)

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
  mock.onGet('/rest/installer/status').reply(function (_) {
    return [200, { success: true, data: { is_running: false } }]
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
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  // await wrapper.find('#btn_backup_confirm').trigger('click')
  await wrapper.find('#backup_confirmation').trigger('confirm')

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
  mock.onGet('/rest/installer/status').reply(function (_) {
    return [200, { success: true, data: { is_running: false } }]
  })

  mock.onPost('/rest/backup/create').reply(function (_) {
    return [200, { success: false, message: 'not ok' }]
  })

  const wrapper = mount(App,
    {
      attachTo: document.body,
      global: {
        stubs: {
          'el-row': ElRow,
          'el-col': ElCol,
          'el-progress': ElProgress,
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
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
  await wrapper.find('#backup_confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('not ok')
  wrapper.unmount()
})
