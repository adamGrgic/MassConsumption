#!/usr/bin/env bun

import { readdir, mkdir, unlink } from "fs/promises";
import { readFileSync, writeFileSync, existsSync, watch } from "fs";
import { join, extname, basename } from "path";
import { argv } from "process";
import type { Dirent } from "node:fs";
import { compileAsync } from "sass";
import { requireEnv } from "../utils/env";
import { getHash, deleteHashedFiles } from "../utils/hash"

// Config
const INPUT_DIR = requireEnv("INPUT_DIR_CSS");
const OUTPUT_DIR = requireEnv("OUTPUT_DIR_CSS");
const MEDIA_MANIFEST_PATH = requireEnv("MEDIA_MANIFEST_PATH");

// Ensure output directory exists
async function ensureOutputDir() {
  if (!existsSync(OUTPUT_DIR)) {
    console.log("📁 Creating output directory...");
    await mkdir(OUTPUT_DIR, { recursive: true });
  }
}

// Compile a single SCSS file
async function compileScssFile(entry: Dirent, fullPath: string): Promise<[string, string] | null> {
  console.log(`🔍 Found file: ${entry.name}`);

  if (entry.name.startsWith("_") || extname(entry.name) !== ".scss") {
    console.log(`⏭️  Skipped partial or unsupported file: ${entry.name}`)
    return null;
  }

  const logicalName = basename(entry.name, ".scss")
  console.log("logical name: ", logicalName)

  console.log("full path to be compiled: ", fullPath)
  let result;
  try {
    result = await compileAsync(fullPath, {
      loadPaths: ["src/styles", "node_modules"],
      style: "compressed",
    });
  } catch (err) {
    console.error(`❌ Error compiling ${entry.name}:`, err);
    return null;
  }

  console.log('compile process completed');
  const cssText = result.css;
  const hash = getHash(cssText);
  const hashedFilename = `${logicalName}.${hash}.css`;
  const outputPath = join(OUTPUT_DIR, hashedFilename);

  console.log("[VERBOSE DEBUG] CSS Content: ", cssText)
  console.log("[VERBOSE DEBUG] hash: ", hash)
  console.log("[VERBOSE DEBUG] hashedFilename: ", hashedFilename)
  console.log("[VERBOSE DEBUG] outputPath: ", outputPath)

  await deleteHashedFiles(OUTPUT_DIR, logicalName, "css", hashedFilename);

  writeFileSync(outputPath, cssText);

  console.log(`✅ Compiled ${entry.name} → ${hashedFilename}`);
  return [logicalName, outputPath];
}

// Walk all directories and compile SCSS files
async function walkAndCompile(dir: string): Promise<[string, string][]> {
  console.log(`📁 Walking directory: ${dir}`);
  const entries = await readdir(dir, { withFileTypes: true });

  let results: [string, string][] = [];

  for (const entry of entries) {
    const fullPath = join(dir, entry.name);

    if (entry.isDirectory()) {
      const nested = await walkAndCompile(fullPath);
      results = results.concat(nested);
    } else {
      const result = await compileScssFile(entry, fullPath);
      if (result) results.push(result);
    }
  }

  return results;
}

// Main entry point
async function run() {
  console.clear();
  console.log("🎨 Starting SCSS compiler with Bun + Sass...\n");

  await ensureOutputDir();
  const compiledResults = await walkAndCompile(INPUT_DIR);

  if (compiledResults.length === 0) {
    console.warn("⚠️ No SCSS files were compiled.");
    return;
  }

  console.log("\n📘 Writing manifest...");
  let manifest: Record<string, string> = {};

  if (existsSync(MEDIA_MANIFEST_PATH)) {
    manifest = JSON.parse(readFileSync(MEDIA_MANIFEST_PATH, "utf-8"));
  }

  for (const [logicalName, outputPath] of compiledResults) {
    console.log(`📝 ${logicalName} → ${outputPath}`);
    manifest[`css:${logicalName}`] = outputPath;
  }

  writeFileSync(MEDIA_MANIFEST_PATH, JSON.stringify(manifest, null, 2));
  console.log("\n✅ Manifest written.");
  console.log("\n🏁 Done compiling all SCSS files.\n");
}

// Watch mode
if (argv.includes("--watch")) {
  console.log("🔁 Watch mode enabled.\n");
  run();
  watch(INPUT_DIR, { recursive: true }, (_, filename) => {
    if (filename && filename.endsWith(".scss")) {
      console.log(`\n🔄 Change detected in ${filename}, recompiling...\n`);
      run();
    }
  });
} else {
  run();
}
