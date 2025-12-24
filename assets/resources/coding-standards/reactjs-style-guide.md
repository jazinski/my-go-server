# ReactJS Coding Standards

## 🎯 Overview

This document defines our team's ReactJS coding standards. All React code must
follow these guidelines to ensure consistency, maintainability, and code quality
across projects.

## 📦 Project Structure

```
my-react-app/
├── public/
│   ├── index.html
│   └── favicon.ico
├── src/
│   ├── components/           # Reusable components
│   │   ├── common/          # Shared UI components (Button, Input, etc.)
│   │   ├── layout/          # Layout components (Header, Footer, Sidebar)
│   │   └── features/        # Feature-specific components
│   ├── pages/               # Page/route components
│   │   ├── Home/
│   │   ├── UserProfile/
│   │   └── Dashboard/
│   ├── hooks/               # Custom React hooks
│   ├── contexts/            # React context providers
│   ├── services/            # API calls and external services
│   ├── utils/               # Utility functions
│   ├── types/               # TypeScript types/interfaces
│   ├── constants/           # App constants
│   ├── assets/              # Images, fonts, etc.
│   ├── styles/              # Global styles
│   ├── App.tsx
│   └── index.tsx
├── package.json
├── tsconfig.json
└── README.md
```

## 🎨 Component Structure & Naming

### File Naming

- **Components**: `PascalCase.tsx` - `UserProfile.tsx`, `NavigationBar.tsx`
- **Hooks**: `camelCase.ts` - `useAuth.ts`, `useFetch.ts`
- **Utils**: `camelCase.ts` - `formatDate.ts`, `validateEmail.ts`
- **Constants**: `UPPER_SNAKE_CASE.ts` - `API_ENDPOINTS.ts`, `ROUTE_PATHS.ts`

### Component Types

**Functional Components (Preferred):**

```tsx
// ✅ DO: Use functional components with TypeScript
import React from "react";

interface UserProfileProps {
  userId: string;
  onEdit?: (userId: string) => void;
}

export const UserProfile: React.FC<UserProfileProps> = ({ userId, onEdit }) => {
  const [user, setUser] = React.useState<User | null>(null);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await userService.getUser(userId);
        setUser(userData);
      } catch (error) {
        console.error("Failed to fetch user:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [userId]);

  if (loading) return <LoadingSpinner />;
  if (!user) return <ErrorMessage message="User not found" />;

  return (
    <div className="user-profile">
      <h1>{user.name}</h1>
      <p>{user.email}</p>
      {onEdit && <button onClick={() => onEdit(userId)}>Edit Profile</button>}
    </div>
  );
};
```

### Component Organization

```tsx
// ✅ DO: Organize component code in this order:
import React, { useEffect, useState } from "react"; // 1. React imports
import { useNavigate } from "react-router-dom"; // 2. Third-party imports
import { Button } from "@/components/common"; // 3. Internal imports
import { formatDate } from "@/utils/dateUtils"; // 4. Utilities
import "./UserProfile.css"; // 5. Styles

// 6. TypeScript interfaces/types
interface UserProfileProps {
  userId: string;
}

interface User {
  id: string;
  name: string;
  email: string;
}

// 7. Constants (component-specific)
const MAX_NAME_LENGTH = 50;

// 8. Component definition
export const UserProfile: React.FC<UserProfileProps> = ({ userId }) => {
  // 9. Hooks (useState, useEffect, custom hooks)
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

  // 10. Event handlers
  const handleEdit = () => {
    navigate(`/users/${userId}/edit`);
  };

  // 11. Render logic
  return (
    <div>
      {/* Component JSX */}
    </div>
  );
};
```

## 🪝 React Hooks Best Practices

### useState

```tsx
// ✅ DO: Use TypeScript with useState
const [count, setCount] = useState<number>(0);
const [user, setUser] = useState<User | null>(null);
const [items, setItems] = useState<Item[]>([]);

// ✅ DO: Use functional updates when new state depends on previous
setCount((prevCount) => prevCount + 1);

// ❌ DON'T: Mutate state directly
items.push(newItem); // ❌
setItems(items);

// ✅ DO: Create new arrays/objects
setItems([...items, newItem]); // ✅
setUser({ ...user, name: "New Name" }); // ✅
```

