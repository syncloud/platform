const state = {
  loggedIn: true,
  credentials: {
    username: '11',
    password: '2'
  },
  jobStatusRunning: false,
  installerIsRunning: false,
  availableAppsSuccess: true,
  activated: false,
  accessSuccess: true,
  diskActionSuccess: true
}

const store = {
  apps: [
    {
      id: 'wordpress',
      name: 'WordPress',
      icon: '/images/wordpress-128.png'
    },
    {
      id: 'diaspora',
      name: 'Diaspora',
      icon: '/images/penguin.png'
    },
    {
      id: 'mail',
      name: 'Mail',
      icon: '/images/penguin.png'
    },
    {
      id: 'talk',
      name: 'Talk',
      icon: '/images/penguin.png'
    },
    {
      id: 'files',
      name: 'Files Browser',
      icon: '/images/penguin.png'
    }
  ]
}

const installedApps = new Set(['wordpress'])

const appCenterDataError = {
  message: 'error',
  success: false
}

const deviceUrl = {
  device_url: 'http://test.syncloud.it',
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
  data: {
    interfaces: [
      {
        ipv4: [
          {
            addr: '172.17.0.2',
            broadcast: '172.17.0.2',
            netmask: '255.255.0.0'
          },
          {
            addr: '172.17.0.3',
            broadcast: '172.17.0.2',
            netmask: '255.255.0.0'
          }
        ],
        name: 'eth0'
      },
      {
        ipv4: [
          {
            addr: '172.17.0.2',
            broadcast: '172.17.0.2',
            netmask: '255.255.0.0'
          },
          {
            addr: '172.17.0.3',
            broadcast: '172.17.0.2',
            netmask: '255.255.0.0'
          }
        ],
        ipv6: [
          {
            addr: 'fe80::42:acff:fe11:2%eth0',
            netmask: 'ffff:ffff:ffff:ffff::'
          },
          {
            addr: 'fe80::42:acff:fe11:11',
            netmask: 'ffff:ffff:ffff:ffff::'
          }
        ],
        name: 'wifi0'
      }
    ]
  },
  success: true
}

const portMappingsData = {
  port_mappings: [
    {
      local_port: 80,
      external_port: 80
    },
    {
      local_port: 443,
      external_port: 10001
    }

  ],
  success: true
}

const accessData = {
  data: {
    external_access: false,
    upnp_available: false,
    upnp_enabled: true,
    upnp_message: 'not used'
    // public_ip: '111.111.111.111'
  },
  success: true
}

