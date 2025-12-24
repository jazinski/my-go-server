# REST API Design Guidelines

> **Comprehensive REST API design standards for ASP.NET Core Web API
> applications**

---

## Table of Contents

1. [Overview](#overview)
2. [Core Principles](#core-principles)
3. [URL Design and Routing](#url-design-and-routing)
4. [HTTP Methods](#http-methods)
5. [Status Codes](#status-codes)
6. [Request and Response Formats](#request-and-response-formats)
7. [Error Handling](#error-handling)
8. [Versioning](#versioning)
9. [Authentication and Authorization](#authentication-and-authorization)
10. [Pagination, Filtering, and Sorting](#pagination-filtering-and-sorting)
11. [Performance and Caching](#performance-and-caching)
12. [Documentation](#documentation)
13. [Security](#security)
14. [Testing](#testing)
15. [Azure DevOps Integration](#azure-devops-integration)

---

## Overview

This guide establishes REST API design standards for our ASP.NET Core Web API
applications. It covers modern REST practices aligned with Microsoft
recommendations, industry standards, and our team's specific requirements for
.NET Core and Azure DevOps integration.

**Target Audience:**

- Backend developers creating new APIs
- Full-stack developers integrating React frontends with APIs
- Team leads reviewing API designs

**Technology Stack:**

- ASP.NET Core 8.0+
- Entity Framework Core
- OpenAPI/Swagger for documentation
- Azure DevOps for work item tracking

---

## Core Principles

### 1. Resource-Oriented Design

APIs should be designed around **resources** (nouns), not actions (verbs).

**✅ Good:**

```
GET    /api/customers
POST   /api/customers
GET    /api/customers/123
PUT    /api/customers/123
DELETE /api/customers/123
```

**❌ Bad:**

```
GET    /api/getCustomers
POST   /api/createCustomer
GET    /api/fetchCustomer/123
POST   /api/updateCustomer
POST   /api/deleteCustomer
```

### 2. Stateless Communication

Each API request must contain all necessary information. The server should not
rely on session state stored on the server.

### 3. Consistent Structure

All APIs should follow the same patterns for:

- URL structure
- Request/response formats
- Error handling
- Authentication

### 4. Predictable Behavior

Use HTTP methods according to their semantic meaning (GET is safe and
idempotent, POST creates resources, etc.).

### 5. Security First

All APIs must:

- Use HTTPS in production
- Require authentication by default
- Validate and sanitize all inputs
- Follow OWASP API Security Top 10

---

## URL Design and Routing

### URL Structure

**General Format:**

```
https://api.example.com/{version}/{resource-collection}/{resource-id}/{sub-resource-collection}/{sub-resource-id}
```

**Examples:**

```
https://api.example.com/v1/customers
https://api.example.com/v1/customers/123
https://api.example.com/v1/customers/123/orders
https://api.example.com/v1/customers/123/orders/456
```

### Naming Conventions

#### Resource Names

**✅ Use plural nouns:**

```
/customers (not /customer)
/orders (not /order)
/products (not /product)
```

**✅ Use kebab-case (lowercase with hyphens):**

```
/customer-orders
/product-categories
/shipping-addresses
```

**❌ Avoid:**

- camelCase: `/customerOrders`
- PascalCase: `/CustomerOrders`
- snake_case: `/customer_orders`
- Actions in URLs: `/getCustomers`, `/createOrder`

#### Route Parameters

**Use camelCase for route parameters** (case-sensitive):

```csharp
[HttpGet("{customerId}/orders/{orderId}")]
public async Task<ActionResult<OrderResponse>> GetOrder(
    int customerId,    // Matches route: camelCase
    int orderId)       // Matches route: camelCase
{
    // ...
}
```

### Attribute Routing

**✅ Always use attribute routing** for REST APIs (not conventional routing):

```csharp
[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[ApiVersion("1.0")]
public class CustomersController : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers(
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        // ...
    }

    [HttpGet("{id}")]
    public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
    {
        // ...
    }

    [HttpPost]
    public async Task<ActionResult<CustomerResponse>> CreateCustomer(
        [FromBody] CreateCustomerRequest request)
    {
        // ...
    }
}
```

### Resource Hierarchies

**Model nested resources in URLs** when there's a clear parent-child
relationship:

```
GET    /api/customers/123/orders           # Get all orders for customer 123
POST   /api/customers/123/orders           # Create order for customer 123
GET    /api/customers/123/orders/456       # Get order 456 for customer 123
DELETE /api/customers/123/orders/456       # Delete order 456 for customer 123
```

**Limit nesting to 2-3 levels** to avoid overly complex URLs:

**✅ Good:**

```
/customers/123/orders/456
```

**❌ Too complex:**

```
/customers/123/orders/456/items/789/details/abc
```

For deeply nested resources, consider:

```
GET /api/order-items/789          # Direct access
GET /api/orders/456/items          # Via parent if needed
```

---

## HTTP Methods

### GET - Retrieve Resources

**Purpose:** Retrieve one or more resources. **Safe and idempotent.**

**Single Resource:**

```csharp
[HttpGet("{id}")]
public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
{
    var customer = await _customerService.GetByIdAsync(id);
    
    if (customer == null)
        return NotFound(new ErrorResponse { Message = "Customer not found" });
    
    return Ok(customer);
}
```

**Collection:**

```csharp
[HttpGet]
public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers(
    [FromQuery] int page = 1,
    [FromQuery] int pageSize = 20,
    [FromQuery] string? search = null)
{
    var result = await _customerService.GetPagedAsync(page, pageSize, search);
    return Ok(result);
}
```

**Rules:**

- ✅ Must not modify server state
- ✅ Should be cacheable
- ✅ Should return 200 OK with response body
- ✅ Should return 404 Not Found if resource doesn't exist
- ❌ Never use GET for operations that change state

### POST - Create Resources

**Purpose:** Create a new resource. **Not idempotent** (multiple requests create
multiple resources unless using idempotency keys).

```csharp
[HttpPost]
public async Task<ActionResult<CustomerResponse>> CreateCustomer(
    [FromBody] CreateCustomerRequest request)
{
    // Validate request
    if (!ModelState.IsValid)
        return BadRequest(ModelState);
    
    var customer = await _customerService.CreateAsync(request);
    
    // Return 201 Created with Location header
    return CreatedAtAction(
        nameof(GetCustomer),
        new { id = customer.Id },
        customer);
}
```

**Rules:**

- ✅ Return 201 Created on success
- ✅ Include `Location` header with URI of created resource
- ✅ Include created resource in response body
- ✅ Return 400 Bad Request if validation fails
- ✅ Return 409 Conflict if resource already exists

**Example Response Headers:**

```
HTTP/1.1 201 Created
Location: /api/v1/customers/123
Content-Type: application/json
```

### PUT - Replace Entire Resource

**Purpose:** Replace an entire resource. **Idempotent.**

```csharp
[HttpPut("{id}")]
public async Task<ActionResult<CustomerResponse>> UpdateCustomer(
    int id,
    [FromBody] UpdateCustomerRequest request)
{
    if (!ModelState.IsValid)
        return BadRequest(ModelState);
    
    var customer = await _customerService.GetByIdAsync(id);
    if (customer == null)
        return NotFound();
    
    var updated = await _customerService.ReplaceAsync(id, request);
    return Ok(updated);
}
```

**Rules:**

- ✅ All fields should be provided (full replacement)
- ✅ Return 200 OK with updated resource
- ✅ Return 404 Not Found if resource doesn't exist
- ✅ Idempotent: Multiple identical requests have same effect

**Note:** PUT requires **all fields** of the resource. Use PATCH for partial
updates.

### PATCH - Partial Update

**Purpose:** Update specific fields of a resource. **Idempotent.**

```csharp
[HttpPatch("{id}")]
public async Task<ActionResult<CustomerResponse>> PatchCustomer(
    int id,
    [FromBody] JsonPatchDocument<CustomerDto> patchDoc)
{
    var customer = await _customerService.GetByIdAsync(id);
    if (customer == null)
        return NotFound();
    
    // Apply patch
    patchDoc.ApplyTo(customer, ModelState);
    
    if (!ModelState.IsValid)
        return BadRequest(ModelState);
    
    await _customerService.UpdateAsync(customer);
    return Ok(customer);
}
```

**Alternative (simpler) approach with partial DTOs:**

```csharp
[HttpPatch("{id}")]
public async Task<ActionResult<CustomerResponse>> PatchCustomer(
    int id,
    [FromBody] PatchCustomerRequest request)  // Contains only fields to update
{
    var customer = await _customerService.PatchAsync(id, request);
    if (customer == null)
        return NotFound();
    
    return Ok(customer);
}
```

**Rules:**

- ✅ Only specified fields are updated
- ✅ Return 200 OK with updated resource
- ✅ Return 404 Not Found if resource doesn't exist
- ✅ Use JSON Patch (RFC 6902) or simple partial DTOs

### DELETE - Remove Resource

**Purpose:** Delete a resource. **Idempotent.**

```csharp
[HttpDelete("{id}")]
public async Task<IActionResult> DeleteCustomer(int id)
{
    var customer = await _customerService.GetByIdAsync(id);
    if (customer == null)
        return NotFound();
    
    await _customerService.DeleteAsync(id);
    
    return NoContent();  // 204 No Content
}
```

**Rules:**

- ✅ Return 204 No Content on successful deletion
- ✅ Return 404 Not Found if resource doesn't exist
- ✅ Idempotent: Deleting same resource multiple times returns 404 after first
  deletion
- ✅ Consider soft deletes for audit trails

### Method Summary Table

| Method | Purpose                 | Idempotent | Safe   | Response Body    |
| ------ | ----------------------- | ---------- | ------ | ---------------- |
| GET    | Retrieve resource(s)    | ✅ Yes     | ✅ Yes | Resource data    |
| POST   | Create resource         | ❌ No      | ❌ No  | Created resource |
| PUT    | Replace entire resource | ✅ Yes     | ❌ No  | Updated resource |
| PATCH  | Partial update          | ✅ Yes     | ❌ No  | Updated resource |
| DELETE | Remove resource         | ✅ Yes     | ❌ No  | None (204)       |

---

## Status Codes

### Success Codes (2xx)

| Code    | Name       | Usage                                                    |
| ------- | ---------- | -------------------------------------------------------- |
| **200** | OK         | Successful GET, PUT, PATCH, or DELETE with response body |
| **201** | Created    | Successful POST creating a resource                      |
| **204** | No Content | Successful DELETE or update without response body        |

**Examples:**

```csharp
// 200 OK - GET with data
return Ok(customer);

// 201 Created - POST
return CreatedAtAction(nameof(GetCustomer), new { id = customer.Id }, customer);

// 204 No Content - DELETE
return NoContent();
```

### Client Error Codes (4xx)

| Code    | Name                 | Usage                                                      |
| ------- | -------------------- | ---------------------------------------------------------- |
| **400** | Bad Request          | Invalid input, validation failure, malformed JSON          |
| **401** | Unauthorized         | Missing or invalid authentication                          |
| **403** | Forbidden            | Authenticated but insufficient permissions                 |
| **404** | Not Found            | Resource does not exist                                    |
| **409** | Conflict             | Resource already exists or conflicts with current state    |
| **422** | Unprocessable Entity | Semantically invalid (valid JSON but business logic fails) |
| **429** | Too Many Requests    | Rate limit exceeded                                        |

**Examples:**

```csharp
// 400 Bad Request - Validation error
if (!ModelState.IsValid)
    return BadRequest(ModelState);

// 401 Unauthorized - Missing authentication
if (user == null)
    return Unauthorized(new ErrorResponse { Message = "Authentication required" });

// 403 Forbidden - No permission
if (!await _authService.HasPermissionAsync(user, "customers:write"))
    return Forbid();

// 404 Not Found - Resource doesn't exist
if (customer == null)
    return NotFound(new ErrorResponse { Message = "Customer not found" });

// 409 Conflict - Resource already exists
if (await _customerService.ExistsByEmailAsync(request.Email))
    return Conflict(new ErrorResponse { Message = "Customer with this email already exists" });

// 429 Too Many Requests - Rate limit
return StatusCode(429, new ErrorResponse 
{ 
    Message = "Too many requests",
    RetryAfter = "60"
});
```

### Server Error Codes (5xx)

| Code    | Name                  | Usage                                            |
| ------- | --------------------- | ------------------------------------------------ |
| **500** | Internal Server Error | Unexpected server failure                        |
| **503** | Service Unavailable   | Temporary unavailability (maintenance, overload) |

**Examples:**

```csharp
try
{
    var result = await _service.ProcessAsync(request);
    return Ok(result);
}
catch (Exception ex)
{
    _logger.LogError(ex, "Unexpected error processing request");
    
    // Don't expose internal details to client
    return StatusCode(500, new ErrorResponse 
    { 
        Message = "An error occurred processing your request",
        Code = "INTERNAL_ERROR"
    });
}
```

**Important:**

- ❌ Never return sensitive error details in production
- ✅ Log full exceptions server-side
- ✅ Return generic error messages to client
- ✅ Include error codes for client-side handling

---

## Request and Response Formats

### Content Types

**Always use JSON** for request and response bodies:

```http
Content-Type: application/json
Accept: application/json
```

### Naming Conventions

**Use camelCase for JSON properties:**

```json
{
  "customerId": 123,
  "firstName": "John",
  "lastName": "Doe",
  "emailAddress": "john.doe@example.com",
  "isActive": true,
  "createdAt": "2024-12-23T10:30:00Z"
}
```

Configure in `Program.cs`:

```csharp
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.CamelCase;
        options.JsonSerializerOptions.DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull;
    });
```

### Request DTOs

**Use dedicated request DTOs with descriptive names:**

```csharp
// Create operation
public class CreateCustomerRequest
{
    [Required]
    [StringLength(100)]
    public string FirstName { get; set; } = string.Empty;
    
    [Required]
    [StringLength(100)]
    public string LastName { get; set; } = string.Empty;
    
    [Required]
    [EmailAddress]
    public string Email { get; set; } = string.Empty;
    
    [Phone]
    public string? PhoneNumber { get; set; }
}

// Update operation
public class UpdateCustomerRequest
{
    [Required]
    [StringLength(100)]
    public string FirstName { get; set; } = string.Empty;
    
    [Required]
    [StringLength(100)]
    public string LastName { get; set; } = string.Empty;
    
    [EmailAddress]
    public string? Email { get; set; }
    
    [Phone]
    public string? PhoneNumber { get; set; }
}

// Patch operation (partial update)
public class PatchCustomerRequest
{
    [StringLength(100)]
    public string? FirstName { get; set; }
    
    [StringLength(100)]
    public string? LastName { get; set; }
    
    [EmailAddress]
    public string? Email { get; set; }
}
```

### Response DTOs

**Use dedicated response DTOs:**

```csharp
public class CustomerResponse
{
    public int Id { get; set; }
    public string FirstName { get; set; } = string.Empty;
    public string LastName { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string? PhoneNumber { get; set; }
    public bool IsActive { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    
    // Links for HATEOAS (optional)
    public Dictionary<string, string>? Links { get; set; }
}
```

**Benefits of separate DTOs:**

- ✅ Decouples API from domain models
- ✅ Prevents over-posting vulnerabilities
- ✅ Enables different validation rules for create/update
- ✅ Allows API evolution without breaking domain

### Date and Time Formats

**Always use ISO 8601 format with UTC timezone:**

```json
{
  "createdAt": "2024-12-23T10:30:00Z",
  "updatedAt": "2024-12-23T15:45:30Z"
}
```

Configure globally:

```csharp
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
        // Dates will use ISO 8601 by default
    });
```

In DTOs:

```csharp
public DateTime CreatedAt { get; set; }  // Serializes to ISO 8601
```

### Enums

**Serialize enums as strings** for readability:

```csharp
public enum OrderStatus
{
    Pending,
    Processing,
    Shipped,
    Delivered,
    Cancelled
}

// Response
{
  "orderId": 123,
  "status": "Processing"  // String, not integer
}
```

Configure:

```csharp
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(
            new JsonStringEnumConverter(JsonNamingPolicy.CamelCase));
    });
```

---

## Error Handling

### Error Response Format

**Use consistent error response structure:**

```csharp
public class ErrorResponse
{
    public string Message { get; set; } = string.Empty;
    public string Code { get; set; } = string.Empty;
    public Dictionary<string, string[]>? Errors { get; set; }
    public string? TraceId { get; set; }
}
```

### Validation Errors (400)

```csharp
[HttpPost]
public async Task<ActionResult<CustomerResponse>> CreateCustomer(
    [FromBody] CreateCustomerRequest request)
{
    if (!ModelState.IsValid)
    {
        var errors = ModelState
            .Where(x => x.Value?.Errors.Any() == true)
            .ToDictionary(
                kvp => kvp.Key,
                kvp => kvp.Value!.Errors.Select(e => e.ErrorMessage).ToArray()
            );
        
        return BadRequest(new ErrorResponse
        {
            Message = "Validation failed",
            Code = "VALIDATION_ERROR",
            Errors = errors,
            TraceId = HttpContext.TraceIdentifier
        });
    }
    
    // ...
}
```

**Example Response:**

```json
{
  "message": "Validation failed",
  "code": "VALIDATION_ERROR",
  "errors": {
    "firstName": ["First name is required"],
    "email": ["Invalid email address format"]
  },
  "traceId": "00-abc123-def456-00"
}
```

### Business Logic Errors (422)

```csharp
if (await _customerService.ExistsByEmailAsync(request.Email))
{
    return UnprocessableEntity(new ErrorResponse
    {
        Message = "Customer with this email already exists",
        Code = "DUPLICATE_EMAIL",
        TraceId = HttpContext.TraceIdentifier
    });
}
```

### Not Found Errors (404)

```csharp
var customer = await _customerService.GetByIdAsync(id);
if (customer == null)
{
    return NotFound(new ErrorResponse
    {
        Message = $"Customer with ID {id} not found",
        Code = "CUSTOMER_NOT_FOUND",
        TraceId = HttpContext.TraceIdentifier
    });
}
```

### Global Exception Handler

**Use middleware for unhandled exceptions:**

```csharp
// Program.cs
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}
else
{
    app.UseExceptionHandler("/error");
}

// ErrorController.cs
[ApiController]
[ApiExplorerSettings(IgnoreApi = true)]
public class ErrorController : ControllerBase
{
    private readonly ILogger<ErrorController> _logger;
    
    public ErrorController(ILogger<ErrorController> logger)
    {
        _logger = logger;
    }
    
    [Route("/error")]
    public IActionResult HandleError()
    {
        var context = HttpContext.Features.Get<IExceptionHandlerFeature>();
        var exception = context?.Error;
        
        _logger.LogError(exception, "Unhandled exception occurred");
        
        // Don't expose internal details in production
        return Problem(
            title: "An error occurred",
            detail: app.Environment.IsDevelopment() ? exception?.Message : null,
            statusCode: 500
        );
    }
}
```

---

## Versioning

### Versioning Strategy

**Use URL path versioning** (recommended for simplicity):

```
https://api.example.com/v1/customers
https://api.example.com/v2/customers
```

**Alternative: Support multiple versioning strategies:**

```csharp
// Program.cs
builder.Services.AddApiVersioning(options =>
{
    options.DefaultApiVersion = new ApiVersion(1, 0);
    options.AssumeDefaultVersionWhenUnspecified = true;
    options.ReportApiVersions = true;
    
    // Support multiple versioning methods
    options.ApiVersionReader = ApiVersionReader.Combine(
        new UrlSegmentApiVersionReader(),                // /api/v1/customers
        new QueryStringApiVersionReader("api-version"),  // ?api-version=1.0
        new HeaderApiVersionReader("X-Api-Version")      // Header: X-Api-Version: 1.0
    );
})
.AddMvc()
.AddApiExplorer(options =>
{
    options.GroupNameFormat = "'v'V";
    options.SubstituteApiVersionInUrl = true;
});
```

### Controller Configuration

**Option 1: Namespace-based organization (recommended):**

```csharp
// Controllers/V1/CustomersController.cs
namespace YourApi.Controllers.V1;

[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[ApiVersion("1.0")]
public class CustomersController : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers()
    {
        // V1 implementation
    }
}

// Controllers/V2/CustomersController.cs
namespace YourApi.Controllers.V2;

[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[ApiVersion("2.0")]
public class CustomersController : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PagedResponse<CustomerResponseV2>>> GetCustomers()
    {
        // V2 implementation with breaking changes
    }
}
```

**Option 2: Method-level versioning:**

```csharp
[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[ApiVersion("1.0")]
[ApiVersion("2.0")]
public class CustomersController : ControllerBase
{
    [HttpGet]
    [MapToApiVersion("1.0")]
    public async Task<ActionResult<CustomerResponse>> GetCustomersV1()
    {
        // V1 implementation
    }
    
    [HttpGet]
    [MapToApiVersion("2.0")]
    public async Task<ActionResult<CustomerResponseV2>> GetCustomersV2()
    {
        // V2 implementation
    }
}
```

### Version Deprecation

**Announce deprecation in advance:**

```csharp
[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[ApiVersion("1.0", Deprecated = true)]
[ApiVersion("2.0")]
public class CustomersController : ControllerBase
{
    // ...
}
```

Response headers will include:

```
api-supported-versions: 2.0
api-deprecated-versions: 1.0
```

### Versioning Best Practices

- ✅ Use semantic versioning (major.minor)
- ✅ Only increment major version for breaking changes
- ✅ Maintain backward compatibility within major versions
- ✅ Support at least 2 major versions simultaneously
- ✅ Announce deprecation at least 6 months before removal
- ✅ Document breaking changes in CHANGELOG and release notes
- ❌ Never change behavior of existing endpoints without version increment

---

## Authentication and Authorization

### Authentication

**Use JWT Bearer tokens** for API authentication:

```csharp
// Program.cs
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.Authority = builder.Configuration["Auth:Authority"];
        options.Audience = builder.Configuration["Auth:Audience"];
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true
        };
    });

app.UseAuthentication();
app.UseAuthorization();
```

**Require authentication by default:**

```csharp
[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[Authorize]  // Require authentication for all actions
public class CustomersController : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers()
    {
        // Only authenticated users can access
    }
    
    [AllowAnonymous]  // Explicitly allow anonymous access
    [HttpGet("public")]
    public async Task<ActionResult<IEnumerable<CustomerResponse>>> GetPublicCustomers()
    {
        // Public endpoint
    }
}
```

### Authorization

**Use policy-based authorization:**

```csharp
// Program.cs
builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("RequireAdminRole", policy =>
        policy.RequireRole("Admin"));
    
    options.AddPolicy("RequireCustomerReadPermission", policy =>
        policy.RequireClaim("permission", "customers:read"));
    
    options.AddPolicy("RequireCustomerWritePermission", policy =>
        policy.RequireClaim("permission", "customers:write"));
});

// Controller
[ApiController]
[Route("api/v{version:apiVersion}/[controller]")]
[Authorize]
public class CustomersController : ControllerBase
{
    [HttpGet]
    [Authorize(Policy = "RequireCustomerReadPermission")]
    public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers()
    {
        // Requires customers:read permission
    }
    
    [HttpPost]
    [Authorize(Policy = "RequireCustomerWritePermission")]
    public async Task<ActionResult<CustomerResponse>> CreateCustomer(
        [FromBody] CreateCustomerRequest request)
    {
        // Requires customers:write permission
    }
    
    [HttpDelete("{id}")]
    [Authorize(Policy = "RequireAdminRole")]
    public async Task<IActionResult> DeleteCustomer(int id)
    {
        // Requires Admin role
    }
}
```

### API Keys (for service-to-service)

```csharp
// Middleware for API key authentication
public class ApiKeyAuthenticationMiddleware
{
    private readonly RequestDelegate _next;
    private readonly IConfiguration _configuration;
    
    public ApiKeyAuthenticationMiddleware(
        RequestDelegate next,
        IConfiguration configuration)
    {
        _next = next;
        _configuration = configuration;
    }
    
    public async Task InvokeAsync(HttpContext context)
    {
        if (!context.Request.Headers.TryGetValue("X-API-Key", out var apiKey))
        {
            context.Response.StatusCode = 401;
            await context.Response.WriteAsJsonAsync(new ErrorResponse
            {
                Message = "API key is required",
                Code = "MISSING_API_KEY"
            });
            return;
        }
        
        var validApiKey = _configuration["ApiKey"];
        if (apiKey != validApiKey)
        {
            context.Response.StatusCode = 401;
            await context.Response.WriteAsJsonAsync(new ErrorResponse
            {
                Message = "Invalid API key",
                Code = "INVALID_API_KEY"
            });
            return;
        }
        
        await _next(context);
    }
}
```

---

## Pagination, Filtering, and Sorting

### Pagination

**Always paginate collection endpoints:**

```csharp
[HttpGet]
public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers(
    [FromQuery] int page = 1,
    [FromQuery] int pageSize = 20)
{
    // Validate
    if (page < 1) page = 1;
    if (pageSize < 1 || pageSize > 100) pageSize = 20;
    
    var (customers, totalCount) = await _customerService.GetPagedAsync(page, pageSize);
    
    var response = new PagedResponse<CustomerResponse>
    {
        Data = customers,
        Pagination = new PaginationMetadata
        {
            Page = page,
            PageSize = pageSize,
            TotalCount = totalCount,
            TotalPages = (int)Math.Ceiling(totalCount / (double)pageSize),
            HasPrevious = page > 1,
            HasNext = page < (int)Math.Ceiling(totalCount / (double)pageSize)
        }
    };
    
    return Ok(response);
}

// Response DTOs
public class PagedResponse<T>
{
    public IEnumerable<T> Data { get; set; } = new List<T>();
    public PaginationMetadata Pagination { get; set; } = new();
}

public class PaginationMetadata
{
    public int Page { get; set; }
    public int PageSize { get; set; }
    public int TotalCount { get; set; }
    public int TotalPages { get; set; }
    public bool HasPrevious { get; set; }
    public bool HasNext { get; set; }
}
```

**Example Response:**

```json
{
  "data": [
    {
      "id": 1,
      "firstName": "John",
      "lastName": "Doe"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalCount": 150,
    "totalPages": 8,
    "hasPrevious": false,
    "hasNext": true
  }
}
```

### Filtering

**Use query parameters for filtering:**

```csharp
[HttpGet]
public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers(
    [FromQuery] int page = 1,
    [FromQuery] int pageSize = 20,
    [FromQuery] string? search = null,
    [FromQuery] bool? isActive = null,
    [FromQuery] string? status = null)
{
    var (customers, totalCount) = await _customerService.GetPagedAsync(
        page, pageSize, search, isActive, status);
    
    // ...
}
```

**Example Requests:**

```
GET /api/v1/customers?isActive=true
GET /api/v1/customers?status=active&search=john
GET /api/v1/customers?isActive=true&page=2&pageSize=50
```

### Sorting

**Use `sort` query parameter:**

```csharp
[HttpGet]
public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers(
    [FromQuery] int page = 1,
    [FromQuery] int pageSize = 20,
    [FromQuery] string sort = "createdAt")  // Default sort
{
    // Parse sort parameter: "-createdAt" means descending
    var descending = sort.StartsWith("-");
    var sortField = descending ? sort[1..] : sort;
    
    var (customers, totalCount) = await _customerService.GetPagedAsync(
        page, pageSize, sortField, descending);
    
    // ...
}
```

**Example Requests:**

```
GET /api/v1/customers?sort=firstName        # Ascending
GET /api/v1/customers?sort=-createdAt       # Descending
GET /api/v1/customers?sort=-lastName,firstName  # Multiple sorts
```

### Field Selection (Sparse Fieldsets)

**Allow clients to request specific fields:**

```csharp
[HttpGet]
public async Task<ActionResult<object>> GetCustomers(
    [FromQuery] string? fields = null)
{
    var customers = await _customerService.GetAllAsync();
    
    if (string.IsNullOrEmpty(fields))
        return Ok(customers);
    
    // Parse fields: "id,firstName,email"
    var fieldList = fields.Split(',', StringSplitOptions.RemoveEmptyEntries);
    
    var result = customers.Select(c => 
        SelectFields(c, fieldList)).ToList();
    
    return Ok(result);
}
```

**Example Request:**

```
GET /api/v1/customers?fields=id,firstName,email
```

---

## Performance and Caching

### Response Caching

**Use response caching for GET endpoints:**

```csharp
// Program.cs
builder.Services.AddResponseCaching();
app.UseResponseCaching();

// Controller
[HttpGet("{id}")]
[ResponseCache(Duration = 300, Location = ResponseCacheLocation.Any)]
public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
{
    // Response cached for 5 minutes
    var customer = await _customerService.GetByIdAsync(id);
    return Ok(customer);
}
```

### ETags for Conditional Requests

```csharp
[HttpGet("{id}")]
public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
{
    var customer = await _customerService.GetByIdAsync(id);
    if (customer == null)
        return NotFound();
    
    // Generate ETag from content
    var eTag = GenerateETag(customer);
    
    // Check If-None-Match header
    if (Request.Headers.TryGetValue("If-None-Match", out var clientETag) 
        && clientETag == eTag)
    {
        return StatusCode(304);  // Not Modified
    }
    
    Response.Headers.Add("ETag", eTag);
    return Ok(customer);
}
```

### Compression

**Enable response compression:**

```csharp
// Program.cs
builder.Services.AddResponseCompression(options =>
{
    options.EnableForHttps = true;
    options.Providers.Add<GzipCompressionProvider>();
    options.Providers.Add<BrotliCompressionProvider>();
});

app.UseResponseCompression();
```

### Async/Await

**Always use async methods for I/O operations:**

```csharp
// ✅ Good - Async all the way
[HttpGet("{id}")]
public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
{
    var customer = await _customerService.GetByIdAsync(id);
    return Ok(customer);
}

// ❌ Bad - Blocking
[HttpGet("{id}")]
public ActionResult<CustomerResponse> GetCustomer(int id)
{
    var customer = _customerService.GetByIdAsync(id).Result;  // Blocks thread
    return Ok(customer);
}
```

---

## Documentation

### OpenAPI/Swagger Configuration

**Configure Swagger with versioning:**

```csharp
// Program.cs
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Customer API",
        Version = "v1",
        Description = "Customer management API",
        Contact = new OpenApiContact
        {
            Name = "API Support",
            Email = "api-support@example.com"
        }
    });
    
    // Include XML comments
    var xmlFile = $"{Assembly.GetExecutingAssembly().GetName().Name}.xml";
    var xmlPath = Path.Combine(AppContext.BaseDirectory, xmlFile);
    c.IncludeXmlComments(xmlPath);
    
    // JWT authentication
    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
    {
        Description = "JWT Authorization header using the Bearer scheme",
        Name = "Authorization",
        In = ParameterLocation.Header,
        Type = SecuritySchemeType.ApiKey,
        Scheme = "Bearer"
    });
    
    c.AddSecurityRequirement(new OpenApiSecurityRequirement
    {
        {
            new OpenApiSecurityScheme
            {
                Reference = new OpenApiReference
                {
                    Type = ReferenceType.SecurityScheme,
                    Id = "Bearer"
                }
            },
            Array.Empty<string>()
        }
    });
});

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/swagger/v1/swagger.json", "Customer API V1");
    });
}
```

### XML Documentation Comments

**Document all public APIs:**

```csharp
/// <summary>
/// Retrieves a customer by ID
/// </summary>
/// <param name="id">The customer's unique identifier</param>
/// <returns>The customer details</returns>
/// <response code="200">Customer found and returned successfully</response>
/// <response code="404">Customer not found</response>
[HttpGet("{id}")]
[ProducesResponseType(typeof(CustomerResponse), StatusCodes.Status200OK)]
[ProducesResponseType(typeof(ErrorResponse), StatusCodes.Status404NotFound)]
public async Task<ActionResult<CustomerResponse>> GetCustomer(int id)
{
    // ...
}
```

**Enable XML documentation in .csproj:**

```xml
<PropertyGroup>
    <GenerateDocumentationFile>true</GenerateDocumentationFile>
    <NoWarn>$(NoWarn);1591</NoWarn>
</PropertyGroup>
```

---

## Security

### Input Validation

**Always validate inputs:**

```csharp
public class CreateCustomerRequest
{
    [Required(ErrorMessage = "First name is required")]
    [StringLength(100, ErrorMessage = "First name cannot exceed 100 characters")]
    [RegularExpression(@"^[a-zA-Z\s'-]+$", ErrorMessage = "First name contains invalid characters")]
    public string FirstName { get; set; } = string.Empty;
    
    [Required]
    [EmailAddress(ErrorMessage = "Invalid email address")]
    public string Email { get; set; } = string.Empty;
    
    [Phone(ErrorMessage = "Invalid phone number")]
    public string? PhoneNumber { get; set; }
}
```

### HTTPS Only

**Require HTTPS in production:**

```csharp
// Program.cs
if (!app.Environment.IsDevelopment())
{
    app.UseHttpsRedirection();
    app.UseHsts();
}
```

### CORS Configuration

**Configure CORS explicitly:**

```csharp
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowFrontend", policy =>
    {
        policy.WithOrigins("https://yourdomain.com")
              .AllowAnyMethod()
              .AllowAnyHeader()
              .AllowCredentials();
    });
});

app.UseCors("AllowFrontend");
```

### Rate Limiting

**Implement rate limiting:**

```csharp
builder.Services.AddRateLimiter(options =>
{
    options.AddFixedWindowLimiter("api", opt =>
    {
        opt.Window = TimeSpan.FromMinutes(1);
        opt.PermitLimit = 100;
        opt.QueueLimit = 0;
    });
});

app.UseRateLimiter();

// Controller
[EnableRateLimiting("api")]
[ApiController]
public class CustomersController : ControllerBase
{
    // ...
}
```

### Security Headers

```csharp
app.Use(async (context, next) =>
{
    context.Response.Headers.Add("X-Content-Type-Options", "nosniff");
    context.Response.Headers.Add("X-Frame-Options", "DENY");
    context.Response.Headers.Add("X-XSS-Protection", "1; mode=block");
    await next();
});
```

For comprehensive security guidelines, see
[Security Expert Prompt](../../prompts/security-expert.md).

---

## Testing

### Unit Tests

```csharp
[Fact]
public async Task GetCustomer_ReturnsOk_WhenCustomerExists()
{
    // Arrange
    var mockService = new Mock<ICustomerService>();
    mockService.Setup(s => s.GetByIdAsync(1))
        .ReturnsAsync(new CustomerResponse { Id = 1, FirstName = "John" });
    
    var controller = new CustomersController(mockService.Object);
    
    // Act
    var result = await controller.GetCustomer(1);
    
    // Assert
    var okResult = Assert.IsType<OkObjectResult>(result.Result);
    var customer = Assert.IsType<CustomerResponse>(okResult.Value);
    Assert.Equal(1, customer.Id);
}

[Fact]
public async Task GetCustomer_ReturnsNotFound_WhenCustomerDoesNotExist()
{
    // Arrange
    var mockService = new Mock<ICustomerService>();
    mockService.Setup(s => s.GetByIdAsync(999))
        .ReturnsAsync((CustomerResponse?)null);
    
    var controller = new CustomersController(mockService.Object);
    
    // Act
    var result = await controller.GetCustomer(999);
    
    // Assert
    Assert.IsType<NotFoundObjectResult>(result.Result);
}
```

### Integration Tests

```csharp
public class CustomersControllerIntegrationTests : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient _client;
    
    public CustomersControllerIntegrationTests(WebApplicationFactory<Program> factory)
    {
        _client = factory.CreateClient();
    }
    
    [Fact]
    public async Task GetCustomers_ReturnsSuccessStatusCode()
    {
        // Act
        var response = await _client.GetAsync("/api/v1/customers");
        
        // Assert
        response.EnsureSuccessStatusCode();
        Assert.Equal("application/json", response.Content.Headers.ContentType?.MediaType);
    }
}
```

---

## Azure DevOps Integration

### Work Item Linking

**Reference work items in commit messages:**

```bash
git commit -m "feat(api): add customer pagination endpoint [AB#1234]"
```

**Link work items in API documentation:**

```csharp
/// <summary>
/// Retrieves paginated list of customers
/// </summary>
/// <remarks>
/// Work Item: AB#1234 - Add customer pagination
/// </remarks>
[HttpGet]
public async Task<ActionResult<PagedResponse<CustomerResponse>>> GetCustomers()
{
    // ...
}
```

### API Documentation in Wiki

- Document all API endpoints in Azure DevOps Wiki
- Include example requests/responses
- Link to Swagger UI for live testing
- Document breaking changes with work item references

---

## Checklist for New APIs

Before deploying a new API endpoint, verify:

- [ ] Follows RESTful design principles
- [ ] Uses appropriate HTTP methods and status codes
- [ ] Implements input validation with data annotations
- [ ] Has dedicated Request/Response DTOs
- [ ] Returns consistent error responses
- [ ] Includes authentication/authorization
- [ ] Supports pagination for collections
- [ ] Has XML documentation comments
- [ ] Appears in Swagger/OpenAPI documentation
- [ ] Has unit and integration tests
- [ ] Uses async/await for all I/O operations
- [ ] Implements appropriate caching
- [ ] Follows team naming conventions
- [ ] Linked to Azure DevOps work item
- [ ] Reviewed for security vulnerabilities (OWASP)
- [ ] Performance tested under expected load

---

## References

- [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- [ASP.NET Core Web API Documentation](https://learn.microsoft.com/en-us/aspnet/core/web-api/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [RFC 7231 - HTTP Semantics](https://datatracker.ietf.org/doc/html/rfc7231)
- [Team .NET Core Style Guide](./dotnet-core-style-guide.md)
- [Team Security Expert Guide](../../prompts/security-expert.md)

---

**Last Updated:** 2024-12-23\
**Version:** 1.0.0
