<div class="form-container" id="createOrderFormContainer" style="display: none;">
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
                <label class="form-label">Estimated Time Needed</label>
                <div class="time-inputs row">
                    <div class="col-4">
                        <div class="input-group">
                            <input type="number" class="form-control" id="hours" min="0" max="24" value="0" placeholder="Hours">
                            <span class="input-group-text">hrs</span>
                        </div>
                    </div>
                    <div class="col-4">
                        <div class="input-group">
                            <input type="number" class="form-control" id="minutes" min="0" max="59" value="30" placeholder="Minutes">
                            <span class="input-group-text">min</span>
                        </div>
                    </div>
                    <div class="col-4">
                        <div class="input-group">
                            <input type="number" class="form-control" id="seconds" min="0" max="59" value="0" placeholder="Seconds">
                            <span class="input-group-text">sec</span>
                        </div>
                    </div>
                    <input type="hidden" id="orderTimeGap" name="orderTimeGap" value="00:30:00">
                </div>
                <div class="form-text mt-2">Set how much time you need to be notified before the order is due</div>
            </div>
        </div>
        
        <div class="d-grid gap-2">
            <button type="submit" class="btn btn-success btn-lg">
                <i class="fas fa-plus-circle me-2"></i>Create Order
            </button>
        </div>
    </form>
</div> 