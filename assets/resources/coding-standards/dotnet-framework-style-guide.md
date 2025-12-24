# .NET Framework (4.x) Coding Standards - Legacy Maintenance

## 🎯 Overview

This guide covers **.NET Framework 4.x** best practices for maintaining legacy
applications. .NET Framework reached end-of-life for active development in 2020,
with long-term support ending in 2029. Focus on sustainable maintenance,
security patching, and eventual migration planning.

**Supported Technologies:**

- ASP.NET Web Forms
- ASP.NET MVC 5
- WCF Services
- Entity Framework 6.x
- Web API 2
- Classic ADO.NET

**Key Principles:**

- Maintain security with regular patching
- Follow Microsoft's Framework Design Guidelines
- Document complex legacy patterns
- Plan for .NET 6/7/8 migration path
- Avoid introducing new technical debt

---

## 🏗️ Project Structure

### Solution Organization

```
Solution/
├── MyApp.Web/              # ASP.NET MVC or Web Forms
├── MyApp.WebApi/           # Web API 2 project
├── MyApp.Domain/           # Domain models and business logic
├── MyApp.Data/             # Entity Framework or ADO.NET
├── MyApp.Services/         # Business services layer
├── MyApp.Common/           # Shared utilities
└── MyApp.Tests/            # Unit tests (NUnit or MSTest)
```

### Namespace Conventions

Follow reverse domain notation:

```csharp
namespace CompanyName.ProjectName.Feature
{
    // Classes here
}

// Examples
namespace Contoso.CustomerPortal.Orders
namespace Contoso.CustomerPortal.Data.Repositories
namespace Contoso.CustomerPortal.Services
```

---

## 📝 C# Coding Conventions

### Naming Conventions

```csharp
// PascalCase for types, methods, properties, constants
public class CustomerRepository { }
public void ProcessOrder() { }
public string FirstName { get; set; }
public const int MaxRetries = 3;

// camelCase for local variables, parameters, private fields
private int orderCount;
private readonly ILogger _logger;  // Prefix with underscore
public void SaveCustomer(Customer customer) { }

// Use meaningful names
// ❌ Avoid
int d; // What is d?
string s;
void Do();

// ✅ Good
int daysSinceLastOrder;
string customerEmail;
void ProcessPayment();
```

### File Organization

One type per file, file name matches type name:

```csharp
// CustomerRepository.cs
public class CustomerRepository : ICustomerRepository
{
    // Private fields
    private readonly DbContext _context;
    private readonly ILogger _logger;
    
    // Constructor
    public CustomerRepository(DbContext context, ILogger logger)
    {
        _context = context ?? throw new ArgumentNullException(nameof(context));
        _logger = logger ?? throw new ArgumentNullException(nameof(logger));
    }
    
    // Public methods
    public Customer GetById(int id)
    {
        // Implementation
    }
    
    // Private methods
    private void LogOperation(string operation)
    {
        // Implementation
    }
}
```

### Using Directives

```csharp
// System namespaces first, then third-party, then local
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web.Mvc;

using Newtonsoft.Json;
using AutoMapper;

using CompanyName.ProjectName.Domain;
using CompanyName.ProjectName.Services;
```

---

## 🎨 ASP.NET MVC 5

### Controller Conventions

