cd ../protocol
make e2etest-build-image
cd ../e2e-testing
docker-compose -f docker-compose-e2e-test.yml build --no-cache
docker-compose -f docker-compose-e2e-test.yml up
