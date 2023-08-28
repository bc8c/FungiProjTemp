const fs = require("fs");
const path = require("path");

const FabricCAServices = require("fabric-ca-client");
const { Wallets } = require("fabric-network");

async function makeAdminWallet(org){
    const id = "admin"
    const pw = "adminpw"

    console.log(id, pw, org);

    try {
        // connection.json 객체화
        const ccpPath = path.resolve(__dirname,"..", "ccp", "connection-"+org+".json");
        const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities["ca."+org+".example.com"];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get(org);
        if (identity) {
            console.log(`An identity for the admin user ${org} already exists in the wallet`);
            return;
        }
        var mspid
        if (org == "org1"){
            mspid = "Org1MSP"
        }else
            mspid = "Org2MSP"

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: id, enrollmentSecret: pw });
        const x509Identity = {
        credentials: {
            certificate: enrollment.certificate,
            privateKey: enrollment.key.toBytes(),
        },
        mspId: mspid,
        type: "X.509",
        };
        await wallet.put(org, x509Identity);
        // response to client
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
    } catch (error) {
        console.error(`Failed to enroll admin user ${org} : ${error}`);        
    }
}
async function makeUsesrWallet(id, org){
    const userrole = "client"

    console.log(id, org);

    try {

        // connection.json 객체화
        const ccpPath = path.resolve(__dirname,"..", "ccp", "connection-"+org+".json");
        const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities["ca."+org+".example.com"];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(id);
        if (userIdentity) {
            console.log(`An identity for the user "${id}" already exists in the wallet`);
            return;
        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get(org);
        if (!adminIdentity) {
            console.log(`An identity for the admin user "${org}" does not exist in the wallet`);
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, org);

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register(
        {
            affiliation: org+".department1",
            enrollmentID: id,
            role: userrole,
        },
        adminUser
        );
            const enrollment = await ca.enroll({
            enrollmentID: id,
            enrollmentSecret: secret,
        });

        var mspid
        if (org == "org1"){
            mspid = "Org1MSP"
        }else
            mspid = "Org2MSP"

        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: mspid,
            type: "X.509",
        };
        await wallet.put(id, x509Identity);

        // response to client
        console.log(`Successfully registered and enrolled admin user "${id}" and imported it into the wallet`);
    } catch (error) {
        console.error(`Failed to enroll admin user ${id} : ${error}`);
    }
}

module.exports.makeAdminWallet = makeAdminWallet;
module.exports.makeUsesrWallet = makeUsesrWallet;