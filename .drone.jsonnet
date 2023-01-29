local name = "platform";
local browser = "chrome";
local go = "1.18.2-buster";
local node = "16.10.0";

local build(arch, testUI, uwsgiDistro) = [{
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
                "echo $DRONE_BUILD_NUMBER > version",
            ]
        },
        {
            name: "build web",
            image: "node:" + node + "-alpine3.12",
            environment: {
                NODE_OPTIONS: '--max_old_space_size=2048',
            },
            commands: [
                "mkdir -p build/snap",
                "cd www",
		        "npm config set fetch-retry-mintimeout 20000",
                "npm config set fetch-retry-maxtimeout 120000",
                "npm install",
                "npm run test",
                "npm run lint",
                "npm run build",
                "cp -r dist ../build/snap/www"
            ]
        },
        {
            name: "build",
            image: "golang:" + go,
            commands: [
                "cd backend",
                "go test ./... -coverprofile cover.out",
                "go tool cover -func cover.out",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/snap/bin/backend ./cmd/backend",
                "../build/snap/bin/backend -h",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/snap/bin/api ./cmd/api",
                "../build/snap/bin/api -h",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/snap/bin/cli ./cmd/cli",
                "../build/snap/bin/cli -h"
            ]
        },
        {
            name: "build api test",
            image: "golang:" + go,
            commands: [
                "cd integration/api",
                "go test -c -o api.test"
            ]
        },
        {
            name: "build uwsgi",
            image: "debian:" + uwsgiDistro + "-slim",
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
        }
    ] + [
        {
            name: "test-intergation-buster",
            image: "python:3.8-slim-buster",
            commands: [
              "cd integration",
              "./deps.sh",
              "py.test -x -s verify.py --distro=buster --domain="+arch+"-buster --app-archive-path=$(realpath ../*.snap) --app=" + name + " --arch=" + arch + " --redirect-user=redirect --redirect-password=redirect"
            ]
        }
    ] + ( if testUI then [
 {
        name: "selenium-video",
        image: "selenium/video:ffmpeg-4.3.1-20220208",
        detach: true,
        environment: {
            DISPLAY_CONTAINER_NAME: "selenium",
            FILE_NAME: "video.mkv"
        },
        volumes: [
            {
                name: "shm",
                path: "/dev/shm"
            },
           {
                name: "videos",
                path: "/videos"
            }
        ]
    }] + [
        {
            name: "test-ui-" + mode + "-" + distro,
            image: "python:3.8-slim-buster",
            commands: [
              "cd integration",
              "./deps.sh",
              "py.test -x -s test-ui.py --distro=" + distro + " --ui-mode=" + mode + " --domain="+arch+"-"+distro+" --redirect-user=redirect --redirect-password=redirect --app=" + name + " --browser=" + browser
            ],
            privileged: true,
            volumes: [{
                name: "videos",
                path: "/videos"
            }]
        } 
        for mode in ["desktop", "mobile"]
        for distro in ["buster"]
    ] else []) + 
   [
    {
        name: "test-upgrade",
        image: "python:3.8-slim-buster",
        commands: [
          "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
          "cd integration",
          "./deps.sh",
          "py.test -x -s test-upgrade.py --distro=buster  --domain="+arch+"-buster --app-archive-path=$APP_ARCHIVE_PATH --app=" + name
        ],
        privileged: true,
        volumes: [{
            name: "videos",
            path: "/videos"
        }]
    } 
] + [
        {
            name: "upload",
            image: "debian:buster-slim",
            environment: {
                AWS_ACCESS_KEY_ID: {
                    from_secret: "AWS_ACCESS_KEY_ID"
                },
                AWS_SECRET_ACCESS_KEY: {
                    from_secret: "AWS_SECRET_ACCESS_KEY"
                }
            },
            commands: [
              "PACKAGE=$(cat package.name)",
              "apt update && apt install -y wget",
              "wget https://github.com/syncloud/snapd/releases/download/1/syncloud-release-" + arch + " -O release --progress=dot:giga",
              "chmod +x release",
              "./release publish -f $PACKAGE -b $DRONE_BRANCH"
            ],
            when: {
                branch: [
                    "stable",
                    "master"
                ],
                event: [
                    "push"
                ]
            }
        },
        {
            name: "artifact",
            image: "appleboy/drone-scp:1.6.4",
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
                privileged: true,
		        strip_components: 1,
                volumes: [
                   {
                        name: "videos",
                        path: "/drone/src/artifact/videos"
                    }
               ]
           },
           when: {
              status: [
                  "failure",
                  "success"
              ],
              event: [
                  "push"
              ]
            }
        }
    ],
    trigger: {
      event: [
        "push",
        "pull_request"
      ]
    },
    services: [
        {
            name: arch + "-buster",
            image: "syncloud/bootstrap-buster-" + arch,
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
            name: "api.redirect",
            image: "syncloud/redirect-test-" + arch,
            environment: {
                SOCKET: "tcp://:80",
                DOMAIN: "redirect"
            }
        }
    ] + ( if testUI then [{
            name: "selenium",
            image: "selenium/standalone-" + browser + ":4.1.2-20220208",
            environment: {
                SE_NODE_SESSION_TIMEOUT: "999999",
                START_XVFB: "true"
            },
               volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        }
    ] else [] ),
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
        },
      {
            name: "videos",
            temp: {}
        }
    ]
},
 {
     kind: "pipeline",
     type: "docker",
     name: "promote-" + arch,
     platform: {
         os: "linux",
         arch: arch
     },
     steps: [
     {
             name: "promote",
             image: "debian:buster-slim",
             environment: {
                 AWS_ACCESS_KEY_ID: {
                     from_secret: "AWS_ACCESS_KEY_ID"
                 },
                 AWS_SECRET_ACCESS_KEY: {
                     from_secret: "AWS_SECRET_ACCESS_KEY"
                 }
             },
             commands: [
               "apt update && apt install -y wget",
               "wget https://github.com/syncloud/snapd/releases/download/1/syncloud-release-" + arch + " -O release --progress=dot:giga",
               "chmod +x release",
               "./release promote -n " + name + " -a $(dpkg --print-architecture)"
             ]
       }
      ],
      trigger: {
       event: [
         "promote"
       ]
     }
 }];

build("amd64", true, "bookworm") +
build("arm64", false, "bookworm") +
build("arm", false, "buster")
