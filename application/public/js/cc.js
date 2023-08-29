
const fs = require("fs");
const path = require("path");


const {Gateway, Wallets} =  require("fabric-network");
// load the network configuration
const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));


async function cc_call(id, fn_name, args) {

    try {
         // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get(id);
        if (!identity) {
            console.log(
                `An identity for the user "${id}" does not exist in the wallet`
            );             
            const res_str = `An identtity for the user(${id}) does not exitsts in th wallet`;
            throw res_str
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: id,
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        // 1. fungusfactory
        // 2. feedfactory
        // const contract = network.getContract("basic");

        var result;

        if (fn_name == "CreateRandomFungus") {
            const contract = network.getContract("fungusfactory");
            result = await contract.submitTransaction("CreateRandomFungus",args);
        } else if (fn_name == "GetFungiByOwner"){
            const contract = network.getContract("fungusfactory");
            result = await contract.evaluateTransaction("GetFungiByOwner");
        }else if (fn_name == "Feed"){
            const contract = network.getContract("fungusfactory");
            result = await contract.submitTransaction("Feed", args[0], args[1]);
        } else if (fn_name == "CreateRandomFeed") {
            const contract = network.getContract("feedfactory");
            result = await contract.submitTransaction("CreateRandomFeed",args);
        } else {
            result = "not supported function!!"
        } 
        // Disconnect from the gateway.
        await gateway.disconnect();
        
        return result;
        
    } catch (error) {
        console.error(error)
    }

   
}

module.exports.cc_call = cc_call;