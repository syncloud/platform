import { mount } from '@vue/test-utils'
import UsersList from '../../src/views/UsersList.vue'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'

function mountList (push = jest.fn()) {
  return mount(UsersList, {
    global: {
      mocks: { $router: { push } },
      stubs: {
        'router-link': { template: '<a :data-testid="$attrs[\'data-testid\']"><slot/></a>' }
      }
    }
  })
}

test('renders users with email, admin badge and groups', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/users').reply(200, {
    success: true,
    data: [
      { username: 'admin', email: 'admin@example.com', admin: true, groups: [] },
      { username: 'alice', email: 'alice@example.com', admin: false, groups: ['family'] }
    ]
  })

  const wrapper = mountList()
  await flushPromises()

  expect(wrapper.find('[data-testid="user-row-admin"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="user-row-alice"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="user-admin-badge-admin"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="user-admin-badge-alice"]').exists()).toBe(false)
  expect(wrapper.text()).toContain('alice@example.com')
  expect(wrapper.text()).toContain('family')
})

test('shows add user button', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/users').reply(200, { success: true, data: [] })

  const wrapper = mountList()
  await flushPromises()

  expect(wrapper.find('[data-testid="users-add"]').exists()).toBe(true)
})

test('add button navigates to create screen', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/users').reply(200, { success: true, data: [] })

  const push = jest.fn()
  const wrapper = mountList(push)
  await flushPromises()

  await wrapper.find('[data-testid="users-add"]').trigger('click')
  expect(push).toHaveBeenCalledWith('/useredit')
})

test('shows empty state when there are no users', async () => {
  const mock = new MockAdapter(axios)
  mock.onGet('/rest/users').reply(200, { success: true, data: [] })

  const wrapper = mountList()
  await flushPromises()

  expect(wrapper.text()).toContain('No users')
})
