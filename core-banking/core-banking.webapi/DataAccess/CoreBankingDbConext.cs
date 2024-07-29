using Microsoft.EntityFrameworkCore.Sqlite;

public class CoreBankingDbContext : DbContext
{
    public CoreBankingDbContext(DbContextOptions<CoreBankingDbContext> options) : base(options)
    {
    }

   protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        IConfiguration configuration = new ConfigurationBuilder()
            .SetBasePath(Directory.GetCurrentDirectory())
            .AddJsonFile("appsettings.json")
            .Build();

        optionsBuilder.UseSqlite(configuration.GetConnectionString("DefaultConnection"));
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<UserAccountModel>(entity =>
        {
            entity.ToTable("UserAccount"); // Set the table name

            // Primary Key
            entity.HasKey(e => e.UserID);

            // Required fields
            entity.Property(e => e.Username).IsRequired().HasMaxLength(255);
            entity.Property(e => e.PasswordHash).IsRequired().HasMaxLength(255);
            entity.Property(e => e.Salt).IsRequired().HasMaxLength(255);

            // Optional fields
            entity.Property(e => e.Email).HasMaxLength(255);
            entity.Property(e => e.PhoneNumber).HasMaxLength(20);
            entity.Property(e => e.Remarks).HasMaxLength(500);

            // DateTime fields
            entity.Property(e => e.CreatedAt).HasColumnType("DATETIME");
            entity.Property(e => e.UpdatedAt).HasColumnType("DATETIME");

        });

         modelBuilder.Entity<UserProfileModel>(entity =>
        {
            entity.ToTable("UserProfile");

            entity.HasKey(e => e.UserID);
            entity.Property(e => e.FirstName).HasMaxLength(255);
            entity.Property(e => e.LastName).HasMaxLength(255);
            entity.Property(e => e.DateOfBirth).HasColumnType("DATE");
            entity.Property(e => e.Address).HasMaxLength(500);
            entity.Property(e => e.Remarks).HasMaxLength(500);
            entity.Property(e => e.CreatedAt).HasColumnType("DATETIME");
            entity.Property(e => e.UpdatedAt).HasColumnType("DATETIME");

            // Define relationships if needed
            // entity.HasOne(e => e.UserAccount).WithOne(u => u.UserProfile).HasForeignKey<UserProfile>(u => u.UserID);
        });


         // BankAccount entity mapping
        modelBuilder.Entity<BankAccountModel>(entity =>
        {
            entity.ToTable("BankAccount");

             entity.HasKey(e => e.AccountID);

            entity.Property(e => e.AccountNumber).IsRequired().HasMaxLength(255);
            entity.Property(e => e.Balance).IsRequired();
            entity.Property(e => e.AccountTypeID);
            entity.Property(e => e.Remarks).HasMaxLength(500);
            
            entity.Property(e => e.CreatedAt).HasColumnType("DATETIME");
            entity.Property(e => e.UpdatedAt).HasColumnType("DATETIME");

            entity.Property(e => e.CreatedBy);
            entity.Property(e => e.UpdatedBy);
            entity.Property(e => e.StatusID);

            // Define relationships if needed
            // entity.HasOne(e => e.UserAccount).WithMany(u => u.BankAccounts).HasForeignKey(e => e.UserID);
        });


         // Deposit entity mapping
        modelBuilder.Entity<DepositModel>(entity =>
        {
            entity.ToTable("Deposit");
            entity.HasKey(e => e.DepositID);

            entity.Property(e => e.Amount).IsRequired();
            entity.Property(e => e.DepositDate).IsRequired().HasColumnType("DATETIME");


            // Define relationships if needed
            // entity.HasOne(e => e.BankAccount).WithMany(b => b.Deposits).HasForeignKey(e => e.AccountID);
        });

        modelBuilder.Entity<TransactionHistoryModel>(entity =>
        {
            entity.ToTable("TransactionHistory"); // Assuming your table name is "TransactionHistory"

            entity.HasKey(e => e.TransactionID);

            entity.Property(e => e.TransactionType).IsRequired().HasMaxLength(255);
            entity.Property(e => e.Amount).IsRequired();
            entity.Property(e => e.TransactionDate).IsRequired().HasColumnType("DATETIME");
            entity.Property(e => e.CreatedAt).IsRequired().HasColumnType("DATETIME");
            entity.Property(e => e.Remarks).HasMaxLength(500);

            // If you have relationships, define them here
            // entity.HasOne(e => e.BankAccount).WithMany(b => b.TransactionHistory).HasForeignKey(e => e.AccountID);
        });

         modelBuilder.Entity<WithdrawalModel>(entity =>
        {
            entity.ToTable("Withdrawal"); // Assuming your table name is "Withdrawal"

            entity.HasKey(e => e.WithdrawalID);

            entity.Property(e => e.Amount).IsRequired();
            entity.Property(e => e.WithdrawalDate).IsRequired().HasColumnType("DATETIME");

            // If you have relationships, define them here
            // entity.HasOne(e => e.BankAccount).WithMany(b => b.Withdrawals).HasForeignKey(e => e.AccountID);
        });

        // Additional configurations for other entities...

        base.OnModelCreating(modelBuilder);
    }

    public DbSet<UserAccountModel> UserAccountTable { get; set; }
    public DbSet<UserProfileModel> UserProfileTable { get; set; }
    public DbSet<BankAccountModel> BankAccountTable { get; set; }
    public DbSet<DepositModel> DepositTable { get; set; }
    public DbSet<TransactionHistoryModel> TransactionHistoryTable { get; set; }
    public DbSet<WithdrawalModel> WithdrawalTable { get; set; }

}