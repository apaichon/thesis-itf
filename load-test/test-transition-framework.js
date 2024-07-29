const autocannon = require('autocannon');
const {faker} = require('@faker-js/faker');
var reporter = require('autocannon-reporter')
const index = 3
var reportOutputPath =  `./reports/transition-framework-reportV2-${index}.html`
const loadData = require('./data/loaddata.json')
// API endpoint
const apiEndpoint = 'http://localhost:8081/api/submit';
const { v4: uuidv4 } = require('uuid');
const topics = ['TransferIntraBank']
const {getAccountMappings} = require('./getAccounts');

let accountMappings = {}
function createMessage() {
    return {
        id: uuidv4(),
        system: "Banking360",
        topic:  topics[0],
        content:  JSON.stringify (genTransferMessage()),
        remark: "TranferIntraBank",
        createdAt: faker.date.recent(),
        createdBy: faker.person.firstName().substring(0,50),
        sign: 1
    };
}

function genTransferMessage() {
    senderAccount = accountMappings[faker.number.int({min:0, max:9999})];
    receiverAccount = accountMappings[faker.number.int({min:0, max:9999})];
    const amount = parseFloat(faker.finance.amount ({min:100, max:1000, dec:2}))
   return {
        transactionDate: faker.date.recent().toISOString(),
        amount:  amount, // Amount between 100 and 1000 with 2 decimal places
        senderAccountId: senderAccount.AccountId,
        "receiverAccountId":receiverAccount.AccountId,
        "sourceInstitution": "Our Bank",
        destinationInstitution: "Our Bank",
        actBy: senderAccount.AccountId,
        actAt: faker.date.recent().toISOString(),
        description: `Transfer Money from ${senderAccount.AccountId} to ${receiverAccount.AccountId} amount is ${amount}.`,
        createdBy: senderAccount.AccountId
    }
  }

// Options for autocannon
const opts = {
  url: apiEndpoint,
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  setupClient: (client) => {
    client.on('request', (requestParams) => {
     
      const request = createMessage();
      client.setBody(JSON.stringify(request));
    });
  },
  // Set a title for the load test
  title: `Integration-Transition-${index}`,
  duration: 300,
  connections:loadData[index-1].connections,
  amount: loadData[index-1].amount
};

async function start() {
  accountMappings = await getAccountMappings();
  // Run load test using autocannon
  autocannon(opts, (err, result) => {
    if (err) {
      console.error(err);
      return;
    }
    console.log(result);
    let html = reporter.buildReport(result) // the html structure
      reporter.writeReport(html, reportOutputPath, (err, res) => {
        if (err) console.err('Error writting report: ', err)
        else console.log('Report written to: ', reportOutputPath)
      }) //write the report
      
  });
}

start();



