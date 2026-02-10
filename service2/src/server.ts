import express from "express";
import type { Request, Response } from "express";
import { startConsumer, stopConsumer } from "./kafka/consumer.js";
import { retrieveLastValueInTheFile } from "./service/service.js";

const app = express();

app.use(express.json());
app.get("/number", async (req: Request, res: Response) => {
  try {
    const result = await retrieveLastValueInTheFile();
    res.status(200).json({ result });
  } catch (error) {
    res.status(500);
  }
});
const server = app.listen(3000, async () => {
  console.log("service 2 server is running on port 80");
  await startConsumer();
});

process.on("SIGINT", async () => {
  console.log("SIGINT recieved");
  await stopConsumer();
  server.close(() => process.exit(0));
});
