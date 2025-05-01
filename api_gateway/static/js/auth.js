// Common functions for authentication forms
function showAlert(alertElement, message, isError = true) {
    alertElement.textContent = message;
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

// Login form handling
function setupLoginForm() {
    const loginForm = document.getElementById('loginForm');
    if (!loginForm) return;

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const errorAlert = document.getElementById('errorAlert');
        
        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Login failed');
            }

            // Store the token in localStorage
            if (data.token) {
                localStorage.setItem('token', data.token);
            }

            // Redirect to dashboard or home page
            window.location.href = '/';
        } catch (error) {
            showAlert(errorAlert, error.message);
        }
    });
}

// Register form handling
function setupRegisterForm() {
    const registerForm = document.getElementById('registerForm');
    if (!registerForm) return;

    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;
        const errorAlert = document.getElementById('errorAlert');
        const successAlert = document.getElementById('successAlert');
        
        // Hide alerts
        hideAlerts(errorAlert, successAlert);
        
        // Validate passwords match
        if (password !== confirmPassword) {
            showAlert(errorAlert, 'Passwords do not match');
            return;
        }
        
        try {
            const response = await fetch('/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Registration failed');
            }

            // Show success message
            showAlert(successAlert, data.message || 'Registration successful!', false);
            
            // Clear form
            registerForm.reset();
            
            // Redirect to login page after 2 seconds
            setTimeout(() => {
                window.location.href = '/auth/login';
            }, 2000);
            
        } catch (error) {
            showAlert(errorAlert, error.message);
        }
    });
}

// Initialize forms when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    setupLoginForm();
    setupRegisterForm();
}); 