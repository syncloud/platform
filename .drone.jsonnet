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
                "echo build-$(uname -m)-$DRONE_BRANCH > domain"
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
    }]
};

[
    build("arm"),
    build("amd64")
]