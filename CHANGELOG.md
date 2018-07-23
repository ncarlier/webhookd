<a name=""></a>
# [](https://github.com/ncarlier/webhookd/compare/v1.2.4...v) (2018-07-23)



<a name="1.2.4"></a>
## [1.2.4](https://github.com/ncarlier/webhookd/compare/v1.2.3...v1.2.4) (2018-07-23)


### Bug Fixes

* **script:** kill script process and sub process on timeout ([92ef4e4](https://github.com/ncarlier/webhookd/commit/92ef4e4))


### Features

* add Docker entrypoint ([4b58677](https://github.com/ncarlier/webhookd/commit/4b58677))
* **config:** improve configuration flags ([eb699dd](https://github.com/ncarlier/webhookd/commit/eb699dd))



<a name="1.2.3"></a>
## [1.2.3](https://github.com/ncarlier/webhookd/compare/v1.2.2...v1.2.3) (2018-05-06)


### Bug Fixes

* **docker:** add bash shell to the Docker image ([d760073](https://github.com/ncarlier/webhookd/commit/d760073))



<a name="1.2.2"></a>
## [1.2.2](https://github.com/ncarlier/webhookd/compare/v1.2.1...v1.2.2) (2018-03-21)


### Bug Fixes

* **server:** remove global server timeouts ([1e38906](https://github.com/ncarlier/webhookd/commit/1e38906))
* fix panic due to writing into closed chan ([2a22680](https://github.com/ncarlier/webhookd/commit/2a22680))



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

* **docker:** add git and ssh client inside the Docker image ([10b82d6](https://github.com/ncarlier/webhookd/commit/10b82d6))
* **docker:** add jq inside the Docker image ([c2a1741](https://github.com/ncarlier/webhookd/commit/c2a1741))
* transmit HTTP headers as env variables to the script ([2e80359](https://github.com/ncarlier/webhookd/commit/2e80359))



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



