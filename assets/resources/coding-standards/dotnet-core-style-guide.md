# .NET Core / C# Coding Standards

## 🎯 Overview

This document defines our team's .NET Core and C# coding standards. All C# code
must follow these guidelines to ensure consistency, maintainability, and code
quality across projects.

## 📦 Project Structure

```
MySolution/
├── src/
│   ├── MyApp.Api/                 # Web API project
│   │   ├── Controllers/
│   │   ├── Models/
│   │   ├── Program.cs
│   │   └── appsettings.json
│   ├── MyApp.Core/                # Domain/Business logic
│   │   ├── Entities/
│   │   ├── Interfaces/
│   │   ├── Services/
│   │   └── Exceptions/
│   ├── MyApp.Infrastructure/      # Data access, external services
│   │   ├── Data/
│   │   ├── Repositories/
│   │   └── ExternalServices/
│   └── MyApp.Shared/              # Shared utilities
│       ├── Constants/
│       ├── Extensions/
│       └── Helpers/
├── tests/
│   ├── MyApp.Tests.Unit/
│   ├── MyApp.Tests.Integration/
│   └── MyApp.Tests.E2E/
├── docs/
├── .editorconfig
├── Directory.Build.props
├── MySolution.sln
└── README.md
```

## 🎨 Code Style

### Formatting

- **Use `.editorconfig`** for consistent formatting across team
- **4 spaces for indentation** (not tabs)
- **Maximum line length: 120 characters** (soft limit)
- **One blank line between methods**
- **Organize usings**: System namespaces first, then third-party, then project

### Naming Conventions

**Classes & Interfaces:**

- `PascalCase` for all types
- Interfaces prefixed with `I`: `IUserRepository`, `IEmailService`
- Abstract classes: `UserServiceBase` or `AbstractUserService`
- Exception classes: Suffix with `Exception`: `UserNotFoundException`

**Methods & Properties:**

- `PascalCase` for all methods and properties
- Use descriptive verb phrases: `GetUserById()`, `ValidateEmail()`,
  `SaveChangesAsync()`
- Boolean properties: prefix with `Is`, `Has`, `Can`: `IsActive`,
  `HasPermission`, `CanEdit`

**Variables & Parameters:**

- `camelCase` for local variables and parameters
- Use descriptive names: `userCount` not `uc`
- Avoid Hungarian notation: use `user` not `objUser` or `strName`

**Constants & Fields:**

- `PascalCase` for `const` and `readonly` fields
- Private fields: `_camelCase` with underscore prefix
- Static readonly: `PascalCase`

**Example:**

```csharp
namespace MyApp.Core.Services
{
    using System;
    using System.Threading;
    using System.Threading.Tasks;
    using Microsoft.Extensions.Logging;
    
    using MyApp.Core.Entities;
    using MyApp.Core.Interfaces;
    using MyApp.Core.Exceptions;

    /// <summary>
    /// Handles user-related business operations.
    /// </summary>
    public class UserService : IUserService
    {
        private readonly IUserRepository _userRepository;
        private readonly ILogger<UserService> _logger;
        private const int MaxLoginAttempts = 5;

        public UserService(
            IUserRepository userRepository,
            ILogger<UserService> logger)
        {
            _userRepository = userRepository ?? throw new ArgumentNullException(nameof(userRepository));
            _logger = logger ?? throw new ArgumentNullException(nameof(logger));
        }

        /// <summary>
        /// Retrieves a user by their unique identifier.
        /// </summary>
        /// <param name="userId">The user's unique identifier.</param>
        /// <param name="cancellationToken">Cancellation token.</param>
        /// <returns>The user if found.</returns>
        /// <exception cref="UserNotFoundException">Thrown when user is not found.</exception>
        public async Task<User> GetUserByIdAsync(
            string userId,
            CancellationToken cancellationToken = default)
        {
            if (string.IsNullOrWhiteSpace(userId))
            {
                throw new ArgumentException("User ID cannot be null or empty.", nameof(userId));
            }

            var user = await _userRepository.GetByIdAsync(userId, cancellationToken);
            if (user == null)
            {
                _logger.LogWarning("User not found: {UserId}", userId);
                throw new UserNotFoundException($"User with ID '{userId}' was not found.");
            }

            return user;
        }
    }
}
```

