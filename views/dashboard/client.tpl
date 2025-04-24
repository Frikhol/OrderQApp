{{template "base.tpl" .}}

{{define "content"}}
<div class="container mt-4">
    <div class="row">
        <div class="col-md-12">
            <div class="d-flex justify-content-between align-items-center mb-4">
                <h2>Client Dashboard</h2>
                <div>
                    <a href="/orders/new" class="btn btn-primary me-2">Create New Order</a>
                </div>
            </div>

            <div class="card">
                <div class="card-header">
                    <h5 class="mb-0">My Orders</h5>
                </div>
                <div class="card-body">
                    <div class="table-responsive">
                        <table class="table table-hover">
                            <thead>
                                <tr>
                                    <th>Order ID</th>
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

<script>
document.addEventListener('DOMContentLoaded', function() {
    // Load user's orders
    fetch('/api/orders/my', {
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
                        <button class="btn btn-sm btn-danger" onclick="cancelOrder(${order.id})">Cancel</button>
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
                window.location.reload();
            } else {
                alert('Failed to cancel order');
            }
        })
        .catch(error => {
            console.error('Error cancelling order:', error);
            alert('An error occurred while cancelling the order');
        });
    }
}
</script>
{{end}} 