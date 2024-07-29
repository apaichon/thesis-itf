[ApiController]
[Route("api/UserProfile")]
public class UserProfileController : ControllerBase
{
    private readonly CoreBankingDbContext _context;

    public UserProfileController(CoreBankingDbContext context)
    {
        _context = context;
    }

    [HttpGet]
    public IActionResult GetAllUserProfileTable()
    {
        var UserProfileTable = _context.UserProfileTable.ToList();
        return Ok(UserProfileTable);
    }

    [HttpGet("{id}")]
    public IActionResult GetUserProfileById(int id)
    {
        var userProfile = _context.UserProfileTable.Find(id);
        if (userProfile == null)
        {
            return NotFound();
        }

        return Ok(userProfile);
    }

    [HttpPost]
    public IActionResult CreateUserProfile(UserProfileModel userProfile)
    {
        _context.UserProfileTable.Add(userProfile);
        _context.SaveChanges();

        return CreatedAtAction(nameof(GetUserProfileById), new { id = userProfile.UserID }, userProfile);
    }

    [HttpPut("{id}")]
    public IActionResult UpdateUserProfile(int id, UserProfileModel updatedUserProfile)
    {
        var userProfile = _context.UserProfileTable.Find(id);
        if (userProfile == null)
        {
            return NotFound();
        }

        userProfile.FirstName = updatedUserProfile.FirstName;
        // Update other properties...

        _context.SaveChanges();

        return NoContent();
    }

    [HttpDelete("{id}")]
    public IActionResult DeleteUserProfile(int id)
    {
        var userProfile = _context.UserProfileTable.Find(id);
        if (userProfile == null)
        {
            return NotFound();
        }

        _context.UserProfileTable.Remove(userProfile);
        _context.SaveChanges();

        return NoContent();
    }
}
