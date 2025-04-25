{{template "base.tpl" .}}

{{define "content"}}
<div class="container mt-4">
    <div class="row">
        <div class="col-md-12">
            <h2 class="mb-4">Client Dashboard</h2>

            <!-- Current Order Section -->
            <div id="currentOrderSection" class="mb-4">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="mb-0">Current Order</h5>
                        <span id="orderStatus" class="badge"></span>
                    </div>
                    <div class="card-body">
                        <div id="currentOrderContent">
                            <!-- Order details will be loaded here -->
                        </div>
                    </div>
                </div>
            </div>

            <!-- New Order Form Section -->
            <div id="newOrderSection" class="card" style="display: none;">
                <div class="card-header">
                    <h5 class="mb-0">Create New Order</h5>
                </div>
                <div class="card-body">
                    <form id="newOrderForm">
                        <div class="mb-3">
                            <label for="location" class="form-label">Location</label>
                            <input type="text" class="form-control" id="location" name="location" required>
                        </div>
                        <div class="mb-3">
                            <label for="start_time" class="form-label">Start Time</label>
                            <input type="datetime-local" class="form-control" id="start_time" name="start_time" required>
                        </div>
                        <div class="mb-3">
                            <label for="time_buffer" class="form-label">Time Buffer (minutes)</label>
                            <input type="number" class="form-control" id="time_buffer" name="time_buffer" min="30" step="30" value="30" required>
                            <small class="text-muted">Minimum 30 minutes, increments of 30 minutes</small>
                        </div>
                        <button type="submit" class="btn btn-primary">Create Order</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
// Error logging function
function logError(error, context) {
    const errorData = {
        timestamp: new Date().toISOString(),
        context: context,
        error: error.message || error,
        stack: error.stack,
        userAgent: navigator.userAgent,
        url: window.location.href
    };

    // Log to console
    console.error('Error in context:', context, errorData);

    // Send to server for logging
    fetch('/api/logs/error', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify(errorData)
    }).catch(logError => {
        console.error('Failed to log error to server:', logError);
    });
}

document.addEventListener('DOMContentLoaded', function() {
    // Load current order
    fetch('/api/orders/current', {
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to load current order');
        }
        return response.json();
    })
    .then(data => {
        console.log('Current order data:', data);
        if (data && data.Id) {
            // Show current order
            displayCurrentOrder(data);
        } else {
            // Show new order form
            showNewOrderForm();
        }
    })
    .catch(error => {
        console.error('Error loading current order:', error);
        logError(error, 'Loading current order');
        showNewOrderForm();
    });
});

function displayCurrentOrder(order) {
    const currentOrderSection = document.getElementById('currentOrderSection');
    const newOrderSection = document.getElementById('newOrderSection');
    const currentOrderContent = document.getElementById('currentOrderContent');
    const orderStatus = document.getElementById('orderStatus');

    // Update status badge
    orderStatus.className = 'badge ' + getStatusBadgeClass(order.Status);
    orderStatus.textContent = order.Status;

    // Format order details
    currentOrderContent.innerHTML = `
        <div class="order-details">
            <div class="mb-3">
                <strong>Order ID:</strong> ${order.Id}
            </div>
            <div class="mb-3">
                <strong>Location:</strong> ${order.Location}
            </div>
            <div class="mb-3">
                <strong>Start Time:</strong> ${new Date(order.StartTime).toLocaleString()}
            </div>
            <div class="mb-3">
                <strong>Price:</strong> $${order.Price ? order.Price.toFixed(2) : '0.00'}
            </div>
            <div class="mb-3">
                <strong>Queue Position:</strong> ${order.QueuePosition || '-'}
            </div>
            <div class="mb-3">
                <strong>Created At:</strong> ${new Date(order.CreatedAt).toLocaleString()}
            </div>
            ${order.Status === 'pending' ? `
                <button class="btn btn-danger" onclick="cancelOrder(${order.Id})">Cancel Order</button>
            ` : ''}
        </div>
    `;

    currentOrderSection.style.display = 'block';
    newOrderSection.style.display = 'none';
}

function showNewOrderForm() {
    const currentOrderSection = document.getElementById('currentOrderSection');
    const newOrderSection = document.getElementById('newOrderSection');

    currentOrderSection.style.display = 'none';
    newOrderSection.style.display = 'block';
}

function getStatusBadgeClass(status) {
    switch(status) {
        case 'pending':
            return 'bg-warning';
        case 'in_progress':
            return 'bg-info';
        case 'completed':
            return 'bg-success';
        case 'cancelled':
            return 'bg-danger';
        default:
            return 'bg-secondary';
    }
}

function cancelOrder(orderId) {
    if (confirm('Are you sure you want to cancel this order?')) {
        fetch(`/api/orders/${orderId}/cancel`, {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            }
        })
        .then(response => {
            if (response.ok) {
                showNewOrderForm();
            } else {
                throw new Error('Failed to cancel order');
            }
        })
        .catch(error => {
            logError(error, 'Cancelling order');
            alert('An error occurred while cancelling the order');
        });
    }
}

// Handle new order form submission
document.getElementById('newOrderForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const location = document.getElementById('location').value;
    const startTime = document.getElementById('start_time').value;
    const timeBuffer = parseInt(document.getElementById('time_buffer').value);

    if (!location || !startTime || !timeBuffer) {
        const error = new Error('Missing required fields');
        logError(error, 'Form validation');
        alert('Please fill in all required fields');
        return;
    }

    // Format the start time to ISO string
    const startTimeDate = new Date(startTime);
    const startTimeISO = startTimeDate.toISOString();

    const formData = {
        location: location,
        start_time: startTimeISO,
        time_buffer: timeBuffer
    };

    console.log('Sending order data:', formData);

    fetch('/api/orders', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        console.log('Response status:', response.status);
        if (!response.ok) {
            return response.json().then(data => {
                console.error('Error response:', data);
                throw new Error(data.error || 'Failed to create order');
            });
        }
        return response.json();
    })
    .then(data => {
        console.log('Order created:', data);
        if (data && data.Id) {  // Changed from data.id to data.Id to match Go struct field
            displayCurrentOrder(data);
        } else {
            console.error('Invalid response data:', data);
            throw new Error('Invalid response from server: Missing order ID');
        }
    })
    .catch(error => {
        console.error('Error creating order:', error);
        logError(error, 'Creating order');
        alert(error.message || 'An error occurred while creating the order');
    });
});

// Add error handling for unhandled promise rejections
window.addEventListener('unhandledrejection', function(event) {
    logError(event.reason, 'Unhandled Promise Rejection');
});

// Add error handling for global errors
window.addEventListener('error', function(event) {
    logError(event.error || event.message, 'Global Error');
});
</script>
{{end}} 