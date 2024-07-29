-- Clear Data For Test
delete from TransactionHistory
delete from Deposit
delete from Withdrawal

with dep as (
    Select count(*) as totalDeposit  
        ,Max(CreatedAt) as maxDepositDate
     ,Min(CreatedAt) as minDepositDate 
    From Deposit
),
withd as ( 
    select count(*) as totalWithdrawal 
    ,Max(CreatedAt) as maxWithdrawalDate
     ,Min(CreatedAt) as minWithdrawalDate    
    From Withdrawal 
),
trans as (
    select count(*) as totalTrans   
    from TransactionHistory
)

select * ,
(strftime('%s', maxDate) - strftime('%s', minDate)) * 1000 
   + 
    (strftime('%f', maxDate) - strftime('%f', minDate)) * 1000
     AS 
    duration_in_milliseconds
from (
select * 
,
 (strftime('%s', dep.maxDepositDate) - strftime('%s', dep.minDepositDate)) * 1000 
   + 
    (strftime('%f', dep.maxDepositDate) - strftime('%f', dep.minDepositDate)) * 1000
     AS 
    deposit_duration_in_milliseconds
    ,
 (strftime('%s', withd.maxWithdrawalDate) - strftime('%s', withd.minWithdrawalDate)) * 1000  + 
    (strftime('%f', withd.maxWithdrawalDate) - strftime('%f', withd.minWithdrawalDate)) * 1000 AS 
    withdrawal_duration_in_milliseconds

,julianday(withd.maxWithdrawalDate) - julianday(withd.minWithdrawalDate) * 86400000 AS milliseconds_difference
,case when minDepositDate < minWithdrawalDate then minDepositDate else minWithdrawalDate end as minDate
,case when maxDepositDate > maxWithdrawalDate then maxDepositDate else maxWithdrawalDate end as maxDate
from dep, withd, trans
)a