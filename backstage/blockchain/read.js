const path = require("path");
const Web3 = require("web3");
const fs = require("fs-extra");
// use the existing Member1 account address or make a new account
const address = "0x2bea587899fd6ca70e8b556d9e434ce0541f9168";
const contractAddress = "0x22Be720AeAaD3C81d5003b277e1843c31f7E8Ee9";
// read in the contracts
const contractJsonPath = path.resolve(__dirname, "Storage.json");
const contractJson = JSON.parse(fs.readFileSync(contractJsonPath));
const contractAbi = contractJson.abi;

async function getValue(
  host,
  key,
  fAddress,
  deployedContractAbi,
  deployedContractAddress,
) {
  const web3 = new Web3(host);
  const contractInstance = new web3.eth.Contract(
    deployedContractAbi,
    deployedContractAddress,
  );
  const res = await contractInstance.methods.getValue(key).call();
  return res;
}

async function main(key) {
	var res = await getValue("http://60.204.205.56:22000",key,address,contractAbi,contractAddress);
	console.log(res);
}
const args = process.argv.slice(2);
main(args[0]);