## ⚙️ Language Features

### Async/Await

**Rules:**

1. **Always use `async`/`await`** for I/O operations (database, HTTP, file
   system)
2. **Use `ConfigureAwait(false)`** in library code (not necessary in ASP.NET
   Core)
3. **Suffix async methods with `Async`**: `GetUserAsync()`, `SaveAsync()`
4. **Pass `CancellationToken`** to all async methods
5. **Don't use `async void`** except for event handlers

**Good Example:**

```csharp
// ✅ DO: Async with CancellationToken
public async Task<User> GetUserAsync(string id, CancellationToken cancellationToken = default)
{
    var user = await _context.Users
        .FirstOrDefaultAsync(u => u.Id == id, cancellationToken);
    
    return user;
}

// ✅ DO: Propagate async all the way up
public async Task<IActionResult> GetUser(string id)
{
    var user = await _userService.GetUserAsync(id);
    return Ok(user);
}

// ✅ DO: Use Task.WhenAll for parallel operations
public async Task<(User user, IEnumerable<Order> orders)> GetUserWithOrdersAsync(string userId)
{
    var userTask = _userService.GetUserAsync(userId);
    var ordersTask = _orderService.GetOrdersByUserAsync(userId);
    
    await Task.WhenAll(userTask, ordersTask);
    
    return (await userTask, await ordersTask);
}
```

**Bad Example:**

```csharp
// ❌ DON'T: Blocking async code
public User GetUser(string id)
{
    return _userService.GetUserAsync(id).Result; // Can cause deadlocks
}

// ❌ DON'T: Async void (except event handlers)
public async void SaveUser(User user) // No way to handle exceptions
{
    await _userRepository.SaveAsync(user);
}

// ❌ DON'T: Unnecessary async
public async Task<int> GetCount()
{
    return await Task.FromResult(5); // Just return 5 directly
}
```

### Null Handling

Use **nullable reference types** (enabled by default in .NET 6+):

```csharp
// ✅ DO: Enable nullable reference types in .csproj
<Nullable>enable</Nullable>

// ✅ DO: Use ? for nullable types
public class User
{
    public string Id { get; set; } = string.Empty; // Non-nullable
    public string? MiddleName { get; set; }        // Nullable
    public Address Address { get; set; } = null!;  // Non-nullable, set later
}

// ✅ DO: Null checking
public void ProcessUser(User? user)
{
    if (user is null)
    {
        throw new ArgumentNullException(nameof(user));
    }
    
    // Use pattern matching
    if (user.Address is { City: "Seattle" })
    {
        // Process Seattle users
    }
}

// ✅ DO: Null-coalescing
var name = user.MiddleName ?? "N/A";
var email = user?.Contact?.Email ?? "unknown@example.com";
```

### LINQ & Collections

```csharp
// ✅ DO: Use LINQ for readability
var activeUsers = users
    .Where(u => u.IsActive)
    .OrderBy(u => u.LastName)
    .ThenBy(u => u.FirstName)
    .ToList();

// ✅ DO: Use async LINQ with Entity Framework
var users = await _context.Users
    .Where(u => u.IsActive)
    .ToListAsync(cancellationToken);

// ✅ DO: Use appropriate collection types
IEnumerable<User> GetUsers();           // For iteration only
IReadOnlyList<User> GetUserList();      // For indexed access, no modification
List<User> GetMutableUsers();           // For modification

// ❌ DON'T: Multiple enumeration
var users = GetUsers().Where(u => u.IsActive); // IEnumerable
var count = users.Count(); // Enumerates
var first = users.First(); // Enumerates again!

// ✅ DO: Materialize when needed
var users = GetUsers().Where(u => u.IsActive).ToList();
var count = users.Count;
var first = users[0];
```

### Exception Handling

