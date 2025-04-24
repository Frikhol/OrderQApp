document.addEventListener('DOMContentLoaded', function() {
    // Handle login form submission
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(loginForm);
            
            fetch('/auth/login', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.token) {
                    localStorage.setItem('token', data.token);
                    window.location.href = '/dashboard';
                } else {
                    alert('Login failed: ' + data.error);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred during login');
            });
        });
    }

    // Handle registration form submission
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(registerForm);
            
            fetch('/auth/register', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.message) {
                    alert('Registration successful! Please login.');
                    window.location.href = '/auth/login';
                } else {
                    alert('Registration failed: ' + data.error);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred during registration');
            });
        });
    }

    // Handle order form submission (client dashboard)
    const orderForm = document.getElementById('orderForm');
    if (orderForm) {
        orderForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(orderForm);
            
            fetch('/orders', {
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + localStorage.getItem('token')
                },
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.id) {
                    alert('Order created successfully!');
                    orderForm.reset();
                    loadActiveOrders();
                } else {
                    alert('Order creation failed: ' + data.error);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while creating the order');
            });
        });
    }

    // Handle queue position update (agent dashboard)
    const queuePositionForm = document.getElementById('queuePositionForm');
    if (queuePositionForm) {
        queuePositionForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const orderId = document.getElementById('orderId').value;
            const position = document.getElementById('position').value;
            
            fetch(`/orders/${orderId}/position`, {
                method: 'PUT',
                headers: {
                    'Authorization': 'Bearer ' + localStorage.getItem('token'),
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ position: position })
            })
            .then(response => response.json())
            .then(data => {
                if (data.id) {
                    alert('Queue position updated successfully!');
                    $('#queuePositionModal').modal('hide');
                    loadActiveOrders();
                } else {
                    alert('Update failed: ' + data.error);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while updating the position');
            });
        });
    }

    // Load active orders for client
    function loadActiveOrders() {
        const activeOrdersDiv = document.getElementById('activeOrders');
        if (activeOrdersDiv) {
            fetch('/orders', {
                headers: {
                    'Authorization': 'Bearer ' + localStorage.getItem('token')
                }
            })
            .then(response => response.json())
            .then(orders => {
                activeOrdersDiv.innerHTML = orders.map(order => `
                    <div class="card mb-2">
                        <div class="card-body">
                            <h5 class="card-title">Order #${order.id}</h5>
                            <p class="card-text">
                                Location: ${order.location}<br>
                                Status: ${order.status}<br>
                                Queue Position: ${order.queuePosition}<br>
                                Price: $${order.price}
                            </p>
                        </div>
                    </div>
                `).join('');
            })
            .catch(error => {
                console.error('Error:', error);
                activeOrdersDiv.innerHTML = '<p>Error loading orders</p>';
            });
        }
    }

    // Load available orders for agent
    function loadAvailableOrders() {
        const availableOrdersDiv = document.getElementById('availableOrders');
        if (availableOrdersDiv) {
            fetch('/orders/available', {
                headers: {
                    'Authorization': 'Bearer ' + localStorage.getItem('token')
                }
            })
            .then(response => response.json())
            .then(orders => {
                availableOrdersDiv.innerHTML = orders.map(order => `
                    <div class="card mb-2">
                        <div class="card-body">
                            <h5 class="card-title">Order #${order.id}</h5>
                            <p class="card-text">
                                Location: ${order.location}<br>
                                Start Time: ${new Date(order.startTime).toLocaleString()}<br>
                                Price: $${order.price}
                            </p>
                            <button class="btn btn-primary" onclick="acceptOrder(${order.id})">Accept Order</button>
                        </div>
                    </div>
                `).join('');
            })
            .catch(error => {
                console.error('Error:', error);
                availableOrdersDiv.innerHTML = '<p>Error loading available orders</p>';
            });
        }
    }

    // Initialize dashboard data
    if (document.getElementById('activeOrders')) {
        loadActiveOrders();
        setInterval(loadActiveOrders, 30000); // Refresh every 30 seconds
    }

    if (document.getElementById('availableOrders')) {
        loadAvailableOrders();
        setInterval(loadAvailableOrders, 30000); // Refresh every 30 seconds
    }
}); 