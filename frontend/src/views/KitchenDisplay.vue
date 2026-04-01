<template>
  <div class="kds-fullscreen">
    <!-- KDS Header -->
    <header class="kds-header">
      <div class="kds-brand">
        <span class="kds-logo">👨‍🍳</span>
        <span class="kds-title">Kitchen Display</span>
      </div>
      <div class="kds-stats">
        <div class="kds-stat pending">
          <span class="kds-stat-count">{{ pendingCount }}</span>
          <span class="kds-stat-label">Pending</span>
        </div>
        <div class="kds-stat cooking">
          <span class="kds-stat-count">{{ cookingCount }}</span>
          <span class="kds-stat-label">Cooking</span>
        </div>
        <div class="kds-stat ready">
          <span class="kds-stat-count">{{ readyCount }}</span>
          <span class="kds-stat-label">Ready</span>
        </div>
      </div>
      <div class="kds-clock">
        <div class="kds-time">{{ currentTime }}</div>
        <div class="kds-controls">
          <button
            :class="['kds-filter-btn', { active: statusFilter === 'active' }]"
            @click="statusFilter = 'active'"
          >Active</button>
          <button
            :class="['kds-filter-btn', { active: statusFilter === 'all' }]"
            @click="statusFilter = 'all'"
          >All</button>
          <button class="kds-refresh-btn" @click="fetchKots" title="Refresh">🔄</button>
          <button class="kds-logout-btn" @click="handleLogout" title="Logout">🚪</button>
        </div>
      </div>
    </header>

    <!-- KOT Grid -->
    <div class="kds-grid">
      <div
        v-for="kot in filteredKots"
        :key="kot.id"
        :class="['kds-card', `status-${kot.status}`]"
      >
        <!-- Card Header -->
        <div class="kds-card-header">
          <div class="kds-card-id">
            <span class="kds-kot-number">KOT #{{ kot.kot_number }}</span>
            <span class="kds-table">{{ kot.table_number ? 'Table ' + kot.table_number : 'Takeaway' }}</span>
          </div>
          <div class="kds-card-meta">
            <span :class="['kds-status-badge', kot.status]">
              {{ statusLabel(kot.status) }}
            </span>
            <span :class="['kds-timer', { urgent: getMinutesAgo(kot.created_at) > 15 }]">
              {{ getMinutesAgo(kot.created_at) }}m
            </span>
          </div>
        </div>

        <!-- Items -->
        <div class="kds-card-items">
          <div
            v-for="item in kot.items"
            :key="item.id"
            :class="['kds-item', { done: item.status === 'completed' }]"
            @click="toggleItem(kot.id, item)"
          >
            <span class="kds-item-qty">{{ item.quantity }}×</span>
            <div class="kds-item-info">
              <span class="kds-item-name">{{ item.menu_item_name }}</span>
              <span v-if="item.special_instructions" class="kds-item-note">
                📝 {{ item.special_instructions }}
              </span>
            </div>
            <span class="kds-item-check">{{ item.status === 'completed' ? '✅' : '⬜' }}</span>
          </div>
        </div>

        <!-- Card Footer -->
        <div class="kds-card-footer">
          <span class="kds-order-ref">Order #{{ kot.order_number }}</span>
          <div class="kds-card-actions">
            <button
              v-if="kot.status === 'pending'"
              class="kds-action-btn start"
              @click="updateStatus(kot, 'in_progress')"
            >🔥 Start</button>
            <button
              v-if="kot.status === 'in_progress'"
              class="kds-action-btn ready"
              @click="updateStatus(kot, 'completed')"
            >✅ Ready</button>
            <span v-if="kot.status === 'completed'" class="kds-done-label">Done</span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="filteredKots.length === 0" class="kds-empty">
        <div class="kds-empty-icon">👨‍🍳</div>
        <div class="kds-empty-text">No orders in the kitchen</div>
        <div class="kds-empty-sub">New orders will appear automatically</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { kotAPI } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const kots = ref([])
const statusFilter = ref('active')
const currentTime = ref('')
let pollInterval = null
let clockInterval = null

const updateClock = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const fetchKots = async () => {
  try {
    const res = await kotAPI.getAll()
    kots.value = res.data || []
  } catch (e) {
    console.error('KDS fetch error', e)
  }
}

onMounted(() => {
  updateClock()
  fetchKots()
  pollInterval = setInterval(fetchKots, 10000)
  clockInterval = setInterval(updateClock, 1000)
})

