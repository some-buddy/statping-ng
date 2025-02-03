import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    services: [],
    incidents: [],
    users: [],
    settings: {},
    loading: false,
    error: null
  },
  mutations: {
    // Set the services state
    SET_SERVICES(state, services) {
      state.services = services
    },
    // Set the incidents state
    SET_INCIDENTS(state, incidents) {
      state.incidents = incidents
    },
    // Set the users state
    SET_USERS(state, users) {
      state.users = users
    },
    // Set the settings state
    SET_SETTINGS(state, settings) {
      state.settings = settings
    },
    // Set the loading state
    SET_LOADING(state, loading) {
      state.loading = loading
    },
    // Set the error state
    SET_ERROR(state, error) {
      state.error = error
    }
  },
  actions: {
    // Fetch services from the API
    async fetchServices({ commit }) {
      commit('SET_LOADING', true)
      try {
        const response = await axios.get('/api/services')
        commit('SET_SERVICES', response.data)
      } catch (error) {
        commit('SET_ERROR', error)
      } finally {
        commit('SET_LOADING', false)
      }
    },
    // Fetch incidents from the API
    async fetchIncidents({ commit }) {
      commit('SET_LOADING', true)
      try {
        const response = await axios.get('/api/incidents')
        commit('SET_INCIDENTS', response.data)
      } catch (error) {
        commit('SET_ERROR', error)
      } finally {
        commit('SET_LOADING', false)
      }
    },
    // Fetch users from the API
    async fetchUsers({ commit }) {
      commit('SET_LOADING', true)
      try {
        const response = await axios.get('/api/users')
        commit('SET_USERS', response.data)
      } catch (error) {
        commit('SET_ERROR', error)
      } finally {
        commit('SET_LOADING', false)
      }
    },
    // Fetch settings from the API
    async fetchSettings({ commit }) {
      commit('SET_LOADING', true)
      try {
        const response = await axios.get('/api/settings')
        commit('SET_SETTINGS', response.data)
      } catch (error) {
        commit('SET_ERROR', error)
      } finally {
        commit('SET_LOADING', false)
      }
    }
  },
  getters: {
    // Get all services
    allServices: state => state.services,
    // Get all incidents
    allIncidents: state => state.incidents,
    // Get all users
    allUsers: state => state.users,
    // Get settings
    allSettings: state => state.settings,
    // Get loading state
    isLoading: state => state.loading,
    // Get error state
    getError: state => state.error
  }
})