```csharp
public class CustomersController : Controller
{
    private readonly ICustomerService _customerService;
    private readonly ILogger _logger;
    
    public CustomersController(ICustomerService customerService, ILogger logger)
    {
        _customerService = customerService;
        _logger = logger;
    }
    
    // GET: Customers
    public ActionResult Index()
    {
        var customers = _customerService.GetAll();
        return View(customers);
    }
    
    // GET: Customers/Details/5
    public ActionResult Details(int? id)
    {
        if (id == null)
        {
            return new HttpStatusCodeResult(HttpStatusCode.BadRequest);
        }
        
        var customer = _customerService.GetById(id.Value);
        if (customer == null)
        {
            return HttpNotFound();
        }
        
        return View(customer);
    }
    
    // GET: Customers/Create
    public ActionResult Create()
    {
        return View();
    }
    
    // POST: Customers/Create
    [HttpPost]
    [ValidateAntiForgeryToken]
    public ActionResult Create(CustomerViewModel model)
    {
        if (!ModelState.IsValid)
        {
            return View(model);
        }
        
        try
        {
            _customerService.Create(model);
            TempData["Success"] = "Customer created successfully.";
            return RedirectToAction("Index");
        }
        catch (Exception ex)
        {
            _logger.Error("Failed to create customer", ex);
            ModelState.AddModelError("", "Failed to create customer. Please try again.");
            return View(model);
        }
    }
    
    // POST: Customers/Delete/5
    [HttpPost, ActionName("Delete")]
    [ValidateAntiForgeryToken]
    public ActionResult DeleteConfirmed(int id)
    {
        try
        {
            _customerService.Delete(id);
            TempData["Success"] = "Customer deleted successfully.";
            return RedirectToAction("Index");
        }
        catch (Exception ex)
        {
            _logger.Error($"Failed to delete customer {id}", ex);
            TempData["Error"] = "Failed to delete customer.";
            return RedirectToAction("Index");
        }
    }
    
    protected override void Dispose(bool disposing)
    {
        if (disposing)
        {
            _customerService?.Dispose();
        }
        base.Dispose(disposing);
    }
}
```

### View Models

Separate view models from domain models:

```csharp
// ViewModel - for UI binding
public class CustomerViewModel
{
    [Required(ErrorMessage = "First name is required")]
    [StringLength(50)]
    [Display(Name = "First Name")]
    public string FirstName { get; set; }
    
    [Required(ErrorMessage = "Last name is required")]
    [StringLength(50)]
    [Display(Name = "Last Name")]
    public string LastName { get; set; }
    
    [Required]
    [EmailAddress]
    public string Email { get; set; }
    
    [Phone]
    [Display(Name = "Phone Number")]
    public string PhoneNumber { get; set; }
    
    [Display(Name = "Is Active")]
    public bool IsActive { get; set; }
}

// Domain Model - for business logic
public class Customer
{
    public int Id { get; set; }
    public string FirstName { get; set; }
    public string LastName { get; set; }
    public string Email { get; set; }
    public string PhoneNumber { get; set; }
    public bool IsActive { get; set; }
    public DateTime CreatedDate { get; set; }
    public DateTime? ModifiedDate { get; set; }
    
    public string FullName => $"{FirstName} {LastName}";
}
```

### Razor Views

```html
@model IEnumerable<CustomerViewModel>

@{
    ViewBag.Title = "Customers";
}

<h2>@ViewBag.Title</h2>

@if (TempData["Success"] != null)
{
    <div class="alert alert-success">
        @TempData["Success"]
    </div>
}

<p>
    @Html.ActionLink("Create New Customer", "Create", null, new { @class = "btn btn-primary" })
</p>

<table class="table table-striped">
    <thead>
        <tr>
            <th>@Html.DisplayNameFor(m => m.FirstName)</th>
            <th>@Html.DisplayNameFor(m => m.LastName)</th>
            <th>@Html.DisplayNameFor(m => m.Email)</th>
            <th>@Html.DisplayNameFor(m => m.IsActive)</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        @foreach (var customer in Model)
        {
            <tr>
                <td>@Html.DisplayFor(m => customer.FirstName)</td>
                <td>@Html.DisplayFor(m => customer.LastName)</td>
                <td>@Html.DisplayFor(m => customer.Email)</td>
                <td>@Html.DisplayFor(m => customer.IsActive)</td>
                <td>
                    @Html.ActionLink("Edit", "Edit", new { id = customer.Id }) |
                    @Html.ActionLink("Details", "Details", new { id = customer.Id }) |
                    @Html.ActionLink("Delete", "Delete", new { id = customer.Id })
                </td>
            </tr>
        }
    </tbody>
</table>
```

### Custom HTML Helpers

