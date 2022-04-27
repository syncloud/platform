import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Access from '@/views/Access'

jest.setTimeout(30000)

test('Private ipv4 disable', async () => {
  let savedIpv4Enabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: true,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    savedIpv4Enabled = JSON.parse(config.data).ipv4_enabled
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedIpv4Enabled).toBe(false)
  wrapper.unmount()
})

test('Private ipv4 enable', async () => {
  let savedIpv4Enabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    savedIpv4Enabled = JSON.parse(config.data).ipv4_enabled
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedIpv4Enabled).toBe(true)
  wrapper.unmount()
})

test('Public ipv4 enable', async () => {
  let savedIpv4Enabled
  let savedIpv4Public
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    const request = JSON.parse(config.data)
    savedIpv4Enabled = request.ipv4_enabled
    savedIpv4Public = request.ipv4_public
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
  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedIpv4Enabled).toBe(true)
  expect(savedIpv4Public).toBe(true)
  wrapper.unmount()
})

test('Public ipv4 auto detect', async () => {
  let ipAutoDetectEnabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    ipAutoDetectEnabled = JSON.parse(config.data).ipv4 === undefined
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
  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#tgl_ip_autodetect').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(ipAutoDetectEnabled).toBe(true)
  wrapper.unmount()
})

test('Ipv6 enable', async () => {
  let savedIpv6Enabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv6_enabled: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    savedIpv6Enabled = JSON.parse(config.data).ipv6_enabled
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
  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv6_enabled').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedIpv6Enabled).toBe(true)
  wrapper.unmount()
})

test('Ipv6 disable', async () => {
  let savedIpv6Enabled
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv6_enabled: true,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    savedIpv6Enabled = JSON.parse(config.data).ipv6_enabled
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
  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv6_enabled').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedIpv6Enabled).toBe(false)
  wrapper.unmount()
})

test('Public access port set', async () => {
  let savedAccessPort
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (config) {
    const request = JSON.parse(config.data)
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
  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(savedAccessPort).toBe(443)
  wrapper.unmount()
})

test('Save http error', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Save service error', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Access port wrong', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(0)
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toBe('Access port (0) has to be between 1 and 65535')
  wrapper.unmount()
})

test('Access port 443 no warning', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: true,
        ipv4_public: true,
        ipv4: '111.111.111.111',
        access_port: 443
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  expect(wrapper.find('#access_port_warning').isVisible()).toBe(false)
  expect(error).toBe('')

  wrapper.unmount()
})

test('Access port non 443 shows warning', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: true,
        ipv4_public: true,
        ipv4: '111.111.111.111',
        access_port: 444
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  expect(wrapper.find('#access_port_warning').isVisible()).toBe(true)
  expect(error).toBe('')

  wrapper.unmount()
})

test('Access port is always 443 in ipv4 private', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false,
        ipv4: '111.111.111.111'
      },
      success: true
    }
  )

  let savedAccessPort
  mock.onPost('/rest/access').reply(function (config) {
    const request = JSON.parse(config.data)
    savedAccessPort = request.access_port
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#access_port').setValue(0)
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toBe('')
  expect(savedAccessPort).toBe(443)
  wrapper.unmount()
})

test('Manual Ipv4 default value', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#tgl_ip_autodetect').trigger('toggle')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toBe('Empty IP')
  wrapper.unmount()
})

test('Manual Ipv4 empty value', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200,
    {
      data: {
        ipv4_enabled: false,
        ipv4_public: false
      },
      success: true
    }
  )

  mock.onPost('/rest/access').reply(function (_) {
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

  await wrapper.find('#tgl_ipv4_enabled').trigger('toggle')
  await wrapper.find('#tgl_ipv4_public').trigger('toggle')
  await wrapper.find('#tgl_ip_autodetect').trigger('toggle')
  await wrapper.find('#ipv4').setValue(' ')
  await wrapper.find('#btn_save').trigger('click')

  await flushPromises()

  expect(error).toBe('Empty IP')
  wrapper.unmount()
})