```csharp
// ✅ DO: Use specific exception types
public async Task<User> GetUserAsync(string id)
{
    if (string.IsNullOrWhiteSpace(id))
    {
        throw new ArgumentException("ID cannot be null or empty.", nameof(id));
    }
    
    var user = await _repository.GetByIdAsync(id);
    if (user == null)
    {
        throw new UserNotFoundException($"User '{id}' not found.");
    }
    
    return user;
}

// ✅ DO: Catch specific exceptions
try
{
    await ProcessPaymentAsync(payment);
}
catch (PaymentDeclinedException ex)
{
    _logger.LogWarning(ex, "Payment declined for user {UserId}", payment.UserId);
    return PaymentResult.Declined(ex.Message);
}
catch (PaymentGatewayException ex)
{
    _logger.LogError(ex, "Payment gateway error");
    throw; // Re-throw infrastructure errors
}

// ❌ DON'T: Catch and ignore
try
{
    await SaveAsync(user);
}
catch { } // ❌ Never do this

// ❌ DON'T: Catch generic Exception unless at boundary
catch (Exception ex) // ❌ Too broad
{
    // ...
}
```

### Dependency Injection

```csharp
// ✅ DO: Constructor injection
public class UserController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly ILogger<UserController> _logger;

    public UserController(
        IUserService userService,
        ILogger<UserController> logger)
    {
        _userService = userService ?? throw new ArgumentNullException(nameof(userService));
        _logger = logger ?? throw new ArgumentNullException(nameof(logger));
    }
}

// ✅ DO: Register services properly in Program.cs
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddSingleton<ICacheService, RedisCacheService>();
builder.Services.AddTransient<IEmailSender, EmailSender>();

// Service lifetimes:
// - Singleton: One instance for application lifetime (caches, config)
// - Scoped: One instance per request (repositories, services, DbContext)
// - Transient: New instance every time (lightweight, stateless services)
```

## 🗄️ Entity Framework Core

### DbContext

```csharp
public class ApplicationDbContext : DbContext
{
    public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
        : base(options)
    {
    }

    public DbSet<User> Users => Set<User>();
    public DbSet<Order> Orders => Set<Order>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
        
        // ✅ DO: Configure entities using Fluent API
        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Email).IsRequired().HasMaxLength(256);
            entity.HasIndex(e => e.Email).IsUnique();
            
            entity.HasMany(e => e.Orders)
                  .WithOne(e => e.User)
                  .HasForeignKey(e => e.UserId);
        });
    }
}
```

### Repository Pattern

```csharp
// ✅ DO: Use repository pattern for data access
public interface IUserRepository
{
    Task<User?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<User?> GetByEmailAsync(string email, CancellationToken cancellationToken = default);
    Task<IEnumerable<User>> GetAllActiveAsync(CancellationToken cancellationToken = default);
    Task AddAsync(User user, CancellationToken cancellationToken = default);
    Task UpdateAsync(User user, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class UserRepository : IUserRepository
{
    private readonly ApplicationDbContext _context;

    public UserRepository(ApplicationDbContext context)
    {
        _context = context ?? throw new ArgumentNullException(nameof(context));
    }

    public async Task<User?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Users
            .Include(u => u.Profile)
            .FirstOrDefaultAsync(u => u.Id == id, cancellationToken);
    }

    public async Task AddAsync(User user, CancellationToken cancellationToken = default)
    {
        await _context.Users.AddAsync(user, cancellationToken);
        await _context.SaveChangesAsync(cancellationToken);
    }
}
```

### Query Best Practices

