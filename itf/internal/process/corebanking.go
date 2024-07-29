package process

import (
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io"
	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
)

type Withdrawal struct {
	WithdrawalID   int       `json:"withdrawalID"`
	AccountID      int       `json:"accountID"`
	Amount         float64   `json:"amount"`
	WithdrawalDate time.Time `json:"withdrawalDate"`
	CreatedBy      int       `json:"createdBy"`
}

type Deposit struct {
	DepositID   int       `json:"depositID"`
	AccountID   int       `json:"accountID"`
	Amount      float64   `json:"amount"`
	DepositDate time.Time `json:"depositDate"`
	CreatedBy   int       `json:"createdBy"`
}

type CoreBanking struct {
	ProcessManagerRepo *repositories.ProcessManagerRepo
	ServiceUrl         string
	HTTPClient         *http.Client
}

func NewCoreBanking() *CoreBanking {
	return &CoreBanking{ServiceUrl: "http://localhost:5055/api",
		ProcessManagerRepo: repositories.NewProcessManagerRepo(),
	}
}

func (c *CoreBanking) TransformWithdrawal(message models.MessageModel) (*Withdrawal, error) {
	var input TransferIntraBankInput
	if err := json.Unmarshal([]byte(message.Content), &input); err != nil {
		return nil, fmt.Errorf("error unmarshalling content for message : %v", err)
	}

	// Convert TransactionDate to time.Time
	transactionDate, err := time.Parse(time.RFC3339, input.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction date: %v", err)
	}

	// Create Withdrawal and Deposit objects
	withdrawal := &Withdrawal{
		WithdrawalID: 0,
		Amount:         input.Amount,
		WithdrawalDate: transactionDate,
	}


	account, err := c.ProcessManagerRepo.GetAccountMapping(input.SenderAccountId)
	if err != nil {
		return nil, err
	}
	withdrawal.AccountID = account.AccountLegacyId
	withdrawal.CreatedBy = account.AccountLegacyId

	return withdrawal, nil
}

func (c *CoreBanking) TransformDeposit(message models.MessageModel) (*Deposit, error) {
	var input TransferIntraBankInput
	if err := json.Unmarshal([]byte(message.Content), &input); err != nil {
		return nil, fmt.Errorf("error unmarshalling content for message : %v", err)
	}

	// Convert TransactionDate to time.Time
	transactionDate, err := time.Parse(time.RFC3339, input.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction date: %v", err)
	}

	// Create Deposit objects
	deposit := &Deposit{
		DepositID:   0,
		Amount:      input.Amount,
		DepositDate: transactionDate,
	}

	// Get Account Mapping
	account, err := c.ProcessManagerRepo.GetAccountMapping(input.ReceiverAccountId)
	if err != nil {
		return nil, err
	}
	deposit.AccountID = account.AccountLegacyId
	deposit.CreatedBy = account.AccountLegacyId

	return deposit, nil
}

func (c *CoreBanking) PostWithdrawal(withdrawal Withdrawal) (string, error) {
	// Post Withdrawal data
	withdrawalURL := c.ServiceUrl + "/withdrawals"
	result, err := c.PostData(withdrawalURL, withdrawal)
	if err != nil {
		return "", fmt.Errorf("failed to post withdrawal data: %v", err)
	}
	return result, nil
}

func (c *CoreBanking) PostDeposit(deposit Deposit) (string, error) {
	// Post Deposit data
	depositURL := c.ServiceUrl + "/deposits"
	result, err := c.PostData(depositURL, deposit)
	if err != nil {
		return "", fmt.Errorf("failed to post deposit data: %v", err)
	}
	return result, nil
}

// postData sends a POST request with JSON data to the specified URL
func (c *CoreBanking) PostData(url string, data interface{}) (string, error) {
	// Marshal the data to JSON
	requestBody, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Convert the byte array to a string
	/* jsonString := string(requestBody)

	// Print the JSON string
	fmt.Printf("/nwithdrawal:%v",jsonString)
	*/
	// Create the HTTP request
	// fmt.Println("PostData: ",url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("Set Header: ")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Error")
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and decode the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("error response from server: %s", respBody)
	}

	// fmt.Printf("\nStatusCode%v",resp.StatusCode)

	return string(respBody), nil
}
