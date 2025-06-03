# Changelog

## [0.21.0](https://github.com/opentdf/otdfctl/compare/v0.20.1...v0.21.0) (2025-06-03)


### Features

* add aliases to improve ux ([#263](https://github.com/opentdf/otdfctl/issues/263)) ([11db3be](https://github.com/opentdf/otdfctl/commit/11db3be2b23e2d91577c28cb8bffd8f15ff93bb1))
* add cli test mode and profile tests ([#313](https://github.com/opentdf/otdfctl/issues/313)) ([e0bc183](https://github.com/opentdf/otdfctl/commit/e0bc1836e8b5f14c87b5d572ad7937924c76d860))
* Add initial Dependency Review configuration ([#551](https://github.com/opentdf/otdfctl/issues/551)) ([b622666](https://github.com/opentdf/otdfctl/commit/b6226660c1d75e133a8ead456efcab74de4b4fc0))
* add mimetypes to ztdf on encrypt ([#209](https://github.com/opentdf/otdfctl/issues/209)) ([a918e12](https://github.com/opentdf/otdfctl/commit/a918e1201bb366d0cc18d7fd047f2c0854558e9e))
* add new global insecure flag ([#152](https://github.com/opentdf/otdfctl/issues/152)) ([777be8c](https://github.com/opentdf/otdfctl/commit/777be8ce7495f7c361a1982482af12c0fec19590))
* add profile support for cli ([#289](https://github.com/opentdf/otdfctl/issues/289)) ([15700f3](https://github.com/opentdf/otdfctl/commit/15700f3375196595e4a0ea3a7a6dea4da06d8612))
* add tdf inspection feature ([#266](https://github.com/opentdf/otdfctl/issues/266)) ([ec6befd](https://github.com/opentdf/otdfctl/commit/ec6befdfb388d85f835b2c29c293249dd2284e94))
* AR-9 first commit ([#1](https://github.com/opentdf/otdfctl/issues/1)) ([95b8723](https://github.com/opentdf/otdfctl/commit/95b87230e45f5d60b536cd879eef1298e994f318))
* Assertion verification ([#452](https://github.com/opentdf/otdfctl/issues/452)) ([5a8fe0d](https://github.com/opentdf/otdfctl/commit/5a8fe0d64088b74c95d3376e4a2a5a47d680d9c0))
* assertion verification disable ([#419](https://github.com/opentdf/otdfctl/issues/419)) ([acf5702](https://github.com/opentdf/otdfctl/commit/acf57028f1481f432b6b0c3c7a3e2c2261ac739f))
* **attributes:** get and delete for attribute CLI ([#5](https://github.com/opentdf/otdfctl/issues/5)) ([e099c0e](https://github.com/opentdf/otdfctl/commit/e099c0e83e62883856898b847c7bb81796970e82))
* **auth:** Add support for clientCredentials OIDC login flow [#63](https://github.com/opentdf/otdfctl/issues/63) ([#65](https://github.com/opentdf/otdfctl/issues/65)) ([b001351](https://github.com/opentdf/otdfctl/commit/b001351439b420425a03d0bf46faa6f1223049f6))
* **ci:** add e2e tests for subject mappings, support for --force delete ([#388](https://github.com/opentdf/otdfctl/issues/388)) ([c1f544b](https://github.com/opentdf/otdfctl/commit/c1f544b1079f52bfccb96c4c9e0b579a6854ad58))
* **ci:** Add e2e workflow setup + bats ([#234](https://github.com/opentdf/otdfctl/issues/234)) ([ca6ce5a](https://github.com/opentdf/otdfctl/commit/ca6ce5ae9aaeaafaea26ece28ba7a0ca9cdac082))
* **ci:** add tests for subject condition sets, and --force delete flag ([#389](https://github.com/opentdf/otdfctl/issues/389)) ([c6d2abc](https://github.com/opentdf/otdfctl/commit/c6d2abcd4afe78d92fd285e5c77fecdfe806ed5d)), closes [#331](https://github.com/opentdf/otdfctl/issues/331)
* **ci:** attr e2e tests with mixed casing ([#315](https://github.com/opentdf/otdfctl/issues/315)) ([50ce712](https://github.com/opentdf/otdfctl/commit/50ce712eab38f6686611e2b306bda5cacd55c28e))
* **ci:** e2e attribute definitions tests ([#384](https://github.com/opentdf/otdfctl/issues/384)) ([2894391](https://github.com/opentdf/otdfctl/commit/28943915f19e0fb565cfb38cfebdd6fde21c019a)), closes [#327](https://github.com/opentdf/otdfctl/issues/327)
* **ci:** make e2e test workflow reusable ([#365](https://github.com/opentdf/otdfctl/issues/365)) ([d94408c](https://github.com/opentdf/otdfctl/commit/d94408cc2898d46b3444e874c035ff2bffe451f4))
* **ci:** namespaces e2e tests and test suite improvements ([#351](https://github.com/opentdf/otdfctl/issues/351)) ([ce28555](https://github.com/opentdf/otdfctl/commit/ce285554866bf89ee8aa2df4a4b426548a58b59a))
* **ci:** reusable platform composite action in e2e tests ([#369](https://github.com/opentdf/otdfctl/issues/369)) ([f7d5a1c](https://github.com/opentdf/otdfctl/commit/f7d5a1c07304bee14dfc92fa81bd65389e76d9f6))
* **core:** add `subject-mappings match` to CLI ([#413](https://github.com/opentdf/otdfctl/issues/413)) ([bc56c19](https://github.com/opentdf/otdfctl/commit/bc56c199a73b12b8c90045d1b6f9cc6fdec16c54))
* **core:** add aliases for profile command ([#510](https://github.com/opentdf/otdfctl/issues/510)) ([45c633d](https://github.com/opentdf/otdfctl/commit/45c633da6b00b04a8c92686521d25144048ac62c))
* **core:** Add base key cmds ([#563](https://github.com/opentdf/otdfctl/issues/563)) ([edfd6c0](https://github.com/opentdf/otdfctl/commit/edfd6c08dc9b84f2cbfc79643ccc266a45ce58fd))
* **core:** add cli attribute values crud ([#49](https://github.com/opentdf/otdfctl/issues/49)) ([a7c8955](https://github.com/opentdf/otdfctl/commit/a7c8955aba68e9c6e883d1ea16329a6a3e455237))
* **core:** Add commands to encrypt and decrypt nano tdfs ([#168](https://github.com/opentdf/otdfctl/issues/168)) ([0659cb8](https://github.com/opentdf/otdfctl/commit/0659cb83754904bb1f6715c11acb2acc51fe15de))
* **core:** add ecdsa-binding encrypt flag ([#360](https://github.com/opentdf/otdfctl/issues/360)) ([8702ec0](https://github.com/opentdf/otdfctl/commit/8702ec007b6d1354b6c0366e6b375f26216dfde1))
* **core:** add kas-registry crud commands ([#57](https://github.com/opentdf/otdfctl/issues/57)) ([c340a22](https://github.com/opentdf/otdfctl/commit/c340a22a450a584962d3aa7eae00b54528644254))
* **core:** add license ([#247](https://github.com/opentdf/otdfctl/issues/247)) ([ebf8d50](https://github.com/opentdf/otdfctl/commit/ebf8d504da5db16954e4e7634b120a95b2ef5104))
* **core:** add metadata behavior and fix subject mappings ([#67](https://github.com/opentdf/otdfctl/issues/67)) ([acb88cf](https://github.com/opentdf/otdfctl/commit/acb88cfdfef8f5c8e997178f0f1ca6a2ca8c7361))
* **core:** add metadata to table output ([#112](https://github.com/opentdf/otdfctl/issues/112)) ([54340a7](https://github.com/opentdf/otdfctl/commit/54340a77faca5056fe9ea9218b4fc3573839806f))
* **core:** add optional name to kas registry CRUD commands ([#429](https://github.com/opentdf/otdfctl/issues/429)) ([f675d86](https://github.com/opentdf/otdfctl/commit/f675d86c83205232db407d6609e80fa865a3998e))
* **core:** add print-access-token auth subcommand for ease of DX and piping to other CLI tools ([#135](https://github.com/opentdf/otdfctl/issues/135)) ([d0ac710](https://github.com/opentdf/otdfctl/commit/d0ac71015cf3a3458a6b78b755af22998876169a)), closes [#136](https://github.com/opentdf/otdfctl/issues/136)
* **core:** add scaffolding and POC for auth code flow ([#144](https://github.com/opentdf/otdfctl/issues/144)) ([03ecbfb](https://github.com/opentdf/otdfctl/commit/03ecbfb4f689f4a9f161a5a03d80efd50f728780))
* **core:** Add support for WithTargetMode encrypt option ([#519](https://github.com/opentdf/otdfctl/issues/519)) ([a0ab213](https://github.com/opentdf/otdfctl/commit/a0ab2136be0b1d39e16a7522210f493fd797089d))
* **core:** add TDF encrypt/decrypt commands and authentication to platform SDK ([#115](https://github.com/opentdf/otdfctl/issues/115)) ([3c50c2d](https://github.com/opentdf/otdfctl/commit/3c50c2decd55a442e74389158159534e6730fd7e))
* **core:** add unsafe attribute definition CLI commands ([#214](https://github.com/opentdf/otdfctl/issues/214)) ([7f7bc70](https://github.com/opentdf/otdfctl/commit/7f7bc70c625bcb429ab1f11d917e7b37a438736e))
* **core:** add unsafe namespace commands ([#213](https://github.com/opentdf/otdfctl/issues/213)) ([a7156d8](https://github.com/opentdf/otdfctl/commit/a7156d8da786e8cbf0119bfe62efe6800b9ab962))
* **core:** Adding examples docs, mainly policy commands ([#461](https://github.com/opentdf/otdfctl/issues/461)) ([04c1743](https://github.com/opentdf/otdfctl/commit/04c17439bb5f68fb5d44ba96cb457ce9ca072250))
* **core:** adds assertions to encrypt subcommand ([#408](https://github.com/opentdf/otdfctl/issues/408)) ([8f0e906](https://github.com/opentdf/otdfctl/commit/8f0e906c1dfe99fe6aa5f2ff43d02f0da90474cf))
* **core:** adds missing long manual output docs ([#362](https://github.com/opentdf/otdfctl/issues/362)) ([8e1390f](https://github.com/opentdf/otdfctl/commit/8e1390f20c17a5900c586f94384af76ffd9a2844)), closes [#359](https://github.com/opentdf/otdfctl/issues/359)
* **core:** adds storeFile to save encrypted profiles to disk and updates auth to propagate tlsNoVerify ([#420](https://github.com/opentdf/otdfctl/issues/420)) ([f709e01](https://github.com/opentdf/otdfctl/commit/f709e014bf3f82a2808eae5df76b3667730c36ef))
* **core:** allow no-cache client credentials options throughout CLI ([#142](https://github.com/opentdf/otdfctl/issues/142)) ([c20dc1e](https://github.com/opentdf/otdfctl/commit/c20dc1e81f9299da6f149f86bef1c2de63b92787))
* **core:** bump SDK and consume new platform connection validation ([#493](https://github.com/opentdf/otdfctl/issues/493)) ([1106b54](https://github.com/opentdf/otdfctl/commit/1106b54e73f9ceb711ff19d15cd08bf1cebbb29f))
* **core:** create subject mapping with subject condition sets ([#79](https://github.com/opentdf/otdfctl/issues/79)) ([17fbde0](https://github.com/opentdf/otdfctl/commit/17fbde09cf5b61670d20036eaf3192b0282c953b))
* **core:** DSP-51 - deprecate PublicKey local field  ([#400](https://github.com/opentdf/otdfctl/issues/400)) ([1955800](https://github.com/opentdf/otdfctl/commit/1955800fcd63c4d5044517ec0355a82c0e687f1b))
* **core:** DSPX-18 clean up Go context usage to follow best practices ([#558](https://github.com/opentdf/otdfctl/issues/558)) ([a2c9f8b](https://github.com/opentdf/otdfctl/commit/a2c9f8b13cbab740b46262f70aecc82a94f3d788))
* **core:** DSPX-608 - Deprecate public_client_id ([#555](https://github.com/opentdf/otdfctl/issues/555)) ([8d396bd](https://github.com/opentdf/otdfctl/commit/8d396bd022126524d9d20daa03ec6ca262cf4406))
* **core:** DSPX-608 - require clientID for login ([#553](https://github.com/opentdf/otdfctl/issues/553)) ([580172e](https://github.com/opentdf/otdfctl/commit/580172e1861b54366f4914a141e459fe3221a16d))
* **core:** DSPX-896 add registered resources CRUD ([#559](https://github.com/opentdf/otdfctl/issues/559)) ([8e7475e](https://github.com/opentdf/otdfctl/commit/8e7475ef8aab91d28ab7efd320af13dc5ab53d3b))
* **core:** enable man build to be dynamic ([#76](https://github.com/opentdf/otdfctl/issues/76)) ([1a03039](https://github.com/opentdf/otdfctl/commit/1a0303970652839803ba465921b59b98fe3b3ec0))
* **core:** enable mounting to another cobra root ([#162](https://github.com/opentdf/otdfctl/issues/162)) ([bb019bd](https://github.com/opentdf/otdfctl/commit/bb019bd944f5b6dbb6d3992725f2a45455d31214))
* **core:** enable setting KAS url path on encrypt ([#225](https://github.com/opentdf/otdfctl/issues/225)) ([0695c96](https://github.com/opentdf/otdfctl/commit/0695c96296cec708f4e39b36b37319c946f51c50))
* **core:** enable styled and json default outputs ([#66](https://github.com/opentdf/otdfctl/issues/66)) ([69b7792](https://github.com/opentdf/otdfctl/commit/69b7792e4b7788ae82d653b64b1e0d68b91729a8))
* **core:** export manual functions for CLI wrappers to consume ([#397](https://github.com/opentdf/otdfctl/issues/397)) ([aa0bf95](https://github.com/opentdf/otdfctl/commit/aa0bf95a39dfc0aec4155e498a2096cbd158efdd))
* **core:** hide the dev command from menu ([#174](https://github.com/opentdf/otdfctl/issues/174)) ([3061a90](https://github.com/opentdf/otdfctl/commit/3061a900c87a30e4b470dd6d8bded442436b9f10))
* **core:** KAS allowlist options ([#539](https://github.com/opentdf/otdfctl/issues/539)) ([af7978f](https://github.com/opentdf/otdfctl/commit/af7978f86ced38543b31b792e008654071333789))
* **core:** kas-grants CRUD ([#80](https://github.com/opentdf/otdfctl/issues/80)) ([f53b61d](https://github.com/opentdf/otdfctl/commit/f53b61dd79c463e363cf5340517fea74c491cdba))
* **core:** kas-grants list ([#346](https://github.com/opentdf/otdfctl/issues/346)) ([7f51282](https://github.com/opentdf/otdfctl/commit/7f512825eab814e3c130e3fe4e8ed85ecbe2d146)), closes [#253](https://github.com/opentdf/otdfctl/issues/253)
* **core:** kasr cached keys to deprecate local ([#318](https://github.com/opentdf/otdfctl/issues/318)) ([5419cc3](https://github.com/opentdf/otdfctl/commit/5419cc39e143eb484f836ca1ee671d626d5e2c60)), closes [#317](https://github.com/opentdf/otdfctl/issues/317)
* **core:** key management operations ([#533](https://github.com/opentdf/otdfctl/issues/533)) ([d4f6aaa](https://github.com/opentdf/otdfctl/commit/d4f6aaac3f6fc1b50fbc988e5d34a32de0ed9f64))
* **core:** list attributes with state ([#90](https://github.com/opentdf/otdfctl/issues/90)) ([01bb38d](https://github.com/opentdf/otdfctl/commit/01bb38d0827c4269fe5e35c5df32b9a37ffbb2e4))
* **core:** make root command exportable and add config bootstrapping ([#120](https://github.com/opentdf/otdfctl/issues/120)) ([564dc68](https://github.com/opentdf/otdfctl/commit/564dc68e6e4f207dcd285c25858d1f2f757f733d))
* **core:** pagination of LIST commands ([#447](https://github.com/opentdf/otdfctl/issues/447)) ([673a064](https://github.com/opentdf/otdfctl/commit/673a06424d30e706798b9a1fa1bbfd9b4601e765))
* **core:** require host flag ([#167](https://github.com/opentdf/otdfctl/issues/167)) ([2f466f1](https://github.com/opentdf/otdfctl/commit/2f466f198ad37ca9fbcef058509c70c97be8e68d))
* **core:** Resource mapping groups ([#567](https://github.com/opentdf/otdfctl/issues/567)) ([03fa307](https://github.com/opentdf/otdfctl/commit/03fa307b3ab91f25baeb74e30fde6eeec6d479a1))
* **core:** resource mappings LIST fix, delete --force support, and e2e tests ([#387](https://github.com/opentdf/otdfctl/issues/387)) ([326e74b](https://github.com/opentdf/otdfctl/commit/326e74b37d0abfb4ad50deadaa1ed46ecf9f8a5d)), closes [#386](https://github.com/opentdf/otdfctl/issues/386)
* **core:** Rotate key. ([#572](https://github.com/opentdf/otdfctl/issues/572)) ([afd0043](https://github.com/opentdf/otdfctl/commit/afd0043f1ea66f0b371a95b556320551f73749bb))
* **core:** Shows SDK version and spec info ([#474](https://github.com/opentdf/otdfctl/issues/474)) ([5a685c4](https://github.com/opentdf/otdfctl/commit/5a685c4e36cf524c4f594fac42cfec30f62a6e83))
* **core:** subject condition set CLI CRUD ([#78](https://github.com/opentdf/otdfctl/issues/78)) ([26f6fcc](https://github.com/opentdf/otdfctl/commit/26f6fcca5115eb9ad87df0e6686a2d632ef0cc35))
* **core:** subject condition set prune ([#439](https://github.com/opentdf/otdfctl/issues/439)) ([c4c8b8b](https://github.com/opentdf/otdfctl/commit/c4c8b8b276b2189df74e6cf30e14abac9369d97e))
* **core:** support for ec-wrapping ([#499](https://github.com/opentdf/otdfctl/issues/499)) ([e839445](https://github.com/opentdf/otdfctl/commit/e839445181c89447d9a2374d54ce5ea4c3f46320))
* **core:** support kas grants to namespaces ([#292](https://github.com/opentdf/otdfctl/issues/292)) ([f2c6689](https://github.com/opentdf/otdfctl/commit/f2c6689d2f775b1aed907d553c42d87c8464e6c7)), closes [#269](https://github.com/opentdf/otdfctl/issues/269)
* **core:** support update of SCS read from JSON file ([#250](https://github.com/opentdf/otdfctl/issues/250)) ([ebc16ea](https://github.com/opentdf/otdfctl/commit/ebc16ea8caf11f71fe04de0ce1e0acb24eb23edc)), closes [#197](https://github.com/opentdf/otdfctl/issues/197)
* **core:** unsafe values CLI functionality ([#218](https://github.com/opentdf/otdfctl/issues/218)) ([77340d1](https://github.com/opentdf/otdfctl/commit/77340d1d3a0c43a376563ef587171b9900970405))
* **core:** Update Resource Mapping delete to use get before delete for cli output ([#398](https://github.com/opentdf/otdfctl/issues/398)) ([79f2a42](https://github.com/opentdf/otdfctl/commit/79f2a423380cbd3f4a7805c4ec35d4657a9c0d5c))
* **core:** update to use the new sdk wrapper ([#13](https://github.com/opentdf/otdfctl/issues/13)) ([9c40a18](https://github.com/opentdf/otdfctl/commit/9c40a186065fa5a8aed9e69a8284aea6349592d9))
* **core:** zip artifacts and generate checksums when releasing ([#84](https://github.com/opentdf/otdfctl/issues/84)) ([d4cd22d](https://github.com/opentdf/otdfctl/commit/d4cd22dabe21e0c3072ddb2bb5581f850bc30795))
* **demo:** adds dev subcommand for jq selectors to generate and test selectors on a subject context JSON or JWT ([#91](https://github.com/opentdf/otdfctl/issues/91)) ([fa5f959](https://github.com/opentdf/otdfctl/commit/fa5f95974999ef51aee88e5b649869afd486cd32))
* **dependabot:** use squash instead of merge commit when dependabot GHA approves PRs ([#104](https://github.com/opentdf/otdfctl/issues/104)) ([1c2b79b](https://github.com/opentdf/otdfctl/commit/1c2b79bac4f719e065d8a811517aa7c1a3656a0c))
* improve auth with client credentials ([#286](https://github.com/opentdf/otdfctl/issues/286)) ([9c4968f](https://github.com/opentdf/otdfctl/commit/9c4968f48d1ba23a61ed5c8ad23a109bf141ba56))
* improve auth with client credentials ([#296](https://github.com/opentdf/otdfctl/issues/296)) ([0f533c7](https://github.com/opentdf/otdfctl/commit/0f533c7278a53ddd90656b3c7efcaee1c5bfd957))
* improve table experience with more dynamic library ([#200](https://github.com/opentdf/otdfctl/issues/200)) ([f199fe3](https://github.com/opentdf/otdfctl/commit/f199fe3a86a7bf7a4d3473ce63edc4156d2d5530))
* **issue 11:** subject mappings CRUD via CLI ([#41](https://github.com/opentdf/otdfctl/issues/41)) ([3db35f6](https://github.com/opentdf/otdfctl/commit/3db35f667956d12c482d8885a41d51e7f384ca90))
* **issue 27:** CLI CRUD for namespaces after updates to policy config schema refactoring ([#28](https://github.com/opentdf/otdfctl/issues/28)) ([d96ab22](https://github.com/opentdf/otdfctl/commit/d96ab227d0923b05972c41d329e149269334170d))
* **issue 39:** quality enforcement pipeline ([#42](https://github.com/opentdf/otdfctl/issues/42)) ([68cabd1](https://github.com/opentdf/otdfctl/commit/68cabd19ba30fa5bec9626bb7b730be6ef2182c7))
* **issue-22:** Implement makefile to build all trucli targets (and other task running needs) ([#45](https://github.com/opentdf/otdfctl/issues/45)) ([11b7f82](https://github.com/opentdf/otdfctl/commit/11b7f82b710c0be24b3ec8bff368b25a79d5ea36))
* **issue#37:** implement ci build workflow ([#53](https://github.com/opentdf/otdfctl/issues/53)) ([bdfcc0f](https://github.com/opentdf/otdfctl/commit/bdfcc0f5b889266f35858da5b31c0b3972472fa9))
* **main:** add actions CRUD and e2e tests ([#523](https://github.com/opentdf/otdfctl/issues/523)) ([2fb9ec7](https://github.com/opentdf/otdfctl/commit/2fb9ec7336da5731b868da94f0bbd5b2f226ede1))
* **main:** refactor actions within existing CLI policy object CRUD ([#543](https://github.com/opentdf/otdfctl/issues/543)) ([9ab1a58](https://github.com/opentdf/otdfctl/commit/9ab1a58418643ea709aefb08e3f5ca8bd06235f4))
* move git checkout before tagging ([#298](https://github.com/opentdf/otdfctl/issues/298)) ([1114e25](https://github.com/opentdf/otdfctl/commit/1114e25a90946e85622c8ff7a7befbf18beb4ba1))
* **policy:** add cli crud for attributes ([#48](https://github.com/opentdf/otdfctl/issues/48)) ([ec70a83](https://github.com/opentdf/otdfctl/commit/ec70a83ff9256835a6273110f797827a09735003))
* refactor encrypt and decrypt + CLI examples ([#418](https://github.com/opentdf/otdfctl/issues/418)) ([e681823](https://github.com/opentdf/otdfctl/commit/e681823ad54ddf70f4aa2215438d69a3d02cf6eb))
* **resource-encodings:** add resource encoding CRUD for CLI ([#56](https://github.com/opentdf/otdfctl/issues/56)) ([0bb961f](https://github.com/opentdf/otdfctl/commit/0bb961f9ea1f88b3277e4cf03317ef112a7cd1bc)), closes [#7](https://github.com/opentdf/otdfctl/issues/7)
* support --with-access-token for auth ([#409](https://github.com/opentdf/otdfctl/issues/409)) ([856efa4](https://github.com/opentdf/otdfctl/commit/856efa4d61bb24b05f3a98943b94600ff77536fa))
* **tui:** abstract away read view ([#141](https://github.com/opentdf/otdfctl/issues/141)) ([e6f44e0](https://github.com/opentdf/otdfctl/commit/e6f44e0b06422b0676ce5f7f96630fc80a839f7b))
* **tui:** add attribute view ([#20](https://github.com/opentdf/otdfctl/issues/20)) ([3e1a4e0](https://github.com/opentdf/otdfctl/commit/3e1a4e0d9bcaca0ce665b9b26efdaa4682845b92))
* **tui:** add CRUD attribute views ([#54](https://github.com/opentdf/otdfctl/issues/54)) ([2cc37af](https://github.com/opentdf/otdfctl/commit/2cc37af0fb5e2f36154565feb49b32cb3353ef7e))
* **tui:** update view ([#156](https://github.com/opentdf/otdfctl/issues/156)) ([09006e7](https://github.com/opentdf/otdfctl/commit/09006e7a1a8b07879ea59729778cb1adc419c1c9))
* update sdk to new refactored version ([#44](https://github.com/opentdf/otdfctl/issues/44)) ([78bc086](https://github.com/opentdf/otdfctl/commit/78bc086727d05c12e8bcf3d235b9ba87b9bf5710)), closes [#43](https://github.com/opentdf/otdfctl/issues/43)
* **update:** add update attributes handler and Cobra CLI cmd ([#4](https://github.com/opentdf/otdfctl/issues/4)) ([d326bd0](https://github.com/opentdf/otdfctl/commit/d326bd075102124820c91303f0600c75a6ddc79a))
* use wellknown endpoint for idp ([#187](https://github.com/opentdf/otdfctl/issues/187)) ([83b0ec8](https://github.com/opentdf/otdfctl/commit/83b0ec85649fb63d77214e1d2b70470128004240))


### Bug Fixes

* add version flag ([#270](https://github.com/opentdf/otdfctl/issues/270)) ([3e20e9e](https://github.com/opentdf/otdfctl/commit/3e20e9eb2b4c541f28527a13c6bc78d81b30f904))
* bump platform/sdk to 0.2.8 ([#206](https://github.com/opentdf/otdfctl/issues/206)) ([bab5151](https://github.com/opentdf/otdfctl/commit/bab515150740eaac26b6a261a5cc3529acc5630b))
* bump sdk version to 0.3.1 ([#230](https://github.com/opentdf/otdfctl/issues/230)) ([b5e73aa](https://github.com/opentdf/otdfctl/commit/b5e73aa557dd636a10af5b729e535b12ef9d41af))
* change name of --insecure flag to --tls-no-verify ([#158](https://github.com/opentdf/otdfctl/issues/158)) ([52adfc3](https://github.com/opentdf/otdfctl/commit/52adfc3c5c77a73198d938570193056e42f404d6))
* **ci:** ci job should run on changes to GHA ([#530](https://github.com/opentdf/otdfctl/issues/530)) ([1d296ca](https://github.com/opentdf/otdfctl/commit/1d296ca8fac889a6e776ad381df999a2fcf9d6ce))
* **ci:** e2e workflow should be fully reusable ([#368](https://github.com/opentdf/otdfctl/issues/368)) ([cc1e2b9](https://github.com/opentdf/otdfctl/commit/cc1e2b938fb0c8c4cf64d735f2961f7c9cae79fa))
* **ci:** enhance lint config and resolve all lint issues ([#363](https://github.com/opentdf/otdfctl/issues/363)) ([5c1dbf1](https://github.com/opentdf/otdfctl/commit/5c1dbf1f5e441ca0ebd8cfcca145a77b623f3638))
* **core:** align kas grant commands with RPCs ([#264](https://github.com/opentdf/otdfctl/issues/264)) ([269171a](https://github.com/opentdf/otdfctl/commit/269171a7624478d1ee50b712287331ee05b5001d))
* **core:** build with latest opentdf releases ([#404](https://github.com/opentdf/otdfctl/issues/404)) ([969b82b](https://github.com/opentdf/otdfctl/commit/969b82b5cf90405002ac2da4a31b022dca9dfa37))
* **core:** bump jwt dep and remove outdated version ([#520](https://github.com/opentdf/otdfctl/issues/520)) ([77bb9ca](https://github.com/opentdf/otdfctl/commit/77bb9ca9a0741ab7b920cc00f264a021064b117c))
* **core:** bump platform deps ([#276](https://github.com/opentdf/otdfctl/issues/276)) ([e4ced99](https://github.com/opentdf/otdfctl/commit/e4ced996ae336b9db6db88906683f6600a2e5bf4))
* **core:** dev selectors employ flattening from platform instead of jq ([#411](https://github.com/opentdf/otdfctl/issues/411)) ([57966ff](https://github.com/opentdf/otdfctl/commit/57966ffadcc61e1611869171bd3fc85723492fb7))
* **core:** do not import unused fmt ([#306](https://github.com/opentdf/otdfctl/issues/306)) ([0dc552d](https://github.com/opentdf/otdfctl/commit/0dc552d3d6814f910c04d5f8cefa35404b4945f5))
* **core:** fix bug where empty string value caused json parse error ([#85](https://github.com/opentdf/otdfctl/issues/85)) ([f78c0e9](https://github.com/opentdf/otdfctl/commit/f78c0e9556e30c2392c785d26fae70dda4e83f47))
* **core:** fix error handler to avoid metadata labels panic ([#219](https://github.com/opentdf/otdfctl/issues/219)) ([2747360](https://github.com/opentdf/otdfctl/commit/2747360e238634eb87bd25089a3349e082a68e42))
* **core:** fix LIST command helper output to include all subcommands up to root ([#201](https://github.com/opentdf/otdfctl/issues/201)) ([b856607](https://github.com/opentdf/otdfctl/commit/b856607ed6c202f367ec21cc9372656bba73cb2b))
* **core:** fix regression where SCS list unavailable ([#249](https://github.com/opentdf/otdfctl/issues/249)) ([732f56b](https://github.com/opentdf/otdfctl/commit/732f56b256f5fda7ad508e776098506f5ad9a70e)), closes [#198](https://github.com/opentdf/otdfctl/issues/198)
* **core:** fix single resource get output ([#212](https://github.com/opentdf/otdfctl/issues/212)) ([5401418](https://github.com/opentdf/otdfctl/commit/54014188af3a2bc7c87c9c9a95046e5928dce29b)), closes [#211](https://github.com/opentdf/otdfctl/issues/211)
* **core:** fix subject-condition-sets Create/Update with protojson marshaling ([#245](https://github.com/opentdf/otdfctl/issues/245)) ([e6afec4](https://github.com/opentdf/otdfctl/commit/e6afec43f1d9591d8c35ffcaa201b92b24f1d81a))
* **core:** Fixes piped input parser on decrypt ([#224](https://github.com/opentdf/otdfctl/issues/224)) ([a375ddb](https://github.com/opentdf/otdfctl/commit/a375ddbcfc801d44af9ad6eddbb98109a532a5a6))
* **core:** GOOS, error message fixes ([#378](https://github.com/opentdf/otdfctl/issues/378)) ([623a82a](https://github.com/opentdf/otdfctl/commit/623a82ad3c1ed698a83eed54cf15a4f552096728)), closes [#380](https://github.com/opentdf/otdfctl/issues/380)
* **core:** improve KASR docs and add spellcheck GHA to pipeline ([#323](https://github.com/opentdf/otdfctl/issues/323)) ([a77cf30](https://github.com/opentdf/otdfctl/commit/a77cf30dc8077d034cb4c9df8cc94712b1a17dff)), closes [#335](https://github.com/opentdf/otdfctl/issues/335) [#337](https://github.com/opentdf/otdfctl/issues/337)
* **core:** improve readability of TDF methods ([#424](https://github.com/opentdf/otdfctl/issues/424)) ([a88d386](https://github.com/opentdf/otdfctl/commit/a88d386b3dfe6e7bf210c632c92eb54069c1c5b8))
* **core:** kas registry get should allow -i 'id' flag shorthand ([#434](https://github.com/opentdf/otdfctl/issues/434)) ([bed3701](https://github.com/opentdf/otdfctl/commit/bed3701d89510ee78c3aed43b1a072e41ee3873f))
* **core:** kas-grants cmd ([#82](https://github.com/opentdf/otdfctl/issues/82)) ([0c0331f](https://github.com/opentdf/otdfctl/commit/0c0331f2b2a95ef44ba9c05fd87fe847ace106c8))
* **core:** kasr creation JSON example ([#453](https://github.com/opentdf/otdfctl/issues/453)) ([192c7b2](https://github.com/opentdf/otdfctl/commit/192c7b2975a4ab6f648ab7924e20e70535ce04b2))
* **core:** mark new algorithm flags experimental ([#501](https://github.com/opentdf/otdfctl/issues/501)) ([95e00bf](https://github.com/opentdf/otdfctl/commit/95e00bf3daa8eb05196a5839488a4718c2230210))
* **core:** metadata rendering cleanup ([#293](https://github.com/opentdf/otdfctl/issues/293)) ([ed21f81](https://github.com/opentdf/otdfctl/commit/ed21f81863450fd6167106711392e713a43c55be))
* **core:** move json global flag to the policy subcommand where it is always relevant ([#117](https://github.com/opentdf/otdfctl/issues/117)) ([2ca6151](https://github.com/opentdf/otdfctl/commit/2ca6151e127775f4712680b747c396bf001aca70))
* **core:** nil panic on set-default ([#304](https://github.com/opentdf/otdfctl/issues/304)) ([92bbfa3](https://github.com/opentdf/otdfctl/commit/92bbfa32ae42b73b68551c2f9d3551d357bc5922))
* **core:** pull correct flag for tls-no-verify on auth cmd ([#165](https://github.com/opentdf/otdfctl/issues/165)) ([d788780](https://github.com/opentdf/otdfctl/commit/d7887805564f717862373832187922125810cc7e)), closes [#164](https://github.com/opentdf/otdfctl/issues/164)
* **core:** remove deprecated policy members ([#231](https://github.com/opentdf/otdfctl/issues/231)) ([038ce1c](https://github.com/opentdf/otdfctl/commit/038ce1c276d16993fb7747ef374210fcdf66d9e4))
* **core:** remove documentation that cached kas pubkey is base64 ([#320](https://github.com/opentdf/otdfctl/issues/320)) ([fce8f44](https://github.com/opentdf/otdfctl/commit/fce8f44f767f35ccc4863f88d46e7ffcbd80f37a)), closes [#321](https://github.com/opentdf/otdfctl/issues/321)
* **core:** remove duplicate titling of help manual ([#391](https://github.com/opentdf/otdfctl/issues/391)) ([cb8db69](https://github.com/opentdf/otdfctl/commit/cb8db69ec4df42c7f230fbd87142bfbcd2d3940f))
* **core:** remove trailing slashes on host/platformEndpoint ([#415](https://github.com/opentdf/otdfctl/issues/415)) ([2ffd3c7](https://github.com/opentdf/otdfctl/commit/2ffd3c7707aa5c610f952d3499a7bfc76e8feca8)), closes [#414](https://github.com/opentdf/otdfctl/issues/414)
* **core:** revert profiles file system storage last commit ([#427](https://github.com/opentdf/otdfctl/issues/427)) ([79f2079](https://github.com/opentdf/otdfctl/commit/79f2079342bfbf210e07ce7cc6714deafea12b29))
* **core:** sm list should provide value fqn instead of just value string ([#438](https://github.com/opentdf/otdfctl/issues/438)) ([9a7cb72](https://github.com/opentdf/otdfctl/commit/9a7cb7242e0e39ccc2b54425028638fa0c5e3f9f))
* **core:** specify cli team rather than all opentdf developers as codeowners of otdfctl ([#110](https://github.com/opentdf/otdfctl/issues/110)) ([f63b9a6](https://github.com/opentdf/otdfctl/commit/f63b9a68f36c954d2fb870a7020a28c5d56c9e5a))
* **core:** sync up auth client credentials flags ([#190](https://github.com/opentdf/otdfctl/issues/190)) ([1503537](https://github.com/opentdf/otdfctl/commit/1503537a121f174a9233a11b31c55ad70e2a402e)), closes [#189](https://github.com/opentdf/otdfctl/issues/189)
* **core:** values list should properly render table output ([#220](https://github.com/opentdf/otdfctl/issues/220)) ([e972db4](https://github.com/opentdf/otdfctl/commit/e972db41f68005e510c698b6f96281dabe0ec9f0))
* **core:** warn and do now allow deletion of default profile ([#308](https://github.com/opentdf/otdfctl/issues/308)) ([fdd8167](https://github.com/opentdf/otdfctl/commit/fdd8167e8e2b22d652b48d796a756f86398bfd3c))
* **core:** wire attribute value FQNs to encrypt ([#370](https://github.com/opentdf/otdfctl/issues/370)) ([21f9b80](https://github.com/opentdf/otdfctl/commit/21f9b80cdee7d695a308937b08dbc768d11fbbd5))
* create new http client to ignore tls verification ([#324](https://github.com/opentdf/otdfctl/issues/324)) ([4d4afb7](https://github.com/opentdf/otdfctl/commit/4d4afb7e5b6411bb08a92bc53181ac5730ca1992))
* **demo:** remove selectors command while platform undergoes changes ([#172](https://github.com/opentdf/otdfctl/issues/172)) ([f3f2a51](https://github.com/opentdf/otdfctl/commit/f3f2a5147a521909143ccd266117fa3d7e9d16bd))
* disable tagging ([#302](https://github.com/opentdf/otdfctl/issues/302)) ([2b5db85](https://github.com/opentdf/otdfctl/commit/2b5db852ed0088e61f1180500135cd1865f9798b))
* **Issue#62:** Fix breakage from sdk update ([#64](https://github.com/opentdf/otdfctl/issues/64)) ([4c16885](https://github.com/opentdf/otdfctl/commit/4c168855e78f0f9af29551f4ab9a91ffa294f949))
* **main:** Pass the full url when building the sdk object ([#544](https://github.com/opentdf/otdfctl/issues/544)) ([8b836f0](https://github.com/opentdf/otdfctl/commit/8b836f0fa3aa414c3ab19d830f4d1f833d3ae61d))
* make file not building correctly ([#307](https://github.com/opentdf/otdfctl/issues/307)) ([64eb821](https://github.com/opentdf/otdfctl/commit/64eb82170fdcc50396194271be358bf9c9d43049))
* reduce prints ([#277](https://github.com/opentdf/otdfctl/issues/277)) ([8b5734a](https://github.com/opentdf/otdfctl/commit/8b5734a18636071566fd8c4cfc808f3f240a02a5))
* refactor to support varying print output ([#350](https://github.com/opentdf/otdfctl/issues/350)) ([d6932f3](https://github.com/opentdf/otdfctl/commit/d6932f30d9f653e46b32761a3257f3555ef0a6eb))
* release-please tweak ([#300](https://github.com/opentdf/otdfctl/issues/300)) ([29fc836](https://github.com/opentdf/otdfctl/commit/29fc8360ae0b701aefe70b25d1838f442fd7eb8d))
* update key mgmt flags to consistent format ([#570](https://github.com/opentdf/otdfctl/issues/570)) ([846f96c](https://github.com/opentdf/otdfctl/commit/846f96cb9adfe03e355c9e64b559f1c11d84a86f))
* update workflow permissions ([#310](https://github.com/opentdf/otdfctl/issues/310)) ([3979fe8](https://github.com/opentdf/otdfctl/commit/3979fe85c9ab6511376d98b672cbfebddbf9bb84))
* updates sdk to 0.3.19 with GetTdfType fixes ([#425](https://github.com/opentdf/otdfctl/issues/425)) ([0a9adfe](https://github.com/opentdf/otdfctl/commit/0a9adfe416b966b09db4b9ee60fa379db93ede76))
* use right kas-grants flag names when retrieving values ([#228](https://github.com/opentdf/otdfctl/issues/228)) ([f8c3e9a](https://github.com/opentdf/otdfctl/commit/f8c3e9ad992ce10b14af467e0441b9d68057beb2)), closes [#227](https://github.com/opentdf/otdfctl/issues/227)

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
