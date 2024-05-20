# RabbitMQ

* rabbitmqctl add_user new_user secret
* rabbitmqctl set_user_tags new_user administrator
* rabbitmqctl delete_user guest

* rabbitmqctl add_vhost customers
* rabbitmqctl set_permissions -p customers new_user ".*" ".*" ".*"
* rabbitmqctl set_topic_permissions -p customers new_user customer_events "^customers.*" "^customers.*"


* rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic -u new_user -p secret durable=true
* rabbitmqadmin delete exchange name=customer_events --vhost=customers -u new_user -p secret

* rabbitmqctl list_exchanges

* rabbitmqctl list_queues name messages_ready messages_unacknowledged

* rabbitmqctl list_bindings