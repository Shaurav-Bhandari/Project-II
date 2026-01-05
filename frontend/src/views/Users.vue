<template>
  <div class="users-page">
    <div class="page-actions">
      <h2>User Management</h2>
      <button class="btn btn-primary" @click="openModal()">+ Add User</button>
    </div>
    
    <div class="card">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>User</th>
              <th>Email</th>
              <th>Role</th>
              <th>Status</th>
              <th>Last Login</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>
                <div style="display: flex; align-items: center; gap: 12px;">
                  <div class="user-avatar" style="width: 40px; height: 40px;">{{ user.first_name[0] }}{{ user.last_name[0] }}</div>
                  <div>
                    <div><strong>{{ user.first_name }} {{ user.last_name }}</strong></div>
                  </div>
                </div>
              </td>
              <td>{{ user.email }}</td>
              <td>
                <span :class="['badge', getRoleBadge(user.role)]">{{ user.role }}</span>
              </td>
              <td>
                <span :class="['badge', user.is_active ? 'badge-success' : 'badge-danger']">
                  {{ user.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td>{{ formatDate(user.last_login) }}</td>
              <td>
                <div class="action-btns">
                  <button class="btn btn-sm btn-secondary" @click="openModal(user)">Edit</button>
                  <button 
                    :class="['btn', 'btn-sm', user.is_active ? 'btn-warning' : 'btn-success']"
                    @click="toggleActive(user)"
                  >
                    {{ user.is_active ? 'Deactivate' : 'Activate' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- User Modal -->
    <div :class="['modal-overlay', { active: showModal }]" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ editingUser ? 'Edit User' : 'Add New User' }}</h3>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
            <div class="form-group">
              <label class="form-label">First Name</label>
              <input v-model="userForm.first_name" type="text" class="form-input" required>
            </div>
            <div class="form-group">
              <label class="form-label">Last Name</label>
              <input v-model="userForm.last_name" type="text" class="form-input" required>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Email</label>
            <input v-model="userForm.email" type="email" class="form-input" required>
          </div>
          <div class="form-group" v-if="!editingUser">
            <label class="form-label">Password</label>
            <input v-model="userForm.password" type="password" class="form-input" required>
          </div>
          <div class="form-group">
            <label class="form-label">Role</label>
            <select v-model="userForm.role" class="form-input form-select">
              <option v-for="role in roles" :key="role.id" :value="role.name">{{ role.name }}</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">Cancel</button>
          <button class="btn btn-primary" @click="saveUser">{{ editingUser ? 'Update' : 'Add' }} User</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const showModal = ref(false)
const editingUser = ref(null)

const roles = ref([
  { id: 1, name: 'admin' },
  { id: 2, name: 'manager' },
  { id: 3, name: 'staff' },
  { id: 4, name: 'kitchen' }
])

const users = ref([
  { id: 1, first_name: 'System', last_name: 'Admin', email: 'admin@restaurant.com', role: 'admin', is_active: true, last_login: new Date() },
  { id: 2, first_name: 'John', last_name: 'Manager', email: 'john@restaurant.com', role: 'manager', is_active: true, last_login: new Date(Date.now() - 86400000) },
  { id: 3, first_name: 'Sarah', last_name: 'Staff', email: 'sarah@restaurant.com', role: 'staff', is_active: true, last_login: new Date(Date.now() - 3600000) },
  { id: 4, first_name: 'Marco', last_name: 'Chef', email: 'marco@restaurant.com', role: 'kitchen', is_active: true, last_login: new Date(Date.now() - 7200000) },
  { id: 5, first_name: 'Jane', last_name: 'Doe', email: 'jane@restaurant.com', role: 'staff', is_active: false, last_login: new Date(Date.now() - 604800000) }
])

const userForm = ref({ first_name: '', last_name: '', email: '', password: '', role: 'staff' })

const getRoleBadge = (role) => {
  const badges = { admin: 'badge-danger', manager: 'badge-warning', staff: 'badge-primary', kitchen: 'badge-info' }
  return badges[role] || 'badge-primary'
}

const formatDate = (date) => {
  if (!date) return 'Never'
  return new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const openModal = (user = null) => {
  if (user) {
    editingUser.value = user
    userForm.value = { ...user, password: '' }
  } else {
    editingUser.value = null
    userForm.value = { first_name: '', last_name: '', email: '', password: '', role: 'staff' }
  }
  showModal.value = true
}

const saveUser = () => {
  if (editingUser.value) {
    Object.assign(editingUser.value, userForm.value)
  } else {
    users.value.push({ ...userForm.value, id: Date.now(), is_active: true, last_login: null })
  }
  showModal.value = false
}

const toggleActive = (user) => {
  user.is_active = !user.is_active
}
</script>

<style scoped>
.action-btns { display: flex; gap: 8px; }
</style>
