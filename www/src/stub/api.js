import { createServer, Model, Response } from 'miragejs'
const domain = 'syncloud.it'
const state = {
  loggedIn: true,
  credentials: {
    username: '11',
    password: '2'
  },
  jobStatusRunning: false,
  installerIsRunning: true,
  availableAppsSuccess: true,
  activated: true,
  accessSuccess: true,
  diskActionSuccess: true,
  diskLastError: true
}

const store = {
  data: [
    {
      app: {
        id: 'wordpress',
        name: 'WordPress',
        icon: '/images/wordpress-128.png',
        required: true,
        ui: false,
        url: 'https://wordpress.odroid-c2.' + domain
      },
      current_version: '2',
      installed_version: '1'
    },
    {
      app: {
        id: 'diaspora',
        name: 'Diaspora',
        icon: '/images/penguin.png',
        required: false,
        ui: true,
        url: 'https://diaspora.odroid-c2.' + domain
      },
      current_version: '1',
      installed_version: '2'
    },
    {
      app: {
        id: 'mail',
        name: 'Mail',
        icon: '/images/penguin.png',
        required: false,
        ui: true,
        url: 'https://mail.odroid-c2.' + domain
      },
      current_version: '1',
      installed_version: '2'
    },
    {
      app: {
        id: 'talk',
        name: 'Talk',
        icon: '/images/penguin.png',
        required: false,
        ui: true,
        url: 'https://talk.odroid-c2.' + domain
      },
      current_version: '1',
      installed_version: '2'
    },
    {
      app: {
        id: 'files',
        name: 'Files Browser',
        icon: '/images/penguin.png',
        required: false,
        ui: true,
        url: 'https://files.odroid-c2.' + domain
      },
      current_version: '1',
      installed_version: '2'
    },
    {
      app: {
        id: 'platform',
        name: 'Platform',
        icon: '/images/penguin.png',
        required: true,
        ui: false,
        url: 'https://platform.odroid-c2.' + domain
      },
      current_version: '880',
      installed_version: '876'
    },
    {
      app: {
        id: 'installer',
        name: 'Installer',
        icon: '/images/penguin.png',
        required: true,
        ui: false,
        url: 'https://installer.odroid-c2.' + domain
      },
      current_version: '78',
      installed_version: '75'
    }
  ]
}

const installer = {
  data: {
    installed_version: 1,
    store_version: 2
  }
}
const installedApps = new Set(['wordpress'])

const appCenterDataError = {
  message: 'error',
  success: false
}

const deviceUrl = {
  data: 'https://test.' + domain,
  success: true
}

let backups = [
  { path: '/data/platform/backup', file: 'files-2019-0515-123506.tar.gz' },
  { path: '/data/platform/backup', file: 'nextcloud-2019-0515-123506.tar.gz' },
  { path: '/data/platform/backup', file: 'diaspora-2019-0512-103501.tar.gz' },
  { path: '/data/platform/backup', file: 'nextcloud-2019-0521-113502.tar.gz' },
  { path: '/data/platform/backup', file: 'nextcloud-2019-0201-122500.tar.gz' },
  { path: '/data/platform/backup', file: 'files-2019-0415-123506.tar.gz' }
]

const networkInterfaces = {
  data: [
    {
      name: 'eth0',
      addresses: [
        '172.17.0.2'
      ]
    },
    {
      name: 'wifi0',
      addresses: [
        // '172.17.0.2',
        // '172.17.0.3',
        // 'fe80::42:acff:fe11:2%eth0',
        // 'fe80::42:acff:fe11:11'
      ]
    }
  ],
  success: true
}

const accessData = {
  data: {
    ipv4_enabled: true,
    ipv4_public: true
    // access_port: 443
    // public_ip: '111.111.111.111'
  },
  success: true
}

const disksData = {
  data: [
    {
      name: 'Disk 1',
      device: '/dev/sda',
      active: true,
      size: '100G',
      has_errors: false,
      raid: 'raid10',
      partitions: []
    },
    {
      name: 'Disk 2',
      device: '/dev/sdb',
      active: true,
      size: '100G',
      raid: 'raid10',
      partitions: []
    },
    {
      name: 'Disk 3',
      device: '/dev/sdc',
      active: true,
      size: '100G',
      raid: 'raid10',
      partitions: []
    },
    {
      name: 'Disk 4',
      device: '/dev/sdd',
      active: true,
      size: '100G',
      raid: 'raid10',
      partitions: []
    }
  ],
  success: true
}

