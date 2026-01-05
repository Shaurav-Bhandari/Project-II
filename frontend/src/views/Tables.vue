<template>
  <div class="tables-page">
    <div class="page-actions">
      <div class="table-stats">
        <span class="stat"><span class="dot available"></span> {{ availableCount }} Available</span>
        <span class="stat"><span class="dot occupied"></span> {{ occupiedCount }} Occupied</span>
        <span class="stat"><span class="dot reserved"></span> {{ reservedCount }} Reserved</span>
      </div>
      <button class="btn btn-primary" @click="openModal()">+ Add Table</button>
    </div>
    
    <div class="tables-grid">
      <div 
        v-for="table in tables" 
        :key="table.id" 
        :class="['table-card-large', table.status]"
        @click="selectTable(table)"
      >
        <div class="table-icon">🪑</div>
        <div class="table-number">{{ table.table_number }}</div>
        <div class="table-capacity">{{ table.capacity }} seats</div>
        <div class="table-location">{{ table.location }}</div>
        <span :class="['badge', getStatusBadge(table.status)]">{{ table.status }}</span>
        
        <div class="table-actions">
          <button 
            v-if="table.status === 'available'" 
            class="btn btn-sm btn-success"
            @click.stop="updateStatus(table, 'occupied')"
          >Seat Guest</button>
          <button 
            v-if="table.status === 'occupied'" 
            class="btn btn-sm btn-warning"
            @click.stop="updateStatus(table, 'available')"
          >Free Table</button>
          <button 
            v-if="table.status === 'reserved'" 
            class="btn btn-sm btn-primary"
            @click.stop="updateStatus(table, 'occupied')"
          >Check In</button>
          <button class="btn btn-sm btn-secondary" @click.stop="openModal(table)">Edit</button>
        </div>
      </div>
    </div>
    
    <!-- Table Allocation Tool -->
    <div class="card" style="margin-top: 32px;">
      <div class="card-header">
        <h3 class="card-title">🎯 Smart Table Allocation</h3>
      </div>
      <div style="display: flex; gap: 16px; align-items: flex-end;">
        <div class="form-group" style="margin: 0; flex: 1;">
          <label class="form-label">Group Size</label>
          <input v-model.number="groupSize" type="number" min="1" max="20" class="form-input" placeholder="Enter party size">
        </div>
        <button class="btn btn-primary" @click="findTables" style="height: 46px;">Find Tables</button>
      </div>
      <div v-if="allocationResult" class="allocation-result" style="margin-top: 16px;">
        <div v-if="allocationResult.found" class="result-success">
          <strong>✓ Recommended Tables:</strong> {{ allocationResult.tables.join(', ') }}
          <div style="color: var(--text-secondary); font-size: 13px; margin-top: 4px;">
            Total capacity: {{ allocationResult.totalCapacity }} seats
          </div>
        </div>
        <div v-else class="result-error">
          ✕ No available table combination found for {{ groupSize }} guests
        </div>
      </div>
    </div>
    
    <!-- Table Modal -->
    <div :class="['modal-overlay', { active: showModal }]" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ editingTable ? 'Edit Table' : 'Add Table' }}</h3>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Table Number</label>
            <input v-model="tableForm.table_number" type="text" class="form-input" placeholder="e.g., T11">
          </div>
          <div class="form-group">
            <label class="form-label">Capacity</label>
            <input v-model.number="tableForm.capacity" type="number" min="1" class="form-input">
          </div>
          <div class="form-group">
            <label class="form-label">Location</label>
            <select v-model="tableForm.location" class="form-input form-select">
              <option>Window</option>
              <option>Center</option>
              <option>Corner</option>
              <option>Patio</option>
              <option>Private</option>
              <option>Bar</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Status</label>
            <select v-model="tableForm.status" class="form-input form-select">
              <option value="available">Available</option>
              <option value="occupied">Occupied</option>
              <option value="reserved">Reserved</option>
              <option value="maintenance">Maintenance</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary" @click="saveTable">{{ editingTable ? 'Update' : 'Add' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const showModal = ref(false)
const editingTable = ref(null)
const groupSize = ref(4)
const allocationResult = ref(null)

