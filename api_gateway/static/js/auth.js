// Common functions for authentication forms
function showAlert(alertElement, message, isError = true) {
    alertElement.innerHTML = `
        <div class="d-flex align-items-center">
            <i class="fas ${isError ? 'fa-exclamation-circle' : 'fa-check-circle'} me-2"></i>
            <span>${message}</span>
        </div>
    `;
    alertElement.style.display = 'block';
    if (isError) {
        alertElement.classList.remove('alert-success');
        alertElement.classList.add('alert-danger');
    } else {
        alertElement.classList.remove('alert-danger');
        alertElement.classList.add('alert-success');
    }
}

function hideAlerts(...alertElements) {
    alertElements.forEach(alert => {
        alert.style.display = 'none';
    });
}

function showFieldError(inputElement, errorElement, message) {
    inputElement.classList.add('is-invalid');
    errorElement.innerHTML = `
        <div class="d-flex align-items-center">
            <i class="fas fa-exclamation-circle me-1"></i>
            <span>${message}</span>
        </div>
    `;
    errorElement.style.display = 'block';
    
    // Add shake animation to the input group
    const inputGroup = inputElement.closest('.input-group');
    inputGroup.classList.add('shake');
    setTimeout(() => {
        inputGroup.classList.remove('shake');
    }, 500);
}

function hideFieldError(inputElement, errorElement) {
    inputElement.classList.remove('is-invalid');
    errorElement.style.display = 'none';
}

function clearFieldErrors() {
    document.querySelectorAll('.is-invalid').forEach(element => {
        element.classList.remove('is-invalid');
    });
    document.querySelectorAll('.invalid-feedback').forEach(element => {
        element.style.display = 'none';
    });
}

function setButtonSubmitting(button, isSubmitting) {
    if (isSubmitting) {
        button.classList.add('submitting');
        button.disabled = true;
        button.innerHTML = `
            <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
            ${button.getAttribute('data-original-text') || button.textContent}
        `;
    } else {
        button.classList.remove('submitting');
        button.disabled = false;
        button.innerHTML = button.getAttribute('data-original-text') || button.innerHTML;
    }
}

// Token handling functions
function getToken() {
    return document.cookie.replace(/(?:(?:^|.*;\s*)token\s*\=\s*([^;]*).*$)|^.*$/, "$1");
}

function clearToken() {
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}

// Function to check if the token is valid
async function validateToken() {
    const token = getToken();
    if (!token) {
        return false;
    }

    try {
        const response = await fetch('/auth/validate', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        return response.ok;
    } catch (error) {
        return false;
    }
}

// Handle unauthorized errors globally
function setupTokenInvalidationHandler() {
    // Handle API response errors that have 401 status
    const originalFetch = window.fetch;
    window.fetch = async function(url, options = {}) {
        const response = await originalFetch(url, options);
        
        if (response.status === 401) {
            const responseText = await response.text();
            if (responseText.includes("Invalid or expired token")) {
                // Clear the invalid token
                clearToken();
                
                // Redirect to login page
                window.location.href = '/auth/login?session_expired=true';
            }
        }
        
        return response;
    };
}

// Login form handling
function setupLoginForm() {
    const loginForm = document.getElementById('loginForm');
    if (!loginForm) return;

    const submitButton = loginForm.querySelector('button[type="submit"]');
    submitButton.setAttribute('data-original-text', submitButton.innerHTML);

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email');
        const password = document.getElementById('password');
        const emailError = document.getElementById('emailError');
        const passwordError = document.getElementById('passwordError');
        const errorAlert = document.getElementById('errorAlert');
        
        // Clear previous errors
        clearFieldErrors();
        hideAlerts(errorAlert);
        
        // Basic client-side validation
        if (!email.value) {
            showFieldError(email, emailError, 'Please enter your email address');
            return;
        }
        if (!password.value) {
            showFieldError(password, passwordError, 'Please enter your password');
            return;
        }
        
        setButtonSubmitting(submitButton, true);
        
        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email.value,
                    password: password.value
                }),
                redirect: 'follow'
            });
            
            // For successful logins, the backend will redirect to /orders (302 status)
            if (response.redirected) {
                window.location.href = response.url;
                return;
            }
            
            // Handle JSON responses for error cases
            if (response.headers.get('content-type')?.includes('application/json')) {
                const data = await response.json();
                
                if (!response.ok) {
                    if (data.errors) {
                        // Handle field-level errors
                        if (data.errors.email) {
                            showFieldError(email, emailError, 'Please check your email address');
                        }
                        if (data.errors.password) {
                            showFieldError(password, passwordError, 'Incorrect password. Please try again');
                        }
                    } else {
                        showAlert(errorAlert, data.error || 'Unable to sign in. Please check your credentials and try again.');
                    }
                    setButtonSubmitting(submitButton, false);
                    return;
                }
            }        
        } catch (error) {
            showAlert(errorAlert, 'Something went wrong. Please try again later.');
            setButtonSubmitting(submitButton, false);
        }
    });
}

