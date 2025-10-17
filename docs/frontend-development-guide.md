# Frontend Development Guide

## Overview

This guide outlines the standardized architecture and development patterns for the DashOps frontend application. Following these patterns ensures consistency, maintainability, and scalability across all modules.

## Architecture Principles

### Functional Programming Approach
- **No Classes**: Use only functional components and pure functions
- **Custom Hooks**: Encapsulate stateful logic in reusable hooks
- **Pure Functions**: Prefer pure functions for data transformation and utilities
- **Immutable State**: Avoid direct state mutations

### Module Structure
Each module follows a consistent directory structure:

```
src/modules/{module-name}/
├── types.ts                    # TypeScript interfaces
├── resources/                  # HTTP request functions
├── adapters/                   # Data transformation functions
├── hooks/                      # Custom React hooks
├── components/                 # React components
├── utils/                      # Utility functions
├── __tests__/                  # Unit tests
└── index.tsx                   # Module exports
```

## Directory Responsibilities

### `types.ts`
- Define TypeScript interfaces that match backend API contracts
- Export all types needed by external consumers
- Keep types focused and specific to the module domain

```typescript
// Example: AWS module types
export interface Instance {
  instance_id: string;
  name: string;
  state: InstanceState;
  // ... other properties
}

export interface InstanceState {
  name: string;
  code: number;
}
```

### `resources/`
- Contains functions that make HTTP requests to backend APIs
- Pure functions that return promises
- Handle API communication and error responses

```typescript
// Example: instanceResource.ts
import { http } from '../../../helpers/http';
import { AWSTypes } from '../types';

export const getInstances = async (
  filter: { accountKey: string },
  config?: { signal?: AbortSignal }
): Promise<{ data: any }> => {
  const response = await http.get(`/aws/accounts/${filter.accountKey}/instances`, config);
  return response.data;
};
```

### `adapters/`
- Pure functions for transforming data between API responses and domain models
- Handle data normalization and formatting
- Keep transformation logic separate from API calls

```typescript
// Example: instanceAdapter.ts
import { AWSTypes } from '../types';

export const transformInstanceToDomain = (apiInstance: any): AWSTypes.Instance => {
  return {
    instance_id: apiInstance.instance_id,
    name: apiInstance.name,
    state: {
      name: apiInstance.state_name,
      code: apiInstance.state_code
    }
    // ... transform other fields
  };
};
```

### `hooks/`
- Custom React hooks for state management and side effects
- Encapsulate complex logic and provide clean interfaces
- Handle loading states, errors, and data fetching

```typescript
// Example: useInstances.ts
import { useState, useEffect, useCallback } from 'react';
import { getInstances } from '../resources/instanceResource';
import { transformInstancesToDomain } from '../adapters/instanceAdapter';

export const useInstances = (filter: { accountKey: string }) => {
  const [instances, setInstances] = useState<AWSTypes.Instance[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchInstances = useCallback(async () => {
    // ... implementation
  }, [filter.accountKey]);

  return { instances, loading, error, fetchInstances };
};
```

### `components/`
- React functional components
- Focus on presentation and user interaction
- Use hooks for state management and side effects
- Keep components small and focused

```typescript
// Example: InstancePage.tsx
import { useInstances } from '../hooks/useInstances';
import { InstanceTag } from './InstanceTag';

export default function InstancePage() {
  const { instances, loading, error } = useInstances(filter);
  
  // ... component implementation
}
```

### `utils/`
- Pure utility functions for formatting, validation, and calculations
- No side effects or external dependencies
- Easy to test and reuse

```typescript
// Example: instanceUtils.ts
export const getInstanceDisplayName = (instance: AWSTypes.Instance): string => {
  return instance.name || instance.instance_id;
};

export const canStartInstance = (instance: AWSTypes.Instance): boolean => {
  return instance.state.name === 'stopped';
};
```

### `__tests__/`
- Comprehensive unit tests for all layers
- Test coverage for adapters, utils, hooks, and components
- Use mocking for external dependencies

