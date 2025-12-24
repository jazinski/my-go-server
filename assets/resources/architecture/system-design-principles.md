# Software Architecture Principles

## 🎯 Overview

This document outlines our team's core software architecture principles. These
principles guide all architectural decisions and system design across our .NET
Core, React, and legacy application projects.

## 🏛️ Core Principles

### 1. **Simplicity First**

> "The simplest solution is usually the best solution."

- **Start simple:** Don't over-engineer for future requirements
- **YAGNI principle:** You Aren't Gonna Need It
- **Delete code:** The best code is code that doesn't exist
- **Boring technology:** Prefer proven solutions over shiny new tech

**Example (.NET Core):**

```csharp
// ❌ Over-engineered
public interface IAbstractUserFactoryProviderInterface
{
    IUserFactory CreateUserFactory();
}

public class ConcreteUserFactoryProvider : IAbstractUserFactoryProviderInterface
{
    // Too many layers of abstraction
}

// ✅ Simple and direct
public class User
{
    public User(string firstName, string email)
    {
        FirstName = firstName;
        Email = email;
    }
    
    public string FirstName { get; }
    public string Email { get; }
}
```

**Example (React):**

```typescript
// ❌ Over-engineered
class AbstractComponentFactoryProvider {
  createComponentFactory(): ComponentFactory {}
}

// ✅ Simple and direct
function UserCard({ name, email }: UserCardProps) {
  return (
    <div className="user-card">
      <h3>{name}</h3>
      <p>{email}</p>
    </div>
  );
}
```

### 2. **Separation of Concerns**

Each component should have a single, well-defined responsibility.

**Layers (.NET Core):**

```
┌─────────────────────────────────────┐
│  Presentation (Controllers/React)  │  ← User interaction
├─────────────────────────────────────┤
│  Application (Services/Use Cases)  │  ← Application logic
├─────────────────────────────────────┤
│  Domain (Business Logic/Entities)  │  ← Core business rules
├─────────────────────────────────────┤
│  Infrastructure (EF Core/Data)     │  ← Database/External APIs
└─────────────────────────────────────┘
```

**Example (.NET Core API):**

```csharp
// ✅ Separated concerns

// Controller - handles HTTP
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;
    
    public UsersController(IUserService userService)
    {
        _userService = userService;
    }
    
    [HttpGet("{id}")]
    public async Task<ActionResult<UserResponse>> GetUser(int id)
    {
        var user = await _userService.GetUserAsync(id);
        if (user == null)
            return NotFound();
        
        return Ok(user);
    }
}

// Service - business logic
public class UserService : IUserService
{
    private readonly IUserRepository _repository;
    
    public UserService(IUserRepository repository)
    {
        _repository = repository;
    }
    
    public async Task<User?> GetUserAsync(int id)
    {
        var user = await _repository.FindByIdAsync(id);
        // Apply business rules...
        return user;
    }
}

// Repository - data access
public class UserRepository : IUserRepository
{
    private readonly AppDbContext _context;
    
    public UserRepository(AppDbContext context)
    {
        _context = context;
    }
    
    public async Task<User?> FindByIdAsync(int id)
    {
        return await _context.Users
            .AsNoTracking()
            .FirstOrDefaultAsync(u => u.Id == id);
    }
}
```

**Example (React):**

```typescript
// ✅ Separated concerns

// Component - presentation only
function UserProfile({ userId }: UserProfileProps) {
  const { user, loading, error } = useUser(userId);

  if (loading) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;
  if (!user) return <NotFound />;

  return (
    <div className="user-profile">
      <h1>{user.name}</h1>
      <p>{user.email}</p>
    </div>
  );
}

// Custom hook - data fetching logic
function useUser(userId: number) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    userService.getUser(userId)
      .then(setUser)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [userId]);

  return { user, loading, error };
}

// Service - API calls
const userService = {
  async getUser(id: number): Promise<User> {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) throw new Error("Failed to fetch user");
    return response.json();
  },
};
```

### 3. **Dependency Injection**

Depend on abstractions (interfaces), not concrete implementations.

**Benefits:**

- ✅ Easier testing (mock dependencies)
- ✅ Loose coupling
- ✅ Flexible configurations
- ✅ Better testability

**Example (.NET Core):**

```csharp
// ✅ Depend on interfaces
public class UserService : IUserService
{
    private readonly IUserRepository _repository;
    private readonly ICacheService _cache;
    private readonly ILogger<UserService> _logger;

    // Constructor injection - .NET Core DI container handles this
    public UserService(
        IUserRepository repository,
        ICacheService cache,
        ILogger<UserService> logger)
    {
        _repository = repository;
        _cache = cache;
        _logger = logger;
    }

    public async Task<User?> GetUserAsync(int id)
    {
        // Implementation...
    }
}

// Register in Program.cs
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddSingleton<ICacheService, RedisCacheService>();

// Easy to test with Moq
[Fact]
public async Task GetUserAsync_ReturnsUser_WhenExists()
{
    // Arrange
    var mockRepo = new Mock<IUserRepository>();
    var mockCache = new Mock<ICacheService>();
    var mockLogger = new Mock<ILogger<UserService>>();
    
    mockRepo.Setup(r => r.FindByIdAsync(1))
        .ReturnsAsync(new User { Id = 1, Name = "John" });
    
    var service = new UserService(mockRepo.Object, mockCache.Object, mockLogger.Object);
    
    // Act & Assert
    var user = await service.GetUserAsync(1);
    Assert.NotNull(user);
    Assert.Equal("John", user.Name);
}
```

**Example (React with Context/Hooks):**

```typescript
// Define dependencies as interfaces
interface UserService {
  getUser(id: number): Promise<User>;
}

interface CacheService {
  get<T>(key: string): T | null;
  set<T>(key: string, value: T): void;
}

// Provide dependencies via Context
const ServiceContext = createContext<
  {
    userService: UserService;
    cacheService: CacheService;
  } | null
>(null);

// Custom hook for dependency injection
function useServices() {
  const context = useContext(ServiceContext);
  if (!context) throw new Error("Services not provided");
  return context;
}

// Component depends on abstractions, not implementations
function UserProfile({ userId }: UserProfileProps) {
  const { userService, cacheService } = useServices();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    userService.getUser(userId).then(setUser);
  }, [userId, userService]);

  return <div>{user?.name}</div>;
}

// Easy to test with mock services
test("UserProfile renders user name", async () => {
  const mockUserService: UserService = {
    getUser: jest.fn().mockResolvedValue({ id: 1, name: "John" }),
  };

  render(
    <ServiceContext.Provider
      value={{ userService: mockUserService, cacheService: mockCache }}
    >
      <UserProfile userId={1} />
    </ServiceContext.Provider>,
  );

  await waitFor(() => {
    expect(screen.getByText("John")).toBeInTheDocument();
  });
});
```

### 4. **Fail Fast, Fail Loud**

Detect errors early and report them clearly.

**Principles:**

- Validate input at boundaries
- Don't hide errors
- Use proper error types
- Log with context
- Fail early in startup for configuration errors

**Example (.NET Core):**

```csharp
// ✅ Fail fast at startup
public class PaymentService
{
    private readonly string _apiKey;
    private readonly int _timeoutSeconds;

    public PaymentService(IConfiguration config)
    {
        // Validate required configuration on startup - fail immediately
        _apiKey = config["Payment:ApiKey"]
            ?? throw new InvalidOperationException("Payment:ApiKey is required");
        
        _timeoutSeconds = config.GetValue<int>("Payment:Timeout");
        if (_timeoutSeconds <= 0)
            throw new InvalidOperationException("Payment:Timeout must be positive");
    }
}

// ✅ Validate at API boundaries
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    [HttpPost]
    public async Task<ActionResult<UserResponse>> CreateUser(
        [FromBody] CreateUserRequest request)
    {
        // Model validation happens automatically via data annotations
        if (!ModelState.IsValid)
        {
            return BadRequest(new ErrorResponse
            {
                Message = "Validation failed",
                Errors = ModelState.ToDictionary(
                    kvp => kvp.Key,
                    kvp => kvp.Value?.Errors.Select(e => e.ErrorMessage).ToArray() ?? Array.Empty<string>()
                )
            });
        }

        // Additional business validation
        if (await _userService.EmailExistsAsync(request.Email))
        {
            return Conflict(new ErrorResponse
            {
                Message = "A user with this email already exists"
            });
        }

        // Continue with valid data...
        var user = await _userService.CreateUserAsync(request);
        return CreatedAtAction(nameof(GetUser), new { id = user.Id }, user);
    }
}

// ✅ Use Data Annotations for validation
public class CreateUserRequest
{
    [Required(ErrorMessage = "Email is required")]
    [EmailAddress(ErrorMessage = "Invalid email format")]
    public string Email { get; set; } = string.Empty;

    [Required(ErrorMessage = "Name is required")]
    [StringLength(100, MinimumLength = 2, ErrorMessage = "Name must be 2-100 characters")]
    public string Name { get; set; } = string.Empty;

    [Range(18, 120, ErrorMessage = "Age must be between 18 and 120")]
    public int Age { get; set; }
}
```