// Register form handling
function setupRegisterForm() {
    const registerForm = document.getElementById('registerForm');
    if (!registerForm) return;

    const submitButton = registerForm.querySelector('button[type="submit"]');
    submitButton.setAttribute('data-original-text', submitButton.innerHTML);

    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email');
        const password = document.getElementById('password');
        const confirmPassword = document.getElementById('confirmPassword');
        const emailError = document.getElementById('emailError');
        const passwordError = document.getElementById('passwordError');
        const confirmPasswordError = document.getElementById('confirmPasswordError');
        const errorAlert = document.getElementById('errorAlert');
        const successAlert = document.getElementById('successAlert');
        
        // Clear previous errors
        clearFieldErrors();
        hideAlerts(errorAlert, successAlert);
        
        // Client-side validation
        if (!email.value) {
            showFieldError(email, emailError, 'Please enter your email address');
            return;
        }
        if (!password.value) {
            showFieldError(password, passwordError, 'Please create a password');
            return;
        }
        if (password.value.length < 8) {
            showFieldError(password, passwordError, 'Password must be at least 8 characters long');
            return;
        }
        if (password.value !== confirmPassword.value) {
            showFieldError(confirmPassword, confirmPasswordError, 'Passwords do not match. Please try again');
            return;
        }
        
        setButtonSubmitting(submitButton, true);
        
        try {
            const response = await fetch('/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email.value,
                    password: password.value
                })
            });

            const data = await response.json();
            
            if (!response.ok) {
                if (data.errors) {
                    // Handle field-level errors
                    if (data.errors.email) {
                        showFieldError(email, emailError, 'This email is already registered or invalid');
                    }
                    if (data.errors.password) {
                        showFieldError(password, passwordError, 'Password must be at least 8 characters long');
                    }
                } else {
                    showAlert(errorAlert, data.error || 'Unable to create account. Please try again.');
                }
                setButtonSubmitting(submitButton, false);
                return;
            }

            showAlert(successAlert, 'Account created successfully! You will be redirected to the login page.', false);
            registerForm.reset();
            
            // If a redirect URL is provided, redirect after a short delay
            if (data.redirect) {
                setTimeout(() => {
                    window.location.href = data.redirect;
                }, 1000); // 2 seconds delay to allow user to see the success message
            } else {
                setButtonSubmitting(submitButton, false);
            }
        } catch (error) {
            showAlert(errorAlert, 'Something went wrong. Please try again later.');
            setButtonSubmitting(submitButton, false);
        }
    });
}

// Initialize forms and handlers when the page loads
document.addEventListener('DOMContentLoaded', function() {
    setupLoginForm();
    setupRegisterForm();
    setupTokenInvalidationHandler();
    
    // Check for session_expired parameter in URL
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('session_expired') === 'true') {
        const errorAlert = document.getElementById('errorAlert');
        if (errorAlert) {
            showAlert(errorAlert, 'Your session has expired. Please sign in again.');
        }
    }
}); 