

docker-compose -f docker-compose-e2e-test.yml build --no-cache
docker-compose -f docker-compose-e2e-test.yml up -d
cd ../protocol
make e2e-setup
cd ../e2e-testing
chmod +x port-forwarder.sh
bash -c "./port-forwarder.sh"