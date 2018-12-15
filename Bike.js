/*
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
*/

'use strict';
const shim = require('fabric-shim');
const util = require('util');

let Chaincode = class {

  // The Init method is called when the Smart Contract 'fabBike' is instantiated by the blockchain network
  // Best practice is to have any Ledger initialization in separate function -- see initLedger()
  async Init(stub) {
    console.info('=========== Instantiated fabBike chaincode ===========');
    return shim.success();
  }

  // The Invoke method is called as a result of an application request to run the Smart Contract
  // 'fabBike'. The calling application program has also specified the particular smart contract
  // function to be called, with arguments
  async Invoke(stub) {
    let ret = stub.getFunctionAndParameters();
    console.info(ret);

    let method = this[ret.fcn];
    if (!method) {
      console.error('no function of name:' + ret.fcn + ' found');
      throw new Error('Received unknown function ' + ret.fcn + ' invocation');
    }
    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  async queryBike(stub, args) {
    if (args.length != 1) {
      throw new Error('Incorrect number of arguments. Expecting BikeNumber ex: Bike01');
    }
    let BikeNumber = args[0];

    let BikeAsBytes = await stub.getState(BikeNumber); //get the Bike from chaincode state
    if (!BikeAsBytes || BikeAsBytes.toString().length <= 0) {
      throw new Error(BikeNumber + ' does not exist: ');
    }
    console.log(BikeAsBytes.toString());
    return BikeAsBytes;
  }

  async initLedger(stub, args) {
    console.info('============= START : Initialize Ledger ===========');
    let Bikes = [];
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
     make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });
    Bikes.push({
      make: 'Honda',
      model: 'Activa',
      color: 'gold',
      vin: '567SDA12',
      EngineCC:'100cc',
      owner: 'Sandhya'
    });

    for (let i = 0; i < Bikes.length; i++) {
      Bikes[i].docType = 'Bike';
      await stub.putState('Bike' + i, Buffer.from(JSON.stringify(Bikes[i])));
      console.info('Added <--> ', Bikes[i]);
    }
    console.info('============= END : Initialize Ledger ===========');
  }

  async createBike(stub, args) {
    console.info('============= START : Create Bike ===========');
    if (args.length != 5) {
      throw new Error('Incorrect number of arguments. Expecting 5');
    }

    var Bike = {
      docType: 'Bike',
      make: args[1],
      model: args[2],
      color: args[3],
      vin: args[4],
      EngineCC: args[5],
      owner: args[6]
    };

    await stub.putState(args[0], Buffer.from(JSON.stringify(Bike)));
    console.info('============= END : Create Bike ===========');
  }

  async queryAllBikes(stub, args) {

    let startKey = 'Bike0';
    let endKey = 'Bike999';

    let iterator = await stub.getStateByRange(startKey, endKey);

    let allResults = [];
    while (true) {
      let res = await iterator.next();

      if (res.value && res.value.value.toString()) {
        let jsonRes = {};
        console.log(res.value.value.toString('utf8'));

        jsonRes.Key = res.value.key;
        try {
          jsonRes.Record = JSON.parse(res.value.value.toString('utf8'));
        } catch (err) {
          console.log(err);
          jsonRes.Record = res.value.value.toString('utf8');
        }
        allResults.push(jsonRes);
      }
      if (res.done) {
        console.log('end of data');
        await iterator.close();
        console.info(allResults);
        return Buffer.from(JSON.stringify(allResults));
      }
    }
  }

  async changeBikeOwner(stub, args) {
    console.info('============= START : changeBikeOwner ===========');
    if (args.length != 2) {
      throw new Error('Incorrect number of arguments. Expecting 2');
    }

    let BikeAsBytes = await stub.getState(args[0]);
    let Bike = JSON.parse(BikeAsBytes);
    Bike.owner = args[1];

    await stub.putState(args[0], Buffer.from(JSON.stringify(Bike)));
    console.info('============= END : changeBikeOwner ===========');
  }
};

shim.start(new Chaincode());
