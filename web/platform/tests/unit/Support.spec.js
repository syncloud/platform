import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Support from '../../src/views/Support.vue'

jest.setTimeout(30000)

function mountSupport (showError) {
  return mount(Support,
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
        }
      }
    }
  )
}

test('Send logs to the owner', async () => {
  let sentToSupport
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onPost('/rest/logs/send').reply(function (config) {
    sentToSupport = config.params.include_support
    return [200, { success: true }]
  })

  const wrapper = mountSupport(showError)

  await flushPromises()

  expect(wrapper.find('[data-testid="support-hint"]').text()).toBe('Logs are emailed only to you, at the address of your syncloud.it account.')
  expect(wrapper.find('[data-testid="support-sent"]').exists()).toBe(false)

  await wrapper.find('[data-testid="support-send"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(sentToSupport).toBe(false)
  expect(wrapper.find('[data-testid="support-sent"]').text()).toContain('Logs sent to your email.')
  wrapper.unmount()
})

test('Send logs to support', async () => {
  let sentToSupport
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onPost('/rest/logs/send').reply(function (config) {
    sentToSupport = config.params.include_support
    return [200, { success: true }]
  })

  const wrapper = mountSupport(showError)

  await flushPromises()

  await wrapper.find('#switch').setValue(true)

  expect(wrapper.find('[data-testid="support-hint"]').text()).toBe('Logs are emailed to you, at the address of your syncloud.it account, and a copy goes to Syncloud support.')

  await wrapper.find('[data-testid="support-send"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(sentToSupport).toBe(true)
  expect(wrapper.find('[data-testid="support-sent"]').text()).toContain('Logs sent to your email and to Syncloud support.')
  wrapper.unmount()
})

test('Send logs failure shows no confirmation', async () => {
  const showError = jest.fn()

  const mock = new MockAdapter(axios)

  mock.onPost('/rest/logs/send').reply(500, { success: false, message: 'not sent' })

  const wrapper = mountSupport(showError)

  await flushPromises()

  await wrapper.find('[data-testid="support-send"]').trigger('click')
  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(1)
  expect(wrapper.find('[data-testid="support-sent"]').exists()).toBe(false)
  wrapper.unmount()
})
