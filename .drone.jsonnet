local name = "platform";
local browser = "firefox";

local build(arch, testUI) = {
    kind: "pipeline",
    name: arch,

    platform: {
        os: "linux",
        arch: arch
    },
    steps: [
        {
            name: "version",
            image: "debian:buster-slim",
            commands: [
                "echo $(date +%y%m%d)$DRONE_BUILD_NUMBER > version",
                "echo " + arch + "-$DRONE_BRANCH > domain"
            ]
        },
        {
            name: "build web",
            image: "node:16.1.0",
            commands: [
                "mkdir -p build/platform",
                "cd www",
                "npm install --unsafe-perm=true",
                "npm run test",
                "npm run lint",
                "npm run build",
                "cp -r dist ../build/platform/www"
            ]
        },
        {
            name: "build backend",
            image: "golang:1.16.4",
            commands: [
                "cd backend",
                "go test ./... -cover",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/platform/bin/backend cmd/backend/main.go",
                "../build/platform/bin/backend -h",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/platform/bin/cli cmd/cli/main.go",
                "../build/platform/bin/cli -h"
            ]
        },
        {
            name: "build uwsgi",
            image: "debian:buster-slim",
            commands: [
                "./build-uwsgi.sh"
            ],
            volumes: [
                {
                    name: "docker",
                    path: "/usr/bin/docker"
                },
                {
                    name: "docker.sock",
                    path: "/var/run/docker.sock"
                }
            ]
        },
        {
            name: "package",
            image: "debian:buster-slim",
            commands: [
                "VERSION=$(cat version)",
                "./build.sh $VERSION",
                "./integration/testapp/build.sh "
            ]
        },
        {
            name: "test-unit",
            image: "python:3.9-buster",
            commands: [
              "apt update",
              "apt install -y build-essential libsasl2-dev libldap2-dev libssl-dev libjansson-dev libltdl7 libnss3",
              "pip install -r requirements.txt",
              "pip install -r dev_requirements.txt",
              "cd src",
              "py.test test"
            ]
        },
        {
            name: "test-intergation-jessie",
            image: "python:3.9-buster",
            environment: {
                REDIRECT_USER: {
                    from_secret: "REDIRECT_USER"
                },
                REDIRECT_PASSWORD: {
                    from_secret: "REDIRECT_PASSWORD"
                }
            },
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client netcat rustc apache2-utils",
              "./integration/wait-ssh.sh device-jessie",
              "mkdir -p /var/snap/platform/common",
              "sshpass -p syncloud ssh -o StrictHostKeyChecking=no -fN -L /var/snap/platform/common/api.socket:/var/snap/platform/common/api.socket root@device-jessie",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s verify.py --distro=jessie --domain=$(cat ../domain) --app-archive-path=$(realpath ../*.snap) --device-host=device-jessie --app=" + name + " --arch=" + arch + " --redirect-user=$REDIRECT_USER --redirect-password=$REDIRECT_PASSWORD"
            ]
        },
        {
            name: "test-intergation-buster",
            image: "python:3.9-buster",
            environment: {
                REDIRECT_USER: {
                    from_secret: "REDIRECT_USER"
                },
                REDIRECT_PASSWORD: {
                    from_secret: "REDIRECT_PASSWORD"
                }
            },
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client netcat rustc apache2-utils",
              "./integration/wait-ssh.sh device-buster",
              "mkdir -p /var/snap/platform/common",
              "sshpass -p syncloud ssh -o StrictHostKeyChecking=no -fN -L /var/snap/platform/common/api.socket:/var/snap/platform/common/api.socket root@device-buster",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s verify.py --distro=buster --domain=$(cat ../domain) --app-archive-path=$(realpath ../*.snap) --device-host=device-buster --app=" + name + " --arch=" + arch + " --redirect-user=$REDIRECT_USER --redirect-password=$REDIRECT_PASSWORD"
            ]
        }
    ] + ( if testUI then [
        {
            name: "test-ui-desktop-jessie",
            image: "python:3.9-buster",
            environment: {
                REDIRECT_USER: {
                    from_secret: "REDIRECT_USER"
                },
                REDIRECT_PASSWORD: {
                    from_secret: "REDIRECT_PASSWORD"
                }
            },
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s test-ui.py --distro=jessie --ui-mode=desktop --domain=$(cat ../domain) --device-host=device-jessie -redirect-user=$REDIRECT_USER --redirect-password=$REDIRECT_PASSWORD --app=" + name + " --browser=" + browser
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        },
        {
            name: "test-ui-desktop-buster",
            image: "python:3.9-buster",
            environment: {
                REDIRECT_USER: {
                    from_secret: "REDIRECT_USER"
                },
                REDIRECT_PASSWORD: {
                    from_secret: "REDIRECT_PASSWORD"
                }
            },
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s test-ui.py --distro=buster --ui-mode=desktop --domain=$(cat ../domain) --device-host=device-buster -redirect-user=$REDIRECT_USER --redirect-password=$REDIRECT_PASSWORD --app=" + name + " --browser=" + browser,
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        },
        {
            name: "test-ui-mobile-jessie",
            image: "python:3.9-buster",
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s test-ui.py --distro=jessie --ui-mode=mobile --domain=$(cat ../domain) --device-host=device-jessie --app=" + name + " --browser=" + browser,
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        },
        {
            name: "test-ui-mobile-buster",
            image: "python:3.9-buster",
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "cd integration",
              "py.test -x -s test-ui.py --distro=buster --ui-mode=mobile --domain=$(cat ../domain) --device-host=device-buster --app=" + name + " --browser=" + browser,
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        }
    ] else []) + [
        {
            name: "upload",
            image: "python:3.9-buster",
            environment: {
                AWS_ACCESS_KEY_ID: {
                    from_secret: "AWS_ACCESS_KEY_ID"
                },
                AWS_SECRET_ACCESS_KEY: {
                    from_secret: "AWS_SECRET_ACCESS_KEY"
                }
            },
            commands: [
              "VERSION=$(cat version)",
              "PACKAGE=$(cat package.name)",
              "pip install syncloud-lib s3cmd",
              "syncloud-upload.sh " + name + " $DRONE_BRANCH $VERSION $PACKAGE"
            ]
        },
     {
            name: "artifact",
            image: "appleboy/drone-scp",
            settings: {
                host: {
                    from_secret: "artifact_host"
                },
                username: "artifact",
                key: {
                    from_secret: "artifact_key"
                },
                timeout: "2m",
                command_timeout: "2m",
                target: "/home/artifact/repo/" + name + "/${DRONE_BUILD_NUMBER}-" + arch,
                source: "artifact/*",
		             strip_components: 1
            },
            when: {
              status: [ "failure", "success" ]
            }
        }
    ],
    services: [
        {
            name: "device-jessie",
            image: "syncloud/platform-jessie-" + arch,
            privileged: true,
            volumes: [
                {
                    name: "dbus",
                    path: "/var/run/dbus"
                },
                {
                    name: "dev",
                    path: "/dev"
                }
            ]
        },
        {
            name: "device-buster",
            image: "syncloud/platform-buster-" + arch,
            privileged: true,
            volumes: [
                {
                    name: "dbus",
                    path: "/var/run/dbus"
                },
                {
                    name: "dev",
                    path: "/dev"
                }
            ]
        }
    ] + if testUI then [{
            name: "selenium",
            image: "selenium/standalone-" + browser + ":4.0.0-beta-3-prerelease-20210402",
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        }
    ] else [],
    volumes: [
        {
            name: "dbus",
            host: {
                path: "/var/run/dbus"
            }
        },
        {
            name: "dev",
            host: {
                path: "/dev"
            }
        },
        {
            name: "shm",
            temp: {}
        },
        {
            name: "docker",
            host: {
                path: "/usr/bin/docker"
            }
        },
        {
            name: "docker.sock",
            host: {
                path: "/var/run/docker.sock"
            }
        }
    ]
};

[
    build("arm", false),
    build("arm64", false),
    build("amd64", true)
]
