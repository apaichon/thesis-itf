package graphql

import (
	"github.com/graphql-go/graphql"
)

var depositInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "DepositInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"transactionDate": &graphql.InputObjectFieldConfig{
				Type: graphql.DateTime,
			},
			"amount": &graphql.InputObjectFieldConfig{
				Type: graphql.Float,
			},
			"accountId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"destinationInstitution": &graphql.InputObjectFieldConfig{
				Type: graphql.String, 
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"actAt": &graphql.InputObjectFieldConfig{
				Type: graphql.DateTime,
			},
			"actBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"createdBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

var DepositArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: depositInputType,
	},
}

var transferIntraBankInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "TransferIntraBankInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"transactionDate": &graphql.InputObjectFieldConfig{
				Type: graphql.DateTime,
			},
			"amount": &graphql.InputObjectFieldConfig{
				Type: graphql.Float,
			},
			"senderAccountId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"receiverAccountId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"sourceInstitution": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"destinationInstitution": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"actBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"createdBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

var TransferIntraBankArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: transferIntraBankInputType,
	},
}

var transferInterBankInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "TransferInterbankInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"transactionDate": &graphql.InputObjectFieldConfig{
				Type: graphql.DateTime,
			},
			"amount": &graphql.InputObjectFieldConfig{
				Type: graphql.Float,
			},
			"senderAccountId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"receiverAccountId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"actBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"createdBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

var TransferInterBankArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: transferInterBankInputType,
	},
}

var TransferInterBankArguments = graphql.FieldConfigArgument{
	"inputs": &graphql.ArgumentConfig{
		Type: graphql.NewList(transferInterBankInputType),
	},
}

var TextSearchPaginationQueryArgument = graphql.FieldConfigArgument{
	"textSearch": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"page": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"pageSize": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}
