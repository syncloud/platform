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

const tiles = ['users', 'access', 'internalmemory', 'storage', 'updates', 'backup', 'certificate', 'health', 'customproxy', 'system', 'activation', 'network', 'support', 'logs', 'locale', 'twofactor']

test('non-admin sees no settings tiles', () => {
  const wrapper = mountSettings(false)
  for (const tile of tiles) {
    expect(wrapper.find('#' + tile).exists()).toBe(false)
  }
})

test('admin sees every settings tile', () => {
  const wrapper = mountSettings(true)
  for (const tile of tiles) {
    expect(wrapper.find('#' + tile).exists()).toBe(true)
  }
})
