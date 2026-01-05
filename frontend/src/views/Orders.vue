<template>
  <div class="orders-page">
    <!-- Filters and Actions -->
    <div class="page-actions">
      <div class="filters">
        <button 
          v-for="status in statusFilters" 
          :key="status.value"
          :class="['filter-btn', { active: activeFilter === status.value }]"
          @click="activeFilter = status.value"
        >
          {{ status.label }}
        </button>
      </div>
      <router-link to="/orders/new" class="btn btn-primary">
        ➕ New Order
      </router-link>
    </div>
    
    <!-- Orders Table -->
    <div class="card">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>Order #</th>
              <th>Table</th>
              <th>Customer</th>
              <th>Items</th>
              <th>Total</th>
              <th>Status</th>
              <th>Time</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in filteredOrders" :key="order.id">
              <td><strong>#{{ order.order_number }}</strong></td>
              <td>{{ order.table_number || 'Takeaway' }}</td>
              <td>{{ order.customer_name || '-' }}</td>
              <td>{{ order.items?.length || 0 }} items</td>
              <td><strong>${{ order.total.toFixed(2) }}</strong></td>
              <td>
                <span :class="['badge', getStatusBadge(order.status)]">
                  {{ order.status }}
                </span>
              </td>
              <td>{{ formatTime(order.created_at) }}</td>
              <td>
                <div class="action-btns">
                  <button 
                    v-if="order.status === 'pending'"
                    class="btn btn-sm btn-success"
                    @click="confirmOrder(order.id)"
                  >
                    Confirm
                  </button>
                  <button 
                    v-if="order.status === 'ready'"
                    class="btn btn-sm btn-primary"
                    @click="updateStatus(order.id, 'served')"
                  >
                    Serve
                  </button>
                  <button 
                    v-if="order.status === 'served'"
                    class="btn btn-sm btn-success"
                    @click="updateStatus(order.id, 'completed')"
                  >
                    Complete
                  </button>
                  <button 
                    class="btn btn-sm btn-secondary"
                    @click="viewOrder(order)"
                  >
                    View
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredOrders.length === 0">
              <td colspan="8">
                <div class="empty-state">
                  <div class="empty-state-icon">📝</div>
                  <div class="empty-state-title">No orders found</div>
                  <div class="empty-state-text">Create a new order to get started</div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- Order Detail Modal -->
    <div :class="['modal-overlay', { active: showModal }]" @click.self="showModal = false">
      <div class="modal" v-if="selectedOrder">
        <div class="modal-header">
          <h3 class="modal-title">Order #{{ selectedOrder.order_number }}</h3>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="order-info-grid">
            <div><strong>Table:</strong> {{ selectedOrder.table_number || 'Takeaway' }}</div>
            <div><strong>Customer:</strong> {{ selectedOrder.customer_name || '-' }}</div>
            <div><strong>Status:</strong> 
              <span :class="['badge', getStatusBadge(selectedOrder.status)]">{{ selectedOrder.status }}</span>
            </div>
            <div><strong>Time:</strong> {{ formatTime(selectedOrder.created_at) }}</div>
          </div>
          
          <h4 style="margin: 20px 0 12px;">Order Items</h4>
          <div class="order-items-list">
            <div v-for="item in selectedOrder.items" :key="item.id" class="order-item-row">
              <span class="item-qty">{{ item.quantity }}x</span>
              <span class="item-name">{{ item.menu_item_name }}</span>
              <span class="item-price">${{ item.total_price.toFixed(2) }}</span>
            </div>
          </div>
          
          <div class="order-totals">
            <div class="total-row">
              <span>Subtotal</span>
              <span>${{ selectedOrder.subtotal.toFixed(2) }}</span>
            </div>
            <div class="total-row">
              <span>Tax (10%)</span>
              <span>${{ selectedOrder.tax.toFixed(2) }}</span>
            </div>
            <div class="total-row total-final">
              <span>Total</span>
              <span>${{ selectedOrder.total.toFixed(2) }}</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">Close</button>
          <router-link 
            v-if="selectedOrder.status === 'completed'" 
            :to="`/billing?order=${selectedOrder.id}`"
            class="btn btn-primary"
          >
            View Bill
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrdersStore } from '../stores/orders'

