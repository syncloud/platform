import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Storage from '../../src/views/Storage.vue'

jest.setTimeout(30000)

const emptyHealth = { disks: [], scrub: null, schedule: { scrub: false, smart: false }, busy: false }

const arrayDisks = [
  {
    name: 'CT2000BX500SSD1',
    device: '/dev/sda',
    active: true,
    size: '1.8T',
    raid: 'raid1',
    raid_metadata: 'raid1',
    serial: '2112E58EB6C6',
    hctl: '0:0:0:0',
    devid: '1',
    errors: { write: 0, read: 0, flush: 0, corruption: 0, generation: 0 },
    has_errors: false,
    partitions: []
  },
  {
    name: 'External_USB_3.0',
    device: '/dev/sdb',
    active: true,
    size: '1.8T',
    raid: 'raid1',
    raid_metadata: 'raid1',
    serial: '20170617007242F',
    hctl: '1:0:0:0',
    devid: '2',
    errors: { write: 0, read: 12, flush: 0, corruption: 3, generation: 0 },
    has_errors: true,
    partitions: []
  }
]

const arrayHealth = (overrides) => Object.assign({
  disks: [
    {
      device: '/dev/sda',
      port: 'ata1',
      smart: {
        available: true,
        passed: true,
        reallocated: 0,
        pending: 0,
        offline_uncorrectable: 0,
        crc_errors: 0,
        power_on_hours: 212,
        temperature: 31,
        self_test: { running: false, remaining_percent: 0, last_type: 'Short offline', last_status: 'Completed without error', last_passed: true }
      },
      link: { count: 0 },
      io: { count: 0 }
    },
    {
      device: '/dev/sdb',
      port: 'ata2',
      smart: { available: false, message: 'Unknown USB bridge', self_test: { running: false, remaining_percent: 0 } },
      link: { count: 4, last: '2026-10-31T02:01:44' },
      io: { count: 3, last: '2026-10-31T02:01:45' }
    }
  ],
  scrub: { state: 'finished', started: 'Sat Nov  7 21:40:51 2026', duration: '1:10:02', errors: 'read=12 csum=3', corrected: 0, uncorrectable: 15 },
  schedule: { scrub: false, smart: true },
  busy: false
}, overrides)

const mountStorage = (showError) => mount(Storage,
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
        Dialog: {
          template: '<span :id="id"><slot name="text"></slot></span>',
          props: { id: String },
          methods: {
            show () {
            }
          }
        }
      }
    }
  }
)

const mockArray = (health) => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/disks').reply(200, { success: true, data: arrayDisks })
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: health })
  mock.onGet('/rest/storage/error/last').reply(200, { success: true, data: 'OK' })
  mock.onGet('/rest/job/status').reply(200, { success: true, data: { name: 'test' } })
  return mock
}

test('Activate partition', async () => {
  let deviceAction = ''
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/activate/partition').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [200, { success: true }]
  })
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#multi').trigger('click')
  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await expect(wrapper.find('#format').isVisible()).toBe(true)
  await wrapper.find('#confirmation').trigger('confirm')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