```csharp
public static class HtmlHelperExtensions
{
    public static MvcHtmlString StatusBadge(this HtmlHelper helper, bool isActive)
    {
        var badgeClass = isActive ? "badge-success" : "badge-secondary";
        var text = isActive ? "Active" : "Inactive";
        
        var badge = new TagBuilder("span");
        badge.AddCssClass("badge");
        badge.AddCssClass(badgeClass);
        badge.SetInnerText(text);
        
        return MvcHtmlString.Create(badge.ToString());
    }
}

// Usage in view
@Html.StatusBadge(customer.IsActive)
```

---

## 🌐 ASP.NET Web Forms (Legacy)

### Code-Behind Pattern

```csharp
// Customers.aspx.cs
public partial class Customers : System.Web.UI.Page
{
    private readonly ICustomerService _customerService;
    
    // Constructor injection with Unity or other DI container
    public Customers(ICustomerService customerService)
    {
        _customerService = customerService;
    }
    
    protected void Page_Load(object sender, EventArgs e)
    {
        if (!IsPostBack)
        {
            BindCustomers();
        }
    }
    
    protected void GridView1_RowCommand(object sender, GridViewCommandEventArgs e)
    {
        if (e.CommandName == "DeleteCustomer")
        {
            int customerId = Convert.ToInt32(e.CommandArgument);
            DeleteCustomer(customerId);
        }
    }
    
    private void BindCustomers()
    {
        try
        {
            GridView1.DataSource = _customerService.GetAll();
            GridView1.DataBind();
        }
        catch (Exception ex)
        {
            LogError("Failed to load customers", ex);
            ShowErrorMessage("Failed to load customers. Please try again.");
        }
    }
    
    private void DeleteCustomer(int id)
    {
        try
        {
            _customerService.Delete(id);
            ShowSuccessMessage("Customer deleted successfully.");
            BindCustomers();
        }
        catch (Exception ex)
        {
            LogError($"Failed to delete customer {id}", ex);
            ShowErrorMessage("Failed to delete customer.");
        }
    }
    
    private void ShowSuccessMessage(string message)
    {
        lblMessage.Text = message;
        lblMessage.CssClass = "alert alert-success";
    }
    
    private void ShowErrorMessage(string message)
    {
        lblMessage.Text = message;
        lblMessage.CssClass = "alert alert-danger";
    }
    
    private void LogError(string message, Exception ex)
    {
        // Use logging framework
        System.Diagnostics.Trace.TraceError($"{message}: {ex}");
    }
}
```

### ASPX Markup

```aspx
<%@ Page Language="C#" AutoEventWireup="true" CodeBehind="Customers.aspx.cs" 
         Inherits="MyApp.Web.Customers" %>

<!DOCTYPE html>
<html>
<head runat="server">
    <title>Customers</title>
    <link href="~/Content/bootstrap.min.css" rel="stylesheet" />
</head>
<body>
    <form id="form1" runat="server">
        <div class="container">
            <h2>Customers</h2>
            
            <asp:Label ID="lblMessage" runat="server" />
            
            <asp:GridView ID="GridView1" runat="server" 
                          CssClass="table table-striped"
                          AutoGenerateColumns="False"
                          OnRowCommand="GridView1_RowCommand">
                <Columns>
                    <asp:BoundField DataField="FirstName" HeaderText="First Name" />
                    <asp:BoundField DataField="LastName" HeaderText="Last Name" />
                    <asp:BoundField DataField="Email" HeaderText="Email" />
                    <asp:TemplateField HeaderText="Actions">
                        <ItemTemplate>
                            <asp:LinkButton ID="btnDelete" runat="server"
                                            CommandName="DeleteCustomer"
                                            CommandArgument='<%# Eval("Id") %>'
                                            Text="Delete"
                                            OnClientClick="return confirm('Are you sure?');"
                                            CssClass="btn btn-danger btn-sm" />
                        </ItemTemplate>
                    </asp:TemplateField>
                </Columns>
            </asp:GridView>
        </div>
    </form>
</body>
</html>
```

