/* using Microsoft.Extensions.DependencyInjection;
namespace core_banking.webapi;
public class Startup
{
    public IConfiguration Configuration { get; }

    public Startup(IConfiguration configuration)
    {
        Configuration = configuration;
    }

    // Other methods...

    public void ConfigureServices(IServiceCollection services)
    {
        // services.AddSingleton(Configuration); // Register IConfiguration as a service if needed
        var appSettings = Configuration.GetSection("AppSettings").Get<AppSettings>();
        services.AddSingleton(appSettings);
         services.AddControllers()
                .AddApplicationPart(typeof(UserAccountController).Assembly);

                 IConfiguration configuration = new ConfigurationBuilder()
        .SetBasePath(Directory.GetCurrentDirectory())
        .AddJsonFile("appsettings.json")
        .Build();

    optionsBuilder.UseSqlServer(configuration.GetConnectionString("MyConnectionString"));
               //  .AddApplicationPart(typeof(ProductController).Assembly);
        // Additional configuration setup...
    }
}
*/