**Example (React/TypeScript):**

```typescript
// ✅ Fail fast with TypeScript type checking
interface CreateUserRequest {
  email: string;
  name: string;
  age: number;
}

// ✅ Validate input at form boundaries
function CreateUserForm() {
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = (data: CreateUserRequest): boolean => {
    const newErrors: Record<string, string> = {};

    // Fail fast - validate all inputs
    if (!data.email) {
      newErrors.email = "Email is required";
    } else if (!/\S+@\S+\.\S+/.test(data.email)) {
      newErrors.email = "Invalid email format";
    }

    if (!data.name || data.name.length < 2) {
      newErrors.name = "Name must be at least 2 characters";
    }

    if (data.age < 18 || data.age > 120) {
      newErrors.age = "Age must be between 18 and 120";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);

    const userData: CreateUserRequest = {
      email: formData.get("email") as string,
      name: formData.get("name") as string,
      age: parseInt(formData.get("age") as string),
    };

    // Fail fast - stop if validation fails
    if (!validateForm(userData)) {
      return;
    }

    try {
      await userService.createUser(userData);
      toast.success("User created successfully");
    } catch (error) {
      // Fail loud - show clear error to user
      if (error instanceof ApiError) {
        toast.error(`Failed to create user: ${error.message}`);
      } else {
        toast.error("An unexpected error occurred");
      }
      console.error("User creation failed:", error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input name="email" />
      {errors.email && <span className="error">{errors.email}</span>}

      <input name="name" />
      {errors.name && <span className="error">{errors.name}</span>}

      <input name="age" type="number" />
      {errors.age && <span className="error">{errors.age}</span>}

      <button type="submit">Create User</button>
    </form>
  );
}
```

### 5. **Explicit Over Implicit**

Code should be clear and obvious, not clever or magical.

**Rules:**

- Avoid hidden dependencies
- Avoid global state
- Be explicit about errors
- Don't use reflection unless necessary
- Avoid "magic" strings/numbers

**Example (.NET Core):**

```csharp
// ❌ Implicit, uses global/static state
public static class Database
{
    public static DbContext Context { get; set; }  // Global - where did this come from?
}

public class UserService
{
    public User GetUser(int id)
    {
        // Hidden dependency on global state
        return Database.Context.Users.Find(id);
    }
}

// ✅ Explicit dependencies via constructor injection
public class UserService : IUserService
{
    private readonly AppDbContext _context;
    private readonly ILogger<UserService> _logger;

    // Clear where dependencies come from
    public UserService(AppDbContext context, ILogger<UserService> logger)
    {
        _context = context ?? throw new ArgumentNullException(nameof(context));
        _logger = logger ?? throw new ArgumentNullException(nameof(logger));
    }

    public async Task<User?> GetUserAsync(int id, CancellationToken cancellationToken)
    {
        // Explicit about async operation and cancellation
        return await _context.Users
            .AsNoTracking()
            .FirstOrDefaultAsync(u => u.Id == id, cancellationToken);
    }
}

// ❌ Magic strings
var connectionString = config["ConnectionStrings:Default"];
var timeout = int.Parse(config["Timeout"]);

// ✅ Explicit configuration classes
public class DatabaseOptions
{
    public const string SectionName = "Database";
    
    public string ConnectionString { get; set; } = string.Empty;
    public int TimeoutSeconds { get; set; } = 30;
    public bool EnableRetryOnFailure { get; set; } = true;
}

// In Program.cs
builder.Services.Configure<DatabaseOptions>(
    builder.Configuration.GetSection(DatabaseOptions.SectionName));

// Usage
public class UserRepository
{
    private readonly DatabaseOptions _options;
    
    public UserRepository(IOptions<DatabaseOptions> options)
    {
        _options = options.Value;
    }
}
```

**Example (React/TypeScript):**

```typescript
// ❌ Implicit global state
let currentUser: User | null = null; // Where did this come from?

function UserProfile() {
  return <div>{currentUser?.name}</div>; // Hidden dependency
}

// ✅ Explicit props and state management
interface UserProfileProps {
  user: User; // Explicit dependency
  onUpdate: (user: User) => void; // Explicit callback
}

function UserProfile({ user, onUpdate }: UserProfileProps) {
  return (
    <div>
      <h1>{user.name}</h1>
      <button onClick={() => onUpdate({ ...user, name: "New Name" })}>
        Update
      </button>
    </div>
  );
}

// ❌ Magic strings and implicit behavior
function fetchUser() {
  return fetch("/api/users/1").then((r) => r.json());
}

// ✅ Explicit configuration and error handling
const API_BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:5000";
const API_TIMEOUT_MS = 30000;

interface FetchOptions {
  timeoutMs?: number;
  retries?: number;
}

async function fetchUser(
  id: number,
  options: FetchOptions = {},
): Promise<User> {
  const { timeoutMs = API_TIMEOUT_MS, retries = 3 } = options;

  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetch(`${API_BASE_URL}/api/users/${id}`, {
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      throw new ApiError(response.status, "Failed to fetch user");
    }

    return await response.json();
  } catch (error) {
    clearTimeout(timeoutId);

    // Explicit error handling
    if (error instanceof DOMException && error.name === "AbortError") {
      throw new TimeoutError(`Request timed out after ${timeoutMs}ms`);
    }
    throw error;
  }
}
```

### 6. **Design for Testability**

If code is hard to test, it's probably poorly designed.

**Guidelines:**

- Small, focused functions (< 50 lines)
- Inject dependencies
- Avoid global state
- Pure functions when possible
- Test business logic independently of infrastructure

**Example (.NET Core with xUnit and Moq):**

