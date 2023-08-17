#!/usr/bin/env sh

PASSWORD="rootpass"
LIMIT=60
NEXT_WAIT_TIME=0

export ARANGO_NO_AUTH=1
until echo '0' | arangoinspect --server.ask-jwt-secret false --server.endpoint tcp://localhost:8529 --server.username root --server.password ${PASSWORD} --quiet true | grep "dignostics collected" || [ $NEXT_WAIT_TIME -eq $LIMIT ]; do
   sleep 1
   echo "retry init arangodb: $(( ++NEXT_WAIT_TIME ))"
done

## database
echo 'var list = db._databases(); if (list.includes("testDb")) {db._dropDatabase("testDb");} db._createDatabase("testDb");' | arangosh --server.password ${PASSWORD} --server.endpoint tcp://localhost:8529
arangosh --server.password ${PASSWORD} --server.endpoint tcp://localhost:8529 --javascript.execute-string '
  const users = require("@arangodb/users");
  users.save("testUser", "123456", true);
  users.grantDatabase("testUser", "testDb", "rw");
'
## collections
arangosh --server.password ${PASSWORD} --server.endpoint tcp://localhost:8529 --server.database testDb --javascript.execute-string '
  db._drop("scheduler");
'
arangoimp --file /opt/tools/scheduler.structure.json --collection scheduler --create-collection true --server.database testDb --server.password ${PASSWORD} --server.endpoint tcp://localhost:8529
