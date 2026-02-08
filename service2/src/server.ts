import express from "express";
import type { Request, Response } from "express";
import { startConsumer, stopConsumer } from "./kafka/consumer.js";

const app = express();

app.use(express.json());
app.use("/", (req: Request, res: Response) => {
  res.json({ helth: "true" });
});
const server = app.listen(80, async () => {
  console.log("service 2 server is running on port 80");
  await startConsumer();
});

process.on("SIGINT", async () => {
  console.log("SIGINT recieved");
  await stopConsumer();
  server.close(() => process.exit(0));
});
