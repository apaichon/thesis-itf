#nullable disable
public class TransactionHistoryModel
{
    public int TransactionID { get; set; }
    public int RefID { get; set; }
    public string TransactionType { get; set; }
    public int AccountID { get; set; }
    public double Amount { get; set; }
    public DateTime TransactionDate { get; set; }
    public DateTime CreatedAt { get; set; }
    public int CreatedBy { get; set; }
    public string Remarks { get; set; }
    public byte StatusID { get; set; }

}