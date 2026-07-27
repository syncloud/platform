import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Access from '../../src/views/Access.vue'

jest.setTimeout(30000)

function mountAccess (getData, onPost) {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200, { data: getData, success: true })
  if (onPost) {
    mock.onPost('/rest/access').reply(onPost)
  }
  const showError = jest.fn()
  const wrapper = mount(Access, {
    attachTo: document.body,
    global: {
      stubs: {
        Error: { template: '<span/>', methods: { showAxios: showError } },
        Dialog: true
      }
    }
  })
  return { wrapper, showError, mock }
}

function captured () {
  const state = {}
  const reply = (config) => {
    state.body = JSON.parse(config.data)
    return [200, { success: true }]
  }
  return { state, reply }
}

async function selectMode (wrapper, mode) {
  await wrapper.find(`[data-testid="ipv4-mode-${mode}"]`).trigger('click')
  await flushPromises()
}

test('ipv4 off saves relay and ipv4 disabled', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: true, ipv4: '111.111.111.111' }, reply)
  await flushPromises()
  await selectMode(wrapper, 'off')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.relay_enabled).toBe(false)
  expect(state.body.ipv4_enabled).toBe(false)
  wrapper.unmount()
})

test('ipv4 public saves ipv4 enabled and public', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'public')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.relay_enabled).toBe(false)
  expect(state.body.ipv4_enabled).toBe(true)
  expect(state.body.ipv4_public).toBe(true)
  wrapper.unmount()
})

test('ipv4 local saves ipv4 enabled not public', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'local')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.relay_enabled).toBe(false)
  expect(state.body.ipv4_enabled).toBe(true)
  expect(state.body.ipv4_public).toBe(false)
  wrapper.unmount()
})

test('ipv4 relay saves relay enabled, ipv4 disabled', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'relay')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.relay_enabled).toBe(true)
  expect(state.body.ipv4_enabled).toBe(false)
  wrapper.unmount()
})

test('existing public config loads into public tab', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: true, ipv4_public: true, ipv4: '111.111.111.111' }, reply)
  await flushPromises()
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv4_enabled).toBe(true)
  expect(state.body.ipv4_public).toBe(true)
  wrapper.unmount()
})

test('existing local config loads into local tab', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: true, ipv4_public: false }, reply)
  await flushPromises()
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv4_enabled).toBe(true)
  expect(state.body.ipv4_public).toBe(false)
  wrapper.unmount()
})

test('relay and ipv6 coexist', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false, ipv6_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'relay')
  await wrapper.find('#tgl_ipv6_enabled').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.relay_enabled).toBe(true)
  expect(state.body.ipv6_enabled).toBe(true)
  wrapper.unmount()
})

test('ipv6 enable', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv6_enabled: false }, reply)
  await flushPromises()
  await wrapper.find('#tgl_ipv6_enabled').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv6_enabled).toBe(true)
  wrapper.unmount()
})

test('ipv6 disable', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv6_enabled: true }, reply)
  await flushPromises()
  await wrapper.find('#tgl_ipv6_enabled').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv6_enabled).toBe(false)
  wrapper.unmount()
})

test('public auto detect omits ipv4', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'public')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv4).toBeUndefined()
  wrapper.unmount()
})

test('public manual ip includes ipv4', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false, ipv4: '111.111.111.111' }, reply)
  await flushPromises()
  await selectMode(wrapper, 'public')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.ipv4).toBe('111.111.111.111')
  wrapper.unmount()
})

test('public access port set', async () => {
  const { state, reply } = captured()
  const { wrapper, showError } = mountAccess({ ipv4_enabled: false }, reply)
  await flushPromises()
  await selectMode(wrapper, 'public')
  await wrapper.find('#access_port').setValue(10000)
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(0)
  expect(state.body.access_port).toBe(10000)
  wrapper.unmount()
})

test('public empty ip shows error', async () => {
  const showError = jest.fn()
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200, { data: { ipv4_enabled: false, ipv4: '111.111.111.111' }, success: true })
  const wrapper = mount(Access, {
    attachTo: document.body,
    global: { stubs: { Error: { template: '<span/>', methods: { showAxios: showError } }, Dialog: true } }
  })
  await flushPromises()
  await selectMode(wrapper, 'public')
  await wrapper.find('#ipv4').setValue(' ')
  await wrapper.find('#access_port').setValue(443)
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('save http error shows error', async () => {
  const showError = jest.fn()
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(200, { data: { ipv4_enabled: false }, success: true })
  mock.onPost('/rest/access').reply(500)
  const wrapper = mount(Access, {
    attachTo: document.body,
    global: { stubs: { Error: { template: '<span/>', methods: { showAxios: showError } }, Dialog: true } }
  })
  await flushPromises()
  await selectMode(wrapper, 'relay')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})

test('get http error shows error', async () => {
  const showError = jest.fn()
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/access').reply(500)
  const wrapper = mount(Access, {
    attachTo: document.body,
    global: { stubs: { Error: { template: '<span/>', methods: { showAxios: showError } }, Dialog: true } }
  })
  await flushPromises()
  expect(showError).toHaveBeenCalledTimes(1)
  wrapper.unmount()
})
