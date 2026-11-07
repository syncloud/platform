import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import StorageDisk from '../../src/views/StorageDisk.vue'
import { arrayDisks, arrayHealth } from './storageFixtures'

jest.setTimeout(30000)

const mountDisk = (device, showError) => mount(StorageDisk,
  {
    attachTo: document.body,
    global: {
      stubs: {
        Error: { template: '<span/>', methods: { showAxios: showError } },
        'router-link': { props: ['to'], template: '<a :href="to"><slot/></a>' }
      },
      mocks: {
        $route: { path: '/storage/disk', query: { device: device } }
      }
    }
  }
)

const mockApi = (health) => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/disks').reply(200, { success: true, data: arrayDisks })
  mock.onGet('/rest/storage/health').reply(200, { success: true, data: health })
  return mock
}

test('Healthy disk shows serial, port, counters, smart and kernel events', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountDisk('/dev/sda', showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="disk-title"]').text()).toBe('CT2000BX500SSD1 - 1.8T')
  expect(wrapper.find('[data-testid="disk-device"]').text()).toBe('/dev/sda')
  expect(wrapper.find('[data-testid="disk-serial"]').text()).toBe('2112E58EB6C6')
  expect(wrapper.find('[data-testid="disk-port"]').text()).toBe('ata1 (0:0:0:0)')
  expect(wrapper.find('[data-testid="disk-devid"]').text()).toBe('1')
  expect(wrapper.find('[data-testid="disk-profile"]').text()).toBe('data raid1, metadata raid1')
  expect(wrapper.find('[data-testid="disk-errors"]').text()).toBe('write 0, read 0, flush 0, corruption 0, generation 0')
  expect(wrapper.find('[data-testid="disk-smart"]').text()).toBe('Passed, 212 h powered on, 31 °C. Reallocated 0, pending 0, uncorrectable 0, cable/connector errors 0')
  expect(wrapper.find('[data-testid="disk-selftest"]').text()).toBe('Short offline: Completed without error')
  expect(wrapper.find('[data-testid="disk-link"]').text()).toBe('None')
  expect(wrapper.find('[data-testid="disk-io"]').text()).toBe('None')
  expect(wrapper.find('[data-testid="storage-back"]').attributes('href')).toBe('/storage')
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Disk without SMART on a bad link shows counters and kernel events', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountDisk('/dev/sdb', showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="disk-errors"]').text()).toBe('write 0, read 12, flush 0, corruption 3, generation 0')
  expect(wrapper.find('[data-testid="disk-errors"]').classes()).toContain('disk-bad')
  expect(wrapper.find('[data-testid="disk-smart"]').text()).toBe('Not available for this disk')
  expect(wrapper.find('[data-testid="disk-smart-test"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="disk-link"]').text()).toBe('4, last at 2026-10-31T02:01:44')
  expect(wrapper.find('[data-testid="disk-io"]').text()).toBe('3, last at 2026-10-31T02:01:45')
  wrapper.unmount()
})

test('Run SMART short test', async () => {
  const showError = jest.fn()
  let device = ''
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check/smart').reply(function (config) {
    device = JSON.parse(config.data).device
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountDisk('/dev/sda', showError)
  await flushPromises()

  await wrapper.find('[data-testid="disk-smart-test"]').trigger('click')
  await flushPromises()

  expect(device).toBe('/dev/sda')
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('SMART test is disabled while another operation runs', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth({ busy: true }))
  const wrapper = mountDisk('/dev/sda', showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="disk-smart-test"]').attributes('disabled')).toBeDefined()
  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(true)
  wrapper.unmount()
})

test('Disconnected disk', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountDisk('/dev/sdz', showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="disk-missing"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="disk-device"]').exists()).toBe(false)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})
