import { mount } from '@vue/test-utils'
import Backup from '../../src/views/Backup.vue'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'

test('show list of backups', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/backup/list').reply(200,
    {
      success: true,
      data: [
        { path: '/data/platform/backup', file: 'files-2019-0515-123506.tar.gz' },
        { path: '/data/platform/backup', file: 'nextcloud-2019-0515-123506.tar.gz' }
      ]
    }
  )

  const wrapper = mount(Backup,
    {
      attachTo: document.body,
      global: {
        stubs: {
          Error: true
        }
      }
    }
  )

  await flushPromises()

  expect(wrapper.text()).toContain('files-2019-0515-123506.tar.gz')

  // await wrapper.findAll('i')[0].trigger('click')

  wrapper.unmount()
})
