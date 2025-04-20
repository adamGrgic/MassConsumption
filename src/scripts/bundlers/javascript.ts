import { readdir, writeFile, mkdir, unlink } from "fs/promises";
import { existsSync } from "fs";
import { join, extname, basename } from "path";
import { writeFileSync, readFileSync, watch} from 'fs';
import path from 'path';
import { argv } from 'process';
import type { Dirent } from "node:fs";
import { requireEnv } from "../utils/env";
import { getHash, deleteHashedFiles } from "../utils/hash";

const INPUT_DIR_JS = requireEnv("INPUT_DIR_JS");
const OUTPUT_DIR_JS = requireEnv("OUTPUT_DIR_JS");
const MANIFEST_PATH = requireEnv("MEDIA_MANIFEST_PATH")

async function ensureOutputDir() {
  if (!existsSync(OUTPUT_DIR_JS)) {
    await mkdir(OUTPUT_DIR_JS, { recursive: true })
  }
}

async function compileFile(entry: Dirent, fullPath: string): Promise<[string, string] | null> {
  const relativePath = fullPath.replace(INPUT_DIR_JS + "/", "");
  const logicalName = basename(entry.name, ".ts");

  const result = await Bun.build({
    entrypoints: [fullPath],
    minify: true,
    target: "bun",
  });

  if (!result.success || result.outputs.length === 0) {
    console.error(`❌ Failed to build: ${relativePath}`);
    return null;
  }

  const outputContent = await result!.outputs[0]!.text();
  const hash = getHash(outputContent);
  const hashedFilename = `${logicalName}.${hash}.js`;
  const outputPath = join(OUTPUT_DIR_JS!, hashedFilename);

  console.log("[VERBOSE DEBUG] JS Content: ", outputContent)
  console.log("[VERBOSE DEBUG] hash: ", hash)
  console.log("[VERBOSE DEBUG] hashedFilename: ", hashedFilename)
  console.log("[VERBOSE DEBUG] outputPath: ", outputPath)

  await deleteHashedFiles(INPUT_DIR_JS!, logicalName, "js", OUTPUT_DIR_JS);
  writeFileSync(outputPath, outputContent);

  console.log(`✅ Compiled ${relativePath} → ${hashedFilename}`);
  return [logicalName, outputPath];
}

async function walkAndCompile(dir: string): Promise<[string, string][]> {
  const entries = await readdir(dir, { withFileTypes: true });
  const compileJobs: Promise<[string, string] | null>[] = [];

  for (const entry of entries) {
    const fullPath = join(dir, entry.name);

    if (entry.isDirectory()) {
      const nestedResults = await walkAndCompile(fullPath); // await here
      for (const result of nestedResults) {
        compileJobs.push(Promise.resolve(result));
      }
    } else if (entry.isFile() && extname(entry.name) === ".ts") {
      compileJobs.push(compileFile(entry, fullPath));
    }
  }

  const results = await Promise.all(compileJobs);
  return results.filter((r): r is [string, string] => r !== null);
}
  
async function run() {
  console.log("🚀 Starting TS compiler...");

  await ensureOutputDir()

  const entries = await walkAndCompile(INPUT_DIR_JS!);

  let manifest: Record<string, string> = {};
  if (existsSync(MANIFEST_PATH)) {
    manifest = JSON.parse(readFileSync(MANIFEST_PATH, "utf8"));
  }

  for (const [logicalName, outputPath] of entries) {
    manifest[`js:${logicalName}`] = outputPath;
  }

  writeFileSync(MANIFEST_PATH, JSON.stringify(manifest, null, 2));
  console.log("✅ All done!");
}

if (argv.includes('--watch')) {
    // console.log(' watch flag given');

    run();
    // console.log(path.dirname(INPUT_DIR));
    watch(path.dirname(INPUT_DIR_JS!), { recursive: true }, (event, filename) => {
    if (filename && filename.endsWith('.ts')) {
      console.log(`🔄 Change detected in ${filename}, recompiling...`);
      run();
    }
  });

} else {
    // console.log('no watch flag given')
    run();
}