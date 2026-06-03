<script>
import notify from '../util/notify'
import i18n from '../i18n'

function t (key) {
  return i18n.global.t(key)
}

function error (error) {
  let message = t('common.serverError')
  if (error.response) {
    const status = error.response.status
    if (status === 401) {
      this.$router.push('/login')
    } else if (status === 0) {
      console.log('user navigated away from the page')
    } else {
      if (error.response.data && error.response.data.message) {
        message = error.response.data.message
      }
    }
  }
  notify({
    title: t('common.error'),
    message: message,
    type: 'error'
  })
}

function info (message) {
  notify({
    title: t('common.info'),
    message: message,
    type: 'info'
  })
}

function success (message) {
  notify({
    title: t('common.success'),
    message: message,
    type: 'success'
  })
}

export default { error, info, success }
</script>
