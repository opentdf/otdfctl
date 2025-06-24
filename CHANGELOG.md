# Changelog

## [0.22.0](https://github.com/opentdf/otdfctl/compare/v0.21.0...v0.22.0) (2025-06-24)


### ⚠ BREAKING CHANGES

* remove the ability to assign grants ([#604](https://github.com/opentdf/otdfctl/issues/604))

### Features

* **core:** dynamic port allocation ([#606](https://github.com/opentdf/otdfctl/issues/606)) ([75552e1](https://github.com/opentdf/otdfctl/commit/75552e187eef204b03b1d13d55920fa43ec3cf30))
* **core:** Uncomment code and pull in new protos. ([#594](https://github.com/opentdf/otdfctl/issues/594)) ([2883e50](https://github.com/opentdf/otdfctl/commit/2883e5060ca1f9d22f9a9500293fc407e7f4bcfd))
* **core:** Unhide key commands. ([#607](https://github.com/opentdf/otdfctl/issues/607)) ([a3660d9](https://github.com/opentdf/otdfctl/commit/a3660d9e8271e3fd179e6521eab02a2b096a01db))
* remove the ability to assign grants ([#604](https://github.com/opentdf/otdfctl/issues/604)) ([c9f0d82](https://github.com/opentdf/otdfctl/commit/c9f0d822747a62a6253c441ede144238715da50b))


### Bug Fixes

* add more Deprecated text to kas-grants ([#605](https://github.com/opentdf/otdfctl/issues/605)) ([2106d2f](https://github.com/opentdf/otdfctl/commit/2106d2f5189de49fe05b94025228474ffdb026ae))
* **ci:** Trigger for release-please (testing) ([#580](https://github.com/opentdf/otdfctl/issues/580)) ([5cd33f9](https://github.com/opentdf/otdfctl/commit/5cd33f9f9b5fb66b2cc9c0c795bd84cf10630298))
* **core:** Change base key name so tests run last. ([#611](https://github.com/opentdf/otdfctl/issues/611)) ([464b179](https://github.com/opentdf/otdfctl/commit/464b179a3134890943d8319bdee41cbad9078d64))
* **core:** Move key management under policy. ([#597](https://github.com/opentdf/otdfctl/issues/597)) ([d657e96](https://github.com/opentdf/otdfctl/commit/d657e96cab3afc516437ae08321ab45aff376460))
* disable kas-registry --public-keys and --publickey-remote flags ([#603](https://github.com/opentdf/otdfctl/issues/603)) ([279bbbd](https://github.com/opentdf/otdfctl/commit/279bbbd8ced14765c97ae3928421d38737ac0a8d))
* enforce hex encoded wrapping-key ([#581](https://github.com/opentdf/otdfctl/issues/581)) ([416e215](https://github.com/opentdf/otdfctl/commit/416e215abf0c910aa4d18dc84729f89ea578fd4d))
* **main:** Use cmd.Context for resource mapping group commands ([#592](https://github.com/opentdf/otdfctl/issues/592)) ([b5d8b6f](https://github.com/opentdf/otdfctl/commit/b5d8b6f6c335483873cec90363d94e0196d18b14))

## [0.21.0](https://github.com/opentdf/otdfctl/compare/v0.20.0...v0.21.0) (2025-05-29)


### Features

* Add initial Dependency Review configuration ([#551](https://github.com/opentdf/otdfctl/issues/551)) ([b622666](https://github.com/opentdf/otdfctl/commit/b6226660c1d75e133a8ead456efcab74de4b4fc0))
* **core:** Add base key cmds ([#563](https://github.com/opentdf/otdfctl/issues/563)) ([edfd6c0](https://github.com/opentdf/otdfctl/commit/edfd6c08dc9b84f2cbfc79643ccc266a45ce58fd))
* **core:** DSPX-18 clean up Go context usage to follow best practices ([#558](https://github.com/opentdf/otdfctl/issues/558)) ([a2c9f8b](https://github.com/opentdf/otdfctl/commit/a2c9f8b13cbab740b46262f70aecc82a94f3d788))
* **core:** DSPX-608 - Deprecate public_client_id ([#555](https://github.com/opentdf/otdfctl/issues/555)) ([8d396bd](https://github.com/opentdf/otdfctl/commit/8d396bd022126524d9d20daa03ec6ca262cf4406))
* **core:** DSPX-608 - require clientID for login ([#553](https://github.com/opentdf/otdfctl/issues/553)) ([580172e](https://github.com/opentdf/otdfctl/commit/580172e1861b54366f4914a141e459fe3221a16d))
* **core:** DSPX-896 add registered resources CRUD ([#559](https://github.com/opentdf/otdfctl/issues/559)) ([8e7475e](https://github.com/opentdf/otdfctl/commit/8e7475ef8aab91d28ab7efd320af13dc5ab53d3b))
* **core:** KAS allowlist options ([#539](https://github.com/opentdf/otdfctl/issues/539)) ([af7978f](https://github.com/opentdf/otdfctl/commit/af7978f86ced38543b31b792e008654071333789))
* **core:** key management operations ([#533](https://github.com/opentdf/otdfctl/issues/533)) ([d4f6aaa](https://github.com/opentdf/otdfctl/commit/d4f6aaac3f6fc1b50fbc988e5d34a32de0ed9f64))
* **main:** add actions CRUD and e2e tests ([#523](https://github.com/opentdf/otdfctl/issues/523)) ([2fb9ec7](https://github.com/opentdf/otdfctl/commit/2fb9ec7336da5731b868da94f0bbd5b2f226ede1))
* **main:** refactor actions within existing CLI policy object CRUD ([#543](https://github.com/opentdf/otdfctl/issues/543)) ([9ab1a58](https://github.com/opentdf/otdfctl/commit/9ab1a58418643ea709aefb08e3f5ca8bd06235f4))
* **core:** Resource mapping groups ([#567](https://github.com/opentdf/otdfctl/issues/567)) ([03fa307](https://github.com/opentdf/otdfctl/commit/03fa307b3ab91f25baeb74e30fde6eeec6d479a1))
* **core:** Update key mgmt flags to consistent format ([#570](https://github.com/opentdf/otdfctl/issues/570)) ([#846f96c](https://github.com/opentdf/otdfctl/commit/846f96cb9adfe03e355c9e64b559f1c11d84a86f))
* **core:** Rotate Key ([#572](https://github.com/opentdf/otdfctl/issues/572)) ([afd0043](https://github.com/opentdf/otdfctl/commit/afd0043f1ea66f0b371a95b556320551f73749bb))


### Bug Fixes

* **ci:** ci job should run on changes to GHA ([#530](https://github.com/opentdf/otdfctl/issues/530)) ([1d296ca](https://github.com/opentdf/otdfctl/commit/1d296ca8fac889a6e776ad381df999a2fcf9d6ce))
* **main:** Pass the full url when building the sdk object ([#544](https://github.com/opentdf/otdfctl/issues/544)) ([8b836f0](https://github.com/opentdf/otdfctl/commit/8b836f0fa3aa414c3ab19d830f4d1f833d3ae61d))

## [0.20.0](https://github.com/opentdf/otdfctl/compare/v0.19.0...v0.20.0) (2025-04-08)


### Features

* **core:** add aliases for profile command ([#510](https://github.com/opentdf/otdfctl/issues/510)) ([45c633d](https://github.com/opentdf/otdfctl/commit/45c633da6b00b04a8c92686521d25144048ac62c))
* **core:** Add support for WithTargetMode encrypt option ([#519](https://github.com/opentdf/otdfctl/issues/519)) ([a0ab213](https://github.com/opentdf/otdfctl/commit/a0ab2136be0b1d39e16a7522210f493fd797089d))


### Bug Fixes

* **core:** bump jwt dep and remove outdated version ([#520](https://github.com/opentdf/otdfctl/issues/520)) ([77bb9ca](https://github.com/opentdf/otdfctl/commit/77bb9ca9a0741ab7b920cc00f264a021064b117c))

## [0.19.0](https://github.com/opentdf/otdfctl/compare/v0.18.0...v0.19.0) (2025-03-05)


### Features

* **core:** support for ec-wrapping ([#499](https://github.com/opentdf/otdfctl/issues/499)) ([e839445](https://github.com/opentdf/otdfctl/commit/e839445181c89447d9a2374d54ce5ea4c3f46320))


### Bug Fixes

* **core:** mark new algorithm flags experimental ([#501](https://github.com/opentdf/otdfctl/issues/501)) ([95e00bf](https://github.com/opentdf/otdfctl/commit/95e00bf3daa8eb05196a5839488a4718c2230210))

## [0.18.0](https://github.com/opentdf/otdfctl/compare/v0.17.1...v0.18.0) (2025-02-25)


### Features

* Assertion verification ([#452](https://github.com/opentdf/otdfctl/issues/452)) ([5a8fe0d](https://github.com/opentdf/otdfctl/commit/5a8fe0d64088b74c95d3376e4a2a5a47d680d9c0))
* **core:** Adding examples docs, mainly policy commands ([#461](https://github.com/opentdf/otdfctl/issues/461)) ([04c1743](https://github.com/opentdf/otdfctl/commit/04c17439bb5f68fb5d44ba96cb457ce9ca072250))
* **core:** bump SDK and consume new platform connection validation ([#493](https://github.com/opentdf/otdfctl/issues/493)) ([1106b54](https://github.com/opentdf/otdfctl/commit/1106b54e73f9ceb711ff19d15cd08bf1cebbb29f))
* **core:** Shows SDK version and spec info ([#474](https://github.com/opentdf/otdfctl/issues/474)) ([5a685c4](https://github.com/opentdf/otdfctl/commit/5a685c4e36cf524c4f594fac42cfec30f62a6e83))

## [0.17.1](https://github.com/opentdf/otdfctl/compare/v0.17.0...v0.17.1) (2024-12-09)


### Bug Fixes

* **core:** kasr creation JSON example ([#453](https://github.com/opentdf/otdfctl/issues/453)) ([192c7b2](https://github.com/opentdf/otdfctl/commit/192c7b2975a4ab6f648ab7924e20e70535ce04b2))

## [0.17.0](https://github.com/opentdf/otdfctl/compare/v0.16.0...v0.17.0) (2024-12-05)


### Features

* **core:** pagination of LIST commands ([#447](https://github.com/opentdf/otdfctl/issues/447)) ([673a064](https://github.com/opentdf/otdfctl/commit/673a06424d30e706798b9a1fa1bbfd9b4601e765))
* **core:** subject condition set prune ([#439](https://github.com/opentdf/otdfctl/issues/439)) ([c4c8b8b](https://github.com/opentdf/otdfctl/commit/c4c8b8b276b2189df74e6cf30e14abac9369d97e))


### Bug Fixes

* **core:** kas registry get should allow -i 'id' flag shorthand ([#434](https://github.com/opentdf/otdfctl/issues/434)) ([bed3701](https://github.com/opentdf/otdfctl/commit/bed3701d89510ee78c3aed43b1a072e41ee3873f))
* **core:** sm list should provide value fqn instead of just value string ([#438](https://github.com/opentdf/otdfctl/issues/438)) ([9a7cb72](https://github.com/opentdf/otdfctl/commit/9a7cb7242e0e39ccc2b54425028638fa0c5e3f9f))

## [0.16.0](https://github.com/opentdf/otdfctl/compare/v0.15.0...v0.16.0) (2024-11-20)


### Features

* assertion verification disable ([#419](https://github.com/opentdf/otdfctl/issues/419)) ([acf5702](https://github.com/opentdf/otdfctl/commit/acf57028f1481f432b6b0c3c7a3e2c2261ac739f))
* **core:** add `subject-mappings match` to CLI ([#413](https://github.com/opentdf/otdfctl/issues/413)) ([bc56c19](https://github.com/opentdf/otdfctl/commit/bc56c199a73b12b8c90045d1b6f9cc6fdec16c54))
* **core:** add optional name to kas registry CRUD commands ([#429](https://github.com/opentdf/otdfctl/issues/429)) ([f675d86](https://github.com/opentdf/otdfctl/commit/f675d86c83205232db407d6609e80fa865a3998e))
* **core:** adds assertions to encrypt subcommand ([#408](https://github.com/opentdf/otdfctl/issues/408)) ([8f0e906](https://github.com/opentdf/otdfctl/commit/8f0e906c1dfe99fe6aa5f2ff43d02f0da90474cf))
* **core:** adds storeFile to save encrypted profiles to disk and updates auth to propagate tlsNoVerify ([#420](https://github.com/opentdf/otdfctl/issues/420)) ([f709e01](https://github.com/opentdf/otdfctl/commit/f709e014bf3f82a2808eae5df76b3667730c36ef))
* refactor encrypt and decrypt + CLI examples ([#418](https://github.com/opentdf/otdfctl/issues/418)) ([e681823](https://github.com/opentdf/otdfctl/commit/e681823ad54ddf70f4aa2215438d69a3d02cf6eb))
* support --with-access-token for auth ([#409](https://github.com/opentdf/otdfctl/issues/409)) ([856efa4](https://github.com/opentdf/otdfctl/commit/856efa4d61bb24b05f3a98943b94600ff77536fa))


### Bug Fixes

* **core:** dev selectors employ flattening from platform instead of jq ([#411](https://github.com/opentdf/otdfctl/issues/411)) ([57966ff](https://github.com/opentdf/otdfctl/commit/57966ffadcc61e1611869171bd3fc85723492fb7))
* **core:** improve readability of TDF methods ([#424](https://github.com/opentdf/otdfctl/issues/424)) ([a88d386](https://github.com/opentdf/otdfctl/commit/a88d386b3dfe6e7bf210c632c92eb54069c1c5b8))
* **core:** remove trailing slashes on host/platformEndpoint ([#415](https://github.com/opentdf/otdfctl/issues/415)) ([2ffd3c7](https://github.com/opentdf/otdfctl/commit/2ffd3c7707aa5c610f952d3499a7bfc76e8feca8)), closes [#414](https://github.com/opentdf/otdfctl/issues/414)
* **core:** revert profiles file system storage last commit ([#427](https://github.com/opentdf/otdfctl/issues/427)) ([79f2079](https://github.com/opentdf/otdfctl/commit/79f2079342bfbf210e07ce7cc6714deafea12b29))
* updates sdk to 0.3.19 with GetTdfType fixes ([#425](https://github.com/opentdf/otdfctl/issues/425)) ([0a9adfe](https://github.com/opentdf/otdfctl/commit/0a9adfe416b966b09db4b9ee60fa379db93ede76))

## [0.15.0](https://github.com/opentdf/otdfctl/compare/v0.14.0...v0.15.0) (2024-10-15)


### Features

* **core:** DSP-51 - deprecate PublicKey local field  ([#400](https://github.com/opentdf/otdfctl/issues/400)) ([1955800](https://github.com/opentdf/otdfctl/commit/1955800fcd63c4d5044517ec0355a82c0e687f1b))
* **core:** Update Resource Mapping delete to use get before delete for cli output ([#398](https://github.com/opentdf/otdfctl/issues/398)) ([79f2a42](https://github.com/opentdf/otdfctl/commit/79f2a423380cbd3f4a7805c4ec35d4657a9c0d5c))


### Bug Fixes

* **core:** build with latest opentdf releases ([#404](https://github.com/opentdf/otdfctl/issues/404)) ([969b82b](https://github.com/opentdf/otdfctl/commit/969b82b5cf90405002ac2da4a31b022dca9dfa37))

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
