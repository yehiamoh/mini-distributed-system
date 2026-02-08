import fs from "fs/promises";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const filePath = path.join(__dirname, "../../result.txt");

export async function sumAndSaveResult(value: number): Promise<void> {
  const incomingValue = value;
  let currentValue = 0;

  try {
    const fileContent = await fs.readFile(filePath, "utf-8");
    const parsed = Number(fileContent);

    if (!Number.isNaN(parsed)) {
      currentValue = parsed;
    }
  } catch (err: any) {
    if (err.code !== "EOENT") {
      throw err;
    }
  }

  const newValue = currentValue + incomingValue;
  await fs.writeFile(filePath, String(newValue));
}

export async function retrieveLastValueInTheFile(): Promise<string> {
  try {
    return await fs.readFile(filePath, "utf-8");
  } catch (err: any) {
    if (err.code === "ENOENT") {
      return "0";
    }
    throw err;
  }
}
