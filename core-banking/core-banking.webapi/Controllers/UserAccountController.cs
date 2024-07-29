using _bcrypt = BCrypt.Net.BCrypt;

[Route("api/[controller]")]
[ApiController]

public class UserController : ControllerBase
{
    private readonly CoreBankingDbContext _context;

    public UserController(CoreBankingDbContext context)
    {
        _context = context;
    }

    // GET api/user: Get all users
    [HttpGet]
    public async Task<ActionResult<IEnumerable<UserAccountModel>>> GetUsers()
    {
        return await _context.UserAccountTable.ToListAsync();
    }

    // GET api/user/{id}: Get a specific user by ID
    [HttpGet("{id}")]
    public async Task<ActionResult<UserAccountModel>> GetUserById(int id)
    {
        var user = await _context.UserAccountTable.FindAsync(id);
        if (user == null)
        {
            return NotFound();
        }

        // Sensitive information like password hash and salt should not be exposed
        // Consider returning a user model with only necessary fields
        return user;
    }

    // POST api/user: Create a new user
    [HttpPost]
    public async Task<ActionResult<UserAccountModel>> CreateUser([FromBody] UserAccountModel user)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest(ModelState);
        }

        // Generate a strong, random salt
        var salt = _bcrypt.GenerateSalt(12);

        // Hash the password using BCryptNet instead of MD5
        user.PasswordHash = _bcrypt.HashPassword(user.PasswordHash, salt);

        // Remove plain text password before saving
        // user.PasswordHash = null;

        _context.UserAccountTable.Add(user);
        await _context.SaveChangesAsync();

        // Return user information without sensitive fields
        return CreatedAtAction(nameof(GetUserById), new { id = user.UserID }, user);
    }



    // PUT api/user/{id}: Update an existing user
    [HttpPut("{id}")]
    public async Task<IActionResult> UpdateUser(int id, [FromBody] UserAccountModel user)
    {
        if (id != user.UserID)
        {
            return BadRequest();
        }

        if (!ModelState.IsValid)
        {
            return BadRequest(ModelState);
        }

        var existingUser = await _context.UserAccountTable.FindAsync(id);
        if (existingUser == null)
        {
            return NotFound();
        }

        // Update properties as needed, but do not update password directly
        existingUser.Username = user.Username;
        existingUser.Email = user.Email;
        existingUser.PhoneNumber = user.PhoneNumber;
        existingUser.Remarks = user.Remarks;
        existingUser.StatusID = user.StatusID;
        existingUser.UpdatedAt = DateTime.UtcNow;

        // If password update is requested, handle it securely:
        if (!string.IsNullOrEmpty(user.PasswordHash))
        {
            // Re-hash the password with same or new salt based on security requirements

            existingUser.PasswordHash = _bcrypt.HashPassword(user.PasswordHash, existingUser.Salt);
        }
        else
        {
            // Generate a new salt and hash the password
            var newSalt = _bcrypt.GenerateSalt(12);
            existingUser.Salt = newSalt;
            existingUser.PasswordHash = _bcrypt.HashPassword(user.PasswordHash, newSalt);
        }

        await _context.SaveChangesAsync();

        return NoContent(); // 204 No Content response is appropriate for successful updates

    }

    [HttpDelete("{id}")]
    public async Task<IActionResult> DeleteUser(int id)
    {
        var user = await _context.UserAccountTable.FindAsync(id);
        if (user == null)
        {
            return NotFound();
        }

        _context.UserAccountTable.Remove(user);
        await _context.SaveChangesAsync();

        return NoContent(); // 204 No Content response is appropriate for successful deletions
    }
}

