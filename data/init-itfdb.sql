CREATE DATABASE QueueManager;

CREATE TABLE QueueManager.MessageQueues
(
    Id UUID,
    SystemOwner String,
    Topic String,
    Content String,
    Remarks String,
    CreatedAt DateTime,
    CreatedBy DateTime,
    Sign Int8
)
ENGINE = CollapsingMergeTree(Sign)
ORDER BY (SystemOwner, Topic, Id)
PARTITION BY toYYYYMM(CreatedAt)
PRIMARY KEY (SystemOwner, Topic)


CREATE TABLE QueueManager.Archive
(
    Id UUID,
    CorrelationId UUID,
    TableName LowCardinality(String),
    JsonData  LowCardinality(String), 
    CreatedAt DateTime,
    CreatedBy LowCardinality(String)
)
ENGINE = MergeTree
ORDER BY (Id, TableName, CorrelationId)
PARTITION BY toYYYYMM(CreatedAt)
PRIMARY KEY (Id,TableName);


CREATE DATABASE OrchestrationManager;

CREATE TABLE OrchestrationManager.ServiceProvider
(
    WorkerGroup LowCardinality(String),
    ServiceName LowCardinality(String),
    ServiceFullName LowCardinality(String),
    ServiceDescription LowCardinality(String),
    Tps UInt32,
    CreatedAt DateTime,
    CreatedBy LowCardinality(String),
    Active Int8
)
ENGINE = ReplacingMergeTree(Active)
ORDER BY (WorkerGroup,ServiceName)
PRIMARY KEY (WorkerGroup,ServiceName)

select * from OrchestrationManager.ServiceProvider
INSERT INTO OrchestrationManager.ServiceProvider
(
    WorkerGroup,
    ServiceName,
    ServiceFullName,
    ServiceDescription,
    Tps,
    CreatedAt,
    CreatedBy,
    Active
)
VALUES
(
    'Banking360',
    'TransferIntraBankService',
    'Banking360.TransferIntraBankService',
    'transfer money to same bank service.',
    10000,
    now(),
    'admin',
    1
),
(
    'Banking360',
    'TransferInterBankService',
    'Banking360.TransferInterBankService',
    'transfer money to another bank service.',
    10000,
    now(),
    'admin',
    1
),
(
    'Legacy',
    'DepositService',
    'Legacy.DepositService',
    'deposit money to bank account.',
    200,
    now(),
    'admin',
    1
),
(
    'Legacy',
    'WithdrawService',
    'Legacy.WithdrawService',
    'withdraw money from bank account.',
    200,
    now(),
    'admin',
    1
)


CREATE TABLE OrchestrationManager.ServiceProcesses
(
    Topic  LowCardinality(String), 
    Services Array(LowCardinality(String)),
    Description LowCardinality(String), 
    CreatedAt DateTime,
    CreatedBy LowCardinality(String),
    Sign Int8,
    Version UInt16
)
ENGINE = VersionedCollapsingMergeTree(Sign, Version)
ORDER BY (Topic, Services)
PRIMARY KEY (Topic)


INSERT INTO OrchestrationManager.ServiceProcesses
(
    Topic,
    Services,
    Description,
    CreatedAt,
    CreatedBy,
    Sign,
    Version
)
VALUES
(
    'TransferIntraBank',
    ['Banking360.TransferIntraBankService', 'Legacy.WithdrawService', 'Legacy.DepositService'],
    'Processes transfer money intra bank',
    now(),
    'admin',
    1,
    1
),
(
    'TransferInterBank',
    ['Banking360.TransferInterBankService', 'Legacy.WithdrawService'],
    'Processes transfer money inter bank',
    now(),
    'admin',
    1,
    1
)


CREATE DATABASE ProcessManager;

CREATE TABLE ProcessManager.MessageQueues
(
    Id UUID,
    SystemOwner LowCardinality(String),
    Topic LowCardinality(String),
    Content LowCardinality(String),
    Remarks LowCardinality(String),
    CreatedAt DateTime,
    CreatedBy LowCardinality(String),
    Sign Int8
)
ENGINE = CollapsingMergeTree(Sign)
ORDER BY (SystemOwner, Topic, Id)
PARTITION BY toYYYYMM(CreatedAt)
PRIMARY KEY (SystemOwner, Topic)


CREATE TABLE ProcessManager.TaskQueues
(
    Id UUID,
    ServiceName LowCardinality(String),
    ServiceVersion UInt16,
    MessageIds Array(UUID),
    CreatedAt DateTime,
    CreatedBy LowCardinality(String),
    Sign Int8
)

ENGINE = CollapsingMergeTree(Sign)
ORDER BY (Id, ServiceName)
PARTITION BY toYYYYMM(CreatedAt)
PRIMARY KEY (Id, ServiceName)


CREATE TABLE ProcessManager.Archive
(
    Id UUID,
    CorrelationId UUID,
    TableName LowCardinality(String),
    JsonData  LowCardinality(String), 
    CreatedAt DateTime,
    CreatedBy LowCardinality(String),
)

ENGINE = MergeTree
ORDER BY (Id, TableName, CorrelationId)
PARTITION BY toYYYYMM(CreatedAt)
PRIMARY KEY (Id, TableName)


CREATE TABLE IF NOT EXISTS ProcessManager.AccountMapping (
		AccountId UUID,
		AccountLegacyId Int32
	) ENGINE = MergeTree()
	ORDER BY AccountId;


-- Insert 10,000 sample records
INSERT INTO ProcessManager.AccountMapping
SELECT
    generateUUIDv4() AS AccountId,
    toInt32(rand() % 10000 + 1) AS AccountLegacyId
FROM
    numbers(10000);