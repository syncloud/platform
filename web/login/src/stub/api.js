import { createServer, Response } from 'miragejs'

const state = {
  authLevel: 0
}

export function mock () {
  createServer({
    routes () {
      this.urlPrefix = ''

      this.get('/api/state', () => {
        return new Response(200, {}, { data: { authentication_level: state.authLevel } })
      })

      this.post('/api/firstfactor', (_schema, request) => {
        const body = JSON.parse(request.requestBody || '{}')
        if (body.username && body.password) {
          return new Response(200, {}, { data: 'OK' })
        }
        return new Response(401, {}, { data: { message: 'Incorrect username or password' } })
      })

      this.post('/login/totp/status', () => {
        return new Response(200, {}, { data: { configured: false } })
      })

      this.post('/login/totp/setup', () => {
        return new Response(200, {}, { data: { uri: 'otpauth://totp/Syncloud:11?secret=STUBSTUBSTUBSTUB&issuer=Syncloud' } })
      })

      this.post('/api/secondfactor/totp', () => {
        return new Response(200, {}, { data: 'OK' })
      })

      this.post('/api/logout', () => {
        return new Response(200, {}, { data: 'OK' })
      })

      this.passthrough()
    }
  })
}
