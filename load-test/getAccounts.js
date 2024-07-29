const { createClient } = require('@clickhouse/client');

// Create a new ClickHouse instance
const client = createClient({
  url: 'http://localhost:8123', // Replace with your ClickHouse server URL
  username: 'default',           // Replace with your ClickHouse username
  password: 'P@ssw0rd',                  // Replace with your ClickHouse password
  database: 'default',           // Replace with your ClickHouse database name
});

async function getAccountMappings() {
  try {
    // Define the SQL query
    const query = 'SELECT * FROM ProcessManager.AccountMapping';

    // Execute the query
    const result = await client.query({query:query, format: 'JSONEachRow'} );

    // Log the results
   return await result.json();
  } catch (error) {
    console.error('Error fetching data from ClickHouse:', error);
  }
}

module.exports = {
  getAccountMappings,
};
// Call the function to get data
// getAccountMappings();
