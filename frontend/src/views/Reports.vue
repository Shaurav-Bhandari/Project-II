<template>
  <div class="reports-page">
    <div class="page-actions">
      <h2>Reports & Analytics</h2>
      <div class="reports-filters">
        <input type="date" v-model="dateFrom" class="form-input date-input">
        <span>to</span>
        <input type="date" v-model="dateTo" class="form-input date-input">
        <button class="btn btn-primary" @click="refreshData">Refresh</button>
      </div>
    </div>
    
    <!-- Summary Cards -->
    <div class="stats-grid">
      <div class="stat-card primary">
        <div class="stat-icon">💰</div>
        <div class="stat-value">NRS {{ summaryStats.totalRevenue.toLocaleString() }}</div>
        <div class="stat-label">Total Revenue</div>
        <div class="stat-change up">▲ 12.5% vs last period</div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">📝</div>
        <div class="stat-value">{{ summaryStats.totalOrders }}</div>
        <div class="stat-label">Total Orders</div>
        <div class="stat-change up">▲ 8.2% vs last period</div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">🍽️</div>
        <div class="stat-value">NRS {{ summaryStats.avgOrderValue.toFixed(2) }}</div>
        <div class="stat-label">Avg Order Value</div>
        <div class="stat-change up">▲ 3.1% vs last period</div>
      </div>
      <div class="stat-card danger">
        <div class="stat-icon">👥</div>
        <div class="stat-value">{{ summaryStats.customersServed }}</div>
        <div class="stat-label">Customers Served</div>
        <div class="stat-change down">▼ 1.4% vs last period</div>
      </div>
    </div>
    
    <div class="reports-grid">
      <!-- Revenue Chart -->
      <div class="card chart-card">
        <div class="card-header">
          <h3 class="card-title">Revenue Trend</h3>
          <div class="chart-legend">
            <span class="legend-item"><span class="legend-dot revenue"></span> Revenue</span>
            <span class="legend-item"><span class="legend-dot orders"></span> Orders</span>
          </div>
        </div>
        <div class="chart-wrapper">
          <Bar :data="revenueChartData" :options="revenueChartOptions" />
        </div>
      </div>
      
      <!-- Popular Items -->
      <div class="card popular-card">
        <div class="card-header">
          <h3 class="card-title">Popular Menu Items</h3>
          <span class="card-subtitle">Top performers</span>
        </div>
        <div class="popular-items-list">
          <div v-for="(item, index) in popularItems" :key="item.name" class="popular-item">
            <span :class="['item-rank', `rank-${index + 1}`]">{{ index + 1 }}</span>
            <div class="item-info">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-category">{{ item.category }}</div>
            </div>
            <div class="item-stats">
              <div class="item-orders">{{ item.orders }} orders</div>
              <div class="item-revenue">NRS {{ item.revenue.toFixed(2) }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Orders by Status -->
      <div class="card status-card">
        <div class="card-header">
          <h3 class="card-title">Orders by Status</h3>
          <span class="card-subtitle">{{ totalOrders }} total</span>
        </div>
        <div class="donut-container">
          <Doughnut :data="statusChartData" :options="statusChartOptions" />
        </div>
        <div class="status-legend">
          <div v-for="status in ordersByStatus" :key="status.name" class="status-legend-item">
            <span class="status-dot" :style="{ background: status.color }"></span>
            <span class="status-name">{{ status.name }}</span>
            <span class="status-count">{{ status.count }}</span>
            <span class="status-pct">{{ status.percentage }}%</span>
          </div>
        </div>
      </div>
      
      <!-- Peak Hours -->
      <div class="card peak-card">
        <div class="card-header">
          <h3 class="card-title">Peak Hours</h3>
          <span class="card-subtitle">Orders per hour today</span>
        </div>
        <div class="peak-hours">
          <div v-for="hour in peakHours" :key="hour.time" class="peak-hour-item">
            <span class="hour-time">{{ hour.time }}</span>
            <div class="hour-bar-track">
              <div 
                class="hour-fill" 
                :style="{ width: (hour.orders / maxOrders * 100) + '%' }"
              >
                <span class="hour-fill-label">{{ hour.orders }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { reportsAPI } from '../api'
import { Bar, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, BarElement, ArcElement, Title, Tooltip, Legend, Filler)

const dateFrom = ref(new Date(Date.now() - 7 * 86400000).toISOString().split('T')[0])
const dateTo = ref(new Date().toISOString().split('T')[0])

const summaryStats = ref({
  totalRevenue: 0,
  totalOrders: 0,
  avgOrderValue: 0,
  customersServed: 0
})

const popularItems = ref([])
const ordersByStatus = ref([])
const revenueData = ref([])
const peakHours = ref([])

const statusColors = {
  completed: '#10b981',
  cancelled: '#ef4444',
  pending: '#f59e0b',
  confirmed: '#3b82f6',
  preparing: '#f97316',
  ready: '#06b6d4',
  served: '#8b5cf6',
  in_progress: '#f59e0b'
}

const fetchData = async () => {
  try {
    const [salesRes, popularRes, statusRes, revenueRes] = await Promise.all([
      reportsAPI.getSalesSummary({ start_date: dateFrom.value, end_date: dateTo.value }),
      reportsAPI.getPopularItems({ limit: 5 }),
      reportsAPI.getOrdersByStatus(),
      reportsAPI.getRevenueByDate({ days: 7 })
    ])

    const sales = salesRes.data || {}
    summaryStats.value = {
      totalRevenue: sales.total_revenue || 0,
      totalOrders: sales.total_orders || 0,
      avgOrderValue: sales.avg_order_value || 0,
      customersServed: sales.tables_served || 0
    }

    popularItems.value = (popularRes.data || []).map(item => ({
      name: item.name,
      category: item.category || '',
      orders: item.orders || 0,
      revenue: item.revenue || 0
    }))

    ordersByStatus.value = (statusRes.data || []).map(s => ({
      name: s.status,
      count: s.count || 0,
      percentage: Math.round(s.percentage || 0),
      color: statusColors[s.status] || '#6b7280'
    }))

    revenueData.value = revenueRes.data || []

  } catch (e) {
    console.error('Failed to load reports', e)
  }
}

onMounted(fetchData)

const totalOrders = computed(() => ordersByStatus.value.reduce((sum, s) => sum + s.count, 0))

const maxOrders = computed(() => {
  if (peakHours.value.length === 0) return 1
  return Math.max(...peakHours.value.map(h => h.orders))
})

/* -------- Chart.js configs -------- */
const revenueChartData = computed(() => ({
  labels: revenueData.value.map(d => d.date ? new Date(d.date).toLocaleDateString('en-US', { weekday: 'short' }) : ''),
  datasets: [
    {
      label: 'Revenue (NRS)',
      data: revenueData.value.map(d => d.revenue || 0),
      backgroundColor: 'rgba(34, 140, 224, 0.75)',
      borderRadius: 6,
      borderSkipped: false,
      barPercentage: 0.6,
      categoryPercentage: 0.7
    }
  ]
}))

const revenueChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#1a1a2e',
      titleFont: { family: 'Inter', size: 13 },
      bodyFont: { family: 'Inter', size: 12 },
      padding: 12,
      cornerRadius: 8,
      callbacks: {
        label: ctx => `NRS ${ctx.parsed.y.toLocaleString()}`
      }
    }
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: '#718096', font: { family: 'Inter', size: 12 } }
    },
    y: {
      border: { dash: [4, 4] },
      grid: { color: 'rgba(0,0,0,0.06)' },
      ticks: {
        color: '#718096',
        font: { family: 'Inter', size: 12 },
        callback: v => `NRS ${(v / 1000).toFixed(1)}k`
      }
    }
  }
}

