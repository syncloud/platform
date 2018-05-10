## Syncloud (syncloud.org)

Syncloud brings popular apps to your home.
It is available as an image or pre-installed devcie.

### Apps

* Nextcloud: File sharing, calendar, contacts.
* Diaspora: Social network.
* Rocketchat: Text, voice and video messaging.
* Mail: Email messaging with Roundcube web.
* GOGS: Git source code hosting.


## For developers

Syncloud platform manages the installation and device settings.

### Running local drone build

Get drone cli binary: http://docs.drone.io/cli-installation/
````
sudo DOCKER_API_VERSION=1.24 arch=amd64 installer=sam /path/to/drone exec
````

Watch drone build processes:
````
watch -n 1 pstree -a $(pgrep -f dockerd)
````

### Build server

http://build.syncloud.org/syncloud/platform

### Build artifacts (screenshots, system logs)

http://artifact.syncloud.org/platform/ci