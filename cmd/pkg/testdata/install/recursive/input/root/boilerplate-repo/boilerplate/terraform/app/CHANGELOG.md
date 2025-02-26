# Changelog

## [9.8.5](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.8.4...app-v9.8.5) (2025-02-06)


### Bug fixes

* remove old boilerplate hooks no longer needed ([#966](https://github.com/oslokommune/golden-path-boilerplate/issues/966)) ([6c08f80](https://github.com/oslokommune/golden-path-boilerplate/commit/6c08f80850f64dfd1d3080c52b3e74cf8911e7e2))

## [9.8.4](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.8.3...app-v9.8.4) (2025-02-06)


### Bug fixes

* terragrunt dependencies ([#926](https://github.com/oslokommune/golden-path-boilerplate/issues/926)) ([e6fb8a1](https://github.com/oslokommune/golden-path-boilerplate/commit/e6fb8a1e7fc0bb377d2ce1725e60a82308b403f9))

## [9.8.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.8.2...app-v9.8.3) (2025-02-04)


### Dependency updates

* update dependency terraform-aws-modules/security-group/aws to v5.3.0 ([#874](https://github.com/oslokommune/golden-path-boilerplate/issues/874)) ([e817d0b](https://github.com/oslokommune/golden-path-boilerplate/commit/e817d0bd82d9de8ad4f309e0d5c8a19952181ead))

## [9.8.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.8.1...app-v9.8.2) (2025-01-20)


### Bug fixes

* Output service URL for apex domain ([#892](https://github.com/oslokommune/golden-path-boilerplate/issues/892)) ([1818031](https://github.com/oslokommune/golden-path-boilerplate/commit/181803135b1eb2da37e3d04b16247137e59e584d))

## [9.8.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.8.0...app-v9.8.1) (2025-01-15)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.52.2 ([#878](https://github.com/oslokommune/golden-path-boilerplate/issues/878)) ([b8199c0](https://github.com/oslokommune/golden-path-boilerplate/commit/b8199c0f3c3b338780888a0d4b44e9a7deb642aa))

## [9.8.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.7.1...app-v9.8.0) (2025-01-09)


### Features

* Use original capacity when autoscaling service back up ([#838](https://github.com/oslokommune/golden-path-boilerplate/issues/838)) ([602f0c3](https://github.com/oslokommune/golden-path-boilerplate/commit/602f0c3dc0564824d3daf2796027f35160068893))
* use wget for nginx example image ([#864](https://github.com/oslokommune/golden-path-boilerplate/issues/864)) ([b5acc6e](https://github.com/oslokommune/golden-path-boilerplate/commit/b5acc6e744b0dd3391e377a0ef90efa33831d427))


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.42.0 ([#833](https://github.com/oslokommune/golden-path-boilerplate/issues/833)) ([dc59a1e](https://github.com/oslokommune/golden-path-boilerplate/commit/dc59a1e463587559599a2e9d7bacf8e9c1a124fa))
* update dependency alb-tg-host-routing-subdomain to v0.8.0 ([#867](https://github.com/oslokommune/golden-path-boilerplate/issues/867)) ([e713360](https://github.com/oslokommune/golden-path-boilerplate/commit/e7133600946e3da6b343eb525d6dabc38b3c2c11))
* update dependency terraform-aws-modules/iam/aws to v5.52.1 ([#862](https://github.com/oslokommune/golden-path-boilerplate/issues/862)) ([b73fe6e](https://github.com/oslokommune/golden-path-boilerplate/commit/b73fe6e855f2dd024dd0b3dfa14b4763e628106c))

## [9.7.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.7.0...app-v9.7.1) (2024-12-30)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.51.0 ([#856](https://github.com/oslokommune/golden-path-boilerplate/issues/856)) ([bc0fc6f](https://github.com/oslokommune/golden-path-boilerplate/commit/bc0fc6f6d7aff74cbb11d3194e8023e11fda66b0))

## [9.7.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.6.0...app-v9.7.0) (2024-12-23)


### Features

* correct dependency for terragrunt ([#840](https://github.com/oslokommune/golden-path-boilerplate/issues/840)) ([a5760b7](https://github.com/oslokommune/golden-path-boilerplate/commit/a5760b7c913bc2ba1c988c25e5deb99194a5afb8))

## [9.6.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.5.0...app-v9.6.0) (2024-12-19)


### Features

* require terraform version &gt;=1.7 for boilerplate ([#831](https://github.com/oslokommune/golden-path-boilerplate/issues/831)) ([77fcd59](https://github.com/oslokommune/golden-path-boilerplate/commit/77fcd5913643bf9af159e8ec6dbc86955825c0c4))

## [9.5.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.4.0...app-v9.5.0) (2024-12-13)


### Features

* make sure terragrunt_custom is not included if Terragrun is not enabled ([#810](https://github.com/oslokommune/golden-path-boilerplate/issues/810)) ([60c2e1d](https://github.com/oslokommune/golden-path-boilerplate/commit/60c2e1df0bbc163ffd29d063018e1548f8c7d615))

## [9.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.3.3...app-v9.4.0) (2024-12-12)


### Features

* Add experimental terragrunt to all terraform packages ([#764](https://github.com/oslokommune/golden-path-boilerplate/issues/764)) ([916c0cc](https://github.com/oslokommune/golden-path-boilerplate/commit/916c0cc2a8d9df21f651646e22d5ef912362710b))

## [9.3.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.3.2...app-v9.3.3) (2024-12-03)


### Dependency updates

* update dependency terraform-aws-modules/ecs/aws to v5.12.0 ([#778](https://github.com/oslokommune/golden-path-boilerplate/issues/778)) ([237fb10](https://github.com/oslokommune/golden-path-boilerplate/commit/237fb10c98df945ae1d4d1586b9671d0e5e704fc))

## [9.3.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.3.1...app-v9.3.2) (2024-11-15)


### Bug fixes

* Respect Alb domain settings in service_url ([#728](https://github.com/oslokommune/golden-path-boilerplate/issues/728)) ([df96d02](https://github.com/oslokommune/golden-path-boilerplate/commit/df96d025a3dddea01ccc088f979ad708239088e2))

## [9.3.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.3.0...app-v9.3.1) (2024-11-14)


### Bug fixes

* Use cross-platform compatible `find`-command ([#722](https://github.com/oslokommune/golden-path-boilerplate/issues/722)) ([eb8764f](https://github.com/oslokommune/golden-path-boilerplate/commit/eb8764fb29ffe0d0390577c19c2e9cdbe48071de))

## [9.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.2.0...app-v9.3.0) (2024-11-14)


### Features

* Rename module to alb_tg_host_routing_apex_domain ([#719](https://github.com/oslokommune/golden-path-boilerplate/issues/719)) ([fa996f1](https://github.com/oslokommune/golden-path-boilerplate/commit/fa996f167ce43f2947b19f054ec8cb4e368626c1))
* Rename module to alb_tg_host_routing_subdomain ([#720](https://github.com/oslokommune/golden-path-boilerplate/issues/720)) ([de8e4b7](https://github.com/oslokommune/golden-path-boilerplate/commit/de8e4b7866889b8e8e24105d191415fbdad853df))

## [9.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.1.0...app-v9.2.0) (2024-11-14)


### Features

* Rename alb-tg-host-routing files ([f5b7f63](https://github.com/oslokommune/golden-path-boilerplate/commit/f5b7f63f26fe4768520a96e4616028529f2f1123))

## [9.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.0.3...app-v9.1.0) (2024-11-14)


### Features

* Change module to alb-tg-host-routing-subdomain ([aef2dc1](https://github.com/oslokommune/golden-path-boilerplate/commit/aef2dc1ff390a318c0f6ed1fd58f7fd21d1b80e8))
* Change module to alb-tg-host-routing-subdomain ([f82e5ab](https://github.com/oslokommune/golden-path-boilerplate/commit/f82e5abee078ff2ace9bf7b2b13f6bd2197b0262))

## [9.0.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.0.2...app-v9.0.3) (2024-11-14)


### Dependency updates

* update dependency alb-tg-host-routing-apex to v0.3.0 ([#713](https://github.com/oslokommune/golden-path-boilerplate/issues/713)) ([b5970d0](https://github.com/oslokommune/golden-path-boilerplate/commit/b5970d00a5ca948004947b84369a9bed76a74b1e))

## [9.0.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.0.1...app-v9.0.2) (2024-11-13)


### Dependency updates

* update dependency alb-tg-host-routing to v0.6.0 ([#695](https://github.com/oslokommune/golden-path-boilerplate/issues/695)) ([2806839](https://github.com/oslokommune/golden-path-boilerplate/commit/28068397ef765c3868182bfa4b29719940614bd3))

## [9.0.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v9.0.0...app-v9.0.1) (2024-11-13)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.48.0 ([#696](https://github.com/oslokommune/golden-path-boilerplate/issues/696)) ([5742515](https://github.com/oslokommune/golden-path-boilerplate/commit/574251568367cd2f221d9fbeca0d1f8cae5cd9d6))

## [9.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.10...app-v9.0.0) (2024-11-13)


### ⚠ BREAKING CHANGES

* Support Apex ("root") domains in AlbHostRouting ([#681](https://github.com/oslokommune/golden-path-boilerplate/issues/681))

### Features

* Support Apex ("root") domains in AlbHostRouting ([#681](https://github.com/oslokommune/golden-path-boilerplate/issues/681)) ([778812a](https://github.com/oslokommune/golden-path-boilerplate/commit/778812af36ec8eb6b4395dcc11ef78350eb02d01))

## [8.3.10](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.9...app-v8.3.10) (2024-11-13)


### Bug fixes

* Add path to IAM roles, shorten IAM role name ([#699](https://github.com/oslokommune/golden-path-boilerplate/issues/699)) ([cea3f78](https://github.com/oslokommune/golden-path-boilerplate/commit/cea3f7826e90bdec8b27f8c0963b5e31022e76f1))

## [8.3.9](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.8...app-v8.3.9) (2024-11-05)


### Bug fixes

* fixes [#609](https://github.com/oslokommune/golden-path-boilerplate/issues/609) by creating GitHub environment when running app/bin/set_role* scripts ([#668](https://github.com/oslokommune/golden-path-boilerplate/issues/668)) ([b49e949](https://github.com/oslokommune/golden-path-boilerplate/commit/b49e9499e9f0b9a00a812f9800fb9c629b912d86))

## [8.3.8](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.7...app-v8.3.8) (2024-10-23)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.47.1 ([#636](https://github.com/oslokommune/golden-path-boilerplate/issues/636)) ([28d4b37](https://github.com/oslokommune/golden-path-boilerplate/commit/28d4b378ee1f6c39f6ad8092ae2df23af567ddff))

## [8.3.7](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.6...app-v8.3.7) (2024-10-22)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.47.0 ([#629](https://github.com/oslokommune/golden-path-boilerplate/issues/629)) ([4f26ad1](https://github.com/oslokommune/golden-path-boilerplate/commit/4f26ad17793ecf49c4ec60109684b1fda8d16524))

## [8.3.6](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.5...app-v8.3.6) (2024-10-17)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.41.1 ([#613](https://github.com/oslokommune/golden-path-boilerplate/issues/613)) ([65f8645](https://github.com/oslokommune/golden-path-boilerplate/commit/65f8645e797e79ca987d8001f725fbfee31b5b70))

## [8.3.5](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.4...app-v8.3.5) (2024-10-16)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.46.0 ([#479](https://github.com/oslokommune/golden-path-boilerplate/issues/479)) ([1d284f6](https://github.com/oslokommune/golden-path-boilerplate/commit/1d284f69d0ac3838cb2686a038b1b3a8a2d5d711))

## [8.3.4](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.3...app-v8.3.4) (2024-10-16)


### Dependency updates

* update dependency terraform-aws-modules/security-group/aws to v5.2.0 ([#483](https://github.com/oslokommune/golden-path-boilerplate/issues/483)) ([95ce33c](https://github.com/oslokommune/golden-path-boilerplate/commit/95ce33c67039292412155ad905d0a6ec22df5ac3))

## [8.3.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.2...app-v8.3.3) (2024-10-16)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.41.0 ([#586](https://github.com/oslokommune/golden-path-boilerplate/issues/586)) ([bdf7067](https://github.com/oslokommune/golden-path-boilerplate/commit/bdf706767c687a83ea9a054334c09c5955c21c74))

## [8.3.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.1...app-v8.3.2) (2024-10-16)


### Dependency updates

* update dependency alb-tg-host-routing to v0.5.2 ([#568](https://github.com/oslokommune/golden-path-boilerplate/issues/568)) ([0644e0d](https://github.com/oslokommune/golden-path-boilerplate/commit/0644e0ddabd5990e167995205798424587147a0b))

## [8.3.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.3.0...app-v8.3.1) (2024-10-14)


### Bug fixes

* revert git@github reference back to local reference ([#556](https://github.com/oslokommune/golden-path-boilerplate/issues/556)) ([66053b1](https://github.com/oslokommune/golden-path-boilerplate/commit/66053b15701ce9ffc5bcd04f643f4c70f830b3d1))

## [8.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.2.2...app-v8.3.0) (2024-10-11)


### Features

* Add attribute TargetGroupTargetStickiness ([#553](https://github.com/oslokommune/golden-path-boilerplate/issues/553)) ([290b922](https://github.com/oslokommune/golden-path-boilerplate/commit/290b9228093ecde0fe9f4c8dc788f19934c1784d))

## [8.2.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.2.1...app-v8.2.2) (2024-10-11)


### Bug fixes

* set versions to a release and not local path ([#539](https://github.com/oslokommune/golden-path-boilerplate/issues/539)) ([cfe0a9d](https://github.com/oslokommune/golden-path-boilerplate/commit/cfe0a9dbc0095f550c9bcdbc9306bf2566fb4d58))

## [8.2.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.2.0...app-v8.2.1) (2024-10-01)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.40.2 ([#513](https://github.com/oslokommune/golden-path-boilerplate/issues/513)) ([e0bd69a](https://github.com/oslokommune/golden-path-boilerplate/commit/e0bd69a7888deff56a6e4545f8779f25e5643f41))

## [8.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.1.0...app-v8.2.0) (2024-10-01)


### Features

* Make runtime platform configurable ([#512](https://github.com/oslokommune/golden-path-boilerplate/issues/512)) ([87b545d](https://github.com/oslokommune/golden-path-boilerplate/commit/87b545dd7841c374bad333ea53a4c9e394c0f4f9))

## [8.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.6...app-v8.1.0) (2024-09-25)


### Features

* Confirm before setting github secrets ([#516](https://github.com/oslokommune/golden-path-boilerplate/issues/516)) ([b369e70](https://github.com/oslokommune/golden-path-boilerplate/commit/b369e70227ff6e044e9139c0101424f52eb9f698))

## [8.0.6](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.5...app-v8.0.6) (2024-09-25)


### Bug fixes

* Delete bin folder if .IamForCicd is disabled ([#427](https://github.com/oslokommune/golden-path-boilerplate/issues/427)) ([aeeb050](https://github.com/oslokommune/golden-path-boilerplate/commit/aeeb0502e1265ce084aeb8778c6662521103d9c6))

## [8.0.5](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.4...app-v8.0.5) (2024-09-10)


### Dependency updates

* update dependency alb-tg-host-routing to v0.4.1 ([#471](https://github.com/oslokommune/golden-path-boilerplate/issues/471)) ([a5d98e9](https://github.com/oslokommune/golden-path-boilerplate/commit/a5d98e97239607e532856e584fe07a8e674e806d))

## [8.0.4](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.3...app-v8.0.4) (2024-09-03)


### Dependency updates

* update dependency terraform-aws-modules/ecs/aws to v5.11.4 ([#441](https://github.com/oslokommune/golden-path-boilerplate/issues/441)) ([6ef4f97](https://github.com/oslokommune/golden-path-boilerplate/commit/6ef4f97779897e7e86cd9cda74464409030d50c6))

## [8.0.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.2...app-v8.0.3) (2024-09-03)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.40.1 ([#453](https://github.com/oslokommune/golden-path-boilerplate/issues/453)) ([f93be86](https://github.com/oslokommune/golden-path-boilerplate/commit/f93be860a4a368fe26db7ed2e2f77dfca1cd8cad))

## [8.0.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.1...app-v8.0.2) (2024-08-19)


### Bug fixes

* Improve description of SG rule ([#422](https://github.com/oslokommune/golden-path-boilerplate/issues/422)) ([ef837cf](https://github.com/oslokommune/golden-path-boilerplate/commit/ef837cfca5b26587d0cd6d1f1d4b822e01851fa4))

## [8.0.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v8.0.0...app-v8.0.1) (2024-08-16)


### Bug fixes

* Include CICD shell files conditionally ([#417](https://github.com/oslokommune/golden-path-boilerplate/issues/417)) ([4a98524](https://github.com/oslokommune/golden-path-boilerplate/commit/4a98524d150efbefedf8726393365b456eb77b50))

## [8.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.3.1...app-v8.0.0) (2024-08-08)


### ⚠ BREAKING CHANGES

* Name data folders based on OutputFolder ([#405](https://github.com/oslokommune/golden-path-boilerplate/issues/405))

### Features

* Name data folders based on OutputFolder ([#405](https://github.com/oslokommune/golden-path-boilerplate/issues/405)) ([2f5039d](https://github.com/oslokommune/golden-path-boilerplate/commit/2f5039da0296e09f6b706b01441d9137089a31e6))

## [7.3.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.3.0...app-v7.3.1) (2024-07-31)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.42.0 ([#333](https://github.com/oslokommune/golden-path-boilerplate/issues/333)) ([b8004ae](https://github.com/oslokommune/golden-path-boilerplate/commit/b8004ae7ccccef3560fffca038080aa3236667aa))

## [7.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.2.2...app-v7.3.0) (2024-07-29)


### Features

* Add gh login validation ([c6f983d](https://github.com/oslokommune/golden-path-boilerplate/commit/c6f983ddb296597cd7662e8001c28cecb2d8e583))

## [7.2.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.2.1...app-v7.2.2) (2024-07-29)


### Bug fixes

* Make role name shorter and unique ([#373](https://github.com/oslokommune/golden-path-boilerplate/issues/373)) ([fbbf376](https://github.com/oslokommune/golden-path-boilerplate/commit/fbbf3763fc9af817217d0e648356131c8fd75aba))

## [7.2.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.2.0...app-v7.2.1) (2024-07-26)


### Bug fixes

* Use correct variables for iac script ([e2d5e67](https://github.com/oslokommune/golden-path-boilerplate/commit/e2d5e67a3b1afdbe7dbfd8450ca52dcc96a75b72))

## [7.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.1.0...app-v7.2.0) (2024-07-26)


### Features

* Remove script replaced by separate scripts ([3c3be2b](https://github.com/oslokommune/golden-path-boilerplate/commit/3c3be2b04cb639afc57b371da201ee3746684536))

## [7.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.0.1...app-v7.1.0) (2024-07-26)


### Features

* Add script for setting AWS_ROLE_ARNs for CICD ([#369](https://github.com/oslokommune/golden-path-boilerplate/issues/369)) ([1636463](https://github.com/oslokommune/golden-path-boilerplate/commit/1636463aec5d7eb3761b4dc4ff7afa9a29d97553))

## [7.0.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v7.0.0...app-v7.0.1) (2024-07-25)


### Dependency updates

* Bump iam-policies-generic from 3.2.0 to 4.0.0 ([#366](https://github.com/oslokommune/golden-path-boilerplate/issues/366)) ([880fcd2](https://github.com/oslokommune/golden-path-boilerplate/commit/880fcd2de37b47aa291d6cafd82a7e3dcb9ba917))

## [7.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.4.3...app-v7.0.0) (2024-07-24)


### ⚠ BREAKING CHANGES

* Add variable to assume CD role ([#347](https://github.com/oslokommune/golden-path-boilerplate/issues/347))

### Features

* Add variable to assume CD role ([#347](https://github.com/oslokommune/golden-path-boilerplate/issues/347)) ([a0611d9](https://github.com/oslokommune/golden-path-boilerplate/commit/a0611d9762fe0c61d5c0baab883033d888dd2627))

## [6.4.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.4.2...app-v6.4.3) (2024-07-24)


### Dependency updates

* update nginx/nginx:1.27-alpine3.19-slim docker digest to 44dbe45 ([#349](https://github.com/oslokommune/golden-path-boilerplate/issues/349)) ([ea5d552](https://github.com/oslokommune/golden-path-boilerplate/commit/ea5d55235d42609bf2ccb3427f937ee44f062b3f))

## [6.4.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.4.1...app-v6.4.2) (2024-07-24)


### Dependency updates

* Bump iam-policies-generic from 3.1.0 to 3.2.0 ([99da2ce](https://github.com/oslokommune/golden-path-boilerplate/commit/99da2cefdbc5d067e28845b848acb3b1c4b225ed))

## [6.4.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.4.0...app-v6.4.1) (2024-07-22)


### Bug fixes

* Resolve ExampleImage-induced errors during Terraform apply ([#345](https://github.com/oslokommune/golden-path-boilerplate/issues/345)) ([2dd70c2](https://github.com/oslokommune/golden-path-boilerplate/commit/2dd70c24500180ed1a1080b6db1e19386d4bc97d))

## [6.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.3.0...app-v6.4.0) (2024-07-12)


### Features

* Remove trailing comma in port mappings ([ae9dd24](https://github.com/oslokommune/golden-path-boilerplate/commit/ae9dd2417fffcfd2f54daaad90fe9e4e19528cec))

## [6.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.2.4...app-v6.3.0) (2024-07-09)


### Features

* Don't include image digest metadata when running example image ([#293](https://github.com/oslokommune/golden-path-boilerplate/issues/293)) ([470b6c2](https://github.com/oslokommune/golden-path-boilerplate/commit/470b6c237ef807e8879b3d8c3bba7b1fdd846a21))

## [6.2.4](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.2.3...app-v6.2.4) (2024-07-08)


### Dependency updates

* update dependency terraform-aws-modules/ecs/aws to v5.11.3 ([#221](https://github.com/oslokommune/golden-path-boilerplate/issues/221)) ([81f3975](https://github.com/oslokommune/golden-path-boilerplate/commit/81f3975789cd260b7d47c2e401289e42df8c06b8))

## [6.2.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.2.2...app-v6.2.3) (2024-07-05)


### Dependency updates

* update dependency terraform-aws-modules/iam/aws to v5.39.1 ([#222](https://github.com/oslokommune/golden-path-boilerplate/issues/222)) ([96ac8a9](https://github.com/oslokommune/golden-path-boilerplate/commit/96ac8a905a08126bf1b19639db3722c50d09b29a))

## [6.2.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.2.1...app-v6.2.2) (2024-07-05)


### Dependency updates

* update nginx/nginx:1.27-alpine3.19-slim docker digest to 97ddddd ([#235](https://github.com/oslokommune/golden-path-boilerplate/issues/235)) ([cf7ebd9](https://github.com/oslokommune/golden-path-boilerplate/commit/cf7ebd96c6c9ea20515e0bfd1256d4cf80ae3304))

## [6.2.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.2.0...app-v6.2.1) (2024-07-05)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.40.0 ([#274](https://github.com/oslokommune/golden-path-boilerplate/issues/274)) ([4bdbcbf](https://github.com/oslokommune/golden-path-boilerplate/commit/4bdbcbfeb06fa3cd2e8b60a3e956106795232d54))

## [6.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.1.1...app-v6.2.0) (2024-07-02)


### Features

* add AlbHostRouting.Internal flag to look up correct ALB in __gp_dependencies.tf ([#255](https://github.com/oslokommune/golden-path-boilerplate/issues/255)) ([33fb4bd](https://github.com/oslokommune/golden-path-boilerplate/commit/33fb4bd54685d4745043663a7c31c8b6c0d5e115))

## [6.1.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.1.0...app-v6.1.1) (2024-06-13)


### Dependency updates

* update aws-observability/aws-otel-collector docker tag to v0.39.1 ([#216](https://github.com/oslokommune/golden-path-boilerplate/issues/216)) ([6967e13](https://github.com/oslokommune/golden-path-boilerplate/commit/6967e13e7f7707b576512bb4fc215d55881168b0))
* update nginx/nginx:1.27-alpine3.19-slim docker digest to c143d9e ([#215](https://github.com/oslokommune/golden-path-boilerplate/issues/215)) ([1d09dc1](https://github.com/oslokommune/golden-path-boilerplate/commit/1d09dc1a4bf775afb0c9387ec4572a7f612a8284))

## [6.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v6.0.0...app-v6.1.0) (2024-06-13)


### Features

* Pin container images to digest and more specific version tags ([#212](https://github.com/oslokommune/golden-path-boilerplate/issues/212)) ([1f3c8d6](https://github.com/oslokommune/golden-path-boilerplate/commit/1f3c8d6f7618cfa7ac3a86611fa9fd82d4faea37))

## [6.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v5.2.2...app-v6.0.0) (2024-06-06)


### ⚠ BREAKING CHANGES

* New directory structure ([#154](https://github.com/oslokommune/golden-path-boilerplate/issues/154))

### Features

* New directory structure ([#154](https://github.com/oslokommune/golden-path-boilerplate/issues/154)) ([3e34dd1](https://github.com/oslokommune/golden-path-boilerplate/commit/3e34dd1e3e5e0e3e1bc35359412809ca16dc199d))

## [5.2.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v5.2.1...app-v5.2.2) (2024-05-31)


### Bug fixes

* Indentation and wording in templates ([#151](https://github.com/oslokommune/golden-path-boilerplate/issues/151)) ([678c37c](https://github.com/oslokommune/golden-path-boilerplate/commit/678c37c7b92f4d6d6794a6d4de8b5c007d30f7dc))

## [5.2.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v5.2.0...app-v5.2.1) (2024-05-31)


### Bug fixes

* Use substr instead of path based role names ([#137](https://github.com/oslokommune/golden-path-boilerplate/issues/137)) ([53618b2](https://github.com/oslokommune/golden-path-boilerplate/commit/53618b2826bd10afbe167afe385255bf1d9f6d69))

## [5.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v5.1.0...app-v5.2.0) (2024-05-31)


### Features

* add name to version file ([#147](https://github.com/oslokommune/golden-path-boilerplate/issues/147)) ([0ceb454](https://github.com/oslokommune/golden-path-boilerplate/commit/0ceb454ede5dff40aa0acf8495d966dc93be219a))
* Move app to new pattern ([#144](https://github.com/oslokommune/golden-path-boilerplate/issues/144)) ([9269f08](https://github.com/oslokommune/golden-path-boilerplate/commit/9269f089012aadddbedd81390aef0fc638118373))

## [5.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v5.0.0...app-v5.1.0) (2024-05-27)


### Features

* Set path for ECS IAM roles ([#122](https://github.com/oslokommune/golden-path-boilerplate/issues/122)) ([a54e3a4](https://github.com/oslokommune/golden-path-boilerplate/commit/a54e3a40212afffad41fff80107484f33c9e9378))

## [5.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.6.0...app-v5.0.0) (2024-05-24)


### ⚠ BREAKING CHANGES

* Make it possible to point to any ALB
* Remove defaults in config ([#117](https://github.com/oslokommune/golden-path-boilerplate/issues/117))

### Features

* Make it possible to point to any ALB ([711ab0a](https://github.com/oslokommune/golden-path-boilerplate/commit/711ab0ae557dbce13fb370e21ebd43944fcf2c1e))
* Remove defaults in config ([#117](https://github.com/oslokommune/golden-path-boilerplate/issues/117)) ([4ae14fa](https://github.com/oslokommune/golden-path-boilerplate/commit/4ae14fa487cab39110f4569a2f6042b346d2217b))


### Bug fixes

* Move autoscaling outside IAM ([921fcef](https://github.com/oslokommune/golden-path-boilerplate/commit/921fcef8fc39b12f9befb06678eda0ac03bd5091))

## [4.6.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.5.1...app-v4.6.0) (2024-05-23)


### Features

* Add parameter for desired count ([#112](https://github.com/oslokommune/golden-path-boilerplate/issues/112)) ([0bfed96](https://github.com/oslokommune/golden-path-boilerplate/commit/0bfed96311badc2a9bb1168ee4ce7485ea34e9a5))

## [4.5.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.5.0...app-v4.5.1) (2024-05-22)


### Bug fixes

* Skip CHANGELOG*.md ([728ed3d](https://github.com/oslokommune/golden-path-boilerplate/commit/728ed3ddf68d29a01368662cde144638b22e4a6e))

## [4.5.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.4.0...app-v4.5.0) (2024-05-22)


### Features

* Bump to data-networking-v0.2.0 ([9388f8c](https://github.com/oslokommune/golden-path-boilerplate/commit/9388f8c89420bed26d2133f3d77ed552614a33ec))


### Bug fixes

* Use conditional for aws_security_group.alb_public ([4cbbedc](https://github.com/oslokommune/golden-path-boilerplate/commit/4cbbedcf3d9ce91fc8c306976bc4db291a2d4963))

## [4.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.3.0...app-v4.4.0) (2024-05-21)


### Features

* Use release-please for all boilerplate templates ([#89](https://github.com/oslokommune/golden-path-boilerplate/issues/89)) ([e380d58](https://github.com/oslokommune/golden-path-boilerplate/commit/e380d58c9a0273bfb4667c6228555784a4e3c6ad))

## [4.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.2.0...app-v4.3.0) (2024-05-16)


### Features

* Comment where to find variables ([f6ab807](https://github.com/oslokommune/golden-path-boilerplate/commit/f6ab807d16494e92b254476bf0c2b9a58b7b1a4d))

## [4.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.1.0...app-v4.2.0) (2024-05-15)


### Features

* Bump CI/CD IAM policies module ([#97](https://github.com/oslokommune/golden-path-boilerplate/issues/97)) ([ce76bf5](https://github.com/oslokommune/golden-path-boilerplate/commit/ce76bf51cd2e7647091405f149e82237725ba284))

## [4.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.0.1...app-v4.1.0) (2024-05-08)


### Features

* Move section "Load balancing" down ([#90](https://github.com/oslokommune/golden-path-boilerplate/issues/90)) ([5ca1207](https://github.com/oslokommune/golden-path-boilerplate/commit/5ca120775d2ab7d33f840b66460b8ccaa17e9744))


### Bug fixes

* Remove accidentally added variable ([dda487d](https://github.com/oslokommune/golden-path-boilerplate/commit/dda487d4e13a817a30e8a20d7e1bea5edb2c417d))
* Use correct principal ARN for trusted role ARN ([#92](https://github.com/oslokommune/golden-path-boilerplate/issues/92)) ([01baef6](https://github.com/oslokommune/golden-path-boilerplate/commit/01baef6512cf69c245c485efdf54d876128f74a8))

## [4.0.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v4.0.0...app-v4.0.1) (2024-05-07)


### Bug fixes

* Support non-defined variable for AssumableCdRole ([d5760cb](https://github.com/oslokommune/golden-path-boilerplate/commit/d5760cb48d05526a0ba869eef7066019bdfa8d82))

## [4.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v3.3.0...app-v4.0.0) (2024-05-07)


### ⚠ BREAKING CHANGES

* Simplify IAM for CI/CD ([#83](https://github.com/oslokommune/golden-path-boilerplate/issues/83))

### Features

* Debug IAM for CD locally ([#85](https://github.com/oslokommune/golden-path-boilerplate/issues/85)) ([6ff69fc](https://github.com/oslokommune/golden-path-boilerplate/commit/6ff69fc75eee1aae787093ca784a0f1c574b57ca))
* Simplify IAM for CI/CD ([#83](https://github.com/oslokommune/golden-path-boilerplate/issues/83)) ([fc6da52](https://github.com/oslokommune/golden-path-boilerplate/commit/fc6da527b8bf7d4f0e7dbbc02f6fa26804a72f47))

## [3.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v3.2.1...app-v3.3.0) (2024-05-06)


### Features

* Add config for multiple load balancers ([#79](https://github.com/oslokommune/golden-path-boilerplate/issues/79)) ([07b0cd8](https://github.com/oslokommune/golden-path-boilerplate/commit/07b0cd806845783a9c87c76e2eb0bb8c24f5a151))

## [3.2.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v3.2.0...app-v3.2.1) (2024-04-30)


### Bug fixes

* Use correct ARNs for AWS Secrets Manager secrets ([8728bb4](https://github.com/oslokommune/golden-path-boilerplate/commit/8728bb4287bda1d162090e77b91c697e6c0dce3b))

## [3.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v3.1.0...app-v3.2.0) (2024-04-29)


### Features

* Use latest IAM module ([714af6c](https://github.com/oslokommune/golden-path-boilerplate/commit/714af6c39db40f0fd1658b90559118613e3c7f1f))
* Use latest IAM module ([8d0ced6](https://github.com/oslokommune/golden-path-boilerplate/commit/8d0ced6b7811a8be6c8bd3ea31867b5f516bf6ee))


### Bug fixes

* Use proper naming conventions for log groups ([05f0daa](https://github.com/oslokommune/golden-path-boilerplate/commit/05f0daa7d5c463d40f9a5911abad30fe3163c737))

## [3.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v3.0.0...app-v3.1.0) (2024-04-29)


### Features

* Rearrange ECS service sections for improved UX ([7598c73](https://github.com/oslokommune/golden-path-boilerplate/commit/7598c739205ed6c0b8af1829232e5a4803664a71))

## [3.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.6.0...app-v3.0.0) (2024-04-26)


### ⚠ BREAKING CHANGES

* New file pattern

### Features

* New file pattern ([0d278d8](https://github.com/oslokommune/golden-path-boilerplate/commit/0d278d86fc460a8ba1c49d5ece3464bf24c34ab3))
* Section code with code banners in config file ([524e5e5](https://github.com/oslokommune/golden-path-boilerplate/commit/524e5e566cd84530acddb3fe9455e9f9f89fa168))


### Bug fixes

* Fix bug where creating a new app stack would fail due to mv ([a561462](https://github.com/oslokommune/golden-path-boilerplate/commit/a561462c401e50a60e861955fc0b8de43bf8723c))

## [2.6.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.5.1...app-v2.6.0) (2024-04-24)


### Features

* Support custom route53 zone name ([04abf2b](https://github.com/oslokommune/golden-path-boilerplate/commit/04abf2be4cffc420e6bac8442abe84a0c8c3674f))

## [2.5.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.5.0...app-v2.5.1) (2024-04-16)


### Bug fixes

* **app:** Add missing parentheses fix comma placement ([f4c661b](https://github.com/oslokommune/golden-path-boilerplate/commit/f4c661b6fcd013c72c30d726e2914fa5cfc48cbb))

## [2.5.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.4.0...app-v2.5.0) (2024-04-16)


### Features

* **app:** Add missing conditionals for DB and ALB ([58c1a95](https://github.com/oslokommune/golden-path-boilerplate/commit/58c1a953aebda1daa2d81dd9445635193ec9d878))
* **app:** Update alb-tg-host-routing module to v0.3.1 ([988d932](https://github.com/oslokommune/golden-path-boilerplate/commit/988d93292192c1f01185a14784f9777d6fb6db13))

## [2.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.3.1...app-v2.4.0) (2024-04-15)


### Features

* Bump to alb-tg-host-routing-v0.2.1 ([276202a](https://github.com/oslokommune/golden-path-boilerplate/commit/276202a99ba6402819324c66755917139968d1d3))

## [2.3.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.3.0...app-v2.3.1) (2024-04-15)


### Bug fixes

* Use healthcheck path / for ExampleImage ([#51](https://github.com/oslokommune/golden-path-boilerplate/issues/51)) ([9b70823](https://github.com/oslokommune/golden-path-boilerplate/commit/9b70823cad9fdd06f4e46d15cd1152bec6ed2cbb))
* Use port 80 for ExampleImage ([#48](https://github.com/oslokommune/golden-path-boilerplate/issues/48)) ([def8a2a](https://github.com/oslokommune/golden-path-boilerplate/commit/def8a2ad6b72609a84bd43f44ee79084a1f4a23d))

## [2.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.2.0...app-v2.3.0) (2024-04-12)


### Features

* Move policies and roles to separate files ([cff9654](https://github.com/oslokommune/golden-path-boilerplate/commit/cff9654c0ae0d5a164f940dbfae9fde2c99976c1))

## [2.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.1.0...app-v2.2.0) (2024-04-11)


### Features

* IAM for CICD ([#43](https://github.com/oslokommune/golden-path-boilerplate/issues/43)) ([f176e11](https://github.com/oslokommune/golden-path-boilerplate/commit/f176e11075c11c3cd9fb036813a040ad1b3eec11))

## [2.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v2.0.0...app-v2.1.0) (2024-04-11)


### Features

* Update default_main_container_port to 8080 instead of using vars ([d4c7742](https://github.com/oslokommune/golden-path-boilerplate/commit/d4c7742cb41d8e91499e456562c3115bc4fb8a7f))


### Bug fixes

* Add missing _config_override.tf ([179696c](https://github.com/oslokommune/golden-path-boilerplate/commit/179696c366c604ddda98b7020704d396f80e9830))

## [2.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.10.0...app-v2.0.0) (2024-04-11)


### ⚠ BREAKING CHANGES

* Rename almost all variables and use maps ([#40](https://github.com/oslokommune/golden-path-boilerplate/issues/40))

### Features

* Rename almost all variables and use maps ([#40](https://github.com/oslokommune/golden-path-boilerplate/issues/40)) ([c9dfab2](https://github.com/oslokommune/golden-path-boilerplate/commit/c9dfab2dd45728e5e758ab79d127ef7d1c245b78))

## [1.10.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.9.0...app-v1.10.0) (2024-04-11)


### Features

* Move app data stores to separate stack ([#38](https://github.com/oslokommune/golden-path-boilerplate/issues/38)) ([b547f3d](https://github.com/oslokommune/golden-path-boilerplate/commit/b547f3d4549815bbcdb31b83b9c8e124f8956186))

## [1.9.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.8.0...app-v1.9.0) (2024-04-08)


### Features

* Add config for CPU and memory ([5c2d087](https://github.com/oslokommune/golden-path-boilerplate/commit/5c2d08795c4e6fff4c36e25503571e0395bd20ba))
* Add custom health check ([#32](https://github.com/oslokommune/golden-path-boilerplate/issues/32)) ([c045a96](https://github.com/oslokommune/golden-path-boilerplate/commit/c045a968237275425293fed317c4fcfb6a07659f))
* Add database to app template ([02df28d](https://github.com/oslokommune/golden-path-boilerplate/commit/02df28dcf19fbe450b2f83c13cc59f37733c1c60))
* Add proper support for mount points ([3650dc8](https://github.com/oslokommune/golden-path-boilerplate/commit/3650dc8c4ecd9b0a1f892d1d3c92a5b94a049a07))
* Add support for ECS Exec with read only file systems ([f1865c8](https://github.com/oslokommune/golden-path-boilerplate/commit/f1865c8ec90c679f6b2b8b48c661305be6c9a705))
* Bump ecs_service module version ([9aced10](https://github.com/oslokommune/golden-path-boilerplate/commit/9aced100a0c03792534d24d0d7326038f19b0c1c))
* Move variables to locals ([#34](https://github.com/oslokommune/golden-path-boilerplate/issues/34)) ([a8ee071](https://github.com/oslokommune/golden-path-boilerplate/commit/a8ee071bc3611b654eb30e58554a66995585a453))
* New system for config overrides ([#36](https://github.com/oslokommune/golden-path-boilerplate/issues/36)) ([7ad41b7](https://github.com/oslokommune/golden-path-boilerplate/commit/7ad41b76f4791406be5e62f9383b1c5d75212182))
* Set alb_listener_priority = null ([d5c05f4](https://github.com/oslokommune/golden-path-boilerplate/commit/d5c05f469637be4040dc5692fe57af77852efe19))
* Support extra container definitions ([f96cdde](https://github.com/oslokommune/golden-path-boilerplate/commit/f96cddee7713ad27800c4633c4fddec4c41669a8))
* Support for volumes ([3a4f729](https://github.com/oslokommune/golden-path-boilerplate/commit/3a4f729c4815b335b40ae52851df19314d8d19c9))
* Support prometheus_* config overrides ([1531c4f](https://github.com/oslokommune/golden-path-boilerplate/commit/1531c4f5f1ec42ceb3cc0ce12741f2179122b894))


### Bug fixes

* More whitespaces and linebreaks for app template ([89a75c3](https://github.com/oslokommune/golden-path-boilerplate/commit/89a75c352cab31e3ec65992044974563e18d8b7f))
* Whitespaces and linebreaks for app template ([a223748](https://github.com/oslokommune/golden-path-boilerplate/commit/a2237483fc1b71a1907d1a4076b33ff369b210cf))

## [1.8.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.7.1...app-v1.8.0) (2024-03-27)


### Features

* Use partials in app template ([64b8483](https://github.com/oslokommune/golden-path-boilerplate/commit/64b8483b0aecd869befe86e1dc2774503da97f1d))

## [1.7.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.7.0...app-v1.7.1) (2024-03-26)


### Bug fixes

* Add --no-run-if-empty to xargs cmd ([24b9a3f](https://github.com/oslokommune/golden-path-boilerplate/commit/24b9a3ff82ccea128eaa6decc1f2c8fb339c514a))

## [1.7.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.6.0...app-v1.7.0) (2024-03-26)


### Features

* Add X-Ray to app template ([fc11344](https://github.com/oslokommune/golden-path-boilerplate/commit/fc11344330ef05a11b6edb584536cf4cbde58d35))
* Bump to aws-otel-collector-ecs-sidecar-v1.1.1 ([99a2bfb](https://github.com/oslokommune/golden-path-boilerplate/commit/99a2bfbf29095e3e839db069410ba9aefbcd573a))

## [1.6.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.5.0...app-v1.6.0) (2024-03-22)


### Features

* Add experimental support for env vars ([325ad82](https://github.com/oslokommune/golden-path-boilerplate/commit/325ad82ee75b35a5028ee78796dc3f75cfe56fba))

## [1.5.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.4.0...app-v1.5.0) (2024-03-21)


### Features

* Add daily shutdown to template ([c3ba09d](https://github.com/oslokommune/golden-path-boilerplate/commit/c3ba09d8861f89bc15c735e749d11d1647647833))

## [1.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.3.0...app-v1.4.0) (2024-03-21)


### Features

* Add AppPort variable ([07c8aa5](https://github.com/oslokommune/golden-path-boilerplate/commit/07c8aa50cba8c34298bb9ca93d6fdcda0d464be8))
* Add support for ECS Service Connect ([83cf71b](https://github.com/oslokommune/golden-path-boilerplate/commit/83cf71b53cb1475c6c2fd2635c10b4c9db5fc70a))
* Add toggle for running with or without example image ([f41ee84](https://github.com/oslokommune/golden-path-boilerplate/commit/f41ee8470f4e5835caa04830dc1be3eb60d2945f))
* Put Boilerplate metadata in a hidden directory ([27e03e5](https://github.com/oslokommune/golden-path-boilerplate/commit/27e03e58f58ef4b0b25c38547458080af42e729a))
* Rename variable ([17f3a96](https://github.com/oslokommune/golden-path-boilerplate/commit/17f3a96c611b72f03427ecfc85ab6e3d7ca0e0c4))


### Bug fixes

* Correct variable type ([fa56ed3](https://github.com/oslokommune/golden-path-boilerplate/commit/fa56ed331715640e152d146a68c03a80c80e858b))
* Fix use of wrong variables ([330cb40](https://github.com/oslokommune/golden-path-boilerplate/commit/330cb40604a2b1548b28c30d029b3aaf868ab506))

## [1.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.2.0...app-v1.3.0) (2024-03-19)


### Features

* Use Nginx Alpine image which seems to include curl that we need for the health checks ([7575183](https://github.com/oslokommune/golden-path-boilerplate/commit/75751831b725bc3fca3b74e80745bb201edcea84))
* Various cleanups for app template ([259f5a9](https://github.com/oslokommune/golden-path-boilerplate/commit/259f5a9bc0bb4f9b5f63cd9900075cd56b09dac8))
* Various improvements to app template ([2ebbcdf](https://github.com/oslokommune/golden-path-boilerplate/commit/2ebbcdf45693a59b152c3f33e463a3b7a4c0d7bc))

## [1.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-v1.1.0...app-v1.2.0) (2024-03-18)


### Features

* Add app template ([5cec464](https://github.com/oslokommune/golden-path-boilerplate/commit/5cec4645e69dd5adca04f8900b3cda8647180590))
* Add support for OpenTelemetry Collector sidecar ([ed7fd44](https://github.com/oslokommune/golden-path-boilerplate/commit/ed7fd44c41d88b210a96af46c6f529c81e74e501))
* Add terraform fmt hook ([f62262a](https://github.com/oslokommune/golden-path-boilerplate/commit/f62262a43b7549a3a90e20f3a5394529fd05454b))
* DRY up the template a bit ([80a7252](https://github.com/oslokommune/golden-path-boilerplate/commit/80a7252ca471ec43376b02939638a56b22f0fa27))
* Make template code easier to read ([093438c](https://github.com/oslokommune/golden-path-boilerplate/commit/093438c2d5b8b0efff0b1a161f74b39143d0d168))
* Makk app template work for simple app ([583d975](https://github.com/oslokommune/golden-path-boilerplate/commit/583d975f9119c3158e54c46031ff18d805cf6189))
* Mark files for deletion and run a hook to remove them ([3f79c90](https://github.com/oslokommune/golden-path-boilerplate/commit/3f79c90d28eb479b3d0849c91e878c1cc7df6ede))
* Rename template files ([77c6fa7](https://github.com/oslokommune/golden-path-boilerplate/commit/77c6fa7031d71305cd2e7baadc283817e83d7c00))
