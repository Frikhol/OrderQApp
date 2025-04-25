<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <meta name="theme-color" content="#0d6efd">
    <title>OrderQApp</title>
    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
    <!-- Custom CSS -->
    <style>
        body {
            background-color: #f8f9fa;
            min-height: 100vh;
            padding-bottom: 60px; /* Space for mobile bottom navigation */
        }
        .navbar {
            margin-bottom: 1rem;
            box-shadow: 0 2px 4px rgba(0,0,0,.1);
        }
        .card {
            border: none;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0,0,0,.1);
            margin-bottom: 1rem;
        }
        .card-body {
            padding: 1.5rem;
        }
        .btn-primary {
            padding: 0.75rem;
            background-color: #0d6efd;
            border-color: #0d6efd;
            width: 100%;
        }
        .btn-primary:hover {
            background-color: #0b5ed7;
            border-color: #0a58ca;
        }
        .form-control:focus {
            border-color: #86b7fe;
            box-shadow: 0 0 0 0.25rem rgba(13,110,253,.25);
        }
        /* Mobile Navigation */
        .mobile-nav {
            display: none;
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            background: #fff;
            box-shadow: 0 -2px 10px rgba(0,0,0,.1);
            z-index: 1000;
        }
        .mobile-nav-item {
            padding: 0.75rem;
            text-align: center;
            color: #6c757d;
            text-decoration: none;
        }
        .mobile-nav-item.active {
            color: #0d6efd;
        }
        .mobile-nav-item i {
            display: block;
            font-size: 1.5rem;
            margin-bottom: 0.25rem;
        }
        @media (max-width: 768px) {
            .mobile-nav {
                display: flex;
                justify-content: space-around;
            }
            .navbar {
                margin-bottom: 0;
            }
            .container {
                padding: 0 1rem;
            }
            .card-body {
                padding: 1rem;
            }
        }
    </style>
</head>
<body class="d-flex flex-column">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand fw-bold" href="/">OrderQApp</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav ms-auto">
                    {{if .IsAuthenticated}}
                        <li class="nav-item">
                            <a class="nav-link" href="/dashboard">Dashboard</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link" href="/auth/logout">Logout</a>
                        </li>
                    {{else}}
                        <li class="nav-item">
                            <a class="nav-link" href="/auth/login">Login</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link" href="/auth/register">Register</a>
                        </li>
                    {{end}}
                </ul>
            </div>
        </div>
    </nav>

    <main class="flex-grow-1">
        <div class="container mt-3">
             {{template "content" .}}
        </div>
    </main>

    <!-- Mobile Navigation -->
    <nav class="mobile-nav">
        <a href="/dashboard" class="mobile-nav-item">
            <i class="fas fa-home"></i>
            <span>Home</span>
        </a>
        <a href="/orders" class="mobile-nav-item">
            <i class="fas fa-list"></i>
            <span>Orders</span>
        </a>
        <a href="/profile" class="mobile-nav-item">
            <i class="fas fa-user"></i>
            <span>Profile</span>
        </a>
    </nav>

    <!-- Font Awesome -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <!-- Bootstrap Bundle with Popper -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"></script>
    <!-- Custom JavaScript -->
    <script src="/static/js/common.js"></script>
</body>
</html> 