const disksDataError = {
  message: 'error',
  success: false
}

const bootDiskData = {
  data: {
    device: '/dev/mmcblk0p2',
    size: '2G',
    extendable: true
  },
  success: true
}

const installerProgress = {
  counter: 0,
  success: true,
  data: {
    progress: {
      wordpress: {
        app: 'wordpress',
        action: 'Install',
        summary: 'Downloading',
        indeterminate: false,
        percentage: 10
      }
    }
  }
}

export function mock () {
  createServer({
    models: {
      author: Model
    },
    routes () {
      this.post('/rest/login', function (_schema, request) {
        const attrs = JSON.parse(request.requestBody)
        if (state.credentials.username === attrs.username && state.credentials.password === attrs.password) {
          state.loggedIn = true
          return new Response(200, {}, { message: 'OK' })
        } else {
          return new Response(400, {}, { message: 'Authentication failed' })
        }
      })
      this.get('/rest/user', function (_schema, _request) {
        if (!state.activated) {
          return new Response(501, {}, { message: 'Not activated' })
        } else {
          if (state.loggedIn) {
            return new Response(200, {}, { message: 'OK' })
          } else {
            return new Response(401, {}, { message: 'Authentication failed' })
          }
        }
      })
      this.post('/rest/logout', function (_schema, _request) {
        state.loggedIn = false
        return new Response(200, {}, { message: 'OK' })
      })
      this.get('/rest/activation/status', function (_schema, _request) {
        console.debug('activated: ' + state.activated)
        return new Response(200, {}, { data: state.activated })
        // return new Response(500, {}, { message: "unknown activation status" })
      })
      this.get('/rest/apps/installed', function (_schema, _request) {
        if (state.activated) {
          const apps = store.data.filter(app => installedApps.has(app.app.id)).map(info => info.app)
          return new Response(200, {}, { data: apps })
        } else {
          return new Response(501, {}, { message: 'Not activated' })
        }
      })
      this.get('/rest/app', function (_schema, request) {
        console.debug(request.queryParams.app_id)
        const info = store.data.find(info => info.app.id === request.queryParams.app_id)
        console.debug(info)
        if (!installedApps.has(info.app.id)) {
          info.installed_version = null
        } else {
          info.installed_version = '1'
        }
        return new Response(200, {}, { data: info })
      })
      this.get('/rest/installer/version', function (_schema, _request) {
        return new Response(200, {}, installer)
      })
      this.post('/rest/app/upgrade', function (_schema, request) {
        installerProgress.counter = 0
        state.installerIsRunning = true
        const attrs = JSON.parse(request.requestBody)
        installerProgress.data.progress[attrs.app_id] = {
          app: attrs.app_id,
          action: 'Upgrade'
        }
        installerProgress.data.progress.app = attrs.app_id
        return new Response(200, {}, { success: true })
      })
      this.post('/rest/app/install', function (_schema, request) {
        installerProgress.counter = 0
        state.installerIsRunning = true
        const attrs = JSON.parse(request.requestBody)
        installedApps.add(attrs.app_id)
        installerProgress.data.progress[attrs.app_id] = {
          app: attrs.app_id,
          action: 'Install'
        }
        return new Response(200, {}, { success: true })
      })
      this.post('/rest/app/remove', function (_schema, request) {
        installerProgress.counter = 0
        state.installerIsRunning = true
        const attrs = JSON.parse(request.requestBody)
        installedApps.delete(attrs.app_id)
        installerProgress.data.progress[attrs.app_id] = {
          app: attrs.app_id,
          action: 'Remove'
        }
        return new Response(200, {}, { success: true })
      })
      this.post('/rest/restart', function (_schema, _request) {
        return new Response(200, {}, { success: true })
      })
      this.post('/rest/shutdown', function (_schema, _request) {
        return new Response(200, {}, { success: true })
      })
      this.get('/rest/installer/status', function (_schema, _request) {
        console.debug('counter: ' + installerProgress.counter)
        installerProgress.counter += 1
        for (const [app, progress] of Object.entries(installerProgress.data.progress)) {
          if (installerProgress.counter <= 5 && progress.action !== 'Remove') {
            progress.indeterminate = false
            progress.summary = 'Downloading'
            progress.percentage = installerProgress.counter * 20
          } else {
            progress.indeterminate = true
            if (progress.action === 'Upgrade') {
              progress.summary = 'Upgrading'
            }
            if (progress.action === 'Remove') {
              progress.summary = 'Removing'
            }
            if (progress.action === 'Install') {
              progress.summary = 'Installing'
            }
          }
          if (installerProgress.counter > 7) {
            installerProgress.counter = 0
            delete installerProgress.data.progress[app]
          }
        }
        return new Response(200, {}, { data: installerProgress.data })
      })
      this.post('/rest/backup/create', function (_schema, _request) {
        return new Response(200, {}, {})
      })
      this.get('/rest/job/status', function (_schema, _request) {
        state.jobStatusRunning = !state.jobStatusRunning
        return new Response(200, {}, { success: true, data: { status: state.jobStatusRunning ? 'Busy' : 'Idle', name: 'storage.activate.disks' } })
      })

      this.get('/rest/apps/available', function (_schema, _request) {
        if (state.availableAppsSuccess) {
          const apps = store.data.map(info => info.app)
          return new Response(200, {}, { data: apps })
        } else {
          return new Response(200, {}, appCenterDataError)
        }
      })

      this.get('/rest/device/url', function (_schema, _request) {
        // return new Response(500, {}, deviceUrl)
        return new Response(200, {}, deviceUrl)
      })

      this.post('/rest/deactivate', function (_schema, _request) {
        state.activated = false
        console.debug('activated: ' + state.activated)
        return new Response(200, {}, {})
      })

      this.get('/rest/backup/list', function (_schema, _request) {
        return new Response(200, {}, {
          success: true,
          data: backups
        })
      })

      this.get('/rest/backup/auto', function (_schema, _request) {
        return new Response(200, {}, {
          success: true,
          data: {
            auto: 'no',
            day: 0,
            hour: 0
          }
        })
      })

      this.post('/rest/backup/auto', function (_schema, request) {
        const attrs = JSON.parse(request.requestBody)
        return new Response(200, {}, {
          success: true,
          data: {
            auto: attrs.auto,
            day: attrs.day,
            hour: attrs.auto
          }
        })
      })

      this.post('/rest/backup/remove', function (_schema, request) {
        const attrs = JSON.parse(request.requestBody)
        backups = backups.filter(v => v.file !== attrs.file)
        return new Response(200, {}, {})
      })

      this.post('/rest/backup/restore', function (_schema, _request) {
        return new Response(200, {}, {})
      })

      this.post('/rest/backup/create', function (_schema, _request) {
        return new Response(200, {}, {})
      })

      this.post('/rest/installer/upgrade', function (_schema, _request) {
        return new Response(200, {}, { success: true })
      })

      this.get('/rest/network/interfaces', function (_schema, _request) {
        return new Response(200, {}, networkInterfaces)
      })

      this.get('/rest/access', function (_schema, _request) {
        return new Response(200, {}, accessData)
      })

      this.post('/rest/access', function (_schema, request) {
        const attrs = JSON.parse(request.requestBody)
        // state.accessSuccess = !state.accessSuccess
        // if (state.accessSuccess) {
        accessData.data = attrs
        // accessData.data.external_access = attrs.external_access
        // if (attrs.public_ip === undefined) {
        //   delete accessData.data.public_ip
        // } else {
        //   accessData.data.public_ip = attrs.public_ip
        // }
        // accessData.access_port = attrs.access_port
        return new Response(200, {}, { success: true })
        // } else {
        //   return new Response(500, {}, { success: false, message: 'error' })
        // }
      })
      this.get('/rest/storage/disks', function (_schema, _request) {
        return new Response(200, {}, disksData)
      })
      this.get('/rest/storage/boot/disk', function (_schema, _request) {
        return new Response(200, {}, bootDiskData)
      })
      this.post('/rest/storage/activate/partition', function (_schema, _request) {
        if (state.diskActionSuccess) {
          return new Response(200, {}, disksData)
        } else {
          return new Response(200, {}, disksDataError)
        }
      })
      this.post('/rest/storage/activate/disk', function (_schema, _request) {
        if (state.diskActionSuccess) {
          return new Response(200, {}, disksData)
        } else {
          return new Response(200, {}, disksDataError)
        }
      })
      this.post('/rest/storage/deactivate', function (_schema, _request) {
        if (state.diskActionSuccess) {
          return new Response(200, {}, disksData)
        } else {
          return new Response(200, {}, disksDataError)
        }
      })
      this.post('/rest/storage/boot_extend', function (_schema, _request) {
        bootDiskData.data.extendable = !bootDiskData.data.extendable
        return new Response(200, {}, { success: true })
      })
      this.get('/rest/storage/error/last', function (_schema, _request) {
        if (state.diskLastError) {
          return new Response(500, {}, { success: false, message: 'Disk format error' })
        } else {
          return new Response(200, {}, { success: true, data: 'OK' })
        }
      })
      this.post('/rest/storage/error/clear', function (_schema, _request) {
        state.diskLastError = false
        return new Response(200, {}, { success: true })
      })

      this.get('/rest/settings/boot_extend_status', function (_schema, _request) {
        return new Response(200, {}, { success: true, is_running: false })
      })

      this.get('/rest/settings/disk_format_status', function (_schema, _request) {
        return new Response(200, {}, { success: true, is_running: false })
      })
      this.post('/rest/logs/send', function (_schema, _request) {
        // return new Response(200, {}, { success: true })
        return new Response(400, {}, { message: 'Cannot send logs' })
      })
      this.post('/rest/activate/managed', function (_schema, _request) {
        state.activated = true
        console.debug('activated: ' + state.activated)
        return new Response(200, {}, { success: true })
        // return new Response(500, {}, {
        //   success: false,
        //   parameters_messages: [
        //     { parameter: 'device_username', messages: ['login is empty'] },
        //     { parameter: 'device_password', messages: ['is too short', 'has no special symbol'] }
        //   ]
        // })
      })
      this.post('/rest/redirect/domain/availability', function (_schema, request) {
        const attrs = JSON.parse(request.requestBody)
        if (attrs.domain === '1') {
          return new Response(400, {}, {
            success: false,
            parameters_messages: [
              { parameter: 'redirect_password', messages: ['wrong password'] },
              { parameter: 'domain', messages: ['domain is already taken'] }
            ]
          })
        } else {
          return new Response(200, {}, {
            success: true
          })
        }
      })
      this.get('/rest/redirect_info', function (_schema, _request) {
        if (state.activated) {
          return new Response(502, {}, { message: 'Device is activated' })
        } else {
          return new Response(200, {}, { success: true, data: { domain: domain } })
        }
      })
      this.get('/rest/certificate', function (_schema, _request) {
        const info = {
          is_valid: true,
          is_real: false,
          valid_for_days: 10
        }
        return new Response(200, {}, { success: true, data: info })
      })
      this.get('/rest/certificate/log', function (_schema, _request) {
        const logs = [
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: ----- {"category": "certificate"}',
          "Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: writing new private key to '/var/snap/platform/current/syncloud.key' {\"category\": \"certificate\"}",
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: .............................+++++ {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: .........................................................+++++ {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: Generating a RSA private key {"category": "certificate"}',
          'Dec 15 08:42:35 syncloud platform.backend[26230]: cert/fake.go:35 generating self signed certificate {"category": "certificate"}',
          'Dec 15 08:42:35 syncloud platform.backend[26230]: cert/generator.go:75 unable to generate certificate: acme: error: 429 :: POST :: https://acme-v02.api.letsencrypt.org/acme/new-acct :: urn:ietf:params:acme:error:rateLimited :: Error creating new account :: too many registrations for this IP: see https://letsencrypt.org/docs/rate-limits/ {"category": "certificate"}'
        ]
        return new Response(200, {}, { success: true, data: logs })
      })
      this.get('/rest/logs', function (_schema, _request) {
        const logs = [
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: ----- {"category": "certificate"}',
          "Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: writing new private key to '/var/snap/platform/current/syncloud.key' {\"category\": \"certificate\"}",
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: .............................+++++ {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: .........................................................+++++ {"category": "certificate"}',
          'Dec 15 08:42:36 syncloud platform.backend[26230]: cert/fake.go:51 output: Generating a RSA private key {"category": "certificate"}',
          'Dec 15 08:42:35 syncloud platform.backend[26230]: cert/fake.go:35 generating self signed certificate {"category": "certificate"}',
          'Dec 15 08:42:35 syncloud platform.backend[26230]: cert/generator.go:75 unable to generate certificate: acme: error: 429 :: POST :: https://acme-v02.api.letsencrypt.org/acme/new-acct :: urn:ietf:params:acme:error:rateLimited :: Error creating new account :: too many registrations for this IP: see https://letsencrypt.org/docs/rate-limits/ {"category": "certificate"}'
        ]
        return new Response(200, {}, { success: true, data: logs })
      })
    }
  })
}
