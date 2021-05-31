import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Access from '@/views/Access'

jest.setTimeout(30000)

const defaultPortMappings = {
  port_mappings: [
    { local_port: 80, external_port: 80 },
    { local_port: 443, external_port: 10001 }
  ],
  success: true
}

test('Disable external access', async () => {
  let savedExternalAccess
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: true,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    savedExternalAccess = JSON.parse(config.data).external_access
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: { id: String }
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedExternalAccess).toBe(false)
  wrapper.unmount()
})

test('Enable external access', async () => {
  let savedExternalAccess
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    savedExternalAccess = JSON.parse(config.data).external_access
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: { id: String }
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedExternalAccess).toBe(true)
  wrapper.unmount()
})

test('Enable external access empty manual ports', async () => {
  let savedExternalAccess
  let savedAccessPort = -1
  let savedCertificatePort = -1
  let errorCount = 0
  const showError = jest.fn()
  const setAccess = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: false,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, {
    port_mappings: [ ],
    success: true
  })
  mock.onPost('/rest/access/set_access').reply(function (config) {
    setAccess()
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: { id: String }
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(1)
  expect(setAccess).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Enable auto port mapping (upnp)', async () => {
  let savedUpnpEnabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: false,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    savedUpnpEnabled = JSON.parse(config.data).upnp_enabled
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: { id: String }
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()
  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#tgl_upnp').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedUpnpEnabled).toBe(true)
  wrapper.unmount()
})

test('Ip auto detect', async () => {
  let ipAutoDetectEnabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: false,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    ipAutoDetectEnabled = JSON.parse(config.data).public_ip === undefined
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()
  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#tgl_ip_autodetect').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(ipAutoDetectEnabled).toBe(true)
  wrapper.unmount()
})

test('Set access and certificate ports', async () => {
  let savedCertificatePort
  let savedAccessPort
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: false,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    const request = JSON.parse(config.data)
    savedCertificatePort = request.certificate_port
    savedAccessPort = request.access_port
    return [200, { success: true }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()
  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#certificate_port').setValue('1')
  await wrapper.find('#access_port').setValue(2)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedCertificatePort).toBe(1)
  expect(savedAccessPort).toBe(2)
  wrapper.unmount()
})

test('Enable external access http error', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    return [400, { }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Enable external access service error', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    return [200, { success: false }]
  })

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Enable external wrong access port', async () => {

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    return [200, { success: true }]
  })

  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#tgl_upnp').trigger('toggle')
  await wrapper.find('#access_port').setValue(0)
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toContain('access port')
  wrapper.unmount()
})

test('Enable external wrong certificate port', async () => {

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access/access').reply(200,
    {
      data: {
        external_access: false,
        upnp_available: false,
        upnp_enabled: true,
        public_ip: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onGet('/rest/access/port_mappings').reply(200, defaultPortMappings)
  mock.onPost('/rest/access/set_access').reply(function (config) {
    return [200, { success: true }]
  })

  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }

  const wrapper = mount(Access,
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
          Switch: {
            template: '<button :id="id" />',
            props: ['id']
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#tgl_external').trigger('toggle')
  await wrapper.find('#tgl_upnp').trigger('toggle')
  await wrapper.find('#certificate_port').setValue(0)
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toContain('certificate port')
  wrapper.unmount()
})
