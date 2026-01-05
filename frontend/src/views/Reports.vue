<template>
  <div class="reports-page">
    <div class="page-actions">
      <h2>Reports & Analytics</h2>
      <div style="display: flex; gap: 12px; align-items: center;">
        <input type="date" v-model="dateFrom" class="form-input" style="width: auto;">
        <span>to</span>
        <input type="date" v-model="dateTo" class="form-input" style="width: auto;">
        <button class="btn btn-primary" @click="refreshData">Refresh</button>
      </div>
    </div>
    
    <!-- Summary Cards -->
    <div class="stats-grid">
      <div class="stat-card primary">
        <div class="stat-icon">💰</div>
        <div class="stat-value">${{ summaryStats.totalRevenue.toLocaleString() }}</div>
        <div class="stat-label">Total Revenue</div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">📝</div>
        <div class="stat-value">{{ summaryStats.totalOrders }}</div>
        <div class="stat-label">Total Orders</div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">🍽️</div>
        <div class="stat-value">${{ summaryStats.avgOrderValue.toFixed(2) }}</div>
        <div class="stat-label">Avg Order Value</div>
      </div>
      <div class="stat-card danger">
        <div class="stat-icon">👥</div>
        <div class="stat-value">{{ summaryStats.customersServed }}</div>
        <div class="stat-label">Customers Served</div>
      </div>
    </div>
    
    <div class="reports-grid">
      <!-- Revenue Chart -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Revenue Trend</h3>
        </div>
        <div class="chart-container">
          <canvas ref="revenueChart"></canvas>
        </div>
      </div>
      
      <!-- Popular Items -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Popular Menu Items</h3>
        </div>
        <div class="popular-items-list">
          <div v-for="(item, index) in popularItems" :key="item.name" class="popular-item">
            <span class="item-rank">{{ index + 1 }}</span>
            <div class="item-info">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-category">{{ item.category }}</div>
            </div>
            <div class="item-stats">
              <div class="item-orders">{{ item.orders }} orders</div>
              <div class="item-revenue">${{ item.revenue.toFixed(2) }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Orders by Status -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Orders by Status</h3>
        </div>
        <div class="status-breakdown">
          <div v-for="status in ordersByStatus" :key="status.name" class="status-item">
            <div class="status-bar">
              <div class="status-fill" :style="{ width: status.percentage + '%', background: status.color }"></div>
            </div>
            <div class="status-info">
              <span :style="{ color: status.color }">{{ status.name }}</span>
              <span>{{ status.count }} ({{ status.percentage }}%)</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Peak Hours -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Peak Hours</h3>
        </div>
        <div class="peak-hours">
          <div v-for="hour in peakHours" :key="hour.time" class="peak-hour-item">
            <span class="hour-time">{{ hour.time }}</span>
            <div class="hour-bar">
              <div class="hour-fill" :style="{ width: (hour.orders / 25 * 100) + '%' }"></div>
            </div>
            <span class="hour-orders">{{ hour.orders }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const dateFrom = ref(new Date().toISOString().split('T')[0])
const dateTo = ref(new Date().toISOString().split('T')[0])
const revenueChart = ref(null)

const summaryStats = ref({
  totalRevenue: 28475.50,
  totalOrders: 472,
  avgOrderValue: 60.33,
  customersServed: 389
})

const popularItems = ref([
  { name: 'Grilled Salmon', category: 'Main Courses', orders: 89, revenue: 2224.11 },
  { name: 'Ribeye Steak', category: 'Main Courses', orders: 67, revenue: 2344.33 },
  { name: 'Chicken Wings', category: 'Appetizers', orders: 124, revenue: 1610.76 },
  { name: 'Vegetable Pasta', category: 'Main Courses', orders: 56, revenue: 951.44 },
  { name: 'Spring Rolls', category: 'Appetizers', orders: 98, revenue: 881.02 }
])

const ordersByStatus = ref([
  { name: 'Completed', count: 412, percentage: 87, color: '#10b981' },
  { name: 'Cancelled', count: 23, percentage: 5, color: '#ef4444' },
  { name: 'In Progress', count: 37, percentage: 8, color: '#f59e0b' }
])

const peakHours = ref([
  { time: '11:00', orders: 12 },
  { time: '12:00', orders: 24 },
  { time: '13:00', orders: 21 },
  { time: '18:00', orders: 18 },
  { time: '19:00', orders: 25 },
  { time: '20:00', orders: 22 },
  { time: '21:00', orders: 15 }
])

const refreshData = () => {
  alert('Refreshing data for selected date range...')
}

onMounted(() => {
  // Chart would be initialized here with Chart.js
})
</script>

<style scoped>
.reports-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 24px; }
.chart-container { height: 250px; display: flex; align-items: center; justify-content: center; background: var(--bg-secondary); border-radius: var(--border-radius-sm); }
.chart-container::before { content: '📊 Revenue Chart'; color: var(--text-muted); }

.popular-items-list { display: flex; flex-direction: column; }
.popular-item { display: flex; align-items: center; gap: 16px; padding: 12px 0; border-bottom: 1px solid var(--border-color); }
.popular-item:last-child { border-bottom: none; }
.item-rank { width: 28px; height: 28px; background: var(--gradient-primary); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 13px; }
.item-info { flex: 1; }
.item-name { font-weight: 600; }
.item-category { font-size: 12px; color: var(--text-muted); }
.item-stats { text-align: right; }
.item-orders { font-size: 13px; color: var(--text-secondary); }
.item-revenue { font-weight: 600; color: var(--accent-success); }

.status-breakdown { display: flex; flex-direction: column; gap: 16px; }
.status-item { display: flex; flex-direction: column; gap: 8px; }
.status-bar { height: 8px; background: var(--bg-secondary); border-radius: 4px; overflow: hidden; }
.status-fill { height: 100%; border-radius: 4px; transition: width 0.5s ease; }
.status-info { display: flex; justify-content: space-between; font-size: 13px; }

.peak-hours { display: flex; flex-direction: column; gap: 12px; }
.peak-hour-item { display: flex; align-items: center; gap: 12px; }
.hour-time { width: 50px; font-size: 13px; color: var(--text-secondary); }
.hour-bar { flex: 1; height: 24px; background: var(--bg-secondary); border-radius: 4px; overflow: hidden; }
.hour-fill { height: 100%; background: var(--gradient-primary); border-radius: 4px; transition: width 0.5s ease; }
.hour-orders { width: 30px; text-align: right; font-weight: 600; }

@media (max-width: 1024px) { .reports-grid { grid-template-columns: 1fr; } }
</style>