```csharp
// ✅ DO: Use AsNoTracking for read-only queries
var users = await _context.Users
    .AsNoTracking()
    .Where(u => u.IsActive)
    .ToListAsync();

// ✅ DO: Project to DTOs to avoid loading entire entities
var userDtos = await _context.Users
    .Where(u => u.IsActive)
    .Select(u => new UserDto
    {
        Id = u.Id,
        Name = u.Name,
        Email = u.Email
    })
    .ToListAsync();

// ✅ DO: Use Include for eager loading related data
var user = await _context.Users
    .Include(u => u.Orders)
    .Include(u => u.Profile)
    .FirstOrDefaultAsync(u => u.Id == id);

// ❌ DON'T: Load all data then filter in memory
var users = await _context.Users.ToListAsync(); // Loads everything!
var activeUsers = users.Where(u => u.IsActive).ToList();

// ✅ DO: Filter at database level
var activeUsers = await _context.Users
    .Where(u => u.IsActive)
    .ToListAsync();
```

## 🧪 Testing

### Unit Tests (xUnit)

```csharp
public class UserServiceTests
{
    private readonly Mock<IUserRepository> _mockRepository;
    private readonly Mock<ILogger<UserService>> _mockLogger;
    private readonly UserService _service;

    public UserServiceTests()
    {
        _mockRepository = new Mock<IUserRepository>();
        _mockLogger = new Mock<ILogger<UserService>>();
        _service = new UserService(_mockRepository.Object, _mockLogger.Object);
    }

    [Fact]
    public async Task GetUserByIdAsync_WithValidId_ReturnsUser()
    {
        // Arrange
        var userId = "123";
        var expectedUser = new User { Id = userId, Email = "test@example.com" };
        _mockRepository
            .Setup(r => r.GetByIdAsync(userId, It.IsAny<CancellationToken>()))
            .ReturnsAsync(expectedUser);

        // Act
        var result = await _service.GetUserByIdAsync(userId);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(userId, result.Id);
        Assert.Equal("test@example.com", result.Email);
    }

    [Theory]
    [InlineData("")]
    [InlineData(null)]
    [InlineData("   ")]
    public async Task GetUserByIdAsync_WithInvalidId_ThrowsArgumentException(string invalidId)
    {
        // Act & Assert
        await Assert.ThrowsAsync<ArgumentException>(() =>
            _service.GetUserByIdAsync(invalidId));
    }

    [Fact]
    public async Task GetUserByIdAsync_UserNotFound_ThrowsUserNotFoundException()
    {
        // Arrange
        _mockRepository
            .Setup(r => r.GetByIdAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync((User?)null);

        // Act & Assert
        await Assert.ThrowsAsync<UserNotFoundException>(() =>
            _service.GetUserByIdAsync("999"));
    }
}
```

### Integration Tests

```csharp
public class UserRepositoryIntegrationTests : IClassFixture<DatabaseFixture>
{
    private readonly DatabaseFixture _fixture;

    public UserRepositoryIntegrationTests(DatabaseFixture fixture)
    {
        _fixture = fixture;
    }

    [Fact]
    public async Task AddAsync_SavesUserToDatabase()
    {
        // Arrange
        using var context = _fixture.CreateContext();
        var repository = new UserRepository(context);
        var user = new User { Id = Guid.NewGuid().ToString(), Email = "test@example.com" };

        // Act
        await repository.AddAsync(user);

        // Assert
        var savedUser = await repository.GetByIdAsync(user.Id);
        Assert.NotNull(savedUser);
        Assert.Equal(user.Email, savedUser.Email);
    }
}
```

## 🔐 Security Best Practices

```csharp
// ✅ DO: Validate and sanitize input
public async Task<User> CreateUserAsync(CreateUserRequest request)
{
    if (string.IsNullOrWhiteSpace(request.Email))
    {
        throw new ValidationException("Email is required.");
    }
    
    if (!IsValidEmail(request.Email))
    {
        throw new ValidationException("Invalid email format.");
    }
    
    // Check for SQL injection patterns if accepting raw SQL
    // Use parameterized queries or EF Core methods
}

// ✅ DO: Use SecureString for sensitive data or clear from memory
// ✅ DO: Hash passwords properly
using var hasher = new PasswordHasher<User>();
var hashedPassword = hasher.HashPassword(user, plainTextPassword);

// ✅ DO: Use HTTPS and validate certificates
builder.Services.AddHttpClient<IExternalService, ExternalService>(client =>
{
    client.BaseAddress = new Uri("https://api.example.com");
});

// ✅ DO: Implement authorization
[Authorize(Roles = "Admin")]
public class AdminController : ControllerBase
{
    [Authorize(Policy = "CanEditUsers")]
    public async Task<IActionResult> UpdateUser(string id, UpdateUserRequest request)
    {
        // Implementation
    }
}

// ✅ DO: Never log sensitive information
_logger.LogInformation("User login attempt: {UserId}", userId); // ✅
_logger.LogInformation("User login: {@User}", user); // ❌ May contain password
```

