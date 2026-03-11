import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import Login from '../../src/views/Login.vue'

jest.setTimeout(30000)

test('Login redirects to OIDC', async () => {
  delete window.location
  window.location = { href: '' }

  const wrapper = mount(Login,
    {
      attachTo: document.body,
      props: {
        checkUserSession: jest.fn(),
        activated: true
      },
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: jest.fn()
            }
          }
        }
      }
    }
  )

  await flushPromises()

  expect(window.location.href).toBe('/rest/oidc/login')

  wrapper.unmount()
})
