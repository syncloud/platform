local build(arch) = {
    kind: "pipeline",
    name: arch,

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
                "echo platform > name",
                "echo build-" + arch + "-$DRONE_BRANCH > domain"
            ]
        },
        {
            name: "build",
            image: "syncloud/build-deps-" + arch,
            commands: [
                "NAME=$(cat name)",
                "VERSION=$(cat version)",
                "./build.sh $NAME $VERSION"
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
            name: "test-js",
            image: "syncloud/build-deps-" + arch,
            commands: [
              "export PATH=$(pwd)/node/bin:$PATH",
              "cd www/public",
              "npm run test"
            ]
        },
        {
            name: "test-intergation",
            image: "syncloud/build-deps-" + arch,
            commands: [
              "pip2 install -r dev_requirements.txt",
              "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
              "DOMAIN=$(cat domain)",
              "NAME=$(cat name)",
              "cd integration",
              "py.test -x -s verify.py --domain=$DOMAIN --app-archive-path=$APP_ARCHIVE_PATH --device-host=device --app=$NAME"
            ]
        },
        if arch == "arm" then {} else
        {
            name: "test-ui",
            image: "syncloud/build-deps-" + arch,
            commands: [
              "pip2 install -r dev_requirements.txt",
              "DOMAIN=$(cat domain)",
              "NAME=$(cat name)",
              "cd integration",
              "xvfb-run -l --server-args='-screen 0, 1024x4096x24' py.test -x -s test-ui.py --domain=$DOMAIN --device-host=device --app=$NAME",
            ],
            volumes: [{
                name: "shm",
                path: "/dev/shm"
            }]
        },
        {
            name: "upload",
            image: "syncloud/build-deps-" + arch,
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
              "NAME=$(cat name)",
              "PACKAGE=$(cat package.name)",
              "pip2 install -r dev_requirements.txt",
              "syncloud-upload.sh $NAME $DRONE_BRANCH $VERSION $PACKAGE"
            ]
        },
        {
            name: "ci-artifact",
            image: "syncloud/build-deps-" + arch,
            environment: {
                ARTIFACT_SSH_KEY: {
                    from_secret: "ARTIFACT_SSH_KEY"
                }
            },
            commands: [
                "NAME=$(cat name)",
                "PACKAGE=$(cat package.name)",
                "pip2 install -r dev_requirements.txt",
                "syncloud-upload-artifact.sh $NAME integration/log $DRONE_BUILD_NUMBER-$(dpkg --print-architecture)",
                "syncloud-upload-artifact.sh $NAME integration/screenshot $DRONE_BUILD_NUMBER-$(dpkg --print-architecture)",
                "syncloud-upload-artifact.sh $NAME $PACKAGE $DRONE_BUILD_NUMBER-$(dpkg --print-architecture)"
            ],
            when: {
              status: [ "failure", "success" ]
            }
        }
    ],
    services: [{
        name: "device",
        image: "syncloud/systemd-" + arch,
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
    }],
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
    build("arm"),
    build("amd64")
]
