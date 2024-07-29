
#nullable disable
public class BankAccountModel
{
    public int AccountID { get; set; }
    public int UserID { get; set; }
    public string AccountNumber { get; set; }
    public double Balance { get; set; }
    public short? AccountTypeID { get; set; }
    public DateTime CreatedAt { get; set; }
    public int CreatedBy { get; set; }
    public DateTime? UpdatedAt { get; set; }
    public int? UpdatedBy { get; set; }
    public string Remarks { get; set; }
    public byte StatusID { get; set; }

}