---

## 🗄️ Entity Framework 6

### DbContext Configuration

```csharp
public class ApplicationDbContext : DbContext
{
    public ApplicationDbContext() : base("name=DefaultConnection")
    {
        // Disable lazy loading for better control
        Configuration.LazyLoadingEnabled = false;
        Configuration.ProxyCreationEnabled = false;
    }
    
    public DbSet<Customer> Customers { get; set; }
    public DbSet<Order> Orders { get; set; }
    public DbSet<Product> Products { get; set; }
    
    protected override void OnModelCreating(DbModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
        
        // Apply configurations
        modelBuilder.Configurations.Add(new CustomerConfiguration());
        modelBuilder.Configurations.Add(new OrderConfiguration());
        
        // Or apply all configurations in assembly
        modelBuilder.Configurations.AddFromAssembly(typeof(ApplicationDbContext).Assembly);
    }
}
```

### Entity Configuration

```csharp
public class CustomerConfiguration : EntityTypeConfiguration<Customer>
{
    public CustomerConfiguration()
    {
        // Table mapping
        ToTable("Customers");
        
        // Primary key
        HasKey(c => c.Id);
        
        // Properties
        Property(c => c.FirstName)
            .IsRequired()
            .HasMaxLength(50);
        
        Property(c => c.LastName)
            .IsRequired()
            .HasMaxLength(50);
        
        Property(c => c.Email)
            .IsRequired()
            .HasMaxLength(100);
        
        Property(c => c.CreatedDate)
            .IsRequired()
            .HasColumnType("datetime2");
        
        // Relationships
        HasMany(c => c.Orders)
            .WithRequired(o => o.Customer)
            .HasForeignKey(o => o.CustomerId)
            .WillCascadeOnDelete(false);
    }
}
```

### Repository Pattern

```csharp
public interface IRepository<T> where T : class
{
    IQueryable<T> GetAll();
    T GetById(int id);
    void Add(T entity);
    void Update(T entity);
    void Delete(int id);
    void SaveChanges();
}

public class Repository<T> : IRepository<T> where T : class
{
    protected readonly DbContext Context;
    protected readonly DbSet<T> DbSet;
    
    public Repository(DbContext context)
    {
        Context = context ?? throw new ArgumentNullException(nameof(context));
        DbSet = context.Set<T>();
    }
    
    public virtual IQueryable<T> GetAll()
    {
        return DbSet;
    }
    
    public virtual T GetById(int id)
    {
        return DbSet.Find(id);
    }
    
    public virtual void Add(T entity)
    {
        DbSet.Add(entity);
    }
    
    public virtual void Update(T entity)
    {
        DbSet.Attach(entity);
        Context.Entry(entity).State = EntityState.Modified;
    }
    
    public virtual void Delete(int id)
    {
        var entity = GetById(id);
        if (entity != null)
        {
            DbSet.Remove(entity);
        }
    }
    
    public void SaveChanges()
    {
        Context.SaveChanges();
    }
}

// Specific repository
public class CustomerRepository : Repository<Customer>, ICustomerRepository
{
    public CustomerRepository(ApplicationDbContext context) : base(context)
    {
    }
    
    public IEnumerable<Customer> GetActiveCustomers()
    {
        return DbSet.Where(c => c.IsActive).ToList();
    }
    
    public Customer GetByEmailWithOrders(string email)
    {
        return DbSet
            .Include(c => c.Orders)
            .FirstOrDefault(c => c.Email == email);
    }
}
```

### Unit of Work Pattern

```csharp
public interface IUnitOfWork : IDisposable
{
    ICustomerRepository Customers { get; }
    IOrderRepository Orders { get; }
    int SaveChanges();
}

public class UnitOfWork : IUnitOfWork
{
    private readonly ApplicationDbContext _context;
    private ICustomerRepository _customers;
    private IOrderRepository _orders;
    
    public UnitOfWork(ApplicationDbContext context)
    {
        _context = context;
    }
    
    public ICustomerRepository Customers
    {
        get { return _customers ?? (_customers = new CustomerRepository(_context)); }
    }
    
    public IOrderRepository Orders
    {
        get { return _orders ?? (_orders = new OrderRepository(_context)); }
    }
    
    public int SaveChanges()
    {
        return _context.SaveChanges();
    }
    
    public void Dispose()
    {
        _context?.Dispose();
    }
}
```

