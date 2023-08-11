#!/usr/bin/env sh

PASSWORD="rootpass"

export ARANGO_NO_AUTH=1

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
