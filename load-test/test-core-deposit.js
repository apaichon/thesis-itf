const autocannon = require('autocannon');
const {faker} = require('@faker-js/faker');
var reporter = require('autocannon-reporter')
const index = 2
var reportOutputPath =  `./reports/deposit-core-report${index}.html`
const loadData = require('./data/loaddata.json')
// API endpoint
const apiEndpoint = 'http://localhost:5055/api/deposits';

// const { AccountID, Amount, DepositDate, CreatedBy}
// Generate random deposit data
function generateDepositData() {
  const AccountID = faker.number.int({ min: 1, max: 2 });
  const Amount = faker.finance.amount({ min: 100, max: 1000000 });
  const DepositDate = faker.date.betweens({ from: '2024-01-01T00:00:00.000Z', to: '2024-05-01T00:00:00.000Z' })[0]
  const CreatedBy = AccountID;
  return { AccountID, Amount,CreatedBy,DepositDate };
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
      // Generate random user data for each request
      const depositData = generateDepositData();
      // console.log('deposit', depositData)
      client.setBody(JSON.stringify(depositData));
    });
  },
  // Set a title for the load test
  title: `Core Bank Deposit Load Test-${index}`,
  duration: 300,
  connections:loadData[index].connections,
  amount: loadData[index].amount
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
