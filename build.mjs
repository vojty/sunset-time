import { execa } from "execa";
import fse from "fs-extra";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Config
const appName = "SunsetTime";
const binaryName = "sunset-time";
const outDir = path.join(__dirname, "dist");

const targets = [
  "windows/amd64",
  "windows/386",
  "windows/arm64",
  "darwin/amd64",
  "darwin/arm64",
];

// https://stackoverflow.com/a/70483079
// -s Omit the symbol table and debug information.
// -w Omit the DWARF symbol table
const ldFlags = `-s -w`;

function getFlags(os) {
  if (os === "windows") {
    // hide terminal window for windows
    return `${ldFlags} -H=windowsgui`;
  }
  return ldFlags;
}

function makeMacApp(dir) {
  /**
   * Structure:
   *
   *  <name>.app
   *    Contents
   *      Info.plist
   *      MacOS
   *        <binaryName>
   *      Resouces
   *        icon.icns
   */
  const appDir = path.join(dir, `${appName}.app`, "Contents");
  fse.ensureDirSync(appDir);
  fse.copyFileSync(path.join(__dirname, "Info.plist"), `${appDir}/Info.plist`);

  const resourcesDir = path.join(appDir, "Resources");
  fse.ensureDirSync(resourcesDir);
  fse.copyFileSync(
    path.join(__dirname, "assets", "icon.icns"),
    path.join(resourcesDir, "icon.icns")
  );

  const binaryDir = path.join(appDir, "MacOS");
  fse.ensureDirSync(binaryDir);
  fse.moveSync(path.join(dir, binaryName), path.join(binaryDir, binaryName));
}

// Prepare directory
if (fse.existsSync(outDir)) {
  fse.removeSync(outDir);
}
fse.ensureDirSync(outDir);

targets.forEach(async (target) => {
  const [os, arch] = target.split("/");

  const dir = `${outDir}/${appName}-${os}-${arch}`;
  fse.ensureDirSync(dir);

  const suffix = os === "windows" ? ".exe" : "";
  const dist = `${dir}/${binaryName}${suffix}`;
  const args = ["build", "-ldflags", getFlags(os), "-o", dist, "."];

  const { stderr } = await execa("go", args, {
    env: {
      CGO_ENABLED: 1,
      GOOS: os,
      GOARCH: arch,
    },
  });
  if (stderr) {
    console.error(`${target} ❌`);
    console.error(stderr);
    process.exit(1);
  }

  if (os === "darwin") {
    makeMacApp(dir);
  }

  console.log(`${target} ✔️`);
});
