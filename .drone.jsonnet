local name = "platform";
local browser = "firefox";

local build(arch, distro) = {
    kind: "pipeline",
    name: arch + " " + distro,

    platform: {
        os: "linux",
        arch: arch
    },
    steps: [
        {
            name: "version",
            image: "syncloud/build-deps-" + arch,
            commands: [
                "echo $(date +%y%m%d)$DRONE_BUILD_NUMBER > version",
                "echo " + distro + "-" + arch + "-$DRONE_BRANCH > domain"
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
            image: "syncloud/build-deps-" + arch,
            commands: [
                "./build-uwsgi.sh"
            ]
        },
        {
            name: "package",
            image: "syncloud/build-deps-" + arch,
            commands: [
                "VERSION=$(cat version)",
                "./build.sh $VERSION",
                "./integration/testapp/build.sh "
            ]
        },
        {
            name: "test-unit",
            image: "syncloud/build-deps-" + arch,
            commands: [
                "./unit-test.sh",
            ]
        },
        {
            name: "test-intergation",
            image: "python:3.9-buster",
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client netcat rustc apache2-utils",
              "./integration/wait-ssh.sh",
              "mkdir -p /var/snap/platform/common",
              "sshpass -p syncloud ssh -o StrictHostKeyChecking=no -fN -L /var/snap/platform/common/api.socket:/var/snap/platform/common/api.socket root@device",
              "pip install -r dev_requirements.txt",
              "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
              "DOMAIN=$(cat domain)",
              "cd integration",
              "py.test -x -s verify.py --domain=$DOMAIN --app-archive-path=$APP_ARCHIVE_PATH --device-host=device --app=" + name
            ]
        }
    ] + ( if arch == "arm" then [] else [
        {
            name: "test-ui-desktop",
            image: "python:3.9-buster",
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "DOMAIN=$(cat domain)",
              "cd integration",
              "py.test -x -s test-ui.py --ui-mode=desktop --domain=$DOMAIN --device-host=device --app=" + name + " --browser=" + browser,
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        },
        {
            name: "test-ui-mobile",
            image: "python:3.9-buster",
            commands: [
              "apt-get update && apt-get install -y sshpass openssh-client",
              "pip install -r dev_requirements.txt",
              "DOMAIN=$(cat domain)",
              "cd integration",
              "py.test -x -s test-ui.py --ui-mode=mobile --domain=$DOMAIN --device-host=device --app=" + name + " --browser=" + browser,
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        }
    ]) + [
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
                target: "/home/artifact/repo/" + name + "/${DRONE_BUILD_NUMBER}-" + distro + "-"  + arch,
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
            name: "device",
            image: "syncloud/platform-" + distro + '-' + arch,
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
    ] + if arch == "arm" then [] else [{
            name: "selenium",
            image: "selenium/standalone-" + browser + ":4.0.0-beta-3-prerelease-20210402",
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        }
    ],
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
        }
    ]
};

[
    build("arm", "jessie"),
    build("amd64", "jessie"),
    build("arm", "buster"),
    build("amd64", "buster")
]
