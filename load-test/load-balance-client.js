// API endpoints
const apiEndpoints = [
    'http://127.0.0.1:8081/submit',
    'http://127.0.0.1:8082/submit',
    'http://127.0.0.1:8083/submit'
  ];
  
  // Function to post message with load balancing
  async function postMessageWithLoadBalancing(message) {
    let endpoint;
  
    // Determine the endpoint based on the topic
    switch (message.topic) {
      case 'Message A':
        endpoint = apiEndpoints[0]; // 8081
        break;
      case 'Message B':
        endpoint = apiEndpoints[1]; // 8082
        break;
      default:
        endpoint = apiEndpoints[2]; // 8083 for all other cases
    }
  
    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(message),
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      const result = await response.json();
      console.log(`Message posted successfully to ${endpoint}`, result);
      return result;
    } catch (error) {
      console.error(`Failed to post message to ${endpoint}`, error);
      throw error;
    }
  }
  
  // Example usage
  const message = {
    id: '123',
    system: 'ExampleSystem',
    topic: 'Message A',
    content: 'This is a test message',
    remark: 'Test remark',
    createdBy: 'User1',
    sign: 1
  };
  
  postMessageWithLoadBalancing(message)
    .then(result => console.log('Operation completed', result))
    .catch(error => console.error('Operation failed', error));