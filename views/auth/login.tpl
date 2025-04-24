{{template "base.tpl" .}}

{{define "content"}}
<div class="row justify-content-center mt-5">
    <div class="col-md-6 col-lg-4">
        <div class="card shadow">
            <div class="card-body p-4">
                <h3 class="text-center mb-4">Login</h3>
                <div id="errorMessage" class="alert alert-danger d-none"></div>
                <div id="successMessage" class="alert alert-success d-none"></div>
                <form id="loginForm" method="post">
                    <div class="mb-3">
                        <label for="email" class="form-label">Email</label>
                        <input type="email" class="form-control" id="email" name="email" required>
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">Password</label>
                        <input type="password" class="form-control" id="password" name="password" required>
                    </div>
                    <div class="d-grid gap-2">
                        <button type="submit" class="btn btn-primary">Login</button>
                    </div>
                </form>
                <div class="text-center mt-3">
                    <p class="mb-0">Don't have an account? <a href="/auth/register">Register</a></p>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
function showMessage(message, type = 'error') {
    const errorDiv = document.getElementById('errorMessage');
    const successDiv = document.getElementById('successMessage');
    
    if (type === 'error') {
        errorDiv.textContent = message;
        errorDiv.classList.remove('d-none');
        successDiv.classList.add('d-none');
    } else {
        successDiv.textContent = message;
        successDiv.classList.remove('d-none');
        errorDiv.classList.add('d-none');
    }
}

document.getElementById('loginForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    // Basic validation
    if (!email || !password) {
        showMessage('Please fill in all fields');
        return;
    }
    
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    
    fetch('/auth/login', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(data => {
                throw new Error(data.error || 'Login failed');
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.error) {
            showMessage(data.error);
        } else if (data.token) {
            // Store the token and user data
            localStorage.setItem('token', data.token);
            localStorage.setItem('role', data.role);
            localStorage.setItem('user_id', data.user_id);
            
            // Show success message
            showMessage('Login successful! Redirecting...', 'success');
            
            // Redirect to dashboard based on role
            setTimeout(() => {
                window.location.href = '/dashboard';
            }, 1000);
        }
    })
    .catch(error => {
        showMessage(error.message || 'An error occurred during login');
    });
});
</script>
{{end}}
