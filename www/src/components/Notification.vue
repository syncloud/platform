<script>
import { ElNotification } from 'element-plus'

function error (error) {
  let message = 'Server Error'
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
  ElNotification({
    title: 'Error',
    message: message,
    type: 'error'
  })
}

function info (message) {
  ElNotification({
    title: 'Info',
    message: message,
    type: 'info'
  })
}

function success (message) {
  ElNotification({
    title: 'Success',
    message: message,
    type: 'success'
  })
}

export default { error, info, success }
</script>