### useEffect

```tsx
// ✅ DO: Specify dependencies correctly
useEffect(() => {
  fetchUser(userId);
}, [userId]); // Re-run when userId changes

// ✅ DO: Clean up side effects
useEffect(() => {
  const subscription = dataService.subscribe(handleData);

  return () => {
    subscription.unsubscribe(); // Cleanup
  };
}, []);

// ✅ DO: Handle async operations properly
useEffect(() => {
  let cancelled = false;

  const fetchData = async () => {
    try {
      const data = await api.getData();
      if (!cancelled) {
        setData(data);
      }
    } catch (error) {
      if (!cancelled) {
        setError(error);
      }
    }
  };

  fetchData();

  return () => {
    cancelled = true; // Prevent state updates if unmounted
  };
}, []);

// ❌ DON'T: Omit dependencies (unless intentional)
useEffect(() => {
  console.log(user.name); // ❌ 'user' should be in dependency array
}, []);
```

### Custom Hooks

```tsx
// ✅ DO: Extract reusable logic into custom hooks
export const useFetch = <T,>(url: string) => {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response = await fetch(url);
        if (!response.ok) throw new Error("Fetch failed");
        const json = await response.json();
        setData(json);
      } catch (err) {
        setError(err as Error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [url]);

  return { data, loading, error };
};

// Usage:
const { data: user, loading, error } = useFetch<User>("/api/users/123");
```

### useCallback & useMemo

```tsx
// ✅ DO: Use useCallback for functions passed to child components
const handleClick = useCallback(() => {
  doSomething(userId);
}, [userId]); // Only recreate if userId changes

// ✅ DO: Use useMemo for expensive calculations
const expensiveValue = useMemo(() => {
  return calculateExpensiveValue(data);
}, [data]);

// ❌ DON'T: Overuse - only when there's a real performance issue
const simpleValue = useMemo(() => x + y, [x, y]); // ❌ Overkill
const simpleValue = x + y; // ✅ Just calculate it
```

## 🎭 Props & TypeScript

### Defining Props

```tsx
// ✅ DO: Define props interface
interface ButtonProps {
  label: string;
  onClick: () => void;
  disabled?: boolean;
  variant?: "primary" | "secondary" | "danger";
  className?: string;
  children?: React.ReactNode;
}

// ✅ DO: Use default props
export const Button: React.FC<ButtonProps> = ({
  label,
  onClick,
  disabled = false,
  variant = "primary",
  className = "",
  children,
}) => {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`btn btn-${variant} ${className}`}
    >
      {children || label}
    </button>
  );
};

// ✅ DO: Destructure props in function signature
// ❌ DON'T:
export const Button: React.FC<ButtonProps> = (props) => {
  return <button onClick={props.onClick}>{props.label}</button>;
};
```

### Props Validation

```tsx
// ✅ DO: Use TypeScript for type safety
interface UserCardProps {
  user: {
    id: string;
    name: string;
    email: string;
    role: "admin" | "user" | "guest";
  };
  onEdit?: (userId: string) => void;
  showDetails?: boolean;
}

// ❌ DON'T: Use any
interface BadProps {
  data: any; // ❌ Loses type safety
}

// ✅ DO: Use generics if type is truly dynamic
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
}
```

## 🎨 JSX Best Practices

### Conditional Rendering

