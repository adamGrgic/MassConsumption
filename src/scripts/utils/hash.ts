import { createHash } from "crypto";
import { existsSync } from "fs";
import { readdir, unlink } from "fs/promises";
import { join } from "path";

export function getHash(content: string): string {
  return createHash("md5").update(content).digest("hex").slice(0, 8);
}

export async function deleteHashedFiles(
  outputDir: string,
  logicalName: string,
  ext: string,
  skipFilename?: string 
) {
  if (!existsSync(outputDir)) return;

  const files = await readdir(outputDir);
  const pattern = new RegExp(`^${logicalName}\\.[a-f0-9]{8}\\.${ext}$`);

  for (const file of files) {
    if (pattern.test(file)) {
      if (skipFilename && file === skipFilename) {
        continue; // Don't delete the file we're about to write
      }

      await unlink(join(outputDir, file));
      console.log(`🗑️  Deleted old build: ${file}`);
    }
  }
}
