cd ../protocol
make e2etest-build-image
cd ../e2e-testing
<<<<<<< clob-updates
#when building w/o  a cache you must split the command into two
=======
>>>>>>> v4-stable
docker-compose -f docker-compose-e2e-test.yml build --no-cache
docker-compose -f docker-compose-e2e-test.yml up
