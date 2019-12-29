## [Syncloud](https://syncloud.org)

Simple self-hosting of popular apps.

It is available as an image or a pre-built [device](https://shop.syncloud.org).

We are open to cooperation with hardware vendors interested in including Syncloud into their products.

### Apps

* [Nextcloud](https://nextcloud.com/): File sharing, calendar, contacts.
* [Diaspora](https://diasporafoundation.org/): Social network.
* [Rocketchat](https://rocket.chat/): Text, voice and video messaging.
* [Mail](https://roundcube.net/): Email messaging with Roundcube web.
* [GOGS](https://gogs.io/): Git source code hosting.
* [Syncthing](https://syncthing.net/): File synchronization between devices.
* [WordPress](https://wordpress.org/): Blogging, mailing lists and forums, media galleries, and online stores.
* [Notes](https://standardnotes.org/): Safe place for your notes, thoughts, and life's work.
* [Users](https://github.com/kakwa/ldapcherr/): User management app.
* [Pi-hole](https://pi-hole.net): Network-wide Ad Blocking.
* [OpenVPN](https://openvpn.net): Secure connection to your Syncloud device over the Internet.

### Download

There are images for various devices and architectures, get one [here](https://github.com/syncloud/platform/wiki).

## For developers

Syncloud image contains the following components:

1. Debian based [linux OS](github.com/syncloud/image).
2. Snap based app [installer](https://github.com/syncloud/snapd).
3. Platform snap package.

Platform provides shared services for all the apps and manages device settings.

### Web UI development

Install [Node.js](https://nodejs.org/en/download)

````
cd www/public
npm install
npm start
````

### Building a package locally
We use Drone build server for automated builds.
The simplest way to build a platform snap package locally is to run [drone cli](http://docs.drone.io/cli-installation):
````
sudo arch=[amd64|arm] /path/to/cli/drone exec
````

### Install a package on a device
````
snap install --devmode /path/to/package.snap
````
