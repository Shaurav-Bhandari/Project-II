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
              <td><strong>NRS {{ order.total.toFixed(2) }}</strong></td>
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
              <span class="item-price">NRS {{ item.total_price.toFixed(2) }}</span>
            </div>
          </div>
          
          <div class="order-totals">
            <div class="total-row">
              <span>Subtotal</span>
              <span>NRS {{ selectedOrder.subtotal.toFixed(2) }}</span>
            </div>
            <div class="total-row">
              <span>Tax (10%)</span>
              <span>NRS {{ selectedOrder.tax.toFixed(2) }}</span>
            </div>
            <div class="total-row total-final">
              <span>Total</span>
              <span>NRS {{ selectedOrder.total.toFixed(2) }}</span>
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
import { ordersAPI } from '../api'

const activeFilter = ref('all')
const showModal = ref(false)
const selectedOrder = ref(null)
const orders = ref([])
const loading = ref(false)

const statusFilters = [
  { value: 'all', label: 'All Orders' },
  { value: 'pending', label: 'Pending' },
  { value: 'confirmed', label: 'Confirmed' },
  { value: 'preparing', label: 'Preparing' },
  { value: 'ready', label: 'Ready' },
  { value: 'served', label: 'Served' },
  { value: 'completed', label: 'Completed' }
]

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await ordersAPI.getAll()
    orders.value = (res.data || []).map(o => ({ ...o, total: o.total || 0, subtotal: o.subtotal || 0, tax: o.tax || 0 }))
  } catch (e) {
    console.error('Failed to fetch orders', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchOrders)

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

const viewOrder = async (order) => {
  try {
    const res = await ordersAPI.getById(order.id)
    selectedOrder.value = { ...res.data, total: res.data.total || 0, subtotal: res.data.subtotal || 0, tax: res.data.tax || 0 }
    showModal.value = true
  } catch (e) {
    console.error('Failed to fetch order details', e)
  }
}

const confirmOrder = async (id) => {
  try {
    await ordersAPI.updateStatus(id, 'confirmed')
    const order = orders.value.find(o => o.id === id)
    if (order) order.status = 'confirmed'
  } catch (e) {
    console.error('Failed to confirm order', e)
  }
}

const updateStatus = async (id, status) => {
  try {
    await ordersAPI.updateStatus(id, status)
    const order = orders.value.find(o => o.id === id)
    if (order) order.status = status
  } catch (e) {
    console.error('Failed to update order status', e)
  }
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

@media (max-width: 768px) {
  .order-info-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  .page-actions {
    flex-direction: column;
    align-items: stretch;
  }
  .filters {
    display: flex;
    overflow-x: auto;
    padding-bottom: 8px;
  }
  .filters::-webkit-scrollbar { height: 4px; }
  .filters::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 2px; }
  .filter-btn { white-space: nowrap; }
}
</style>
