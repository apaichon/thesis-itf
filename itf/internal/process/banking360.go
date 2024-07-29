package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"io"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/google/uuid"
)

type TransferIntraBankInput struct {
	TransactionDate   string    `json:"transactionDate"`
	Amount            float64   `json:"amount"`
	SenderAccountId   uuid.UUID `json:"senderAccountId"`
	ReceiverAccountId uuid.UUID `json:"receiverAccountId"`
	ActBy             uuid.UUID `json:"actBy"`
	CreatedBy         uuid.UUID `json:"createdBy"`
	Description       string    `json:"description"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

type Banking360 struct {
	Query      string
	Variables  map[string]interface{}
	ServiceUrl string
}

func NewBanking360() *Banking360 {
	return &Banking360{ServiceUrl: "http://localhost:4009/graphql"}
}

func (b *Banking360) TransformIntraBank(messages []models.MessageModel) (string, map[string]interface{}) {
	query, variables := b.convertToGraphQL(messages)
	return query, variables
}

func (b *Banking360) convertToGraphQL(messages []models.MessageModel) (string, map[string]interface{}) {
	query := "mutation TransferIntraBank("
	vars := make(map[string]interface{})
	inputs := make([]string, len(messages))
	mutations := make([]string, len(messages))

	for i, message := range messages {

		inputName := fmt.Sprintf("input%d", i+1)
		inputs[i] = fmt.Sprintf("$%s: TransferIntraBankInput!", inputName)
		mutations[i] = fmt.Sprintf("transfer%d: transferIntraBank(input: $%s) { code status message }", i+1, inputName)

		var input TransferIntraBankInput
		if err := json.Unmarshal([]byte(message.Content), &input); err != nil {
			return "", map[string]interface{}{"error": fmt.Errorf("error unmarshalling content for message %d: %w", i, err)}
		}
		vars[inputName] = input
	}

	query += strings.Join(inputs, ", ") + ") {\n  banking360Mutations {\n    "
	query += strings.Join(mutations, "\n    ")
	query += "\n  }\n}"

	return query, vars
}

func (b *Banking360) ConvertMessageToTransferInput(msg models.MessageModel) (TransferIntraBankInput, error) {
	var input TransferIntraBankInput
	err := json.Unmarshal([]byte(msg.Content), &input)
	if err != nil {
		return TransferIntraBankInput{}, err
	}

	return input, nil
}

func (banking360 *Banking360) PostToGraphQL(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	// Create the GraphQL request payload
	requestPayload := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// fmt.Printf("payload:%v",requestPayload)

	// Marshal the request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", banking360.ServiceUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var graphQLResponse GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphQLResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &graphQLResponse, nil
}

func (b *Banking360) PostTransferIntraBankInput(input TransferIntraBankInput) error {
	url := "http://localhost:4500/transfer"

	// Convert input to JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post data, status code: %d", resp.StatusCode)
	}

	return nil
}
