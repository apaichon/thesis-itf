package repositories

import (
	"banking360/config"
	"banking360/internal/data"
	"banking360/internal/data/models"
	"fmt"
	"log"
	"time"
    "strings"
    "sync"

	"github.com/google/uuid"
)

type FinancialTransactionRepo struct {
	DBPool *data.DBPool
}

var (
	instance *FinancialTransactionRepo
	once     sync.Once
)

// NewMessageRepo creates a new instance of MessageRepo
/*func NewFinancialTransactionRepo() *FinancialTransactionRepo {
	cfg := config.NewConfig()
	dsn := cfg.GetDns()
	pool, err := data.GetDBPool(dsn, cfg.MaxConnection, cfg.MaxIdleConnection, time.Duration(cfg.ConnectionLifeTime)*time.Minute)
	if err != nil {
		log.Fatalf("Failed to get database pool: %v", err)
	}
	return &FinancialTransactionRepo{DBPool: pool}
}*/

func NewFinancialTransactionRepo() *FinancialTransactionRepo {
	once.Do(func() {
		cfg := config.NewConfig()
		dsn := cfg.GetDns()
		pool, err := data.GetDBPool(dsn, cfg.MaxConnection, cfg.MaxIdleConnection, time.Duration(cfg.ConnectionLifeTime)*time.Minute)
		if err != nil {
			log.Fatalf("Failed to get database pool: %v", err)
		}
		instance = &FinancialTransactionRepo{DBPool: pool}
	})
	return instance
}

func (fin *FinancialTransactionRepo) Deposit(trans models.FinancialTransaction) (uuid.UUID, error) {

	command := `INSERT INTO banking360.financial_transactions
(
    transaction_id,
    transaction_type,
    transaction_date,
    amount,
    currency,
    account_id,
    status,
    description,
	act_at,
    act_by,
    created_at,
    created_by,
    sign
)
VALUES
(
    ?,  -- Generate a new UUID for the transaction
    'deposit',         -- transaction_type for deposit
    ?,             -- Current timestamp for transaction_date
    ?,           -- Amount of the deposit (example: 1000.00)
    ?,             -- Currency (example: USD)
    ?,  -- account_id (replace with actual account UUID)
    'completed', -- status (assuming the deposit is completed)
    ?, -- description
	?, -- act at
	?, -- act by
    now(),             -- Current timestamp for created_at
    ?,  -- created_by (replace with actual user UUID who created this transaction)
    1                  -- sign (1 for insertion, -1 would be for deletion in a collapsing merge tree)
)`
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}

	_, err =
		fin.DBPool.Insert(command,
			id,
			trans.TransactionDate,
			trans.Amount,
			trans.Currency,
			trans.AccountID,
			trans.Description,
			trans.ActAt,
			trans.ActBy,
			trans.CreatedBy,
		)
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
	return id, nil
}

func (fin *FinancialTransactionRepo) GetDepositsByTextSearchPagination(searchText string, page, pageSize int) ([]*models.FinancialTransaction, *models.PaginationModel, error) {
	var deposits []*models.FinancialTransaction

	query := fmt.Sprintf(`
        SELECT 
            transaction_id,
            transaction_date,
            amount,
            currency,
            account_id,
            status,
            description,
            created_at,
            created_by
        FROM banking360.financial_transactions final
        WHERE transaction_type = 'deposit'
        AND (
            cast(amount as String) LIKE '%%%s%%'
            OR currency LIKE '%%%s%%'
            OR cast(account_id as String) LIKE '%%%s%%'
            OR description LIKE '%%%s%%'
        )
    `, searchText, searchText, searchText, searchText)

	offset := (page - 1) * pageSize
	limit := pageSize

	pagination := data.NewPagination(page, pageSize, query, limit, offset)

	pager, err := pagination.GetPageData(fin.DBPool)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting page data: %v", err)
	}

	query = query + " ORDER BY transaction_date DESC LIMIT ? OFFSET ?"

	rows, err := fin.DBPool.Query(query, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deposit models.FinancialTransaction
		var statusStr string
		err := rows.Scan(
			&deposit.TransactionID,
			&deposit.TransactionDate,
			&deposit.Amount,
			&deposit.Currency,
			&deposit.AccountID,
			&statusStr,
			&deposit.Description,
			&deposit.CreatedAt,
			&deposit.CreatedBy,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %v", err)
		}
		deposit.TransactionType = models.Deposit

		status, err := data.StringToTransactionStatus(statusStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error converting status: %v", err)
		}
		deposit.Status = status

		deposits = append(deposits, &deposit)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	return deposits, pager, nil
}


