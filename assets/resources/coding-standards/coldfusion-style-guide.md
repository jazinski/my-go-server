# ColdFusion (CFML) Coding Standards - Legacy Maintenance

## 🎯 Overview

This guide covers **Adobe ColdFusion** and **Lucee CFML** best practices for
maintaining legacy applications. ColdFusion remains widely used for enterprise
applications, though new development has declined. Focus on maintainable, secure
code that follows modern CFML patterns.

**Supported Versions:**

- Adobe ColdFusion 2016, 2018, 2021, 2023
- Lucee 5.x, 6.x
- ColdFusion 11 (legacy, end-of-life)

**Key Principles:**

- Prefer **CFScript** over tag-based syntax for business logic
- Follow **component-based architecture** (CFCs)
- Implement proper **error handling** and logging
- Use **parameterized queries** to prevent SQL injection
- Document complex business logic
- Plan for modernization or migration when appropriate

---

## 📁 Project Structure

### Directory Organization

```
webroot/
├── Application.cfc          # Application configuration
├── index.cfm                # Entry point
├── components/              # CFCs (business logic)
│   ├── Customer.cfc
│   ├── Order.cfc
│   └── services/
│       ├── CustomerService.cfc
│       └── OrderService.cfc
├── models/                  # Data access components
│   ├── CustomerDAO.cfc
│   └── OrderDAO.cfc
├── views/                   # Display templates
│   ├── customers/
│   │   ├── list.cfm
│   │   └── edit.cfm
│   └── layouts/
│       └── main.cfm
├── includes/                # Reusable includes
│   ├── header.cfm
│   └── footer.cfm
├── assets/                  # Static files
│   ├── css/
│   ├── js/
│   └── images/
└── config/                  # Configuration files
    └── settings.cfm
```

### Naming Conventions

```coldfusion
// Files
Customer.cfc              // PascalCase for components
CustomerService.cfc       // Service suffix
CustomerDAO.cfc          // DAO suffix
customer-list.cfm        // kebab-case for templates

// Variables
var firstName = "John";           // camelCase
var customerId = 123;
var isActive = true;

// Constants (Application/Server scope)
application.MAX_LOGIN_ATTEMPTS = 3;  // SCREAMING_SNAKE_CASE

// Components and functions
component Customer {}                 // PascalCase
function getCustomerById() {}        // camelCase

// Database
Customers                            // PascalCase table names
CustomerId                           // PascalCase column names
FirstName
```

---

## 🔧 Application.cfc

### Modern Application Configuration

```coldfusion
component {
    // Application settings
    this.name = "MyApplication";
    this.applicationTimeout = createTimeSpan(0, 2, 0, 0);  // 2 hours
    this.sessionManagement = true;
    this.sessionTimeout = createTimeSpan(0, 0, 30, 0);     // 30 minutes
    this.setClientCookies = true;
    this.sessionCookie.httpOnly = true;
    this.sessionCookie.secure = true;  // HTTPS only
    
    // Data source
    this.datasource = "myAppDB";
    
    // Mappings
    this.mappings = {
        "/components" = expandPath("./components"),
        "/models" = expandPath("./models")
    };
    
    // Custom tag paths
    this.customTagPaths = expandPath("./customtags");
    
    // ORM settings (if using ORM)
    this.ormEnabled = false;  // Only if actively used
    
    // Application initialization
    function onApplicationStart() {
        // Initialize application-scope variables
        application.settings = {};
        application.settings.siteName = "My Application";
        application.settings.maxLoginAttempts = 3;
        
        // Load configuration
        include "config/settings.cfm";
        
        return true;
    }
    
    // Request initialization
    function onRequestStart(targetPage) {
        // Reload application (for development only)
        if (structKeyExists(url, "reinit") && url.reinit eq "true") {
            applicationStop();
            location(url=cgi.script_name, addToken=false);
        }
        
        // Security checks
        if (!isUserLoggedIn() && !isPublicPage(targetPage)) {
            location(url="/login.cfm", addToken=false);
            return false;
        }
        
        return true;
    }
    
    // Session initialization
    function onSessionStart() {
        session.user = {};
        session.user.isAuthenticated = false;
        session.user.id = 0;
        session.user.role = "";
    }
    
    // Error handling
    function onError(exception, eventName) {
        // Log error
        writeLog(
            file = "application",
            type = "error",
            text = "Error in #eventName#: #exception.message#"
        );
        
        // Display user-friendly error
        if (application.environment eq "production") {
            include "views/errors/500.cfm";
        } else {
            writeDump(var=exception, label="Error Details");
            abort;
        }
    }
    
    // Missing template handler
    function onMissingTemplate(targetPage) {
        writeLog(
            file = "application",
            type = "warning",
            text = "Missing template: #targetPage#"
        );
        include "views/errors/404.cfm";
        return true;
    }
    
    // Helper functions
    private boolean function isUserLoggedIn() {
        return (
            structKeyExists(session, "user") && 
            session.user.isAuthenticated
        );
    }
    
    private boolean function isPublicPage(targetPage) {
        var publicPages = [
            "/login.cfm",
            "/register.cfm",
            "/forgot-password.cfm"
        ];
        return arrayContains(publicPages, targetPage);
    }
}
```