```csharp
// ✅ Testable design - all dependencies are interfaces
public class OrderService : IOrderService
{
    private readonly IPaymentProcessor _paymentProcessor;
    private readonly IInventoryService _inventory;
    private readonly INotificationService _notifier;
    private readonly ILogger<OrderService> _logger;

    public OrderService(
        IPaymentProcessor paymentProcessor,
        IInventoryService inventory,
        INotificationService notifier,
        ILogger<OrderService> logger)
    {
        _paymentProcessor = paymentProcessor;
        _inventory = inventory;
        _notifier = notifier;
        _logger = logger;
    }

    public async Task<Result<Order>> PlaceOrderAsync(
        PlaceOrderRequest request,
        CancellationToken cancellationToken)
    {
        // Pure business logic - easy to test
        var validation = ValidateOrder(request);
        if (!validation.IsSuccess)
            return Result<Order>.Failure(validation.Error);

        var order = new Order
        {
            UserId = request.UserId,
            Items = request.Items,
            TotalAmount = request.Items.Sum(i => i.Price * i.Quantity)
        };

        // All dependencies are mockable
        var reservationResult = await _inventory.ReserveItemsAsync(
            order.Items, cancellationToken);
        
        if (!reservationResult.IsSuccess)
            return Result<Order>.Failure("Insufficient inventory");

        try
        {
            var paymentResult = await _paymentProcessor.ChargeAsync(
                order.TotalAmount, request.PaymentMethod, cancellationToken);

            if (!paymentResult.IsSuccess)
            {
                await _inventory.ReleaseItemsAsync(order.Items, cancellationToken);
                return Result<Order>.Failure("Payment failed");
            }

            order.Status = OrderStatus.Confirmed;
            
            await _notifier.SendOrderConfirmationAsync(
                order.UserId, order.Id, cancellationToken);

            return Result<Order>.Success(order);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to place order for user {UserId}", order.UserId);
            await _inventory.ReleaseItemsAsync(order.Items, cancellationToken);
            throw;
        }
    }

    // Pure function - easy to test without dependencies
    private static Result ValidateOrder(PlaceOrderRequest request)
    {
        if (request.Items.Count == 0)
            return Result.Failure("Order must contain at least one item");
        
        if (request.Items.Any(i => i.Quantity <= 0))
            return Result.Failure("Item quantities must be positive");

        return Result.Success();
    }
}

// Unit test - all dependencies are mocked
public class OrderServiceTests
{
    private readonly Mock<IPaymentProcessor> _mockPayment;
    private readonly Mock<IInventoryService> _mockInventory;
    private readonly Mock<INotificationService> _mockNotifier;
    private readonly Mock<ILogger<OrderService>> _mockLogger;
    private readonly OrderService _service;

    public OrderServiceTests()
    {
        _mockPayment = new Mock<IPaymentProcessor>();
        _mockInventory = new Mock<IInventoryService>();
        _mockNotifier = new Mock<INotificationService>();
        _mockLogger = new Mock<ILogger<OrderService>>();
        
        _service = new OrderService(
            _mockPayment.Object,
            _mockInventory.Object,
            _mockNotifier.Object,
            _mockLogger.Object
        );
    }

    [Fact]
    public async Task PlaceOrderAsync_ReturnsSuccess_WhenOrderIsValid()
    {
        // Arrange
        var request = new PlaceOrderRequest
        {
            UserId = 1,
            Items = new List<OrderItem> { new() { ProductId = 1, Quantity = 2, Price = 10 } },
            PaymentMethod = "CreditCard"
        };

        _mockInventory
            .Setup(x => x.ReserveItemsAsync(It.IsAny<List<OrderItem>>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Success());

        _mockPayment
            .Setup(x => x.ChargeAsync(20m, "CreditCard", It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Success());

        // Act
        var result = await _service.PlaceOrderAsync(request, CancellationToken.None);

        // Assert
        Assert.True(result.IsSuccess);
        Assert.Equal(OrderStatus.Confirmed, result.Value.Status);
        _mockNotifier.Verify(
            x => x.SendOrderConfirmationAsync(1, It.IsAny<int>(), It.IsAny<CancellationToken>()),
            Times.Once
        );
    }

    [Fact]
    public async Task PlaceOrderAsync_ReleasesInventory_WhenPaymentFails()
    {
        // Arrange
        var request = new PlaceOrderRequest
        {
            UserId = 1,
            Items = new List<OrderItem> { new() { ProductId = 1, Quantity = 2, Price = 10 } }
        };

        _mockInventory
            .Setup(x => x.ReserveItemsAsync(It.IsAny<List<OrderItem>>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Success());

        _mockPayment
            .Setup(x => x.ChargeAsync(It.IsAny<decimal>(), It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Failure("Insufficient funds"));

        // Act
        var result = await _service.PlaceOrderAsync(request, CancellationToken.None);

        // Assert
        Assert.False(result.IsSuccess);
        _mockInventory.Verify(
            x => x.ReleaseItemsAsync(It.IsAny<List<OrderItem>>(), It.IsAny<CancellationToken>()),
            Times.Once,
            "Inventory should be released when payment fails"
        );
    }
}
```

**Example (React with Jest and React Testing Library):**

```typescript
// ✅ Testable component - dependencies injected via props/context
interface CheckoutFormProps {
  cartItems: CartItem[];
  onSubmit: (order: Order) => Promise<void>;
  paymentService: PaymentService;
}

function CheckoutForm(
  { cartItems, onSubmit, paymentService }: CheckoutFormProps,
) {
  const [processing, setProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Pure calculation - easy to test separately
  const calculateTotal = (items: CartItem[]): number => {
    return items.reduce((sum, item) => sum + item.price * item.quantity, 0);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setProcessing(true);
    setError(null);

    try {
      const total = calculateTotal(cartItems);
      await paymentService.processPayment(total);

      const order: Order = {
        items: cartItems,
        total,
        status: "confirmed",
      };

      await onSubmit(order);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Payment failed");
    } finally {
      setProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>Total: ${calculateTotal(cartItems).toFixed(2)}</div>
      {error && <div className="error">{error}</div>}
      <button type="submit" disabled={processing || cartItems.length === 0}>
        {processing ? "Processing..." : "Place Order"}
      </button>
    </form>
  );
}

// Unit tests
describe("CheckoutForm", () => {
  const mockPaymentService: PaymentService = {
    processPayment: jest.fn().mockResolvedValue(undefined),
  };

  const mockCartItems: CartItem[] = [
    { id: 1, name: "Product 1", price: 10, quantity: 2 },
    { id: 2, name: "Product 2", price: 5, quantity: 1 },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("calculates total correctly", () => {
    const onSubmit = jest.fn();

    render(
      <CheckoutForm
        cartItems={mockCartItems}
        onSubmit={onSubmit}
        paymentService={mockPaymentService}
      />,
    );

    expect(screen.getByText("Total: $25.00")).toBeInTheDocument();
  });

  it("processes payment and submits order on form submit", async () => {
    const onSubmit = jest.fn().mockResolvedValue(undefined);

    render(
      <CheckoutForm
        cartItems={mockCartItems}
        onSubmit={onSubmit}
        paymentService={mockPaymentService}
      />,
    );

    const submitButton = screen.getByRole("button", { name: /place order/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockPaymentService.processPayment).toHaveBeenCalledWith(25);
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          items: mockCartItems,
          total: 25,
          status: "confirmed",
        }),
      );
    });
  });

  it("displays error when payment fails", async () => {
    const failingPaymentService: PaymentService = {
      processPayment: jest.fn().mockRejectedValue(new Error("Card declined")),
    };

    render(
      <CheckoutForm
        cartItems={mockCartItems}
        onSubmit={jest.fn()}
        paymentService={failingPaymentService}
      />,
    );

    const submitButton = screen.getByRole("button", { name: /place order/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText("Card declined")).toBeInTheDocument();
    });
  });

  it("disables button when cart is empty", () => {
    render(
      <CheckoutForm
        cartItems={[]}
        onSubmit={jest.fn()}
        paymentService={mockPaymentService}
      />,
    );

    const submitButton = screen.getByRole("button", { name: /place order/i });
    expect(submitButton).toBeDisabled();
  });
});
```

### 7. **Configuration Over Code**

Behavior should be configurable without code changes.

**What to configure:**

- URLs and endpoints
- Timeouts and retries
- Feature flags
- Resource limits
- Log levels

**Example (.NET Core):**

```csharp
// ✅ Strongly-typed configuration with IOptions pattern
public class ApiOptions
{
    public const string SectionName = "Api";

    public string BaseUrl { get; set; } = "https://api.example.com";
    public int TimeoutSeconds { get; set; } = 30;
    public int RetryAttempts { get; set; } = 3;
    public int RateLimitRequestsPerSecond { get; set; } = 100;
}

public class DatabaseOptions
{
    public const string SectionName = "Database";

    public string ConnectionString { get; set; } = string.Empty;
    public int MaxPoolSize { get; set; } = 100;
    public int CommandTimeoutSeconds { get; set; } = 30;
    public bool EnableSensitiveDataLogging { get; set; } = false;
}

public class FeatureFlags
{
    public const string SectionName = "Features";

    public bool EnableNewCheckoutFlow { get; set; } = false;
    public bool EnableCaching { get; set; } = true;
    public bool EnableDetailedLogging { get; set; } = false;
    public int MaxUploadSizeMB { get; set; } = 10;
}

// appsettings.json
{
  "Api": {
    "BaseUrl": "https://api.production.com",
    "TimeoutSeconds": 30,
    "RetryAttempts": 3,
    "RateLimitRequestsPerSecond": 100
  },
  "Database": {
    "ConnectionString": "Server=localhost;Database=MyApp;",
    "MaxPoolSize": 100,
    "CommandTimeoutSeconds": 30,
    "EnableSensitiveDataLogging": false
  },
  "Features": {
    "EnableNewCheckoutFlow": false,
    "EnableCaching": true,
    "EnableDetailedLogging": false,
    "MaxUploadSizeMB": 10
  }
}

// Program.cs - register configurations
builder.Services.Configure<ApiOptions>(
    builder.Configuration.GetSection(ApiOptions.SectionName));
builder.Services.Configure<DatabaseOptions>(
    builder.Configuration.GetSection(DatabaseOptions.SectionName));
builder.Services.Configure<FeatureFlags>(
    builder.Configuration.GetSection(FeatureFlags.SectionName));

// Validate configuration on startup
builder.Services.AddOptions<ApiOptions>()
    .Bind(builder.Configuration.GetSection(ApiOptions.SectionName))
    .ValidateDataAnnotations()
    .ValidateOnStart();

// Usage in services
public class ExternalApiService
{
    private readonly ApiOptions _options;
    private readonly ILogger<ExternalApiService> _logger;

    public ExternalApiService(
        IOptions<ApiOptions> options,
        ILogger<ExternalApiService> logger)
    {
        _options = options.Value;
        _logger = logger;
    }

    public async Task<T> GetAsync<T>(string endpoint)
    {
        using var client = new HttpClient
        {
            BaseAddress = new Uri(_options.BaseUrl),
            Timeout = TimeSpan.FromSeconds(_options.TimeoutSeconds)
        };

        // Use configured retry attempts
        for (int attempt = 1; attempt <= _options.RetryAttempts; attempt++)
        {
            try
            {
                var response = await client.GetAsync(endpoint);
                response.EnsureSuccessStatusCode();
                return await response.Content.ReadFromJsonAsync<T>();
            }
            catch (Exception ex) when (attempt < _options.RetryAttempts)
            {
                _logger.LogWarning(ex, 
                    "API request failed, attempt {Attempt}/{MaxAttempts}", 
                    attempt, _options.RetryAttempts);
                await Task.Delay(TimeSpan.FromSeconds(Math.Pow(2, attempt))); // Exponential backoff
            }
        }

        throw new InvalidOperationException("API request failed after all retry attempts");
    }
}

// Feature flag usage
public class CheckoutController : ControllerBase
{
    private readonly FeatureFlags _features;

    public CheckoutController(IOptions<FeatureFlags> features)
    {
        _features = features.Value;
    }

    [HttpPost]
    public IActionResult ProcessCheckout(CheckoutRequest request)
    {
        if (_features.EnableNewCheckoutFlow)
        {
            return ProcessNewCheckout(request);
        }
        return ProcessLegacyCheckout(request);
    }
}
```