```tsx
// ✅ DO: Use ternary for simple conditions
{
  isLoading ? <Spinner /> : <Content />;
}

// ✅ DO: Use && for single condition
{
  error && <ErrorMessage error={error} />;
}
{
  user?.isPremium && <PremiumBadge />;
}

// ✅ DO: Extract complex conditions
const showWarning = user && !user.isVerified && user.createdAt < thirtyDaysAgo;
{
  showWarning && <WarningMessage />;
}

// ❌ DON'T: Use long ternary chains
{
  loading ? <Spinner /> : error ? <Error /> : data ? <Content /> : null;
} // ❌

// ✅ DO: Use early returns for complex conditions
if (loading) return <Spinner />;
if (error) return <ErrorMessage error={error} />;
if (!data) return <NoDataMessage />;
return <Content data={data} />;
```

### Lists & Keys

```tsx
// ✅ DO: Use unique, stable IDs for keys
{
  users.map((user) => <UserCard key={user.id} user={user} />);
}

// ❌ DON'T: Use array index as key (unless list is static)
{
  users.map((user, index) => (
    <UserCard key={index} user={user} /> // ❌ Can cause issues
  ));
}

// ✅ DO: Extract list items into components for clarity
const UserList: React.FC<{ users: User[] }> = ({ users }) => (
  <ul>
    {users.map((user) => <UserListItem key={user.id} user={user} />)}
  </ul>
);
```

### Event Handlers

```tsx
// ✅ DO: Use arrow functions in component body
const handleClick = () => {
  doSomething();
};

// ✅ DO: Pass parameters correctly
<button onClick={() => handleDelete(user.id)}>Delete</button>;

// ✅ DO: Prevent default when needed
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  submitForm();
};

// ❌ DON'T: Create inline functions for complex logic
<button
  onClick={() => {
    // 20 lines of code here ❌
  }}
>
  Submit
</button>;
```

## 🏪 State Management

### Context API

```tsx
// ✅ DO: Create context with TypeScript
interface AuthContextType {
  user: User | null;
  login: (credentials: Credentials) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ✅ DO: Create custom hook for context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
};

// ✅ DO: Create provider component
export const AuthProvider: React.FC<{ children: React.ReactNode }> = (
  { children },
) => {
  const [user, setUser] = useState<User | null>(null);

  const login = async (credentials: Credentials) => {
    const userData = await authService.login(credentials);
    setUser(userData);
  };

  const logout = () => {
    authService.logout();
    setUser(null);
  };

  const value = {
    user,
    login,
    logout,
    isAuthenticated: !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
```

### Local vs Global State

```tsx
// ✅ DO: Use local state for component-specific data
const [isOpen, setIsOpen] = useState(false); // Only this component needs it

// ✅ DO: Use context/global state for shared data
const { user, theme } = useAuth(); // Multiple components need it

// ❌ DON'T: Put everything in global state
// Only use global state when data is needed by multiple components
```

## 🔄 Side Effects & API Calls

### API Service Pattern

```tsx
// ✅ DO: Create service modules for API calls
// services/userService.ts
export const userService = {
  async getUser(id: string): Promise<User> {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) {
      throw new Error("Failed to fetch user");
    }
    return response.json();
  },

  async updateUser(id: string, data: Partial<User>): Promise<User> {
    const response = await fetch(`/api/users/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error("Failed to update user");
    }
    return response.json();
  },
};

// Component usage:
const UserProfile: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await userService.getUser("123");
        setUser(userData);
      } catch (error) {
        console.error("Failed to fetch user:", error);
      }
    };

    fetchUser();
  }, []);

  return <div>{user?.name}</div>;
};
```

## ♿ Accessibility (A11y)

```tsx
// ✅ DO: Use semantic HTML
<button onClick={handleClick}>Submit</button> // ✅
<div onClick={handleClick}>Submit</div> // ❌

// ✅ DO: Add ARIA labels where needed
<button aria-label="Close dialog" onClick={onClose}>
  <CloseIcon />
</button>

// ✅ DO: Use form labels
<label htmlFor="email">Email:</label>
<input id="email" type="email" name="email" />

// ✅ DO: Provide alt text for images
<img src={user.avatar} alt={`${user.name}'s profile picture`} />

// ✅ DO: Handle keyboard navigation
const handleKeyDown = (e: React.KeyboardEvent) => {
  if (e.key === 'Enter' || e.key === ' ') {
    handleClick();
  }
};