---

## 🏗️ Components (CFCs)

### CFScript Component Syntax (Preferred)

```coldfusion
/**
 * Customer component
 * Handles customer business logic
 * 
 * @author Your Name
 * @date 2025-12-23
 */
component output="false" {
    
    // Properties
    property name="customerId" type="numeric" default="0";
    property name="firstName" type="string";
    property name="lastName" type="string";
    property name="email" type="string";
    property name="isActive" type="boolean" default="true";
    
    // Constructor
    public Customer function init(
        numeric customerId = 0,
        string firstName = "",
        string lastName = "",
        string email = ""
    ) {
        variables.customerId = arguments.customerId;
        variables.firstName = arguments.firstName;
        variables.lastName = arguments.lastName;
        variables.email = arguments.email;
        variables.isActive = true;
        
        return this;
    }
    
    // Getters and Setters
    public numeric function getCustomerId() {
        return variables.customerId;
    }
    
    public void function setCustomerId(required numeric customerId) {
        variables.customerId = arguments.customerId;
    }
    
    public string function getFirstName() {
        return variables.firstName;
    }
    
    public void function setFirstName(required string firstName) {
        variables.firstName = trim(arguments.firstName);
    }
    
    public string function getEmail() {
        return variables.email;
    }
    
    public void function setEmail(required string email) {
        if (!isValid("email", arguments.email)) {
            throw(type="ValidationException", message="Invalid email address");
        }
        variables.email = arguments.email;
    }
    
    // Business methods
    public string function getFullName() {
        return "#variables.firstName# #variables.lastName#";
    }
    
    public boolean function isActive() {
        return variables.isActive;
    }
    
    public void function activate() {
        variables.isActive = true;
    }
    
    public void function deactivate() {
        variables.isActive = false;
    }
    
    // Validation
    public boolean function isValid() {
        return (
            len(trim(variables.firstName)) > 0 &&
            len(trim(variables.lastName)) > 0 &&
            isValid("email", variables.email)
        );
    }
}
```

### Service Layer Pattern

