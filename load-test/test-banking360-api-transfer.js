const autocannon = require('autocannon');
const {faker} = require('@faker-js/faker');
var reporter = require('autocannon-reporter')
const index = 22
var reportOutputPath =  `./reports/banking360-transfer-report${index}.html`
const loadData = require('./data/loaddata.json')
// API endpoint
const apiEndpoint = 'http://localhost:4500/transfer';

// Generate random transfer data
function generateTransferData() {
  const senderAccountId = faker.string.uuid();
  const receiverAccountId = faker.string.uuid();
  const actBy = faker.string.uuid();
  const createdBy = faker.string.uuid();
  const amount = parseFloat (faker.finance.amount({min:100, max:10000, dec:2}));
  const transactionDate = faker.date.betweens({ from: '2024-01-01T00:00:00.000Z', to: '2024-08-01T00:00:00.000Z' })[0]
  const description = `Transfer Money Intra Bank #${faker.number.int({ min: 1000, max: 9999 })}`;
  
  return {
    transactionDate,
    amount,
    senderAccountId,
    receiverAccountId,
    actBy,
    createdBy,
    description,
  };
}

{
  "transaction_id": "550e8400-e29b-41d4-a716-446655440000",
  "transaction_type": "transfer_intra",
  "transaction_date": "2024-07-27T14:30:00Z",
  "amount": 1000.00,
  "currency": "THB",
  "account_id": "a1b2c3d4-e5f6-4321-8765-abcdef123456",
  "recipient_account_id": "98765432-dcba-4321-abcd-ef1234567890",
  "source_institution": "BankA",
  "destination_institution": "BankA",
  "status": "completed",
  "description": "Transfer from Mr. A to Mr. B",
  "act_by": "11112222-3333-4444-5555-666677778888",
  "act_at": "2024-07-27T14:30:05Z",
  "created_at": "2024-07-27T14:30:00Z",
  "created_by": "11112222-3333-4444-5555-666677778888",
  "sign": 1
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
      const request = generateTransferData();
      client.setBody(JSON.stringify(request));
    });
  },
  // Set a title for the load test
  title: `Banking360-Transfer-${index}`,
  duration: 300,
  connections: loadData[index-1].connections,
  amount: loadData[index-1].amount
};

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
