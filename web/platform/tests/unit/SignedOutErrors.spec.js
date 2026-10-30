import { shallowMount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Login from '../../src/views/Login.vue'
import Activation from '../../src/views/Activation.vue'
import TwoFactor from '../../src/views/TwoFactor.vue'

jest.setTimeout(30000)

test.each([
  ['Login', Login],
  ['Activation', Activation],
  ['TwoFactor', TwoFactor]
])('%s does not offer to send logs', async (name, view) => {
  delete window.location
  window.location = { href: '' }
  const mock = new MockAdapter(axios)
  mock.onAny().reply(200, { success: true, data: {} })
  const wrapper = shallowMount(view, {
    global: {
      mocks: { $route: { query: { error: 'session' } }, $router: { push: jest.fn() }, $t: key => key }
    }
  })
  await flushPromises()
  expect(wrapper.findComponent({ name: 'Error' }).props('enableLogs')).toBe(false)
  wrapper.unmount()
})
