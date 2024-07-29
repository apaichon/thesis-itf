[ApiController]
[Route("api/withdrawals")]
public class WithdrawalController : ControllerBase
{
    private readonly CoreBankingDbContext _context;

    public WithdrawalController(CoreBankingDbContext context)
    {
        _context = context;
    }

    [HttpGet]
    public IActionResult GetAllWithdrawals()
    {
        var withdrawals = _context.WithdrawalTable.ToList();
        return Ok(withdrawals);
    }

    [HttpGet("{id}")]
    public IActionResult GetWithdrawalById(int id)
    {
        var withdrawal = _context.WithdrawalTable.Find(id);
        if (withdrawal == null)
        {
            return NotFound();
        }

        return Ok(withdrawal);
    }


[HttpPost]
    public async Task<IActionResult> AddWithdrawal([FromBody] WithdrawalModel withdrawal)
    {
        using (var transaction = _context.Database.BeginTransaction())
        {
            try
            {

                // Step 1: Check balance and transfer amount data
                var withdrawalAccount = await _context.BankAccountTable
                    .FirstOrDefaultAsync(a => a.AccountID == withdrawal.AccountID);

                if (withdrawalAccount == null || withdrawal.Amount > withdrawalAccount.Balance || withdrawal.Amount <= 0)
                {
                    // Insufficient balance or invalid withdrawal amount
                    return BadRequest("Invalid withdrawal amount or insufficient balance.");
                }
                // Step 2: Insert withdrawal data
                if (withdrawal.Amount > 0)
                {
                    _context.WithdrawalTable.Add(withdrawal);
                    await _context.SaveChangesAsync();
                }

                // Step 3: Insert deposit transaction to TransactionHistory table
                var transactionHistory = new TransactionHistoryModel
                {
                    TransactionType = "Withdrawal",
                    RefID = withdrawal.WithdrawalID,
                    AccountID = withdrawal.AccountID,
                    Amount = withdrawal.Amount * -1,
                    TransactionDate = withdrawal.WithdrawalDate,
                    CreatedAt = DateTime.UtcNow,
                    CreatedBy = withdrawal.CreatedBy,
                    Remarks = $"[{withdrawal.WithdrawalDate}]{withdrawal.AccountID} withdraw {withdrawal.Amount}",
                    StatusID = 1 // Assuming 1 represents a successful transaction
                };
                _context.TransactionHistoryTable.Add(transactionHistory);
                await _context.SaveChangesAsync();

                // Step 4: Upsert balance to BankAccount table
                var bankAccount = await _context.BankAccountTable.FindAsync(withdrawal.AccountID);
                if (bankAccount != null)
                {
                    bankAccount.Balance -= withdrawal.Amount;
                    _context.Entry(bankAccount).State = EntityState.Modified;
                    await _context.SaveChangesAsync();
                }

                transaction.Commit();
                return Ok("Withdrawal added successfully.");
            }
            catch (Exception ex)
            {
                transaction.Rollback();
                return BadRequest($"Error: {ex.Message}");
            }
        }
    }

    [HttpDelete("{id}")]
    public IActionResult DeleteWithdrawal(int id)
    {
        var withdrawal = _context.WithdrawalTable.Find(id);
        if (withdrawal == null)
        {
            return NotFound();
        }

        _context.WithdrawalTable.Remove(withdrawal);
        _context.SaveChanges();

        return NoContent();
    }
}
