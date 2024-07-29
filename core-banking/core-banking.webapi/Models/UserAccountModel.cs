#nullable disable
public class UserAccountModel
{
    public int UserID { get; set; }
    public string Username { get; set; }
    public string PasswordHash { get; set; } // Store hashed password instead of plain text
    public string Salt { get; set; }
    public string Email { get; set; }
    public string PhoneNumber { get; set; }
    public DateTime? CreatedAt { get; set; }
    public int? CreatedBy { get; set; }
    public DateTime? UpdatedAt { get; set; }
    public int? UpdatedBy { get; set; }
    public string Remarks { get; set; }
    public int StatusID { get; set; }
}
