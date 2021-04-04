/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke coins ====
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin1","blue","35","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin2","red","50","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["initCoin","coin3","blue","70","tom"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["transferCoin","coin2","jerry"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["transferCoinsBasedOnColor","blue","jerry"]}'
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["delete","coin1"]}'

// ==== Query coins ====
// peer chaincode query -C myc1 -n coins -c '{"Args":["readCoin","coin1"]}'
// peer chaincode query -C myc1 -n coins -c '{"Args":["getCoinsByRange","coin1","coin3"]}'
// peer chaincode query -C myc1 -n coins -c '{"Args":["getHistoryForCoin","coin1"]}'

// Rich Query (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoinsByOwner","tom"]}'
//   peer chaincode query -C myc1 -n coins -c '{"Args":["queryCoins","{\"selector\":{\"owner\":\"tom\"}}"]}'

// INDEXES TO SUPPORT COUCHDB RICH QUERIES
//
// Indexes in CouchDB are required in order to make JSON queries efficient and are required for
// any JSON query with a sort. As of Hyperledger Fabric 1.1, indexes may be packaged alongside
// chaincode in a META-INF/statedb/couchdb/indexes directory. Each index must be defined in its own
// text file with extension *.json with the index definition formatted in JSON following the
// CouchDB index JSON syntax as documented at:
// http://docs.couchdb.org/en/2.1.1/api/database/find.html#db-index
//
// This coins02 example chaincode demonstrates a packaged
// index which you can find in META-INF/statedb/couchdb/indexes/indexOwner.json.
// For deployment of chaincode to production environments, it is recommended
// to define any indexes alongside chaincode so that the chaincode and supporting indexes
// are deployed automatically as a unit, once the chaincode has been installed on a peer and
// instantiated on a channel. See Hyperledger Fabric documentation for more details.
//
// If you have access to the your peer's CouchDB state database in a development environment,
// you may want to iteratively test various indexes in support of your chaincode queries.  You
// can use the CouchDB Fauxton interface or a command line curl utility to create and update
// indexes. Then once you finalize an index, include the index definition alongside your
// chaincode in the META-INF/statedb/couchdb/indexes directory, for packaging and deployment
// to managed environments.
//
// In the examples below you can find index definitions that support coins02
// chaincode queries, along with the syntax that you can use in development environments
// to create the indexes in the CouchDB Fauxton interface or a curl command line utility.
//
