public class WithdrawalModel
{
    public int WithdrawalID { get; set; }
    public int AccountID { get; set; }
    public double Amount { get; set; }
    public DateTime WithdrawalDate { get; set; }
    public int CreatedBy { get; set; }

}