---

## 🔌 Web API 2

### API Controller

```csharp
[RoutePrefix("api/customers")]
public class CustomersApiController : ApiController
{
    private readonly ICustomerService _customerService;
    private readonly ILogger _logger;
    
    public CustomersApiController(ICustomerService customerService, ILogger logger)
    {
        _customerService = customerService;
        _logger = logger;
    }
    
    // GET: api/customers
    [HttpGet]
    [Route("")]
    public IHttpActionResult Get()
    {
        try
        {
            var customers = _customerService.GetAll();
            return Ok(customers);
        }
        catch (Exception ex)
        {
            _logger.Error("Failed to get customers", ex);
            return InternalServerError(ex);
        }
    }
    
    // GET: api/customers/5
    [HttpGet]
    [Route("{id:int}")]
    public IHttpActionResult Get(int id)
    {
        try
        {
            var customer = _customerService.GetById(id);
            if (customer == null)
            {
                return NotFound();
            }
            return Ok(customer);
        }
        catch (Exception ex)
        {
            _logger.Error($"Failed to get customer {id}", ex);
            return InternalServerError(ex);
        }
    }
    
    // POST: api/customers
    [HttpPost]
    [Route("")]
    [ValidateModelState]
    public IHttpActionResult Post([FromBody] CustomerDto dto)
    {
        try
        {
            var customer = _customerService.Create(dto);
            return CreatedAtRoute("DefaultApi", new { id = customer.Id }, customer);
        }
        catch (Exception ex)
        {
            _logger.Error("Failed to create customer", ex);
            return InternalServerError(ex);
        }
    }
    
    // PUT: api/customers/5
    [HttpPut]
    [Route("{id:int}")]
    [ValidateModelState]
    public IHttpActionResult Put(int id, [FromBody] CustomerDto dto)
    {
        try
        {
            var customer = _customerService.GetById(id);
            if (customer == null)
            {
                return NotFound();
            }
            
            _customerService.Update(id, dto);
            return Ok();
        }
        catch (Exception ex)
        {
            _logger.Error($"Failed to update customer {id}", ex);
            return InternalServerError(ex);
        }
    }
    
    // DELETE: api/customers/5
    [HttpDelete]
    [Route("{id:int}")]
    public IHttpActionResult Delete(int id)
    {
        try
        {
            var customer = _customerService.GetById(id);
            if (customer == null)
            {
                return NotFound();
            }
            
            _customerService.Delete(id);
            return Ok();
        }
        catch (Exception ex)
        {
            _logger.Error($"Failed to delete customer {id}", ex);
            return InternalServerError(ex);
        }
    }
    
    protected override void Dispose(bool disposing)
    {
        if (disposing)
        {
            _customerService?.Dispose();
        }
        base.Dispose(disposing);
    }
}

// Model state validation attribute
public class ValidateModelStateAttribute : ActionFilterAttribute
{
    public override void OnActionExecuting(HttpActionContext actionContext)
    {
        if (!actionContext.ModelState.IsValid)
        {
            actionContext.Response = actionContext.Request.CreateErrorResponse(
                HttpStatusCode.BadRequest,
                actionContext.ModelState);
        }
    }
}
```

---

## 🛡️ Security Best Practices

### Input Validation

```csharp
// Always validate input
public ActionResult Create(CustomerViewModel model)
{
    // Model validation
    if (!ModelState.IsValid)
    {
        return View(model);
    }
    
    // Business rule validation
    if (_customerService.EmailExists(model.Email))
    {
        ModelState.AddModelError("Email", "Email already exists");
        return View(model);
    }
    
    // Sanitize input before processing
    model.FirstName = model.FirstName?.Trim();
    model.LastName = model.LastName?.Trim();
    
    _customerService.Create(model);
    return RedirectToAction("Index");
}
```