const statusChartData = computed(() => ({
  labels: ordersByStatus.value.map(s => s.name),
  datasets: [{
    data: ordersByStatus.value.map(s => s.count),
    backgroundColor: ordersByStatus.value.map(s => s.color),
    borderWidth: 0,
    cutout: '72%',
    spacing: 4,
    borderRadius: 6
  }]
}))

const statusChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#1a1a2e',
      titleFont: { family: 'Inter', size: 13 },
      bodyFont: { family: 'Inter', size: 12 },
      padding: 12,
      cornerRadius: 8
    }
  }
}

const refreshData = () => {
  fetchData()
}
</script>

<style scoped>
/* ===== GRID LAYOUT ===== */
.reports-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

/* ===== CHART CARD ===== */
.chart-card {
  position: relative;
}
.chart-wrapper {
  height: 280px;
  position: relative;
}
.chart-legend {
  display: flex;
  gap: 16px;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}
.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
}
.legend-dot.revenue { background: rgba(34, 140, 224, 0.75); }
.legend-dot.orders  { background: rgba(128, 161, 232, 0.5); }

/* ===== POPULAR ITEMS ===== */
.popular-items-list { display: flex; flex-direction: column; }
.popular-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 0;
  border-bottom: 1px solid var(--border-color);
  transition: background 0.15s ease;
}
.popular-item:last-child { border-bottom: none; }
.popular-item:hover { background: var(--bg-hover); margin: 0 -24px; padding-left: 24px; padding-right: 24px; border-radius: var(--border-radius-sm); }

