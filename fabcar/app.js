
'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Chaincode Invoke
 */

var Fabric_Client = require('fabric-client');
var path = require('path');
var util = require('util');
var os = require('os');
let bodyParser = require('body-parser');
const express = require('express');
const { raw } = require('body-parser');
const { request } = require('http');
const app = express();
app.use(express.json());
//
var fabric_client = new Fabric_Client();

// setup the fabric network
var channel = fabric_client.newChannel('mychannel');
var peer = fabric_client.newPeer('grpc://localhost:7051');
channel.addPeer(peer);
var order = fabric_client.newOrderer('grpc://localhost:7050')
channel.addOrderer(order);

//
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:'+store_path);
var tx_id = null;

// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
Fabric_Client.newDefaultKeyValueStore({ path: store_path
}).then((state_store) => {
	// assign the store to the fabric client
	fabric_client.setStateStore(state_store);
	var crypto_suite = Fabric_Client.newCryptoSuite();
	// use the same location for the state store (where the users' certificate are kept)
	// and the crypto store (where the users' keys are kept)
	var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
	crypto_suite.setCryptoKeyStore(crypto_store);
	fabric_client.setCryptoSuite(crypto_suite);

	// get the enrolled user from persistence, this user will sign all requests
	return fabric_client.getUserContext('user1', true);
}).then((user_from_store) => {
	if (user_from_store && user_from_store.isEnrolled()) {
		console.log('Successfully loaded user1 from persistence');
		member_user = user_from_store;
	} else {
		throw new Error('Failed to get user1.... run registerUser.js');
	}
});

app.post('/storemetadata', (req, res) => {
    let user = req.body.user
    let metaData = req.body.metaData

    // get a transaction id object based on the current user assigned to fabric client
    tx_id = fabric_client.newTransactionID();
    console.log("Assigning transaction_id: ", tx_id._transaction_id);

    var request = {
		//targets: let default to the peer assigned to the client
		chaincodeId: 'fabcar',
		fcn: 'storeMetaData',
		args: [user, metaData],
		chainId: 'mychannel',
		txId: tx_id
	};

    console.log("request has been set");
    // send the transaction proposal to the peers
    ledgerAPI.slimInvoke(channel, request, peer).then(() => {
        res.redirect('/storemetadata');
    }).catch((err) => {
        console.error('Failed to invoke successfully :: ' + err);
    });
});

app.post('/storeTalList', (req, res) => {
    let entityId = req.body.entityId
    let tal = req.body.tal

    // get a transaction id object based on the current user assigned to fabric client
    tx_id = fabric_client.newTransactionID();
    console.log("Assigning transaction_id: ", tx_id._transaction_id);

    var request = {
		//targets: let default to the peer assigned to the client
		chaincodeId: 'fabcar',
		fcn: 'storeTalList',
		args: [entityId, tal],
		chainId: 'mychannel',
		txId: tx_id
	};

    console.log("request has been set");
    // send the transaction proposal to the peers
    ledgerAPI.slimInvoke(channel, request, peer).then(() => {
        res.redirect('/storeTalList');
    }).catch((err) => {
        console.error('Failed to invoke successfully :: ' + err);
    });
});