### SQL Injection Prevention

```csharp
// ✅ Good - Parameterized queries with EF
var customers = context.Customers
    .Where(c => c.Email == email)
    .ToList();

// ✅ Good - Parameterized raw SQL
var customers = context.Database.SqlQuery<Customer>(
    "SELECT * FROM Customers WHERE Email = @email",
    new SqlParameter("@email", email)).ToList();

// ❌ NEVER - String concatenation
var query = $"SELECT * FROM Customers WHERE Email = '{email}'";  // SQL INJECTION!
```

### XSS Prevention

```csharp
// In MVC, Razor automatically HTML-encodes
@Model.CustomerName  // Automatically encoded

// When you need raw HTML (be careful!)
@Html.Raw(Model.Content)  // Only use with sanitized content

// Sanitize HTML with HtmlSanitizer
using Ganss.XSS;

var sanitizer = new HtmlSanitizer();
var cleanHtml = sanitizer.Sanitize(dirtyHtml);
```

### Authentication & Authorization

```csharp
// Controller authorization
[Authorize]
public class CustomersController : Controller
{
    // All actions require authentication
}

// Action-level authorization
[Authorize(Roles = "Admin")]
public ActionResult Delete(int id)
{
    // Only admins can delete
}

// Check authorization in code
if (User.IsInRole("Admin"))
{
    // Admin-specific logic
}
```

### CSRF Protection

```csharp
// In MVC, use ValidateAntiForgeryToken
[HttpPost]
[ValidateAntiForgeryToken]
public ActionResult Create(CustomerViewModel model)
{
    // Protected against CSRF
}

// In view
@using (Html.BeginForm())
{
    @Html.AntiForgeryToken()
    // Form fields
}
```

---

## ⚡ Performance Best Practices

### Async/Await (Framework 4.5+)

```csharp
// Async controller actions
public async Task<ActionResult> Index()
{
    var customers = await _customerService.GetAllAsync();
    return View(customers);
}

// Async service methods
public async Task<IEnumerable<Customer>> GetAllAsync()
{
    return await _context.Customers.ToListAsync();
}

// Async API methods
public async Task<IHttpActionResult> Get()
{
    var customers = await _customerService.GetAllAsync();
    return Ok(customers);
}
```

### Entity Framework Optimization

```csharp
// ❌ Avoid - N+1 query problem
var customers = context.Customers.ToList();
foreach (var customer in customers)
{
    var orders = customer.Orders.ToList();  // Separate query per customer!
}

// ✅ Good - Eager loading
var customers = context.Customers
    .Include(c => c.Orders)
    .ToList();

// ✅ Good - Explicit loading
var customer = context.Customers.Find(id);
context.Entry(customer).Collection(c => c.Orders).Load();

// ✅ Good - Projection (select only needed fields)
var customerNames = context.Customers
    .Select(c => new { c.Id, c.FirstName, c.LastName })
    .ToList();

// Use AsNoTracking for read-only queries
var customers = context.Customers
    .AsNoTracking()
    .ToList();
```

### Caching

```csharp
// Output caching in MVC
[OutputCache(Duration = 60, VaryByParam = "none")]
public ActionResult Index()
{
    var customers = _customerService.GetAll();
    return View(customers);
}

// Application cache
public class CacheService
{
    private static readonly ObjectCache Cache = MemoryCache.Default;
    
    public T Get<T>(string key)
    {
        return (T)Cache.Get(key);
    }
    
    public void Set<T>(string key, T value, int minutes = 30)
    {
        var policy = new CacheItemPolicy
        {
            AbsoluteExpiration = DateTimeOffset.Now.AddMinutes(minutes)
        };
        Cache.Set(key, value, policy);
    }
    
    public void Remove(string key)
    {
        Cache.Remove(key);
    }
}
```

