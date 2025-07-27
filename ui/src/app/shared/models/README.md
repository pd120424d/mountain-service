# Generated Models

This directory contains TypeScript models automatically generated from backend Swagger/OpenAPI specifications.

## Structure

```
shared/models/
├── generated/           # Auto-generated models from swagger specs
│   ├── employee/        # Employee service models
│   │   ├── data-contracts.ts
│   │   ├── http-client.ts
│   │   └── index.ts
│   ├── urgency/         # Urgency service models
│   │   ├── data-contracts.ts
│   │   ├── http-client.ts
│   │   └── index.ts
│   └── activity/        # Activity service models
│       ├── data-contracts.ts
│       ├── http-client.ts
│       └── index.ts
├── extensions/          # Manual extensions and utilities
│   ├── employee-extensions.ts
│   └── index.ts
├── index.ts            # Main barrel export
└── README.md           # This file
```

## Usage

Import models from the shared location:

```typescript
// Good - Use shared models
import { Employee, EmployeeCreateRequest, MedicRole } from '../shared/models';

// Bad - Don't import from old locations
import { Employee } from '../employee/employee.model';
```

## Available Models

### Employee Service
- `Employee` (alias for `EmployeeResponse`)
- `EmployeeCreateRequest`
- `EmployeeUpdateRequest`
- `EmployeeLogin`
- `TokenResponse`
- `ErrorResponse`
- `MessageResponse`

### Urgency Service
- `Urgency` (alias for `UrgencyResponse`)
- `UrgencyCreateRequest`
- `UrgencyUpdateRequest`
- `UrgencyLevel` (legacy enum)
- `Status` (legacy enum)
- `GeneratedUrgencyLevel` (generated enum)
- `GeneratedUrgencyStatus` (generated enum)

### Activity Service
- `Activity` (alias for `ActivityResponse`)
- `ActivityCreateRequest`
- `ActivityListRequest`
- `ActivityListResponse`
- `ActivityStatsResponse`
- `ActivityType`
- `ActivityLevel`

### Role Constants
- `MedicRole`
- `TechnicalRole`
- `AdministratorRole`
- `EmployeeRole` (union type)

### Utility Functions

#### Employee Utilities
- `createDisplayName(employee)` - Creates full name
- `isAdmin(employee)` - Checks if employee is admin
- `isMedic(employee)` - Checks if employee is medic
- `isTechnical(employee)` - Checks if employee is technical

#### Urgency Utilities
- `getUrgencyLevelColor(level)` - Returns color for urgency level
- `getStatusColor(status)` - Returns color for urgency status
- `mapGeneratedLevelToLegacy(level)` - Converts generated to legacy enum
- `mapLegacyLevelToGenerated(level)` - Converts legacy to generated enum

#### Activity Utilities
- `getActivityLevelColor(level)` - Returns color for activity level
- `getActivityTypeIcon(type)` - Returns icon for activity type
- `getActivityTypeDisplayName(type)` - Returns display name for activity type
- `isSystemActivity(type)` - Checks if activity is system-related
- `isEmployeeActivity(type)` - Checks if activity is employee-related
- `isUrgencyActivity(type)` - Checks if activity is urgency-related

## Generating Models

### Automatic Generation
Run the generation script to fetch from live APIs:
```bash
npm run generate-models
```

### Manual Generation
For individual services:
```bash
# Employee service (from live API)
npm run generate-employee-models

# Employee service (from local file)
npm run generate-employee-models-local

# Urgency service (from live API)
npm run generate-urgency-models

# Urgency service (from local file)
npm run generate-urgency-models-local

# Activity service (from live API)
npm run generate-activity-models

# Activity service (from local file)
npm run generate-activity-models-local
```

### Adding New Services
1. Add service configuration to `scripts/generate-models.js`
2. Create fallback swagger file if needed
3. Run generation script
4. Update exports in `index.ts`

## Best Practices

1. **Always use generated models** - Don't create manual interfaces that duplicate backend contracts
2. **Use extensions for frontend-specific logic** - Add computed properties and utilities in `extensions/`
3. **Keep generated files untouched** - Never edit files in `generated/` directories
4. **Regenerate after backend changes** - Run generation scripts when backend APIs change
5. **Use type aliases for clarity** - Export commonly used types with cleaner names

## Troubleshooting

### Generation Fails
- Check if backend services are running
- Verify swagger endpoints are accessible
- Use fallback files for offline development

### Type Errors
- Ensure all imports use shared models
- Check if generated models match expected structure
- Regenerate models if backend contracts changed

### Missing Models
- Add new services to generation script
- Update barrel exports in index files
- Run full generation script
