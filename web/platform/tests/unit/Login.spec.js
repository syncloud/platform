import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import Login from '../../src/views/Login.vue'

jest.setTimeout(30000)

function mountLogin (query) {
  return mount(Login,
    {
      attachTo: document.body,
      global: {
        mocks: {
          $route: { query: query || {} },
          $t: key => key
        },
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
}

test('Login redirects to OIDC', async () => {
  delete window.location
  window.location = { href: '' }

  const wrapper = mountLogin()

  await flushPromises()

  expect(window.location.href).toBe('/rest/oidc/login')

  wrapper.unmount()
})

test('Login explains a failed sign in instead of looping', async () => {
  delete window.location
  window.location = { href: '' }

  const wrapper = mountLogin({ error: 'session' })

  await flushPromises()

  expect(window.location.href).toBe('')
  expect(wrapper.find('[data-testid="login-reason"]').text()).toBe('login.session')

  await wrapper.find('#btn_login_retry').trigger('click')
  expect(window.location.href).toBe('/rest/oidc/login')

  wrapper.unmount()
})
