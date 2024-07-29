
CREATE Database banking360

CREATE TABLE banking360.financial_transactions (
    transaction_id UUID,
    transaction_type Enum8('deposit' = 1, 'withdraw' = 2, 'transfer_intra' = 3, 'transfer_inter' = 4),
    transaction_date DateTime,
    amount Decimal(18, 2),
    currency LowCardinality(String),
    -- For all transaction types
    account_id UUID,
    -- For transfers (both intra and inter)
    recipient_account_id UUID,
    -- For inter-institution transfers
    source_institution LowCardinality(String),
    destination_institution LowCardinality(String),
     -- Additional fields
    status Enum8('pending' = 1, 'completed' = 2, 'failed' = 3),
    description String,
    -- Metadata
    act_by UUID,
    act_at DateTime,
    created_at DateTime Default Now(),
    created_by UUID,
    sign Int8
) ENGINE = CollapsingMergeTree(sign)

ORDER BY (transaction_id, account_id,transaction_date, transaction_type)
PARTITION BY toYYYYMM(transaction_date)
PRIMARY KEY (transaction_id, account_id);

CREATE TABLE banking360.bank_accounts (
    account_id UUID,
    account_number String,
    account_type Enum8('checking' = 1, 'savings' = 2, 'investment' = 3),
    currency LowCardinality(String),
    owner_id UUID,
    status Enum8('active' = 1, 'inactive' = 2, 'closed' = 3),
    created_at DateTime,
    sign Int8,
    version UInt16
) ENGINE = VersionedCollapsingMergeTree(sign, version)
ORDER BY (account_id, created_at)
PRIMARY KEY (account_id,created_at);

CREATE TABLE banking360.account_balances (
    balance_id UUID,
    account_id UUID,
    transaction_id UUID,
    timestamp DateTime,
    operation_type Enum8('debit' = 1, 'credit' = 2),
    amount Decimal(18, 2),
    running_balance Decimal(18, 2),
    description String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (account_id, timestamp);

CREATE TABLE banking360.bank_institutions (
    institution_id UUID,
    name String,
    short_name LowCardinality(String),
    swift_code String,
    country LowCardinality(String),
    type Enum8('commercial' = 1, 'investment' = 2, 'central' = 3, 'cooperative' = 4, 'savings' = 5),
    status Enum8('active' = 1, 'inactive' = 2, 'suspended' = 3),
    founded_date Date,
    website String,
    headquarters_address String,
    regulatory_body LowCardinality(String),
    created_at DateTime,
    sign Int8
) ENGINE = ReplacingMergeTree(sign)
ORDER BY (institution_id, name, country)
 PRIMARY KEY(institution_id, name );

 CREATE TABLE banking360.bank_owners (
    owner_id UUID,
    owner_type Enum8('individual' = 1, 'corporate' = 2, 'government' = 3),
    -- For individuals
    first_name String,
    last_name String,
    date_of_birth Date,
    nationality LowCardinality(String),
    -- For corporate entities
    company_name String,
    registration_number String,
    incorporation_date Date,
    -- Common fields
    tax_id String,
    address String,
    city String,
    state LowCardinality(String),
    country LowCardinality(String),
    postal_code String,
    email String,
    phone_number String,
    created_at DateTime,
    sign Int8,
    version UInt16
) ENGINE = VersionedCollapsingMergeTree(sign, version)
ORDER BY (owner_id, tax_id, owner_type)
PRIMARY KEY (owner_id, tax_id);


-- Create the financial_transactions table with Memory engine
CREATE TABLE banking360.financial_transactions_memory (
    transaction_id UUID,
    transaction_type Enum8('deposit' = 1, 'withdraw' = 2, 'transfer_intra' = 3, 'transfer_inter' = 4),
    transaction_date DateTime,
    amount Decimal(18, 2),
    currency LowCardinality(String),
    account_id UUID,
    recipient_account_id UUID,
    source_institution LowCardinality(String),
    destination_institution LowCardinality(String),
    status Enum8('pending' = 1, 'completed' = 2, 'failed' = 3),
    description String,
    act_by UUID,
    act_at DateTime,
    created_at DateTime Default Now(),
    created_by UUID,
    sign Int8
) ENGINE = Memory
ORDER BY (transaction_id, account_id, transaction_date, transaction_type);

