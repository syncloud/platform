# Debugging CI failures

When a CI build fails, always start by identifying the failing step:
```
curl -s "http://ci.syncloud.org:8080/api/repos/syncloud/platform/builds/{N}" | python3 -c "
import json,sys
b=json.load(sys.stdin)
for stage in b.get('stages',[]):
    for step in stage.get('steps',[]):
        if step.get('status') == 'failure':
            print(step.get('name'), '-', step.get('status'))
"
```

Then get the step log (stage=pipeline index, step=step number):
```
curl -s "http://ci.syncloud.org:8080/api/repos/syncloud/platform/builds/{N}/logs/{stage}/{step}" | python3 -c "
import json,sys; [print(l.get('out',''), end='') for l in json.load(sys.stdin)]
" | tail -80
```

# CI

http://ci.syncloud.org:8080/syncloud/platform

CI is Drone CI (JS SPA). Check builds via API:
```
curl -s "http://ci.syncloud.org:8080/api/repos/syncloud/platform/builds?limit=5"
```

## CI Artifacts

Artifacts are served at `http://ci.syncloud.org:8081` (returns JSON directory listings).

Browse the top level for a build (returns distro subdirs + snap file):
```
curl -s "http://ci.syncloud.org:8081/files/platform/{build}-{arch}/"
```

Each distro dir contains `app/`, `platform/`, and for upgrade/UI tests also `desktop/`, `refresh.journalctl.log`, `video.mkv`:
```
curl -s "http://ci.syncloud.org:8081/files/platform/{build}-{arch}/{distro}/"
curl -s "http://ci.syncloud.org:8081/files/platform/{build}-{arch}/{distro}/app/"
curl -s "http://ci.syncloud.org:8081/files/platform/{build}-{arch}/{distro}/desktop/"
```

Directory structure:
```
{build}-{arch}/
  {distro}/
    app/
      journalctl.log          # full journal from integration test teardown
      ps.log, netstat.log     # process/network state at teardown
    platform/                 # platform logs
    refresh.journalctl.log    # full journal from upgrade test (pre/post-refresh)
  distro/
    desktop/                  # UI test artifacts (amd64 only)
      journalctl.log
      screenshot/
        {test-name}-desktop.png
        {test-name}-desktop.html.log
    video.mkv                 # selenium recording
```

Download a file directly:
```
curl -O "http://ci.syncloud.org:8081/files/platform/282-amd64/bookworm/refresh.journalctl.log"
curl -O "http://ci.syncloud.org:8081/files/platform/282-amd64/bookworm/app/journalctl.log"
curl -O "http://ci.syncloud.org:8081/files/platform/282-amd64/distro/desktop/journalctl.log"
curl -O "http://ci.syncloud.org:8081/files/platform/282-amd64/distro/desktop/screenshot/exception-desktop.png"
```

# Project Structure

- **Snap-based platform** providing self-hosting OS, app installer, and platform services for Syncloud
- Architectures: amd64, arm64, arm
- Distros: bookworm, buster
- CI pipelines defined in `.drone.jsonnet`

## Key directories

- `backend/` — Go backend services (API server, backend server, CLI, snap hooks)
  - `cmd/` — Executables: api, backend, cli, install, post-refresh
  - `rest/` — REST API endpoints (gorilla/mux)
  - `auth/` — Authentication: OIDC, LDAP, Authelia integration
  - `config/` — Configuration management (SQLite)
  - `storage/` — Disk/btrfs management
  - `installer/` — App installation service
  - `backup/` — Backup/restore
  - Built with `CGO_ENABLED` and static linking
- `web/platform/` — Vue 3 frontend (Element Plus, Vite, TypeScript)
- `config/` — Configuration templates (authelia, ldap, nginx, errors)
- `authelia/` — Authelia auth server packaging
- `nginx/` — Nginx build/test scripts
- `bin/` — Platform utility scripts
- `meta/snap.yaml` — Snap metadata (services: nginx-public, openldap, backend, api, authelia, cli)
- `test/` — Python integration tests (pytest), Selenium UI tests, Go API tests
- `package.sh` — Creates snap package

## Build pipeline steps (per arch)

1. `version` — writes build number
2. `nginx` / `nginx test {distro}` — build and test nginx (tested on both bookworm and buster)
3. `authelia` / `authelia test` — package and test Authelia auth server
4. `build web` — npm install, test, lint, build Vue frontend
5. `build` — compile Go backend binaries (with `go test ./...` coverage)
6. `build api test` — compile Go API integration test binary
7. `package` — create `.snap` file + test app
8. `test {distro}` — integration tests per distro against bootstrap service containers
9. (amd64 only) `selenium` + `test-ui-desktop` + `test-ui-mobile` — Selenium UI tests
10. `test-upgrade` — upgrade path testing
11. `upload` / `promote` — publish to release repo (stable/master branches only)
12. `artifact` — upload test artifacts via SCP

# Visual Diff

Tools in `visual-diff/` for comparing UI test screenshots between branches.

## Download artifacts

Download screenshots and videos from a branch's latest CI build:
```
visual-diff/download-artifacts.sh [branch]
```
Defaults to current git branch. Downloads desktop and mobile screenshots + video to `visual-diff/output/{branch}/`. On Termux, also copies to `Pictures/syncloud-{branch}/` and `Movies/syncloud-{branch}/`.

## Screenshot diff

Compare screenshots between two branches:
```
visual-diff/screenshot-diff.sh [base-branch] [compare-branch] [-f filter]
```
Defaults: base=master, compare=current branch. Downloads screenshots from both branches and compares them pixel-by-pixel.

- Uses ImageMagick (`magick compare` or `compare`) with 5% fuzz tolerance when available, falls back to binary `cmp`
- Filter with `-f` to compare only screenshots matching a substring (e.g. `-f settings_access`)
- Output: `visual-diff/output/diff/` contains diff images with changed pixels highlighted in red
- On Termux, copies base/compare/diff images to `Pictures/screenshot-diff/`

# Running Drone builds locally

Generate `.drone.yml` from jsonnet (run from project root):
```
drone jsonnet --stdout --stream > .drone.yml
```

Run a specific pipeline with selected steps:
```
drone exec --pipeline amd64 --trusted \
  --include version \
  --include nginx \
  --include authelia \
  --include "build web" \
  --include build \
  --include package \
  --include "test bookworm" \
  .drone.yml
```