**Example (React with environment variables):**

```typescript
// ✅ Configuration via environment variables
// .env.production
REACT_APP_API_URL=https://api.production.com
REACT_APP_API_TIMEOUT_MS=30000
REACT_APP_ENABLE_ANALYTICS=true
REACT_APP_ENABLE_DEBUG_LOGS=false
REACT_APP_MAX_UPLOAD_SIZE_MB=10

// config.ts - centralized configuration
export interface AppConfig {
  api: {
    baseUrl: string;
    timeoutMs: number;
  };
  features: {
    enableAnalytics: boolean;
    enableDebugLogs: boolean;
    maxUploadSizeMB: number;
  };
}

const config: AppConfig = {
  api: {
    baseUrl: process.env.REACT_APP_API_URL || "http://localhost:5000",
    timeoutMs: parseInt(process.env.REACT_APP_API_TIMEOUT_MS || "30000"),
  },
  features: {
    enableAnalytics: process.env.REACT_APP_ENABLE_ANALYTICS === "true",
    enableDebugLogs: process.env.REACT_APP_ENABLE_DEBUG_LOGS === "true",
    maxUploadSizeMB: parseInt(process.env.REACT_APP_MAX_UPLOAD_SIZE_MB || "10"),
  },
};

// Validate configuration at startup
function validateConfig(config: AppConfig): void {
  if (!config.api.baseUrl) {
    throw new Error("API base URL is required");
  }
  if (config.api.timeoutMs <= 0) {
    throw new Error("API timeout must be positive");
  }
}

validateConfig(config);

export default config;

// Usage in services
import config from "./config";

class ApiService {
  private baseUrl = config.api.baseUrl;
  private timeout = config.api.timeoutMs;

  async get<T>(endpoint: string): Promise<T> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      return await response.json();
    } finally {
      clearTimeout(timeoutId);
    }
  }
}

// Feature flag usage in components
function UploadForm() {
  const maxSize = config.features.maxUploadSizeMB * 1024 * 1024;

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file && file.size > maxSize) {
      alert(`File size must not exceed ${config.features.maxUploadSizeMB}MB`);
      e.target.value = "";
    }
  };

  return <input type="file" onChange={handleFileChange} />;
}

// Analytics wrapper respects feature flag
function trackEvent(eventName: string, properties?: Record<string, any>) {
  if (config.features.enableAnalytics) {
    // Send to analytics service
    analytics.track(eventName, properties);
  }
  
  if (config.features.enableDebugLogs) {
    console.log("Event tracked:", eventName, properties);
  }
}
```

### 8. **Graceful Degradation**

System should degrade gracefully when dependencies fail.

**Strategies:**

- Circuit breakers
- Fallbacks
- Timeouts
- Retries with backoff
- Caching

**Example (.NET Core with Polly):**

```csharp
// ✅ Graceful degradation with Polly resilience patterns
public class UserProfileService
{
    private readonly IHttpClientFactory _httpClientFactory;
    private readonly IMemoryCache _cache;
    private readonly ILogger<UserProfileService> _logger;
    private readonly IAsyncPolicy<UserProfile> _resiliencePolicy;

    public UserProfileService(
        IHttpClientFactory httpClientFactory,
        IMemoryCache cache,
        ILogger<UserProfileService> logger)
    {
        _httpClientFactory = httpClientFactory;
        _cache = cache;
        _logger = logger;

        // Define resilience policy with circuit breaker, retry, and fallback
        _resiliencePolicy = Policy<UserProfile>
            .Handle<HttpRequestException>()
            .Or<TimeoutException>()
            .FallbackAsync(
                fallbackValue: GetFallbackProfile(),
                onFallbackAsync: async (outcome, context) =>
                {
                    _logger.LogWarning("Using fallback profile due to: {Exception}", 
                        outcome.Exception?.Message);
                    await Task.CompletedTask;
                })
            .WrapAsync(Policy<UserProfile>
                .Handle<HttpRequestException>()
                .CircuitBreakerAsync(
                    handledEventsAllowedBeforeBreaking: 3,
                    durationOfBreak: TimeSpan.FromMinutes(1),
                    onBreak: (outcome, duration) =>
                    {
                        _logger.LogError("Circuit breaker opened for {Duration}", duration);
                    },
                    onReset: () =>
                    {
                        _logger.LogInformation("Circuit breaker reset");
                    }))
            .WrapAsync(Policy<UserProfile>
                .Handle<HttpRequestException>()
                .WaitAndRetryAsync(
                    retryCount: 3,
                    sleepDurationProvider: attempt => TimeSpan.FromSeconds(Math.Pow(2, attempt)),
                    onRetry: (outcome, timespan, retryAttempt, context) =>
                    {
                        _logger.LogWarning(
                            "Retry {RetryAttempt} after {Delay}ms due to: {Exception}",
                            retryAttempt, timespan.TotalMilliseconds, outcome.Exception?.Message);
                    }));
    }

    public async Task<UserProfile> GetUserProfileAsync(int userId)
    {
        var cacheKey = $"user_profile_{userId}";

        // 1. Try cache first (fastest)
        if (_cache.TryGetValue<UserProfile>(cacheKey, out var cachedProfile))
        {
            _logger.LogDebug("Retrieved profile from cache for user {UserId}", userId);
            return cachedProfile;
        }

        // 2. Try primary service with resilience policies
        try
        {
            var profile = await _resiliencePolicy.ExecuteAsync(async () =>
            {
                var client = _httpClientFactory.CreateClient("ProfileService");
                var response = await client.GetAsync($"/api/users/{userId}/profile");
                response.EnsureSuccessStatusCode();
                return await response.Content.ReadFromJsonAsync<UserProfile>()
                    ?? throw new InvalidOperationException("Profile is null");
            });

            // Update cache on success
            _cache.Set(cacheKey, profile, TimeSpan.FromMinutes(15));
            return profile;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get profile for user {UserId}", userId);

            // 3. Try stale cache as last resort
            if (_cache.TryGetValue<UserProfile>($"{cacheKey}_stale", out var staleProfile))
            {
                _logger.LogWarning("Using stale cache for user {UserId}", userId);
                return staleProfile;
            }

            // 4. Return minimal fallback profile
            _logger.LogWarning("Returning fallback profile for user {UserId}", userId);
            return GetFallbackProfile();
        }
    }

    private static UserProfile GetFallbackProfile()
    {
        return new UserProfile
        {
            DisplayName = "User",
            Email = "N/A",
            IsLimitedProfile = true
        };
    }
}

// Program.cs - configure HTTP client with timeout and retry
builder.Services.AddHttpClient("ProfileService", client =>
{
    client.BaseAddress = new Uri("https://profile-service.example.com");
    client.Timeout = TimeSpan.FromSeconds(10);
})
.AddTransientHttpErrorPolicy(policyBuilder =>
    policyBuilder.WaitAndRetryAsync(
        retryCount: 2,
        sleepDurationProvider: attempt => TimeSpan.FromSeconds(Math.Pow(2, attempt))));

// Add memory cache for resilience
builder.Services.AddMemoryCache();
```

