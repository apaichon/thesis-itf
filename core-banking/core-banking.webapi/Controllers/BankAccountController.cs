
[ApiController]
[Route("api/[controller]")]
public class BankAccountController : ControllerBase
{
    private readonly CoreBankingDbContext _dbContext;

    public BankAccountController(CoreBankingDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    // GET: api/BankAccount
    [HttpGet]
    public async Task<ActionResult<IEnumerable<BankAccountModel>>> GetBankAccounts()
    {
        var bankAccounts = await _dbContext.BankAccountTable.ToListAsync();
        return Ok(bankAccounts);
    }

    // GET: api/BankAccount/5
    [HttpGet("{id}")]
    public async Task<ActionResult<BankAccountModel>> GetBankAccount(int id)
    {
        var bankAccount = await _dbContext.BankAccountTable.FindAsync(id);

        if (bankAccount == null)
        {
            return NotFound();
        }

        return Ok(bankAccount);
    }

    // POST: api/BankAccount
    [HttpPost]
    public async Task<ActionResult<BankAccountModel>> CreateBankAccount(BankAccountModel bankAccount)
    {
        _dbContext.BankAccountTable.Add(bankAccount);
        await _dbContext.SaveChangesAsync();

        return CreatedAtAction(nameof(GetBankAccount), new { id = bankAccount.AccountID }, bankAccount);
    }

    // PUT: api/BankAccount/5
    [HttpPut("{id}")]
    public async Task<IActionResult> UpdateBankAccount(int id, BankAccountModel bankAccount)
    {
        if (id != bankAccount.AccountID)
        {
            return BadRequest();
        }

        _dbContext.Entry(bankAccount).State = EntityState.Modified;

        try
        {
            await _dbContext.SaveChangesAsync();
        }
        catch (DbUpdateConcurrencyException)
        {
            if (!BankAccountExists(id))
            {
                return NotFound();
            }
            else
            {
                throw;
            }
        }

        return NoContent();
    }

    // DELETE: api/BankAccount/5
    [HttpDelete("{id}")]
    public async Task<IActionResult> DeleteBankAccount(int id)
    {
        var bankAccount = await _dbContext.BankAccountTable.FindAsync(id);
        if (bankAccount == null)
        {
            return NotFound();
        }

        _dbContext.BankAccountTable.Remove(bankAccount);
        await _dbContext.SaveChangesAsync();

        return NoContent();
    }

    private bool BankAccountExists(int id)
    {
        return _dbContext.BankAccountTable.Any(e => e.AccountID == id);
    }
}
