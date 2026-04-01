<template>
  <div class="billing-page">
    <div class="page-actions">
      <div class="filters">
        <button 
          v-for="status in statusFilters" 
          :key="status.value"
          :class="['filter-btn', { active: activeFilter === status.value }]"
          @click="activeFilter = status.value"
        >{{ status.label }}</button>
      </div>
      <div style="display: flex; gap: 12px; align-items: center;">
        <input type="date" v-model="dateFilter" class="form-input" style="width: auto;">
        <select v-model="paymentMethodFilter" class="form-input form-select" style="width: auto;">
          <option value="">All Methods</option>
          <option v-for="method in paymentMethods" :key="method.id" :value="method.id">{{ method.name }}</option>
        </select>
      </div>
    </div>
    
    <div class="billing-grid">
      <div class="card" style="grid-column: span 2;">
        <div class="card-header">
          <h3 class="card-title">Orders & Payments</h3>
        </div>
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th>Order #</th>
                <th>Table</th>
                <th>Items</th>
                <th>Total</th>
                <th>Payment</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in filteredOrders" :key="order.id">
                <td><strong>#{{ order.order_number }}</strong></td>
                <td>{{ order.table_number || 'Takeaway' }}</td>
                <td>{{ order.item_count }} items</td>
                <td><strong>NRS {{ order.total.toFixed(2) }}</strong></td>
                <td>{{ order.payment_method || '-' }}</td>
                <td>
                  <span :class="['badge', order.payment_status === 'paid' ? 'badge-success' : 'badge-warning']">
                    {{ order.payment_status }}
                  </span>
                </td>
                <td>
                  <div class="action-btns">
                    <button 
                      v-if="order.payment_status === 'unpaid'"
                      class="btn btn-sm btn-success"
                      @click="openPaymentModal(order)"
                    >Process Payment</button>
                    <button class="btn btn-sm btn-secondary" @click="generateReceipt(order)">Receipt</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">Today's Summary</h3>
        </div>
        <div class="summary-stats">
          <div class="summary-item">
            <div class="summary-label">Total Revenue</div>
            <div class="summary-value success">NRS {{ todaySummary.revenue.toLocaleString() }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">Orders Completed</div>
            <div class="summary-value">{{ todaySummary.ordersCompleted }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">Avg Order Value</div>
            <div class="summary-value">NRS {{ todaySummary.avgOrderValue.toFixed(2) }}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">Pending Payments</div>
            <div class="summary-value warning">{{ todaySummary.pendingPayments }}</div>
          </div>
        </div>
        
        <h4 style="margin: 24px 0 16px;">Payment Methods</h4>
        <div class="payment-breakdown">
          <div v-for="method in paymentBreakdown" :key="method.name" class="payment-method-row">
            <span>{{ method.name }}</span>
            <span><strong>NRS {{ method.amount.toFixed(2) }}</strong></span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Payment Modal -->
    <div :class="['modal-overlay', { active: showPaymentModal }]" @click.self="showPaymentModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">Process Payment - Order #{{ selectedOrder?.order_number }}</h3>
          <button class="modal-close" @click="showPaymentModal = false">✕</button>
        </div>
        <div class="modal-body" v-if="selectedOrder">
          <div class="payment-summary">
            <div class="payment-row"><span>Subtotal</span><span>NRS {{ selectedOrder.subtotal.toFixed(2) }}</span></div>
            <div class="payment-row"><span>Tax</span><span>NRS {{ selectedOrder.tax.toFixed(2) }}</span></div>
            <div class="payment-row total"><span>Total Due</span><span>NRS {{ selectedOrder.total.toFixed(2) }}</span></div>
          </div>
          
          <div class="form-group" style="margin-top: 24px;">
            <label class="form-label">Payment Method</label>
            <div class="payment-methods-grid">
              <button 
                v-for="method in paymentMethods" 
                :key="method.id"
                :class="['payment-method-btn', { active: paymentForm.method === method.id }]"
                @click="paymentForm.method = method.id"
              >
                <span class="method-icon">{{ method.icon }}</span>
                <span>{{ method.name }}</span>
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">Tip (optional)</label>
            <input v-model.number="paymentForm.tip" type="number" step="0.01" class="form-input" placeholder="0.00">
          </div>
          
          <div class="payment-row total" style="margin-top: 16px;">
            <span>Total with Tip</span>
            <span>NRS {{ (selectedOrder.total + (paymentForm.tip || 0)).toFixed(2) }}</span>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showPaymentModal = false">Cancel</button>
          <button class="btn btn-success" @click="processPayment">Complete Payment</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ordersAPI, billingAPI, reportsAPI } from '../api'

const activeFilter = ref('all')
const dateFilter = ref('')
const paymentMethodFilter = ref('')
const showPaymentModal = ref(false)
const selectedOrder = ref(null)
const paymentForm = ref({ method: '', tip: 0 })

const statusFilters = [
  { value: 'all', label: 'All' },
  { value: 'unpaid', label: 'Unpaid' },
  { value: 'paid', label: 'Paid' }
]

const paymentMethods = ref([])
const orders = ref([])
const payments = ref([])

const todaySummary = ref({
  revenue: 0,
  ordersCompleted: 0,
  avgOrderValue: 0,
  pendingPayments: 0
})

const paymentBreakdown = ref([])

const fetchData = async () => {
  try {
    const [ordersRes, methodsRes, paymentsRes, salesRes] = await Promise.all([
      ordersAPI.getAll(),
      billingAPI.getPaymentMethods(),
      billingAPI.getPayments(),
      reportsAPI.getSalesSummary()
    ])

    const allOrders = ordersRes.data || []
    const allPayments = paymentsRes.data || []
    const paidOrderIds = new Set(allPayments.filter(p => p.status === 'completed').map(p => p.order_id))

    orders.value = allOrders.map(o => ({
      ...o,
      total: o.total || 0,
      subtotal: o.subtotal || 0,
      tax: o.tax || 0,
      item_count: o.items?.length || 0,
      payment_status: paidOrderIds.has(o.id) ? 'paid' : 'unpaid',
      payment_method: allPayments.find(p => p.order_id === o.id && p.status === 'completed')?.payment_method_name || null
    }))

    const methods = methodsRes.data || []
    paymentMethods.value = methods.map(m => ({ ...m, icon: m.name === 'Cash' ? '💵' : m.name.includes('Mobile') ? '📱' : '💳' }))

    const sales = salesRes.data
    todaySummary.value = {
      revenue: sales.total_revenue || 0,
      ordersCompleted: sales.total_orders || 0,
      avgOrderValue: sales.avg_order_value || 0,
      pendingPayments: orders.value.filter(o => o.payment_status === 'unpaid').length
    }

    // Build payment breakdown
    const breakdown = {}
    for (const p of allPayments.filter(p => p.status === 'completed')) {
      const name = p.payment_method_name || 'Other'
      breakdown[name] = (breakdown[name] || 0) + (p.amount || 0)
    }
    paymentBreakdown.value = Object.entries(breakdown).map(([name, amount]) => ({ name, amount }))

  } catch (e) {
    console.error('Failed to load billing data', e)
  }
}

onMounted(fetchData)

const filteredOrders = computed(() => {
  let result = orders.value
  if (activeFilter.value !== 'all') {
    result = result.filter(o => o.payment_status === activeFilter.value)
  }
  return result
})

const openPaymentModal = (order) => {
  selectedOrder.value = order
  paymentForm.value = { method: paymentMethods.value[0]?.id || '', tip: 0 }
  showPaymentModal.value = true
}

const processPayment = async () => {
  try {
    await billingAPI.createPayment({
      order_id: selectedOrder.value.id,
      payment_method_id: paymentForm.value.method,
      amount: selectedOrder.value.total + (paymentForm.value.tip || 0),
      tip: paymentForm.value.tip || 0
    })
    showPaymentModal.value = false
    await fetchData()
  } catch (e) {
    alert('Payment failed: ' + (e.response?.data?.error || e.message))
  }
}

const generateReceipt = async (order) => {
  try {
    const res = await billingAPI.generateReceipt(order.id)
    const r = res.data
    alert(`Receipt for Order #${r.order_number}\nTable: ${r.table}\nTotal: NRS ${r.total}\nItems: ${(r.items || []).map(i => `${i.quantity}x ${i.name}`).join(', ')}`)
  } catch (e) {
    alert('Failed to generate receipt: ' + (e.response?.data?.error || e.message))
  }
}
</script>

<style scoped>
.billing-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 24px; }
.action-btns { display: flex; gap: 8px; }
.summary-stats { display: grid; gap: 16px; }
.summary-item { padding: 16px; background: var(--bg-secondary); border-radius: var(--border-radius-sm); }
.summary-label { font-size: 13px; color: var(--text-secondary); margin-bottom: 4px; }
.summary-value { font-size: 24px; font-weight: 700; }
.summary-value.success { color: var(--accent-success); }
.summary-value.warning { color: var(--accent-warning); }
.payment-breakdown { display: flex; flex-direction: column; gap: 8px; }
.payment-method-row { display: flex; justify-content: space-between; padding: 12px; background: var(--bg-secondary); border-radius: var(--border-radius-sm); }
.payment-summary { background: var(--bg-secondary); border-radius: var(--border-radius-sm); padding: 16px; }
.payment-row { display: flex; justify-content: space-between; padding: 8px 0; }
.payment-row.total { font-size: 20px; font-weight: 700; border-top: 1px solid var(--border-color); margin-top: 8px; padding-top: 16px; }
.payment-methods-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.payment-method-btn { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 16px; background: var(--bg-secondary); border: 2px solid var(--border-color); border-radius: var(--border-radius-sm); color: var(--text-primary); transition: var(--transition-fast); }
.payment-method-btn:hover, .payment-method-btn.active { border-color: var(--accent-primary); background: rgba(99, 102, 241, 0.1); }
.method-icon { font-size: 24px; }
@media (max-width: 1024px) { .billing-grid { grid-template-columns: 1fr; } }
</style>
