<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - OrderQ API Gateway</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/css/auth.css" rel="stylesheet">
</head>
<body class="bg-light">
    <div class="container">
        <div class="auth-container">
            <div class="auth-form">
                <div class="auth-header">
                    <h2 class="mb-3">Login</h2>
                    <p class="text-muted">Please sign in to access the OrderQ API Gateway</p>
                </div>
                <div class="alert alert-danger" role="alert" id="errorAlert"></div>
                <form id="loginForm">
                    <div class="mb-3">
                        <label for="email" class="form-label">Email</label>
                        <input type="email" class="form-control" id="email" name="email" required>
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">Password</label>
                        <input type="password" class="form-control" id="password" name="password" required>
                    </div>
                    <div class="d-grid gap-2">
                        <button type="submit" class="btn btn-primary">Sign In</button>
                    </div>
                </form>
                <div class="auth-link">
                    <p class="mb-0">Don't have an account? <a href="/auth/register">Create one</a></p>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/auth.js"></script>
</body>
</html> 