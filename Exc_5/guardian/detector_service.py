"""The anomaly detector as a Kafka service. YOUR CODE (Part 2, Stage C2).

The plumbing is here, the decisions (marked TODO) are yours. Contract with
the maintenance team: every reading goes through model.score(); every
reading with score >= model.threshold becomes ONE anomaly event,
model.anomaly_event(reading, score), as JSON, published to the topic
`anomalies`.

    python3 detector_service.py --name det-1
    python3 detector_service.py --name det-2          # one more instance

My decisions:
- in-topic "readings": the topic the gateways write to (key = machine_id).
- group "line-guardian": all instances share ONE group, so Kafka gives every
  partition (and so every machine) to exactly one instance.
- "latest": a live detector should start with new readings, not with the
  history of old runs. It only applies when the group has no committed offset yet.
- anomalies keyed by machine_id: all anomalies of one machine stay in order
  for the maintenance team.
- when a partition is taken away (rebalance), the model forgets the machines of
  that partition. If the partition comes back later, the model starts fresh
  instead of using old state and counting the gap as "missing".
"""
import argparse
import json
import time

from confluent_kafka import Consumer, Producer

from model import BearingModel

ap = argparse.ArgumentParser()
ap.add_argument("--name", default="det-1", help="instance name (shows up in the scoreboard)")
ap.add_argument("--in-topic", default="readings", help="the topic the gateways write to")
ap.add_argument("--group", default="line-guardian", help="the consumer group of this detector")
ap.add_argument("--out-topic", default="anomalies")
ap.add_argument("--bootstrap", default="localhost:9092")
args = ap.parse_args()

model = BearingModel(instance=args.name)
consumer = Consumer({
    "bootstrap.servers": args.bootstrap,
    "group.id": args.group,
    "client.id": args.name,
    "auto.offset.reset": "latest",     # live detector: start with new readings, not old runs
})
producer = Producer({
    "bootstrap.servers": args.bootstrap,
    "acks": "all",
    "enable.idempotence": True,
})

# which machines came from which partition, so we can forget them on a rebalance
partition_machines = {}


def on_assign(consumer, partitions):
    print(f"[{args.name}] got partitions {sorted(p.partition for p in partitions)}", flush=True)


def on_revoke(consumer, partitions):
    for p in partitions:
        for machine_id in partition_machines.pop(p.partition, set()):
            model.forget(machine_id)
    print(f"[{args.name}] lost partitions {sorted(p.partition for p in partitions)}", flush=True)


consumer.subscribe([args.in_topic], on_assign=on_assign, on_revoke=on_revoke)
print(f"[{args.name}] {args.in_topic} -> {args.out_topic} (Ctrl+C to stop)", flush=True)

readings = anomalies = 0
report_at = time.time() + 5
try:
    while True:
        msg = consumer.poll(1.0)
        if msg is None:
            continue
        if msg.error():
            print(f"[{args.name}] {msg.error()}")
            continue
        reading = json.loads(msg.value())
        partition_machines.setdefault(msg.partition(), set()).add(reading["machine_id"])
        score = model.score(reading)
        readings += 1
        if score is not None and score >= model.threshold:
            anomalies += 1
            event = model.anomaly_event(reading, score)
            producer.produce(args.out_topic, key=reading["machine_id"],
                             value=json.dumps(event))
        producer.poll(0)
        if time.time() >= report_at:
            print(f"[{args.name}] {readings:>9,} readings  {anomalies:>6,} anomalies  "
                  f"missing {model.missing}  out-of-order {model.out_of_order}", flush=True)
            report_at += 5
except KeyboardInterrupt:
    pass
finally:
    producer.flush(10)
    consumer.close()
    print(f"[{args.name}] stopped after {readings:,} readings, {anomalies:,} anomalies")
