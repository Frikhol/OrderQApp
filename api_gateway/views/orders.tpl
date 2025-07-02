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
    <!-- Add hidden input for user ID -->
    <input type="hidden" id="userId" value="{{.user_id}}">
    
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

        <!-- Order Details Section -->
        <div id="orderDetailsSection" style="display: none;">
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2 class="mb-0">Order Details</h2>
                <button class="btn btn-outline-secondary" id="backToList">
                    <i class="fas fa-arrow-left me-2"></i>Back to Orders
                </button>
            </div>
            <div class="row">
                <div class="col-md-6 col-lg-4 mx-auto">
                    <div class="card h-100">
                        <div class="card-body">
                            <h5 class="card-title" id="orderLocation">Location</h5>
                            <h6 class="card-subtitle mb-3 text-muted" id="orderAddress">Address</h6>
                            <div class="d-flex align-items-center mb-2">
                                <i class="fas fa-calendar-alt text-success me-2"></i>
                                <span id="orderDate">Date</span>
                            </div>
                            <div class="d-flex align-items-center mb-2">
                                <i class="fas fa-clock text-success me-2"></i>
                                <span id="orderTimeGap">Time</span>
                            </div>
                            <div class="d-flex align-items-center">
                                <i class="fas fa-tag text-success me-2"></i>
                                <span id="orderStatus">Status</span>
                            </div>
                        </div>
                        <div class="card-footer bg-white">
                            <div class="d-grid gap-2">
                                <button type="button" class="btn btn-outline-danger" id="cancelOrderBtn">
                                    <i class="fas fa-times-circle me-2"></i>Cancel Order
                                </button>
                                <button type="button" class="btn btn-outline-primary" id="finishOrderBtn">
                                    <i class="fas fa-check-circle me-2"></i>Finish Order
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
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