import { mount, RouterLinkStub } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Certificate from '@/views/Certificate'

jest.setTimeout(30000)

test('Certificate', async () => {
  const showError = jest.fn()
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/certificate').reply(200,
    {
      data: {
        is_valid: true,
        is_real: false,
        valid_for_days: 10
      },
      success: true
    }
  )
  const wrapper = mount(Certificate,
    {
      attachTo: document.body,
      global: {
        components: {
          RouterLink: RouterLinkStub
        },
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Confirmation: true
        },
        mocks: {
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  await expect(wrapper.find('#valid_days').text()).toBe('10')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})
