import { mount } from '@vue/test-utils'
import UserEdit from '../../src/views/UserEdit.vue'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'

function mountEdit (query = {}, push = jest.fn()) {
  return mount(UserEdit, {
    global: {
      mocks: {
        $route: { query },
        $router: { push }
      }
    }
  })
}

function mockGroups (mock, groups) {
  mock.onGet('/rest/groups').reply(200, { success: true, data: groups })
}

test('create: title is Add user and username editable', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [])
  const wrapper = mountEdit({})
  await flushPromises()

  expect(wrapper.find('[data-testid="user-edit-title"]').text()).toContain('Add user')
  expect(wrapper.find('#user_username').attributes().disabled).toBeUndefined()
})

test('create: save disabled until username and a strong password', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [])
  const wrapper = mountEdit({})
  await flushPromises()

  const save = wrapper.find('#btn_save')
  expect(save.attributes().disabled).toBeDefined()

  await wrapper.find('#user_username').setValue('bob')
  expect(save.attributes().disabled).toBeDefined()

  await wrapper.find('#user_password').setValue('short')
  expect(wrapper.find('[data-testid="pwrule-length"]').classes()).not.toContain('pw-ok')
  expect(wrapper.find('[data-testid="pwrule-number"]').classes()).not.toContain('pw-ok')
  expect(wrapper.find('[data-testid="pwrule-letter"]').classes()).toContain('pw-ok')
  expect(save.attributes().disabled).toBeDefined()

  await wrapper.find('#user_password').setValue('strongpass1')
  expect(wrapper.find('[data-testid="pwrule-length"]').classes()).toContain('pw-ok')
  expect(wrapper.find('[data-testid="pwrule-number"]').classes()).toContain('pw-ok')
  expect(save.attributes().disabled).toBeUndefined()
})

test('create: syncloud group is hidden from chips', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [{ name: 'syncloud', members: [] }, { name: 'family', members: [] }])
  const wrapper = mountEdit({})
  await flushPromises()

  expect(wrapper.find('[data-testid="user-group-syncloud"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="user-group-family"]').exists()).toBe(true)
})

test('create: saves user then assigns selected groups and navigates', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [{ name: 'family', members: [] }])
  let added = null
  let member = null
  mock.onPost('/rest/users/add').reply(c => { added = JSON.parse(c.data); return [200, { success: true }] })
  mock.onPost('/rest/groups/member').reply(c => { member = JSON.parse(c.data); return [200, { success: true }] })

  const push = jest.fn()
  const wrapper = mountEdit({}, push)
  await flushPromises()

  await wrapper.find('#user_username').setValue('bob')
  await wrapper.find('#user_password').setValue('strongpass1')
  await wrapper.find('[data-testid="user-group-family"]').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()

  expect(added).toEqual({ username: 'bob', password: 'strongpass1', email: '', admin: false })
  expect(member).toEqual({ group: 'family', username: 'bob', member: true })
  expect(push).toHaveBeenCalledWith('/users')
})

test('create: inline new group is created and auto-selected', async () => {
  const mock = new MockAdapter(axios)
  let call = 0
  mock.onGet('/rest/groups').reply(() => {
    call++
    const data = call === 1 ? [] : [{ name: 'team', members: [] }]
    return [200, { success: true, data }]
  })
  let addReq = null
  mock.onPost('/rest/groups/add').reply(c => { addReq = JSON.parse(c.data); return [200, { success: true }] })

  const wrapper = mountEdit({})
  await flushPromises()

  await wrapper.find('#new_group').setValue('team')
  await wrapper.find('[data-testid="group-create"]').trigger('click')
  await flushPromises()

  expect(addReq).toEqual({ name: 'team' })
  expect(wrapper.find('[data-testid="user-group-team"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="user-group-team"]').classes()).toContain('is-on')
})

test('edit: loads user, username disabled, sole admin locked', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [{ name: 'syncloud', members: ['admin'] }, { name: 'family', members: [] }])
  mock.onGet('/rest/users').reply(200, {
    success: true,
    data: [{ username: 'admin', email: 'admin@example.com', admin: true, groups: [] }]
  })

  const wrapper = mountEdit({ username: 'admin' })
  await flushPromises()

  expect(wrapper.find('[data-testid="user-edit-title"]').text()).toContain('Edit user')
  expect(wrapper.find('#user_username').element.value).toBe('admin')
  expect(wrapper.find('#user_username').attributes().disabled).toBeDefined()
  expect(wrapper.find('#user_email').element.value).toBe('admin@example.com')
  expect(wrapper.find('[data-testid="user-admin-last"]').exists()).toBe(true)
  expect(wrapper.find('#user_admin').attributes().disabled).toBeDefined()
})

