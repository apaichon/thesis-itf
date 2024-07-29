const autocannon = require('autocannon');
const {faker} = require('@faker-js/faker');
var reporter = require('autocannon-reporter')
const reportindex = 2
var reportOutputPaths = []  
const loadData = require('./data/loaddata.json')

// API endpoint
const apiEndpoints = [] 

for (var i = 1; i <= 3; i++) {
 reportOutputPaths.push(`./reports/Integration-Transition-${reportindex}-${i}.html`)
 apiEndpoints.push(`http://127.0.0.1:${8080 +i}/api/submit`)
}

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

   return {
        transactionDate: faker.date.recent().toISOString(),
        amount:  parseFloat(faker.finance.amount ({min:100, max:1000, dec:2})), // Amount between 100 and 1000 with 2 decimal places
        senderAccountId: senderAccount.AccountId,
        "receiverAccountId":receiverAccount.AccountId,
        "sourceInstitution": "Bank A",
        destinationInstitution: "Bank A",
        actBy: senderAccount.AccountId,
        actAt: faker.date.recent().toISOString(),
        description: faker.lorem.words(2),
        createdBy: senderAccount.AccountId
    }
  }


  async function getOptions () {
    let opts = []
    apiEndpoints.forEach((api, index) => {
      opts.push (
        {
          url: api,
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
          title: `transition-framework-report${reportindex}-${index+1}`,
          duration: 300,
          connections:loadData[index-1].connections,
          amount: loadData[index-1].amount
        
        }
      )
    })
  
    return opts
  }


async function start() {
  accountMappings = await getAccountMappings();
  const opts = await getOptions();
  // console.log('accountmapping',accountMappings)
    // Run load test using autocannon
    opts.forEach( (o,i) => {
      autocannon(o, (err, result) => {
        if (err) {
          console.error(err);
          return;
        }
        console.log(result);
        let html = reporter.buildReport(result) // the html structure
          reporter.writeReport(html, reportOutputPaths[i], (err, res) => {
            if (err) console.err('Error writting report: ', err)
            else console.log('Report written to: ', reportOutputPaths[i])
          }) //write the report
          
      });
  });
}

start();



