// Error Logger module
const ErrorLogger = (function() {
    // Log error to server
    function log(error, context) {
        const errorData = {
            timestamp: new Date().toISOString(),
            context: context,
            error: error.message || error,
            stack: error.stack,
            userAgent: navigator.userAgent,
            url: window.location.href
        };

        // Log to console
        console.error('Error in context:', context, errorData);

        // Send to server for logging
        fetch('/api/logs/error', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            },
            body: JSON.stringify(errorData)
        }).catch(logError => {
            console.error('Failed to log error to server:', logError);
        });
    }

    // Initialize error handlers
    function init() {
        // Handle unhandled promise rejections
        window.addEventListener('unhandledrejection', function(event) {
            log(event.reason, 'Unhandled Promise Rejection');
        });

        // Handle global errors
        window.addEventListener('error', function(event) {
            log(event.error || event.message, 'Global Error');
        });
    }

    // Public API
    return {
        log: log,
        init: init
    };
})();

// Initialize error logger when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    ErrorLogger.init();
}); 