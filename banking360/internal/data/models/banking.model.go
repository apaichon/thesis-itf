package models

import (
    "github.com/google/uuid"
    "time"
)

type TransactionType int8
const (
    Deposit TransactionType = iota + 1
    Withdraw
    TransferIntra
    TransferInter
)

type TransactionStatus int8
const (
    Pending TransactionStatus = iota + 1
    Completed
    Failed
)

type FinancialTransaction struct {
    TransactionID          uuid.UUID         `json:"transaction_id"`
    TransactionType        TransactionType   `json:"transaction_type"`
    TransactionDate        time.Time         `json:"transaction_date"`
    Amount                 float64           `json:"amount"`
    Currency               string            `json:"currency"`
    AccountID              uuid.UUID         `json:"account_id"`
    RecipientAccountID     *uuid.UUID        `json:"recipient_account_id"`
    SourceInstitution      *string           `json:"source_institution"`
    DestinationInstitution *string           `json:"destination_institution"`
    Status                 TransactionStatus `json:"status"`
    Description            string            `json:"description"`
    ActAt              time.Time         `json:"act_at"`
    ActBy              uuid.UUID         `json:"act_by"`
    CreatedAt              time.Time         `json:"created_at"`
    CreatedBy              uuid.UUID         `json:"created_by"`
    Sign                   int8              `json:"sign"`
}

type FinancialTransactionPaginationModel struct {
	FinancialTransactions []*FinancialTransaction `json:"financialTransactions"`
	Pagination *PaginationModel `json:"pagination"`
}

type AccountType int8
const (
    AccountTypeChecking AccountType = iota + 1
    AccountTypeSavings
    AccountTypeInvestment
)

type AccountStatus int8
const (
    Active AccountStatus = iota + 1
    Inactive
    Closed
)

type BankAccount struct {
    AccountID     uuid.UUID     `json:"account_id"`
    AccountNumber string        `json:"account_number"`
    AccountType   AccountType   `json:"account_type"`
    Currency      string        `json:"currency"`
    OwnerID       uuid.UUID     `json:"owner_id"`
    Status        AccountStatus `json:"status"`
    CreatedAt     time.Time     `json:"created_at"`
    Sign          int8          `json:"sign"`
    Version       uint16        `json:"version"`
}

type OperationType int8
const (
    Debit OperationType = iota + 1
    Credit
)

type AccountBalance struct {
    BalanceID      uuid.UUID     `json:"balance_id"`
    AccountID      uuid.UUID     `json:"account_id"`
    TransactionID  uuid.UUID     `json:"transaction_id"`
    Timestamp      time.Time     `json:"timestamp"`
    OperationType  OperationType `json:"operation_type"`
    Amount         float64       `json:"amount"`
    RunningBalance float64       `json:"running_balance"`
    Description    string        `json:"description"`
}

type InstitutionType int8
const (
    InstitutionTypeCommercial InstitutionType = iota + 1
    InstitutionTypeInvestment
    InstitutionTypeCentral
    InstitutionTypeCooperative
    InstitutionTypeSavings
)

type InstitutionStatus int8
const (
    ActiveInst InstitutionStatus = iota + 1
    InactiveInst
    Suspended
)

type BankInstitution struct {
    InstitutionID       uuid.UUID         `json:"institution_id"`
    Name                string            `json:"name"`
    ShortName           string            `json:"short_name"`
    SwiftCode           string            `json:"swift_code"`
    Country             string            `json:"country"`
    Type                InstitutionType   `json:"type"`
    Status              InstitutionStatus `json:"status"`
    FoundedDate         time.Time         `json:"founded_date"`
    Website             string            `json:"website"`
    HeadquartersAddress string            `json:"headquarters_address"`
    RegulatoryBody      string            `json:"regulatory_body"`
    CreatedAt           time.Time         `json:"created_at"`
    Sign                int8              `json:"sign"`
}

type OwnerType int8
const (
    Individual OwnerType = iota + 1
    Corporate
    Government
)

type BankOwner struct {
    OwnerID           uuid.UUID `json:"owner_id"`
    OwnerType         OwnerType `json:"owner_type"`
    FirstName         *string   `json:"first_name"`
    LastName          *string   `json:"last_name"`
    DateOfBirth       *time.Time `json:"date_of_birth"`
    Nationality       *string   `json:"nationality"`
    CompanyName       *string   `json:"company_name"`
    RegistrationNumber *string   `json:"registration_number"`
    IncorporationDate *time.Time `json:"incorporation_date"`
    TaxID             string    `json:"tax_id"`
    Address           string    `json:"address"`
    City              string    `json:"city"`
    State             string    `json:"state"`
    Country           string    `json:"country"`
    PostalCode        string    `json:"postal_code"`
    Email             string    `json:"email"`
    PhoneNumber       string    `json:"phone_number"`
    CreatedAt         time.Time `json:"created_at"`
    Sign              int8      `json:"sign"`
    Version           uint16    `json:"version"`
}