const ordersStore = useOrdersStore()

const activeFilter = ref('all')
const showModal = ref(false)
const selectedOrder = ref(null)

const statusFilters = [
  { value: 'all', label: 'All Orders' },
  { value: 'pending', label: 'Pending' },
  { value: 'confirmed', label: 'Confirmed' },
  { value: 'preparing', label: 'Preparing' },
  { value: 'ready', label: 'Ready' },
  { value: 'served', label: 'Served' },
  { value: 'completed', label: 'Completed' }
]

// Demo orders data
const orders = ref([
  { id: 1, order_number: 1047, table_number: 'T3', customer_name: 'John Doe', status: 'preparing', total: 89.96, subtotal: 81.78, tax: 8.18, created_at: new Date(), items: [
    { id: 1, menu_item_name: 'Grilled Salmon', quantity: 2, total_price: 49.98 },
    { id: 2, menu_item_name: 'Spring Rolls', quantity: 1, total_price: 8.99 },
    { id: 3, menu_item_name: 'Fresh Lemonade', quantity: 2, total_price: 9.98 }
  ]},
  { id: 2, order_number: 1046, table_number: 'T7', customer_name: 'Jane Smith', status: 'ready', total: 156.94, subtotal: 142.67, tax: 14.27, created_at: new Date(Date.now() - 30*60000), items: []},
  { id: 3, order_number: 1045, table_number: null, customer_name: 'Mike Johnson', status: 'completed', total: 34.98, subtotal: 31.80, tax: 3.18, created_at: new Date(Date.now() - 60*60000), items: []},
  { id: 4, order_number: 1044, table_number: 'T1', customer_name: null, status: 'served', total: 52.97, subtotal: 48.15, tax: 4.82, created_at: new Date(Date.now() - 90*60000), items: []},
  { id: 5, order_number: 1043, table_number: 'T5', customer_name: 'Sarah Wilson', status: 'pending', total: 124.95, subtotal: 113.59, tax: 11.36, created_at: new Date(Date.now() - 5*60000), items: []}
])

const filteredOrders = computed(() => {
  if (activeFilter.value === 'all') return orders.value
  return orders.value.filter(o => o.status === activeFilter.value)
})

const getStatusBadge = (status) => {
  const badges = {
    pending: 'badge-warning',
    confirmed: 'badge-info',
    preparing: 'badge-warning',
    ready: 'badge-success',
    served: 'badge-primary',
    completed: 'badge-success',
    cancelled: 'badge-danger'
  }
  return badges[status] || 'badge-primary'
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}

const viewOrder = (order) => {
  selectedOrder.value = order
  showModal.value = true
}

const confirmOrder = async (id) => {
  const order = orders.value.find(o => o.id === id)
  if (order) order.status = 'confirmed'
}

const updateStatus = async (id, status) => {
  const order = orders.value.find(o => o.id === id)
  if (order) order.status = status
}
</script>

<style scoped>
.page-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.action-btns {
  display: flex;
  gap: 8px;
}

.order-info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: var(--border-radius-sm);
}

.order-items-list {
  background: var(--bg-secondary);
  border-radius: var(--border-radius-sm);
  padding: 12px;
}

.order-item-row {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.order-item-row:last-child {
  border-bottom: none;
}

.item-qty {
  font-weight: 600;
  color: var(--accent-primary);
  min-width: 40px;
}

.item-name {
  flex: 1;
}

.item-price {
  font-weight: 500;
}

.order-totals {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.total-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  color: var(--text-secondary);
}

.total-final {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
  margin-top: 8px;
}
</style>