---

## 🧪 Testing

### Unit Testing with MSTest/NUnit

```csharp
[TestClass]
public class CustomerServiceTests
{
    private Mock<ICustomerRepository> _mockRepository;
    private CustomerService _service;
    
    [TestInitialize]
    public void Setup()
    {
        _mockRepository = new Mock<ICustomerRepository>();
        _service = new CustomerService(_mockRepository.Object);
    }
    
    [TestMethod]
    public void GetById_ExistingCustomer_ReturnsCustomer()
    {
        // Arrange
        var expectedCustomer = new Customer { Id = 1, FirstName = "John" };
        _mockRepository.Setup(r => r.GetById(1))
            .Returns(expectedCustomer);
        
        // Act
        var result = _service.GetById(1);
        
        // Assert
        Assert.IsNotNull(result);
        Assert.AreEqual(expectedCustomer.Id, result.Id);
        Assert.AreEqual(expectedCustomer.FirstName, result.FirstName);
    }
    
    [TestMethod]
    [ExpectedException(typeof(ArgumentException))]
    public void GetById_InvalidId_ThrowsException()
    {
        // Act
        _service.GetById(-1);
    }
    
    [TestCleanup]
    public void Cleanup()
    {
        _service?.Dispose();
    }
}
```

---

## 📚 Best Practices Summary

### Do's ✅

- Use async/await for I/O operations (Framework 4.5+)
- Implement proper error handling and logging
- Use parameterized queries to prevent SQL injection
- Apply [ValidateAntiForgeryToken] to POST actions
- Use dependency injection (Unity, Autofac, Ninject)
- Follow repository and unit of work patterns with EF
- Implement proper disposal (IDisposable pattern)
- Use DTOs/ViewModels, not domain models in views
- Write unit tests with mocking frameworks
- Apply caching for expensive operations

### Don'ts ❌

- Don't use string concatenation for SQL queries
- Don't store sensitive data in ViewState or Session
- Don't ignore ModelState validation
- Don't use synchronous I/O in async methods
- Don't enable lazy loading in Entity Framework
- Don't return domain entities directly from Web API
- Don't hardcode connection strings
- Don't use Server.Transfer or Response.Redirect(false)
- Don't ignore exception handling
- Don't use magic strings (use constants)

---

## 🚀 Migration Planning

### Preparing for .NET Core/.NET 6+ Migration

```csharp
// Use .NET Standard libraries where possible
// These work in both Framework and Core

// Avoid Framework-specific APIs
// ❌ Avoid
System.Web.HttpContext.Current

// ✅ Good (works in both)
using Microsoft.AspNetCore.Http;  // Use abstractions

// Document dependencies on Framework-only features
// TODO: [MIGRATION] This uses System.Web - needs refactoring for .NET Core
```

### Compatibility Analyzers

Use **.NET Portability Analyzer** to assess migration readiness:

```bash
# Install as Visual Studio extension or CLI tool
dotnet tool install -g apiport
apiport analyze -f MyApp.dll
```

---

## 📖 Resources

- [.NET Framework Design Guidelines](https://docs.microsoft.com/en-us/dotnet/standard/design-guidelines/)
- [ASP.NET MVC 5 Documentation](https://docs.microsoft.com/en-us/aspnet/mvc/overview/getting-started/)
- [Entity Framework 6 Documentation](https://docs.microsoft.com/en-us/ef/ef6/)
- [Web API 2 Documentation](https://docs.microsoft.com/en-us/aspnet/web-api/)

---

## 🏁 Maintenance Philosophy

**Remember:** .NET Framework is in **long-term support** until 2029. Focus on:

1. **Security** - Keep security patches up to date
2. **Stability** - Avoid unnecessary changes to working code
3. **Documentation** - Document complex patterns and business logic
4. **Migration Path** - Plan eventual migration to .NET 6/7/8
5. **Minimize Debt** - Don't add new technical debt to legacy code

When in doubt, **prioritize stability over innovation** in legacy applications.
