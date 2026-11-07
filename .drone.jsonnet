local name = 'platform';
local playwright = 'v1.49.1-jammy';
local go = '1.27.3';
local node = '22.16.0';
local publisher_image = 'syncloud/store-publisher:stable-346';
local authelia = '4.39.15';
local distro_default = 'buster';
local distros = ['bookworm', 'buster'];
local tools = ['gptfdisk', 'btrfs', 'smartmontools'];
local bootstrap = '25.02';
local nginx = '1.24.0';
local caddy = '2.11.7';
local caddy_l4 = '0.1.2-syncloud.1';
local xcaddy = '0.4.7';
local python = '3.12-slim-bookworm';
local alpine = '3.21';
local visual_diff_skip_build = '3144';

local shared = [{
  kind: 'pipeline',
  name: 'shared',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [
    {
      name: 'build web',
      image: 'node:' + node,
      environment: {
        NODE_OPTIONS: '--max_old_space_size=2048',
      },
      commands: [
        './web/build.sh',
      ],
    },
    {
      name: 'build',
      image: 'golang:' + go,
      commands: [
        './shared/build.sh ' + caddy + ' ' + caddy_l4 + ' ' + xcaddy,
      ],
    },
    {
      name: 'upload',
      image: 'appleboy/drone-scp:1.6.4',
      settings: {
        host: {
          from_secret: 'artifact_host',
        },
        username: 'artifact',
        key: {
          from_secret: 'artifact_key',
        },
        timeout: '2m',
        command_timeout: '2m',
        target: '/home/artifact/repo/' + name + '/${DRONE_BUILD_NUMBER}-shared',
        source: 'artifact/shared/*',
        strip_components: 2,
      },
    },
  ],
  trigger: {
    event: [
      'push',
    ],
  },
}];

