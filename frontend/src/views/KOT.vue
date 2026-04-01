<template>
  <div class="kot-page">
    <!-- Filters -->
    <div class="filters">
      <button 
        v-for="status in statusFilters" 
        :key="status.value"
        :class="['filter-btn', { active: activeFilter === status.value }]"
        @click="activeFilter = status.value"
      >
        {{ status.label }}
        <span v-if="getCount(status.value) > 0" class="filter-count">{{ getCount(status.value) }}</span>
      </button>
    </div>
    
    <!-- KOT Grid -->
    <div class="kot-grid">
      <div v-for="kot in filteredKots" :key="kot.id" class="kot-card">
        <div class="kot-header" :class="kot.status">
          <div>
            <div class="kot-number">KOT #{{ kot.kot_number }}</div>
            <div class="kot-table">Table {{ kot.table_number }} • Order #{{ kot.order_number }}</div>
          </div>
          <span :class="['badge', getStatusBadge(kot.status)]">
            {{ kot.status === 'in_progress' ? 'Cooking' : kot.status }}
          </span>
        </div>
        
        <div class="kot-body">
          <div v-for="item in kot.items" :key="item.id" class="kot-item">
            <div class="kot-item-check">
              <input 
                type="checkbox" 
                :checked="item.status === 'completed'"
                @change="toggleItemStatus(kot.id, item)"
                :id="`item-${item.id}`"
              >
            </div>
            <div class="kot-qty">{{ item.quantity }}x</div>
            <div class="kot-item-details">
              <div :class="['kot-item-name', { completed: item.status === 'completed' }]">
                {{ item.menu_item_name }}
              </div>
              <div v-if="item.special_instructions" class="kot-item-note">
                📝 {{ item.special_instructions }}
              </div>
            </div>
          </div>
        </div>
        
        <div class="kot-meta">
          <div class="kot-time">
            <span class="time-icon">⏱️</span>
            <span>{{ getTimeAgo(kot.created_at) }}</span>
          </div>
          <div v-if="kot.assigned_chef" class="kot-chef">
            <span>👨‍🍳</span>
            <span>{{ kot.assigned_chef }}</span>
          </div>
        </div>
        
        <div class="kot-footer">
          <button 
            v-if="kot.status === 'pending'"
            class="btn btn-warning" 
            style="flex: 1;"
            @click="updateKotStatus(kot, 'in_progress')"
          >
            Start Cooking
          </button>
          <button 
            v-if="kot.status === 'in_progress'"
            class="btn btn-success" 
            style="flex: 1;"
            @click="updateKotStatus(kot, 'completed')"
          >
            Mark Ready
          </button>
          <button 
            v-if="kot.status === 'completed'"
            class="btn btn-secondary" 
            style="flex: 1;"
            disabled
          >
            ✓ Completed
          </button>
        </div>
      </div>
      
      <div v-if="filteredKots.length === 0" class="empty-state" style="grid-column: 1 / -1; padding: 60px;">
        <div class="empty-state-icon">👨‍🍳</div>
        <div class="empty-state-title">No kitchen orders</div>
        <div class="empty-state-text">All caught up! New orders will appear here.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { kotAPI } from '../api'

const activeFilter = ref('all')
const kots = ref([])

const statusFilters = [
  { value: 'all', label: 'All KOTs' },
  { value: 'pending', label: 'Pending' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' }
]

const fetchKots = async () => {
  try {
    const res = await kotAPI.getAll()
    kots.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch KOTs', e)
  }
}

onMounted(fetchKots)

const filteredKots = computed(() => {
  if (activeFilter.value === 'all') return kots.value
  return kots.value.filter(k => k.status === activeFilter.value)
})

const getCount = (status) => {
  if (status === 'all') return kots.value.length
  return kots.value.filter(k => k.status === status).length
}

const getStatusBadge = (status) => {
  const badges = {
    pending: 'badge-danger',
    in_progress: 'badge-warning',
    completed: 'badge-success'
  }
  return badges[status] || 'badge-primary'
}

const getTimeAgo = (date) => {
  const minutes = Math.floor((Date.now() - new Date(date).getTime()) / 60000)
  if (minutes < 1) return 'Just now'
  if (minutes === 1) return '1 min ago'
  return `${minutes} mins ago`
}

const updateKotStatus = async (kot, status) => {
  try {
    await kotAPI.updateStatus(kot.id, status)
    kot.status = status
  } catch (e) {
    console.error('Failed to update KOT status', e)
  }
}

const toggleItemStatus = async (kotId, item) => {
  const newStatus = item.status === 'completed' ? 'in_progress' : 'completed'
  try {
    await kotAPI.updateItemStatus(kotId, item.id, newStatus)
    item.status = newStatus
  } catch (e) {
    console.error('Failed to toggle item status', e)
  }
}
</script>

<style scoped>
.kot-header.pending { border-left: 4px solid var(--accent-danger); }
.kot-header.in_progress { border-left: 4px solid var(--accent-warning); }
.kot-header.completed { border-left: 4px solid var(--accent-success); }

.filter-count {
  background: var(--accent-danger);
  color: white;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 10px;
  margin-left: 6px;
}

.kot-item-check {
  margin-right: 8px;
}

.kot-item-check input {
  width: 18px;
  height: 18px;
  accent-color: var(--accent-success);
}

.kot-item-name.completed {
  text-decoration: line-through;
  color: var(--text-muted);
}

.kot-item-details {
  flex: 1;
}

.kot-meta {
  display: flex;
  justify-content: space-between;
  padding: 12px 20px;
  background: var(--bg-secondary);
  font-size: 13px;
  color: var(--text-secondary);
}

.kot-time, .kot-chef {
  display: flex;
  align-items: center;
  gap: 6px;
}

@media (max-width: 768px) {
  .filters {
    display: flex;
    overflow-x: auto;
    padding-bottom: 8px;
    margin-bottom: 16px;
  }
  .filters::-webkit-scrollbar { height: 4px; }
  .filters::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 2px; }
  .filter-btn { white-space: nowrap; }
}
</style>
