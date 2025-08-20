# DashOPS Frontend

A modern, TypeScript-based React application for cloud operations management, built with Vite and featuring a comprehensive UI component system.

## 🏗️ Architecture Overview

### **Tech Stack**

- **React 18.3** - Modern React with hooks and concurrent features
- **TypeScript** - Strict type checking and enhanced developer experience
- **Vite** - Lightning-fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **shadcn/ui** - High-quality, accessible React components
- **React Router 7** - Client-side routing
- **Axios** - HTTP client with interceptors
- **Sonner** - Toast notifications

### **Key Features**

- 🔐 **OAuth2 Authentication** with GitHub integration
- ☁️ **AWS Management** - EC2 instances, accounts, permissions
- ⚙️ **Kubernetes Operations** - Clusters, deployments, pods, logs
- 📊 **Dashboard** - Centralized monitoring and metrics
- 🎨 **Modern UI** - Responsive design with dark/light themes
- 🔄 **Real-time Updates** - Auto-refresh capabilities
- 🧩 **Plugin Architecture** - Modular and extensible

## 🚀 Getting Started

### Prerequisites

- **Node.js** >= 18.0.0
- **Yarn** >= 1.22.0 (preferred package manager)

### Installation & Development

```bash
# Clone and navigate to frontend
cd front

# Install dependencies
yarn

# Start development server
yarn dev

# Open browser at http://localhost:5173
```

### Build for Production

```bash
# Type check + build optimized bundle
yarn build

# Preview production build
yarn preview
```

## 📦 Available Scripts

### **Development**

```bash
yarn dev              # Start development server with hot reload
yarn preview          # Preview production build locally
```

### **Quality Assurance**

```bash
yarn quality          # Run complete quality gate (type + lint + format)
yarn fix              # Auto-fix linting and formatting issues
```

### **Type Checking**

```bash
yarn type-check       # Run TypeScript type checking
yarn type-check:watch # Run type checking in watch mode
```

### **Linting**

```bash
yarn lint             # Run ESLint
yarn lint:fix         # Run ESLint with auto-fix
yarn lint:check       # Run ESLint with zero warnings policy (CI/CD)
```

### **Code Formatting**

```bash
yarn format           # Format all files with Prettier
yarn format:check     # Check if files are properly formatted
```

### **Testing**

```bash
yarn test             # Run unit tests with Vitest
yarn coverage         # Run tests with coverage report
```

### **Build**

```bash
yarn build            # Type check + production build
```

## 🎯 Project Structure

### **Directory Organization**

```
src/
├── components/          # Shared React components
│   ├── ui/             # shadcn/ui components (Button, Card, etc.)
│   ├── AppSidebar.tsx  # Main navigation sidebar
│   ├── ContentWithMenu.tsx # Layout with navigation menu
│   └── ...
├── hooks/              # Custom React hooks
│   ├── use-interval.ts # Auto-refresh interval hook
│   └── use-mobile.ts   # Mobile device detection
├── helpers/            # Utility functions
│   ├── http.ts         # Axios configuration and interceptors
│   ├── oauth.ts        # OAuth2 authentication helpers
│   ├── localStorage.ts # Browser storage utilities
│   └── loadModules.ts  # Dynamic module loading
├── lib/
│   └── utils.ts        # General utility functions
├── modules/            # Feature modules
│   ├── aws/           # AWS operations
│   ├── kubernetes/    # Kubernetes management
│   ├── oauth2/        # Authentication & user management
│   ├── dashboard/     # Main dashboard
│   └── config/        # Application configuration
├── types/             # TypeScript type definitions
│   ├── index.ts       # Centralized exports
│   ├── common.ts      # Shared types (Menu, Router, Filter)
│   ├── api.ts         # API and state management types
│   └── ui.ts          # UI component types
└── main.tsx           # Application entry point
```

### **Module Architecture**

Each module follows a consistent structure:

