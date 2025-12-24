# AngularJS (1.x) Coding Standards - Legacy Maintenance

## 🎯 Overview

This guide covers **AngularJS 1.x** best practices for maintaining legacy
applications. AngularJS reached end-of-life in December 2021, but many
applications continue in maintenance mode. This guide focuses on sustainable
maintenance patterns, not new feature development.

**Key Principles:**

- Follow John Papa's AngularJS Style Guide patterns
- Maintain consistency with existing codebase
- Avoid introducing technical debt
- Document complex patterns for future maintainers
- Consider migration path to modern frameworks when appropriate

---

## 📁 Project Structure

### Single Responsibility Principle

**One component per file** - Makes code easier to maintain and test.

```javascript
// ❌ Avoid - Multiple components in one file
angular.module("app").controller("CustomerCtrl", function () {});
angular.module("app").factory("CustomerService", function () {});

// ✅ Good - One component per file
// customer.controller.js
angular.module("app").controller("CustomerCtrl", function () {});

// customer.service.js
angular.module("app").factory("CustomerService", function () {});
```

### Folder Structure

Organize by **feature**, not by type:

```
app/
├── customers/
│   ├── customer-list.component.js
│   ├── customer-detail.component.js
│   ├── customer.service.js
│   ├── customer.routes.js
│   └── customers.module.js
├── orders/
│   ├── order-list.component.js
│   ├── order.service.js
│   └── orders.module.js
└── shared/
    ├── logger.service.js
    ├── http-interceptor.service.js
    └── shared.module.js
```

---

## 🏗️ Module Definition

### Module Naming

Use **dot notation** for nested modules:

```javascript
// Main module
angular.module("app", ["ngRoute", "ngAnimate"]);

// Feature modules
angular.module("app.customers", []);
angular.module("app.orders", []);
angular.module("app.shared", []);
```

### Module Getters vs Setters

**Setter** (with dependencies) - Use once to create:

```javascript
// customers.module.js
angular.module("app.customers", ["app.shared"]);
```

**Getter** (no dependencies) - Use everywhere else:

```javascript
// customer-list.component.js
angular.module("app.customers")
  .component("customerList", {
    // ...
  });
```

### Avoid Anonymous Functions

Use **named functions** for better stack traces:

```javascript
// ❌ Avoid
angular.module("app")
  .controller("CustomerCtrl", function ($http) {
    // ...
  });

// ✅ Good
angular.module("app")
  .controller("CustomerCtrl", CustomerController);

function CustomerController($http) {
  // ...
}
```

---

## 🎮 Controllers

### Controller As Syntax

Always use `controllerAs` syntax:

```javascript
// ❌ Avoid - $scope usage
function CustomerController($scope) {
  $scope.name = "John";
  $scope.save = function () {/* ... */};
}

// ✅ Good - controllerAs
function CustomerController() {
  var vm = this;
  vm.name = "John";
  vm.save = save;

  function save() {
    // Implementation
  }
}
```

### View Model Convention

Use `vm` (view model) for the `this` reference:

```javascript
function CustomerController() {
  var vm = this; // ✅ Consistent convention

  vm.title = "Customer Details";
  vm.customer = {};
  vm.save = save;
  vm.cancel = cancel;

  activate();

  function activate() {
    // Controller initialization
  }

  function save() {
    // Save logic
  }

  function cancel() {
    // Cancel logic
  }
}
```

### Bindable Members at Top

Put bindable properties and functions at the **top** of the controller:

```javascript
function CustomerController(customerService) {
  var vm = this;

  // Bindable members
  vm.customers = [];
  vm.selectedCustomer = null;
  vm.search = search;
  vm.refresh = refresh;

  activate();

  // Implementation details below
  function activate() {
    return loadCustomers();
  }

  function search(query) {
    // ...
  }

  function refresh() {
    return loadCustomers();
  }

  function loadCustomers() {
    return customerService.getAll()
      .then(function (data) {
        vm.customers = data;
      });
  }
}
```

### Defer Controller Logic to Services

Keep controllers **thin** - move business logic to services:

```javascript
// ❌ Avoid - Business logic in controller
function CustomerController($http) {
  var vm = this;
  vm.save = save;

  function save(customer) {
    return $http.post("/api/customers", customer)
      .then(function (response) {
        // Complex validation and transformation
        if (response.data.status === "success") {
          customer.id = response.data.id;
          customer.updatedAt = new Date();
        }
        return customer;
      });
  }
}

// ✅ Good - Business logic in service
function CustomerController(customerService) {
  var vm = this;
  vm.save = save;

  function save(customer) {
    return customerService.save(customer);
  }
}
```

---

## 🏭 Services & Factories

### Prefer Factories Over Services

