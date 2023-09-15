const path = require("path");
const Web3 = require("web3");
const fs = require("fs-extra");
// use the existing Member1 account address or make a new account
const address = "0x2bea587899fd6ca70e8b556d9e434ce0541f9168";
// read in the contracts
const contractAddress = "0x22Be720AeAaD3C81d5003b277e1843c31f7E8Ee9";
const contractJsonPath = path.resolve(__dirname, "Storage.json");
const contractJson = JSON.parse(fs.readFileSync(contractJsonPath));
const contractAbi = contractJson.abi;
const contractByteCode = contractJson.bytecode;

// You need to use the accountAddress details provided to GoQuorum to send/interact with contracts
async function setValue(
  host,
  accountAddress,
  key,
  value,
  deployedContractAbi,
  deployedContractAddress,
) {
  const web3 = new Web3(host);
  const contractInstance = new web3.eth.Contract(
    deployedContractAbi,
    deployedContractAddress,
  );
  const res = await contractInstance.methods
    .setValue(key,value)
    .send({ from: accountAddress,gasPrice: "0x0",  gasLimit: "0x24A22" });
  return res;
}


async function main(key, value) {
	var res = await setValue("http://60.204.205.56:22000",address,key,value,contractAbi,contractAddress);
	//var res = await getValue(new Web3.providers.HttpProvider("http://localhost:22000"),456,contractAbi,contractAddress);
	console.log(res.status);
	return res.status;
}

const args = process.argv.slice(2)
var res = main(args[0], args[1]);

