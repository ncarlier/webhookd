<a name=""></a>
# [](https://github.com/ncarlier/webhookd/compare/v1.11.0...v) (2021-05-29)


### Features

* output hook execution logs to server logs ([07fbb6e](https://github.com/ncarlier/webhookd/commit/07fbb6e)), closes [#44](https://github.com/ncarlier/webhookd/issues/44)
* **notification:** email subject customization ([ed67fc7](https://github.com/ncarlier/webhookd/commit/ed67fc7))


### BREAKING CHANGES

* for configuration

- rename `-log-dir` to `-hook-log-dir`
- add `-hook-log-output`
- add `-log-level`
- remove `-debug`



<a name="1.11.0"></a>
# [1.11.0](https://github.com/ncarlier/webhookd/compare/v1.10.0...v1.11.0) (2020-10-14)


### Bug Fixes

* fix timeout parameter name ([f5f4838](https://github.com/ncarlier/webhookd/commit/f5f4838)), closes [#37](https://github.com/ncarlier/webhookd/issues/37)


### Features

* **notif:** add TLS and password support to the SMTP notifier ([f29a174](https://github.com/ncarlier/webhookd/commit/f29a174)), closes [#39](https://github.com/ncarlier/webhookd/issues/39)
* serve static assets ([b049677](https://github.com/ncarlier/webhookd/commit/b049677))



<a name="1.10.0"></a>
# [1.10.0](https://github.com/ncarlier/webhookd/compare/v1.9.0...v1.10.0) (2020-08-22)


### Bug Fixes

* **scripts:** fix line breaks ([#33](https://github.com/ncarlier/webhookd/issues/33)) ([552b683](https://github.com/ncarlier/webhookd/commit/552b683))


### Features

* **server:** simplify TLS usage ([1ee71be](https://github.com/ncarlier/webhookd/commit/1ee71be))
* improve HTTP signature support ([296ab6a](https://github.com/ncarlier/webhookd/commit/296ab6a))



<a name="1.9.0"></a>
# [1.9.0](https://github.com/ncarlier/webhookd/compare/1.8.0...v1.9.0) (2020-03-25)


### Bug Fixes

* **test:** add missing assets ([d96be9c](https://github.com/ncarlier/webhookd/commit/d96be9c))
* fix tracing id log ([a257e82](https://github.com/ncarlier/webhookd/commit/a257e82))
* **logger:** print colors only for TTY ([548149a](https://github.com/ncarlier/webhookd/commit/548149a))


### Features

* **api:** add info endpoint ([a5fe96d](https://github.com/ncarlier/webhookd/commit/a5fe96d))
* **docker:** use Docker Compose container wrapper ([8a393cc](https://github.com/ncarlier/webhookd/commit/8a393cc))
* **signature:** multi entries for a PEM file ([0e2f580](https://github.com/ncarlier/webhookd/commit/0e2f580))
* **signature:** refactore the trust store system ([d91e84d](https://github.com/ncarlier/webhookd/commit/d91e84d))



<a name="1.8.0"></a>
# [1.8.0](https://github.com/ncarlier/webhookd/compare/v1.7.0...1.8.0) (2020-02-26)


### Bug Fixes

* fix password default configuration ([4022007](https://github.com/ncarlier/webhookd/commit/4022007))
* improve configflag stability ([3044b39](https://github.com/ncarlier/webhookd/commit/3044b39))


### Features

* finalize HTTP signature support ([4320467](https://github.com/ncarlier/webhookd/commit/4320467))
* HTTP signature support ([c16ec83](https://github.com/ncarlier/webhookd/commit/c16ec83))
* **auth:** simplify validate method ([5948b60](https://github.com/ncarlier/webhookd/commit/5948b60))



<a name="1.7.0"></a>
# [1.7.0](https://github.com/ncarlier/webhookd/compare/v1.6.1...v1.7.0) (2020-02-04)


### Bug Fixes

* fix error message ([84524e9](https://github.com/ncarlier/webhookd/commit/84524e9))
* remove "done" statement at execution  end ([57e5b79](https://github.com/ncarlier/webhookd/commit/57e5b79))
* typo on HC route ([4c7b73b](https://github.com/ncarlier/webhookd/commit/4c7b73b))


### Features

* ACME support + configuration refactoring ([c7ea370](https://github.com/ncarlier/webhookd/commit/c7ea370))
* allow scripts with extensions ([2828873](https://github.com/ncarlier/webhookd/commit/2828873))
* improve logger ([e663336](https://github.com/ncarlier/webhookd/commit/e663336))
* logs refactoring ([d793c78](https://github.com/ncarlier/webhookd/commit/d793c78))
* refactoring of the config flag system ([6a01127](https://github.com/ncarlier/webhookd/commit/6a01127))



<a name="1.6.1"></a>
## [1.6.1](https://github.com/ncarlier/webhookd/compare/v1.6.0...v1.6.1) (2019-01-09)


### Bug Fixes

* catch SIGTERM signal for clean shutdown ([5c01d87](https://github.com/ncarlier/webhookd/commit/5c01d87))



<a name="1.6.0"></a>
# [1.6.0](https://github.com/ncarlier/webhookd/compare/v1.5.1...v1.6.0) (2019-01-07)


### Features

* **api:** add method whitelist ([d11da6f](https://github.com/ncarlier/webhookd/commit/d11da6f))
* **api:** add varz endpoint with metrics ([35a2321](https://github.com/ncarlier/webhookd/commit/35a2321))
* **api:** use GET and POST requests for hooks ([e7fac82](https://github.com/ncarlier/webhookd/commit/e7fac82))
* safer script resolution ([682b265](https://github.com/ncarlier/webhookd/commit/682b265))



<a name="1.5.1"></a>
## [1.5.1](https://github.com/ncarlier/webhookd/compare/v1.5.0...v1.5.1) (2019-01-06)


### Bug Fixes

* **api:** fix nil pointer ([4d2c75e](https://github.com/ncarlier/webhookd/commit/4d2c75e))


### Features

* **api:** add basic CORS support ([7a6af73](https://github.com/ncarlier/webhookd/commit/7a6af73))
* **cli:** improve parameters desc ([d8d440a](https://github.com/ncarlier/webhookd/commit/d8d440a))



<a name="1.5.0"></a>
# [1.5.0](https://github.com/ncarlier/webhookd/compare/v1.4.0...v1.5.0) (2018-12-31)


### Features

* **api:** add API endpoint to retrieve logs ([2ca5d67](https://github.com/ncarlier/webhookd/commit/2ca5d67))
* **notification:** complete refactoring of the notification system ([1dab1e9](https://github.com/ncarlier/webhookd/commit/1dab1e9))
* **worker:** add worker status lifecycle ([adead6d](https://github.com/ncarlier/webhookd/commit/adead6d))



<a name="1.4.0"></a>
# [1.4.0](https://github.com/ncarlier/webhookd/compare/v1.3.2...v1.4.0) (2018-12-18)


### Features

* use htpasswd to manage basic auth ([aab844c](https://github.com/ncarlier/webhookd/commit/aab844c))



<a name="1.3.2"></a>
## [1.3.2](https://github.com/ncarlier/webhookd/compare/v1.3.1...v1.3.2) (2018-12-13)


### Bug Fixes

* **runner:** fix concurrency and log file creation ([c5e393e](https://github.com/ncarlier/webhookd/commit/c5e393e))



<a name="1.3.1"></a>
## [1.3.1](https://github.com/ncarlier/webhookd/compare/v1.3.0...v1.3.1) (2018-11-17)


### Bug Fixes

* **worker:** use snakecase for log filename ([3021c19](https://github.com/ncarlier/webhookd/commit/3021c19)), closes [#8](https://github.com/ncarlier/webhookd/issues/8)



<a name="1.3.0"></a>
# [1.3.0](https://github.com/ncarlier/webhookd/compare/v1.2.1...v1.3.0) (2018-09-04)


### Bug Fixes

* **docker:** add bash shell to the Docker image ([c723387](https://github.com/ncarlier/webhookd/commit/c723387))
* **runner:** fix concurent access to the work request channel ([e8d1c6e](https://github.com/ncarlier/webhookd/commit/e8d1c6e))
* **script:** kill script process and sub process on timeout ([5f32a4f](https://github.com/ncarlier/webhookd/commit/5f32a4f))
* fix panic due to writing into closed chan ([43820cd](https://github.com/ncarlier/webhookd/commit/43820cd))
* **server:** remove global server timeouts ([82346b0](https://github.com/ncarlier/webhookd/commit/82346b0))


### Features

* **cli:** add print version command ([6565f6f](https://github.com/ncarlier/webhookd/commit/6565f6f))
* **config:** improve configuration flags ([fbf8794](https://github.com/ncarlier/webhookd/commit/fbf8794))
* **logging:** improve log outputs ([5cd5547](https://github.com/ncarlier/webhookd/commit/5cd5547))
* **security:** add http basic auth (fix [#6](https://github.com/ncarlier/webhookd/issues/6)) ([#7](https://github.com/ncarlier/webhookd/issues/7)) ([513e6d7](https://github.com/ncarlier/webhookd/commit/513e6d7))
* add Docker entrypoint ([67bfe07](https://github.com/ncarlier/webhookd/commit/67bfe07))



<a name="1.2.1"></a>
## [1.2.1](https://github.com/ncarlier/webhookd/compare/v1.2.0...v1.2.1) (2018-01-10)



<a name="1.2.0"></a>
# [1.2.0](https://github.com/ncarlier/webhookd/compare/v1.1.0...v1.2.0) (2018-01-09)


### Features

* add custom logger ([f500911](https://github.com/ncarlier/webhookd/commit/f500911))
* add webhook timeout ([7154828](https://github.com/ncarlier/webhookd/commit/7154828))



<a name="1.1.0"></a>
# [1.1.0](https://github.com/ncarlier/webhookd/compare/v1.0.1...v1.1.0) (2018-01-05)


### Features

* transmit HTTP headers as env variables to the script ([2e80359](https://github.com/ncarlier/webhookd/commit/2e80359))
* **docker:** add git and ssh client inside the Docker image ([10b82d6](https://github.com/ncarlier/webhookd/commit/10b82d6))
* **docker:** add jq inside the Docker image ([c2a1741](https://github.com/ncarlier/webhookd/commit/c2a1741))



<a name="1.0.1"></a>
## [1.0.1](https://github.com/ncarlier/webhookd/compare/v1.0.0...v1.0.1) (2018-01-05)


### Bug Fixes

* **api:** set SSE headers correctly ([d7b65e6](https://github.com/ncarlier/webhookd/commit/d7b65e6))



<a name="1.0.0"></a>
# [1.0.0](https://github.com/ncarlier/webhookd/compare/v0.0.3...v1.0.0) (2018-01-02)


### Bug Fixes

* **ci:** keep binary for release phase ([e6b0206](https://github.com/ncarlier/webhookd/commit/e6b0206))



<a name="0.0.3"></a>
## [0.0.3](https://github.com/ncarlier/webhookd/compare/v0.0.2...v0.0.3) (2015-04-20)


### Features

* Redirect script output in the console. ([d1fbdb1](https://github.com/ncarlier/webhookd/commit/d1fbdb1))



<a name="0.0.2"></a>
## [0.0.2](https://github.com/ncarlier/webhookd/compare/v0.0.1...v0.0.2) (2015-04-07)


### Bug Fixes

* Fix attachment in http notifier. ([deba0ef](https://github.com/ncarlier/webhookd/commit/deba0ef))
* Fix method definitions. ([bc8c93c](https://github.com/ncarlier/webhookd/commit/bc8c93c))
* Fix post request without attachment. ([80b82e4](https://github.com/ncarlier/webhookd/commit/80b82e4))
* Fix typo. ([79def29](https://github.com/ncarlier/webhookd/commit/79def29))
* Improve error logs. ([6149bca](https://github.com/ncarlier/webhookd/commit/6149bca))
* Merge stdout and stderr. ([539b674](https://github.com/ncarlier/webhookd/commit/539b674))
* remove specific scripts ([56ae93a](https://github.com/ncarlier/webhookd/commit/56ae93a))
* Some correction ([185f423](https://github.com/ncarlier/webhookd/commit/185f423))


### Features

* Add basic auth to http notifier. ([759126c](https://github.com/ncarlier/webhookd/commit/759126c))
* Add Gitlab hook. ([66be4ef](https://github.com/ncarlier/webhookd/commit/66be4ef))
* Add Gitlab hook. ([baf50c9](https://github.com/ncarlier/webhookd/commit/baf50c9))
* Add unit file. ([8978ddf](https://github.com/ncarlier/webhookd/commit/8978ddf))
* Add unit tests for the API. ([31c14f1](https://github.com/ncarlier/webhookd/commit/31c14f1))
* Create worker queue. ([b71c506](https://github.com/ncarlier/webhookd/commit/b71c506))
* Move JSON decoder inside each hook. ([0906ae1](https://github.com/ncarlier/webhookd/commit/0906ae1))



<a name="0.0.1"></a>
## [0.0.1](https://github.com/ncarlier/webhookd/compare/dbdd9f5...v0.0.1) (2014-09-23)


### Features

* Add build script for bitbucket hub. ([057ddd1](https://github.com/ncarlier/webhookd/commit/057ddd1))
* Add env configuration sample. ([9f1f9d0](https://github.com/ncarlier/webhookd/commit/9f1f9d0))
* Add minimal test script. ([dbdd9f5](https://github.com/ncarlier/webhookd/commit/dbdd9f5))
* Add notification system and docker hook support. ([474610e](https://github.com/ncarlier/webhookd/commit/474610e))
* Big refactoring. ([eb4b9ba](https://github.com/ncarlier/webhookd/commit/eb4b9ba))
* Create docker container. ([cb73848](https://github.com/ncarlier/webhookd/commit/cb73848))



