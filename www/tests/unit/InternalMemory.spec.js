import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import InternalMemory from '@/views/InternalMemory'

jest.setTimeout(30000)

test('Extend', async () => {
  const showError = jest.fn()
  let extended = false
  const mockRouter = { push: jest.fn() }

  const mock = new MockAdapter(axios)

  mock.onGet('/rest/storage/boot/disk').reply(function (_) {
    return [200, {
      data: {
        device: '/dev/mmcblk0p2',
        size: '2G',
        extendable: true
      },
      success: true
    }]
  })

  mock.onPost('/rest/storage/boot_extend').reply(function (_) {
    extended = true
    return [200, { success: true }]
  })

  const wrapper = mount(InternalMemory,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: {
            template: '<span/>',
            methods: {
              showAxios: showError
            }
          }
        },
        mocks: {
          $route: { path: '/app', query: { id: 'files' } },
          $router: mockRouter
        }
      }
    }
  )

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)

  await wrapper.find('#btn_boot_extend').trigger('click')
  // await wrapper.find('#btn_confirm').trigger('click')

  await flushPromises()

  expect(showError).toHaveBeenCalledTimes(0)
  expect(extended).toBe(true)
  wrapper.unmount()
})
