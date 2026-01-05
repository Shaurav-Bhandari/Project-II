import { defineStore } from 'pinia'
import { menuAPI } from '../api'

export const useMenuStore = defineStore('menu', {
    state: () => ({
        categories: [],
        items: [],
        loading: false,
        error: null
    }),

    getters: {
        availableItems: (state) => state.items.filter(item => item.is_available),
        itemsByCategory: (state) => (categoryId) =>
            state.items.filter(item => item.category_id === categoryId),
        getItemById: (state) => (id) => state.items.find(item => item.id === id)
    },

    actions: {
        async fetchCategories() {
            this.loading = true
            try {
                const response = await menuAPI.getCategories()
                this.categories = response.data
            } catch (error) {
                this.error = error.response?.data?.message || 'Failed to fetch categories'
            } finally {
                this.loading = false
            }
        },

        async fetchItems(categoryId = null) {
            this.loading = true
            try {
                const response = await menuAPI.getItems(categoryId)
                this.items = response.data
            } catch (error) {
                this.error = error.response?.data?.message || 'Failed to fetch menu items'
            } finally {
                this.loading = false
            }
        },

        async createCategory(data) {
            try {
                const response = await menuAPI.createCategory(data)
                this.categories.push(response.data)
                return { success: true, category: response.data }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to create category'
                }
            }
        },

        async updateCategory(id, data) {
            try {
                const response = await menuAPI.updateCategory(id, data)
                const index = this.categories.findIndex(c => c.id === id)
                if (index !== -1) this.categories[index] = response.data
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to update category'
                }
            }
        },

        async createItem(data) {
            try {
                const response = await menuAPI.createItem(data)
                this.items.push(response.data)
                return { success: true, item: response.data }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to create menu item'
                }
            }
        },

        async updateItem(id, data) {
            try {
                const response = await menuAPI.updateItem(id, data)
                const index = this.items.findIndex(i => i.id === id)
                if (index !== -1) this.items[index] = response.data
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to update menu item'
                }
            }
        },

        async toggleItemAvailability(id) {
            try {
                await menuAPI.toggleAvailability(id)
                const item = this.items.find(i => i.id === id)
                if (item) item.is_available = !item.is_available
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to toggle availability'
                }
            }
        }
    }
})
