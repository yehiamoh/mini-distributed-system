import client from "prom-client";
import { Counter, Histogram } from "prom-client";
client.collectDefaultMetrics();

export const cumulativeSumCounter = new Counter({
  name: "cumulative_sum",
  help: "Current cumulative sum of all processed numbers",
});

export const kafkaMessageCounter = new Counter({
  name: "kafka_messages_processed_total",
  help: "Total number of Kafka messages processed",
});

export const kafkaProcessingDuration = new Histogram({
  name: "kafka_message_processing_duration_seconds",
  help: "Time from message creation until processing completes",
  labelNames: ["topic", "consumer_group"],
  buckets: [0.05, 0.1, 0.3, 0.5, 1, 2, 5, 10],
});
export const register = client.register;
