local name = 'platform';
local playwright = 'v1.49.1-jammy';
local go = '1.25.8';
local node = '22.16.0';
local deployer = 'https://github.com/syncloud/store/releases/download/4/syncloud-release';
local authelia = '4.39.15';
local distro_default = 'buster';
local distros = ['bookworm', 'buster'];
local bootstrap = '25.02';
local nginx = '1.24.0';
local python = '3.12-slim-bookworm';
local alpine = '3.21';
local visual_diff_skip_build = '2758';

local build(arch, testUI) = [{
  kind: 'pipeline',
  name: arch,

  platform: {
    os: 'linux',
    arch: arch,
  },
  steps: [
           {
             name: 'version',
             image: 'debian:bookworm-slim',
             commands: [
               'echo $DRONE_BUILD_NUMBER > version',
             ],
           },
           {
             name: 'gptfdisk',
             image: 'debian:bookworm-slim',
             commands: [
               './gptfdisk/build.sh',
             ],
           },
         ] + [
           {
             name: 'gptfdisk test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './gptfdisk/test.sh',
             ],
           }
           for distro in distros
         ] + [
           {
             name: 'nginx',
             image: 'nginx:' + nginx,
             commands: [
               './nginx/build.sh',
             ],
           },
         ] + [
           {
             name: 'nginx test ' + distro,
             image: 'syncloud/bootstrap-' + distro + '-' + arch + ':' + bootstrap,
             commands: [
               './nginx/test.sh',
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
             name: 'build web',
             image: 'node:' + node,
             environment: {
               NODE_OPTIONS: '--max_old_space_size=2048',
             },
             commands: [
               'mkdir -p build/snap',
               'cd web/platform',
               'npm config ls -l',
               'npm config set fetch-retry-mintimeout 200000',
               'npm config set fetch-retry-maxtimeout 1200000',
               'npm install',
               'npm update browserslist',
               'npm run test',
               'npm run lint',
               'npm run build',
               'mkdir -p ../../build/snap/web',
               'cp -r dist ../../build/snap/web/platform',
               'cd ../login',
               'npm install',
               'npm run build',
               'cp -r dist ../../build/snap/web/login',
             ],
           },
           {
             name: 'build',
             image: 'golang:' + go,
             commands: [
               'cd backend',
               'for i in 1 2 3; do go mod download && break || sleep 5; done',
               'go test ./... -coverprofile cover.out',
               'go tool cover -func cover.out',
               "CGO_ENABLED=0 go build -o ../build/snap/bin/backend ./cmd/backend",
               '../build/snap/bin/backend -h',
               "CGO_ENABLED=0 go build -o ../build/snap/bin/api ./cmd/api",
               '../build/snap/bin/api -h',
               "CGO_ENABLED=0 go build -o ../build/snap/bin/cli ./cmd/cli",
               '../build/snap/bin/cli -h',
               "CGO_ENABLED=0 go build -o ../build/snap/meta/hooks/install ./cmd/install",
               '../build/snap/meta/hooks/install -h',
               "CGO_ENABLED=0 go build -o ../build/snap/meta/hooks/post-refresh ./cmd/post-refresh",
               '../build/snap/meta/hooks/post-refresh -h',
               "CGO_ENABLED=0 go build -o ../build/snap/meta/hooks/configure ./cmd/configure",
               '../build/snap/meta/hooks/configure -h',
               "CGO_ENABLED=0 go build -o ../build/snap/bin/login ./cmd/login",
               'cd ../visual-diff',
               'CGO_ENABLED=0 go build -o visual-diff ./cmd',
             ],
           },
           {
             name: 'build api test',
             image: 'golang:' + go,
             commands: [
               'cd test/api',
               "CGO_ENABLED=0 go test -c -o api.test",
             ],
           },
           {
             name: 'build external app',
             image: 'golang:' + go,
             commands: [
               'cd test/externalapp',
               "CGO_ENABLED=0 go build -o externalapp",
             ],
           },
           {
             name: 'build testapp cli',
             image: 'golang:' + go,
             commands: [
               'cd test/testapp/cli',
               'CGO_ENABLED=0 go build -o ../build/meta/hooks/install ./cmd/install',
               'CGO_ENABLED=0 go build -o ../build/meta/hooks/configure ./cmd/configure',
               'CGO_ENABLED=0 go build -o ../build/bin/cli ./cmd/cli',
               'CGO_ENABLED=0 go build -o ../build/bin/backend ./cmd/backend',
             ],
           },
           {
             name: 'package',
             image: 'debian:bookworm-slim',
             commands: [
               'VERSION=$(cat version)',
               './package.sh $VERSION',
               './test/testapp/build.sh ',
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
                    'apt-get update && apt-get install -y sshpass sudo',
                    'DOMAIN="' + distro_default + '-' + arch + '"',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/auth.$DOMAIN.redirect/g" >> /etc/hosts',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/$DOMAIN.redirect/g" >> /etc/hosts',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/unknown.$DOMAIN.redirect/g" >> /etc/hosts',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/externalapp.$DOMAIN.redirect/g" >> /etc/hosts',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/testapp.$DOMAIN.redirect/g" >> /etc/hosts',
                    'getent hosts $DOMAIN | sed "s/$DOMAIN/files.$DOMAIN.redirect/g" >> /etc/hosts',
                    'cat /etc/hosts',
                    'cd web/e2e',
                    'npm ci',
                    'npx playwright test --project=' + mode,
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
             name: 'upload',
             image: 'debian:bookworm-slim',
             environment: {
               AWS_ACCESS_KEY_ID: {
                 from_secret: 'AWS_ACCESS_KEY_ID',
               },
               AWS_SECRET_ACCESS_KEY: {
                 from_secret: 'AWS_SECRET_ACCESS_KEY',
               },
               SYNCLOUD_TOKEN: {
                 from_secret: 'SYNCLOUD_TOKEN',
               },
             },
             commands: [
               'PACKAGE=$(cat package.name)',
               'apt update && apt install -y wget',
               'wget ' + deployer + '-' + arch + ' -O release --progress=dot:giga',
               'chmod +x release',
               './release publish -f $PACKAGE -b $DRONE_BRANCH',
             ],
             when: {
               branch: ['stable', 'master'],
               event: ['push'],
             },
           },
           {
             name: 'promote',
             image: 'debian:bookworm-slim',
             environment: {
               AWS_ACCESS_KEY_ID: {
                 from_secret: 'AWS_ACCESS_KEY_ID',
               },
               AWS_SECRET_ACCESS_KEY: {
                 from_secret: 'AWS_SECRET_ACCESS_KEY',
               },
               SYNCLOUD_TOKEN: {
                 from_secret: 'SYNCLOUD_TOKEN',
               },
             },
             commands: [
               'apt update && apt install -y wget',
               'wget ' + deployer + '-' + arch + ' -O release --progress=dot:giga',
               'chmod +x release',
               './release promote -n ' + name + ' -a $(dpkg --print-architecture)',
             ],
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
      'pull_request',
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
  ],
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
    {
      name: 'docker',
      host: {
        path: '/usr/bin/docker',
      },
    },
    {
      name: 'docker.sock',
      host: {
        path: '/var/run/docker.sock',
      },
    },
  ],
}];

build('amd64', true) +
build('arm64', false) +
build('arm', false)
