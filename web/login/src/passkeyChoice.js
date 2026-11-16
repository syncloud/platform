const key = 'passkey'

export function readPasskeyChoice () {
  try {
    return window.localStorage.getItem(key) || ''
  } catch {
    return ''
  }
}

export function writePasskeyChoice (choice) {
  try {
    if (choice) {
      window.localStorage.setItem(key, choice)
    } else {
      window.localStorage.removeItem(key)
    }
  } catch (err) {
    console.log('login: cannot store passkey choice: ' + err)
  }
}
