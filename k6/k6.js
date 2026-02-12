import grpc from "k6/net/grpc";
import { check, sleep } from "k6";

const client = new grpc.Client();
client.load(["./proto"], "message.proto");

export const options = {
  vus: 5,
  duration: "30s",
};

export default () => {
  try {
    client.connect("localhost:50051", { plaintext: true });

    const res = client.invoke("message.MessageService/AddNumber", {
      number_1: Math.floor(Math.random() * 100),
      number_2: Math.floor(Math.random() * 100),
    });

    check(res, {
      "gRPC OK": (r) => r && r.status === grpc.StatusOK,
    });
  } catch (err) {
    console.error("k6 gRPC error:", err);
  } finally {
    client.close();
    sleep(0.1);
  }
};
