docker run --rm --network=host confluentinc/cp-kafka:7.5.0 \
  kafka-topics \
  --create \
  --topic my-topic \
  --bootstrap-server localhost:29092 \
  --partitions 1 \
  --replication-factor 1
