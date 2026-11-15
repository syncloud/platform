import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import Login from '../../src/views/Login.vue'
import Error from '../../src/components/Error.vue'

jest.setTimeout(30000)

function mountError () {
  return mount(Error, {
    attachTo: document.body,
    global: {
      mocks: {
        $route: { path: '/' },
        $router: { push: jest.fn() },
        $t: key => key,
        $te: () => false
      }
    }
  })
}

function mountLogin () {
  return mount(Login, {
    attachTo: document.body,
    global: {
      mocks: {
        $route: { query: {} },
        $t: key => key
      }
    }
  })
}

const aborted = { message: 'Request aborted', code: 'ECONNABORTED' }

beforeEach(() => {
  delete window.location
  window.location = { href: '' }
  window.dispatchEvent(new Event('pageshow'))
})

test('request cut off by the sign in redirect shows no error', async () => {
  const error = mountError()
  const login = mountLogin()
  await flushPromises()
  expect(window.location.href).toBe('/rest/oidc/login')

  await error.vm.showAxios(aborted)

  expect(error.find('[data-testid="error-dialog"]').exists()).toBe(false)
  login.unmount()
  error.unmount()
})

test('request without a response still shows an error while staying on the page', async () => {
  const error = mountError()

  await error.vm.showAxios(aborted)

  expect(error.find('[data-testid="error-dialog"]').exists()).toBe(true)
  error.unmount()
})
