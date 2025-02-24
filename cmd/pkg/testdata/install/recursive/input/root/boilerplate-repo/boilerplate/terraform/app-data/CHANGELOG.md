# Changelog

## [4.5.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.4.1...app-data-v4.5.0) (2025-02-06)


### Features

* fixes [#960](https://github.com/oslokommune/golden-path-boilerplate/issues/960) by adding a new Ecr.Enable option ([#961](https://github.com/oslokommune/golden-path-boilerplate/issues/961)) ([9d593c5](https://github.com/oslokommune/golden-path-boilerplate/commit/9d593c52b8aa06d952fc1c0fbaf307e425137556))

## [4.4.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.4.0...app-data-v4.4.1) (2025-02-06)


### Bug fixes

* terragrunt dependencies ([#926](https://github.com/oslokommune/golden-path-boilerplate/issues/926)) ([e6fb8a1](https://github.com/oslokommune/golden-path-boilerplate/commit/e6fb8a1e7fc0bb377d2ce1725e60a82308b403f9))

## [4.4.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.3.0...app-data-v4.4.0) (2024-12-19)


### Features

* require terraform version &gt;=1.7 for boilerplate ([#831](https://github.com/oslokommune/golden-path-boilerplate/issues/831)) ([77fcd59](https://github.com/oslokommune/golden-path-boilerplate/commit/77fcd5913643bf9af159e8ec6dbc86955825c0c4))

## [4.3.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.2.0...app-data-v4.3.0) (2024-12-13)


### Features

* make sure terragrunt_custom is not included if Terragrun is not enabled ([#810](https://github.com/oslokommune/golden-path-boilerplate/issues/810)) ([60c2e1d](https://github.com/oslokommune/golden-path-boilerplate/commit/60c2e1df0bbc163ffd29d063018e1548f8c7d615))

## [4.2.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.1.3...app-data-v4.2.0) (2024-12-12)


### Features

* Add experimental terragrunt to all terraform packages ([#764](https://github.com/oslokommune/golden-path-boilerplate/issues/764)) ([916c0cc](https://github.com/oslokommune/golden-path-boilerplate/commit/916c0cc2a8d9df21f651646e22d5ef912362710b))

## [4.1.3](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.1.2...app-data-v4.1.3) (2024-11-14)


### Bug fixes

* Use cross-platform compatible `find`-command ([#722](https://github.com/oslokommune/golden-path-boilerplate/issues/722)) ([eb8764f](https://github.com/oslokommune/golden-path-boilerplate/commit/eb8764fb29ffe0d0390577c19c2e9cdbe48071de))

## [4.1.2](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.1.1...app-data-v4.1.2) (2024-10-14)


### Bug fixes

* revert git@github reference back to local reference ([#556](https://github.com/oslokommune/golden-path-boilerplate/issues/556)) ([66053b1](https://github.com/oslokommune/golden-path-boilerplate/commit/66053b15701ce9ffc5bcd04f643f4c70f830b3d1))

## [4.1.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.1.0...app-data-v4.1.1) (2024-10-11)


### Bug fixes

* set versions to a release and not local path ([#539](https://github.com/oslokommune/golden-path-boilerplate/issues/539)) ([cfe0a9d](https://github.com/oslokommune/golden-path-boilerplate/commit/cfe0a9dbc0095f550c9bcdbc9306bf2566fb4d58))

## [4.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v4.0.0...app-data-v4.1.0) (2024-10-10)


### Features

* Add local variable `ecr_max_image_count` ([#547](https://github.com/oslokommune/golden-path-boilerplate/issues/547)) ([b1bb138](https://github.com/oslokommune/golden-path-boilerplate/commit/b1bb138ec795ead371efd8ed60b9b835f5212d98))

## [4.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v3.0.0...app-data-v4.0.0) (2024-07-24)


### ⚠ BREAKING CHANGES

* Add variable to assume CD role ([#347](https://github.com/oslokommune/golden-path-boilerplate/issues/347))

### Features

* Add variable to assume CD role ([#347](https://github.com/oslokommune/golden-path-boilerplate/issues/347)) ([a0611d9](https://github.com/oslokommune/golden-path-boilerplate/commit/a0611d9762fe0c61d5c0baab883033d888dd2627))

## [3.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v2.1.1...app-data-v3.0.0) (2024-06-06)


### ⚠ BREAKING CHANGES

* New directory structure ([#154](https://github.com/oslokommune/golden-path-boilerplate/issues/154))

### Features

* New directory structure ([#154](https://github.com/oslokommune/golden-path-boilerplate/issues/154)) ([3e34dd1](https://github.com/oslokommune/golden-path-boilerplate/commit/3e34dd1e3e5e0e3e1bc35359412809ca16dc199d))

## [2.1.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v2.1.0...app-data-v2.1.1) (2024-05-31)


### Bug fixes

* Indentation and wording in templates ([#151](https://github.com/oslokommune/golden-path-boilerplate/issues/151)) ([678c37c](https://github.com/oslokommune/golden-path-boilerplate/commit/678c37c7b92f4d6d6794a6d4de8b5c007d30f7dc))

## [2.1.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v2.0.0...app-data-v2.1.0) (2024-05-31)


### Features

* add name to version file ([#147](https://github.com/oslokommune/golden-path-boilerplate/issues/147)) ([0ceb454](https://github.com/oslokommune/golden-path-boilerplate/commit/0ceb454ede5dff40aa0acf8495d966dc93be219a))

## [2.0.0](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v1.0.1...app-data-v2.0.0) (2024-05-30)


### ⚠ BREAKING CHANGES

* Move app-data files to new pattern ([#129](https://github.com/oslokommune/golden-path-boilerplate/issues/129))

### Features

* Move app-data files to new pattern ([#129](https://github.com/oslokommune/golden-path-boilerplate/issues/129)) ([3960bda](https://github.com/oslokommune/golden-path-boilerplate/commit/3960bdadd2701b42118a17f3d7dacb9b2b958642))

## [1.0.1](https://github.com/oslokommune/golden-path-boilerplate/compare/app-data-v1.0.0...app-data-v1.0.1) (2024-05-22)


### Bug fixes

* Skip CHANGELOG*.md ([728ed3d](https://github.com/oslokommune/golden-path-boilerplate/commit/728ed3ddf68d29a01368662cde144638b22e4a6e))

## 1.0.0 (2024-05-21)


### ⚠ BREAKING CHANGES

* New file pattern

### Features

* **app-data:** Remove obsolete import.tf ([4fb598d](https://github.com/oslokommune/golden-path-boilerplate/commit/4fb598d02d83584686f596b0f791170a96526d2d))
* New file pattern ([0d278d8](https://github.com/oslokommune/golden-path-boilerplate/commit/0d278d86fc460a8ba1c49d5ece3464bf24c34ab3))
* Use release-please for all boilerplate templates ([#89](https://github.com/oslokommune/golden-path-boilerplate/issues/89)) ([e380d58](https://github.com/oslokommune/golden-path-boilerplate/commit/e380d58c9a0273bfb4667c6228555784a4e3c6ad))


### Bug fixes

* Add missing header ([#86](https://github.com/oslokommune/golden-path-boilerplate/issues/86)) ([77fa6fc](https://github.com/oslokommune/golden-path-boilerplate/commit/77fa6fc8c17064b4220575ad3763d2fe5dc27198))
* Use proper naming conventions for log groups ([05f0daa](https://github.com/oslokommune/golden-path-boilerplate/commit/05f0daa7d5c463d40f9a5911abad30fe3163c737))