func (fin *FinancialTransactionRepo) PerformIntraBankTransferOld(transaction models.FinancialTransaction) (uuid.UUID, error) {
    // Generate IDs
	transactionID, err := uuid.NewV7()
	if err != nil {
		transactionID = uuid.New()
	}

    // Get the current time
    now := time.Now()

    // Create the financial transaction
    transaction.TransactionID = transactionID
	transaction.TransactionType = models.TransferIntra
	transaction.TransactionDate= now
	transaction.Status = models.Completed
	transaction.Sign = 1
	transaction.CreatedAt = now
	transaction.Description =  "Intra Bank Transfer"

    // Insert the financial transaction
    _, err = fin.DBPool.Insert(
		`INSERT INTO banking360.financial_transactions_memory (
            transaction_id, transaction_type, transaction_date, amount, currency, account_id,
            recipient_account_id, status, description, act_by, act_at, created_at, created_by, sign
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        transaction.TransactionID, transaction.TransactionType, transaction.TransactionDate, transaction.Amount, transaction.Currency, transaction.AccountID,
        transaction.RecipientAccountID, transaction.Status, transaction.Description, transaction.ActBy, transaction.ActAt, transaction.CreatedAt, transaction.CreatedBy, transaction.Sign)
    if err != nil {
		return transactionID, err
    }

	senderBalanceID, err := uuid.NewV7()
	if err != nil {
		senderBalanceID = uuid.New()
	}
    // Create the sender account balance
    senderBalance := models.AccountBalance{
        BalanceID:      senderBalanceID,
        AccountID:      transaction.AccountID,
        TransactionID:  transactionID,
        Timestamp:      now,
        OperationType:  models.Debit,
        Amount:         transaction.Amount *-1,
        RunningBalance: 0, // Assume this is calculated elsewhere
        Description:    transaction.Description,
    }

	
    // Insert the sender account balance
    _, err = fin.DBPool.Insert (`
        INSERT INTO banking360.account_balances_memory (
            balance_id, account_id, transaction_id, timestamp, operation_type, amount, running_balance, description
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
        senderBalanceID, senderBalance.AccountID, senderBalance.TransactionID, senderBalance.Timestamp, senderBalance.OperationType, senderBalance.Amount, senderBalance.RunningBalance, senderBalance.Description)
    if err != nil {
        return transactionID, err
    }

	receiverBalanceID, err := uuid.NewV7()
	if err != nil {
		receiverBalanceID = uuid.New()
	}
    // Create the receiver account balance
    receiverBalance := models.AccountBalance{
        BalanceID:      receiverBalanceID,
       
        TransactionID:  transactionID,
        Timestamp:      now,
        OperationType:  models.Credit,
        Amount:         transaction.Amount,
        RunningBalance: 0, // Assume this is calculated elsewhere
        Description:    transaction.Description,
    }

	if transaction.RecipientAccountID !=nil{
		receiverBalance.AccountID = *transaction.RecipientAccountID
	}

    // Insert the receiver account balance
    _, err = fin.DBPool.Insert(`
        INSERT INTO banking360.account_balances_memory (
            balance_id, account_id, transaction_id, timestamp, operation_type, amount, running_balance, description
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
        receiverBalance.BalanceID, receiverBalance.AccountID, receiverBalance.TransactionID, receiverBalance.Timestamp, receiverBalance.OperationType, receiverBalance.Amount, receiverBalance.RunningBalance, receiverBalance.Description)
    if err != nil {
        return transactionID, err
    }

    return transactionID, nil
}

func (fin *FinancialTransactionRepo) PerformIntraBankTransfer(transaction models.FinancialTransaction) (uuid.UUID, error) {
    // Generate IDs
	transactionID, err := uuid.NewV7()
	if err != nil {
		transactionID = uuid.New()
	}

    // Get the current time
    now := time.Now()

    // Create the financial transaction
    transaction.TransactionID = transactionID
	transaction.TransactionType = models.TransferIntra
	transaction.TransactionDate= now
	transaction.Status = models.Completed
	transaction.Sign = 1
	transaction.CreatedAt = now
	transaction.Description =  "Intra Bank Transfer"


	senderBalanceID, err := uuid.NewV7()
	if err != nil {
		senderBalanceID = uuid.New()
	}
    // Create the sender account balance
    senderBalance := models.AccountBalance{
        BalanceID:      senderBalanceID,
        AccountID:      transaction.AccountID,
        TransactionID:  transactionID,
        Timestamp:      now,
        OperationType:  models.Debit,
        Amount:         transaction.Amount *-1,
        RunningBalance: 0, // Assume this is calculated elsewhere
        Description:    transaction.Description,
    }

	

	receiverBalanceID, err := uuid.NewV7()
	if err != nil {
		receiverBalanceID = uuid.New()
	}
    // Create the receiver account balance
    receiverBalance := models.AccountBalance{
        BalanceID:      receiverBalanceID,
       
        TransactionID:  transactionID,
        Timestamp:      now,
        OperationType:  models.Credit,
        Amount:         transaction.Amount,
        RunningBalance: 0, // Assume this is calculated elsewhere
        Description:    transaction.Description,
    }

	if transaction.RecipientAccountID !=nil{
		receiverBalance.AccountID = *transaction.RecipientAccountID
	}

    transactionResultChan := make(chan error)
	balanceResultChan := make(chan error)

    go insertTransaction(fin.DBPool, transaction, transactionResultChan)
	go insertAccountBalance(fin.DBPool, senderBalance, balanceResultChan)
    go insertAccountBalance(fin.DBPool, receiverBalance, balanceResultChan)

    return transactionID, nil
}

func (fin *FinancialTransactionRepo) PerformIntraBankTransfers(transactions []models.FinancialTransaction) ([]uuid.UUID, error) {
    var transactionIDs []uuid.UUID
    var financialTransactionSQLs []string
    var financialTransactionValues []interface{}
    var accountBalanceSQLs []string
    var accountBalanceValues []interface{}

    now := time.Now()

    for _, transaction := range transactions {
        // Generate IDs
        transactionID, err := uuid.NewV7()
        if err != nil {
            transactionID = uuid.New()
        }
        transactionIDs = append(transactionIDs, transactionID)

        // Create the financial transaction
        transaction.TransactionID = transactionID
        transaction.TransactionType = models.TransferIntra
        transaction.TransactionDate = now
        transaction.Status = models.Completed
        transaction.Sign = 1
        transaction.CreatedAt = now
        transaction.Description = "Intra Bank Transfer"

        // Prepare SQL and values for financial transaction insert
        financialTransactionSQLs = append(financialTransactionSQLs, `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
        financialTransactionValues = append(financialTransactionValues,
            transaction.TransactionID, transaction.TransactionType, transaction.TransactionDate, transaction.Amount, transaction.Currency, transaction.AccountID,
            transaction.RecipientAccountID, transaction.Status, transaction.Description, transaction.ActBy, transaction.ActAt, transaction.CreatedAt, transaction.CreatedBy, transaction.Sign,
        )

       // fmt.Printf ("sqlinsert:%v", strings.Join(financialTransactionSQLs, ", "))

        // Create the sender account balance
        senderBalanceID, err := uuid.NewV7()
        if err != nil {
            senderBalanceID = uuid.New()
        }
        senderBalance := models.AccountBalance{
            BalanceID:      senderBalanceID,
            AccountID:      transaction.AccountID,
            TransactionID:  transactionID,
            Timestamp:      now,
            OperationType:  models.Debit,
            Amount:         transaction.Amount * -1,
            RunningBalance: 0, // Assume this is calculated elsewhere
            Description:    transaction.Description,
        }

        // Prepare SQL and values for sender account balance insert
        accountBalanceSQLs = append(accountBalanceSQLs, `(?, ?, ?, ?, ?, ?, ?, ?)`)
        accountBalanceValues = append(accountBalanceValues,
            senderBalance.BalanceID, senderBalance.AccountID, senderBalance.TransactionID, senderBalance.Timestamp, senderBalance.OperationType, senderBalance.Amount, senderBalance.RunningBalance, senderBalance.Description,
        )

        // Create the receiver account balance
        receiverBalanceID, err := uuid.NewV7()
        if err != nil {
            receiverBalanceID = uuid.New()
        }
        receiverBalance := models.AccountBalance{
            BalanceID:      receiverBalanceID,
            TransactionID:  transactionID,
            Timestamp:      now,
            OperationType:  models.Credit,
            Amount:         transaction.Amount,
            RunningBalance: 0, // Assume this is calculated elsewhere
            Description:    transaction.Description,
        }

        if transaction.RecipientAccountID != nil {
            receiverBalance.AccountID = *transaction.RecipientAccountID
        }

        // Prepare SQL and values for receiver account balance insert
        accountBalanceSQLs = append(accountBalanceSQLs, `(?, ?, ?, ?, ?, ?, ?, ?)`)
        accountBalanceValues = append(accountBalanceValues,
            receiverBalance.BalanceID, receiverBalance.AccountID, receiverBalance.TransactionID, receiverBalance.Timestamp, receiverBalance.OperationType, receiverBalance.Amount, receiverBalance.RunningBalance, receiverBalance.Description,
        )
    }

    /*// Combine SQL statements for batch insert
    batchSQL := strings.Join(financialTransactionSQLs, "; ") + "; " + strings.Join(accountBalanceSQLs, "; ")
    batchValues := append(financialTransactionValues, accountBalanceValues...)

    // Execute batch insert
    _, err := fin.DBPool.Insert(batchSQL, batchValues...)
    if err != nil {
        return transactionIDs, err
    }
    */    

    _, err := fin.DBPool.Insert(`INSERT INTO banking360.financial_transactions (
            transaction_id, transaction_type, transaction_date, amount, currency, account_id,
            recipient_account_id, status, description, act_by, act_at, created_at, created_by, sign
        ) VALUES ` + strings.Join(financialTransactionSQLs, ", ") , financialTransactionValues...)
    if err != nil {
        fmt.Println("Error inserting financial transactions")
        return transactionIDs, err
    }

    _, err = fin.DBPool.Insert(` INSERT INTO banking360.account_balances (
                balance_id, account_id, transaction_id, timestamp, operation_type, amount, running_balance, description
            ) VALUES ` +strings.Join(accountBalanceSQLs, ", "), accountBalanceValues...)
    if err != nil {
        fmt.Println("Error inserting account balance transactions")
        return transactionIDs, err
    }

    return transactionIDs, nil
}


func (fin *FinancialTransactionRepo) PerformInterBankTransfer(transaction models.FinancialTransaction) (uuid.UUID, error) {
    // Generate IDs
	transactionID, err := uuid.NewV7()
	if err != nil {
		transactionID = uuid.New()
	}

    // Get the current time
    now := time.Now()

    // Create the financial transaction
    transaction.TransactionID = transactionID
	transaction.TransactionType = models.TransferInter
	transaction.TransactionDate= now
	transaction.Status = models.Completed
	transaction.Sign = 1
	transaction.CreatedAt = now
	transaction.Description =  "Inter Bank Transfer"

    // Insert the financial transaction
    _, err = fin.DBPool.Insert(
		`INSERT INTO banking360.financial_transactions (
            transaction_id, transaction_type, transaction_date, amount, currency, account_id,
            recipient_account_id, status, description, act_by, act_at, created_at, created_by, sign
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        transaction.TransactionID, transaction.TransactionType, transaction.TransactionDate, transaction.Amount, transaction.Currency, transaction.AccountID,
        transaction.RecipientAccountID, transaction.Status, transaction.Description, transaction.ActBy, transaction.ActAt, transaction.CreatedAt, transaction.CreatedBy, transaction.Sign)
    if err != nil {
		return transactionID, err
    }

	senderBalanceID, err := uuid.NewV7()
	if err != nil {
		senderBalanceID = uuid.New()
	}
    // Create the sender account balance
    senderBalance := models.AccountBalance{
        BalanceID:      senderBalanceID,
        AccountID:      transaction.AccountID,
        TransactionID:  transactionID,
        Timestamp:      now,
        OperationType:  models.Debit,
        Amount:         transaction.Amount *-1,
        RunningBalance: 0, // Assume this is calculated elsewhere
        Description:    transaction.Description,
    }

	
    // Insert the sender account balance
    _, err = fin.DBPool.Insert (`
        INSERT INTO banking360.account_balances (
            balance_id, account_id, transaction_id, timestamp, operation_type, amount, running_balance, description
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
        senderBalanceID, senderBalance.AccountID, senderBalance.TransactionID, senderBalance.Timestamp, senderBalance.OperationType, senderBalance.Amount, senderBalance.RunningBalance, senderBalance.Description)
    if err != nil {
		return transactionID, err
    }


    return transactionID, nil
}

func insertTransaction(db *data.DBPool, transaction models.FinancialTransaction, resultChan chan<- error) {
	_, err := db.Insert(
		`INSERT INTO banking360.financial_transactions (
            transaction_id, transaction_type, transaction_date, amount, currency, account_id,
            recipient_account_id, status, description, act_by, act_at, created_at, created_by, sign
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.TransactionID, transaction.TransactionType, transaction.TransactionDate, transaction.Amount, transaction.Currency, transaction.AccountID,
		transaction.RecipientAccountID, transaction.Status, transaction.Description, transaction.ActBy, transaction.ActAt, transaction.CreatedAt, transaction.CreatedBy, transaction.Sign)

	resultChan <- err
}

func insertAccountBalance(db *data.DBPool, balance models.AccountBalance, resultChan chan<- error) {
	_, err := db.Insert(
		`INSERT INTO banking360.account_balances (
            balance_id, account_id, transaction_id, timestamp, operation_type, amount, running_balance, description
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		balance.BalanceID, balance.AccountID, balance.TransactionID, balance.Timestamp, balance.OperationType, balance.Amount, balance.RunningBalance, balance.Description)

	resultChan <- err
}

