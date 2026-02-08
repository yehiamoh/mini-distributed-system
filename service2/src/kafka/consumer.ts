import { Kafka } from "kafkajs";
import type { EachMessagePayload, Consumer } from "kafkajs";

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
      const value = message.value?.toString();

      if (!value) return;

      console.log(
        "Recieved:",
        value,
        " from partition:",
        partition,
        "topic: ",
        topic,
      );

      const payload = JSON.parse(value) as {
        id: string;
        content: string;
      };
      console.log("Recieved", payload);

      // Logic
    },
  });
}

export async function stopConsumer(): Promise<void> {
  if (consumer) await consumer.stop();
}