## 📝 Documentation

### XML Documentation Comments

```csharp
/// <summary>
/// Retrieves a user by their unique identifier.
/// </summary>
/// <param name="userId">The user's unique identifier.</param>
/// <param name="cancellationToken">Optional cancellation token.</param>
/// <returns>A task representing the asynchronous operation, containing the user if found.</returns>
/// <exception cref="ArgumentException">Thrown when userId is null or empty.</exception>
/// <exception cref="UserNotFoundException">Thrown when the user is not found.</exception>
public async Task<User> GetUserByIdAsync(
    string userId,
    CancellationToken cancellationToken = default)
{
    // Implementation
}
```

## 🚫 Common Anti-Patterns to Avoid

### ❌ DON'T:

```csharp
// Catch and rethrow without adding value
try
{
    await DoSomethingAsync();
}
catch (Exception ex)
{
    throw ex; // ❌ Loses stack trace, use 'throw;'
}

// Using async void (except event handlers)
public async void ProcessData() // ❌
{
    await _service.ProcessAsync();
}

// Blocking on async code
var result = _service.GetDataAsync().Result; // ❌ Can deadlock

// Not disposing IDisposable
var stream = File.OpenRead("file.txt"); // ❌ No using statement

// Over-using dynamic
dynamic user = GetUser(); // ❌ Loses type safety

// Ignoring CancellationToken
public async Task ProcessAsync(CancellationToken cancellationToken)
{
    await Task.Delay(1000); // ❌ Should pass cancellationToken
}
```

### ✅ DO:

```csharp
// Re-throw to preserve stack trace
try
{
    await DoSomethingAsync();
}
catch (Exception)
{
    // Log or handle
    throw; // ✅ Preserves stack trace
}

// Use async Task
public async Task ProcessDataAsync()
{
    await _service.ProcessAsync();
}

// Always await async code
var result = await _service.GetDataAsync();

// Use using statements or declarations
using var stream = File.OpenRead("file.txt"); // ✅

// Use strong typing
User user = GetUser(); // ✅

// Always pass CancellationToken
public async Task ProcessAsync(CancellationToken cancellationToken)
{
    await Task.Delay(1000, cancellationToken); // ✅
}
```

## 🛠️ Tools & Configuration

### Required Tools

```bash
dotnet tool install --global dotnet-ef
dotnet tool install --global dotnet-format
dotnet tool install --global dotnet-outdated-tool
```

### .editorconfig

Create `.editorconfig` in solution root:

```ini
root = true

[*.cs]
indent_style = space
indent_size = 4
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

# Code style rules
csharp_prefer_braces = true:warning
csharp_prefer_simple_using_statement = true:suggestion
csharp_style_namespace_declarations = file_scoped:warning
csharp_style_var_for_built_in_types = false:suggestion
csharp_style_var_when_type_is_apparent = true:suggestion
```

## 📚 Additional Resources

- [Microsoft C# Coding Conventions](https://learn.microsoft.com/en-us/dotnet/csharp/fundamentals/coding-style/coding-conventions)
- [.NET API Design Guidelines](https://learn.microsoft.com/en-us/dotnet/standard/design-guidelines/)
- [Entity Framework Core Best Practices](https://learn.microsoft.com/en-us/ef/core/)
- [Async/Await Best Practices](https://learn.microsoft.com/en-us/archive/msdn-magazine/2013/march/async-await-best-practices-in-asynchronous-programming)

---

**Last Updated:** 2025-12-23\
**Version:** 1.0.0
