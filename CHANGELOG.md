# Changelog

## [0.14.0](https://github.com/opentdf/otdfctl/compare/v0.13.0...v0.14.0) (2024-10-01)


### Features

* **ci:** add e2e tests for subject mappings, support for --force delete ([#388](https://github.com/opentdf/otdfctl/issues/388)) ([c1f544b](https://github.com/opentdf/otdfctl/commit/c1f544b1079f52bfccb96c4c9e0b579a6854ad58))
* **ci:** add tests for subject condition sets, and --force delete flag ([#389](https://github.com/opentdf/otdfctl/issues/389)) ([c6d2abc](https://github.com/opentdf/otdfctl/commit/c6d2abcd4afe78d92fd285e5c77fecdfe806ed5d)), closes [#331](https://github.com/opentdf/otdfctl/issues/331)
* **ci:** e2e attribute definitions tests ([#384](https://github.com/opentdf/otdfctl/issues/384)) ([2894391](https://github.com/opentdf/otdfctl/commit/28943915f19e0fb565cfb38cfebdd6fde21c019a)), closes [#327](https://github.com/opentdf/otdfctl/issues/327)
* **core:** export manual functions for CLI wrappers to consume ([#397](https://github.com/opentdf/otdfctl/issues/397)) ([aa0bf95](https://github.com/opentdf/otdfctl/commit/aa0bf95a39dfc0aec4155e498a2096cbd158efdd))
* **core:** resource mappings LIST fix, delete --force support, and e2e tests ([#387](https://github.com/opentdf/otdfctl/issues/387)) ([326e74b](https://github.com/opentdf/otdfctl/commit/326e74b37d0abfb4ad50deadaa1ed46ecf9f8a5d)), closes [#386](https://github.com/opentdf/otdfctl/issues/386)


### Bug Fixes

* **core:** remove duplicate titling of help manual ([#391](https://github.com/opentdf/otdfctl/issues/391)) ([cb8db69](https://github.com/opentdf/otdfctl/commit/cb8db69ec4df42c7f230fbd87142bfbcd2d3940f))

## [0.13.0](https://github.com/opentdf/otdfctl/compare/v0.12.2...v0.13.0) (2024-09-12)


### Features

* add cli test mode and profile tests ([#313](https://github.com/opentdf/otdfctl/issues/313)) ([e0bc183](https://github.com/opentdf/otdfctl/commit/e0bc1836e8b5f14c87b5d572ad7937924c76d860))
* **ci:** make e2e test workflow reusable ([#365](https://github.com/opentdf/otdfctl/issues/365)) ([d94408c](https://github.com/opentdf/otdfctl/commit/d94408cc2898d46b3444e874c035ff2bffe451f4))
* **ci:** namespaces e2e tests and test suite improvements ([#351](https://github.com/opentdf/otdfctl/issues/351)) ([ce28555](https://github.com/opentdf/otdfctl/commit/ce285554866bf89ee8aa2df4a4b426548a58b59a))
* **ci:** reusable platform composite action in e2e tests ([#369](https://github.com/opentdf/otdfctl/issues/369)) ([f7d5a1c](https://github.com/opentdf/otdfctl/commit/f7d5a1c07304bee14dfc92fa81bd65389e76d9f6))
* **core:** add ecdsa-binding encrypt flag ([#360](https://github.com/opentdf/otdfctl/issues/360)) ([8702ec0](https://github.com/opentdf/otdfctl/commit/8702ec007b6d1354b6c0366e6b375f26216dfde1))
* **core:** adds missing long manual output docs ([#362](https://github.com/opentdf/otdfctl/issues/362)) ([8e1390f](https://github.com/opentdf/otdfctl/commit/8e1390f20c17a5900c586f94384af76ffd9a2844)), closes [#359](https://github.com/opentdf/otdfctl/issues/359)
* **core:** kas-grants list ([#346](https://github.com/opentdf/otdfctl/issues/346)) ([7f51282](https://github.com/opentdf/otdfctl/commit/7f512825eab814e3c130e3fe4e8ed85ecbe2d146)), closes [#253](https://github.com/opentdf/otdfctl/issues/253)


### Bug Fixes

* **ci:** e2e workflow should be fully reusable ([#368](https://github.com/opentdf/otdfctl/issues/368)) ([cc1e2b9](https://github.com/opentdf/otdfctl/commit/cc1e2b938fb0c8c4cf64d735f2961f7c9cae79fa))
* **ci:** enhance lint config and resolve all lint issues ([#363](https://github.com/opentdf/otdfctl/issues/363)) ([5c1dbf1](https://github.com/opentdf/otdfctl/commit/5c1dbf1f5e441ca0ebd8cfcca145a77b623f3638))
* **core:** GOOS, error message fixes ([#378](https://github.com/opentdf/otdfctl/issues/378)) ([623a82a](https://github.com/opentdf/otdfctl/commit/623a82ad3c1ed698a83eed54cf15a4f552096728)), closes [#380](https://github.com/opentdf/otdfctl/issues/380)
* **core:** metadata rendering cleanup ([#293](https://github.com/opentdf/otdfctl/issues/293)) ([ed21f81](https://github.com/opentdf/otdfctl/commit/ed21f81863450fd6167106711392e713a43c55be))
* **core:** wire attribute value FQNs to encrypt ([#370](https://github.com/opentdf/otdfctl/issues/370)) ([21f9b80](https://github.com/opentdf/otdfctl/commit/21f9b80cdee7d695a308937b08dbc768d11fbbd5))
* refactor to support varying print output ([#350](https://github.com/opentdf/otdfctl/issues/350)) ([d6932f3](https://github.com/opentdf/otdfctl/commit/d6932f30d9f653e46b32761a3257f3555ef0a6eb))

## [0.12.2](https://github.com/opentdf/otdfctl/compare/v0.12.1...v0.12.2) (2024-08-27)


### Bug Fixes

* **core:** improve KASR docs and add spellcheck GHA to pipeline ([#323](https://github.com/opentdf/otdfctl/issues/323)) ([a77cf30](https://github.com/opentdf/otdfctl/commit/a77cf30dc8077d034cb4c9df8cc94712b1a17dff)), closes [#335](https://github.com/opentdf/otdfctl/issues/335) [#337](https://github.com/opentdf/otdfctl/issues/337)
* create new http client to ignore tls verification ([#324](https://github.com/opentdf/otdfctl/issues/324)) ([4d4afb7](https://github.com/opentdf/otdfctl/commit/4d4afb7e5b6411bb08a92bc53181ac5730ca1992))

## [0.12.1](https://github.com/opentdf/otdfctl/compare/v0.12.0...v0.12.1) (2024-08-26)


### Bug Fixes

* **core:** remove documentation that cached kas pubkey is base64 ([#320](https://github.com/opentdf/otdfctl/issues/320)) ([fce8f44](https://github.com/opentdf/otdfctl/commit/fce8f44f767f35ccc4863f88d46e7ffcbd80f37a)), closes [#321](https://github.com/opentdf/otdfctl/issues/321)

## [0.12.0](https://github.com/opentdf/otdfctl/compare/v0.11.4...v0.12.0) (2024-08-23)


### Features

* **ci:** attr e2e tests with mixed casing ([#315](https://github.com/opentdf/otdfctl/issues/315)) ([50ce712](https://github.com/opentdf/otdfctl/commit/50ce712eab38f6686611e2b306bda5cacd55c28e))
* **core:** kasr cached keys to deprecate local ([#318](https://github.com/opentdf/otdfctl/issues/318)) ([5419cc3](https://github.com/opentdf/otdfctl/commit/5419cc39e143eb484f836ca1ee671d626d5e2c60)), closes [#317](https://github.com/opentdf/otdfctl/issues/317)

## [0.11.4](https://github.com/opentdf/otdfctl/compare/v0.11.3...v0.11.4) (2024-08-22)


### Bug Fixes

* update workflow permissions ([#310](https://github.com/opentdf/otdfctl/issues/310)) ([3979fe8](https://github.com/opentdf/otdfctl/commit/3979fe85c9ab6511376d98b672cbfebddbf9bb84))

## [0.11.3](https://github.com/opentdf/otdfctl/compare/v0.11.2...v0.11.3) (2024-08-22)


### Bug Fixes

* **core:** do not import unused fmt ([#306](https://github.com/opentdf/otdfctl/issues/306)) ([0dc552d](https://github.com/opentdf/otdfctl/commit/0dc552d3d6814f910c04d5f8cefa35404b4945f5))
* **core:** nil panic on set-default ([#304](https://github.com/opentdf/otdfctl/issues/304)) ([92bbfa3](https://github.com/opentdf/otdfctl/commit/92bbfa32ae42b73b68551c2f9d3551d357bc5922))
* **core:** warn and do now allow deletion of default profile ([#308](https://github.com/opentdf/otdfctl/issues/308)) ([fdd8167](https://github.com/opentdf/otdfctl/commit/fdd8167e8e2b22d652b48d796a756f86398bfd3c))
* make file not building correctly ([#307](https://github.com/opentdf/otdfctl/issues/307)) ([64eb821](https://github.com/opentdf/otdfctl/commit/64eb82170fdcc50396194271be358bf9c9d43049))

## [0.11.2](https://github.com/opentdf/otdfctl/compare/v0.11.1...v0.11.2) (2024-08-22)


### Bug Fixes

* disable tagging ([#302](https://github.com/opentdf/otdfctl/issues/302)) ([2b5db85](https://github.com/opentdf/otdfctl/commit/2b5db852ed0088e61f1180500135cd1865f9798b))

## [0.11.1](https://github.com/opentdf/otdfctl/compare/v0.11.0...v0.11.1) (2024-08-22)


### Bug Fixes

* release-please tweak ([#300](https://github.com/opentdf/otdfctl/issues/300)) ([29fc836](https://github.com/opentdf/otdfctl/commit/29fc8360ae0b701aefe70b25d1838f442fd7eb8d))

## [0.11.0](https://github.com/opentdf/otdfctl/compare/v0.10.0...v0.11.0) (2024-08-22)


### Features

* move git checkout before tagging ([#298](https://github.com/opentdf/otdfctl/issues/298)) ([1114e25](https://github.com/opentdf/otdfctl/commit/1114e25a90946e85622c8ff7a7befbf18beb4ba1))

## [0.10.0](https://github.com/opentdf/otdfctl/compare/v0.9.4...v0.10.0) (2024-08-22)


### Features

* add profile support for cli ([#289](https://github.com/opentdf/otdfctl/issues/289)) ([15700f3](https://github.com/opentdf/otdfctl/commit/15700f3375196595e4a0ea3a7a6dea4da06d8612))
* **core:** add scaffolding and POC for auth code flow ([#144](https://github.com/opentdf/otdfctl/issues/144)) ([03ecbfb](https://github.com/opentdf/otdfctl/commit/03ecbfb4f689f4a9f161a5a03d80efd50f728780))
* **core:** support kas grants to namespaces ([#292](https://github.com/opentdf/otdfctl/issues/292)) ([f2c6689](https://github.com/opentdf/otdfctl/commit/f2c6689d2f775b1aed907d553c42d87c8464e6c7)), closes [#269](https://github.com/opentdf/otdfctl/issues/269)
* improve auth with client credentials ([#286](https://github.com/opentdf/otdfctl/issues/286)) ([9c4968f](https://github.com/opentdf/otdfctl/commit/9c4968f48d1ba23a61ed5c8ad23a109bf141ba56))
* improve auth with client credentials ([#296](https://github.com/opentdf/otdfctl/issues/296)) ([0f533c7](https://github.com/opentdf/otdfctl/commit/0f533c7278a53ddd90656b3c7efcaee1c5bfd957))


### Bug Fixes

* **core:** bump platform deps ([#276](https://github.com/opentdf/otdfctl/issues/276)) ([e4ced99](https://github.com/opentdf/otdfctl/commit/e4ced996ae336b9db6db88906683f6600a2e5bf4))
* reduce prints ([#277](https://github.com/opentdf/otdfctl/issues/277)) ([8b5734a](https://github.com/opentdf/otdfctl/commit/8b5734a18636071566fd8c4cfc808f3f240a02a5))
