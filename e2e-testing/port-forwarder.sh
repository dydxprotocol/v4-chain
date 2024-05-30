
docker exec -it interchain-security-instance bash -c "yum install -y socat && socat TCP-LISTEN:26658,fork TCP:7.7.8.253:26658"