**Example (React with fallback UI and error boundaries):**

```typescript
// ✅ Graceful degradation in React
interface UserProfileProps {
  userId: number;
}

function UserProfile({ userId }: UserProfileProps) {
  const {
    data: profile,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["userProfile", userId],
    queryFn: () => fetchUserProfile(userId),
    // Retry failed requests
    retry: 3,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
    // Use stale data while revalidating
    staleTime: 5 * 60 * 1000, // 5 minutes
    // Keep cache for fallback
    cacheTime: 60 * 60 * 1000, // 1 hour
    // Enable background refetch
    refetchOnWindowFocus: true,
  });

  // Loading state with skeleton UI
  if (isLoading) {
    return <ProfileSkeleton />;
  }

  // Error state with retry option
  if (error) {
    return (
      <ErrorFallback
        message="Failed to load user profile"
        onRetry={refetch}
        fallbackContent={<MinimalProfile userId={userId} />}
      />
    );
  }

  // Success state
  return profile
    ? <ProfileView profile={profile} />
    : <MinimalProfile userId={userId} />;
}

// Error Boundary for component-level failures
class ProfileErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error("Profile component error:", error, errorInfo);
    // Log to error tracking service
    logErrorToService(error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      // Graceful fallback UI
      return (
        <div className="error-container">
          <p>Something went wrong loading the profile.</p>
          <button
            onClick={() => this.setState({ hasError: false, error: null })}
          >
            Try Again
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

// Fetch with timeout and fallback
async function fetchUserProfile(userId: number): Promise<UserProfile> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 10000); // 10s timeout

  try {
    // Try primary API
    const response = await fetch(`/api/users/${userId}/profile`, {
      signal: controller.signal,
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    // Try fallback API endpoint
    try {
      const fallbackResponse = await fetch(`/api/users/${userId}/basic`, {
        signal: controller.signal,
      });

      if (fallbackResponse.ok) {
        console.warn("Using fallback API for user profile");
        return await fallbackResponse.json();
      }
    } catch (fallbackError) {
      console.error("Fallback API also failed:", fallbackError);
    }

    // Re-throw original error if fallback also fails
    throw error;
  } finally {
    clearTimeout(timeoutId);
  }
}

// Minimal profile component as last resort
function MinimalProfile({ userId }: { userId: number }) {
  return (
    <div className="minimal-profile">
      <div className="avatar-placeholder" />
      <p>User #{userId}</p>
      <p className="warning">Limited profile information available</p>
    </div>
  );
}

// Reusable error fallback component
interface ErrorFallbackProps {
  message: string;
  onRetry?: () => void;
  fallbackContent?: React.ReactNode;
}

function ErrorFallback(
  { message, onRetry, fallbackContent }: ErrorFallbackProps,
) {
  return (
    <div className="error-fallback">
      <p className="error-message">{message}</p>
      {onRetry && (
        <button onClick={onRetry} className="retry-button">
          Retry
        </button>
      )}
      {fallbackContent && (
        <div className="fallback-content">
          <p>Showing limited information:</p>
          {fallbackContent}
        </div>
      )}
    </div>
  );
}

// Usage with error boundary
function App() {
  return (
    <ProfileErrorBoundary>
      <UserProfile userId={123} />
    </ProfileErrorBoundary>
  );
}
```

### 9. **Security by Design**

Security should be built-in, not bolted on.

**Principles:**

- Defense in depth
- Least privilege
- Fail securely
- No security through obscurity
- Validate all input

**Example (.NET Core):**

```csharp
// ✅ Security built into the design
[ApiController]
[Route("api/[controller]")]
[Authorize] // Require authentication for entire controller
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly IAuthorizationService _authorizationService;
    private readonly ILogger<UsersController> _logger;

    public UsersController(
        IUserService userService,
        IAuthorizationService authorizationService,
        ILogger<UsersController> logger)
    {
        _userService = userService;
        _authorizationService = authorizationService;
        _logger = logger;
    }

    [HttpPut("{id}")]
    [ValidateAntiForgeryToken] // CSRF protection
    public async Task<IActionResult> UpdateUser(
        int id,
        [FromBody] UpdateUserRequest request)
    {
        // 1. Authenticate - handled by [Authorize] attribute
        var currentUserId = int.Parse(User.FindFirst(ClaimTypes.NameIdentifier)?.Value 
            ?? throw new UnauthorizedAccessException());

        // 2. Authorize - check if user can update this resource
        var user = await _userService.GetUserAsync(id);
        if (user == null)
            return NotFound();

        var authResult = await _authorizationService.AuthorizeAsync(
            User, user, "CanUpdateUser");

        if (!authResult.Succeeded)
        {
            _logger.LogWarning(
                "User {UserId} attempted unauthorized update of user {TargetId}",
                currentUserId, id);
            return Forbid();
        }

        // 3. Validate input - automatic via data annotations
        if (!ModelState.IsValid)
        {
            return BadRequest(new ErrorResponse
            {
                Message = "Invalid input",
                Errors = ModelState.ToDictionary(
                    kvp => kvp.Key,
                    kvp => kvp.Value?.Errors.Select(e => e.ErrorMessage).ToArray() 
                        ?? Array.Empty<string>())
            });
        }

        // 4. Sanitize input
        request.Email = request.Email.Trim().ToLowerInvariant();
        request.Name = SecurityHelper.SanitizeInput(request.Name);

        // 5. Apply business rules with security in mind
        if (request.Role == UserRole.Admin && !User.IsInRole("SuperAdmin"))
        {
            _logger.LogWarning(
                "User {UserId} attempted to assign Admin role without permission",
                currentUserId);
            return Forbid();
        }

        // 6. Process securely
        var updatedUser = await _userService.UpdateUserAsync(id, request);

        // 7. Log security-relevant actions
        _logger.LogInformation(
            "User {ActorId} updated user {TargetId}",
            currentUserId, id);

        return Ok(updatedUser);
    }
}

// Input validation with data annotations
public class UpdateUserRequest
{
    [Required]
    [EmailAddress]
    [MaxLength(255)]
    public string Email { get; set; } = string.Empty;

    [Required]
    [StringLength(100, MinimumLength = 2)]
    [RegularExpression(@"^[a-zA-Z\s'-]+$", ErrorMessage = "Name contains invalid characters")]
    public string Name { get; set; } = string.Empty;

    [EnumDataType(typeof(UserRole))]
    public UserRole Role { get; set; }
}

// Authorization policy-based access control
public class CanUpdateUserRequirement : IAuthorizationRequirement { }

public class CanUpdateUserHandler : AuthorizationHandler<CanUpdateUserRequirement, User>
{
    protected override Task HandleRequirementAsync(
        AuthorizationHandlerContext context,
        CanUpdateUserRequirement requirement,
        User targetUser)
    {
        var currentUserId = context.User.FindFirst(ClaimTypes.NameIdentifier)?.Value;

        // Users can update their own profile
        if (targetUser.Id.ToString() == currentUserId)
        {
            context.Succeed(requirement);
            return Task.CompletedTask;
        }

        // Admins can update any user
        if (context.User.IsInRole("Admin"))
        {
            context.Succeed(requirement);
            return Task.CompletedTask;
        }

        // Otherwise, deny
        context.Fail();
        return Task.CompletedTask;
    }
}

// Program.cs - configure security
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ClockSkew = TimeSpan.Zero // No clock skew tolerance
        };
    });

builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("CanUpdateUser", policy =>
        policy.Requirements.Add(new CanUpdateUserRequirement()));
});

builder.Services.AddSingleton<IAuthorizationHandler, CanUpdateUserHandler>();

// Security headers
app.Use(async (context, next) =>
{
    context.Response.Headers.Add("X-Content-Type-Options", "nosniff");
    context.Response.Headers.Add("X-Frame-Options", "DENY");
    context.Response.Headers.Add("X-XSS-Protection", "1; mode=block");
    context.Response.Headers.Add("Referrer-Policy", "strict-origin-when-cross-origin");
    context.Response.Headers.Add(
        "Content-Security-Policy",
        "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'");
    await next();
});

// Enable HTTPS redirection
app.UseHttpsRedirection();

// Rate limiting
app.UseRateLimiter();

// Database queries - always use parameterized queries (EF Core does this automatically)
public async Task<User?> FindByEmailAsync(string email)
{
    // ✅ Safe - EF Core uses parameterized queries
    return await _context.Users
        .AsNoTracking()
        .FirstOrDefaultAsync(u => u.Email == email);
    
    // ❌ NEVER do this - SQL injection vulnerability
    // return await _context.Users
    //     .FromSqlRaw($"SELECT * FROM Users WHERE Email = '{email}'")
    //     .FirstOrDefaultAsync();
}
```

