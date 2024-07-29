-- clear data for load test
truncate table banking360.financial_transactions
truncate table banking360.account_balances 
truncate table QueueManager.MessageQueues
truncate table ProcessManager.MessageQueues
truncate table ProcessManager.TaskQueues

-- Validate load test results
With fin as (
    select count(*)  as totalFin
    --  , dateDiff('second',Min(created_at), Max(created_at))
     -- ,  (Max(created_at) - Min(created_at)) *1000 as duration
    ,  Min(created_at) MinFin, Max(created_at) MaxFin
     from banking360.financial_transactions final
),
queue as (
    select count(*) as totalMqs ,Max(CreatedAt) MaxQueue , Min(CreatedAt) MinQueue  from QueueManager.MessageQueues
),
process as (
    select count(*) as totalProcQs 
    ,Max(CreatedAt) MaxQueueProcess , Min(CreatedAt) MinQueueProcess
    from ProcessManager.MessageQueues
),
task as (
    select count(*) as totalTasks
     ,Max(CreatedAt) MaxTaskProcess , Min(CreatedAt) MinTaskProcess
      from ProcessManager.TaskQueues
)
select *
, MaxFin - MinFin as TotalFinDuraitons
, MaxFin -  ( Select Min(CreatedAt) from QueueManager.MessageQueues) as TotalDuraitons
from  queue, process, task, fin