```coldfusion
/**
 * Customer Service
 * Business logic for customer operations
 */
component output="false" {
    
    // Dependencies (inject via constructor)
    property name="customerDAO" inject="CustomerDAO";
    property name="logger" inject="Logger";
    
    /**
     * Initialize service
     */
    public CustomerService function init(
        required CustomerDAO customerDAO,
        required Logger logger
    ) {
        variables.customerDAO = arguments.customerDAO;
        variables.logger = arguments.logger;
        return this;
    }
    
    /**
     * Get all active customers
     */
    public array function getAllActive() {
        try {
            return variables.customerDAO.getActive();
        } catch (any e) {
            variables.logger.error("Failed to get active customers", e);
            rethrow;
        }
    }
    
    /**
     * Get customer by ID
     */
    public Customer function getById(required numeric customerId) {
        if (arguments.customerId <= 0) {
            throw(
                type = "ValidationException",
                message = "Invalid customer ID"
            );
        }
        
        try {
            return variables.customerDAO.getById(arguments.customerId);
        } catch (any e) {
            variables.logger.error("Failed to get customer #arguments.customerId#", e);
            rethrow;
        }
    }
    
    /**
     * Create new customer
     */
    public numeric function create(required struct customerData) {
        // Validate input
        validateCustomerData(arguments.customerData);
        
        // Check for duplicate email
        if (variables.customerDAO.emailExists(arguments.customerData.email)) {
            throw(
                type = "DuplicateException",
                message = "Email address already exists"
            );
        }
        
        try {
            return variables.customerDAO.create(arguments.customerData);
        } catch (any e) {
            variables.logger.error("Failed to create customer", e);
            rethrow;
        }
    }
    
    /**
     * Update existing customer
     */
    public void function update(
        required numeric customerId,
        required struct customerData
    ) {
        // Validate input
        validateCustomerData(arguments.customerData);
        
        // Ensure customer exists
        var existingCustomer = getById(arguments.customerId);
        
        try {
            variables.customerDAO.update(arguments.customerId, arguments.customerData);
            variables.logger.info("Updated customer #arguments.customerId#");
        } catch (any e) {
            variables.logger.error("Failed to update customer #arguments.customerId#", e);
            rethrow;
        }
    }
    
    /**
     * Delete customer
     */
    public void function delete(required numeric customerId) {
        // Ensure customer exists
        var existingCustomer = getById(arguments.customerId);
        
        try {
            variables.customerDAO.delete(arguments.customerId);
            variables.logger.info("Deleted customer #arguments.customerId#");
        } catch (any e) {
            variables.logger.error("Failed to delete customer #arguments.customerId#", e);
            rethrow;
        }
    }
    
    /**
     * Validate customer data
     */
    private void function validateCustomerData(required struct data) {
        var errors = [];
        
        if (!structKeyExists(arguments.data, "firstName") || 
            len(trim(arguments.data.firstName)) == 0) {
            arrayAppend(errors, "First name is required");
        }
        
        if (!structKeyExists(arguments.data, "lastName") || 
            len(trim(arguments.data.lastName)) == 0) {
            arrayAppend(errors, "Last name is required");
        }
        
        if (!structKeyExists(arguments.data, "email") || 
            !isValid("email", arguments.data.email)) {
            arrayAppend(errors, "Valid email is required");
        }
        
        if (arrayLen(errors) > 0) {
            throw(
                type = "ValidationException",
                message = arrayToList(errors, "; ")
            );
        }
    }
}
```

### Data Access Object (DAO) Pattern

```coldfusion
/**
 * Customer DAO
 * Database operations for customers
 */
component output="false" {
    
    // Data source name
    variables.dsn = "myAppDB";
    
    /**
     * Get all active customers
     */
    public array function getActive() {
        var sql = "
            SELECT CustomerId, FirstName, LastName, Email, IsActive
            FROM Customers
            WHERE IsActive = :isActive
            ORDER BY LastName, FirstName
        ";
        
        var qCustomers = queryExecute(
            sql,
            { isActive = { value = 1, cfsqltype = "cf_sql_bit" } },
            { datasource = variables.dsn }
        );
        
        return queryToArray(qCustomers);
    }
    
    /**
     * Get customer by ID
     */
    public struct function getById(required numeric customerId) {
        var sql = "
            SELECT CustomerId, FirstName, LastName, Email, IsActive, CreatedDate
            FROM Customers
            WHERE CustomerId = :customerId
        ";
        
        var qCustomer = queryExecute(
            sql,
            { customerId = { value = arguments.customerId, cfsqltype = "cf_sql_integer" } },
            { datasource = variables.dsn }
        );
        
        if (qCustomer.recordCount == 0) {
            throw(
                type = "NotFoundException",
                message = "Customer not found"
            );
        }
        
        return queryRowToStruct(qCustomer, 1);
    }
    
    /**
     * Create new customer
     */
    public numeric function create(required struct customerData) {
        var sql = "
            INSERT INTO Customers (FirstName, LastName, Email, IsActive, CreatedDate)
            VALUES (:firstName, :lastName, :email, :isActive, :createdDate)
        ";
        
        var params = {
            firstName = { value = arguments.customerData.firstName, cfsqltype = "cf_sql_varchar" },
            lastName = { value = arguments.customerData.lastName, cfsqltype = "cf_sql_varchar" },
            email = { value = arguments.customerData.email, cfsqltype = "cf_sql_varchar" },
            isActive = { value = 1, cfsqltype = "cf_sql_bit" },
            createdDate = { value = now(), cfsqltype = "cf_sql_timestamp" }
        };
        
        var result = queryExecute(
            sql,
            params,
            { datasource = variables.dsn, result = "insertResult" }
        );
        
        return insertResult.generatedKey;
    }
    
    /**
     * Update existing customer
     */
    public void function update(
        required numeric customerId,
        required struct customerData
    ) {
        var sql = "
            UPDATE Customers
            SET FirstName = :firstName,
                LastName = :lastName,
                Email = :email,
                ModifiedDate = :modifiedDate
            WHERE CustomerId = :customerId
        ";
        
        var params = {
            customerId = { value = arguments.customerId, cfsqltype = "cf_sql_integer" },
            firstName = { value = arguments.customerData.firstName, cfsqltype = "cf_sql_varchar" },
            lastName = { value = arguments.customerData.lastName, cfsqltype = "cf_sql_varchar" },
            email = { value = arguments.customerData.email, cfsqltype = "cf_sql_varchar" },
            modifiedDate = { value = now(), cfsqltype = "cf_sql_timestamp" }
        };
        
        queryExecute(sql, params, { datasource = variables.dsn });
    }
    
    /**
     * Delete customer
     */
    public void function delete(required numeric customerId) {
        var sql = "DELETE FROM Customers WHERE CustomerId = :customerId";
        
        queryExecute(
            sql,
            { customerId = { value = arguments.customerId, cfsqltype = "cf_sql_integer" } },
            { datasource = variables.dsn }
        );
    }
    
    /**
     * Check if email exists
     */
    public boolean function emailExists(required string email) {
        var sql = "SELECT COUNT(*) as cnt FROM Customers WHERE Email = :email";
        
        var qCheck = queryExecute(
            sql,
            { email = { value = arguments.email, cfsqltype = "cf_sql_varchar" } },
            { datasource = variables.dsn }
        );
        
        return (qCheck.cnt > 0);
    }
    
    /**
     * Helper: Convert query row to struct
     */
    private struct function queryRowToStruct(required query qry, required numeric row) {
        var result = {};
        var columns = listToArray(arguments.qry.columnList);
        
        for (var col in columns) {
            result[col] = arguments.qry[col][arguments.row];
        }
        
        return result;
    }
    
    /**
     * Helper: Convert query to array of structs
     */
    private array function queryToArray(required query qry) {
        var result = [];
        
        for (var i = 1; i <= arguments.qry.recordCount; i++) {
            arrayAppend(result, queryRowToStruct(arguments.qry, i));
        }
        
        return result;
    }
}
```

