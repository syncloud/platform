import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import CertificateLog from '../../src/views/CertificateLog.vue'

jest.setTimeout(30000)

test('Certificate logs', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/certificate/log').reply(200,
    {
      data: [
        'log 1',
        'log 2'
      ],
      success: true
    }
  )
  const wrapper = mount(CertificateLog,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          },
          Dialog: true
        }
      }
    }
  )

  await flushPromises()

  await expect(wrapper.find('#logs').text()).toBe('log 1log 2')

  expect(showError).toHaveBeenCalledTimes(0)
  wrapper.unmount()
})