```
modules/[module-name]/
├── index.tsx          # Module configuration and routing
├── types.ts           # Module-specific TypeScript types
├── [Name]Page.tsx     # Main page components
├── [name]Resource.ts  # API functions and data fetching
├── [Name]Actions.tsx  # Action components (buttons, forms)
└── __tests__/         # Unit tests
```

## 🎨 UI Components

### **shadcn/ui Integration**

The project uses [shadcn/ui](https://ui.shadcn.com/) for consistent, accessible components:

- **Layout**: Card, Separator, Sheet, Sidebar
- **Navigation**: Button, Navigation Menu, Dropdown Menu
- **Forms**: Input, Label, Select, Checkbox, Form
- **Feedback**: Badge, Progress, Toast (Sonner), Tooltip
- **Data Display**: Table, Avatar, Collapsible
- **Overlays**: Dialog, Tooltip, Hover Card

### **Theming**

- **CSS Variables** - Dynamic theming support
- **Tailwind Integration** - Utility-first styling
- **Responsive Design** - Mobile-first approach
- **Dark/Light Mode** - Built-in theme switching

## 🔧 TypeScript Configuration

### **Strict Mode Enabled**

```json
{
  "strict": true,
  "noUnusedLocals": true,
  "noUnusedParameters": true,
  "noUncheckedIndexedAccess": true,
  "exactOptionalPropertyTypes": true
}
```

### **Type Organization**

- **Global Types** (`/src/types/`) - Shared across modules
- **Module Types** (`modules/*/types.ts`) - Module-specific definitions
- **Centralized Imports** - `import { AWSTypes, KubernetesTypes } from '@/types'`

### **Benefits**

- ✅ **Zero `any` types** - Strict typing throughout
- ✅ **Compile-time safety** - Catch errors before runtime
- ✅ **IntelliSense** - Enhanced developer experience
- ✅ **Refactoring safety** - Confident code changes

## 🔍 Code Quality

### **ESLint Configuration**

- **@typescript-eslint** - TypeScript-specific rules
- **react-hooks** - React hooks validation
- **prettier** - Code formatting integration
- **unused-imports** - Automatic import cleanup

### **Prettier Setup**

- **Consistent formatting** - 80-character line width
- **Single quotes** - JavaScript/TypeScript preference
- **Trailing commas** - ES5 compatibility
- **Auto-formatting** - Format on save (VS Code)

### **Quality Gates**

```bash
# Complete quality check (required for CI/CD)
yarn quality

# Automated fixes
yarn fix
```

## 🧩 Plugin Architecture

### **Dynamic Module Loading**

The application supports dynamic plugin loading:

```typescript
// Auto-discovery of enabled plugins
const plugins = await getPlugins();

// Dynamic import based on plugin configuration
const module = await import(`../modules/${pluginName}/index.tsx`);
```

### **Module Configuration**

```typescript
interface ModuleConfig {
  menus?: Menu[]; // Sidebar navigation items
  routers?: Router[]; // Route definitions
  oAuth2?: OAuth2Config; // Authentication settings
}
```

### **Supported Modules**

- **Dashboard** - Overview and metrics
- **AWS** - EC2 instance management
- **Kubernetes** - Cluster operations
- **OAuth2** - User authentication and profiles

## 🔐 Authentication

### **OAuth2 Flow**

- **Provider**: GitHub OAuth2
- **Token Storage**: Browser localStorage with namespace
- **Auto-refresh**: Token validation and cleanup
- **Protected Routes**: Automatic redirection to login

### **User Management**

- **Profile Page** - User information and permissions
- **Team Management** - Organization and team membership
- **Permission System** - Role-based access control

## 🌐 API Integration

### **HTTP Client**

- **Base URL**: Configurable via `VITE_API_URL`
- **Interceptors**: Automatic token injection and error handling
- **Cancel Tokens**: Request cancellation for component unmount
- **Error Handling**: Centralized error management with toast notifications

### **Resource Pattern**

```typescript
// Consistent API resource functions
export function getResources(
  filter: ResourceFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<Resource[]>>;
```

## 🎛️ State Management

### **Local State with Reducers**

```typescript
interface ResourceState {
  data: Resource[];
  loading: boolean;
}

type ResourceAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Resource[] };
```

### **Patterns**

- **Loading States** - Consistent loading indicators
- **Error Boundaries** - Graceful error handling
- **Data Fetching** - AbortController for cleanup
- **Auto-refresh** - Configurable intervals

## 📱 Responsive Design

### **Breakpoints**

- **Mobile** - `< 768px`
- **Tablet** - `768px - 1024px`
- **Desktop** - `> 1024px`

### **Layout Features**

- **Collapsible Sidebar** - Space-efficient navigation
- **Responsive Tables** - Horizontal scrolling on mobile
- **Grid Layouts** - Adaptive column counts
- **Touch-friendly** - Mobile interaction patterns

## 🧪 Testing

### **Test Setup**

- **Vitest** - Fast unit test runner
- **Testing Library** - React component testing
- **jsdom** - Browser environment simulation
- **Coverage Reports** - Code coverage analysis

### **Test Organization**

```
__tests__/              # Test files co-located with source
├── Component.test.jsx  # Component tests
├── resource.test.js    # API function tests
└── hook.test.jsx       # Custom hook tests
```

## 🛠️ Development Tools

### **VS Code Integration**

```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": "explicit",
    "source.organizeImports": "explicit"
  }
}
```

### **Recommended Extensions**

- **Prettier** - Code formatting
- **ESLint** - Linting and code quality
- **Tailwind CSS IntelliSense** - CSS class completion
- **TypeScript Hero** - Import management
- **Auto Rename Tag** - HTML/JSX tag renaming

## 📋 Configuration Files

### **Build & Development**

- `vite.config.js` - Vite build configuration
- `tsconfig.json` - TypeScript project config
- `tsconfig.app.json` - App-specific TypeScript settings

### **Code Quality**

- `eslint.config.js` - ESLint rules and plugins
- `.prettierrc` - Prettier formatting rules
- `.editorconfig` - Editor consistency settings

### **Package Management**

- `package.json` - Dependencies and scripts
- `yarn.lock` - Locked dependency versions
- `.npmrc` - NPM configuration (forces Yarn usage)

## 🚀 Deployment

### **Build Process**

1. **Type Check** - `tsc --noEmit` validates all TypeScript
2. **Vite Build** - Optimized production bundle
3. **Asset Optimization** - Minification and compression

### **Environment Variables**

```bash
VITE_API_URL=https://api.dashops.example.com  # Backend API endpoint
```

### **CI/CD Integration**

```bash
# Quality gate for pull requests
yarn quality

# Production build
yarn build
```

## 📖 Developer Guide

### **Adding New Components**

```typescript
// 1. Create component with proper typing
interface ComponentProps {
  data: DataType;
  onAction: (id: string) => void;
}

export function Component({ data, onAction }: ComponentProps): JSX.Element {
  return <div>{/* component JSX */}</div>;
}

// 2. Add types to appropriate types file
// 3. Export from module if needed
```

### **Adding New Modules**

```typescript
// 1. Create module directory structure
modules/new-module/
├── index.tsx          # Module configuration
├── types.ts           # Module-specific types
├── NewModulePage.tsx  # Main page component
├── newModuleResource.ts # API functions
└── __tests__/         # Unit tests

// 2. Export from module index
export default {
  menus: Menu[],
  routers: Router[],
};
```

### **State Management Pattern**

```typescript
// 1. Define state and actions in types
interface EntityState {
  data: Entity[];
  loading: boolean;
}

type EntityAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Entity[] };

// 2. Create reducer
function reducer(state: EntityState, action: EntityAction): EntityState {
  switch (action.type) {
    case 'LOADING':
      return { ...state, loading: true, data: [] };
    case 'SET_DATA':
      return { ...state, loading: false, data: action.response };
    default:
      return state;
  }
}

// 3. Use in component
const [state, dispatch] = useReducer(reducer, INITIAL_STATE);
```

## 🔧 Configuration

### **Environment Setup**

```bash
# .env.local
VITE_API_URL=http://localhost:8080/api
```

### **VS Code Settings**

The project includes optimized VS Code settings for:

- Auto-formatting on save
- Import organization
- ESLint auto-fix
- TypeScript IntelliSense

### **Tailwind Configuration**

- Custom color schemes
- Component-specific utilities
- Responsive breakpoints
- Dark mode support

## 🧪 Testing Strategy

### **Unit Tests**

```bash
yarn test              # Run all tests
yarn coverage          # Generate coverage report
```

### **Test Patterns**

- **Component Testing** - User interaction testing
- **Hook Testing** - Custom hook validation
- **API Testing** - Resource function mocking
- **Integration Testing** - End-to-end workflows

## 📈 Performance Monitoring

### **Bundle Analysis**

- **Webpack Bundle Analyzer** - Dependency visualization
- **Build Metrics** - Size tracking over time
- **Code Splitting** - Optimal loading strategies

### **Runtime Monitoring**

- **Error Boundaries** - Graceful error handling
- **Performance Profiling** - React DevTools integration
- **Memory Management** - Cleanup on unmount

## 🔒 Security

### **Authentication Security**

- **Token Storage** - Secure localStorage with namespace
- **Auto-logout** - Token expiration handling
- **CSRF Protection** - Token-based request validation

### **Code Security**

- **TypeScript Strict** - Compile-time safety
- **ESLint Security Rules** - Security best practices
- **Dependency Auditing** - Regular security updates

## 🚀 Production Deployment

### **Build Optimization**

```bash
# Production build with all optimizations
yarn build

# Output: dist/ directory with optimized assets
dist/
├── index.html
├── assets/
│   ├── index-[hash].js   # Main application bundle
│   ├── index-[hash].css  # Compiled styles
│   └── [chunk]-[hash].js # Code-split chunks
```

### **Performance Features**

- **Tree Shaking** - Dead code elimination
- **Code Splitting** - Dynamic imports
- **Asset Optimization** - Minification and compression
- **Caching Strategy** - Long-term caching with hash names

## 🤝 Contributing

### **Code Quality Standards**

- **TypeScript** - All new code must be properly typed
- **Testing** - Unit tests required for new features
- **Formatting** - Prettier formatting enforced
- **Linting** - ESLint rules must pass

### **Development Workflow**

```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Develop with quality checks
yarn type-check:watch  # Terminal 1 - Type checking
yarn dev               # Terminal 2 - Development server

# 3. Quality gate before commit
yarn quality

# 4. Auto-fix any issues
yarn fix

# 5. Commit and push
git commit -m "feat: add new feature"
git push origin feature/new-feature
```

### **Pull Request Requirements**

- ✅ All tests passing
- ✅ Type checking successful
- ✅ Zero ESLint warnings
- ✅ Code properly formatted
- ✅ Build successful

## 📚 Additional Resources

- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Guide](https://vitejs.dev/guide/)
- [shadcn/ui Components](https://ui.shadcn.com/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)

## 🎊 Project Status

- ✅ **TypeScript Migration** - 100% complete
- ✅ **Type Safety** - Zero `any` types in production code
- ✅ **UI Components** - Full shadcn/ui integration
- ✅ **Quality Tools** - ESLint + Prettier + Strict TS
- ✅ **Build Pipeline** - Type checking integrated
- ✅ **Documentation** - Comprehensive setup guide

**Ready for production deployment and team collaboration!** 🚀
