[ApiController]
[Route("api/deposits")]
public class DepositController : ControllerBase
{
    private readonly CoreBankingDbContext _context;
    private readonly ILogger<DepositController> _logger;

    public DepositController(CoreBankingDbContext context, ILogger<DepositController> logger)
    {
        _context = context;
        _logger = logger;
    }

    [HttpGet]
    public IActionResult GetAllDeposits()
    {
        var deposits = _context.DepositTable.ToList();
        return Ok(deposits);
    }

    [HttpGet("{id}")]
    public IActionResult GetDepositById(int id)
    {
        var deposit = _context.DepositTable.Find(id);
        if (deposit == null)
        {
            return NotFound();
        }

        return Ok(deposit);
    }

   [HttpPost("depositBatch")]
    public async Task<IActionResult> AddDepositBatch([FromBody] DepositModel deposit)
    {
        try
        {
            DepositBatchService.EnqueueDeposit(deposit);
            _logger.LogInformation("Deposit request for AccountID: {AccountID} has been enqueued successfully.", deposit.AccountID);
            return Ok("Deposit received successfully.");
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error occurred while enqueuing deposit request for AccountID: {AccountID}", deposit.AccountID);
            return BadRequest($"Error: {ex.Message}");
        }
    }


    [HttpPost]
    public async Task<IActionResult> AddDeposit([FromBody] DepositModel deposit)
    {
        using (var transaction = _context.Database.BeginTransaction())
        {
            try
            {
                // Step 2: Insert deposit data
                if (deposit.Amount > 0)
                {
                    _context.DepositTable.Add(deposit);
                    await _context.SaveChangesAsync();
                }

                // Step 3: Insert deposit transaction to TransactionHistory table
                var transactionHistory = new TransactionHistoryModel
                {
                    TransactionType = "Deposit",
                    RefID = deposit.DepositID ,
                    AccountID = deposit.AccountID,
                    Amount = deposit.Amount,
                    TransactionDate = deposit.DepositDate,
                    CreatedAt = DateTime.UtcNow,
                    CreatedBy = deposit.CreatedBy,
                    Remarks = $"[{deposit.DepositDate}]{deposit.AccountID} deposit {deposit.Amount}. ",
                    StatusID = 1 // Assuming 1 represents a successful transaction
                };
            
                _context.TransactionHistoryTable.Add(transactionHistory);
                await _context.SaveChangesAsync();

                // Step 4: Upsert balance to BankAccount table
                var bankAccount = await _context.BankAccountTable.FindAsync(deposit.AccountID);
                if (bankAccount != null)
                {
                    bankAccount.Balance += deposit.Amount;
                    _context.Entry(bankAccount).State = EntityState.Modified;
                    await _context.SaveChangesAsync();
                }

                transaction.Commit();
                return Ok("Deposit added successfully.");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return BadRequest($"Error: {ex.Message}");
            }
        }
    }

    [HttpDelete("{id}")]
    public IActionResult DeleteDeposit(int id)
    {
        var deposit = _context.DepositTable.Find(id);
        if (deposit == null)
        {
            return NotFound();
        }

        _context.DepositTable.Remove(deposit);
        _context.SaveChanges();

        return NoContent();
    }
}
