# Database Design and Conventions

> **Comprehensive database design standards for Entity Framework Core and SQL
> Server**

---

## Table of Contents

1. [Overview](#overview)
2. [Naming Conventions](#naming-conventions)
3. [Schema Design](#schema-design)
4. [Data Types](#data-types)
5. [Entity Framework Core Configuration](#entity-framework-core-configuration)
6. [Migrations](#migrations)
7. [Indexes and Performance](#indexes-and-performance)
8. [Relationships](#relationships)
9. [Queries and Performance](#queries-and-performance)
10. [Security](#security)
11. [Legacy Database Considerations](#legacy-database-considerations)
12. [Best Practices](#best-practices)

---

## Overview

This guide establishes database design standards for our applications using
Entity Framework Core with SQL Server. It covers modern EF Core patterns as well
as considerations for legacy Entity Framework 6 databases.

**Technology Stack:**

- Entity Framework Core 8.0+
- SQL Server 2019+
- Entity Framework 6.x (legacy systems)
- Azure SQL Database

**Guiding Principles:**

- **Clarity**: Names should be self-explanatory
- **Consistency**: Follow conventions across all databases
- **Performance**: Design for efficient queries and indexing
- **Maintainability**: Use migrations for version control
- **Security**: Parameterized queries, least privilege

---

## Naming Conventions

### General Rules

**Use PascalCase for all database objects:**

```sql
-- Tables
Customers
Orders
ProductCategories

-- Columns
FirstName
LastName
EmailAddress
CreatedAt

-- Stored Procedures
uspGetCustomerOrders
uspUpdateInventory
```

**Avoid:**

- ❌ snake_case: `customer_orders`
- ❌ camelCase: `customerOrders`
- ❌ Hungarian notation: `tblCustomers`, `intCustomerId`
- ❌ Abbreviations: `Cust`, `Ord`, `Prod`

### Tables

**Use singular nouns:**

```
Customer     (not Customers)
Order        (not Orders)
Product      (not Products)
OrderItem    (not OrderItems)
```

**Rationale:** Each row represents a single entity. In EF Core, the entity class
is singular, and the DbSet property is pluralized.

```csharp
public class Customer { }                    // Entity (singular)
public DbSet<Customer> Customers { get; set; }  // DbSet (plural)
```

**Use composite names for junction/join tables:**

```
OrderItem           (links Order and Product)
CustomerAddress     (links Customer and Address)
UserRole            (links User and Role)
```

### Columns

**Use descriptive PascalCase names:**

```sql
-- ✅ Good
FirstName
LastName
EmailAddress
PhoneNumber
DateOfBirth
IsActive
CreatedAt
UpdatedAt

-- ❌ Bad
FName
LName
Email
Phone
DOB
Active
CreateDate
ModDate
```

**Avoid repeating table name in column:**

```sql
-- ✅ Good
Customer.FirstName
Customer.EmailAddress

-- ❌ Bad
Customer.CustomerFirstName
Customer.CustomerEmailAddress
```

**Boolean columns:**

Use `Is`, `Has`, `Can`, or `Should` prefixes:

```sql
IsActive
IsDeleted
HasShipped
CanRefund
ShouldNotify
```

**Date/Time columns:**

Use `At` or `On` suffixes:

```sql
CreatedAt       -- DateTime with time component
UpdatedAt
DeletedAt
ScheduledAt
BirthDate       -- Date only (no time)
OrderDate
```

### Primary Keys

**Use `Id` for single-column primary keys:**

```sql
CREATE TABLE Customer (
    Id INT IDENTITY(1,1) PRIMARY KEY,
    FirstName NVARCHAR(100) NOT NULL,
    -- ...
);
```

In EF Core entities:

```csharp
public class Customer
{
    public int Id { get; set; }  // EF Core convention: property named "Id" or "{ClassName}Id"
    public string FirstName { get; set; } = string.Empty;
}
```

**For composite keys, use descriptive names:**

```sql
CREATE TABLE OrderItem (
    OrderId INT NOT NULL,
    ProductId INT NOT NULL,
    Quantity INT NOT NULL,
    PRIMARY KEY (OrderId, ProductId)
);
```

### Foreign Keys

**Use `{ReferencedTable}Id` format:**

```sql
CREATE TABLE Order (
    Id INT IDENTITY(1,1) PRIMARY KEY,
    CustomerId INT NOT NULL,
    OrderDate DATE NOT NULL,
    FOREIGN KEY (CustomerId) REFERENCES Customer(Id)
);
```

```csharp
public class Order
{
    public int Id { get; set; }
    public int CustomerId { get; set; }     // Foreign key property
    public Customer Customer { get; set; }  // Navigation property
}
```

### Indexes

**Naming pattern: `IX_{TableName}_{Column(s)}`**

```sql
CREATE INDEX IX_Customer_EmailAddress ON Customer(EmailAddress);
CREATE INDEX IX_Order_CustomerId ON Order(CustomerId);
CREATE INDEX IX_Order_OrderDate_CustomerId ON Order(OrderDate, CustomerId);
```

### Constraints

**Check constraints: `CK_{TableName}_{Column}_{Description}`**

```sql
ALTER TABLE Product
ADD CONSTRAINT CK_Product_Price_Positive CHECK (Price >= 0);

ALTER TABLE Customer
ADD CONSTRAINT CK_Customer_EmailAddress_Format CHECK (EmailAddress LIKE '%@%.%');
```

**Unique constraints: `UQ_{TableName}_{Column(s)}`**

```sql
ALTER TABLE Customer
ADD CONSTRAINT UQ_Customer_EmailAddress UNIQUE (EmailAddress);
```

**Default constraints: `DF_{TableName}_{Column}`**

```sql
ALTER TABLE Order
ADD CONSTRAINT DF_Order_Status DEFAULT ('Pending') FOR Status;
```

---

## Schema Design

### Database Schemas (SQL Server)

**Organize tables by domain using schemas:**

```sql
-- Identity/Auth
CREATE SCHEMA Auth;
GO

CREATE TABLE Auth.User ( ... );
CREATE TABLE Auth.Role ( ... );
CREATE TABLE Auth.UserRole ( ... );

-- Business domains
CREATE SCHEMA Sales;
CREATE SCHEMA Inventory;
CREATE SCHEMA Shipping;

CREATE TABLE Sales.Order ( ... );
CREATE TABLE Inventory.Product ( ... );
CREATE TABLE Shipping.Shipment ( ... );
```

In EF Core:

```csharp
[Table("User", Schema = "Auth")]
public class User
{
    public int Id { get; set; }
    // ...
}

// Or in Fluent API
modelBuilder.Entity<User>()
    .ToTable("User", "Auth");
```

### Required vs Optional Fields

**Mark required fields explicitly:**

```csharp
public class Customer
{
    public int Id { get; set; }
    
    [Required]
    [StringLength(100)]
    public string FirstName { get; set; } = string.Empty;  // Required
    
    [StringLength(100)]
    public string? MiddleName { get; set; }  // Optional (nullable)
    
    [Required]
    [EmailAddress]
    public string EmailAddress { get; set; } = string.Empty;  // Required
}
```

**In Fluent API:**

```csharp
modelBuilder.Entity<Customer>(entity =>
{
    entity.Property(e => e.FirstName)
        .IsRequired()
        .HasMaxLength(100);
    
    entity.Property(e => e.MiddleName)
        .IsRequired(false)  // Optional
        .HasMaxLength(100);
});
```

### Audit Fields

**Include standard audit fields on all tables:**

```csharp
public abstract class AuditableEntity
{
    public int Id { get; set; }
    
    [Required]
    public DateTime CreatedAt { get; set; }
    
    [Required]
    [StringLength(100)]
    public string CreatedBy { get; set; } = string.Empty;
    
    [Required]
    public DateTime UpdatedAt { get; set; }
    
    [Required]
    [StringLength(100)]
    public string UpdatedBy { get; set; } = string.Empty;
}

public class Customer : AuditableEntity
{
    public string FirstName { get; set; } = string.Empty;
    public string LastName { get; set; } = string.Empty;
    // ...
}
```

**Automatically populate audit fields:**

```csharp
public override int SaveChanges()
{
    var entries = ChangeTracker.Entries()
        .Where(e => e.Entity is AuditableEntity && 
                   (e.State == EntityState.Added || e.State == EntityState.Modified));
    
    foreach (var entry in entries)
    {
        var entity = (AuditableEntity)entry.Entity;
        var currentUser = _httpContextAccessor.HttpContext?.User?.Identity?.Name ?? "System";
        
        if (entry.State == EntityState.Added)
        {
            entity.CreatedAt = DateTime.UtcNow;
            entity.CreatedBy = currentUser;
        }
        
        entity.UpdatedAt = DateTime.UtcNow;
        entity.UpdatedBy = currentUser;
    }
    
    return base.SaveChanges();
}
```

### Soft Deletes

**Use soft deletes for important data:**

```csharp
public abstract class SoftDeletableEntity : AuditableEntity
{
    public bool IsDeleted { get; set; }
    public DateTime? DeletedAt { get; set; }
    public string? DeletedBy { get; set; }
}

// Configure global query filter
modelBuilder.Entity<Customer>()
    .HasQueryFilter(c => !c.IsDeleted);
```

**Soft delete implementation:**

```csharp
public async Task SoftDeleteAsync(int id)
{
    var customer = await _context.Customers.FindAsync(id);
    if (customer != null)
    {
        customer.IsDeleted = true;
        customer.DeletedAt = DateTime.UtcNow;
        customer.DeletedBy = _currentUser;
        await _context.SaveChangesAsync();
    }
}

// To include deleted records in a query:
var allCustomers = await _context.Customers
    .IgnoreQueryFilters()
    .ToListAsync();
```

---

## Data Types

### String Columns

**Use appropriate sizes:**

```csharp
[StringLength(100)]
public string FirstName { get; set; }         // NVARCHAR(100)

[StringLength(500)]
public string Description { get; set; }       // NVARCHAR(500)

[StringLength(4000)]
public string Notes { get; set; }             // NVARCHAR(4000)

public string FullText { get; set; }          // NVARCHAR(MAX) - use sparingly
```

**Avoid NVARCHAR(MAX) unless necessary:**

```csharp
// ✅ Good - Specific size
[StringLength(200)]
public string EmailAddress { get; set; }

// ❌ Bad - Unlimited size when not needed
public string EmailAddress { get; set; }      // Results in NVARCHAR(MAX)
```

**Use VARCHAR for ASCII-only data:**

```csharp
// For data that's guaranteed to be ASCII (e.g., codes, system identifiers)
[Column(TypeName = "varchar(50)")]
public string ProductCode { get; set; }
```

### Numeric Types

**Use appropriate precision:**

```csharp
public int Quantity { get; set; }                          // INT

public long TotalSales { get; set; }                       // BIGINT

[Column(TypeName = "decimal(18,2)")]
public decimal Price { get; set; }                         // DECIMAL(18,2) - Money

[Column(TypeName = "decimal(5,2)")]
public decimal DiscountPercent { get; set; }               // DECIMAL(5,2) - Percentage

public double Latitude { get; set; }                       // FLOAT(53)
public double Longitude { get; set; }                      // FLOAT(53)
```

**Always specify precision for decimal:**

```csharp
// ✅ Good
[Column(TypeName = "decimal(18,2)")]
public decimal Amount { get; set; }

// ❌ Bad - Uses default decimal(18,0)
public decimal Amount { get; set; }
```

### Date and Time

**Use appropriate types:**

```csharp
public DateTime CreatedAt { get; set; }        // DATETIME2(7) - Full date and time

[Column(TypeName = "date")]
public DateTime BirthDate { get; set; }        // DATE - Date only

[Column(TypeName = "time")]
public TimeSpan StartTime { get; set; }        // TIME - Time only

public DateTimeOffset ScheduledAt { get; set; } // DATETIMEOFFSET - With timezone
```

**Store times in UTC:**

```csharp
public class Order
{
    public int Id { get; set; }
    
    public DateTime OrderDate { get; set; }  // Always store in UTC
    
    // Convert to local time in application layer
    public DateTime GetLocalOrderDate(TimeZoneInfo timeZone)
    {
        return TimeZoneInfo.ConvertTimeFromUtc(OrderDate, timeZone);
    }
}

// When saving:
order.OrderDate = DateTime.UtcNow;
```

### Boolean

```csharp
public bool IsActive { get; set; }             // BIT
public bool? IsVerified { get; set; }          // BIT NULL (nullable)
```

### GUID

```csharp
public Guid Id { get; set; }                   // UNIQUEIDENTIFIER

// Auto-generate in database
[DatabaseGenerated(DatabaseGeneratedOption.Identity)]
public Guid Id { get; set; }
```

### Enums

**Store as string for readability:**

```csharp
public enum OrderStatus
{
    Pending,
    Processing,
    Shipped,
    Delivered,
    Cancelled
}

public class Order
{
    public int Id { get; set; }
    
    [Column(TypeName = "nvarchar(50)")]
    public OrderStatus Status { get; set; }
}

// Configure in Fluent API
modelBuilder.Entity<Order>()
    .Property(o => o.Status)
    .HasConversion<string>()
    .HasMaxLength(50);
```

**Alternative: Store as integer (more storage-efficient but less readable):**

```csharp
// Default EF Core behavior - stores as INT
public OrderStatus Status { get; set; }
```

---

## Entity Framework Core Configuration

### Fluent API vs Data Annotations

**Prefer Fluent API for complex configurations:**

```csharp
public class AppDbContext : DbContext
{
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        // Extract to separate configuration classes
        modelBuilder.ApplyConfiguration(new CustomerConfiguration());
        modelBuilder.ApplyConfiguration(new OrderConfiguration());
    }
}

// Separate configuration file
public class CustomerConfiguration : IEntityTypeConfiguration<Customer>
{
    public void Configure(EntityTypeBuilder<Customer> builder)
    {
        builder.ToTable("Customer");
        
        builder.HasKey(c => c.Id);
        
        builder.Property(c => c.FirstName)
            .IsRequired()
            .HasMaxLength(100);
        
        builder.Property(c => c.EmailAddress)
            .IsRequired()
            .HasMaxLength(200);
        
        builder.HasIndex(c => c.EmailAddress)
            .IsUnique()
            .HasDatabaseName("IX_Customer_EmailAddress");
        
        builder.HasMany(c => c.Orders)
            .WithOne(o => o.Customer)
            .HasForeignKey(o => o.CustomerId)
            .OnDelete(DeleteBehavior.Restrict);
    }
}
```

**Benefits of separate configuration classes:**

- ✅ Keeps `OnModelCreating` clean
- ✅ Organizes configuration by entity
- ✅ Easier to test and maintain
- ✅ Better separation of concerns

### DbContext Configuration

**Configure DbContext properly:**

```csharp
public class AppDbContext : DbContext
{
    private readonly IHttpContextAccessor _httpContextAccessor;
    
    public AppDbContext(
        DbContextOptions<AppDbContext> options,
        IHttpContextAccessor httpContextAccessor) 
        : base(options)
    {
        _httpContextAccessor = httpContextAccessor;
    }
    
    public DbSet<Customer> Customers => Set<Customer>();
    public DbSet<Order> Orders => Set<Order>();
    public DbSet<Product> Products => Set<Product>();
    
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
        
        // Apply all configurations in current assembly
        modelBuilder.ApplyConfigurationsFromAssembly(typeof(AppDbContext).Assembly);
        
        // Global query filters
        foreach (var entityType in modelBuilder.Model.GetEntityTypes())
        {
            if (typeof(ISoftDeletable).IsAssignableFrom(entityType.ClrType))
            {
                modelBuilder.Entity(entityType.ClrType)
                    .HasQueryFilter(GetSoftDeleteFilter(entityType.ClrType));
            }
        }
    }
    
    public override int SaveChanges()
    {
        UpdateAuditFields();
        return base.SaveChanges();
    }
    
    public override Task<int> SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        UpdateAuditFields();
        return base.SaveChangesAsync(cancellationToken);
    }
    
    private void UpdateAuditFields()
    {
        var entries = ChangeTracker.Entries()
            .Where(e => e.Entity is AuditableEntity &&
                       (e.State == EntityState.Added || e.State == EntityState.Modified));
        
        var currentUser = _httpContextAccessor.HttpContext?.User?.Identity?.Name ?? "System";
        var now = DateTime.UtcNow;
        
        foreach (var entry in entries)
        {
            var entity = (AuditableEntity)entry.Entity;
            
            if (entry.State == EntityState.Added)
            {
                entity.CreatedAt = now;
                entity.CreatedBy = currentUser;
            }
            
            entity.UpdatedAt = now;
            entity.UpdatedBy = currentUser;
        }
    }
}
```

### Connection Configuration

```csharp
// Program.cs
builder.Services.AddDbContext<AppDbContext>(options =>
{
    options.UseSqlServer(
        builder.Configuration.GetConnectionString("DefaultConnection"),
        sqlOptions =>
        {
            sqlOptions.EnableRetryOnFailure(
                maxRetryCount: 3,
                maxRetryDelay: TimeSpan.FromSeconds(5),
                errorNumbersToAdd: null);
            
            sqlOptions.CommandTimeout(30);  // 30 seconds
        });
    
    if (builder.Environment.IsDevelopment())
    {
        options.EnableSensitiveDataLogging();
        options.EnableDetailedErrors();
    }
});
```

---

## Migrations

### Creating Migrations

**Use descriptive migration names:**

```bash
# ✅ Good
dotnet ef migrations add AddCustomerEmailIndex
dotnet ef migrations add CreateOrderTable
dotnet ef migrations add UpdateProductPriceColumn

# ❌ Bad
dotnet ef migrations add Migration1
dotnet ef migrations add Update
dotnet ef migrations add Fix
```

### Migration Best Practices

**1. Review generated migration code before applying:**

```csharp
public partial class AddCustomerEmailIndex : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        // Review this carefully!
        migrationBuilder.CreateIndex(
            name: "IX_Customer_EmailAddress",
            table: "Customer",
            column: "EmailAddress",
            unique: true);
    }
    
    protected override void Down(MigrationBuilder migrationBuilder)
    {
        // Ensure Down reverses Up
        migrationBuilder.DropIndex(
            name: "IX_Customer_EmailAddress",
            table: "Customer");
    }
}
```

**2. Add data migrations when needed:**

```csharp
public partial class SeedDefaultRoles : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.InsertData(
            table: "Role",
            columns: new[] { "Id", "Name", "Description" },
            values: new object[,]
            {
                { 1, "Admin", "Administrator role" },
                { 2, "User", "Standard user role" },
                { 3, "Guest", "Guest user role" }
            });
    }
    
    protected override void Down(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.DeleteData(
            table: "Role",
            keyColumn: "Id",
            keyValues: new object[] { 1, 2, 3 });
    }
}
```

**3. Handle breaking changes carefully:**

```csharp
// Example: Renaming a column
public partial class RenameCustomerNameColumns : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        // Use EXEC for SQL Server rename
        migrationBuilder.Sql("EXEC sp_rename 'Customer.Name', 'FirstName', 'COLUMN'");
        
        // Add new column
        migrationBuilder.AddColumn<string>(
            name: "LastName",
            table: "Customer",
            maxLength: 100,
            nullable: false,
            defaultValue: "");
        
        // Migrate data
        migrationBuilder.Sql(@"
            UPDATE Customer
            SET LastName = ''
            WHERE LastName IS NULL
        ");
    }
    
    protected override void Down(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.Sql("EXEC sp_rename 'Customer.FirstName', 'Name', 'COLUMN'");
        migrationBuilder.DropColumn("LastName", "Customer");
    }
}
```

**4. Test migrations in non-production first:**

```bash
# Apply to development
dotnet ef database update --context AppDbContext

# Generate SQL script for production review
dotnet ef migrations script --context AppDbContext --output migration.sql

# Apply specific migration
dotnet ef database update AddCustomerEmailIndex
```

**5. Never modify applied migrations:**

```
❌ DO NOT edit migration files after they've been applied
✅ Create a new migration to make additional changes
```

### Migration Checklist

Before applying a migration:

- [ ] Migration name is descriptive
- [ ] Reviewed generated `Up` method
- [ ] Reviewed generated `Down` method
- [ ] Breaking changes documented
- [ ] Data migrations tested with sample data
- [ ] Migration tested in development environment
- [ ] SQL script reviewed by DBA (for production)
- [ ] Rollback plan documented
- [ ] Linked to Azure DevOps work item

---

## Indexes and Performance

### When to Add Indexes

**Add indexes for:**

- ✅ Primary keys (automatic)
- ✅ Foreign keys
- ✅ Columns frequently used in WHERE clauses
- ✅ Columns used in JOIN conditions
- ✅ Columns used in ORDER BY
- ✅ Unique constraints

**Example:**

```csharp
public class OrderConfiguration : IEntityTypeConfiguration<Order>
{
    public void Configure(EntityTypeBuilder<Order> builder)
    {
        // Index on foreign key
        builder.HasIndex(o => o.CustomerId)
            .HasDatabaseName("IX_Order_CustomerId");
        
        // Index on frequently queried column
        builder.HasIndex(o => o.OrderDate)
            .HasDatabaseName("IX_Order_OrderDate");
        
        // Composite index for common query
        builder.HasIndex(o => new { o.CustomerId, o.OrderDate })
            .HasDatabaseName("IX_Order_CustomerId_OrderDate");
        
        // Unique index
        builder.HasIndex(o => o.OrderNumber)
            .IsUnique()
            .HasDatabaseName("IX_Order_OrderNumber");
    }
}
```

### Composite Indexes

**Order columns by selectivity (most selective first):**

```csharp
// Good: Email is more selective than Status
builder.HasIndex(c => new { c.EmailAddress, c.Status });

// Less optimal: Status is less selective
builder.HasIndex(c => new { c.Status, c.EmailAddress });
```

### Index Maintenance

**Monitor index usage:**

```sql
-- Find unused indexes
SELECT 
    OBJECT_NAME(i.object_id) AS TableName,
    i.name AS IndexName,
    i.type_desc AS IndexType,
    s.user_seeks,
    s.user_scans,
    s.user_lookups,
    s.user_updates
FROM sys.indexes i
LEFT JOIN sys.dm_db_index_usage_stats s 
    ON i.object_id = s.object_id AND i.index_id = s.index_id
WHERE i.type_desc <> 'HEAP'
    AND OBJECTPROPERTY(i.object_id, 'IsUserTable') = 1
ORDER BY s.user_updates DESC;
```

---

## Relationships

### One-to-Many

```csharp
public class Customer
{
    public int Id { get; set; }
    public string FirstName { get; set; } = string.Empty;
    
    // Navigation property (collection)
    public ICollection<Order> Orders { get; set; } = new List<Order>();
}

public class Order
{
    public int Id { get; set; }
    
    // Foreign key
    public int CustomerId { get; set; }
    
    // Navigation property (reference)
    public Customer Customer { get; set; } = null!;
}

// Configuration
public class OrderConfiguration : IEntityTypeConfiguration<Order>
{
    public void Configure(EntityTypeBuilder<Order> builder)
    {
        builder.HasOne(o => o.Customer)
            .WithMany(c => c.Orders)
            .HasForeignKey(o => o.CustomerId)
            .OnDelete(DeleteBehavior.Restrict);  // Prevent cascade delete
    }
}
```

### Many-to-Many

**Modern EF Core (5.0+):**

```csharp
public class Student
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    
    public ICollection<Course> Courses { get; set; } = new List<Course>();
}

public class Course
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    
    public ICollection<Student> Students { get; set; } = new List<Student>();
}

// Configuration - EF Core creates junction table automatically
public class StudentConfiguration : IEntityTypeConfiguration<Student>
{
    public void Configure(EntityTypeBuilder<Student> builder)
    {
        builder.HasMany(s => s.Courses)
            .WithMany(c => c.Students)
            .UsingEntity<Dictionary<string, object>>(
                "StudentCourse",  // Junction table name
                j => j.HasOne<Course>().WithMany().HasForeignKey("CourseId"),
                j => j.HasOne<Student>().WithMany().HasForeignKey("StudentId"));
    }
}
```

**With explicit junction entity (for additional properties):**

```csharp
public class Student
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    
    public ICollection<StudentCourse> StudentCourses { get; set; } = new List<StudentCourse>();
}

public class Course
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    
    public ICollection<StudentCourse> StudentCourses { get; set; } = new List<StudentCourse>();
}

public class StudentCourse
{
    public int StudentId { get; set; }
    public Student Student { get; set; } = null!;
    
    public int CourseId { get; set; }
    public Course Course { get; set; } = null!;
    
    public DateTime EnrolledAt { get; set; }
    public string? Grade { get; set; }
}

// Configuration
public class StudentCourseConfiguration : IEntityTypeConfiguration<StudentCourse>
{
    public void Configure(EntityTypeBuilder<StudentCourse> builder)
    {
        builder.HasKey(sc => new { sc.StudentId, sc.CourseId });
        
        builder.HasOne(sc => sc.Student)
            .WithMany(s => s.StudentCourses)
            .HasForeignKey(sc => sc.StudentId);
        
        builder.HasOne(sc => sc.Course)
            .WithMany(c => c.StudentCourses)
            .HasForeignKey(sc => sc.CourseId);
    }
}
```

### One-to-One

```csharp
public class User
{
    public int Id { get; set; }
    public string Username { get; set; } = string.Empty;
    
    public UserProfile Profile { get; set; } = null!;
}

public class UserProfile
{
    public int Id { get; set; }
    
    public int UserId { get; set; }
    public User User { get; set; } = null!;
    
    public string Bio { get; set; } = string.Empty;
    public string? Website { get; set; }
}

// Configuration
public class UserProfileConfiguration : IEntityTypeConfiguration<UserProfile>
{
    public void Configure(EntityTypeBuilder<UserProfile> builder)
    {
        builder.HasOne(p => p.User)
            .WithOne(u => u.Profile)
            .HasForeignKey<UserProfile>(p => p.UserId);
    }
}
```

### Cascade Delete Behavior

```csharp
// Restrict - Prevents delete if related entities exist (RECOMMENDED)
builder.HasOne(o => o.Customer)
    .WithMany(c => c.Orders)
    .OnDelete(DeleteBehavior.Restrict);

// Cascade - Deletes related entities (use with caution)
builder.HasOne(o => o.OrderItems)
    .WithMany(oi => oi.Order)
    .OnDelete(DeleteBehavior.Cascade);

// SetNull - Sets foreign key to NULL
builder.HasOne(o => o.ShippingAddress)
    .WithMany()
    .OnDelete(DeleteBehavior.SetNull);
```

---

## Queries and Performance

### Always Use Async

```csharp
// ✅ Good - Async
var customers = await _context.Customers
    .Where(c => c.IsActive)
    .ToListAsync();

// ❌ Bad - Synchronous (blocks thread)
var customers = _context.Customers
    .Where(c => c.IsActive)
    .ToList();
```

### Use AsNoTracking for Read-Only Queries

```csharp
// Read-only query - No need to track changes
var customers = await _context.Customers
    .AsNoTracking()
    .Where(c => c.IsActive)
    .ToListAsync();

// When you need to update
var customer = await _context.Customers
    .FirstOrDefaultAsync(c => c.Id == id);  // Tracked by default

customer.EmailAddress = "new@example.com";
await _context.SaveChangesAsync();
```

### Avoid N+1 Queries - Use Include

**❌ Bad - N+1 Query Problem:**

```csharp
var customers = await _context.Customers.ToListAsync();

foreach (var customer in customers)
{
    // Each loop triggers a separate query!
    Console.WriteLine($"{customer.Name} has {customer.Orders.Count} orders");
}
```

**✅ Good - Eager Loading:**

```csharp
var customers = await _context.Customers
    .Include(c => c.Orders)
    .ToListAsync();

foreach (var customer in customers)
{
    Console.WriteLine($"{customer.Name} has {customer.Orders.Count} orders");
}
```

**✅ Better - Projection (select only needed data):**

```csharp
var customers = await _context.Customers
    .Select(c => new 
    {
        c.Id,
        c.FirstName,
        c.LastName,
        OrderCount = c.Orders.Count
    })
    .ToListAsync();
```

### Use Projections Instead of Loading Full Entities

```csharp
// ❌ Bad - Loads all columns
var customers = await _context.Customers
    .Include(c => c.Orders)
    .ToListAsync();

// ✅ Good - Only load needed fields
var customers = await _context.Customers
    .Select(c => new CustomerDto
    {
        Id = c.Id,
        FullName = c.FirstName + " " + c.LastName,
        EmailAddress = c.EmailAddress,
        OrderCount = c.Orders.Count
    })
    .ToListAsync();
```

### Split Complex Queries

**For complex includes with multiple levels:**

```csharp
// ❌ May cause Cartesian explosion
var orders = await _context.Orders
    .Include(o => o.OrderItems)
    .Include(o => o.Customer)
    .Include(o => o.ShippingAddress)
    .ToListAsync();

// ✅ Split query
var orders = await _context.Orders
    .Include(o => o.OrderItems)
    .Include(o => o.Customer)
    .Include(o => o.ShippingAddress)
    .AsSplitQuery()
    .ToListAsync();
```

### Compiled Queries for Frequent Queries

```csharp
private static readonly Func<AppDbContext, int, Task<Customer?>> _getCustomerById =
    EF.CompileAsyncQuery((AppDbContext context, int id) =>
        context.Customers.FirstOrDefault(c => c.Id == id));

public async Task<Customer?> GetByIdAsync(int id)
{
    return await _getCustomerById(_context, id);
}
```

### DbContext Pooling

```csharp
// Program.cs - Improves performance for high-throughput scenarios
builder.Services.AddDbContextPool<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")),
    poolSize: 128);
```

### Query Performance Monitoring

```csharp
// Log slow queries in development
builder.Services.AddDbContext<AppDbContext>(options =>
{
    options.UseSqlServer(connectionString);
    
    if (builder.Environment.IsDevelopment())
    {
        options.LogTo(
            Console.WriteLine,
            new[] { DbLoggerCategory.Database.Command.Name },
            LogLevel.Information)
        .EnableSensitiveDataLogging();
    }
});
```

---

## Security

### Always Use Parameterized Queries

**✅ EF Core automatically parameterizes:**

```csharp
// Safe - Parameterized
var email = "user@example.com";
var customer = await _context.Customers
    .FirstOrDefaultAsync(c => c.EmailAddress == email);
```

**❌ Never use string concatenation:**

```csharp
// UNSAFE - SQL Injection vulnerability!
var email = "user@example.com";
var customers = await _context.Customers
    .FromSqlRaw($"SELECT * FROM Customer WHERE EmailAddress = '{email}'")
    .ToListAsync();
```

**✅ Use FromSqlRaw with parameters:**

```csharp
// Safe - Parameterized
var email = "user@example.com";
var customers = await _context.Customers
    .FromSqlRaw("SELECT * FROM Customer WHERE EmailAddress = {0}", email)
    .ToListAsync();
```

### Connection String Security

**Never hardcode connection strings:**

```csharp
// ❌ Bad
var connectionString = "Server=prod-db;Database=MyDb;User=sa;Password=P@ssw0rd!";

// ✅ Good - Use configuration
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");

// ✅ Better - Use Azure Key Vault or Managed Identity
builder.Configuration.AddAzureKeyVault(
    new Uri(builder.Configuration["KeyVaultUrl"]!),
    new DefaultAzureCredential());
```

**Connection strings in appsettings.json:**

```json
{
  "ConnectionStrings": {
    "DefaultConnection": "Server=(localdb)\\mssqllocaldb;Database=MyDb;Trusted_Connection=True;"
  }
}
```

**For production, use environment variables or Azure:**

```bash
# Environment variable
export ConnectionStrings__DefaultConnection="Server=prod;Database=MyDb;..."

# Azure App Service configuration
az webapp config connection-string set \
  --name myapp \
  --resource-group mygroup \
  --connection-string-type SQLAzure \
  --settings DefaultConnection="Server=..."
```

### Least Privilege

**Database user should have minimum required permissions:**

```sql
-- Create application-specific user
CREATE LOGIN [MyAppUser] WITH PASSWORD = 'ComplexPassword123!';
CREATE USER [MyAppUser] FOR LOGIN [MyAppUser];

-- Grant only necessary permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [MyAppUser];
DENY ALTER, DROP ON SCHEMA::dbo TO [MyAppUser];

-- For read-only operations
CREATE LOGIN [MyAppReadOnlyUser] WITH PASSWORD = 'ComplexPassword123!';
CREATE USER [MyAppReadOnlyUser] FOR LOGIN [MyAppReadOnlyUser];
GRANT SELECT ON SCHEMA::dbo TO [MyAppReadOnlyUser];
```

---

## Legacy Database Considerations

### Working with Entity Framework 6

**Key differences from EF Core:**

```csharp
// EF6
public class LegacyDbContext : DbContext
{
    public LegacyDbContext() : base("name=DefaultConnection") { }
    
    public DbSet<Customer> Customers { get; set; }
    
    protected override void OnModelCreating(DbModelBuilder modelBuilder)
    {
        // Different API than EF Core
        modelBuilder.Entity<Customer>()
            .Property(c => c.FirstName)
            .IsRequired()
            .HasMaxLength(100);
    }
}
```

**Migration path to EF Core:**

1. Create new EF Core DbContext alongside EF6
2. Use same database schema
3. Migrate queries gradually
4. Test thoroughly before removing EF6

### Handling Legacy Naming Conventions

**Map EF Core entities to legacy tables:**

```csharp
[Table("tblCustomers")]  // Legacy table name
public class Customer
{
    [Column("customer_id")]
    public int Id { get; set; }
    
    [Column("first_name")]
    public string FirstName { get; set; } = string.Empty;
    
    [Column("email_addr")]
    public string EmailAddress { get; set; } = string.Empty;
}
```

### Database Views

**Map to views for complex legacy queries:**

```csharp
[Keyless]
[Table("vw_CustomerOrderSummary")]
public class CustomerOrderSummary
{
    public int CustomerId { get; set; }
    public string CustomerName { get; set; } = string.Empty;
    public int TotalOrders { get; set; }
    public decimal TotalSpent { get; set; }
}

// Query
var summary = await _context.Set<CustomerOrderSummary>()
    .AsNoTracking()
    .ToListAsync();
```

---

## Best Practices

### Summary Checklist

#### Design

- [ ] Use PascalCase for all database objects
- [ ] Use singular table names
- [ ] Include audit fields (CreatedAt, UpdatedAt, CreatedBy, UpdatedBy)
- [ ] Use soft deletes for important data
- [ ] Choose appropriate data types and sizes
- [ ] Define relationships explicitly in Fluent API

#### Configuration

- [ ] Use `IEntityTypeConfiguration<T>` for entity configuration
- [ ] Keep DbContext clean - apply configurations from assembly
- [ ] Configure global query filters
- [ ] Implement audit field auto-population
- [ ] Enable retry logic for transient failures

#### Migrations

- [ ] Use descriptive migration names
- [ ] Review generated migrations before applying
- [ ] Test migrations in development first
- [ ] Generate SQL scripts for production review
- [ ] Never modify applied migrations

#### Performance

- [ ] Always use async/await
- [ ] Use AsNoTracking() for read-only queries
- [ ] Use Include() for eager loading (avoid N+1)
- [ ] Use projections to select only needed fields
- [ ] Add indexes on foreign keys and frequently queried columns
- [ ] Consider DbContext pooling for high-throughput scenarios

#### Security

- [ ] Always use parameterized queries
- [ ] Store connection strings in configuration
- [ ] Use least privilege for database users
- [ ] Validate all user inputs
- [ ] Log security-relevant events

---

## References

- [Entity Framework Core Documentation](https://learn.microsoft.com/en-us/ef/core/)
- [EF Core Performance Best Practices](https://learn.microsoft.com/en-us/ef/core/performance/)
- [SQL Server Best Practices](https://learn.microsoft.com/en-us/sql/sql-server/)
- [Team .NET Core Style Guide](./dotnet-core-style-guide.md)
- [Team API Design Guide](./api-design-guide.md)

---

**Last Updated:** 2024-12-23\
**Version:** 1.0.0
