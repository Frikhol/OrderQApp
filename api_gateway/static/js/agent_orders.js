document.addEventListener('DOMContentLoaded', function() {
    const searchOrdersBtn = document.getElementById('searchOrders');
    const availableOrdersList = document.getElementById('availableOrdersList');
    const ordersContainer = document.getElementById('ordersContainer');
    const orderDetailsSection = document.getElementById('orderDetailsSection');
    const backToListBtn = document.getElementById('backToList');
    const acceptOrderBtn = document.getElementById('acceptOrderBtn');

    let currentOrderId = null;

    // Search for available orders
    searchOrdersBtn.addEventListener('click', async function() {
        try {
            const response = await fetch('/api/orders/start_search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch available orders');
            }

            const data = await response.json();
            displayOrders(data.orders || []);
            availableOrdersList.style.display = 'block';
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to load available orders. Please try again.');
        }
    });

    // Display orders in the list
    function displayOrders(orders) {
        ordersContainer.innerHTML = '';
        
        if (orders.length === 0) {
            ordersContainer.innerHTML = `
                <div class="col-12 text-center">
                    <div class="p-5 bg-light rounded-3">
                        <h4 class="text-muted"><i class="fas fa-info-circle me-2"></i>No Orders Found</h4>
                        <p>No available orders at the moment.</p>
                    </div>
                </div>`;
            return;
        }

        orders.forEach(order => {
            const orderCard = createOrderCard(order);
            ordersContainer.appendChild(orderCard);
        });
    }

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

    // Function to get status badge color
    function getStatusBadgeColor(status) {
        if (!status) return 'secondary';
        
        switch(status.toLowerCase()) {
            case 'pending':
                return 'warning';
            case 'matching':
                return 'info';
            case 'signed':
                return 'primary';
            case 'cancelled':
                return 'danger';
            case 'finished':
                return 'success';
            default:
                return 'secondary';
        }
    }

    // Format time gap from nanoseconds to minutes
    function formatTimeGap(timeGap) {
        if (!timeGap) return 'Not specified';
        
        try {
            console.log('Time gap value:', timeGap);
            const minutes = Math.floor(Number(timeGap) / 60000000000); // Convert nanoseconds to minutes
            console.log('Calculated minutes:', minutes);
            
            if (isNaN(minutes)) {
                console.error('Invalid time gap value:', timeGap);
                return 'Time unavailable';
            }
            
            return `${minutes} minutes`;
        } catch (error) {
            console.error('Error formatting time gap:', error, timeGap);
            return 'Time unavailable';
        }
    }

    // Create order card element
    function createOrderCard(order) {
        const formattedDate = formatDate(order.order_date);
        const statusText = getStatusText(order.order_status);
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

        // Add event listener to view details button
        orderCard.querySelector('.view-order-btn').addEventListener('click', () => {
            showOrderDetails(order);
        });

        return orderCard;
    }

    // Show order details
    function showOrderDetails(order) {
        currentOrderId = order.order_id;
        
        document.getElementById('orderLocation').textContent = order.order_location || 'No location';
        document.getElementById('orderAddress').textContent = order.order_address || 'No address';
        document.getElementById('orderDate').textContent = formatDate(order.order_date);
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
        document.getElementById('orderTimeGap').textContent = timeGap;
        document.getElementById('orderStatus').textContent = getStatusText(order.order_status);

        availableOrdersList.style.display = 'none';
        orderDetailsSection.style.display = 'block';
    }

    // Back to list button handler
    backToListBtn.addEventListener('click', function() {
        orderDetailsSection.style.display = 'none';
        availableOrdersList.style.display = 'block';
        currentOrderId = null;
    });

    // Accept order button handler
    acceptOrderBtn.addEventListener('click', async function() {
        if (!currentOrderId) return;

        try {
            const response = await fetch(`/api/orders/${currentOrderId}/accept`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error('Failed to accept order');
            }

            alert('Order accepted successfully!');
            // Refresh the orders list
            searchOrdersBtn.click();
            orderDetailsSection.style.display = 'none';
            availableOrdersList.style.display = 'block';
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to accept order. Please try again.');
        }
    });
}); 