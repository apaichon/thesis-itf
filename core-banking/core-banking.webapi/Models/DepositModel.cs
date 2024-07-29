
#nullable disable
public class DepositModel
{
    public int DepositID { get; set; }
   
    public int AccountID { get; set; }
    public double Amount { get; set; }
    public DateTime DepositDate { get; set; }
    public int CreatedBy { get; set; }

}