**Example (React/TypeScript):**

```typescript
// ✅ Security in React applications
import DOMPurify from "dompurify";

// 1. Input validation and sanitization
interface CommentFormProps {
  onSubmit: (comment: string) => void;
}

function CommentForm({ onSubmit }: CommentFormProps) {
  const [comment, setComment] = useState("");
  const [error, setError] = useState<string | null>(null);

  const validateComment = (text: string): boolean => {
    setError(null);

    // Length validation
    if (text.length < 3) {
      setError("Comment must be at least 3 characters");
      return false;
    }
    if (text.length > 500) {
      setError("Comment must not exceed 500 characters");
      return false;
    }

    // Content validation - no scripts or suspicious patterns
    if (/<script|javascript:|onerror=/i.test(text)) {
      setError("Invalid content detected");
      return false;
    }

    return true;
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();

    if (!validateComment(comment)) {
      return;
    }

    // Sanitize before sending
    const sanitized = DOMPurify.sanitize(comment, {
      ALLOWED_TAGS: [], // Strip all HTML
      ALLOWED_ATTR: [],
    });

    onSubmit(sanitized);
    setComment("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <textarea
        value={comment}
        onChange={(e) => setComment(e.target.value)}
        maxLength={500}
        placeholder="Enter your comment..."
      />
      {error && <span className="error">{error}</span>}
      <button type="submit">Submit</button>
    </form>
  );
}

// 2. Secure API client with CSRF protection
class SecureApiClient {
  private baseUrl: string;
  private csrfToken: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async getCsrfToken(): Promise<string> {
    if (this.csrfToken) return this.csrfToken;

    const response = await fetch(`${this.baseUrl}/api/csrf-token`, {
      credentials: "include", // Send cookies
    });

    const data = await response.json();
    this.csrfToken = data.token;
    return this.csrfToken;
  }

  async post<T>(endpoint: string, body: any): Promise<T> {
    const csrfToken = await this.getCsrfToken();

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": csrfToken, // CSRF protection
      },
      credentials: "include", // Include authentication cookies
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      if (response.status === 401) {
        // Redirect to login
        window.location.href = "/login";
        throw new Error("Unauthorized");
      }
      throw new Error(`HTTP ${response.status}`);
    }

    return await response.json();
  }
}

// 3. Render user content safely
interface UserCommentProps {
  comment: string;
  allowHtml?: boolean;
}

function UserComment({ comment, allowHtml = false }: UserCommentProps) {
  if (allowHtml) {
    // Sanitize HTML before rendering
    const sanitized = DOMPurify.sanitize(comment, {
      ALLOWED_TAGS: ["b", "i", "em", "strong", "a"],
      ALLOWED_ATTR: ["href"],
    });

    return <div dangerouslySetInnerHTML={{ __html: sanitized }} />;
  }

  // Default: render as plain text (safe)
  return <div>{comment}</div>;
}

// 4. Protected routes with authentication check
interface ProtectedRouteProps {
  children: React.ReactNode;
}

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (!user) {
    // Redirect to login
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

// 5. Role-based access control
interface AdminRouteProps {
  children: React.ReactNode;
}

function AdminRoute({ children }: AdminRouteProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (!user || !user.roles.includes("Admin")) {
    return <Navigate to="/unauthorized" replace />;
  }

  return <>{children}</>;
}

// 6. Secure local storage wrapper
class SecureStorage {
  // Never store sensitive data in localStorage
  static set(key: string, value: any): void {
    if (this.isSensitiveKey(key)) {
      console.error("Attempted to store sensitive data in localStorage");
      return;
    }
    localStorage.setItem(key, JSON.stringify(value));
  }

  static get<T>(key: string): T | null {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : null;
  }

  private static isSensitiveKey(key: string): boolean {
    const sensitive = ["password", "token", "secret", "apiKey", "creditCard"];
    return sensitive.some((s) => key.toLowerCase().includes(s));
  }
}

// Usage in App
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />

        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin"
          element={
            <AdminRoute>
              <AdminPanel />
            </AdminRoute>
          }
        />
      </Routes>
    </Router>
  );
}
```

### 10. **Monitoring & Observability**

Systems should be designed to be observable.

**The Three Pillars:**

- **Logs:** What happened? (events, errors)
- **Metrics:** How much/how many? (counts, durations, gauges)
- **Traces:** Where did time go? (distributed tracing)

**Example (.NET Core with Application Insights and Serilog):**

```csharp
// ✅ Comprehensive observability
public class OrderService : IOrderService
{
    private readonly ILogger<OrderService> _logger;
    private readonly TelemetryClient _telemetryClient;
    private readonly IOrderRepository _repository;

    public OrderService(
        ILogger<OrderService> logger,
        TelemetryClient telemetryClient,
        IOrderRepository repository)
    {
        _logger = logger;
        _telemetryClient = telemetryClient;
        _repository = repository;
    }

    public async Task<Result<Order>> ProcessOrderAsync(
        PlaceOrderRequest request,
        CancellationToken cancellationToken)
    {
        var stopwatch = Stopwatch.StartNew();
        var orderId = Guid.NewGuid();

        // Structured logging with context
        using (_logger.BeginScope(new Dictionary<string, object>
        {
            ["OrderId"] = orderId,
            ["UserId"] = request.UserId,
            ["OrderAmount"] = request.TotalAmount
        }))
        {
            _logger.LogInformation(
                "Processing order {OrderId} for user {UserId} with amount {Amount:C}",
                orderId, request.UserId, request.TotalAmount);

            try
            {
                // Custom metric - order count
                _telemetryClient.GetMetric("Orders.Processed").TrackValue(1);
                
                // Custom metric - order value
                _telemetryClient.GetMetric("Orders.Value").TrackValue(request.TotalAmount);

                // Process order
                var order = await _repository.CreateOrderAsync(request, cancellationToken);

                stopwatch.Stop();
                var duration = stopwatch.ElapsedMilliseconds;

                // Track operation duration
                _telemetryClient.TrackMetric(
                    "Order.Processing.Duration",
                    duration,
                    new Dictionary<string, string>
                    {
                        ["OrderId"] = orderId.ToString(),
                        ["Status"] = "Success"
                    });

                // Structured success log
                _logger.LogInformation(
                    "Order {OrderId} processed successfully in {Duration}ms",
                    orderId, duration);

                // Custom event for business analytics
                _telemetryClient.TrackEvent("OrderPlaced",
                    properties: new Dictionary<string, string>
                    {
                        ["OrderId"] = orderId.ToString(),
                        ["UserId"] = request.UserId.ToString(),
                        ["ItemCount"] = request.Items.Count.ToString()
                    },
                    metrics: new Dictionary<string, double>
                    {
                        ["Amount"] = (double)request.TotalAmount,
                        ["Duration"] = duration
                    });

                return Result<Order>.Success(order);
            }
            catch (InsufficientInventoryException ex)
            {
                stopwatch.Stop();
                
                // Business exception - warning level
                _logger.LogWarning(ex,
                    "Order {OrderId} failed due to insufficient inventory. Duration: {Duration}ms",
                    orderId, stopwatch.ElapsedMilliseconds);

                _telemetryClient.GetMetric("Orders.Failed.Inventory").TrackValue(1);

                return Result<Order>.Failure("Insufficient inventory");
            }
            catch (Exception ex)
            {
                stopwatch.Stop();

                // System exception - error level
                _logger.LogError(ex,
                    "Order {OrderId} failed with exception. Duration: {Duration}ms",
                    orderId, stopwatch.ElapsedMilliseconds);

                // Track exception with context
                _telemetryClient.TrackException(ex,
                    properties: new Dictionary<string, string>
                    {
                        ["OrderId"] = orderId.ToString(),
                        ["UserId"] = request.UserId.ToString(),
                        ["Stage"] = "Processing"
                    });

                _telemetryClient.GetMetric("Orders.Failed.Error").TrackValue(1);

                throw;
            }
        }
    }
}

// Program.cs - configure logging and monitoring
builder.Services.AddLogging(logging =>
{
    logging.ClearProviders();
    
    // Serilog for structured logging
    var logger = new LoggerConfiguration()
        .MinimumLevel.Information()
        .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
        .Enrich.FromLogContext()
        .Enrich.WithProperty("Application", "MyApp")
        .Enrich.WithProperty("Environment", builder.Environment.EnvironmentName)
        .WriteTo.Console(
            outputTemplate: "[{Timestamp:HH:mm:ss} {Level:u3}] {Message:lj} {Properties:j}{NewLine}{Exception}")
        .WriteTo.File(
            path: "logs/app-.log",
            rollingInterval: RollingInterval.Day,
            outputTemplate: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff zzz} [{Level:u3}] {Message:lj} {Properties:j}{NewLine}{Exception}")
        .WriteTo.ApplicationInsights(
            builder.Configuration["ApplicationInsights:ConnectionString"],
            TelemetryConverter.Traces)
        .CreateLogger();

    logging.AddSerilog(logger);
});

// Application Insights for metrics and tracing
builder.Services.AddApplicationInsightsTelemetry(options =>
{
    options.ConnectionString = builder.Configuration["ApplicationInsights:ConnectionString"];
    options.EnableAdaptiveSampling = true;
    options.EnableQuickPulseMetricStream = true;
});

// Add health checks
builder.Services.AddHealthChecks()
    .AddDbContextCheck<AppDbContext>()
    .AddUrlGroup(new Uri("https://external-api.example.com/health"), "External API")
    .AddCheck<CustomHealthCheck>("Custom Check");

var app = builder.Build();

// Health check endpoints
app.MapHealthChecks("/health", new HealthCheckOptions
{
    ResponseWriter = async (context, report) =>
    {
        context.Response.ContentType = "application/json";
        var result = JsonSerializer.Serialize(new
        {
            status = report.Status.ToString(),
            checks = report.Entries.Select(e => new
            {
                name = e.Key,
                status = e.Value.Status.ToString(),
                description = e.Value.Description,
                duration = e.Value.Duration.TotalMilliseconds
            }),
            totalDuration = report.TotalDuration.TotalMilliseconds
        });
        await context.Response.WriteAsync(result);
    }
});

// Custom middleware for request logging
app.Use(async (context, next) =>
{
    var sw = Stopwatch.StartNew();
    var requestId = Guid.NewGuid().ToString();
    
    context.Response.Headers.Add("X-Request-Id", requestId);
    
    using (LogContext.PushProperty("RequestId", requestId))
    {
        try
        {
            await next();
        }
        finally
        {
            sw.Stop();
            
            var logger = context.RequestServices.GetRequiredService<ILogger<Program>>();
            logger.LogInformation(
                "HTTP {Method} {Path} responded {StatusCode} in {Duration}ms",
                context.Request.Method,
                context.Request.Path,
                context.Response.StatusCode,
                sw.ElapsedMilliseconds);
        }
    }
});
```