.item-rank {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  color: #fff;
  flex-shrink: 0;
}
.rank-1 { background: linear-gradient(135deg, #228CE0, #80A1E8); box-shadow: 0 2px 8px rgba(34,140,224,0.35); }
.rank-2 { background: linear-gradient(135deg, #6366f1, #a5b4fc); box-shadow: 0 2px 8px rgba(99,102,241,0.3); }
.rank-3 { background: linear-gradient(135deg, #10b981, #34d399); box-shadow: 0 2px 8px rgba(16,185,129,0.3); }
.rank-4 { background: var(--accent-warning); }
.rank-5 { background: var(--accent-info); color: #1a1a2e; }

.item-info { flex: 1; }
.item-name { font-weight: 600; font-size: 14px; }
.item-category { font-size: 12px; color: var(--text-muted); margin-top: 2px; }
.item-stats { text-align: right; }
.item-orders { font-size: 12px; color: var(--text-secondary); }
.item-revenue { font-weight: 700; color: var(--accent-success); font-size: 15px; }
.card-subtitle { font-size: 12px; color: var(--text-muted); font-weight: 400; }

/* ===== DONUT / STATUS ===== */
.donut-container {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}
.status-legend {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.status-legend-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  padding: 8px 12px;
  border-radius: var(--border-radius-sm);
  background: var(--bg-secondary);
  transition: var(--transition-fast);
}
.status-legend-item:hover { background: var(--bg-hover); }
.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-name { flex: 1; font-weight: 500; }
.status-count { font-weight: 700; }
.status-pct { color: var(--text-muted); font-size: 12px; min-width: 36px; text-align: right; }

/* ===== PEAK HOURS ===== */
.peak-hours { display: flex; flex-direction: column; gap: 14px; }
.peak-hour-item { display: flex; align-items: center; gap: 14px; }
.hour-time {
  width: 50px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  flex-shrink: 0;
}
.hour-bar-track {
  flex: 1;
  height: 32px;
  background: var(--bg-secondary);
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}
.hour-fill {
  height: 100%;
  background: var(--gradient-primary);
  border-radius: 8px;
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 10px;
  min-width: 40px;
}
.hour-fill-label {
  font-size: 12px;
  font-weight: 700;
  color: #fff;
}

/* ===== STAT CHANGE INDICATORS ===== */
.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  margin-top: 12px;
  font-weight: 500;
}
.stat-change.up { color: var(--accent-success); }
.stat-change.down { color: var(--accent-danger); }

.stat-change.down { color: var(--accent-danger); }

@media (max-width: 1024px) { .reports-grid { grid-template-columns: 1fr; } }

.reports-filters {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}
.date-input {
  width: auto;
}

@media (max-width: 768px) {
  .page-actions {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .reports-filters {
    width: 100%;
  }
  .date-input {
    flex: 1;
    min-width: 120px;
  }
  .chart-wrapper {
    height: 220px;
  }
}
</style>
