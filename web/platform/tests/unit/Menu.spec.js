import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Menu from '../../src/components/Menu.vue'
import { useAuthStore } from '../../src/stores/auth'

function mountMenu (admin, space) {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/space').reply(200, { success: true, data: space })
  const pinia = createPinia()
  setActivePinia(pinia)
  const auth = useAuthStore()
  auth.loggedIn = true
  auth.admin = admin
  const wrapper = mount(Menu, {
    global: {
      plugins: [pinia],
      stubs: {
        'router-link': { props: ['to'], template: '<a :href="to"><slot/></a>' }
      }
    }
  })
  return { wrapper, mock, auth }
}

const enough = {
  low: false,
  mounts: [{ kind: 'system', path: '/', total_kb: 15 * 1024 * 1024, free_kb: 8 * 1024 * 1024, low: false }]
}

const full = {
  low: true,
  mounts: [
    { kind: 'system', path: '/', total_kb: 15 * 1024 * 1024, free_kb: 855 * 1024, low: true },
    { kind: 'data', path: '/data', total_kb: 900 * 1024 * 1024, free_kb: 400 * 1024 * 1024, low: false }
  ]
}

test('no warning when there is enough space', async () => {
  const { wrapper } = mountMenu(true, enough)
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').exists()).toBe(false)
  wrapper.unmount()
})

test('warning names the internal memory and how much is left', async () => {
  const { wrapper } = mountMenu(true, full)
  await flushPromises()
  const warning = wrapper.find('[data-testid="disk-space-warning"]')
  expect(warning.exists()).toBe(true)
  expect(warning.text()).toContain('855 MB')
  expect(warning.text()).toContain('internal memory')
  expect(warning.text()).not.toContain('/data')
  expect(warning.attributes('href')).toBe('/internalmemory')
  wrapper.unmount()
})

test('warning names the data disk when that is the full one', async () => {
  const { wrapper } = mountMenu(true, {
    low: true,
    mounts: [
      { kind: 'system', path: '/', total_kb: 15 * 1024 * 1024, free_kb: 8 * 1024 * 1024, low: false },
      { kind: 'data', path: '/data', total_kb: 900 * 1024 * 1024, free_kb: 512 * 1024, low: true }
    ]
  })
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').text()).toContain('data disk')
  wrapper.unmount()
})

test('free space in gigabytes is rounded', async () => {
  const { wrapper } = mountMenu(true, {
    low: true,
    mounts: [{ kind: 'data', path: '/data', total_kb: 900 * 1024 * 1024, free_kb: 1.5 * 1024 * 1024, low: true }]
  })
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').text()).toContain('1.5 GB')
  wrapper.unmount()
})

test('non admin is not shown the warning', async () => {
  const { wrapper } = mountMenu(false, full)
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').exists()).toBe(false)
  wrapper.unmount()
})

test('space is loaded once the session turns out to be admin', async () => {
  const { wrapper, auth } = mountMenu(false, full)
  await flushPromises()
  auth.admin = true
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').exists()).toBe(true)
  wrapper.unmount()
})

test('failing space request does not show a warning', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/storage/space').reply(500)
  const pinia = createPinia()
  setActivePinia(pinia)
  const auth = useAuthStore()
  auth.loggedIn = true
  auth.admin = true
  const wrapper = mount(Menu, {
    global: {
      plugins: [pinia],
      stubs: { 'router-link': { props: ['to'], template: '<a :href="to"><slot/></a>' } }
    }
  })
  await flushPromises()
  expect(wrapper.find('[data-testid="disk-space-warning"]').exists()).toBe(false)
  wrapper.unmount()
})