**Example (React with monitoring):**

```typescript
// ✅ Frontend observability
import * as Sentry from "@sentry/react";
import { BrowserTracing } from "@sentry/tracing";

// Initialize error tracking
Sentry.init({
  dsn: process.env.REACT_APP_SENTRY_DSN,
  integrations: [new BrowserTracing()],
  tracesSampleRate: 0.1, // 10% of transactions
  environment: process.env.NODE_ENV,
  beforeSend(event, hint) {
    // Filter sensitive data
    if (event.request?.headers) {
      delete event.request.headers["Authorization"];
    }
    return event;
  },
});

// Performance monitoring
class PerformanceMonitor {
  static measureRender(componentName: string): () => void {
    const startTime = performance.now();

    return () => {
      const duration = performance.now() - startTime;

      // Log slow renders
      if (duration > 100) {
        console.warn(
          `Slow render: ${componentName} took ${duration.toFixed(2)}ms`,
        );

        // Send to analytics
        this.trackMetric("Component.Render.Duration", duration, {
          component: componentName,
          slow: "true",
        });
      }

      // Track all render times
      this.trackMetric("Component.Render.Duration", duration, {
        component: componentName,
      });
    };
  }

  static trackMetric(
    name: string,
    value: number,
    properties?: Record<string, string>,
  ): void {
    // Send to Application Insights or similar
    if (window.appInsights) {
      window.appInsights.trackMetric({ name, average: value }, properties);
    }
  }

  static trackEvent(
    name: string,
    properties?: Record<string, string>,
  ): void {
    if (window.appInsights) {
      window.appInsights.trackEvent({ name }, properties);
    }
  }
}

// Component with performance monitoring
function ProductList({ category }: ProductListProps) {
  const endMeasure = PerformanceMonitor.measureRender("ProductList");

  const { data, isLoading, error } = useQuery({
    queryKey: ["products", category],
    queryFn: () => fetchProducts(category),
    onError: (error) => {
      // Track errors
      Sentry.captureException(error, {
        tags: { component: "ProductList", category },
      });

      PerformanceMonitor.trackEvent("Products.Load.Error", {
        category,
        error: error.message,
      });
    },
    onSuccess: (data) => {
      // Track successful loads
      PerformanceMonitor.trackEvent("Products.Load.Success", {
        category,
        count: data.length.toString(),
      });
    },
  });

  useEffect(() => {
    return endMeasure; // Measure on unmount
  }, [endMeasure]);

  if (error) {
    return <ErrorMessage error={error} />;
  }

  return (
    <div>
      {data?.map((product) => (
        <ProductCard
          key={product.id}
          product={product}
        />
      ))}
    </div>
  );
}

// API client with observability
class ObservableApiClient {
  async request<T>(
    method: string,
    endpoint: string,
    body?: any,
  ): Promise<T> {
    const startTime = performance.now();
    const requestId = crypto.randomUUID();

    try {
      const response = await fetch(endpoint, {
        method,
        headers: {
          "Content-Type": "application/json",
          "X-Request-Id": requestId,
        },
        body: body ? JSON.stringify(body) : undefined,
      });

      const duration = performance.now() - startTime;

      // Log request details
      console.log(`API ${method} ${endpoint}:`, {
        status: response.status,
        duration: `${duration.toFixed(2)}ms`,
        requestId,
      });

      // Track metrics
      PerformanceMonitor.trackMetric("API.Request.Duration", duration, {
        method,
        endpoint,
        status: response.status.toString(),
      });

      // Track slow requests
      if (duration > 1000) {
        PerformanceMonitor.trackEvent("API.Request.Slow", {
          method,
          endpoint,
          duration: duration.toString(),
        });
      }

      if (!response.ok) {
        const error = new ApiError(response.status, `HTTP ${response.status}`);

        // Track API errors
        Sentry.captureException(error, {
          tags: {
            method,
            endpoint,
            status: response.status.toString(),
          },
          contexts: {
            request: { requestId },
          },
        });

        throw error;
      }

      return await response.json();
    } catch (error) {
      const duration = performance.now() - startTime;

      // Log and track errors
      console.error(`API ${method} ${endpoint} failed:`, error);

      PerformanceMonitor.trackEvent("API.Request.Error", {
        method,
        endpoint,
        error: error instanceof Error ? error.message : "Unknown error",
        duration: duration.toString(),
      });

      throw error;
    }
  }
}

// Custom error boundary with monitoring
function MonitoredErrorBoundary({ children }: { children: React.ReactNode }) {
  return (
    <Sentry.ErrorBoundary
      fallback={({ error, resetError }) => (
        <ErrorFallback error={error} onReset={resetError} />
      )}
      onError={(error, errorInfo) => {
        // Additional logging
        console.error("Component error:", error, errorInfo);

        PerformanceMonitor.trackEvent("Component.Error", {
          error: error.message,
          componentStack: errorInfo.componentStack?.slice(0, 500),
        });
      }}
    >
      {children}
    </Sentry.ErrorBoundary>
  );
}

// Usage
function App() {
  // Track page views
  useEffect(() => {
    PerformanceMonitor.trackEvent("PageView", {
      path: window.location.pathname,
    });
  }, []);

  return (
    <MonitoredErrorBoundary>
      <Router>
        <Routes>
          {/* Routes */}
        </Routes>
      </Router>
    </MonitoredErrorBoundary>
  );
}

// Web Vitals monitoring
import { getCLS, getFCP, getFID, getLCP, getTTFB } from "web-vitals";

function reportWebVitals() {
  getCLS((metric) =>
    PerformanceMonitor.trackMetric("WebVitals.CLS", metric.value)
  );
  getFID((metric) =>
    PerformanceMonitor.trackMetric("WebVitals.FID", metric.value)
  );
  getFCP((metric) =>
    PerformanceMonitor.trackMetric("WebVitals.FCP", metric.value)
  );
  getLCP((metric) =>
    PerformanceMonitor.trackMetric("WebVitals.LCP", metric.value)
  );
  getTTFB((metric) =>
    PerformanceMonitor.trackMetric("WebVitals.TTFB", metric.value)
  );
}

reportWebVitals();
```