Use **factories** for consistency (they're more flexible):

```javascript
// ✅ Factory pattern (recommended)
angular.module("app")
  .factory("customerService", customerService);

function customerService($http, logger) {
  var service = {
    getAll: getAll,
    getById: getById,
    save: save,
    remove: remove,
  };

  return service;

  function getAll() {
    return $http.get("/api/customers")
      .then(getComplete)
      .catch(getFailed);
  }

  function getById(id) {
    return $http.get("/api/customers/" + id)
      .then(getComplete)
      .catch(getFailed);
  }

  function save(customer) {
    if (customer.id) {
      return $http.put("/api/customers/" + customer.id, customer);
    }
    return $http.post("/api/customers", customer);
  }

  function remove(id) {
    return $http.delete("/api/customers/" + id);
  }

  function getComplete(response) {
    return response.data;
  }

  function getFailed(error) {
    logger.error("XHR Failed: " + error.data);
    return $q.reject(error);
  }
}
```

### Singleton Services

Services are **singletons** - use them for shared state:

```javascript
angular.module("app")
  .factory("authService", authService);

function authService($http, $window) {
  var service = {
    currentUser: null,
    login: login,
    logout: logout,
    isAuthenticated: isAuthenticated,
  };

  return service;

  function login(credentials) {
    return $http.post("/api/auth/login", credentials)
      .then(function (response) {
        service.currentUser = response.data.user;
        $window.localStorage.setItem("token", response.data.token);
        return service.currentUser;
      });
  }

  function logout() {
    service.currentUser = null;
    $window.localStorage.removeItem("token");
  }

  function isAuthenticated() {
    return !!service.currentUser;
  }
}
```

---

## 🎨 Directives

### Restrict to Elements and Attributes

Avoid `class` and `comment` restrictions:

```javascript
// ✅ Good - Element directive
angular.module("app")
  .directive("customerCard", customerCard);

function customerCard() {
  return {
    restrict: "E", // <customer-card>
    templateUrl: "customers/customer-card.html",
    scope: {
      customer: "=",
    },
    controller: CustomerCardController,
    controllerAs: "vm",
    bindToController: true,
  };
}

// ✅ Good - Attribute directive
angular.module("app")
  .directive("validEmail", validEmail);

function validEmail() {
  return {
    restrict: "A", // <input valid-email>
    require: "ngModel",
    link: function (scope, element, attrs, ngModel) {
      ngModel.$validators.validEmail = function (modelValue) {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(modelValue);
      };
    },
  };
}
```

### Component Directive Pattern (Pre-1.5)

For AngularJS < 1.5, use this component pattern:

```javascript
angular.module("app")
  .directive("customerList", customerList);

function customerList() {
  return {
    restrict: "E",
    templateUrl: "customers/customer-list.html",
    scope: {},
    controller: CustomerListController,
    controllerAs: "vm",
    bindToController: true,
  };
}

function CustomerListController(customerService) {
  var vm = this;
  vm.customers = [];
  vm.selectCustomer = selectCustomer;

  activate();

  function activate() {
    return customerService.getAll()
      .then(function (data) {
        vm.customers = data;
      });
  }

  function selectCustomer(customer) {
    // Handle selection
  }
}
```

### Use .component() for AngularJS 1.5+

If using AngularJS 1.5 or higher, prefer **components**:

```javascript
angular.module("app")
  .component("customerDetail", {
    templateUrl: "customers/customer-detail.html",
    bindings: {
      customer: "<",
      onSave: "&",
      onCancel: "&",
    },
    controller: CustomerDetailController,
  });

function CustomerDetailController() {
  var vm = this;

  vm.$onInit = function () {
    // Component initialization
    vm.editableCustomer = angular.copy(vm.customer);
  };

  vm.save = function () {
    vm.onSave({ customer: vm.editableCustomer });
  };

  vm.cancel = function () {
    vm.onCancel();
  };
}
```

---

## 🔄 Routing

### UI-Router Recommended

Prefer **ui-router** over ngRoute for complex applications:

```javascript
angular.module("app")
  .config(routerConfig);

function routerConfig($stateProvider, $urlRouterProvider) {
  $urlRouterProvider.otherwise("/customers");

  $stateProvider
    .state("customers", {
      url: "/customers",
      templateUrl: "customers/customer-list.html",
      controller: "CustomerListCtrl",
      controllerAs: "vm",
      resolve: {
        customers: function (customerService) {
          return customerService.getAll();
        },
      },
    })
    .state("customer-detail", {
      url: "/customers/:id",
      templateUrl: "customers/customer-detail.html",
      controller: "CustomerDetailCtrl",
      controllerAs: "vm",
      resolve: {
        customer: function ($stateParams, customerService) {
          return customerService.getById($stateParams.id);
        },
      },
    });
}
```

### Route Resolve

Use **resolve** to fetch data before route activates:

```javascript
.state('order-detail', {
  url: '/orders/:id',
  templateUrl: 'orders/order-detail.html',
  controller: 'OrderDetailCtrl',
  controllerAs: 'vm',
  resolve: {
    order: function($stateParams, orderService) {
      return orderService.getById($stateParams.id);
    },
    customer: function(order, customerService) {
      return customerService.getById(order.customerId);
    }
  }
})
```

Inject resolved data in controller:

```javascript
function OrderDetailController(order, customer) {
  var vm = this;
  vm.order = order;
  vm.customer = customer;
}
```

---

## 🎁 Promises

### Use $q for Promises

Angular's **$q** integrates with digest cycle:

```javascript
angular.module("app")
  .factory("dataService", dataService);

function dataService($q, $timeout) {
  var service = {
    getData: getData,
    processData: processData,
  };

  return service;

  function getData() {
    var deferred = $q.defer();

    $timeout(function () {
      deferred.resolve({ id: 1, name: "Data" });
    }, 1000);

    return deferred.promise;
  }

  function processData() {
    return getData()
      .then(function (data) {
        return transformData(data);
      })
      .catch(function (error) {
        return $q.reject(error);
      });
  }

  function transformData(data) {
    return {
      id: data.id,
      name: data.name.toUpperCase(),
      timestamp: new Date(),
    };
  }
}
```

### Chain Promises

```javascript
function CustomerController(customerService, orderService) {
  var vm = this;
  vm.customer = null;
  vm.orders = [];

  activate();

  function activate() {
    var customerId = 1;

    return customerService.getById(customerId)
      .then(function (customer) {
        vm.customer = customer;
        return orderService.getByCustomerId(customer.id);
      })
      .then(function (orders) {
        vm.orders = orders;
      })
      .catch(function (error) {
        logger.error("Failed to load customer data");
      });
  }
}
```

---

## 🔧 Dependency Injection

### Manual Annotation for Minification

Use **$inject** array to make code minification-safe:

```javascript
// ✅ Good - Minification safe
angular.module("app")
  .controller("CustomerCtrl", CustomerController);

CustomerController.$inject = ["$http", "$q", "customerService", "logger"];

function CustomerController($http, $q, customerService, logger) {
  var vm = this;
  // ...
}
```

Or use array syntax:

```javascript
angular.module("app")
  .controller("CustomerCtrl", [
    "$http",
    "$q",
    "customerService",
    "logger",
    CustomerController,
  ]);

function CustomerController($http, $q, customerService, logger) {
  var vm = this;
  // ...
}
```

**Never** rely on implicit annotation (breaks with minification):

```javascript
// ❌ Avoid - Breaks when minified
function CustomerController($http, customerService) {
  // ...
}
```

---

## 🧪 Testing

### Unit Testing Controllers

Use **Jasmine** or **Mocha** with **karma**:

```javascript
describe("CustomerController", function () {
  var controller;
  var customerService;
  var $q;
  var $rootScope;

  beforeEach(module("app"));

  beforeEach(inject(function ($controller, _$q_, _$rootScope_) {
    $q = _$q_;
    $rootScope = _$rootScope_;

    // Mock service
    customerService = {
      getAll: jasmine.createSpy("getAll").and.returnValue(
        $q.resolve([{ id: 1, name: "Customer 1" }]),
      ),
    };

    controller = $controller("CustomerController", {
      customerService: customerService,
    });
  }));

  it("should load customers on activation", function () {
    $rootScope.$apply(); // Trigger digest cycle

    expect(customerService.getAll).toHaveBeenCalled();
    expect(controller.customers.length).toBe(1);
    expect(controller.customers[0].name).toBe("Customer 1");
  });
});
```

### Unit Testing Services

```javascript
describe("customerService", function () {
  var service;
  var $httpBackend;

  beforeEach(module("app"));

  beforeEach(inject(function (_customerService_, _$httpBackend_) {
    service = _customerService_;
    $httpBackend = _$httpBackend_;
  }));

  afterEach(function () {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
  });

  it("should get all customers", function () {
    var mockData = [{ id: 1, name: "Customer 1" }];

    $httpBackend.expectGET("/api/customers")
      .respond(200, mockData);

    service.getAll().then(function (data) {
      expect(data).toEqual(mockData);
    });

    $httpBackend.flush();
  });
});
```

---

## ⚡ Performance

### One-Time Binding

Use **`::`** for data that doesn't change:

```html
<!-- ❌ Avoid - Creates watcher -->
<div>{{vm.customer.name}}</div>

<!-- ✅ Good - One-time binding, no watcher -->
<div>{{::vm.customer.name}}</div>
```

### $digest Optimization

Minimize watchers:

```javascript
// ❌ Avoid - Creates unnecessary watchers
<div ng-repeat="customer in vm.customers">
  <span>{{customer.name}}</span>
  <span>{{customer.email}}</span>
  <span>{{customer.phone}}</span>
</div>

// ✅ Good - One-time binding where possible
<div ng-repeat="customer in vm.customers track by customer.id">
  <span>{{::customer.name}}</span>
  <span>{{::customer.email}}</span>
  <span>{{customer.status}}</span> <!-- Only this watches -->
</div>
```

### Track By in ng-repeat

Always use **track by** for better performance:

```html
<!-- ❌ Avoid -->
<div ng-repeat="item in vm.items">

<!-- ✅ Good -->
<div ng-repeat="item in vm.items track by item.id">
```

### Debounce Inputs

For search inputs, use **ng-model-options**:

```html
<input
  type="text"
  ng-model="vm.searchQuery"
  ng-model-options="{ debounce: 300 }"
  ng-change="vm.search()"
>
```

---

## 🛡️ Security

### Sanitize User Input

Always sanitize HTML content:

```javascript
// Use $sanitize service
angular.module("app", ["ngSanitize"]);

// In controller
function CommentController($sanitize) {
  var vm = this;

  vm.saveComment = function (comment) {
    vm.comment = $sanitize(comment);
  };
}
```

### Use ng-bind-html Carefully

```html
<!-- ✅ Good - Sanitized -->
<div ng-bind-html="vm.trustedHtml"></div>

<!-- In controller -->
<script>
  function ContentController($sce) {
    var vm = this;

    // Only trust content you control
    vm.trustedHtml = $sce.trustAsHtml("<p>Safe content</p>");
  }
</script>
```

### CSRF Protection

Include CSRF tokens in HTTP requests:

```javascript
angular.module("app")
  .config(function ($httpProvider) {
    $httpProvider.defaults.xsrfCookieName = "XSRF-TOKEN";
    $httpProvider.defaults.xsrfHeaderName = "X-XSRF-TOKEN";
  });
```

---

## 📚 Best Practices Summary

### Do's ✅

- Use `controllerAs` syntax consistently
- Keep controllers thin, business logic in services
- Use named functions for better debugging
- Implement proper error handling
- Use one-time binding where possible
- Always use `track by` in `ng-repeat`
- Manually annotate dependencies for minification
- Write unit tests for controllers and services
- Follow single responsibility principle
- Organize by feature, not by type

### Don'ts ❌

- Don't use `$scope` unless necessary ($watch, events)
- Don't put business logic in controllers
- Don't use anonymous functions
- Don't forget dependency injection annotations
- Don't manipulate DOM in controllers
- Don't use global variables
- Don't ignore error handling
- Don't create unnecessary watchers
- Don't use `ng-repeat` without `track by`

---

## 🚀 Migration Considerations

Since AngularJS is EOL, consider:

### Document Technical Debt

```javascript
// TODO: [MIGRATION] This controller uses $scope events
// Consider refactoring to service-based communication
// when migrating to Angular/React
function LegacyController($scope) {
  $scope.$on("someEvent", function (event, data) {
    // ...
  });
}
```

### Incremental Modernization

- **Components over directives** - If on 1.5+, use .component()
- **Services over factories** - Prepare for ES6 classes
- **TypeScript** - Consider adding types incrementally
- **Build tools** - Use Webpack/Babel for modern JS features

### Hybrid Applications

Consider using **ngUpgrade** for gradual Angular migration:

```javascript
// Bootstrap hybrid app
import { UpgradeModule } from '@angular/upgrade/static';

@NgModule({
  imports: [
    BrowserModule,
    UpgradeModule
  ]
})
export class AppModule {
  constructor(private upgrade: UpgradeModule) {}
  
  ngDoBootstrap() {
    this.upgrade.bootstrap(document.body, ['legacyApp'], { strictDi: true });
  }
}
```

---

## 📖 Resources

- [John Papa's AngularJS Style Guide](https://github.com/johnpapa/angular-styleguide/blob/master/a1/README.md)
- [AngularJS Official Documentation](https://docs.angularjs.org)
- [Todd Motto's AngularJS Style Guide](https://github.com/toddmotto/angularjs-styleguide)
- [AngularJS Patterns (ng-book)](https://www.ng-book.com/p/The-Complete-Book-on-AngularJS/)

---

## 🏁 Maintenance Philosophy

**Remember:** This is a **legacy** framework. Focus on:

1. **Stability** - Don't introduce breaking changes
2. **Documentation** - Help future maintainers understand the code
3. **Consistency** - Follow existing patterns in the codebase
4. **Security** - Keep dependencies patched where possible
5. **Exit Strategy** - Plan for eventual migration to modern frameworks

When in doubt, prefer **boring, proven patterns** over clever solutions. Future
maintainers (including yourself) will thank you.
