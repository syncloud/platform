import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import PasskeyList from '../../src/components/PasskeyList.vue'

const passkeys = [
  { id: 1, description: 'Firefox on Android', created: '2026-11-14T21:45:00Z', last_used: '2026-11-15T09:12:00Z' },
  { id: 2, description: 'Chrome on Windows', created: '2026-10-02T18:20:00Z', last_used: null }
]

function mountList (props = {}) {
  return mount(PasskeyList, { props, attachTo: document.body })
}

test('lists own passkeys', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/passkeys').reply(200, { success: true, data: passkeys })
  const wrapper = mountList()
  await flushPromises()

  expect(wrapper.find('[data-testid="passkey-name-1"]').text()).toBe('Firefox on Android')
  expect(wrapper.find('[data-testid="passkey-name-2"]').text()).toBe('Chrome on Windows')
  expect(wrapper.find('[data-testid="passkeys-remove-all"]').exists()).toBe(true)
  wrapper.unmount()
})

test('shows empty state', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/passkeys').reply(200, { success: true, data: [] })
  const wrapper = mountList()
  await flushPromises()

  expect(wrapper.find('[data-testid="passkeys-empty"]').exists()).toBe(true)
  wrapper.unmount()
})

test('removes own passkey after confirmation', async () => {
  const mock = new MockAdapter(axios)
  let current = [...passkeys]
  mock.onGet('/rest/passkeys').reply(() => [200, { success: true, data: current }])
  let removed
  mock.onPost('/rest/passkeys/remove').reply(config => {
    removed = JSON.parse(config.data)
    current = current.filter(p => p.id !== removed.id)
    return [200, { success: true }]
  })
  const wrapper = mountList()
  await flushPromises()

  await wrapper.find('[data-testid="passkey-remove-1"]').trigger('click')
  await document.querySelector('[data-testid="btn_confirm"]').click()
  await flushPromises()

  expect(removed).toEqual({ id: 1 })
  expect(wrapper.find('[data-testid="passkey-name-1"]').exists()).toBe(false)
  wrapper.unmount()
})

test('admin lists and resets a user passkeys', async () => {
  const mock = new MockAdapter(axios)
  let current = [...passkeys]
  let params
  mock.onGet('/rest/users/passkeys').reply(config => {
    params = config.params
    return [200, { success: true, data: current }]
  })
  let reset
  mock.onPost('/rest/users/passkeys/reset').reply(config => {
    reset = JSON.parse(config.data)
    current = []
    return [200, { success: true }]
  })
  const wrapper = mountList({ username: 'alice' })
  await flushPromises()
  expect(params).toEqual({ username: 'alice' })

  await wrapper.find('[data-testid="passkeys-remove-all"]').trigger('click')
  await document.querySelector('[data-testid="btn_confirm"]').click()
  await flushPromises()

  expect(reset).toEqual({ username: 'alice' })
  expect(wrapper.find('[data-testid="passkeys-empty"]').exists()).toBe(true)
  wrapper.unmount()
})