onUnmounted(() => {
  clearInterval(pollInterval)
  clearInterval(clockInterval)
})

const pendingCount = computed(() => kots.value.filter(k => k.status === 'pending').length)
const cookingCount = computed(() => kots.value.filter(k => k.status === 'in_progress').length)
const readyCount = computed(() => kots.value.filter(k => k.status === 'completed').length)

const filteredKots = computed(() => {
  if (statusFilter.value === 'active') {
    return kots.value.filter(k => k.status !== 'completed')
  }
  return kots.value
})

const statusLabel = (s) => {
  const map = { pending: 'Pending', in_progress: 'Cooking', completed: 'Ready' }
  return map[s] || s
}

const getMinutesAgo = (dateStr) => {
  if (!dateStr) return 0
  return Math.floor((Date.now() - new Date(dateStr).getTime()) / 60000)
}

const updateStatus = async (kot, newStatus) => {
  try {
    await kotAPI.updateStatus(kot.id, newStatus)
    kot.status = newStatus
  } catch (e) {
    console.error('Failed to update KOT', e)
  }
}

const toggleItem = async (kotId, item) => {
  const newStatus = item.status === 'completed' ? 'in_progress' : 'completed'
  try {
    await kotAPI.updateItemStatus(kotId, item.id, newStatus)
    item.status = newStatus
  } catch (e) {
    console.error('Failed to toggle item', e)
  }
}
const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
/* ===== FULLSCREEN LAYOUT (LIGHT MODE) ===== */
.kds-fullscreen {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: var(--bg-primary, #F9F8F8);
  color: var(--text-primary, #1a202c);
  display: flex;
  flex-direction: column;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  overflow: hidden;
}

/* ===== HEADER ===== */
.kds-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: var(--bg-card, #ffffff);
  border-bottom: 1px solid var(--border-color, #D2D8E7);
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
}
.kds-brand { display: flex; align-items: center; gap: 12px; }
.kds-logo { font-size: 28px; }
.kds-title { font-size: 20px; font-weight: 700; color: var(--text-primary, #1a202c); }

.kds-stats { display: flex; gap: 20px; }
.kds-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 16px;
  border-radius: 10px;
}
.kds-stat.pending { background: rgba(239, 68, 68, 0.08); }
.kds-stat.cooking { background: rgba(245, 158, 11, 0.08); }
.kds-stat.ready   { background: rgba(16, 185, 129, 0.08); }
.kds-stat-count { font-size: 22px; font-weight: 800; }
.kds-stat.pending .kds-stat-count { color: #dc2626; }
.kds-stat.cooking .kds-stat-count { color: #d97706; }
.kds-stat.ready .kds-stat-count   { color: #059669; }
.kds-stat-label { font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; color: var(--text-muted, #9ca3af); }

.kds-clock { display: flex; align-items: center; gap: 16px; }
.kds-time { font-size: 18px; font-weight: 600; font-variant-numeric: tabular-nums; color: var(--text-secondary, #6b7280); }
.kds-controls { display: flex; gap: 8px; }
.kds-filter-btn {
  padding: 6px 16px;
  border-radius: 6px;
  background: var(--bg-secondary, #D2D8E7);
  color: var(--text-secondary, #6b7280);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid var(--border-color, #D2D8E7);
  transition: all 0.15s;
  cursor: pointer;
}
.kds-filter-btn:hover { background: var(--border-color, #c0c8d9); }
.kds-filter-btn.active {
  background: var(--accent-primary, #228CE0);
  color: white;
  border-color: var(--accent-primary, #228CE0);
}
.kds-refresh-btn {
  padding: 6px 12px;
  border-radius: 6px;
  background: var(--bg-secondary, #D2D8E7);
  color: var(--text-secondary, #6b7280);
  font-size: 16px;
  border: 1px solid var(--border-color, #D2D8E7);
  transition: all 0.15s;
  cursor: pointer;
}
.kds-refresh-btn:hover { background: var(--border-color, #c0c8d9); }
.kds-logout-btn {
  padding: 6px 12px;
  border-radius: 6px;
  background: rgba(220, 38, 38, 0.08);
  color: #dc2626;
  font-size: 16px;
  border: 1px solid rgba(220, 38, 38, 0.2);
  transition: all 0.15s;
  cursor: pointer;
}
.kds-logout-btn:hover { background: #dc2626; color: white; }

/* ===== KOT GRID ===== */
.kds-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  padding: 20px 24px;
  overflow-y: auto;
  align-content: start;
}
.kds-grid::-webkit-scrollbar { width: 6px; }
.kds-grid::-webkit-scrollbar-thumb { background: var(--border-color, #D2D8E7); border-radius: 3px; }

/* ===== KOT CARD ===== */
.kds-card {
  background: var(--bg-card, #ffffff);
  border-radius: 12px;
  overflow: hidden;
  border-left: 5px solid;
  transition: transform 0.15s, box-shadow 0.15s;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}
.kds-card:hover { transform: scale(1.01); box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
.kds-card.status-pending     { border-left-color: #dc2626; }
.kds-card.status-in_progress { border-left-color: #d97706; }
.kds-card.status-completed   { border-left-color: #059669; opacity: 0.65; }

.kds-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 14px 16px 10px;
  border-bottom: 1px solid var(--border-color, #e5e7eb);
}
.kds-card-id { display: flex; flex-direction: column; gap: 2px; }
.kds-kot-number { font-size: 16px; font-weight: 800; color: var(--text-primary, #1a202c); }
.kds-table { font-size: 12px; color: var(--text-muted, #9ca3af); }
.kds-card-meta { display: flex; align-items: center; gap: 10px; }

.kds-status-badge {
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.kds-status-badge.pending     { background: rgba(220, 38, 38, 0.1); color: #dc2626; }
.kds-status-badge.in_progress { background: rgba(217, 119, 6, 0.1); color: #d97706; }
.kds-status-badge.completed   { background: rgba(5, 150, 105, 0.1); color: #059669; }

.kds-timer {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-muted, #9ca3af);
  font-variant-numeric: tabular-nums;
}
.kds-timer.urgent {
  color: #dc2626;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* ===== ITEMS ===== */
.kds-card-items {
  padding: 8px 16px;
  flex: 1;
}
.kds-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color, #f0f0f0);
  cursor: pointer;
  transition: background 0.1s;
}
.kds-item:last-child { border-bottom: none; }
.kds-item:hover { background: rgba(34, 140, 224, 0.04); margin: 0 -16px; padding-left: 16px; padding-right: 16px; }
.kds-item.done { opacity: 0.4; }

.kds-item-qty {
  font-weight: 800;
  color: var(--accent-primary, #228CE0);
  font-size: 16px;
  min-width: 30px;
}
.kds-item-info { flex: 1; }
.kds-item-name {
  font-weight: 600;
  font-size: 14px;
  display: block;
  color: var(--text-primary, #1a202c);
}
.kds-item.done .kds-item-name { text-decoration: line-through; }
.kds-item-note {
  font-size: 11px;
  color: #d97706;
  display: block;
  margin-top: 2px;
}
.kds-item-check { font-size: 16px; }

/* ===== FOOTER ===== */
.kds-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: var(--bg-secondary, #f3f4f6);
  border-top: 1px solid var(--border-color, #e5e7eb);
}
.kds-order-ref {
  font-size: 12px;
  color: var(--text-muted, #9ca3af);
}
.kds-card-actions { display: flex; gap: 8px; }

.kds-action-btn {
  padding: 8px 20px;
  border-radius: 8px;
  font-weight: 700;
  font-size: 13px;
  transition: all 0.15s;
  border: none;
  cursor: pointer;
}
.kds-action-btn.start {
  background: rgba(217, 119, 6, 0.1);
  color: #d97706;
  border: 1px solid rgba(217, 119, 6, 0.25);
}
.kds-action-btn.start:hover {
  background: #d97706;
  color: white;
}
.kds-action-btn.ready {
  background: rgba(5, 150, 105, 0.1);
  color: #059669;
  border: 1px solid rgba(5, 150, 105, 0.25);
}
.kds-action-btn.ready:hover {
  background: #059669;
  color: white;
}
.kds-done-label {
  font-size: 12px;
  color: #059669;
  font-weight: 600;
}

/* ===== EMPTY STATE ===== */
.kds-empty {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}
.kds-empty-icon { font-size: 64px; margin-bottom: 16px; }
.kds-empty-text { font-size: 20px; font-weight: 700; color: var(--text-secondary, #6b7280); }
.kds-empty-sub { font-size: 14px; color: var(--text-muted, #9ca3af); margin-top: 4px; }
</style>
