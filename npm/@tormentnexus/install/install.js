#!/usr/bin/env node
const { execSync } = require("child_process");
const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");

console.log("╔══════════════════════════════════════════╗");
console.log("║   HyperNexus Universal Installer        ║");
console.log("║   38 AI Clients • One Command           ║");
console.log("╚══════════════════════════════════════════╝\n");

const platform = os.platform();
const arch = os.arch();

// Determine download URL based on platform
function getDownloadUrl() {
  const baseUrl = "https://releases.hypernexus.site/v1.0.0";

  if (platform === "win32") {
    return `${baseUrl}/hypernexus-setup.exe`;
  } else if (platform === "darwin") {
    return arch === "arm64"
      ? `${baseUrl}/hypernexus-darwin-arm64`
      : `${baseUrl}/hypernexus-darwin-amd64`;
  } else {
    return arch === "arm64"
      ? `${baseUrl}/hypernexus-linux-arm64`
      : `${baseUrl}/hypernexus-linux-amd64`;
  }
}

// Download file
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https
      .get(url, (response) => {
        // Handle redirects
        if (
          response.statusCode >= 300 &&
          response.statusCode < 400 &&
          response.headers.location
        ) {
          file.close();
          fs.unlinkSync(dest);
          download(response.headers.location, dest).then(resolve).catch(reject);
          return;
        }

        if (response.statusCode !== 200) {
          file.close();
          fs.unlinkSync(dest);
          reject(new Error(`Download failed: ${response.statusCode}`));
          return;
        }

        response.pipe(file);
        file.on("finish", () => {
          file.close();
          resolve();
        });
      })
      .on("error", (err) => {
        file.close();
        fs.unlinkSync(dest);
        reject(err);
      });
  });
}

async function main() {
  const url = getDownloadUrl();
  const isWindows = platform === "win32";
  const ext = isWindows ? ".exe" : "";
  const dest = path.join(os.tmpdir(), `hypernexus-install${ext}`);

  console.log(`Platform: ${platform} (${arch})`);
  console.log(`Downloading: ${url}\n`);

  try {
    await download(url, dest);

    // Make executable on Unix
    if (!isWindows) {
      fs.chmodSync(dest, "755");
    }

    console.log("Download complete. Running installer...\n");

    // Run the installer
    if (isWindows) {
      // Run .exe installer (silent mode)
      execSync(`"${dest}" /S`, { stdio: "inherit" });
    } else {
      // Run binary installer
      execSync(`"${dest}"`, { stdio: "inherit" });
    }

    // Cleanup
    try {
      fs.unlinkSync(dest);
    } catch {}

    console.log("\n✅ HyperNexus installed successfully!");
    console.log("   Run 'hypernexus serve' to start the server.");
    console.log("   Run 'hypernexus --help' for more options.\n");
  } catch (err) {
    console.error("\n❌ Installation failed:", err.message);
    console.error("\nPlease download manually from:");
    console.error("  https://hypernexus.site/download\n");
    process.exit(1);
  }
}

main();
