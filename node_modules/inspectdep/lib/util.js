"use strict";

const path = require("path");
const { promisify } = require("util");
const fs = require("fs");

const readFile = promisify(fs.readFile);
const readJson = async (filePath) => JSON.parse(await readFile(filePath));
const stat = promisify(fs.stat);
const exists = (filePath) => stat(filePath)
  .then(() => true)
  .catch((err) => {
    if (err.code === "ENOENT") { return false; }
    throw err;
  });

// Return location and data from package json on disk.
//
// Iterates from curPath up to rootPath.
const findPkg = async ({ rootPath, curPath, name }) => {
  let loc;
  let pkg;

  while (curPath.length >= rootPath.length) {
    loc = path.join(curPath, "node_modules", name);
    try {
      pkg = await readJson(path.join(loc, "package.json"));
      return { pkg, loc };
    } catch (err) {
      // Decrement path and try again if not found.
      if (err.code === "ENOENT") {
        curPath = path.dirname(curPath);
        continue;
      }

      // Otherwise, we have a real error.
      throw err;
    }
  }

  return null;
};

module.exports = {
  exists,
  readJson,
  findPkg
};
