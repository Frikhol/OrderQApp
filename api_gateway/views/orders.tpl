<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Orders - OrderQ</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/orders.css" rel="stylesheet">
</head>
<body>
    <div class="hero-section text-center">
        <div class="container">
            <h1 class="display-4 mb-4">Manage Your Orders</h1>
            <p class="lead mb-4">Create and track your queue service orders</p>
        </div>
    </div>

    <div class="container mb-5">
        <div class="action-buttons text-center">
            <button class="btn btn-success btn-lg me-2" id="showCreateForm">
                <i class="fas fa-plus-circle me-2"></i>Create New Order
            </button>
            <button class="btn btn-outline-primary btn-lg" id="showOrdersList">
                <i class="fas fa-list me-2"></i>View My Orders
            </button>
        </div>

        <!-- Include form template -->
        {{template "order_create_form.tpl" .}}

        <!-- Include list template -->
        {{template "order_list.tpl" .}}
    </div>

    <footer class="bg-light py-4 mt-auto">
        <div class="container text-center">
            <p class="mb-0">© 2024 OrderQ. All rights reserved.</p>
        </div>
    </footer>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script src="/static/js/orders.js"></script>
</body>
</html> 