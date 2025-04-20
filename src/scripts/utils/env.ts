export function requireEnv(varName: string) {
    console.log(process.env[varName]);
    const value = process.env[varName];
    if (!value) {
      console.error(`❌ Required environment variable ${varName} is not set.`);
      process.exit(1);
    }
    return value;
  }