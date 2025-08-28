include .env

kafka-produce:
	cat $(FILE) | docker exec -i kafka \
		kafka-console-producer --broker-list $(KAFKA_BROKER) --topic $(KAFKA_TOPIC)