test('Deactivate partition', async () => {
  let deactivated = false
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/deactivate').reply(function (_) {
    deactivated = true
    return [200, { success: true }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#none').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await expect(wrapper.find('#format').isVisible()).toBe(false)
  await wrapper.find('#confirmation').trigger('confirm')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(deactivated).toBe(true)
  wrapper.unmount()
})

test('Activate partition error', async () => {
  let deviceAction = ''
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/activate/partition').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [500, { message: 'not ok' }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('not ok')
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

test('Activate partition service error', async () => {
  let deviceAction = ''
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/activate/partition').reply(function (config) {
    deviceAction = JSON.parse(config.data).device
    return [200, { success: false, message: 'not ok' }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#partition_1_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('not ok')
  expect(deviceAction).toBe('/dev/sdc1')
  wrapper.unmount()
})

test('Activate disks', async () => {
  let devices = []
  let error = ''
  const showError = (err) => {
    error = err
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/activate/disk').reply(function (config) {
    devices = JSON.parse(config.data).devices
    return [200, { success: true }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#disk_0').trigger('click')
  await wrapper.find('#disk_1').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  await expect(wrapper.find('#format').isVisible()).toBe(true)
  await wrapper.find('#confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('')
  expect(devices).toEqual(['/dev/sdb', '/dev/sdc'])
  wrapper.unmount()
})

test('Activate disks error', async () => {
  let error = ''
  const showError = (err) => {
    error = err.response.data.message
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').replyOnce(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: true,
          size: '2G',
          partitions: []
        }
      ],
      success: true
    }
  ).onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: []
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/activate/disk').reply(function (_) {
    return [500, { success: false, message: 'not ok' }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  expect(wrapper.find('#disk_0').classes()).toContain('is-selected')
  await wrapper.find('#btn_save').trigger('click')
  await wrapper.find('#confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('not ok')
  expect(wrapper.find('#disk_0').classes()).not.toContain('is-selected')
  wrapper.unmount()
})

test('Deactivate disks', async () => {
  let deactivated = false
  let error = ''
  const showError = (err) => {
    error = err
  }
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: true,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onPost('/rest/storage/deactivate').reply(function (_) {
    deactivated = true
    return [200, { success: true }]
  })
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<span :id="id"><slot name="text"></slot></span>',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#disk_0').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  await expect(wrapper.find('#format').isVisible()).toBe(false)
  await wrapper.find('#confirmation').trigger('confirm')

  await flushPromises()

  expect(error).toBe('')
  expect(deactivated).toBe(true)
  wrapper.unmount()
})

test('Show single partition', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: true, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await expect(wrapper.find('#partition_0_0').classes()).toContain('is-selected')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Show single partition none', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#multi').trigger('click')
  await expect(wrapper.find('#none').classes()).toContain('is-selected')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Show multi disk', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: [
        {
          name: 'Name1',
          device: '/dev/sdb',
          active: true,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdb1', size: '931.5G' }
          ]
        },
        {
          name: 'Name2',
          device: '/dev/sdc',
          active: false,
          size: '2G',
          partitions: [
            { active: false, device: '/dev/sdc1', size: '931.5G' }
          ]
        }
      ],
      success: true
    }
  )
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await expect(wrapper.find('#multi').attributes('aria-checked')).toBe('true')
  await expect(wrapper.find('#disk_0').classes()).toContain('is-selected')
  await expect(wrapper.find('#disk_1').classes()).not.toContain('is-selected')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Show null disk', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: emptyHealth })
  mock.onGet('/rest/storage/disks').reply(200,
    {
      data: null,
      success: true
    }
  )
  mock.onGet('/rest/storage/error/last').reply(200,
    { success: true, data: 'OK' }
  )
  mock.onGet('/rest/job/status').reply(200,
    { success: true, data: { name: 'test' } }
  )
  const wrapper = mount(Storage,
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
          Dialog: {
            template: '<button :id="id" />',
            props: { id: String },
            methods: {
              show () {
              }
            }
          }
        }
      }
    }
  )

  await flushPromises()

  await expect(wrapper.find('#multi').attributes('aria-checked')).toBe('true')
  await expect(wrapper.find('#no_disks').text()).toBe('No external disks found')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Array view shows each disk with serial, port, counters, smart and kernel events', async () => {
  const showError = jest.fn()
  mockArray(arrayHealth())
  const wrapper = mountStorage(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="disk-serial-0"]').text()).toBe('2112E58EB6C6')
  expect(wrapper.find('[data-testid="disk-port-0"]').text()).toBe('ata1 (0:0:0:0)')
  expect(wrapper.find('[data-testid="disk-devid-1"]').text()).toBe('2')
  expect(wrapper.find('[data-testid="disk-profile-0"]').text()).toBe('data raid1, metadata raid1')
  expect(wrapper.find('[data-testid="disk-errors-1"]').text()).toBe('write 0, read 12, flush 0, corruption 3, generation 0')
  expect(wrapper.find('[data-testid="disk-errors-1"]').classes()).toContain('disk-bad')
  expect(wrapper.find('[data-testid="disk-has-errors-1"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="disk-has-errors-0"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="disk-smart-0"]').text()).toBe('Passed, 212 h powered on, 31 °C. Reallocated 0, pending 0, uncorrectable 0, cable/connector errors 0')
  expect(wrapper.find('[data-testid="disk-selftest-0"]').text()).toContain('Short offline: Completed without error')
  expect(wrapper.find('[data-testid="disk-smart-1"]').text()).toBe('Not available for this disk')
  expect(wrapper.find('[data-testid="disk-smart-test-1"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="disk-link-0"]').text()).toBe('None')
  expect(wrapper.find('[data-testid="disk-link-1"]').text()).toBe('4, last at 2026-10-31T02:01:44')
  expect(wrapper.find('[data-testid="disk-io-1"]').text()).toBe('3, last at 2026-10-31T02:01:45')
  expect(wrapper.find('[data-testid="scrub-status"]').text()).toBe('Finished, started Sat Nov  7 21:40:51 2026, took 1:10:02. Errors: read=12 csum=3, corrected 0, uncorrectable 15')
  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(false)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Run read-only scrub', async () => {
  const showError = jest.fn()
  let scrubs = 0
  const mock = mockArray(arrayHealth())
  mock.onPost('/rest/storage/check/scrub').reply(function () {
    scrubs++
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountStorage(showError)
  await flushPromises()

  await wrapper.find('[data-testid="scrub-start"]').trigger('click')
  await flushPromises()

  expect(scrubs).toBe(1)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Scrub refused during an app update shows the error', async () => {
  const showError = jest.fn()
  const mock = mockArray(arrayHealth())
  mock.onPost('/rest/storage/check/scrub').reply(500, { success: false, message: 'an app install or update is in progress, try again when it finishes' })
  const wrapper = mountStorage(showError)
  await flushPromises()

  await wrapper.find('[data-testid="scrub-start"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Run SMART short test on one disk', async () => {
  const showError = jest.fn()
  let device = ''
  const mock = mockArray(arrayHealth())
  mock.onPost('/rest/storage/check/smart').reply(function (config) {
    device = JSON.parse(config.data).device
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountStorage(showError)
  await flushPromises()

  await wrapper.find('[data-testid="disk-smart-test-0"]').trigger('click')
  await flushPromises()

  expect(device).toBe('/dev/sda')
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Reset error counters after confirmation', async () => {
  const showError = jest.fn()
  let resets = 0
  const mock = mockArray(arrayHealth())
  mock.onPost('/rest/storage/errors/reset').reply(function () {
    resets++
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountStorage(showError)
  await flushPromises()

  await wrapper.find('[data-testid="errors-reset"]').trigger('click')
  expect(resets).toBe(0)
  await wrapper.find('[data-testid="reset-confirmation"]').trigger('confirm')
  await flushPromises()

  expect(resets).toBe(1)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Schedule monthly scrub', async () => {
  const showError = jest.fn()
  let saved = null
  const mock = mockArray(arrayHealth())
  mock.onPost('/rest/storage/check/schedule').reply(function (config) {
    saved = JSON.parse(config.data)
    return [200, { success: true, data: saved }]
  })
  const wrapper = mountStorage(showError)
  await flushPromises()

  await wrapper.find('[data-testid="schedule-scrub"]').trigger('click')
  await flushPromises()

  expect(saved).toEqual({ scrub: true, smart: true })
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Checks are disabled while another operation runs', async () => {
  const showError = jest.fn()
  mockArray(arrayHealth({ busy: true }))
  const wrapper = mountStorage(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="scrub-start"]').attributes('disabled')).toBeDefined()
  expect(wrapper.find('[data-testid="disk-smart-test-0"]').attributes('disabled')).toBeDefined()
  wrapper.unmount()
})

test('No array section without an active btrfs array', async () => {
  const showError = jest.fn()
  mockArray(arrayHealth({ scrub: null }))
  const wrapper = mountStorage(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="array-section"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="schedule-smart-row"]').exists()).toBe(true)
  wrapper.unmount()
})
