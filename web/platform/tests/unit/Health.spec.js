import { mount } from '@vue/test-utils'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import flushPromises from 'flush-promises'
import Health from '../../src/views/Health.vue'

jest.setTimeout(30000)

function metrics (swapFreeKb, swapInPages, swapOutPages) {
  return {
    cpu: { user: 1000, nice: 0, system: 100, idle: 9000, iowait: 0, irq: 0, softirq: 0, steal: 0 },
    memory: {
      total_kb: 8052952,
      available_kb: 1605392,
      free_kb: 265000,
      buffers_kb: 240000,
      cached_kb: 2054000,
      swap_total_kb: 3073016,
      swap_free_kb: swapFreeKb,
      swap_in_pages: swapInPages,
      swap_out_pages: swapOutPages
    },
    disks: [],
    mounts: [],
    net: []
  }
}

async function mountHealth (samples) {
  const mock = new MockAdapter(axios)
  let index = 0
  mock.onGet('/rest/health/metrics').reply(() => {
    const sample = samples[Math.min(index, samples.length - 1)]
    index++
    return [200, { success: true, data: sample }]
  })
  mock.onGet(/\/rest\/health\/events/).reply(200, { success: true, data: [] })
  const wrapper = mount(Health, { attachTo: document.body })
  await flushPromises()
  for (let i = 0; i < samples.length; i++) {
    wrapper.vm.fetchMetrics()
    await flushPromises()
  }
  return { wrapper, mock }
}

test('full but idle swap is not flagged', async () => {
  const { wrapper, mock } = await mountHealth([
    metrics(0, 1000000, 2000000),
    metrics(0, 1000000, 2000000)
  ])
  expect(Math.round(wrapper.vm.swapPct)).toBe(100)
  expect(wrapper.vm.swapRate).toEqual({ inKBs: 0, outKBs: 0 })
  expect(wrapper.vm.swapStatus).toBe('success')
  wrapper.unmount()
  mock.restore()
})

test('sustained paging is flagged even when swap is mostly free', async () => {
  const { wrapper, mock } = await mountHealth([
    metrics(3000000, 1000000, 2000000),
    metrics(3000000, 1002000, 2002000)
  ])
  expect(wrapper.vm.swapPct).toBeLessThan(10)
  expect(wrapper.vm.swapStatus).toBe('exception')
  wrapper.unmount()
  mock.restore()
})

test('light paging warns', async () => {
  const { wrapper, mock } = await mountHealth([
    metrics(1500000, 1000000, 2000000),
    metrics(1500000, 1000200, 2000000)
  ])
  expect(wrapper.vm.swapStatus).toBe('warning')
  wrapper.unmount()
  mock.restore()
})

test('memory and disk bars still key off fill', async () => {
  const { wrapper, mock } = await mountHealth([metrics(0, 0, 0)])
  expect(wrapper.vm.pctStatus(95)).toBe('exception')
  expect(wrapper.vm.pctStatus(80)).toBe('warning')
  expect(wrapper.vm.pctStatus(10)).toBe('success')
  wrapper.unmount()
  mock.restore()
})
