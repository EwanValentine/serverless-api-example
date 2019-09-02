"use strict";

const path = require("path");
const { readJson, exists, findPkg } = require("./util");

const binOrNull = (rootPath, binName) => {
  const binPath = path.normalize(`node_modules/.bin/${binName}`);
  return exists(path.resolve(rootPath, binPath))
    .then((isFound) => isFound ? binPath : null);
};

// Recursively traverse package.json + path to resolve all on-disk locations
const resolveLocations = async ({ rootPath, curPath, pkg, visited }) => {
  // Track visited paths to detect cycles.
  visited = visited || new Set();

  // Check if previously visited for cycles, and start tracking if not.
  // TODO(2): Treating circular dependency as a "not found" upon cycle.
  if (visited.has(curPath)) { return null; }
  visited.add(curPath);

  // Production dependencies include both production and optional if found.
  const names = []
    .concat(Object.keys(pkg.dependencies || {}))
    .concat(Object.keys(pkg.optionalDependencies || {}));

  const locs = await Promise.all(names.map(async (name) => {
    // Find current package.
    const found = await findPkg({ rootPath, curPath, name });
    if (!found) {
      // TODO(2): Handle not found and/or filter nulls.
      // https://github.com/FormidableLabs/inspectdep/issues/2
      // Potential information:
      // - `const isOptional = !(pkg.dependencies || {})[name];`
      // - `const parent = curPath;`
      // - `name`
      return null;
    }

    // Recurse into dependencies.
    const deps = await resolveLocations({ rootPath, curPath: found.loc, pkg: found.pkg, visited });

    return [path.relative(rootPath, found.loc)].concat(deps);
  }));

  // Coerce to bin "object" (can be string or object).
  const bin = typeof pkg.bin === "string" ? { [pkg.name]: pkg.bin } : pkg.bin || {};

  // Check for binary at root path.
  const bins = await Promise.all(Object.keys(bin)
    // Check if normal and/or windows binaries exist.
    .map((binName) => Promise.all([
      binOrNull(rootPath, binName),
      binOrNull(rootPath, `${binName}.cmd`)
    ])
      // Filter out missing binary paths.
      .then((arr) => arr.filter((binPath) => !!binPath))
    )
  );

  // Start with arrays of arrays for recursion and bins.
  return []
    .concat(bins)
    .concat(locs)
    // Flatten sub-arrays.
    .reduce((m, a) => m.concat(a), [])
    // Remove dependencies we didn't find (`null`)
    .filter((item) => !!item)
    // By sorting, we can filter duplicates just looking one behind.
    .sort()
    .filter((item, i, items) => item !== items[i - 1]);
};

/**
 * Find on-disk locations of all production dependencies in `node_modules`.
 *
 * **Note**: This includes all `dependencies` and `optionalDependencies`.
 *
 * @param {*}       opts          options object
 * @param {string}  opts.rootPath `node_modules` root location (default: `process.cwd()`)
 * @param {string}  opts.curPath  starting inference location (default: `rootPath`)
 * @returns {Promise<Array<String>>} list of relative paths to on-disk dependencies.
 */
const findProdInstalls = async ({ rootPath, curPath } = {}) => {
  rootPath = path.resolve(rootPath || ".");
  curPath = curPath || rootPath;

  let pkg;
  try {
    pkg = await readJson(path.resolve(curPath, "package.json"));
  } catch (err) {
    // Enhance the error for consumers.
    if (err.code === "ENOENT") {
      throw new Error(`Unable to find package.json at: ${curPath}`);
    }
    throw err;
  }

  return resolveLocations({ pkg, rootPath, curPath });
};

module.exports = {
  findProdInstalls
};
