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
// peer chaincode invoke -C myc1 -n coins -c '{"Args":["transferCoinsBasedO