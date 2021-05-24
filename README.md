# map-live-stream
A Live Stream made using Kafka Go

## Commands for running Kafka via Docker and Docker Compose

### Docker run zookeeper
```
docker run -d \
--net=host \
--name=zookeeper \
-e ZOOKEEPER_CLIENT_PORT=32181 \
-e ZOOKEEPER_TICK_TIME=2000 \
-e ZOOKEEPER_SYNC_LIMIT=2 \
confluentinc/cp-zookeeper:6.1.1

```
### Docker run kafka
```
docker run -d \
    --net=host \
    --name=kafka \
    -e KAFKA_ZOOKEEPER_CONNECT=localhost:32181 \
    -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:29092 \
    -e KAFKA_BROKER_ID=2 \
    -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
    -e CONFLUENT_SUPPORT_CUSTOMER_ID=c0 \
    confluentinc/cp-server:6.1.1
```

### Run Docker Compose 

```
cd scripts/Kafka
docker-compose up
```

### Checking running processes

```
docker ps
docker-compose ps
```
### Removing Docker containers

```
sudo docker container rm --force [container names]
```

