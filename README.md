## Syncloud (https://syncloud.org)

Simple self-hosting of popular apps.

It is available as an image or a pre-built [device](https://shop.syncloud.org).

We are open to cooperation with hardware vendors interested in including Syncloud into their products.

### Apps

https://syncloud.org/apps.html

### Download

There are images for various devices and architectures, get one [here](https://github.com/syncloud/platform/wiki).

## For developers

Syncloud image contains the following components:

1. Debian based [linux OS](https://github.com/syncloud/image).
2. Snap based app [installer](https://github.com/syncloud/snapd).
3. Platform snap package.

Platform provides shared services for all the apps and manages device settings.

### Web UI development

Install [Node.js](https://nodejs.org/en/download)

````
cd www
npm i
npm run dev
````

### Building a package locally
We use Drone build server for automated builds.
The simplest way to build a platform snap package locally is to run [drone cli](http://docs.drone.io/cli-installation):
````
/path/to/cli/drone jsonnet --stream
sudo /path/to/cli/drone exec --pipeline=[amd64|arm64|arm] --trusted
````

### Install a package on a device
````
snap install --devmode /path/to/package.snap
````
