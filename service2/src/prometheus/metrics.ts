import client from "prom-client";
import { Counter } from "prom-client";
client.collectDefaultMetrics();

export const cumlativeSumGauge = new Counter({
  name: "cumulative_sum",
  help: "Current cumulative sum of all processed numbers",
});

export const kafkaMessageCounter = new Counter({
  name: "kafka_messages_processed_total",
  help: "Total number of Kafka messages processed",
});

export const register = client.register;
