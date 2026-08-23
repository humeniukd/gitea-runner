
## [3.3.0](https://gitea.com/gitea/runner/compare/v3.2.0...v3.3.0) (2026-08-21)

### Features

* add `log.job.dir` ([#1165](https://gitea.com/gitea/runner/issues/1165)) ([a0c4de7](https://gitea.com/gitea/runner/commits/a0c4de79f7331189c6ab414d9575bf1187edccd9))

### Enhancements

* add `runner.default_image` for jobs matching no label ([#1164](https://gitea.com/gitea/runner/issues/1164)) ([11ac12e](https://gitea.com/gitea/runner/commits/11ac12efa4dc3aff94d251cdce9c728486c7fa4a))

### Bug Fixes

* stop a failed step disabling the toolkit patch ([#1177](https://gitea.com/gitea/runner/issues/1177)) ([dfe979e](https://gitea.com/gitea/runner/commits/dfe979e1d039f6407499fca34171080077fc0d46))

## [3.2.0](https://gitea.com/gitea/runner/compare/v3.1.0...v3.2.0) (2026-08-19)

### Features

* add size-based cache eviction ([#1170](https://gitea.com/gitea/runner/issues/1170)) ([be90c01](https://gitea.com/gitea/runner/commits/be90c01468ee9937dddda2980c82414198ee5b8a))
* support `$/` prefix in action `uses:` ([#1150](https://gitea.com/gitea/runner/issues/1150)) ([9dd9204](https://gitea.com/gitea/runner/commits/9dd92049376248652a99b214e7456c398395256b))

### Enhancements

* add `runner.tool_cache_mode` and default it to `none` ([#1171](https://gitea.com/gitea/runner/issues/1171)) ([3f70822](https://gitea.com/gitea/runner/commits/3f70822458b435ca406fda96671289f2d846d339))
* show workflow_dispatch inputs in the "Set up job" section ([#1167](https://gitea.com/gitea/runner/issues/1167)) ([6c6a878](https://gitea.com/gitea/runner/commits/6c6a878403137f04c380ff6bf5bf6a5e323f83b5))

### Bug Fixes

* bound blocking calls and stop failing silently ([#1174](https://gitea.com/gitea/runner/issues/1174)) ([bc8161c](https://gitea.com/gitea/runner/commits/bc8161c673da8b6e72083940ffe2449513fd5d86))
* **cache:** build job URLs on the address its runner registered ([#1153](https://gitea.com/gitea/runner/issues/1153)) ([e178c03](https://gitea.com/gitea/runner/commits/e178c03adc10da456e9d310e4eea4285d04d8001))
* **deps:** update module github.com/go-git/go-git/v5 to v5.19.2 [security] ([#1156](https://gitea.com/gitea/runner/issues/1156)) ([b66433e](https://gitea.com/gitea/runner/commits/b66433e667d1499d9da66c2ca49e10ff61ee72a8)), references [#8203](https://gitea.com/gitea/runner/issues/8203) [#8203](https://gitea.com/gitea/runner/issues/8203)

### Refactoring

* move `act/model` and `act/exprparser` to actionslib ([#1143](https://gitea.com/gitea/runner/issues/1143)) ([a8dcd5b](https://gitea.com/gitea/runner/commits/a8dcd5b67c52059db85363d3b5b12534b1b7e8a3))

### Chores

* **deps:** update dependencies ([#1160](https://gitea.com/gitea/runner/issues/1160)) ([dbd9a89](https://gitea.com/gitea/runner/commits/dbd9a892f8310d4191ecca9ce29ce085ac70f974))
* drop AWS S3 release upload, keep Cloudflare R2 only ([#1169](https://gitea.com/gitea/runner/issues/1169)) ([f97680a](https://gitea.com/gitea/runner/commits/f97680a68d91d988cd9c96e0071e6ee0bee55cc2))
* revert docker 29.7.0 workaround ([#1155](https://gitea.com/gitea/runner/issues/1155)) ([da9b559](https://gitea.com/gitea/runner/commits/da9b559fb5d462f4523a4fa813934534b693ebae))

## [3.1.0](https://gitea.com/gitea/runner/compare/v3.0.2...v3.1.0) (2026-08-06)

### Features

* `config` command to edit config files ([#1140](https://gitea.com/gitea/runner/issues/1140)) ([8700adc](https://gitea.com/gitea/runner/commits/8700adc933ecdb5bcb06ab3972c3d58362242672))
* gate `set-env`/`add-path` and render annotation locations ([#1109](https://gitea.com/gitea/runner/issues/1109)) ([b70ff68](https://gitea.com/gitea/runner/commits/b70ff6893adc78879df84f6ae587e9a3f7d3f40b))
* wait for healthy services and fill the `job` context ([#1107](https://gitea.com/gitea/runner/issues/1107)) ([6854788](https://gitea.com/gitea/runner/commits/68547886a58094325220c309220abaef016ff4a1))

### Enhancements

* improve config, comment out values in example file ([#1145](https://gitea.com/gitea/runner/issues/1145)) ([1d6c6ff](https://gitea.com/gitea/runner/commits/1d6c6ffef968aa80fe5654374178920803f06d96))

### Bug Fixes

* evaluate each `${{ }}` part on its own ([#1146](https://gitea.com/gitea/runner/issues/1146)) ([20497aa](https://gitea.com/gitea/runner/commits/20497aaf4f029801486e4ed9ad34c9fb6f7d9488))
* trim whitespace from register inputs ([#1147](https://gitea.com/gitea/runner/issues/1147)) ([09b643b](https://gitea.com/gitea/runner/commits/09b643bc147915928876e50d03ebd08b97eb3632))

### Chores

* align go version handling with gitea ([#1144](https://gitea.com/gitea/runner/issues/1144)) ([70387cc](https://gitea.com/gitea/runner/commits/70387cca44268ba14374e0692f32a89ea4b9caaa))
* bump the module path to `/v3`, take the version from the VCS stamp ([#1136](https://gitea.com/gitea/runner/issues/1136)) ([4c2ab94](https://gitea.com/gitea/runner/commits/4c2ab943a8b85f940ce7e0a4cff23e852cee02c5))
* revert 4c2ab943a8b85f940ce7e0a4cff23e852cee02c5 ([#1148](https://gitea.com/gitea/runner/issues/1148)) ([24c13a1](https://gitea.com/gitea/runner/commits/24c13a1fd0cac7981c0bb6af4aeeedb03f44ff46)), references [#1136](https://gitea.com/gitea/runner/issues/1136) [#1136](https://gitea.com/gitea/runner/issues/1136)

## [3.0.2](https://gitea.com/gitea/runner/compare/v3.0.1...v3.0.2) (2026-08-03)

### Bug Fixes

* serve the whole results service from the cache server ([#1141](https://gitea.com/gitea/runner/issues/1141)) ([3618385](https://gitea.com/gitea/runner/commits/3618385b28a9fa51536860fce50db2da80641d9e))

### Chores

* **deps:** update dependencies ([#1138](https://gitea.com/gitea/runner/issues/1138)) ([55a625f](https://gitea.com/gitea/runner/commits/55a625f733111a107ec05749eccfca290a68208b))

## [3.0.1](https://gitea.com/gitea/runner/compare/v3.0.0...v3.0.1) (2026-08-02)

### Bug Fixes

* coerce every expression value kind to a string and drop the `format()` panic ([#1135](https://gitea.com/gitea/runner/issues/1135)) ([47366f8](https://gitea.com/gitea/runner/commits/47366f8f34169c6b577ac3b748100b318146647c))

### Documentation

* sync `AGENTS.md` with gitea/gitea ([#1134](https://gitea.com/gitea/runner/issues/1134)) ([b7aeda6](https://gitea.com/gitea/runner/commits/b7aeda6e7f77b2d4359a590d3463d0309871d488))

### Chores

* drop `Masterminds/semver` and `pkg/errors`, re-sync the `docker/cli` copies ([#1132](https://gitea.com/gitea/runner/issues/1132)) ([aced51b](https://gitea.com/gitea/runner/commits/aced51b4d58a0237a82df49cc05abf2923cda15e))

## [3.0.0](https://gitea.com/gitea/runner/compare/v2.3.0...v3.0.0) (2026-07-31)

### ⚠ BREAKING CHANGES

* add cache service v2, add toolkit patches ([#1110](https://gitea.com/gitea/runner/issues/1110))
* guard against two runner processes sharing one runner file ([#1099](https://gitea.com/gitea/runner/issues/1099))
* strip host-escape container options when privileged mode is disabled ([#1058](https://gitea.com/gitea/runner/issues/1058))

### Features

* add cache service v2, add toolkit patches ([#1110](https://gitea.com/gitea/runner/issues/1110)) ([47d5b5a](https://gitea.com/gitea/runner/commits/47d5b5ad03ab1b82cf565fe300796f7baa697bc8))
* add job hooks ([#1111](https://gitea.com/gitea/runner/issues/1111)) ([e6c7ba3](https://gitea.com/gitea/runner/commits/e6c7ba3a15b75adcee1ec7257a7c146c2b186609))
* mask secrets that reach the log in an encoded form ([#1108](https://gitea.com/gitea/runner/issues/1108)) ([68c6a5b](https://gitea.com/gitea/runner/commits/68c6a5b4f1399622bbf945798ae349de4abadfaf))
* propagate proxy variables to jobs, services and builds ([#1112](https://gitea.com/gitea/runner/issues/1112)) ([41c7221](https://gitea.com/gitea/runner/commits/41c72216bf67dc4170f4d54cdb599acad9d1be08))

### Bug Fixes

* allow relocating the tool cache and mounting over runner paths ([#1122](https://gitea.com/gitea/runner/issues/1122)) ([0192861](https://gitea.com/gitea/runner/commits/0192861155a075db357987e5111af33ff0738940))
* **deps:** update module github.com/ulikunitz/xz to v0.5.15 [security] ([#1127](https://gitea.com/gitea/runner/issues/1127)) ([2398d4a](https://gitea.com/gitea/runner/commits/2398d4a527e88eb4f3d89175a1d4ba37b4eb0e6d)), references [#8203](https://gitea.com/gitea/runner/issues/8203)
* fix panics, enable the `forcetypeassert` lint ([#1123](https://gitea.com/gitea/runner/issues/1123)) ([3f7fd16](https://gitea.com/gitea/runner/commits/3f7fd16ea1dfdd98b02dcdf3101daff8218b26a2))
* guard against two runner processes sharing one runner file ([#1099](https://gitea.com/gitea/runner/issues/1099)) ([0cd0e52](https://gitea.com/gitea/runner/commits/0cd0e52a24b941e8e3ca1dc35e081bb10ce5d293))
* resolve symlinked container paths before building tar entries ([#1130](https://gitea.com/gitea/runner/issues/1130)) ([34bfa19](https://gitea.com/gitea/runner/commits/34bfa19150229f8daee531a55ebd4d1e8eb585b8))
* stop leaking per-job docker networks ([#1124](https://gitea.com/gitea/runner/issues/1124)) ([b7a3bf9](https://gitea.com/gitea/runner/commits/b7a3bf98bcc084c82b3ecf9c44cfbb9693c089bd))
* strip host-escape container options when privileged mode is disabled ([#1058](https://gitea.com/gitea/runner/issues/1058)) ([14ec00b](https://gitea.com/gitea/runner/commits/14ec00b66eae973f7745b614d122e1e0eb7275a0))
* strip leading slash from mkdir tarball in CopyTarStream ([#1129](https://gitea.com/gitea/runner/issues/1129)) ([96d9f49](https://gitea.com/gitea/runner/commits/96d9f491db57a2a49e405de85443622145c820d8)), closes [#1128](https://gitea.com/gitea/runner/issues/1128)

### Tests

* speed up ([#1121](https://gitea.com/gitea/runner/issues/1121)) ([61f0cfa](https://gitea.com/gitea/runner/commits/61f0cfa951b11562b9050e48945fd1784b87fc47))

### Chores

* add renovate markers to `Makefile` tool packages ([#1126](https://gitea.com/gitea/runner/issues/1126)) ([da40378](https://gitea.com/gitea/runner/commits/da4037899a6221cd165ac87adabdad7a375cbe7b))
* scope `fmt-check` to Go sources ([#1125](https://gitea.com/gitea/runner/issues/1125)) ([e4fe49d](https://gitea.com/gitea/runner/commits/e4fe49dba402721440531e353c918b65c80209ac))

## [2.3.0](https://gitea.com/gitea/runner/compare/v2.2.0...v2.3.0) (2026-07-28)

### Features

* report runner name, environment, workspace and debug to jobs ([#1105](https://gitea.com/gitea/runner/issues/1105)) ([333eb17](https://gitea.com/gitea/runner/commits/333eb17d19340a7619b5a83ef01e1a31fa03d584))
* run pre-entrypoint and post-entrypoint of docker actions ([#1106](https://gitea.com/gitea/runner/issues/1106)) ([fc0e03e](https://gitea.com/gitea/runner/commits/fc0e03e5a92bd234e38dac6773f1ad7749658c37))
* support `--platform` and `--pull` in `container.options` ([#1104](https://gitea.com/gitea/runner/issues/1104)) ([b4a64b9](https://gitea.com/gitea/runner/commits/b4a64b97dde557de7451c3e3ee9915427138a4ba))
* support reading `cache.external_secret` from a file ([#1100](https://gitea.com/gitea/runner/issues/1100)) ([c9c4957](https://gitea.com/gitea/runner/commits/c9c4957e3833a54c2c97fb50d4553bca755b8292))

### Bug Fixes

* clean service containers after failed job setup ([#1066](https://gitea.com/gitea/runner/issues/1066)) ([26f9fb1](https://gitea.com/gitea/runner/commits/26f9fb12af1063f2e2453eb291bb8ee747fe77d5))
* **deps:** update module github.com/prometheus/client_golang to v1.24.0 ([#1113](https://gitea.com/gitea/runner/issues/1113)) ([8c519ce](https://gitea.com/gitea/runner/commits/8c519ce31879a3490674d5bb26c88a88bd25455f))
* **deps:** update module go.opentelemetry.io/otel to v1.44.0 [security] ([#1115](https://gitea.com/gitea/runner/issues/1115)) ([de43c84](https://gitea.com/gitea/runner/commits/de43c842032ee150d00b1b1f0b0cbdb6d45e8e31))
* **deps:** update module golang.org/x/net to v0.56.0 [security] ([#1116](https://gitea.com/gitea/runner/issues/1116)) ([3c5ef17](https://gitea.com/gitea/runner/commits/3c5ef1721a6e1c97852016671bcf118ab543c6d3))
* escape command data the runner writes itself ([#1120](https://gitea.com/gitea/runner/issues/1120)) ([c3b39e0](https://gitea.com/gitea/runner/commits/c3b39e0d99ff69dbbee0b7ca4ec823038d11c636))

### Continuous Integration

* mirror release artifacts to Cloudflare R2 ([#1114](https://gitea.com/gitea/runner/issues/1114)) ([94ab020](https://gitea.com/gitea/runner/commits/94ab020204340a6df03404c1e46802c141bd5f84))

### Chores

* **deps:** pin dependencies ([#1117](https://gitea.com/gitea/runner/issues/1117)) ([78a74f7](https://gitea.com/gitea/runner/commits/78a74f78f819d70d4faaf284e86f3017f81926ce))
* **deps:** update actions/setup-go action to v7 ([#1102](https://gitea.com/gitea/runner/issues/1102)) ([0fd8602](https://gitea.com/gitea/runner/commits/0fd8602ac3a785bf88e6b8e53dc1c947568ed530))
* **deps:** update docker docker tag to v29.6.2 ([#1101](https://gitea.com/gitea/runner/issues/1101)) ([b1a02cd](https://gitea.com/gitea/runner/commits/b1a02cdd5da16cf4c4fded839178e8166ec86b09))

## [2.2.0](https://gitea.com/gitea/runner/compare/v2.1.0...v2.2.0) (2026-07-22)

### Features

* add runner health admission checks ([#1090](https://gitea.com/gitea/runner/issues/1090)) ([c43cbe8](https://gitea.com/gitea/runner/commits/c43cbe87ca95efd12eee221b5640646f24e6a511))

### Enhancements

* report a GitHub-style "Set up job" section ([#1089](https://gitea.com/gitea/runner/issues/1089)) ([8af385d](https://gitea.com/gitea/runner/commits/8af385d14729ae4f7ffb8a8573454f0fa6001cd6))

### Bug Fixes

* classify a cancelled step as an interruption, not a failure ([#1095](https://gitea.com/gitea/runner/issues/1095)) ([0e8896c](https://gitea.com/gitea/runner/commits/0e8896c52a97fd28debe34a32ddb5bde49894fa4))
* **deps:** update module github.com/docker/cli to v29.6.2+incompatible ([#1096](https://gitea.com/gitea/runner/issues/1096)) ([068afc3](https://gitea.com/gitea/runner/commits/068afc399654a109c340d116f58e71eb45d38a07))
* **deps:** update module github.com/mattn/go-isatty to v0.0.23 ([#1097](https://gitea.com/gitea/runner/issues/1097)) ([46f22c7](https://gitea.com/gitea/runner/commits/46f22c78d23534d1481bccda43b030253a03cfd4))
* guard status-check functions against a nil job context ([#1092](https://gitea.com/gitea/runner/issues/1092)) ([aa7a29a](https://gitea.com/gitea/runner/commits/aa7a29a157f99c39bceaa6dc4ffe067f5cfa33c0))
* repair free-disk-space build on FreeBSD ([#1098](https://gitea.com/gitea/runner/issues/1098)) ([6133d64](https://gitea.com/gitea/runner/commits/6133d64270edfa7c211d2f23c1677af95bf93290))
* stop host-mode jobs from leaking processes on Windows ([#1080](https://gitea.com/gitea/runner/issues/1080)) ([7bec310](https://gitea.com/gitea/runner/commits/7bec31000208c1f60e4a41862db52d907db782ff))
* stop racing the daemon when removing containers ([#1093](https://gitea.com/gitea/runner/issues/1093)) ([89467c9](https://gitea.com/gitea/runner/commits/89467c9dd07d017f05797e6a23d9fcc937a7e823))

### Chores

* **deps:** update actions/setup-node action to v7 ([#1094](https://gitea.com/gitea/runner/issues/1094)) ([0c08b0f](https://gitea.com/gitea/runner/commits/0c08b0f2da5f3fcaba36d687931f4b9f3db9068b))

## [2.1.0](https://gitea.com/gitea/runner/compare/v2.0.1...v2.1.0) (2026-07-16)

### Features

* honour GITEA_RUNNER_LABELS on daemon start and accept labels containing a colon ([#1085](https://gitea.com/gitea/runner/issues/1085)) ([58c5eb8](https://gitea.com/gitea/runner/commits/58c5eb8d21d28d8bb2de7be028f2e3bb1f1d1c7c)), closes [#648](https://gitea.com/gitea/runner/issues/648) [#656](https://gitea.com/gitea/runner/issues/656) [#664](https://gitea.com/gitea/runner/issues/664)

### Enhancements

* add `--token-file` flag to `register` command ([#1076](https://gitea.com/gitea/runner/issues/1076)) ([be9b450](https://gitea.com/gitea/runner/commits/be9b4502d61c55e7ac8ba60292d96333b4d61a54)), references [#362](https://gitea.com/gitea/runner/issues/362)

### Bug Fixes

* accept natively typed boolean workflow inputs ([#1087](https://gitea.com/gitea/runner/issues/1087)) ([16357a3](https://gitea.com/gitea/runner/commits/16357a34b2c6c143f58f797159870ddf2c61ebd0)), closes [#1082](https://gitea.com/gitea/runner/issues/1082)
* **deps:** update module golang.org/x/term to v0.45.0 ([#1081](https://gitea.com/gitea/runner/issues/1081)) ([554b3b7](https://gitea.com/gitea/runner/commits/554b3b767144902dab9f78a8d7742086f83bedd9))
* **deps:** update module golang.org/x/text to v0.40.0 ([#1091](https://gitea.com/gitea/runner/issues/1091)) ([ad96733](https://gitea.com/gitea/runner/commits/ad967330a8788c9b8ab723abbc1a86d53c3bc5e6))
* drop action outputs whose value exceeds the size limit ([#1070](https://gitea.com/gitea/runner/issues/1070)) ([d53538a](https://gitea.com/gitea/runner/commits/d53538ac38ab07e0e17e2a85aff0fc04e4f5ce43))
* ignore blank lines and decode UTF-16 in the runner env files ([#1084](https://gitea.com/gitea/runner/issues/1084)) ([6017700](https://gitea.com/gitea/runner/commits/60177008a5920899c871dc960eb1baaed7d87a24)), closes [#496](https://gitea.com/gitea/runner/issues/496) [#552](https://gitea.com/gitea/runner/issues/552)
* Minor fixes ([#1075](https://gitea.com/gitea/runner/issues/1075)) ([65756d6](https://gitea.com/gitea/runner/commits/65756d60b364797c449d8335bf367686358a6811)), closes [#625](https://gitea.com/gitea/runner/issues/625) [#804](https://gitea.com/gitea/runner/issues/804) [#571](https://gitea.com/gitea/runner/issues/571) [#650](https://gitea.com/gitea/runner/issues/650) [#766](https://gitea.com/gitea/runner/issues/766) [#256](https://gitea.com/gitea/runner/issues/256) [#617](https://gitea.com/gitea/runner/issues/617) [#635](https://gitea.com/gitea/runner/issues/635) [#1005](https://gitea.com/gitea/runner/issues/1005) [#448](https://gitea.com/gitea/runner/issues/448) [#209](https://gitea.com/gitea/runner/issues/209) [#757](https://gitea.com/gitea/runner/issues/757) [#474](https://gitea.com/gitea/runner/issues/474) [#776](https://gitea.com/gitea/runner/issues/776)
* stop service containers from clobbering the job container's credentials ([#1083](https://gitea.com/gitea/runner/issues/1083)) ([7e7e3ef](https://gitea.com/gitea/runner/commits/7e7e3ef1a6f725e210bba203166757f9a07e5b81)), closes [#835](https://gitea.com/gitea/runner/issues/835) [#643](https://gitea.com/gitea/runner/issues/643)

### Documentation

* explain labels, the docker image cache volume, and the dind-rootless UID ([#1086](https://gitea.com/gitea/runner/issues/1086)) ([d6882b3](https://gitea.com/gitea/runner/commits/d6882b3df5d87bd682be4abf7b4df67b62551f50)), closes [#106](https://gitea.com/gitea/runner/issues/106) [#570](https://gitea.com/gitea/runner/issues/570) [#627](https://gitea.com/gitea/runner/issues/627)

## [2.0.1](https://gitea.com/gitea/runner/compare/v2.0.0...v2.0.1) (2026-07-09)

### Bug Fixes

* attach task token when cloning actions from self-hosted instance on a different host ([#1056](https://gitea.com/gitea/runner/issues/1056)) ([745b0ab](https://gitea.com/gitea/runner/commits/745b0ab6e4953ed77e23beca2d32077371723954))
* **deps:** update module github.com/docker/cli to v29.6.1+incompatible ([#1060](https://gitea.com/gitea/runner/issues/1060)) ([eba33e1](https://gitea.com/gitea/runner/commits/eba33e178d7450239a5960e8e42db36282324e21))
* **deps:** update module golang.org/x/sys to v0.47.0 ([#1073](https://gitea.com/gitea/runner/issues/1073)) ([b12d02c](https://gitea.com/gitea/runner/commits/b12d02c25f9e8eeb95b8ffaa74ed9cb50f9b7865))
* install `nftables` in dind images to silence `nft` cleanup errors ([#1064](https://gitea.com/gitea/runner/issues/1064)) ([0ee4643](https://gitea.com/gitea/runner/commits/0ee4643d4ac503665424c7bf531d9523503b5a93))
* prevent exponential growth of RunContext masks in composite actions ([#1059](https://gitea.com/gitea/runner/issues/1059)) ([eeb479e](https://gitea.com/gitea/runner/commits/eeb479ea89a54e42a8ede6da9b00e5c1da98fca7)), references [#1001](https://gitea.com/gitea/runner/issues/1001) [#883](https://gitea.com/gitea/runner/issues/883)
* skip service containers with an empty image ([#1074](https://gitea.com/gitea/runner/issues/1074)) ([1d74ae6](https://gitea.com/gitea/runner/commits/1d74ae636a07ec3df70f3a25d8bc05b890824a20))

### Documentation

* clarify cache reachability for dockerized runners ([#1069](https://gitea.com/gitea/runner/issues/1069)) ([f2e0cf9](https://gitea.com/gitea/runner/commits/f2e0cf9131e2b956229a923c48d24f41835254b1))

### Tests

* Enhance Coverage + CI ([#1055](https://gitea.com/gitea/runner/issues/1055)) ([3396021](https://gitea.com/gitea/runner/commits/3396021e0f4d78f9975cca9fd88ee72341bc83fd))

### Chores

* **deps:** update docker docker tag to v29.6.1 ([#1063](https://gitea.com/gitea/runner/issues/1063)) ([e774003](https://gitea.com/gitea/runner/commits/e774003c18812473ce60e2b68f88c38679ec3b47))

## [2.0.0](https://gitea.com/gitea/runner/compare/v1.0.8...v2.0.0) (2026-06-29)

### ⚠ BREAKING CHANGES

* stop implicitly using `DOCKER_USERNAME`/`DOCKER_PASSWORD` secrets for image pulls ([#1007](https://gitea.com/gitea/runner/issues/1007))

### Features

* Add optional `runner.post_task_script` hook after task cleanup ([#1026](https://gitea.com/gitea/runner/issues/1026)) ([0077179](https://gitea.com/gitea/runner/commits/007717956aec00246b672a0f475defb77a303279))
* complete runner-side cancellation handling ([#1016](https://gitea.com/gitea/runner/issues/1016)) ([822af50](https://gitea.com/gitea/runner/commits/822af5029f28f8f04eaf52cae5c75d57a4c15bcd)), references [#825](https://gitea.com/gitea/runner/issues/825) [#825](https://gitea.com/gitea/runner/issues/825)
* Enable `jobs.<job_id>.timeout-minutes` and `jobs.<job_id>.continue-on-error` ([#1032](https://gitea.com/gitea/runner/issues/1032)) ([6bdcb54](https://gitea.com/gitea/runner/commits/6bdcb548282d464eed228183b9e3e6bbb1ab6e27))
* ipv6 options for network container creation ([#1029](https://gitea.com/gitea/runner/issues/1029)) ([2963716](https://gitea.com/gitea/runner/commits/2963716953f6d28825a9e8833d7312d4d1469b40))
* shallow clone action repositories ([#1053](https://gitea.com/gitea/runner/issues/1053)) ([99bc50d](https://gitea.com/gitea/runner/commits/99bc50d538e69018371b461d73bfa6027bb456d7))
* Support `ssh://` action URLs ([#1035](https://gitea.com/gitea/runner/issues/1035)) ([5f0636f](https://gitea.com/gitea/runner/commits/5f0636faaddb83fc3230a9ed567476e198f92e0d)), closes [#841](https://gitea.com/gitea/runner/issues/841)
* upload job summary when supported ([#917](https://gitea.com/gitea/runner/issues/917)) ([53c4db6](https://gitea.com/gitea/runner/commits/53c4db6a4bed89ba19e2491c3723d2a6b368ba77)), closes [go-gitea/gitea#23721](https://gitea.com/go-gitea/gitea/issues/23721)

### Bug Fixes

* **cleanup:** kill Unix step process group on cancel to avoid hang ([#1025](https://gitea.com/gitea/runner/issues/1025)) ([3996d6d](https://gitea.com/gitea/runner/commits/3996d6d032b4fd1eb4579c677dfb6ae18ce2235d)), references [#1011](https://gitea.com/gitea/runner/issues/1011)
* composite nested uses clone token ([#1041](https://gitea.com/gitea/runner/issues/1041)) ([b7f6b6d](https://gitea.com/gitea/runner/commits/b7f6b6d90a181d3f78db88c6a44e77bc39786993))
* **deps:** update module code.gitea.io/actions-proto-go to gitea.dev/actions-proto-go v0.5.0 ([#1009](https://gitea.com/gitea/runner/issues/1009)) ([984b47c](https://gitea.com/gitea/runner/commits/984b47c716755de10fa54180571808be4651153e))
* **deps:** update module github.com/docker/cli to v29.5.3+incompatible ([#1018](https://gitea.com/gitea/runner/issues/1018)) ([8ad84cd](https://gitea.com/gitea/runner/commits/8ad84cd96a00580e3136cc736be7ea02dd8b08f3))
* **deps:** update module github.com/docker/cli to v29.6.0+incompatible ([#1045](https://gitea.com/gitea/runner/issues/1045)) ([bd41a36](https://gitea.com/gitea/runner/commits/bd41a367fe9ce23014c8df5ef5a46dd7f4451188))
* **deps:** update module github.com/moby/moby/client to v0.5.0 ([#1049](https://gitea.com/gitea/runner/issues/1049)) ([4e7fd1c](https://gitea.com/gitea/runner/commits/4e7fd1c68a0d32a307137669af085c6f5f081650))
* **deps:** update module github.com/opencontainers/selinux to v1.15.1 ([#1017](https://gitea.com/gitea/runner/issues/1017)) ([443b0e3](https://gitea.com/gitea/runner/commits/443b0e336cce21e511cc95491f27770995a8039e))
* **deps:** update module go.etcd.io/bbolt to v1.5.0 ([#1040](https://gitea.com/gitea/runner/issues/1040)) ([d3b3519](https://gitea.com/gitea/runner/commits/d3b3519dea2dd1d85c861e1ae57b94d077f62d42))
* **deps:** update module golang.org/x/sys to v0.45.0 ([#1012](https://gitea.com/gitea/runner/issues/1012)) ([ff7d9ca](https://gitea.com/gitea/runner/commits/ff7d9ca8d00aa0498db21502f1ae33b9439e4e44))
* **deps:** update module golang.org/x/sys to v0.46.0 ([#1019](https://gitea.com/gitea/runner/issues/1019)) ([e583b07](https://gitea.com/gitea/runner/commits/e583b0706b70b58a1182c0a3516c0e9c4b1777bf))
* **deps:** update module golang.org/x/term to v0.44.0 ([#1031](https://gitea.com/gitea/runner/issues/1031)) ([56979e6](https://gitea.com/gitea/runner/commits/56979e6ab88c968f148fa03c384c8e2f96a52c8a))
* do not update cached actions with stale origin URL ([#1014](https://gitea.com/gitea/runner/issues/1014)) ([1073c8b](https://gitea.com/gitea/runner/commits/1073c8bfec2093e87d784eb4e16d0edcc0874dcf)), closes [#1010](https://gitea.com/gitea/runner/issues/1010)
* guard SetJobError against missing job-error container ([#1050](https://gitea.com/gitea/runner/issues/1050)) ([8f72c60](https://gitea.com/gitea/runner/commits/8f72c60afa9b3ccdaf768d3111722ab594599fc0))
* **host:** bound host-environment cleanup and reclaim leaked scratch dirs ([#1024](https://gitea.com/gitea/runner/issues/1024)) ([33e6d1d](https://gitea.com/gitea/runner/commits/33e6d1d8ff9f96495df5331c32bf4868fbd1cde3)), closes [#1023](https://gitea.com/gitea/runner/issues/1023)
* Interpolate job container.volume ([#1036](https://gitea.com/gitea/runner/issues/1036)) ([df0370f](https://gitea.com/gitea/runner/commits/df0370f8bfd3563cf207098ad4c2cde4d98d371b))
* namespace local docker action image tags per repository ([#1051](https://gitea.com/gitea/runner/issues/1051)) ([cdcea87](https://gitea.com/gitea/runner/commits/cdcea87a459394efb082b223b85c1b4f1ed2c25e)), closes [#1039](https://gitea.com/gitea/runner/issues/1039)
* prevent loss of step log output at end of step ([#1028](https://gitea.com/gitea/runner/issues/1028)) ([205af7c](https://gitea.com/gitea/runner/commits/205af7cd0118c60587848d89b2569498093da94a))
* run always()/cancelled() and post steps correctly on cancellation ([#1043](https://gitea.com/gitea/runner/issues/1043)) ([e22d3fa](https://gitea.com/gitea/runner/commits/e22d3fa26383f997c8aa5080c2d3ff9a5f397777)), closes [#657](https://gitea.com/gitea/runner/issues/657), references [#1016](https://gitea.com/gitea/runner/issues/1016) [#1016](https://gitea.com/gitea/runner/issues/1016)
* stop implicitly using `DOCKER_USERNAME`/`DOCKER_PASSWORD` secrets for image pulls ([#1007](https://gitea.com/gitea/runner/issues/1007)) ([0a2f282](https://gitea.com/gitea/runner/commits/0a2f28244dcc31a63bedeebdcc3518d4466692ca)), closes [#386](https://gitea.com/gitea/runner/issues/386)

### Build

* cover Windows Go files ([#1052](https://gitea.com/gitea/runner/issues/1052)) ([3c4bcf3](https://gitea.com/gitea/runner/commits/3c4bcf3ebf56b75f36122cea26a22d5d7bfe1d00))

### Documentation

* **docker-images:** Update docs ([#1020](https://gitea.com/gitea/runner/issues/1020)) ([355289b](https://gitea.com/gitea/runner/commits/355289bc54dacc9c1a871ba509ec8f644dc1aef0))
* Improve the documentation for cache ([#1034](https://gitea.com/gitea/runner/issues/1034)) ([4997f33](https://gitea.com/gitea/runner/commits/4997f33b5f41a1e067fd61ccc5ad723f5158e8d6))

### Chores

* **deps:** update actions/checkout action to v7 ([#1042](https://gitea.com/gitea/runner/issues/1042)) ([40e0213](https://gitea.com/gitea/runner/commits/40e021309aec9e11606038265719d9d0d4e4482d))
* **deps:** update alpine docker tag to v3.24 ([#1030](https://gitea.com/gitea/runner/issues/1030)) ([bf99e6a](https://gitea.com/gitea/runner/commits/bf99e6a758014cc3dc02ac9a1b3d54d3065278bf))
* **deps:** update docker docker tag to v29.5.3 ([#1021](https://gitea.com/gitea/runner/issues/1021)) ([526c46b](https://gitea.com/gitea/runner/commits/526c46b485fe5fc6f618b29bf6516db7fe2e24e3))
* **deps:** update docker docker tag to v29.6.0 ([#1044](https://gitea.com/gitea/runner/issues/1044)) ([c566013](https://gitea.com/gitea/runner/commits/c566013db41ec3d77e5b3cf0c9f11fa86949b326))
* **deps:** update golang.org/x/crypto to v0.52.0 ([#1027](https://gitea.com/gitea/runner/issues/1027)) ([740a3d4](https://gitea.com/gitea/runner/commits/740a3d4db4c2900ccaa67e28eccde7672d4cb4f3))

## [1.0.8](https://gitea.com/gitea/runner/compare/v1.0.7...v1.0.8) (2026-06-02)

### Bug Fixes

* **cleanup:** kill Windows step process tree on cancel to avoid hang ([#1011](https://gitea.com/gitea/runner/issues/1011)) ([c749e52](https://gitea.com/gitea/runner/commits/c749e52bb712bf8029bc8d9193297e32740305c6))

## [1.0.7](https://gitea.com/gitea/runner/compare/v1.0.6...v1.0.7) (2026-05-29)

### Bug Fixes

* **container:** re-validate cached container id before reuse ([#1003](https://gitea.com/gitea/runner/issues/1003)) ([f17b6b9](https://gitea.com/gitea/runner/commits/f17b6b9fc3cdf32b046e96d6aa1e85398d91c9a4))
* deliver cancel ack and reap leftover Windows job processes ([#996](https://gitea.com/gitea/runner/issues/996)) ([0b9f251](https://gitea.com/gitea/runner/commits/0b9f251b6abb30d5f292a49cfe0c611f7c26d857))
* matrix-job data races + outputs, leaner offline test suite ([#994](https://gitea.com/gitea/runner/issues/994)) ([270ea41](https://gitea.com/gitea/runner/commits/270ea4123275fd3aaa685fc4584ede4aa1de8281))
* **reporter:** respect configured log level for job log forwarding ([#989](https://gitea.com/gitea/runner/issues/989)) ([273f6b4](https://gitea.com/gitea/runner/commits/273f6b4247862aabfdbf3c9f8c87d2717b7c2944)), references [#677](https://gitea.com/gitea/runner/issues/677) [#677](https://gitea.com/gitea/runner/issues/677) [#409](https://gitea.com/gitea/runner/issues/409)
* restore global docker config dir and socket env in tests ([#1004](https://gitea.com/gitea/runner/issues/1004)) ([abec931](https://gitea.com/gitea/runner/commits/abec931d98ba969c671a71b38e0cc1b60c15671d))
* support multiline secret masking ([#1001](https://gitea.com/gitea/runner/issues/1001)) ([c7c4bd6](https://gitea.com/gitea/runner/commits/c7c4bd600ab2b36658a1f39624e78a6638f02bf5)), closes [#998](https://gitea.com/gitea/runner/issues/998)

## [1.0.6](https://gitea.com/gitea/runner/compare/v1.0.5...v1.0.6) (2026-05-22)

### Bug Fixes

* clean up job network and container when container start fails ([#986](https://gitea.com/gitea/runner/issues/986)) ([b30204a](https://gitea.com/gitea/runner/commits/b30204aa9432b0bc1da46fcfca8438b76a68e222))
* **deps:** update module github.com/opencontainers/selinux to v1.15.0 ([#990](https://gitea.com/gitea/runner/issues/990)) ([47ee454](https://gitea.com/gitea/runner/commits/47ee45412a969de0144794cdc8265f4da09c0325))

### Continuous Integration

* add PR title linting against Conventional Commits ([#988](https://gitea.com/gitea/runner/issues/988)) ([d6fbe75](https://gitea.com/gitea/runner/commits/d6fbe757218254ec764a4bfa573303b1292f861c))

### Tests

* make TestRunEvent integration suite runnable locally ([#987](https://gitea.com/gitea/runner/issues/987)) ([0e0c54b](https://gitea.com/gitea/runner/commits/0e0c54b272a2e31bd5f5ab34b5e53de42bb8e1a1))

### Chores

* **deps:** update actions/setup-node action to v6 ([#991](https://gitea.com/gitea/runner/issues/991)) ([1c62c06](https://gitea.com/gitea/runner/commits/1c62c0635f38c8663137a2267cbe2ca7ade9182c))
* pin Docker base images to explicit versions ([#992](https://gitea.com/gitea/runner/issues/992)) ([38b69bb](https://gitea.com/gitea/runner/commits/38b69bb214fcd1070ae19995833b48f9611d7aa5))

## [1.0.5](https://gitea.com/gitea/runner/compare/v1.0.4...v1.0.5) (2026-05-21)

### Features

* add cache.offline_mode to reuse cached actions ([#966](https://gitea.com/gitea/runner/issues/966)) ([2208e7e](https://gitea.com/gitea/runner/commits/2208e7ec63c3266fd3c566810d697ffa49bb10be))

### Bug Fixes

* **deps:** update module connectrpc.com/connect to v1.20.0 ([#985](https://gitea.com/gitea/runner/issues/985)) ([7b5ebe9](https://gitea.com/gitea/runner/commits/7b5ebe9618c34ecd0bee180268837f619c889ba6))
* **deps:** update module github.com/docker/cli to v29.5.1+incompatible ([#979](https://gitea.com/gitea/runner/issues/979)) ([10475db](https://gitea.com/gitea/runner/commits/10475db58a83ba30a8975c0bb96fabf93e526521))
* **deps:** update module github.com/go-git/go-git/v5 to v5.19.1 ([#980](https://gitea.com/gitea/runner/issues/980)) ([9e738c2](https://gitea.com/gitea/runner/commits/9e738c203c1662167d4ff9376e6fb31e83d13b5e))

## [1.0.4](https://gitea.com/gitea/runner/compare/v1.0.3...v1.0.4) (2026-05-18)

### Features

* make pseudo-TTY allocation opt-in ([#961](https://gitea.com/gitea/runner/issues/961)) ([3c5f03f](https://gitea.com/gitea/runner/commits/3c5f03ff8f38d5f030e45881a0ea1d180589d997)), closes [#956](https://gitea.com/gitea/runner/issues/956)

### Bug Fixes

* **deps:** bump `docker` deps, switch to `moby/moby` ([#943](https://gitea.com/gitea/runner/issues/943)) ([32bed52](https://gitea.com/gitea/runner/commits/32bed526861bce58430e81b627c0af464b678906))
* **deps:** update module github.com/docker/cli to v29.5.0+incompatible ([#969](https://gitea.com/gitea/runner/issues/969)) ([8d7cf48](https://gitea.com/gitea/runner/commits/8d7cf48a6f4526b55a6dd0af1846ff55e08e431d))
* **parse_env_file:** support env-file lines larger than 64 KiB ([#974](https://gitea.com/gitea/runner/issues/974)) ([cf7e29c](https://gitea.com/gitea/runner/commits/cf7e29c10d2444f8074457a6399d86bdca688906))
* respect proxy env vars in runner client ([#962](https://gitea.com/gitea/runner/issues/962)) ([a7e972d](https://gitea.com/gitea/runner/commits/a7e972d8deaaf1d77b5b7a2fa6d19661afe1ecc7)), closes [#957](https://gitea.com/gitea/runner/issues/957)
* Return if executors length is zero in ParallelExecutor ([#960](https://gitea.com/gitea/runner/issues/960)) ([5464d33](https://gitea.com/gitea/runner/commits/5464d33eef3caa5897e08ecef548c1e9622995d1))

### Chores

* **deps:** bump `retry-go`, `golangci-lint`, `govulncheck` ([#965](https://gitea.com/gitea/runner/issues/965)) ([dda5841](https://gitea.com/gitea/runner/commits/dda5841af88866f48fba86530f3fe927fbe26106))
* **deps:** update workflow dependencies (major) ([#967](https://gitea.com/gitea/runner/issues/967)) ([f23605c](https://gitea.com/gitea/runner/commits/f23605c614d70e2231b3e6d683f76eeaf76c856d))
* **deps:** update workflow dependencies (major) ([#968](https://gitea.com/gitea/runner/issues/968)) ([880e975](https://gitea.com/gitea/runner/commits/880e9755d96181897b24ca643655bc1e16419c97))

## [1.0.3](https://gitea.com/gitea/runner/compare/v1.0.2...v1.0.3) (2026-05-12)

### Bug Fixes

* **artifactcache:** preserve cache key case to stop redundant uploads ([#947](https://gitea.com/gitea/runner/issues/947)) ([ef6ca95](https://gitea.com/gitea/runner/commits/ef6ca957b5acea963c12628116c3fcf52ffbb2d2))
* **deps:** update mergo to v1.0.2 (now `dario.cat/mergo`) ([#954](https://gitea.com/gitea/runner/issues/954)) ([1e3ab0c](https://gitea.com/gitea/runner/commits/1e3ab0c40ac429eff95642d65ac4402679aa3b03))
* **deps:** update module github.com/opencontainers/selinux to v1.14.1 ([#955](https://gitea.com/gitea/runner/issues/955)) ([a1f13cb](https://gitea.com/gitea/runner/commits/a1f13cb97056d17e93580b5b33108a9a86032841))
* **deps:** update module golang.org/x/term to v0.43.0 ([#948](https://gitea.com/gitea/runner/issues/948)) ([8088df5](https://gitea.com/gitea/runner/commits/8088df52b94ba6b19703e896ffc79a366372a887))
* ensure `dbfs_data` is cleaned up after task completion ([#952](https://gitea.com/gitea/runner/issues/952)) ([295eecb](https://gitea.com/gitea/runner/commits/295eecb9afd1f784866b7d9d79102406c812e393)), closes [#950](https://gitea.com/gitea/runner/issues/950), references [#819](https://gitea.com/gitea/runner/issues/819) [#819](https://gitea.com/gitea/runner/issues/819)
* isolate per-task runner envs ([#959](https://gitea.com/gitea/runner/issues/959)) ([763b38e](https://gitea.com/gitea/runner/commits/763b38ece3fd8b6959d79336efa0fea75701ba74)), closes [#958](https://gitea.com/gitea/runner/issues/958)

## [1.0.2](https://gitea.com/gitea/runner/compare/v1.0.1...v1.0.2) (2026-05-08)

### Features

* remove emojis from runner logging, add `Starting job container` group ([#940](https://gitea.com/gitea/runner/issues/940)) ([7564364](https://gitea.com/gitea/runner/commits/75643645f07e49f588852f4d4add96cf9e133c07)), closes [#935](https://gitea.com/gitea/runner/issues/935)

### Bug Fixes

* overwrite read-only files when copying action directories ([#942](https://gitea.com/gitea/runner/issues/942)) ([3ea7d39](https://gitea.com/gitea/runner/commits/3ea7d39690d789aa2627573e7614bf88d81dc3cf))
* serialize action-cache reads to prevent worktree race ([#938](https://gitea.com/gitea/runner/issues/938)) ([cce8543](https://gitea.com/gitea/runner/commits/cce8543d06c203da92e06ac72988d54dda2f1f2c))

## [1.0.1](https://gitea.com/gitea/runner/compare/v1.0.0...v1.0.1) (2026-05-07)

### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.19.0 ([#934](https://gitea.com/gitea/runner/issues/934)) ([dff63b3](https://gitea.com/gitea/runner/commits/dff63b3ecc4320d340d3b5511f91a65eb2e2e6c2))
* **deps:** update module github.com/opencontainers/selinux to v1.14.0 ([#928](https://gitea.com/gitea/runner/issues/928)) ([a5d9fe9](https://gitea.com/gitea/runner/commits/a5d9fe965120c0b398a6349005dd9a73eb6e511c))
* re-fetch cached reusable workflow on every run ([#930](https://gitea.com/gitea/runner/issues/930)) ([5e59402](https://gitea.com/gitea/runner/commits/5e59402fb2ad128c3f3008d4c41c7d3d03ed3f19))

### Tests

* clean up dead/stale fixtures and bump test container images ([#932](https://gitea.com/gitea/runner/issues/932)) ([d607f3b](https://gitea.com/gitea/runner/commits/d607f3b342e351197691a08795189972071d6374))

### Chores

* **deps:** update docker docker tag to v29 ([#924](https://gitea.com/gitea/runner/issues/924)) ([dfeb463](https://gitea.com/gitea/runner/commits/dfeb463904e768bb33086563d620a1925c278618))

## [1.0.0](https://gitea.com/gitea/runner/compare/v0.6.1...v1.0.0) (2026-05-05)

### Features

* add startup janitor for stale bind-workdir task workspaces ([#870](https://gitea.com/gitea/runner/issues/870)) ([2a4d56c](https://gitea.com/gitea/runner/commits/2a4d56c6501d7fa73a064b91075668466b494890))

### Bug Fixes

* **deps:** update module github.com/rhysd/actionlint to v1.7.12 ([#873](https://gitea.com/gitea/runner/issues/873)) ([1d6e787](https://gitea.com/gitea/runner/commits/1d6e7879c82d4161e3b4f0a80290127467446538))
* **deps:** update module golang.org/x/term to v0.42.0 ([#920](https://gitea.com/gitea/runner/issues/920)) ([d143423](https://gitea.com/gitea/runner/commits/d1434237c2d1e30cdf61c96c740dc8b02cfd211f))
* **host:** correct host workspace cleanup on Windows ([#883](https://gitea.com/gitea/runner/issues/883)) ([a22119c](https://gitea.com/gitea/runner/commits/a22119cf883f5dc89c00fc226d00056a936b3c7a))

### Continuous Integration

* Fix triggers ([#882](https://gitea.com/gitea/runner/issues/882)) ([c45a4e6](https://gitea.com/gitea/runner/commits/c45a4e6d32f2a713a30fb226ea2e8c79d524421e))

### Documentation

* Update docs ([#872](https://gitea.com/gitea/runner/issues/872)) ([8e6b3be](https://gitea.com/gitea/runner/commits/8e6b3be96af083be7ff753019f12e080ba60c05b))

### Chores

* **deps:** update actions/hello-world-docker-action action to v2 ([#921](https://gitea.com/gitea/runner/issues/921)) ([35c65e2](https://gitea.com/gitea/runner/commits/35c65e2b14f8e88d9172d51f1c1692ce1cf9ee6b))
* **deps:** update crazy-max/ghaction-import-gpg action to v7 ([#923](https://gitea.com/gitea/runner/issues/923)) ([b68ecf2](https://gitea.com/gitea/runner/commits/b68ecf2580b0be3a577b803fc19925fc57b2a2cb))
* **deps:** update dependency @actions/core to v1.11.1 ([#880](https://gitea.com/gitea/runner/issues/880)) ([b1c873a](https://gitea.com/gitea/runner/commits/b1c873a66b4fd1eae8f52fe568b2c6958f2df36f)), references [#8203](https://gitea.com/gitea/runner/issues/8203) [README.md#logging](https://gitea.com/README.md/issues/logging)
* **deps:** update dependency @vercel/ncc to ^0.38.0 ([#881](https://gitea.com/gitea/runner/issues/881)) ([68d9fc4](https://gitea.com/gitea/runner/commits/68d9fc45c9c04ecac239fa28a28174760b2c7948)), references [#8203](https://gitea.com/gitea/runner/issues/8203)

## [0.6.1](https://gitea.com/gitea/runner/compare/v0.6.0...v0.6.1) (2026-04-30)

### Bug Fixes

* **deps:** update module github.com/docker/docker to v25.0.15+incompatible ([#867](https://gitea.com/gitea/runner/issues/867)) ([11a5dc8](https://gitea.com/gitea/runner/commits/11a5dc893646272d2a8c73eb0f4d9c4abb9cbc11))
* **deps:** update module github.com/moby/patternmatcher to v0.6.1 ([#868](https://gitea.com/gitea/runner/issues/868)) ([35834bf](https://gitea.com/gitea/runner/commits/35834bf8178b7330b4f17315cdb35fa95d4b4f1e))

### Performance Improvements

* **config:** lower default fetch_interval_max from 60s to 5s ([#875](https://gitea.com/gitea/runner/issues/875)) ([e5e53c7](https://gitea.com/gitea/runner/commits/e5e53c732ed825791e2a1325dc1ae3845358acbb)), references [#819](https://gitea.com/gitea/runner/issues/819)

### Chores

* clean up `nolint` directives in `act` package ([#864](https://gitea.com/gitea/runner/issues/864)) ([2516573](https://gitea.com/gitea/runner/commits/251657359217f8750fc3ddb60a8f116965f6f585))

## [0.6.0](https://gitea.com/gitea/runner/compare/v0.5.0...v0.6.0) (2026-04-28)

### Features

* show run command, shell and env in collapsible group before step output ([#847](https://gitea.com/gitea/runner/issues/847)) ([547a0ff](https://gitea.com/gitea/runner/commits/547a0ff297fe077ef823ed1cc583a5dcf7f80479))

### Bug Fixes

* avoid 'filename too long' on matrix jobs by hashing container names ([#853](https://gitea.com/gitea/runner/issues/853)) ([801e5cf](https://gitea.com/gitea/runner/commits/801e5cf4d51baf9de3c843a8ede953a91af3e243)), closes [#686](https://gitea.com/gitea/runner/issues/686)
* **deps:** update module github.com/docker/docker to v25.0.14+incompatible ([#862](https://gitea.com/gitea/runner/issues/862)) ([59d90bf](https://gitea.com/gitea/runner/commits/59d90bff26a37e70e8969817e675bc36cb25fbb7))
* **deps:** update module github.com/mattn/go-isatty to v0.0.22 ([#863](https://gitea.com/gitea/runner/issues/863)) ([3f05040](https://gitea.com/gitea/runner/commits/3f05040438c86dc6661c53d26cf41dfdd84a0e69))

### Chores

* **deps:** drop unused `distribution/reference` replace directive ([#858](https://gitea.com/gitea/runner/issues/858)) ([bad4239](https://gitea.com/gitea/runner/commits/bad4239d18fd12b64f2a10c4aeca495271f4a24a))

## [0.5.0](https://gitea.com/gitea/runner/compare/v0.4.1...v0.5.0) (2026-04-27)

### Features

* Add support for Dynamic Matrix Strategies with Job Outputs ([#149](https://gitea.com/gitea/runner/issues/149)) ([15dd63a](https://gitea.com/gitea/runner/commits/15dd63a839b4936bbb920f9dc5a7a2e8ff5b01be))

### Bug Fixes

* **deps:** update module connectrpc.com/connect to v1.19.2 ([#854](https://gitea.com/gitea/runner/issues/854)) ([1032f85](https://gitea.com/gitea/runner/commits/1032f857a19a2c1d58ca152e319c105daada3a9d))
* **deps:** update module github.com/docker/cli to v25.0.7+incompatible ([#855](https://gitea.com/gitea/runner/issues/855)) ([589db33](https://gitea.com/gitea/runner/commits/589db33e70ca2eff43224c69b31e1056890b56fd))
* fetch when other refs get force-pushed ([#846](https://gitea.com/gitea/runner/issues/846)) ([7c6f126](https://gitea.com/gitea/runner/commits/7c6f1261d4bcba4ebc123a694b75c125dac76d88)), closes [#726](https://gitea.com/gitea/runner/issues/726)
* Heartbeat `ReportState` for long-running silent jobs ([#852](https://gitea.com/gitea/runner/issues/852)) ([fa5334e](https://gitea.com/gitea/runner/commits/fa5334eb246f6be9b3bf0141dc37a1bb7f509a64)), closes [#826](https://gitea.com/gitea/runner/issues/826) [#819](https://gitea.com/gitea/runner/issues/819), references [#819](https://gitea.com/gitea/runner/issues/819)
* wait for docker supervise dir before `s6-svwait` ([#851](https://gitea.com/gitea/runner/issues/851)) ([e56b984](https://gitea.com/gitea/runner/commits/e56b984c04de9609177d90b9690246a7657bfe96)), closes [#760](https://gitea.com/gitea/runner/issues/760)

## [0.261.8](https://gitea.com/gitea/runner/compare/v0.2.13...v0.261.8) (2026-02-18)

### Features

* `--container-options` ([#1462](https://gitea.com/gitea/runner/issues/1462)) ([4c2524a](https://gitea.com/gitea/runner/commits/4c2524ab4d0d923f241f2649702e5dd8927846a4))
* add `go` install instruction ([#672](https://gitea.com/gitea/runner/issues/672)) ([37a98f1](https://gitea.com/gitea/runner/commits/37a98f134a7e69bd45b0da40ea175d4cb5b41211))
* add bug-report flag ([#1056](https://gitea.com/gitea/runner/issues/1056)) ([5d7027d](https://gitea.com/gitea/runner/commits/5d7027dc3fc6a16e2023ffa1ebac93f1c667786e))
* add check for newer versions ([#1562](https://gitea.com/gitea/runner/issues/1562)) ([9390793](https://gitea.com/gitea/runner/commits/93907931df737bdd21ecc624f4983fb1b365357a))
* Add CopyTarStream to Container Interface ([#1912](https://gitea.com/gitea/runner/issues/1912)) ([6468dd7](https://gitea.com/gitea/runner/commits/6468dd7fc8c2013dc76e8ea6e9ab590dc6dbb56e))
* add Dependabot config ([#797](https://gitea.com/gitea/runner/issues/797)) ([12928a0](https://gitea.com/gitea/runner/commits/12928a0ac61b2b3278aead1e42923c876a8cf5a4))
* add flag to always run the checkout action ([#1049](https://gitea.com/gitea/runner/issues/1049)) ([f09cdb8](https://gitea.com/gitea/runner/commits/f09cdb844367189aafe2f3af5795d3136db6b290))
* add fromJSON support ([#352](https://gitea.com/gitea/runner/issues/352)) ([d521fa5](https://gitea.com/gitea/runner/commits/d521fa5bbae1d6fd5b775e0924907e882a4d8655))
* add json logger output ([#1026](https://gitea.com/gitea/runner/issues/1026)) ([14c9801](https://gitea.com/gitea/runner/commits/14c98015c57bf577619101ac4770f90c6435470b))
* add logger hook for standalone job logger ([0fa1fe0](https://gitea.com/gitea/runner/commits/0fa1fe03109cb5e41aa7a074fe0d837089d65c1e))
* Add new Action Cache ([#1913](https://gitea.com/gitea/runner/issues/1913)) ([a42f3cf](https://gitea.com/gitea/runner/commits/a42f3cf1cdaa1d93da181bd5f556371d61198169))
* add node16 support ([#922](https://gitea.com/gitea/runner/issues/922)) ([5580812](https://gitea.com/gitea/runner/commits/558081242ca322997437de126324ae9977eb7da0))
* add option for custom socket path ([#698](https://gitea.com/gitea/runner/issues/698)) ([799098b](https://gitea.com/gitea/runner/commits/799098b0e630caff38f220dd0ab5e216667ec916))
* add option for docker image rebuild ([#878](https://gitea.com/gitea/runner/issues/878)) ([8ad6c07](https://gitea.com/gitea/runner/commits/8ad6c07083c7cd7f45ea9e0808e3297bf2370ab5)), references [#714](https://gitea.com/gitea/runner/issues/714) [#716](https://gitea.com/gitea/runner/issues/716)
* add option to bypass GHE for actions checkout ([#1162](https://gitea.com/gitea/runner/issues/1162)) ([de37f75](https://gitea.com/gitea/runner/commits/de37f750773575651a91363eb6b1ce97be964845)), references [#1161](https://gitea.com/gitea/runner/issues/1161)
* add option to specify git remote name ([#1104](https://gitea.com/gitea/runner/issues/1104)) ([d4272bd](https://gitea.com/gitea/runner/commits/d4272bd9fe13566ab76e5ad091aedbded74ff4e2))
* add remote reusable workflows ([#1525](https://gitea.com/gitea/runner/issues/1525)) ([82a8c1e](https://gitea.com/gitea/runner/commits/82a8c1e80dfcd73842bc0359ca8447303b6fcc08)), closes [#826](https://gitea.com/gitea/runner/issues/826)
* add skipped status as step result ([#950](https://gitea.com/gitea/runner/issues/950)) ([2eda7c6](https://gitea.com/gitea/runner/commits/2eda7c6037b921ae88742b2ff29b20bb171497a6))
* add step name to logger field ([#1027](https://gitea.com/gitea/runner/issues/1027)) ([6b05957](https://gitea.com/gitea/runner/commits/6b059572ff51e30f4cfbde76cc7b82af313967b6))
* add support for building docker actions with private registries ([#1557](https://gitea.com/gitea/runner/issues/1557)) ([3f3b25a](https://gitea.com/gitea/runner/commits/3f3b25ae84a25f8488f3801e3342bb3f96c6cd8a))
* Add support for github action variables ([#1833](https://gitea.com/gitea/runner/issues/1833)) ([6ce45e3](https://gitea.com/gitea/runner/commits/6ce45e3f246f12d6617691de9a2423920d5fdbe6))
* add ubuntu-22.04 ([#1150](https://gitea.com/gitea/runner/issues/1150)) ([69691c8](https://gitea.com/gitea/runner/commits/69691c88b301ee695f5712d2ce8ad2b760725c80))
* adjust container name format ([#1](https://gitea.com/gitea/runner/issues/1)) ([4cacc14](https://gitea.com/gitea/runner/commits/4cacc14d22447358f1d34b4ef556fe1fd9cb9dd0))
* Allow building without docker support ([#1507](https://gitea.com/gitea/runner/issues/1507)) ([7e8d070](https://gitea.com/gitea/runner/commits/7e8d0708111125b11166494dc22da6d67f9e83e7))
* allow existing logger from context ([#859](https://gitea.com/gitea/runner/issues/859)) ([603cf56](https://gitea.com/gitea/runner/commits/603cf568780752efbe587dd50825b2ee5c3013ae))
* allow existing logger from context ([#898](https://gitea.com/gitea/runner/issues/898)) ([6517d04](https://gitea.com/gitea/runner/commits/6517d04b9a4a677fc03277b6cf32fcd56a4dcf36))
* allow overriding of `GITHUB_` env variables ([#1582](https://gitea.com/gitea/runner/issues/1582)) ([ce168f9](https://gitea.com/gitea/runner/commits/ce168f959502feab3ec815810fac7a3b5a1f470c))
* allow to spawn and run a local reusable workflow ([#1423](https://gitea.com/gitea/runner/issues/1423)) ([a8e05cd](https://gitea.com/gitea/runner/commits/a8e05cded6be8bfb77536d823933725144918bcb))
* cache notices to reduce frequency of upgrade notifications ([#1592](https://gitea.com/gitea/runner/issues/1592)) ([42b9b73](https://gitea.com/gitea/runner/commits/42b9b73d388544693d44d419457151c10e66472a))
* cli option to enable the new action cache ([#1954](https://gitea.com/gitea/runner/issues/1954)) ([f7a846d](https://gitea.com/gitea/runner/commits/f7a846d2f53ad47bdd98483ae25829497f8a1f0a))
* cmd support for windows ([#1941](https://gitea.com/gitea/runner/issues/1941)) ([8314095](https://gitea.com/gitea/runner/commits/83140951bf96c500f0c1a745e4333b27f2820e49))
* correctly use the xdg library, which has the side effect to fix the config survey ([#2195](https://gitea.com/gitea/runner/issues/2195)) ([12c0c42](https://gitea.com/gitea/runner/commits/12c0c4277af1e3d07a6389c3260d28579ba31e01))
* docker 19.03.12 ([#347](https://gitea.com/gitea/runner/issues/347)) ([b3b8961](https://gitea.com/gitea/runner/commits/b3b8961122d06c17bb9bbc063bf56e12aba60106))
* Don't push `pkg/runner/act/` where actions cache is stored ([#620](https://gitea.com/gitea/runner/issues/620)) ([b04d762](https://gitea.com/gitea/runner/commits/b04d762614433695202ee1c12d16654b558165bc))
* GITHUB_STATE and GITHUB_OUTPUT file commands ([#1391](https://gitea.com/gitea/runner/issues/1391)) ([e520382](https://gitea.com/gitea/runner/commits/e520382d2fdc5872ca8e939c6c0f2e1e74749182))
* **goreleaser:** format changelog, add prereleases ([#942](https://gitea.com/gitea/runner/issues/942)) ([4e6cddf](https://gitea.com/gitea/runner/commits/4e6cddf80a712469695640058bcfceee7adaeda6))
* handle context cancelation during docker exec ([#1170](https://gitea.com/gitea/runner/issues/1170)) ([4ef50ee](https://gitea.com/gitea/runner/commits/4ef50eeae76f76af3ab0c8c7377c9388d28a0c38))
* Host environment ([#1293](https://gitea.com/gitea/runner/issues/1293)) ([f2b98ed](https://gitea.com/gitea/runner/commits/f2b98ed301290f53e96c64814ab9949640e10307))
* implement steps.timeout-minutes ([#1776](https://gitea.com/gitea/runner/issues/1776)) ([8913375](https://gitea.com/gitea/runner/commits/8913375af89792be37603c3c6429372964bf86b1))
* improve GetOutboundIP ([#1707](https://gitea.com/gitea/runner/issues/1707)) ([8518d70](https://gitea.com/gitea/runner/commits/8518d70bdf3ee04f08e5d096fe07c2a85e8f21ea))
* improve list ([#786](https://gitea.com/gitea/runner/issues/786)) ([37aaec8](https://gitea.com/gitea/runner/commits/37aaec81f444bd5e00de4e8cedb35c73e131b41d))
* interpolate the step names ([#1422](https://gitea.com/gitea/runner/issues/1422)) ([809da71](https://gitea.com/gitea/runner/commits/809da7198cc2bed5b758c77274d47cb5360ee848)), closes [#1353](https://gitea.com/gitea/runner/issues/1353)
* JobLoggerFactory ([#1496](https://gitea.com/gitea/runner/issues/1496)) ([7d99512](https://gitea.com/gitea/runner/commits/7d9951200f366ee7a207844e185bf5de7e3f7ab5))
* load every environment from --env-file to workflow ([#184](https://gitea.com/gitea/runner/issues/184)) ([2f39547](https://gitea.com/gitea/runner/commits/2f395475b0ac4f5c0fb3d7c8b97a23d9214ea77b))
* non root user container ([#1202](https://gitea.com/gitea/runner/issues/1202)) ([b59e6de](https://gitea.com/gitea/runner/commits/b59e6dee6d01253cb96505f8effc0516c2a390d9))
* offline mode for new action cache ([#2173](https://gitea.com/gitea/runner/issues/2173)) ([852959e](https://gitea.com/gitea/runner/commits/852959e1e14bf5cd9c0284d15ca899bef68ff4af))
* parse types of reusable workflows ([#1414](https://gitea.com/gitea/runner/issues/1414)) ([9b95a72](https://gitea.com/gitea/runner/commits/9b95a728d84920a75d22c0dc6088391fbbc4002a)), references [#826](https://gitea.com/gitea/runner/issues/826)
* pass current matrix information to job logger ([#1364](https://gitea.com/gitea/runner/issues/1364)) ([00acf68](https://gitea.com/gitea/runner/commits/00acf681880b0fe7cea27a29d61c261d624966c3))
* release extension ([78e1ceb](https://gitea.com/gitea/runner/commits/78e1ceb0a8ad911b04f821ca707c3b352fe2c74b))
* release extension ([d690a5f](https://gitea.com/gitea/runner/commits/d690a5fee95491d519085990d5c700ff2c9081ce))
* release extension ([19e6929](https://gitea.com/gitea/runner/commits/19e692939861bf3770e4c0ea7a7d8d5f27d84561))
* run jobs in parallel ([#1003](https://gitea.com/gitea/runner/issues/1003)) ([18792f9](https://gitea.com/gitea/runner/commits/18792f9620fd2fd7dacc53ae30a234fbfc7c81b1)), references [#CPU](https://gitea.com/gitea/runner/issues/CPU)
* SELinux lables for --bind ([#906](https://gitea.com/gitea/runner/issues/906)) ([f726339](https://gitea.com/gitea/runner/commits/f7263399b9146a43b988e011d986d18d180414d2))
* specify matrix on command line ([#1675](https://gitea.com/gitea/runner/issues/1675)) ([636c8a3](https://gitea.com/gitea/runner/commits/636c8a34aedf9b6251df6729fe673d828f972030))
* Speeding up workflows ([#643](https://gitea.com/gitea/runner/issues/643)) ([3e768cd](https://gitea.com/gitea/runner/commits/3e768cd91695af12bdf860ba49a0f992854b9c87))
* split job steps into its own files/structs ([#1004](https://gitea.com/gitea/runner/issues/1004)) ([2bb3e74](https://gitea.com/gitea/runner/commits/2bb3e74616f86ece6a11cd854ad54701542410fd))
* step summary of test results ([#1580](https://gitea.com/gitea/runner/issues/1580)) ([932863b](https://gitea.com/gitea/runner/commits/932863bef508868168954ef69572750cf87f463f))
* Support "result" on "needs" context. ([#1497](https://gitea.com/gitea/runner/issues/1497)) ([bef9b5c](https://gitea.com/gitea/runner/commits/bef9b5c3c77bcdbc6c15169c0d55f635e27350ae)), closes [#1367](https://gitea.com/gitea/runner/issues/1367) [#1466](https://gitea.com/gitea/runner/issues/1466), references [#1423](https://gitea.com/gitea/runner/issues/1423) [#1506](https://gitea.com/gitea/runner/issues/1506)
* support `GITHUB_REF_NAME` & `GITHUB_REF_TYPE` ([#1142](https://gitea.com/gitea/runner/issues/1142)) ([8e216c6](https://gitea.com/gitea/runner/commits/8e216c642e11dfa8dba7aab040bb09274c48b5a6))
* support branches as GITHUB_REF ([03d57a8](https://gitea.com/gitea/runner/commits/03d57a80c9efe130fc0dde98ab7895ac312945c5))
* support config env expansion ([#2063](https://gitea.com/gitea/runner/issues/2063)) ([0806c8b](https://gitea.com/gitea/runner/commits/0806c8b10946de509889eb723aaba45cc100c8de))
* support custom GITHUB_RUN_ID, GITHUB_RUN_NUMBER ([#369](https://gitea.com/gitea/runner/issues/369)) ([1070c58](https://gitea.com/gitea/runner/commits/1070c585381980abe8f636af5212457711a51f03))
* support interpolation in `<job>.container.options` ([#1958](https://gitea.com/gitea/runner/issues/1958)) ([19764bc](https://gitea.com/gitea/runner/commits/19764bcb06ac9acdd7100d95c72055b697820705))
* support more options of containers ([b9c20dc](https://gitea.com/gitea/runner/commits/b9c20dcaa43899cb3bb327619d447248303170e0))
* support node20 runtime ([#1988](https://gitea.com/gitea/runner/issues/1988)) ([935e37c](https://gitea.com/gitea/runner/commits/935e37c25b6c48b50882024b3c4235d12f23f469))
* support offline mode ([#2128](https://gitea.com/gitea/runner/issues/2128)) ([f055d4a](https://gitea.com/gitea/runner/commits/f055d4ae607a1d68695bbee1d9b444d3056fb387))
* support PlatformPicker ([c6b8548](https://gitea.com/gitea/runner/commits/c6b8548d35ab4a5fc0f4478ef9f001ee4fbda16f))
* support runs-on labels and group ([#2062](https://gitea.com/gitea/runner/issues/2062)) ([04011b6](https://gitea.com/gitea/runner/commits/04011b6b78c70d36db18419861c737cd0d72d146))
* support yaml env/secrets/inputs file ([#1733](https://gitea.com/gitea/runner/issues/1733)) ([50dcc57](https://gitea.com/gitea/runner/commits/50dcc57e4babe816692d342307183801ddd78356))
* Test writing remote Actions locally ([#2155](https://gitea.com/gitea/runner/issues/2155)) ([ac6a465](https://gitea.com/gitea/runner/commits/ac6a465e27dccae89e3fb98de355ae79f00f2d6e))
* try to read ref and sha from event payload if available ([#889](https://gitea.com/gitea/runner/issues/889)) ([edd0fb9](https://gitea.com/gitea/runner/commits/edd0fb92a78f2fb20ea1afe14a047fdc026d2643))
* uses in composite ([#793](https://gitea.com/gitea/runner/issues/793)) ([9868e13](https://gitea.com/gitea/runner/commits/9868e1377274a7657d4aac6bd2b66c0d564405cb))
* workflowpattern package ([#1618](https://gitea.com/gitea/runner/issues/1618)) ([05eaeaa](https://gitea.com/gitea/runner/commits/05eaeaa528f791cdb796e3fae4cabd833dcb7a40))
* **workflow:** support schedule event ([#4](https://gitea.com/gitea/runner/issues/4)) ([88cce47](https://gitea.com/gitea/runner/commits/88cce47022d34d5cbc49771b3753541fe905f5b6))

### Bug Fixes

* --version, --help do not need check platform config ([#508](https://gitea.com/gitea/runner/issues/508)) ([22d56c3](https://gitea.com/gitea/runner/commits/22d56c3517483c63dc98dc4d280d60eb5a625a3d))
* ::stop-commands:: should continue to print the lines ([#920](https://gitea.com/gitea/runner/issues/920)) ([fec0c0c](https://gitea.com/gitea/runner/commits/fec0c0c529f5510464c858ff9a6664e0647245a8))
* ([#2075](https://gitea.com/gitea/runner/issues/2075)) ([74b0fe8](https://gitea.com/gitea/runner/commits/74b0fe8ba9c04cbf3d55f6536249d959599ed676))
* [#170](https://gitea.com/gitea/runner/issues/170) [#171](https://gitea.com/gitea/runner/issues/171) - add RUNNER_OS and RUNNER_TEMP ([a9533b0](https://gitea.com/gitea/runner/commits/a9533b05ce708a84e931d63cb2cd7b795e70bbe2))
* [#170](https://gitea.com/gitea/runner/issues/170) add RUNNER_TEMP ([a48cce3](https://gitea.com/gitea/runner/commits/a48cce3a78d915d7c9e6bd28248926578caaf1c4))
* [#182](https://gitea.com/gitea/runner/issues/182) - support symlinks in workspace ([8c29099](https://gitea.com/gitea/runner/commits/8c290994c1b84153d2b9f911ac3bb64858fb5e35))
* 168 include error message from docker SDK ([6c4e7b9](https://gitea.com/gitea/runner/commits/6c4e7b9fdedeeeb1eda78cefaff55750033d1327))
* add `ImageOS` env var based on running platform ([#571](https://gitea.com/gitea/runner/issues/571)) ([41b03b5](https://gitea.com/gitea/runner/commits/41b03b581c7aec251d96e9a5ea27009542a798d4))
* add `server_url` attribute to github context ([#1727](https://gitea.com/gitea/runner/issues/1727)) ([c0130ed](https://gitea.com/gitea/runner/commits/c0130ed030e0f267b7b7b5b3ad5e1dfcc08980d8)), closes [#1726](https://gitea.com/gitea/runner/issues/1726)
* add default values in vmInput ([cc369e2](https://gitea.com/gitea/runner/commits/cc369e2f73e84596a936d602a7e112d849e6e6fa))
* add GITHUB_STEP_SUMMARY ([#1607](https://gitea.com/gitea/runner/issues/1607)) ([8790c9b](https://gitea.com/gitea/runner/commits/8790c9b8e6b67f1552b824ba56d5a851f7b394e5))
* add parent step id in composite action step id ([#1268](https://gitea.com/gitea/runner/issues/1268)) ([5f5b895](https://gitea.com/gitea/runner/commits/5f5b8959d6b12b034dfa64e91f8b17c97808eda5))
* add simple concurrency limit ([#823](https://gitea.com/gitea/runner/issues/823)) ([4d552e6](https://gitea.com/gitea/runner/commits/4d552e65cedafda89c11a36abe3a218f5d090e4b))
* add-matcher fails github workflow ([#1532](https://gitea.com/gitea/runner/issues/1532)) ([8b4f210](https://gitea.com/gitea/runner/commits/8b4f21087209e027ca0e2ec225e0c74d03552f86))
* added tests and improved code ([8d181a6](https://gitea.com/gitea/runner/commits/8d181a6a1ea1eca1076ae9d47c3fded31a6099d2))
* align github.ref to GitHub Action ([#1362](https://gitea.com/gitea/runner/issues/1362)) ([97c083e](https://gitea.com/gitea/runner/commits/97c083e902a86010ce6267347b830e91fe2435e1))
* align runner.os / runner.arch to known values ([#1510](https://gitea.com/gitea/runner/issues/1510)) ([6ab71ec](https://gitea.com/gitea/runner/commits/6ab71ecebbc0d01a0e72aecd06a15277723b6031))
* allow override of artifact server bind address ([#1560](https://gitea.com/gitea/runner/issues/1560)) ([d064863](https://gitea.com/gitea/runner/commits/d064863f9bb65cde863c8cb20fea50e0f139032c)), references [#1559](https://gitea.com/gitea/runner/issues/1559)
* allow TimeoutMinutes to be expression ([#1217](https://gitea.com/gitea/runner/issues/1217)) ([52f5c45](https://gitea.com/gitea/runner/commits/52f5c4592cbbf2e4dde3b35ad42f4e7475ad52df))
* allow TimeoutMinutes to be expression in Jobs ([#1247](https://gitea.com/gitea/runner/issues/1247)) ([91296bd](https://gitea.com/gitea/runner/commits/91296bd5eb022dce9e7891540257232b3e9906e1)), references [#1217](https://gitea.com/gitea/runner/issues/1217)
* always execute closeContainer() executor ([#988](https://gitea.com/gitea/runner/issues/988)) ([9abc87b](https://gitea.com/gitea/runner/commits/9abc87b41634cab2ad9325287fbed9135c2664c4)), references [#984](https://gitea.com/gitea/runner/issues/984)
* Apply forcePull only for prebuild docker actions ([#1599](https://gitea.com/gitea/runner/issues/1599)) ([f91b2aa](https://gitea.com/gitea/runner/commits/f91b2aa5dbcb16075f49c96a6108d0d418382f4d))
* avoid race condition by not running exec/start twice ([#702](https://gitea.com/gitea/runner/issues/702)) ([8a9167d](https://gitea.com/gitea/runner/commits/8a9167da82b3824272fdbfcc25b3671cd1260e7d)), references [#627](https://gitea.com/gitea/runner/issues/627) [#627](https://gitea.com/gitea/runner/issues/627)
* cache adjust restore order of exact key matches ([#2267](https://gitea.com/gitea/runner/issues/2267)) ([f825e42](https://gitea.com/gitea/runner/commits/f825e42ce200fc4973c3f28797ffba942d322d38))
* call of reflect.Value.Interface on zero Value ([#1067](https://gitea.com/gitea/runner/issues/1067)) ([380bd58](https://gitea.com/gitea/runner/commits/380bd581b1a56683a77b9eb017a72129f65251f8))
* checkout container path should always use `/` ([#1226](https://gitea.com/gitea/runner/issues/1226)) ([f71f277](https://gitea.com/gitea/runner/commits/f71f2778f0e75e3e28a81f8f726c64f27b7fffa3))
* ci is failing since 2022-10-17 ([#1397](https://gitea.com/gitea/runner/issues/1397)) ([5d2eb1f](https://gitea.com/gitea/runner/commits/5d2eb1f2388a7d99aae811721c70e75ed8b8a76c))
* CI is unstable due to new linting rules ([#1292](https://gitea.com/gitea/runner/issues/1292)) ([3387fd7](https://gitea.com/gitea/runner/commits/3387fd72ee138893b953e2b8b23efceb23bfefef))
* clarify --reuse and --rm flags ([#929](https://gitea.com/gitea/runner/issues/929)) ([4e0ba61](https://gitea.com/gitea/runner/commits/4e0ba618d34a198fa30904b1bf8549010513a931))
* clean code ([9fd6310](https://gitea.com/gitea/runner/commits/9fd63103632d5288cc7fa6c6f7e175e2ca17ad2f))
* cleanup unused actions ([#873](https://gitea.com/gitea/runner/issues/873)) ([3054f3e](https://gitea.com/gitea/runner/commits/3054f3ea1486267424d71c03c943eded2b95f50f))
* coerce booleans to numbers for comparison in exprparser ([#1030](https://gitea.com/gitea/runner/issues/1030)) ([aab2af0](https://gitea.com/gitea/runner/commits/aab2af09194a77a6b93d81088baf981cd5e27c17))
* conclusion and outcome are no integers ([#1136](https://gitea.com/gitea/runner/issues/1136)) ([77ddc89](https://gitea.com/gitea/runner/commits/77ddc89444d8cca5be47a3c83074dcb992382e68))
* continue jobs + steps after failure ([#840](https://gitea.com/gitea/runner/issues/840)) ([1891c72](https://gitea.com/gitea/runner/commits/1891c72ab158508e36009d16b24913fa5836422b)), closes [#442](https://gitea.com/gitea/runner/issues/442)
* copy ignored tracked files ([#1074](https://gitea.com/gitea/runner/issues/1074)) ([c27e079](https://gitea.com/gitea/runner/commits/c27e079407647e0ab455f068061b5e61e31f670f))
* copy whole action dir to container ([#774](https://gitea.com/gitea/runner/issues/774)) ([de44a50](https://gitea.com/gitea/runner/commits/de44a505da37f9405db195902b277821df2742b6))
* correct ref and ref_name ([#1672](https://gitea.com/gitea/runner/issues/1672)) ([6744e68](https://gitea.com/gitea/runner/commits/6744e68ee29d3bcb6878f895d6ccb5f9f8e2c478))
* correctly unmarshal entrypoint ([#876](https://gitea.com/gitea/runner/issues/876)) ([12fa4d7](https://gitea.com/gitea/runner/commits/12fa4d703da8d05930303acbe6bbc8fd00c77de7))
* crash if the id tool fails to run in the container (1660) ([f3c88b5](https://gitea.com/gitea/runner/commits/f3c88b50915a7234af9b06e28b4a45a267172fe1))
* crash malformed composite action ([#1616](https://gitea.com/gitea/runner/issues/1616)) ([53095d7](https://gitea.com/gitea/runner/commits/53095d76f48747f9f6249ba01e9410e9098d3977))
* crash on invalid docker credentials ([#1348](https://gitea.com/gitea/runner/issues/1348)) ([7d0407f](https://gitea.com/gitea/runner/commits/7d0407fa4abfe81308f8787cbd07c4c72b2aa676))
* deep evaluate matrix strategy ([#964](https://gitea.com/gitea/runner/issues/964)) ([d1e0216](https://gitea.com/gitea/runner/commits/d1e0216039d92a4253fecc2b11a775a37ba43ceb))
* do not error that much ([#1137](https://gitea.com/gitea/runner/issues/1137)) ([9357512](https://gitea.com/gitea/runner/commits/93575124d35f7b84860cb8cc553c257e0945d0af))
* docker buildx cache restore not working ([#2236](https://gitea.com/gitea/runner/issues/2236)) ([75e4ad9](https://gitea.com/gitea/runner/commits/75e4ad93f4c6ba51a7a64f12335dd0268dd7465a))
* docker cp of dangling symlink ([#943](https://gitea.com/gitea/runner/issues/943)) ([#948](https://gitea.com/gitea/runner/issues/948)) ([c802064](https://gitea.com/gitea/runner/commits/c802064975099a00c93b0c57f46b1e3964d33eef))
* don't allow `-self-hosted` mode as container image ([#1783](https://gitea.com/gitea/runner/issues/1783)) ([ca9b783](https://gitea.com/gitea/runner/commits/ca9b783491b4ea1c5d67dbffdb95f34fd4a22573))
* don't override env ([#1629](https://gitea.com/gitea/runner/issues/1629)) ([b0a5068](https://gitea.com/gitea/runner/commits/b0a5068f6d8ec06db4488b4b981d39f476d3abf6))
* don't overwrite with empty cmd/entrypoint ([#1076](https://gitea.com/gitea/runner/issues/1076)) ([24422bf](https://gitea.com/gitea/runner/commits/24422bf4e4e1e12b9a1e076b2f9e0a7b14c6f890))
* don't replace empty secrets ([#725](https://gitea.com/gitea/runner/issues/725)) ([96a0564](https://gitea.com/gitea/runner/commits/96a05645269241587fdcf0f1ccc49ce38047627c))
* dryrun skip local actions + enable Tests ([#1199](https://gitea.com/gitea/runner/issues/1199)) ([b7d380b](https://gitea.com/gitea/runner/commits/b7d380b3f096e079f2c1ba91a9377752fc4c9751))
* **editorconfig:** correct braces ([#1094](https://gitea.com/gitea/runner/issues/1094)) ([8d41f7a](https://gitea.com/gitea/runner/commits/8d41f7aa6b6702d9c1c3bfc372dd4c5533e88c8c))
* ensure all post steps are executed ([#1286](https://gitea.com/gitea/runner/issues/1286)) ([3364f95](https://gitea.com/gitea/runner/commits/3364f95569d8a39d440b18af69d29421fcc77b55))
* ensure networkmode "host" unless explicitly specified ([#1739](https://gitea.com/gitea/runner/issues/1739)) ([de06444](https://gitea.com/gitea/runner/commits/de0644499a555d240bd065f7c6d839e5162df46f))
* environment handling windows (host mode) ([#1732](https://gitea.com/gitea/runner/issues/1732)) ([9884da0](https://gitea.com/gitea/runner/commits/9884da0122b7fc45595c4f2950c1044fcce2c3e6))
* extra path lost in composite actions ([#1531](https://gitea.com/gitea/runner/issues/1531)) ([a53a1c2](https://gitea.com/gitea/runner/commits/a53a1c20005614cefd3d6519abd7f76d91467132)), references [#1528](https://gitea.com/gitea/runner/issues/1528)
* fail if no stages were found ([#1970](https://gitea.com/gitea/runner/issues/1970)) ([7286b43](https://gitea.com/gitea/runner/commits/7286b43b0ea0110cc2052848c8c1c4ab66bcd068))
* fail workflow if the job name is invalid ([#596](https://gitea.com/gitea/runner/issues/596)) ([94d736a](https://gitea.com/gitea/runner/commits/94d736a602dcfc7db2568d41a16e222bf1bcc416))
* fallback to unauthenticated pull ([#1774](https://gitea.com/gitea/runner/issues/1774)) ([aa21277](https://gitea.com/gitea/runner/commits/aa212773804561eaee04c3b4a7d9399dfbb9cad1))
* fromJSON should work with any JSON ([#883](https://gitea.com/gitea/runner/issues/883)) ([09693ec](https://gitea.com/gitea/runner/commits/09693ec5b9f0d2643b0a57826e2c415dfc32de2d))
* ghc assignment typo ([#1729](https://gitea.com/gitea/runner/issues/1729)) ([97749a2](https://gitea.com/gitea/runner/commits/97749a27b948268a921ef3e75d679a0f1f5c26ce))
* GITHUB_PATH should prepend ([#690](https://gitea.com/gitea/runner/issues/690)) ([29ea8cf](https://gitea.com/gitea/runner/commits/29ea8cfc4e00401b144cd29643a61ca089efd7e9))
* github.job property is empty, GITHUB_JOB should be job id ([#1646](https://gitea.com/gitea/runner/issues/1646)) ([5500c92](https://gitea.com/gitea/runner/commits/5500c928ebb6c4e7dcfb37c29c7100dd7e721b68)), references [#1621](https://gitea.com/gitea/runner/issues/1621) [#1473](https://gitea.com/gitea/runner/issues/1473)
* **goreleaser:** add append mode for release notes ([#962](https://gitea.com/gitea/runner/issues/962)) ([518148d](https://gitea.com/gitea/runner/commits/518148d1629381b0abd73e8f681d9fb6e0242205))
* handle env-vars case sensitive ([#1493](https://gitea.com/gitea/runner/issues/1493)) ([0f8861b](https://gitea.com/gitea/runner/commits/0f8861b9e37bfdf5e0fabb188d03365c0401823a)), closes [#1488](https://gitea.com/gitea/runner/issues/1488)
* handle go closure iteration ([#1374](https://gitea.com/gitea/runner/issues/1374)) ([1bade27](https://gitea.com/gitea/runner/commits/1bade27534209f2fd2d7a0b6437efb4965b648c8))
* handle subdirectors and exclusions ([#1012](https://gitea.com/gitea/runner/issues/1012)) ([#1446](https://gitea.com/gitea/runner/issues/1446)) ([5200c49](https://gitea.com/gitea/runner/commits/5200c491555ab1665cf280523f97bf92401c0415))
* handle zero size ([#1888](https://gitea.com/gitea/runner/issues/1888)) ([8c7c0f5](https://gitea.com/gitea/runner/commits/8c7c0f53c158867b446d1ef21356f49087c238a3))
* hardcode `ubuntu-latest` for `ImageOS` envvar ([#579](https://gitea.com/gitea/runner/issues/579)) ([29b3d43](https://gitea.com/gitea/runner/commits/29b3d43988e1354519f667830eef06662cadd9e4))
* hide masked value ([#668](https://gitea.com/gitea/runner/issues/668)) ([3db5037](https://gitea.com/gitea/runner/commits/3db50376aa4e0fbf5c9d206cb955c6f660c7779a))
* hide masked values in log files ([#1011](https://gitea.com/gitea/runner/issues/1011)) ([c22d833](https://gitea.com/gitea/runner/commits/c22d833830b3a1903a12c590b0b2ee637c3639e1))
* ignore `act` binary built when using `make` ([#654](https://gitea.com/gitea/runner/issues/654)) ([780f60a](https://gitea.com/gitea/runner/commits/780f60a1e6a07f495ad9a2d0e3181e85b36b0537))
* ignore symlink errors ([c44b71c](https://gitea.com/gitea/runner/commits/c44b71c9962f527e9c8847a7a3b959c71ac41b21))
* improve action not found error ([#2171](https://gitea.com/gitea/runner/issues/2171)) ([a6ec2c1](https://gitea.com/gitea/runner/commits/a6ec2c129a94a9cd488750a3be8b961cff08fd0c))
* improve new-action-cache fetch failure error ([#2172](https://gitea.com/gitea/runner/issues/2172)) ([6dd6725](https://gitea.com/gitea/runner/commits/6dd67253bc4ca18ee7ec5e96b8967332fe94dda1))
* improve warning about remote not found ([#2169](https://gitea.com/gitea/runner/issues/2169)) ([5e0d29d](https://gitea.com/gitea/runner/commits/5e0d29d6655ffb51ba394edd236ddf8e362698e0))
* interpolate action input defaults ([#1376](https://gitea.com/gitea/runner/issues/1376)) ([3631013](https://gitea.com/gitea/runner/commits/36310139e07b29768700e312c92ea49635f0d1ac))
* interpolate env for step actions ([#872](https://gitea.com/gitea/runner/issues/872)) ([099af5e](https://gitea.com/gitea/runner/commits/099af5e6a8e882cab7aa94e11d0616faceb0951f))
* invoke login shells to source `/etc/environment` ([#570](https://gitea.com/gitea/runner/issues/570)) ([b438b83](https://gitea.com/gitea/runner/commits/b438b836ea5e031c1e48cb8b2a777f1647d70f8a))
* IsHost is defined as false on windows ([#2093](https://gitea.com/gitea/runner/issues/2093)) ([1b10028](https://gitea.com/gitea/runner/commits/1b1002844772c5107f16321bd26ccfe7f6f0b48b))
* keep action inputs ([#1215](https://gitea.com/gitea/runner/issues/1215)) ([1d4c2aa](https://gitea.com/gitea/runner/commits/1d4c2aaa3fb1fce7a0fc70a3923391f20ecd6792))
* keep path to event json file in composite actions ([#1428](https://gitea.com/gitea/runner/issues/1428)) ([2614b3e](https://gitea.com/gitea/runner/commits/2614b3eb0c62772b17ba11509287fdc6f843f9d3))
* localcheckout mock ([#1198](https://gitea.com/gitea/runner/issues/1198)) ([c3fb686](https://gitea.com/gitea/runner/commits/c3fb6864e86542f5ff1f025bec5eff4f81497c96))
* logo and demo asset urls ([#2226](https://gitea.com/gitea/runner/issues/2226)) ([e1cd7c9](https://gitea.com/gitea/runner/commits/e1cd7c915f79120d77bdee6f9ec11eadb2da25bf))
* match cache `restore-keys` in creation reverse order ([#2153](https://gitea.com/gitea/runner/issues/2153)) ([7f7d84b](https://gitea.com/gitea/runner/commits/7f7d84b10ffbfbb699c0d95c17e94ce9a610ac44))
* Max parallel Support for Matrix Jobs and Remote Action Tests ([#150](https://gitea.com/gitea/runner/issues/150)) ([c0f19d9](https://gitea.com/gitea/runner/commits/c0f19d9a26134d5b4222aff3194f9d2b6f51d1d7))
* **mergify:** needs-work label handling ([#810](https://gitea.com/gitea/runner/issues/810)) ([a6aea44](https://gitea.com/gitea/runner/commits/a6aea44fb023f33543e58a5fcdcff937cf0e8ec6))
* missing defaults in nodejs pre ([#1349](https://gitea.com/gitea/runner/issues/1349)) ([6037182](https://gitea.com/gitea/runner/commits/6037182487e39d9f50c6887ab3453986fbf29a9b))
* mount volume for workflow info ([#776](https://gitea.com/gitea/runner/issues/776)) ([3c14548](https://gitea.com/gitea/runner/commits/3c1454825d8d4c622a4a94a9b3f796461d06409c))
* nil pointer access ( workflow_dispatch ) ([d97481d](https://gitea.com/gitea/runner/commits/d97481d5671a4906b8b9fe7b951333fa938ffe2b))
* outputs espace ([#404](https://gitea.com/gitea/runner/issues/404)) ([92067eb](https://gitea.com/gitea/runner/commits/92067eb1a5f271278373d6b2c2df9b6f9c373446))
* panic ([6791833](https://gitea.com/gitea/runner/commits/67918333fa7fac871f9de707970f0ef5f564ab61))
* panic ([#2071](https://gitea.com/gitea/runner/issues/2071)) ([5a79256](https://gitea.com/gitea/runner/commits/5a79256ee4d65ca068687bd0c85cd112d438a5d8))
* panic if a step in a job is nil ([#1145](https://gitea.com/gitea/runner/issues/1145)) ([91fd412](https://gitea.com/gitea/runner/commits/91fd412c51667dd10c434e9eb16ba10685581fa6))
* panic on regex in main.workflow ([a04141c](https://gitea.com/gitea/runner/commits/a04141c4440be5c0f003d5441e4e2a97bc09060e)), closes [#12](https://gitea.com/gitea/runner/issues/12)
* panic while running DRYRUN with local actions ([#1141](https://gitea.com/gitea/runner/issues/1141)) ([562265b](https://gitea.com/gitea/runner/commits/562265bb295f618123f65af2b9d5877638b2bda7))
* panic: reflect: slice index out of range ([#1066](https://gitea.com/gitea/runner/issues/1066)) ([a970145](https://gitea.com/gitea/runner/commits/a970145c95b9cc072905ca4f64c0640736aeff1f))
* pass `sc.Env` through `common.Logger` to hide secrets ([#578](https://gitea.com/gitea/runner/issues/578)) ([f3b53d8](https://gitea.com/gitea/runner/commits/f3b53d8eca1142fb587d984d8c95584e236b612a))
* preserve job result state in case of failure ([#1519](https://gitea.com/gitea/runner/issues/1519)) ([b14398e](https://gitea.com/gitea/runner/commits/b14398eff3d3a2692c976e250c4b91fbbeed00a5)), closes [#1518](https://gitea.com/gitea/runner/issues/1518)
* print line after password prompt ([#1014](https://gitea.com/gitea/runner/issues/1014)) ([ad06f5d](https://gitea.com/gitea/runner/commits/ad06f5dfb8a658521c855f8c374a58b02227e8fa))
* processing of strategy.matrix.include ([#1200](https://gitea.com/gitea/runner/issues/1200)) ([c30bc82](https://gitea.com/gitea/runner/commits/c30bc824b23335f2c99a32578f8d907f0be6647e))
* re-evaluate env for remote composite actions ([#1385](https://gitea.com/gitea/runner/issues/1385)) ([f0b1845](https://gitea.com/gitea/runner/commits/f0b1845802ca6bea4a75f4e3b514fad85f070b4f))
* regression run after failure ([#971](https://gitea.com/gitea/runner/issues/971)) ([7dbf3fc](https://gitea.com/gitea/runner/commits/7dbf3fcb96e2d70b3436b718c1f4a9fe4b169a4f))
* remove composite restrictions ([#1128](https://gitea.com/gitea/runner/issues/1128)) ([ebb408f](https://gitea.com/gitea/runner/commits/ebb408f373a4fc6508d1d75aa6fddcba5fc7df81))
* remove hardcoded reference to github.com when using reusable remote workflows and remote actions ([#1784](https://gitea.com/gitea/runner/issues/1784)) ([4721abf](https://gitea.com/gitea/runner/commits/4721abfa6dbc661922d34d20b9f60aa14eb30b34))
* remove HOME env var ([#594](https://gitea.com/gitea/runner/issues/594)) ([5044ec6](https://gitea.com/gitea/runner/commits/5044ec6c4364cb74e623953b8c5d39790c8af126))
* remove unused TestContext ([#1013](https://gitea.com/gitea/runner/issues/1013)) ([1e92d87](https://gitea.com/gitea/runner/commits/1e92d87b75f1691a531c76e941e614b521e308f2))
* Remove volume with same name as the container again ([#1203](https://gitea.com/gitea/runner/issues/1203)) ([28f5b37](https://gitea.com/gitea/runner/commits/28f5b37fd1941e00e77f33e3964011f83789c45c))
* remove Windows issue from `README.md` ([#675](https://gitea.com/gitea/runner/issues/675)) ([cc4e23d](https://gitea.com/gitea/runner/commits/cc4e23d96c24effff203c86d679c22c9ab092f27))
* replace `gopass` lib with `golang.org/x/term` ([#674](https://gitea.com/gitea/runner/issues/674)) ([a6a865e](https://gitea.com/gitea/runner/commits/a6a865e973d7138c70e9699d2e7a9a11b542f16c))
* restore the job logger setup after job cancelation ([#1365](https://gitea.com/gitea/runner/issues/1365)) ([13d5303](https://gitea.com/gitea/runner/commits/13d530302e5f8203b0f104155eb6e136f5aaf910))
* return err in walk ([#1667](https://gitea.com/gitea/runner/issues/1667)) ([ac5dd8f](https://gitea.com/gitea/runner/commits/ac5dd8feb876d37ae483376a137c57383577dace))
* return error if both `run:`/`uses:` keys are used ([#593](https://gitea.com/gitea/runner/issues/593)) ([25779af](https://gitea.com/gitea/runner/commits/25779af4bf1e7762d8ff2d3be985c2fcde7bb1e2))
* return error on reusable workflows ([#1093](https://gitea.com/gitea/runner/issues/1093)) ([b3bd268](https://gitea.com/gitea/runner/commits/b3bd268e4585b32673635685d68bf86c5b2d0bfa))
* return invalid step type ([#1157](https://gitea.com/gitea/runner/issues/1157)) ([1e72c59](https://gitea.com/gitea/runner/commits/1e72c594aca1854854d40a24f3e038537f392fc8))
* reusable workflow panic ([#1728](https://gitea.com/gitea/runner/issues/1728)) ([d70b225](https://gitea.com/gitea/runner/commits/d70b225e85b44efdcd97452959a3bbabec5eb6ed))
* rework setupShellCommand ([#930](https://gitea.com/gitea/runner/issues/930)) ([adabf2a](https://gitea.com/gitea/runner/commits/adabf2a2027fa9da3410febc474d55327e2abeb7))
* reworked container architecture ([#619](https://gitea.com/gitea/runner/issues/619)) ([f32babb](https://gitea.com/gitea/runner/commits/f32babb51db6761fd5dff1bcc03bbdf9eb7aa354))
* rewrite how image env is merged ([#828](https://gitea.com/gitea/runner/issues/828)) ([6c60af7](https://gitea.com/gitea/runner/commits/6c60af76779b739b7487e3fbe83ea15cd2ab8ad5))
* rewrite StepContext env ([#832](https://gitea.com/gitea/runner/issues/832)) ([a0d3602](https://gitea.com/gitea/runner/commits/a0d360236ef5466a441631c2f93794b55d4b76c5))
* Rootless error concerning `/var/run/docker.sock` ([#2181](https://gitea.com/gitea/runner/issues/2181)) ([f2e65e1](https://gitea.com/gitea/runner/commits/f2e65e1d406255cad05b562138758b73219fa74d))
* rootless permission bits (new actions cache) ([#2242](https://gitea.com/gitea/runner/issues/2242)) ([119ceb8](https://gitea.com/gitea/runner/commits/119ceb81d906be5e1b23bd085f642ffde7144ef5))
* safe file name ([#1651](https://gitea.com/gitea/runner/issues/1651)) ([19abab6](https://gitea.com/gitea/runner/commits/19abab637549dae95d4ea3becd77852e2787195b))
* set `sh` as default shell for containers  ([#853](https://gitea.com/gitea/runner/issues/853)) ([e9ae789](https://gitea.com/gitea/runner/commits/e9ae7894e38ed276c41e6cd7b958a81da79249ad))
* set composite outputs on failure ([#945](https://gitea.com/gitea/runner/issues/945)) ([4be9062](https://gitea.com/gitea/runner/commits/4be9062dd2745c1657ad7b98683397a393159436))
* set default version string ([#1070](https://gitea.com/gitea/runner/issues/1070)) ([de81979](https://gitea.com/gitea/runner/commits/de8197980079129395c67b1e2865cff8be4ebca8))
* set logger with trace level ([97629ae](https://gitea.com/gitea/runner/commits/97629ae8af8dbe76387c34e2b053b3645ceb1bde))
* show workflow info even if on.push is not defined ([#1329](https://gitea.com/gitea/runner/issues/1329)) ([#1335](https://gitea.com/gitea/runner/issues/1335)) ([21484b5](https://gitea.com/gitea/runner/commits/21484b5c1ecbd9e1ddb6a57debd7731a3e1fc25f)), closes [#1332](https://gitea.com/gitea/runner/issues/1332) [#1318](https://gitea.com/gitea/runner/issues/1318)
* show workflow name when get error on yml parsing ([#1133](https://gitea.com/gitea/runner/issues/1133)) ([#1328](https://gitea.com/gitea/runner/issues/1328)) ([b514649](https://gitea.com/gitea/runner/commits/b514649c3da64a1be92b8e22b50e296a07d759f6))
* skip local actions pre step in any case ([#1204](https://gitea.com/gitea/runner/issues/1204)) ([8a47394](https://gitea.com/gitea/runner/commits/8a473943c3dffbbf47c2a2ef8b261e845646fa6e)), closes [#1193](https://gitea.com/gitea/runner/issues/1193)
* spelling mistake in readme ([#1854](https://gitea.com/gitea/runner/issues/1854)) ([80f6de5](https://gitea.com/gitea/runner/commits/80f6de51d544e2c06d730640d683814e247ec0c2))
* step env is unavailable in with property expr ([#1458](https://gitea.com/gitea/runner/issues/1458)) ([8c5748a](https://gitea.com/gitea/runner/commits/8c5748a55c6190ba34354ca29fdb473e5d3724a1))
* subpath actions via new artifact cache ([#2170](https://gitea.com/gitea/runner/issues/2170)) ([09d4b5d](https://gitea.com/gitea/runner/commits/09d4b5d6ade1a06fa006a0fb202f9ff8c19c33f6))
* support docker create arguments from container.options ([#1022](https://gitea.com/gitea/runner/issues/1022)) ([#1351](https://gitea.com/gitea/runner/issues/1351)) ([48188a6](https://gitea.com/gitea/runner/commits/48188a6280f367e63936a239976be65317ebabd2))
* support expression for step's continue-on-error field ([#900](https://gitea.com/gitea/runner/issues/900)) ([#1331](https://gitea.com/gitea/runner/issues/1331)) ([d1daf2f](https://gitea.com/gitea/runner/commits/d1daf2f28d8c2256e5712e3064ce6ac45fca2562))
* tail (not absolute) as entrypoint of job container ([#1506](https://gitea.com/gitea/runner/issues/1506)) ([0c8c082](https://gitea.com/gitea/runner/commits/0c8c082ac0baf25909fd79b87897426816c3c4eb))
* **test:** adapt to sjw change ([ae39a4b](https://gitea.com/gitea/runner/commits/ae39a4b1d354e13891ef8bd7eb8bcd4a0916f0e6))
* **test:** needs condition. ([#8](https://gitea.com/gitea/runner/issues/8)) ([0efa2d5](https://gitea.com/gitea/runner/commits/0efa2d5e63fdc03360b4eb6aac0b9ac5eac9ca42))
* the number in the github event is of type number ([#1252](https://gitea.com/gitea/runner/issues/1252)) ([4094462](https://gitea.com/gitea/runner/commits/409446211fe99420d5d3761b04f491f8ef52412d))
* tojson(null) == 'null' ([#1068](https://gitea.com/gitea/runner/issues/1068)) ([5083284](https://gitea.com/gitea/runner/commits/5083284deb804db5b9d5997ab4cf345b1c7ec069))
* tolerate workflow that needs a missing job ([#1595](https://gitea.com/gitea/runner/issues/1595)) ([#1619](https://gitea.com/gitea/runner/issues/1619)) ([be4a147](https://gitea.com/gitea/runner/commits/be4a1477a5b73b3e17f79891ec6147413acb4df7))
* trim whitespace from actrc ([#1024](https://gitea.com/gitea/runner/issues/1024)) ([3db3c73](https://gitea.com/gitea/runner/commits/3db3c737230adc457e08a8d60b34458977d01e5c)), closes [#985](https://gitea.com/gitea/runner/issues/985)
* try finding a socket, otherwise fail, respect user choice ([#1745](https://gitea.com/gitea/runner/issues/1745)) ([b6718fd](https://gitea.com/gitea/runner/commits/b6718fdf5d2fdc3e61798f4ac749b53548cbbfae))
* update `ubuntu-latest` to point to `ubuntu-20.04` ([#580](https://gitea.com/gitea/runner/issues/580)) ([72f2d2d](https://gitea.com/gitea/runner/commits/72f2d2de51fcc309a2f0381a9b5098c81802598b))
* Update ARCH environment variable used in runners/actions ([#1818](https://gitea.com/gitea/runner/issues/1818)) ([11dd2ac](https://gitea.com/gitea/runner/commits/11dd2ac7457bd22080b9d8e5277327110563919a))
* update artifact server to address GHSL-2023-004 ([#1565](https://gitea.com/gitea/runner/issues/1565)) ([63ae215](https://gitea.com/gitea/runner/commits/63ae215071f94569d910964bdee866d91d6e3a10))
* update codecov action ([#727](https://gitea.com/gitea/runner/issues/727)) ([05bce00](https://gitea.com/gitea/runner/commits/05bce00e932e450ec3cd0d73c064e2a396ac9500))
* update codecov action (again) ([#731](https://gitea.com/gitea/runner/issues/731)) ([6550eb7](https://gitea.com/gitea/runner/commits/6550eb7f2515b635c1210fc218e03132008463c6))
* update install.sh ([#937](https://gitea.com/gitea/runner/issues/937)) ([e9c4609](https://gitea.com/gitea/runner/commits/e9c4609dcadf70739e34f79680e69df973b24d8e))
* update issue template ([#722](https://gitea.com/gitea/runner/issues/722)) ([fa93f4d](https://gitea.com/gitea/runner/commits/fa93f4d5e7e9ccfe68741a555e3137c1dfd6e997))
* update output handling for reusable workflows ([#1521](https://gitea.com/gitea/runner/issues/1521)) ([89cb558](https://gitea.com/gitea/runner/commits/89cb55855802fd7997bf9128140ee4e48b7e89c1))
* update readme with more relevant/updated information ([#630](https://gitea.com/gitea/runner/issues/630)) ([08141f5](https://gitea.com/gitea/runner/commits/08141f5b0095f0ddbbf3a62b03338a9570e27658))
* use actions/runner hashfiles in container ([#1940](https://gitea.com/gitea/runner/issues/1940)) ([7c7d80e](https://gitea.com/gitea/runner/commits/7c7d80ebdd4aee01bb858a828ad864aa120b4a75))
* use auth when fetching (same as with cloning) ([#687](https://gitea.com/gitea/runner/issues/687)) ([38f6dfb](https://gitea.com/gitea/runner/commits/38f6dfb49a7e3cee826911e36f0c990813282d28))
* use correct path to toolcache ([#1494](https://gitea.com/gitea/runner/issues/1494)) ([054caec](https://gitea.com/gitea/runner/commits/054caec791b2e98a917210ad4acdb46181372e09))
* use docker images from dockerhub ([#1249](https://gitea.com/gitea/runner/issues/1249)) ([e70b968](https://gitea.com/gitea/runner/commits/e70b9689241785df1ff509678d99f1fdcb138d91))
* use docker lib for image ref parsing ([#877](https://gitea.com/gitea/runner/issues/877)) ([ec34eb9](https://gitea.com/gitea/runner/commits/ec34eb9532d314ba8d3c23472acf4037f2c084cd))
* use github auth token for pulling a repository ([#693](https://gitea.com/gitea/runner/issues/693)) ([60b9606](https://gitea.com/gitea/runner/commits/60b9606fc687549c7f55343a4d52ff3f2c32cf4d))
* use new platforms after survey, check working dir for `.actrc` ([#577](https://gitea.com/gitea/runner/issues/577)) ([780a8a0](https://gitea.com/gitea/runner/commits/780a8a061c2c7a3a6958e444e105d7d2441893a9))
* use os.UserHomeDir ([#1706](https://gitea.com/gitea/runner/issues/1706)) ([d3dfde0](https://gitea.com/gitea/runner/commits/d3dfde055aad9c0e9631fd6bbf06f4b0db706eff))
* using filepath instead of path ([01f7ff3](https://gitea.com/gitea/runner/commits/01f7ff315c1b3eaa5e3a4eee06bd6d88ebf60ce4))
* warn Apple M1 (arm64) users ([#743](https://gitea.com/gitea/runner/issues/743)) ([1cf422e](https://gitea.com/gitea/runner/commits/1cf422e41108bf72ec837bea66e080d400ca1192))
* write default config in XDG config dir to avoid cluttering the HOME directory by default ([#2140](https://gitea.com/gitea/runner/issues/2140)) ([6091094](https://gitea.com/gitea/runner/commits/6091094e14ab4e2940bc4bde12525d2dfd569e0b))
* wrong perms ([#726](https://gitea.com/gitea/runner/issues/726)) ([2261204](https://gitea.com/gitea/runner/commits/2261204c65b71d661495de95df6610c451157a67))

### Refactoring

* `NewWorkflowPlanner` ([#648](https://gitea.com/gitea/runner/issues/648)) ([806bc4d](https://gitea.com/gitea/runner/commits/806bc4d999cafab455bbda9b966e3ed6080a4104))
* **cmd/root:** simplify `parseEnvs` ([#2162](https://gitea.com/gitea/runner/issues/2162)) ([424fd5e](https://gitea.com/gitea/runner/commits/424fd5e02b4be9fce8aae518420b0499c08cb31b))
* default empty if: statement to if: success() in evaluator & remove `FixIfStatement()` ([#990](https://gitea.com/gitea/runner/issues/990)) ([ff13844](https://gitea.com/gitea/runner/commits/ff13844b8643abaff9daa47a5bc96aa91fddcf5e))
* docker build BuildContext field ([#1914](https://gitea.com/gitea/runner/issues/1914)) ([f0ca0ab](https://gitea.com/gitea/runner/commits/f0ca0abc4003e1d99a534564fc1ad62f47eaf8a4))
* export and move shared contexts into pkg/model ([#931](https://gitea.com/gitea/runner/issues/931)) ([ed01f46](https://gitea.com/gitea/runner/commits/ed01f464ede3515c94dedf5739b54576177bf80c)), references [#908](https://gitea.com/gitea/runner/issues/908)
* extract RunContext Executor in JobExecutor ([#984](https://gitea.com/gitea/runner/issues/984)) ([e23223a](https://gitea.com/gitea/runner/commits/e23223ad02ddeafe581df97c7e9462df28d8e104))
* extract setupAction into ActionReader ([#986](https://gitea.com/gitea/runner/issues/986)) ([e4f0080](https://gitea.com/gitea/runner/commits/e4f0080a1efc9887fe8aa892669b2002e70b6040))
* filecollector into new package ([#2174](https://gitea.com/gitea/runner/issues/2174)) ([5a80a04](https://gitea.com/gitea/runner/commits/5a80a044f986a48b4d2ed380815e9bb47bca2b66))
* fix add-path / GITHUB_PATH commands ([#1472](https://gitea.com/gitea/runner/issues/1472)) ([d281230](https://gitea.com/gitea/runner/commits/d281230ccef6438085180866e50a026025a05f67))
* fix savestate in pre steps ([#1466](https://gitea.com/gitea/runner/issues/1466)) ([b8d7e94](https://gitea.com/gitea/runner/commits/b8d7e947cf91e24202326b63821a3ff82a3e2b57))
* GITHUB_ENV command / remove env.PATH ([#1503](https://gitea.com/gitea/runner/issues/1503)) ([e775fea](https://gitea.com/gitea/runner/commits/e775fea2655d9e3342bd06c2743cb4bab791e2ab))
* ignore already closed error ([#1285](https://gitea.com/gitea/runner/issues/1285)) ([92ddcda](https://gitea.com/gitea/runner/commits/92ddcdae098873ad6d9c05748ce5b6c7be83b811))
* move autoremove into the jobexecutor ([#1463](https://gitea.com/gitea/runner/issues/1463)) ([d9fe63e](https://gitea.com/gitea/runner/commits/d9fe63ec249d4ee3468151851b8db642dc32bf2c))
* move from io/ioutil to io and os packages ([#1417](https://gitea.com/gitea/runner/issues/1417)) ([2f1c5a1](https://gitea.com/gitea/runner/commits/2f1c5a19f12ad0fc784924b4585d4ace061c7ef5))
* open boltdb only while using it ([#1879](https://gitea.com/gitea/runner/issues/1879)) ([e597046](https://gitea.com/gitea/runner/commits/e597046195774991a75ff66506fdfcc3fc903364))
* pull and rebuild docker by default ([#1569](https://gitea.com/gitea/runner/issues/1569)) ([b00babd](https://gitea.com/gitea/runner/commits/b00babd0d9fe8538aced303346cc7139ef51d0cd))
* re-implement `embed` without "unused" import ([#830](https://gitea.com/gitea/runner/issues/830)) ([7a426a0](https://gitea.com/gitea/runner/commits/7a426a0f37d21e32a01e85d9403edba17a54d545))
* remove `gotest.tools` ([#688](https://gitea.com/gitea/runner/issues/688)) ([f571290](https://gitea.com/gitea/runner/commits/f571290b25cdbb34e863cb98eb155369fff5d572))
* remove composite action runcontext workaround ([#1085](https://gitea.com/gitea/runner/issues/1085)) ([e360811](https://gitea.com/gitea/runner/commits/e3608115706382233c3b943b4dabea09bab9d7d2))
* remove docker image list reference filter ([#1501](https://gitea.com/gitea/runner/issues/1501)) ([b2fb9e6](https://gitea.com/gitea/runner/commits/b2fb9e64ac1b263088003b94963b36b6edeaba31))
* remove github.com/pkg/errors dependency ([#1077](https://gitea.com/gitea/runner/issues/1077)) ([2aa0699](https://gitea.com/gitea/runner/commits/2aa0699aecff2e596e280165ea48eac85c5f9a6a))
* remove unnecessary nil check in RunContext ([#1955](https://gitea.com/gitea/runner/issues/1955)) ([8a9e4f9](https://gitea.com/gitea/runner/commits/8a9e4f9f38340ff9b7f9727add3dae62f3f6c7bb))
* return more errors, add more tests ([#679](https://gitea.com/gitea/runner/issues/679)) ([aa68181](https://gitea.com/gitea/runner/commits/aa68181f6b3df2ff24a523800871c3c5ccdf76be))
* share UpdateFromEnv logic ([#1457](https://gitea.com/gitea/runner/issues/1457)) ([57bf4d2](https://gitea.com/gitea/runner/commits/57bf4d27a2b6d03a2c61f39ba57c24138ae94370))
* simpilfy go-git cache ([#2208](https://gitea.com/gitea/runner/issues/2208)) ([5601fb0](https://gitea.com/gitea/runner/commits/5601fb0e136fad09282705ef65cbc9b2c6f0bebf))
* simplify action function signatures ([#1083](https://gitea.com/gitea/runner/issues/1083)) ([5f673cb](https://gitea.com/gitea/runner/commits/5f673cbb6db5569334cc806cf0536a97ad6b899b))
* simplify adding new node versions add node 24 ([#140](https://gitea.com/gitea/runner/issues/140)) ([91852fa](https://gitea.com/gitea/runner/commits/91852faf934e3436fe427d70017b9a630d1421df))

### Build

* **deps:** bump actions/cache from 2 to 3 ([#1082](https://gitea.com/gitea/runner/issues/1082)) ([4f0bfd2](https://gitea.com/gitea/runner/commits/4f0bfd2ad1b16dc0a5f34844069513af178fe86d))
* **deps:** bump actions/checkout from 2 to 3 ([#1028](https://gitea.com/gitea/runner/issues/1028)) ([0579f7c](https://gitea.com/gitea/runner/commits/0579f7c23598fd0b76500b95e5d69d55a3fb0b98))
* **deps:** bump actions/checkout from 3 to 4 ([#1998](https://gitea.com/gitea/runner/issues/1998)) ([0dfb067](https://gitea.com/gitea/runner/commits/0dfb06748ebd077c25451e07a55e563d7bdc6b4a))
* **deps:** bump actions/github-script from 5 to 6 ([#1572](https://gitea.com/gitea/runner/issues/1572)) ([75d19d0](https://gitea.com/gitea/runner/commits/75d19d0ff58e54cafa924b9bf72a89176a602f2b))
* **deps:** bump actions/github-script from 6 to 7 ([#2097](https://gitea.com/gitea/runner/issues/2097)) ([3d0cb3d](https://gitea.com/gitea/runner/commits/3d0cb3d82bf9e0448a245a17aad7a80e672e76a8))
* **deps:** bump actions/setup-go from 2 to 3 ([#1120](https://gitea.com/gitea/runner/issues/1120)) ([4e1cf1d](https://gitea.com/gitea/runner/commits/4e1cf1d7cc84a6da069401ba1ff08d14cf71044b))
* **deps:** bump actions/setup-go from 3 to 4 ([#1689](https://gitea.com/gitea/runner/issues/1689)) ([8c28c9f](https://gitea.com/gitea/runner/commits/8c28c9fd8f2ae53af23367be2cda094ecc4fbfe0))
* **deps:** bump actions/setup-go from 4 to 5 ([#2118](https://gitea.com/gitea/runner/issues/2118)) ([9cecf94](https://gitea.com/gitea/runner/commits/9cecf94039eed005aa76bb4886d9ebfa16b6d184))
* **deps:** bump actions/stale from 4 to 5 ([#1117](https://gitea.com/gitea/runner/issues/1117)) ([f551165](https://gitea.com/gitea/runner/commits/f551165970914f2541f692716c522e95d06a0608))
* **deps:** bump actions/stale from 5 to 6 ([#1357](https://gitea.com/gitea/runner/issues/1357)) ([9a494ca](https://gitea.com/gitea/runner/commits/9a494cac1ed5f22539d727f6958f4d24df7104b8))
* **deps:** bump actions/stale from 6 to 7 ([#1535](https://gitea.com/gitea/runner/issues/1535)) ([d3ce2b0](https://gitea.com/gitea/runner/commits/d3ce2b0a467b39616059b19c5ba0e69e145099aa))
* **deps:** bump actions/stale from 7 to 8 ([#1700](https://gitea.com/gitea/runner/issues/1700)) ([d65cf86](https://gitea.com/gitea/runner/commits/d65cf869e6824ac8aca448df8d7e0bc8b8fc21af))
* **deps:** bump actions/stale from 8 to 9 ([#2120](https://gitea.com/gitea/runner/issues/2120)) ([6b0ef97](https://gitea.com/gitea/runner/commits/6b0ef97c52e7cbab64513b3de7e6618bc2579f28))
* **deps:** bump actions/upload-artifact from 2 to 3 ([#1119](https://gitea.com/gitea/runner/issues/1119)) ([39cc9b0](https://gitea.com/gitea/runner/commits/39cc9b0840f07953b44416b62d4817ca375fe932))
* **deps:** bump actions/upload-artifact from 3 to 4 ([#2133](https://gitea.com/gitea/runner/issues/2133)) ([15bb54f](https://gitea.com/gitea/runner/commits/15bb54f14e96734963b78310f0cb8ed0c67e25bb))
* **deps:** bump codecov/codecov-action from 1.5.2 to 2.1.0 ([#819](https://gitea.com/gitea/runner/issues/819)) ([e82f17e](https://gitea.com/gitea/runner/commits/e82f17e230a7ae5b715b041b6827825de671dc9e))
* **deps:** bump codecov/codecov-action from 2.1.0 to 3.0.0 ([#1118](https://gitea.com/gitea/runner/issues/1118)) ([74cda72](https://gitea.com/gitea/runner/commits/74cda72c2672a1f665d5a55416cd93ca6292180a))
* **deps:** bump codecov/codecov-action from 3.0.0 to 3.1.0 ([#1139](https://gitea.com/gitea/runner/issues/1139)) ([3f2469b](https://gitea.com/gitea/runner/commits/3f2469bf5385b90da5f6a5d9922ccf9c9d7274a6))
* **deps:** bump codecov/codecov-action from 3.1.0 to 3.1.1 ([#1359](https://gitea.com/gitea/runner/issues/1359)) ([d137ec8](https://gitea.com/gitea/runner/commits/d137ec82852083eef45009dcaf368c3dcdd0b9a6))
* **deps:** bump codecov/codecov-action from 3.1.1 to 3.1.2 ([#1735](https://gitea.com/gitea/runner/issues/1735)) ([6d527bf](https://gitea.com/gitea/runner/commits/6d527bf1a356b69ee4ee35e43159c9d0a21c7718))
* **deps:** bump codecov/codecov-action from 3.1.2 to 3.1.3 ([#1752](https://gitea.com/gitea/runner/issues/1752)) ([60bcfb5](https://gitea.com/gitea/runner/commits/60bcfb5d4f717bd0a581de8e97866d58a5593152))
* **deps:** bump codecov/codecov-action from 3.1.3 to 3.1.4 ([#1813](https://gitea.com/gitea/runner/issues/1813)) ([1371215](https://gitea.com/gitea/runner/commits/1371215aabfd82b9288132042aecc9eacc0ab380))
* **deps:** bump codecov/codecov-action from 3.1.4 to 3.1.5 ([#2175](https://gitea.com/gitea/runner/issues/2175)) ([4ca35d2](https://gitea.com/gitea/runner/commits/4ca35d2192ccd6bddf4ebb9519d5fe27a48a6132))
* **deps:** bump deps which dependabot doesn't want to ([#1107](https://gitea.com/gitea/runner/issues/1107)) ([afb23a2](https://gitea.com/gitea/runner/commits/afb23a289187f4a90241b8b923555b7c3d844ec9))
* **deps:** bump docker/setup-qemu-action from 1 to 2 ([#1156](https://gitea.com/gitea/runner/issues/1156)) ([7974a17](https://gitea.com/gitea/runner/commits/7974a17afc613431461f9a8c629597962b1f54b4))
* **deps:** bump docker/setup-qemu-action from 2 to 3 ([#2007](https://gitea.com/gitea/runner/issues/2007)) ([1891bef](https://gitea.com/gitea/runner/commits/1891bef4335ba55573de1c7c5824acd44b3a4994))
* **deps:** bump fregante/setup-git-user from 1 to 2 ([#1664](https://gitea.com/gitea/runner/issues/1664)) ([fad986a](https://gitea.com/gitea/runner/commits/fad986af5f95a2878c9d9f15ded8d815fdf5e061))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.2.7 to 2.3.1 ([#803](https://gitea.com/gitea/runner/issues/803)) ([a7a2659](https://gitea.com/gitea/runner/commits/a7a2659c0e90a53fee4e55817d0eedcad1ab9967))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.2 to 2.3.4 ([#1079](https://gitea.com/gitea/runner/issues/1079)) ([efd34f1](https://gitea.com/gitea/runner/commits/efd34f1e10e56b251a5b9eb75f269d2e7c15681f))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.4 to 2.3.5 ([#1214](https://gitea.com/gitea/runner/issues/1214)) ([084c4e6](https://gitea.com/gitea/runner/commits/084c4e61d99c1bfd109ea5bc381342d94a6e65cd))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.5 to 2.3.6 ([#1341](https://gitea.com/gitea/runner/issues/1341)) ([3a0fe69](https://gitea.com/gitea/runner/commits/3a0fe6967fd8ecd3cb86550d0225d55a0bb37ac9))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.6 to 2.3.7 ([#1872](https://gitea.com/gitea/runner/issues/1872)) ([df44dff](https://gitea.com/gitea/runner/commits/df44dffd30b6503b91bc63130742c31b1ce7badf))
* **deps:** bump github.com/cloudflare/circl from 1.1.0 to 1.3.3 ([#1800](https://gitea.com/gitea/runner/issues/1800)) ([046f5aa](https://gitea.com/gitea/runner/commits/046f5aa71587ba5e81a3f755c913da2f0644cdfb))
* **deps:** bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 ([#2149](https://gitea.com/gitea/runner/issues/2149)) ([96d6cf8](https://gitea.com/gitea/runner/commits/96d6cf8b2c9587912e32a2a42f212fa710efb08e))
* **deps:** bump github.com/containerd/containerd from 1.6.16 to 1.6.18 ([#1637](https://gitea.com/gitea/runner/issues/1637)) ([4eba04b](https://gitea.com/gitea/runner/commits/4eba04b229388c4fc8bbbcffb5e4d62cd927eac9))
* **deps:** bump github.com/containerd/containerd from 1.6.19 to 1.6.26 ([#2189](https://gitea.com/gitea/runner/issues/2189)) ([df61c7f](https://gitea.com/gitea/runner/commits/df61c7fcdb1bb95a68f68238d1637d39ca0329f7))
* **deps:** bump github.com/containerd/containerd from 1.6.6 to 1.6.12 ([#1554](https://gitea.com/gitea/runner/issues/1554)) ([3ac756b](https://gitea.com/gitea/runner/commits/3ac756b37a496020ec8bb7b5852272d01c471057))
* **deps:** bump github.com/containerd/containerd from 1.7.2 to 1.7.11 ([#2136](https://gitea.com/gitea/runner/issues/2136)) ([74dcce4](https://gitea.com/gitea/runner/commits/74dcce467d84b5507db338e98cc8129dd86e90b1))
* **deps:** bump github.com/creack/pty from 1.1.17 to 1.1.18 ([#1447](https://gitea.com/gitea/runner/issues/1447)) ([ccb3b0e](https://gitea.com/gitea/runner/commits/ccb3b0e8757158bb945cf9d3b12a4a9299eea9c4))
* **deps:** bump github.com/creack/pty from 1.1.18 to 1.1.20 ([#2068](https://gitea.com/gitea/runner/issues/2068)) ([fb4f29f](https://gitea.com/gitea/runner/commits/fb4f29fd6d38c6ce57654be5799cb0c911820de4))
* **deps:** bump github.com/creack/pty from 1.1.20 to 1.1.21 ([#2099](https://gitea.com/gitea/runner/issues/2099)) ([603b44b](https://gitea.com/gitea/runner/commits/603b44b585ee35090a13498c7af15f79a8a12a51))
* **deps:** bump github.com/docker/cli ([#1038](https://gitea.com/gitea/runner/issues/1038)) ([933413e](https://gitea.com/gitea/runner/commits/933413e2c0e1d600e09291ba6e63211a3c294385))
* **deps:** bump github.com/docker/cli ([#1081](https://gitea.com/gitea/runner/issues/1081)) ([34706f4](https://gitea.com/gitea/runner/commits/34706f49fe8c82b3a17c9b22bb4fc5d8a70bd25e))
* **deps:** bump github.com/docker/cli ([#1154](https://gitea.com/gitea/runner/issues/1154)) ([d5e9721](https://gitea.com/gitea/runner/commits/d5e972159f9ee43b9229bf028853755c7f23c4d6))
* **deps:** bump github.com/docker/cli ([#1168](https://gitea.com/gitea/runner/issues/1168)) ([1fed666](https://gitea.com/gitea/runner/commits/1fed666165079c8fcab76ff11b66a60463ee0f1f))
* **deps:** bump github.com/docker/cli ([#1211](https://gitea.com/gitea/runner/issues/1211)) ([87d9687](https://gitea.com/gitea/runner/commits/87d9687716a279525fd952b80ba6f044bce88492))
* **deps:** bump github.com/docker/cli ([#1340](https://gitea.com/gitea/runner/issues/1340)) ([ceac1f2](https://gitea.com/gitea/runner/commits/ceac1f2429073270e984ce7da39e28f31e2c56c6))
* **deps:** bump github.com/docker/cli ([#1393](https://gitea.com/gitea/runner/issues/1393)) ([bd5e0d2](https://gitea.com/gitea/runner/commits/bd5e0d24d68cddab105ee76bc2fb8000dce76a11))
* **deps:** bump github.com/docker/cli ([#1405](https://gitea.com/gitea/runner/issues/1405)) ([ca69531](https://gitea.com/gitea/runner/commits/ca69531096e88bec7f8fbd6aecef9b3769271cb7))
* **deps:** bump github.com/docker/cli ([#1419](https://gitea.com/gitea/runner/issues/1419)) ([ac6571c](https://gitea.com/gitea/runner/commits/ac6571c8af87ca3e26b4a5eed1d0a919b7a7aed3))
* **deps:** bump github.com/docker/cli ([#1516](https://gitea.com/gitea/runner/issues/1516)) ([7679695](https://gitea.com/gitea/runner/commits/767969502aa0cc95e65916981a950250ddb6c71d))
* **deps:** bump github.com/docker/cli ([#1626](https://gitea.com/gitea/runner/issues/1626)) ([98ad62f](https://gitea.com/gitea/runner/commits/98ad62fcefc3b760e69040b25d15d34379376319))
* **deps:** bump github.com/docker/cli ([#1712](https://gitea.com/gitea/runner/issues/1712)) ([24b04df](https://gitea.com/gitea/runner/commits/24b04dfa55ebd31c7b4f2fb114350d14783a533a))
* **deps:** bump github.com/docker/cli ([#1722](https://gitea.com/gitea/runner/issues/1722)) ([f2bd194](https://gitea.com/gitea/runner/commits/f2bd194c7fa282fb35f31bcafe9bf1875da71fc7))
* **deps:** bump github.com/docker/cli ([#1753](https://gitea.com/gitea/runner/issues/1753)) ([ded31bd](https://gitea.com/gitea/runner/commits/ded31bdbc0c82ddb737535fab0fc4d1800ae2b00))
* **deps:** bump github.com/docker/cli ([#1780](https://gitea.com/gitea/runner/issues/1780)) ([3eb6e83](https://gitea.com/gitea/runner/commits/3eb6e83ea4c589a0dfcb2a17e93a04c23909f770))
* **deps:** bump github.com/docker/cli ([#1805](https://gitea.com/gitea/runner/issues/1805)) ([39abbce](https://gitea.com/gitea/runner/commits/39abbce4e6cdf3988a689a404f2f0082ce0125bf))
* **deps:** bump github.com/docker/cli ([#1816](https://gitea.com/gitea/runner/issues/1816)) ([9c1f1f8](https://gitea.com/gitea/runner/commits/9c1f1f8d844148079cf80c49413d7f3357c42025))
* **deps:** bump github.com/docker/cli ([#1843](https://gitea.com/gitea/runner/issues/1843)) ([e48c079](https://gitea.com/gitea/runner/commits/e48c07988ed70d94860c6d6276e134f5f51bd802))
* **deps:** bump github.com/docker/cli ([#1932](https://gitea.com/gitea/runner/issues/1932)) ([8e6c7c1](https://gitea.com/gitea/runner/commits/8e6c7c11fe35bd117579277c29238e4e45996a7c))
* **deps:** bump github.com/docker/cli ([#1999](https://gitea.com/gitea/runner/issues/1999)) ([bd9032d](https://gitea.com/gitea/runner/commits/bd9032de0a79e9c0755596b4a1a10f5de8847639))
* **deps:** bump github.com/docker/cli ([#2069](https://gitea.com/gitea/runner/issues/2069)) ([3d65b0f](https://gitea.com/gitea/runner/commits/3d65b0f73f3c2015f75911a1696e1e8c6d844bd5))
* **deps:** bump github.com/docker/cli ([#854](https://gitea.com/gitea/runner/issues/854)) ([50dd0c0](https://gitea.com/gitea/runner/commits/50dd0c0db149afd3f92d64b79b6eb8d2d0b0ec16))
* **deps:** bump github.com/docker/cli ([#956](https://gitea.com/gitea/runner/issues/956)) ([bbf6b7e](https://gitea.com/gitea/runner/commits/bbf6b7e080959ebb13d3500fed637e822e2b8ec0))
* **deps:** bump github.com/docker/cli from 23.0.0-rc.1+incompatible to 23.0.0+incompatible ([#1611](https://gitea.com/gitea/runner/issues/1611)) ([36dbbc1](https://gitea.com/gitea/runner/commits/36dbbc1dfa049676e17e1f4712c6b918ce4ce5b7))
* **deps:** bump github.com/docker/distribution ([#1037](https://gitea.com/gitea/runner/issues/1037)) ([5c0023c](https://gitea.com/gitea/runner/commits/5c0023c3466b7d5d8a900a1b5a3798e4d1dabe75))
* **deps:** bump github.com/docker/distribution ([#1801](https://gitea.com/gitea/runner/issues/1801)) ([b5fa245](https://gitea.com/gitea/runner/commits/b5fa24537f70e88166c7e43d45b5f6ff41d0b12d))
* **deps:** bump github.com/docker/distribution ([#2037](https://gitea.com/gitea/runner/issues/2037)) ([e7e158c](https://gitea.com/gitea/runner/commits/e7e158cd7eed7c78d7111111aa1d005cadf1df0c))
* **deps:** bump github.com/docker/distribution ([#994](https://gitea.com/gitea/runner/issues/994)) ([0fae967](https://gitea.com/gitea/runner/commits/0fae96792c77b395587c4d8129d0f680bffad643))
* **deps:** bump github.com/docker/docker ([#1062](https://gitea.com/gitea/runner/issues/1062)) ([87175ec](https://gitea.com/gitea/runner/commits/87175ecb6882f0735c1de8eeb884ced21f978873))
* **deps:** bump github.com/docker/docker ([#1080](https://gitea.com/gitea/runner/issues/1080)) ([c8a40bf](https://gitea.com/gitea/runner/commits/c8a40bf8f41a68efe55f28e8b17a9b5f3d004d5e))
* **deps:** bump github.com/docker/docker ([#1155](https://gitea.com/gitea/runner/issues/1155)) ([01d85c9](https://gitea.com/gitea/runner/commits/01d85c9a0150092afa061ba1518a96cae6dcb73b))
* **deps:** bump github.com/docker/docker ([#1167](https://gitea.com/gitea/runner/issues/1167)) ([018c846](https://gitea.com/gitea/runner/commits/018c846d8d92979e158f0618b9f3d068fecf08f1))
* **deps:** bump github.com/docker/docker ([#1212](https://gitea.com/gitea/runner/issues/1212)) ([63df5dd](https://gitea.com/gitea/runner/commits/63df5ddf429881b88d766919db2f404d1dc35c9c))
* **deps:** bump github.com/docker/docker ([#1339](https://gitea.com/gitea/runner/issues/1339)) ([ca8767f](https://gitea.com/gitea/runner/commits/ca8767f08057b176ea7f6261c3f557bde1817a8c))
* **deps:** bump github.com/docker/docker ([#1392](https://gitea.com/gitea/runner/issues/1392)) ([6ff4491](https://gitea.com/gitea/runner/commits/6ff44912381a15fbc27ce62937bae89682cfaa2f))
* **deps:** bump github.com/docker/docker ([#1403](https://gitea.com/gitea/runner/issues/1403)) ([0d1dccc](https://gitea.com/gitea/runner/commits/0d1dccca3a5d663860e57cc53a50e61bbbc003fa))
* **deps:** bump github.com/docker/docker ([#1418](https://gitea.com/gitea/runner/issues/1418)) ([e03b816](https://gitea.com/gitea/runner/commits/e03b816c5880b60af28f7dda33e4361d30451092))
* **deps:** bump github.com/docker/docker ([#1517](https://gitea.com/gitea/runner/issues/1517)) ([469d024](https://gitea.com/gitea/runner/commits/469d024c1fb6f376c56bcf380911761d9ad829db))
* **deps:** bump github.com/docker/docker ([#1573](https://gitea.com/gitea/runner/issues/1573)) ([612e08f](https://gitea.com/gitea/runner/commits/612e08f113b0d5c97a8c1ed664fdb2b19b1a418e))
* **deps:** bump github.com/docker/docker ([#1613](https://gitea.com/gitea/runner/issues/1613)) ([f41e912](https://gitea.com/gitea/runner/commits/f41e9127a8e4516e5d78d58702d00c55de1a6b96))
* **deps:** bump github.com/docker/docker ([#1624](https://gitea.com/gitea/runner/issues/1624)) ([7c3c534](https://gitea.com/gitea/runner/commits/7c3c5349ab54ac3a2c779b199310b8337f6e8fbf))
* **deps:** bump github.com/docker/docker ([#1711](https://gitea.com/gitea/runner/issues/1711)) ([68c72b9](https://gitea.com/gitea/runner/commits/68c72b9a51f4fa586af62619f9e6361e023a5d48))
* **deps:** bump github.com/docker/docker ([#1719](https://gitea.com/gitea/runner/issues/1719)) ([65a925b](https://gitea.com/gitea/runner/commits/65a925b4ed5f0bebab9314433fcfdac4018482ac))
* **deps:** bump github.com/docker/docker ([#1754](https://gitea.com/gitea/runner/issues/1754)) ([16f35f6](https://gitea.com/gitea/runner/commits/16f35f64eb992d96992c43ab1ad0e2fae5fbc8b5))
* **deps:** bump github.com/docker/docker ([#1779](https://gitea.com/gitea/runner/issues/1779)) ([1f9bbe1](https://gitea.com/gitea/runner/commits/1f9bbe12dd5a86ee6401000bcc654b9b66910bae))
* **deps:** bump github.com/docker/docker ([#1806](https://gitea.com/gitea/runner/issues/1806)) ([c00810b](https://gitea.com/gitea/runner/commits/c00810bdf4ca9e5cdf0c34e7ede5d4ebc344f558))
* **deps:** bump github.com/docker/docker ([#1933](https://gitea.com/gitea/runner/issues/1933)) ([d720ff0](https://gitea.com/gitea/runner/commits/d720ff09a281aae453fe6cc2970e6f833b1da49d))
* **deps:** bump github.com/docker/docker ([#2000](https://gitea.com/gitea/runner/issues/2000)) ([c84a3ef](https://gitea.com/gitea/runner/commits/c84a3ef6d0e4c545f627572aef0c0313bc9d4187))
* **deps:** bump github.com/docker/docker ([#2067](https://gitea.com/gitea/runner/issues/2067)) ([84a4025](https://gitea.com/gitea/runner/commits/84a4025bc844eee060ad0d2d71af7b4bd2595a94))
* **deps:** bump github.com/docker/docker ([#2252](https://gitea.com/gitea/runner/issues/2252)) ([069720a](https://gitea.com/gitea/runner/commits/069720abff4e169834dda534da5b0013c03775a3))
* **deps:** bump github.com/docker/docker ([#855](https://gitea.com/gitea/runner/issues/855)) ([6a3a47c](https://gitea.com/gitea/runner/commits/6a3a47c2175c5f6438370d28febf8933097cdc55))
* **deps:** bump github.com/docker/docker ([#957](https://gitea.com/gitea/runner/issues/957)) ([b0f5557](https://gitea.com/gitea/runner/commits/b0f55571f092c47e75f746ecb8aeade718db9495))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.0.0 to 5.3.1 ([#802](https://gitea.com/gitea/runner/issues/802)) ([6d7f574](https://gitea.com/gitea/runner/commits/6d7f574859a1a358db7418dbb96fcd07b983eb31))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.3.1 to 5.4.0 ([#1541](https://gitea.com/gitea/runner/issues/1541)) ([9c70851](https://gitea.com/gitea/runner/commits/9c7085195c4fcfe1f37b53dd9a329824a1753b83))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.0 to 5.4.1 ([#1612](https://gitea.com/gitea/runner/issues/1612)) ([22dc1e0](https://gitea.com/gitea/runner/commits/22dc1e0e7ef5291adbc80a50574ad91a07c64754))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.1 to 5.5.0 ([#2012](https://gitea.com/gitea/runner/issues/2012)) ([d134079](https://gitea.com/gitea/runner/commits/d134079807fe35a2d8ebcfd24b139f48c8123caa))
* **deps:** bump github.com/go-git/go-git/v5 ([#1830](https://gitea.com/gitea/runner/issues/1830)) ([b0996e0](https://gitea.com/gitea/runner/commits/b0996e057750223a14321ff92c1c6d6fbb8e1a10))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.0 to 5.10.1 ([#2114](https://gitea.com/gitea/runner/issues/2114)) ([3c405a0](https://gitea.com/gitea/runner/commits/3c405a0d94a7f0c905b51161be0ac695653e48fc))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.1 to 5.11.0 ([#2117](https://gitea.com/gitea/runner/issues/2117)) ([85c3b3b](https://gitea.com/gitea/runner/commits/85c3b3b5419bba9ac6434aefbae1407d34f02654))
* **deps:** bump github.com/go-git/go-git/v5 from 5.7.0 to 5.8.0 ([#1925](https://gitea.com/gitea/runner/issues/1925)) ([17bf4fc](https://gitea.com/gitea/runner/commits/17bf4fc5afd147ea967445b9d0dc2e684f57377c))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.0 to 5.8.1 ([#1934](https://gitea.com/gitea/runner/issues/1934)) ([8bf10cf](https://gitea.com/gitea/runner/commits/8bf10cf87618e3ebe1c0f77c1add9e3a08b3f514))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.1 to 5.9.0 ([#2011](https://gitea.com/gitea/runner/issues/2011)) ([b637d79](https://gitea.com/gitea/runner/commits/b637d79ec3a0468e0ff968cdf57420c7020265c8))
* **deps:** bump github.com/go-git/go-git/v5 from 5.9.0 to 5.10.0 ([#2066](https://gitea.com/gitea/runner/issues/2066)) ([854e3e9](https://gitea.com/gitea/runner/commits/854e3e9ec533132b92838aaf6d801c29064716dd))
* **deps:** bump github.com/go-ini/ini from 1.63.2 to 1.64.0 ([#882](https://gitea.com/gitea/runner/issues/882)) ([c26e661](https://gitea.com/gitea/runner/commits/c26e661889059d3afb2a962f86403e29a8b7c213))
* **deps:** bump github.com/go-ini/ini from 1.64.0 to 1.66.4 ([#997](https://gitea.com/gitea/runner/issues/997)) ([7d43396](https://gitea.com/gitea/runner/commits/7d433968e525098b506aae549d5fc5da4fec571e))
* **deps:** bump github.com/go-ini/ini from 1.66.4 to 1.66.5 ([#1189](https://gitea.com/gitea/runner/issues/1189)) ([64387bc](https://gitea.com/gitea/runner/commits/64387bcf7b92f07b42a4c93f7e4504219af5df2c))
* **deps:** bump github.com/go-ini/ini from 1.66.5 to 1.66.6 ([#1201](https://gitea.com/gitea/runner/issues/1201)) ([859445f](https://gitea.com/gitea/runner/commits/859445fb94e12cc600e5a59d9c9cd687214cb801))
* **deps:** bump github.com/go-ini/ini from 1.66.6 to 1.67.0 ([#1305](https://gitea.com/gitea/runner/issues/1305)) ([12029e3](https://gitea.com/gitea/runner/commits/12029e3df5841538a806293881f836f88ec07a2c))
* **deps:** bump github.com/imdario/mergo from 0.3.12 to 0.3.13 ([#1378](https://gitea.com/gitea/runner/issues/1378)) ([9f9f644](https://gitea.com/gitea/runner/commits/9f9f644d36dc8a7ba73701bdd18cf0d23eaa736b))
* **deps:** bump github.com/imdario/mergo from 0.3.13 to 0.3.14 ([#1688](https://gitea.com/gitea/runner/issues/1688)) ([2fa0a5f](https://gitea.com/gitea/runner/commits/2fa0a5f76979689945215fc392f88f836bc71103))
* **deps:** bump github.com/imdario/mergo from 0.3.14 to 0.3.15 ([#1701](https://gitea.com/gitea/runner/issues/1701)) ([f764ecb](https://gitea.com/gitea/runner/commits/f764ecb283febc36da0fa302a8b8aa8d6d350054))
* **deps:** bump github.com/imdario/mergo from 0.3.15 to 0.3.16 ([#1840](https://gitea.com/gitea/runner/issues/1840)) ([0b05ee8](https://gitea.com/gitea/runner/commits/0b05ee8e0b44454f56baaf23f0f00d7e946b8f37))
* **deps:** bump github.com/joho/godotenv from 1.3.0 to 1.4.0 ([#824](https://gitea.com/gitea/runner/issues/824)) ([bfcf0ab](https://gitea.com/gitea/runner/commits/bfcf0abd73e558521367e36b323afad0dc1f1a71))
* **deps:** bump github.com/joho/godotenv from 1.4.0 to 1.5.1 ([#1614](https://gitea.com/gitea/runner/issues/1614)) ([281a52f](https://gitea.com/gitea/runner/commits/281a52f0d6ed70bd31feaedcbbb46e941cc8bdcc))
* **deps:** bump github.com/mattn/go-isatty from 0.0.14 to 0.0.16 ([#1312](https://gitea.com/gitea/runner/issues/1312)) ([b23bbef](https://gitea.com/gitea/runner/commits/b23bbefc365012886192e477a88d38a0909ecba1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.16 to 0.0.17 ([#1539](https://gitea.com/gitea/runner/issues/1539)) ([0d38857](https://gitea.com/gitea/runner/commits/0d388572a89fabdbd229536b7ecca9db4ea88571))
* **deps:** bump github.com/mattn/go-isatty from 0.0.17 to 0.0.18 ([#1702](https://gitea.com/gitea/runner/issues/1702)) ([385d71a](https://gitea.com/gitea/runner/commits/385d71a215b532dc983c640452e64f38426999b1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.18 to 0.0.19 ([#1859](https://gitea.com/gitea/runner/issues/1859)) ([e912ab3](https://gitea.com/gitea/runner/commits/e912ab32c246692f29d9ef520e92f23edd073a4a))
* **deps:** bump github.com/mattn/go-isatty from 0.0.19 to 0.0.20 ([#2059](https://gitea.com/gitea/runner/issues/2059)) ([db71c41](https://gitea.com/gitea/runner/commits/db71c41d17e967bc12ec24c122aad9d314c7d693))
* **deps:** bump github.com/moby/buildkit from 0.10.0 to 0.10.1 ([#1116](https://gitea.com/gitea/runner/issues/1116)) ([bb6ead4](https://gitea.com/gitea/runner/commits/bb6ead459a8e7b46e817369da3565b5ac86767a1))
* **deps:** bump github.com/moby/buildkit from 0.10.1 to 0.10.2 ([#1147](https://gitea.com/gitea/runner/issues/1147)) ([dfca2c5](https://gitea.com/gitea/runner/commits/dfca2c57df978aae5ec8b29a1a80c0fba09b56a8))
* **deps:** bump github.com/moby/buildkit from 0.10.2 to 0.10.3 ([#1153](https://gitea.com/gitea/runner/issues/1153)) ([a76c349](https://gitea.com/gitea/runner/commits/a76c34987233b065109457893a334879a4615baf))
* **deps:** bump github.com/moby/buildkit from 0.10.3 to 0.10.4 ([#1322](https://gitea.com/gitea/runner/issues/1322)) ([8bad6ac](https://gitea.com/gitea/runner/commits/8bad6aca50b7595fcc5253c72a2936fe37157a63))
* **deps:** bump github.com/moby/buildkit from 0.10.4 to 0.10.5 ([#1406](https://gitea.com/gitea/runner/issues/1406)) ([07a1890](https://gitea.com/gitea/runner/commits/07a18902ffa08ef1f209d8b976142805e2afb8f0))
* **deps:** bump github.com/moby/buildkit from 0.10.5 to 0.10.6 ([#1437](https://gitea.com/gitea/runner/issues/1437)) ([64e68bd](https://gitea.com/gitea/runner/commits/64e68bd7f21afdbbe2f40676c2a8c08701057753))
* **deps:** bump github.com/moby/buildkit from 0.10.6 to 0.11.0 ([#1563](https://gitea.com/gitea/runner/issues/1563)) ([d8ba8cb](https://gitea.com/gitea/runner/commits/d8ba8cbf17d816ed2be8a2b8a031c84b3d6fd8e9))
* **deps:** bump github.com/moby/buildkit from 0.11.0 to 0.11.1 ([#1576](https://gitea.com/gitea/runner/issues/1576)) ([09e9f6a](https://gitea.com/gitea/runner/commits/09e9f6a1ac00592b0341d73d6c642dd4aeb83a02))
* **deps:** bump github.com/moby/buildkit from 0.11.1 to 0.11.2 ([#1589](https://gitea.com/gitea/runner/issues/1589)) ([e259dc2](https://gitea.com/gitea/runner/commits/e259dc2835b6b759de7f4dea2daa490438f41bbe))
* **deps:** bump github.com/moby/buildkit from 0.11.2 to 0.11.3 ([#1635](https://gitea.com/gitea/runner/issues/1635)) ([2f55276](https://gitea.com/gitea/runner/commits/2f55276cea863b82aefc806c52773c5ba65774e2))
* **deps:** bump github.com/moby/buildkit from 0.11.3 to 0.11.4 ([#1669](https://gitea.com/gitea/runner/issues/1669)) ([aeee205](https://gitea.com/gitea/runner/commits/aeee2052de45b3ceff3fb66d3eff461379b89d92))
* **deps:** bump github.com/moby/buildkit from 0.11.4 to 0.11.5 ([#1703](https://gitea.com/gitea/runner/issues/1703)) ([351ae99](https://gitea.com/gitea/runner/commits/351ae99bc1b205c61b02d4fb04c06f8287b57f9f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1755](https://gitea.com/gitea/runner/issues/1755)) ([2ea7891](https://gitea.com/gitea/runner/commits/2ea78917879854128154124d3b7c8fa8d552e6e8))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1791](https://gitea.com/gitea/runner/issues/1791)) ([65ef31f](https://gitea.com/gitea/runner/commits/65ef31f102ceb75623973921099454637bab55b0))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1817](https://gitea.com/gitea/runner/issues/1817)) ([16c574c](https://gitea.com/gitea/runner/commits/16c574cd266ffe11ae48a37ccf3f6582ce90de21))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1858](https://gitea.com/gitea/runner/issues/1858)) ([b92d95f](https://gitea.com/gitea/runner/commits/b92d95f89931549a235e5f141ecb0a053f4216d3))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1880](https://gitea.com/gitea/runner/issues/1880)) ([a0307d3](https://gitea.com/gitea/runner/commits/a0307d3b7cd20437000d3931149fe600fe515a0f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1882](https://gitea.com/gitea/runner/issues/1882)) ([3d1d8a9](https://gitea.com/gitea/runner/commits/3d1d8a9acac74363f09790f0efa450a27dbeb1ce))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1884](https://gitea.com/gitea/runner/issues/1884)) ([4c8da85](https://gitea.com/gitea/runner/commits/4c8da8558d0a741bb2cecd6d6edf832dc1a055a9))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.12.5 ([#2188](https://gitea.com/gitea/runner/issues/2188)) ([6dcf9bc](https://gitea.com/gitea/runner/commits/6dcf9bc6e62d093fe8152edced9c70ba686d603d))
* **deps:** bump github.com/moby/buildkit from 0.11.6 to 0.12.0 ([#1917](https://gitea.com/gitea/runner/issues/1917)) ([4810f69](https://gitea.com/gitea/runner/commits/4810f69367a141f9cc282d21ab61bb598d5ec584))
* **deps:** bump github.com/moby/buildkit from 0.12.0 to 0.12.1 ([#1947](https://gitea.com/gitea/runner/issues/1947)) ([c812715](https://gitea.com/gitea/runner/commits/c8127155bc64224a1392d2264b8d432328edea57))
* **deps:** bump github.com/moby/buildkit from 0.12.1 to 0.12.2 ([#1986](https://gitea.com/gitea/runner/issues/1986)) ([0c12273](https://gitea.com/gitea/runner/commits/0c12273ebaa937c52dea06d91b8985304874ab14))
* **deps:** bump github.com/moby/buildkit from 0.12.2 to 0.12.3 ([#2060](https://gitea.com/gitea/runner/issues/2060)) ([db6e477](https://gitea.com/gitea/runner/commits/db6e477e25808fcab235fa28b941befdd83c5c03))
* **deps:** bump github.com/moby/buildkit from 0.12.3 to 0.12.4 ([#2113](https://gitea.com/gitea/runner/issues/2113)) ([899a1f2](https://gitea.com/gitea/runner/commits/899a1f206ee6c8636b14a0bc7dc823602365ed9c))
* **deps:** bump github.com/moby/buildkit from 0.12.4 to 0.12.5 ([#2186](https://gitea.com/gitea/runner/issues/2186)) ([46dc2ff](https://gitea.com/gitea/runner/commits/46dc2ffe803f783c3b9fc2190c6566a81d8645dc))
* **deps:** bump github.com/moby/buildkit from 0.9.0 to 0.9.1 ([#838](https://gitea.com/gitea/runner/issues/838)) ([b5e8a18](https://gitea.com/gitea/runner/commits/b5e8a18683c77c7035c704bb7cec980401c73e18))
* **deps:** bump github.com/moby/buildkit from 0.9.1 to 0.9.2 ([#864](https://gitea.com/gitea/runner/issues/864)) ([f8b3563](https://gitea.com/gitea/runner/commits/f8b35631028c7705eb90b1d97300dc19c739cc72))
* **deps:** bump github.com/moby/buildkit from 0.9.2 to 0.9.3 ([#958](https://gitea.com/gitea/runner/issues/958)) ([dba7f83](https://gitea.com/gitea/runner/commits/dba7f8379f327232a1726b336d2cf4a6c3e80ccc))
* **deps:** bump github.com/moby/buildkit from 0.9.3 to 0.10.0 ([#1036](https://gitea.com/gitea/runner/issues/1036)) ([56486f1](https://gitea.com/gitea/runner/commits/56486f14b4a094542d2cf4a903676cd9f14220ec))
* **deps:** bump github.com/moby/patternmatcher from 0.5.0 to 0.6.0 ([#1985](https://gitea.com/gitea/runner/issues/1985)) ([323bee9](https://gitea.com/gitea/runner/commits/323bee9ab5a73473236432e0b89e4d22100cd24d))
* **deps:** bump github.com/opencontainers/image-spec ([#1649](https://gitea.com/gitea/runner/issues/1649)) ([04d12b0](https://gitea.com/gitea/runner/commits/04d12b0206a62cc48a765ed29f2a1ad732d33b3b))
* **deps:** bump github.com/opencontainers/image-spec ([#1781](https://gitea.com/gitea/runner/issues/1781)) ([15eb1fa](https://gitea.com/gitea/runner/commits/15eb1fa92a494f0abde33ad1bfe63496195db9ec))
* **deps:** bump github.com/opencontainers/image-spec ([#1807](https://gitea.com/gitea/runner/issues/1807)) ([1eacf23](https://gitea.com/gitea/runner/commits/1eacf23dcbbea2e5a7d17e8c325ce94ee0b1d3d1))
* **deps:** bump github.com/opencontainers/image-spec ([#1829](https://gitea.com/gitea/runner/issues/1829)) ([26f1b1f](https://gitea.com/gitea/runner/commits/26f1b1f4b6df887a08af0ed6e692962989f5a57c))
* **deps:** bump github.com/opencontainers/image-spec ([#1871](https://gitea.com/gitea/runner/issues/1871)) ([ef60885](https://gitea.com/gitea/runner/commits/ef608854d0b42704a3316c1b09c74c4fff5c40ae))
* **deps:** bump github.com/opencontainers/image-spec ([#1881](https://gitea.com/gitea/runner/issues/1881)) ([f9dcb58](https://gitea.com/gitea/runner/commits/f9dcb58db23061c016bba6712080275219f7008c))
* **deps:** bump github.com/opencontainers/image-spec ([#1883](https://gitea.com/gitea/runner/issues/1883)) ([70e1e37](https://gitea.com/gitea/runner/commits/70e1e372804054374c222c7de1e066bf7ab92b31))
* **deps:** bump github.com/opencontainers/image-spec ([#1893](https://gitea.com/gitea/runner/issues/1893)) ([f40d0b8](https://gitea.com/gitea/runner/commits/f40d0b873da435b6a1dc9fbb5391ce31b2d2dc0e))
* **deps:** bump github.com/opencontainers/image-spec ([#2010](https://gitea.com/gitea/runner/issues/2010)) ([83af8f8](https://gitea.com/gitea/runner/commits/83af8f8767091542dcc4263ebea06e6cc4edd7ce))
* **deps:** bump github.com/opencontainers/image-spec ([#2167](https://gitea.com/gitea/runner/issues/2167)) ([36e0261](https://gitea.com/gitea/runner/commits/36e0261150e54c2ebe727d526a9c82729d6023f7))
* **deps:** bump github.com/opencontainers/image-spec ([#2214](https://gitea.com/gitea/runner/issues/2214)) ([27eb79b](https://gitea.com/gitea/runner/commits/27eb79bb5a6bddb67393432371eed24a4fe242fe))
* **deps:** bump github.com/opencontainers/image-spec ([#2223](https://gitea.com/gitea/runner/issues/2223)) ([9e6ce12](https://gitea.com/gitea/runner/commits/9e6ce121bc1f41e1ecb2c3db7b35f79af70d4ae2))
* **deps:** bump github.com/opencontainers/image-spec ([#954](https://gitea.com/gitea/runner/issues/954)) ([eedac17](https://gitea.com/gitea/runner/commits/eedac179c9f314ee38eda7482c45664df276e546))
* **deps:** bump github.com/opencontainers/runc from 1.1.3 to 1.1.5 ([#1708](https://gitea.com/gitea/runner/issues/1708)) ([6745999](https://gitea.com/gitea/runner/commits/6745999c105a1bb4429275103dc89b6a704b421c))
* **deps:** bump github.com/opencontainers/runc from 1.1.7 to 1.1.12 ([#2187](https://gitea.com/gitea/runner/issues/2187)) ([0dbf44c](https://gitea.com/gitea/runner/commits/0dbf44c657d28b6bdcb7b70c1da59bc925a658db))
* **deps:** bump github.com/opencontainers/selinux ([#1102](https://gitea.com/gitea/runner/issues/1102)) ([e0044ba](https://gitea.com/gitea/runner/commits/e0044ba913149be76acd64ca7d1a66f27532f4c7))
* **deps:** bump github.com/opencontainers/selinux ([#1372](https://gitea.com/gitea/runner/issues/1372)) ([d587ce0](https://gitea.com/gitea/runner/commits/d587ce08e9b04738fd1d3a2433e1f8b4d44c4aed))
* **deps:** bump github.com/opencontainers/selinux ([#1625](https://gitea.com/gitea/runner/issues/1625)) ([8048baf](https://gitea.com/gitea/runner/commits/8048baf435665cb3293f593a06526dfff1691b54))
* **deps:** bump github.com/rhysd/actionlint from 1.6.10 to 1.6.11 ([#1106](https://gitea.com/gitea/runner/issues/1106)) ([065d630](https://gitea.com/gitea/runner/commits/065d630af01b4ed2a0c797d54933ea3b65c3cb38))
* **deps:** bump github.com/rhysd/actionlint from 1.6.11 to 1.6.12 ([#1126](https://gitea.com/gitea/runner/issues/1126)) ([de1dea9](https://gitea.com/gitea/runner/commits/de1dea93795ebf78e0ca531c35969f125035e8c4))
* **deps:** bump github.com/rhysd/actionlint from 1.6.12 to 1.6.13 ([#1173](https://gitea.com/gitea/runner/issues/1173)) ([166d063](https://gitea.com/gitea/runner/commits/166d0630593b40a288cc0190d4aa3e20b31d1461))
* **deps:** bump github.com/rhysd/actionlint from 1.6.13 to 1.6.14 ([#1234](https://gitea.com/gitea/runner/issues/1234)) ([9499612](https://gitea.com/gitea/runner/commits/9499612b5818774c5c6d73fbf2596ed464b7b01d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.14 to 1.6.15 ([#1239](https://gitea.com/gitea/runner/issues/1239)) ([4d9d6ec](https://gitea.com/gitea/runner/commits/4d9d6ecc924bef6c5976acd7f224558e21fa5b5f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.15 to 1.6.16 ([#1311](https://gitea.com/gitea/runner/issues/1311)) ([cf9d82f](https://gitea.com/gitea/runner/commits/cf9d82f69daab9dc907b893a38307a17f444c529))
* **deps:** bump github.com/rhysd/actionlint from 1.6.16 to 1.6.17 ([#1321](https://gitea.com/gitea/runner/issues/1321)) ([9b14698](https://gitea.com/gitea/runner/commits/9b146980dfabd14e9b9defea8f894ad8ab24132a))
* **deps:** bump github.com/rhysd/actionlint from 1.6.17 to 1.6.18 ([#1352](https://gitea.com/gitea/runner/issues/1352)) ([053bff1](https://gitea.com/gitea/runner/commits/053bff19a710669308da5e356ab028bab2a70b5d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.18 to 1.6.19 ([#1360](https://gitea.com/gitea/runner/issues/1360)) ([3fccb5a](https://gitea.com/gitea/runner/commits/3fccb5ae3467fbd013ca72ed08732d2916beb1c2))
* **deps:** bump github.com/rhysd/actionlint from 1.6.19 to 1.6.20 ([#1371](https://gitea.com/gitea/runner/issues/1371)) ([bec45d0](https://gitea.com/gitea/runner/commits/bec45d0dc46e2ccae81c521f3db313b04b0c0831))
* **deps:** bump github.com/rhysd/actionlint from 1.6.20 to 1.6.21 ([#1380](https://gitea.com/gitea/runner/issues/1380)) ([c4be1cf](https://gitea.com/gitea/runner/commits/c4be1cfd6d8ab6664c3a3ffd7c4e5e253fa0b77f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.21 to 1.6.22 ([#1427](https://gitea.com/gitea/runner/issues/1427)) ([8036f49](https://gitea.com/gitea/runner/commits/8036f49fb220bb2315ecba9707242cc0a8533aba))
* **deps:** bump github.com/rhysd/actionlint from 1.6.22 to 1.6.23 ([#1574](https://gitea.com/gitea/runner/issues/1574)) ([1b7212f](https://gitea.com/gitea/runner/commits/1b7212f5fd1e85d97a7119ab098941bd4cf5255b))
* **deps:** bump github.com/rhysd/actionlint from 1.6.23 to 1.6.24 ([#1721](https://gitea.com/gitea/runner/issues/1721)) ([5a331d2](https://gitea.com/gitea/runner/commits/5a331d2d9935565eabdf5a09226b7ba59a029133))
* **deps:** bump github.com/rhysd/actionlint from 1.6.24 to 1.6.25 ([#1870](https://gitea.com/gitea/runner/issues/1870)) ([4fc176f](https://gitea.com/gitea/runner/commits/4fc176f556dcf41e69854570ce65680f127a9357))
* **deps:** bump github.com/rhysd/actionlint from 1.6.25 to 1.6.26 ([#2026](https://gitea.com/gitea/runner/issues/2026)) ([2be4def](https://gitea.com/gitea/runner/commits/2be4def7be0deda86b7e0a00aa06be5e87ac599c))
* **deps:** bump github.com/rhysd/actionlint from 1.6.26 to 1.6.27 ([#2228](https://gitea.com/gitea/runner/issues/2228)) ([68597d6](https://gitea.com/gitea/runner/commits/68597d68f6cc3c384bfd458ab4ec6d51f63a23f6))
* **deps:** bump github.com/rhysd/actionlint from 1.6.8 to 1.6.9 ([#1018](https://gitea.com/gitea/runner/issues/1018)) ([9a54c99](https://gitea.com/gitea/runner/commits/9a54c99ac76a2667551bfd85c5aabcabc969e9dd))
* **deps:** bump github.com/rhysd/actionlint from 1.6.9 to 1.6.10 ([#1039](https://gitea.com/gitea/runner/issues/1039)) ([ee92e07](https://gitea.com/gitea/runner/commits/ee92e07daa0b82e4758807f424f79194e5cb857b))
* **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1273](https://gitea.com/gitea/runner/issues/1273)) ([17dd54d](https://gitea.com/gitea/runner/commits/17dd54d69231884ea561fe2bd0c0e059debd2658))
* **deps:** bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2 ([#1814](https://gitea.com/gitea/runner/issues/1814)) ([06054d2](https://gitea.com/gitea/runner/commits/06054d2931861d2f02cd4028d04164a4d52bbc1c))
* **deps:** bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3 ([#1841](https://gitea.com/gitea/runner/issues/1841)) ([e150310](https://gitea.com/gitea/runner/commits/e150310da77918e9ca904363f6cd80907877ff12))
* **deps:** bump github.com/spf13/cobra from 1.2.1 to 1.3.0 ([#955](https://gitea.com/gitea/runner/issues/955)) ([e9e6dda](https://gitea.com/gitea/runner/commits/e9e6ddadf567beb62ddb94b3fe472297081d4077))
* **deps:** bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#1040](https://gitea.com/gitea/runner/issues/1040)) ([fa71f0b](https://gitea.com/gitea/runner/commits/fa71f0b1e04a6391b77ef091c0db4d1a8f0ad1c1))
* **deps:** bump github.com/spf13/cobra from 1.4.0 to 1.5.0 ([#1233](https://gitea.com/gitea/runner/issues/1233)) ([4deb213](https://gitea.com/gitea/runner/commits/4deb21344d40dd4c029d5c1048a7f1d557ab8774))
* **deps:** bump github.com/spf13/cobra from 1.5.0 to 1.6.0 ([#1394](https://gitea.com/gitea/runner/issues/1394)) ([52c6edf](https://gitea.com/gitea/runner/commits/52c6edf0a48fa8287eaaa1ea28d89f65d05e8e7b))
* **deps:** bump github.com/spf13/cobra from 1.6.0 to 1.6.1 ([#1420](https://gitea.com/gitea/runner/issues/1420)) ([a9298d9](https://gitea.com/gitea/runner/commits/a9298d936c5265461dd8f6322961232b08fc6d12))
* **deps:** bump github.com/spf13/cobra from 1.6.1 to 1.7.0 ([#1724](https://gitea.com/gitea/runner/issues/1724)) ([c5ce502](https://gitea.com/gitea/runner/commits/c5ce502f9f9e98ea76a09858075e3ee91da6bd10))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([#2079](https://gitea.com/gitea/runner/issues/2079)) ([55b09a0](https://gitea.com/gitea/runner/commits/55b09a04cd14006b34eef02f0d19120759ee167b))
* **deps:** bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#1061](https://gitea.com/gitea/runner/issues/1061)) ([7d946c4](https://gitea.com/gitea/runner/commits/7d946c43e208aa33a0f6183a2b84a8e9f33425ae))
* **deps:** bump github.com/stretchr/testify from 1.7.1 to 1.7.2 ([#1213](https://gitea.com/gitea/runner/issues/1213)) ([3415347](https://gitea.com/gitea/runner/commits/3415347efcdf6d6d7544c5df7cb33e6db043ff4a))
* **deps:** bump github.com/stretchr/testify from 1.7.2 to 1.7.5 ([#1232](https://gitea.com/gitea/runner/issues/1232)) ([c67abf2](https://gitea.com/gitea/runner/commits/c67abf2401ba57097e862b524b86de7171f1c609))
* **deps:** bump github.com/stretchr/testify from 1.7.5 to 1.8.0 ([#1240](https://gitea.com/gitea/runner/issues/1240)) ([50f0b0e](https://gitea.com/gitea/runner/commits/50f0b0e7f49d368583a30e515fb16db58c357bbc))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1404](https://gitea.com/gitea/runner/issues/1404)) ([db850e8](https://gitea.com/gitea/runner/commits/db850e818c2dd36f0021ac0a227ed99d2daafff2))
* **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#1650](https://gitea.com/gitea/runner/issues/1650)) ([a6c95ef](https://gitea.com/gitea/runner/commits/a6c95ef2a771afdc9f83f6750b642491e52f681c))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://gitea.com/gitea/runner/issues/1842)) ([c70a674](https://gitea.com/gitea/runner/commits/c70a6743f62af2c223f62c2fddd6a1a68271a50e))
* **deps:** bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#2235](https://gitea.com/gitea/runner/issues/2235)) ([934b13a](https://gitea.com/gitea/runner/commits/934b13a7a126198019800e79fb2ae7fe57a1c2bd))
* **deps:** bump go.etcd.io/bbolt from 1.3.7 to 1.3.8 ([#2065](https://gitea.com/gitea/runner/issues/2065)) ([83bfbcd](https://gitea.com/gitea/runner/commits/83bfbcdafd55c5dc6989f6a6b0aa07cab5b99ea8))
* **deps:** bump go.etcd.io/bbolt from 1.3.8 to 1.3.9 ([#2229](https://gitea.com/gitea/runner/issues/2229)) ([c79f59f](https://gitea.com/gitea/runner/commits/c79f59f802673f00911bea93db15b83f5bf3507b))
* **deps:** bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#2134](https://gitea.com/gitea/runner/issues/2134)) ([238a495](https://gitea.com/gitea/runner/commits/238a495579abc12665c1fbf01e8a4f16c2c96fe1))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([#2045](https://gitea.com/gitea/runner/issues/2045)) ([99067a9](https://gitea.com/gitea/runner/commits/99067a9c1ed3b9cdb75797a7e87e6ef0471105f8))
* **deps:** bump golang.org/x/net from 0.4.0 to 0.7.0 ([#1636](https://gitea.com/gitea/runner/issues/1636)) ([b7a9eb9](https://gitea.com/gitea/runner/commits/b7a9eb9fbf86da76cbad35296ff11b64ee416d18))
* **deps:** bump golang.org/x/term from 0.10.0 to 0.11.0 ([#1948](https://gitea.com/gitea/runner/issues/1948)) ([74b74e8](https://gitea.com/gitea/runner/commits/74b74e847bad3c71ad7a1df76384a804bb46d644))
* **deps:** bump golang.org/x/term from 0.11.0 to 0.12.0 ([#1992](https://gitea.com/gitea/runner/issues/1992)) ([1d32507](https://gitea.com/gitea/runner/commits/1d32507b52c8a3c502415d1350aeca0eb9dcfe81))
* **deps:** bump golang.org/x/term from 0.12.0 to 0.13.0 ([#2036](https://gitea.com/gitea/runner/issues/2036)) ([3c730d7](https://gitea.com/gitea/runner/commits/3c730d7924ed9da3ce318bba982873b4484106a7))
* **deps:** bump golang.org/x/term from 0.13.0 to 0.14.0 ([#2091](https://gitea.com/gitea/runner/issues/2091)) ([4dcb9b7](https://gitea.com/gitea/runner/commits/4dcb9b7a139f7399b2778575e61a398aefb41d85))
* **deps:** bump golang.org/x/term from 0.14.0 to 0.15.0 ([#2112](https://gitea.com/gitea/runner/issues/2112)) ([2b47c99](https://gitea.com/gitea/runner/commits/2b47c99bb7b5d73c852ca7241cda3e1e70ba3670))
* **deps:** bump golang.org/x/term from 0.15.0 to 0.16.0 ([#2148](https://gitea.com/gitea/runner/issues/2148)) ([ef5746b](https://gitea.com/gitea/runner/commits/ef5746ba747a4a39770769725a62b519d8d919aa))
* **deps:** bump golang.org/x/term from 0.16.0 to 0.17.0 ([#2213](https://gitea.com/gitea/runner/issues/2213)) ([8277584](https://gitea.com/gitea/runner/commits/8277584f00a0fa82598fc795f35b165e0bfa3bdf))
* **deps:** bump golang.org/x/term from 0.17.0 to 0.18.0 ([#2244](https://gitea.com/gitea/runner/issues/2244)) ([8c83d57](https://gitea.com/gitea/runner/commits/8c83d5721258537b3a47dfc060be67d180ba6a0e))
* **deps:** bump golang.org/x/term from 0.3.0 to 0.4.0 ([#1577](https://gitea.com/gitea/runner/issues/1577)) ([a2856ae](https://gitea.com/gitea/runner/commits/a2856aed75e6d4e5607dc59ff074c9612553360f))
* **deps:** bump golang.org/x/term from 0.5.0 to 0.6.0 ([#1666](https://gitea.com/gitea/runner/issues/1666)) ([8dbd151](https://gitea.com/gitea/runner/commits/8dbd151fa71031c737d3f35580906f4028c355ca))
* **deps:** bump golang.org/x/term from 0.6.0 to 0.7.0 ([#1723](https://gitea.com/gitea/runner/issues/1723)) ([3ae0a9d](https://gitea.com/gitea/runner/commits/3ae0a9dfb7d9e4eb2adb3dc4a778290784fc5ce8))
* **deps:** bump golang.org/x/term from 0.8.0 to 0.9.0 ([#1869](https://gitea.com/gitea/runner/issues/1869)) ([b2191ae](https://gitea.com/gitea/runner/commits/b2191ae204a713a5e9bd4b3e1ec34b23b7fd34a6))
* **deps:** bump golang.org/x/term from 0.9.0 to 0.10.0 ([#1906](https://gitea.com/gitea/runner/issues/1906)) ([e343ea9](https://gitea.com/gitea/runner/commits/e343ea9d5f6fa478a364b1f502b4859eed63fa38))
* **deps:** bump golangci/golangci-lint-action from 2 to 3.1.0 ([#1019](https://gitea.com/gitea/runner/issues/1019)) ([7d403b8](https://gitea.com/gitea/runner/commits/7d403b83075932a18007210a62c6d0ce64df11a6))
* **deps:** bump golangci/golangci-lint-action from 3.1.0 to 3.2.0 ([#1169](https://gitea.com/gitea/runner/issues/1169)) ([8994773](https://gitea.com/gitea/runner/commits/89947735d508751b9e84d653635dec91299e9eec))
* **deps:** bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 ([#1408](https://gitea.com/gitea/runner/issues/1408)) ([925ac20](https://gitea.com/gitea/runner/commits/925ac2098b569f95beca47d5aac7e2287922963b))
* **deps:** bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 ([#1436](https://gitea.com/gitea/runner/issues/1436)) ([1b554ae](https://gitea.com/gitea/runner/commits/1b554aeff651cf801c11277bacb162ba31483b08))
* **deps:** bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 ([#1590](https://gitea.com/gitea/runner/issues/1590)) ([70ae63b](https://gitea.com/gitea/runner/commits/70ae63bd10eb4ac3a9a64633b88242f9a10122e3))
* **deps:** bump golangci/golangci-lint-action from 3.4.0 to 3.5.0 ([#1839](https://gitea.com/gitea/runner/issues/1839)) ([34db861](https://gitea.com/gitea/runner/commits/34db86138ba3f375c3797dfe5b2952f32397f7b1))
* **deps:** bump golangci/golangci-lint-action from 3.5.0 to 3.6.0 ([#1868](https://gitea.com/gitea/runner/issues/1868)) ([79f93be](https://gitea.com/gitea/runner/commits/79f93beef21650b25ab20a24aa91032e3a6b597e))
* **deps:** bump golangci/golangci-lint-action from 3.6.0 to 3.7.0 ([#1978](https://gitea.com/gitea/runner/issues/1978)) ([f64c267](https://gitea.com/gitea/runner/commits/f64c267dacd026b306da8ee405e84df56d124e69))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1179](https://gitea.com/gitea/runner/issues/1179)) ([5054455](https://gitea.com/gitea/runner/commits/50544556f8672c5553fd5241f4adb157698c781b))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1515](https://gitea.com/gitea/runner/issues/1515)) ([c044035](https://gitea.com/gitea/runner/commits/c044035d9ef358482efc996fba0a4c3f7f5f3847))
* **deps:** bump goreleaser/goreleaser-action from 4 to 5 ([#2008](https://gitea.com/gitea/runner/issues/2008)) ([2911b21](https://gitea.com/gitea/runner/commits/2911b2172cee57abbfdafd8ffd04fe4b6ccdd9dc))
* **deps:** bump gotest.tools/v3 from 3.0.3 to 3.4.0 ([#1379](https://gitea.com/gitea/runner/issues/1379)) ([c7ccb0d](https://gitea.com/gitea/runner/commits/c7ccb0dde9ab6332af7940e6f61e3ab73685fd07))
* **deps:** bump gotest.tools/v3 from 3.4.0 to 3.5.0 ([#1892](https://gitea.com/gitea/runner/issues/1892)) ([70956f2](https://gitea.com/gitea/runner/commits/70956f2929eddc7a65ffc16ac335db9472924d66))
* **deps:** bump gotest.tools/v3 from 3.5.0 to 3.5.1 ([#2009](https://gitea.com/gitea/runner/issues/2009)) ([60060a7](https://gitea.com/gitea/runner/commits/60060a7c9c7730d653bd8d980f1a17abd068e739))
* **deps:** bump megalinter/megalinter from 5 to 6.11.0 ([#1370](https://gitea.com/gitea/runner/issues/1370)) ([79384c3](https://gitea.com/gitea/runner/commits/79384c35ab42715144388b7d3e0b0290fb93a962))
* **deps:** bump megalinter/megalinter from 6.11.0 to 6.11.1 ([#1381](https://gitea.com/gitea/runner/issues/1381)) ([4541139](https://gitea.com/gitea/runner/commits/4541139109f2838bcac25ddd3d435a5c0661bf25))
* **deps:** bump megalinter/megalinter from 6.11.1 to 6.12.0 ([#1395](https://gitea.com/gitea/runner/issues/1395)) ([37f5b7f](https://gitea.com/gitea/runner/commits/37f5b7f4a2bc363be670ee65f623c6b54ec9d8c5))
* **deps:** bump megalinter/megalinter from 6.12.0 to 6.13.0 ([#1407](https://gitea.com/gitea/runner/issues/1407)) ([e81dc0d](https://gitea.com/gitea/runner/commits/e81dc0d295ad58d13927e61c91d3c4b4a6def942))
* **deps:** bump megalinter/megalinter from 6.13.0 to 6.14.0 ([#1426](https://gitea.com/gitea/runner/issues/1426)) ([c232981](https://gitea.com/gitea/runner/commits/c2329815b8bbc739043e1584e465322a602475e9))
* **deps:** bump megalinter/megalinter from 6.14.0 to 6.15.0 ([#1475](https://gitea.com/gitea/runner/issues/1475)) ([7754ba7](https://gitea.com/gitea/runner/commits/7754ba7fcf3f80b76bb4d7cf1edbe5935c8a6bdc))
* **deps:** bump megalinter/megalinter from 6.15.0 to 6.16.0 ([#1534](https://gitea.com/gitea/runner/issues/1534)) ([e3a722e](https://gitea.com/gitea/runner/commits/e3a722e7c6c697a1b0a38aa314bfb227da4deb5f))
* **deps:** bump megalinter/megalinter from 6.16.0 to 6.17.0 ([#1540](https://gitea.com/gitea/runner/issues/1540)) ([cfc1b91](https://gitea.com/gitea/runner/commits/cfc1b91aa11220f102470622cacd105cedd931c9))
* **deps:** bump megalinter/megalinter from 6.17.0 to 6.18.0 ([#1550](https://gitea.com/gitea/runner/issues/1550)) ([f0c6fa1](https://gitea.com/gitea/runner/commits/f0c6fa11be88093fbcb7ecb1cbebfee8874008cc))
* **deps:** bump megalinter/megalinter from 6.18.0 to 6.19.0 ([#1610](https://gitea.com/gitea/runner/issues/1610)) ([460c78d](https://gitea.com/gitea/runner/commits/460c78da0773c226162d93623e8932715a22e443))
* **deps:** bump megalinter/megalinter from 6.19.0 to 6.20.0 ([#1665](https://gitea.com/gitea/runner/issues/1665)) ([7ee2138](https://gitea.com/gitea/runner/commits/7ee2138418822460ed576391201db55803bd7069))
* **deps:** bump megalinter/megalinter from 6.20.0 to 6.20.1 ([#1679](https://gitea.com/gitea/runner/issues/1679)) ([5a2112a](https://gitea.com/gitea/runner/commits/5a2112a7f8e282fcb75f1c73caa1582811ff4890))
* **deps:** bump megalinter/megalinter from 6.20.1 to 6.21.0 ([#1699](https://gitea.com/gitea/runner/issues/1699)) ([775a128](https://gitea.com/gitea/runner/commits/775a128cd4660be4f98da019c52a2ab9a5080433))
* **deps:** bump megalinter/megalinter from 6.21.0 to 6.22.1 ([#1710](https://gitea.com/gitea/runner/issues/1710)) ([ad7a8a3](https://gitea.com/gitea/runner/commits/ad7a8a3559d937d79e9425b8d29714752879b823))
* **deps:** bump megalinter/megalinter from 6.22.1 to 6.22.2 ([#1720](https://gitea.com/gitea/runner/issues/1720)) ([4ab812c](https://gitea.com/gitea/runner/commits/4ab812c1b10bc16a54dd208963d6b60e2835d328))
* **deps:** bump megalinter/megalinter from 6.22.2 to 7.0.2 ([#1827](https://gitea.com/gitea/runner/issues/1827)) ([481999f](https://gitea.com/gitea/runner/commits/481999f59da465d2d8743fd3d97e24bb6e1a2311))
* **deps:** bump megalinter/megalinter from 7.0.2 to 7.0.4 ([#1838](https://gitea.com/gitea/runner/issues/1838)) ([04b5c73](https://gitea.com/gitea/runner/commits/04b5c739d2649b67565dd7addbe2bf5fc910e2a6))
* **deps:** bump megalinter/megalinter from 7.0.4 to 7.1.0 ([#1857](https://gitea.com/gitea/runner/issues/1857)) ([3555d65](https://gitea.com/gitea/runner/commits/3555d65d08bc9c66a6e4541b2064a082af1b6b3e))
* **deps:** bump megalinter/megalinter from 7.1.0 to 7.2.1 ([#1931](https://gitea.com/gitea/runner/issues/1931)) ([3939f48](https://gitea.com/gitea/runner/commits/3939f48e6def83e8a6ad7358a848ce11912c3798))
* **deps:** bump megalinter/megalinter from 7.2.1 to 7.3.0 ([#1965](https://gitea.com/gitea/runner/issues/1965)) ([2a0a0a1](https://gitea.com/gitea/runner/commits/2a0a0a1a62b8c72cd891208f085c00cb00402b57))
* **deps:** bump megalinter/megalinter from 7.3.0 to 7.4.0 ([#2025](https://gitea.com/gitea/runner/issues/2025)) ([3d47885](https://gitea.com/gitea/runner/commits/3d47885894b804abaf5cebde36d9b9f260ae2c50))
* **deps:** bump megalinter/megalinter from 7.4.0 to 7.5.0 ([#2070](https://gitea.com/gitea/runner/issues/2070)) ([3e5c629](https://gitea.com/gitea/runner/commits/3e5c62977fac53931267a86493e8c5c292508e30))
* **deps:** bump megalinter/megalinter from 7.5.0 to 7.6.0 ([#2098](https://gitea.com/gitea/runner/issues/2098)) ([bd10c9a](https://gitea.com/gitea/runner/commits/bd10c9a80191e7f3c1eed0ee42847f12bc4a81ca))
* **deps:** bump megalinter/megalinter from 7.6.0 to 7.7.0 ([#2119](https://gitea.com/gitea/runner/issues/2119)) ([29f4123](https://gitea.com/gitea/runner/commits/29f4123b5ced89ad39f9f80889576810cfdecf7f))
* **deps:** bump megalinter/megalinter from 7.7.0 to 7.8.0 ([#2164](https://gitea.com/gitea/runner/issues/2164)) ([c215e08](https://gitea.com/gitea/runner/commits/c215e0888aa29780903c113e35c46f495d4157fc))
* **deps:** update go-git to ce62f3e9ff86270538a514a68d3bd5563a733e3b ([#1725](https://gitea.com/gitea/runner/issues/1725)) ([148a545](https://gitea.com/gitea/runner/commits/148a54502167880f90dbf2f68ff44478679f41dd))

### Continuous Integration

* Add snapshot artifacts to PR builds ([#574](https://gitea.com/gitea/runner/issues/574)) ([da07324](https://gitea.com/gitea/runner/commits/da07324779854fda8bff4bd45452a04ce5e68a1e))
* add Winget Releaser workflow ([#1623](https://gitea.com/gitea/runner/issues/1623)) ([34ab815](https://gitea.com/gitea/runner/commits/34ab8150bf8b438964b7e52796f2123466e0cf76))
* always grab latest go version ([#1124](https://gitea.com/gitea/runner/issues/1124)) ([25ebaf4](https://gitea.com/gitea/runner/commits/25ebaf4e7b507ffc7cabdb2a560348af5bec4ac3))
* automatic merge PRs created by a maintainer and approved by 1 other maintainer ([#2156](https://gitea.com/gitea/runner/issues/2156)) ([cd40f3f](https://gitea.com/gitea/runner/commits/cd40f3fe9bd86f7c2c5496686e119d047e9c04a9))
* bump go to 1.18 ([#1087](https://gitea.com/gitea/runner/issues/1087)) ([11c3fb3](https://gitea.com/gitea/runner/commits/11c3fb39dd5d3907baf8f13ef6953d0a21321cc2))
* change winget job runner to `ubuntu-latest` ([#1959](https://gitea.com/gitea/runner/issues/1959)) ([f55ae1a](https://gitea.com/gitea/runner/commits/f55ae1a0bc4fb14d463b79a62ff9a37d7e5bda68))
* **choco:** fix chocolatey docker build ([#903](https://gitea.com/gitea/runner/issues/903)) ([6ebcac3](https://gitea.com/gitea/runner/commits/6ebcac37710ae5d4eb87f63f31d7e1be8fe3d217))
* **choco:** update chocolatey to 1.1.0 ([#1164](https://gitea.com/gitea/runner/issues/1164)) ([507ae61](https://gitea.com/gitea/runner/commits/507ae61d1b3ac3f66273a00fad1321af8c6df260))
* deduplicate running workflows ([#1751](https://gitea.com/gitea/runner/issues/1751)) ([7c6237d](https://gitea.com/gitea/runner/commits/7c6237d93f01b2c50a98a488e40ab7f9a719da35))
* **golangci-lint:** switch to latest ([#870](https://gitea.com/gitea/runner/issues/870)) ([cd24df5](https://gitea.com/gitea/runner/commits/cd24df5727c54de6671f83b44b7527f3c74f122f))
* make stalebot less annoying ([#1587](https://gitea.com/gitea/runner/issues/1587)) ([01952d0](https://gitea.com/gitea/runner/commits/01952d05c858c70ec0f007d90e394cf4cdfcd83b))
* **Mergify:** configuration update ([#924](https://gitea.com/gitea/runner/issues/924)) ([cad4bc8](https://gitea.com/gitea/runner/commits/cad4bc8c367fd4e623dfe2820bd5f3daca80bf58))
* **Mergify:** configuration update ([#966](https://gitea.com/gitea/runner/issues/966)) ([557dc75](https://gitea.com/gitea/runner/commits/557dc755b8fe6b286d15bd7a1091f3731dc44791))
* **Mergify:** configuration update ([#972](https://gitea.com/gitea/runner/issues/972)) ([4f8da0a](https://gitea.com/gitea/runner/commits/4f8da0a51cd44753f3e4ea0364877f6d2aa11f1d))
* **mergify:** fix deprecated strict mode ([#862](https://gitea.com/gitea/runner/issues/862)) ([80db062](https://gitea.com/gitea/runner/commits/80db062472a72735f5fa454a1f3977ea8702aad4))
* pin goreleaser version ([#825](https://gitea.com/gitea/runner/issues/825)) ([b1e4844](https://gitea.com/gitea/runner/commits/b1e4844aac4b8a143312c3ec9d3859362c669f88))
* remove macOS tests ([#860](https://gitea.com/gitea/runner/issues/860)) ([aad724c](https://gitea.com/gitea/runner/commits/aad724c87abb04184c9205643ec4659043637b07))
* replace superlinter with megalinter ([#923](https://gitea.com/gitea/runner/issues/923)) ([9be6a58](https://gitea.com/gitea/runner/commits/9be6a58c0eb706d514e355aa10f2b4c638cfec4a))
* **stalebot:** ignore feature requests ([#1017](https://gitea.com/gitea/runner/issues/1017)) ([932ad0e](https://gitea.com/gitea/runner/commits/932ad0ef6a2483b86f331f933b9c9081326ecbac))
* switch macos job to macos-10.15 ([#849](https://gitea.com/gitea/runner/issues/849)) ([977b526](https://gitea.com/gitea/runner/commits/977b526384a00e4f0804db6ddbaa3248c1918c75))

### Documentation

* add brew installation from HEAD ([#886](https://gitea.com/gitea/runner/issues/886)) ([df4ef4d](https://gitea.com/gitea/runner/commits/df4ef4de808edb8e3b657e4e24c96f8643ffcd32))
* add naming rule for git ref ([70cc6c0](https://gitea.com/gitea/runner/commits/70cc6c017bc8382b2efe4a50266f733037094da8))
* clarifying skipping steps / jobs ([#1480](https://gitea.com/gitea/runner/issues/1480)) ([7073eac](https://gitea.com/gitea/runner/commits/7073eac24057020cd54bcb00402281be55df3246))
* Expand the GITHUB_TOKEN section ([#968](https://gitea.com/gitea/runner/issues/968)) ([b1f5963](https://gitea.com/gitea/runner/commits/b1f5963c8623d5a5a92d2783b30da73c339232e2))
* fix brew formula link ([#846](https://gitea.com/gitea/runner/issues/846)) ([2319eb2](https://gitea.com/gitea/runner/commits/2319eb210f8d097da470331d2d12e7c6db0f8056))
* **README:** add --env-file flag ([#147](https://gitea.com/gitea/runner/issues/147)) ([f6b7613](https://gitea.com/gitea/runner/commits/f6b761378af17296aba8ff3be54e26d56acfeae5))
* **README:** remove go install ([#1123](https://gitea.com/gitea/runner/issues/1123)) ([e15f7a5](https://gitea.com/gitea/runner/commits/e15f7a51c6cf3e76cdc8ef8401b87080b7e9f95e))
* remove help section ([#1648](https://gitea.com/gitea/runner/issues/1648)) ([4b56ace](https://gitea.com/gitea/runner/commits/4b56aced159489639e5ca8226e2b82f5c7e9823e))

### Tests

* add test for networking setup in act ([#1375](https://gitea.com/gitea/runner/issues/1375)) ([ff5e289](https://gitea.com/gitea/runner/commits/ff5e2898041230978cef9af37dd87c09b07a77e8))
* Do not leak step env in composite ([#1585](https://gitea.com/gitea/runner/issues/1585)) ([3748772](https://gitea.com/gitea/runner/commits/3748772201458442afde9a582d70b40db05deaf8))
* Enshure ForcePull config doesn't break docker actions (1661) ([24440d9](https://gitea.com/gitea/runner/commits/24440d9f15673993dc6986f2c64dbd056e4286fe))
* env.ACT in if condition of the workflow ([#965](https://gitea.com/gitea/runner/issues/965)) ([3f0375a](https://gitea.com/gitea/runner/commits/3f0375aeff1973b674b33a0e4c44bc05a41714c5))
* make sure workflow_call is not a github event calling our workflow ([#1520](https://gitea.com/gitea/runner/issues/1520)) ([3f4a6dc](https://gitea.com/gitea/runner/commits/3f4a6dc4626eed03991e207b4d0e343b634edf9e))
* use ping to improve network test ([#2266](https://gitea.com/gitea/runner/issues/2266)) ([f56dd65](https://gitea.com/gitea/runner/commits/f56dd65ff69f38f68bf82d485c3cb1b90e9ffc3e))

### Chores

* add Winget installation instructions ([#1622](https://gitea.com/gitea/runner/issues/1622)) ([c378a7d](https://gitea.com/gitea/runner/commits/c378a7d28b5c20de9fbb4456544329d04ffb62a4))
* bump version ([0d07805](https://gitea.com/gitea/runner/commits/0d0780580cfe63786c21a31500ccb69797cb9cf6))
* bump version ([33cb8bc](https://gitea.com/gitea/runner/commits/33cb8bca649a1b63ac59596b20d4001ca9cbd3e6))
* bump VERSION to 0.2.39 ([efb12b7](https://gitea.com/gitea/runner/commits/efb12b7f122e1daa3fce410a0808098a23dad910))
* bump VERSION to 0.2.40 ([3e23b4f](https://gitea.com/gitea/runner/commits/3e23b4fbb558bc4e8bdc72d4f75e0d94e5be3bf1))
* bump VERSION to 0.2.41 ([b91e4b0](https://gitea.com/gitea/runner/commits/b91e4b0d5576ee0efeb381c9427b791e5272ece9))
* bump VERSION to 0.2.42 ([72d0321](https://gitea.com/gitea/runner/commits/72d03214c5ed8985bf23dc15440fa5301e9d13bf))
* bump VERSION to 0.2.43 ([44b510f](https://gitea.com/gitea/runner/commits/44b510f48cac6a58654bf920ee26cc6e1d36ac3a))
* bump VERSION to 0.2.44 ([220d6f1](https://gitea.com/gitea/runner/commits/220d6f1251aea75c378e5a98546f43cd918e1892))
* bump VERSION to 0.2.45 ([f1df2ca](https://gitea.com/gitea/runner/commits/f1df2ca5d6b315935cf9bb0d2a989175adf4e894))
* bump VERSION to 0.2.46 ([b0d0cec](https://gitea.com/gitea/runner/commits/b0d0cec71f2240fc629dcaa1889ed8c4ad52beb9))
* bump VERSION to 0.2.47 ([310cb79](https://gitea.com/gitea/runner/commits/310cb79e81e07ef11985b2f934092b82b167ccf1))
* bump VERSION to 0.2.48 ([3a0a642](https://gitea.com/gitea/runner/commits/3a0a6425a80f3405670d8d7b11ed18ea56a09af4))
* bump VERSION to 0.2.49 ([d5d8548](https://gitea.com/gitea/runner/commits/d5d854854651c151ecd914bb6e2e370e0057929a))
* bump VERSION to 0.2.50 ([80b0955](https://gitea.com/gitea/runner/commits/80b0955303888742c3ab73af5758bb7b01f5f57c))
* bump VERSION to 0.2.51 ([c241ecd](https://gitea.com/gitea/runner/commits/c241ecda318e4977501971c267f0a2fad70ea5cd))
* bump VERSION to 0.2.52 ([44ea01c](https://gitea.com/gitea/runner/commits/44ea01c2097eaad1596f4dafd6ec84872bd00402))
* bump VERSION to 0.2.53 ([1bb2ee7](https://gitea.com/gitea/runner/commits/1bb2ee7098f338765ce5169c850d0f19e59d3147))
* bump VERSION to 0.2.54 ([7693697](https://gitea.com/gitea/runner/commits/7693697f4cc60143c0561668b1c6fcb4bf733b94))
* bump VERSION to 0.2.55 ([95ff5bf](https://gitea.com/gitea/runner/commits/95ff5bf299a18833b19f157386184539ad6f4e64))
* bump VERSION to 0.2.56 ([0c60f97](https://gitea.com/gitea/runner/commits/0c60f9749a35a451d40e0741839f7bc403080c8d))
* bump VERSION to 0.2.57 ([4fae81e](https://gitea.com/gitea/runner/commits/4fae81efe4cdd9e09e7ef8e874a2d63b1ed98524))
* bump VERSION to 0.2.58 ([3ed38d8](https://gitea.com/gitea/runner/commits/3ed38d8e8b63f435cb845fbe9cdccc7ff032f402))
* bump VERSION to 0.2.59 ([b7a8145](https://gitea.com/gitea/runner/commits/b7a8145d09a7469e03c57e24230cbba109218faf))
* bump VERSION to 0.2.60 ([d3c8664](https://gitea.com/gitea/runner/commits/d3c8664d3d6acb855ae5418f6d02930c98bf877e))
* bump VERSION to 0.2.61 ([361b7e9](https://gitea.com/gitea/runner/commits/361b7e9f1abe9081edf594aa35dede4eab972c70))
* **deps:** cleanup go.mod a little bit ([#136](https://gitea.com/gitea/runner/issues/136)) ([ca648a3](https://gitea.com/gitea/runner/commits/ca648a37c8df80a4e96ffa7c3360c13c1a0715be))
* fix promote script ([b3c06dc](https://gitea.com/gitea/runner/commits/b3c06dc19d5ba396336615cb404cd394bd3e9f30))
* fix release script to trigger gh-act ([787388d](https://gitea.com/gitea/runner/commits/787388daf54f3d78c56da29a37062b019e10e84f))
* fold integration test outputs ([#1111](https://gitea.com/gitea/runner/issues/1111)) ([a85e89d](https://gitea.com/gitea/runner/commits/a85e89d3fb8e1a3b2e9dda843e18aa819095e56e))
* makefile update ([872a695](https://gitea.com/gitea/runner/commits/872a69548c77b07c01ba18fb507ffe0bfe278f5e))
* push new version ([003947e](https://gitea.com/gitea/runner/commits/003947ea33bfabcaa8c3787aa78780fd88b9c1ff))
* Remove obsolete Container.UpdateFromPath ([#1631](https://gitea.com/gitea/runner/issues/1631)) ([21ea3d0](https://gitea.com/gitea/runner/commits/21ea3d0d5fc92cf9ae93c62c0af7ec3a5da36541))
* remove repetitive words ([#2259](https://gitea.com/gitea/runner/issues/2259)) ([3949d74](https://gitea.com/gitea/runner/commits/3949d74af58d72309706e8fc5353de080fae6f3e))
* run act from cli on linux ([#1758](https://gitea.com/gitea/runner/issues/1758)) ([c81a770](https://gitea.com/gitea/runner/commits/c81a770bc5f8041c2c169e227996ea34ae62ec48)), references [#1756](https://gitea.com/gitea/runner/issues/1756)
* update docs for installing act as GH CLI extension ([7ada9d3](https://gitea.com/gitea/runner/commits/7ada9d3f74f4d8705ba01ce110f3d611989dd442))
* update go-git dependency in go.mod ([#30](https://gitea.com/gitea/runner/issues/30)) ([8f12a6c](https://gitea.com/gitea/runner/commits/8f12a6c947deaa3488013503a518144c9ecdd035))
* update LICENSE ([7fb84a5](https://gitea.com/gitea/runner/commits/7fb84a54a8dd93b5874c7edb89861f1c7d23b469))
* upgrade docker v19.03.12 ([#348](https://gitea.com/gitea/runner/issues/348)) ([3eef111](https://gitea.com/gitea/runner/commits/3eef111e46f42f43e62b64e6cdfbade4fe0935e2))
* upgrade golangci-lint and address findings ([#1904](https://gitea.com/gitea/runner/issues/1904)) ([724ec91](https://gitea.com/gitea/runner/commits/724ec918c92c162e537ba879e6d1ccc997f63483))
* use new style octal ([#1630](https://gitea.com/gitea/runner/issues/1630)) ([1316307](https://gitea.com/gitea/runner/commits/131630731374a16dbb7c8ed4893397c52d9c4c98))
* **yaml:** Improve `ParseRawOn` ([#28](https://gitea.com/gitea/runner/issues/28)) ([3daf313](https://gitea.com/gitea/runner/commits/3daf3132057d33bee325857b016a60db332fd050))

### Reverts

* Revert "Erase needs of job in SingleWorkflow (#9)" (#32) ([568f053](https://gitea.com/gitea/runner/commits/568f0537239f4566ad2b05b8866611d5af612b94)), references [#9](https://gitea.com/gitea/runner/issues/9) [#32](https://gitea.com/gitea/runner/issues/32)
* Revert "feat: allow existing logger from context (#859)" (#890) ([4ae71b5](https://gitea.com/gitea/runner/commits/4ae71b5878a8f6906ad3d2cd1d8cef0791527f20)), references [#859](https://gitea.com/gitea/runner/issues/859) [#890](https://gitea.com/gitea/runner/issues/890)
* deprecation of containerArchitecture ([#1514](https://gitea.com/gitea/runner/issues/1514)) ([4989f44](https://gitea.com/gitea/runner/commits/4989f444f1085ac0ec6b9bc35de92dc22f3c9f5d))
* Don't stop JobContainer on failure ([#1041](https://gitea.com/gitea/runner/issues/1041)) ([4d71071](https://gitea.com/gitea/runner/commits/4d7107161e231406bee856258ec847844623224e))

## [0.4.1](https://gitea.com/gitea/runner/compare/v0.4.0...v0.4.1) (2026-04-19)

### Performance Improvements

* use single poller with semaphore-based capacity control ([#822](https://gitea.com/gitea/runner/issues/822)) ([9aafec1](https://gitea.com/gitea/runner/commits/9aafec169b7ebf3d71df2b06780868fb552b435e)), references [#819](https://gitea.com/gitea/runner/issues/819) [#819](https://gitea.com/gitea/runner/issues/819) [pre-#819](https://gitea.com/pre-/issues/819) [post-#819](https://gitea.com/post-/issues/819) [Post-#819](https://gitea.com/Post-/issues/819) [pre-#819](https://gitea.com/pre-/issues/819) [pre-#819](https://gitea.com/pre-/issues/819) [post-#819](https://gitea.com/post-/issues/819) [#819](https://gitea.com/gitea/runner/issues/819) [#819](https://gitea.com/gitea/runner/issues/819)

## [0.4.0](https://gitea.com/gitea/runner/compare/v0.3.1...v0.4.0) (2026-04-15)

### Features

* add Prometheus metrics endpoint for runner observability ([#820](https://gitea.com/gitea/runner/issues/820)) ([f33e5a6](https://gitea.com/gitea/runner/commits/f33e5a6245696af3dbd3a2c6004da35c720abfed))

### Performance Improvements

* reduce runner-to-server connection load with adaptive reporting and polling ([#819](https://gitea.com/gitea/runner/issues/819)) ([f2d5455](https://gitea.com/gitea/runner/commits/f2d545565fb93df2c06cb7aeb791d9b0a24b67e4))

### Chores

* **deps:** upgrade golangci-lint from v2.10.1 to v2.11.4 ([#821](https://gitea.com/gitea/runner/issues/821)) ([40dcee0](https://gitea.com/gitea/runner/commits/40dcee0991c3bd33b657bb77aa1f2f46d69cc0e2))

## [0.3.1](https://gitea.com/gitea/runner/compare/v0.3.0...v0.3.1) (2026-03-26)

### Features

* add configurable bind_workdir option with workspace cleanup for DinD setups ([#810](https://gitea.com/gitea/runner/issues/810)) ([9933ea0](https://gitea.com/gitea/runner/commits/9933ea0d9270c9c9e687aed7953a7a95e0c5c4bd))

### Bug Fixes

* composite action log result reported as step result ([#801](https://gitea.com/gitea/runner/issues/801)) ([f0f9f0c](https://gitea.com/gitea/runner/commits/f0f9f0c8ab4bbbbcbe5ce8843e6624caaa197417))
* race condition in reporter between RunDaemon and Close ([#796](https://gitea.com/gitea/runner/issues/796)) ([ae6e5df](https://gitea.com/gitea/runner/commits/ae6e5dfcf7da01951d742e52a5979551285fc39f)), closes [#793](https://gitea.com/gitea/runner/issues/793)

### Chores

* **deps:** update go dependencies ([#802](https://gitea.com/gitea/runner/issues/802)) ([f4418ef](https://gitea.com/gitea/runner/commits/f4418eff18d00dd77546142cda73af4b6d1ffa59))
* **deps:** update workflow action versions ([#792](https://gitea.com/gitea/runner/issues/792)) ([0f7efae](https://gitea.com/gitea/runner/commits/0f7efae806df1750083452f8211688e472c0c12d))
* **lint:** add golangci-lint v2 and fix all lint issues ([#803](https://gitea.com/gitea/runner/issues/803)) ([658101d](https://gitea.com/gitea/runner/commits/658101d9cb423b9ef5da3526e8d36cc256b3381d))

## [0.3.0](https://gitea.com/gitea/runner/compare/v0.261.8...v0.3.0) (2026-02-18)

### Bug Fixes

* **dockerfile:** Pin docker dind images to version 28, as version 29 has breaking changes in the API that are currently causing problems with the `act` fork of Gitea ([#769](https://gitea.com/gitea/runner/issues/769)) ([774f316](https://gitea.com/gitea/runner/commits/774f316c8b85b8123bc9ea388243c6f0f4167e1a)), closes [#768](https://gitea.com/gitea/runner/issues/768)

### Continuous Integration

* release binary for linux/loong64 ([#756](https://gitea.com/gitea/runner/issues/756)) ([3b11bac](https://gitea.com/gitea/runner/commits/3b11bac2ad744fd8709726fef2248d1e171d8b88))

### Chores

* **deps:** bump Go version to 1.26 in Dockerfile and Makefile ([#788](https://gitea.com/gitea/runner/issues/788)) ([9f91fdd](https://gitea.com/gitea/runner/commits/9f91fdd0eb070b19ade342e54f078a0c63685599)), references [#787](https://gitea.com/gitea/runner/issues/787)
* **deps:** bump Go version to 1.26.0 ([#787](https://gitea.com/gitea/runner/issues/787)) ([7a7a5d9](https://gitea.com/gitea/runner/commits/7a7a5d90514b51a9a58a905b741da692abaa70f9))
* **deps:** update act to v0.261.8 ([#791](https://gitea.com/gitea/runner/issues/791)) ([aab2490](https://gitea.com/gitea/runner/commits/aab249000c238cce71a38c685561e88a1e1014f1))
* **deps:** upgrade dependencies ([#775](https://gitea.com/gitea/runner/issues/775)) ([c4b57fb](https://gitea.com/gitea/runner/commits/c4b57fbcb250b543848cdfedc89ebb6be20b50f4))
* update config.example.yaml for https://gitea.com/gitea/act_runner/pulls/724 ([#763](https://gitea.com/gitea/runner/issues/763)) ([96b866a](https://gitea.com/gitea/runner/commits/96b866a3a87168255e81e855ff9bbd1a7f2aa567))

## [0.2.13](https://gitea.com/gitea/runner/compare/v0.2.12...v0.2.13) (2025-08-28)

### Features

* **docker:** add TZ env variable working ([#738](https://gitea.com/gitea/runner/issues/738)) ([aa28f8d](https://gitea.com/gitea/runner/commits/aa28f8d99c88fe7d662e49efcb8caa6215213058))
* support github mirror ([#716](https://gitea.com/gitea/runner/issues/716)) ([bbf9d7e](https://gitea.com/gitea/runner/commits/bbf9d7e90f3fee9e9968fc49c1ee6cfdb11f55e9))

## [0.2.12](https://gitea.com/gitea/runner/compare/v0.2.11...v0.2.12) (2025-06-11)

### Features

* add once flag to daemon command ([#19](https://gitea.com/gitea/runner/issues/19)) ([#598](https://gitea.com/gitea/runner/issues/598)) ([8bc0275](https://gitea.com/gitea/runner/commits/8bc0275e74af3516cb9d952e144e48405c084aa5))
* docker env vars for ephemeral and once ([#685](https://gitea.com/gitea/runner/issues/685)) ([44b4736](https://gitea.com/gitea/runner/commits/44b473670330342d66300d6ac96bccd0e39c37e7))
* register interactive with values from cli ([#682](https://gitea.com/gitea/runner/issues/682)) ([a616ed1](https://gitea.com/gitea/runner/commits/a616ed1a10efcb91f8ff3f67b35e16686cfa28b9))

### Bug Fixes

* invalid label NoInteractive exit code ([#683](https://gitea.com/gitea/runner/issues/683)) ([f0b5aff](https://gitea.com/gitea/runner/commits/f0b5aff3bbeda469a682fe21144c16485ea2189a)), closes [#665](https://gitea.com/gitea/runner/issues/665)

## [0.2.11](https://gitea.com/gitea/runner/compare/v0.2.10...v0.2.11) (2024-09-24)

### Features

* allow graceful shutdowns ([#546](https://gitea.com/gitea/runner/issues/546)) ([d1d3cad](https://gitea.com/gitea/runner/commits/d1d3cad4b0bfd0bfc5df344306f304043ff63223))

### Continuous Integration

* standardize code style and update version extraction ([#566](https://gitea.com/gitea/runner/issues/566)) ([2f78411](https://gitea.com/gitea/runner/commits/2f78411c3d871f0046832e4d082ba67c0655baf7))

## [0.2.10](https://gitea.com/gitea/runner/compare/v0.2.9...v0.2.10) (2024-04-10)

## [0.2.9](https://gitea.com/gitea/runner/compare/v0.2.8...v0.2.9) (2024-04-02)

## [0.2.8](https://gitea.com/gitea/runner/compare/v0.2.7...v0.2.8) (2024-03-27)

## [0.2.7](https://gitea.com/gitea/runner/compare/v0.2.6...v0.2.7) (2024-03-27)

### Bug Fixes

* **deps:** update module code.gitea.io/gitea-vet to v0.2.3 ([#410](https://gitea.com/gitea/runner/issues/410)) ([feb3966](https://gitea.com/gitea/runner/commits/feb39666cc7f72d59a776ac301ea66abe849403e))
* **deps:** update module github.com/avast/retry-go/v4 to v4.5.1 ([#411](https://gitea.com/gitea/runner/issues/411)) ([8c8a8ce](https://gitea.com/gitea/runner/commits/8c8a8ce40151ec5c28eeb80b24a83c9a509a3839))
* **deps:** update module github.com/docker/docker to v24.0.7+incompatible ([#413](https://gitea.com/gitea/runner/issues/413)) ([2f3e5c7](https://gitea.com/gitea/runner/commits/2f3e5c712523079ac9677158339b714234a835f2))
* **deps:** update module github.com/mattn/go-isatty to v0.0.20 ([#414](https://gitea.com/gitea/runner/issues/414)) ([825c6f9](https://gitea.com/gitea/runner/commits/825c6f97b79b33f3b103ab45c6d7a328cf2a464f))
* **deps:** update module github.com/spf13/cobra to v1.8.0 ([#416](https://gitea.com/gitea/runner/issues/416)) ([4d9de6c](https://gitea.com/gitea/runner/commits/4d9de6ca8ce74c2718ce6e3417cd3970317fa11a))
* **deps:** update module golang.org/x/time to v0.4.0 ([#424](https://gitea.com/gitea/runner/issues/424)) ([91bfe4c](https://gitea.com/gitea/runner/commits/91bfe4c186bf38daad2e3596912657bfce68849a))
* **deps:** update module golang.org/x/time to v0.5.0 ([#429](https://gitea.com/gitea/runner/issues/429)) ([08c681b](https://gitea.com/gitea/runner/commits/08c681be0c3e7e1f10eace654b8161cd5d9b5d03))
* Exit docker container on any process exit ([#462](https://gitea.com/gitea/runner/issues/462)) ([f1f9142](https://gitea.com/gitea/runner/commits/f1f9142a3c9ec062dd0ca222456b24c4a105dac9)), references [#460](https://gitea.com/gitea/runner/issues/460)

### Chores

* **deps:** update actions/setup-go action to v4 ([#418](https://gitea.com/gitea/runner/issues/418)) ([0adfc1c](https://gitea.com/gitea/runner/commits/0adfc1c7cc315c564071a95ca19b923616355dc5))
* **deps:** update docker/build-push-action action to v5 ([#419](https://gitea.com/gitea/runner/issues/419)) ([e3c6866](https://gitea.com/gitea/runner/commits/e3c68668fa3a97d88613ce2a73d73da1836d62d7))
* **deps:** update docker/login-action action to v3 ([#420](https://gitea.com/gitea/runner/issues/420)) ([f1b27d5](https://gitea.com/gitea/runner/commits/f1b27d52747c63da5643b7bdde8996208faff208))
* **deps:** update docker/setup-buildx-action action to v3 ([#421](https://gitea.com/gitea/runner/issues/421)) ([655a39f](https://gitea.com/gitea/runner/commits/655a39fd619c103c87277b4fa5c4fa7e2d23b5c5))
* **deps:** update docker/setup-qemu-action action to v3 ([#422](https://gitea.com/gitea/runner/issues/422)) ([cca7d54](https://gitea.com/gitea/runner/commits/cca7d54117a78bb84c2efc78e0cf8f63876962c6))

## [0.2.6](https://gitea.com/gitea/runner/compare/v0.2.51...v0.2.6) (2023-09-24)

## [0.2.5](https://gitea.com/gitea/runner/compare/v0.2.4...v0.2.5) (2023-08-02)

## [0.2.4](https://gitea.com/gitea/runner/compare/v0.2.3...v0.2.4) (2023-07-25)

### Bug Fixes

* **register:** refactor context usage in registration functions ([#286](https://gitea.com/gitea/runner/issues/286)) ([8705871](https://gitea.com/gitea/runner/commits/87058716fb95ffb715472ee7fab2e8fb8a6e75b0))

### Refactoring

* **register:** refactor registration functions and improve logging ([#288](https://gitea.com/gitea/runner/issues/288)) ([cf92a97](https://gitea.com/gitea/runner/commits/cf92a979e274dace1d545cbe140f4f4db064f921))

### Continuous Integration

* **lint:** refactor code for clarity and linting compliance ([#289](https://gitea.com/gitea/runner/issues/289)) ([db662b3](https://gitea.com/gitea/runner/commits/db662b36906a3aa0bd1211cee64e4ce73288f485))

## [0.2.3](https://gitea.com/gitea/runner/compare/v0.2.2...v0.2.3) (2023-07-03)

## [0.2.2](https://gitea.com/gitea/runner/compare/v0.2.1...v0.2.2) (2023-07-01)

## [0.2.1](https://gitea.com/gitea/runner/compare/v0.2.0...v0.2.1) (2023-06-30)

### Features

* improve Docker configuration and detection handling ([#242](https://gitea.com/gitea/runner/issues/242)) ([9e4a5f7](https://gitea.com/gitea/runner/commits/9e4a5f7363a5559ddaa3a48e91bb8e69db861de4))

## [0.2.0](https://gitea.com/gitea/runner/compare/v0.1.8...v0.2.0) (2023-06-15)

### Refactoring

* daemon, config, and logging for better clarity ([#225](https://gitea.com/gitea/runner/issues/225)) ([69c55ee](https://gitea.com/gitea/runner/commits/69c55ee003a3553d4e8b50d9bf98d2c31c8832e4))

### Documentation

* improve examples README and organization ([#230](https://gitea.com/gitea/runner/issues/230)) ([a83f29d](https://gitea.com/gitea/runner/commits/a83f29d5a96a78b61a9dc6967d09a31ef0ad5d7e))

## [0.1.8](https://gitea.com/gitea/runner/compare/v0.1.7...v0.1.8) (2023-05-09)

## [0.1.7](https://gitea.com/gitea/runner/compare/v0.1.6...v0.1.7) (2023-05-06)

### Continuous Integration

* improve release process and test robustness ([#173](https://gitea.com/gitea/runner/issues/173)) ([eef3c32](https://gitea.com/gitea/runner/commits/eef3c32eb2418dcbcd40cb13076c2887d4f6913b))

### Chores

* improve Dockerfile, README, and testing settings ([#172](https://gitea.com/gitea/runner/issues/172)) ([c40b651](https://gitea.com/gitea/runner/commits/c40b6518736636db6edf7bbecf4d650f3a2117c4))

## [0.1.6](https://gitea.com/gitea/runner/compare/v0.1.5...v0.1.6) (2023-05-01)

### Build

* improve compression and update GitHub actions ([#168](https://gitea.com/gitea/runner/issues/168)) ([b498341](https://gitea.com/gitea/runner/commits/b49834185777fffed9e0b4d7c3b08ad82fcdbcfa))
* optimize Dockerfile and update dependencies ([#162](https://gitea.com/gitea/runner/issues/162)) ([0d727eb](https://gitea.com/gitea/runner/commits/0d727eb2621c706d973987490bb3b14163693cf3))

## [0.1.5](https://gitea.com/gitea/runner/compare/v0.1.4...v0.1.5) (2023-04-28)

### Continuous Integration

* improve API usage and test robustness across platforms ([#159](https://gitea.com/gitea/runner/issues/159)) ([49d2cb0](https://gitea.com/gitea/runner/commits/49d2cb0cb53ebe6f67284b2fd27e11fe81884d85))

## [0.1.4](https://gitea.com/gitea/runner/compare/v0.1.3...v0.1.4) (2023-04-27)

## [0.1.3](https://gitea.com/gitea/runner/compare/v0.1.2...v0.1.3) (2023-04-27)

### Continuous Integration

* **actions:** add build docker image workflow ([#118](https://gitea.com/gitea/runner/issues/118)) ([66bab3d](https://gitea.com/gitea/runner/commits/66bab3d8059fde868b59d1aeb9e0bdecb6c0d55e))

## [0.1.2](https://gitea.com/gitea/runner/compare/v0.1.1...v0.1.2) (2023-04-13)

## [0.1.1](https://gitea.com/gitea/runner/compare/v0.1.0...v0.1.1) (2023-04-11)

### Features

* add artifact api from gitea server ([#103](https://gitea.com/gitea/runner/issues/103)) ([b79c3aa](https://gitea.com/gitea/runner/commits/b79c3aa1a325ab49902fd662174a5ea0f8eed196))

## [0.1.0](https://gitea.com/gitea/runner/compare/v0.0.1...v0.1.0) (2023-04-04)

## [0.0.1](https://gitea.com/gitea/runner/compare/e54de3724292fa3fc4a50f1e1186964ce1e7a816...v0.0.1) (2023-02-26)

### Features

* add task log output hook ([dfcb104](https://gitea.com/gitea/runner/commits/dfcb104af02b536df7dfd7dc6d5911b654038ae9))
* auto remove containers ([8f950cc](https://gitea.com/gitea/runner/commits/8f950ccec0a9aae4a6e66d10dc28edd5fdf04ccb))
* avoid reuse of containers ([cbd21b8](https://gitea.com/gitea/runner/commits/cbd21b8ec539a5def5649b10702e55f6656e2ee9))
* **cli:** add register command basic functions ([561bfad](https://gitea.com/gitea/runner/commits/561bfad7c5413e29f61f44fc4526ad6123705c4c))
* **cli:** make register interactive stages working ([bfcc6f5](https://gitea.com/gitea/runner/commits/bfcc6f56f320d891ba61e886acf3b49ac6d766e2))
* **cli:** set interactive stages in register command ([8f9d750](https://gitea.com/gitea/runner/commits/8f9d7506dc71aaa55db74b80a7118bc18d6947f2))
* don't require docker ([#16](https://gitea.com/gitea/runner/issues/16)) ([990cf93](https://gitea.com/gitea/runner/commits/990cf93c7136669408eb1832cd05df3ad4dd81b3))
* fill default task runner ([3281c66](https://gitea.com/gitea/runner/commits/3281c6637d5b49f761e80584235d67384541a04a))
* move main task logic to runtime package ([c5edbbd](https://gitea.com/gitea/runner/commits/c5edbbd004418890b593dd9217160b02f5f2143d))
* parse register arguments and do register when non-interactive mode ([45de619](https://gitea.com/gitea/runner/commits/45de6199d12626de6c934444624c24e8f887e8f6))
* **poller:** support scheduler to fetch task and dispatch to worker ([a13ed0c](https://gitea.com/gitea/runner/commits/a13ed0c52e2879ef5fa2ec8632da2efd535c56d0))
* **register:** add label format validation ([3f676a7](https://gitea.com/gitea/runner/commits/3f676a742b72f550b7900a1c0d4faab19b11750f))
* **register:** add labels to local runner file ([86f66ea](https://gitea.com/gitea/runner/commits/86f66eaa4508603d9359398e8dcf866d8689bb89))
* replace ciphertext in log ([d650bae](https://gitea.com/gitea/runner/commits/d650bae0b6212b5a38aef62a6fab13c04b971f45))
* support post step log ([ab4e06f](https://gitea.com/gitea/runner/commits/ab4e06f977ca78181077e5de43d7aad070204192))
* task can report step and final result ([6030610](https://gitea.com/gitea/runner/commits/6030610c04cc75219d6047ce51e57b45ea3c1fb9))
* use specified labels ([2354f5b](https://gitea.com/gitea/runner/commits/2354f5bb18902124c0ccc07c632888bac8b4c99b))

### Bug Fixes

* call RunDaemon ([93e42d1](https://gitea.com/gitea/runner/commits/93e42d14650c8e1c55fbf37287c6f6e8ef2f24c5))
* ignore stuck because of wrong token ([88ae188](https://gitea.com/gitea/runner/commits/88ae188699428e10d0d5d886f00f07c100f4cd7b))
* limit workers ([5781e23](https://gitea.com/gitea/runner/commits/5781e233c14c6ff2c953be0727b1a9ed7977f4dc))
* need more time to run task ([cce60d4](https://gitea.com/gitea/runner/commits/cce60d466bb189ee83c92659aefb6354e34b9c56))
* **poller:** graceful shutdown ([9830f34](https://gitea.com/gitea/runner/commits/9830f34d364aa289e98b03f6d501a5afbb8379eb))
* **poller:** initial ready channel ([cab5699](https://gitea.com/gitea/runner/commits/cab56996b7b2f84723e57b0c2f0d80751570fc80))
* report job failure on error return ([#14](https://gitea.com/gitea/runner/issues/14)) ([f7a5278](https://gitea.com/gitea/runner/commits/f7a52789d950a57cf0c985d02715714740426e5c))
* **runner:** check task state field exist. ([2442cdd](https://gitea.com/gitea/runner/commits/2442cdd8adbd002ebadd982fd303879cbc73454b))
* set more image ([64c9856](https://gitea.com/gitea/runner/commits/64c9856ba7b05fe4e9cd77173a9b96332f7e0fec))
* set NoSkipCheckout ([2e1326b](https://gitea.com/gitea/runner/commits/2e1326b725bda03923e4f9f6d6b24a025efc5e3c))
* **typo:** Use uppercase N for default input prompt ([#9](https://gitea.com/gitea/runner/issues/9)) ([9ba9604](https://gitea.com/gitea/runner/commits/9ba9604256964f52ff102b16acf782d8c56d87ff))
* update container config ([b3e7d6a](https://gitea.com/gitea/runner/commits/b3e7d6a62d6c00850616fddc199bef6bb65e8a85))
* update step result ([41b5fa6](https://gitea.com/gitea/runner/commits/41b5fa6b557d5cb937690635306aea6749ee6c91))

### Refactoring

* Add graceful shutdown signal notify func ([e54de37](https://gitea.com/gitea/runner/commits/e54de3724292fa3fc4a50f1e1186964ce1e7a816))
* docker engine with options ([14a3455](https://gitea.com/gitea/runner/commits/14a345504f03c02a70b9323d157deccab236d86a))
* **task:** execute single task with gRPC data ([20c3d85](https://gitea.com/gitea/runner/commits/20c3d85ba9caef27fbf83118ce4b6126fed97580))

### Documentation

* Add Prerequisites ([#11](https://gitea.com/gitea/runner/issues/11)) ([7abe039](https://gitea.com/gitea/runner/commits/7abe039c311dd257f7ac4f0c569afe701c1ae6d0)), closes [#10](https://gitea.com/gitea/runner/issues/10)
* add quickstart to readme ([03c4c4e](https://gitea.com/gitea/runner/commits/03c4c4e1164e18306652ed5d18f7e9dcfb3b139a))
* **readme:** update some header. ([#1](https://gitea.com/gitea/runner/issues/1)) ([474683c](https://gitea.com/gitea/runner/commits/474683c0e8bca9bef135fc49d44e650144591dfb))

### Chores

* add .golangci config ([518aa1b](https://gitea.com/gitea/runner/commits/518aa1b2a76822e413293701af87be5f94b37c82))
* Add makefile ([a3e9bbe](https://gitea.com/gitea/runner/commits/a3e9bbed25680fcf46bfad48dd04e50565412b47))
* Add time sleep 1 second. ([bca586f](https://gitea.com/gitea/runner/commits/bca586ffd0bb9a726aed8b3d86e7c76b6ca0c3ef))
* **client:** support gRPC and gRPC web protocol ([d359276](https://gitea.com/gitea/runner/commits/d359276fe1d4203f84f21249f40bb942a29fc7a5))
* **config:** remove zerolog ([9225a3a](https://gitea.com/gitea/runner/commits/9225a3a856e85f63e1160bff0acbd3cb21f3d9d4))
* **engine:** ping docker daemon ([4eb5213](https://gitea.com/gitea/runner/commits/4eb5213294e4be894712f5a081cc18f01db308ac))
* fix lint error ([827bd95](https://gitea.com/gitea/runner/commits/827bd953e7d9a0dbe7caab149b898a926546a044))
* **go:** rename module name ([aa765ea](https://gitea.com/gitea/runner/commits/aa765ea4d0fe51a1211b755f31f96a1d5748aa8a))
* **gRPC:** add request interface client ([f2fb879](https://gitea.com/gitea/runner/commits/f2fb8798fa1556ad80d78a9649588a6089dfedcd))
* **gRPC:** add update and update step for client ([9be39b8](https://gitea.com/gitea/runner/commits/9be39b8cd4629678dac490eb5f0710cdd12ae00e))
* **gRPC:** fetch new stage data. ([631f801](https://gitea.com/gitea/runner/commits/631f801aa658ed0fa07df5025b067e20f4db4749))
* **gRPC:** register new runner ([449388f](https://gitea.com/gitea/runner/commits/449388f3ab1dae8944dd955288c77373036a96ac))
* load config from env file ([aa9ed05](https://gitea.com/gitea/runner/commits/aa9ed0590359c29181811b9f503921e414e465b2))
* **piepline:** add runtime package. ([7d55fd5](https://gitea.com/gitea/runner/commits/7d55fd57c94d1862f91ff7e15b892ab5f4b0af88))
* **poller:** add error retry limit ([45b0429](https://gitea.com/gitea/runner/commits/45b0429b2151c1dddaa722206bbd518f315b2006))
* **poller:** add metric to track the worker number ([abdb547](https://gitea.com/gitea/runner/commits/abdb547b1b5f5f3e5addc3c9c07b8eb0d9d1fc22))
* **poller:** Add poller package ([0b885c5](https://gitea.com/gitea/runner/commits/0b885c5e5fe33e2781755a8a02fa4bddfa997306))
* **proto:** Add ping request. ([4d7ef95](https://gitea.com/gitea/runner/commits/4d7ef95d402d3e53d175e262e01af224f4c58da7))
* **proto:** replace import path. ([307c2c6](https://gitea.com/gitea/runner/commits/307c2c63de7e486c785f191ec2e66e7a9b49df80))
* **register:** update proto file. ([c161a48](https://gitea.com/gitea/runner/commits/c161a48a0a5cdea220acc96b9e5df6c0b9c27e54))
* remove tracing codes ([2f879c4](https://gitea.com/gitea/runner/commits/2f879c41c4a9e44dded5e947b08f45b98976c50d))
* remove unused logrous package ([5213612](https://gitea.com/gitea/runner/commits/52136120333c5738feae5fc5f43c9577fdcbac3c))
* replace with code.gitea.io/bots-proto-go ([44988db](https://gitea.com/gitea/runner/commits/44988db9f89d2c55128688ec69309050aa140dcc))
* **runner:** add default label ""self-hosted" ([dada073](https://gitea.com/gitea/runner/commits/dada0730b04723483552337a2fbf1c40489571bb))
* **runner:** add new register command ([a503f74](https://gitea.com/gitea/runner/commits/a503f7429f82e1f55c8c13293469c85db5dfa39b))
* **runner:** add new token in header ([d178051](https://gitea.com/gitea/runner/commits/d178051832efa7b617b117c0d4c4dbe543d79c06))
* **runner:** add runner uuid in http header. ([bf5e3dc](https://gitea.com/gitea/runner/commits/bf5e3dc302e195dbe0b65688f5a642c372af0411))
* **runner:** cancel task if get the cancel from server ([08c94bb](https://gitea.com/gitea/runner/commits/08c94bb564709d87fe62a094d07f30578b16f594))
* **runner:** refactor register flow ([3a15031](https://gitea.com/gitea/runner/commits/3a1503138b7360ff143d7e48b09effc5ef45f3ab))
* **runner:** remove client secret and add UUID in runner ([e08495a](https://gitea.com/gitea/runner/commits/e08495af090634d69704f2d9667bdddce7c47684))
* **runner:** remove update runner status ([e9e42e8](https://gitea.com/gitea/runner/commits/e9e42e850b526fa87e5edde09b3cb60d57e9fbf2))
* **runner:** support update log and task ([82431d8](https://gitea.com/gitea/runner/commits/82431d8e114b75c9d4a02f009c74503f1dca0789))
* **runner:** update fetch job request ([7486d2a](https://gitea.com/gitea/runner/commits/7486d2ab918aacdb6c0bb3b85a3c7aa6230d1afb))
* **runner:** update runner status if start program ([d4c1515](https://gitea.com/gitea/runner/commits/d4c1515f4ef6ebe732d8bee9cd6f3e019d1560fb))
* **runner:** update runner status when start job ([7bebd2b](https://gitea.com/gitea/runner/commits/7bebd2bbadd10171b355306a25294cb9576ffb9b))
* **runner:** wait workflow done before shutdown the service ([d1114da](https://gitea.com/gitea/runner/commits/d1114da299be65d69a4d1bb7574851c1f6cd5d1a))
* **runtime:** check error message ([5051e4a](https://gitea.com/gitea/runner/commits/5051e4aebdbb51575fa77bf472f4ee4ce03dfd0c))
* **runtime:** fetch build data ([d3d56ed](https://gitea.com/gitea/runner/commits/d3d56ed0ef86b83745555adb849852ccd3caacb8))
* update grpc router path ([f05e08a](https://gitea.com/gitea/runner/commits/f05e08a76750d1f8459b75a3b8f1461f54e4aa1c))
* upgrade act ([66ec50e](https://gitea.com/gitea/runner/commits/66ec50e427cb8e1b7b217803636e80704ae27b06))

## [0.261.8](https://gitea.com/gitea/runner/compare/v0.2.13...v0.261.8) (2026-02-18)

### Features

* `--container-options` ([#1462](https://gitea.com/gitea/runner/issues/1462)) ([4c2524a](https://gitea.com/gitea/runner/commits/4c2524ab4d0d923f241f2649702e5dd8927846a4))
* add `go` install instruction ([#672](https://gitea.com/gitea/runner/issues/672)) ([37a98f1](https://gitea.com/gitea/runner/commits/37a98f134a7e69bd45b0da40ea175d4cb5b41211))
* add bug-report flag ([#1056](https://gitea.com/gitea/runner/issues/1056)) ([5d7027d](https://gitea.com/gitea/runner/commits/5d7027dc3fc6a16e2023ffa1ebac93f1c667786e))
* add check for newer versions ([#1562](https://gitea.com/gitea/runner/issues/1562)) ([9390793](https://gitea.com/gitea/runner/commits/93907931df737bdd21ecc624f4983fb1b365357a))
* Add CopyTarStream to Container Interface ([#1912](https://gitea.com/gitea/runner/issues/1912)) ([6468dd7](https://gitea.com/gitea/runner/commits/6468dd7fc8c2013dc76e8ea6e9ab590dc6dbb56e))
* add Dependabot config ([#797](https://gitea.com/gitea/runner/issues/797)) ([12928a0](https://gitea.com/gitea/runner/commits/12928a0ac61b2b3278aead1e42923c876a8cf5a4))
* add flag to always run the checkout action ([#1049](https://gitea.com/gitea/runner/issues/1049)) ([f09cdb8](https://gitea.com/gitea/runner/commits/f09cdb844367189aafe2f3af5795d3136db6b290))
* add fromJSON support ([#352](https://gitea.com/gitea/runner/issues/352)) ([d521fa5](https://gitea.com/gitea/runner/commits/d521fa5bbae1d6fd5b775e0924907e882a4d8655))
* add json logger output ([#1026](https://gitea.com/gitea/runner/issues/1026)) ([14c9801](https://gitea.com/gitea/runner/commits/14c98015c57bf577619101ac4770f90c6435470b))
* add logger hook for standalone job logger ([0fa1fe0](https://gitea.com/gitea/runner/commits/0fa1fe03109cb5e41aa7a074fe0d837089d65c1e))
* Add new Action Cache ([#1913](https://gitea.com/gitea/runner/issues/1913)) ([a42f3cf](https://gitea.com/gitea/runner/commits/a42f3cf1cdaa1d93da181bd5f556371d61198169))
* add node16 support ([#922](https://gitea.com/gitea/runner/issues/922)) ([5580812](https://gitea.com/gitea/runner/commits/558081242ca322997437de126324ae9977eb7da0))
* add option for custom socket path ([#698](https://gitea.com/gitea/runner/issues/698)) ([799098b](https://gitea.com/gitea/runner/commits/799098b0e630caff38f220dd0ab5e216667ec916))
* add option for docker image rebuild ([#878](https://gitea.com/gitea/runner/issues/878)) ([8ad6c07](https://gitea.com/gitea/runner/commits/8ad6c07083c7cd7f45ea9e0808e3297bf2370ab5)), references [#714](https://gitea.com/gitea/runner/issues/714) [#716](https://gitea.com/gitea/runner/issues/716)
* add option to bypass GHE for actions checkout ([#1162](https://gitea.com/gitea/runner/issues/1162)) ([de37f75](https://gitea.com/gitea/runner/commits/de37f750773575651a91363eb6b1ce97be964845)), references [#1161](https://gitea.com/gitea/runner/issues/1161)
* add option to specify git remote name ([#1104](https://gitea.com/gitea/runner/issues/1104)) ([d4272bd](https://gitea.com/gitea/runner/commits/d4272bd9fe13566ab76e5ad091aedbded74ff4e2))
* add remote reusable workflows ([#1525](https://gitea.com/gitea/runner/issues/1525)) ([82a8c1e](https://gitea.com/gitea/runner/commits/82a8c1e80dfcd73842bc0359ca8447303b6fcc08)), closes [#826](https://gitea.com/gitea/runner/issues/826)
* add skipped status as step result ([#950](https://gitea.com/gitea/runner/issues/950)) ([2eda7c6](https://gitea.com/gitea/runner/commits/2eda7c6037b921ae88742b2ff29b20bb171497a6))
* add step name to logger field ([#1027](https://gitea.com/gitea/runner/issues/1027)) ([6b05957](https://gitea.com/gitea/runner/commits/6b059572ff51e30f4cfbde76cc7b82af313967b6))
* add support for building docker actions with private registries ([#1557](https://gitea.com/gitea/runner/issues/1557)) ([3f3b25a](https://gitea.com/gitea/runner/commits/3f3b25ae84a25f8488f3801e3342bb3f96c6cd8a))
* Add support for github action variables ([#1833](https://gitea.com/gitea/runner/issues/1833)) ([6ce45e3](https://gitea.com/gitea/runner/commits/6ce45e3f246f12d6617691de9a2423920d5fdbe6))
* add ubuntu-22.04 ([#1150](https://gitea.com/gitea/runner/issues/1150)) ([69691c8](https://gitea.com/gitea/runner/commits/69691c88b301ee695f5712d2ce8ad2b760725c80))
* adjust container name format ([#1](https://gitea.com/gitea/runner/issues/1)) ([4cacc14](https://gitea.com/gitea/runner/commits/4cacc14d22447358f1d34b4ef556fe1fd9cb9dd0))
* Allow building without docker support ([#1507](https://gitea.com/gitea/runner/issues/1507)) ([7e8d070](https://gitea.com/gitea/runner/commits/7e8d0708111125b11166494dc22da6d67f9e83e7))
* allow existing logger from context ([#859](https://gitea.com/gitea/runner/issues/859)) ([603cf56](https://gitea.com/gitea/runner/commits/603cf568780752efbe587dd50825b2ee5c3013ae))
* allow existing logger from context ([#898](https://gitea.com/gitea/runner/issues/898)) ([6517d04](https://gitea.com/gitea/runner/commits/6517d04b9a4a677fc03277b6cf32fcd56a4dcf36))
* allow overriding of `GITHUB_` env variables ([#1582](https://gitea.com/gitea/runner/issues/1582)) ([ce168f9](https://gitea.com/gitea/runner/commits/ce168f959502feab3ec815810fac7a3b5a1f470c))
* allow to spawn and run a local reusable workflow ([#1423](https://gitea.com/gitea/runner/issues/1423)) ([a8e05cd](https://gitea.com/gitea/runner/commits/a8e05cded6be8bfb77536d823933725144918bcb))
* cache notices to reduce frequency of upgrade notifications ([#1592](https://gitea.com/gitea/runner/issues/1592)) ([42b9b73](https://gitea.com/gitea/runner/commits/42b9b73d388544693d44d419457151c10e66472a))
* cli option to enable the new action cache ([#1954](https://gitea.com/gitea/runner/issues/1954)) ([f7a846d](https://gitea.com/gitea/runner/commits/f7a846d2f53ad47bdd98483ae25829497f8a1f0a))
* cmd support for windows ([#1941](https://gitea.com/gitea/runner/issues/1941)) ([8314095](https://gitea.com/gitea/runner/commits/83140951bf96c500f0c1a745e4333b27f2820e49))
* correctly use the xdg library, which has the side effect to fix the config survey ([#2195](https://gitea.com/gitea/runner/issues/2195)) ([12c0c42](https://gitea.com/gitea/runner/commits/12c0c4277af1e3d07a6389c3260d28579ba31e01))
* docker 19.03.12 ([#347](https://gitea.com/gitea/runner/issues/347)) ([b3b8961](https://gitea.com/gitea/runner/commits/b3b8961122d06c17bb9bbc063bf56e12aba60106))
* Don't push `pkg/runner/act/` where actions cache is stored ([#620](https://gitea.com/gitea/runner/issues/620)) ([b04d762](https://gitea.com/gitea/runner/commits/b04d762614433695202ee1c12d16654b558165bc))
* GITHUB_STATE and GITHUB_OUTPUT file commands ([#1391](https://gitea.com/gitea/runner/issues/1391)) ([e520382](https://gitea.com/gitea/runner/commits/e520382d2fdc5872ca8e939c6c0f2e1e74749182))
* **goreleaser:** format changelog, add prereleases ([#942](https://gitea.com/gitea/runner/issues/942)) ([4e6cddf](https://gitea.com/gitea/runner/commits/4e6cddf80a712469695640058bcfceee7adaeda6))
* handle context cancelation during docker exec ([#1170](https://gitea.com/gitea/runner/issues/1170)) ([4ef50ee](https://gitea.com/gitea/runner/commits/4ef50eeae76f76af3ab0c8c7377c9388d28a0c38))
* Host environment ([#1293](https://gitea.com/gitea/runner/issues/1293)) ([f2b98ed](https://gitea.com/gitea/runner/commits/f2b98ed301290f53e96c64814ab9949640e10307))
* implement steps.timeout-minutes ([#1776](https://gitea.com/gitea/runner/issues/1776)) ([8913375](https://gitea.com/gitea/runner/commits/8913375af89792be37603c3c6429372964bf86b1))
* improve GetOutboundIP ([#1707](https://gitea.com/gitea/runner/issues/1707)) ([8518d70](https://gitea.com/gitea/runner/commits/8518d70bdf3ee04f08e5d096fe07c2a85e8f21ea))
* improve list ([#786](https://gitea.com/gitea/runner/issues/786)) ([37aaec8](https://gitea.com/gitea/runner/commits/37aaec81f444bd5e00de4e8cedb35c73e131b41d))
* interpolate the step names ([#1422](https://gitea.com/gitea/runner/issues/1422)) ([809da71](https://gitea.com/gitea/runner/commits/809da7198cc2bed5b758c77274d47cb5360ee848)), closes [#1353](https://gitea.com/gitea/runner/issues/1353)
* JobLoggerFactory ([#1496](https://gitea.com/gitea/runner/issues/1496)) ([7d99512](https://gitea.com/gitea/runner/commits/7d9951200f366ee7a207844e185bf5de7e3f7ab5))
* load every environment from --env-file to workflow ([#184](https://gitea.com/gitea/runner/issues/184)) ([2f39547](https://gitea.com/gitea/runner/commits/2f395475b0ac4f5c0fb3d7c8b97a23d9214ea77b))
* non root user container ([#1202](https://gitea.com/gitea/runner/issues/1202)) ([b59e6de](https://gitea.com/gitea/runner/commits/b59e6dee6d01253cb96505f8effc0516c2a390d9))
* offline mode for new action cache ([#2173](https://gitea.com/gitea/runner/issues/2173)) ([852959e](https://gitea.com/gitea/runner/commits/852959e1e14bf5cd9c0284d15ca899bef68ff4af))
* parse types of reusable workflows ([#1414](https://gitea.com/gitea/runner/issues/1414)) ([9b95a72](https://gitea.com/gitea/runner/commits/9b95a728d84920a75d22c0dc6088391fbbc4002a)), references [#826](https://gitea.com/gitea/runner/issues/826)
* pass current matrix information to job logger ([#1364](https://gitea.com/gitea/runner/issues/1364)) ([00acf68](https://gitea.com/gitea/runner/commits/00acf681880b0fe7cea27a29d61c261d624966c3))
* release extension ([78e1ceb](https://gitea.com/gitea/runner/commits/78e1ceb0a8ad911b04f821ca707c3b352fe2c74b))
* release extension ([d690a5f](https://gitea.com/gitea/runner/commits/d690a5fee95491d519085990d5c700ff2c9081ce))
* release extension ([19e6929](https://gitea.com/gitea/runner/commits/19e692939861bf3770e4c0ea7a7d8d5f27d84561))
* run jobs in parallel ([#1003](https://gitea.com/gitea/runner/issues/1003)) ([18792f9](https://gitea.com/gitea/runner/commits/18792f9620fd2fd7dacc53ae30a234fbfc7c81b1)), references [#CPU](https://gitea.com/gitea/runner/issues/CPU)
* SELinux lables for --bind ([#906](https://gitea.com/gitea/runner/issues/906)) ([f726339](https://gitea.com/gitea/runner/commits/f7263399b9146a43b988e011d986d18d180414d2))
* specify matrix on command line ([#1675](https://gitea.com/gitea/runner/issues/1675)) ([636c8a3](https://gitea.com/gitea/runner/commits/636c8a34aedf9b6251df6729fe673d828f972030))
* Speeding up workflows ([#643](https://gitea.com/gitea/runner/issues/643)) ([3e768cd](https://gitea.com/gitea/runner/commits/3e768cd91695af12bdf860ba49a0f992854b9c87))
* split job steps into its own files/structs ([#1004](https://gitea.com/gitea/runner/issues/1004)) ([2bb3e74](https://gitea.com/gitea/runner/commits/2bb3e74616f86ece6a11cd854ad54701542410fd))
* step summary of test results ([#1580](https://gitea.com/gitea/runner/issues/1580)) ([932863b](https://gitea.com/gitea/runner/commits/932863bef508868168954ef69572750cf87f463f))
* Support "result" on "needs" context. ([#1497](https://gitea.com/gitea/runner/issues/1497)) ([bef9b5c](https://gitea.com/gitea/runner/commits/bef9b5c3c77bcdbc6c15169c0d55f635e27350ae)), closes [#1367](https://gitea.com/gitea/runner/issues/1367) [#1466](https://gitea.com/gitea/runner/issues/1466), references [#1423](https://gitea.com/gitea/runner/issues/1423) [#1506](https://gitea.com/gitea/runner/issues/1506)
* support `GITHUB_REF_NAME` & `GITHUB_REF_TYPE` ([#1142](https://gitea.com/gitea/runner/issues/1142)) ([8e216c6](https://gitea.com/gitea/runner/commits/8e216c642e11dfa8dba7aab040bb09274c48b5a6))
* support branches as GITHUB_REF ([03d57a8](https://gitea.com/gitea/runner/commits/03d57a80c9efe130fc0dde98ab7895ac312945c5))
* support config env expansion ([#2063](https://gitea.com/gitea/runner/issues/2063)) ([0806c8b](https://gitea.com/gitea/runner/commits/0806c8b10946de509889eb723aaba45cc100c8de))
* support custom GITHUB_RUN_ID, GITHUB_RUN_NUMBER ([#369](https://gitea.com/gitea/runner/issues/369)) ([1070c58](https://gitea.com/gitea/runner/commits/1070c585381980abe8f636af5212457711a51f03))
* support interpolation in `<job>.container.options` ([#1958](https://gitea.com/gitea/runner/issues/1958)) ([19764bc](https://gitea.com/gitea/runner/commits/19764bcb06ac9acdd7100d95c72055b697820705))
* support more options of containers ([b9c20dc](https://gitea.com/gitea/runner/commits/b9c20dcaa43899cb3bb327619d447248303170e0))
* support node20 runtime ([#1988](https://gitea.com/gitea/runner/issues/1988)) ([935e37c](https://gitea.com/gitea/runner/commits/935e37c25b6c48b50882024b3c4235d12f23f469))
* support offline mode ([#2128](https://gitea.com/gitea/runner/issues/2128)) ([f055d4a](https://gitea.com/gitea/runner/commits/f055d4ae607a1d68695bbee1d9b444d3056fb387))
* support PlatformPicker ([c6b8548](https://gitea.com/gitea/runner/commits/c6b8548d35ab4a5fc0f4478ef9f001ee4fbda16f))
* support runs-on labels and group ([#2062](https://gitea.com/gitea/runner/issues/2062)) ([04011b6](https://gitea.com/gitea/runner/commits/04011b6b78c70d36db18419861c737cd0d72d146))
* support yaml env/secrets/inputs file ([#1733](https://gitea.com/gitea/runner/issues/1733)) ([50dcc57](https://gitea.com/gitea/runner/commits/50dcc57e4babe816692d342307183801ddd78356))
* Test writing remote Actions locally ([#2155](https://gitea.com/gitea/runner/issues/2155)) ([ac6a465](https://gitea.com/gitea/runner/commits/ac6a465e27dccae89e3fb98de355ae79f00f2d6e))
* try to read ref and sha from event payload if available ([#889](https://gitea.com/gitea/runner/issues/889)) ([edd0fb9](https://gitea.com/gitea/runner/commits/edd0fb92a78f2fb20ea1afe14a047fdc026d2643))
* uses in composite ([#793](https://gitea.com/gitea/runner/issues/793)) ([9868e13](https://gitea.com/gitea/runner/commits/9868e1377274a7657d4aac6bd2b66c0d564405cb))
* workflowpattern package ([#1618](https://gitea.com/gitea/runner/issues/1618)) ([05eaeaa](https://gitea.com/gitea/runner/commits/05eaeaa528f791cdb796e3fae4cabd833dcb7a40))
* **workflow:** support schedule event ([#4](https://gitea.com/gitea/runner/issues/4)) ([88cce47](https://gitea.com/gitea/runner/commits/88cce47022d34d5cbc49771b3753541fe905f5b6))

### Bug Fixes

* --version, --help do not need check platform config ([#508](https://gitea.com/gitea/runner/issues/508)) ([22d56c3](https://gitea.com/gitea/runner/commits/22d56c3517483c63dc98dc4d280d60eb5a625a3d))
* ::stop-commands:: should continue to print the lines ([#920](https://gitea.com/gitea/runner/issues/920)) ([fec0c0c](https://gitea.com/gitea/runner/commits/fec0c0c529f5510464c858ff9a6664e0647245a8))
* ([#2075](https://gitea.com/gitea/runner/issues/2075)) ([74b0fe8](https://gitea.com/gitea/runner/commits/74b0fe8ba9c04cbf3d55f6536249d959599ed676))
* [#170](https://gitea.com/gitea/runner/issues/170) [#171](https://gitea.com/gitea/runner/issues/171) - add RUNNER_OS and RUNNER_TEMP ([a9533b0](https://gitea.com/gitea/runner/commits/a9533b05ce708a84e931d63cb2cd7b795e70bbe2))
* [#170](https://gitea.com/gitea/runner/issues/170) add RUNNER_TEMP ([a48cce3](https://gitea.com/gitea/runner/commits/a48cce3a78d915d7c9e6bd28248926578caaf1c4))
* [#182](https://gitea.com/gitea/runner/issues/182) - support symlinks in workspace ([8c29099](https://gitea.com/gitea/runner/commits/8c290994c1b84153d2b9f911ac3bb64858fb5e35))
* 168 include error message from docker SDK ([6c4e7b9](https://gitea.com/gitea/runner/commits/6c4e7b9fdedeeeb1eda78cefaff55750033d1327))
* add `ImageOS` env var based on running platform ([#571](https://gitea.com/gitea/runner/issues/571)) ([41b03b5](https://gitea.com/gitea/runner/commits/41b03b581c7aec251d96e9a5ea27009542a798d4))
* add `server_url` attribute to github context ([#1727](https://gitea.com/gitea/runner/issues/1727)) ([c0130ed](https://gitea.com/gitea/runner/commits/c0130ed030e0f267b7b7b5b3ad5e1dfcc08980d8)), closes [#1726](https://gitea.com/gitea/runner/issues/1726)
* add default values in vmInput ([cc369e2](https://gitea.com/gitea/runner/commits/cc369e2f73e84596a936d602a7e112d849e6e6fa))
* add GITHUB_STEP_SUMMARY ([#1607](https://gitea.com/gitea/runner/issues/1607)) ([8790c9b](https://gitea.com/gitea/runner/commits/8790c9b8e6b67f1552b824ba56d5a851f7b394e5))
* add parent step id in composite action step id ([#1268](https://gitea.com/gitea/runner/issues/1268)) ([5f5b895](https://gitea.com/gitea/runner/commits/5f5b8959d6b12b034dfa64e91f8b17c97808eda5))
* add simple concurrency limit ([#823](https://gitea.com/gitea/runner/issues/823)) ([4d552e6](https://gitea.com/gitea/runner/commits/4d552e65cedafda89c11a36abe3a218f5d090e4b))
* add-matcher fails github workflow ([#1532](https://gitea.com/gitea/runner/issues/1532)) ([8b4f210](https://gitea.com/gitea/runner/commits/8b4f21087209e027ca0e2ec225e0c74d03552f86))
* added tests and improved code ([8d181a6](https://gitea.com/gitea/runner/commits/8d181a6a1ea1eca1076ae9d47c3fded31a6099d2))
* align github.ref to GitHub Action ([#1362](https://gitea.com/gitea/runner/issues/1362)) ([97c083e](https://gitea.com/gitea/runner/commits/97c083e902a86010ce6267347b830e91fe2435e1))
* align runner.os / runner.arch to known values ([#1510](https://gitea.com/gitea/runner/issues/1510)) ([6ab71ec](https://gitea.com/gitea/runner/commits/6ab71ecebbc0d01a0e72aecd06a15277723b6031))
* allow override of artifact server bind address ([#1560](https://gitea.com/gitea/runner/issues/1560)) ([d064863](https://gitea.com/gitea/runner/commits/d064863f9bb65cde863c8cb20fea50e0f139032c)), references [#1559](https://gitea.com/gitea/runner/issues/1559)
* allow TimeoutMinutes to be expression ([#1217](https://gitea.com/gitea/runner/issues/1217)) ([52f5c45](https://gitea.com/gitea/runner/commits/52f5c4592cbbf2e4dde3b35ad42f4e7475ad52df))
* allow TimeoutMinutes to be expression in Jobs ([#1247](https://gitea.com/gitea/runner/issues/1247)) ([91296bd](https://gitea.com/gitea/runner/commits/91296bd5eb022dce9e7891540257232b3e9906e1)), references [#1217](https://gitea.com/gitea/runner/issues/1217)
* always execute closeContainer() executor ([#988](https://gitea.com/gitea/runner/issues/988)) ([9abc87b](https://gitea.com/gitea/runner/commits/9abc87b41634cab2ad9325287fbed9135c2664c4)), references [#984](https://gitea.com/gitea/runner/issues/984)
* Apply forcePull only for prebuild docker actions ([#1599](https://gitea.com/gitea/runner/issues/1599)) ([f91b2aa](https://gitea.com/gitea/runner/commits/f91b2aa5dbcb16075f49c96a6108d0d418382f4d))
* avoid race condition by not running exec/start twice ([#702](https://gitea.com/gitea/runner/issues/702)) ([8a9167d](https://gitea.com/gitea/runner/commits/8a9167da82b3824272fdbfcc25b3671cd1260e7d)), references [#627](https://gitea.com/gitea/runner/issues/627) [#627](https://gitea.com/gitea/runner/issues/627)
* cache adjust restore order of exact key matches ([#2267](https://gitea.com/gitea/runner/issues/2267)) ([f825e42](https://gitea.com/gitea/runner/commits/f825e42ce200fc4973c3f28797ffba942d322d38))
* call of reflect.Value.Interface on zero Value ([#1067](https://gitea.com/gitea/runner/issues/1067)) ([380bd58](https://gitea.com/gitea/runner/commits/380bd581b1a56683a77b9eb017a72129f65251f8))
* checkout container path should always use `/` ([#1226](https://gitea.com/gitea/runner/issues/1226)) ([f71f277](https://gitea.com/gitea/runner/commits/f71f2778f0e75e3e28a81f8f726c64f27b7fffa3))
* ci is failing since 2022-10-17 ([#1397](https://gitea.com/gitea/runner/issues/1397)) ([5d2eb1f](https://gitea.com/gitea/runner/commits/5d2eb1f2388a7d99aae811721c70e75ed8b8a76c))
* CI is unstable due to new linting rules ([#1292](https://gitea.com/gitea/runner/issues/1292)) ([3387fd7](https://gitea.com/gitea/runner/commits/3387fd72ee138893b953e2b8b23efceb23bfefef))
* clarify --reuse and --rm flags ([#929](https://gitea.com/gitea/runner/issues/929)) ([4e0ba61](https://gitea.com/gitea/runner/commits/4e0ba618d34a198fa30904b1bf8549010513a931))
* clean code ([9fd6310](https://gitea.com/gitea/runner/commits/9fd63103632d5288cc7fa6c6f7e175e2ca17ad2f))
* cleanup unused actions ([#873](https://gitea.com/gitea/runner/issues/873)) ([3054f3e](https://gitea.com/gitea/runner/commits/3054f3ea1486267424d71c03c943eded2b95f50f))
* coerce booleans to numbers for comparison in exprparser ([#1030](https://gitea.com/gitea/runner/issues/1030)) ([aab2af0](https://gitea.com/gitea/runner/commits/aab2af09194a77a6b93d81088baf981cd5e27c17))
* conclusion and outcome are no integers ([#1136](https://gitea.com/gitea/runner/issues/1136)) ([77ddc89](https://gitea.com/gitea/runner/commits/77ddc89444d8cca5be47a3c83074dcb992382e68))
* continue jobs + steps after failure ([#840](https://gitea.com/gitea/runner/issues/840)) ([1891c72](https://gitea.com/gitea/runner/commits/1891c72ab158508e36009d16b24913fa5836422b)), closes [#442](https://gitea.com/gitea/runner/issues/442)
* copy ignored tracked files ([#1074](https://gitea.com/gitea/runner/issues/1074)) ([c27e079](https://gitea.com/gitea/runner/commits/c27e079407647e0ab455f068061b5e61e31f670f))
* copy whole action dir to container ([#774](https://gitea.com/gitea/runner/issues/774)) ([de44a50](https://gitea.com/gitea/runner/commits/de44a505da37f9405db195902b277821df2742b6))
* correct ref and ref_name ([#1672](https://gitea.com/gitea/runner/issues/1672)) ([6744e68](https://gitea.com/gitea/runner/commits/6744e68ee29d3bcb6878f895d6ccb5f9f8e2c478))
* correctly unmarshal entrypoint ([#876](https://gitea.com/gitea/runner/issues/876)) ([12fa4d7](https://gitea.com/gitea/runner/commits/12fa4d703da8d05930303acbe6bbc8fd00c77de7))
* crash if the id tool fails to run in the container (1660) ([f3c88b5](https://gitea.com/gitea/runner/commits/f3c88b50915a7234af9b06e28b4a45a267172fe1))
* crash malformed composite action ([#1616](https://gitea.com/gitea/runner/issues/1616)) ([53095d7](https://gitea.com/gitea/runner/commits/53095d76f48747f9f6249ba01e9410e9098d3977))
* crash on invalid docker credentials ([#1348](https://gitea.com/gitea/runner/issues/1348)) ([7d0407f](https://gitea.com/gitea/runner/commits/7d0407fa4abfe81308f8787cbd07c4c72b2aa676))
* deep evaluate matrix strategy ([#964](https://gitea.com/gitea/runner/issues/964)) ([d1e0216](https://gitea.com/gitea/runner/commits/d1e0216039d92a4253fecc2b11a775a37ba43ceb))
* do not error that much ([#1137](https://gitea.com/gitea/runner/issues/1137)) ([9357512](https://gitea.com/gitea/runner/commits/93575124d35f7b84860cb8cc553c257e0945d0af))
* docker buildx cache restore not working ([#2236](https://gitea.com/gitea/runner/issues/2236)) ([75e4ad9](https://gitea.com/gitea/runner/commits/75e4ad93f4c6ba51a7a64f12335dd0268dd7465a))
* docker cp of dangling symlink ([#943](https://gitea.com/gitea/runner/issues/943)) ([#948](https://gitea.com/gitea/runner/issues/948)) ([c802064](https://gitea.com/gitea/runner/commits/c802064975099a00c93b0c57f46b1e3964d33eef))
* don't allow `-self-hosted` mode as container image ([#1783](https://gitea.com/gitea/runner/issues/1783)) ([ca9b783](https://gitea.com/gitea/runner/commits/ca9b783491b4ea1c5d67dbffdb95f34fd4a22573))
* don't override env ([#1629](https://gitea.com/gitea/runner/issues/1629)) ([b0a5068](https://gitea.com/gitea/runner/commits/b0a5068f6d8ec06db4488b4b981d39f476d3abf6))
* don't overwrite with empty cmd/entrypoint ([#1076](https://gitea.com/gitea/runner/issues/1076)) ([24422bf](https://gitea.com/gitea/runner/commits/24422bf4e4e1e12b9a1e076b2f9e0a7b14c6f890))
* don't replace empty secrets ([#725](https://gitea.com/gitea/runner/issues/725)) ([96a0564](https://gitea.com/gitea/runner/commits/96a05645269241587fdcf0f1ccc49ce38047627c))
* dryrun skip local actions + enable Tests ([#1199](https://gitea.com/gitea/runner/issues/1199)) ([b7d380b](https://gitea.com/gitea/runner/commits/b7d380b3f096e079f2c1ba91a9377752fc4c9751))
* **editorconfig:** correct braces ([#1094](https://gitea.com/gitea/runner/issues/1094)) ([8d41f7a](https://gitea.com/gitea/runner/commits/8d41f7aa6b6702d9c1c3bfc372dd4c5533e88c8c))
* ensure all post steps are executed ([#1286](https://gitea.com/gitea/runner/issues/1286)) ([3364f95](https://gitea.com/gitea/runner/commits/3364f95569d8a39d440b18af69d29421fcc77b55))
* ensure networkmode "host" unless explicitly specified ([#1739](https://gitea.com/gitea/runner/issues/1739)) ([de06444](https://gitea.com/gitea/runner/commits/de0644499a555d240bd065f7c6d839e5162df46f))
* environment handling windows (host mode) ([#1732](https://gitea.com/gitea/runner/issues/1732)) ([9884da0](https://gitea.com/gitea/runner/commits/9884da0122b7fc45595c4f2950c1044fcce2c3e6))
* extra path lost in composite actions ([#1531](https://gitea.com/gitea/runner/issues/1531)) ([a53a1c2](https://gitea.com/gitea/runner/commits/a53a1c20005614cefd3d6519abd7f76d91467132)), references [#1528](https://gitea.com/gitea/runner/issues/1528)
* fail if no stages were found ([#1970](https://gitea.com/gitea/runner/issues/1970)) ([7286b43](https://gitea.com/gitea/runner/commits/7286b43b0ea0110cc2052848c8c1c4ab66bcd068))
* fail workflow if the job name is invalid ([#596](https://gitea.com/gitea/runner/issues/596)) ([94d736a](https://gitea.com/gitea/runner/commits/94d736a602dcfc7db2568d41a16e222bf1bcc416))
* fallback to unauthenticated pull ([#1774](https://gitea.com/gitea/runner/issues/1774)) ([aa21277](https://gitea.com/gitea/runner/commits/aa212773804561eaee04c3b4a7d9399dfbb9cad1))
* fromJSON should work with any JSON ([#883](https://gitea.com/gitea/runner/issues/883)) ([09693ec](https://gitea.com/gitea/runner/commits/09693ec5b9f0d2643b0a57826e2c415dfc32de2d))
* ghc assignment typo ([#1729](https://gitea.com/gitea/runner/issues/1729)) ([97749a2](https://gitea.com/gitea/runner/commits/97749a27b948268a921ef3e75d679a0f1f5c26ce))
* GITHUB_PATH should prepend ([#690](https://gitea.com/gitea/runner/issues/690)) ([29ea8cf](https://gitea.com/gitea/runner/commits/29ea8cfc4e00401b144cd29643a61ca089efd7e9))
* github.job property is empty, GITHUB_JOB should be job id ([#1646](https://gitea.com/gitea/runner/issues/1646)) ([5500c92](https://gitea.com/gitea/runner/commits/5500c928ebb6c4e7dcfb37c29c7100dd7e721b68)), references [#1621](https://gitea.com/gitea/runner/issues/1621) [#1473](https://gitea.com/gitea/runner/issues/1473)
* **goreleaser:** add append mode for release notes ([#962](https://gitea.com/gitea/runner/issues/962)) ([518148d](https://gitea.com/gitea/runner/commits/518148d1629381b0abd73e8f681d9fb6e0242205))
* handle env-vars case sensitive ([#1493](https://gitea.com/gitea/runner/issues/1493)) ([0f8861b](https://gitea.com/gitea/runner/commits/0f8861b9e37bfdf5e0fabb188d03365c0401823a)), closes [#1488](https://gitea.com/gitea/runner/issues/1488)
* handle go closure iteration ([#1374](https://gitea.com/gitea/runner/issues/1374)) ([1bade27](https://gitea.com/gitea/runner/commits/1bade27534209f2fd2d7a0b6437efb4965b648c8))
* handle subdirectors and exclusions ([#1012](https://gitea.com/gitea/runner/issues/1012)) ([#1446](https://gitea.com/gitea/runner/issues/1446)) ([5200c49](https://gitea.com/gitea/runner/commits/5200c491555ab1665cf280523f97bf92401c0415))
* handle zero size ([#1888](https://gitea.com/gitea/runner/issues/1888)) ([8c7c0f5](https://gitea.com/gitea/runner/commits/8c7c0f53c158867b446d1ef21356f49087c238a3))
* hardcode `ubuntu-latest` for `ImageOS` envvar ([#579](https://gitea.com/gitea/runner/issues/579)) ([29b3d43](https://gitea.com/gitea/runner/commits/29b3d43988e1354519f667830eef06662cadd9e4))
* hide masked value ([#668](https://gitea.com/gitea/runner/issues/668)) ([3db5037](https://gitea.com/gitea/runner/commits/3db50376aa4e0fbf5c9d206cb955c6f660c7779a))
* hide masked values in log files ([#1011](https://gitea.com/gitea/runner/issues/1011)) ([c22d833](https://gitea.com/gitea/runner/commits/c22d833830b3a1903a12c590b0b2ee637c3639e1))
* ignore `act` binary built when using `make` ([#654](https://gitea.com/gitea/runner/issues/654)) ([780f60a](https://gitea.com/gitea/runner/commits/780f60a1e6a07f495ad9a2d0e3181e85b36b0537))
* ignore symlink errors ([c44b71c](https://gitea.com/gitea/runner/commits/c44b71c9962f527e9c8847a7a3b959c71ac41b21))
* improve action not found error ([#2171](https://gitea.com/gitea/runner/issues/2171)) ([a6ec2c1](https://gitea.com/gitea/runner/commits/a6ec2c129a94a9cd488750a3be8b961cff08fd0c))
* improve new-action-cache fetch failure error ([#2172](https://gitea.com/gitea/runner/issues/2172)) ([6dd6725](https://gitea.com/gitea/runner/commits/6dd67253bc4ca18ee7ec5e96b8967332fe94dda1))
* improve warning about remote not found ([#2169](https://gitea.com/gitea/runner/issues/2169)) ([5e0d29d](https://gitea.com/gitea/runner/commits/5e0d29d6655ffb51ba394edd236ddf8e362698e0))
* interpolate action input defaults ([#1376](https://gitea.com/gitea/runner/issues/1376)) ([3631013](https://gitea.com/gitea/runner/commits/36310139e07b29768700e312c92ea49635f0d1ac))
* interpolate env for step actions ([#872](https://gitea.com/gitea/runner/issues/872)) ([099af5e](https://gitea.com/gitea/runner/commits/099af5e6a8e882cab7aa94e11d0616faceb0951f))
* invoke login shells to source `/etc/environment` ([#570](https://gitea.com/gitea/runner/issues/570)) ([b438b83](https://gitea.com/gitea/runner/commits/b438b836ea5e031c1e48cb8b2a777f1647d70f8a))
* IsHost is defined as false on windows ([#2093](https://gitea.com/gitea/runner/issues/2093)) ([1b10028](https://gitea.com/gitea/runner/commits/1b1002844772c5107f16321bd26ccfe7f6f0b48b))
* keep action inputs ([#1215](https://gitea.com/gitea/runner/issues/1215)) ([1d4c2aa](https://gitea.com/gitea/runner/commits/1d4c2aaa3fb1fce7a0fc70a3923391f20ecd6792))
* keep path to event json file in composite actions ([#1428](https://gitea.com/gitea/runner/issues/1428)) ([2614b3e](https://gitea.com/gitea/runner/commits/2614b3eb0c62772b17ba11509287fdc6f843f9d3))
* localcheckout mock ([#1198](https://gitea.com/gitea/runner/issues/1198)) ([c3fb686](https://gitea.com/gitea/runner/commits/c3fb6864e86542f5ff1f025bec5eff4f81497c96))
* logo and demo asset urls ([#2226](https://gitea.com/gitea/runner/issues/2226)) ([e1cd7c9](https://gitea.com/gitea/runner/commits/e1cd7c915f79120d77bdee6f9ec11eadb2da25bf))
* match cache `restore-keys` in creation reverse order ([#2153](https://gitea.com/gitea/runner/issues/2153)) ([7f7d84b](https://gitea.com/gitea/runner/commits/7f7d84b10ffbfbb699c0d95c17e94ce9a610ac44))
* Max parallel Support for Matrix Jobs and Remote Action Tests ([#150](https://gitea.com/gitea/runner/issues/150)) ([c0f19d9](https://gitea.com/gitea/runner/commits/c0f19d9a26134d5b4222aff3194f9d2b6f51d1d7))
* **mergify:** needs-work label handling ([#810](https://gitea.com/gitea/runner/issues/810)) ([a6aea44](https://gitea.com/gitea/runner/commits/a6aea44fb023f33543e58a5fcdcff937cf0e8ec6))
* missing defaults in nodejs pre ([#1349](https://gitea.com/gitea/runner/issues/1349)) ([6037182](https://gitea.com/gitea/runner/commits/6037182487e39d9f50c6887ab3453986fbf29a9b))
* mount volume for workflow info ([#776](https://gitea.com/gitea/runner/issues/776)) ([3c14548](https://gitea.com/gitea/runner/commits/3c1454825d8d4c622a4a94a9b3f796461d06409c))
* nil pointer access ( workflow_dispatch ) ([d97481d](https://gitea.com/gitea/runner/commits/d97481d5671a4906b8b9fe7b951333fa938ffe2b))
* outputs espace ([#404](https://gitea.com/gitea/runner/issues/404)) ([92067eb](https://gitea.com/gitea/runner/commits/92067eb1a5f271278373d6b2c2df9b6f9c373446))
* panic ([6791833](https://gitea.com/gitea/runner/commits/67918333fa7fac871f9de707970f0ef5f564ab61))
* panic ([#2071](https://gitea.com/gitea/runner/issues/2071)) ([5a79256](https://gitea.com/gitea/runner/commits/5a79256ee4d65ca068687bd0c85cd112d438a5d8))
* panic if a step in a job is nil ([#1145](https://gitea.com/gitea/runner/issues/1145)) ([91fd412](https://gitea.com/gitea/runner/commits/91fd412c51667dd10c434e9eb16ba10685581fa6))
* panic on regex in main.workflow ([a04141c](https://gitea.com/gitea/runner/commits/a04141c4440be5c0f003d5441e4e2a97bc09060e)), closes [#12](https://gitea.com/gitea/runner/issues/12)
* panic while running DRYRUN with local actions ([#1141](https://gitea.com/gitea/runner/issues/1141)) ([562265b](https://gitea.com/gitea/runner/commits/562265bb295f618123f65af2b9d5877638b2bda7))
* panic: reflect: slice index out of range ([#1066](https://gitea.com/gitea/runner/issues/1066)) ([a970145](https://gitea.com/gitea/runner/commits/a970145c95b9cc072905ca4f64c0640736aeff1f))
* pass `sc.Env` through `common.Logger` to hide secrets ([#578](https://gitea.com/gitea/runner/issues/578)) ([f3b53d8](https://gitea.com/gitea/runner/commits/f3b53d8eca1142fb587d984d8c95584e236b612a))
* preserve job result state in case of failure ([#1519](https://gitea.com/gitea/runner/issues/1519)) ([b14398e](https://gitea.com/gitea/runner/commits/b14398eff3d3a2692c976e250c4b91fbbeed00a5)), closes [#1518](https://gitea.com/gitea/runner/issues/1518)
* print line after password prompt ([#1014](https://gitea.com/gitea/runner/issues/1014)) ([ad06f5d](https://gitea.com/gitea/runner/commits/ad06f5dfb8a658521c855f8c374a58b02227e8fa))
* processing of strategy.matrix.include ([#1200](https://gitea.com/gitea/runner/issues/1200)) ([c30bc82](https://gitea.com/gitea/runner/commits/c30bc824b23335f2c99a32578f8d907f0be6647e))
* re-evaluate env for remote composite actions ([#1385](https://gitea.com/gitea/runner/issues/1385)) ([f0b1845](https://gitea.com/gitea/runner/commits/f0b1845802ca6bea4a75f4e3b514fad85f070b4f))
* regression run after failure ([#971](https://gitea.com/gitea/runner/issues/971)) ([7dbf3fc](https://gitea.com/gitea/runner/commits/7dbf3fcb96e2d70b3436b718c1f4a9fe4b169a4f))
* remove composite restrictions ([#1128](https://gitea.com/gitea/runner/issues/1128)) ([ebb408f](https://gitea.com/gitea/runner/commits/ebb408f373a4fc6508d1d75aa6fddcba5fc7df81))
* remove hardcoded reference to github.com when using reusable remote workflows and remote actions ([#1784](https://gitea.com/gitea/runner/issues/1784)) ([4721abf](https://gitea.com/gitea/runner/commits/4721abfa6dbc661922d34d20b9f60aa14eb30b34))
* remove HOME env var ([#594](https://gitea.com/gitea/runner/issues/594)) ([5044ec6](https://gitea.com/gitea/runner/commits/5044ec6c4364cb74e623953b8c5d39790c8af126))
* remove unused TestContext ([#1013](https://gitea.com/gitea/runner/issues/1013)) ([1e92d87](https://gitea.com/gitea/runner/commits/1e92d87b75f1691a531c76e941e614b521e308f2))
* Remove volume with same name as the container again ([#1203](https://gitea.com/gitea/runner/issues/1203)) ([28f5b37](https://gitea.com/gitea/runner/commits/28f5b37fd1941e00e77f33e3964011f83789c45c))
* remove Windows issue from `README.md` ([#675](https://gitea.com/gitea/runner/issues/675)) ([cc4e23d](https://gitea.com/gitea/runner/commits/cc4e23d96c24effff203c86d679c22c9ab092f27))
* replace `gopass` lib with `golang.org/x/term` ([#674](https://gitea.com/gitea/runner/issues/674)) ([a6a865e](https://gitea.com/gitea/runner/commits/a6a865e973d7138c70e9699d2e7a9a11b542f16c))
* restore the job logger setup after job cancelation ([#1365](https://gitea.com/gitea/runner/issues/1365)) ([13d5303](https://gitea.com/gitea/runner/commits/13d530302e5f8203b0f104155eb6e136f5aaf910))
* return err in walk ([#1667](https://gitea.com/gitea/runner/issues/1667)) ([ac5dd8f](https://gitea.com/gitea/runner/commits/ac5dd8feb876d37ae483376a137c57383577dace))
* return error if both `run:`/`uses:` keys are used ([#593](https://gitea.com/gitea/runner/issues/593)) ([25779af](https://gitea.com/gitea/runner/commits/25779af4bf1e7762d8ff2d3be985c2fcde7bb1e2))
* return error on reusable workflows ([#1093](https://gitea.com/gitea/runner/issues/1093)) ([b3bd268](https://gitea.com/gitea/runner/commits/b3bd268e4585b32673635685d68bf86c5b2d0bfa))
* return invalid step type ([#1157](https://gitea.com/gitea/runner/issues/1157)) ([1e72c59](https://gitea.com/gitea/runner/commits/1e72c594aca1854854d40a24f3e038537f392fc8))
* reusable workflow panic ([#1728](https://gitea.com/gitea/runner/issues/1728)) ([d70b225](https://gitea.com/gitea/runner/commits/d70b225e85b44efdcd97452959a3bbabec5eb6ed))
* rework setupShellCommand ([#930](https://gitea.com/gitea/runner/issues/930)) ([adabf2a](https://gitea.com/gitea/runner/commits/adabf2a2027fa9da3410febc474d55327e2abeb7))
* reworked container architecture ([#619](https://gitea.com/gitea/runner/issues/619)) ([f32babb](https://gitea.com/gitea/runner/commits/f32babb51db6761fd5dff1bcc03bbdf9eb7aa354))
* rewrite how image env is merged ([#828](https://gitea.com/gitea/runner/issues/828)) ([6c60af7](https://gitea.com/gitea/runner/commits/6c60af76779b739b7487e3fbe83ea15cd2ab8ad5))
* rewrite StepContext env ([#832](https://gitea.com/gitea/runner/issues/832)) ([a0d3602](https://gitea.com/gitea/runner/commits/a0d360236ef5466a441631c2f93794b55d4b76c5))
* Rootless error concerning `/var/run/docker.sock` ([#2181](https://gitea.com/gitea/runner/issues/2181)) ([f2e65e1](https://gitea.com/gitea/runner/commits/f2e65e1d406255cad05b562138758b73219fa74d))
* rootless permission bits (new actions cache) ([#2242](https://gitea.com/gitea/runner/issues/2242)) ([119ceb8](https://gitea.com/gitea/runner/commits/119ceb81d906be5e1b23bd085f642ffde7144ef5))
* safe file name ([#1651](https://gitea.com/gitea/runner/issues/1651)) ([19abab6](https://gitea.com/gitea/runner/commits/19abab637549dae95d4ea3becd77852e2787195b))
* set `sh` as default shell for containers  ([#853](https://gitea.com/gitea/runner/issues/853)) ([e9ae789](https://gitea.com/gitea/runner/commits/e9ae7894e38ed276c41e6cd7b958a81da79249ad))
* set composite outputs on failure ([#945](https://gitea.com/gitea/runner/issues/945)) ([4be9062](https://gitea.com/gitea/runner/commits/4be9062dd2745c1657ad7b98683397a393159436))
* set default version string ([#1070](https://gitea.com/gitea/runner/issues/1070)) ([de81979](https://gitea.com/gitea/runner/commits/de8197980079129395c67b1e2865cff8be4ebca8))
* set logger with trace level ([97629ae](https://gitea.com/gitea/runner/commits/97629ae8af8dbe76387c34e2b053b3645ceb1bde))
* show workflow info even if on.push is not defined ([#1329](https://gitea.com/gitea/runner/issues/1329)) ([#1335](https://gitea.com/gitea/runner/issues/1335)) ([21484b5](https://gitea.com/gitea/runner/commits/21484b5c1ecbd9e1ddb6a57debd7731a3e1fc25f)), closes [#1332](https://gitea.com/gitea/runner/issues/1332) [#1318](https://gitea.com/gitea/runner/issues/1318)
* show workflow name when get error on yml parsing ([#1133](https://gitea.com/gitea/runner/issues/1133)) ([#1328](https://gitea.com/gitea/runner/issues/1328)) ([b514649](https://gitea.com/gitea/runner/commits/b514649c3da64a1be92b8e22b50e296a07d759f6))
* skip local actions pre step in any case ([#1204](https://gitea.com/gitea/runner/issues/1204)) ([8a47394](https://gitea.com/gitea/runner/commits/8a473943c3dffbbf47c2a2ef8b261e845646fa6e)), closes [#1193](https://gitea.com/gitea/runner/issues/1193)
* spelling mistake in readme ([#1854](https://gitea.com/gitea/runner/issues/1854)) ([80f6de5](https://gitea.com/gitea/runner/commits/80f6de51d544e2c06d730640d683814e247ec0c2))
* step env is unavailable in with property expr ([#1458](https://gitea.com/gitea/runner/issues/1458)) ([8c5748a](https://gitea.com/gitea/runner/commits/8c5748a55c6190ba34354ca29fdb473e5d3724a1))
* subpath actions via new artifact cache ([#2170](https://gitea.com/gitea/runner/issues/2170)) ([09d4b5d](https://gitea.com/gitea/runner/commits/09d4b5d6ade1a06fa006a0fb202f9ff8c19c33f6))
* support docker create arguments from container.options ([#1022](https://gitea.com/gitea/runner/issues/1022)) ([#1351](https://gitea.com/gitea/runner/issues/1351)) ([48188a6](https://gitea.com/gitea/runner/commits/48188a6280f367e63936a239976be65317ebabd2))
* support expression for step's continue-on-error field ([#900](https://gitea.com/gitea/runner/issues/900)) ([#1331](https://gitea.com/gitea/runner/issues/1331)) ([d1daf2f](https://gitea.com/gitea/runner/commits/d1daf2f28d8c2256e5712e3064ce6ac45fca2562))
* tail (not absolute) as entrypoint of job container ([#1506](https://gitea.com/gitea/runner/issues/1506)) ([0c8c082](https://gitea.com/gitea/runner/commits/0c8c082ac0baf25909fd79b87897426816c3c4eb))
* **test:** adapt to sjw change ([ae39a4b](https://gitea.com/gitea/runner/commits/ae39a4b1d354e13891ef8bd7eb8bcd4a0916f0e6))
* **test:** needs condition. ([#8](https://gitea.com/gitea/runner/issues/8)) ([0efa2d5](https://gitea.com/gitea/runner/commits/0efa2d5e63fdc03360b4eb6aac0b9ac5eac9ca42))
* the number in the github event is of type number ([#1252](https://gitea.com/gitea/runner/issues/1252)) ([4094462](https://gitea.com/gitea/runner/commits/409446211fe99420d5d3761b04f491f8ef52412d))
* tojson(null) == 'null' ([#1068](https://gitea.com/gitea/runner/issues/1068)) ([5083284](https://gitea.com/gitea/runner/commits/5083284deb804db5b9d5997ab4cf345b1c7ec069))
* tolerate workflow that needs a missing job ([#1595](https://gitea.com/gitea/runner/issues/1595)) ([#1619](https://gitea.com/gitea/runner/issues/1619)) ([be4a147](https://gitea.com/gitea/runner/commits/be4a1477a5b73b3e17f79891ec6147413acb4df7))
* trim whitespace from actrc ([#1024](https://gitea.com/gitea/runner/issues/1024)) ([3db3c73](https://gitea.com/gitea/runner/commits/3db3c737230adc457e08a8d60b34458977d01e5c)), closes [#985](https://gitea.com/gitea/runner/issues/985)
* try finding a socket, otherwise fail, respect user choice ([#1745](https://gitea.com/gitea/runner/issues/1745)) ([b6718fd](https://gitea.com/gitea/runner/commits/b6718fdf5d2fdc3e61798f4ac749b53548cbbfae))
* update `ubuntu-latest` to point to `ubuntu-20.04` ([#580](https://gitea.com/gitea/runner/issues/580)) ([72f2d2d](https://gitea.com/gitea/runner/commits/72f2d2de51fcc309a2f0381a9b5098c81802598b))
* Update ARCH environment variable used in runners/actions ([#1818](https://gitea.com/gitea/runner/issues/1818)) ([11dd2ac](https://gitea.com/gitea/runner/commits/11dd2ac7457bd22080b9d8e5277327110563919a))
* update artifact server to address GHSL-2023-004 ([#1565](https://gitea.com/gitea/runner/issues/1565)) ([63ae215](https://gitea.com/gitea/runner/commits/63ae215071f94569d910964bdee866d91d6e3a10))
* update codecov action ([#727](https://gitea.com/gitea/runner/issues/727)) ([05bce00](https://gitea.com/gitea/runner/commits/05bce00e932e450ec3cd0d73c064e2a396ac9500))
* update codecov action (again) ([#731](https://gitea.com/gitea/runner/issues/731)) ([6550eb7](https://gitea.com/gitea/runner/commits/6550eb7f2515b635c1210fc218e03132008463c6))
* update install.sh ([#937](https://gitea.com/gitea/runner/issues/937)) ([e9c4609](https://gitea.com/gitea/runner/commits/e9c4609dcadf70739e34f79680e69df973b24d8e))
* update issue template ([#722](https://gitea.com/gitea/runner/issues/722)) ([fa93f4d](https://gitea.com/gitea/runner/commits/fa93f4d5e7e9ccfe68741a555e3137c1dfd6e997))
* update output handling for reusable workflows ([#1521](https://gitea.com/gitea/runner/issues/1521)) ([89cb558](https://gitea.com/gitea/runner/commits/89cb55855802fd7997bf9128140ee4e48b7e89c1))
* update readme with more relevant/updated information ([#630](https://gitea.com/gitea/runner/issues/630)) ([08141f5](https://gitea.com/gitea/runner/commits/08141f5b0095f0ddbbf3a62b03338a9570e27658))
* use actions/runner hashfiles in container ([#1940](https://gitea.com/gitea/runner/issues/1940)) ([7c7d80e](https://gitea.com/gitea/runner/commits/7c7d80ebdd4aee01bb858a828ad864aa120b4a75))
* use auth when fetching (same as with cloning) ([#687](https://gitea.com/gitea/runner/issues/687)) ([38f6dfb](https://gitea.com/gitea/runner/commits/38f6dfb49a7e3cee826911e36f0c990813282d28))
* use correct path to toolcache ([#1494](https://gitea.com/gitea/runner/issues/1494)) ([054caec](https://gitea.com/gitea/runner/commits/054caec791b2e98a917210ad4acdb46181372e09))
* use docker images from dockerhub ([#1249](https://gitea.com/gitea/runner/issues/1249)) ([e70b968](https://gitea.com/gitea/runner/commits/e70b9689241785df1ff509678d99f1fdcb138d91))
* use docker lib for image ref parsing ([#877](https://gitea.com/gitea/runner/issues/877)) ([ec34eb9](https://gitea.com/gitea/runner/commits/ec34eb9532d314ba8d3c23472acf4037f2c084cd))
* use github auth token for pulling a repository ([#693](https://gitea.com/gitea/runner/issues/693)) ([60b9606](https://gitea.com/gitea/runner/commits/60b9606fc687549c7f55343a4d52ff3f2c32cf4d))
* use new platforms after survey, check working dir for `.actrc` ([#577](https://gitea.com/gitea/runner/issues/577)) ([780a8a0](https://gitea.com/gitea/runner/commits/780a8a061c2c7a3a6958e444e105d7d2441893a9))
* use os.UserHomeDir ([#1706](https://gitea.com/gitea/runner/issues/1706)) ([d3dfde0](https://gitea.com/gitea/runner/commits/d3dfde055aad9c0e9631fd6bbf06f4b0db706eff))
* using filepath instead of path ([01f7ff3](https://gitea.com/gitea/runner/commits/01f7ff315c1b3eaa5e3a4eee06bd6d88ebf60ce4))
* warn Apple M1 (arm64) users ([#743](https://gitea.com/gitea/runner/issues/743)) ([1cf422e](https://gitea.com/gitea/runner/commits/1cf422e41108bf72ec837bea66e080d400ca1192))
* write default config in XDG config dir to avoid cluttering the HOME directory by default ([#2140](https://gitea.com/gitea/runner/issues/2140)) ([6091094](https://gitea.com/gitea/runner/commits/6091094e14ab4e2940bc4bde12525d2dfd569e0b))
* wrong perms ([#726](https://gitea.com/gitea/runner/issues/726)) ([2261204](https://gitea.com/gitea/runner/commits/2261204c65b71d661495de95df6610c451157a67))

### Refactoring

* `NewWorkflowPlanner` ([#648](https://gitea.com/gitea/runner/issues/648)) ([806bc4d](https://gitea.com/gitea/runner/commits/806bc4d999cafab455bbda9b966e3ed6080a4104))
* **cmd/root:** simplify `parseEnvs` ([#2162](https://gitea.com/gitea/runner/issues/2162)) ([424fd5e](https://gitea.com/gitea/runner/commits/424fd5e02b4be9fce8aae518420b0499c08cb31b))
* default empty if: statement to if: success() in evaluator & remove `FixIfStatement()` ([#990](https://gitea.com/gitea/runner/issues/990)) ([ff13844](https://gitea.com/gitea/runner/commits/ff13844b8643abaff9daa47a5bc96aa91fddcf5e))
* docker build BuildContext field ([#1914](https://gitea.com/gitea/runner/issues/1914)) ([f0ca0ab](https://gitea.com/gitea/runner/commits/f0ca0abc4003e1d99a534564fc1ad62f47eaf8a4))
* export and move shared contexts into pkg/model ([#931](https://gitea.com/gitea/runner/issues/931)) ([ed01f46](https://gitea.com/gitea/runner/commits/ed01f464ede3515c94dedf5739b54576177bf80c)), references [#908](https://gitea.com/gitea/runner/issues/908)
* extract RunContext Executor in JobExecutor ([#984](https://gitea.com/gitea/runner/issues/984)) ([e23223a](https://gitea.com/gitea/runner/commits/e23223ad02ddeafe581df97c7e9462df28d8e104))
* extract setupAction into ActionReader ([#986](https://gitea.com/gitea/runner/issues/986)) ([e4f0080](https://gitea.com/gitea/runner/commits/e4f0080a1efc9887fe8aa892669b2002e70b6040))
* filecollector into new package ([#2174](https://gitea.com/gitea/runner/issues/2174)) ([5a80a04](https://gitea.com/gitea/runner/commits/5a80a044f986a48b4d2ed380815e9bb47bca2b66))
* fix add-path / GITHUB_PATH commands ([#1472](https://gitea.com/gitea/runner/issues/1472)) ([d281230](https://gitea.com/gitea/runner/commits/d281230ccef6438085180866e50a026025a05f67))
* fix savestate in pre steps ([#1466](https://gitea.com/gitea/runner/issues/1466)) ([b8d7e94](https://gitea.com/gitea/runner/commits/b8d7e947cf91e24202326b63821a3ff82a3e2b57))
* GITHUB_ENV command / remove env.PATH ([#1503](https://gitea.com/gitea/runner/issues/1503)) ([e775fea](https://gitea.com/gitea/runner/commits/e775fea2655d9e3342bd06c2743cb4bab791e2ab))
* ignore already closed error ([#1285](https://gitea.com/gitea/runner/issues/1285)) ([92ddcda](https://gitea.com/gitea/runner/commits/92ddcdae098873ad6d9c05748ce5b6c7be83b811))
* move autoremove into the jobexecutor ([#1463](https://gitea.com/gitea/runner/issues/1463)) ([d9fe63e](https://gitea.com/gitea/runner/commits/d9fe63ec249d4ee3468151851b8db642dc32bf2c))
* move from io/ioutil to io and os packages ([#1417](https://gitea.com/gitea/runner/issues/1417)) ([2f1c5a1](https://gitea.com/gitea/runner/commits/2f1c5a19f12ad0fc784924b4585d4ace061c7ef5))
* open boltdb only while using it ([#1879](https://gitea.com/gitea/runner/issues/1879)) ([e597046](https://gitea.com/gitea/runner/commits/e597046195774991a75ff66506fdfcc3fc903364))
* pull and rebuild docker by default ([#1569](https://gitea.com/gitea/runner/issues/1569)) ([b00babd](https://gitea.com/gitea/runner/commits/b00babd0d9fe8538aced303346cc7139ef51d0cd))
* re-implement `embed` without "unused" import ([#830](https://gitea.com/gitea/runner/issues/830)) ([7a426a0](https://gitea.com/gitea/runner/commits/7a426a0f37d21e32a01e85d9403edba17a54d545))
* remove `gotest.tools` ([#688](https://gitea.com/gitea/runner/issues/688)) ([f571290](https://gitea.com/gitea/runner/commits/f571290b25cdbb34e863cb98eb155369fff5d572))
* remove composite action runcontext workaround ([#1085](https://gitea.com/gitea/runner/issues/1085)) ([e360811](https://gitea.com/gitea/runner/commits/e3608115706382233c3b943b4dabea09bab9d7d2))
* remove docker image list reference filter ([#1501](https://gitea.com/gitea/runner/issues/1501)) ([b2fb9e6](https://gitea.com/gitea/runner/commits/b2fb9e64ac1b263088003b94963b36b6edeaba31))
* remove github.com/pkg/errors dependency ([#1077](https://gitea.com/gitea/runner/issues/1077)) ([2aa0699](https://gitea.com/gitea/runner/commits/2aa0699aecff2e596e280165ea48eac85c5f9a6a))
* remove unnecessary nil check in RunContext ([#1955](https://gitea.com/gitea/runner/issues/1955)) ([8a9e4f9](https://gitea.com/gitea/runner/commits/8a9e4f9f38340ff9b7f9727add3dae62f3f6c7bb))
* return more errors, add more tests ([#679](https://gitea.com/gitea/runner/issues/679)) ([aa68181](https://gitea.com/gitea/runner/commits/aa68181f6b3df2ff24a523800871c3c5ccdf76be))
* share UpdateFromEnv logic ([#1457](https://gitea.com/gitea/runner/issues/1457)) ([57bf4d2](https://gitea.com/gitea/runner/commits/57bf4d27a2b6d03a2c61f39ba57c24138ae94370))
* simpilfy go-git cache ([#2208](https://gitea.com/gitea/runner/issues/2208)) ([5601fb0](https://gitea.com/gitea/runner/commits/5601fb0e136fad09282705ef65cbc9b2c6f0bebf))
* simplify action function signatures ([#1083](https://gitea.com/gitea/runner/issues/1083)) ([5f673cb](https://gitea.com/gitea/runner/commits/5f673cbb6db5569334cc806cf0536a97ad6b899b))
* simplify adding new node versions add node 24 ([#140](https://gitea.com/gitea/runner/issues/140)) ([91852fa](https://gitea.com/gitea/runner/commits/91852faf934e3436fe427d70017b9a630d1421df))

### Build

* **deps:** bump actions/cache from 2 to 3 ([#1082](https://gitea.com/gitea/runner/issues/1082)) ([4f0bfd2](https://gitea.com/gitea/runner/commits/4f0bfd2ad1b16dc0a5f34844069513af178fe86d))
* **deps:** bump actions/checkout from 2 to 3 ([#1028](https://gitea.com/gitea/runner/issues/1028)) ([0579f7c](https://gitea.com/gitea/runner/commits/0579f7c23598fd0b76500b95e5d69d55a3fb0b98))
* **deps:** bump actions/checkout from 3 to 4 ([#1998](https://gitea.com/gitea/runner/issues/1998)) ([0dfb067](https://gitea.com/gitea/runner/commits/0dfb06748ebd077c25451e07a55e563d7bdc6b4a))
* **deps:** bump actions/github-script from 5 to 6 ([#1572](https://gitea.com/gitea/runner/issues/1572)) ([75d19d0](https://gitea.com/gitea/runner/commits/75d19d0ff58e54cafa924b9bf72a89176a602f2b))
* **deps:** bump actions/github-script from 6 to 7 ([#2097](https://gitea.com/gitea/runner/issues/2097)) ([3d0cb3d](https://gitea.com/gitea/runner/commits/3d0cb3d82bf9e0448a245a17aad7a80e672e76a8))
* **deps:** bump actions/setup-go from 2 to 3 ([#1120](https://gitea.com/gitea/runner/issues/1120)) ([4e1cf1d](https://gitea.com/gitea/runner/commits/4e1cf1d7cc84a6da069401ba1ff08d14cf71044b))
* **deps:** bump actions/setup-go from 3 to 4 ([#1689](https://gitea.com/gitea/runner/issues/1689)) ([8c28c9f](https://gitea.com/gitea/runner/commits/8c28c9fd8f2ae53af23367be2cda094ecc4fbfe0))
* **deps:** bump actions/setup-go from 4 to 5 ([#2118](https://gitea.com/gitea/runner/issues/2118)) ([9cecf94](https://gitea.com/gitea/runner/commits/9cecf94039eed005aa76bb4886d9ebfa16b6d184))
* **deps:** bump actions/stale from 4 to 5 ([#1117](https://gitea.com/gitea/runner/issues/1117)) ([f551165](https://gitea.com/gitea/runner/commits/f551165970914f2541f692716c522e95d06a0608))
* **deps:** bump actions/stale from 5 to 6 ([#1357](https://gitea.com/gitea/runner/issues/1357)) ([9a494ca](https://gitea.com/gitea/runner/commits/9a494cac1ed5f22539d727f6958f4d24df7104b8))
* **deps:** bump actions/stale from 6 to 7 ([#1535](https://gitea.com/gitea/runner/issues/1535)) ([d3ce2b0](https://gitea.com/gitea/runner/commits/d3ce2b0a467b39616059b19c5ba0e69e145099aa))
* **deps:** bump actions/stale from 7 to 8 ([#1700](https://gitea.com/gitea/runner/issues/1700)) ([d65cf86](https://gitea.com/gitea/runner/commits/d65cf869e6824ac8aca448df8d7e0bc8b8fc21af))
* **deps:** bump actions/stale from 8 to 9 ([#2120](https://gitea.com/gitea/runner/issues/2120)) ([6b0ef97](https://gitea.com/gitea/runner/commits/6b0ef97c52e7cbab64513b3de7e6618bc2579f28))
* **deps:** bump actions/upload-artifact from 2 to 3 ([#1119](https://gitea.com/gitea/runner/issues/1119)) ([39cc9b0](https://gitea.com/gitea/runner/commits/39cc9b0840f07953b44416b62d4817ca375fe932))
* **deps:** bump actions/upload-artifact from 3 to 4 ([#2133](https://gitea.com/gitea/runner/issues/2133)) ([15bb54f](https://gitea.com/gitea/runner/commits/15bb54f14e96734963b78310f0cb8ed0c67e25bb))
* **deps:** bump codecov/codecov-action from 1.5.2 to 2.1.0 ([#819](https://gitea.com/gitea/runner/issues/819)) ([e82f17e](https://gitea.com/gitea/runner/commits/e82f17e230a7ae5b715b041b6827825de671dc9e))
* **deps:** bump codecov/codecov-action from 2.1.0 to 3.0.0 ([#1118](https://gitea.com/gitea/runner/issues/1118)) ([74cda72](https://gitea.com/gitea/runner/commits/74cda72c2672a1f665d5a55416cd93ca6292180a))
* **deps:** bump codecov/codecov-action from 3.0.0 to 3.1.0 ([#1139](https://gitea.com/gitea/runner/issues/1139)) ([3f2469b](https://gitea.com/gitea/runner/commits/3f2469bf5385b90da5f6a5d9922ccf9c9d7274a6))
* **deps:** bump codecov/codecov-action from 3.1.0 to 3.1.1 ([#1359](https://gitea.com/gitea/runner/issues/1359)) ([d137ec8](https://gitea.com/gitea/runner/commits/d137ec82852083eef45009dcaf368c3dcdd0b9a6))
* **deps:** bump codecov/codecov-action from 3.1.1 to 3.1.2 ([#1735](https://gitea.com/gitea/runner/issues/1735)) ([6d527bf](https://gitea.com/gitea/runner/commits/6d527bf1a356b69ee4ee35e43159c9d0a21c7718))
* **deps:** bump codecov/codecov-action from 3.1.2 to 3.1.3 ([#1752](https://gitea.com/gitea/runner/issues/1752)) ([60bcfb5](https://gitea.com/gitea/runner/commits/60bcfb5d4f717bd0a581de8e97866d58a5593152))
* **deps:** bump codecov/codecov-action from 3.1.3 to 3.1.4 ([#1813](https://gitea.com/gitea/runner/issues/1813)) ([1371215](https://gitea.com/gitea/runner/commits/1371215aabfd82b9288132042aecc9eacc0ab380))
* **deps:** bump codecov/codecov-action from 3.1.4 to 3.1.5 ([#2175](https://gitea.com/gitea/runner/issues/2175)) ([4ca35d2](https://gitea.com/gitea/runner/commits/4ca35d2192ccd6bddf4ebb9519d5fe27a48a6132))
* **deps:** bump deps which dependabot doesn't want to ([#1107](https://gitea.com/gitea/runner/issues/1107)) ([afb23a2](https://gitea.com/gitea/runner/commits/afb23a289187f4a90241b8b923555b7c3d844ec9))
* **deps:** bump docker/setup-qemu-action from 1 to 2 ([#1156](https://gitea.com/gitea/runner/issues/1156)) ([7974a17](https://gitea.com/gitea/runner/commits/7974a17afc613431461f9a8c629597962b1f54b4))
* **deps:** bump docker/setup-qemu-action from 2 to 3 ([#2007](https://gitea.com/gitea/runner/issues/2007)) ([1891bef](https://gitea.com/gitea/runner/commits/1891bef4335ba55573de1c7c5824acd44b3a4994))
* **deps:** bump fregante/setup-git-user from 1 to 2 ([#1664](https://gitea.com/gitea/runner/issues/1664)) ([fad986a](https://gitea.com/gitea/runner/commits/fad986af5f95a2878c9d9f15ded8d815fdf5e061))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.2.7 to 2.3.1 ([#803](https://gitea.com/gitea/runner/issues/803)) ([a7a2659](https://gitea.com/gitea/runner/commits/a7a2659c0e90a53fee4e55817d0eedcad1ab9967))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.2 to 2.3.4 ([#1079](https://gitea.com/gitea/runner/issues/1079)) ([efd34f1](https://gitea.com/gitea/runner/commits/efd34f1e10e56b251a5b9eb75f269d2e7c15681f))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.4 to 2.3.5 ([#1214](https://gitea.com/gitea/runner/issues/1214)) ([084c4e6](https://gitea.com/gitea/runner/commits/084c4e61d99c1bfd109ea5bc381342d94a6e65cd))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.5 to 2.3.6 ([#1341](https://gitea.com/gitea/runner/issues/1341)) ([3a0fe69](https://gitea.com/gitea/runner/commits/3a0fe6967fd8ecd3cb86550d0225d55a0bb37ac9))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.6 to 2.3.7 ([#1872](https://gitea.com/gitea/runner/issues/1872)) ([df44dff](https://gitea.com/gitea/runner/commits/df44dffd30b6503b91bc63130742c31b1ce7badf))
* **deps:** bump github.com/cloudflare/circl from 1.1.0 to 1.3.3 ([#1800](https://gitea.com/gitea/runner/issues/1800)) ([046f5aa](https://gitea.com/gitea/runner/commits/046f5aa71587ba5e81a3f755c913da2f0644cdfb))
* **deps:** bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 ([#2149](https://gitea.com/gitea/runner/issues/2149)) ([96d6cf8](https://gitea.com/gitea/runner/commits/96d6cf8b2c9587912e32a2a42f212fa710efb08e))
* **deps:** bump github.com/containerd/containerd from 1.6.16 to 1.6.18 ([#1637](https://gitea.com/gitea/runner/issues/1637)) ([4eba04b](https://gitea.com/gitea/runner/commits/4eba04b229388c4fc8bbbcffb5e4d62cd927eac9))
* **deps:** bump github.com/containerd/containerd from 1.6.19 to 1.6.26 ([#2189](https://gitea.com/gitea/runner/issues/2189)) ([df61c7f](https://gitea.com/gitea/runner/commits/df61c7fcdb1bb95a68f68238d1637d39ca0329f7))
* **deps:** bump github.com/containerd/containerd from 1.6.6 to 1.6.12 ([#1554](https://gitea.com/gitea/runner/issues/1554)) ([3ac756b](https://gitea.com/gitea/runner/commits/3ac756b37a496020ec8bb7b5852272d01c471057))
* **deps:** bump github.com/containerd/containerd from 1.7.2 to 1.7.11 ([#2136](https://gitea.com/gitea/runner/issues/2136)) ([74dcce4](https://gitea.com/gitea/runner/commits/74dcce467d84b5507db338e98cc8129dd86e90b1))
* **deps:** bump github.com/creack/pty from 1.1.17 to 1.1.18 ([#1447](https://gitea.com/gitea/runner/issues/1447)) ([ccb3b0e](https://gitea.com/gitea/runner/commits/ccb3b0e8757158bb945cf9d3b12a4a9299eea9c4))
* **deps:** bump github.com/creack/pty from 1.1.18 to 1.1.20 ([#2068](https://gitea.com/gitea/runner/issues/2068)) ([fb4f29f](https://gitea.com/gitea/runner/commits/fb4f29fd6d38c6ce57654be5799cb0c911820de4))
* **deps:** bump github.com/creack/pty from 1.1.20 to 1.1.21 ([#2099](https://gitea.com/gitea/runner/issues/2099)) ([603b44b](https://gitea.com/gitea/runner/commits/603b44b585ee35090a13498c7af15f79a8a12a51))
* **deps:** bump github.com/docker/cli ([#1038](https://gitea.com/gitea/runner/issues/1038)) ([933413e](https://gitea.com/gitea/runner/commits/933413e2c0e1d600e09291ba6e63211a3c294385))
* **deps:** bump github.com/docker/cli ([#1081](https://gitea.com/gitea/runner/issues/1081)) ([34706f4](https://gitea.com/gitea/runner/commits/34706f49fe8c82b3a17c9b22bb4fc5d8a70bd25e))
* **deps:** bump github.com/docker/cli ([#1154](https://gitea.com/gitea/runner/issues/1154)) ([d5e9721](https://gitea.com/gitea/runner/commits/d5e972159f9ee43b9229bf028853755c7f23c4d6))
* **deps:** bump github.com/docker/cli ([#1168](https://gitea.com/gitea/runner/issues/1168)) ([1fed666](https://gitea.com/gitea/runner/commits/1fed666165079c8fcab76ff11b66a60463ee0f1f))
* **deps:** bump github.com/docker/cli ([#1211](https://gitea.com/gitea/runner/issues/1211)) ([87d9687](https://gitea.com/gitea/runner/commits/87d9687716a279525fd952b80ba6f044bce88492))
* **deps:** bump github.com/docker/cli ([#1340](https://gitea.com/gitea/runner/issues/1340)) ([ceac1f2](https://gitea.com/gitea/runner/commits/ceac1f2429073270e984ce7da39e28f31e2c56c6))
* **deps:** bump github.com/docker/cli ([#1393](https://gitea.com/gitea/runner/issues/1393)) ([bd5e0d2](https://gitea.com/gitea/runner/commits/bd5e0d24d68cddab105ee76bc2fb8000dce76a11))
* **deps:** bump github.com/docker/cli ([#1405](https://gitea.com/gitea/runner/issues/1405)) ([ca69531](https://gitea.com/gitea/runner/commits/ca69531096e88bec7f8fbd6aecef9b3769271cb7))
* **deps:** bump github.com/docker/cli ([#1419](https://gitea.com/gitea/runner/issues/1419)) ([ac6571c](https://gitea.com/gitea/runner/commits/ac6571c8af87ca3e26b4a5eed1d0a919b7a7aed3))
* **deps:** bump github.com/docker/cli ([#1516](https://gitea.com/gitea/runner/issues/1516)) ([7679695](https://gitea.com/gitea/runner/commits/767969502aa0cc95e65916981a950250ddb6c71d))
* **deps:** bump github.com/docker/cli ([#1626](https://gitea.com/gitea/runner/issues/1626)) ([98ad62f](https://gitea.com/gitea/runner/commits/98ad62fcefc3b760e69040b25d15d34379376319))
* **deps:** bump github.com/docker/cli ([#1712](https://gitea.com/gitea/runner/issues/1712)) ([24b04df](https://gitea.com/gitea/runner/commits/24b04dfa55ebd31c7b4f2fb114350d14783a533a))
* **deps:** bump github.com/docker/cli ([#1722](https://gitea.com/gitea/runner/issues/1722)) ([f2bd194](https://gitea.com/gitea/runner/commits/f2bd194c7fa282fb35f31bcafe9bf1875da71fc7))
* **deps:** bump github.com/docker/cli ([#1753](https://gitea.com/gitea/runner/issues/1753)) ([ded31bd](https://gitea.com/gitea/runner/commits/ded31bdbc0c82ddb737535fab0fc4d1800ae2b00))
* **deps:** bump github.com/docker/cli ([#1780](https://gitea.com/gitea/runner/issues/1780)) ([3eb6e83](https://gitea.com/gitea/runner/commits/3eb6e83ea4c589a0dfcb2a17e93a04c23909f770))
* **deps:** bump github.com/docker/cli ([#1805](https://gitea.com/gitea/runner/issues/1805)) ([39abbce](https://gitea.com/gitea/runner/commits/39abbce4e6cdf3988a689a404f2f0082ce0125bf))
* **deps:** bump github.com/docker/cli ([#1816](https://gitea.com/gitea/runner/issues/1816)) ([9c1f1f8](https://gitea.com/gitea/runner/commits/9c1f1f8d844148079cf80c49413d7f3357c42025))
* **deps:** bump github.com/docker/cli ([#1843](https://gitea.com/gitea/runner/issues/1843)) ([e48c079](https://gitea.com/gitea/runner/commits/e48c07988ed70d94860c6d6276e134f5f51bd802))
* **deps:** bump github.com/docker/cli ([#1932](https://gitea.com/gitea/runner/issues/1932)) ([8e6c7c1](https://gitea.com/gitea/runner/commits/8e6c7c11fe35bd117579277c29238e4e45996a7c))
* **deps:** bump github.com/docker/cli ([#1999](https://gitea.com/gitea/runner/issues/1999)) ([bd9032d](https://gitea.com/gitea/runner/commits/bd9032de0a79e9c0755596b4a1a10f5de8847639))
* **deps:** bump github.com/docker/cli ([#2069](https://gitea.com/gitea/runner/issues/2069)) ([3d65b0f](https://gitea.com/gitea/runner/commits/3d65b0f73f3c2015f75911a1696e1e8c6d844bd5))
* **deps:** bump github.com/docker/cli ([#854](https://gitea.com/gitea/runner/issues/854)) ([50dd0c0](https://gitea.com/gitea/runner/commits/50dd0c0db149afd3f92d64b79b6eb8d2d0b0ec16))
* **deps:** bump github.com/docker/cli ([#956](https://gitea.com/gitea/runner/issues/956)) ([bbf6b7e](https://gitea.com/gitea/runner/commits/bbf6b7e080959ebb13d3500fed637e822e2b8ec0))
* **deps:** bump github.com/docker/cli from 23.0.0-rc.1+incompatible to 23.0.0+incompatible ([#1611](https://gitea.com/gitea/runner/issues/1611)) ([36dbbc1](https://gitea.com/gitea/runner/commits/36dbbc1dfa049676e17e1f4712c6b918ce4ce5b7))
* **deps:** bump github.com/docker/distribution ([#1037](https://gitea.com/gitea/runner/issues/1037)) ([5c0023c](https://gitea.com/gitea/runner/commits/5c0023c3466b7d5d8a900a1b5a3798e4d1dabe75))
* **deps:** bump github.com/docker/distribution ([#1801](https://gitea.com/gitea/runner/issues/1801)) ([b5fa245](https://gitea.com/gitea/runner/commits/b5fa24537f70e88166c7e43d45b5f6ff41d0b12d))
* **deps:** bump github.com/docker/distribution ([#2037](https://gitea.com/gitea/runner/issues/2037)) ([e7e158c](https://gitea.com/gitea/runner/commits/e7e158cd7eed7c78d7111111aa1d005cadf1df0c))
* **deps:** bump github.com/docker/distribution ([#994](https://gitea.com/gitea/runner/issues/994)) ([0fae967](https://gitea.com/gitea/runner/commits/0fae96792c77b395587c4d8129d0f680bffad643))
* **deps:** bump github.com/docker/docker ([#1062](https://gitea.com/gitea/runner/issues/1062)) ([87175ec](https://gitea.com/gitea/runner/commits/87175ecb6882f0735c1de8eeb884ced21f978873))
* **deps:** bump github.com/docker/docker ([#1080](https://gitea.com/gitea/runner/issues/1080)) ([c8a40bf](https://gitea.com/gitea/runner/commits/c8a40bf8f41a68efe55f28e8b17a9b5f3d004d5e))
* **deps:** bump github.com/docker/docker ([#1155](https://gitea.com/gitea/runner/issues/1155)) ([01d85c9](https://gitea.com/gitea/runner/commits/01d85c9a0150092afa061ba1518a96cae6dcb73b))
* **deps:** bump github.com/docker/docker ([#1167](https://gitea.com/gitea/runner/issues/1167)) ([018c846](https://gitea.com/gitea/runner/commits/018c846d8d92979e158f0618b9f3d068fecf08f1))
* **deps:** bump github.com/docker/docker ([#1212](https://gitea.com/gitea/runner/issues/1212)) ([63df5dd](https://gitea.com/gitea/runner/commits/63df5ddf429881b88d766919db2f404d1dc35c9c))
* **deps:** bump github.com/docker/docker ([#1339](https://gitea.com/gitea/runner/issues/1339)) ([ca8767f](https://gitea.com/gitea/runner/commits/ca8767f08057b176ea7f6261c3f557bde1817a8c))
* **deps:** bump github.com/docker/docker ([#1392](https://gitea.com/gitea/runner/issues/1392)) ([6ff4491](https://gitea.com/gitea/runner/commits/6ff44912381a15fbc27ce62937bae89682cfaa2f))
* **deps:** bump github.com/docker/docker ([#1403](https://gitea.com/gitea/runner/issues/1403)) ([0d1dccc](https://gitea.com/gitea/runner/commits/0d1dccca3a5d663860e57cc53a50e61bbbc003fa))
* **deps:** bump github.com/docker/docker ([#1418](https://gitea.com/gitea/runner/issues/1418)) ([e03b816](https://gitea.com/gitea/runner/commits/e03b816c5880b60af28f7dda33e4361d30451092))
* **deps:** bump github.com/docker/docker ([#1517](https://gitea.com/gitea/runner/issues/1517)) ([469d024](https://gitea.com/gitea/runner/commits/469d024c1fb6f376c56bcf380911761d9ad829db))
* **deps:** bump github.com/docker/docker ([#1573](https://gitea.com/gitea/runner/issues/1573)) ([612e08f](https://gitea.com/gitea/runner/commits/612e08f113b0d5c97a8c1ed664fdb2b19b1a418e))
* **deps:** bump github.com/docker/docker ([#1613](https://gitea.com/gitea/runner/issues/1613)) ([f41e912](https://gitea.com/gitea/runner/commits/f41e9127a8e4516e5d78d58702d00c55de1a6b96))
* **deps:** bump github.com/docker/docker ([#1624](https://gitea.com/gitea/runner/issues/1624)) ([7c3c534](https://gitea.com/gitea/runner/commits/7c3c5349ab54ac3a2c779b199310b8337f6e8fbf))
* **deps:** bump github.com/docker/docker ([#1711](https://gitea.com/gitea/runner/issues/1711)) ([68c72b9](https://gitea.com/gitea/runner/commits/68c72b9a51f4fa586af62619f9e6361e023a5d48))
* **deps:** bump github.com/docker/docker ([#1719](https://gitea.com/gitea/runner/issues/1719)) ([65a925b](https://gitea.com/gitea/runner/commits/65a925b4ed5f0bebab9314433fcfdac4018482ac))
* **deps:** bump github.com/docker/docker ([#1754](https://gitea.com/gitea/runner/issues/1754)) ([16f35f6](https://gitea.com/gitea/runner/commits/16f35f64eb992d96992c43ab1ad0e2fae5fbc8b5))
* **deps:** bump github.com/docker/docker ([#1779](https://gitea.com/gitea/runner/issues/1779)) ([1f9bbe1](https://gitea.com/gitea/runner/commits/1f9bbe12dd5a86ee6401000bcc654b9b66910bae))
* **deps:** bump github.com/docker/docker ([#1806](https://gitea.com/gitea/runner/issues/1806)) ([c00810b](https://gitea.com/gitea/runner/commits/c00810bdf4ca9e5cdf0c34e7ede5d4ebc344f558))
* **deps:** bump github.com/docker/docker ([#1933](https://gitea.com/gitea/runner/issues/1933)) ([d720ff0](https://gitea.com/gitea/runner/commits/d720ff09a281aae453fe6cc2970e6f833b1da49d))
* **deps:** bump github.com/docker/docker ([#2000](https://gitea.com/gitea/runner/issues/2000)) ([c84a3ef](https://gitea.com/gitea/runner/commits/c84a3ef6d0e4c545f627572aef0c0313bc9d4187))
* **deps:** bump github.com/docker/docker ([#2067](https://gitea.com/gitea/runner/issues/2067)) ([84a4025](https://gitea.com/gitea/runner/commits/84a4025bc844eee060ad0d2d71af7b4bd2595a94))
* **deps:** bump github.com/docker/docker ([#2252](https://gitea.com/gitea/runner/issues/2252)) ([069720a](https://gitea.com/gitea/runner/commits/069720abff4e169834dda534da5b0013c03775a3))
* **deps:** bump github.com/docker/docker ([#855](https://gitea.com/gitea/runner/issues/855)) ([6a3a47c](https://gitea.com/gitea/runner/commits/6a3a47c2175c5f6438370d28febf8933097cdc55))
* **deps:** bump github.com/docker/docker ([#957](https://gitea.com/gitea/runner/issues/957)) ([b0f5557](https://gitea.com/gitea/runner/commits/b0f55571f092c47e75f746ecb8aeade718db9495))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.0.0 to 5.3.1 ([#802](https://gitea.com/gitea/runner/issues/802)) ([6d7f574](https://gitea.com/gitea/runner/commits/6d7f574859a1a358db7418dbb96fcd07b983eb31))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.3.1 to 5.4.0 ([#1541](https://gitea.com/gitea/runner/issues/1541)) ([9c70851](https://gitea.com/gitea/runner/commits/9c7085195c4fcfe1f37b53dd9a329824a1753b83))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.0 to 5.4.1 ([#1612](https://gitea.com/gitea/runner/issues/1612)) ([22dc1e0](https://gitea.com/gitea/runner/commits/22dc1e0e7ef5291adbc80a50574ad91a07c64754))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.1 to 5.5.0 ([#2012](https://gitea.com/gitea/runner/issues/2012)) ([d134079](https://gitea.com/gitea/runner/commits/d134079807fe35a2d8ebcfd24b139f48c8123caa))
* **deps:** bump github.com/go-git/go-git/v5 ([#1830](https://gitea.com/gitea/runner/issues/1830)) ([b0996e0](https://gitea.com/gitea/runner/commits/b0996e057750223a14321ff92c1c6d6fbb8e1a10))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.0 to 5.10.1 ([#2114](https://gitea.com/gitea/runner/issues/2114)) ([3c405a0](https://gitea.com/gitea/runner/commits/3c405a0d94a7f0c905b51161be0ac695653e48fc))
* **deps:** bump github.com/go-git/go-git/v5 from 5.10.1 to 5.11.0 ([#2117](https://gitea.com/gitea/runner/issues/2117)) ([85c3b3b](https://gitea.com/gitea/runner/commits/85c3b3b5419bba9ac6434aefbae1407d34f02654))
* **deps:** bump github.com/go-git/go-git/v5 from 5.7.0 to 5.8.0 ([#1925](https://gitea.com/gitea/runner/issues/1925)) ([17bf4fc](https://gitea.com/gitea/runner/commits/17bf4fc5afd147ea967445b9d0dc2e684f57377c))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.0 to 5.8.1 ([#1934](https://gitea.com/gitea/runner/issues/1934)) ([8bf10cf](https://gitea.com/gitea/runner/commits/8bf10cf87618e3ebe1c0f77c1add9e3a08b3f514))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.1 to 5.9.0 ([#2011](https://gitea.com/gitea/runner/issues/2011)) ([b637d79](https://gitea.com/gitea/runner/commits/b637d79ec3a0468e0ff968cdf57420c7020265c8))
* **deps:** bump github.com/go-git/go-git/v5 from 5.9.0 to 5.10.0 ([#2066](https://gitea.com/gitea/runner/issues/2066)) ([854e3e9](https://gitea.com/gitea/runner/commits/854e3e9ec533132b92838aaf6d801c29064716dd))
* **deps:** bump github.com/go-ini/ini from 1.63.2 to 1.64.0 ([#882](https://gitea.com/gitea/runner/issues/882)) ([c26e661](https://gitea.com/gitea/runner/commits/c26e661889059d3afb2a962f86403e29a8b7c213))
* **deps:** bump github.com/go-ini/ini from 1.64.0 to 1.66.4 ([#997](https://gitea.com/gitea/runner/issues/997)) ([7d43396](https://gitea.com/gitea/runner/commits/7d433968e525098b506aae549d5fc5da4fec571e))
* **deps:** bump github.com/go-ini/ini from 1.66.4 to 1.66.5 ([#1189](https://gitea.com/gitea/runner/issues/1189)) ([64387bc](https://gitea.com/gitea/runner/commits/64387bcf7b92f07b42a4c93f7e4504219af5df2c))
* **deps:** bump github.com/go-ini/ini from 1.66.5 to 1.66.6 ([#1201](https://gitea.com/gitea/runner/issues/1201)) ([859445f](https://gitea.com/gitea/runner/commits/859445fb94e12cc600e5a59d9c9cd687214cb801))
* **deps:** bump github.com/go-ini/ini from 1.66.6 to 1.67.0 ([#1305](https://gitea.com/gitea/runner/issues/1305)) ([12029e3](https://gitea.com/gitea/runner/commits/12029e3df5841538a806293881f836f88ec07a2c))
* **deps:** bump github.com/imdario/mergo from 0.3.12 to 0.3.13 ([#1378](https://gitea.com/gitea/runner/issues/1378)) ([9f9f644](https://gitea.com/gitea/runner/commits/9f9f644d36dc8a7ba73701bdd18cf0d23eaa736b))
* **deps:** bump github.com/imdario/mergo from 0.3.13 to 0.3.14 ([#1688](https://gitea.com/gitea/runner/issues/1688)) ([2fa0a5f](https://gitea.com/gitea/runner/commits/2fa0a5f76979689945215fc392f88f836bc71103))
* **deps:** bump github.com/imdario/mergo from 0.3.14 to 0.3.15 ([#1701](https://gitea.com/gitea/runner/issues/1701)) ([f764ecb](https://gitea.com/gitea/runner/commits/f764ecb283febc36da0fa302a8b8aa8d6d350054))
* **deps:** bump github.com/imdario/mergo from 0.3.15 to 0.3.16 ([#1840](https://gitea.com/gitea/runner/issues/1840)) ([0b05ee8](https://gitea.com/gitea/runner/commits/0b05ee8e0b44454f56baaf23f0f00d7e946b8f37))
* **deps:** bump github.com/joho/godotenv from 1.3.0 to 1.4.0 ([#824](https://gitea.com/gitea/runner/issues/824)) ([bfcf0ab](https://gitea.com/gitea/runner/commits/bfcf0abd73e558521367e36b323afad0dc1f1a71))
* **deps:** bump github.com/joho/godotenv from 1.4.0 to 1.5.1 ([#1614](https://gitea.com/gitea/runner/issues/1614)) ([281a52f](https://gitea.com/gitea/runner/commits/281a52f0d6ed70bd31feaedcbbb46e941cc8bdcc))
* **deps:** bump github.com/mattn/go-isatty from 0.0.14 to 0.0.16 ([#1312](https://gitea.com/gitea/runner/issues/1312)) ([b23bbef](https://gitea.com/gitea/runner/commits/b23bbefc365012886192e477a88d38a0909ecba1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.16 to 0.0.17 ([#1539](https://gitea.com/gitea/runner/issues/1539)) ([0d38857](https://gitea.com/gitea/runner/commits/0d388572a89fabdbd229536b7ecca9db4ea88571))
* **deps:** bump github.com/mattn/go-isatty from 0.0.17 to 0.0.18 ([#1702](https://gitea.com/gitea/runner/issues/1702)) ([385d71a](https://gitea.com/gitea/runner/commits/385d71a215b532dc983c640452e64f38426999b1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.18 to 0.0.19 ([#1859](https://gitea.com/gitea/runner/issues/1859)) ([e912ab3](https://gitea.com/gitea/runner/commits/e912ab32c246692f29d9ef520e92f23edd073a4a))
* **deps:** bump github.com/mattn/go-isatty from 0.0.19 to 0.0.20 ([#2059](https://gitea.com/gitea/runner/issues/2059)) ([db71c41](https://gitea.com/gitea/runner/commits/db71c41d17e967bc12ec24c122aad9d314c7d693))
* **deps:** bump github.com/moby/buildkit from 0.10.0 to 0.10.1 ([#1116](https://gitea.com/gitea/runner/issues/1116)) ([bb6ead4](https://gitea.com/gitea/runner/commits/bb6ead459a8e7b46e817369da3565b5ac86767a1))
* **deps:** bump github.com/moby/buildkit from 0.10.1 to 0.10.2 ([#1147](https://gitea.com/gitea/runner/issues/1147)) ([dfca2c5](https://gitea.com/gitea/runner/commits/dfca2c57df978aae5ec8b29a1a80c0fba09b56a8))
* **deps:** bump github.com/moby/buildkit from 0.10.2 to 0.10.3 ([#1153](https://gitea.com/gitea/runner/issues/1153)) ([a76c349](https://gitea.com/gitea/runner/commits/a76c34987233b065109457893a334879a4615baf))
* **deps:** bump github.com/moby/buildkit from 0.10.3 to 0.10.4 ([#1322](https://gitea.com/gitea/runner/issues/1322)) ([8bad6ac](https://gitea.com/gitea/runner/commits/8bad6aca50b7595fcc5253c72a2936fe37157a63))
* **deps:** bump github.com/moby/buildkit from 0.10.4 to 0.10.5 ([#1406](https://gitea.com/gitea/runner/issues/1406)) ([07a1890](https://gitea.com/gitea/runner/commits/07a18902ffa08ef1f209d8b976142805e2afb8f0))
* **deps:** bump github.com/moby/buildkit from 0.10.5 to 0.10.6 ([#1437](https://gitea.com/gitea/runner/issues/1437)) ([64e68bd](https://gitea.com/gitea/runner/commits/64e68bd7f21afdbbe2f40676c2a8c08701057753))
* **deps:** bump github.com/moby/buildkit from 0.10.6 to 0.11.0 ([#1563](https://gitea.com/gitea/runner/issues/1563)) ([d8ba8cb](https://gitea.com/gitea/runner/commits/d8ba8cbf17d816ed2be8a2b8a031c84b3d6fd8e9))
* **deps:** bump github.com/moby/buildkit from 0.11.0 to 0.11.1 ([#1576](https://gitea.com/gitea/runner/issues/1576)) ([09e9f6a](https://gitea.com/gitea/runner/commits/09e9f6a1ac00592b0341d73d6c642dd4aeb83a02))
* **deps:** bump github.com/moby/buildkit from 0.11.1 to 0.11.2 ([#1589](https://gitea.com/gitea/runner/issues/1589)) ([e259dc2](https://gitea.com/gitea/runner/commits/e259dc2835b6b759de7f4dea2daa490438f41bbe))
* **deps:** bump github.com/moby/buildkit from 0.11.2 to 0.11.3 ([#1635](https://gitea.com/gitea/runner/issues/1635)) ([2f55276](https://gitea.com/gitea/runner/commits/2f55276cea863b82aefc806c52773c5ba65774e2))
* **deps:** bump github.com/moby/buildkit from 0.11.3 to 0.11.4 ([#1669](https://gitea.com/gitea/runner/issues/1669)) ([aeee205](https://gitea.com/gitea/runner/commits/aeee2052de45b3ceff3fb66d3eff461379b89d92))
* **deps:** bump github.com/moby/buildkit from 0.11.4 to 0.11.5 ([#1703](https://gitea.com/gitea/runner/issues/1703)) ([351ae99](https://gitea.com/gitea/runner/commits/351ae99bc1b205c61b02d4fb04c06f8287b57f9f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1755](https://gitea.com/gitea/runner/issues/1755)) ([2ea7891](https://gitea.com/gitea/runner/commits/2ea78917879854128154124d3b7c8fa8d552e6e8))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1791](https://gitea.com/gitea/runner/issues/1791)) ([65ef31f](https://gitea.com/gitea/runner/commits/65ef31f102ceb75623973921099454637bab55b0))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1817](https://gitea.com/gitea/runner/issues/1817)) ([16c574c](https://gitea.com/gitea/runner/commits/16c574cd266ffe11ae48a37ccf3f6582ce90de21))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1858](https://gitea.com/gitea/runner/issues/1858)) ([b92d95f](https://gitea.com/gitea/runner/commits/b92d95f89931549a235e5f141ecb0a053f4216d3))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1880](https://gitea.com/gitea/runner/issues/1880)) ([a0307d3](https://gitea.com/gitea/runner/commits/a0307d3b7cd20437000d3931149fe600fe515a0f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1882](https://gitea.com/gitea/runner/issues/1882)) ([3d1d8a9](https://gitea.com/gitea/runner/commits/3d1d8a9acac74363f09790f0efa450a27dbeb1ce))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1884](https://gitea.com/gitea/runner/issues/1884)) ([4c8da85](https://gitea.com/gitea/runner/commits/4c8da8558d0a741bb2cecd6d6edf832dc1a055a9))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.12.5 ([#2188](https://gitea.com/gitea/runner/issues/2188)) ([6dcf9bc](https://gitea.com/gitea/runner/commits/6dcf9bc6e62d093fe8152edced9c70ba686d603d))
* **deps:** bump github.com/moby/buildkit from 0.11.6 to 0.12.0 ([#1917](https://gitea.com/gitea/runner/issues/1917)) ([4810f69](https://gitea.com/gitea/runner/commits/4810f69367a141f9cc282d21ab61bb598d5ec584))
* **deps:** bump github.com/moby/buildkit from 0.12.0 to 0.12.1 ([#1947](https://gitea.com/gitea/runner/issues/1947)) ([c812715](https://gitea.com/gitea/runner/commits/c8127155bc64224a1392d2264b8d432328edea57))
* **deps:** bump github.com/moby/buildkit from 0.12.1 to 0.12.2 ([#1986](https://gitea.com/gitea/runner/issues/1986)) ([0c12273](https://gitea.com/gitea/runner/commits/0c12273ebaa937c52dea06d91b8985304874ab14))
* **deps:** bump github.com/moby/buildkit from 0.12.2 to 0.12.3 ([#2060](https://gitea.com/gitea/runner/issues/2060)) ([db6e477](https://gitea.com/gitea/runner/commits/db6e477e25808fcab235fa28b941befdd83c5c03))
* **deps:** bump github.com/moby/buildkit from 0.12.3 to 0.12.4 ([#2113](https://gitea.com/gitea/runner/issues/2113)) ([899a1f2](https://gitea.com/gitea/runner/commits/899a1f206ee6c8636b14a0bc7dc823602365ed9c))
* **deps:** bump github.com/moby/buildkit from 0.12.4 to 0.12.5 ([#2186](https://gitea.com/gitea/runner/issues/2186)) ([46dc2ff](https://gitea.com/gitea/runner/commits/46dc2ffe803f783c3b9fc2190c6566a81d8645dc))
* **deps:** bump github.com/moby/buildkit from 0.9.0 to 0.9.1 ([#838](https://gitea.com/gitea/runner/issues/838)) ([b5e8a18](https://gitea.com/gitea/runner/commits/b5e8a18683c77c7035c704bb7cec980401c73e18))
* **deps:** bump github.com/moby/buildkit from 0.9.1 to 0.9.2 ([#864](https://gitea.com/gitea/runner/issues/864)) ([f8b3563](https://gitea.com/gitea/runner/commits/f8b35631028c7705eb90b1d97300dc19c739cc72))
* **deps:** bump github.com/moby/buildkit from 0.9.2 to 0.9.3 ([#958](https://gitea.com/gitea/runner/issues/958)) ([dba7f83](https://gitea.com/gitea/runner/commits/dba7f8379f327232a1726b336d2cf4a6c3e80ccc))
* **deps:** bump github.com/moby/buildkit from 0.9.3 to 0.10.0 ([#1036](https://gitea.com/gitea/runner/issues/1036)) ([56486f1](https://gitea.com/gitea/runner/commits/56486f14b4a094542d2cf4a903676cd9f14220ec))
* **deps:** bump github.com/moby/patternmatcher from 0.5.0 to 0.6.0 ([#1985](https://gitea.com/gitea/runner/issues/1985)) ([323bee9](https://gitea.com/gitea/runner/commits/323bee9ab5a73473236432e0b89e4d22100cd24d))
* **deps:** bump github.com/opencontainers/image-spec ([#1649](https://gitea.com/gitea/runner/issues/1649)) ([04d12b0](https://gitea.com/gitea/runner/commits/04d12b0206a62cc48a765ed29f2a1ad732d33b3b))
* **deps:** bump github.com/opencontainers/image-spec ([#1781](https://gitea.com/gitea/runner/issues/1781)) ([15eb1fa](https://gitea.com/gitea/runner/commits/15eb1fa92a494f0abde33ad1bfe63496195db9ec))
* **deps:** bump github.com/opencontainers/image-spec ([#1807](https://gitea.com/gitea/runner/issues/1807)) ([1eacf23](https://gitea.com/gitea/runner/commits/1eacf23dcbbea2e5a7d17e8c325ce94ee0b1d3d1))
* **deps:** bump github.com/opencontainers/image-spec ([#1829](https://gitea.com/gitea/runner/issues/1829)) ([26f1b1f](https://gitea.com/gitea/runner/commits/26f1b1f4b6df887a08af0ed6e692962989f5a57c))
* **deps:** bump github.com/opencontainers/image-spec ([#1871](https://gitea.com/gitea/runner/issues/1871)) ([ef60885](https://gitea.com/gitea/runner/commits/ef608854d0b42704a3316c1b09c74c4fff5c40ae))
* **deps:** bump github.com/opencontainers/image-spec ([#1881](https://gitea.com/gitea/runner/issues/1881)) ([f9dcb58](https://gitea.com/gitea/runner/commits/f9dcb58db23061c016bba6712080275219f7008c))
* **deps:** bump github.com/opencontainers/image-spec ([#1883](https://gitea.com/gitea/runner/issues/1883)) ([70e1e37](https://gitea.com/gitea/runner/commits/70e1e372804054374c222c7de1e066bf7ab92b31))
* **deps:** bump github.com/opencontainers/image-spec ([#1893](https://gitea.com/gitea/runner/issues/1893)) ([f40d0b8](https://gitea.com/gitea/runner/commits/f40d0b873da435b6a1dc9fbb5391ce31b2d2dc0e))
* **deps:** bump github.com/opencontainers/image-spec ([#2010](https://gitea.com/gitea/runner/issues/2010)) ([83af8f8](https://gitea.com/gitea/runner/commits/83af8f8767091542dcc4263ebea06e6cc4edd7ce))
* **deps:** bump github.com/opencontainers/image-spec ([#2167](https://gitea.com/gitea/runner/issues/2167)) ([36e0261](https://gitea.com/gitea/runner/commits/36e0261150e54c2ebe727d526a9c82729d6023f7))
* **deps:** bump github.com/opencontainers/image-spec ([#2214](https://gitea.com/gitea/runner/issues/2214)) ([27eb79b](https://gitea.com/gitea/runner/commits/27eb79bb5a6bddb67393432371eed24a4fe242fe))
* **deps:** bump github.com/opencontainers/image-spec ([#2223](https://gitea.com/gitea/runner/issues/2223)) ([9e6ce12](https://gitea.com/gitea/runner/commits/9e6ce121bc1f41e1ecb2c3db7b35f79af70d4ae2))
* **deps:** bump github.com/opencontainers/image-spec ([#954](https://gitea.com/gitea/runner/issues/954)) ([eedac17](https://gitea.com/gitea/runner/commits/eedac179c9f314ee38eda7482c45664df276e546))
* **deps:** bump github.com/opencontainers/runc from 1.1.3 to 1.1.5 ([#1708](https://gitea.com/gitea/runner/issues/1708)) ([6745999](https://gitea.com/gitea/runner/commits/6745999c105a1bb4429275103dc89b6a704b421c))
* **deps:** bump github.com/opencontainers/runc from 1.1.7 to 1.1.12 ([#2187](https://gitea.com/gitea/runner/issues/2187)) ([0dbf44c](https://gitea.com/gitea/runner/commits/0dbf44c657d28b6bdcb7b70c1da59bc925a658db))
* **deps:** bump github.com/opencontainers/selinux ([#1102](https://gitea.com/gitea/runner/issues/1102)) ([e0044ba](https://gitea.com/gitea/runner/commits/e0044ba913149be76acd64ca7d1a66f27532f4c7))
* **deps:** bump github.com/opencontainers/selinux ([#1372](https://gitea.com/gitea/runner/issues/1372)) ([d587ce0](https://gitea.com/gitea/runner/commits/d587ce08e9b04738fd1d3a2433e1f8b4d44c4aed))
* **deps:** bump github.com/opencontainers/selinux ([#1625](https://gitea.com/gitea/runner/issues/1625)) ([8048baf](https://gitea.com/gitea/runner/commits/8048baf435665cb3293f593a06526dfff1691b54))
* **deps:** bump github.com/rhysd/actionlint from 1.6.10 to 1.6.11 ([#1106](https://gitea.com/gitea/runner/issues/1106)) ([065d630](https://gitea.com/gitea/runner/commits/065d630af01b4ed2a0c797d54933ea3b65c3cb38))
* **deps:** bump github.com/rhysd/actionlint from 1.6.11 to 1.6.12 ([#1126](https://gitea.com/gitea/runner/issues/1126)) ([de1dea9](https://gitea.com/gitea/runner/commits/de1dea93795ebf78e0ca531c35969f125035e8c4))
* **deps:** bump github.com/rhysd/actionlint from 1.6.12 to 1.6.13 ([#1173](https://gitea.com/gitea/runner/issues/1173)) ([166d063](https://gitea.com/gitea/runner/commits/166d0630593b40a288cc0190d4aa3e20b31d1461))
* **deps:** bump github.com/rhysd/actionlint from 1.6.13 to 1.6.14 ([#1234](https://gitea.com/gitea/runner/issues/1234)) ([9499612](https://gitea.com/gitea/runner/commits/9499612b5818774c5c6d73fbf2596ed464b7b01d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.14 to 1.6.15 ([#1239](https://gitea.com/gitea/runner/issues/1239)) ([4d9d6ec](https://gitea.com/gitea/runner/commits/4d9d6ecc924bef6c5976acd7f224558e21fa5b5f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.15 to 1.6.16 ([#1311](https://gitea.com/gitea/runner/issues/1311)) ([cf9d82f](https://gitea.com/gitea/runner/commits/cf9d82f69daab9dc907b893a38307a17f444c529))
* **deps:** bump github.com/rhysd/actionlint from 1.6.16 to 1.6.17 ([#1321](https://gitea.com/gitea/runner/issues/1321)) ([9b14698](https://gitea.com/gitea/runner/commits/9b146980dfabd14e9b9defea8f894ad8ab24132a))
* **deps:** bump github.com/rhysd/actionlint from 1.6.17 to 1.6.18 ([#1352](https://gitea.com/gitea/runner/issues/1352)) ([053bff1](https://gitea.com/gitea/runner/commits/053bff19a710669308da5e356ab028bab2a70b5d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.18 to 1.6.19 ([#1360](https://gitea.com/gitea/runner/issues/1360)) ([3fccb5a](https://gitea.com/gitea/runner/commits/3fccb5ae3467fbd013ca72ed08732d2916beb1c2))
* **deps:** bump github.com/rhysd/actionlint from 1.6.19 to 1.6.20 ([#1371](https://gitea.com/gitea/runner/issues/1371)) ([bec45d0](https://gitea.com/gitea/runner/commits/bec45d0dc46e2ccae81c521f3db313b04b0c0831))
* **deps:** bump github.com/rhysd/actionlint from 1.6.20 to 1.6.21 ([#1380](https://gitea.com/gitea/runner/issues/1380)) ([c4be1cf](https://gitea.com/gitea/runner/commits/c4be1cfd6d8ab6664c3a3ffd7c4e5e253fa0b77f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.21 to 1.6.22 ([#1427](https://gitea.com/gitea/runner/issues/1427)) ([8036f49](https://gitea.com/gitea/runner/commits/8036f49fb220bb2315ecba9707242cc0a8533aba))
* **deps:** bump github.com/rhysd/actionlint from 1.6.22 to 1.6.23 ([#1574](https://gitea.com/gitea/runner/issues/1574)) ([1b7212f](https://gitea.com/gitea/runner/commits/1b7212f5fd1e85d97a7119ab098941bd4cf5255b))
* **deps:** bump github.com/rhysd/actionlint from 1.6.23 to 1.6.24 ([#1721](https://gitea.com/gitea/runner/issues/1721)) ([5a331d2](https://gitea.com/gitea/runner/commits/5a331d2d9935565eabdf5a09226b7ba59a029133))
* **deps:** bump github.com/rhysd/actionlint from 1.6.24 to 1.6.25 ([#1870](https://gitea.com/gitea/runner/issues/1870)) ([4fc176f](https://gitea.com/gitea/runner/commits/4fc176f556dcf41e69854570ce65680f127a9357))
* **deps:** bump github.com/rhysd/actionlint from 1.6.25 to 1.6.26 ([#2026](https://gitea.com/gitea/runner/issues/2026)) ([2be4def](https://gitea.com/gitea/runner/commits/2be4def7be0deda86b7e0a00aa06be5e87ac599c))
* **deps:** bump github.com/rhysd/actionlint from 1.6.26 to 1.6.27 ([#2228](https://gitea.com/gitea/runner/issues/2228)) ([68597d6](https://gitea.com/gitea/runner/commits/68597d68f6cc3c384bfd458ab4ec6d51f63a23f6))
* **deps:** bump github.com/rhysd/actionlint from 1.6.8 to 1.6.9 ([#1018](https://gitea.com/gitea/runner/issues/1018)) ([9a54c99](https://gitea.com/gitea/runner/commits/9a54c99ac76a2667551bfd85c5aabcabc969e9dd))
* **deps:** bump github.com/rhysd/actionlint from 1.6.9 to 1.6.10 ([#1039](https://gitea.com/gitea/runner/issues/1039)) ([ee92e07](https://gitea.com/gitea/runner/commits/ee92e07daa0b82e4758807f424f79194e5cb857b))
* **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1273](https://gitea.com/gitea/runner/issues/1273)) ([17dd54d](https://gitea.com/gitea/runner/commits/17dd54d69231884ea561fe2bd0c0e059debd2658))
* **deps:** bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2 ([#1814](https://gitea.com/gitea/runner/issues/1814)) ([06054d2](https://gitea.com/gitea/runner/commits/06054d2931861d2f02cd4028d04164a4d52bbc1c))
* **deps:** bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3 ([#1841](https://gitea.com/gitea/runner/issues/1841)) ([e150310](https://gitea.com/gitea/runner/commits/e150310da77918e9ca904363f6cd80907877ff12))
* **deps:** bump github.com/spf13/cobra from 1.2.1 to 1.3.0 ([#955](https://gitea.com/gitea/runner/issues/955)) ([e9e6dda](https://gitea.com/gitea/runner/commits/e9e6ddadf567beb62ddb94b3fe472297081d4077))
* **deps:** bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#1040](https://gitea.com/gitea/runner/issues/1040)) ([fa71f0b](https://gitea.com/gitea/runner/commits/fa71f0b1e04a6391b77ef091c0db4d1a8f0ad1c1))
* **deps:** bump github.com/spf13/cobra from 1.4.0 to 1.5.0 ([#1233](https://gitea.com/gitea/runner/issues/1233)) ([4deb213](https://gitea.com/gitea/runner/commits/4deb21344d40dd4c029d5c1048a7f1d557ab8774))
* **deps:** bump github.com/spf13/cobra from 1.5.0 to 1.6.0 ([#1394](https://gitea.com/gitea/runner/issues/1394)) ([52c6edf](https://gitea.com/gitea/runner/commits/52c6edf0a48fa8287eaaa1ea28d89f65d05e8e7b))
* **deps:** bump github.com/spf13/cobra from 1.6.0 to 1.6.1 ([#1420](https://gitea.com/gitea/runner/issues/1420)) ([a9298d9](https://gitea.com/gitea/runner/commits/a9298d936c5265461dd8f6322961232b08fc6d12))
* **deps:** bump github.com/spf13/cobra from 1.6.1 to 1.7.0 ([#1724](https://gitea.com/gitea/runner/issues/1724)) ([c5ce502](https://gitea.com/gitea/runner/commits/c5ce502f9f9e98ea76a09858075e3ee91da6bd10))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([#2079](https://gitea.com/gitea/runner/issues/2079)) ([55b09a0](https://gitea.com/gitea/runner/commits/55b09a04cd14006b34eef02f0d19120759ee167b))
* **deps:** bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#1061](https://gitea.com/gitea/runner/issues/1061)) ([7d946c4](https://gitea.com/gitea/runner/commits/7d946c43e208aa33a0f6183a2b84a8e9f33425ae))
* **deps:** bump github.com/stretchr/testify from 1.7.1 to 1.7.2 ([#1213](https://gitea.com/gitea/runner/issues/1213)) ([3415347](https://gitea.com/gitea/runner/commits/3415347efcdf6d6d7544c5df7cb33e6db043ff4a))
* **deps:** bump github.com/stretchr/testify from 1.7.2 to 1.7.5 ([#1232](https://gitea.com/gitea/runner/issues/1232)) ([c67abf2](https://gitea.com/gitea/runner/commits/c67abf2401ba57097e862b524b86de7171f1c609))
* **deps:** bump github.com/stretchr/testify from 1.7.5 to 1.8.0 ([#1240](https://gitea.com/gitea/runner/issues/1240)) ([50f0b0e](https://gitea.com/gitea/runner/commits/50f0b0e7f49d368583a30e515fb16db58c357bbc))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1404](https://gitea.com/gitea/runner/issues/1404)) ([db850e8](https://gitea.com/gitea/runner/commits/db850e818c2dd36f0021ac0a227ed99d2daafff2))
* **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#1650](https://gitea.com/gitea/runner/issues/1650)) ([a6c95ef](https://gitea.com/gitea/runner/commits/a6c95ef2a771afdc9f83f6750b642491e52f681c))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://gitea.com/gitea/runner/issues/1842)) ([c70a674](https://gitea.com/gitea/runner/commits/c70a6743f62af2c223f62c2fddd6a1a68271a50e))
* **deps:** bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#2235](https://gitea.com/gitea/runner/issues/2235)) ([934b13a](https://gitea.com/gitea/runner/commits/934b13a7a126198019800e79fb2ae7fe57a1c2bd))
* **deps:** bump go.etcd.io/bbolt from 1.3.7 to 1.3.8 ([#2065](https://gitea.com/gitea/runner/issues/2065)) ([83bfbcd](https://gitea.com/gitea/runner/commits/83bfbcdafd55c5dc6989f6a6b0aa07cab5b99ea8))
* **deps:** bump go.etcd.io/bbolt from 1.3.8 to 1.3.9 ([#2229](https://gitea.com/gitea/runner/issues/2229)) ([c79f59f](https://gitea.com/gitea/runner/commits/c79f59f802673f00911bea93db15b83f5bf3507b))
* **deps:** bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#2134](https://gitea.com/gitea/runner/issues/2134)) ([238a495](https://gitea.com/gitea/runner/commits/238a495579abc12665c1fbf01e8a4f16c2c96fe1))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([#2045](https://gitea.com/gitea/runner/issues/2045)) ([99067a9](https://gitea.com/gitea/runner/commits/99067a9c1ed3b9cdb75797a7e87e6ef0471105f8))
* **deps:** bump golang.org/x/net from 0.4.0 to 0.7.0 ([#1636](https://gitea.com/gitea/runner/issues/1636)) ([b7a9eb9](https://gitea.com/gitea/runner/commits/b7a9eb9fbf86da76cbad35296ff11b64ee416d18))
* **deps:** bump golang.org/x/term from 0.10.0 to 0.11.0 ([#1948](https://gitea.com/gitea/runner/issues/1948)) ([74b74e8](https://gitea.com/gitea/runner/commits/74b74e847bad3c71ad7a1df76384a804bb46d644))
* **deps:** bump golang.org/x/term from 0.11.0 to 0.12.0 ([#1992](https://gitea.com/gitea/runner/issues/1992)) ([1d32507](https://gitea.com/gitea/runner/commits/1d32507b52c8a3c502415d1350aeca0eb9dcfe81))
* **deps:** bump golang.org/x/term from 0.12.0 to 0.13.0 ([#2036](https://gitea.com/gitea/runner/issues/2036)) ([3c730d7](https://gitea.com/gitea/runner/commits/3c730d7924ed9da3ce318bba982873b4484106a7))
* **deps:** bump golang.org/x/term from 0.13.0 to 0.14.0 ([#2091](https://gitea.com/gitea/runner/issues/2091)) ([4dcb9b7](https://gitea.com/gitea/runner/commits/4dcb9b7a139f7399b2778575e61a398aefb41d85))
* **deps:** bump golang.org/x/term from 0.14.0 to 0.15.0 ([#2112](https://gitea.com/gitea/runner/issues/2112)) ([2b47c99](https://gitea.com/gitea/runner/commits/2b47c99bb7b5d73c852ca7241cda3e1e70ba3670))
* **deps:** bump golang.org/x/term from 0.15.0 to 0.16.0 ([#2148](https://gitea.com/gitea/runner/issues/2148)) ([ef5746b](https://gitea.com/gitea/runner/commits/ef5746ba747a4a39770769725a62b519d8d919aa))
* **deps:** bump golang.org/x/term from 0.16.0 to 0.17.0 ([#2213](https://gitea.com/gitea/runner/issues/2213)) ([8277584](https://gitea.com/gitea/runner/commits/8277584f00a0fa82598fc795f35b165e0bfa3bdf))
* **deps:** bump golang.org/x/term from 0.17.0 to 0.18.0 ([#2244](https://gitea.com/gitea/runner/issues/2244)) ([8c83d57](https://gitea.com/gitea/runner/commits/8c83d5721258537b3a47dfc060be67d180ba6a0e))
* **deps:** bump golang.org/x/term from 0.3.0 to 0.4.0 ([#1577](https://gitea.com/gitea/runner/issues/1577)) ([a2856ae](https://gitea.com/gitea/runner/commits/a2856aed75e6d4e5607dc59ff074c9612553360f))
* **deps:** bump golang.org/x/term from 0.5.0 to 0.6.0 ([#1666](https://gitea.com/gitea/runner/issues/1666)) ([8dbd151](https://gitea.com/gitea/runner/commits/8dbd151fa71031c737d3f35580906f4028c355ca))
* **deps:** bump golang.org/x/term from 0.6.0 to 0.7.0 ([#1723](https://gitea.com/gitea/runner/issues/1723)) ([3ae0a9d](https://gitea.com/gitea/runner/commits/3ae0a9dfb7d9e4eb2adb3dc4a778290784fc5ce8))
* **deps:** bump golang.org/x/term from 0.8.0 to 0.9.0 ([#1869](https://gitea.com/gitea/runner/issues/1869)) ([b2191ae](https://gitea.com/gitea/runner/commits/b2191ae204a713a5e9bd4b3e1ec34b23b7fd34a6))
* **deps:** bump golang.org/x/term from 0.9.0 to 0.10.0 ([#1906](https://gitea.com/gitea/runner/issues/1906)) ([e343ea9](https://gitea.com/gitea/runner/commits/e343ea9d5f6fa478a364b1f502b4859eed63fa38))
* **deps:** bump golangci/golangci-lint-action from 2 to 3.1.0 ([#1019](https://gitea.com/gitea/runner/issues/1019)) ([7d403b8](https://gitea.com/gitea/runner/commits/7d403b83075932a18007210a62c6d0ce64df11a6))
* **deps:** bump golangci/golangci-lint-action from 3.1.0 to 3.2.0 ([#1169](https://gitea.com/gitea/runner/issues/1169)) ([8994773](https://gitea.com/gitea/runner/commits/89947735d508751b9e84d653635dec91299e9eec))
* **deps:** bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 ([#1408](https://gitea.com/gitea/runner/issues/1408)) ([925ac20](https://gitea.com/gitea/runner/commits/925ac2098b569f95beca47d5aac7e2287922963b))
* **deps:** bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 ([#1436](https://gitea.com/gitea/runner/issues/1436)) ([1b554ae](https://gitea.com/gitea/runner/commits/1b554aeff651cf801c11277bacb162ba31483b08))
* **deps:** bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 ([#1590](https://gitea.com/gitea/runner/issues/1590)) ([70ae63b](https://gitea.com/gitea/runner/commits/70ae63bd10eb4ac3a9a64633b88242f9a10122e3))
* **deps:** bump golangci/golangci-lint-action from 3.4.0 to 3.5.0 ([#1839](https://gitea.com/gitea/runner/issues/1839)) ([34db861](https://gitea.com/gitea/runner/commits/34db86138ba3f375c3797dfe5b2952f32397f7b1))
* **deps:** bump golangci/golangci-lint-action from 3.5.0 to 3.6.0 ([#1868](https://gitea.com/gitea/runner/issues/1868)) ([79f93be](https://gitea.com/gitea/runner/commits/79f93beef21650b25ab20a24aa91032e3a6b597e))
* **deps:** bump golangci/golangci-lint-action from 3.6.0 to 3.7.0 ([#1978](https://gitea.com/gitea/runner/issues/1978)) ([f64c267](https://gitea.com/gitea/runner/commits/f64c267dacd026b306da8ee405e84df56d124e69))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1179](https://gitea.com/gitea/runner/issues/1179)) ([5054455](https://gitea.com/gitea/runner/commits/50544556f8672c5553fd5241f4adb157698c781b))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1515](https://gitea.com/gitea/runner/issues/1515)) ([c044035](https://gitea.com/gitea/runner/commits/c044035d9ef358482efc996fba0a4c3f7f5f3847))
* **deps:** bump goreleaser/goreleaser-action from 4 to 5 ([#2008](https://gitea.com/gitea/runner/issues/2008)) ([2911b21](https://gitea.com/gitea/runner/commits/2911b2172cee57abbfdafd8ffd04fe4b6ccdd9dc))
* **deps:** bump gotest.tools/v3 from 3.0.3 to 3.4.0 ([#1379](https://gitea.com/gitea/runner/issues/1379)) ([c7ccb0d](https://gitea.com/gitea/runner/commits/c7ccb0dde9ab6332af7940e6f61e3ab73685fd07))
* **deps:** bump gotest.tools/v3 from 3.4.0 to 3.5.0 ([#1892](https://gitea.com/gitea/runner/issues/1892)) ([70956f2](https://gitea.com/gitea/runner/commits/70956f2929eddc7a65ffc16ac335db9472924d66))
* **deps:** bump gotest.tools/v3 from 3.5.0 to 3.5.1 ([#2009](https://gitea.com/gitea/runner/issues/2009)) ([60060a7](https://gitea.com/gitea/runner/commits/60060a7c9c7730d653bd8d980f1a17abd068e739))
* **deps:** bump megalinter/megalinter from 5 to 6.11.0 ([#1370](https://gitea.com/gitea/runner/issues/1370)) ([79384c3](https://gitea.com/gitea/runner/commits/79384c35ab42715144388b7d3e0b0290fb93a962))
* **deps:** bump megalinter/megalinter from 6.11.0 to 6.11.1 ([#1381](https://gitea.com/gitea/runner/issues/1381)) ([4541139](https://gitea.com/gitea/runner/commits/4541139109f2838bcac25ddd3d435a5c0661bf25))
* **deps:** bump megalinter/megalinter from 6.11.1 to 6.12.0 ([#1395](https://gitea.com/gitea/runner/issues/1395)) ([37f5b7f](https://gitea.com/gitea/runner/commits/37f5b7f4a2bc363be670ee65f623c6b54ec9d8c5))
* **deps:** bump megalinter/megalinter from 6.12.0 to 6.13.0 ([#1407](https://gitea.com/gitea/runner/issues/1407)) ([e81dc0d](https://gitea.com/gitea/runner/commits/e81dc0d295ad58d13927e61c91d3c4b4a6def942))
* **deps:** bump megalinter/megalinter from 6.13.0 to 6.14.0 ([#1426](https://gitea.com/gitea/runner/issues/1426)) ([c232981](https://gitea.com/gitea/runner/commits/c2329815b8bbc739043e1584e465322a602475e9))
* **deps:** bump megalinter/megalinter from 6.14.0 to 6.15.0 ([#1475](https://gitea.com/gitea/runner/issues/1475)) ([7754ba7](https://gitea.com/gitea/runner/commits/7754ba7fcf3f80b76bb4d7cf1edbe5935c8a6bdc))
* **deps:** bump megalinter/megalinter from 6.15.0 to 6.16.0 ([#1534](https://gitea.com/gitea/runner/issues/1534)) ([e3a722e](https://gitea.com/gitea/runner/commits/e3a722e7c6c697a1b0a38aa314bfb227da4deb5f))
* **deps:** bump megalinter/megalinter from 6.16.0 to 6.17.0 ([#1540](https://gitea.com/gitea/runner/issues/1540)) ([cfc1b91](https://gitea.com/gitea/runner/commits/cfc1b91aa11220f102470622cacd105cedd931c9))
* **deps:** bump megalinter/megalinter from 6.17.0 to 6.18.0 ([#1550](https://gitea.com/gitea/runner/issues/1550)) ([f0c6fa1](https://gitea.com/gitea/runner/commits/f0c6fa11be88093fbcb7ecb1cbebfee8874008cc))
* **deps:** bump megalinter/megalinter from 6.18.0 to 6.19.0 ([#1610](https://gitea.com/gitea/runner/issues/1610)) ([460c78d](https://gitea.com/gitea/runner/commits/460c78da0773c226162d93623e8932715a22e443))
* **deps:** bump megalinter/megalinter from 6.19.0 to 6.20.0 ([#1665](https://gitea.com/gitea/runner/issues/1665)) ([7ee2138](https://gitea.com/gitea/runner/commits/7ee2138418822460ed576391201db55803bd7069))
* **deps:** bump megalinter/megalinter from 6.20.0 to 6.20.1 ([#1679](https://gitea.com/gitea/runner/issues/1679)) ([5a2112a](https://gitea.com/gitea/runner/commits/5a2112a7f8e282fcb75f1c73caa1582811ff4890))
* **deps:** bump megalinter/megalinter from 6.20.1 to 6.21.0 ([#1699](https://gitea.com/gitea/runner/issues/1699)) ([775a128](https://gitea.com/gitea/runner/commits/775a128cd4660be4f98da019c52a2ab9a5080433))
* **deps:** bump megalinter/megalinter from 6.21.0 to 6.22.1 ([#1710](https://gitea.com/gitea/runner/issues/1710)) ([ad7a8a3](https://gitea.com/gitea/runner/commits/ad7a8a3559d937d79e9425b8d29714752879b823))
* **deps:** bump megalinter/megalinter from 6.22.1 to 6.22.2 ([#1720](https://gitea.com/gitea/runner/issues/1720)) ([4ab812c](https://gitea.com/gitea/runner/commits/4ab812c1b10bc16a54dd208963d6b60e2835d328))
* **deps:** bump megalinter/megalinter from 6.22.2 to 7.0.2 ([#1827](https://gitea.com/gitea/runner/issues/1827)) ([481999f](https://gitea.com/gitea/runner/commits/481999f59da465d2d8743fd3d97e24bb6e1a2311))
* **deps:** bump megalinter/megalinter from 7.0.2 to 7.0.4 ([#1838](https://gitea.com/gitea/runner/issues/1838)) ([04b5c73](https://gitea.com/gitea/runner/commits/04b5c739d2649b67565dd7addbe2bf5fc910e2a6))
* **deps:** bump megalinter/megalinter from 7.0.4 to 7.1.0 ([#1857](https://gitea.com/gitea/runner/issues/1857)) ([3555d65](https://gitea.com/gitea/runner/commits/3555d65d08bc9c66a6e4541b2064a082af1b6b3e))
* **deps:** bump megalinter/megalinter from 7.1.0 to 7.2.1 ([#1931](https://gitea.com/gitea/runner/issues/1931)) ([3939f48](https://gitea.com/gitea/runner/commits/3939f48e6def83e8a6ad7358a848ce11912c3798))
* **deps:** bump megalinter/megalinter from 7.2.1 to 7.3.0 ([#1965](https://gitea.com/gitea/runner/issues/1965)) ([2a0a0a1](https://gitea.com/gitea/runner/commits/2a0a0a1a62b8c72cd891208f085c00cb00402b57))
* **deps:** bump megalinter/megalinter from 7.3.0 to 7.4.0 ([#2025](https://gitea.com/gitea/runner/issues/2025)) ([3d47885](https://gitea.com/gitea/runner/commits/3d47885894b804abaf5cebde36d9b9f260ae2c50))
* **deps:** bump megalinter/megalinter from 7.4.0 to 7.5.0 ([#2070](https://gitea.com/gitea/runner/issues/2070)) ([3e5c629](https://gitea.com/gitea/runner/commits/3e5c62977fac53931267a86493e8c5c292508e30))
* **deps:** bump megalinter/megalinter from 7.5.0 to 7.6.0 ([#2098](https://gitea.com/gitea/runner/issues/2098)) ([bd10c9a](https://gitea.com/gitea/runner/commits/bd10c9a80191e7f3c1eed0ee42847f12bc4a81ca))
* **deps:** bump megalinter/megalinter from 7.6.0 to 7.7.0 ([#2119](https://gitea.com/gitea/runner/issues/2119)) ([29f4123](https://gitea.com/gitea/runner/commits/29f4123b5ced89ad39f9f80889576810cfdecf7f))
* **deps:** bump megalinter/megalinter from 7.7.0 to 7.8.0 ([#2164](https://gitea.com/gitea/runner/issues/2164)) ([c215e08](https://gitea.com/gitea/runner/commits/c215e0888aa29780903c113e35c46f495d4157fc))
* **deps:** update go-git to ce62f3e9ff86270538a514a68d3bd5563a733e3b ([#1725](https://gitea.com/gitea/runner/issues/1725)) ([148a545](https://gitea.com/gitea/runner/commits/148a54502167880f90dbf2f68ff44478679f41dd))

### Continuous Integration

* Add snapshot artifacts to PR builds ([#574](https://gitea.com/gitea/runner/issues/574)) ([da07324](https://gitea.com/gitea/runner/commits/da07324779854fda8bff4bd45452a04ce5e68a1e))
* add Winget Releaser workflow ([#1623](https://gitea.com/gitea/runner/issues/1623)) ([34ab815](https://gitea.com/gitea/runner/commits/34ab8150bf8b438964b7e52796f2123466e0cf76))
* always grab latest go version ([#1124](https://gitea.com/gitea/runner/issues/1124)) ([25ebaf4](https://gitea.com/gitea/runner/commits/25ebaf4e7b507ffc7cabdb2a560348af5bec4ac3))
* automatic merge PRs created by a maintainer and approved by 1 other maintainer ([#2156](https://gitea.com/gitea/runner/issues/2156)) ([cd40f3f](https://gitea.com/gitea/runner/commits/cd40f3fe9bd86f7c2c5496686e119d047e9c04a9))
* bump go to 1.18 ([#1087](https://gitea.com/gitea/runner/issues/1087)) ([11c3fb3](https://gitea.com/gitea/runner/commits/11c3fb39dd5d3907baf8f13ef6953d0a21321cc2))
* change winget job runner to `ubuntu-latest` ([#1959](https://gitea.com/gitea/runner/issues/1959)) ([f55ae1a](https://gitea.com/gitea/runner/commits/f55ae1a0bc4fb14d463b79a62ff9a37d7e5bda68))
* **choco:** fix chocolatey docker build ([#903](https://gitea.com/gitea/runner/issues/903)) ([6ebcac3](https://gitea.com/gitea/runner/commits/6ebcac37710ae5d4eb87f63f31d7e1be8fe3d217))
* **choco:** update chocolatey to 1.1.0 ([#1164](https://gitea.com/gitea/runner/issues/1164)) ([507ae61](https://gitea.com/gitea/runner/commits/507ae61d1b3ac3f66273a00fad1321af8c6df260))
* deduplicate running workflows ([#1751](https://gitea.com/gitea/runner/issues/1751)) ([7c6237d](https://gitea.com/gitea/runner/commits/7c6237d93f01b2c50a98a488e40ab7f9a719da35))
* **golangci-lint:** switch to latest ([#870](https://gitea.com/gitea/runner/issues/870)) ([cd24df5](https://gitea.com/gitea/runner/commits/cd24df5727c54de6671f83b44b7527f3c74f122f))
* make stalebot less annoying ([#1587](https://gitea.com/gitea/runner/issues/1587)) ([01952d0](https://gitea.com/gitea/runner/commits/01952d05c858c70ec0f007d90e394cf4cdfcd83b))
* **Mergify:** configuration update ([#924](https://gitea.com/gitea/runner/issues/924)) ([cad4bc8](https://gitea.com/gitea/runner/commits/cad4bc8c367fd4e623dfe2820bd5f3daca80bf58))
* **Mergify:** configuration update ([#966](https://gitea.com/gitea/runner/issues/966)) ([557dc75](https://gitea.com/gitea/runner/commits/557dc755b8fe6b286d15bd7a1091f3731dc44791))
* **Mergify:** configuration update ([#972](https://gitea.com/gitea/runner/issues/972)) ([4f8da0a](https://gitea.com/gitea/runner/commits/4f8da0a51cd44753f3e4ea0364877f6d2aa11f1d))
* **mergify:** fix deprecated strict mode ([#862](https://gitea.com/gitea/runner/issues/862)) ([80db062](https://gitea.com/gitea/runner/commits/80db062472a72735f5fa454a1f3977ea8702aad4))
* pin goreleaser version ([#825](https://gitea.com/gitea/runner/issues/825)) ([b1e4844](https://gitea.com/gitea/runner/commits/b1e4844aac4b8a143312c3ec9d3859362c669f88))
* remove macOS tests ([#860](https://gitea.com/gitea/runner/issues/860)) ([aad724c](https://gitea.com/gitea/runner/commits/aad724c87abb04184c9205643ec4659043637b07))
* replace superlinter with megalinter ([#923](https://gitea.com/gitea/runner/issues/923)) ([9be6a58](https://gitea.com/gitea/runner/commits/9be6a58c0eb706d514e355aa10f2b4c638cfec4a))
* **stalebot:** ignore feature requests ([#1017](https://gitea.com/gitea/runner/issues/1017)) ([932ad0e](https://gitea.com/gitea/runner/commits/932ad0ef6a2483b86f331f933b9c9081326ecbac))
* switch macos job to macos-10.15 ([#849](https://gitea.com/gitea/runner/issues/849)) ([977b526](https://gitea.com/gitea/runner/commits/977b526384a00e4f0804db6ddbaa3248c1918c75))

### Documentation

* add brew installation from HEAD ([#886](https://gitea.com/gitea/runner/issues/886)) ([df4ef4d](https://gitea.com/gitea/runner/commits/df4ef4de808edb8e3b657e4e24c96f8643ffcd32))
* add naming rule for git ref ([70cc6c0](https://gitea.com/gitea/runner/commits/70cc6c017bc8382b2efe4a50266f733037094da8))
* clarifying skipping steps / jobs ([#1480](https://gitea.com/gitea/runner/issues/1480)) ([7073eac](https://gitea.com/gitea/runner/commits/7073eac24057020cd54bcb00402281be55df3246))
* Expand the GITHUB_TOKEN section ([#968](https://gitea.com/gitea/runner/issues/968)) ([b1f5963](https://gitea.com/gitea/runner/commits/b1f5963c8623d5a5a92d2783b30da73c339232e2))
* fix brew formula link ([#846](https://gitea.com/gitea/runner/issues/846)) ([2319eb2](https://gitea.com/gitea/runner/commits/2319eb210f8d097da470331d2d12e7c6db0f8056))
* **README:** add --env-file flag ([#147](https://gitea.com/gitea/runner/issues/147)) ([f6b7613](https://gitea.com/gitea/runner/commits/f6b761378af17296aba8ff3be54e26d56acfeae5))
* **README:** remove go install ([#1123](https://gitea.com/gitea/runner/issues/1123)) ([e15f7a5](https://gitea.com/gitea/runner/commits/e15f7a51c6cf3e76cdc8ef8401b87080b7e9f95e))
* remove help section ([#1648](https://gitea.com/gitea/runner/issues/1648)) ([4b56ace](https://gitea.com/gitea/runner/commits/4b56aced159489639e5ca8226e2b82f5c7e9823e))

### Tests

* add test for networking setup in act ([#1375](https://gitea.com/gitea/runner/issues/1375)) ([ff5e289](https://gitea.com/gitea/runner/commits/ff5e2898041230978cef9af37dd87c09b07a77e8))
* Do not leak step env in composite ([#1585](https://gitea.com/gitea/runner/issues/1585)) ([3748772](https://gitea.com/gitea/runner/commits/3748772201458442afde9a582d70b40db05deaf8))
* Enshure ForcePull config doesn't break docker actions (1661) ([24440d9](https://gitea.com/gitea/runner/commits/24440d9f15673993dc6986f2c64dbd056e4286fe))
* env.ACT in if condition of the workflow ([#965](https://gitea.com/gitea/runner/issues/965)) ([3f0375a](https://gitea.com/gitea/runner/commits/3f0375aeff1973b674b33a0e4c44bc05a41714c5))
* make sure workflow_call is not a github event calling our workflow ([#1520](https://gitea.com/gitea/runner/issues/1520)) ([3f4a6dc](https://gitea.com/gitea/runner/commits/3f4a6dc4626eed03991e207b4d0e343b634edf9e))
* use ping to improve network test ([#2266](https://gitea.com/gitea/runner/issues/2266)) ([f56dd65](https://gitea.com/gitea/runner/commits/f56dd65ff69f38f68bf82d485c3cb1b90e9ffc3e))

### Chores

* add Winget installation instructions ([#1622](https://gitea.com/gitea/runner/issues/1622)) ([c378a7d](https://gitea.com/gitea/runner/commits/c378a7d28b5c20de9fbb4456544329d04ffb62a4))
* bump version ([0d07805](https://gitea.com/gitea/runner/commits/0d0780580cfe63786c21a31500ccb69797cb9cf6))
* bump version ([33cb8bc](https://gitea.com/gitea/runner/commits/33cb8bca649a1b63ac59596b20d4001ca9cbd3e6))
* bump VERSION to 0.2.39 ([efb12b7](https://gitea.com/gitea/runner/commits/efb12b7f122e1daa3fce410a0808098a23dad910))
* bump VERSION to 0.2.40 ([3e23b4f](https://gitea.com/gitea/runner/commits/3e23b4fbb558bc4e8bdc72d4f75e0d94e5be3bf1))
* bump VERSION to 0.2.41 ([b91e4b0](https://gitea.com/gitea/runner/commits/b91e4b0d5576ee0efeb381c9427b791e5272ece9))
* bump VERSION to 0.2.42 ([72d0321](https://gitea.com/gitea/runner/commits/72d03214c5ed8985bf23dc15440fa5301e9d13bf))
* bump VERSION to 0.2.43 ([44b510f](https://gitea.com/gitea/runner/commits/44b510f48cac6a58654bf920ee26cc6e1d36ac3a))
* bump VERSION to 0.2.44 ([220d6f1](https://gitea.com/gitea/runner/commits/220d6f1251aea75c378e5a98546f43cd918e1892))
* bump VERSION to 0.2.45 ([f1df2ca](https://gitea.com/gitea/runner/commits/f1df2ca5d6b315935cf9bb0d2a989175adf4e894))
* bump VERSION to 0.2.46 ([b0d0cec](https://gitea.com/gitea/runner/commits/b0d0cec71f2240fc629dcaa1889ed8c4ad52beb9))
* bump VERSION to 0.2.47 ([310cb79](https://gitea.com/gitea/runner/commits/310cb79e81e07ef11985b2f934092b82b167ccf1))
* bump VERSION to 0.2.48 ([3a0a642](https://gitea.com/gitea/runner/commits/3a0a6425a80f3405670d8d7b11ed18ea56a09af4))
* bump VERSION to 0.2.49 ([d5d8548](https://gitea.com/gitea/runner/commits/d5d854854651c151ecd914bb6e2e370e0057929a))
* bump VERSION to 0.2.50 ([80b0955](https://gitea.com/gitea/runner/commits/80b0955303888742c3ab73af5758bb7b01f5f57c))
* bump VERSION to 0.2.51 ([c241ecd](https://gitea.com/gitea/runner/commits/c241ecda318e4977501971c267f0a2fad70ea5cd))
* bump VERSION to 0.2.52 ([44ea01c](https://gitea.com/gitea/runner/commits/44ea01c2097eaad1596f4dafd6ec84872bd00402))
* bump VERSION to 0.2.53 ([1bb2ee7](https://gitea.com/gitea/runner/commits/1bb2ee7098f338765ce5169c850d0f19e59d3147))
* bump VERSION to 0.2.54 ([7693697](https://gitea.com/gitea/runner/commits/7693697f4cc60143c0561668b1c6fcb4bf733b94))
* bump VERSION to 0.2.55 ([95ff5bf](https://gitea.com/gitea/runner/commits/95ff5bf299a18833b19f157386184539ad6f4e64))
* bump VERSION to 0.2.56 ([0c60f97](https://gitea.com/gitea/runner/commits/0c60f9749a35a451d40e0741839f7bc403080c8d))
* bump VERSION to 0.2.57 ([4fae81e](https://gitea.com/gitea/runner/commits/4fae81efe4cdd9e09e7ef8e874a2d63b1ed98524))
* bump VERSION to 0.2.58 ([3ed38d8](https://gitea.com/gitea/runner/commits/3ed38d8e8b63f435cb845fbe9cdccc7ff032f402))
* bump VERSION to 0.2.59 ([b7a8145](https://gitea.com/gitea/runner/commits/b7a8145d09a7469e03c57e24230cbba109218faf))
* bump VERSION to 0.2.60 ([d3c8664](https://gitea.com/gitea/runner/commits/d3c8664d3d6acb855ae5418f6d02930c98bf877e))
* bump VERSION to 0.2.61 ([361b7e9](https://gitea.com/gitea/runner/commits/361b7e9f1abe9081edf594aa35dede4eab972c70))
* **deps:** cleanup go.mod a little bit ([#136](https://gitea.com/gitea/runner/issues/136)) ([ca648a3](https://gitea.com/gitea/runner/commits/ca648a37c8df80a4e96ffa7c3360c13c1a0715be))
* fix promote script ([b3c06dc](https://gitea.com/gitea/runner/commits/b3c06dc19d5ba396336615cb404cd394bd3e9f30))
* fix release script to trigger gh-act ([787388d](https://gitea.com/gitea/runner/commits/787388daf54f3d78c56da29a37062b019e10e84f))
* fold integration test outputs ([#1111](https://gitea.com/gitea/runner/issues/1111)) ([a85e89d](https://gitea.com/gitea/runner/commits/a85e89d3fb8e1a3b2e9dda843e18aa819095e56e))
* makefile update ([872a695](https://gitea.com/gitea/runner/commits/872a69548c77b07c01ba18fb507ffe0bfe278f5e))
* push new version ([003947e](https://gitea.com/gitea/runner/commits/003947ea33bfabcaa8c3787aa78780fd88b9c1ff))
* Remove obsolete Container.UpdateFromPath ([#1631](https://gitea.com/gitea/runner/issues/1631)) ([21ea3d0](https://gitea.com/gitea/runner/commits/21ea3d0d5fc92cf9ae93c62c0af7ec3a5da36541))
* remove repetitive words ([#2259](https://gitea.com/gitea/runner/issues/2259)) ([3949d74](https://gitea.com/gitea/runner/commits/3949d74af58d72309706e8fc5353de080fae6f3e))
* run act from cli on linux ([#1758](https://gitea.com/gitea/runner/issues/1758)) ([c81a770](https://gitea.com/gitea/runner/commits/c81a770bc5f8041c2c169e227996ea34ae62ec48)), references [#1756](https://gitea.com/gitea/runner/issues/1756)
* update docs for installing act as GH CLI extension ([7ada9d3](https://gitea.com/gitea/runner/commits/7ada9d3f74f4d8705ba01ce110f3d611989dd442))
* update go-git dependency in go.mod ([#30](https://gitea.com/gitea/runner/issues/30)) ([8f12a6c](https://gitea.com/gitea/runner/commits/8f12a6c947deaa3488013503a518144c9ecdd035))
* update LICENSE ([7fb84a5](https://gitea.com/gitea/runner/commits/7fb84a54a8dd93b5874c7edb89861f1c7d23b469))
* upgrade docker v19.03.12 ([#348](https://gitea.com/gitea/runner/issues/348)) ([3eef111](https://gitea.com/gitea/runner/commits/3eef111e46f42f43e62b64e6cdfbade4fe0935e2))
* upgrade golangci-lint and address findings ([#1904](https://gitea.com/gitea/runner/issues/1904)) ([724ec91](https://gitea.com/gitea/runner/commits/724ec918c92c162e537ba879e6d1ccc997f63483))
* use new style octal ([#1630](https://gitea.com/gitea/runner/issues/1630)) ([1316307](https://gitea.com/gitea/runner/commits/131630731374a16dbb7c8ed4893397c52d9c4c98))
* **yaml:** Improve `ParseRawOn` ([#28](https://gitea.com/gitea/runner/issues/28)) ([3daf313](https://gitea.com/gitea/runner/commits/3daf3132057d33bee325857b016a60db332fd050))

### Reverts

* Revert "Erase needs of job in SingleWorkflow (#9)" (#32) ([568f053](https://gitea.com/gitea/runner/commits/568f0537239f4566ad2b05b8866611d5af612b94)), references [#9](https://gitea.com/gitea/runner/issues/9) [#32](https://gitea.com/gitea/runner/issues/32)
* Revert "feat: allow existing logger from context (#859)" (#890) ([4ae71b5](https://gitea.com/gitea/runner/commits/4ae71b5878a8f6906ad3d2cd1d8cef0791527f20)), references [#859](https://gitea.com/gitea/runner/issues/859) [#890](https://gitea.com/gitea/runner/issues/890)
* deprecation of containerArchitecture ([#1514](https://gitea.com/gitea/runner/issues/1514)) ([4989f44](https://gitea.com/gitea/runner/commits/4989f444f1085ac0ec6b9bc35de92dc22f3c9f5d))
* Don't stop JobContainer on failure ([#1041](https://gitea.com/gitea/runner/issues/1041)) ([4d71071](https://gitea.com/gitea/runner/commits/4d7107161e231406bee856258ec847844623224e))

## [0.2.13](https://gitea.com/gitea/runner/compare/v0.2.12...v0.2.13) (2025-08-28)

### Features

* **docker:** add TZ env variable working ([#738](https://gitea.com/gitea/runner/issues/738)) ([aa28f8d](https://gitea.com/gitea/runner/commits/aa28f8d99c88fe7d662e49efcb8caa6215213058))
* support github mirror ([#716](https://gitea.com/gitea/runner/issues/716)) ([bbf9d7e](https://gitea.com/gitea/runner/commits/bbf9d7e90f3fee9e9968fc49c1ee6cfdb11f55e9))

## [0.2.12](https://gitea.com/gitea/runner/compare/v0.2.11...v0.2.12) (2025-06-11)

### Features

* add once flag to daemon command ([#19](https://gitea.com/gitea/runner/issues/19)) ([#598](https://gitea.com/gitea/runner/issues/598)) ([8bc0275](https://gitea.com/gitea/runner/commits/8bc0275e74af3516cb9d952e144e48405c084aa5))
* docker env vars for ephemeral and once ([#685](https://gitea.com/gitea/runner/issues/685)) ([44b4736](https://gitea.com/gitea/runner/commits/44b473670330342d66300d6ac96bccd0e39c37e7))
* register interactive with values from cli ([#682](https://gitea.com/gitea/runner/issues/682)) ([a616ed1](https://gitea.com/gitea/runner/commits/a616ed1a10efcb91f8ff3f67b35e16686cfa28b9))

### Bug Fixes

* invalid label NoInteractive exit code ([#683](https://gitea.com/gitea/runner/issues/683)) ([f0b5aff](https://gitea.com/gitea/runner/commits/f0b5aff3bbeda469a682fe21144c16485ea2189a)), closes [#665](https://gitea.com/gitea/runner/issues/665)

## [0.2.11](https://gitea.com/gitea/runner/compare/v0.2.10...v0.2.11) (2024-09-24)

### Features

* allow graceful shutdowns ([#546](https://gitea.com/gitea/runner/issues/546)) ([d1d3cad](https://gitea.com/gitea/runner/commits/d1d3cad4b0bfd0bfc5df344306f304043ff63223))

### Continuous Integration

* standardize code style and update version extraction ([#566](https://gitea.com/gitea/runner/issues/566)) ([2f78411](https://gitea.com/gitea/runner/commits/2f78411c3d871f0046832e4d082ba67c0655baf7))

## [0.2.10](https://gitea.com/gitea/runner/compare/v0.2.9...v0.2.10) (2024-04-10)

## [0.2.9](https://gitea.com/gitea/runner/compare/v0.2.8...v0.2.9) (2024-04-02)

## [0.2.8](https://gitea.com/gitea/runner/compare/v0.2.7...v0.2.8) (2024-03-27)

## [0.2.7](https://gitea.com/gitea/runner/compare/v0.2.6...v0.2.7) (2024-03-27)

### Bug Fixes

* **deps:** update module code.gitea.io/gitea-vet to v0.2.3 ([#410](https://gitea.com/gitea/runner/issues/410)) ([feb3966](https://gitea.com/gitea/runner/commits/feb39666cc7f72d59a776ac301ea66abe849403e))
* **deps:** update module github.com/avast/retry-go/v4 to v4.5.1 ([#411](https://gitea.com/gitea/runner/issues/411)) ([8c8a8ce](https://gitea.com/gitea/runner/commits/8c8a8ce40151ec5c28eeb80b24a83c9a509a3839))
* **deps:** update module github.com/docker/docker to v24.0.7+incompatible ([#413](https://gitea.com/gitea/runner/issues/413)) ([2f3e5c7](https://gitea.com/gitea/runner/commits/2f3e5c712523079ac9677158339b714234a835f2))
* **deps:** update module github.com/mattn/go-isatty to v0.0.20 ([#414](https://gitea.com/gitea/runner/issues/414)) ([825c6f9](https://gitea.com/gitea/runner/commits/825c6f97b79b33f3b103ab45c6d7a328cf2a464f))
* **deps:** update module github.com/spf13/cobra to v1.8.0 ([#416](https://gitea.com/gitea/runner/issues/416)) ([4d9de6c](https://gitea.com/gitea/runner/commits/4d9de6ca8ce74c2718ce6e3417cd3970317fa11a))
* **deps:** update module golang.org/x/time to v0.4.0 ([#424](https://gitea.com/gitea/runner/issues/424)) ([91bfe4c](https://gitea.com/gitea/runner/commits/91bfe4c186bf38daad2e3596912657bfce68849a))
* **deps:** update module golang.org/x/time to v0.5.0 ([#429](https://gitea.com/gitea/runner/issues/429)) ([08c681b](https://gitea.com/gitea/runner/commits/08c681be0c3e7e1f10eace654b8161cd5d9b5d03))
* Exit docker container on any process exit ([#462](https://gitea.com/gitea/runner/issues/462)) ([f1f9142](https://gitea.com/gitea/runner/commits/f1f9142a3c9ec062dd0ca222456b24c4a105dac9)), references [#460](https://gitea.com/gitea/runner/issues/460)

### Chores

* **deps:** update actions/setup-go action to v4 ([#418](https://gitea.com/gitea/runner/issues/418)) ([0adfc1c](https://gitea.com/gitea/runner/commits/0adfc1c7cc315c564071a95ca19b923616355dc5))
* **deps:** update docker/build-push-action action to v5 ([#419](https://gitea.com/gitea/runner/issues/419)) ([e3c6866](https://gitea.com/gitea/runner/commits/e3c68668fa3a97d88613ce2a73d73da1836d62d7))
* **deps:** update docker/login-action action to v3 ([#420](https://gitea.com/gitea/runner/issues/420)) ([f1b27d5](https://gitea.com/gitea/runner/commits/f1b27d52747c63da5643b7bdde8996208faff208))
* **deps:** update docker/setup-buildx-action action to v3 ([#421](https://gitea.com/gitea/runner/issues/421)) ([655a39f](https://gitea.com/gitea/runner/commits/655a39fd619c103c87277b4fa5c4fa7e2d23b5c5))
* **deps:** update docker/setup-qemu-action action to v3 ([#422](https://gitea.com/gitea/runner/issues/422)) ([cca7d54](https://gitea.com/gitea/runner/commits/cca7d54117a78bb84c2efc78e0cf8f63876962c6))

## [0.2.6](https://gitea.com/gitea/runner/compare/v0.2.51...v0.2.6) (2023-09-24)

## [0.2.5](https://gitea.com/gitea/runner/compare/v0.2.4...v0.2.5) (2023-08-02)

## [0.2.4](https://gitea.com/gitea/runner/compare/v0.2.3...v0.2.4) (2023-07-25)

### Bug Fixes

* **register:** refactor context usage in registration functions ([#286](https://gitea.com/gitea/runner/issues/286)) ([8705871](https://gitea.com/gitea/runner/commits/87058716fb95ffb715472ee7fab2e8fb8a6e75b0))

### Refactoring

* **register:** refactor registration functions and improve logging ([#288](https://gitea.com/gitea/runner/issues/288)) ([cf92a97](https://gitea.com/gitea/runner/commits/cf92a979e274dace1d545cbe140f4f4db064f921))

### Continuous Integration

* **lint:** refactor code for clarity and linting compliance ([#289](https://gitea.com/gitea/runner/issues/289)) ([db662b3](https://gitea.com/gitea/runner/commits/db662b36906a3aa0bd1211cee64e4ce73288f485))

## [0.2.3](https://gitea.com/gitea/runner/compare/v0.2.2...v0.2.3) (2023-07-03)

## [0.2.2](https://gitea.com/gitea/runner/compare/v0.2.1...v0.2.2) (2023-07-01)

## [0.2.1](https://gitea.com/gitea/runner/compare/v0.2.0...v0.2.1) (2023-06-30)

### Features

* improve Docker configuration and detection handling ([#242](https://gitea.com/gitea/runner/issues/242)) ([9e4a5f7](https://gitea.com/gitea/runner/commits/9e4a5f7363a5559ddaa3a48e91bb8e69db861de4))

## [0.2.0](https://gitea.com/gitea/runner/compare/v0.1.8...v0.2.0) (2023-06-15)

### Refactoring

* daemon, config, and logging for better clarity ([#225](https://gitea.com/gitea/runner/issues/225)) ([69c55ee](https://gitea.com/gitea/runner/commits/69c55ee003a3553d4e8b50d9bf98d2c31c8832e4))

### Documentation

* improve examples README and organization ([#230](https://gitea.com/gitea/runner/issues/230)) ([a83f29d](https://gitea.com/gitea/runner/commits/a83f29d5a96a78b61a9dc6967d09a31ef0ad5d7e))

## [0.1.8](https://gitea.com/gitea/runner/compare/v0.1.7...v0.1.8) (2023-05-09)

## [0.1.7](https://gitea.com/gitea/runner/compare/v0.1.6...v0.1.7) (2023-05-06)

### Continuous Integration

* improve release process and test robustness ([#173](https://gitea.com/gitea/runner/issues/173)) ([eef3c32](https://gitea.com/gitea/runner/commits/eef3c32eb2418dcbcd40cb13076c2887d4f6913b))

### Chores

* improve Dockerfile, README, and testing settings ([#172](https://gitea.com/gitea/runner/issues/172)) ([c40b651](https://gitea.com/gitea/runner/commits/c40b6518736636db6edf7bbecf4d650f3a2117c4))

## [0.1.6](https://gitea.com/gitea/runner/compare/v0.1.5...v0.1.6) (2023-05-01)

### Build

* improve compression and update GitHub actions ([#168](https://gitea.com/gitea/runner/issues/168)) ([b498341](https://gitea.com/gitea/runner/commits/b49834185777fffed9e0b4d7c3b08ad82fcdbcfa))
* optimize Dockerfile and update dependencies ([#162](https://gitea.com/gitea/runner/issues/162)) ([0d727eb](https://gitea.com/gitea/runner/commits/0d727eb2621c706d973987490bb3b14163693cf3))

## [0.1.5](https://gitea.com/gitea/runner/compare/v0.1.4...v0.1.5) (2023-04-28)

### Continuous Integration

* improve API usage and test robustness across platforms ([#159](https://gitea.com/gitea/runner/issues/159)) ([49d2cb0](https://gitea.com/gitea/runner/commits/49d2cb0cb53ebe6f67284b2fd27e11fe81884d85))

## [0.1.4](https://gitea.com/gitea/runner/compare/v0.1.3...v0.1.4) (2023-04-27)

## [0.1.3](https://gitea.com/gitea/runner/compare/v0.1.2...v0.1.3) (2023-04-27)

### Continuous Integration

* **actions:** add build docker image workflow ([#118](https://gitea.com/gitea/runner/issues/118)) ([66bab3d](https://gitea.com/gitea/runner/commits/66bab3d8059fde868b59d1aeb9e0bdecb6c0d55e))

## [0.1.2](https://gitea.com/gitea/runner/compare/v0.1.1...v0.1.2) (2023-04-13)

## [0.1.1](https://gitea.com/gitea/runner/compare/v0.1.0...v0.1.1) (2023-04-11)

### Features

* add artifact api from gitea server ([#103](https://gitea.com/gitea/runner/issues/103)) ([b79c3aa](https://gitea.com/gitea/runner/commits/b79c3aa1a325ab49902fd662174a5ea0f8eed196))

## [0.1.0](https://gitea.com/gitea/runner/compare/v0.0.1...v0.1.0) (2023-04-04)

## [0.0.1](https://gitea.com/gitea/runner/compare/a04141c4440be5c0f003d5441e4e2a97bc09060e...v0.0.1) (2023-02-26)

### Features

* `--container-options` ([#1462](https://gitea.com/gitea/runner/issues/1462)) ([4c2524a](https://gitea.com/gitea/runner/commits/4c2524ab4d0d923f241f2649702e5dd8927846a4))
* add `go` install instruction ([#672](https://gitea.com/gitea/runner/issues/672)) ([37a98f1](https://gitea.com/gitea/runner/commits/37a98f134a7e69bd45b0da40ea175d4cb5b41211))
* add bug-report flag ([#1056](https://gitea.com/gitea/runner/issues/1056)) ([5d7027d](https://gitea.com/gitea/runner/commits/5d7027dc3fc6a16e2023ffa1ebac93f1c667786e))
* add check for newer versions ([#1562](https://gitea.com/gitea/runner/issues/1562)) ([9390793](https://gitea.com/gitea/runner/commits/93907931df737bdd21ecc624f4983fb1b365357a))
* Add CopyTarStream to Container Interface ([#1912](https://gitea.com/gitea/runner/issues/1912)) ([6468dd7](https://gitea.com/gitea/runner/commits/6468dd7fc8c2013dc76e8ea6e9ab590dc6dbb56e))
* add Dependabot config ([#797](https://gitea.com/gitea/runner/issues/797)) ([12928a0](https://gitea.com/gitea/runner/commits/12928a0ac61b2b3278aead1e42923c876a8cf5a4))
* add flag to always run the checkout action ([#1049](https://gitea.com/gitea/runner/issues/1049)) ([f09cdb8](https://gitea.com/gitea/runner/commits/f09cdb844367189aafe2f3af5795d3136db6b290))
* add fromJSON support ([#352](https://gitea.com/gitea/runner/issues/352)) ([d521fa5](https://gitea.com/gitea/runner/commits/d521fa5bbae1d6fd5b775e0924907e882a4d8655))
* add json logger output ([#1026](https://gitea.com/gitea/runner/issues/1026)) ([14c9801](https://gitea.com/gitea/runner/commits/14c98015c57bf577619101ac4770f90c6435470b))
* add logger hook for standalone job logger ([0fa1fe0](https://gitea.com/gitea/runner/commits/0fa1fe03109cb5e41aa7a074fe0d837089d65c1e))
* Add new Action Cache ([#1913](https://gitea.com/gitea/runner/issues/1913)) ([a42f3cf](https://gitea.com/gitea/runner/commits/a42f3cf1cdaa1d93da181bd5f556371d61198169))
* add node16 support ([#922](https://gitea.com/gitea/runner/issues/922)) ([5580812](https://gitea.com/gitea/runner/commits/558081242ca322997437de126324ae9977eb7da0))
* add option for custom socket path ([#698](https://gitea.com/gitea/runner/issues/698)) ([799098b](https://gitea.com/gitea/runner/commits/799098b0e630caff38f220dd0ab5e216667ec916))
* add option for docker image rebuild ([#878](https://gitea.com/gitea/runner/issues/878)) ([8ad6c07](https://gitea.com/gitea/runner/commits/8ad6c07083c7cd7f45ea9e0808e3297bf2370ab5)), references [#714](https://gitea.com/gitea/runner/issues/714) [#716](https://gitea.com/gitea/runner/issues/716)
* add option to bypass GHE for actions checkout ([#1162](https://gitea.com/gitea/runner/issues/1162)) ([de37f75](https://gitea.com/gitea/runner/commits/de37f750773575651a91363eb6b1ce97be964845)), references [#1161](https://gitea.com/gitea/runner/issues/1161)
* add option to specify git remote name ([#1104](https://gitea.com/gitea/runner/issues/1104)) ([d4272bd](https://gitea.com/gitea/runner/commits/d4272bd9fe13566ab76e5ad091aedbded74ff4e2))
* add remote reusable workflows ([#1525](https://gitea.com/gitea/runner/issues/1525)) ([82a8c1e](https://gitea.com/gitea/runner/commits/82a8c1e80dfcd73842bc0359ca8447303b6fcc08)), closes [#826](https://gitea.com/gitea/runner/issues/826)
* add skipped status as step result ([#950](https://gitea.com/gitea/runner/issues/950)) ([2eda7c6](https://gitea.com/gitea/runner/commits/2eda7c6037b921ae88742b2ff29b20bb171497a6))
* add step name to logger field ([#1027](https://gitea.com/gitea/runner/issues/1027)) ([6b05957](https://gitea.com/gitea/runner/commits/6b059572ff51e30f4cfbde76cc7b82af313967b6))
* add support for building docker actions with private registries ([#1557](https://gitea.com/gitea/runner/issues/1557)) ([3f3b25a](https://gitea.com/gitea/runner/commits/3f3b25ae84a25f8488f3801e3342bb3f96c6cd8a))
* Add support for github action variables ([#1833](https://gitea.com/gitea/runner/issues/1833)) ([6ce45e3](https://gitea.com/gitea/runner/commits/6ce45e3f246f12d6617691de9a2423920d5fdbe6))
* add task log output hook ([dfcb104](https://gitea.com/gitea/runner/commits/dfcb104af02b536df7dfd7dc6d5911b654038ae9))
* add ubuntu-22.04 ([#1150](https://gitea.com/gitea/runner/issues/1150)) ([69691c8](https://gitea.com/gitea/runner/commits/69691c88b301ee695f5712d2ce8ad2b760725c80))
* adjust container name format ([#1](https://gitea.com/gitea/runner/issues/1)) ([4cacc14](https://gitea.com/gitea/runner/commits/4cacc14d22447358f1d34b4ef556fe1fd9cb9dd0))
* Allow building without docker support ([#1507](https://gitea.com/gitea/runner/issues/1507)) ([7e8d070](https://gitea.com/gitea/runner/commits/7e8d0708111125b11166494dc22da6d67f9e83e7))
* allow existing logger from context ([#859](https://gitea.com/gitea/runner/issues/859)) ([603cf56](https://gitea.com/gitea/runner/commits/603cf568780752efbe587dd50825b2ee5c3013ae))
* allow existing logger from context ([#898](https://gitea.com/gitea/runner/issues/898)) ([6517d04](https://gitea.com/gitea/runner/commits/6517d04b9a4a677fc03277b6cf32fcd56a4dcf36))
* allow overriding of `GITHUB_` env variables ([#1582](https://gitea.com/gitea/runner/issues/1582)) ([ce168f9](https://gitea.com/gitea/runner/commits/ce168f959502feab3ec815810fac7a3b5a1f470c))
* allow to spawn and run a local reusable workflow ([#1423](https://gitea.com/gitea/runner/issues/1423)) ([a8e05cd](https://gitea.com/gitea/runner/commits/a8e05cded6be8bfb77536d823933725144918bcb))
* auto remove containers ([8f950cc](https://gitea.com/gitea/runner/commits/8f950ccec0a9aae4a6e66d10dc28edd5fdf04ccb))
* avoid reuse of containers ([cbd21b8](https://gitea.com/gitea/runner/commits/cbd21b8ec539a5def5649b10702e55f6656e2ee9))
* cache notices to reduce frequency of upgrade notifications ([#1592](https://gitea.com/gitea/runner/issues/1592)) ([42b9b73](https://gitea.com/gitea/runner/commits/42b9b73d388544693d44d419457151c10e66472a))
* **cli:** add register command basic functions ([561bfad](https://gitea.com/gitea/runner/commits/561bfad7c5413e29f61f44fc4526ad6123705c4c))
* **cli:** make register interactive stages working ([bfcc6f5](https://gitea.com/gitea/runner/commits/bfcc6f56f320d891ba61e886acf3b49ac6d766e2))
* **cli:** set interactive stages in register command ([8f9d750](https://gitea.com/gitea/runner/commits/8f9d7506dc71aaa55db74b80a7118bc18d6947f2))
* cmd support for windows ([#1941](https://gitea.com/gitea/runner/issues/1941)) ([8314095](https://gitea.com/gitea/runner/commits/83140951bf96c500f0c1a745e4333b27f2820e49))
* docker 19.03.12 ([#347](https://gitea.com/gitea/runner/issues/347)) ([b3b8961](https://gitea.com/gitea/runner/commits/b3b8961122d06c17bb9bbc063bf56e12aba60106))
* Don't push `pkg/runner/act/` where actions cache is stored ([#620](https://gitea.com/gitea/runner/issues/620)) ([b04d762](https://gitea.com/gitea/runner/commits/b04d762614433695202ee1c12d16654b558165bc))
* don't require docker ([#16](https://gitea.com/gitea/runner/issues/16)) ([990cf93](https://gitea.com/gitea/runner/commits/990cf93c7136669408eb1832cd05df3ad4dd81b3))
* fill default task runner ([3281c66](https://gitea.com/gitea/runner/commits/3281c6637d5b49f761e80584235d67384541a04a))
* GITHUB_STATE and GITHUB_OUTPUT file commands ([#1391](https://gitea.com/gitea/runner/issues/1391)) ([e520382](https://gitea.com/gitea/runner/commits/e520382d2fdc5872ca8e939c6c0f2e1e74749182))
* **goreleaser:** format changelog, add prereleases ([#942](https://gitea.com/gitea/runner/issues/942)) ([4e6cddf](https://gitea.com/gitea/runner/commits/4e6cddf80a712469695640058bcfceee7adaeda6))
* handle context cancelation during docker exec ([#1170](https://gitea.com/gitea/runner/issues/1170)) ([4ef50ee](https://gitea.com/gitea/runner/commits/4ef50eeae76f76af3ab0c8c7377c9388d28a0c38))
* Host environment ([#1293](https://gitea.com/gitea/runner/issues/1293)) ([f2b98ed](https://gitea.com/gitea/runner/commits/f2b98ed301290f53e96c64814ab9949640e10307))
* implement steps.timeout-minutes ([#1776](https://gitea.com/gitea/runner/issues/1776)) ([8913375](https://gitea.com/gitea/runner/commits/8913375af89792be37603c3c6429372964bf86b1))
* improve GetOutboundIP ([#1707](https://gitea.com/gitea/runner/issues/1707)) ([8518d70](https://gitea.com/gitea/runner/commits/8518d70bdf3ee04f08e5d096fe07c2a85e8f21ea))
* improve list ([#786](https://gitea.com/gitea/runner/issues/786)) ([37aaec8](https://gitea.com/gitea/runner/commits/37aaec81f444bd5e00de4e8cedb35c73e131b41d))
* interpolate the step names ([#1422](https://gitea.com/gitea/runner/issues/1422)) ([809da71](https://gitea.com/gitea/runner/commits/809da7198cc2bed5b758c77274d47cb5360ee848)), closes [#1353](https://gitea.com/gitea/runner/issues/1353)
* JobLoggerFactory ([#1496](https://gitea.com/gitea/runner/issues/1496)) ([7d99512](https://gitea.com/gitea/runner/commits/7d9951200f366ee7a207844e185bf5de7e3f7ab5))
* load every environment from --env-file to workflow ([#184](https://gitea.com/gitea/runner/issues/184)) ([2f39547](https://gitea.com/gitea/runner/commits/2f395475b0ac4f5c0fb3d7c8b97a23d9214ea77b))
* move main task logic to runtime package ([c5edbbd](https://gitea.com/gitea/runner/commits/c5edbbd004418890b593dd9217160b02f5f2143d))
* non root user container ([#1202](https://gitea.com/gitea/runner/issues/1202)) ([b59e6de](https://gitea.com/gitea/runner/commits/b59e6dee6d01253cb96505f8effc0516c2a390d9))
* parse register arguments and do register when non-interactive mode ([45de619](https://gitea.com/gitea/runner/commits/45de6199d12626de6c934444624c24e8f887e8f6))
* parse types of reusable workflows ([#1414](https://gitea.com/gitea/runner/issues/1414)) ([9b95a72](https://gitea.com/gitea/runner/commits/9b95a728d84920a75d22c0dc6088391fbbc4002a)), references [#826](https://gitea.com/gitea/runner/issues/826)
* pass current matrix information to job logger ([#1364](https://gitea.com/gitea/runner/issues/1364)) ([00acf68](https://gitea.com/gitea/runner/commits/00acf681880b0fe7cea27a29d61c261d624966c3))
* **poller:** support scheduler to fetch task and dispatch to worker ([a13ed0c](https://gitea.com/gitea/runner/commits/a13ed0c52e2879ef5fa2ec8632da2efd535c56d0))
* **register:** add label format validation ([3f676a7](https://gitea.com/gitea/runner/commits/3f676a742b72f550b7900a1c0d4faab19b11750f))
* **register:** add labels to local runner file ([86f66ea](https://gitea.com/gitea/runner/commits/86f66eaa4508603d9359398e8dcf866d8689bb89))
* release extension ([78e1ceb](https://gitea.com/gitea/runner/commits/78e1ceb0a8ad911b04f821ca707c3b352fe2c74b))
* release extension ([d690a5f](https://gitea.com/gitea/runner/commits/d690a5fee95491d519085990d5c700ff2c9081ce))
* release extension ([19e6929](https://gitea.com/gitea/runner/commits/19e692939861bf3770e4c0ea7a7d8d5f27d84561))
* replace ciphertext in log ([d650bae](https://gitea.com/gitea/runner/commits/d650bae0b6212b5a38aef62a6fab13c04b971f45))
* run jobs in parallel ([#1003](https://gitea.com/gitea/runner/issues/1003)) ([18792f9](https://gitea.com/gitea/runner/commits/18792f9620fd2fd7dacc53ae30a234fbfc7c81b1)), references [#CPU](https://gitea.com/gitea/runner/issues/CPU)
* SELinux lables for --bind ([#906](https://gitea.com/gitea/runner/issues/906)) ([f726339](https://gitea.com/gitea/runner/commits/f7263399b9146a43b988e011d986d18d180414d2))
* specify matrix on command line ([#1675](https://gitea.com/gitea/runner/issues/1675)) ([636c8a3](https://gitea.com/gitea/runner/commits/636c8a34aedf9b6251df6729fe673d828f972030))
* Speeding up workflows ([#643](https://gitea.com/gitea/runner/issues/643)) ([3e768cd](https://gitea.com/gitea/runner/commits/3e768cd91695af12bdf860ba49a0f992854b9c87))
* split job steps into its own files/structs ([#1004](https://gitea.com/gitea/runner/issues/1004)) ([2bb3e74](https://gitea.com/gitea/runner/commits/2bb3e74616f86ece6a11cd854ad54701542410fd))
* step summary of test results ([#1580](https://gitea.com/gitea/runner/issues/1580)) ([932863b](https://gitea.com/gitea/runner/commits/932863bef508868168954ef69572750cf87f463f))
* Support "result" on "needs" context. ([#1497](https://gitea.com/gitea/runner/issues/1497)) ([bef9b5c](https://gitea.com/gitea/runner/commits/bef9b5c3c77bcdbc6c15169c0d55f635e27350ae)), closes [#1367](https://gitea.com/gitea/runner/issues/1367) [#1466](https://gitea.com/gitea/runner/issues/1466), references [#1423](https://gitea.com/gitea/runner/issues/1423) [#1506](https://gitea.com/gitea/runner/issues/1506)
* support `GITHUB_REF_NAME` & `GITHUB_REF_TYPE` ([#1142](https://gitea.com/gitea/runner/issues/1142)) ([8e216c6](https://gitea.com/gitea/runner/commits/8e216c642e11dfa8dba7aab040bb09274c48b5a6))
* support branches as GITHUB_REF ([03d57a8](https://gitea.com/gitea/runner/commits/03d57a80c9efe130fc0dde98ab7895ac312945c5))
* support custom GITHUB_RUN_ID, GITHUB_RUN_NUMBER ([#369](https://gitea.com/gitea/runner/issues/369)) ([1070c58](https://gitea.com/gitea/runner/commits/1070c585381980abe8f636af5212457711a51f03))
* support interpolation in `<job>.container.options` ([#1958](https://gitea.com/gitea/runner/issues/1958)) ([19764bc](https://gitea.com/gitea/runner/commits/19764bcb06ac9acdd7100d95c72055b697820705))
* support more options of containers ([b9c20dc](https://gitea.com/gitea/runner/commits/b9c20dcaa43899cb3bb327619d447248303170e0))
* support node20 runtime ([#1988](https://gitea.com/gitea/runner/issues/1988)) ([935e37c](https://gitea.com/gitea/runner/commits/935e37c25b6c48b50882024b3c4235d12f23f469))
* support PlatformPicker ([c6b8548](https://gitea.com/gitea/runner/commits/c6b8548d35ab4a5fc0f4478ef9f001ee4fbda16f))
* support post step log ([ab4e06f](https://gitea.com/gitea/runner/commits/ab4e06f977ca78181077e5de43d7aad070204192))
* support yaml env/secrets/inputs file ([#1733](https://gitea.com/gitea/runner/issues/1733)) ([50dcc57](https://gitea.com/gitea/runner/commits/50dcc57e4babe816692d342307183801ddd78356))
* task can report step and final result ([6030610](https://gitea.com/gitea/runner/commits/6030610c04cc75219d6047ce51e57b45ea3c1fb9))
* try to read ref and sha from event payload if available ([#889](https://gitea.com/gitea/runner/issues/889)) ([edd0fb9](https://gitea.com/gitea/runner/commits/edd0fb92a78f2fb20ea1afe14a047fdc026d2643))
* use specified labels ([2354f5b](https://gitea.com/gitea/runner/commits/2354f5bb18902124c0ccc07c632888bac8b4c99b))
* uses in composite ([#793](https://gitea.com/gitea/runner/issues/793)) ([9868e13](https://gitea.com/gitea/runner/commits/9868e1377274a7657d4aac6bd2b66c0d564405cb))
* workflowpattern package ([#1618](https://gitea.com/gitea/runner/issues/1618)) ([05eaeaa](https://gitea.com/gitea/runner/commits/05eaeaa528f791cdb796e3fae4cabd833dcb7a40))
* **workflow:** support schedule event ([#4](https://gitea.com/gitea/runner/issues/4)) ([88cce47](https://gitea.com/gitea/runner/commits/88cce47022d34d5cbc49771b3753541fe905f5b6))

### Bug Fixes

* --version, --help do not need check platform config ([#508](https://gitea.com/gitea/runner/issues/508)) ([22d56c3](https://gitea.com/gitea/runner/commits/22d56c3517483c63dc98dc4d280d60eb5a625a3d))
* ::stop-commands:: should continue to print the lines ([#920](https://gitea.com/gitea/runner/issues/920)) ([fec0c0c](https://gitea.com/gitea/runner/commits/fec0c0c529f5510464c858ff9a6664e0647245a8))
* [#170](https://gitea.com/gitea/runner/issues/170) [#171](https://gitea.com/gitea/runner/issues/171) - add RUNNER_OS and RUNNER_TEMP ([a9533b0](https://gitea.com/gitea/runner/commits/a9533b05ce708a84e931d63cb2cd7b795e70bbe2))
* [#170](https://gitea.com/gitea/runner/issues/170) add RUNNER_TEMP ([a48cce3](https://gitea.com/gitea/runner/commits/a48cce3a78d915d7c9e6bd28248926578caaf1c4))
* [#182](https://gitea.com/gitea/runner/issues/182) - support symlinks in workspace ([8c29099](https://gitea.com/gitea/runner/commits/8c290994c1b84153d2b9f911ac3bb64858fb5e35))
* 168 include error message from docker SDK ([6c4e7b9](https://gitea.com/gitea/runner/commits/6c4e7b9fdedeeeb1eda78cefaff55750033d1327))
* add `ImageOS` env var based on running platform ([#571](https://gitea.com/gitea/runner/issues/571)) ([41b03b5](https://gitea.com/gitea/runner/commits/41b03b581c7aec251d96e9a5ea27009542a798d4))
* add `server_url` attribute to github context ([#1727](https://gitea.com/gitea/runner/issues/1727)) ([c0130ed](https://gitea.com/gitea/runner/commits/c0130ed030e0f267b7b7b5b3ad5e1dfcc08980d8)), closes [#1726](https://gitea.com/gitea/runner/issues/1726)
* add default values in vmInput ([cc369e2](https://gitea.com/gitea/runner/commits/cc369e2f73e84596a936d602a7e112d849e6e6fa))
* add GITHUB_STEP_SUMMARY ([#1607](https://gitea.com/gitea/runner/issues/1607)) ([8790c9b](https://gitea.com/gitea/runner/commits/8790c9b8e6b67f1552b824ba56d5a851f7b394e5))
* add parent step id in composite action step id ([#1268](https://gitea.com/gitea/runner/issues/1268)) ([5f5b895](https://gitea.com/gitea/runner/commits/5f5b8959d6b12b034dfa64e91f8b17c97808eda5))
* add simple concurrency limit ([#823](https://gitea.com/gitea/runner/issues/823)) ([4d552e6](https://gitea.com/gitea/runner/commits/4d552e65cedafda89c11a36abe3a218f5d090e4b))
* add-matcher fails github workflow ([#1532](https://gitea.com/gitea/runner/issues/1532)) ([8b4f210](https://gitea.com/gitea/runner/commits/8b4f21087209e027ca0e2ec225e0c74d03552f86))
* added tests and improved code ([8d181a6](https://gitea.com/gitea/runner/commits/8d181a6a1ea1eca1076ae9d47c3fded31a6099d2))
* align github.ref to GitHub Action ([#1362](https://gitea.com/gitea/runner/issues/1362)) ([97c083e](https://gitea.com/gitea/runner/commits/97c083e902a86010ce6267347b830e91fe2435e1))
* align runner.os / runner.arch to known values ([#1510](https://gitea.com/gitea/runner/issues/1510)) ([6ab71ec](https://gitea.com/gitea/runner/commits/6ab71ecebbc0d01a0e72aecd06a15277723b6031))
* allow override of artifact server bind address ([#1560](https://gitea.com/gitea/runner/issues/1560)) ([d064863](https://gitea.com/gitea/runner/commits/d064863f9bb65cde863c8cb20fea50e0f139032c)), references [#1559](https://gitea.com/gitea/runner/issues/1559)
* allow TimeoutMinutes to be expression ([#1217](https://gitea.com/gitea/runner/issues/1217)) ([52f5c45](https://gitea.com/gitea/runner/commits/52f5c4592cbbf2e4dde3b35ad42f4e7475ad52df))
* allow TimeoutMinutes to be expression in Jobs ([#1247](https://gitea.com/gitea/runner/issues/1247)) ([91296bd](https://gitea.com/gitea/runner/commits/91296bd5eb022dce9e7891540257232b3e9906e1)), references [#1217](https://gitea.com/gitea/runner/issues/1217)
* always execute closeContainer() executor ([#988](https://gitea.com/gitea/runner/issues/988)) ([9abc87b](https://gitea.com/gitea/runner/commits/9abc87b41634cab2ad9325287fbed9135c2664c4)), references [#984](https://gitea.com/gitea/runner/issues/984)
* Apply forcePull only for prebuild docker actions ([#1599](https://gitea.com/gitea/runner/issues/1599)) ([f91b2aa](https://gitea.com/gitea/runner/commits/f91b2aa5dbcb16075f49c96a6108d0d418382f4d))
* avoid race condition by not running exec/start twice ([#702](https://gitea.com/gitea/runner/issues/702)) ([8a9167d](https://gitea.com/gitea/runner/commits/8a9167da82b3824272fdbfcc25b3671cd1260e7d)), references [#627](https://gitea.com/gitea/runner/issues/627) [#627](https://gitea.com/gitea/runner/issues/627)
* call of reflect.Value.Interface on zero Value ([#1067](https://gitea.com/gitea/runner/issues/1067)) ([380bd58](https://gitea.com/gitea/runner/commits/380bd581b1a56683a77b9eb017a72129f65251f8))
* call RunDaemon ([93e42d1](https://gitea.com/gitea/runner/commits/93e42d14650c8e1c55fbf37287c6f6e8ef2f24c5))
* checkout container path should always use `/` ([#1226](https://gitea.com/gitea/runner/issues/1226)) ([f71f277](https://gitea.com/gitea/runner/commits/f71f2778f0e75e3e28a81f8f726c64f27b7fffa3))
* ci is failing since 2022-10-17 ([#1397](https://gitea.com/gitea/runner/issues/1397)) ([5d2eb1f](https://gitea.com/gitea/runner/commits/5d2eb1f2388a7d99aae811721c70e75ed8b8a76c))
* CI is unstable due to new linting rules ([#1292](https://gitea.com/gitea/runner/issues/1292)) ([3387fd7](https://gitea.com/gitea/runner/commits/3387fd72ee138893b953e2b8b23efceb23bfefef))
* clarify --reuse and --rm flags ([#929](https://gitea.com/gitea/runner/issues/929)) ([4e0ba61](https://gitea.com/gitea/runner/commits/4e0ba618d34a198fa30904b1bf8549010513a931))
* clean code ([9fd6310](https://gitea.com/gitea/runner/commits/9fd63103632d5288cc7fa6c6f7e175e2ca17ad2f))
* cleanup unused actions ([#873](https://gitea.com/gitea/runner/issues/873)) ([3054f3e](https://gitea.com/gitea/runner/commits/3054f3ea1486267424d71c03c943eded2b95f50f))
* coerce booleans to numbers for comparison in exprparser ([#1030](https://gitea.com/gitea/runner/issues/1030)) ([aab2af0](https://gitea.com/gitea/runner/commits/aab2af09194a77a6b93d81088baf981cd5e27c17))
* conclusion and outcome are no integers ([#1136](https://gitea.com/gitea/runner/issues/1136)) ([77ddc89](https://gitea.com/gitea/runner/commits/77ddc89444d8cca5be47a3c83074dcb992382e68))
* continue jobs + steps after failure ([#840](https://gitea.com/gitea/runner/issues/840)) ([1891c72](https://gitea.com/gitea/runner/commits/1891c72ab158508e36009d16b24913fa5836422b)), closes [#442](https://gitea.com/gitea/runner/issues/442)
* copy ignored tracked files ([#1074](https://gitea.com/gitea/runner/issues/1074)) ([c27e079](https://gitea.com/gitea/runner/commits/c27e079407647e0ab455f068061b5e61e31f670f))
* copy whole action dir to container ([#774](https://gitea.com/gitea/runner/issues/774)) ([de44a50](https://gitea.com/gitea/runner/commits/de44a505da37f9405db195902b277821df2742b6))
* correct ref and ref_name ([#1672](https://gitea.com/gitea/runner/issues/1672)) ([6744e68](https://gitea.com/gitea/runner/commits/6744e68ee29d3bcb6878f895d6ccb5f9f8e2c478))
* correctly unmarshal entrypoint ([#876](https://gitea.com/gitea/runner/issues/876)) ([12fa4d7](https://gitea.com/gitea/runner/commits/12fa4d703da8d05930303acbe6bbc8fd00c77de7))
* crash if the id tool fails to run in the container (1660) ([f3c88b5](https://gitea.com/gitea/runner/commits/f3c88b50915a7234af9b06e28b4a45a267172fe1))
* crash malformed composite action ([#1616](https://gitea.com/gitea/runner/issues/1616)) ([53095d7](https://gitea.com/gitea/runner/commits/53095d76f48747f9f6249ba01e9410e9098d3977))
* crash on invalid docker credentials ([#1348](https://gitea.com/gitea/runner/issues/1348)) ([7d0407f](https://gitea.com/gitea/runner/commits/7d0407fa4abfe81308f8787cbd07c4c72b2aa676))
* deep evaluate matrix strategy ([#964](https://gitea.com/gitea/runner/issues/964)) ([d1e0216](https://gitea.com/gitea/runner/commits/d1e0216039d92a4253fecc2b11a775a37ba43ceb))
* do not error that much ([#1137](https://gitea.com/gitea/runner/issues/1137)) ([9357512](https://gitea.com/gitea/runner/commits/93575124d35f7b84860cb8cc553c257e0945d0af))
* docker cp of dangling symlink ([#943](https://gitea.com/gitea/runner/issues/943)) ([#948](https://gitea.com/gitea/runner/issues/948)) ([c802064](https://gitea.com/gitea/runner/commits/c802064975099a00c93b0c57f46b1e3964d33eef))
* don't allow `-self-hosted` mode as container image ([#1783](https://gitea.com/gitea/runner/issues/1783)) ([ca9b783](https://gitea.com/gitea/runner/commits/ca9b783491b4ea1c5d67dbffdb95f34fd4a22573))
* don't override env ([#1629](https://gitea.com/gitea/runner/issues/1629)) ([b0a5068](https://gitea.com/gitea/runner/commits/b0a5068f6d8ec06db4488b4b981d39f476d3abf6))
* don't overwrite with empty cmd/entrypoint ([#1076](https://gitea.com/gitea/runner/issues/1076)) ([24422bf](https://gitea.com/gitea/runner/commits/24422bf4e4e1e12b9a1e076b2f9e0a7b14c6f890))
* don't replace empty secrets ([#725](https://gitea.com/gitea/runner/issues/725)) ([96a0564](https://gitea.com/gitea/runner/commits/96a05645269241587fdcf0f1ccc49ce38047627c))
* dryrun skip local actions + enable Tests ([#1199](https://gitea.com/gitea/runner/issues/1199)) ([b7d380b](https://gitea.com/gitea/runner/commits/b7d380b3f096e079f2c1ba91a9377752fc4c9751))
* **editorconfig:** correct braces ([#1094](https://gitea.com/gitea/runner/issues/1094)) ([8d41f7a](https://gitea.com/gitea/runner/commits/8d41f7aa6b6702d9c1c3bfc372dd4c5533e88c8c))
* ensure all post steps are executed ([#1286](https://gitea.com/gitea/runner/issues/1286)) ([3364f95](https://gitea.com/gitea/runner/commits/3364f95569d8a39d440b18af69d29421fcc77b55))
* ensure networkmode "host" unless explicitly specified ([#1739](https://gitea.com/gitea/runner/issues/1739)) ([de06444](https://gitea.com/gitea/runner/commits/de0644499a555d240bd065f7c6d839e5162df46f))
* environment handling windows (host mode) ([#1732](https://gitea.com/gitea/runner/issues/1732)) ([9884da0](https://gitea.com/gitea/runner/commits/9884da0122b7fc45595c4f2950c1044fcce2c3e6))
* extra path lost in composite actions ([#1531](https://gitea.com/gitea/runner/issues/1531)) ([a53a1c2](https://gitea.com/gitea/runner/commits/a53a1c20005614cefd3d6519abd7f76d91467132)), references [#1528](https://gitea.com/gitea/runner/issues/1528)
* fail if no stages were found ([#1970](https://gitea.com/gitea/runner/issues/1970)) ([7286b43](https://gitea.com/gitea/runner/commits/7286b43b0ea0110cc2052848c8c1c4ab66bcd068))
* fail workflow if the job name is invalid ([#596](https://gitea.com/gitea/runner/issues/596)) ([94d736a](https://gitea.com/gitea/runner/commits/94d736a602dcfc7db2568d41a16e222bf1bcc416))
* fallback to unauthenticated pull ([#1774](https://gitea.com/gitea/runner/issues/1774)) ([aa21277](https://gitea.com/gitea/runner/commits/aa212773804561eaee04c3b4a7d9399dfbb9cad1))
* fromJSON should work with any JSON ([#883](https://gitea.com/gitea/runner/issues/883)) ([09693ec](https://gitea.com/gitea/runner/commits/09693ec5b9f0d2643b0a57826e2c415dfc32de2d))
* ghc assignment typo ([#1729](https://gitea.com/gitea/runner/issues/1729)) ([97749a2](https://gitea.com/gitea/runner/commits/97749a27b948268a921ef3e75d679a0f1f5c26ce))
* GITHUB_PATH should prepend ([#690](https://gitea.com/gitea/runner/issues/690)) ([29ea8cf](https://gitea.com/gitea/runner/commits/29ea8cfc4e00401b144cd29643a61ca089efd7e9))
* github.job property is empty, GITHUB_JOB should be job id ([#1646](https://gitea.com/gitea/runner/issues/1646)) ([5500c92](https://gitea.com/gitea/runner/commits/5500c928ebb6c4e7dcfb37c29c7100dd7e721b68)), references [#1621](https://gitea.com/gitea/runner/issues/1621) [#1473](https://gitea.com/gitea/runner/issues/1473)
* **goreleaser:** add append mode for release notes ([#962](https://gitea.com/gitea/runner/issues/962)) ([518148d](https://gitea.com/gitea/runner/commits/518148d1629381b0abd73e8f681d9fb6e0242205))
* handle env-vars case sensitive ([#1493](https://gitea.com/gitea/runner/issues/1493)) ([0f8861b](https://gitea.com/gitea/runner/commits/0f8861b9e37bfdf5e0fabb188d03365c0401823a)), closes [#1488](https://gitea.com/gitea/runner/issues/1488)
* handle go closure iteration ([#1374](https://gitea.com/gitea/runner/issues/1374)) ([1bade27](https://gitea.com/gitea/runner/commits/1bade27534209f2fd2d7a0b6437efb4965b648c8))
* handle subdirectors and exclusions ([#1012](https://gitea.com/gitea/runner/issues/1012)) ([#1446](https://gitea.com/gitea/runner/issues/1446)) ([5200c49](https://gitea.com/gitea/runner/commits/5200c491555ab1665cf280523f97bf92401c0415))
* handle zero size ([#1888](https://gitea.com/gitea/runner/issues/1888)) ([8c7c0f5](https://gitea.com/gitea/runner/commits/8c7c0f53c158867b446d1ef21356f49087c238a3))
* hardcode `ubuntu-latest` for `ImageOS` envvar ([#579](https://gitea.com/gitea/runner/issues/579)) ([29b3d43](https://gitea.com/gitea/runner/commits/29b3d43988e1354519f667830eef06662cadd9e4))
* hide masked value ([#668](https://gitea.com/gitea/runner/issues/668)) ([3db5037](https://gitea.com/gitea/runner/commits/3db50376aa4e0fbf5c9d206cb955c6f660c7779a))
* hide masked values in log files ([#1011](https://gitea.com/gitea/runner/issues/1011)) ([c22d833](https://gitea.com/gitea/runner/commits/c22d833830b3a1903a12c590b0b2ee637c3639e1))
* ignore `act` binary built when using `make` ([#654](https://gitea.com/gitea/runner/issues/654)) ([780f60a](https://gitea.com/gitea/runner/commits/780f60a1e6a07f495ad9a2d0e3181e85b36b0537))
* ignore stuck because of wrong token ([88ae188](https://gitea.com/gitea/runner/commits/88ae188699428e10d0d5d886f00f07c100f4cd7b))
* ignore symlink errors ([c44b71c](https://gitea.com/gitea/runner/commits/c44b71c9962f527e9c8847a7a3b959c71ac41b21))
* interpolate action input defaults ([#1376](https://gitea.com/gitea/runner/issues/1376)) ([3631013](https://gitea.com/gitea/runner/commits/36310139e07b29768700e312c92ea49635f0d1ac))
* interpolate env for step actions ([#872](https://gitea.com/gitea/runner/issues/872)) ([099af5e](https://gitea.com/gitea/runner/commits/099af5e6a8e882cab7aa94e11d0616faceb0951f))
* invoke login shells to source `/etc/environment` ([#570](https://gitea.com/gitea/runner/issues/570)) ([b438b83](https://gitea.com/gitea/runner/commits/b438b836ea5e031c1e48cb8b2a777f1647d70f8a))
* keep action inputs ([#1215](https://gitea.com/gitea/runner/issues/1215)) ([1d4c2aa](https://gitea.com/gitea/runner/commits/1d4c2aaa3fb1fce7a0fc70a3923391f20ecd6792))
* keep path to event json file in composite actions ([#1428](https://gitea.com/gitea/runner/issues/1428)) ([2614b3e](https://gitea.com/gitea/runner/commits/2614b3eb0c62772b17ba11509287fdc6f843f9d3))
* limit workers ([5781e23](https://gitea.com/gitea/runner/commits/5781e233c14c6ff2c953be0727b1a9ed7977f4dc))
* localcheckout mock ([#1198](https://gitea.com/gitea/runner/issues/1198)) ([c3fb686](https://gitea.com/gitea/runner/commits/c3fb6864e86542f5ff1f025bec5eff4f81497c96))
* **mergify:** needs-work label handling ([#810](https://gitea.com/gitea/runner/issues/810)) ([a6aea44](https://gitea.com/gitea/runner/commits/a6aea44fb023f33543e58a5fcdcff937cf0e8ec6))
* missing defaults in nodejs pre ([#1349](https://gitea.com/gitea/runner/issues/1349)) ([6037182](https://gitea.com/gitea/runner/commits/6037182487e39d9f50c6887ab3453986fbf29a9b))
* mount volume for workflow info ([#776](https://gitea.com/gitea/runner/issues/776)) ([3c14548](https://gitea.com/gitea/runner/commits/3c1454825d8d4c622a4a94a9b3f796461d06409c))
* need more time to run task ([cce60d4](https://gitea.com/gitea/runner/commits/cce60d466bb189ee83c92659aefb6354e34b9c56))
* nil pointer access ( workflow_dispatch ) ([d97481d](https://gitea.com/gitea/runner/commits/d97481d5671a4906b8b9fe7b951333fa938ffe2b))
* outputs espace ([#404](https://gitea.com/gitea/runner/issues/404)) ([92067eb](https://gitea.com/gitea/runner/commits/92067eb1a5f271278373d6b2c2df9b6f9c373446))
* panic if a step in a job is nil ([#1145](https://gitea.com/gitea/runner/issues/1145)) ([91fd412](https://gitea.com/gitea/runner/commits/91fd412c51667dd10c434e9eb16ba10685581fa6))
* panic on regex in main.workflow ([a04141c](https://gitea.com/gitea/runner/commits/a04141c4440be5c0f003d5441e4e2a97bc09060e)), closes [#12](https://gitea.com/gitea/runner/issues/12)
* panic while running DRYRUN with local actions ([#1141](https://gitea.com/gitea/runner/issues/1141)) ([562265b](https://gitea.com/gitea/runner/commits/562265bb295f618123f65af2b9d5877638b2bda7))
* panic: reflect: slice index out of range ([#1066](https://gitea.com/gitea/runner/issues/1066)) ([a970145](https://gitea.com/gitea/runner/commits/a970145c95b9cc072905ca4f64c0640736aeff1f))
* pass `sc.Env` through `common.Logger` to hide secrets ([#578](https://gitea.com/gitea/runner/issues/578)) ([f3b53d8](https://gitea.com/gitea/runner/commits/f3b53d8eca1142fb587d984d8c95584e236b612a))
* **poller:** graceful shutdown ([9830f34](https://gitea.com/gitea/runner/commits/9830f34d364aa289e98b03f6d501a5afbb8379eb))
* **poller:** initial ready channel ([cab5699](https://gitea.com/gitea/runner/commits/cab56996b7b2f84723e57b0c2f0d80751570fc80))
* preserve job result state in case of failure ([#1519](https://gitea.com/gitea/runner/issues/1519)) ([b14398e](https://gitea.com/gitea/runner/commits/b14398eff3d3a2692c976e250c4b91fbbeed00a5)), closes [#1518](https://gitea.com/gitea/runner/issues/1518)
* print line after password prompt ([#1014](https://gitea.com/gitea/runner/issues/1014)) ([ad06f5d](https://gitea.com/gitea/runner/commits/ad06f5dfb8a658521c855f8c374a58b02227e8fa))
* processing of strategy.matrix.include ([#1200](https://gitea.com/gitea/runner/issues/1200)) ([c30bc82](https://gitea.com/gitea/runner/commits/c30bc824b23335f2c99a32578f8d907f0be6647e))
* re-evaluate env for remote composite actions ([#1385](https://gitea.com/gitea/runner/issues/1385)) ([f0b1845](https://gitea.com/gitea/runner/commits/f0b1845802ca6bea4a75f4e3b514fad85f070b4f))
* regression run after failure ([#971](https://gitea.com/gitea/runner/issues/971)) ([7dbf3fc](https://gitea.com/gitea/runner/commits/7dbf3fcb96e2d70b3436b718c1f4a9fe4b169a4f))
* remove composite restrictions ([#1128](https://gitea.com/gitea/runner/issues/1128)) ([ebb408f](https://gitea.com/gitea/runner/commits/ebb408f373a4fc6508d1d75aa6fddcba5fc7df81))
* remove hardcoded reference to github.com when using reusable remote workflows and remote actions ([#1784](https://gitea.com/gitea/runner/issues/1784)) ([4721abf](https://gitea.com/gitea/runner/commits/4721abfa6dbc661922d34d20b9f60aa14eb30b34))
* remove HOME env var ([#594](https://gitea.com/gitea/runner/issues/594)) ([5044ec6](https://gitea.com/gitea/runner/commits/5044ec6c4364cb74e623953b8c5d39790c8af126))
* remove unused TestContext ([#1013](https://gitea.com/gitea/runner/issues/1013)) ([1e92d87](https://gitea.com/gitea/runner/commits/1e92d87b75f1691a531c76e941e614b521e308f2))
* Remove volume with same name as the container again ([#1203](https://gitea.com/gitea/runner/issues/1203)) ([28f5b37](https://gitea.com/gitea/runner/commits/28f5b37fd1941e00e77f33e3964011f83789c45c))
* remove Windows issue from `README.md` ([#675](https://gitea.com/gitea/runner/issues/675)) ([cc4e23d](https://gitea.com/gitea/runner/commits/cc4e23d96c24effff203c86d679c22c9ab092f27))
* replace `gopass` lib with `golang.org/x/term` ([#674](https://gitea.com/gitea/runner/issues/674)) ([a6a865e](https://gitea.com/gitea/runner/commits/a6a865e973d7138c70e9699d2e7a9a11b542f16c))
* report job failure on error return ([#14](https://gitea.com/gitea/runner/issues/14)) ([f7a5278](https://gitea.com/gitea/runner/commits/f7a52789d950a57cf0c985d02715714740426e5c))
* restore the job logger setup after job cancelation ([#1365](https://gitea.com/gitea/runner/issues/1365)) ([13d5303](https://gitea.com/gitea/runner/commits/13d530302e5f8203b0f104155eb6e136f5aaf910))
* return err in walk ([#1667](https://gitea.com/gitea/runner/issues/1667)) ([ac5dd8f](https://gitea.com/gitea/runner/commits/ac5dd8feb876d37ae483376a137c57383577dace))
* return error if both `run:`/`uses:` keys are used ([#593](https://gitea.com/gitea/runner/issues/593)) ([25779af](https://gitea.com/gitea/runner/commits/25779af4bf1e7762d8ff2d3be985c2fcde7bb1e2))
* return error on reusable workflows ([#1093](https://gitea.com/gitea/runner/issues/1093)) ([b3bd268](https://gitea.com/gitea/runner/commits/b3bd268e4585b32673635685d68bf86c5b2d0bfa))
* return invalid step type ([#1157](https://gitea.com/gitea/runner/issues/1157)) ([1e72c59](https://gitea.com/gitea/runner/commits/1e72c594aca1854854d40a24f3e038537f392fc8))
* reusable workflow panic ([#1728](https://gitea.com/gitea/runner/issues/1728)) ([d70b225](https://gitea.com/gitea/runner/commits/d70b225e85b44efdcd97452959a3bbabec5eb6ed))
* rework setupShellCommand ([#930](https://gitea.com/gitea/runner/issues/930)) ([adabf2a](https://gitea.com/gitea/runner/commits/adabf2a2027fa9da3410febc474d55327e2abeb7))
* reworked container architecture ([#619](https://gitea.com/gitea/runner/issues/619)) ([f32babb](https://gitea.com/gitea/runner/commits/f32babb51db6761fd5dff1bcc03bbdf9eb7aa354))
* rewrite how image env is merged ([#828](https://gitea.com/gitea/runner/issues/828)) ([6c60af7](https://gitea.com/gitea/runner/commits/6c60af76779b739b7487e3fbe83ea15cd2ab8ad5))
* rewrite StepContext env ([#832](https://gitea.com/gitea/runner/issues/832)) ([a0d3602](https://gitea.com/gitea/runner/commits/a0d360236ef5466a441631c2f93794b55d4b76c5))
* **runner:** check task state field exist. ([2442cdd](https://gitea.com/gitea/runner/commits/2442cdd8adbd002ebadd982fd303879cbc73454b))
* safe file name ([#1651](https://gitea.com/gitea/runner/issues/1651)) ([19abab6](https://gitea.com/gitea/runner/commits/19abab637549dae95d4ea3becd77852e2787195b))
* set `sh` as default shell for containers  ([#853](https://gitea.com/gitea/runner/issues/853)) ([e9ae789](https://gitea.com/gitea/runner/commits/e9ae7894e38ed276c41e6cd7b958a81da79249ad))
* set composite outputs on failure ([#945](https://gitea.com/gitea/runner/issues/945)) ([4be9062](https://gitea.com/gitea/runner/commits/4be9062dd2745c1657ad7b98683397a393159436))
* set default version string ([#1070](https://gitea.com/gitea/runner/issues/1070)) ([de81979](https://gitea.com/gitea/runner/commits/de8197980079129395c67b1e2865cff8be4ebca8))
* set logger with trace level ([97629ae](https://gitea.com/gitea/runner/commits/97629ae8af8dbe76387c34e2b053b3645ceb1bde))
* set more image ([64c9856](https://gitea.com/gitea/runner/commits/64c9856ba7b05fe4e9cd77173a9b96332f7e0fec))
* set NoSkipCheckout ([2e1326b](https://gitea.com/gitea/runner/commits/2e1326b725bda03923e4f9f6d6b24a025efc5e3c))
* show workflow info even if on.push is not defined ([#1329](https://gitea.com/gitea/runner/issues/1329)) ([#1335](https://gitea.com/gitea/runner/issues/1335)) ([21484b5](https://gitea.com/gitea/runner/commits/21484b5c1ecbd9e1ddb6a57debd7731a3e1fc25f)), closes [#1332](https://gitea.com/gitea/runner/issues/1332) [#1318](https://gitea.com/gitea/runner/issues/1318)
* show workflow name when get error on yml parsing ([#1133](https://gitea.com/gitea/runner/issues/1133)) ([#1328](https://gitea.com/gitea/runner/issues/1328)) ([b514649](https://gitea.com/gitea/runner/commits/b514649c3da64a1be92b8e22b50e296a07d759f6))
* skip local actions pre step in any case ([#1204](https://gitea.com/gitea/runner/issues/1204)) ([8a47394](https://gitea.com/gitea/runner/commits/8a473943c3dffbbf47c2a2ef8b261e845646fa6e)), closes [#1193](https://gitea.com/gitea/runner/issues/1193)
* spelling mistake in readme ([#1854](https://gitea.com/gitea/runner/issues/1854)) ([80f6de5](https://gitea.com/gitea/runner/commits/80f6de51d544e2c06d730640d683814e247ec0c2))
* step env is unavailable in with property expr ([#1458](https://gitea.com/gitea/runner/issues/1458)) ([8c5748a](https://gitea.com/gitea/runner/commits/8c5748a55c6190ba34354ca29fdb473e5d3724a1))
* support docker create arguments from container.options ([#1022](https://gitea.com/gitea/runner/issues/1022)) ([#1351](https://gitea.com/gitea/runner/issues/1351)) ([48188a6](https://gitea.com/gitea/runner/commits/48188a6280f367e63936a239976be65317ebabd2))
* support expression for step's continue-on-error field ([#900](https://gitea.com/gitea/runner/issues/900)) ([#1331](https://gitea.com/gitea/runner/issues/1331)) ([d1daf2f](https://gitea.com/gitea/runner/commits/d1daf2f28d8c2256e5712e3064ce6ac45fca2562))
* tail (not absolute) as entrypoint of job container ([#1506](https://gitea.com/gitea/runner/issues/1506)) ([0c8c082](https://gitea.com/gitea/runner/commits/0c8c082ac0baf25909fd79b87897426816c3c4eb))
* **test:** adapt to sjw change ([ae39a4b](https://gitea.com/gitea/runner/commits/ae39a4b1d354e13891ef8bd7eb8bcd4a0916f0e6))
* **test:** needs condition. ([#8](https://gitea.com/gitea/runner/issues/8)) ([0efa2d5](https://gitea.com/gitea/runner/commits/0efa2d5e63fdc03360b4eb6aac0b9ac5eac9ca42))
* the number in the github event is of type number ([#1252](https://gitea.com/gitea/runner/issues/1252)) ([4094462](https://gitea.com/gitea/runner/commits/409446211fe99420d5d3761b04f491f8ef52412d))
* tojson(null) == 'null' ([#1068](https://gitea.com/gitea/runner/issues/1068)) ([5083284](https://gitea.com/gitea/runner/commits/5083284deb804db5b9d5997ab4cf345b1c7ec069))
* tolerate workflow that needs a missing job ([#1595](https://gitea.com/gitea/runner/issues/1595)) ([#1619](https://gitea.com/gitea/runner/issues/1619)) ([be4a147](https://gitea.com/gitea/runner/commits/be4a1477a5b73b3e17f79891ec6147413acb4df7))
* trim whitespace from actrc ([#1024](https://gitea.com/gitea/runner/issues/1024)) ([3db3c73](https://gitea.com/gitea/runner/commits/3db3c737230adc457e08a8d60b34458977d01e5c)), closes [#985](https://gitea.com/gitea/runner/issues/985)
* try finding a socket, otherwise fail, respect user choice ([#1745](https://gitea.com/gitea/runner/issues/1745)) ([b6718fd](https://gitea.com/gitea/runner/commits/b6718fdf5d2fdc3e61798f4ac749b53548cbbfae))
* **typo:** Use uppercase N for default input prompt ([#9](https://gitea.com/gitea/runner/issues/9)) ([9ba9604](https://gitea.com/gitea/runner/commits/9ba9604256964f52ff102b16acf782d8c56d87ff))
* update `ubuntu-latest` to point to `ubuntu-20.04` ([#580](https://gitea.com/gitea/runner/issues/580)) ([72f2d2d](https://gitea.com/gitea/runner/commits/72f2d2de51fcc309a2f0381a9b5098c81802598b))
* Update ARCH environment variable used in runners/actions ([#1818](https://gitea.com/gitea/runner/issues/1818)) ([11dd2ac](https://gitea.com/gitea/runner/commits/11dd2ac7457bd22080b9d8e5277327110563919a))
* update artifact server to address GHSL-2023-004 ([#1565](https://gitea.com/gitea/runner/issues/1565)) ([63ae215](https://gitea.com/gitea/runner/commits/63ae215071f94569d910964bdee866d91d6e3a10))
* update codecov action ([#727](https://gitea.com/gitea/runner/issues/727)) ([05bce00](https://gitea.com/gitea/runner/commits/05bce00e932e450ec3cd0d73c064e2a396ac9500))
* update codecov action (again) ([#731](https://gitea.com/gitea/runner/issues/731)) ([6550eb7](https://gitea.com/gitea/runner/commits/6550eb7f2515b635c1210fc218e03132008463c6))
* update container config ([b3e7d6a](https://gitea.com/gitea/runner/commits/b3e7d6a62d6c00850616fddc199bef6bb65e8a85))
* update install.sh ([#937](https://gitea.com/gitea/runner/issues/937)) ([e9c4609](https://gitea.com/gitea/runner/commits/e9c4609dcadf70739e34f79680e69df973b24d8e))
* update issue template ([#722](https://gitea.com/gitea/runner/issues/722)) ([fa93f4d](https://gitea.com/gitea/runner/commits/fa93f4d5e7e9ccfe68741a555e3137c1dfd6e997))
* update output handling for reusable workflows ([#1521](https://gitea.com/gitea/runner/issues/1521)) ([89cb558](https://gitea.com/gitea/runner/commits/89cb55855802fd7997bf9128140ee4e48b7e89c1))
* update readme with more relevant/updated information ([#630](https://gitea.com/gitea/runner/issues/630)) ([08141f5](https://gitea.com/gitea/runner/commits/08141f5b0095f0ddbbf3a62b03338a9570e27658))
* update step result ([41b5fa6](https://gitea.com/gitea/runner/commits/41b5fa6b557d5cb937690635306aea6749ee6c91))
* use auth when fetching (same as with cloning) ([#687](https://gitea.com/gitea/runner/issues/687)) ([38f6dfb](https://gitea.com/gitea/runner/commits/38f6dfb49a7e3cee826911e36f0c990813282d28))
* use docker images from dockerhub ([#1249](https://gitea.com/gitea/runner/issues/1249)) ([e70b968](https://gitea.com/gitea/runner/commits/e70b9689241785df1ff509678d99f1fdcb138d91))
* use docker lib for image ref parsing ([#877](https://gitea.com/gitea/runner/issues/877)) ([ec34eb9](https://gitea.com/gitea/runner/commits/ec34eb9532d314ba8d3c23472acf4037f2c084cd))
* use github auth token for pulling a repository ([#693](https://gitea.com/gitea/runner/issues/693)) ([60b9606](https://gitea.com/gitea/runner/commits/60b9606fc687549c7f55343a4d52ff3f2c32cf4d))
* use new platforms after survey, check working dir for `.actrc` ([#577](https://gitea.com/gitea/runner/issues/577)) ([780a8a0](https://gitea.com/gitea/runner/commits/780a8a061c2c7a3a6958e444e105d7d2441893a9))
* use os.UserHomeDir ([#1706](https://gitea.com/gitea/runner/issues/1706)) ([d3dfde0](https://gitea.com/gitea/runner/commits/d3dfde055aad9c0e9631fd6bbf06f4b0db706eff))
* using filepath instead of path ([01f7ff3](https://gitea.com/gitea/runner/commits/01f7ff315c1b3eaa5e3a4eee06bd6d88ebf60ce4))
* warn Apple M1 (arm64) users ([#743](https://gitea.com/gitea/runner/issues/743)) ([1cf422e](https://gitea.com/gitea/runner/commits/1cf422e41108bf72ec837bea66e080d400ca1192))
* wrong perms ([#726](https://gitea.com/gitea/runner/issues/726)) ([2261204](https://gitea.com/gitea/runner/commits/2261204c65b71d661495de95df6610c451157a67))

### Refactoring

* `NewWorkflowPlanner` ([#648](https://gitea.com/gitea/runner/issues/648)) ([806bc4d](https://gitea.com/gitea/runner/commits/806bc4d999cafab455bbda9b966e3ed6080a4104))
* Add graceful shutdown signal notify func ([e54de37](https://gitea.com/gitea/runner/commits/e54de3724292fa3fc4a50f1e1186964ce1e7a816))
* default empty if: statement to if: success() in evaluator & remove `FixIfStatement()` ([#990](https://gitea.com/gitea/runner/issues/990)) ([ff13844](https://gitea.com/gitea/runner/commits/ff13844b8643abaff9daa47a5bc96aa91fddcf5e))
* docker build BuildContext field ([#1914](https://gitea.com/gitea/runner/issues/1914)) ([f0ca0ab](https://gitea.com/gitea/runner/commits/f0ca0abc4003e1d99a534564fc1ad62f47eaf8a4))
* docker engine with options ([14a3455](https://gitea.com/gitea/runner/commits/14a345504f03c02a70b9323d157deccab236d86a))
* export and move shared contexts into pkg/model ([#931](https://gitea.com/gitea/runner/issues/931)) ([ed01f46](https://gitea.com/gitea/runner/commits/ed01f464ede3515c94dedf5739b54576177bf80c)), references [#908](https://gitea.com/gitea/runner/issues/908)
* extract RunContext Executor in JobExecutor ([#984](https://gitea.com/gitea/runner/issues/984)) ([e23223a](https://gitea.com/gitea/runner/commits/e23223ad02ddeafe581df97c7e9462df28d8e104))
* extract setupAction into ActionReader ([#986](https://gitea.com/gitea/runner/issues/986)) ([e4f0080](https://gitea.com/gitea/runner/commits/e4f0080a1efc9887fe8aa892669b2002e70b6040))
* fix add-path / GITHUB_PATH commands ([#1472](https://gitea.com/gitea/runner/issues/1472)) ([d281230](https://gitea.com/gitea/runner/commits/d281230ccef6438085180866e50a026025a05f67))
* fix savestate in pre steps ([#1466](https://gitea.com/gitea/runner/issues/1466)) ([b8d7e94](https://gitea.com/gitea/runner/commits/b8d7e947cf91e24202326b63821a3ff82a3e2b57))
* GITHUB_ENV command / remove env.PATH ([#1503](https://gitea.com/gitea/runner/issues/1503)) ([e775fea](https://gitea.com/gitea/runner/commits/e775fea2655d9e3342bd06c2743cb4bab791e2ab))
* ignore already closed error ([#1285](https://gitea.com/gitea/runner/issues/1285)) ([92ddcda](https://gitea.com/gitea/runner/commits/92ddcdae098873ad6d9c05748ce5b6c7be83b811))
* move autoremove into the jobexecutor ([#1463](https://gitea.com/gitea/runner/issues/1463)) ([d9fe63e](https://gitea.com/gitea/runner/commits/d9fe63ec249d4ee3468151851b8db642dc32bf2c))
* move from io/ioutil to io and os packages ([#1417](https://gitea.com/gitea/runner/issues/1417)) ([2f1c5a1](https://gitea.com/gitea/runner/commits/2f1c5a19f12ad0fc784924b4585d4ace061c7ef5))
* open boltdb only while using it ([#1879](https://gitea.com/gitea/runner/issues/1879)) ([e597046](https://gitea.com/gitea/runner/commits/e597046195774991a75ff66506fdfcc3fc903364))
* pull and rebuild docker by default ([#1569](https://gitea.com/gitea/runner/issues/1569)) ([b00babd](https://gitea.com/gitea/runner/commits/b00babd0d9fe8538aced303346cc7139ef51d0cd))
* re-implement `embed` without "unused" import ([#830](https://gitea.com/gitea/runner/issues/830)) ([7a426a0](https://gitea.com/gitea/runner/commits/7a426a0f37d21e32a01e85d9403edba17a54d545))
* remove `gotest.tools` ([#688](https://gitea.com/gitea/runner/issues/688)) ([f571290](https://gitea.com/gitea/runner/commits/f571290b25cdbb34e863cb98eb155369fff5d572))
* remove composite action runcontext workaround ([#1085](https://gitea.com/gitea/runner/issues/1085)) ([e360811](https://gitea.com/gitea/runner/commits/e3608115706382233c3b943b4dabea09bab9d7d2))
* remove docker image list reference filter ([#1501](https://gitea.com/gitea/runner/issues/1501)) ([b2fb9e6](https://gitea.com/gitea/runner/commits/b2fb9e64ac1b263088003b94963b36b6edeaba31))
* remove github.com/pkg/errors dependency ([#1077](https://gitea.com/gitea/runner/issues/1077)) ([2aa0699](https://gitea.com/gitea/runner/commits/2aa0699aecff2e596e280165ea48eac85c5f9a6a))
* remove unnecessary nil check in RunContext ([#1955](https://gitea.com/gitea/runner/issues/1955)) ([8a9e4f9](https://gitea.com/gitea/runner/commits/8a9e4f9f38340ff9b7f9727add3dae62f3f6c7bb))
* return more errors, add more tests ([#679](https://gitea.com/gitea/runner/issues/679)) ([aa68181](https://gitea.com/gitea/runner/commits/aa68181f6b3df2ff24a523800871c3c5ccdf76be))
* share UpdateFromEnv logic ([#1457](https://gitea.com/gitea/runner/issues/1457)) ([57bf4d2](https://gitea.com/gitea/runner/commits/57bf4d27a2b6d03a2c61f39ba57c24138ae94370))
* simplify action function signatures ([#1083](https://gitea.com/gitea/runner/issues/1083)) ([5f673cb](https://gitea.com/gitea/runner/commits/5f673cbb6db5569334cc806cf0536a97ad6b899b))
* **task:** execute single task with gRPC data ([20c3d85](https://gitea.com/gitea/runner/commits/20c3d85ba9caef27fbf83118ce4b6126fed97580))

### Build

* **deps:** bump actions/cache from 2 to 3 ([#1082](https://gitea.com/gitea/runner/issues/1082)) ([4f0bfd2](https://gitea.com/gitea/runner/commits/4f0bfd2ad1b16dc0a5f34844069513af178fe86d))
* **deps:** bump actions/checkout from 2 to 3 ([#1028](https://gitea.com/gitea/runner/issues/1028)) ([0579f7c](https://gitea.com/gitea/runner/commits/0579f7c23598fd0b76500b95e5d69d55a3fb0b98))
* **deps:** bump actions/github-script from 5 to 6 ([#1572](https://gitea.com/gitea/runner/issues/1572)) ([75d19d0](https://gitea.com/gitea/runner/commits/75d19d0ff58e54cafa924b9bf72a89176a602f2b))
* **deps:** bump actions/setup-go from 2 to 3 ([#1120](https://gitea.com/gitea/runner/issues/1120)) ([4e1cf1d](https://gitea.com/gitea/runner/commits/4e1cf1d7cc84a6da069401ba1ff08d14cf71044b))
* **deps:** bump actions/setup-go from 3 to 4 ([#1689](https://gitea.com/gitea/runner/issues/1689)) ([8c28c9f](https://gitea.com/gitea/runner/commits/8c28c9fd8f2ae53af23367be2cda094ecc4fbfe0))
* **deps:** bump actions/stale from 4 to 5 ([#1117](https://gitea.com/gitea/runner/issues/1117)) ([f551165](https://gitea.com/gitea/runner/commits/f551165970914f2541f692716c522e95d06a0608))
* **deps:** bump actions/stale from 5 to 6 ([#1357](https://gitea.com/gitea/runner/issues/1357)) ([9a494ca](https://gitea.com/gitea/runner/commits/9a494cac1ed5f22539d727f6958f4d24df7104b8))
* **deps:** bump actions/stale from 6 to 7 ([#1535](https://gitea.com/gitea/runner/issues/1535)) ([d3ce2b0](https://gitea.com/gitea/runner/commits/d3ce2b0a467b39616059b19c5ba0e69e145099aa))
* **deps:** bump actions/stale from 7 to 8 ([#1700](https://gitea.com/gitea/runner/issues/1700)) ([d65cf86](https://gitea.com/gitea/runner/commits/d65cf869e6824ac8aca448df8d7e0bc8b8fc21af))
* **deps:** bump actions/upload-artifact from 2 to 3 ([#1119](https://gitea.com/gitea/runner/issues/1119)) ([39cc9b0](https://gitea.com/gitea/runner/commits/39cc9b0840f07953b44416b62d4817ca375fe932))
* **deps:** bump codecov/codecov-action from 1.5.2 to 2.1.0 ([#819](https://gitea.com/gitea/runner/issues/819)) ([e82f17e](https://gitea.com/gitea/runner/commits/e82f17e230a7ae5b715b041b6827825de671dc9e))
* **deps:** bump codecov/codecov-action from 2.1.0 to 3.0.0 ([#1118](https://gitea.com/gitea/runner/issues/1118)) ([74cda72](https://gitea.com/gitea/runner/commits/74cda72c2672a1f665d5a55416cd93ca6292180a))
* **deps:** bump codecov/codecov-action from 3.0.0 to 3.1.0 ([#1139](https://gitea.com/gitea/runner/issues/1139)) ([3f2469b](https://gitea.com/gitea/runner/commits/3f2469bf5385b90da5f6a5d9922ccf9c9d7274a6))
* **deps:** bump codecov/codecov-action from 3.1.0 to 3.1.1 ([#1359](https://gitea.com/gitea/runner/issues/1359)) ([d137ec8](https://gitea.com/gitea/runner/commits/d137ec82852083eef45009dcaf368c3dcdd0b9a6))
* **deps:** bump codecov/codecov-action from 3.1.1 to 3.1.2 ([#1735](https://gitea.com/gitea/runner/issues/1735)) ([6d527bf](https://gitea.com/gitea/runner/commits/6d527bf1a356b69ee4ee35e43159c9d0a21c7718))
* **deps:** bump codecov/codecov-action from 3.1.2 to 3.1.3 ([#1752](https://gitea.com/gitea/runner/issues/1752)) ([60bcfb5](https://gitea.com/gitea/runner/commits/60bcfb5d4f717bd0a581de8e97866d58a5593152))
* **deps:** bump codecov/codecov-action from 3.1.3 to 3.1.4 ([#1813](https://gitea.com/gitea/runner/issues/1813)) ([1371215](https://gitea.com/gitea/runner/commits/1371215aabfd82b9288132042aecc9eacc0ab380))
* **deps:** bump deps which dependabot doesn't want to ([#1107](https://gitea.com/gitea/runner/issues/1107)) ([afb23a2](https://gitea.com/gitea/runner/commits/afb23a289187f4a90241b8b923555b7c3d844ec9))
* **deps:** bump docker/setup-qemu-action from 1 to 2 ([#1156](https://gitea.com/gitea/runner/issues/1156)) ([7974a17](https://gitea.com/gitea/runner/commits/7974a17afc613431461f9a8c629597962b1f54b4))
* **deps:** bump docker/setup-qemu-action from 2 to 3 ([#2007](https://gitea.com/gitea/runner/issues/2007)) ([1891bef](https://gitea.com/gitea/runner/commits/1891bef4335ba55573de1c7c5824acd44b3a4994))
* **deps:** bump fregante/setup-git-user from 1 to 2 ([#1664](https://gitea.com/gitea/runner/issues/1664)) ([fad986a](https://gitea.com/gitea/runner/commits/fad986af5f95a2878c9d9f15ded8d815fdf5e061))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.2.7 to 2.3.1 ([#803](https://gitea.com/gitea/runner/issues/803)) ([a7a2659](https://gitea.com/gitea/runner/commits/a7a2659c0e90a53fee4e55817d0eedcad1ab9967))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.2 to 2.3.4 ([#1079](https://gitea.com/gitea/runner/issues/1079)) ([efd34f1](https://gitea.com/gitea/runner/commits/efd34f1e10e56b251a5b9eb75f269d2e7c15681f))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.4 to 2.3.5 ([#1214](https://gitea.com/gitea/runner/issues/1214)) ([084c4e6](https://gitea.com/gitea/runner/commits/084c4e61d99c1bfd109ea5bc381342d94a6e65cd))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.5 to 2.3.6 ([#1341](https://gitea.com/gitea/runner/issues/1341)) ([3a0fe69](https://gitea.com/gitea/runner/commits/3a0fe6967fd8ecd3cb86550d0225d55a0bb37ac9))
* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.3.6 to 2.3.7 ([#1872](https://gitea.com/gitea/runner/issues/1872)) ([df44dff](https://gitea.com/gitea/runner/commits/df44dffd30b6503b91bc63130742c31b1ce7badf))
* **deps:** bump github.com/cloudflare/circl from 1.1.0 to 1.3.3 ([#1800](https://gitea.com/gitea/runner/issues/1800)) ([046f5aa](https://gitea.com/gitea/runner/commits/046f5aa71587ba5e81a3f755c913da2f0644cdfb))
* **deps:** bump github.com/containerd/containerd from 1.6.16 to 1.6.18 ([#1637](https://gitea.com/gitea/runner/issues/1637)) ([4eba04b](https://gitea.com/gitea/runner/commits/4eba04b229388c4fc8bbbcffb5e4d62cd927eac9))
* **deps:** bump github.com/containerd/containerd from 1.6.6 to 1.6.12 ([#1554](https://gitea.com/gitea/runner/issues/1554)) ([3ac756b](https://gitea.com/gitea/runner/commits/3ac756b37a496020ec8bb7b5852272d01c471057))
* **deps:** bump github.com/creack/pty from 1.1.17 to 1.1.18 ([#1447](https://gitea.com/gitea/runner/issues/1447)) ([ccb3b0e](https://gitea.com/gitea/runner/commits/ccb3b0e8757158bb945cf9d3b12a4a9299eea9c4))
* **deps:** bump github.com/docker/cli ([#1038](https://gitea.com/gitea/runner/issues/1038)) ([933413e](https://gitea.com/gitea/runner/commits/933413e2c0e1d600e09291ba6e63211a3c294385))
* **deps:** bump github.com/docker/cli ([#1081](https://gitea.com/gitea/runner/issues/1081)) ([34706f4](https://gitea.com/gitea/runner/commits/34706f49fe8c82b3a17c9b22bb4fc5d8a70bd25e))
* **deps:** bump github.com/docker/cli ([#1154](https://gitea.com/gitea/runner/issues/1154)) ([d5e9721](https://gitea.com/gitea/runner/commits/d5e972159f9ee43b9229bf028853755c7f23c4d6))
* **deps:** bump github.com/docker/cli ([#1168](https://gitea.com/gitea/runner/issues/1168)) ([1fed666](https://gitea.com/gitea/runner/commits/1fed666165079c8fcab76ff11b66a60463ee0f1f))
* **deps:** bump github.com/docker/cli ([#1211](https://gitea.com/gitea/runner/issues/1211)) ([87d9687](https://gitea.com/gitea/runner/commits/87d9687716a279525fd952b80ba6f044bce88492))
* **deps:** bump github.com/docker/cli ([#1340](https://gitea.com/gitea/runner/issues/1340)) ([ceac1f2](https://gitea.com/gitea/runner/commits/ceac1f2429073270e984ce7da39e28f31e2c56c6))
* **deps:** bump github.com/docker/cli ([#1393](https://gitea.com/gitea/runner/issues/1393)) ([bd5e0d2](https://gitea.com/gitea/runner/commits/bd5e0d24d68cddab105ee76bc2fb8000dce76a11))
* **deps:** bump github.com/docker/cli ([#1405](https://gitea.com/gitea/runner/issues/1405)) ([ca69531](https://gitea.com/gitea/runner/commits/ca69531096e88bec7f8fbd6aecef9b3769271cb7))
* **deps:** bump github.com/docker/cli ([#1419](https://gitea.com/gitea/runner/issues/1419)) ([ac6571c](https://gitea.com/gitea/runner/commits/ac6571c8af87ca3e26b4a5eed1d0a919b7a7aed3))
* **deps:** bump github.com/docker/cli ([#1516](https://gitea.com/gitea/runner/issues/1516)) ([7679695](https://gitea.com/gitea/runner/commits/767969502aa0cc95e65916981a950250ddb6c71d))
* **deps:** bump github.com/docker/cli ([#1626](https://gitea.com/gitea/runner/issues/1626)) ([98ad62f](https://gitea.com/gitea/runner/commits/98ad62fcefc3b760e69040b25d15d34379376319))
* **deps:** bump github.com/docker/cli ([#1712](https://gitea.com/gitea/runner/issues/1712)) ([24b04df](https://gitea.com/gitea/runner/commits/24b04dfa55ebd31c7b4f2fb114350d14783a533a))
* **deps:** bump github.com/docker/cli ([#1722](https://gitea.com/gitea/runner/issues/1722)) ([f2bd194](https://gitea.com/gitea/runner/commits/f2bd194c7fa282fb35f31bcafe9bf1875da71fc7))
* **deps:** bump github.com/docker/cli ([#1753](https://gitea.com/gitea/runner/issues/1753)) ([ded31bd](https://gitea.com/gitea/runner/commits/ded31bdbc0c82ddb737535fab0fc4d1800ae2b00))
* **deps:** bump github.com/docker/cli ([#1780](https://gitea.com/gitea/runner/issues/1780)) ([3eb6e83](https://gitea.com/gitea/runner/commits/3eb6e83ea4c589a0dfcb2a17e93a04c23909f770))
* **deps:** bump github.com/docker/cli ([#1805](https://gitea.com/gitea/runner/issues/1805)) ([39abbce](https://gitea.com/gitea/runner/commits/39abbce4e6cdf3988a689a404f2f0082ce0125bf))
* **deps:** bump github.com/docker/cli ([#1816](https://gitea.com/gitea/runner/issues/1816)) ([9c1f1f8](https://gitea.com/gitea/runner/commits/9c1f1f8d844148079cf80c49413d7f3357c42025))
* **deps:** bump github.com/docker/cli ([#1843](https://gitea.com/gitea/runner/issues/1843)) ([e48c079](https://gitea.com/gitea/runner/commits/e48c07988ed70d94860c6d6276e134f5f51bd802))
* **deps:** bump github.com/docker/cli ([#1932](https://gitea.com/gitea/runner/issues/1932)) ([8e6c7c1](https://gitea.com/gitea/runner/commits/8e6c7c11fe35bd117579277c29238e4e45996a7c))
* **deps:** bump github.com/docker/cli ([#1999](https://gitea.com/gitea/runner/issues/1999)) ([bd9032d](https://gitea.com/gitea/runner/commits/bd9032de0a79e9c0755596b4a1a10f5de8847639))
* **deps:** bump github.com/docker/cli ([#854](https://gitea.com/gitea/runner/issues/854)) ([50dd0c0](https://gitea.com/gitea/runner/commits/50dd0c0db149afd3f92d64b79b6eb8d2d0b0ec16))
* **deps:** bump github.com/docker/cli ([#956](https://gitea.com/gitea/runner/issues/956)) ([bbf6b7e](https://gitea.com/gitea/runner/commits/bbf6b7e080959ebb13d3500fed637e822e2b8ec0))
* **deps:** bump github.com/docker/cli from 23.0.0-rc.1+incompatible to 23.0.0+incompatible ([#1611](https://gitea.com/gitea/runner/issues/1611)) ([36dbbc1](https://gitea.com/gitea/runner/commits/36dbbc1dfa049676e17e1f4712c6b918ce4ce5b7))
* **deps:** bump github.com/docker/distribution ([#1037](https://gitea.com/gitea/runner/issues/1037)) ([5c0023c](https://gitea.com/gitea/runner/commits/5c0023c3466b7d5d8a900a1b5a3798e4d1dabe75))
* **deps:** bump github.com/docker/distribution ([#1801](https://gitea.com/gitea/runner/issues/1801)) ([b5fa245](https://gitea.com/gitea/runner/commits/b5fa24537f70e88166c7e43d45b5f6ff41d0b12d))
* **deps:** bump github.com/docker/distribution ([#994](https://gitea.com/gitea/runner/issues/994)) ([0fae967](https://gitea.com/gitea/runner/commits/0fae96792c77b395587c4d8129d0f680bffad643))
* **deps:** bump github.com/docker/docker ([#1062](https://gitea.com/gitea/runner/issues/1062)) ([87175ec](https://gitea.com/gitea/runner/commits/87175ecb6882f0735c1de8eeb884ced21f978873))
* **deps:** bump github.com/docker/docker ([#1080](https://gitea.com/gitea/runner/issues/1080)) ([c8a40bf](https://gitea.com/gitea/runner/commits/c8a40bf8f41a68efe55f28e8b17a9b5f3d004d5e))
* **deps:** bump github.com/docker/docker ([#1155](https://gitea.com/gitea/runner/issues/1155)) ([01d85c9](https://gitea.com/gitea/runner/commits/01d85c9a0150092afa061ba1518a96cae6dcb73b))
* **deps:** bump github.com/docker/docker ([#1167](https://gitea.com/gitea/runner/issues/1167)) ([018c846](https://gitea.com/gitea/runner/commits/018c846d8d92979e158f0618b9f3d068fecf08f1))
* **deps:** bump github.com/docker/docker ([#1212](https://gitea.com/gitea/runner/issues/1212)) ([63df5dd](https://gitea.com/gitea/runner/commits/63df5ddf429881b88d766919db2f404d1dc35c9c))
* **deps:** bump github.com/docker/docker ([#1339](https://gitea.com/gitea/runner/issues/1339)) ([ca8767f](https://gitea.com/gitea/runner/commits/ca8767f08057b176ea7f6261c3f557bde1817a8c))
* **deps:** bump github.com/docker/docker ([#1392](https://gitea.com/gitea/runner/issues/1392)) ([6ff4491](https://gitea.com/gitea/runner/commits/6ff44912381a15fbc27ce62937bae89682cfaa2f))
* **deps:** bump github.com/docker/docker ([#1403](https://gitea.com/gitea/runner/issues/1403)) ([0d1dccc](https://gitea.com/gitea/runner/commits/0d1dccca3a5d663860e57cc53a50e61bbbc003fa))
* **deps:** bump github.com/docker/docker ([#1418](https://gitea.com/gitea/runner/issues/1418)) ([e03b816](https://gitea.com/gitea/runner/commits/e03b816c5880b60af28f7dda33e4361d30451092))
* **deps:** bump github.com/docker/docker ([#1517](https://gitea.com/gitea/runner/issues/1517)) ([469d024](https://gitea.com/gitea/runner/commits/469d024c1fb6f376c56bcf380911761d9ad829db))
* **deps:** bump github.com/docker/docker ([#1573](https://gitea.com/gitea/runner/issues/1573)) ([612e08f](https://gitea.com/gitea/runner/commits/612e08f113b0d5c97a8c1ed664fdb2b19b1a418e))
* **deps:** bump github.com/docker/docker ([#1613](https://gitea.com/gitea/runner/issues/1613)) ([f41e912](https://gitea.com/gitea/runner/commits/f41e9127a8e4516e5d78d58702d00c55de1a6b96))
* **deps:** bump github.com/docker/docker ([#1624](https://gitea.com/gitea/runner/issues/1624)) ([7c3c534](https://gitea.com/gitea/runner/commits/7c3c5349ab54ac3a2c779b199310b8337f6e8fbf))
* **deps:** bump github.com/docker/docker ([#1711](https://gitea.com/gitea/runner/issues/1711)) ([68c72b9](https://gitea.com/gitea/runner/commits/68c72b9a51f4fa586af62619f9e6361e023a5d48))
* **deps:** bump github.com/docker/docker ([#1719](https://gitea.com/gitea/runner/issues/1719)) ([65a925b](https://gitea.com/gitea/runner/commits/65a925b4ed5f0bebab9314433fcfdac4018482ac))
* **deps:** bump github.com/docker/docker ([#1754](https://gitea.com/gitea/runner/issues/1754)) ([16f35f6](https://gitea.com/gitea/runner/commits/16f35f64eb992d96992c43ab1ad0e2fae5fbc8b5))
* **deps:** bump github.com/docker/docker ([#1779](https://gitea.com/gitea/runner/issues/1779)) ([1f9bbe1](https://gitea.com/gitea/runner/commits/1f9bbe12dd5a86ee6401000bcc654b9b66910bae))
* **deps:** bump github.com/docker/docker ([#1806](https://gitea.com/gitea/runner/issues/1806)) ([c00810b](https://gitea.com/gitea/runner/commits/c00810bdf4ca9e5cdf0c34e7ede5d4ebc344f558))
* **deps:** bump github.com/docker/docker ([#1933](https://gitea.com/gitea/runner/issues/1933)) ([d720ff0](https://gitea.com/gitea/runner/commits/d720ff09a281aae453fe6cc2970e6f833b1da49d))
* **deps:** bump github.com/docker/docker ([#2000](https://gitea.com/gitea/runner/issues/2000)) ([c84a3ef](https://gitea.com/gitea/runner/commits/c84a3ef6d0e4c545f627572aef0c0313bc9d4187))
* **deps:** bump github.com/docker/docker ([#855](https://gitea.com/gitea/runner/issues/855)) ([6a3a47c](https://gitea.com/gitea/runner/commits/6a3a47c2175c5f6438370d28febf8933097cdc55))
* **deps:** bump github.com/docker/docker ([#957](https://gitea.com/gitea/runner/issues/957)) ([b0f5557](https://gitea.com/gitea/runner/commits/b0f55571f092c47e75f746ecb8aeade718db9495))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.0.0 to 5.3.1 ([#802](https://gitea.com/gitea/runner/issues/802)) ([6d7f574](https://gitea.com/gitea/runner/commits/6d7f574859a1a358db7418dbb96fcd07b983eb31))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.3.1 to 5.4.0 ([#1541](https://gitea.com/gitea/runner/issues/1541)) ([9c70851](https://gitea.com/gitea/runner/commits/9c7085195c4fcfe1f37b53dd9a329824a1753b83))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.0 to 5.4.1 ([#1612](https://gitea.com/gitea/runner/issues/1612)) ([22dc1e0](https://gitea.com/gitea/runner/commits/22dc1e0e7ef5291adbc80a50574ad91a07c64754))
* **deps:** bump github.com/go-git/go-billy/v5 from 5.4.1 to 5.5.0 ([#2012](https://gitea.com/gitea/runner/issues/2012)) ([d134079](https://gitea.com/gitea/runner/commits/d134079807fe35a2d8ebcfd24b139f48c8123caa))
* **deps:** bump github.com/go-git/go-git/v5 ([#1830](https://gitea.com/gitea/runner/issues/1830)) ([b0996e0](https://gitea.com/gitea/runner/commits/b0996e057750223a14321ff92c1c6d6fbb8e1a10))
* **deps:** bump github.com/go-git/go-git/v5 from 5.7.0 to 5.8.0 ([#1925](https://gitea.com/gitea/runner/issues/1925)) ([17bf4fc](https://gitea.com/gitea/runner/commits/17bf4fc5afd147ea967445b9d0dc2e684f57377c))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.0 to 5.8.1 ([#1934](https://gitea.com/gitea/runner/issues/1934)) ([8bf10cf](https://gitea.com/gitea/runner/commits/8bf10cf87618e3ebe1c0f77c1add9e3a08b3f514))
* **deps:** bump github.com/go-git/go-git/v5 from 5.8.1 to 5.9.0 ([#2011](https://gitea.com/gitea/runner/issues/2011)) ([b637d79](https://gitea.com/gitea/runner/commits/b637d79ec3a0468e0ff968cdf57420c7020265c8))
* **deps:** bump github.com/go-ini/ini from 1.63.2 to 1.64.0 ([#882](https://gitea.com/gitea/runner/issues/882)) ([c26e661](https://gitea.com/gitea/runner/commits/c26e661889059d3afb2a962f86403e29a8b7c213))
* **deps:** bump github.com/go-ini/ini from 1.64.0 to 1.66.4 ([#997](https://gitea.com/gitea/runner/issues/997)) ([7d43396](https://gitea.com/gitea/runner/commits/7d433968e525098b506aae549d5fc5da4fec571e))
* **deps:** bump github.com/go-ini/ini from 1.66.4 to 1.66.5 ([#1189](https://gitea.com/gitea/runner/issues/1189)) ([64387bc](https://gitea.com/gitea/runner/commits/64387bcf7b92f07b42a4c93f7e4504219af5df2c))
* **deps:** bump github.com/go-ini/ini from 1.66.5 to 1.66.6 ([#1201](https://gitea.com/gitea/runner/issues/1201)) ([859445f](https://gitea.com/gitea/runner/commits/859445fb94e12cc600e5a59d9c9cd687214cb801))
* **deps:** bump github.com/go-ini/ini from 1.66.6 to 1.67.0 ([#1305](https://gitea.com/gitea/runner/issues/1305)) ([12029e3](https://gitea.com/gitea/runner/commits/12029e3df5841538a806293881f836f88ec07a2c))
* **deps:** bump github.com/imdario/mergo from 0.3.12 to 0.3.13 ([#1378](https://gitea.com/gitea/runner/issues/1378)) ([9f9f644](https://gitea.com/gitea/runner/commits/9f9f644d36dc8a7ba73701bdd18cf0d23eaa736b))
* **deps:** bump github.com/imdario/mergo from 0.3.13 to 0.3.14 ([#1688](https://gitea.com/gitea/runner/issues/1688)) ([2fa0a5f](https://gitea.com/gitea/runner/commits/2fa0a5f76979689945215fc392f88f836bc71103))
* **deps:** bump github.com/imdario/mergo from 0.3.14 to 0.3.15 ([#1701](https://gitea.com/gitea/runner/issues/1701)) ([f764ecb](https://gitea.com/gitea/runner/commits/f764ecb283febc36da0fa302a8b8aa8d6d350054))
* **deps:** bump github.com/imdario/mergo from 0.3.15 to 0.3.16 ([#1840](https://gitea.com/gitea/runner/issues/1840)) ([0b05ee8](https://gitea.com/gitea/runner/commits/0b05ee8e0b44454f56baaf23f0f00d7e946b8f37))
* **deps:** bump github.com/joho/godotenv from 1.3.0 to 1.4.0 ([#824](https://gitea.com/gitea/runner/issues/824)) ([bfcf0ab](https://gitea.com/gitea/runner/commits/bfcf0abd73e558521367e36b323afad0dc1f1a71))
* **deps:** bump github.com/joho/godotenv from 1.4.0 to 1.5.1 ([#1614](https://gitea.com/gitea/runner/issues/1614)) ([281a52f](https://gitea.com/gitea/runner/commits/281a52f0d6ed70bd31feaedcbbb46e941cc8bdcc))
* **deps:** bump github.com/mattn/go-isatty from 0.0.14 to 0.0.16 ([#1312](https://gitea.com/gitea/runner/issues/1312)) ([b23bbef](https://gitea.com/gitea/runner/commits/b23bbefc365012886192e477a88d38a0909ecba1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.16 to 0.0.17 ([#1539](https://gitea.com/gitea/runner/issues/1539)) ([0d38857](https://gitea.com/gitea/runner/commits/0d388572a89fabdbd229536b7ecca9db4ea88571))
* **deps:** bump github.com/mattn/go-isatty from 0.0.17 to 0.0.18 ([#1702](https://gitea.com/gitea/runner/issues/1702)) ([385d71a](https://gitea.com/gitea/runner/commits/385d71a215b532dc983c640452e64f38426999b1))
* **deps:** bump github.com/mattn/go-isatty from 0.0.18 to 0.0.19 ([#1859](https://gitea.com/gitea/runner/issues/1859)) ([e912ab3](https://gitea.com/gitea/runner/commits/e912ab32c246692f29d9ef520e92f23edd073a4a))
* **deps:** bump github.com/moby/buildkit from 0.10.0 to 0.10.1 ([#1116](https://gitea.com/gitea/runner/issues/1116)) ([bb6ead4](https://gitea.com/gitea/runner/commits/bb6ead459a8e7b46e817369da3565b5ac86767a1))
* **deps:** bump github.com/moby/buildkit from 0.10.1 to 0.10.2 ([#1147](https://gitea.com/gitea/runner/issues/1147)) ([dfca2c5](https://gitea.com/gitea/runner/commits/dfca2c57df978aae5ec8b29a1a80c0fba09b56a8))
* **deps:** bump github.com/moby/buildkit from 0.10.2 to 0.10.3 ([#1153](https://gitea.com/gitea/runner/issues/1153)) ([a76c349](https://gitea.com/gitea/runner/commits/a76c34987233b065109457893a334879a4615baf))
* **deps:** bump github.com/moby/buildkit from 0.10.3 to 0.10.4 ([#1322](https://gitea.com/gitea/runner/issues/1322)) ([8bad6ac](https://gitea.com/gitea/runner/commits/8bad6aca50b7595fcc5253c72a2936fe37157a63))
* **deps:** bump github.com/moby/buildkit from 0.10.4 to 0.10.5 ([#1406](https://gitea.com/gitea/runner/issues/1406)) ([07a1890](https://gitea.com/gitea/runner/commits/07a18902ffa08ef1f209d8b976142805e2afb8f0))
* **deps:** bump github.com/moby/buildkit from 0.10.5 to 0.10.6 ([#1437](https://gitea.com/gitea/runner/issues/1437)) ([64e68bd](https://gitea.com/gitea/runner/commits/64e68bd7f21afdbbe2f40676c2a8c08701057753))
* **deps:** bump github.com/moby/buildkit from 0.10.6 to 0.11.0 ([#1563](https://gitea.com/gitea/runner/issues/1563)) ([d8ba8cb](https://gitea.com/gitea/runner/commits/d8ba8cbf17d816ed2be8a2b8a031c84b3d6fd8e9))
* **deps:** bump github.com/moby/buildkit from 0.11.0 to 0.11.1 ([#1576](https://gitea.com/gitea/runner/issues/1576)) ([09e9f6a](https://gitea.com/gitea/runner/commits/09e9f6a1ac00592b0341d73d6c642dd4aeb83a02))
* **deps:** bump github.com/moby/buildkit from 0.11.1 to 0.11.2 ([#1589](https://gitea.com/gitea/runner/issues/1589)) ([e259dc2](https://gitea.com/gitea/runner/commits/e259dc2835b6b759de7f4dea2daa490438f41bbe))
* **deps:** bump github.com/moby/buildkit from 0.11.2 to 0.11.3 ([#1635](https://gitea.com/gitea/runner/issues/1635)) ([2f55276](https://gitea.com/gitea/runner/commits/2f55276cea863b82aefc806c52773c5ba65774e2))
* **deps:** bump github.com/moby/buildkit from 0.11.3 to 0.11.4 ([#1669](https://gitea.com/gitea/runner/issues/1669)) ([aeee205](https://gitea.com/gitea/runner/commits/aeee2052de45b3ceff3fb66d3eff461379b89d92))
* **deps:** bump github.com/moby/buildkit from 0.11.4 to 0.11.5 ([#1703](https://gitea.com/gitea/runner/issues/1703)) ([351ae99](https://gitea.com/gitea/runner/commits/351ae99bc1b205c61b02d4fb04c06f8287b57f9f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1755](https://gitea.com/gitea/runner/issues/1755)) ([2ea7891](https://gitea.com/gitea/runner/commits/2ea78917879854128154124d3b7c8fa8d552e6e8))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1791](https://gitea.com/gitea/runner/issues/1791)) ([65ef31f](https://gitea.com/gitea/runner/commits/65ef31f102ceb75623973921099454637bab55b0))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1817](https://gitea.com/gitea/runner/issues/1817)) ([16c574c](https://gitea.com/gitea/runner/commits/16c574cd266ffe11ae48a37ccf3f6582ce90de21))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1858](https://gitea.com/gitea/runner/issues/1858)) ([b92d95f](https://gitea.com/gitea/runner/commits/b92d95f89931549a235e5f141ecb0a053f4216d3))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1880](https://gitea.com/gitea/runner/issues/1880)) ([a0307d3](https://gitea.com/gitea/runner/commits/a0307d3b7cd20437000d3931149fe600fe515a0f))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1882](https://gitea.com/gitea/runner/issues/1882)) ([3d1d8a9](https://gitea.com/gitea/runner/commits/3d1d8a9acac74363f09790f0efa450a27dbeb1ce))
* **deps:** bump github.com/moby/buildkit from 0.11.5 to 0.11.6 ([#1884](https://gitea.com/gitea/runner/issues/1884)) ([4c8da85](https://gitea.com/gitea/runner/commits/4c8da8558d0a741bb2cecd6d6edf832dc1a055a9))
* **deps:** bump github.com/moby/buildkit from 0.11.6 to 0.12.0 ([#1917](https://gitea.com/gitea/runner/issues/1917)) ([4810f69](https://gitea.com/gitea/runner/commits/4810f69367a141f9cc282d21ab61bb598d5ec584))
* **deps:** bump github.com/moby/buildkit from 0.12.0 to 0.12.1 ([#1947](https://gitea.com/gitea/runner/issues/1947)) ([c812715](https://gitea.com/gitea/runner/commits/c8127155bc64224a1392d2264b8d432328edea57))
* **deps:** bump github.com/moby/buildkit from 0.12.1 to 0.12.2 ([#1986](https://gitea.com/gitea/runner/issues/1986)) ([0c12273](https://gitea.com/gitea/runner/commits/0c12273ebaa937c52dea06d91b8985304874ab14))
* **deps:** bump github.com/moby/buildkit from 0.9.0 to 0.9.1 ([#838](https://gitea.com/gitea/runner/issues/838)) ([b5e8a18](https://gitea.com/gitea/runner/commits/b5e8a18683c77c7035c704bb7cec980401c73e18))
* **deps:** bump github.com/moby/buildkit from 0.9.1 to 0.9.2 ([#864](https://gitea.com/gitea/runner/issues/864)) ([f8b3563](https://gitea.com/gitea/runner/commits/f8b35631028c7705eb90b1d97300dc19c739cc72))
* **deps:** bump github.com/moby/buildkit from 0.9.2 to 0.9.3 ([#958](https://gitea.com/gitea/runner/issues/958)) ([dba7f83](https://gitea.com/gitea/runner/commits/dba7f8379f327232a1726b336d2cf4a6c3e80ccc))
* **deps:** bump github.com/moby/buildkit from 0.9.3 to 0.10.0 ([#1036](https://gitea.com/gitea/runner/issues/1036)) ([56486f1](https://gitea.com/gitea/runner/commits/56486f14b4a094542d2cf4a903676cd9f14220ec))
* **deps:** bump github.com/moby/patternmatcher from 0.5.0 to 0.6.0 ([#1985](https://gitea.com/gitea/runner/issues/1985)) ([323bee9](https://gitea.com/gitea/runner/commits/323bee9ab5a73473236432e0b89e4d22100cd24d))
* **deps:** bump github.com/opencontainers/image-spec ([#1649](https://gitea.com/gitea/runner/issues/1649)) ([04d12b0](https://gitea.com/gitea/runner/commits/04d12b0206a62cc48a765ed29f2a1ad732d33b3b))
* **deps:** bump github.com/opencontainers/image-spec ([#1781](https://gitea.com/gitea/runner/issues/1781)) ([15eb1fa](https://gitea.com/gitea/runner/commits/15eb1fa92a494f0abde33ad1bfe63496195db9ec))
* **deps:** bump github.com/opencontainers/image-spec ([#1807](https://gitea.com/gitea/runner/issues/1807)) ([1eacf23](https://gitea.com/gitea/runner/commits/1eacf23dcbbea2e5a7d17e8c325ce94ee0b1d3d1))
* **deps:** bump github.com/opencontainers/image-spec ([#1829](https://gitea.com/gitea/runner/issues/1829)) ([26f1b1f](https://gitea.com/gitea/runner/commits/26f1b1f4b6df887a08af0ed6e692962989f5a57c))
* **deps:** bump github.com/opencontainers/image-spec ([#1871](https://gitea.com/gitea/runner/issues/1871)) ([ef60885](https://gitea.com/gitea/runner/commits/ef608854d0b42704a3316c1b09c74c4fff5c40ae))
* **deps:** bump github.com/opencontainers/image-spec ([#1881](https://gitea.com/gitea/runner/issues/1881)) ([f9dcb58](https://gitea.com/gitea/runner/commits/f9dcb58db23061c016bba6712080275219f7008c))
* **deps:** bump github.com/opencontainers/image-spec ([#1883](https://gitea.com/gitea/runner/issues/1883)) ([70e1e37](https://gitea.com/gitea/runner/commits/70e1e372804054374c222c7de1e066bf7ab92b31))
* **deps:** bump github.com/opencontainers/image-spec ([#1893](https://gitea.com/gitea/runner/issues/1893)) ([f40d0b8](https://gitea.com/gitea/runner/commits/f40d0b873da435b6a1dc9fbb5391ce31b2d2dc0e))
* **deps:** bump github.com/opencontainers/image-spec ([#2010](https://gitea.com/gitea/runner/issues/2010)) ([83af8f8](https://gitea.com/gitea/runner/commits/83af8f8767091542dcc4263ebea06e6cc4edd7ce))
* **deps:** bump github.com/opencontainers/image-spec ([#954](https://gitea.com/gitea/runner/issues/954)) ([eedac17](https://gitea.com/gitea/runner/commits/eedac179c9f314ee38eda7482c45664df276e546))
* **deps:** bump github.com/opencontainers/runc from 1.1.3 to 1.1.5 ([#1708](https://gitea.com/gitea/runner/issues/1708)) ([6745999](https://gitea.com/gitea/runner/commits/6745999c105a1bb4429275103dc89b6a704b421c))
* **deps:** bump github.com/opencontainers/selinux ([#1102](https://gitea.com/gitea/runner/issues/1102)) ([e0044ba](https://gitea.com/gitea/runner/commits/e0044ba913149be76acd64ca7d1a66f27532f4c7))
* **deps:** bump github.com/opencontainers/selinux ([#1372](https://gitea.com/gitea/runner/issues/1372)) ([d587ce0](https://gitea.com/gitea/runner/commits/d587ce08e9b04738fd1d3a2433e1f8b4d44c4aed))
* **deps:** bump github.com/opencontainers/selinux ([#1625](https://gitea.com/gitea/runner/issues/1625)) ([8048baf](https://gitea.com/gitea/runner/commits/8048baf435665cb3293f593a06526dfff1691b54))
* **deps:** bump github.com/rhysd/actionlint from 1.6.10 to 1.6.11 ([#1106](https://gitea.com/gitea/runner/issues/1106)) ([065d630](https://gitea.com/gitea/runner/commits/065d630af01b4ed2a0c797d54933ea3b65c3cb38))
* **deps:** bump github.com/rhysd/actionlint from 1.6.11 to 1.6.12 ([#1126](https://gitea.com/gitea/runner/issues/1126)) ([de1dea9](https://gitea.com/gitea/runner/commits/de1dea93795ebf78e0ca531c35969f125035e8c4))
* **deps:** bump github.com/rhysd/actionlint from 1.6.12 to 1.6.13 ([#1173](https://gitea.com/gitea/runner/issues/1173)) ([166d063](https://gitea.com/gitea/runner/commits/166d0630593b40a288cc0190d4aa3e20b31d1461))
* **deps:** bump github.com/rhysd/actionlint from 1.6.13 to 1.6.14 ([#1234](https://gitea.com/gitea/runner/issues/1234)) ([9499612](https://gitea.com/gitea/runner/commits/9499612b5818774c5c6d73fbf2596ed464b7b01d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.14 to 1.6.15 ([#1239](https://gitea.com/gitea/runner/issues/1239)) ([4d9d6ec](https://gitea.com/gitea/runner/commits/4d9d6ecc924bef6c5976acd7f224558e21fa5b5f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.15 to 1.6.16 ([#1311](https://gitea.com/gitea/runner/issues/1311)) ([cf9d82f](https://gitea.com/gitea/runner/commits/cf9d82f69daab9dc907b893a38307a17f444c529))
* **deps:** bump github.com/rhysd/actionlint from 1.6.16 to 1.6.17 ([#1321](https://gitea.com/gitea/runner/issues/1321)) ([9b14698](https://gitea.com/gitea/runner/commits/9b146980dfabd14e9b9defea8f894ad8ab24132a))
* **deps:** bump github.com/rhysd/actionlint from 1.6.17 to 1.6.18 ([#1352](https://gitea.com/gitea/runner/issues/1352)) ([053bff1](https://gitea.com/gitea/runner/commits/053bff19a710669308da5e356ab028bab2a70b5d))
* **deps:** bump github.com/rhysd/actionlint from 1.6.18 to 1.6.19 ([#1360](https://gitea.com/gitea/runner/issues/1360)) ([3fccb5a](https://gitea.com/gitea/runner/commits/3fccb5ae3467fbd013ca72ed08732d2916beb1c2))
* **deps:** bump github.com/rhysd/actionlint from 1.6.19 to 1.6.20 ([#1371](https://gitea.com/gitea/runner/issues/1371)) ([bec45d0](https://gitea.com/gitea/runner/commits/bec45d0dc46e2ccae81c521f3db313b04b0c0831))
* **deps:** bump github.com/rhysd/actionlint from 1.6.20 to 1.6.21 ([#1380](https://gitea.com/gitea/runner/issues/1380)) ([c4be1cf](https://gitea.com/gitea/runner/commits/c4be1cfd6d8ab6664c3a3ffd7c4e5e253fa0b77f))
* **deps:** bump github.com/rhysd/actionlint from 1.6.21 to 1.6.22 ([#1427](https://gitea.com/gitea/runner/issues/1427)) ([8036f49](https://gitea.com/gitea/runner/commits/8036f49fb220bb2315ecba9707242cc0a8533aba))
* **deps:** bump github.com/rhysd/actionlint from 1.6.22 to 1.6.23 ([#1574](https://gitea.com/gitea/runner/issues/1574)) ([1b7212f](https://gitea.com/gitea/runner/commits/1b7212f5fd1e85d97a7119ab098941bd4cf5255b))
* **deps:** bump github.com/rhysd/actionlint from 1.6.23 to 1.6.24 ([#1721](https://gitea.com/gitea/runner/issues/1721)) ([5a331d2](https://gitea.com/gitea/runner/commits/5a331d2d9935565eabdf5a09226b7ba59a029133))
* **deps:** bump github.com/rhysd/actionlint from 1.6.24 to 1.6.25 ([#1870](https://gitea.com/gitea/runner/issues/1870)) ([4fc176f](https://gitea.com/gitea/runner/commits/4fc176f556dcf41e69854570ce65680f127a9357))
* **deps:** bump github.com/rhysd/actionlint from 1.6.8 to 1.6.9 ([#1018](https://gitea.com/gitea/runner/issues/1018)) ([9a54c99](https://gitea.com/gitea/runner/commits/9a54c99ac76a2667551bfd85c5aabcabc969e9dd))
* **deps:** bump github.com/rhysd/actionlint from 1.6.9 to 1.6.10 ([#1039](https://gitea.com/gitea/runner/issues/1039)) ([ee92e07](https://gitea.com/gitea/runner/commits/ee92e07daa0b82e4758807f424f79194e5cb857b))
* **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1273](https://gitea.com/gitea/runner/issues/1273)) ([17dd54d](https://gitea.com/gitea/runner/commits/17dd54d69231884ea561fe2bd0c0e059debd2658))
* **deps:** bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2 ([#1814](https://gitea.com/gitea/runner/issues/1814)) ([06054d2](https://gitea.com/gitea/runner/commits/06054d2931861d2f02cd4028d04164a4d52bbc1c))
* **deps:** bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3 ([#1841](https://gitea.com/gitea/runner/issues/1841)) ([e150310](https://gitea.com/gitea/runner/commits/e150310da77918e9ca904363f6cd80907877ff12))
* **deps:** bump github.com/spf13/cobra from 1.2.1 to 1.3.0 ([#955](https://gitea.com/gitea/runner/issues/955)) ([e9e6dda](https://gitea.com/gitea/runner/commits/e9e6ddadf567beb62ddb94b3fe472297081d4077))
* **deps:** bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#1040](https://gitea.com/gitea/runner/issues/1040)) ([fa71f0b](https://gitea.com/gitea/runner/commits/fa71f0b1e04a6391b77ef091c0db4d1a8f0ad1c1))
* **deps:** bump github.com/spf13/cobra from 1.4.0 to 1.5.0 ([#1233](https://gitea.com/gitea/runner/issues/1233)) ([4deb213](https://gitea.com/gitea/runner/commits/4deb21344d40dd4c029d5c1048a7f1d557ab8774))
* **deps:** bump github.com/spf13/cobra from 1.5.0 to 1.6.0 ([#1394](https://gitea.com/gitea/runner/issues/1394)) ([52c6edf](https://gitea.com/gitea/runner/commits/52c6edf0a48fa8287eaaa1ea28d89f65d05e8e7b))
* **deps:** bump github.com/spf13/cobra from 1.6.0 to 1.6.1 ([#1420](https://gitea.com/gitea/runner/issues/1420)) ([a9298d9](https://gitea.com/gitea/runner/commits/a9298d936c5265461dd8f6322961232b08fc6d12))
* **deps:** bump github.com/spf13/cobra from 1.6.1 to 1.7.0 ([#1724](https://gitea.com/gitea/runner/issues/1724)) ([c5ce502](https://gitea.com/gitea/runner/commits/c5ce502f9f9e98ea76a09858075e3ee91da6bd10))
* **deps:** bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#1061](https://gitea.com/gitea/runner/issues/1061)) ([7d946c4](https://gitea.com/gitea/runner/commits/7d946c43e208aa33a0f6183a2b84a8e9f33425ae))
* **deps:** bump github.com/stretchr/testify from 1.7.1 to 1.7.2 ([#1213](https://gitea.com/gitea/runner/issues/1213)) ([3415347](https://gitea.com/gitea/runner/commits/3415347efcdf6d6d7544c5df7cb33e6db043ff4a))
* **deps:** bump github.com/stretchr/testify from 1.7.2 to 1.7.5 ([#1232](https://gitea.com/gitea/runner/issues/1232)) ([c67abf2](https://gitea.com/gitea/runner/commits/c67abf2401ba57097e862b524b86de7171f1c609))
* **deps:** bump github.com/stretchr/testify from 1.7.5 to 1.8.0 ([#1240](https://gitea.com/gitea/runner/issues/1240)) ([50f0b0e](https://gitea.com/gitea/runner/commits/50f0b0e7f49d368583a30e515fb16db58c357bbc))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1404](https://gitea.com/gitea/runner/issues/1404)) ([db850e8](https://gitea.com/gitea/runner/commits/db850e818c2dd36f0021ac0a227ed99d2daafff2))
* **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#1650](https://gitea.com/gitea/runner/issues/1650)) ([a6c95ef](https://gitea.com/gitea/runner/commits/a6c95ef2a771afdc9f83f6750b642491e52f681c))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://gitea.com/gitea/runner/issues/1842)) ([c70a674](https://gitea.com/gitea/runner/commits/c70a6743f62af2c223f62c2fddd6a1a68271a50e))
* **deps:** bump golang.org/x/net from 0.4.0 to 0.7.0 ([#1636](https://gitea.com/gitea/runner/issues/1636)) ([b7a9eb9](https://gitea.com/gitea/runner/commits/b7a9eb9fbf86da76cbad35296ff11b64ee416d18))
* **deps:** bump golang.org/x/term from 0.10.0 to 0.11.0 ([#1948](https://gitea.com/gitea/runner/issues/1948)) ([74b74e8](https://gitea.com/gitea/runner/commits/74b74e847bad3c71ad7a1df76384a804bb46d644))
* **deps:** bump golang.org/x/term from 0.11.0 to 0.12.0 ([#1992](https://gitea.com/gitea/runner/issues/1992)) ([1d32507](https://gitea.com/gitea/runner/commits/1d32507b52c8a3c502415d1350aeca0eb9dcfe81))
* **deps:** bump golang.org/x/term from 0.3.0 to 0.4.0 ([#1577](https://gitea.com/gitea/runner/issues/1577)) ([a2856ae](https://gitea.com/gitea/runner/commits/a2856aed75e6d4e5607dc59ff074c9612553360f))
* **deps:** bump golang.org/x/term from 0.5.0 to 0.6.0 ([#1666](https://gitea.com/gitea/runner/issues/1666)) ([8dbd151](https://gitea.com/gitea/runner/commits/8dbd151fa71031c737d3f35580906f4028c355ca))
* **deps:** bump golang.org/x/term from 0.6.0 to 0.7.0 ([#1723](https://gitea.com/gitea/runner/issues/1723)) ([3ae0a9d](https://gitea.com/gitea/runner/commits/3ae0a9dfb7d9e4eb2adb3dc4a778290784fc5ce8))
* **deps:** bump golang.org/x/term from 0.8.0 to 0.9.0 ([#1869](https://gitea.com/gitea/runner/issues/1869)) ([b2191ae](https://gitea.com/gitea/runner/commits/b2191ae204a713a5e9bd4b3e1ec34b23b7fd34a6))
* **deps:** bump golang.org/x/term from 0.9.0 to 0.10.0 ([#1906](https://gitea.com/gitea/runner/issues/1906)) ([e343ea9](https://gitea.com/gitea/runner/commits/e343ea9d5f6fa478a364b1f502b4859eed63fa38))
* **deps:** bump golangci/golangci-lint-action from 2 to 3.1.0 ([#1019](https://gitea.com/gitea/runner/issues/1019)) ([7d403b8](https://gitea.com/gitea/runner/commits/7d403b83075932a18007210a62c6d0ce64df11a6))
* **deps:** bump golangci/golangci-lint-action from 3.1.0 to 3.2.0 ([#1169](https://gitea.com/gitea/runner/issues/1169)) ([8994773](https://gitea.com/gitea/runner/commits/89947735d508751b9e84d653635dec91299e9eec))
* **deps:** bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 ([#1408](https://gitea.com/gitea/runner/issues/1408)) ([925ac20](https://gitea.com/gitea/runner/commits/925ac2098b569f95beca47d5aac7e2287922963b))
* **deps:** bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 ([#1436](https://gitea.com/gitea/runner/issues/1436)) ([1b554ae](https://gitea.com/gitea/runner/commits/1b554aeff651cf801c11277bacb162ba31483b08))
* **deps:** bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 ([#1590](https://gitea.com/gitea/runner/issues/1590)) ([70ae63b](https://gitea.com/gitea/runner/commits/70ae63bd10eb4ac3a9a64633b88242f9a10122e3))
* **deps:** bump golangci/golangci-lint-action from 3.4.0 to 3.5.0 ([#1839](https://gitea.com/gitea/runner/issues/1839)) ([34db861](https://gitea.com/gitea/runner/commits/34db86138ba3f375c3797dfe5b2952f32397f7b1))
* **deps:** bump golangci/golangci-lint-action from 3.5.0 to 3.6.0 ([#1868](https://gitea.com/gitea/runner/issues/1868)) ([79f93be](https://gitea.com/gitea/runner/commits/79f93beef21650b25ab20a24aa91032e3a6b597e))
* **deps:** bump golangci/golangci-lint-action from 3.6.0 to 3.7.0 ([#1978](https://gitea.com/gitea/runner/issues/1978)) ([f64c267](https://gitea.com/gitea/runner/commits/f64c267dacd026b306da8ee405e84df56d124e69))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1179](https://gitea.com/gitea/runner/issues/1179)) ([5054455](https://gitea.com/gitea/runner/commits/50544556f8672c5553fd5241f4adb157698c781b))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1515](https://gitea.com/gitea/runner/issues/1515)) ([c044035](https://gitea.com/gitea/runner/commits/c044035d9ef358482efc996fba0a4c3f7f5f3847))
* **deps:** bump goreleaser/goreleaser-action from 4 to 5 ([#2008](https://gitea.com/gitea/runner/issues/2008)) ([2911b21](https://gitea.com/gitea/runner/commits/2911b2172cee57abbfdafd8ffd04fe4b6ccdd9dc))
* **deps:** bump gotest.tools/v3 from 3.0.3 to 3.4.0 ([#1379](https://gitea.com/gitea/runner/issues/1379)) ([c7ccb0d](https://gitea.com/gitea/runner/commits/c7ccb0dde9ab6332af7940e6f61e3ab73685fd07))
* **deps:** bump gotest.tools/v3 from 3.4.0 to 3.5.0 ([#1892](https://gitea.com/gitea/runner/issues/1892)) ([70956f2](https://gitea.com/gitea/runner/commits/70956f2929eddc7a65ffc16ac335db9472924d66))
* **deps:** bump gotest.tools/v3 from 3.5.0 to 3.5.1 ([#2009](https://gitea.com/gitea/runner/issues/2009)) ([60060a7](https://gitea.com/gitea/runner/commits/60060a7c9c7730d653bd8d980f1a17abd068e739))
* **deps:** bump megalinter/megalinter from 5 to 6.11.0 ([#1370](https://gitea.com/gitea/runner/issues/1370)) ([79384c3](https://gitea.com/gitea/runner/commits/79384c35ab42715144388b7d3e0b0290fb93a962))
* **deps:** bump megalinter/megalinter from 6.11.0 to 6.11.1 ([#1381](https://gitea.com/gitea/runner/issues/1381)) ([4541139](https://gitea.com/gitea/runner/commits/4541139109f2838bcac25ddd3d435a5c0661bf25))
* **deps:** bump megalinter/megalinter from 6.11.1 to 6.12.0 ([#1395](https://gitea.com/gitea/runner/issues/1395)) ([37f5b7f](https://gitea.com/gitea/runner/commits/37f5b7f4a2bc363be670ee65f623c6b54ec9d8c5))
* **deps:** bump megalinter/megalinter from 6.12.0 to 6.13.0 ([#1407](https://gitea.com/gitea/runner/issues/1407)) ([e81dc0d](https://gitea.com/gitea/runner/commits/e81dc0d295ad58d13927e61c91d3c4b4a6def942))
* **deps:** bump megalinter/megalinter from 6.13.0 to 6.14.0 ([#1426](https://gitea.com/gitea/runner/issues/1426)) ([c232981](https://gitea.com/gitea/runner/commits/c2329815b8bbc739043e1584e465322a602475e9))
* **deps:** bump megalinter/megalinter from 6.14.0 to 6.15.0 ([#1475](https://gitea.com/gitea/runner/issues/1475)) ([7754ba7](https://gitea.com/gitea/runner/commits/7754ba7fcf3f80b76bb4d7cf1edbe5935c8a6bdc))
* **deps:** bump megalinter/megalinter from 6.15.0 to 6.16.0 ([#1534](https://gitea.com/gitea/runner/issues/1534)) ([e3a722e](https://gitea.com/gitea/runner/commits/e3a722e7c6c697a1b0a38aa314bfb227da4deb5f))
* **deps:** bump megalinter/megalinter from 6.16.0 to 6.17.0 ([#1540](https://gitea.com/gitea/runner/issues/1540)) ([cfc1b91](https://gitea.com/gitea/runner/commits/cfc1b91aa11220f102470622cacd105cedd931c9))
* **deps:** bump megalinter/megalinter from 6.17.0 to 6.18.0 ([#1550](https://gitea.com/gitea/runner/issues/1550)) ([f0c6fa1](https://gitea.com/gitea/runner/commits/f0c6fa11be88093fbcb7ecb1cbebfee8874008cc))
* **deps:** bump megalinter/megalinter from 6.18.0 to 6.19.0 ([#1610](https://gitea.com/gitea/runner/issues/1610)) ([460c78d](https://gitea.com/gitea/runner/commits/460c78da0773c226162d93623e8932715a22e443))
* **deps:** bump megalinter/megalinter from 6.19.0 to 6.20.0 ([#1665](https://gitea.com/gitea/runner/issues/1665)) ([7ee2138](https://gitea.com/gitea/runner/commits/7ee2138418822460ed576391201db55803bd7069))
* **deps:** bump megalinter/megalinter from 6.20.0 to 6.20.1 ([#1679](https://gitea.com/gitea/runner/issues/1679)) ([5a2112a](https://gitea.com/gitea/runner/commits/5a2112a7f8e282fcb75f1c73caa1582811ff4890))
* **deps:** bump megalinter/megalinter from 6.20.1 to 6.21.0 ([#1699](https://gitea.com/gitea/runner/issues/1699)) ([775a128](https://gitea.com/gitea/runner/commits/775a128cd4660be4f98da019c52a2ab9a5080433))
* **deps:** bump megalinter/megalinter from 6.21.0 to 6.22.1 ([#1710](https://gitea.com/gitea/runner/issues/1710)) ([ad7a8a3](https://gitea.com/gitea/runner/commits/ad7a8a3559d937d79e9425b8d29714752879b823))
* **deps:** bump megalinter/megalinter from 6.22.1 to 6.22.2 ([#1720](https://gitea.com/gitea/runner/issues/1720)) ([4ab812c](https://gitea.com/gitea/runner/commits/4ab812c1b10bc16a54dd208963d6b60e2835d328))
* **deps:** bump megalinter/megalinter from 6.22.2 to 7.0.2 ([#1827](https://gitea.com/gitea/runner/issues/1827)) ([481999f](https://gitea.com/gitea/runner/commits/481999f59da465d2d8743fd3d97e24bb6e1a2311))
* **deps:** bump megalinter/megalinter from 7.0.2 to 7.0.4 ([#1838](https://gitea.com/gitea/runner/issues/1838)) ([04b5c73](https://gitea.com/gitea/runner/commits/04b5c739d2649b67565dd7addbe2bf5fc910e2a6))
* **deps:** bump megalinter/megalinter from 7.0.4 to 7.1.0 ([#1857](https://gitea.com/gitea/runner/issues/1857)) ([3555d65](https://gitea.com/gitea/runner/commits/3555d65d08bc9c66a6e4541b2064a082af1b6b3e))
* **deps:** bump megalinter/megalinter from 7.1.0 to 7.2.1 ([#1931](https://gitea.com/gitea/runner/issues/1931)) ([3939f48](https://gitea.com/gitea/runner/commits/3939f48e6def83e8a6ad7358a848ce11912c3798))
* **deps:** bump megalinter/megalinter from 7.2.1 to 7.3.0 ([#1965](https://gitea.com/gitea/runner/issues/1965)) ([2a0a0a1](https://gitea.com/gitea/runner/commits/2a0a0a1a62b8c72cd891208f085c00cb00402b57))
* **deps:** update go-git to ce62f3e9ff86270538a514a68d3bd5563a733e3b ([#1725](https://gitea.com/gitea/runner/issues/1725)) ([148a545](https://gitea.com/gitea/runner/commits/148a54502167880f90dbf2f68ff44478679f41dd))

### Continuous Integration

* Add snapshot artifacts to PR builds ([#574](https://gitea.com/gitea/runner/issues/574)) ([da07324](https://gitea.com/gitea/runner/commits/da07324779854fda8bff4bd45452a04ce5e68a1e))
* add Winget Releaser workflow ([#1623](https://gitea.com/gitea/runner/issues/1623)) ([34ab815](https://gitea.com/gitea/runner/commits/34ab8150bf8b438964b7e52796f2123466e0cf76))
* always grab latest go version ([#1124](https://gitea.com/gitea/runner/issues/1124)) ([25ebaf4](https://gitea.com/gitea/runner/commits/25ebaf4e7b507ffc7cabdb2a560348af5bec4ac3))
* bump go to 1.18 ([#1087](https://gitea.com/gitea/runner/issues/1087)) ([11c3fb3](https://gitea.com/gitea/runner/commits/11c3fb39dd5d3907baf8f13ef6953d0a21321cc2))
* change winget job runner to `ubuntu-latest` ([#1959](https://gitea.com/gitea/runner/issues/1959)) ([f55ae1a](https://gitea.com/gitea/runner/commits/f55ae1a0bc4fb14d463b79a62ff9a37d7e5bda68))
* **choco:** fix chocolatey docker build ([#903](https://gitea.com/gitea/runner/issues/903)) ([6ebcac3](https://gitea.com/gitea/runner/commits/6ebcac37710ae5d4eb87f63f31d7e1be8fe3d217))
* **choco:** update chocolatey to 1.1.0 ([#1164](https://gitea.com/gitea/runner/issues/1164)) ([507ae61](https://gitea.com/gitea/runner/commits/507ae61d1b3ac3f66273a00fad1321af8c6df260))
* deduplicate running workflows ([#1751](https://gitea.com/gitea/runner/issues/1751)) ([7c6237d](https://gitea.com/gitea/runner/commits/7c6237d93f01b2c50a98a488e40ab7f9a719da35))
* **golangci-lint:** switch to latest ([#870](https://gitea.com/gitea/runner/issues/870)) ([cd24df5](https://gitea.com/gitea/runner/commits/cd24df5727c54de6671f83b44b7527f3c74f122f))
* make stalebot less annoying ([#1587](https://gitea.com/gitea/runner/issues/1587)) ([01952d0](https://gitea.com/gitea/runner/commits/01952d05c858c70ec0f007d90e394cf4cdfcd83b))
* **Mergify:** configuration update ([#924](https://gitea.com/gitea/runner/issues/924)) ([cad4bc8](https://gitea.com/gitea/runner/commits/cad4bc8c367fd4e623dfe2820bd5f3daca80bf58))
* **Mergify:** configuration update ([#966](https://gitea.com/gitea/runner/issues/966)) ([557dc75](https://gitea.com/gitea/runner/commits/557dc755b8fe6b286d15bd7a1091f3731dc44791))
* **Mergify:** configuration update ([#972](https://gitea.com/gitea/runner/issues/972)) ([4f8da0a](https://gitea.com/gitea/runner/commits/4f8da0a51cd44753f3e4ea0364877f6d2aa11f1d))
* **mergify:** fix deprecated strict mode ([#862](https://gitea.com/gitea/runner/issues/862)) ([80db062](https://gitea.com/gitea/runner/commits/80db062472a72735f5fa454a1f3977ea8702aad4))
* pin goreleaser version ([#825](https://gitea.com/gitea/runner/issues/825)) ([b1e4844](https://gitea.com/gitea/runner/commits/b1e4844aac4b8a143312c3ec9d3859362c669f88))
* remove macOS tests ([#860](https://gitea.com/gitea/runner/issues/860)) ([aad724c](https://gitea.com/gitea/runner/commits/aad724c87abb04184c9205643ec4659043637b07))
* replace superlinter with megalinter ([#923](https://gitea.com/gitea/runner/issues/923)) ([9be6a58](https://gitea.com/gitea/runner/commits/9be6a58c0eb706d514e355aa10f2b4c638cfec4a))
* **stalebot:** ignore feature requests ([#1017](https://gitea.com/gitea/runner/issues/1017)) ([932ad0e](https://gitea.com/gitea/runner/commits/932ad0ef6a2483b86f331f933b9c9081326ecbac))
* switch macos job to macos-10.15 ([#849](https://gitea.com/gitea/runner/issues/849)) ([977b526](https://gitea.com/gitea/runner/commits/977b526384a00e4f0804db6ddbaa3248c1918c75))

### Documentation

* add brew installation from HEAD ([#886](https://gitea.com/gitea/runner/issues/886)) ([df4ef4d](https://gitea.com/gitea/runner/commits/df4ef4de808edb8e3b657e4e24c96f8643ffcd32))
* add naming rule for git ref ([70cc6c0](https://gitea.com/gitea/runner/commits/70cc6c017bc8382b2efe4a50266f733037094da8))
* Add Prerequisites ([#11](https://gitea.com/gitea/runner/issues/11)) ([7abe039](https://gitea.com/gitea/runner/commits/7abe039c311dd257f7ac4f0c569afe701c1ae6d0)), closes [#10](https://gitea.com/gitea/runner/issues/10)
* add quickstart to readme ([03c4c4e](https://gitea.com/gitea/runner/commits/03c4c4e1164e18306652ed5d18f7e9dcfb3b139a))
* clarifying skipping steps / jobs ([#1480](https://gitea.com/gitea/runner/issues/1480)) ([7073eac](https://gitea.com/gitea/runner/commits/7073eac24057020cd54bcb00402281be55df3246))
* Expand the GITHUB_TOKEN section ([#968](https://gitea.com/gitea/runner/issues/968)) ([b1f5963](https://gitea.com/gitea/runner/commits/b1f5963c8623d5a5a92d2783b30da73c339232e2))
* fix brew formula link ([#846](https://gitea.com/gitea/runner/issues/846)) ([2319eb2](https://gitea.com/gitea/runner/commits/2319eb210f8d097da470331d2d12e7c6db0f8056))
* **README:** add --env-file flag ([#147](https://gitea.com/gitea/runner/issues/147)) ([f6b7613](https://gitea.com/gitea/runner/commits/f6b761378af17296aba8ff3be54e26d56acfeae5))
* **README:** remove go install ([#1123](https://gitea.com/gitea/runner/issues/1123)) ([e15f7a5](https://gitea.com/gitea/runner/commits/e15f7a51c6cf3e76cdc8ef8401b87080b7e9f95e))
* **readme:** update some header. ([#1](https://gitea.com/gitea/runner/issues/1)) ([474683c](https://gitea.com/gitea/runner/commits/474683c0e8bca9bef135fc49d44e650144591dfb))
* remove help section ([#1648](https://gitea.com/gitea/runner/issues/1648)) ([4b56ace](https://gitea.com/gitea/runner/commits/4b56aced159489639e5ca8226e2b82f5c7e9823e))

### Tests

* add test for networking setup in act ([#1375](https://gitea.com/gitea/runner/issues/1375)) ([ff5e289](https://gitea.com/gitea/runner/commits/ff5e2898041230978cef9af37dd87c09b07a77e8))
* Do not leak step env in composite ([#1585](https://gitea.com/gitea/runner/issues/1585)) ([3748772](https://gitea.com/gitea/runner/commits/3748772201458442afde9a582d70b40db05deaf8))
* Enshure ForcePull config doesn't break docker actions (1661) ([24440d9](https://gitea.com/gitea/runner/commits/24440d9f15673993dc6986f2c64dbd056e4286fe))
* env.ACT in if condition of the workflow ([#965](https://gitea.com/gitea/runner/issues/965)) ([3f0375a](https://gitea.com/gitea/runner/commits/3f0375aeff1973b674b33a0e4c44bc05a41714c5))
* make sure workflow_call is not a github event calling our workflow ([#1520](https://gitea.com/gitea/runner/issues/1520)) ([3f4a6dc](https://gitea.com/gitea/runner/commits/3f4a6dc4626eed03991e207b4d0e343b634edf9e))

### Chores

* add .golangci config ([518aa1b](https://gitea.com/gitea/runner/commits/518aa1b2a76822e413293701af87be5f94b37c82))
* Add makefile ([a3e9bbe](https://gitea.com/gitea/runner/commits/a3e9bbed25680fcf46bfad48dd04e50565412b47))
* Add time sleep 1 second. ([bca586f](https://gitea.com/gitea/runner/commits/bca586ffd0bb9a726aed8b3d86e7c76b6ca0c3ef))
* add Winget installation instructions ([#1622](https://gitea.com/gitea/runner/issues/1622)) ([c378a7d](https://gitea.com/gitea/runner/commits/c378a7d28b5c20de9fbb4456544329d04ffb62a4))
* bump version ([0d07805](https://gitea.com/gitea/runner/commits/0d0780580cfe63786c21a31500ccb69797cb9cf6))
* bump version ([33cb8bc](https://gitea.com/gitea/runner/commits/33cb8bca649a1b63ac59596b20d4001ca9cbd3e6))
* bump VERSION to 0.2.39 ([efb12b7](https://gitea.com/gitea/runner/commits/efb12b7f122e1daa3fce410a0808098a23dad910))
* bump VERSION to 0.2.40 ([3e23b4f](https://gitea.com/gitea/runner/commits/3e23b4fbb558bc4e8bdc72d4f75e0d94e5be3bf1))
* bump VERSION to 0.2.41 ([b91e4b0](https://gitea.com/gitea/runner/commits/b91e4b0d5576ee0efeb381c9427b791e5272ece9))
* bump VERSION to 0.2.42 ([72d0321](https://gitea.com/gitea/runner/commits/72d03214c5ed8985bf23dc15440fa5301e9d13bf))
* bump VERSION to 0.2.43 ([44b510f](https://gitea.com/gitea/runner/commits/44b510f48cac6a58654bf920ee26cc6e1d36ac3a))
* bump VERSION to 0.2.44 ([220d6f1](https://gitea.com/gitea/runner/commits/220d6f1251aea75c378e5a98546f43cd918e1892))
* bump VERSION to 0.2.45 ([f1df2ca](https://gitea.com/gitea/runner/commits/f1df2ca5d6b315935cf9bb0d2a989175adf4e894))
* bump VERSION to 0.2.46 ([b0d0cec](https://gitea.com/gitea/runner/commits/b0d0cec71f2240fc629dcaa1889ed8c4ad52beb9))
* bump VERSION to 0.2.47 ([310cb79](https://gitea.com/gitea/runner/commits/310cb79e81e07ef11985b2f934092b82b167ccf1))
* bump VERSION to 0.2.48 ([3a0a642](https://gitea.com/gitea/runner/commits/3a0a6425a80f3405670d8d7b11ed18ea56a09af4))
* bump VERSION to 0.2.49 ([d5d8548](https://gitea.com/gitea/runner/commits/d5d854854651c151ecd914bb6e2e370e0057929a))
* bump VERSION to 0.2.50 ([80b0955](https://gitea.com/gitea/runner/commits/80b0955303888742c3ab73af5758bb7b01f5f57c))
* bump VERSION to 0.2.51 ([c241ecd](https://gitea.com/gitea/runner/commits/c241ecda318e4977501971c267f0a2fad70ea5cd))
* **client:** support gRPC and gRPC web protocol ([d359276](https://gitea.com/gitea/runner/commits/d359276fe1d4203f84f21249f40bb942a29fc7a5))
* **config:** remove zerolog ([9225a3a](https://gitea.com/gitea/runner/commits/9225a3a856e85f63e1160bff0acbd3cb21f3d9d4))
* **deps:** cleanup go.mod a little bit ([#136](https://gitea.com/gitea/runner/issues/136)) ([ca648a3](https://gitea.com/gitea/runner/commits/ca648a37c8df80a4e96ffa7c3360c13c1a0715be))
* **engine:** ping docker daemon ([4eb5213](https://gitea.com/gitea/runner/commits/4eb5213294e4be894712f5a081cc18f01db308ac))
* fix lint error ([827bd95](https://gitea.com/gitea/runner/commits/827bd953e7d9a0dbe7caab149b898a926546a044))
* fix promote script ([b3c06dc](https://gitea.com/gitea/runner/commits/b3c06dc19d5ba396336615cb404cd394bd3e9f30))
* fix release script to trigger gh-act ([787388d](https://gitea.com/gitea/runner/commits/787388daf54f3d78c56da29a37062b019e10e84f))
* fold integration test outputs ([#1111](https://gitea.com/gitea/runner/issues/1111)) ([a85e89d](https://gitea.com/gitea/runner/commits/a85e89d3fb8e1a3b2e9dda843e18aa819095e56e))
* **go:** rename module name ([aa765ea](https://gitea.com/gitea/runner/commits/aa765ea4d0fe51a1211b755f31f96a1d5748aa8a))
* **gRPC:** add request interface client ([f2fb879](https://gitea.com/gitea/runner/commits/f2fb8798fa1556ad80d78a9649588a6089dfedcd))
* **gRPC:** add update and update step for client ([9be39b8](https://gitea.com/gitea/runner/commits/9be39b8cd4629678dac490eb5f0710cdd12ae00e))
* **gRPC:** fetch new stage data. ([631f801](https://gitea.com/gitea/runner/commits/631f801aa658ed0fa07df5025b067e20f4db4749))
* **gRPC:** register new runner ([449388f](https://gitea.com/gitea/runner/commits/449388f3ab1dae8944dd955288c77373036a96ac))
* load config from env file ([aa9ed05](https://gitea.com/gitea/runner/commits/aa9ed0590359c29181811b9f503921e414e465b2))
* makefile update ([872a695](https://gitea.com/gitea/runner/commits/872a69548c77b07c01ba18fb507ffe0bfe278f5e))
* **piepline:** add runtime package. ([7d55fd5](https://gitea.com/gitea/runner/commits/7d55fd57c94d1862f91ff7e15b892ab5f4b0af88))
* **poller:** add error retry limit ([45b0429](https://gitea.com/gitea/runner/commits/45b0429b2151c1dddaa722206bbd518f315b2006))
* **poller:** add metric to track the worker number ([abdb547](https://gitea.com/gitea/runner/commits/abdb547b1b5f5f3e5addc3c9c07b8eb0d9d1fc22))
* **poller:** Add poller package ([0b885c5](https://gitea.com/gitea/runner/commits/0b885c5e5fe33e2781755a8a02fa4bddfa997306))
* **proto:** Add ping request. ([4d7ef95](https://gitea.com/gitea/runner/commits/4d7ef95d402d3e53d175e262e01af224f4c58da7))
* **proto:** replace import path. ([307c2c6](https://gitea.com/gitea/runner/commits/307c2c63de7e486c785f191ec2e66e7a9b49df80))
* push new version ([003947e](https://gitea.com/gitea/runner/commits/003947ea33bfabcaa8c3787aa78780fd88b9c1ff))
* **register:** update proto file. ([c161a48](https://gitea.com/gitea/runner/commits/c161a48a0a5cdea220acc96b9e5df6c0b9c27e54))
* Remove obsolete Container.UpdateFromPath ([#1631](https://gitea.com/gitea/runner/issues/1631)) ([21ea3d0](https://gitea.com/gitea/runner/commits/21ea3d0d5fc92cf9ae93c62c0af7ec3a5da36541))
* remove tracing codes ([2f879c4](https://gitea.com/gitea/runner/commits/2f879c41c4a9e44dded5e947b08f45b98976c50d))
* remove unused logrous package ([5213612](https://gitea.com/gitea/runner/commits/52136120333c5738feae5fc5f43c9577fdcbac3c))
* replace with code.gitea.io/bots-proto-go ([44988db](https://gitea.com/gitea/runner/commits/44988db9f89d2c55128688ec69309050aa140dcc))
* run act from cli on linux ([#1758](https://gitea.com/gitea/runner/issues/1758)) ([c81a770](https://gitea.com/gitea/runner/commits/c81a770bc5f8041c2c169e227996ea34ae62ec48)), references [#1756](https://gitea.com/gitea/runner/issues/1756)
* **runner:** add default label ""self-hosted" ([dada073](https://gitea.com/gitea/runner/commits/dada0730b04723483552337a2fbf1c40489571bb))
* **runner:** add new register command ([a503f74](https://gitea.com/gitea/runner/commits/a503f7429f82e1f55c8c13293469c85db5dfa39b))
* **runner:** add new token in header ([d178051](https://gitea.com/gitea/runner/commits/d178051832efa7b617b117c0d4c4dbe543d79c06))
* **runner:** add runner uuid in http header. ([bf5e3dc](https://gitea.com/gitea/runner/commits/bf5e3dc302e195dbe0b65688f5a642c372af0411))
* **runner:** cancel task if get the cancel from server ([08c94bb](https://gitea.com/gitea/runner/commits/08c94bb564709d87fe62a094d07f30578b16f594))
* **runner:** refactor register flow ([3a15031](https://gitea.com/gitea/runner/commits/3a1503138b7360ff143d7e48b09effc5ef45f3ab))
* **runner:** remove client secret and add UUID in runner ([e08495a](https://gitea.com/gitea/runner/commits/e08495af090634d69704f2d9667bdddce7c47684))
* **runner:** remove update runner status ([e9e42e8](https://gitea.com/gitea/runner/commits/e9e42e850b526fa87e5edde09b3cb60d57e9fbf2))
* **runner:** support update log and task ([82431d8](https://gitea.com/gitea/runner/commits/82431d8e114b75c9d4a02f009c74503f1dca0789))
* **runner:** update fetch job request ([7486d2a](https://gitea.com/gitea/runner/commits/7486d2ab918aacdb6c0bb3b85a3c7aa6230d1afb))
* **runner:** update runner status if start program ([d4c1515](https://gitea.com/gitea/runner/commits/d4c1515f4ef6ebe732d8bee9cd6f3e019d1560fb))
* **runner:** update runner status when start job ([7bebd2b](https://gitea.com/gitea/runner/commits/7bebd2bbadd10171b355306a25294cb9576ffb9b))
* **runner:** wait workflow done before shutdown the service ([d1114da](https://gitea.com/gitea/runner/commits/d1114da299be65d69a4d1bb7574851c1f6cd5d1a))
* **runtime:** check error message ([5051e4a](https://gitea.com/gitea/runner/commits/5051e4aebdbb51575fa77bf472f4ee4ce03dfd0c))
* **runtime:** fetch build data ([d3d56ed](https://gitea.com/gitea/runner/commits/d3d56ed0ef86b83745555adb849852ccd3caacb8))
* update docs for installing act as GH CLI extension ([7ada9d3](https://gitea.com/gitea/runner/commits/7ada9d3f74f4d8705ba01ce110f3d611989dd442))
* update go-git dependency in go.mod ([#30](https://gitea.com/gitea/runner/issues/30)) ([8f12a6c](https://gitea.com/gitea/runner/commits/8f12a6c947deaa3488013503a518144c9ecdd035))
* update grpc router path ([f05e08a](https://gitea.com/gitea/runner/commits/f05e08a76750d1f8459b75a3b8f1461f54e4aa1c))
* update LICENSE ([7fb84a5](https://gitea.com/gitea/runner/commits/7fb84a54a8dd93b5874c7edb89861f1c7d23b469))
* upgrade act ([66ec50e](https://gitea.com/gitea/runner/commits/66ec50e427cb8e1b7b217803636e80704ae27b06))
* upgrade docker v19.03.12 ([#348](https://gitea.com/gitea/runner/issues/348)) ([3eef111](https://gitea.com/gitea/runner/commits/3eef111e46f42f43e62b64e6cdfbade4fe0935e2))
* upgrade golangci-lint and address findings ([#1904](https://gitea.com/gitea/runner/issues/1904)) ([724ec91](https://gitea.com/gitea/runner/commits/724ec918c92c162e537ba879e6d1ccc997f63483))
* use new style octal ([#1630](https://gitea.com/gitea/runner/issues/1630)) ([1316307](https://gitea.com/gitea/runner/commits/131630731374a16dbb7c8ed4893397c52d9c4c98))
* **yaml:** Improve `ParseRawOn` ([#28](https://gitea.com/gitea/runner/issues/28)) ([3daf313](https://gitea.com/gitea/runner/commits/3daf3132057d33bee325857b016a60db332fd050))

### Reverts

* Revert "Erase needs of job in SingleWorkflow (#9)" (#32) ([568f053](https://gitea.com/gitea/runner/commits/568f0537239f4566ad2b05b8866611d5af612b94)), references [#9](https://gitea.com/gitea/runner/issues/9) [#32](https://gitea.com/gitea/runner/issues/32)
* Revert "feat: allow existing logger from context (#859)" (#890) ([4ae71b5](https://gitea.com/gitea/runner/commits/4ae71b5878a8f6906ad3d2cd1d8cef0791527f20)), references [#859](https://gitea.com/gitea/runner/issues/859) [#890](https://gitea.com/gitea/runner/issues/890)
* deprecation of containerArchitecture ([#1514](https://gitea.com/gitea/runner/issues/1514)) ([4989f44](https://gitea.com/gitea/runner/commits/4989f444f1085ac0ec6b9bc35de92dc22f3c9f5d))
* Don't stop JobContainer on failure ([#1041](https://gitea.com/gitea/runner/issues/1041)) ([4d71071](https://gitea.com/gitea/runner/commits/4d7107161e231406bee856258ec847844623224e))

## [0.2.5](https://gitea.com/gitea/runner/compare/v0.2.4...v0.2.5) (2023-08-02)

## [0.2.4](https://gitea.com/gitea/runner/compare/v0.2.3...v0.2.4) (2023-07-25)

### Bug Fixes

* **register:** refactor context usage in registration functions ([#286](https://gitea.com/gitea/runner/issues/286)) ([8705871](https://gitea.com/gitea/runner/commits/87058716fb95ffb715472ee7fab2e8fb8a6e75b0))

### Refactoring

* **register:** refactor registration functions and improve logging ([#288](https://gitea.com/gitea/runner/issues/288)) ([cf92a97](https://gitea.com/gitea/runner/commits/cf92a979e274dace1d545cbe140f4f4db064f921))

### Continuous Integration

* **lint:** refactor code for clarity and linting compliance ([#289](https://gitea.com/gitea/runner/issues/289)) ([db662b3](https://gitea.com/gitea/runner/commits/db662b36906a3aa0bd1211cee64e4ce73288f485))

## [0.2.3](https://gitea.com/gitea/runner/compare/v0.2.2...v0.2.3) (2023-07-03)

## [0.2.2](https://gitea.com/gitea/runner/compare/v0.2.1...v0.2.2) (2023-07-01)

## [0.2.1](https://gitea.com/gitea/runner/compare/v0.2.0...v0.2.1) (2023-06-30)

### Features

* improve Docker configuration and detection handling ([#242](https://gitea.com/gitea/runner/issues/242)) ([9e4a5f7](https://gitea.com/gitea/runner/commits/9e4a5f7363a5559ddaa3a48e91bb8e69db861de4))

## [0.2.0](https://gitea.com/gitea/runner/compare/v0.1.8...v0.2.0) (2023-06-15)

### Refactoring

* daemon, config, and logging for better clarity ([#225](https://gitea.com/gitea/runner/issues/225)) ([69c55ee](https://gitea.com/gitea/runner/commits/69c55ee003a3553d4e8b50d9bf98d2c31c8832e4))

### Documentation

* improve examples README and organization ([#230](https://gitea.com/gitea/runner/issues/230)) ([a83f29d](https://gitea.com/gitea/runner/commits/a83f29d5a96a78b61a9dc6967d09a31ef0ad5d7e))

## [0.1.8](https://gitea.com/gitea/runner/compare/v0.1.7...v0.1.8) (2023-05-09)

## [0.1.7](https://gitea.com/gitea/runner/compare/v0.1.6...v0.1.7) (2023-05-06)

### Continuous Integration

* improve release process and test robustness ([#173](https://gitea.com/gitea/runner/issues/173)) ([eef3c32](https://gitea.com/gitea/runner/commits/eef3c32eb2418dcbcd40cb13076c2887d4f6913b))

### Chores

* improve Dockerfile, README, and testing settings ([#172](https://gitea.com/gitea/runner/issues/172)) ([c40b651](https://gitea.com/gitea/runner/commits/c40b6518736636db6edf7bbecf4d650f3a2117c4))

## [0.1.6](https://gitea.com/gitea/runner/compare/v0.1.5...v0.1.6) (2023-05-01)

### Build

* improve compression and update GitHub actions ([#168](https://gitea.com/gitea/runner/issues/168)) ([b498341](https://gitea.com/gitea/runner/commits/b49834185777fffed9e0b4d7c3b08ad82fcdbcfa))
* optimize Dockerfile and update dependencies ([#162](https://gitea.com/gitea/runner/issues/162)) ([0d727eb](https://gitea.com/gitea/runner/commits/0d727eb2621c706d973987490bb3b14163693cf3))

## [0.1.5](https://gitea.com/gitea/runner/compare/v0.1.4...v0.1.5) (2023-04-28)

### Continuous Integration

* improve API usage and test robustness across platforms ([#159](https://gitea.com/gitea/runner/issues/159)) ([49d2cb0](https://gitea.com/gitea/runner/commits/49d2cb0cb53ebe6f67284b2fd27e11fe81884d85))

## [0.1.4](https://gitea.com/gitea/runner/compare/v0.1.3...v0.1.4) (2023-04-27)

## [0.1.3](https://gitea.com/gitea/runner/compare/v0.1.2...v0.1.3) (2023-04-27)

### Continuous Integration

* **actions:** add build docker image workflow ([#118](https://gitea.com/gitea/runner/issues/118)) ([66bab3d](https://gitea.com/gitea/runner/commits/66bab3d8059fde868b59d1aeb9e0bdecb6c0d55e))

## [0.1.2](https://gitea.com/gitea/runner/compare/v0.1.1...v0.1.2) (2023-04-13)

## [0.1.1](https://gitea.com/gitea/runner/compare/v0.1.0...v0.1.1) (2023-04-11)

### Features

* add artifact api from gitea server ([#103](https://gitea.com/gitea/runner/issues/103)) ([b79c3aa](https://gitea.com/gitea/runner/commits/b79c3aa1a325ab49902fd662174a5ea0f8eed196))

## [0.1.0](https://gitea.com/gitea/runner/compare/v0.0.1...v0.1.0) (2023-04-04)

## [0.0.1](https://gitea.com/gitea/runner/compare/e54de3724292fa3fc4a50f1e1186964ce1e7a816...v0.0.1) (2023-02-26)

### Features

* add task log output hook ([dfcb104](https://gitea.com/gitea/runner/commits/dfcb104af02b536df7dfd7dc6d5911b654038ae9))
* auto remove containers ([8f950cc](https://gitea.com/gitea/runner/commits/8f950ccec0a9aae4a6e66d10dc28edd5fdf04ccb))
* avoid reuse of containers ([cbd21b8](https://gitea.com/gitea/runner/commits/cbd21b8ec539a5def5649b10702e55f6656e2ee9))
* **cli:** add register command basic functions ([561bfad](https://gitea.com/gitea/runner/commits/561bfad7c5413e29f61f44fc4526ad6123705c4c))
* **cli:** make register interactive stages working ([bfcc6f5](https://gitea.com/gitea/runner/commits/bfcc6f56f320d891ba61e886acf3b49ac6d766e2))
* **cli:** set interactive stages in register command ([8f9d750](https://gitea.com/gitea/runner/commits/8f9d7506dc71aaa55db74b80a7118bc18d6947f2))
* don't require docker ([#16](https://gitea.com/gitea/runner/issues/16)) ([990cf93](https://gitea.com/gitea/runner/commits/990cf93c7136669408eb1832cd05df3ad4dd81b3))
* fill default task runner ([3281c66](https://gitea.com/gitea/runner/commits/3281c6637d5b49f761e80584235d67384541a04a))
* move main task logic to runtime package ([c5edbbd](https://gitea.com/gitea/runner/commits/c5edbbd004418890b593dd9217160b02f5f2143d))
* parse register arguments and do register when non-interactive mode ([45de619](https://gitea.com/gitea/runner/commits/45de6199d12626de6c934444624c24e8f887e8f6))
* **poller:** support scheduler to fetch task and dispatch to worker ([a13ed0c](https://gitea.com/gitea/runner/commits/a13ed0c52e2879ef5fa2ec8632da2efd535c56d0))
* **register:** add label format validation ([3f676a7](https://gitea.com/gitea/runner/commits/3f676a742b72f550b7900a1c0d4faab19b11750f))
* **register:** add labels to local runner file ([86f66ea](https://gitea.com/gitea/runner/commits/86f66eaa4508603d9359398e8dcf866d8689bb89))
* replace ciphertext in log ([d650bae](https://gitea.com/gitea/runner/commits/d650bae0b6212b5a38aef62a6fab13c04b971f45))
* support post step log ([ab4e06f](https://gitea.com/gitea/runner/commits/ab4e06f977ca78181077e5de43d7aad070204192))
* task can report step and final result ([6030610](https://gitea.com/gitea/runner/commits/6030610c04cc75219d6047ce51e57b45ea3c1fb9))
* use specified labels ([2354f5b](https://gitea.com/gitea/runner/commits/2354f5bb18902124c0ccc07c632888bac8b4c99b))

### Bug Fixes

* call RunDaemon ([93e42d1](https://gitea.com/gitea/runner/commits/93e42d14650c8e1c55fbf37287c6f6e8ef2f24c5))
* ignore stuck because of wrong token ([88ae188](https://gitea.com/gitea/runner/commits/88ae188699428e10d0d5d886f00f07c100f4cd7b))
* limit workers ([5781e23](https://gitea.com/gitea/runner/commits/5781e233c14c6ff2c953be0727b1a9ed7977f4dc))
* need more time to run task ([cce60d4](https://gitea.com/gitea/runner/commits/cce60d466bb189ee83c92659aefb6354e34b9c56))
* **poller:** graceful shutdown ([9830f34](https://gitea.com/gitea/runner/commits/9830f34d364aa289e98b03f6d501a5afbb8379eb))
* **poller:** initial ready channel ([cab5699](https://gitea.com/gitea/runner/commits/cab56996b7b2f84723e57b0c2f0d80751570fc80))
* report job failure on error return ([#14](https://gitea.com/gitea/runner/issues/14)) ([f7a5278](https://gitea.com/gitea/runner/commits/f7a52789d950a57cf0c985d02715714740426e5c))
* **runner:** check task state field exist. ([2442cdd](https://gitea.com/gitea/runner/commits/2442cdd8adbd002ebadd982fd303879cbc73454b))
* set more image ([64c9856](https://gitea.com/gitea/runner/commits/64c9856ba7b05fe4e9cd77173a9b96332f7e0fec))
* set NoSkipCheckout ([2e1326b](https://gitea.com/gitea/runner/commits/2e1326b725bda03923e4f9f6d6b24a025efc5e3c))
* **typo:** Use uppercase N for default input prompt ([#9](https://gitea.com/gitea/runner/issues/9)) ([9ba9604](https://gitea.com/gitea/runner/commits/9ba9604256964f52ff102b16acf782d8c56d87ff))
* update container config ([b3e7d6a](https://gitea.com/gitea/runner/commits/b3e7d6a62d6c00850616fddc199bef6bb65e8a85))
* update step result ([41b5fa6](https://gitea.com/gitea/runner/commits/41b5fa6b557d5cb937690635306aea6749ee6c91))

### Refactoring

* Add graceful shutdown signal notify func ([e54de37](https://gitea.com/gitea/runner/commits/e54de3724292fa3fc4a50f1e1186964ce1e7a816))
* docker engine with options ([14a3455](https://gitea.com/gitea/runner/commits/14a345504f03c02a70b9323d157deccab236d86a))
* **task:** execute single task with gRPC data ([20c3d85](https://gitea.com/gitea/runner/commits/20c3d85ba9caef27fbf83118ce4b6126fed97580))

### Documentation

* Add Prerequisites ([#11](https://gitea.com/gitea/runner/issues/11)) ([7abe039](https://gitea.com/gitea/runner/commits/7abe039c311dd257f7ac4f0c569afe701c1ae6d0)), closes [#10](https://gitea.com/gitea/runner/issues/10)
* add quickstart to readme ([03c4c4e](https://gitea.com/gitea/runner/commits/03c4c4e1164e18306652ed5d18f7e9dcfb3b139a))
* **readme:** update some header. ([#1](https://gitea.com/gitea/runner/issues/1)) ([474683c](https://gitea.com/gitea/runner/commits/474683c0e8bca9bef135fc49d44e650144591dfb))

### Chores

* add .golangci config ([518aa1b](https://gitea.com/gitea/runner/commits/518aa1b2a76822e413293701af87be5f94b37c82))
* Add makefile ([a3e9bbe](https://gitea.com/gitea/runner/commits/a3e9bbed25680fcf46bfad48dd04e50565412b47))
* Add time sleep 1 second. ([bca586f](https://gitea.com/gitea/runner/commits/bca586ffd0bb9a726aed8b3d86e7c76b6ca0c3ef))
* **client:** support gRPC and gRPC web protocol ([d359276](https://gitea.com/gitea/runner/commits/d359276fe1d4203f84f21249f40bb942a29fc7a5))
* **config:** remove zerolog ([9225a3a](https://gitea.com/gitea/runner/commits/9225a3a856e85f63e1160bff0acbd3cb21f3d9d4))
* **engine:** ping docker daemon ([4eb5213](https://gitea.com/gitea/runner/commits/4eb5213294e4be894712f5a081cc18f01db308ac))
* fix lint error ([827bd95](https://gitea.com/gitea/runner/commits/827bd953e7d9a0dbe7caab149b898a926546a044))
* **go:** rename module name ([aa765ea](https://gitea.com/gitea/runner/commits/aa765ea4d0fe51a1211b755f31f96a1d5748aa8a))
* **gRPC:** add request interface client ([f2fb879](https://gitea.com/gitea/runner/commits/f2fb8798fa1556ad80d78a9649588a6089dfedcd))
* **gRPC:** add update and update step for client ([9be39b8](https://gitea.com/gitea/runner/commits/9be39b8cd4629678dac490eb5f0710cdd12ae00e))
* **gRPC:** fetch new stage data. ([631f801](https://gitea.com/gitea/runner/commits/631f801aa658ed0fa07df5025b067e20f4db4749))
* **gRPC:** register new runner ([449388f](https://gitea.com/gitea/runner/commits/449388f3ab1dae8944dd955288c77373036a96ac))
* load config from env file ([aa9ed05](https://gitea.com/gitea/runner/commits/aa9ed0590359c29181811b9f503921e414e465b2))
* **piepline:** add runtime package. ([7d55fd5](https://gitea.com/gitea/runner/commits/7d55fd57c94d1862f91ff7e15b892ab5f4b0af88))
* **poller:** add error retry limit ([45b0429](https://gitea.com/gitea/runner/commits/45b0429b2151c1dddaa722206bbd518f315b2006))
* **poller:** add metric to track the worker number ([abdb547](https://gitea.com/gitea/runner/commits/abdb547b1b5f5f3e5addc3c9c07b8eb0d9d1fc22))
* **poller:** Add poller package ([0b885c5](https://gitea.com/gitea/runner/commits/0b885c5e5fe33e2781755a8a02fa4bddfa997306))
* **proto:** Add ping request. ([4d7ef95](https://gitea.com/gitea/runner/commits/4d7ef95d402d3e53d175e262e01af224f4c58da7))
* **proto:** replace import path. ([307c2c6](https://gitea.com/gitea/runner/commits/307c2c63de7e486c785f191ec2e66e7a9b49df80))
* **register:** update proto file. ([c161a48](https://gitea.com/gitea/runner/commits/c161a48a0a5cdea220acc96b9e5df6c0b9c27e54))
* remove tracing codes ([2f879c4](https://gitea.com/gitea/runner/commits/2f879c41c4a9e44dded5e947b08f45b98976c50d))
* remove unused logrous package ([5213612](https://gitea.com/gitea/runner/commits/52136120333c5738feae5fc5f43c9577fdcbac3c))
* replace with code.gitea.io/bots-proto-go ([44988db](https://gitea.com/gitea/runner/commits/44988db9f89d2c55128688ec69309050aa140dcc))
* **runner:** add default label ""self-hosted" ([dada073](https://gitea.com/gitea/runner/commits/dada0730b04723483552337a2fbf1c40489571bb))
* **runner:** add new register command ([a503f74](https://gitea.com/gitea/runner/commits/a503f7429f82e1f55c8c13293469c85db5dfa39b))
* **runner:** add new token in header ([d178051](https://gitea.com/gitea/runner/commits/d178051832efa7b617b117c0d4c4dbe543d79c06))
* **runner:** add runner uuid in http header. ([bf5e3dc](https://gitea.com/gitea/runner/commits/bf5e3dc302e195dbe0b65688f5a642c372af0411))
* **runner:** cancel task if get the cancel from server ([08c94bb](https://gitea.com/gitea/runner/commits/08c94bb564709d87fe62a094d07f30578b16f594))
* **runner:** refactor register flow ([3a15031](https://gitea.com/gitea/runner/commits/3a1503138b7360ff143d7e48b09effc5ef45f3ab))
* **runner:** remove client secret and add UUID in runner ([e08495a](https://gitea.com/gitea/runner/commits/e08495af090634d69704f2d9667bdddce7c47684))
* **runner:** remove update runner status ([e9e42e8](https://gitea.com/gitea/runner/commits/e9e42e850b526fa87e5edde09b3cb60d57e9fbf2))
* **runner:** support update log and task ([82431d8](https://gitea.com/gitea/runner/commits/82431d8e114b75c9d4a02f009c74503f1dca0789))
* **runner:** update fetch job request ([7486d2a](https://gitea.com/gitea/runner/commits/7486d2ab918aacdb6c0bb3b85a3c7aa6230d1afb))
* **runner:** update runner status if start program ([d4c1515](https://gitea.com/gitea/runner/commits/d4c1515f4ef6ebe732d8bee9cd6f3e019d1560fb))
* **runner:** update runner status when start job ([7bebd2b](https://gitea.com/gitea/runner/commits/7bebd2bbadd10171b355306a25294cb9576ffb9b))
* **runner:** wait workflow done before shutdown the service ([d1114da](https://gitea.com/gitea/runner/commits/d1114da299be65d69a4d1bb7574851c1f6cd5d1a))
* **runtime:** check error message ([5051e4a](https://gitea.com/gitea/runner/commits/5051e4aebdbb51575fa77bf472f4ee4ce03dfd0c))
* **runtime:** fetch build data ([d3d56ed](https://gitea.com/gitea/runner/commits/d3d56ed0ef86b83745555adb849852ccd3caacb8))
* update grpc router path ([f05e08a](https://gitea.com/gitea/runner/commits/f05e08a76750d1f8459b75a3b8f1461f54e4aa1c))
* upgrade act ([66ec50e](https://gitea.com/gitea/runner/commits/66ec50e427cb8e1b7b217803636e80704ae27b06))