---

## 🎨 CFML Templates

### Modern CFM Template Structure

```coldfusion
<!--- customer-list.cfm --->
<cfscript>
    // Controller logic
    param name="url.page" default="1";
    param name="url.search" default="";
    
    // Initialize service
    customerService = new components.services.CustomerService(
        customerDAO = new models.CustomerDAO(),
        logger = application.logger
    );
    
    // Get data
    try {
        customers = customerService.getAllActive();
        pageTitle = "Customer List";
    } catch (any e) {
        writeLog(file="application", type="error", text=e.message);
        customers = [];
        errorMessage = "Failed to load customers. Please try again.";
    }
</cfscript>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title><cfoutput>#pageTitle#</cfoutput></title>
    <link rel="stylesheet" href="/assets/css/main.css">
</head>
<body>
    <cfinclude template="/includes/header.cfm">
    
    <main class="container">
        <h1><cfoutput>#pageTitle#</cfoutput></h1>
        
        <!--- Error message --->
        <cfif isDefined("errorMessage")>
            <div class="alert alert-error">
                <cfoutput>#encodeForHTML(errorMessage)#</cfoutput>
            </div>
        </cfif>
        
        <!--- Success message --->
        <cfif structKeyExists(session, "successMessage")>
            <div class="alert alert-success">
                <cfoutput>#encodeForHTML(session.successMessage)#</cfoutput>
            </div>
            <cfset structDelete(session, "successMessage")>
        </cfif>
        
        <!--- Customer table --->
        <table class="table">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Status</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                <cfloop array="#customers#" index="customer">
                    <tr>
                        <td>
                            <cfoutput>
                                #encodeForHTML(customer.firstName)# 
                                #encodeForHTML(customer.lastName)#
                            </cfoutput>
                        </td>
                        <td>
                            <cfoutput>#encodeForHTML(customer.email)#</cfoutput>
                        </td>
                        <td>
                            <cfif customer.isActive>
                                <span class="badge badge-success">Active</span>
                            <cfelse>
                                <span class="badge badge-secondary">Inactive</span>
                            </cfif>
                        </td>
                        <td>
                            <cfoutput>
                                <a href="customer-edit.cfm?id=#customer.customerId#">Edit</a> |
                                <a href="customer-delete.cfm?id=#customer.customerId#" 
                                   onclick="return confirm('Are you sure?')">Delete</a>
                            </cfoutput>
                        </td>
                    </tr>
                </cfloop>
            </tbody>
        </table>
    </main>
    
    <cfinclude template="/includes/footer.cfm">
</body>
</html>
```