## 🚫 Anti-Patterns to Avoid

### 1. **God Objects**

Classes that do too much.

```csharp
// ❌ God object - knows and does everything
public class Application
{
    private readonly DbContext _db;
    private readonly IMemoryCache _cache;
    private readonly IEmailService _mailer;
    private readonly IJobScheduler _scheduler;
    private readonly IHttpClientFactory _httpFactory;
    // ... 50 more dependencies

    public async Task DoEverything()
    {
        // Handles users, orders, payments, notifications, etc.
    }
}

// ✅ Split responsibilities into focused services
public class UserService { /* user-related operations only */ }
public class OrderService { /* order-related operations only */ }
public class PaymentService { /* payment-related operations only */ }
public class NotificationService { /* notification-related operations only */ }
```

### 2. **Circular Dependencies**

Package A depends on B, B depends on A.

```csharp
// ❌ Circular dependency
// UserService depends on OrderService
public class UserService
{
    private readonly OrderService _orderService;
}

// OrderService depends on UserService
public class OrderService
{
    private readonly UserService _userService;
}

// ✅ Extract shared interfaces to break the cycle
public interface IUserProvider
{
    Task<User> GetUserAsync(int id);
}

public class UserService : IUserProvider
{
    // No dependency on OrderService
}

public class OrderService
{
    private readonly IUserProvider _userProvider; // Depends on interface, not UserService
}
```

### 3. **Premature Optimization**

Optimizing before measuring.

**Rule:** Make it work → Make it right → Make it fast (in that order)

```csharp
// ❌ Premature optimization - complex caching before profiling
public class UserService
{
    private readonly ConcurrentDictionary<int, CacheEntry<User>> _l1Cache;
    private readonly IDistributedCache _l2Cache;
    private readonly IMemoryMappedFileCache _l3Cache;
    // Overly complex before knowing if caching is even needed
}

// ✅ Start simple, optimize when you have data
public class UserService
{
    private readonly AppDbContext _context;

    // Simple, clear implementation
    public async Task<User?> GetUserAsync(int id)
    {
        return await _context.Users.FindAsync(id);
    }
    
    // Add caching AFTER profiling shows it's needed
}
```

### 4. **Not Invented Here (NIH) Syndrome**

Reinventing the wheel unnecessarily.

**When to use libraries:**

- ✅ Well-maintained, popular library exists (e.g., Newtonsoft.Json, Dapper)
- ✅ Complex domain (crypto, date/time, parsing)
- ✅ Standard implementation (OAuth, JWT, etc.)

**When to write custom:**

- ✅ Core business logic unique to your domain
- ✅ Simple utility specific to your needs
- ✅ Library has security issues or is abandoned

```csharp
// ❌ Reinventing the wheel
public class MyCustomJsonParser
{
    // 5000 lines of custom JSON parsing code
}

// ✅ Use proven libraries
using System.Text.Json;

var user = JsonSerializer.Deserialize<User>(json);
```

### 5. **Tight Coupling**

Components that can't be changed independently.

**Solution:** Depend on interfaces, not implementations.

```csharp
// ❌ Tightly coupled to concrete class
public class OrderService
{
    private readonly SqlServerOrderRepository _repository; // Concrete type

    public OrderService()
    {
        _repository = new SqlServerOrderRepository(); // Direct instantiation
    }
}

// ✅ Loosely coupled via interface
public class OrderService
{
    private readonly IOrderRepository _repository; // Interface

    public OrderService(IOrderRepository repository) // Injected
    {
        _repository = repository;
    }
}

// Can easily swap implementations
builder.Services.AddScoped<IOrderRepository, SqlServerOrderRepository>();
// or
builder.Services.AddScoped<IOrderRepository, CosmosDbOrderRepository>();
```

### 6. **Magic Strings and Numbers**

Hard-coded values scattered throughout code.

```csharp
// ❌ Magic strings and numbers everywhere
if (user.Role == "admin") // What roles exist?
{
    connection.Timeout = 30; // Why 30?
}

// ✅ Use constants or enums
public enum UserRole
{
    User,
    Admin,
    SuperAdmin
}

public class DatabaseOptions
{
    public const int DefaultTimeoutSeconds = 30;
}

if (user.Role == UserRole.Admin)
{
    connection.Timeout = DatabaseOptions.DefaultTimeoutSeconds;
}
```

## 📋 Architecture Decision Records (ADRs)

For significant architectural decisions, we create ADRs:

```markdown
# ADR-001: Use Entity Framework Core with Repository Pattern

## Status

Accepted

## Context

We need a data access strategy for our .NET Core application that:

- Supports multiple database providers (SQL Server, Azure SQL, PostgreSQL)
- Provides good performance for our read-heavy workload
- Allows easy unit testing without hitting the database
- Integrates well with our existing .NET Core stack
- Supports LINQ for type-safe queries

Options considered:

1. Entity Framework Core with repository pattern
2. Dapper (micro-ORM)
3. ADO.NET with stored procedures
4. Entity Framework Core without repository pattern (DbContext as repository)

## Decision

We will use Entity Framework Core as our ORM with the repository pattern for
data access.

**Implementation details:**

- Use code-first migrations for schema management
- Implement repository pattern for testability and abstraction
- Use `AsNoTracking()` for read-only queries
- Configure DbContext pooling for performance
- Apply Fluent API configuration via `IEntityTypeConfiguration<T>`

## Consequences

**Positive:**

- Type-safe LINQ queries reduce SQL injection risk
- Code-first migrations provide version control for schema
- Repository pattern makes business logic unit testable
- DbContext pooling improves performance
- Rich ecosystem of extensions and tooling
- Familiar to .NET developers

**Negative:**

- Slight performance overhead compared to Dapper/ADO.NET
- Learning curve for complex query optimization
- Generated SQL may not always be optimal
- Requires careful attention to N+1 query issues
- Migration strategy needed for production deployments

**Mitigations:**

- Use `AsNoTracking()` for read operations
- Profile queries with SQL Server Profiler
- Use `Include()` and projection to avoid N+1 queries
- Consider Dapper for performance-critical read-heavy queries
- Implement query result caching where appropriate

## Alternatives Considered

**Dapper:**

- Pros: Better performance, full SQL control
- Cons: No change tracking, more boilerplate, manual mapping

**ADO.NET:**

- Pros: Maximum performance, full control
- Cons: Very verbose, no LINQ, manual mapping, SQL injection risk

## Review Date

2025-07-01 (6 months)

## Related Documents

- [Database Conventions Guide](./database-conventions.md)
- [API Design Guidelines](./api-design-guide.md)
```

Store ADRs in `docs/adr/` directory with format: `ADR-NNN-short-title.md`

**Example ADR topics:**

- Choice of database (SQL Server vs PostgreSQL vs Cosmos DB)
- Authentication strategy (JWT vs cookie-based)
- API versioning approach (URL vs header)
- Logging framework (Serilog vs NLog)
- Caching strategy (Redis vs memory cache)
- Frontend state management (Redux vs Context API)

## 🎯 Architecture Review Checklist

Use this checklist when reviewing architectural decisions:

- [ ] **Simplicity:** Is this the simplest solution?
- [ ] **Separation:** Are concerns properly separated?
- [ ] **Dependencies:** Are dependencies injected and testable?
- [ ] **Error Handling:** Do errors fail fast and provide context?
- [ ] **Testability:** Can this be easily tested?
- [ ] **Configuration:** Is behavior configurable?
- [ ] **Resilience:** Does it handle failures gracefully?
- [ ] **Security:** Are security principles followed?
- [ ] **Observability:** Can we monitor and debug this?
- [ ] **Documentation:** Is the design documented (ADR)?

## 📚 Recommended Reading

- **Books:**
  - _Clean Architecture_ by Robert C. Martin
  - _Domain-Driven Design_ by Eric Evans
  - _Building Microservices_ by Sam Newman
  - _Designing Data-Intensive Applications_ by Martin Kleppmann

- **Articles:**
  - [The Twelve-Factor App](https://12factor.net/)
  - [C4 Model for Software Architecture](https://c4model.com/)
  - [Microservices Patterns](https://microservices.io/patterns/)

## 🤝 Discussion & Feedback

Architecture is never perfect and always evolving. If you:

- Have questions about these principles
- Want to propose changes
- Found a better approach
- Need clarification

Create an issue or start a discussion in our team chat.

---

**Last Updated:** 2025-01-23 **Version:** 1.0.0
