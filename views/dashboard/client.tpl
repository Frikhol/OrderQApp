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

<script src="/static/js/error-logger.js"></script>
<script src="/static/js/dashboard.js"></script>
{{end}} 