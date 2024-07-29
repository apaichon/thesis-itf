package graphql

import (
    // "time"
    // "github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql"
    "banking360/internal/data/models"
)


/* var DateTimeType = graphql.NewScalar(graphql.ScalarConfig{
	Name: "CustomDateTime",
	Serialize: func(value interface{}) interface{} {
		switch t := value.(type) {
		case time.Time:
			return t.Format(time.RFC3339)
		case string:
			return t
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch t := value.(type) {
		case string:
			return t
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST.(type) {
		case *ast.StringValue:
			return valueAST.(*ast.StringValue).Value
		default:
			return nil
		}
	},
})
    */

var ResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Result",
	Fields: graphql.Fields{
		"code":    &graphql.Field{Type: graphql.Int},
		"status":  &graphql.Field{Type: graphql.String},
		"message": &graphql.Field{Type: graphql.String},
		"data":    &graphql.Field{Type: graphql.String},
	},
})

// Enums
var TransactionTypeEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "TransactionType",
    Values: graphql.EnumValueConfigMap{
        "DEPOSIT":       &graphql.EnumValueConfig{Value: 1},
        "WITHDRAW":      &graphql.EnumValueConfig{Value: 2},
        "TRANSFER_INTRA": &graphql.EnumValueConfig{Value: 3},
        "TRANSFER_INTER": &graphql.EnumValueConfig{Value: 4},
    },
})

var TransactionStatusEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "TransactionStatus",
    Values: graphql.EnumValueConfigMap{
        "PENDING":   &graphql.EnumValueConfig{Value: 1},
        "COMPLETED": &graphql.EnumValueConfig{Value: 2},
        "FAILED":    &graphql.EnumValueConfig{Value: 3},
    },
})

var AccountTypeEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "AccountType",
    Values: graphql.EnumValueConfigMap{
        "CHECKING":   &graphql.EnumValueConfig{Value: 1},
        "SAVINGS":    &graphql.EnumValueConfig{Value: 2},
        "INVESTMENT": &graphql.EnumValueConfig{Value: 3},
    },
})

var AccountStatusEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "AccountStatus",
    Values: graphql.EnumValueConfigMap{
        "ACTIVE":   &graphql.EnumValueConfig{Value: 1},
        "INACTIVE": &graphql.EnumValueConfig{Value: 2},
        "CLOSED":   &graphql.EnumValueConfig{Value: 3},
    },
})

var OperationTypeEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "OperationType",
    Values: graphql.EnumValueConfigMap{
        "DEBIT":  &graphql.EnumValueConfig{Value: 1},
        "CREDIT": &graphql.EnumValueConfig{Value: 2},
    },
})

var InstitutionTypeEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "InstitutionType",
    Values: graphql.EnumValueConfigMap{
        "COMMERCIAL":  &graphql.EnumValueConfig{Value: 1},
        "INVESTMENT":  &graphql.EnumValueConfig{Value: 2},
        "CENTRAL":     &graphql.EnumValueConfig{Value: 3},
        "COOPERATIVE": &graphql.EnumValueConfig{Value: 4},
        "SAVINGS":     &graphql.EnumValueConfig{Value: 5},
    },
})

var InstitutionStatusEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "InstitutionStatus",
    Values: graphql.EnumValueConfigMap{
        "ACTIVE":    &graphql.EnumValueConfig{Value: 1},
        "INACTIVE":  &graphql.EnumValueConfig{Value: 2},
        "SUSPENDED": &graphql.EnumValueConfig{Value: 3},
    },
})

var OwnerTypeEnum = graphql.NewEnum(graphql.EnumConfig{
    Name: "OwnerType",
    Values: graphql.EnumValueConfigMap{
        "INDIVIDUAL": &graphql.EnumValueConfig{Value: 1},
        "CORPORATE":  &graphql.EnumValueConfig{Value: 2},
        "GOVERNMENT": &graphql.EnumValueConfig{Value: 3},
    },
})

// Object Types
var FinancialTransactionType = graphql.NewObject(graphql.ObjectConfig{
    Name: "FinancialTransaction",
    Fields: graphql.Fields{
        "transaction_id":           &graphql.Field{Type: graphql.String},
        "transaction_type":         &graphql.Field{Type: TransactionTypeEnum},
        "transaction_date":         &graphql.Field{Type: graphql.DateTime},
        "amount":                   &graphql.Field{Type: graphql.Float},
        "currency":                 &graphql.Field{Type: graphql.String},
        "account_id":               &graphql.Field{Type: graphql.String},
        "recipient_account_id":     &graphql.Field{Type: graphql.String},
        "source_institution":       &graphql.Field{Type: graphql.String},
        "destination_institution":  &graphql.Field{Type: graphql.String},
        "status":                   &graphql.Field{Type: TransactionStatusEnum},
        "description":              &graphql.Field{Type: graphql.String},
		"act_at":               &graphql.Field{Type: graphql.DateTime},
        "act_by":               &graphql.Field{Type: graphql.String},
        "created_at":               &graphql.Field{Type: graphql.DateTime},
        "created_by":               &graphql.Field{Type: graphql.String},
    },
})

var FinancialTransactionPaginationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FinancialTransactionPagination",
	Fields: graphql.Fields{
		"financialTransactions":   &graphql.Field{Type: graphql.NewList(FinancialTransactionType)},
		"pagination": &graphql.Field{Type: PaginationGraphQLType},
		// Add field here
	},
})

