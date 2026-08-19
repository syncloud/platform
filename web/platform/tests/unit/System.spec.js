import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import System from '../../src/views/System.vue'

jest.setTimeout(30000)

let reload

beforeEach(() => {
  reload = jest.fn()
  Object.defineProperty(window, 'location', { configurable: true, value: { reload } })
})

function mountSystem (probes, postStatus = 200, attempts = 10) {
  const mock = new MockAdapter(axios)
  let index = 0
  mock.onGet('/rest/uptime').reply(() => {
    const probe = probes[Math.min(index, probes.length - 1)]
    index++
    if (probe === 'down') {
      return Promise.reject(new Error('network error'))
    }
    if (typeof probe === 'object') {
      return [probe.status, {}]
    }
    return [200, { success: true, data: probe }]
  })
  mock.onPost('/rest/restart').reply(postStatus, { success: postStatus === 200 })
  mock.onPost('/rest/shutdown').reply(postStatus, { success: postStatus === 200 })
  const showError = jest.fn()
  const wrapper = mount(System, {
    attachTo: document.body,
    props: { pollIntervalMs: 0, pollAttempts: attempts },
    global: {
      stubs: {
        Error: { template: '<span/>', methods: { showAxios: showError } }
      }
    }
  })
  return { wrapper, showError, mock }
}

async function confirm (wrapper, button) {
  await wrapper.find(button).trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')
  for (let i = 0; i < 100; i++) {
    await flushPromises()
  }
}

test('restart shows progress and reloads when the device is back', async () => {
  const { wrapper, showError } = mountSystem([1000, 1000, 'down', 50, 50])
  await confirm(wrapper, '#restart')
  expect(showError).toHaveBeenCalledTimes(0)
  expect(wrapper.find('[data-testid="progress"]').exists()).toBe(true)
  expect(wrapper.vm.stage).toBe('done')
  expect(reload).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('restart is detected by uptime drop without seeing the device offline', async () => {
  const { wrapper } = mountSystem([1000, 1000, 20])
  await confirm(wrapper, '#restart')
  expect(wrapper.vm.stage).toBe('done')
  expect(reload).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('progress dialog is shown while the device is restarting', async () => {
  const { wrapper } = mountSystem([1000])
  await wrapper.find('#restart').trigger('click')
  await wrapper.find('#btn_confirm').trigger('click')
  await flushPromises()
  expect(wrapper.find('[data-testid="progress"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="progress-spinner"]').exists()).toBe(true)
  expect(wrapper.find('#btn_close').exists()).toBe(false)
  expect(wrapper.find('[data-testid="progress-text"]').text()).toBe('The device is going offline.')
  wrapper.unmount()
})

test('unauthorized response means the device is back online', async () => {
  const { wrapper } = mountSystem([1000, 'down', { status: 401 }])
  await confirm(wrapper, '#restart')
  expect(wrapper.vm.stage).toBe('done')
  expect(reload).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('restart offers a reload when the device does not come back', async () => {
  const { wrapper, showError } = mountSystem([1000], 200, 2)
  await confirm(wrapper, '#restart')
  expect(wrapper.vm.stage).toBe('timeout')
  expect(reload).toHaveBeenCalledTimes(0)
  expect(wrapper.find('#btn_reload').exists()).toBe(true)
  expect(showError).toHaveBeenCalledTimes(0)
  await wrapper.find('#btn_reload').trigger('click')
  expect(reload).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('shutdown waits until the device stops responding', async () => {
  const { wrapper } = mountSystem([1000, 'down'])
  await confirm(wrapper, '#shutdown')
  expect(wrapper.vm.stage).toBe('off')
  expect(reload).toHaveBeenCalledTimes(0)
  expect(wrapper.find('[data-testid="progress-spinner"]').exists()).toBe(false)
  expect(wrapper.find('#btn_close').exists()).toBe(true)
  await wrapper.find('#btn_close').trigger('click')
  expect(wrapper.find('[data-testid="progress"]').exists()).toBe(false)
  wrapper.unmount()
})

test('restart error hides the progress and shows the error', async () => {
  const { wrapper, showError } = mountSystem([1000], 500)
  await confirm(wrapper, '#restart')
  expect(showError).toHaveBeenCalledTimes(1)
  expect(wrapper.find('[data-testid="progress"]').exists()).toBe(false)
  expect(reload).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})