### Form Handling

```coldfusion
<!--- customer-edit.cfm --->
<cfscript>
    // Get customer ID
    param name="url.id" default="0" type="numeric";
    
    // Initialize service
    customerService = new components.services.CustomerService(
        customerDAO = new models.CustomerDAO(),
        logger = application.logger
    );
    
    // Handle form submission
    if (structKeyExists(form, "submit")) {
        // CSRF protection
        if (!csrfVerifyToken(form.csrfToken)) {
            errorMessage = "Invalid security token. Please try again.";
        } else {
            // Validate and save
            try {
                customerData = {
                    firstName = trim(form.firstName),
                    lastName = trim(form.lastName),
                    email = trim(form.email)
                };
                
                if (url.id > 0) {
                    customerService.update(url.id, customerData);
                    session.successMessage = "Customer updated successfully.";
                } else {
                    newId = customerService.create(customerData);
                    session.successMessage = "Customer created successfully.";
                }
                
                location(url="customer-list.cfm", addToken=false);
            } catch (ValidationException e) {
                errorMessage = e.message;
            } catch (any e) {
                writeLog(file="application", type="error", text=e.message);
                errorMessage = "Failed to save customer. Please try again.";
            }
        }
    }
    
    // Load existing customer or initialize empty
    if (url.id > 0) {
        try {
            customer = customerService.getById(url.id);
            pageTitle = "Edit Customer";
        } catch (any e) {
            location(url="customer-list.cfm", addToken=false);
        }
    } else {
        customer = {
            customerId = 0,
            firstName = "",
            lastName = "",
            email = ""
        };
        pageTitle = "Create Customer";
    }
    
    // Generate CSRF token
    csrfToken = csrfGenerateToken();
</cfscript>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title><cfoutput>#pageTitle#</cfoutput></title>
    <link rel="stylesheet" href="/assets/css/main.css">
</head>
<body>
    <cfinclude template="/includes/header.cfm">
    
    <main class="container">
        <h1><cfoutput>#pageTitle#</cfoutput></h1>
        
        <cfif isDefined("errorMessage")>
            <div class="alert alert-error">
                <cfoutput>#encodeForHTML(errorMessage)#</cfoutput>
            </div>
        </cfif>
        
        <form method="post" action="customer-edit.cfm?id=#url.id#">
            <input type="hidden" name="csrfToken" value="<cfoutput>#csrfToken#</cfoutput>">
            
            <div class="form-group">
                <label for="firstName">First Name *</label>
                <input type="text" 
                       id="firstName" 
                       name="firstName" 
                       value="<cfoutput>#encodeForHTMLAttribute(customer.firstName)#</cfoutput>"
                       required
                       maxlength="50">
            </div>
            
            <div class="form-group">
                <label for="lastName">Last Name *</label>
                <input type="text" 
                       id="lastName" 
                       name="lastName" 
                       value="<cfoutput>#encodeForHTMLAttribute(customer.lastName)#</cfoutput>"
                       required
                       maxlength="50">
            </div>
            
            <div class="form-group">
                <label for="email">Email *</label>
                <input type="email" 
                       id="email" 
                       name="email" 
                       value="<cfoutput>#encodeForHTMLAttribute(customer.email)#</cfoutput>"
                       required
                       maxlength="100">
            </div>
            
            <div class="form-actions">
                <button type="submit" name="submit" class="btn btn-primary">Save</button>
                <a href="customer-list.cfm" class="btn btn-secondary">Cancel</a>
            </div>
        </form>
    </main>
    
    <cfinclude template="/includes/footer.cfm">
</body>
</html>
```

---

## 🛡️ Security Best Practices

### SQL Injection Prevention

