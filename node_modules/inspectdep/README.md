inspectdep 🔎
============

[![npm version][npm_img]][npm_site]
[![Travis Status][trav_img]][trav_site]
[![AppVeyor Status][appveyor_img]][appveyor_site]
[![Coverage Status][cov_img]][cov_site]

An inspection tool for dependencies in `node_modules`.

## API

### `findProdInstalls({ rootPath })`

Find on-disk locations of all production dependencies in `node_modules`.

_Notes_:

* This includes all `dependencies` and `optionalDependencies`, simulating what would happen during a `yarn|npm install --production`.
* Paths are relative to `rootPath` and separated with `path.sep` native OS separators
* If dependencies are not found installed on-disk they are simply ignored.
  [#2](https://github.com/FormidableLabs/inspectdep/issues/2)

_Parameters_:

* `rootPath` (`string`): `node_modules` root location (default: `process.cwd()`)
* `curPath` (`string`): location to start inferring dependencies from (default: `rootPath`)

_Returns_:

* (`Promise<Array<String>>`): list of relative paths to on-disk dependencies

[npm_img]: https://badge.fury.io/js/inspectdep.svg
[npm_site]: http://badge.fury.io/js/inspectdep
[trav_img]: https://api.travis-ci.com/FormidableLabs/inspectdep.svg
[trav_site]: https://travis-ci.com/FormidableLabs/inspectdep
[appveyor_img]: https://ci.appveyor.com/api/projects/status/github/formidablelabs/inspectdep?branch=master&svg=true
[appveyor_site]: https://ci.appveyor.com/project/FormidableLabs/inspectdep
[cov_img]: https://codecov.io/gh/FormidableLabs/inspectdep/branch/master/graph/badge.svg
[cov_site]: https://codecov.io/gh/FormidableLabs/inspectdep