```typescript
// Example: instanceAdapter.test.ts
import { describe, it, expect } from 'vitest';
import { transformInstanceToDomain } from '../adapters/instanceAdapter';

describe('instanceAdapter', () => {
  it('should transform API response to domain model', () => {
    const apiInstance = { instance_id: 'i-123', name: 'test' };
    const result = transformInstanceToDomain(apiInstance);
    
    expect(result.instance_id).toBe('i-123');
    expect(result.name).toBe('test');
  });
});
```

### `index.tsx`
- Main module export file
- Export only what external consumers need
- Maintain proper encapsulation

```typescript
// Example: AWS module index.tsx
export * from './types';
// Export only necessary components and hooks
export { default as InstancePage } from './components/instances/InstancePage';
```

## Testing Strategy

### Test Structure
- **Unit Tests**: Test individual functions and components in isolation
- **Integration Tests**: Test hooks and component interactions
- **Mocking**: Use mocks for external dependencies and API calls

### Test Coverage
- Aim for 100% test coverage on utility functions and adapters
- Test all hook behaviors including loading, error, and success states
- Test component rendering and user interactions

### Global Test Setup
- Use `src/setupTests.js` for global test configuration
- Create `src/__mocks__/` for shared mocks
- Avoid module-specific test configuration files

## Development Workflow

### Creating a New Module
1. Create the directory structure following the standard pattern
2. Start with `types.ts` defining the interfaces
3. Implement `resources/` for API communication
4. Create `adapters/` for data transformation
5. Build `hooks/` for state management
6. Develop `components/` for UI
7. Add `utils/` for helper functions
8. Write comprehensive tests
9. Export only necessary items in `index.tsx`

### Refactoring Existing Modules
1. Analyze current structure and dependencies
2. Create new directory structure
3. Move files to appropriate directories
4. Refactor class components to functional components
5. Extract custom hooks from component logic
6. Add comprehensive test coverage
7. Update exports and imports
8. Verify all functionality works correctly

## Best Practices

### Code Organization
- Keep functions small and focused on a single responsibility
- Use descriptive names for functions and variables
- Group related functionality in the same directory
- Maintain consistent file naming conventions

### Type Safety
- Always define TypeScript interfaces for data structures
- Use strict TypeScript configuration
- Avoid `any` types; prefer specific interfaces
- Align frontend types with backend API contracts

### Performance
- Use `useCallback` and `useMemo` for expensive operations
- Implement proper cleanup in useEffect hooks
- Avoid unnecessary re-renders with proper dependency arrays
- Use React.memo for pure components when appropriate

### Error Handling
- Implement proper error boundaries
- Handle API errors gracefully
- Provide meaningful error messages to users
- Log errors for debugging purposes

## Module Examples

### AWS Module (Reference Implementation)
The AWS module serves as the reference implementation of this architecture:

- **67 tests passing (100% success rate)**
- Complete separation of concerns
- Comprehensive test coverage
- Proper encapsulation and exports

This module can be used as a template for refactoring other modules (Kubernetes, Service Catalog, OAuth2, Config).

## Migration Guide

### From Class-Based to Functional Architecture
1. Convert class components to functional components
2. Extract state logic into custom hooks
3. Move data transformation to adapter functions
4. Separate API calls into resource functions
5. Create utility functions for common operations
6. Add comprehensive test coverage
7. Update module exports

### Benefits of Migration
- **Better Testability**: Pure functions are easier to test
- **Improved Reusability**: Hooks and utilities can be reused
- **Enhanced Maintainability**: Clear separation of concerns
- **Type Safety**: Better TypeScript support
- **Performance**: Optimized re-rendering with hooks

## Conclusion

This architecture provides a solid foundation for building maintainable, testable, and scalable frontend modules. By following these patterns consistently across all modules, we ensure code quality and developer productivity.

For questions or clarifications, refer to the AWS module implementation or reach out to the development team.