```coldfusion
// ❌ NEVER - String concatenation (SQL INJECTION!)
<cfquery name="qCustomer" datasource="myAppDB">
    SELECT * FROM Customers
    WHERE Email = '#form.email#'  <!--- DANGEROUS! --->
</cfquery>

// ✅ ALWAYS - Use cfqueryparam
<cfquery name="qCustomer" datasource="myAppDB">
    SELECT * FROM Customers
    WHERE Email = <cfqueryparam value="#form.email#" cfsqltype="cf_sql_varchar">
</cfquery>

// ✅ BEST - Use queryExecute with params
<cfscript>
    qCustomer = queryExecute(
        "SELECT * FROM Customers WHERE Email = :email",
        { email = { value = form.email, cfsqltype = "cf_sql_varchar" } },
        { datasource = "myAppDB" }
    );
</cfscript>
```

### XSS Prevention

```coldfusion
<!--- ❌ Avoid - Unescaped output --->
<cfoutput>#form.userName#</cfoutput>

<!--- ✅ Good - HTML encoding --->
<cfoutput>#encodeForHTML(form.userName)#</cfoutput>

<!--- ✅ For attributes --->
<input type="text" value="<cfoutput>#encodeForHTMLAttribute(form.userName)#</cfoutput>">

<!--- ✅ For JavaScript --->
<script>
    var userName = '<cfoutput>#encodeForJavaScript(form.userName)#</cfoutput>';
</script>

<!--- ✅ For URLs --->
<a href="profile.cfm?name=<cfoutput>#encodeForURL(form.userName)#</cfoutput>">Profile</a>
```

### CSRF Protection

```coldfusion
<!--- Generate token (in Application.cfc or helper) --->
<cfscript>
    function csrfGenerateToken() {
        var token = hash(createUUID() & now(), "SHA-256");
        session.csrfToken = token;
        return token;
    }
    
    function csrfVerifyToken(required string token) {
        return (
            structKeyExists(session, "csrfToken") && 
            session.csrfToken == arguments.token
        );
    }
</cfscript>

<!--- In form --->
<form method="post">
    <input type="hidden" name="csrfToken" value="<cfoutput>#csrfGenerateToken()#</cfoutput>">
    <!--- Form fields --->
</form>

<!--- Verify on submission --->
<cfif structKeyExists(form, "submit")>
    <cfif !csrfVerifyToken(form.csrfToken)>
        <cfthrow type="SecurityException" message="Invalid CSRF token">
    </cfif>
    <!--- Process form --->
</cfif>
```

### Authentication

```coldfusion
<!--- Login function --->
<cfscript>
    function authenticateUser(required string email, required string password) {
        // Get user from database
        var sql = "
            SELECT UserId, Email, PasswordHash, Role
            FROM Users
            WHERE Email = :email AND IsActive = 1
        ";
        
        var qUser = queryExecute(
            sql,
            { email = { value = arguments.email, cfsqltype = "cf_sql_varchar" } },
            { datasource = "myAppDB" }
        );
        
        if (qUser.recordCount == 0) {
            return false;
        }
        
        // Verify password (using BCrypt)
        if (!bcryptCheckPassword(arguments.password, qUser.passwordHash)) {
            return false;
        }
        
        // Set session
        session.user = {
            id = qUser.userId,
            email = qUser.email,
            role = qUser.role,
            isAuthenticated = true
        };
        
        return true;
    }
    
    function isLoggedIn() {
        return (
            structKeyExists(session, "user") && 
            session.user.isAuthenticated
        );
    }
    
    function requireAuth() {
        if (!isLoggedIn()) {
            location(url="/login.cfm", addToken=false);
            abort;
        }
    }
</cfscript>
```

---

## ⚡ Performance Best Practices

### Query Caching

```coldfusion
<!--- Cache query results --->
<cfquery name="qStates" datasource="myAppDB" cachedWithin="#createTimeSpan(0, 1, 0, 0)#">
    SELECT StateCode, StateName
    FROM States
    ORDER BY StateName
</cfquery>

<!--- Or with queryExecute --->
<cfscript>
    qStates = queryExecute(
        "SELECT StateCode, StateName FROM States ORDER BY StateName",
        {},
        { 
            datasource = "myAppDB",
            cachedWithin = createTimeSpan(0, 1, 0, 0)  // Cache for 1 hour
        }
    );
</cfscript>
```

### Variable Scoping

```coldfusion
<!--- ❌ Avoid - Unscoped variables (slow lookup) --->
<cfloop from="1" to="100" index="i">
    <cfset total = total + i>  <!--- Where is total? --->
</cfloop>

<!--- ✅ Good - Scoped variables (fast lookup) --->
<cfset var total = 0>
<cfloop from="1" to="100" index="local.i">
    <cfset local.total = local.total + local.i>
</cfloop>

<!--- ✅ In CFScript --->
<cfscript>
    var total = 0;
    for (var i = 1; i <= 100; i++) {
        total += i;
    }
</cfscript>
```

