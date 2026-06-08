#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const path = require("node:path");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows"
};

const archMap = {
  x64: "amd64",
  arm64: "arm64"
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(`Unsupported platform: ${process.platform} ${process.arch}`);
  process.exit(1);
}

const extension = platform === "windows" ? ".exe" : "";
const binaryName = `vole-clean-${platform}-${arch}${extension}`;
const binaryPath = path.join(__dirname, "..", "dist", binaryName);

const result = spawnSync(binaryPath, process.argv.slice(2), {
  stdio: "inherit"
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 0);