<div
  role="button"
  tabIndex={0}
  onClick={handleClick}
  onKeyDown={handleKeyDown}
>
  Custom Button
</div>
```

## 🧪 Testing

### Component Testing (React Testing Library)

```tsx
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { UserProfile } from "./UserProfile";

describe("UserProfile", () => {
  it("renders user information", async () => {
    // Arrange
    const mockUser = { id: "1", name: "John Doe", email: "john@example.com" };
    jest.spyOn(userService, "getUser").mockResolvedValue(mockUser);

    // Act
    render(<UserProfile userId="1" />);

    // Assert
    await waitFor(() => {
      expect(screen.getByText("John Doe")).toBeInTheDocument();
      expect(screen.getByText("john@example.com")).toBeInTheDocument();
    });
  });

  it("calls onEdit when edit button is clicked", () => {
    // Arrange
    const mockOnEdit = jest.fn();
    const mockUser = { id: "1", name: "John Doe", email: "john@example.com" };

    // Act
    render(<UserProfile user={mockUser} onEdit={mockOnEdit} />);
    fireEvent.click(screen.getByRole("button", { name: /edit/i }));

    // Assert
    expect(mockOnEdit).toHaveBeenCalledWith("1");
  });
});
```

## ⚡ Performance Optimization

```tsx
// ✅ DO: Code splitting with lazy loading
import { lazy, Suspense } from "react";

const HeavyComponent = lazy(() => import("./HeavyComponent"));

const App = () => (
  <Suspense fallback={<LoadingSpinner />}>
    <HeavyComponent />
  </Suspense>
);

// ✅ DO: Memoize expensive components
const ExpensiveComponent = React.memo<ExpensiveProps>(({ data }) => {
  return <div>{/* Complex rendering */}</div>;
});

// ✅ DO: Use React.memo selectively for components that re-render often
// ❌ DON'T: Wrap every component in React.memo
```

## 🚫 Common Anti-Patterns to Avoid

### ❌ DON'T:

```tsx
// Mutate state directly
state.items.push(newItem); // ❌
setState(state);

// Use index as key with dynamic lists
{
  items.map((item, index) => <Item key={index} />);
} // ❌

// Call hooks conditionally
if (condition) {
  useState(0); // ❌ Breaks rules of hooks
}

// Forget dependencies in useEffect
useEffect(() => {
  fetchData(userId); // ❌ userId should be in deps
}, []);

// Create components inside components
const Parent = () => {
  const Child = () => <div>Child</div>; // ❌ Recreated every render
  return <Child />;
};
```

### ✅ DO:

```tsx
// Create new state objects
setState({ ...state, items: [...state.items, newItem] }); // ✅

// Use stable unique IDs
{
  items.map((item) => <Item key={item.id} />);
} // ✅

// Always call hooks at top level
const [count, setCount] = useState(0); // ✅
if (condition) {
  // Use hook result conditionally
}

// Include all dependencies
useEffect(() => {
  fetchData(userId);
}, [userId]); // ✅

// Define components outside
const Child = () => <div>Child</div>; // ✅
const Parent = () => <Child />;
```

## 🛠️ Tools & Configuration

### ESLint Configuration

```json
{
  "extends": [
    "react-app",
    "plugin:jsx-a11y/recommended"
  ],
  "plugins": ["jsx-a11y"],
  "rules": {
    "react-hooks/rules-of-hooks": "error",
    "react-hooks/exhaustive-deps": "warn",
    "jsx-a11y/anchor-is-valid": "warn"
  }
}
```

### TypeScript Configuration

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "jsx": "react-jsx",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  }
}
```

## 📚 Additional Resources

- [React Documentation](https://react.dev/)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [React Testing Library](https://testing-library.com/react)
- [React Accessibility](https://react.dev/learn/accessibility)
- [React Performance](https://react.dev/learn/render-and-commit)

---

**Last Updated:** 2025-12-23\
**Version:** 1.0.0
