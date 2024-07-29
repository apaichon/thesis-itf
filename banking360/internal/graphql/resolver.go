package graphql

import (
	// "banking360/internal/data"
	"banking360/internal/data/models"
	"banking360/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func DepositResolver(params graphql.ResolveParams) (interface{}, error) {
	input := params.Args["input"].(map[string]interface{})

	accountId, err := uuid.Parse(input["accountId"].(string))
	if err != nil {
		fmt.Println("Invalid UUID:", accountId)
		return nil, err
	}

	actBy, err := uuid.Parse(input["actBy"].(string))
	if err != nil {
		fmt.Println("Invalid UUID:", actBy)
		return nil, err
	}

	createdBy, err := uuid.Parse(input["createdBy"].(string))
	if err != nil {
		fmt.Println("Invalid UUID:", createdBy)
		return nil, err
	}

	destination := input["destinationInstitution"].(string)
	description := fmt.Sprintf("user:%v deposit %v to account:%v", actBy, input["amount"].(float64), accountId)
	depositInput := models.FinancialTransaction{
			TransactionDate:        input["transactionDate"].(time.Time),
			Amount:                 input["amount"].(float64),
			Currency:               "THB",
			AccountID:              accountId,
			DestinationInstitution: &destination,
			Status:                 models.Completed,
			Description:            description,
			ActBy:                  actBy,
			ActAt:                  time.Now(),
			CreatedBy:              createdBy,
			CreatedAt:              time.Now(),
		}
		finRepo := repositories.NewFinancialTransactionRepo()
		id, err := finRepo.Deposit(depositInput)
		if err != nil {
			return nil, err
		}
		return models.NewSuccessModel("deposit:" + id.String() + " is successfully."), nil
}

func GetDepositsByTextSearchPaginationResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	page := defaultValue(params.Args, "page", 1).(int)
	pageSize := defaultValue(params.Args, "pageSize", 10).(int)
	textSearch := defaultValue(params.Args, "textSearch", "").(string)

	repo := repositories.NewFinancialTransactionRepo()

	// Fetch contacts from the database
	deposits, pager, err := repo.GetDepositsByTextSearchPagination(textSearch, page, pageSize)
	var depositPagination = models.FinancialTransactionPaginationModel{
		FinancialTransactions: deposits,
		Pagination:            pager,
	}

	if err != nil {
		return nil, err
	}

	return depositPagination, nil
}

func defaultValue(args map[string]interface{}, key string, defaultVal interface{}) interface{} {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	return val
}

func TransferIntraBankResolver(params graphql.ResolveParams) (interface{}, error) {
    input := params.Args["input"].(map[string]interface{})

    senderAccountID := parseUUID(input, "senderAccountId")
    receiverAccountID := parseUUID(input, "receiverAccountId")
    actBy := parseUUID(input, "actBy")
    createdBy := parseUUID(input, "createdBy")

    if senderAccountID == uuid.Nil || receiverAccountID == uuid.Nil || actBy == uuid.Nil || createdBy == uuid.Nil {
        return nil, fmt.Errorf("invalid UUID input")
    }

    amount := input["amount"].(float64)
    currency := "THB"
    transactionDate := input["transactionDate"].(time.Time)
    description := fmt.Sprintf("user:%v transfers %v to account:%v", actBy, amount, receiverAccountID)

    // Create financial transaction for sender
    senderTransaction := models.FinancialTransaction{
        TransactionDate:        transactionDate,
        Amount:                 -amount,
        Currency:               currency,
        AccountID:              senderAccountID,
		RecipientAccountID:		&receiverAccountID,
        DestinationInstitution: nil,
        Status:                 models.Completed,
        Description:            description,
        ActBy:                  actBy,
        ActAt:                  time.Now(),
        CreatedBy:              createdBy,
        CreatedAt:              time.Now(),
    }
    // Initialize repository
    finRepo := repositories.NewFinancialTransactionRepo()

	_ , err := finRepo.PerformIntraBankTransfer(senderTransaction)

    if err != nil {
        return nil, err
    }

    return models.NewSuccessModel(fmt.Sprintf("transfer:%s and %s are successfully processed", senderAccountID.String(),  receiverAccountID.String())), nil
}

func TransferInterBankResolver(params graphql.ResolveParams) (interface{}, error) {
    input := params.Args["input"].(map[string]interface{})

    senderAccountID := parseUUID(input, "senderAccountId")
    receiverAccountID := parseUUID(input, "receiverAccountId")
    actBy := parseUUID(input, "actBy")
    createdBy := parseUUID(input, "createdBy")

    if senderAccountID == uuid.Nil || receiverAccountID == uuid.Nil || actBy == uuid.Nil || createdBy == uuid.Nil {
        return nil, fmt.Errorf("invalid UUID input")
    }

    amount := input["amount"].(float64)
    currency := "THB"
    transactionDate := input["transactionDate"].(time.Time)
    description := fmt.Sprintf("user:%v transfers %v to account:%v", actBy, amount, receiverAccountID)

    // Create financial transaction for sender
    senderTransaction := models.FinancialTransaction{
        TransactionDate:        transactionDate,
        Amount:                 -amount,
        Currency:               currency,
        AccountID:              senderAccountID,
		RecipientAccountID:     &receiverAccountID,
        DestinationInstitution: nil,
        Status:                 models.Completed,
        Description:            description,
        ActBy:                  actBy,
        ActAt:                  time.Now(),
        CreatedBy:              createdBy,
        CreatedAt:              time.Now(),
    }
    // Initialize repository
    finRepo := repositories.NewFinancialTransactionRepo()

	 _, err := finRepo.PerformInterBankTransfer(senderTransaction)

   
    if err != nil {
        return nil, err
    }

    return models.NewSuccessModel(fmt.Sprintf("transfer:%s and %s are successfully processed", senderAccountID.String(),  receiverAccountID.String())), nil
}

func parseUUID(input map[string]interface{}, key string) uuid.UUID {
    value, ok := input[key].(string)
    if !ok {
        fmt.Printf("Invalid type for key %s\n", key)
        return uuid.Nil
    }

    parsedUUID, err := uuid.Parse(value)
    if err != nil {
        fmt.Printf("Invalid UUID for key %s: %v\n", key, err)
        return uuid.Nil
    }

    return parsedUUID
}