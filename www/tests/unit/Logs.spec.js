import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Logs from '../../src/views/Logs.vue'

jest.setTimeout(30000)

test('Logs', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)
  mock.onGet('/rest/logs').reply(200,
    {
      data: [
        'log 1',
        'log 2'
      ],
      success: true
    }
  )
  const wrapper = mount(Logs,
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
