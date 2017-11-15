## Syncloud Platform

Platform provides device Web UI and manages app installation using Sam (Syncloud app manager)

### Running local drone build

Get drone cli binary: http://docs.drone.io/cli-installation/
````
sudo DOCKER_API_VERSION=1.24 arch=amd64 installer=sam /path/to/drone exec
````

Watch drone build processes:
````
watch -n 1 pstree -a $(pgrep -f dockerd)
````

#### Build server

http://build.syncloud.org/syncloud/platform

### Build artifacts

http://artifact.syncloud.org/platform/ci