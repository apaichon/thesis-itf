const {faker} = require('@faker-js/faker');

// Function to generate a single record
function generateRecord(index) {
    return {
        [`input${index}`]: {
            transactionDate: faker.date.recent().toISOString(),
            amount: faker.finance.amount ({min:100, max:1000, dec:2}), // Amount between 100 and 1000 with 2 decimal places
            accountId: faker.string.uuid(),
            destinationInstitution: faker.company.name(),
            actBy: faker.string.uuid(),
            actAt: faker.date.recent().toISOString(),
            description: faker.lorem.words(2),
            createdBy: faker.string.uuid()
        }
    };
}

// Function to generate multiple records
function generateRecords(numRecords) {
    const records = {};
    for (let i = 1; i <= numRecords; i++) {
        Object.assign(records, generateRecord(i));
    }
    return records;
}

function generateDepositMutation(numInputs) {
    let inputs = [];
    let deposits = [];

    for (let i = 1; i <= numInputs; i++) {
        inputs.push(`$input${i}:DepositInput`);
        deposits.push(`deposit${i}:deposit(input:$input${i}) {
            code
            status
            message
        }`);
    }

    return `
mutation Deposit(${inputs.join(', ')}) {
  banking360Mutations {
    ${deposits.join('\n    ')}
  }
}`;
}

export { generateRecords, generateDepositMutation};

// Specify the number of records to generate
// const numberOfRecords = 100;
// const generatedData = generateRecords(numberOfRecords);
// const mutationString = generateDepositMutation(numberOfRecords);

// console.log(mutationString);

// console.log(JSON.stringify(generatedData, null, 2));
