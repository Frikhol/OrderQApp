document.addEventListener('DOMContentLoaded', function() {
    // WebSocket connection
    const userId = document.getElementById('userId').value;
    if (userId) {
        // Create WebSocket connection
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws?user_id=${userId}`;
        const ws = new WebSocket(wsUrl);

        // Connection opened
        ws.addEventListener('open', function (event) {
            console.log('WebSocket connection established');
        });

        // Listen for messages
        ws.addEventListener('message', function (event) {
            console.log('Message from server:', event.data);
        });

        // // Connection closed
        // ws.addEventListener('close', function (event) {
        //     console.log('WebSocket connection closed');
        //     // Attempt to reconnect after 5 seconds
        //     setTimeout(() => {
        //         window.location.reload();
        //     }, 5000);
        // });

        // Connection error
        ws.addEventListener('error', function (event) {
            console.error('WebSocket error:', event);
        });
    }

    // Get references to buttons and containers
    const createFormBtn = document.getElementById('showCreateForm');
    const listOrdersBtn = document.getElementById('showOrdersList');
    const createFormContainer = document.getElementById('createOrderFormContainer');
    const ordersListContainer = document.getElementById('ordersListContainer');
    const ordersContainer = document.getElementById('ordersContainer');
    const loadingSpinner = document.getElementById('loadingSpinner');
    const errorAlert = document.getElementById('errorAlert');
    const successAlert = document.getElementById('successAlert');
    const createOrderForm = document.getElementById('createOrderForm');
    
    // Get time input references
    const hoursInput = document.getElementById('hours');
    const minutesInput = document.getElementById('minutes');
    const secondsInput = document.getElementById('seconds');
    const timeGapHidden = document.getElementById('orderTimeGap');
    
    // Get order details section elements
    const orderDetailsSection = document.getElementById('orderDetailsSection');
    const backToListBtn = document.getElementById('backToList');
    const cancelOrderBtn = document.getElementById('cancelOrderBtn');
    const finishOrderBtn = document.getElementById('finishOrderBtn');
    let currentOrderId = null;
    
    // Add event listeners to time inputs to update the hidden field
    [hoursInput, minutesInput, secondsInput].forEach(input => {
        input.addEventListener('change', updateTimeGapValue);
        input.addEventListener('input', updateTimeGapValue);
    });
    
    // Function to update the hidden time gap value
    function updateTimeGapValue() {
        // Ensure valid values
        const hours = Math.min(Math.max(parseInt(hoursInput.value) || 0, 0), 24);
        const minutes = Math.min(Math.max(parseInt(minutesInput.value) || 0, 0), 59);
        const seconds = Math.min(Math.max(parseInt(secondsInput.value) || 0, 0), 59);
        
        // Update the input values (in case they were adjusted)
        hoursInput.value = hours;
        minutesInput.value = minutes;
        secondsInput.value = seconds;
        
        // Format the time string as HH:MM:SS
        const formattedHours = hours.toString().padStart(2, '0');
        const formattedMinutes = minutes.toString().padStart(2, '0');
        const formattedSeconds = seconds.toString().padStart(2, '0');
        
        // Set the hidden field value
        timeGapHidden.value = `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
    }
    
    // Initialize the hidden field with the default values
    updateTimeGapValue();

    // Button event listeners
    createFormBtn.addEventListener('click', function() {
        createFormContainer.style.display = 'block';
        ordersListContainer.style.display = 'none';
        orderDetailsSection.style.display = 'none';
    });

    listOrdersBtn.addEventListener('click', function() {
        createFormContainer.style.display = 'none';
        ordersListContainer.style.display = 'block';
        orderDetailsSection.style.display = 'none';
        loadOrders();
    });

    // Format date safely
    function formatDate(dateString) {
        if (!dateString) return 'Not specified';
        
        try {
            // For timestamp objects with seconds and nanos fields (protobuf timestamp format)
            if (dateString.seconds) {
                const seconds = parseInt(dateString.seconds);
                const milliseconds = seconds * 1000;
                const date = new Date(milliseconds);
                
                if (isNaN(date.getTime())) {
                    return 'Date unavailable';
                }
                
                return date.toLocaleString();
            }
            
            // For ISO string dates
            const date = new Date(dateString);
            if (isNaN(date.getTime())) {
                return 'Date unavailable';
            }
            
            return date.toLocaleString();
        } catch (error) {
            console.error('Error formatting date:', error, dateString);
            return 'Date unavailable';
        }
    }
    
    // Function to get the status display text
    function getStatusText(status) {
        if (!status) return 'Unknown';
        
        switch(status.toLowerCase()) {
            case 'pending':
                return 'Waiting for Executor';
            case 'matching':
                return 'Executor Confirmation';
            case 'signed':
                return 'In Progress';
            case 'cancelled':
                return 'Cancelled';
            case 'finished':
                return 'Completed';
            default:
                return status;
        }
    }

    // Load orders function
    function loadOrders() {
        // Show loading spinner
        loadingSpinner.style.display = 'flex';
        ordersContainer.innerHTML = '';

        fetch('/api/orders/list', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load orders');
            }
            return response.json();
        })
        .then(data => {
            loadingSpinner.style.display = 'none';
            
            if (data.error) {
                ordersContainer.innerHTML = `<div class="col-12 text-center"><p class="text-danger">${data.error}</p></div>`;
                return;
            }

            if (!data.orders || data.orders.length === 0) {
                ordersContainer.innerHTML = `
                    <div class="col-12 text-center">
                        <div class="p-5 bg-light rounded-3">
                            <h4 class="text-muted"><i class="fas fa-info-circle me-2"></i>No Orders Found</h4>
                            <p>You don't have any orders yet. Create one!</p>
                        </div>
                    </div>`;
                return;
            }

            // Render each order
            data.orders.forEach(order => {
                const formattedDate = formatDate(order.order_date);
                const statusText = getStatusText(order.order_status);
                
                // Convert time gap (duration) to hours, minutes and seconds
                let timeGap = '';
                if (order.order_time_gap) {
                    const seconds = parseInt(order.order_time_gap.seconds);
                    const hours = Math.floor(seconds / 3600);
                    const minutes = Math.floor((seconds % 3600) / 60);
                    const remainingSeconds = seconds % 60;
                    
                    if (hours > 0) {
                        timeGap += `${hours} hour${hours > 1 ? 's' : ''}`;
                    }
                    if (minutes > 0) {
                        timeGap += `${hours > 0 ? ' ' : ''}${minutes} minute${minutes > 1 ? 's' : ''}`;
                    }
                    if (remainingSeconds > 0) {
                        timeGap += `${(hours > 0 || minutes > 0) ? ' ' : ''}${remainingSeconds} second${remainingSeconds > 1 ? 's' : ''}`;
                    }
                    
                    if (timeGap === '') {
                        timeGap = '0 seconds';
                    }
                } else {
                    timeGap = 'Not specified';
                }

                const orderCard = document.createElement('div');
                orderCard.className = 'col-md-6 col-lg-4 mb-4';
                orderCard.innerHTML = `
                    <div class="card h-100">
                        <div class="card-body">
                            <h5 class="card-title">${order.order_location || 'No location'}</h5>
                            <h6 class="card-subtitle mb-3 text-muted">${order.order_address || 'No address'}</h6>
                            <div class="d-flex align-items-center mb-2">
                                <i class="fas fa-calendar-alt text-success me-2"></i>
                                <span>${formattedDate}</span>
                            </div>
                            <div class="d-flex align-items-center mb-2">
                                <i class="fas fa-clock text-success me-2"></i>
                                <span>${timeGap}</span>
                            </div>
                            <div class="d-flex align-items-center">
                                <i class="fas fa-tag text-success me-2"></i>
                                <span class="badge bg-${getStatusBadgeColor(order.order_status)}">${statusText}</span>
                            </div>
                        </div>
                        <div class="card-footer bg-white">
                            <div class="d-grid">
                                <button class="btn btn-outline-success view-order-btn" data-order-id="${order.order_id}">
                                    <i class="fas fa-eye me-1"></i> View Details
                                </button>
                            </div>
                        </div>
                    </div>
                `;
                ordersContainer.appendChild(orderCard);
            });

            // Add event listeners to view buttons
            document.querySelectorAll('.view-order-btn').forEach(btn => {
                btn.addEventListener('click', function() {
                    const orderId = this.getAttribute('data-order-id');
                    viewOrderDetails(orderId);
                });
            });
        })
        .catch(error => {
            loadingSpinner.style.display = 'none';
            ordersContainer.innerHTML = `
                <div class="col-12 text-center">
                    <div class="p-5 bg-light rounded-3">
                        <h4 class="text-danger"><i class="fas fa-exclamation-circle me-2"></i>Error</h4>
                        <p>${error.message}</p>
                    </div>
                </div>`;
        });
    }

    // Function to view order details
    function viewOrderDetails(orderId) {
        currentOrderId = orderId;
        fetch(`/api/orders/${orderId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load order details');
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                alert('Error: ' + data.error);
                return;
            }
            
            const order = data.order;
            if (!order) {
                alert('Error: Order not found');
                return;
            }

            const formattedDate = formatDate(order.order_date);
            const statusText = getStatusText(order.order_status);
            
            // Convert time gap (duration) to hours, minutes and seconds
            let timeGap = '';
            if (order.order_time_gap) {
                const seconds = Math.floor(order.order_time_gap / 1000000000); // Convert nanoseconds to seconds
                const hours = Math.floor(seconds / 3600);
                const minutes = Math.floor((seconds % 3600) / 60);
                const remainingSeconds = seconds % 60;
                
                if (hours > 0) {
                    timeGap += `${hours} hour${hours > 1 ? 's' : ''}`;
                }
                if (minutes > 0) {
                    timeGap += `${hours > 0 ? ' ' : ''}${minutes} minute${minutes > 1 ? 's' : ''}`;
                }
                if (remainingSeconds > 0) {
                    timeGap += `${(hours > 0 || minutes > 0) ? ' ' : ''}${remainingSeconds} second${remainingSeconds > 1 ? 's' : ''}`;
                }
                
                if (timeGap === '') {
                    timeGap = '0 seconds';
                }
            } else {
                timeGap = 'Not specified';
            }

            // Create order details HTML
            const orderDetailsHTML = `
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
                                <h5 class="card-title">${order.order_location || 'Not specified'}</h5>
                                <h6 class="card-subtitle mb-3 text-muted">${order.order_address || 'Not specified'}</h6>
                                <div class="d-flex align-items-center mb-2">
                                    <i class="fas fa-calendar-alt text-success me-2"></i>
                                    <span>${formattedDate}</span>
                                </div>
                                <div class="d-flex align-items-center mb-2">
                                    <i class="fas fa-clock text-success me-2"></i>
                                    <span>${timeGap}</span>
                                </div>
                                <div class="d-flex align-items-center">
                                    <i class="fas fa-tag text-success me-2"></i>
                                    <span class="badge bg-${getStatusBadgeColor(order.order_status)}">${statusText}</span>
                                </div>
                            </div>
                            <div class="card-footer bg-white">
                                <div class="d-grid gap-2">
                                    ${order.order_status !== 'cancelled' ? `
                                        <button type="button" class="btn btn-outline-danger" id="cancelOrderBtn">
                                            <i class="fas fa-times-circle me-2"></i>Cancel Order
                                        </button>
                                    ` : ''}
                                    ${order.order_status === 'active' ? `
                                        <button type="button" class="btn btn-outline-primary" id="finishOrderBtn">
                                            <i class="fas fa-check-circle me-2"></i>Finish Order
                                        </button>
                                    ` : ''}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;

            // Update the order details section
            orderDetailsSection.innerHTML = orderDetailsHTML;
            
            // Add event listeners for buttons
            document.getElementById('backToList').addEventListener('click', function() {
                orderDetailsSection.style.display = 'none';
                ordersListContainer.style.display = 'block';
            });

            const cancelOrderBtn = document.getElementById('cancelOrderBtn');
            if (cancelOrderBtn) {
                cancelOrderBtn.addEventListener('click', function() {
                    if (confirm('Are you sure you want to cancel this order?')) {
                        cancelOrder(currentOrderId);
                    }
                });
            }

            const finishOrderBtn = document.getElementById('finishOrderBtn');
            if (finishOrderBtn) {
                finishOrderBtn.addEventListener('click', function() {
                    if (confirm('Are you sure you want to mark this order as finished?')) {
                        finishOrder(currentOrderId);
                    }
                });
            }

            // Show order details section and hide orders list
            ordersListContainer.style.display = 'none';
            orderDetailsSection.style.display = 'block';
        })
        .catch(error => {
            console.error('Error in viewOrderDetails:', error);
            alert('Error: ' + error.message);
        });
    }

    // Function to get badge color based on order status
    function getStatusBadgeColor(status) {
        if (!status) return 'secondary';
        
        switch(status.toLowerCase()) {
            case 'pending':
                return 'warning';
            case 'matching':
                return 'success';
            case 'signed':
                return 'primary';
            case 'cancelled':
            case 'canceled':
                return 'danger';
            case 'finished':
                return 'info';
            default:
                return 'secondary';
        }
    }

    // Function to convert time format (HH:MM:SS) to seconds
    function convertTimeToSeconds(timeStr) {
        if (!timeStr) return 0;
        
        const [hours, minutes, seconds] = timeStr.split(':').map(num => parseInt(num || 0, 10));
        return (hours * 3600) + (minutes * 60) + (seconds || 0);
    }

    // Handle form submission
    createOrderForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Hide previous alerts
        errorAlert.style.display = 'none';
        successAlert.style.display = 'none';
        
        // Update the time gap value one last time before submission
        updateTimeGapValue();
        
        // Get form data
        const formData = new FormData(createOrderForm);
        const timeGapStr = formData.get('orderTimeGap');
        const totalSeconds = convertTimeToSeconds(timeGapStr);
        
        // Check if at least some time is specified
        if (totalSeconds <= 0) {
            errorAlert.textContent = 'Please specify a time greater than zero';
            errorAlert.style.display = 'block';
            return;
        }

        // Validate the date
        const dateStr = formData.get('orderDate');
        const dateObj = new Date(dateStr);
        if (isNaN(dateObj.getTime())) {
            errorAlert.textContent = 'Please enter a valid date and time';
            errorAlert.style.display = 'block';
            return;
        }
        
        const orderData = {
            order_location: formData.get('orderLocation'),
            order_address: formData.get('orderAddress'),
            order_date: dateObj.toISOString(),
            order_time_gap: `${totalSeconds}s`
        };
        
        // Send request to create order
        fetch('/api/orders/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(orderData)
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
                
                // Reset the time inputs to default values
                hoursInput.value = 0;
                minutesInput.value = 30;
                secondsInput.value = 0;
                updateTimeGapValue();
                
                // After a delay, hide the form and show the orders
                setTimeout(() => {
                    createFormContainer.style.display = 'none';
                    ordersListContainer.style.display = 'block';
                    loadOrders();
                }, 2000);
            }
        })
        .catch(error => {
            errorAlert.textContent = error.message;
            errorAlert.style.display = 'block';
        });
    });

    // Function to cancel an order
    function cancelOrder(orderId) {
        fetch(`/api/orders/${orderId}/cancel`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to cancel order');
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                alert('Error: ' + data.error);
                return;
            }
            alert('Order cancelled successfully');
            // Go back to orders list and refresh
            orderDetailsSection.style.display = 'none';
            ordersListContainer.style.display = 'block';
            loadOrders();
        })
        .catch(error => {
            alert('Error: ' + error.message);
        });
    }

    // Function to finish an order
    function finishOrder(orderId) {
        fetch(`/api/orders/${orderId}/complete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to finish order');
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                alert('Error: ' + data.error);
                return;
            }
            alert('Order marked as finished');
            // Go back to orders list and refresh
            orderDetailsSection.style.display = 'none';
            ordersListContainer.style.display = 'block';
            loadOrders();
        })
        .catch(error => {
            alert('Error: ' + error.message);
        });
    }

    // Show order list by default
    listOrdersBtn.click();
});
