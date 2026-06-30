# Changelog

All notable changes to this project are documented in this file.

## [2.0.0] - 2026-06-29

### Features

- Upload job summary when supported ([#917](https://gitea.com/gitea/runner/pulls/917))
- Complete runner-side cancellation handling ([#1016](https://gitea.com/gitea/runner/pulls/1016))
- Ipv6 options for network container creation ([#1029](https://gitea.com/gitea/runner/pulls/1029))
- Support `ssh://` action URLs ([#1035](https://gitea.com/gitea/runner/pulls/1035))
- Add optional `runner.post_task_script` hook after task cleanup ([#1026](https://gitea.com/gitea/runner/pulls/1026))
- Enable `jobs.<job_id>.timeout-minutes` and `jobs.<job_id>.continue-on-error` ([#1032](https://gitea.com/gitea/runner/pulls/1032))
- Shallow clone action repositories ([#1053](https://gitea.com/gitea/runner/pulls/1053))

### Dependencies

- **deps:** Update docker docker tag to v29.5.3 ([#1021](https://gitea.com/gitea/runner/pulls/1021))
- **deps:** Update golang.org/x/crypto to v0.52.0 ([#1027](https://gitea.com/gitea/runner/pulls/1027))
- **deps:** Update alpine docker tag to v3.24 ([#1030](https://gitea.com/gitea/runner/pulls/1030))
- **deps:** Update actions/checkout action to v7 ([#1042](https://gitea.com/gitea/runner/pulls/1042))
- **deps:** Update docker docker tag to v29.6.0 ([#1044](https://gitea.com/gitea/runner/pulls/1044))

### Bug fixes

- **deps:** Update module code.gitea.io/actions-proto-go to gitea.dev/actions-proto-go v0.5.0 ([#1009](https://gitea.com/gitea/runner/pulls/1009))
- **deps:** Update module golang.org/x/sys to v0.45.0 ([#1012](https://gitea.com/gitea/runner/pulls/1012))
- Do not update cached actions with stale origin URL ([#1014](https://gitea.com/gitea/runner/pulls/1014))
- **deps:** Update module github.com/opencontainers/selinux to v1.15.1 ([#1017](https://gitea.com/gitea/runner/pulls/1017))
- **Breaking:** Stop implicitly using `DOCKER_USERNAME`/`DOCKER_PASSWORD` secrets for image pulls ([#1007](https://gitea.com/gitea/runner/pulls/1007))
- **deps:** Update module github.com/docker/cli to v29.5.3+incompatible ([#1018](https://gitea.com/gitea/runner/pulls/1018))
- **deps:** Update module golang.org/x/sys to v0.46.0 ([#1019](https://gitea.com/gitea/runner/pulls/1019))
- **deps:** Update module golang.org/x/term to v0.44.0 ([#1031](https://gitea.com/gitea/runner/pulls/1031))
- **host:** Bound host-environment cleanup and reclaim leaked scratch dirs ([#1024](https://gitea.com/gitea/runner/pulls/1024))
- Prevent loss of step log output at end of step ([#1028](https://gitea.com/gitea/runner/pulls/1028))
- **cleanup:** Kill Unix step process group on cancel to avoid hang ([#1025](https://gitea.com/gitea/runner/pulls/1025))
- Interpolate job container.volume ([#1036](https://gitea.com/gitea/runner/pulls/1036))
- **deps:** Update module go.etcd.io/bbolt to v1.5.0 ([#1040](https://gitea.com/gitea/runner/pulls/1040))
- **deps:** Update module github.com/docker/cli to v29.6.0+incompatible ([#1045](https://gitea.com/gitea/runner/pulls/1045))
- **deps:** Update module github.com/moby/moby/client to v0.5.0 ([#1049](https://gitea.com/gitea/runner/pulls/1049))
- Guard SetJobError against missing job-error container ([#1050](https://gitea.com/gitea/runner/pulls/1050))
- Run always()/cancelled() and post steps correctly on cancellation ([#1043](https://gitea.com/gitea/runner/pulls/1043))
- Namespace local docker action image tags per repository ([#1051](https://gitea.com/gitea/runner/pulls/1051))
- Composite nested uses clone token ([#1041](https://gitea.com/gitea/runner/pulls/1041))

### Build

- Cover Windows Go files ([#1052](https://gitea.com/gitea/runner/pulls/1052))

### Documentation

- **docker-images:** Update docs ([#1020](https://gitea.com/gitea/runner/pulls/1020))
- Improve the documentation for cache ([#1034](https://gitea.com/gitea/runner/pulls/1034))
## [1.0.8] - 2026-06-02

### Bug fixes

- **cleanup:** Kill Windows step process tree on cancel to avoid hang ([#1011](https://gitea.com/gitea/runner/pulls/1011))
## [1.0.7] - 2026-05-29

### Bug fixes

- **reporter:** Respect configured log level for job log forwarding ([#989](https://gitea.com/gitea/runner/pulls/989))
- Deliver cancel ack and reap leftover Windows job processes ([#996](https://gitea.com/gitea/runner/pulls/996))
- Matrix-job data races + outputs, leaner offline test suite ([#994](https://gitea.com/gitea/runner/pulls/994))
- Restore global docker config dir and socket env in tests ([#1004](https://gitea.com/gitea/runner/pulls/1004))
- Support multiline secret masking ([#1001](https://gitea.com/gitea/runner/pulls/1001))
- **container:** Re-validate cached container id before reuse ([#1003](https://gitea.com/gitea/runner/pulls/1003))
## [1.0.6] - 2026-05-22

### Dependencies

- **deps:** Update actions/setup-node action to v6 ([#991](https://gitea.com/gitea/runner/pulls/991))

### Chores

- Pin Docker base images to explicit versions ([#992](https://gitea.com/gitea/runner/pulls/992))

### Bug fixes

- Clean up job network and container when container start fails ([#986](https://gitea.com/gitea/runner/pulls/986))
- **deps:** Update module github.com/opencontainers/selinux to v1.15.0 ([#990](https://gitea.com/gitea/runner/pulls/990))

### CI

- Add PR title linting against Conventional Commits ([#988](https://gitea.com/gitea/runner/pulls/988))

### Tests

- Make TestRunEvent integration suite runnable locally ([#987](https://gitea.com/gitea/runner/pulls/987))
## [1.0.5] - 2026-05-21

### Features

- Add cache.offline_mode to reuse cached actions ([#966](https://gitea.com/gitea/runner/pulls/966))

### Bug fixes

- **deps:** Update module github.com/go-git/go-git/v5 to v5.19.1 ([#980](https://gitea.com/gitea/runner/pulls/980))
- **deps:** Update module github.com/docker/cli to v29.5.1+incompatible ([#979](https://gitea.com/gitea/runner/pulls/979))
- **deps:** Update module connectrpc.com/connect to v1.20.0 ([#985](https://gitea.com/gitea/runner/pulls/985))
## [1.0.4] - 2026-05-18

### Features

- Make pseudo-TTY allocation opt-in ([#961](https://gitea.com/gitea/runner/pulls/961))

### Dependencies

- **deps:** Bump `retry-go`, `golangci-lint`, `govulncheck` ([#965](https://gitea.com/gitea/runner/pulls/965))
- **deps:** Update workflow dependencies (major) ([#967](https://gitea.com/gitea/runner/pulls/967))
- **deps:** Update workflow dependencies (major) ([#968](https://gitea.com/gitea/runner/pulls/968))

### Bug fixes

- Respect proxy env vars in runner client ([#962](https://gitea.com/gitea/runner/pulls/962))
- **deps:** Bump `docker` deps, switch to `moby/moby` ([#943](https://gitea.com/gitea/runner/pulls/943))
- **deps:** Update module github.com/docker/cli to v29.5.0+incompatible ([#969](https://gitea.com/gitea/runner/pulls/969))
- Return if executors length is zero in ParallelExecutor ([#960](https://gitea.com/gitea/runner/pulls/960))
- **parse_env_file:** Support env-file lines larger than 64 KiB ([#974](https://gitea.com/gitea/runner/pulls/974))
## [1.0.3] - 2026-05-12

### Bug fixes

- **deps:** Update module golang.org/x/term to v0.43.0 ([#948](https://gitea.com/gitea/runner/pulls/948))
- **artifactcache:** Preserve cache key case to stop redundant uploads ([#947](https://gitea.com/gitea/runner/pulls/947))
- Ensure `dbfs_data` is cleaned up after task completion ([#952](https://gitea.com/gitea/runner/pulls/952))
- **deps:** Update mergo to v1.0.2 (now `dario.cat/mergo`) ([#954](https://gitea.com/gitea/runner/pulls/954))
- **deps:** Update module github.com/opencontainers/selinux to v1.14.1 ([#955](https://gitea.com/gitea/runner/pulls/955))
- Isolate per-task runner envs ([#959](https://gitea.com/gitea/runner/pulls/959))
## [1.0.2] - 2026-05-08

### Features

- Remove emojis from runner logging, add `Starting job container` group ([#940](https://gitea.com/gitea/runner/pulls/940))

### Bug fixes

- Serialize action-cache reads to prevent worktree race ([#938](https://gitea.com/gitea/runner/pulls/938))
- Overwrite read-only files when copying action directories ([#942](https://gitea.com/gitea/runner/pulls/942))
## [1.0.1] - 2026-05-07

### Dependencies

- **deps:** Update docker docker tag to v29 ([#924](https://gitea.com/gitea/runner/pulls/924))

### Bug fixes

- Re-fetch cached reusable workflow on every run ([#930](https://gitea.com/gitea/runner/pulls/930))
- **deps:** Update module github.com/opencontainers/selinux to v1.14.0 ([#928](https://gitea.com/gitea/runner/pulls/928))
- **deps:** Update module github.com/go-git/go-git/v5 to v5.19.0 ([#934](https://gitea.com/gitea/runner/pulls/934))

### Tests

- Clean up dead/stale fixtures and bump test container images ([#932](https://gitea.com/gitea/runner/pulls/932))
## [1.0.0] - 2026-05-05

### Features

- Add startup janitor for stale bind-workdir task workspaces ([#870](https://gitea.com/gitea/runner/pulls/870))

### Dependencies

- **deps:** Update dependency @actions/core to v1.11.1 ([#880](https://gitea.com/gitea/runner/pulls/880))
- **deps:** Update dependency @vercel/ncc to ^0.38.0 ([#881](https://gitea.com/gitea/runner/pulls/881))
- **deps:** Update actions/hello-world-docker-action action to v2 ([#921](https://gitea.com/gitea/runner/pulls/921))
- **deps:** Update crazy-max/ghaction-import-gpg action to v7 ([#923](https://gitea.com/gitea/runner/pulls/923))

### Bug fixes

- **deps:** Update module github.com/rhysd/actionlint to v1.7.12 ([#873](https://gitea.com/gitea/runner/pulls/873))
- **deps:** Update module golang.org/x/term to v0.42.0 ([#920](https://gitea.com/gitea/runner/pulls/920))
- **host:** Correct host workspace cleanup on Windows ([#883](https://gitea.com/gitea/runner/pulls/883))

### CI

- Fix triggers ([#882](https://gitea.com/gitea/runner/pulls/882))

### Documentation

- Update docs ([#872](https://gitea.com/gitea/runner/pulls/872))
## [0.6.1] - 2026-04-30

### Chores

- Clean up `nolint` directives in `act` package ([#864](https://gitea.com/gitea/runner/pulls/864))

### Bug fixes

- **deps:** Update module github.com/docker/docker to v25.0.15+incompatible ([#867](https://gitea.com/gitea/runner/pulls/867))
- **deps:** Update module github.com/moby/patternmatcher to v0.6.1 ([#868](https://gitea.com/gitea/runner/pulls/868))

### Performance

- **config:** Lower default fetch_interval_max from 60s to 5s ([#875](https://gitea.com/gitea/runner/pulls/875))
## [0.6.0] - 2026-04-28

### Features

- Show run command, shell and env in collapsible group before step output ([#847](https://gitea.com/gitea/runner/pulls/847))

### Dependencies

- **deps:** Drop unused `distribution/reference` replace directive ([#858](https://gitea.com/gitea/runner/pulls/858))

### Bug fixes

- **deps:** Update module github.com/docker/docker to v25.0.14+incompatible ([#862](https://gitea.com/gitea/runner/pulls/862))
- **deps:** Update module github.com/mattn/go-isatty to v0.0.22 ([#863](https://gitea.com/gitea/runner/pulls/863))
- Avoid 'filename too long' on matrix jobs by hashing container names ([#853](https://gitea.com/gitea/runner/pulls/853))
## [0.5.0] - 2026-04-27

### Features

- Add support for Dynamic Matrix Strategies with Job Outputs ([#149](https://gitea.com/gitea/runner/pulls/149))

### Bug fixes

- Fetch when other refs get force-pushed ([#846](https://gitea.com/gitea/runner/pulls/846))
- Heartbeat `ReportState` for long-running silent jobs ([#852](https://gitea.com/gitea/runner/pulls/852))
- Wait for docker supervise dir before `s6-svwait` ([#851](https://gitea.com/gitea/runner/pulls/851))
- **deps:** Update module connectrpc.com/connect to v1.19.2 ([#854](https://gitea.com/gitea/runner/pulls/854))
- **deps:** Update module github.com/docker/cli to v25.0.7+incompatible ([#855](https://gitea.com/gitea/runner/pulls/855))
## [0.261.8] - 2026-02-18

### Features

- Support runs-on labels and group ([#2062](https://gitea.com/gitea/runner/pulls/2062))
- Support config env expansion ([#2063](https://gitea.com/gitea/runner/pulls/2063))
- Cli option to enable the new action cache ([#1954](https://gitea.com/gitea/runner/pulls/1954))
- Support offline mode ([#2128](https://gitea.com/gitea/runner/pulls/2128))
- Correctly use the xdg library, which has the side effect to fix the config survey ([#2195](https://gitea.com/gitea/runner/pulls/2195))
- Offline mode for new action cache ([#2173](https://gitea.com/gitea/runner/pulls/2173))
- Test writing remote Actions locally ([#2155](https://gitea.com/gitea/runner/pulls/2155))

### Chores

- Bump VERSION to 0.2.52
- Bump VERSION to 0.2.53
- Bump VERSION to 0.2.54
- Bump VERSION to 0.2.55
- Bump VERSION to 0.2.56
- Bump VERSION to 0.2.57
- Bump VERSION to 0.2.58
- Bump VERSION to 0.2.59
- Bump VERSION to 0.2.60
- Remove repetitive words ([#2259](https://gitea.com/gitea/runner/pulls/2259))
- Bump VERSION to 0.2.61

### Other

- Reduce log spam. ([#2256](https://gitea.com/gitea/runner/pulls/2256))

### Bug fixes

- Use actions/runner hashfiles in container ([#1940](https://gitea.com/gitea/runner/pulls/1940))
- Panic
- Panic ([#2071](https://gitea.com/gitea/runner/pulls/2071))
- ([#2075](https://gitea.com/gitea/runner/pulls/2075))
- IsHost is defined as false on windows ([#2093](https://gitea.com/gitea/runner/pulls/2093))
- Match cache `restore-keys` in creation reverse order ([#2153](https://gitea.com/gitea/runner/pulls/2153))
- Write default config in XDG config dir to avoid cluttering the HOME directory by default ([#2140](https://gitea.com/gitea/runner/pulls/2140))
- Improve action not found error ([#2171](https://gitea.com/gitea/runner/pulls/2171))
- Subpath actions via new artifact cache ([#2170](https://gitea.com/gitea/runner/pulls/2170))
- Improve new-action-cache fetch failure error ([#2172](https://gitea.com/gitea/runner/pulls/2172))
- Improve warning about remote not found ([#2169](https://gitea.com/gitea/runner/pulls/2169))
- Use correct path to toolcache ([#1494](https://gitea.com/gitea/runner/pulls/1494))
- Rootless error concerning `/var/run/docker.sock` ([#2181](https://gitea.com/gitea/runner/pulls/2181))
- Logo and demo asset urls ([#2226](https://gitea.com/gitea/runner/pulls/2226))
- Docker buildx cache restore not working ([#2236](https://gitea.com/gitea/runner/pulls/2236))
- Rootless permission bits (new actions cache) ([#2242](https://gitea.com/gitea/runner/pulls/2242))
- Cache adjust restore order of exact key matches ([#2267](https://gitea.com/gitea/runner/pulls/2267))
- Max parallel Support for Matrix Jobs and Remote Action Tests ([#150](https://gitea.com/gitea/runner/pulls/150))

### Refactoring

- **cmd/root:** Simplify `parseEnvs` ([#2162](https://gitea.com/gitea/runner/pulls/2162))
- Filecollector into new package ([#2174](https://gitea.com/gitea/runner/pulls/2174))
- Simpilfy go-git cache ([#2208](https://gitea.com/gitea/runner/pulls/2208))
- Simplify adding new node versions add node 24 ([#140](https://gitea.com/gitea/runner/pulls/140))

### Build

- **deps:** Bump megalinter/megalinter from 7.3.0 to 7.4.0 ([#2025](https://gitea.com/gitea/runner/pulls/2025))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.25 to 1.6.26 ([#2026](https://gitea.com/gitea/runner/pulls/2026))
- **deps:** Bump golang.org/x/term from 0.12.0 to 0.13.0 ([#2036](https://gitea.com/gitea/runner/pulls/2036))
- **deps:** Bump github.com/docker/distribution ([#2037](https://gitea.com/gitea/runner/pulls/2037))
- **deps:** Bump golang.org/x/net from 0.15.0 to 0.17.0 ([#2045](https://gitea.com/gitea/runner/pulls/2045))
- **deps:** Bump github.com/moby/buildkit from 0.12.2 to 0.12.3 ([#2060](https://gitea.com/gitea/runner/pulls/2060))
- **deps:** Bump github.com/mattn/go-isatty from 0.0.19 to 0.0.20 ([#2059](https://gitea.com/gitea/runner/pulls/2059))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.9.0 to 5.10.0 ([#2066](https://gitea.com/gitea/runner/pulls/2066))
- **deps:** Bump github.com/docker/cli ([#2069](https://gitea.com/gitea/runner/pulls/2069))
- **deps:** Bump go.etcd.io/bbolt from 1.3.7 to 1.3.8 ([#2065](https://gitea.com/gitea/runner/pulls/2065))
- **deps:** Bump megalinter/megalinter from 7.4.0 to 7.5.0 ([#2070](https://gitea.com/gitea/runner/pulls/2070))
- **deps:** Bump github.com/creack/pty from 1.1.18 to 1.1.20 ([#2068](https://gitea.com/gitea/runner/pulls/2068))
- **deps:** Bump github.com/docker/docker ([#2067](https://gitea.com/gitea/runner/pulls/2067))
- **deps:** Bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([#2079](https://gitea.com/gitea/runner/pulls/2079))
- **deps:** Bump golang.org/x/term from 0.13.0 to 0.14.0 ([#2091](https://gitea.com/gitea/runner/pulls/2091))
- **deps:** Bump actions/github-script from 6 to 7 ([#2097](https://gitea.com/gitea/runner/pulls/2097))
- **deps:** Bump megalinter/megalinter from 7.5.0 to 7.6.0 ([#2098](https://gitea.com/gitea/runner/pulls/2098))
- **deps:** Bump github.com/moby/buildkit from 0.12.3 to 0.12.4 ([#2113](https://gitea.com/gitea/runner/pulls/2113))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.10.0 to 5.10.1 ([#2114](https://gitea.com/gitea/runner/pulls/2114))
- **deps:** Bump golang.org/x/term from 0.14.0 to 0.15.0 ([#2112](https://gitea.com/gitea/runner/pulls/2112))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.10.1 to 5.11.0 ([#2117](https://gitea.com/gitea/runner/pulls/2117))
- **deps:** Bump megalinter/megalinter from 7.6.0 to 7.7.0 ([#2119](https://gitea.com/gitea/runner/pulls/2119))
- **deps:** Bump actions/setup-go from 4 to 5 ([#2118](https://gitea.com/gitea/runner/pulls/2118))
- **deps:** Bump github.com/creack/pty from 1.1.20 to 1.1.21 ([#2099](https://gitea.com/gitea/runner/pulls/2099))
- **deps:** Bump actions/checkout from 3 to 4 ([#1998](https://gitea.com/gitea/runner/pulls/1998))
- **deps:** Bump actions/stale from 8 to 9 ([#2120](https://gitea.com/gitea/runner/pulls/2120))
- **deps:** Bump github.com/containerd/containerd from 1.7.2 to 1.7.11 ([#2136](https://gitea.com/gitea/runner/pulls/2136))
- **deps:** Bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#2134](https://gitea.com/gitea/runner/pulls/2134))
- **deps:** Bump golang.org/x/term from 0.15.0 to 0.16.0 ([#2148](https://gitea.com/gitea/runner/pulls/2148))
- **deps:** Bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 ([#2149](https://gitea.com/gitea/runner/pulls/2149))
- **deps:** Bump actions/upload-artifact from 3 to 4 ([#2133](https://gitea.com/gitea/runner/pulls/2133))
- **deps:** Bump megalinter/megalinter from 7.7.0 to 7.8.0 ([#2164](https://gitea.com/gitea/runner/pulls/2164))
- **deps:** Bump codecov/codecov-action from 3.1.4 to 3.1.5 ([#2175](https://gitea.com/gitea/runner/pulls/2175))
- **deps:** Bump github.com/moby/buildkit from 0.12.4 to 0.12.5 ([#2186](https://gitea.com/gitea/runner/pulls/2186))
- **deps:** Bump github.com/opencontainers/image-spec ([#2167](https://gitea.com/gitea/runner/pulls/2167))
- **deps:** Bump github.com/containerd/containerd from 1.6.19 to 1.6.26 ([#2189](https://gitea.com/gitea/runner/pulls/2189))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.12.5 ([#2188](https://gitea.com/gitea/runner/pulls/2188))
- **deps:** Bump github.com/opencontainers/runc from 1.1.7 to 1.1.12 ([#2187](https://gitea.com/gitea/runner/pulls/2187))
- **deps:** Bump golang.org/x/term from 0.16.0 to 0.17.0 ([#2213](https://gitea.com/gitea/runner/pulls/2213))
- **deps:** Bump github.com/opencontainers/image-spec ([#2214](https://gitea.com/gitea/runner/pulls/2214))
- **deps:** Bump github.com/opencontainers/image-spec ([#2223](https://gitea.com/gitea/runner/pulls/2223))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.26 to 1.6.27 ([#2228](https://gitea.com/gitea/runner/pulls/2228))
- **deps:** Bump go.etcd.io/bbolt from 1.3.8 to 1.3.9 ([#2229](https://gitea.com/gitea/runner/pulls/2229))
- **deps:** Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#2235](https://gitea.com/gitea/runner/pulls/2235))
- **deps:** Bump golang.org/x/term from 0.17.0 to 0.18.0 ([#2244](https://gitea.com/gitea/runner/pulls/2244))
- **deps:** Bump github.com/docker/docker ([#2252](https://gitea.com/gitea/runner/pulls/2252))

### CI

- Automatic merge PRs created by a maintainer and approved by 1 other maintainer ([#2156](https://gitea.com/gitea/runner/pulls/2156))

### Tests

- Use ping to improve network test ([#2266](https://gitea.com/gitea/runner/pulls/2266))
## [0.2.51] - 2023-09-24

### Features

- Support branches as GITHUB_REF
- Load every environment from --env-file to workflow ([#184](https://gitea.com/gitea/runner/pulls/184))
- Docker 19.03.12 ([#347](https://gitea.com/gitea/runner/pulls/347))
- Add fromJSON support ([#352](https://gitea.com/gitea/runner/pulls/352))
- Support custom GITHUB_RUN_ID, GITHUB_RUN_NUMBER ([#369](https://gitea.com/gitea/runner/pulls/369))
- Speeding up workflows ([#643](https://gitea.com/gitea/runner/pulls/643))
- Don't push `pkg/runner/act/` where actions cache is stored ([#620](https://gitea.com/gitea/runner/pulls/620))
- Add `go` install instruction ([#672](https://gitea.com/gitea/runner/pulls/672))
- Add option for custom socket path ([#698](https://gitea.com/gitea/runner/pulls/698))
- Improve list ([#786](https://gitea.com/gitea/runner/pulls/786))
- Add Dependabot config ([#797](https://gitea.com/gitea/runner/pulls/797))
- Allow existing logger from context ([#859](https://gitea.com/gitea/runner/pulls/859))
- Add option for docker image rebuild ([#878](https://gitea.com/gitea/runner/pulls/878))
- SELinux lables for --bind ([#906](https://gitea.com/gitea/runner/pulls/906))
- Allow existing logger from context ([#898](https://gitea.com/gitea/runner/pulls/898))
- Add node16 support ([#922](https://gitea.com/gitea/runner/pulls/922))
- **goreleaser:** Format changelog, add prereleases ([#942](https://gitea.com/gitea/runner/pulls/942))
- Try to read ref and sha from event payload if available ([#889](https://gitea.com/gitea/runner/pulls/889))
- Add skipped status as step result ([#950](https://gitea.com/gitea/runner/pulls/950))
- Run jobs in parallel ([#1003](https://gitea.com/gitea/runner/pulls/1003))
- Add json logger output ([#1026](https://gitea.com/gitea/runner/pulls/1026))
- Add step name to logger field ([#1027](https://gitea.com/gitea/runner/pulls/1027))
- Add flag to always run the checkout action ([#1049](https://gitea.com/gitea/runner/pulls/1049))
- Add bug-report flag ([#1056](https://gitea.com/gitea/runner/pulls/1056))
- Split job steps into its own files/structs ([#1004](https://gitea.com/gitea/runner/pulls/1004))
- Add option to specify git remote name ([#1104](https://gitea.com/gitea/runner/pulls/1104))
- Support `GITHUB_REF_NAME` & `GITHUB_REF_TYPE` ([#1142](https://gitea.com/gitea/runner/pulls/1142))
- Add ubuntu-22.04 ([#1150](https://gitea.com/gitea/runner/pulls/1150))
- Handle context cancelation during docker exec ([#1170](https://gitea.com/gitea/runner/pulls/1170))
- Non root user container ([#1202](https://gitea.com/gitea/runner/pulls/1202))
- Add option to bypass GHE for actions checkout ([#1162](https://gitea.com/gitea/runner/pulls/1162))
- Pass current matrix information to job logger ([#1364](https://gitea.com/gitea/runner/pulls/1364))
- Add logger hook for standalone job logger
- Support more options of containers
- Support PlatformPicker
- Adjust container name format ([#1](https://gitea.com/gitea/runner/pulls/1))
- Parse types of reusable workflows ([#1414](https://gitea.com/gitea/runner/pulls/1414))
- Host environment ([#1293](https://gitea.com/gitea/runner/pulls/1293))
- GITHUB_STATE and GITHUB_OUTPUT file commands ([#1391](https://gitea.com/gitea/runner/pulls/1391))
- Interpolate the step names ([#1422](https://gitea.com/gitea/runner/pulls/1422))
- **workflow:** Support schedule event ([#4](https://gitea.com/gitea/runner/pulls/4))
- `--container-options` ([#1462](https://gitea.com/gitea/runner/pulls/1462))
- JobLoggerFactory ([#1496](https://gitea.com/gitea/runner/pulls/1496))
- Allow to spawn and run a local reusable workflow ([#1423](https://gitea.com/gitea/runner/pulls/1423))
- Support "result" on "needs" context. ([#1497](https://gitea.com/gitea/runner/pulls/1497))
- Allow building without docker support ([#1507](https://gitea.com/gitea/runner/pulls/1507))
- Add support for building docker actions with private registries ([#1557](https://gitea.com/gitea/runner/pulls/1557))
- Add check for newer versions ([#1562](https://gitea.com/gitea/runner/pulls/1562))
- Add remote reusable workflows ([#1525](https://gitea.com/gitea/runner/pulls/1525))
- Release extension
- Release extension
- Release extension
- Step summary of test results ([#1580](https://gitea.com/gitea/runner/pulls/1580))
- Cache notices to reduce frequency of upgrade notifications ([#1592](https://gitea.com/gitea/runner/pulls/1592))
- Allow overriding of `GITHUB_` env variables ([#1582](https://gitea.com/gitea/runner/pulls/1582))
- Workflowpattern package ([#1618](https://gitea.com/gitea/runner/pulls/1618))
- Specify matrix on command line ([#1675](https://gitea.com/gitea/runner/pulls/1675))
- Improve GetOutboundIP ([#1707](https://gitea.com/gitea/runner/pulls/1707))
- Support yaml env/secrets/inputs file ([#1733](https://gitea.com/gitea/runner/pulls/1733))
- Implement steps.timeout-minutes ([#1776](https://gitea.com/gitea/runner/pulls/1776))
- Add support for github action variables ([#1833](https://gitea.com/gitea/runner/pulls/1833))
- Add CopyTarStream to Container Interface ([#1912](https://gitea.com/gitea/runner/pulls/1912))
- Cmd support for windows ([#1941](https://gitea.com/gitea/runner/pulls/1941))
- Add new Action Cache ([#1913](https://gitea.com/gitea/runner/pulls/1913))
- Support interpolation in `<job>.container.options` ([#1958](https://gitea.com/gitea/runner/pulls/1958))
- Support node20 runtime ([#1988](https://gitea.com/gitea/runner/pulls/1988))

### Dependencies

- **deps:** Cleanup go.mod a little bit ([#136](https://gitea.com/gitea/runner/pulls/136))

### Chores

- Upgrade docker v19.03.12 ([#348](https://gitea.com/gitea/runner/pulls/348))
- Fold integration test outputs ([#1111](https://gitea.com/gitea/runner/pulls/1111))
- Update LICENSE
- Push new version
- Bump version
- Makefile update
- Bump version
- Fix promote script
- Bump VERSION to 0.2.39
- Bump VERSION to 0.2.40
- Update docs for installing act as GH CLI extension
- Bump VERSION to 0.2.41
- Fix release script to trigger gh-act
- Bump VERSION to 0.2.42
- Add Winget installation instructions ([#1622](https://gitea.com/gitea/runner/pulls/1622))
- Use new style octal ([#1630](https://gitea.com/gitea/runner/pulls/1630))
- Remove obsolete Container.UpdateFromPath ([#1631](https://gitea.com/gitea/runner/pulls/1631))
- Bump VERSION to 0.2.43
- **yaml:** Improve `ParseRawOn` ([#28](https://gitea.com/gitea/runner/pulls/28))
- Update go-git dependency in go.mod ([#30](https://gitea.com/gitea/runner/pulls/30))
- Bump VERSION to 0.2.44
- Run act from cli on linux ([#1758](https://gitea.com/gitea/runner/pulls/1758))
- Bump VERSION to 0.2.45
- Bump VERSION to 0.2.46
- Bump VERSION to 0.2.47
- Bump VERSION to 0.2.48
- Upgrade golangci-lint and address findings ([#1904](https://gitea.com/gitea/runner/pulls/1904))
- Bump VERSION to 0.2.49
- Bump VERSION to 0.2.50
- Bump VERSION to 0.2.51

### Reverts

- Deprecation of containerArchitecture ([#1514](https://gitea.com/gitea/runner/pulls/1514))

### Other

- Update go-git import ([#194](https://gitea.com/gitea/runner/pulls/194))
- Add snapshot artifacts to PR builds ([#574](https://gitea.com/gitea/runner/pulls/574))
- Configuration update ([#659](https://gitea.com/gitea/runner/pulls/659))
- Configuration update ([#661](https://gitea.com/gitea/runner/pulls/661))
- Configuration update ([#662](https://gitea.com/gitea/runner/pulls/662))
- Configuration update ([#663](https://gitea.com/gitea/runner/pulls/663))
- Skip Checkout Regression ([#680](https://gitea.com/gitea/runner/pulls/680))
- GITHUB_PATH should prepend ([#690](https://gitea.com/gitea/runner/pulls/690))
- Add MacPorts info ([#762](https://gitea.com/gitea/runner/pulls/762))
- Switch to latest and drop Homebrew config ([#848](https://gitea.com/gitea/runner/pulls/848))
- **megalinter:** Disable markdown link check ([#938](https://gitea.com/gitea/runner/pulls/938))
- Uses in composite ([#793](https://gitea.com/gitea/runner/pulls/793))
- Env.ACT in if condition of the workflow ([#965](https://gitea.com/gitea/runner/pulls/965))
- Regression run after failure ([#971](https://gitea.com/gitea/runner/pulls/971))
- Don't stop JobContainer on failure ([#1041](https://gitea.com/gitea/runner/pulls/1041))
- Tojson(null) == 'null' ([#1068](https://gitea.com/gitea/runner/pulls/1068))
- Call of reflect.Value.Interface on zero Value ([#1067](https://gitea.com/gitea/runner/pulls/1067))
- Panic: reflect: slice index out of range ([#1066](https://gitea.com/gitea/runner/pulls/1066))
- Bump go-git ([#1057](https://gitea.com/gitea/runner/pulls/1057))
- Remove shlex ([#1086](https://gitea.com/gitea/runner/pulls/1086))
- Bump go.mod to go1.18 ([#1088](https://gitea.com/gitea/runner/pulls/1088))
- Add MacPorts badge ([#1129](https://gitea.com/gitea/runner/pulls/1129))
- New issue templates ([#1048](https://gitea.com/gitea/runner/pulls/1048))
- Fix bug report rendering ([#1186](https://gitea.com/gitea/runner/pulls/1186))
- Add support to replace GitHub's env ([#1197](https://gitea.com/gitea/runner/pulls/1197))
- Add ReadHeaderTimeout ([#1277](https://gitea.com/gitea/runner/pulls/1277))
- Fix shell installation ([#1547](https://gitea.com/gitea/runner/pulls/1547))
- Make stalebot less annoying ([#1587](https://gitea.com/gitea/runner/pulls/1587))
- Fix expression of warning message on macOS ([#1693](https://gitea.com/gitea/runner/pulls/1693))

### Bug fixes

- Panic on regex in main.workflow
- Added tests and improved code
- Using filepath instead of path
- Clean code
- 168 include error message from docker SDK
- #182 - support symlinks in workspace
- Ignore symlink errors
- #170 #171 - add RUNNER_OS and RUNNER_TEMP
- #170 add RUNNER_TEMP
- Add default values in vmInput
- **test:** Adapt to sjw change
- Outputs espace ([#404](https://gitea.com/gitea/runner/pulls/404))
- --version, --help do not need check platform config ([#508](https://gitea.com/gitea/runner/pulls/508))
- Add `ImageOS` env var based on running platform ([#571](https://gitea.com/gitea/runner/pulls/571))
- Invoke login shells to source `/etc/environment` ([#570](https://gitea.com/gitea/runner/pulls/570))
- Pass `sc.Env` through `common.Logger` to hide secrets ([#578](https://gitea.com/gitea/runner/pulls/578))
- Hardcode `ubuntu-latest` for `ImageOS` envvar ([#579](https://gitea.com/gitea/runner/pulls/579))
- Update `ubuntu-latest` to point to `ubuntu-20.04` ([#580](https://gitea.com/gitea/runner/pulls/580))
- Use new platforms after survey, check working dir for `.actrc` ([#577](https://gitea.com/gitea/runner/pulls/577))
- Return error if both `run:`/`uses:` keys are used ([#593](https://gitea.com/gitea/runner/pulls/593))
- Remove HOME env var ([#594](https://gitea.com/gitea/runner/pulls/594))
- Fail workflow if the job name is invalid ([#596](https://gitea.com/gitea/runner/pulls/596))
- Reworked container architecture ([#619](https://gitea.com/gitea/runner/pulls/619))
- Update readme with more relevant/updated information ([#630](https://gitea.com/gitea/runner/pulls/630))
- Ignore `act` binary built when using `make` ([#654](https://gitea.com/gitea/runner/pulls/654))
- Replace `gopass` lib with `golang.org/x/term` ([#674](https://gitea.com/gitea/runner/pulls/674))
- Remove Windows issue from `README.md` ([#675](https://gitea.com/gitea/runner/pulls/675))
- Hide masked value ([#668](https://gitea.com/gitea/runner/pulls/668))
- Use auth when fetching (same as with cloning) ([#687](https://gitea.com/gitea/runner/pulls/687))
- Use github auth token for pulling a repository ([#693](https://gitea.com/gitea/runner/pulls/693))
- Avoid race condition by not running exec/start twice ([#702](https://gitea.com/gitea/runner/pulls/702))
- Update issue template ([#722](https://gitea.com/gitea/runner/pulls/722))
- Don't replace empty secrets ([#725](https://gitea.com/gitea/runner/pulls/725))
- Update codecov action ([#727](https://gitea.com/gitea/runner/pulls/727))
- Update codecov action (again) ([#731](https://gitea.com/gitea/runner/pulls/731))
- Wrong perms ([#726](https://gitea.com/gitea/runner/pulls/726))
- Warn Apple M1 (arm64) users ([#743](https://gitea.com/gitea/runner/pulls/743))
- Copy whole action dir to container ([#774](https://gitea.com/gitea/runner/pulls/774))
- Mount volume for workflow info ([#776](https://gitea.com/gitea/runner/pulls/776))
- **mergify:** Needs-work label handling ([#810](https://gitea.com/gitea/runner/pulls/810))
- Add simple concurrency limit ([#823](https://gitea.com/gitea/runner/pulls/823))
- Rewrite how image env is merged ([#828](https://gitea.com/gitea/runner/pulls/828))
- Rewrite StepContext env ([#832](https://gitea.com/gitea/runner/pulls/832))
- Set `sh` as default shell for containers  ([#853](https://gitea.com/gitea/runner/pulls/853))
- Interpolate env for step actions ([#872](https://gitea.com/gitea/runner/pulls/872))
- Cleanup unused actions ([#873](https://gitea.com/gitea/runner/pulls/873))
- Use docker lib for image ref parsing ([#877](https://gitea.com/gitea/runner/pulls/877))
- Correctly unmarshal entrypoint ([#876](https://gitea.com/gitea/runner/pulls/876))
- FromJSON should work with any JSON ([#883](https://gitea.com/gitea/runner/pulls/883))
- Continue jobs + steps after failure ([#840](https://gitea.com/gitea/runner/pulls/840))
- ::stop-commands:: should continue to print the lines ([#920](https://gitea.com/gitea/runner/pulls/920))
- Clarify --reuse and --rm flags ([#929](https://gitea.com/gitea/runner/pulls/929))
- Rework setupShellCommand ([#930](https://gitea.com/gitea/runner/pulls/930))
- Update install.sh ([#937](https://gitea.com/gitea/runner/pulls/937))
- **goreleaser:** Add append mode for release notes ([#962](https://gitea.com/gitea/runner/pulls/962))
- Set composite outputs on failure ([#945](https://gitea.com/gitea/runner/pulls/945))
- Docker cp of dangling symlink ([#943](https://gitea.com/gitea/runner/pulls/943)) ([#948](https://gitea.com/gitea/runner/pulls/948))
- Always execute closeContainer() executor ([#988](https://gitea.com/gitea/runner/pulls/988))
- Deep evaluate matrix strategy ([#964](https://gitea.com/gitea/runner/pulls/964))
- Print line after password prompt ([#1014](https://gitea.com/gitea/runner/pulls/1014))
- Hide masked values in log files ([#1011](https://gitea.com/gitea/runner/pulls/1011))
- Trim whitespace from actrc ([#1024](https://gitea.com/gitea/runner/pulls/1024))
- Coerce booleans to numbers for comparison in exprparser ([#1030](https://gitea.com/gitea/runner/pulls/1030))
- Remove unused TestContext ([#1013](https://gitea.com/gitea/runner/pulls/1013))
- Set default version string ([#1070](https://gitea.com/gitea/runner/pulls/1070))
- Don't overwrite with empty cmd/entrypoint ([#1076](https://gitea.com/gitea/runner/pulls/1076))
- **editorconfig:** Correct braces ([#1094](https://gitea.com/gitea/runner/pulls/1094))
- Return error on reusable workflows ([#1093](https://gitea.com/gitea/runner/pulls/1093))
- Copy ignored tracked files ([#1074](https://gitea.com/gitea/runner/pulls/1074))
- Do not error that much ([#1137](https://gitea.com/gitea/runner/pulls/1137))
- Conclusion and outcome are no integers ([#1136](https://gitea.com/gitea/runner/pulls/1136))
- Panic while running DRYRUN with local actions ([#1141](https://gitea.com/gitea/runner/pulls/1141))
- Return invalid step type ([#1157](https://gitea.com/gitea/runner/pulls/1157))
- Panic if a step in a job is nil ([#1145](https://gitea.com/gitea/runner/pulls/1145))
- Remove composite restrictions ([#1128](https://gitea.com/gitea/runner/pulls/1128))
- Remove volume with same name as the container again ([#1203](https://gitea.com/gitea/runner/pulls/1203))
- Dryrun skip local actions + enable Tests ([#1199](https://gitea.com/gitea/runner/pulls/1199))
- Skip local actions pre step in any case ([#1204](https://gitea.com/gitea/runner/pulls/1204))
- Allow TimeoutMinutes to be expression ([#1217](https://gitea.com/gitea/runner/pulls/1217))
- Keep action inputs ([#1215](https://gitea.com/gitea/runner/pulls/1215))
- Localcheckout mock ([#1198](https://gitea.com/gitea/runner/pulls/1198))
- Processing of strategy.matrix.include ([#1200](https://gitea.com/gitea/runner/pulls/1200))
- Checkout container path should always use `/` ([#1226](https://gitea.com/gitea/runner/pulls/1226))
- Use docker images from dockerhub ([#1249](https://gitea.com/gitea/runner/pulls/1249))
- Allow TimeoutMinutes to be expression in Jobs ([#1247](https://gitea.com/gitea/runner/pulls/1247))
- The number in the github event is of type number ([#1252](https://gitea.com/gitea/runner/pulls/1252))
- Add parent step id in composite action step id ([#1268](https://gitea.com/gitea/runner/pulls/1268))
- Ensure all post steps are executed ([#1286](https://gitea.com/gitea/runner/pulls/1286))
- CI is unstable due to new linting rules ([#1292](https://gitea.com/gitea/runner/pulls/1292))
- Show workflow name when get error on yml parsing ([#1133](https://gitea.com/gitea/runner/pulls/1133)) ([#1328](https://gitea.com/gitea/runner/pulls/1328))
- Support expression for step's continue-on-error field ([#900](https://gitea.com/gitea/runner/pulls/900)) ([#1331](https://gitea.com/gitea/runner/pulls/1331))
- Missing defaults in nodejs pre ([#1349](https://gitea.com/gitea/runner/pulls/1349))
- Crash on invalid docker credentials ([#1348](https://gitea.com/gitea/runner/pulls/1348))
- Align github.ref to GitHub Action ([#1362](https://gitea.com/gitea/runner/pulls/1362))
- Restore the job logger setup after job cancelation ([#1365](https://gitea.com/gitea/runner/pulls/1365))
- Show workflow info even if on.push is not defined ([#1329](https://gitea.com/gitea/runner/pulls/1329)) ([#1335](https://gitea.com/gitea/runner/pulls/1335))
- Handle go closure iteration ([#1374](https://gitea.com/gitea/runner/pulls/1374))
- Support docker create arguments from container.options ([#1022](https://gitea.com/gitea/runner/pulls/1022)) ([#1351](https://gitea.com/gitea/runner/pulls/1351))
- Interpolate action input defaults ([#1376](https://gitea.com/gitea/runner/pulls/1376))
- Re-evaluate env for remote composite actions ([#1385](https://gitea.com/gitea/runner/pulls/1385))
- Ci is failing since 2022-10-17 ([#1397](https://gitea.com/gitea/runner/pulls/1397))
- Set logger with trace level
- Nil pointer access ( workflow_dispatch )
- Keep path to event json file in composite actions ([#1428](https://gitea.com/gitea/runner/pulls/1428))
- Handle subdirectors and exclusions ([#1012](https://gitea.com/gitea/runner/pulls/1012)) ([#1446](https://gitea.com/gitea/runner/pulls/1446))
- **test:** Needs condition. ([#8](https://gitea.com/gitea/runner/pulls/8))
- Step env is unavailable in with property expr ([#1458](https://gitea.com/gitea/runner/pulls/1458))
- Handle env-vars case sensitive ([#1493](https://gitea.com/gitea/runner/pulls/1493))
- Tail (not absolute) as entrypoint of job container ([#1506](https://gitea.com/gitea/runner/pulls/1506))
- Align runner.os / runner.arch to known values ([#1510](https://gitea.com/gitea/runner/pulls/1510))
- Preserve job result state in case of failure ([#1519](https://gitea.com/gitea/runner/pulls/1519))
- Extra path lost in composite actions ([#1531](https://gitea.com/gitea/runner/pulls/1531))
- Add-matcher fails github workflow ([#1532](https://gitea.com/gitea/runner/pulls/1532))
- Allow override of artifact server bind address ([#1560](https://gitea.com/gitea/runner/pulls/1560))
- Update artifact server to address GHSL-2023-004 ([#1565](https://gitea.com/gitea/runner/pulls/1565))
- Apply forcePull only for prebuild docker actions ([#1599](https://gitea.com/gitea/runner/pulls/1599))
- Don't override env ([#1629](https://gitea.com/gitea/runner/pulls/1629))
- Tolerate workflow that needs a missing job ([#1595](https://gitea.com/gitea/runner/pulls/1595)) ([#1619](https://gitea.com/gitea/runner/pulls/1619))
- Add GITHUB_STEP_SUMMARY ([#1607](https://gitea.com/gitea/runner/pulls/1607))
- Crash malformed composite action ([#1616](https://gitea.com/gitea/runner/pulls/1616))
- Update output handling for reusable workflows ([#1521](https://gitea.com/gitea/runner/pulls/1521))
- Github.job property is empty, GITHUB_JOB should be job id ([#1646](https://gitea.com/gitea/runner/pulls/1646))
- Safe file name ([#1651](https://gitea.com/gitea/runner/pulls/1651))
- Crash if the id tool fails to run in the container (1660)
- Return err in walk ([#1667](https://gitea.com/gitea/runner/pulls/1667))
- Correct ref and ref_name ([#1672](https://gitea.com/gitea/runner/pulls/1672))
- Use os.UserHomeDir ([#1706](https://gitea.com/gitea/runner/pulls/1706))
- Add `server_url` attribute to github context ([#1727](https://gitea.com/gitea/runner/pulls/1727))
- Reusable workflow panic ([#1728](https://gitea.com/gitea/runner/pulls/1728))
- Ghc assignment typo ([#1729](https://gitea.com/gitea/runner/pulls/1729))
- Ensure networkmode "host" unless explicitly specified ([#1739](https://gitea.com/gitea/runner/pulls/1739))
- Environment handling windows (host mode) ([#1732](https://gitea.com/gitea/runner/pulls/1732))
- Try finding a socket, otherwise fail, respect user choice ([#1745](https://gitea.com/gitea/runner/pulls/1745))
- Remove hardcoded reference to github.com when using reusable remote workflows and remote actions ([#1784](https://gitea.com/gitea/runner/pulls/1784))
- Fallback to unauthenticated pull ([#1774](https://gitea.com/gitea/runner/pulls/1774))
- Don't allow `-self-hosted` mode as container image ([#1783](https://gitea.com/gitea/runner/pulls/1783))
- Update ARCH environment variable used in runners/actions ([#1818](https://gitea.com/gitea/runner/pulls/1818))
- Spelling mistake in readme ([#1854](https://gitea.com/gitea/runner/pulls/1854))
- Handle zero size ([#1888](https://gitea.com/gitea/runner/pulls/1888))
- Fail if no stages were found ([#1970](https://gitea.com/gitea/runner/pulls/1970))

### Refactoring

- `NewWorkflowPlanner` ([#648](https://gitea.com/gitea/runner/pulls/648))
- Return more errors, add more tests ([#679](https://gitea.com/gitea/runner/pulls/679))
- Remove `gotest.tools` ([#688](https://gitea.com/gitea/runner/pulls/688))
- Re-implement `embed` without "unused" import ([#830](https://gitea.com/gitea/runner/pulls/830))
- Export and move shared contexts into pkg/model ([#931](https://gitea.com/gitea/runner/pulls/931))
- Extract RunContext Executor in JobExecutor ([#984](https://gitea.com/gitea/runner/pulls/984))
- Extract setupAction into ActionReader ([#986](https://gitea.com/gitea/runner/pulls/986))
- Default empty if: statement to if: success() in evaluator & remove `FixIfStatement()` ([#990](https://gitea.com/gitea/runner/pulls/990))
- Simplify action function signatures ([#1083](https://gitea.com/gitea/runner/pulls/1083))
- Remove composite action runcontext workaround ([#1085](https://gitea.com/gitea/runner/pulls/1085))
- Remove github.com/pkg/errors dependency ([#1077](https://gitea.com/gitea/runner/pulls/1077))
- Ignore already closed error ([#1285](https://gitea.com/gitea/runner/pulls/1285))
- Move from io/ioutil to io and os packages ([#1417](https://gitea.com/gitea/runner/pulls/1417))
- Move autoremove into the jobexecutor ([#1463](https://gitea.com/gitea/runner/pulls/1463))
- Share UpdateFromEnv logic ([#1457](https://gitea.com/gitea/runner/pulls/1457))
- Fix add-path / GITHUB_PATH commands ([#1472](https://gitea.com/gitea/runner/pulls/1472))
- Fix savestate in pre steps ([#1466](https://gitea.com/gitea/runner/pulls/1466))
- Remove docker image list reference filter ([#1501](https://gitea.com/gitea/runner/pulls/1501))
- Pull and rebuild docker by default ([#1569](https://gitea.com/gitea/runner/pulls/1569))
- GITHUB_ENV command / remove env.PATH ([#1503](https://gitea.com/gitea/runner/pulls/1503))
- Open boltdb only while using it ([#1879](https://gitea.com/gitea/runner/pulls/1879))
- Docker build BuildContext field ([#1914](https://gitea.com/gitea/runner/pulls/1914))
- Remove unnecessary nil check in RunContext ([#1955](https://gitea.com/gitea/runner/pulls/1955))

### Build

- **deps:** Bump github.com/AlecAivazis/survey/v2 from 2.2.7 to 2.3.1 ([#803](https://gitea.com/gitea/runner/pulls/803))
- **deps:** Bump github.com/go-git/go-billy/v5 from 5.0.0 to 5.3.1 ([#802](https://gitea.com/gitea/runner/pulls/802))
- **deps:** Bump codecov/codecov-action from 1.5.2 to 2.1.0 ([#819](https://gitea.com/gitea/runner/pulls/819))
- **deps:** Bump github.com/joho/godotenv from 1.3.0 to 1.4.0 ([#824](https://gitea.com/gitea/runner/pulls/824))
- **deps:** Bump github.com/moby/buildkit from 0.9.0 to 0.9.1 ([#838](https://gitea.com/gitea/runner/pulls/838))
- **deps:** Bump github.com/docker/cli ([#854](https://gitea.com/gitea/runner/pulls/854))
- **deps:** Bump github.com/docker/docker ([#855](https://gitea.com/gitea/runner/pulls/855))
- **deps:** Bump github.com/moby/buildkit from 0.9.1 to 0.9.2 ([#864](https://gitea.com/gitea/runner/pulls/864))
- **deps:** Bump github.com/go-ini/ini from 1.63.2 to 1.64.0 ([#882](https://gitea.com/gitea/runner/pulls/882))
- **deps:** Bump github.com/docker/cli ([#956](https://gitea.com/gitea/runner/pulls/956))
- **deps:** Bump github.com/opencontainers/image-spec ([#954](https://gitea.com/gitea/runner/pulls/954))
- **deps:** Bump github.com/docker/docker ([#957](https://gitea.com/gitea/runner/pulls/957))
- **deps:** Bump github.com/moby/buildkit from 0.9.2 to 0.9.3 ([#958](https://gitea.com/gitea/runner/pulls/958))
- **deps:** Bump github.com/docker/distribution ([#994](https://gitea.com/gitea/runner/pulls/994))
- **deps:** Bump github.com/go-ini/ini from 1.64.0 to 1.66.4 ([#997](https://gitea.com/gitea/runner/pulls/997))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.8 to 1.6.9 ([#1018](https://gitea.com/gitea/runner/pulls/1018))
- **deps:** Bump github.com/spf13/cobra from 1.2.1 to 1.3.0 ([#955](https://gitea.com/gitea/runner/pulls/955))
- **deps:** Bump github.com/moby/buildkit from 0.9.3 to 0.10.0 ([#1036](https://gitea.com/gitea/runner/pulls/1036))
- **deps:** Bump github.com/docker/distribution ([#1037](https://gitea.com/gitea/runner/pulls/1037))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.9 to 1.6.10 ([#1039](https://gitea.com/gitea/runner/pulls/1039))
- **deps:** Bump github.com/docker/cli ([#1038](https://gitea.com/gitea/runner/pulls/1038))
- **deps:** Bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#1040](https://gitea.com/gitea/runner/pulls/1040))
- **deps:** Bump actions/checkout from 2 to 3 ([#1028](https://gitea.com/gitea/runner/pulls/1028))
- **deps:** Bump golangci/golangci-lint-action from 2 to 3.1.0 ([#1019](https://gitea.com/gitea/runner/pulls/1019))
- **deps:** Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#1061](https://gitea.com/gitea/runner/pulls/1061))
- **deps:** Bump github.com/docker/docker ([#1062](https://gitea.com/gitea/runner/pulls/1062))
- **deps:** Bump github.com/docker/cli ([#1081](https://gitea.com/gitea/runner/pulls/1081))
- **deps:** Bump github.com/AlecAivazis/survey/v2 from 2.3.2 to 2.3.4 ([#1079](https://gitea.com/gitea/runner/pulls/1079))
- **deps:** Bump github.com/docker/docker ([#1080](https://gitea.com/gitea/runner/pulls/1080))
- **deps:** Bump actions/cache from 2 to 3 ([#1082](https://gitea.com/gitea/runner/pulls/1082))
- **deps:** Bump github.com/opencontainers/selinux ([#1102](https://gitea.com/gitea/runner/pulls/1102))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.10 to 1.6.11 ([#1106](https://gitea.com/gitea/runner/pulls/1106))
- **deps:** Bump deps which dependabot doesn't want to ([#1107](https://gitea.com/gitea/runner/pulls/1107))
- **deps:** Bump actions/stale from 4 to 5 ([#1117](https://gitea.com/gitea/runner/pulls/1117))
- **deps:** Bump actions/setup-go from 2 to 3 ([#1120](https://gitea.com/gitea/runner/pulls/1120))
- **deps:** Bump actions/upload-artifact from 2 to 3 ([#1119](https://gitea.com/gitea/runner/pulls/1119))
- **deps:** Bump github.com/moby/buildkit from 0.10.0 to 0.10.1 ([#1116](https://gitea.com/gitea/runner/pulls/1116))
- **deps:** Bump codecov/codecov-action from 2.1.0 to 3.0.0 ([#1118](https://gitea.com/gitea/runner/pulls/1118))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.11 to 1.6.12 ([#1126](https://gitea.com/gitea/runner/pulls/1126))
- **deps:** Bump codecov/codecov-action from 3.0.0 to 3.1.0 ([#1139](https://gitea.com/gitea/runner/pulls/1139))
- **deps:** Bump github.com/moby/buildkit from 0.10.1 to 0.10.2 ([#1147](https://gitea.com/gitea/runner/pulls/1147))
- **deps:** Bump github.com/docker/cli ([#1154](https://gitea.com/gitea/runner/pulls/1154))
- **deps:** Bump github.com/docker/docker ([#1155](https://gitea.com/gitea/runner/pulls/1155))
- **deps:** Bump docker/setup-qemu-action from 1 to 2 ([#1156](https://gitea.com/gitea/runner/pulls/1156))
- **deps:** Bump github.com/moby/buildkit from 0.10.2 to 0.10.3 ([#1153](https://gitea.com/gitea/runner/pulls/1153))
- **deps:** Bump github.com/docker/cli ([#1168](https://gitea.com/gitea/runner/pulls/1168))
- **deps:** Bump github.com/docker/docker ([#1167](https://gitea.com/gitea/runner/pulls/1167))
- **deps:** Bump golangci/golangci-lint-action from 3.1.0 to 3.2.0 ([#1169](https://gitea.com/gitea/runner/pulls/1169))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.12 to 1.6.13 ([#1173](https://gitea.com/gitea/runner/pulls/1173))
- **deps:** Bump goreleaser/goreleaser-action from 2 to 3 ([#1179](https://gitea.com/gitea/runner/pulls/1179))
- **deps:** Bump github.com/go-ini/ini from 1.66.4 to 1.66.5 ([#1189](https://gitea.com/gitea/runner/pulls/1189))
- **deps:** Bump github.com/go-ini/ini from 1.66.5 to 1.66.6 ([#1201](https://gitea.com/gitea/runner/pulls/1201))
- **deps:** Bump github.com/docker/docker ([#1212](https://gitea.com/gitea/runner/pulls/1212))
- **deps:** Bump github.com/docker/cli ([#1211](https://gitea.com/gitea/runner/pulls/1211))
- **deps:** Bump github.com/AlecAivazis/survey/v2 from 2.3.4 to 2.3.5 ([#1214](https://gitea.com/gitea/runner/pulls/1214))
- **deps:** Bump github.com/stretchr/testify from 1.7.1 to 1.7.2 ([#1213](https://gitea.com/gitea/runner/pulls/1213))
- **deps:** Bump github.com/stretchr/testify from 1.7.2 to 1.7.5 ([#1232](https://gitea.com/gitea/runner/pulls/1232))
- **deps:** Bump github.com/spf13/cobra from 1.4.0 to 1.5.0 ([#1233](https://gitea.com/gitea/runner/pulls/1233))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.13 to 1.6.14 ([#1234](https://gitea.com/gitea/runner/pulls/1234))
- **deps:** Bump github.com/stretchr/testify from 1.7.5 to 1.8.0 ([#1240](https://gitea.com/gitea/runner/pulls/1240))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.14 to 1.6.15 ([#1239](https://gitea.com/gitea/runner/pulls/1239))
- **deps:** Bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1273](https://gitea.com/gitea/runner/pulls/1273))
- **deps:** Bump github.com/go-ini/ini from 1.66.6 to 1.67.0 ([#1305](https://gitea.com/gitea/runner/pulls/1305))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.15 to 1.6.16 ([#1311](https://gitea.com/gitea/runner/pulls/1311))
- **deps:** Bump github.com/mattn/go-isatty from 0.0.14 to 0.0.16 ([#1312](https://gitea.com/gitea/runner/pulls/1312))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.16 to 1.6.17 ([#1321](https://gitea.com/gitea/runner/pulls/1321))
- **deps:** Bump github.com/moby/buildkit from 0.10.3 to 0.10.4 ([#1322](https://gitea.com/gitea/runner/pulls/1322))
- **deps:** Bump github.com/docker/docker ([#1339](https://gitea.com/gitea/runner/pulls/1339))
- **deps:** Bump github.com/docker/cli ([#1340](https://gitea.com/gitea/runner/pulls/1340))
- **deps:** Bump github.com/AlecAivazis/survey/v2 from 2.3.5 to 2.3.6 ([#1341](https://gitea.com/gitea/runner/pulls/1341))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.17 to 1.6.18 ([#1352](https://gitea.com/gitea/runner/pulls/1352))
- **deps:** Bump actions/stale from 5 to 6 ([#1357](https://gitea.com/gitea/runner/pulls/1357))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.18 to 1.6.19 ([#1360](https://gitea.com/gitea/runner/pulls/1360))
- **deps:** Bump codecov/codecov-action from 3.1.0 to 3.1.1 ([#1359](https://gitea.com/gitea/runner/pulls/1359))
- **deps:** Bump github.com/opencontainers/selinux ([#1372](https://gitea.com/gitea/runner/pulls/1372))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.19 to 1.6.20 ([#1371](https://gitea.com/gitea/runner/pulls/1371))
- **deps:** Bump megalinter/megalinter from 5 to 6.11.0 ([#1370](https://gitea.com/gitea/runner/pulls/1370))
- **deps:** Bump github.com/imdario/mergo from 0.3.12 to 0.3.13 ([#1378](https://gitea.com/gitea/runner/pulls/1378))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.20 to 1.6.21 ([#1380](https://gitea.com/gitea/runner/pulls/1380))
- **deps:** Bump gotest.tools/v3 from 3.0.3 to 3.4.0 ([#1379](https://gitea.com/gitea/runner/pulls/1379))
- **deps:** Bump megalinter/megalinter from 6.11.0 to 6.11.1 ([#1381](https://gitea.com/gitea/runner/pulls/1381))
- **deps:** Bump github.com/docker/docker ([#1392](https://gitea.com/gitea/runner/pulls/1392))
- **deps:** Bump github.com/spf13/cobra from 1.5.0 to 1.6.0 ([#1394](https://gitea.com/gitea/runner/pulls/1394))
- **deps:** Bump github.com/docker/cli ([#1393](https://gitea.com/gitea/runner/pulls/1393))
- **deps:** Bump megalinter/megalinter from 6.11.1 to 6.12.0 ([#1395](https://gitea.com/gitea/runner/pulls/1395))
- **deps:** Bump github.com/moby/buildkit from 0.10.4 to 0.10.5 ([#1406](https://gitea.com/gitea/runner/pulls/1406))
- **deps:** Bump github.com/docker/cli ([#1405](https://gitea.com/gitea/runner/pulls/1405))
- **deps:** Bump github.com/docker/docker ([#1403](https://gitea.com/gitea/runner/pulls/1403))
- **deps:** Bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 ([#1408](https://gitea.com/gitea/runner/pulls/1408))
- **deps:** Bump megalinter/megalinter from 6.12.0 to 6.13.0 ([#1407](https://gitea.com/gitea/runner/pulls/1407))
- **deps:** Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1404](https://gitea.com/gitea/runner/pulls/1404))
- **deps:** Bump github.com/docker/cli ([#1419](https://gitea.com/gitea/runner/pulls/1419))
- **deps:** Bump github.com/docker/docker ([#1418](https://gitea.com/gitea/runner/pulls/1418))
- **deps:** Bump github.com/spf13/cobra from 1.6.0 to 1.6.1 ([#1420](https://gitea.com/gitea/runner/pulls/1420))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.21 to 1.6.22 ([#1427](https://gitea.com/gitea/runner/pulls/1427))
- **deps:** Bump megalinter/megalinter from 6.13.0 to 6.14.0 ([#1426](https://gitea.com/gitea/runner/pulls/1426))
- **deps:** Bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 ([#1436](https://gitea.com/gitea/runner/pulls/1436))
- **deps:** Bump github.com/moby/buildkit from 0.10.5 to 0.10.6 ([#1437](https://gitea.com/gitea/runner/pulls/1437))
- **deps:** Bump github.com/creack/pty from 1.1.17 to 1.1.18 ([#1447](https://gitea.com/gitea/runner/pulls/1447))
- **deps:** Bump megalinter/megalinter from 6.14.0 to 6.15.0 ([#1475](https://gitea.com/gitea/runner/pulls/1475))
- **deps:** Bump github.com/docker/docker ([#1517](https://gitea.com/gitea/runner/pulls/1517))
- **deps:** Bump github.com/docker/cli ([#1516](https://gitea.com/gitea/runner/pulls/1516))
- **deps:** Bump goreleaser/goreleaser-action from 3 to 4 ([#1515](https://gitea.com/gitea/runner/pulls/1515))
- **deps:** Bump actions/stale from 6 to 7 ([#1535](https://gitea.com/gitea/runner/pulls/1535))
- **deps:** Bump megalinter/megalinter from 6.15.0 to 6.16.0 ([#1534](https://gitea.com/gitea/runner/pulls/1534))
- **deps:** Bump github.com/go-git/go-billy/v5 from 5.3.1 to 5.4.0 ([#1541](https://gitea.com/gitea/runner/pulls/1541))
- **deps:** Bump github.com/mattn/go-isatty from 0.0.16 to 0.0.17 ([#1539](https://gitea.com/gitea/runner/pulls/1539))
- **deps:** Bump megalinter/megalinter from 6.16.0 to 6.17.0 ([#1540](https://gitea.com/gitea/runner/pulls/1540))
- **deps:** Bump megalinter/megalinter from 6.17.0 to 6.18.0 ([#1550](https://gitea.com/gitea/runner/pulls/1550))
- **deps:** Bump github.com/containerd/containerd from 1.6.6 to 1.6.12 ([#1554](https://gitea.com/gitea/runner/pulls/1554))
- **deps:** Bump github.com/moby/buildkit from 0.10.6 to 0.11.0 ([#1563](https://gitea.com/gitea/runner/pulls/1563))
- **deps:** Bump actions/github-script from 5 to 6 ([#1572](https://gitea.com/gitea/runner/pulls/1572))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.22 to 1.6.23 ([#1574](https://gitea.com/gitea/runner/pulls/1574))
- **deps:** Bump github.com/docker/docker ([#1573](https://gitea.com/gitea/runner/pulls/1573))
- **deps:** Bump golang.org/x/term from 0.3.0 to 0.4.0 ([#1577](https://gitea.com/gitea/runner/pulls/1577))
- **deps:** Bump github.com/moby/buildkit from 0.11.0 to 0.11.1 ([#1576](https://gitea.com/gitea/runner/pulls/1576))
- **deps:** Bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 ([#1590](https://gitea.com/gitea/runner/pulls/1590))
- **deps:** Bump github.com/moby/buildkit from 0.11.1 to 0.11.2 ([#1589](https://gitea.com/gitea/runner/pulls/1589))
- **deps:** Bump github.com/joho/godotenv from 1.4.0 to 1.5.1 ([#1614](https://gitea.com/gitea/runner/pulls/1614))
- **deps:** Bump github.com/go-git/go-billy/v5 from 5.4.0 to 5.4.1 ([#1612](https://gitea.com/gitea/runner/pulls/1612))
- **deps:** Bump megalinter/megalinter from 6.18.0 to 6.19.0 ([#1610](https://gitea.com/gitea/runner/pulls/1610))
- **deps:** Bump github.com/docker/docker ([#1613](https://gitea.com/gitea/runner/pulls/1613))
- **deps:** Bump github.com/docker/cli from 23.0.0-rc.1+incompatible to 23.0.0+incompatible ([#1611](https://gitea.com/gitea/runner/pulls/1611))
- **deps:** Bump github.com/docker/cli ([#1626](https://gitea.com/gitea/runner/pulls/1626))
- **deps:** Bump github.com/docker/docker ([#1624](https://gitea.com/gitea/runner/pulls/1624))
- **deps:** Bump github.com/opencontainers/selinux ([#1625](https://gitea.com/gitea/runner/pulls/1625))
- **deps:** Bump github.com/moby/buildkit from 0.11.2 to 0.11.3 ([#1635](https://gitea.com/gitea/runner/pulls/1635))
- **deps:** Bump golang.org/x/net from 0.4.0 to 0.7.0 ([#1636](https://gitea.com/gitea/runner/pulls/1636))
- **deps:** Bump github.com/containerd/containerd from 1.6.16 to 1.6.18 ([#1637](https://gitea.com/gitea/runner/pulls/1637))
- **deps:** Bump github.com/opencontainers/image-spec ([#1649](https://gitea.com/gitea/runner/pulls/1649))
- **deps:** Bump golang.org/x/term from 0.5.0 to 0.6.0 ([#1666](https://gitea.com/gitea/runner/pulls/1666))
- **deps:** Bump megalinter/megalinter from 6.19.0 to 6.20.0 ([#1665](https://gitea.com/gitea/runner/pulls/1665))
- **deps:** Bump github.com/moby/buildkit from 0.11.3 to 0.11.4 ([#1669](https://gitea.com/gitea/runner/pulls/1669))
- **deps:** Bump fregante/setup-git-user from 1 to 2 ([#1664](https://gitea.com/gitea/runner/pulls/1664))
- **deps:** Bump megalinter/megalinter from 6.20.0 to 6.20.1 ([#1679](https://gitea.com/gitea/runner/pulls/1679))
- **deps:** Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#1650](https://gitea.com/gitea/runner/pulls/1650))
- **deps:** Bump github.com/imdario/mergo from 0.3.13 to 0.3.14 ([#1688](https://gitea.com/gitea/runner/pulls/1688))
- **deps:** Bump actions/setup-go from 3 to 4 ([#1689](https://gitea.com/gitea/runner/pulls/1689))
- **deps:** Bump megalinter/megalinter from 6.20.1 to 6.21.0 ([#1699](https://gitea.com/gitea/runner/pulls/1699))
- **deps:** Bump actions/stale from 7 to 8 ([#1700](https://gitea.com/gitea/runner/pulls/1700))
- **deps:** Bump github.com/imdario/mergo from 0.3.14 to 0.3.15 ([#1701](https://gitea.com/gitea/runner/pulls/1701))
- **deps:** Bump github.com/mattn/go-isatty from 0.0.17 to 0.0.18 ([#1702](https://gitea.com/gitea/runner/pulls/1702))
- **deps:** Bump github.com/moby/buildkit from 0.11.4 to 0.11.5 ([#1703](https://gitea.com/gitea/runner/pulls/1703))
- **deps:** Bump github.com/opencontainers/runc from 1.1.3 to 1.1.5 ([#1708](https://gitea.com/gitea/runner/pulls/1708))
- **deps:** Bump megalinter/megalinter from 6.21.0 to 6.22.1 ([#1710](https://gitea.com/gitea/runner/pulls/1710))
- **deps:** Bump github.com/docker/docker ([#1711](https://gitea.com/gitea/runner/pulls/1711))
- **deps:** Bump github.com/docker/cli ([#1712](https://gitea.com/gitea/runner/pulls/1712))
- **deps:** Bump golang.org/x/term from 0.6.0 to 0.7.0 ([#1723](https://gitea.com/gitea/runner/pulls/1723))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.23 to 1.6.24 ([#1721](https://gitea.com/gitea/runner/pulls/1721))
- **deps:** Bump megalinter/megalinter from 6.22.1 to 6.22.2 ([#1720](https://gitea.com/gitea/runner/pulls/1720))
- **deps:** Bump github.com/docker/cli ([#1722](https://gitea.com/gitea/runner/pulls/1722))
- **deps:** Bump github.com/spf13/cobra from 1.6.1 to 1.7.0 ([#1724](https://gitea.com/gitea/runner/pulls/1724))
- **deps:** Bump github.com/docker/docker ([#1719](https://gitea.com/gitea/runner/pulls/1719))
- **deps:** Update go-git to ce62f3e9ff86270538a514a68d3bd5563a733e3b ([#1725](https://gitea.com/gitea/runner/pulls/1725))
- **deps:** Bump codecov/codecov-action from 3.1.1 to 3.1.2 ([#1735](https://gitea.com/gitea/runner/pulls/1735))
- **deps:** Bump codecov/codecov-action from 3.1.2 to 3.1.3 ([#1752](https://gitea.com/gitea/runner/pulls/1752))
- **deps:** Bump github.com/docker/cli ([#1753](https://gitea.com/gitea/runner/pulls/1753))
- **deps:** Bump github.com/docker/docker ([#1754](https://gitea.com/gitea/runner/pulls/1754))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1755](https://gitea.com/gitea/runner/pulls/1755))
- **deps:** Bump github.com/docker/docker ([#1779](https://gitea.com/gitea/runner/pulls/1779))
- **deps:** Bump github.com/opencontainers/image-spec ([#1781](https://gitea.com/gitea/runner/pulls/1781))
- **deps:** Bump github.com/docker/cli ([#1780](https://gitea.com/gitea/runner/pulls/1780))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1791](https://gitea.com/gitea/runner/pulls/1791))
- **deps:** Bump github.com/cloudflare/circl from 1.1.0 to 1.3.3 ([#1800](https://gitea.com/gitea/runner/pulls/1800))
- **deps:** Bump github.com/docker/distribution ([#1801](https://gitea.com/gitea/runner/pulls/1801))
- **deps:** Bump github.com/docker/cli ([#1805](https://gitea.com/gitea/runner/pulls/1805))
- **deps:** Bump github.com/docker/docker ([#1806](https://gitea.com/gitea/runner/pulls/1806))
- **deps:** Bump github.com/opencontainers/image-spec ([#1807](https://gitea.com/gitea/runner/pulls/1807))
- **deps:** Bump codecov/codecov-action from 3.1.3 to 3.1.4 ([#1813](https://gitea.com/gitea/runner/pulls/1813))
- **deps:** Bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2 ([#1814](https://gitea.com/gitea/runner/pulls/1814))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1817](https://gitea.com/gitea/runner/pulls/1817))
- **deps:** Bump megalinter/megalinter from 6.22.2 to 7.0.2 ([#1827](https://gitea.com/gitea/runner/pulls/1827))
- **deps:** Bump github.com/go-git/go-git/v5 ([#1830](https://gitea.com/gitea/runner/pulls/1830))
- **deps:** Bump github.com/docker/cli ([#1816](https://gitea.com/gitea/runner/pulls/1816))
- **deps:** Bump github.com/opencontainers/image-spec ([#1829](https://gitea.com/gitea/runner/pulls/1829))
- **deps:** Bump golangci/golangci-lint-action from 3.4.0 to 3.5.0 ([#1839](https://gitea.com/gitea/runner/pulls/1839))
- **deps:** Bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3 ([#1841](https://gitea.com/gitea/runner/pulls/1841))
- **deps:** Bump github.com/imdario/mergo from 0.3.15 to 0.3.16 ([#1840](https://gitea.com/gitea/runner/pulls/1840))
- **deps:** Bump github.com/docker/cli ([#1843](https://gitea.com/gitea/runner/pulls/1843))
- **deps:** Bump megalinter/megalinter from 7.0.2 to 7.0.4 ([#1838](https://gitea.com/gitea/runner/pulls/1838))
- **deps:** Bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://gitea.com/gitea/runner/pulls/1842))
- **deps:** Bump github.com/mattn/go-isatty from 0.0.18 to 0.0.19 ([#1859](https://gitea.com/gitea/runner/pulls/1859))
- **deps:** Bump megalinter/megalinter from 7.0.4 to 7.1.0 ([#1857](https://gitea.com/gitea/runner/pulls/1857))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1858](https://gitea.com/gitea/runner/pulls/1858))
- **deps:** Bump github.com/opencontainers/image-spec ([#1871](https://gitea.com/gitea/runner/pulls/1871))
- **deps:** Bump golang.org/x/term from 0.8.0 to 0.9.0 ([#1869](https://gitea.com/gitea/runner/pulls/1869))
- **deps:** Bump github.com/AlecAivazis/survey/v2 from 2.3.6 to 2.3.7 ([#1872](https://gitea.com/gitea/runner/pulls/1872))
- **deps:** Bump github.com/rhysd/actionlint from 1.6.24 to 1.6.25 ([#1870](https://gitea.com/gitea/runner/pulls/1870))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1880](https://gitea.com/gitea/runner/pulls/1880))
- **deps:** Bump github.com/opencontainers/image-spec ([#1881](https://gitea.com/gitea/runner/pulls/1881))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1882](https://gitea.com/gitea/runner/pulls/1882))
- **deps:** Bump github.com/opencontainers/image-spec ([#1883](https://gitea.com/gitea/runner/pulls/1883))
- **deps:** Bump github.com/opencontainers/image-spec ([#1893](https://gitea.com/gitea/runner/pulls/1893))
- **deps:** Bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1884](https://gitea.com/gitea/runner/pulls/1884))
- **deps:** Bump gotest.tools/v3 from 3.4.0 to 3.5.0 ([#1892](https://gitea.com/gitea/runner/pulls/1892))
- **deps:** Bump golangci/golangci-lint-action from 3.5.0 to 3.6.0 ([#1868](https://gitea.com/gitea/runner/pulls/1868))
- **deps:** Bump golang.org/x/term from 0.9.0 to 0.10.0 ([#1906](https://gitea.com/gitea/runner/pulls/1906))
- **deps:** Bump github.com/moby/buildkit from 0.11.6 to 0.12.0 ([#1917](https://gitea.com/gitea/runner/pulls/1917))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.7.0 to 5.8.0 ([#1925](https://gitea.com/gitea/runner/pulls/1925))
- **deps:** Bump github.com/docker/docker ([#1933](https://gitea.com/gitea/runner/pulls/1933))
- **deps:** Bump github.com/docker/cli ([#1932](https://gitea.com/gitea/runner/pulls/1932))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.8.0 to 5.8.1 ([#1934](https://gitea.com/gitea/runner/pulls/1934))
- **deps:** Bump github.com/moby/buildkit from 0.12.0 to 0.12.1 ([#1947](https://gitea.com/gitea/runner/pulls/1947))
- **deps:** Bump golang.org/x/term from 0.10.0 to 0.11.0 ([#1948](https://gitea.com/gitea/runner/pulls/1948))
- **deps:** Bump megalinter/megalinter from 7.1.0 to 7.2.1 ([#1931](https://gitea.com/gitea/runner/pulls/1931))
- **deps:** Bump megalinter/megalinter from 7.2.1 to 7.3.0 ([#1965](https://gitea.com/gitea/runner/pulls/1965))
- **deps:** Bump golangci/golangci-lint-action from 3.6.0 to 3.7.0 ([#1978](https://gitea.com/gitea/runner/pulls/1978))
- **deps:** Bump github.com/moby/patternmatcher from 0.5.0 to 0.6.0 ([#1985](https://gitea.com/gitea/runner/pulls/1985))
- **deps:** Bump github.com/moby/buildkit from 0.12.1 to 0.12.2 ([#1986](https://gitea.com/gitea/runner/pulls/1986))
- **deps:** Bump golang.org/x/term from 0.11.0 to 0.12.0 ([#1992](https://gitea.com/gitea/runner/pulls/1992))
- **deps:** Bump github.com/docker/cli ([#1999](https://gitea.com/gitea/runner/pulls/1999))
- **deps:** Bump github.com/docker/docker ([#2000](https://gitea.com/gitea/runner/pulls/2000))
- **deps:** Bump goreleaser/goreleaser-action from 4 to 5 ([#2008](https://gitea.com/gitea/runner/pulls/2008))
- **deps:** Bump docker/setup-qemu-action from 2 to 3 ([#2007](https://gitea.com/gitea/runner/pulls/2007))
- **deps:** Bump github.com/go-git/go-billy/v5 from 5.4.1 to 5.5.0 ([#2012](https://gitea.com/gitea/runner/pulls/2012))
- **deps:** Bump gotest.tools/v3 from 3.5.0 to 3.5.1 ([#2009](https://gitea.com/gitea/runner/pulls/2009))
- **deps:** Bump github.com/opencontainers/image-spec ([#2010](https://gitea.com/gitea/runner/pulls/2010))
- **deps:** Bump github.com/go-git/go-git/v5 from 5.8.1 to 5.9.0 ([#2011](https://gitea.com/gitea/runner/pulls/2011))

### CI

- Pin goreleaser version ([#825](https://gitea.com/gitea/runner/pulls/825))
- Switch macos job to macos-10.15 ([#849](https://gitea.com/gitea/runner/pulls/849))
- Remove macOS tests ([#860](https://gitea.com/gitea/runner/pulls/860))
- **mergify:** Fix deprecated strict mode ([#862](https://gitea.com/gitea/runner/pulls/862))
- **golangci-lint:** Switch to latest ([#870](https://gitea.com/gitea/runner/pulls/870))
- **choco:** Fix chocolatey docker build ([#903](https://gitea.com/gitea/runner/pulls/903))
- **Mergify:** Configuration update ([#924](https://gitea.com/gitea/runner/pulls/924))
- Replace superlinter with megalinter ([#923](https://gitea.com/gitea/runner/pulls/923))
- **Mergify:** Configuration update ([#966](https://gitea.com/gitea/runner/pulls/966))
- **Mergify:** Configuration update ([#972](https://gitea.com/gitea/runner/pulls/972))
- **stalebot:** Ignore feature requests ([#1017](https://gitea.com/gitea/runner/pulls/1017))
- Bump go to 1.18 ([#1087](https://gitea.com/gitea/runner/pulls/1087))
- Always grab latest go version ([#1124](https://gitea.com/gitea/runner/pulls/1124))
- **choco:** Update chocolatey to 1.1.0 ([#1164](https://gitea.com/gitea/runner/pulls/1164))
- Add Winget Releaser workflow ([#1623](https://gitea.com/gitea/runner/pulls/1623))
- Deduplicate running workflows ([#1751](https://gitea.com/gitea/runner/pulls/1751))
- Change winget job runner to `ubuntu-latest` ([#1959](https://gitea.com/gitea/runner/pulls/1959))

### Documentation

- Remove unnecessary brew command
- **README:** Add --env-file flag ([#147](https://gitea.com/gitea/runner/pulls/147))
- Add --secret-file example ([#395](https://gitea.com/gitea/runner/pulls/395))
- Fix brew formula link ([#846](https://gitea.com/gitea/runner/pulls/846))
- Add brew installation from HEAD ([#886](https://gitea.com/gitea/runner/pulls/886))
- Expand the GITHUB_TOKEN section ([#968](https://gitea.com/gitea/runner/pulls/968))
- **README:** Remove go install ([#1123](https://gitea.com/gitea/runner/pulls/1123))
- Add naming rule for git ref
- Clarifying skipping steps / jobs ([#1480](https://gitea.com/gitea/runner/pulls/1480))
- Remove help section ([#1648](https://gitea.com/gitea/runner/pulls/1648))

### Tests

- **runner:** Move arm64 test to own func ([#1095](https://gitea.com/gitea/runner/pulls/1095))
- Remove runTestJobFile ([#1127](https://gitea.com/gitea/runner/pulls/1127))
- Add test for networking setup in act ([#1375](https://gitea.com/gitea/runner/pulls/1375))
- Make sure workflow_call is not a github event calling our workflow ([#1520](https://gitea.com/gitea/runner/pulls/1520))
- Do not leak step env in composite ([#1585](https://gitea.com/gitea/runner/pulls/1585))
- Enshure ForcePull config doesn't break docker actions (1661)
## [0.4.1] - 2026-04-19

### Performance

- Use single poller with semaphore-based capacity control ([#822](https://gitea.com/gitea/runner/pulls/822))
## [0.4.0] - 2026-04-15

### Features

- Add Prometheus metrics endpoint for runner observability ([#820](https://gitea.com/gitea/runner/pulls/820))

### Dependencies

- **deps:** Upgrade golangci-lint from v2.10.1 to v2.11.4 ([#821](https://gitea.com/gitea/runner/pulls/821))

### Performance

- Reduce runner-to-server connection load with adaptive reporting and polling ([#819](https://gitea.com/gitea/runner/pulls/819))
## [0.3.1] - 2026-03-26

### Features

- Add configurable bind_workdir option with workspace cleanup for DinD setups ([#810](https://gitea.com/gitea/runner/pulls/810))

### Dependencies

- **deps:** Update workflow action versions ([#792](https://gitea.com/gitea/runner/pulls/792))
- **deps:** Update go dependencies ([#802](https://gitea.com/gitea/runner/pulls/802))

### Chores

- **lint:** Add golangci-lint v2 and fix all lint issues ([#803](https://gitea.com/gitea/runner/pulls/803))

### Bug fixes

- Race condition in reporter between RunDaemon and Close ([#796](https://gitea.com/gitea/runner/pulls/796))
- Composite action log result reported as step result ([#801](https://gitea.com/gitea/runner/pulls/801))
## [0.3.0] - 2026-02-18

### Dependencies

- **deps:** Upgrade dependencies ([#775](https://gitea.com/gitea/runner/pulls/775))
- **deps:** Bump Go version to 1.26.0 ([#787](https://gitea.com/gitea/runner/pulls/787))
- **deps:** Bump Go version to 1.26 in Dockerfile and Makefile ([#788](https://gitea.com/gitea/runner/pulls/788))
- **deps:** Update act to v0.261.8 ([#791](https://gitea.com/gitea/runner/pulls/791))

### Chores

- Update config.example.yaml for https://gitea.com/gitea/act_runner/pulls/724 ([#763](https://gitea.com/gitea/runner/pulls/763))

### Other

- Improve DIND rootless network performance ([#786](https://gitea.com/gitea/runner/pulls/786))

### Bug fixes

- **dockerfile:** Pin docker dind images to version 28, as version 29 has breaking changes in the API that are currently causing problems with the `act` fork of Gitea ([#769](https://gitea.com/gitea/runner/pulls/769))

### CI

- Release binary for linux/loong64 ([#756](https://gitea.com/gitea/runner/pulls/756))
## [0.2.13] - 2025-08-28

### Features

- **docker:** Add TZ env variable working ([#738](https://gitea.com/gitea/runner/pulls/738))
- Support github mirror ([#716](https://gitea.com/gitea/runner/pulls/716))

### Other

- **workflow:** Use matrix to compile different docker images ([#740](https://gitea.com/gitea/runner/pulls/740))
## [0.2.12] - 2025-06-11

### Features

- Add once flag to daemon command ([#19](https://gitea.com/gitea/runner/pulls/19)) ([#598](https://gitea.com/gitea/runner/pulls/598))
- Docker env vars for ephemeral and once ([#685](https://gitea.com/gitea/runner/pulls/685))
- Register interactive with values from cli ([#682](https://gitea.com/gitea/runner/pulls/682))

### Bug fixes

- Invalid label NoInteractive exit code ([#683](https://gitea.com/gitea/runner/pulls/683))
## [0.2.11] - 2024-09-24

### Features

- Allow graceful shutdowns ([#546](https://gitea.com/gitea/runner/pulls/546))

### CI

- Standardize code style and update version extraction ([#566](https://gitea.com/gitea/runner/pulls/566))
## [0.2.7] - 2024-03-27

### Bug fixes

- Exit docker container on any process exit ([#462](https://gitea.com/gitea/runner/pulls/462))
## [0.2.4] - 2023-07-25

### Other

- Default container workspace set to host path ([#279](https://gitea.com/gitea/runner/pulls/279))

### Bug fixes

- **register:** Refactor context usage in registration functions ([#286](https://gitea.com/gitea/runner/pulls/286))

### Refactoring

- **register:** Refactor registration functions and improve logging ([#288](https://gitea.com/gitea/runner/pulls/288))

### CI

- **lint:** Refactor code for clarity and linting compliance ([#289](https://gitea.com/gitea/runner/pulls/289))
## [0.2.1] - 2023-06-30

### Features

- Improve Docker configuration and detection handling ([#242](https://gitea.com/gitea/runner/pulls/242))
## [0.2.0] - 2023-06-15

### Refactoring

- Daemon, config, and logging for better clarity ([#225](https://gitea.com/gitea/runner/pulls/225))

### Documentation

- Improve examples README and organization ([#230](https://gitea.com/gitea/runner/pulls/230))
## [0.1.8] - 2023-05-09

### Other

- Improve signal handling by adding a runtime init ([#180](https://gitea.com/gitea/runner/pulls/180))
## [0.1.7] - 2023-05-06

### Chores

- Improve Dockerfile, README, and testing settings ([#172](https://gitea.com/gitea/runner/pulls/172))

### CI

- Improve release process and test robustness ([#173](https://gitea.com/gitea/runner/pulls/173))
## [0.1.6] - 2023-05-01

### Build

- Optimize Dockerfile and update dependencies ([#162](https://gitea.com/gitea/runner/pulls/162))
- Improve compression and update GitHub actions ([#168](https://gitea.com/gitea/runner/pulls/168))
## [0.1.5] - 2023-04-28

### CI

- Improve API usage and test robustness across platforms ([#159](https://gitea.com/gitea/runner/pulls/159))
## [0.1.3] - 2023-04-27

### Other

- Skip --disable-content-trust at docker buildx ([#139](https://gitea.com/gitea/runner/pulls/139))

### CI

- **actions:** Add build docker image workflow ([#118](https://gitea.com/gitea/runner/pulls/118))
## [0.1.1] - 2023-04-11

### Features

- Add artifact api from gitea server ([#103](https://gitea.com/gitea/runner/pulls/103))
## [0.0.1] - 2023-02-26

### Features

- Add task log output hook
- Move main task logic to runtime package
- Fill default task runner
- Task can report step and final result
- Replace ciphertext in log
- Avoid reuse of containers
- **cli:** Add register command basic functions
- **cli:** Set interactive stages in register command
- **cli:** Make register interactive stages working
- **poller:** Support scheduler to fetch task and dispatch to worker
- Support post step log
- Parse register arguments and do register when non-interactive mode
- Auto remove containers
- **register:** Add label format validation
- **register:** Add labels to local runner file
- Use specified labels
- Don't require docker ([#16](https://gitea.com/gitea/runner/pulls/16))

### Chores

- Fix lint error
- Add .golangci config
- Remove unused logrous package
- Load config from env file
- **config:** Remove zerolog
- **engine:** Ping docker daemon
- **proto:** Add ping request.
- **poller:** Add poller package
- **client:** Support gRPC and gRPC web protocol
- Add time sleep 1 second.
- **piepline:** Add runtime package.
- **proto:** Replace import path.
- **go:** Rename module name
- Add makefile
- **gRPC:** Register new runner
- **gRPC:** Add request interface client
- **gRPC:** Add update and update step for client
- **gRPC:** Fetch new stage data.
- **runtime:** Fetch build data
- **runtime:** Check error message
- **runner:** Support update log and task
- **runner:** Update fetch job request
- **runner:** Refactor register flow
- **register:** Update proto file.
- **runner:** Add default label ""self-hosted"
- **poller:** Add error retry limit
- **runner:** Update runner status if start program
- **runner:** Update runner status when start job
- **runner:** Add new register command
- **runner:** Add runner uuid in http header.
- **runner:** Remove client secret and add UUID in runner
- Remove tracing codes
- **runner:** Add new token in header
- **runner:** Cancel task if get the cancel from server
- **runner:** Wait workflow done before shutdown the service
- **poller:** Add metric to track the worker number
- **runner:** Remove update runner status
- Replace with code.gitea.io/bots-proto-go
- Upgrade act
- Update grpc router path

### Bug fixes

- Need more time to run task
- Update step result
- **runner:** Check task state field exist.
- Call RunDaemon
- Set NoSkipCheckout
- **poller:** Graceful shutdown
- Ignore stuck because of wrong token
- **poller:** Initial ready channel
- Set more image
- Limit workers
- Update container config
- **typo:** Use uppercase N for default input prompt ([#9](https://gitea.com/gitea/runner/pulls/9))
- Report job failure on error return ([#14](https://gitea.com/gitea/runner/pulls/14))

### Refactoring

- Add graceful shutdown signal notify func
- Docker engine with options
- **task:** Execute single task with gRPC data

### Documentation

- Add quickstart to readme
- **readme:** Update some header. ([#1](https://gitea.com/gitea/runner/pulls/1))
- Add Prerequisites ([#11](https://gitea.com/gitea/runner/pulls/11))

