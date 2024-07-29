using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

public class DepositBatchService : BackgroundService
{
    private readonly IServiceProvider _serviceProvider;
    private static readonly ConcurrentQueue<DepositModel> _depositQueue = new ConcurrentQueue<DepositModel>();
    private static readonly TimeSpan BatchDuration = TimeSpan.FromMilliseconds(500);
    private static readonly SemaphoreSlim _semaphore = new SemaphoreSlim(1, 1);
    private static readonly int BatchSize = 1000; // Define the batch size
    private readonly ILogger<DepositBatchService> _logger;

    public DepositBatchService(IServiceProvider serviceProvider, ILogger<DepositBatchService> logger)
    {
        _serviceProvider = serviceProvider;
         _logger = logger;
    }

    public static void EnqueueDeposit(DepositModel deposit)
    {
        _depositQueue.Enqueue(deposit);
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            await Task.Delay(BatchDuration, stoppingToken);
            await ProcessDepositsAsync();
        }
    }

    private async Task ProcessDepositsAsync()
    {
        await _semaphore.WaitAsync();
        try
        {
            if (_depositQueue.IsEmpty)
                return;

            List<DepositModel> deposits = new List<DepositModel>();

            for (int i = 0; i < BatchSize && _depositQueue.TryDequeue(out var deposit); i++)
            {
                deposits.Add(deposit);
            }

            if (deposits.Count > 0)
            {
                using (var scope = _serviceProvider.CreateScope())
                {
                    var context = scope.ServiceProvider.GetRequiredService<CoreBankingDbContext>();

                    using (var transaction = await context.Database.BeginTransactionAsync())
                    {
                        try
                        {
                            context.DepositTable.AddRange(deposits);
                            await context.SaveChangesAsync();

                            List<TransactionHistoryModel> transactionHistories = deposits.Select(deposit => new TransactionHistoryModel
                            {
                                TransactionType = "Deposit",
                                RefID = deposit.DepositID,
                                AccountID = deposit.AccountID,
                                Amount = deposit.Amount,
                                TransactionDate = deposit.DepositDate,
                                CreatedAt = DateTime.UtcNow,
                                CreatedBy = deposit.CreatedBy,
                                Remarks = $"[{deposit.DepositDate}]{deposit.AccountID} deposit {deposit.Amount}. ",
                                StatusID = 1
                            }).ToList();

                            context.TransactionHistoryTable.AddRange(transactionHistories);
                            await context.SaveChangesAsync();

                            foreach (var deposit in deposits)
                            {
                                var bankAccount = await context.BankAccountTable.FindAsync(deposit.AccountID);
                                if (bankAccount != null)
                                {
                                    bankAccount.Balance += deposit.Amount;
                                    context.Entry(bankAccount).State = EntityState.Modified;
                                }
                            }

                            await context.SaveChangesAsync();
                            await transaction.CommitAsync();
                             _logger.LogInformation("Transaction committed successfully");
                        }
                        catch (Exception ex)
                        {
                            await transaction.RollbackAsync();
                              _logger.LogError(ex, "Error processing deposit batch. Transaction rolled back.");
                        }
                    }
                }
            }
        }
        finally
        {
            _semaphore.Release();
        }
    }
}