const disksData = {
  disks: [
    {
      name: 'My Passport 0837',
      device: '/dev/sdb',
      active: true,
      size: '931.5G',
      partitions: [
        {
          active: true,
          device: '/dev/sdb1',
          fs_type: 'ntfs',
          mount_point: '/opt/disk/external',
          mountable: true,
          size: '931.5G'
        }
      ]
    },
    {
      name: 'My Passport 0990',
      device: '/dev/sdc',
      active: false,
      size: '931.5G',
      partitions: [
        {
          active: false,
          device: '/dev/sdc1',
          fs_type: 'ntfs',
          mount_point: '',
          mountable: true,
          size: '931.5G'
        }
      ]
    },
    {
      name: 'Blank Disk',
      device: '/dev/sdb',
      size: '100 TB',
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

const versionsData = {
  data: [
    {
      app: {
        id: 'platform',
        name: 'Platform',
        required: true,
        ui: false,
        url: 'http://platform.odroid-c2.syncloud.it'
      },
      current_version: '880',
      installed_version: '876'
    },
    {
      app: {
        id: 'installer',
        name: 'Installer',
        required: true,
        ui: false,
        url: 'http://installer.odroid-c2.syncloud.it'
      },
      current_version: '78',
      installed_version: '75'
    }
  ],
  success: true
}

function appCenterToInstalledApp (app) {
  app.url = 'http://' + app.id + '.odroid-c2.syncloud.it'
  return app
}

function appToInfo (app, installed) {
  const info = {
    app: app,
    current_version: '2',
    installed_version: installed ? '1' : null
  }

  info.app.required = true
  info.app.ui = false
  info.app.url = 'http://' + app.id + '.odroid-c2.syncloud.it'

  return info
}

const express = require('express')
const bodyparser = require('body-parser')
const mock = function (app, server, compiler) {
  app.use(express.urlencoded())
  app.use(bodyparser.json())
  app.post('/rest/login', function (req, res) {
    if (state.credentials.username === req.body.username && state.credentials.password === req.body.password) {
      state.loggedIn = true
      res.json({ message: 'OK' })
    } else {
      res.status(400).json({ message: 'Authentication failed' })
    }
  })
  app.get('/rest/user', function (req, res) {
    if (!state.activated) {
      res.status(501).json({ message: 'Not activated' })
    } else {
      if (state.loggedIn) {
        res.json({ message: 'OK' })
      } else {
        res.status(401).json({ message: 'Authentication failed' })
      }
    }
  })
  app.post('/rest/logout', function (req, res) {
    state.loggedIn = false
    res.json({ message: 'OK' })
  })
  app.get('/rest/activation_status', function (req, res) {
    res.json({ activated: state.activated })
    // res.status(500).json({ message: "unknown activation status" })
  })
  app.get('/rest/installed_apps', function (req, res) {
    if (state.activated) {
      const apps = store.apps.filter(app => installedApps.has(app.id)).map(appCenterToInstalledApp)
      res.json({ apps: apps })
    } else {
      res.status(501).json({ message: 'Not activated' })
    }
  })
  app.get('/rest/app', function (req, res) {
    const app = store.apps.find(app => app.id === req.query.app_id)
    res.json({ info: appToInfo(app, installedApps.has(app.id)) })
  })
  app.post('/rest/upgrade', function (req, res) {
    res.json({ success: true })
  })
  app.post('/rest/install', function (req, res) {
    installedApps.add(req.body.app_id)
    res.json({ success: true })
  })
  app.post('/rest/remove', function (req, res) {
    installedApps.delete(req.body.app_id)
    res.json({ success: true })
  })
  app.post('/rest/restart', function (req, res) {
    res.json({ success: true })
  })
  app.post('/rest/shutdown', function (req, res) {
    res.json({ success: true })
  })
  app.get('/rest/settings/installer_status', function (req, res) {
    res.json({ success: true, is_running: state.installerIsRunning })
    state.installerIsRunning = !state.installerIsRunning
  })
  app.post('/rest/backup/create', function (req, res) {
    res.json({})
  })
  app.get('/rest/job/status', function (req, res) {
    res.json({ success: true, data: state.jobStatusRunning ? 'JobStatusBusy' : 'JobStatusIdle' })
    state.jobStatusRunning = !state.jobStatusRunning
  })

  app.get('/rest/available_apps', function (req, res) {
    let response = {}
    if (state.availableAppsSuccess) {
      response = store
    } else {
      response = appCenterDataError
    }
    res.json(response)
  })

  app.get('/rest/settings/device_url', function (req, res) {
    // res.status(500).json(deviceUrl)
    res.json(deviceUrl)
  })

  app.post('/rest/settings/deactivate', function (req, res) {
    state.activated = false
    res.json({})
  })

  app.get('/rest/backup/list', function (req, res) {
    res.json({
      success: true,
      data: backups
    })
  })

  app.post('/rest/backup/remove', function (req, res) {
    backups = backups.filter(v => v.file !== req.query.file)
    res.json({})
  })

  app.post('/rest/backup/restore', function (req, res) {
    res.json({})
  })

  app.post('/rest/backup/create', function (req, res) {
    res.json({})
  })

  app.post('/rest/installer/upgrade', function (req, res) {
    res.json({ success: true })
  })

  app.get('/rest/access/network_interfaces', function (req, res) {
    res.json(networkInterfaces)
  })

  app.get('/rest/access/access', function (req, res) {
    res.json(accessData)
  })

  app.get('/rest/access/port_mappings', function (req, res) {
    res.json(portMappingsData)
  })

  app.post('/rest/access/set_access', function (req, res) {
    if (state.accessSuccess) {
      accessData.data.external_access = req.body.external_access
      accessData.data.upnp_enabled = req.body.upnp_enabled
      if (req.body.public_ip === undefined) {
        delete accessData.data.public_ip
      } else {
        accessData.data.public_ip = req.body.public_ip
      }
      if (req.query.upnp_enabled) {
        portMappingsData.port_mappings[0].external_port = 81
        portMappingsData.port_mappings[1].external_port = 444
      } else {
        portMappingsData.port_mappings[0].external_port = req.body.certificate_port
        portMappingsData.port_mappings[1].external_port = req.body.access_port
      }
      res.json({ success: true })
    } else {
      res.status(500).json({ success: false, message: 'error' })
    }
    state.accessSuccess = !state.accessSuccess
  })
  app.get('/rest/settings/disks', function (req, res) {
    res.json(disksData)
  })
  app.get('/rest/settings/boot_disk', function (req, res) {
    res.json(bootDiskData)
  })
  app.post('/rest/settings/disk_activate', function (req, res) {
    if (state.diskActionSuccess) {
      res.json(disksData)
    } else {
      res.json(disksDataError)
    }
  })
  app.post('/rest/settings/disk_deactivate', function (req, res) {
    if (state.diskActionSuccess) {
      res.json(disksData)
    } else {
      res.json(disksDataError)
    }
  })
  app.post('/rest/storage/boot_extend', function (req, res) {
    res.json({ success: true })
  })

  app.post('/rest/storage/disk_format', function (req, res) {
    res.json({ success: true })
  })

  app.get('/rest/settings/boot_extend_status', function (req, res) {
    res.json({ success: true, is_running: false })
  })

  app.get('/rest/settings/disk_format_status', function (req, res) {
    res.json({ success: true, is_running: false })
  })

  app.get('/rest/settings/versions', function (req, res) {
    res.json(versionsData)
  })
  app.post('/rest/send_log', function (req, res) {
    res.json({ success: true })
  })
  app.post('/rest/activate/free', function (req, res) {
    state.activated = true
    res.json({ success: true })
    // res.status(500).json({
    //   success: false,
    //   parameters_messages: [
    //     { parameter: 'device_username', messages: ['login is empty'] },
    //     { parameter: 'device_password', messages: ['is too short', 'has no special symbol'] }
    //   ]
    // })
  })
  app.post('/rest/activate/custom', function (req, res) {
    state.activated = true
    res.json({ success: true })
  })
  app.post('/rest/redirect/domain/availability', function (req, res) {
    if (req.body.domain === '1') {
      res.status(400).json({
        success: false,
        parameters_messages: [
          { parameter: 'redirect_password', messages: ['wrong password'] },
          { parameter: 'domain', messages: ['domain is already taken'] }
        ]
      })
    } else {
      res.json({
        success: true
      })
    }
  })
}

exports.mock = mock
