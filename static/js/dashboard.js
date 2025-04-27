// Dashboard module
const Dashboard = (function() {
    // DOM Elements
    const elements = {
        currentOrderSection: document.getElementById('currentOrderSection'),
        newOrderSection: document.getElementById('newOrderSection'),
        currentOrderContent: document.getElementById('currentOrderContent'),
        orderStatus: document.getElementById('orderStatus'),
        newOrderForm: document.getElementById('newOrderForm')
    };

    // Status badge classes
    const statusClasses = {
        pending: 'bg-warning',
        in_progress: 'bg-info',
        completed: 'bg-success',
        cancelled: 'bg-danger',
        default: 'bg-secondary'
    };

    // Initialize dashboard
    function init() {
        loadCurrentOrder();
        if (elements.newOrderForm) {
            elements.newOrderForm.addEventListener('submit', handleOrderSubmit);
        }
    }

    // Load current order
    function loadCurrentOrder() {
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
                displayCurrentOrder(data);
            } else {
                showNewOrderForm();
            }
        })
        .catch(error => {
            console.error('Error loading current order:', error);
            ErrorLogger.log(error, 'Loading current order');
            showNewOrderForm();
        });
    }

    // Display current order
    function displayCurrentOrder(order) {
        elements.orderStatus.className = 'badge ' + getStatusBadgeClass(order.Status);
        elements.orderStatus.textContent = order.Status;

        elements.currentOrderContent.innerHTML = `
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
                    <button class="btn btn-danger" onclick="Dashboard.cancelOrder(${order.Id})">Cancel Order</button>
                ` : ''}
            </div>
        `;

        elements.currentOrderSection.style.display = 'block';
        elements.newOrderSection.style.display = 'none';
    }

    // Show new order form
    function showNewOrderForm() {
        elements.currentOrderSection.style.display = 'none';
        elements.newOrderSection.style.display = 'block';
    }

    // Get status badge class
    function getStatusBadgeClass(status) {
        return statusClasses[status] || statusClasses.default;
    }

    // Handle order form submission
    function handleOrderSubmit(e) {
        e.preventDefault();
        
        const location = document.getElementById('location').value;
        const startTime = document.getElementById('start_time').value;
        const timeBuffer = parseInt(document.getElementById('time_buffer').value);

        if (!location || !startTime || !timeBuffer) {
            const error = new Error('Missing required fields');
            ErrorLogger.log(error, 'Form validation');
            alert('Please fill in all required fields');
            return;
        }

        const formData = {
            location: location,
            start_time: new Date(startTime).toISOString(),
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
            if (data && data.Id) {
                displayCurrentOrder(data);
            } else {
                console.error('Invalid response data:', data);
                throw new Error('Invalid response from server: Missing order ID');
            }
        })
        .catch(error => {
            console.error('Error creating order:', error);
            ErrorLogger.log(error, 'Creating order');
            alert(error.message || 'An error occurred while creating the order');
        });
    }

    // Cancel order
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
                console.error('Error cancelling order:', error);
                ErrorLogger.log(error, 'Cancelling order');
                alert('An error occurred while cancelling the order');
            });
        }
    }

    // Public API
    return {
        init: init,
        cancelOrder: cancelOrder
    };
})();

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    Dashboard.init();
}); 