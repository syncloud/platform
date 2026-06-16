import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Settings from '../../src/views/Settings.vue'
import { useAuthStore } from '../../src/stores/auth'

function mountSettings (admin) {
  const pinia = createPinia()
  setActivePinia(pinia)
  useAuthStore().admin = admin
  return mount(Settings, {
    global: {
      plugins: [pinia],
      stubs: {
        'router-link': { template: '<a><slot/></a>' }
      }
    }
  })
}

const adminTiles = ['users', 'access', 'internalmemory', 'storage', 'updates', 'backup', 'certificate', 'health', 'customproxy', 'system']
const userTiles = ['activation', 'network', 'support', 'twofactor', 'logs', 'locale']

test('regular user sees locale and two-factor tiles', () => {
  const wrapper = mountSettings(false)
  expect(wrapper.find('#locale').exists()).toBe(true)
  expect(wrapper.find('#twofactor').exists()).toBe(true)
})

test('regular user does not see admin-only tiles', () => {
  const wrapper = mountSettings(false)
  for (const tile of adminTiles) {
    expect(wrapper.find('#' + tile).exists()).toBe(false)
  }
})

test('regular user still sees all user-facing tiles', () => {
  const wrapper = mountSettings(false)
  for (const tile of userTiles) {
    expect(wrapper.find('#' + tile).exists()).toBe(true)
  }
})

test('admin sees every tile', () => {
  const wrapper = mountSettings(true)
  for (const tile of [...adminTiles, ...userTiles]) {
    expect(wrapper.find('#' + tile).exists()).toBe(true)
  }
})
