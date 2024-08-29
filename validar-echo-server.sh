#!/bin/bash
TEST_MESSAGE=$(openssl rand -base64 12)

RESPONSE=$(sudo docker run --rm --network tp0_testing_net busybox sh -c "echo $TEST_MESSAGE | nc server 12345")

if [ "$RESPONSE" == "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi