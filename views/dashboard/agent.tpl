{{template "base.tpl" .}}

{{define "content"}}
<div class="container mt-4">
    <div class="row">
        <div class="col-md-12">
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Agent Dashboard</h2>
                <div>
                    <a href="/auth/logout" class="btn btn-outline-danger">Logout</a>
                </div>
            </div>

            <div class="row">
                <div class="col-md-4">
                    <div class="card mb-4">
                        <div class="card-body">
                            <h5 class="card-title">Queue Statistics</h5>
                            <div class="d-flex justify-content-between mb-2">
                                <span>Total Orders:</span>
                                <span id="totalOrders">-</span>
                            </div>
                            <div class="d-flex justify-content-between mb-2">
                                <span>Pending Orders:</span>
                                <span id="pendingOrders">-</span>
                            </div>
                            <div class="d-flex justify-content-between">
                                <span>In Progress:</span>
                                <span id="inProgressOrders">-</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="col-md-8">
                    <div class="card">
                        <div class="card-header">
                            <h5 class="mb-0">Order Queue</h5>
                        </div>
                        <div class="card-body">
                            <div class="table-responsive">
                                <table class="table table-hover">
                                    <thead>
                                        <tr>
                                            <th>Order ID</th>
                                            <th>Client</th>
                                            <th>Location</th>
                                            <th>Status</th>
                                            <th>Queue Position</th>
                                            <th>Created At</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody id="ordersTableBody">
                                        <!-- Orders will be loaded here via JavaScript -->
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
document.addEventListener('DOMContentLoaded', function() {
    // Load queue statistics
    fetch('/api/orders/queue/stats', {
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('totalOrders').textContent = data.total || 0;
        document.getElementById('pendingOrders').textContent = data.pending || 0;
        document.getElementById('inProgressOrders').textContent = data.in_progress || 0;
    })
    .catch(error => {
        console.error('Error loading statistics:', error);
    });

    // Load orders
    fetch('/api/orders/queue', {
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    })
    .then(response => response.json())
    .then(data => {
        const tbody = document.getElementById('ordersTableBody');
        tbody.innerHTML = '';

        data.forEach(order => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${order.id}</td>
                <td>${order.client_email}</td>
                <td>${order.location}</td>
                <td>
                    <span class="badge ${getStatusBadgeClass(order.status)}">
                        ${order.status}
                    </span>
                </td>
                <td>${order.queue_position || '-'}</td>
                <td>${new Date(order.created_at).toLocaleString()}</td>
                <td>
                    <button class="btn btn-sm btn-info" onclick="viewOrder(${order.id})">View</button>
                    ${order.status === 'pending' ? `
                        <button class="btn btn-sm btn-success" onclick="startOrder(${order.id})">Start</button>
                    ` : ''}
                    ${order.status === 'in_progress' ? `
                        <button class="btn btn-sm btn-primary" onclick="completeOrder(${order.id})">Complete</button>
                    ` : ''}
                </td>
            `;
            tbody.appendChild(row);
        });
    })
    .catch(error => {
        console.error('Error loading orders:', error);
    });
});

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

function viewOrder(orderId) {
    window.location.href = `/orders/${orderId}`;
}

function startOrder(orderId) {
    fetch(`/api/orders/${orderId}/start`, {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    })
    .then(response => {
        if (response.ok) {
            window.location.reload();
        } else {
            alert('Failed to start order');
        }
    })
    .catch(error => {
        console.error('Error starting order:', error);
        alert('An error occurred while starting the order');
    });
}

function completeOrder(orderId) {
    fetch(`/api/orders/${orderId}/complete`, {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    })
    .then(response => {
        if (response.ok) {
            window.location.reload();
        } else {
            alert('Failed to complete order');
        }
    })
    .catch(error => {
        console.error('Error completing order:', error);
        alert('An error occurred while completing the order');
    });
}
</script>
{{end}} 