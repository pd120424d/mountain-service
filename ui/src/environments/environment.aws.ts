export const environment = {
    production: true,
    staging: false,
    apiUrl: '/api/v1',
    useMockApi: false,
    
    // AWS-specific configuration
    aws: {
        region: 'us-east-1', // Update with your AWS region
        enabled: true
    },
    
    // Feature flags for AWS deployment
    features: {
        enableAnalytics: true,
        enableErrorReporting: true,
        enablePerformanceMonitoring: true,
        enableOfflineMode: false
    },
    
    // API endpoints
    endpoints: {
        employee: '/api/v1/employees',
        urgency: '/api/v1/urgencies',
        activity: '/api/v1/activities',
        version: '/api/v1/version',
        health: '/api/v1/health',
        swagger: {
            employee: '/employee-swagger/',
            urgency: '/urgency-swagger/',
            activity: '/activity-swagger/'
        }
    },
    
    // Security configuration
    security: {
        enableCSP: true,
        enableHTTPS: false, // Set to true if using SSL
        tokenRefreshThreshold: 300000, // 5 minutes in milliseconds
        sessionTimeout: 3600000 // 1 hour in milliseconds
    },
    
    // Logging configuration
    logging: {
        level: 'warn', // 'debug', 'info', 'warn', 'error'
        enableConsoleLogging: false,
        enableRemoteLogging: true
    },
    
    // Performance configuration
    performance: {
        enableLazyLoading: true,
        enableServiceWorker: false, // Set to true if you want PWA features
        cacheTimeout: 300000 // 5 minutes
    }
};
