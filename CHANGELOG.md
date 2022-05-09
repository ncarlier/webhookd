# [](https://github.com/ncarlier/webhookd/compare/v1.14.0...v) (2022-05-09)



# [1.14.0](https://github.com/ncarlier/webhookd/compare/v1.13.0...v1.14.0) (2022-05-09)


### Features

* improve control on streaming protocol ([4e6298d](https://github.com/ncarlier/webhookd/commit/4e6298dda713d816855818282cc32caca108c597)), closes [#51](https://github.com/ncarlier/webhookd/issues/51)
* support application/x-www-form-urlencoded ([173ba6c](https://github.com/ncarlier/webhookd/commit/173ba6c347e3a4a158fa976cb8a3a7913aa4cede))



# [1.13.0](https://github.com/ncarlier/webhookd/compare/v1.12.0...v1.13.0) (2021-12-05)


### Features

* **api:** refactore router ([6b3623f](https://github.com/ncarlier/webhookd/commit/6b3623f67a402649b22490475e93e8567559cdf9))
* configure static path ([9fa96ac](https://github.com/ncarlier/webhookd/commit/9fa96acfb264185f773e677858813eea1dcff552)), closes [#45](https://github.com/ncarlier/webhookd/issues/45)



# [1.12.0](https://github.com/ncarlier/webhookd/compare/v1.11.0...v1.12.0) (2021-05-29)


### Features

* **notification:** email subject customization ([ed67fc7](https://github.com/ncarlier/webhookd/commit/ed67fc72f6a9512ca040edfce3f4119ac113eb5a))
* output hook execution logs to server logs ([07fbb6e](https://github.com/ncarlier/webhookd/commit/07fbb6ee3a382d85cbb8a8243c01f751496ec85f)), closes [#44](https://github.com/ncarlier/webhookd/issues/44)



# [1.11.0](https://github.com/ncarlier/webhookd/compare/v1.10.0...v1.11.0) (2020-10-14)


### Bug Fixes

* fix timeout parameter name ([f5f4838](https://github.com/ncarlier/webhookd/commit/f5f48381aff6d80cd2da10f7d2fa76a8e5683704)), closes [#37](https://github.com/ncarlier/webhookd/issues/37)


### Features

* **notif:** add TLS and password support to the SMTP notifier ([f29a174](https://github.com/ncarlier/webhookd/commit/f29a1748ef5afde7de2f60a89d2883d337008853)), closes [#39](https://github.com/ncarlier/webhookd/issues/39)
* serve static assets ([b049677](https://github.com/ncarlier/webhookd/commit/b0496778e82016e5d2aab100a2f58784d106e230))



# [1.10.0](https://github.com/ncarlier/webhookd/compare/v1.9.0...v1.10.0) (2020-08-22)


### Bug Fixes

* **scripts:** fix line breaks ([#33](https://github.com/ncarlier/webhookd/issues/33)) ([552b683](https://github.com/ncarlier/webhookd/commit/552b683f9ec7f6f0e340d80d74269bc038af0fc7))


### Features

* improve HTTP signature support ([296ab6a](https://github.com/ncarlier/webhookd/commit/296ab6aaa31b747fcf33d4647e2cb0f4d8c585d4))
* **server:** simplify TLS usage ([1ee71be](https://github.com/ncarlier/webhookd/commit/1ee71be4c534daa01dce1c23b40b3b4fb3e65881))



# [1.9.0](https://github.com/ncarlier/webhookd/compare/1.8.0...v1.9.0) (2020-03-25)


### Bug Fixes

* fix tracing id log ([a257e82](https://github.com/ncarlier/webhookd/commit/a257e82274255e9574a58d9f95d4b297f78ae9c9))
* **logger:** print colors only for TTY ([548149a](https://github.com/ncarlier/webhookd/commit/548149a63d200076f0bf93e267dff797b44c4658))
* **test:** add missing assets ([d96be9c](https://github.com/ncarlier/webhookd/commit/d96be9cd8d59003f85a702e9a046aae7f9d3d062))


### Features

* **api:** add info endpoint ([a5fe96d](https://github.com/ncarlier/webhookd/commit/a5fe96d2e8772b05cf45fa18bca6d13e22390929))
* **docker:** use Docker Compose container wrapper ([8a393cc](https://github.com/ncarlier/webhookd/commit/8a393cc0c32340abfbaf2c2693e38764ac41a0e3))
* **signature:** multi entries for a PEM file ([0e2f580](https://github.com/ncarlier/webhookd/commit/0e2f58012d9ae40dc85a26dc3edb9eabea49bd71))
* **signature:** refactore the trust store system ([d91e84d](https://github.com/ncarlier/webhookd/commit/d91e84d1be3f0a15c590ebdc51f6e2694cb6f645))



# [1.8.0](https://github.com/ncarlier/webhookd/compare/v1.7.0...1.8.0) (2020-02-26)


### Bug Fixes

* fix password default configuration ([4022007](https://github.com/ncarlier/webhookd/commit/40220077d292cc760218f4288e3addac1cc6d9ce))
* improve configflag stability ([3044b39](https://github.com/ncarlier/webhookd/commit/3044b3951b7cb18a425310ddfe24cb44f5939f15))


### Features

* **auth:** simplify validate method ([5948b60](https://github.com/ncarlier/webhookd/commit/5948b6001fb15276c3125efab3bdf98b8e5d0fe1))
* finalize HTTP signature support ([4320467](https://github.com/ncarlier/webhookd/commit/43204677d34e11d145ae26012a2702c9ec4ddfb5))
* HTTP signature support ([c16ec83](https://github.com/ncarlier/webhookd/commit/c16ec83a5ae0921e8b06506650d839109ef69fdc))



# [1.7.0](https://github.com/ncarlier/webhookd/compare/v1.6.1...v1.7.0) (2020-02-04)


### Bug Fixes

* fix error message ([84524e9](https://github.com/ncarlier/webhookd/commit/84524e91752fbe4573320f26891e097ad65d4b0e))
* remove "done" statement at execution  end ([57e5b79](https://github.com/ncarlier/webhookd/commit/57e5b79011dcd3fd88f5b3c889551ea895a96e5c))
* typo on HC route ([4c7b73b](https://github.com/ncarlier/webhookd/commit/4c7b73b987c7381bde8387722aaa69081605156b))


### Features

* ACME support + configuration refactoring ([c7ea370](https://github.com/ncarlier/webhookd/commit/c7ea370de124c9f5eba9e18db0234be43b8efc64))
* allow scripts with extensions ([2828873](https://github.com/ncarlier/webhookd/commit/28288739f35bbfc0253bdcb9d7d106698d0d3e92))
* improve logger ([e663336](https://github.com/ncarlier/webhookd/commit/e663336ecb3f88dd7becd6d62b48f0628abd6b42))
* logs refactoring ([d793c78](https://github.com/ncarlier/webhookd/commit/d793c7813d254db8218b4f4bed22beb003d8db24))
* refactoring of the config flag system ([6a01127](https://github.com/ncarlier/webhookd/commit/6a011272fdd10354948fab670cead2c1deced317))



## [1.6.1](https://github.com/ncarlier/webhookd/compare/v1.6.0...v1.6.1) (2019-01-09)


### Bug Fixes

* catch SIGTERM signal for clean shutdown ([5c01d87](https://github.com/ncarlier/webhookd/commit/5c01d87aa3f8a6fb3843691cac5b8bb9871679ca))



# [1.6.0](https://github.com/ncarlier/webhookd/compare/v1.5.1...v1.6.0) (2019-01-07)


### Features

* **api:** add method whitelist ([d11da6f](https://github.com/ncarlier/webhookd/commit/d11da6fa54cb5477366f8711fc88de498df1408a))
* **api:** add varz endpoint with metrics ([35a2321](https://github.com/ncarlier/webhookd/commit/35a2321f808c75eb19fddffa39ee5e6dd2259d32))
* **api:** use GET and POST requests for hooks ([e7fac82](https://github.com/ncarlier/webhookd/commit/e7fac829aa20d93b8724504000dc80d890c3f22c))
* safer script resolution ([682b265](https://github.com/ncarlier/webhookd/commit/682b265d3e4a878a2bea7662d33e039f86e9731e))



## [1.5.1](https://github.com/ncarlier/webhookd/compare/v1.5.0...v1.5.1) (2019-01-06)


### Bug Fixes

* **api:** fix nil pointer ([4d2c75e](https://github.com/ncarlier/webhookd/commit/4d2c75e70b3ce780414fb884fcd102da88eef2c7))


### Features

* **api:** add basic CORS support ([7a6af73](https://github.com/ncarlier/webhookd/commit/7a6af7312ac545a74d18288408916f57c71cb0c9))
* **cli:** improve parameters desc ([d8d440a](https://github.com/ncarlier/webhookd/commit/d8d440a6b10042a5a7dd7bc796434ac1d8994ffc))



# [1.5.0](https://github.com/ncarlier/webhookd/compare/v1.4.0...v1.5.0) (2018-12-31)


### Features

* **api:** add API endpoint to retrieve logs ([2ca5d67](https://github.com/ncarlier/webhookd/commit/2ca5d671b9264e7e52d2db576ae039e3991c813a))
* **notification:** complete refactoring of the notification system ([1dab1e9](https://github.com/ncarlier/webhookd/commit/1dab1e968d5069b5cebb38c187e0cdbe13bc84ce))
* **worker:** add worker status lifecycle ([adead6d](https://github.com/ncarlier/webhookd/commit/adead6d3b379c05a2ee83aa69b6f825a0b2b77bd))



# [1.4.0](https://github.com/ncarlier/webhookd/compare/v1.3.2...v1.4.0) (2018-12-18)


### Features

* use htpasswd to manage basic auth ([aab844c](https://github.com/ncarlier/webhookd/commit/aab844cee7086cd0492957061c6234935bf908e1))



## [1.3.2](https://github.com/ncarlier/webhookd/compare/v1.3.1...v1.3.2) (2018-12-13)


### Bug Fixes

* **runner:** fix concurrency and log file creation ([c5e393e](https://github.com/ncarlier/webhookd/commit/c5e393eb928f8805bb4a04129eb2d5c2ab8bc3d4))



## [1.3.1](https://github.com/ncarlier/webhookd/compare/v1.3.0...v1.3.1) (2018-11-17)


### Bug Fixes

* **worker:** use snakecase for log filename ([3021c19](https://github.com/ncarlier/webhookd/commit/3021c19551975286be43170f79cfcf90b57b3b2b)), closes [#8](https://github.com/ncarlier/webhookd/issues/8)



# [1.3.0](https://github.com/ncarlier/webhookd/compare/v1.2.1...v1.3.0) (2018-09-04)


### Bug Fixes

* **docker:** add bash shell to the Docker image ([c723387](https://github.com/ncarlier/webhookd/commit/c723387d8b9dd5bdd6e63bf243a079ac3c188af7))
* fix panic due to writing into closed chan ([43820cd](https://github.com/ncarlier/webhookd/commit/43820cd9f0331e4212769876d186348bd5153c43))
* **runner:** fix concurent access to the work request channel ([e8d1c6e](https://github.com/ncarlier/webhookd/commit/e8d1c6e5816dbdeca4b6704e0ae7b02616b4ab5f))
* **script:** kill script process and sub process on timeout ([5f32a4f](https://github.com/ncarlier/webhookd/commit/5f32a4f7f807959927f282083d0e9d83ccde15b6))
* **server:** remove global server timeouts ([82346b0](https://github.com/ncarlier/webhookd/commit/82346b08da34671ce5db6afb8a55d866358677b2))


### Features

* add Docker entrypoint ([67bfe07](https://github.com/ncarlier/webhookd/commit/67bfe0786ffec1452b690c4f31fd061b74754367))
* **cli:** add print version command ([6565f6f](https://github.com/ncarlier/webhookd/commit/6565f6f6ba6b0ef13a9280e5cfbac116d6e06290))
* **config:** improve configuration flags ([fbf8794](https://github.com/ncarlier/webhookd/commit/fbf8794d0a891c7a934f0d08f9ee5af591329cd1))
* **logging:** improve log outputs ([5cd5547](https://github.com/ncarlier/webhookd/commit/5cd5547aa54a7cb57f329b0723af7035ddcac48f))
* **security:** add http basic auth (fix [#6](https://github.com/ncarlier/webhookd/issues/6)) ([#7](https://github.com/ncarlier/webhookd/issues/7)) ([513e6d7](https://github.com/ncarlier/webhookd/commit/513e6d78dd33d84fc313fdca07b06581595a519d))



## [1.2.1](https://github.com/ncarlier/webhookd/compare/v1.2.0...v1.2.1) (2018-01-10)



# [1.2.0](https://github.com/ncarlier/webhookd/compare/v1.1.0...v1.2.0) (2018-01-09)


### Features

* add custom logger ([f500911](https://github.com/ncarlier/webhookd/commit/f50091131b84c896016bbef3f62b00aaf93f3edd))
* add webhook timeout ([7154828](https://github.com/ncarlier/webhookd/commit/7154828ecab83105f1f7ee78b124ef9ac76b2548))



# [1.1.0](https://github.com/ncarlier/webhookd/compare/v1.0.1...v1.1.0) (2018-01-05)


### Features

* **docker:** add git and ssh client inside the Docker image ([10b82d6](https://github.com/ncarlier/webhookd/commit/10b82d67a1167cac69a1f175fbfcd9e08a53f602))
* **docker:** add jq inside the Docker image ([c2a1741](https://github.com/ncarlier/webhookd/commit/c2a17414f9dc35cdad0301f065d998bd7ef53481))
* transmit HTTP headers as env variables to the script ([2e80359](https://github.com/ncarlier/webhookd/commit/2e803598d2b70190cb8613e0911f1ac404056c33))



## [1.0.1](https://github.com/ncarlier/webhookd/compare/v1.0.0...v1.0.1) (2018-01-05)


### Bug Fixes

* **api:** set SSE headers correctly ([d7b65e6](https://github.com/ncarlier/webhookd/commit/d7b65e68ae61897cf64666855227184941adbb1c))



# [1.0.0](https://github.com/ncarlier/webhookd/compare/v0.0.3...v1.0.0) (2018-01-02)


### Bug Fixes

* **ci:** keep binary for release phase ([e6b0206](https://github.com/ncarlier/webhookd/commit/e6b02069e708521a06c21baba67ba5d733102d54))



## [0.0.3](https://github.com/ncarlier/webhookd/compare/v0.0.2...v0.0.3) (2015-04-20)


### Features

* Redirect script output in the console. ([d1fbdb1](https://github.com/ncarlier/webhookd/commit/d1fbdb139bf2d8bad8ee684c021ccfc2842a0803))



## [0.0.2](https://github.com/ncarlier/webhookd/compare/v0.0.1...v0.0.2) (2015-04-07)


### Bug Fixes

* Fix attachment in http notifier. ([deba0ef](https://github.com/ncarlier/webhookd/commit/deba0ef4624279b7932cc37fcd36b1d808cdcdd2))
* Fix method definitions. ([bc8c93c](https://github.com/ncarlier/webhookd/commit/bc8c93c9904b260463699527672d7b0fad22c20b))
* Fix post request without attachment. ([80b82e4](https://github.com/ncarlier/webhookd/commit/80b82e4ec7cecfcd82672f2435691bd8fb58e267))
* Fix typo. ([79def29](https://github.com/ncarlier/webhookd/commit/79def293f0a7f5b43302df634575c48a51546b57))
* Improve error logs. ([6149bca](https://github.com/ncarlier/webhookd/commit/6149bca5ab4f07dbe9fd38f93995dfc9d1fff92e))
* Merge stdout and stderr. ([539b674](https://github.com/ncarlier/webhookd/commit/539b674cc774db643a1c70b3f2500e662dbd2002))
* remove specific scripts ([56ae93a](https://github.com/ncarlier/webhookd/commit/56ae93ac33bdf13b0a4695fc2b0a7e64f4f91344))
* Some correction ([185f423](https://github.com/ncarlier/webhookd/commit/185f4239d8a726f37968457844b4761619ef3e67))


### Features

* Add basic auth to http notifier. ([759126c](https://github.com/ncarlier/webhookd/commit/759126cbc6f68b27501ec505188f763b44b98d0f))
* Add Gitlab hook. ([66be4ef](https://github.com/ncarlier/webhookd/commit/66be4efba026a3fdd0c0de6ccb2b04c730af117e))
* Add Gitlab hook. ([baf50c9](https://github.com/ncarlier/webhookd/commit/baf50c9709aeaa27771250b1dc554f80612edaf9))
* Add unit file. ([8978ddf](https://github.com/ncarlier/webhookd/commit/8978ddf4f484caca2d7c975324488fb76c7ac993))
* Add unit tests for the API. ([31c14f1](https://github.com/ncarlier/webhookd/commit/31c14f17dc321103f0de17077bb664006db3ed75))
* Create worker queue. ([b71c506](https://github.com/ncarlier/webhookd/commit/b71c506587571dc3dcdbb9d5cc14311f5cd2ddd4))
* Move JSON decoder inside each hook. ([0906ae1](https://github.com/ncarlier/webhookd/commit/0906ae1217434e256ee09a42f2ba72c1f339e4df))



## [0.0.1](https://github.com/ncarlier/webhookd/compare/dbdd9f57767c0b236a2ac5c097689ad64c255ed9...v0.0.1) (2014-09-23)


### Features

* Add build script for bitbucket hub. ([057ddd1](https://github.com/ncarlier/webhookd/commit/057ddd1e5ecccd5c52b286ffa3067ed2c7e89ee0))
* Add env configuration sample. ([9f1f9d0](https://github.com/ncarlier/webhookd/commit/9f1f9d0f783f1b9b315ce8484b9de8e90c5e8d48))
* Add minimal test script. ([dbdd9f5](https://github.com/ncarlier/webhookd/commit/dbdd9f57767c0b236a2ac5c097689ad64c255ed9))
* Add notification system and docker hook support. ([474610e](https://github.com/ncarlier/webhookd/commit/474610e25b08a18defea41b013718aac44b6b5e5))
* Big refactoring. ([eb4b9ba](https://github.com/ncarlier/webhookd/commit/eb4b9ba7ffaf6617c687f01af1b3f99402b20fc5))
* Create docker container. ([cb73848](https://github.com/ncarlier/webhookd/commit/cb738486f3085a515beaf0a397d12eba0305f51c))



