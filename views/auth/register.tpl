{{template "base.tpl" .}}

{{define "content"}}
<div class="row justify-content-center mt-5">
    <div class="col-md-6 col-lg-4">
        <div class="card shadow">
            <div class="card-body p-4">
                <h3 class="text-center mb-4">Register</h3>
                <div id="errorMessage" class="alert alert-danger d-none"></div>
                <div id="successMessage" class="alert alert-success d-none"></div>
                <form id="registerForm" method="post">
                    <div class="mb-3">
                        <label for="email" class="form-label">Email</label>
                        <input type="email" class="form-control" id="email" name="email" required>
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">Password</label>
                        <input type="password" class="form-control" id="password" name="password" required>
                    </div>
                    <div class="mb-3">
                        <label for="role" class="form-label">Role</label>
                        <select class="form-select" id="role" name="role" required>
                            <option value="">Select role</option>
                            <option value="client">Client</option>
                            <option value="agent">Agent</option>
                        </select>
                    </div>
                    <div class="d-grid gap-2">
                        <button type="submit" class="btn btn-primary">Register</button>
                    </div>
                </form>
                <div class="text-center mt-3">
                    <p class="mb-0">Already have an account? <a href="/auth/login">Login</a></p>
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

function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

document.getElementById('registerForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const role = document.getElementById('role').value;
    
    // Client-side validation
    if (!email || !password || !role) {
        showMessage('Please fill in all fields');
        return;
    }
    
    if (!validateEmail(email)) {
        showMessage('Please enter a valid email address');
        return;
    }
    
    if (password.length < 6) {
        showMessage('Password must be at least 6 characters long');
        return;
    }
    
    if (role !== 'client' && role !== 'agent') {
        showMessage('Please select a valid role');
        return;
    }
    
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    formData.append('role', role);
    
    fetch('/auth/register', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(data => {
                throw new Error(data.error || 'Registration failed');
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.error) {
            showMessage(data.error);
        } else {
            showMessage('Registration successful! Redirecting to login...', 'success');
            setTimeout(() => {
                window.location.href = '/auth/login';
            }, 1500);
        }
    })
    .catch(error => {
        showMessage(error.message || 'An error occurred during registration');
    });
});
</script>
{{end}} 