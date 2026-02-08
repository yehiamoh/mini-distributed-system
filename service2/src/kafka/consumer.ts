import { Kafka } from "kafkajs";
import type { EachMessagePayload, Consumer } from "kafkajs";
import { sumAndSaveResult } from "../service/service.js";

const kafka = new Kafka({
  clientId: "service2",
  brokers: ["localhost:9092"],
});

let consumer: Consumer;

export async function startConsumer(): Promise<void> {
  consumer = kafka.consumer({ groupId: "service2-group" });

  await consumer.connect();
  await consumer.subscribe({
    topic: "message.created",
    fromBeginning: false,
  });

  await consumer.run({
    eachMessage: async ({ topic, partition, message }: EachMessagePayload) => {
      try {
        const value = message.value?.toString();

        if (!value) return;

        const payload = JSON.parse(value) as {
          Num_1: number;
          Nums_2: number;
          Result: number;
          CreatedAt: string;
          RequestID: string;
        };
        console.log("Recieved", payload);

        await sumAndSaveResult(payload.Result);
      } catch (err) {
        console.error("Failed to process message", err);
      }
    },
  });
}

export async function stopConsumer(): Promise<void> {
  if (consumer) await consumer.stop();
}
