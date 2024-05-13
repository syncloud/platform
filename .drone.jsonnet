local name = "platform";
local browser = "chrome";
local selenium = '4.19.0-20240328';
local go = "1.21.4";
local node = "16.10.0";
local deployer = "https://github.com/syncloud/store/releases/download/4/syncloud-release";
local authelia = "4.38.8";

local build(arch, testUI) = [{
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
              name: 'authelia',
              image: "authelia/authelia:" + authelia,
              commands: [
                './authelia/package.sh',
              ],

            },
        {
            name: 'authelia test',
            image: "debian:buster-slim",
            commands: [
                './authelia/test.sh',
            ],
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
                "npm config ls -l",
		        "npm config set fetch-retry-mintimeout 200000",
                "npm config set fetch-retry-maxtimeout 1200000",
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
                "../build/snap/bin/cli -h",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/snap/meta/hooks/install ./cmd/install",
                "../build/snap/meta/hooks/install -h",
                "go build -ldflags '-linkmode external -extldflags -static' -o ../build/snap/meta/hooks/post-refresh ./cmd/post-refresh",
                "../build/snap/meta/hooks/post-refresh -h"
            ]
        },
        {
            name: "build api test",
            image: "golang:" + go,
            commands: [
                "cd test/api",
                "go test -c -ldflags '-linkmode external -extldflags -static' -o api.test"
            ]
        },
        {
            name: "package",
            image: "debian:buster-slim",
            commands: [
                "VERSION=$(cat version)",
                "./package.sh $VERSION",
                "./test/testapp/build.sh "
            ]
        },
        {
            name: "test",
            image: "python:3.8-slim-buster",
            commands: [
              "cd test",
              "./deps.sh",
              "py.test -x -s test.py --domain="+arch+" --app-archive-path=$(realpath ../*.snap) --app=" + name + " --arch=" + arch + " --redirect-user=redirect --redirect-password=redirect"
            ]
        }
    ] + ( if testUI then [
        {
            name: "selenium",
            image: "selenium/standalone-" + browser + ":" + selenium,
            detach: true,
            environment: {
                SE_NODE_SESSION_TIMEOUT: "999999",
                START_XVFB: "true"
            },
               volumes: [{
                name: "shm",
                path: "/dev/shm"
            }],
            commands: [
                "cat /etc/hosts",
                "getent hosts " + arch + " | sed 's/" + arch +"/auth." + arch + ".redirect/g' | sudo tee -a /etc/hosts",
                "getent hosts " + arch + " | sed 's/" + arch +"/" + arch + ".redirect/g' | sudo tee -a /etc/hosts",
                "cat /etc/hosts",
                "/opt/bin/entry_point.sh"
            ]
         },
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
            name: "test-ui-" + mode,
            image: "python:3.8-slim-buster",
            commands: [
              "cd test",
              "./deps.sh",
              "py.test -x -s test-ui.py --ui-mode=" + mode + " --domain="+arch+" --redirect-user=redirect --redirect-password=redirect --app=" + name + " --browser=" + browser
            ],
            privileged: true,
            volumes: [{
                name: "videos",
                path: "/videos"
            }]
        } 
        for mode in ["desktop", "mobile"]
    ] else []) + 
   [
    {
        name: "test-upgrade",
        image: "python:3.8-slim-buster",
        commands: [
          "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
          "cd test",
          "./deps.sh",
          "py.test -x -s test-upgrade.py --domain="+arch+" --app-archive-path=$APP_ARCHIVE_PATH --app=" + name
        ],
        privileged: true,
        volumes: [{
            name: "videos",
            path: "/videos"
        }]
    },
    {
             name: "upload",
             image: "debian:buster-slim",
             environment: {
                 AWS_ACCESS_KEY_ID: {
                     from_secret: "AWS_ACCESS_KEY_ID"
                 },
                 AWS_SECRET_ACCESS_KEY: {
                     from_secret: "AWS_SECRET_ACCESS_KEY"
                 },
                 SYNCLOUD_TOKEN: {
                          from_secret: "SYNCLOUD_TOKEN"
                      }
             },
             commands: [
                 "PACKAGE=$(cat package.name)",
                 "apt update && apt install -y wget",
                 "wget " + deployer + "-" + arch + " -O release --progress=dot:giga",
                 "chmod +x release",
                 "./release publish -f $PACKAGE -b $DRONE_BRANCH"
             ],
             when: {
                 branch: ["stable", "master"],
     	    event: [ "push" ]
     }
         },
         {
                 name: "promote",
                 image: "debian:buster-slim",
                 environment: {
                     AWS_ACCESS_KEY_ID: {
                         from_secret: "AWS_ACCESS_KEY_ID"
                     },
                     AWS_SECRET_ACCESS_KEY: {
                         from_secret: "AWS_SECRET_ACCESS_KEY"
                     },
                      SYNCLOUD_TOKEN: {
                          from_secret: "SYNCLOUD_TOKEN"
                      }
                 },
                 commands: [
                   "apt update && apt install -y wget",
                   "wget " + deployer + "-" + arch + " -O release --progress=dot:giga",
                   "chmod +x release",
                   "./release promote -n " + name + " -a $(dpkg --print-architecture)"
                 ],
                 when: {
                     branch: ["stable"],
                     event: ["push"]
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
            name: arch,
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
}];

build("amd64", true) +
build("arm64", false) +
build("arm", false)
