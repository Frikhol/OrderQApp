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
                })
            });

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
                    showAlert(errorAlert, 'Unable to sign in. Please check your credentials and try again.');
                }
                setButtonSubmitting(submitButton, false);
                return;
            }

            // Store the token in localStorage
            if (data.token) {
                localStorage.setItem('token', data.token);
            }

            // Show success message before redirect
            showAlert(errorAlert, 'Successfully signed in! Redirecting...', false);
            setTimeout(() => {
                window.location.href = '/';
            }, 1000);
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
                    showAlert(errorAlert, 'Unable to create account. Please try again.');
                }
                setButtonSubmitting(submitButton, false);
                return;
            }

            showAlert(successAlert, 'Account created successfully! You can now sign in.', false);
            registerForm.reset();
            setButtonSubmitting(submitButton, false);
        } catch (error) {
            showAlert(errorAlert, 'Something went wrong. Please try again later.');
            setButtonSubmitting(submitButton, false);
        }
    });
}

// Initialize forms when the page loads
document.addEventListener('DOMContentLoaded', function() {
    setupLoginForm();
    setupRegisterForm();
}); 