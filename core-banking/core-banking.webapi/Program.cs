var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddControllers()
.AddNewtonsoftJson(options =>
            {
                options.SerializerSettings.ReferenceLoopHandling = ReferenceLoopHandling.Ignore; // Avoid reference loop issues
                options.SerializerSettings.NullValueHandling = NullValueHandling.Ignore; // Ignore null values in JSON
            });
builder.Services.AddMvc();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Configure configuration
var configuration = new ConfigurationBuilder()
    .SetBasePath(Directory.GetCurrentDirectory())
    .AddJsonFile("appsettings.json", optional: true)
    .AddEnvironmentVariables() // Read environment variables
    .Build();

// Read the connection string from configuration
var connectionString = configuration.GetConnectionString("DefaultConnection");

// Add DbContext services
builder.Services.AddDbContext<CoreBankingDbContext>(options =>
{
    options.UseSqlite(connectionString);
});

builder.Services.AddHostedService<DepositBatchService>();



var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}


app.UseHttpsRedirection();
app.UseRouting();

app.UseEndpoints(endpoints =>
{
    endpoints.MapControllers();
});



/*
  app.MapControllers(); // Top-level route registration

    app.UseEndpoints(endpoints =>
    {
        endpoints.MapControllerRoute(
            name: "default",
            pattern: "{controller=Home}/{action=Index}/{id?}");
    });
*/


app.Run();



