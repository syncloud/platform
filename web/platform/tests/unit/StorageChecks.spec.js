import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import StorageChecks from '../../src/views/StorageChecks.vue'
import { arrayDisks, arrayHealth } from './storageFixtures'

jest.setTimeout(30000)

const mountChecks = (showError) => mount(StorageChecks,
  {
    attachTo: document.body,
    global: {
      stubs: {
        Error: { template: '<span/>', methods: { showAxios: showError } },
        Dialog: {
          template: '<span v-if="visible" :id="id"><slot name="text"></slot></span>',
          props: { id: String, visible: Boolean }
        },
        'router-link': { props: ['to'], template: '<a :href="to"><slot/></a>' }
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

test('Self-test and scrub results', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="selftest-status"]').text()).toBe('Passed on all disks')
  expect(wrapper.find('[data-testid="scrub-status"]').text()).toBe('Finished, started Sat Nov  7 21:40:51 2026, took 1:10:02. Errors: read=12 csum=3, corrected 0, uncorrectable 15')
  expect(wrapper.find('[data-testid="errors-section"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(false)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Failed self-test is shown in red', async () => {
  const showError = jest.fn()
  const health = arrayHealth()
  health.disks[0].smart.self_test = { running: false, remaining_percent: 0, last_type: 'Short offline', last_status: 'Completed: read failure', last_passed: false }
  mockApi(health)
  const wrapper = mountChecks(showError)
  await flushPromises()

  const status = wrapper.find('[data-testid="selftest-status"]')
  expect(status.text()).toBe('Failed on /dev/sda')
  expect(status.classes()).toContain('disk-bad')
  wrapper.unmount()
})

test('Check now starts both checks', async () => {
  const showError = jest.fn()
  let checks = 0
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check').reply(function () {
    checks++
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="checks-start"]').trigger('click')
  await flushPromises()

  expect(checks).toBe(1)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Check refused during an app update shows the error', async () => {
  const showError = jest.fn()
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check').reply(500, { success: false, message: 'an app install or update is in progress, try again when it finishes' })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="checks-start"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('Check now is disabled while a self-test runs', async () => {
  const showError = jest.fn()
  const health = arrayHealth()
  health.disks[0].smart.self_test = { running: true, remaining_percent: 60 }
  mockApi(health)
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="selftest-status"]').text()).toBe('Running on /dev/sda')
  expect(wrapper.find('[data-testid="checks-start"]').attributes('disabled')).toBeDefined()
  wrapper.unmount()
})

test('Reset error counters after confirmation', async () => {
  const showError = jest.fn()
  let resets = 0
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/errors/reset').reply(function () {
    resets++
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="errors-reset"]').trigger('click')
  expect(resets).toBe(0)
  await wrapper.find('[data-testid="reset-confirmation"]').trigger('confirm')
  await flushPromises()

  expect(resets).toBe(1)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Schedule automatic checks', async () => {
  const showError = jest.fn()
  let saved = null
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check/schedule').reply(function (config) {
    saved = JSON.parse(config.data)
    return [200, { success: true, data: saved }]
  })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="schedule-checks"]').trigger('click')
  await flushPromises()
  expect(saved).toBe(null)

  await wrapper.find('[data-testid="checks-save"]').trigger('click')
  await flushPromises()

  expect(saved).toEqual({ enabled: true })
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Checks are disabled while another operation runs', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth({ busy: true }))
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="checks-start"]').attributes('disabled')).toBeDefined()
  wrapper.unmount()
})

test('Without a btrfs array only the self-test is shown', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth({ scrub: null }))
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="scrub-status"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="errors-section"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="selftest-status"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="schedule-checks"]').exists()).toBe(true)
  wrapper.unmount()
})

test('Disk check explains what it does', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="checks-help-dialog"]').exists()).toBe(false)
  await wrapper.find('[data-testid="checks-help"]').trigger('click')
  const text = wrapper.find('[data-testid="checks-help-dialog"]').text()
  expect(text).toContain('SMART self-test')
  expect(text).toContain('read-only')
  wrapper.unmount()
})
