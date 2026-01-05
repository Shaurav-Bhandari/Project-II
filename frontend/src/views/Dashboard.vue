<template>
  <div class="dashboard">
    <!-- Quick Actions -->
    <div class="quick-actions">
      <router-link to="/orders/new" class="quick-action-btn">
        <span class="icon" style="background: rgba(99, 102, 241, 0.15); color: var(--accent-primary);">➕</span>
        <span>New Order</span>
      </router-link>
      <router-link to="/kot" class="quick-action-btn">
        <span class="icon" style="background: rgba(245, 158, 11, 0.15); color: var(--accent-warning);">👨‍🍳</span>
        <span>View KOTs</span>
      </router-link>
      <router-link to="/billing" class="quick-action-btn">
        <span class="icon" style="background: rgba(16, 185, 129, 0.15); color: var(--accent-success);">💳</span>
        <span>Process Payment</span>
      </router-link>
      <router-link to="/reports" class="quick-action-btn">
        <span class="icon" style="background: rgba(59, 130, 246, 0.15); color: var(--accent-info);">📈</span>
        <span>View Reports</span>
      </router-link>
    </div>
    
    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card primary">
        <div class="stat-icon">📝</div>
        <div class="stat-value">{{ stats.totalOrders }}</div>
        <div class="stat-label">Total Orders Today</div>
        <div class="stat-change up">
          <span>↑</span>
          <span>12% from yesterday</span>
        </div>
      </div>
      
      <div class="stat-card success">
        <div class="stat-icon">🪑</div>
        <div class="stat-value">{{ stats.activeTables }}</div>
        <div class="stat-label">Active Tables</div>
        <div class="stat-change">
          <span style="color: var(--text-muted);">of {{ stats.totalTables }} tables</span>
        </div>
      </div>
      
      <div class="stat-card warning">
        <div class="stat-icon">💰</div>
        <div class="stat-value">${{ stats.dailyRevenue.toLocaleString() }}</div>
        <div class="stat-label">Daily Revenue</div>
        <div class="stat-change up">
          <span>↑</span>
          <span>8% from yesterday</span>
        </div>
      </div>
      
      <div class="stat-card danger">
        <div class="stat-icon">⏳</div>
        <div class="stat-value">{{ stats.pendingKots }}</div>
        <div class="stat-label">Pending KOTs</div>
        <div class="stat-change">
          <span style="color: var(--text-muted);">{{ stats.inProgressKots }} in progress</span>
        </div>
      </div>
    </div>
    
    <!-- Main Grid -->
    <div class="dashboard-grid">
      <!-- Recent Orders -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Recent Orders</h3>
          <router-link to="/orders" class="btn btn-sm btn-secondary">View All</router-link>
        </div>
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th>Order #</th>
                <th>Table</th>
                <th>Items</th>
                <th>Total</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in recentOrders" :key="order.id">
                <td><strong>#{{ order.order_number }}</strong></td>
                <td>{{ order.table_number || 'Takeaway' }}</td>
                <td>{{ order.item_count }} items</td>
                <td>${{ order.total.toFixed(2) }}</td>
                <td>
                  <span :class="['badge', getStatusBadge(order.status)]">
                    {{ order.status }}
                  </span>
                </td>
              </tr>
              <tr v-if="recentOrders.length === 0">
                <td colspan="5" class="empty-state">
                  <div class="empty-state-icon">📝</div>
                  <div class="empty-state-text">No orders yet today</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- Active KOTs -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Active Kitchen Orders</h3>
          <router-link to="/kot" class="btn btn-sm btn-secondary">View All</router-link>
        </div>
        <div class="kot-list">
          <div v-for="kot in activeKots" :key="kot.id" class="kot-mini-card">
            <div class="kot-mini-header">
              <span class="kot-mini-number">KOT #{{ kot.kot_number }}</span>
              <span :class="['badge', kot.status === 'in_progress' ? 'badge-warning' : 'badge-danger']">
                {{ kot.status === 'in_progress' ? 'Cooking' : 'Pending' }}
              </span>
            </div>
            <div class="kot-mini-info">
              <span>Table {{ kot.table_number }}</span>
              <span>{{ kot.item_count }} items</span>
            </div>
            <div class="kot-mini-time">
              {{ getTimeAgo(kot.created_at) }}
            </div>
          </div>
          <div v-if="activeKots.length === 0" class="empty-state">
            <div class="empty-state-icon">👨‍🍳</div>
            <div class="empty-state-text">No pending kitchen orders</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Table Overview -->
    <div class="card" style="margin-top: 24px;">
      <div class="card-header">
        <h3 class="card-title">Table Overview</h3>
        <router-link to="/tables" class="btn btn-sm btn-secondary">Manage Tables</router-link>
      </div>
      <div class="table-grid">
        <div 
          v-for="table in tables" 
          :key="table.id" 
          :class="['table-card', table.status]"
        >
          <div class="table-card-number">{{ table.table_number }}</div>
          <div class="table-card-capacity">{{ table.capacity }} seats</div>
          <div class="table-card-status">{{ table.status }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { reportsAPI, tablesAPI } from '../api'

const stats = ref({
  totalOrders: 0,
  activeTables: 0,
  totalTables: 10,
  dailyRevenue: 0,
  pendingKots: 0,
  inProgressKots: 0
})

const recentOrders = ref([])
const activeKots = ref([])
const tables = ref([])

// Demo data
onMounted(async () => {
  // In production, these would be API calls
  stats.value = {
    totalOrders: 47,
    activeTables: 6,
    totalTables: 10,
    dailyRevenue: 2847.50,
    pendingKots: 4,
    inProgressKots: 2
  }
  
  recentOrders.value = [
    { id: 1, order_number: 1047, table_number: 'T3', item_count: 4, total: 89.96, status: 'preparing' },
    { id: 2, order_number: 1046, table_number: 'T7', item_count: 6, total: 156.94, status: 'ready' },
    { id: 3, order_number: 1045, table_number: null, item_count: 2, total: 34.98, status: 'completed' },
    { id: 4, order_number: 1044, table_number: 'T1', item_count: 3, total: 52.97, status: 'served' },
    { id: 5, order_number: 1043, table_number: 'T5', item_count: 5, total: 124.95, status: 'completed' }
  ]
  
  activeKots.value = [
    { id: 1, kot_number: 89, table_number: 'T3', item_count: 4, status: 'in_progress', created_at: new Date(Date.now() - 12 * 60000) },
    { id: 2, kot_number: 90, table_number: 'T7', item_count: 3, status: 'pending', created_at: new Date(Date.now() - 5 * 60000) },
    { id: 3, kot_number: 91, table_number: 'T1', item_count: 2, status: 'pending', created_at: new Date(Date.now() - 2 * 60000) }
  ]
  
  tables.value = [
    { id: 1, table_number: 'T1', capacity: 2, status: 'occupied' },
    { id: 2, table_number: 'T2', capacity: 2, status: 'available' },
    { id: 3, table_number: 'T3', capacity: 4, status: 'occupied' },
    { id: 4, table_number: 'T4', capacity: 4, status: 'available' },
    { id: 5, table_number: 'T5', capacity: 6, status: 'reserved' },
    { id: 6, table_number: 'T6', capacity: 6, status: 'available' },
    { id: 7, table_number: 'T7', capacity: 8, status: 'occupied' },
    { id: 8, table_number: 'T8', capacity: 4, status: 'available' },
    { id: 9, table_number: 'T9', capacity: 4, status: 'maintenance' },
    { id: 10, table_number: 'T10', capacity: 2, status: 'available' }
  ]
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

const getTimeAgo = (date) => {
  const minutes = Math.floor((Date.now() - new Date(date).getTime()) / 60000)
  if (minutes < 1) return 'Just now'
  if (minutes === 1) return '1 min ago'
  return `${minutes} mins ago`
}
</script>

<style scoped>
.dashboard-grid {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 24px;
}

.kot-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.kot-mini-card {
  background: var(--bg-secondary);
  border-radius: var(--border-radius-sm);
  padding: 16px;
}

.kot-mini-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.kot-mini-number {
  font-weight: 600;
}

.kot-mini-info {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: var(--text-secondary);
}

.kot-mini-time {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 8px;
}

.table-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.table-card {
  background: var(--bg-secondary);
  border-radius: var(--border-radius-sm);
  padding: 16px;
  text-align: center;
  border: 2px solid transparent;
}

.table-card.available {
  border-color: var(--accent-success);
}

.table-card.occupied {
  border-color: var(--accent-danger);
  background: rgba(239, 68, 68, 0.1);
}

.table-card.reserved {
  border-color: var(--accent-warning);
  background: rgba(245, 158, 11, 0.1);
}

.table-card.maintenance {
  border-color: var(--text-muted);
  opacity: 0.5;
}

.table-card-number {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 4px;
}

.table-card-capacity {
  font-size: 12px;
  color: var(--text-secondary);
}

.table-card-status {
  font-size: 11px;
  text-transform: uppercase;
  margin-top: 8px;
  font-weight: 500;
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
