import axios from 'axios'

export function checkForServiceError (data, onComplete, onError) {
  if ('success' in data && !data.success) {
    const err = {
      response: {
        status: 200,
        data: data
      }
    }
    onError(err)
  } else {
    onComplete()
  }
}

export async function post (url, body) {
  const resp = await axios.post(url, body)
  if (resp.data && 'success' in resp.data && !resp.data.success) {
    const err = new Error('service error')
    err.response = { status: 200, data: resp.data }
    throw err
  }
  return resp.data
}

export const INSTALLER_STATUS_URL = '/rest/installer/status'
export const DEFAULT_STATUS_PREDICATE = (response) => {
  return response.data.data.is_running
}

export const JOB_STATUS_URL = '/rest/job/status'
export const JOB_STATUS_PREDICATE = (response) => {
  return response.data.data.status !== 'Idle'
}

export function runAfterJobIsComplete (timeoutFunc, onComplete, onError, statusUrl, statusPredicate) {
  const recheckFunc = function () {
    runAfterJobIsComplete(timeoutFunc, onComplete, onError, statusUrl, statusPredicate)
  }

  const recheckTimeout = 2000
  axios.get(statusUrl)
    .then(response => {
      if (statusPredicate(response)) {
        timeoutFunc(recheckFunc, recheckTimeout)
      } else {
        onComplete()
      }
    })
    .catch(err => {
      console.log('status err: ' + err)
      // Auth error means job is finished
      if (err.response !== undefined && err.response.status === 401) {
        onComplete()
      } else {
        timeoutFunc(recheckFunc, recheckTimeout)
      }
    })
}