test('edit: admin switch editable when more than one admin', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [{ name: 'syncloud', members: ['admin', 'bob'] }])
  mock.onGet('/rest/users').reply(200, {
    success: true,
    data: [
      { username: 'admin', email: 'a@example.com', admin: true, groups: [] },
      { username: 'bob', email: 'b@example.com', admin: true, groups: [] }
    ]
  })

  const wrapper = mountEdit({ username: 'bob' })
  await flushPromises()

  expect(wrapper.find('[data-testid="user-admin-last"]').exists()).toBe(false)
  expect(wrapper.find('#user_admin').attributes().disabled).toBeUndefined()
})

test('edit: applies email, password, admin and group diffs', async () => {
  const mock = new MockAdapter(axios)
  mockGroups(mock, [
    { name: 'syncloud', members: ['admin'] },
    { name: 'family', members: [] },
    { name: 'work', members: ['alice'] }
  ])
  mock.onGet('/rest/users').reply(200, {
    success: true,
    data: [
      { username: 'admin', email: 'a@example.com', admin: true, groups: [] },
      { username: 'alice', email: 'alice@example.com', admin: false, groups: ['work'] }
    ]
  })
  let emailBody = null
  let passwordBody = null
  let adminBody = null
  const members = []
  mock.onPost('/rest/users/email').reply(c => { emailBody = JSON.parse(c.data); return [200, { success: true }] })
  mock.onPost('/rest/users/password').reply(c => { passwordBody = JSON.parse(c.data); return [200, { success: true }] })
  mock.onPost('/rest/users/admin').reply(c => { adminBody = JSON.parse(c.data); return [200, { success: true }] })
  mock.onPost('/rest/groups/member').reply(c => { members.push(JSON.parse(c.data)); return [200, { success: true }] })

  const push = jest.fn()
  const wrapper = mountEdit({ username: 'alice' }, push)
  await flushPromises()

  await wrapper.find('#user_email').setValue('alice@new.org')
  await wrapper.find('#user_password').setValue('newpass123')
  await wrapper.find('#user_admin').setValue(true)
  await wrapper.find('[data-testid="user-group-family"]').trigger('click')
  await wrapper.find('[data-testid="user-group-work"]').trigger('click')
  await wrapper.find('#btn_save').trigger('click')
  await flushPromises()

  expect(emailBody).toEqual({ username: 'alice', email: 'alice@new.org' })
  expect(passwordBody).toEqual({ username: 'alice', password: 'newpass123' })
  expect(adminBody).toEqual({ username: 'alice', admin: true })
  expect(members).toContainEqual({ group: 'family', username: 'alice', member: true })
  expect(members).toContainEqual({ group: 'work', username: 'alice', member: false })
  expect(push).toHaveBeenCalledWith('/users')
})

function mockAliceEdit (mock) {
  mockGroups(mock, [{ name: 'syncloud', members: ['admin', 'alice'] }])
  mock.onGet('/rest/users').reply(200, {
    success: true,
    data: [
      { username: 'admin', email: 'a@example.com', admin: true, groups: [] },
      { username: 'alice', email: 'alice@example.com', admin: false, groups: [] }
    ]
  })
}

test('edit: delete asks for confirmation before removing', async () => {
  const mock = new MockAdapter(axios)
  mockAliceEdit(mock)
  let removeCalled = false
  mock.onPost('/rest/users/remove').reply(() => { removeCalled = true; return [200, { success: true }] })

  const push = jest.fn()
  const wrapper = mountEdit({ username: 'alice' }, push)
  await flushPromises()

  expect(wrapper.find('[data-testid="btn_confirm"]').exists()).toBe(false)
  await wrapper.find('#btn_delete').trigger('click')
  expect(wrapper.find('[data-testid="btn_confirm"]').exists()).toBe(true)
  expect(removeCalled).toBe(false)
})

test('edit: confirming delete posts remove and navigates', async () => {
  const mock = new MockAdapter(axios)
  mockAliceEdit(mock)
  let removeBody = null
  mock.onPost('/rest/users/remove').reply(c => { removeBody = JSON.parse(c.data); return [200, { success: true }] })

  const push = jest.fn()
  const wrapper = mountEdit({ username: 'alice' }, push)
  await flushPromises()

  await wrapper.find('#btn_delete').trigger('click')
  await wrapper.find('[data-testid="btn_confirm"]').trigger('click')
  await flushPromises()

  expect(removeBody).toEqual({ username: 'alice' })
  expect(push).toHaveBeenCalledWith('/users')
})