var BankAccountType = graphql.NewObject(graphql.ObjectConfig{
    Name: "BankAccount",
    Fields: graphql.Fields{
        "account_id":     &graphql.Field{Type: graphql.String},
        "account_number": &graphql.Field{Type: graphql.String},
        "account_type":   &graphql.Field{Type: AccountTypeEnum},
        "currency":       &graphql.Field{Type: graphql.String},
        "owner_id":       &graphql.Field{Type: graphql.String},
        "status":         &graphql.Field{Type: AccountStatusEnum},
        "created_at":     &graphql.Field{Type: graphql.DateTime},
    },
})

var AccountBalanceType = graphql.NewObject(graphql.ObjectConfig{
    Name: "AccountBalance",
    Fields: graphql.Fields{
        "balance_id":      &graphql.Field{Type: graphql.String},
        "account_id":      &graphql.Field{Type: graphql.String},
        "transaction_id":  &graphql.Field{Type: graphql.String},
        "timestamp":       &graphql.Field{Type: graphql.DateTime},
        "operation_type":  &graphql.Field{Type: OperationTypeEnum},
        "amount":          &graphql.Field{Type: graphql.Float},
        "running_balance": &graphql.Field{Type: graphql.Float},
        "description":     &graphql.Field{Type: graphql.String},
    },
})

var BankInstitutionType = graphql.NewObject(graphql.ObjectConfig{
    Name: "BankInstitution",
    Fields: graphql.Fields{
        "institution_id":       &graphql.Field{Type: graphql.String},
        "name":                 &graphql.Field{Type: graphql.String},
        "short_name":           &graphql.Field{Type: graphql.String},
        "swift_code":           &graphql.Field{Type: graphql.String},
        "country":              &graphql.Field{Type: graphql.String},
        "type":                 &graphql.Field{Type: InstitutionTypeEnum},
        "status":               &graphql.Field{Type: InstitutionStatusEnum},
        "founded_date":         &graphql.Field{Type: graphql.DateTime},
        "website":              &graphql.Field{Type: graphql.String},
        "headquarters_address": &graphql.Field{Type: graphql.String},
        "regulatory_body":      &graphql.Field{Type: graphql.String},
        "created_at":           &graphql.Field{Type: graphql.DateTime},
    },
})

var BankOwnerType = graphql.NewObject(graphql.ObjectConfig{
    Name: "BankOwner",
    Fields: graphql.Fields{
        "owner_id":           &graphql.Field{Type: graphql.String},
        "owner_type":         &graphql.Field{Type: OwnerTypeEnum},
        "first_name":         &graphql.Field{Type: graphql.String},
        "last_name":          &graphql.Field{Type: graphql.String},
        "date_of_birth":      &graphql.Field{Type: graphql.DateTime},
        "nationality":        &graphql.Field{Type: graphql.String},
        "company_name":       &graphql.Field{Type: graphql.String},
        "registration_number": &graphql.Field{Type: graphql.String},
        "incorporation_date": &graphql.Field{Type: graphql.DateTime},
        "tax_id":             &graphql.Field{Type: graphql.String},
        "address":            &graphql.Field{Type: graphql.String},
        "city":               &graphql.Field{Type: graphql.String},
        "state":              &graphql.Field{Type: graphql.String},
        "country":            &graphql.Field{Type: graphql.String},
        "postal_code":        &graphql.Field{Type: graphql.String},
        "email":              &graphql.Field{Type: graphql.String},
        "phone_number":       &graphql.Field{Type: graphql.String},
        "created_at":         &graphql.Field{Type: graphql.DateTime},
    },
})

/*
Pagination Type
*/
var PaginationGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Pagination",
	Fields: graphql.Fields{
		"page":        &graphql.Field{Type: graphql.Int},
		"pageSize":    &graphql.Field{Type: graphql.Int},
		"totalPages":  &graphql.Field{Type: graphql.Int},
		"totalItems":  &graphql.Field{Type: graphql.Int},
		"hasNext":     &graphql.Field{Type: graphql.Boolean},
		"hasPrevious": &graphql.Field{Type: graphql.Boolean},
		// Add field here
	},
})


 // Query
type Banking360Queries struct {
	GetDepositsByTextSearch func(string) (*models.FinancialTransactionPaginationModel, error) `json:"getDepositsByTextSearch"`
}

// Define the TicketQueries type
var Banking360QueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Bank360Queries",
	Fields: graphql.Fields{
		"getDepositsByTextSearch": &graphql.Field{
			Type:    FinancialTransactionPaginationType,
			Args:    TextSearchPaginationQueryArgument,
			Resolve: GetDepositsByTextSearchPaginationResolve,
		},
	},
})


// Mutations
type Banking360Mutations struct {
	Deposit  func(map[string]interface{}) (*models.ResultModel, error) `json:"deposit"`
}

// Define the DepositMutations type
var Banking360MutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Banking360Mutations",
	Fields: graphql.Fields{
		"deposit": &graphql.Field{
			Type:    ResultType,
			Args:    DepositArgument,
			Resolve: DepositResolver,
		},
        "transferIntraBank": &graphql.Field{
			Type:    ResultType,
			Args:    TransferIntraBankArgument,
			Resolve: TransferIntraBankResolver,
		},
        "transferInterBank": &graphql.Field{
			Type:    ResultType,
			Args:    TransferInterBankArgument,
			Resolve: TransferInterBankResolver,
		},
	},
})