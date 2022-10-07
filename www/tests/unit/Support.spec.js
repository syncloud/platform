import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Support from '../../src/views/Support.vue'

jest.setTimeout(30000)

test('Send logs to the owner', async () => {
  let sentToSupport
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onPost('/rest/send_log').reply(function (config) {
    sentToSupport = config.params.include_support
    return [200, { success: true }]
  })

  const wrapper = mount(Support,
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
          Switch: {
            template: '<button id="switch" />'
          }
        }
      }
    }
  )

  await flushPromises()

  // await wrapper.find('#switch').trigger('toggle')
  await wrapper.find('#send').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(sentToSupport).toBe(false)
  wrapper.unmount()
})

test('Send logs to support', async () => {
  let sentToSupport
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onPost('/rest/send_log').reply(function (config) {
    sentToSupport = config.params.include_support
    return [200, { success: true }]
  })

  const wrapper = mount(Support,
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
          Switch: {
            template: '<button id="switch" />'
          }
        }
      }
    }
  )

  await flushPromises()

  await wrapper.find('#switch').trigger('toggle')
  await wrapper.find('#send').trigger('click')

  expect(showError).toHaveBeenCalledTimes(0)
  expect(sentToSupport).toBe(true)
  wrapper.unmount()
})
