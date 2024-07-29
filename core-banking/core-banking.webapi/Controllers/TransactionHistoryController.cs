[ApiController]
[Route("api/transactionhistory")]
public class TransactionHistoryController : ControllerBase
{
    private readonly CoreBankingDbContext _context;

    public TransactionHistoryController(CoreBankingDbContext context)
    {
        _context = context;
    }

    [HttpGet]
    public IActionResult GetAllTransactionHistory()
    {
        var transactionHistory = _context.TransactionHistoryTable.ToList();
        return Ok(transactionHistory);
    }

    [HttpGet("{id}")]
    public IActionResult GetTransactionHistoryById(int id)
    {
        var transaction = _context.TransactionHistoryTable.Find(id);
        if (transaction == null)
        {
            return NotFound();
        }

        return Ok(transaction);
    }

    [HttpPost]
    public IActionResult CreateTransactionHistory(TransactionHistoryModel transaction)
    {
        _context.TransactionHistoryTable.Add(transaction);
        _context.SaveChanges();

        return CreatedAtAction(nameof(GetTransactionHistoryById), new { id = transaction.TransactionID }, transaction);
    }

    [HttpDelete("{id}")]
    public IActionResult DeleteTransactionHistory(int id)
    {
        var transaction = _context.TransactionHistoryTable.Find(id);
        if (transaction == null)
        {
            return NotFound();
        }

        _context.TransactionHistoryTable.Remove(transaction);
        _context.SaveChanges();

        return NoContent();
    }
}
