<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Orders - OrderQ</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        .page-header {
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
            color: white;
            padding: 60px 0;
            margin-bottom: 40px;
        }
        .order-card {
            border: none;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            transition: transform 0.3s ease;
            margin-bottom: 20px;
        }
        .order-card:hover {
            transform: translateY(-5px);
        }
        .form-container {
            max-width: 800px;
            margin: 0 auto;
            padding: 30px;
            background-color: #fff;
            border-radius: 10px;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        }
        .btn-submit {
            padding: 10px 30px;
            font-weight: 500;
        }
        .input-group-text {
            background-color: #f8f9fa;
        }
    </style>
</head>
<body class="bg-light">
    <div class="page-header text-center">
        <div class="container">
            <h1 class="display-4 mb-4">Orders</h1>
            <p class="lead">Create and manage your queue service orders</p>
        </div>
    </div>

    <div class="container mb-5">
        <div class="form-container">
            <h2 class="mb-4 text-center">Create New Order</h2>
            <div class="alert alert-danger" role="alert" id="errorAlert" style="display: none;"></div>
            <div class="alert alert-success" role="alert" id="successAlert" style="display: none;"></div>
            
            <form id="createOrderForm">
                <div class="mb-4">
                    <label for="orderLocation" class="form-label">Location Name</label>
                    <div class="input-group">
                        <span class="input-group-text"><i class="fas fa-building"></i></span>
                        <input type="text" class="form-control" id="orderLocation" name="orderLocation" 
                               placeholder="Enter the location name (e.g., Restaurant, Bank, Government Office)" required>
                    </div>
                    <div class="form-text">Specify the name of the place where you need someone to stand in line</div>
                </div>
                
                <div class="mb-4">
                    <label for="orderAddress" class="form-label">Address</label>
                    <div class="input-group">
                        <span class="input-group-text"><i class="fas fa-map-marker-alt"></i></span>
                        <input type="text" class="form-control" id="orderAddress" name="orderAddress" 
                               placeholder="Enter the full address" required>
                    </div>
                    <div class="form-text">Provide the complete address of the location</div>
                </div>
                
                <div class="row mb-4">
                    <div class="col-md-6">
                        <label for="orderDate" class="form-label">Date and Time</label>
                        <div class="input-group">
                            <span class="input-group-text"><i class="fas fa-calendar-alt"></i></span>
                            <input type="datetime-local" class="form-control" id="orderDate" name="orderDate" required>
                        </div>
                        <div class="form-text">When do you need someone to be in line</div>
                    </div>
                    
                    <div class="col-md-6">
                        <label for="orderTimeGap" class="form-label">Estimated Time Needed</label>
                        <div class="input-group">
                            <span class="input-group-text"><i class="fas fa-clock"></i></span>
                            <input type="time" class="form-control" id="orderTimeGap" name="orderTimeGap" required>
                        </div>
                        <div class="form-text">How long do you expect to need someone in line</div>
                    </div>
                </div>
                
                <div class="d-grid gap-2">
                    <button type="submit" class="btn btn-success btn-submit">
                        <i class="fas fa-plus-circle me-2"></i>Create Order
                    </button>
                </div>
            </form>
        </div>
    </div>

    <footer class="bg-light py-4 mt-auto">
        <div class="container text-center">
            <p class="mb-0">© 2024 OrderQ. All rights reserved.</p>
        </div>
    </footer>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const createOrderForm = document.getElementById('createOrderForm');
            const errorAlert = document.getElementById('errorAlert');
            const successAlert = document.getElementById('successAlert');
            
            // Check if token exists in localStorage
            const token = localStorage.getItem('token');
            if (!token) {
                // If no token, redirect to login
                window.location.href = '/auth/login';
                return;
            }
            
            // Helper function to get auth headers
            function getAuthHeaders() {
                return {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                };
            }

            createOrderForm.addEventListener('submit', function(e) {
                e.preventDefault();
                
                // Hide any existing alerts
                errorAlert.style.display = 'none';
                successAlert.style.display = 'none';
                
                // Get form data
                const orderLocation = document.getElementById('orderLocation').value;
                const orderAddress = document.getElementById('orderAddress').value;
                const orderDate = document.getElementById('orderDate').value;
                const orderTimeGap = document.getElementById('orderTimeGap').value;
                
                // Create request payload
                const payload = {
                    order_location: orderLocation,
                    order_address: orderAddress,
                    order_date: new Date(orderDate).toISOString(),
                    order_time_gap: createTimeGapTimestamp(orderTimeGap)
                };
                
                // Send API request with auth header
                fetch('/api/orders/create', {
                    method: 'POST',
                    headers: getAuthHeaders(),
                    body: JSON.stringify(payload)
                })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        errorAlert.textContent = data.error;
                        errorAlert.style.display = 'block';
                    } else {
                        successAlert.textContent = 'Order created successfully!';
                        successAlert.style.display = 'block';
                        createOrderForm.reset();
                        
                        // Optionally reload orders list
                        setTimeout(() => {
                            loadOrders();
                        }, 1000);
                    }
                })
                .catch(error => {
                    errorAlert.textContent = 'Failed to create order. Please try again.';
                    errorAlert.style.display = 'block';
                    console.error('Error:', error);
                });
            });
            
            // Helper function to create time gap timestamp
            function createTimeGapTimestamp(timeString) {
                // Parse the time string (HH:MM) and create a timestamp
                const [hours, minutes] = timeString.split(':').map(Number);
                const date = new Date();
                date.setHours(hours, minutes, 0, 0);
                return date.toISOString();
            }
            
            // Load existing orders
            function loadOrders() {
                fetch('/api/orders', {
                    headers: getAuthHeaders()
                })
                .then(response => response.json())
                .then(data => {
                    const ordersContainer = document.getElementById('ordersContainer');
                    ordersContainer.innerHTML = '';
                    
                    if (data.orders && data.orders.length > 0) {
                        data.orders.forEach(order => {
                            const orderDate = new Date(order.order_date);
                            const orderCard = document.createElement('div');
                            orderCard.className = 'col-md-6 col-lg-4 mb-4';
                            orderCard.innerHTML = `
                                <div class="card order-card">
                                    <div class="card-body">
                                        <h5 class="card-title">${order.order_location}</h5>
                                        <p class="card-text">
                                            <i class="fas fa-map-marker-alt me-2 text-success"></i>${order.order_address}<br>
                                            <i class="fas fa-calendar-alt me-2 text-success"></i>${orderDate.toLocaleString()}<br>
                                            <span class="badge ${getStatusBadgeClass(order.order_status)}">${order.order_status}</span>
                                        </p>
                                        <div class="d-flex justify-content-between mt-3">
                                            <button class="btn btn-sm btn-outline-danger" onclick="cancelOrder('${order.order_id}')">
                                                <i class="fas fa-times me-1"></i>Cancel
                                            </button>
                                            ${order.order_status === "active" ? 
                                                `<button class="btn btn-sm btn-outline-success" onclick="finishOrder('${order.order_id}')">
                                                    <i class="fas fa-check me-1"></i>Mark as Complete
                                                </button>` : ''}
                                        </div>
                                    </div>
                                </div>
                            `;
                            ordersContainer.appendChild(orderCard);
                        });
                    } else {
                        ordersContainer.innerHTML = '<div class="col-12"><p class="text-center">No orders found. Create your first order above!</p></div>';
                    }
                })
                .catch(error => {
                    console.error('Error loading orders:', error);
                });
            }
            
            // Helper function to get badge class based on status
            function getStatusBadgeClass(status) {
                switch (status) {
                    case 'pending':
                        return 'bg-warning text-dark';
                    case 'active':
                        return 'bg-primary';
                    case 'finished':
                        return 'bg-success';
                    case 'cancelled':
                        return 'bg-danger';
                    default:
                        return 'bg-secondary';
                }
            }
            
            // Global functions for order actions
            window.cancelOrder = function(orderId) {
                if (confirm('Are you sure you want to cancel this order?')) {
                    fetch(`/api/orders/${orderId}/cancel`, {
                        method: 'POST',
                        headers: getAuthHeaders()
                    })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            loadOrders();
                        } else {
                            alert('Failed to cancel order.');
                        }
                    })
                    .catch(error => {
                        console.error('Error cancelling order:', error);
                    });
                }
            };
            
            window.finishOrder = function(orderId) {
                if (confirm('Are you sure this order is complete?')) {
                    fetch(`/api/orders/${orderId}/finish`, {
                        method: 'POST',
                        headers: getAuthHeaders()
                    })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            loadOrders();
                        } else {
                            alert('Failed to mark order as complete.');
                        }
                    })
                    .catch(error => {
                        console.error('Error completing order:', error);
                    });
                }
            };
            
            // Load orders on page load
            loadOrders();
        });
    </script>
</body>
</html> 