### Output Buffering

```coldfusion
<!--- Use cfsavecontent for complex output --->
<cfsavecontent variable="emailBody">
    <h1>Welcome, <cfoutput>#customer.firstName#</cfoutput>!</h1>
    <p>Thank you for signing up.</p>
    <!--- More content --->
</cfsavecontent>

<cfmail to="#customer.email#" 
        from="noreply@example.com" 
        subject="Welcome"
        type="html">
    #emailBody#
</cfmail>
```

---

## 🧪 Testing

### Unit Testing with TestBox

```coldfusion
/**
 * Customer Service Tests
 */
component extends="testbox.system.BaseSpec" {
    
    function beforeAll() {
        // Setup
        variables.mockDAO = createMock("models.CustomerDAO");
        variables.mockLogger = createMock("components.Logger");
        variables.service = new components.services.CustomerService(
            customerDAO = variables.mockDAO,
            logger = variables.mockLogger
        );
    }
    
    function run() {
        describe("Customer Service", function() {
            
            it("should get customer by ID", function() {
                // Arrange
                var expectedCustomer = {
                    customerId = 1,
                    firstName = "John",
                    lastName = "Doe",
                    email = "john@example.com"
                };
                variables.mockDAO.$("getById").$args(1).$results(expectedCustomer);
                
                // Act
                var result = variables.service.getById(1);
                
                // Assert
                expect(result).toBeStruct();
                expect(result.customerId).toBe(1);
                expect(result.firstName).toBe("John");
            });
            
            it("should throw error for invalid ID", function() {
                // Act & Assert
                expect(function() {
                    variables.service.getById(-1);
                }).toThrow("ValidationException");
            });
            
        });
    }
}
```

---

## 📚 Best Practices Summary

### Do's ✅

- Use **CFScript** for business logic (more readable)
- Use **parameterized queries** (cfqueryparam/queryExecute)
- **Encode output** (encodeForHTML, encodeForHTMLAttribute)
- Implement **CSRF protection** for forms
- Use **component-based architecture** (CFCs)
- **Scope all variables** (var, local, variables)
- Implement proper **error handling** (try/catch)
- **Log errors** (writeLog)
- Use **BCrypt** for password hashing
- Cache expensive queries when appropriate

### Don'ts ❌

- Don't use string concatenation in SQL queries
- Don't output unescaped user input
- Don't use unscoped variables in loops
- Don't store passwords in plain text
- Don't ignore error handling
- Don't use deprecated tags (cfobject, cfwddx)
- Don't enable debugging in production
- Don't use application/session scope without locking (CF9 and earlier)
- Don't forget to validate input
- Don't use cfinclude for business logic (use CFCs)

---

## 🚀 Modernization Path

### Consider Migration Options

- **Lucee** - Open-source CFML engine with better performance
- **CommandBox** - Modern development tooling for CFML
- **ColdBox MVC** - Full-featured MVC framework
- **Node.js/Java** - For greenfield rewrites

### Document Legacy Patterns

```coldfusion
<!---
    TODO: [MODERNIZATION] This uses legacy cfform tags
    Consider refactoring to standard HTML forms with validation
--->
<cfform action="submit.cfm" method="post">
    <!--- Legacy cfform code --->
</cfform>
```

---

## 📖 Resources

- [Adobe ColdFusion Documentation](https://helpx.adobe.com/coldfusion/home.html)
- [Lucee Documentation](https://docs.lucee.org/)
- [CFML Slack Community](https://cfml-slack.herokuapp.com/)
- [Modern ColdFusion (CFML) in 100 Minutes](https://github.com/mhenke/Modern-ColdFusion-CFML-In-100-Minutes)

---

## 🏁 Maintenance Philosophy

**Remember:** ColdFusion is mature and stable. Focus on:

1. **Security** - Parameterized queries, output encoding, CSRF protection
2. **Maintainability** - Use CFCs, service layers, proper structure
3. **Documentation** - Help future maintainers understand business logic
4. **Performance** - Cache wisely, scope variables, optimize queries
5. **Modernization** - Consider Lucee, CommandBox, and modern CFML practices

ColdFusion applications can run reliably for decades with proper maintenance.
Keep it simple, secure, and well-documented.