const tables = ref([
  { id: 1, table_number: 'T1', capacity: 2, location: 'Window', status: 'occupied' },
  { id: 2, table_number: 'T2', capacity: 2, location: 'Window', status: 'available' },
  { id: 3, table_number: 'T3', capacity: 4, location: 'Center', status: 'occupied' },
  { id: 4, table_number: 'T4', capacity: 4, location: 'Center', status: 'available' },
  { id: 5, table_number: 'T5', capacity: 6, location: 'Corner', status: 'reserved' },
  { id: 6, table_number: 'T6', capacity: 6, location: 'Corner', status: 'available' },
  { id: 7, table_number: 'T7', capacity: 8, location: 'Private', status: 'occupied' },
  { id: 8, table_number: 'T8', capacity: 4, location: 'Patio', status: 'available' },
  { id: 9, table_number: 'T9', capacity: 4, location: 'Patio', status: 'maintenance' },
  { id: 10, table_number: 'T10', capacity: 2, location: 'Bar', status: 'available' }
])

const tableForm = ref({ table_number: '', capacity: 4, location: 'Center', status: 'available' })

const availableCount = computed(() => tables.value.filter(t => t.status === 'available').length)
const occupiedCount = computed(() => tables.value.filter(t => t.status === 'occupied').length)
const reservedCount = computed(() => tables.value.filter(t => t.status === 'reserved').length)

const getStatusBadge = (status) => {
  const badges = { available: 'badge-success', occupied: 'badge-danger', reserved: 'badge-warning', maintenance: 'badge-info' }
  return badges[status] || 'badge-primary'
}

const openModal = (table = null) => {
  if (table) {
    editingTable.value = table
    tableForm.value = { ...table }
  } else {
    editingTable.value = null
    tableForm.value = { table_number: '', capacity: 4, location: 'Center', status: 'available' }
  }
  showModal.value = true
}

const saveTable = () => {
  if (editingTable.value) {
    Object.assign(editingTable.value, tableForm.value)
  } else {
    tables.value.push({ ...tableForm.value, id: Date.now() })
  }
  showModal.value = false
}

const updateStatus = (table, status) => {
  table.status = status
}

const selectTable = (table) => {
  // Could show table details or start an order
}

// Backtracking algorithm for table allocation
const findTables = () => {
  const available = tables.value.filter(t => t.status === 'available').sort((a, b) => b.capacity - a.capacity)
  const result = findCombination(available, groupSize.value, [], 0)
  
  if (result) {
    allocationResult.value = {
      found: true,
      tables: result.map(t => t.table_number),
      totalCapacity: result.reduce((sum, t) => sum + t.capacity, 0)
    }
  } else {
    allocationResult.value = { found: false }
  }
}

// Backtracking helper
const findCombination = (tables, target, current, start) => {
  const currentCapacity = current.reduce((sum, t) => sum + t.capacity, 0)
  
  if (currentCapacity >= target) return current
  if (start >= tables.length) return null
  
  // Try including current table
  const withCurrent = findCombination(tables, target, [...current, tables[start]], start + 1)
  if (withCurrent) return withCurrent
  
  // Try without current table
  return findCombination(tables, target, current, start + 1)
}
</script>

<style scoped>
.table-stats { display: flex; gap: 24px; }
.stat { display: flex; align-items: center; gap: 8px; font-size: 14px; }
.dot { width: 10px; height: 10px; border-radius: 50%; }
.dot.available { background: var(--accent-success); }
.dot.occupied { background: var(--accent-danger); }
.dot.reserved { background: var(--accent-warning); }

.tables-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 20px; }

.table-card-large {
  background: var(--bg-card);
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: var(--transition-fast);
}

.table-card-large:hover { transform: translateY(-4px); box-shadow: var(--shadow-lg); }
.table-card-large.available { border-color: var(--accent-success); }
.table-card-large.occupied { border-color: var(--accent-danger); background: rgba(239, 68, 68, 0.05); }
.table-card-large.reserved { border-color: var(--accent-warning); background: rgba(245, 158, 11, 0.05); }
.table-card-large.maintenance { border-color: var(--text-muted); opacity: 0.6; }

.table-icon { font-size: 32px; margin-bottom: 8px; }
.table-number { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.table-capacity { font-size: 14px; color: var(--text-secondary); }
.table-location { font-size: 12px; color: var(--text-muted); margin-bottom: 12px; }

.table-actions { display: flex; flex-direction: column; gap: 8px; margin-top: 16px; }

.allocation-result { padding: 16px; border-radius: var(--border-radius-sm); }
.result-success { background: rgba(16, 185, 129, 0.1); border: 1px solid var(--accent-success); color: var(--accent-success); }
.result-error { background: rgba(239, 68, 68, 0.1); border: 1px solid var(--accent-danger); color: var(--accent-danger); }
</style>
