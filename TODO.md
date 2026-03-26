# TODO

## Migrate 3rdparty dependencies to in-tree Docker builds

Currently `btrfs` and `openldap` are downloaded as prebuilt binaries from
`github.com/syncloud/3rdparty` releases during the package step (`package.sh`).
This makes version upgrades harder because each dependency lives in a separate
repo with its own release workflow.

**Goal:** Build btrfs and openldap inside Docker containers as part of the
platform CI pipeline (similar to how nginx and authelia are already built
in-tree), so upgrading a dependency is just a version bump in this repo.

**Current downloads in `package.sh`:**
- `https://github.com/syncloud/3rdparty/releases/download/openldap/openldap-${ARCH}.tar.gz`
- `https://github.com/syncloud/3rdparty/releases/download/btrfs/btrfs-${ARCH}.tar.gz`
- `https://github.com/syncloud/3rdparty/releases/download/gptfdisk/gptfdisk-${ARCH}.tar.gz`
