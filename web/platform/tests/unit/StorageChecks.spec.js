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

test('Scrub status and schedules', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="scrub-status"]').text()).toBe('Finished, started Sat Nov  7 21:40:51 2026, took 1:10:02. Errors: read=12 csum=3, corrected 0, uncorrectable 15')
  expect(wrapper.find('[data-testid="errors-section"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(false)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Run read-only scrub', async () => {
  const showError = jest.fn()
  let scrubs = 0
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check/scrub').reply(function () {
    scrubs++
    return [200, { success: true, data: 'OK' }]
  })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="scrub-start"]').trigger('click')
  await flushPromises()

  expect(scrubs).toBe(1)
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Scrub refused during an app update shows the error', async () => {
  const showError = jest.fn()
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check/scrub').reply(500, { success: false, message: 'an app install or update is in progress, try again when it finishes' })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="scrub-start"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
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

test('Schedule monthly scrub', async () => {
  const showError = jest.fn()
  let saved = null
  const mock = mockApi(arrayHealth())
  mock.onPost('/rest/storage/check/schedule').reply(function (config) {
    saved = JSON.parse(config.data)
    return [200, { success: true, data: saved }]
  })
  const wrapper = mountChecks(showError)
  await flushPromises()

  await wrapper.find('[data-testid="schedule-scrub"]').trigger('click')
  await flushPromises()
  expect(saved).toBe(null)

  await wrapper.find('[data-testid="checks-save"]').trigger('click')
  await flushPromises()

  expect(saved).toEqual({ scrub: true, smart: true })
  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})

test('Checks are disabled while another operation runs', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth({ busy: true }))
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="checks-busy"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="scrub-start"]').attributes('disabled')).toBeDefined()
  wrapper.unmount()
})

test('Without a btrfs array only the SMART schedule is offered', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth({ scrub: null }))
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="array-section"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="errors-section"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="schedule-smart"]').exists()).toBe(true)
  wrapper.unmount()
})

test('Scrub and SMART explain what they do', async () => {
  const showError = jest.fn()
  mockApi(arrayHealth())
  const wrapper = mountChecks(showError)
  await flushPromises()

  expect(wrapper.find('[data-testid="scrub-help-dialog"]').exists()).toBe(false)
  await wrapper.find('[data-testid="scrub-help"]').trigger('click')
  expect(wrapper.find('[data-testid="scrub-help-dialog"]').text()).toContain('read-only')

  await wrapper.find('[data-testid="smart-help"]').trigger('click')
  expect(wrapper.find('[data-testid="smart-help-dialog"]').text()).toContain('short self-test')
  wrapper.unmount()
})