local build(arch, testUI) = [{
  kind: 'pipeline',
  name: arch,

  platform: {
    os: 'linux',
    arch: arch,
  },
  depends_on: ['shared'],
  steps: [
           {
             name: 'shared',
             image: 'alpine:' + alpine,
             environment: {
               ARTIFACT_HOST: {
                 from_secret: 'artifact_host',
               },
             },
             commands: [
               './shared/fetch.sh ${DRONE_BUILD_NUMBER} ' + arch,
             ],
           },
         ] + [
           {
             name: tool,
             image: 'debian:bookworm-slim',
             commands: [
               './' + tool + '/build.sh',
             ],
           }
           for tool in tools
         ] + [
           {
             name: tool + ' test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './' + tool + '/test.sh',
             ],
           }
           for tool in tools
           for distro in distros
         ] + [
           {
             name: 'caddy test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './caddy/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'testapp nginx',
             image: 'nginx:' + nginx,
             commands: [
               './test/testapp/nginx/build.sh',
             ],
           },
         ] + [
           {
             name: 'testapp nginx test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './test/testapp/nginx/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'authelia',
             image: 'authelia/authelia:' + authelia,
             commands: [
               './authelia/build.sh',
             ],
           },
         ] + [
           {
             name: 'authelia test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './authelia/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'frp test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './frp/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'backend test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './backend/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'package',
             image: 'debian:bookworm-slim',
             commands: [
               './package.sh ${DRONE_BUILD_NUMBER}',
               './test/testapp/build.sh ',
               './test/streamapp/build.sh',
             ],
           },
         ] + [
           {
             name: 'test ' + distro,
             image: 'python:' + python,
             commands: [
               'cd test',
               './deps.sh',
               'py.test -x -s test.py --distro=' + distro + ' --domain=' + distro + '-' + arch + ' --device-host=' + distro + '-' + arch + ' --app-archive-path=$(realpath ../*.snap) --app=' + name + ' --arch=' + arch + ' --redirect-user=redirect --redirect-password=redirect',
             ],
           }
           for distro in distros
         ] + (if testUI then [
                {
                  name: 'test-ui-' + mode,
                  image: 'mcr.microsoft.com/playwright:' + playwright,
                  environment: {
                    PLAYWRIGHT_APP: name,
                    PLAYWRIGHT_DISTRO: distro_default,
                    PLAYWRIGHT_DEVICE_HOST: distro_default + '-' + arch,
                    PLAYWRIGHT_DOMAIN: distro_default + '-' + arch + '.redirect',
                    PLAYWRIGHT_FULL_DOMAIN: distro_default + '-' + arch + '.redirect',
                    PLAYWRIGHT_DOMAIN_SHORT: distro_default + '-' + arch,
                    PLAYWRIGHT_MAIN_DOMAIN: 'redirect',
                    PLAYWRIGHT_REDIRECT_USER: 'redirect',
                    PLAYWRIGHT_REDIRECT_PASSWORD: 'redirect',
                    PLAYWRIGHT_DEVICE_USER: 'user',
                    PLAYWRIGHT_DEVICE_PASSWORD: 'Password1',
                    PLAYWRIGHT_SSH_USER: 'root',
                    PLAYWRIGHT_SSH_PASSWORD: 'Password1',
                    PLAYWRIGHT_TESTAPP_SNAP: '/testapp.snap',
                    PLAYWRIGHT_PROJECT: mode,
                    PLAYWRIGHT_ARTIFACT_DIR: '/drone/src/artifact',
                    CI: 'true',
                  },
                  commands: [
                    './web/e2e/test.sh ' + mode + ' ' + distro_default + '-' + arch,
                  ],
                  privileged: true,
                }
                for mode in ['desktop', 'mobile']
              ] + [
                {
                  name: 'visual-diff',
                  image: 'alpine:' + alpine,
                  commands: [
                    './visual-diff/visual-diff ci-diff artifact/playwright ' + visual_diff_skip_build,
                  ],
                },
              ] else []) +
         [
           {
             name: 'test-upgrade',
             image: 'python:' + python,
             commands: [
               'APP_ARCHIVE_PATH=$(realpath $(cat package.name))',
               'cd test',
               './deps.sh',
               'py.test -x -s upgrade.py --domain=' + distro_default + '-' + arch + ' --device-host=' + distro_default + '-' + arch + ' --app-archive-path=$APP_ARCHIVE_PATH --app=' + name,
             ],
             privileged: true,
           },
           {
             name: 'publish',
             image: publisher_image,
             environment: {
               SYNCLOUD_TOKEN: {
                 from_secret: 'SYNCLOUD_TOKEN',
               },
             },
             command: ['snap', '-c', '${DRONE_BRANCH}'],
             when: {
               branch: ['stable'],
               event: ['push'],
             },
           },
           {
             name: 'artifact',
             image: 'appleboy/drone-scp:1.6.4',
             settings: {
               host: {
                 from_secret: 'artifact_host',
               },
               username: 'artifact',
               key: {
                 from_secret: 'artifact_key',
               },
               timeout: '2m',
               command_timeout: '2m',
               target: '/home/artifact/repo/' + name + '/${DRONE_BUILD_NUMBER}-' + arch,
               source: 'artifact/*',
               privileged: true,
               strip_components: 1,
             },
             when: {
               status: [
                 'failure',
                 'success',
               ],
               event: [
                 'push',
               ],
             },
           },
         ],
  trigger: {
    event: [
      'push',
    ],
  },
  services: [
    {
      name: distro + '-' + arch,
      image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
      privileged: true,
      volumes: [
        {
          name: 'dbus',
          path: '/var/run/dbus',
        },
        {
          name: 'dev',
          path: '/dev',
        },
      ],
    }
    for distro in distros
  ] + [
    {
      name: 'api.redirect',
      image: 'syncloud/redirect-test-' + arch,
      environment: {
        SOCKET: 'tcp://:80',
        DOMAIN: 'redirect',
      },
    },
  ] + (if testUI then [
    {
      name: 'relay.redirect',
      image: 'snowdreamtech/frps:0.61.1',
      commands: [
        'printf "bindPort = 443\nvhostHTTPSPort = 4443\ntcpMuxHTTPConnectPort = 1337\n" > /etc/frp/frps.toml',
        '/usr/bin/frps -c /etc/frp/frps.toml',
      ],
    },
  ] else []),
  volumes: [
    {
      name: 'dbus',
      host: {
        path: '/var/run/dbus',
      },
    },
    {
      name: 'dev',
      host: {
        path: '/dev',
      },
    },
  ],
}];

shared +
build('amd64', true) +
build('arm64', false) +
build('arm', false)
