import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

const store = new Vuex.Store({
	state: {
		currency: localStorage.getItem('displayCurrency') || 'usd',
		exchangeRateSC: {},
		exchangeRateSF: {}
	},
	mutations: {
		setCurrency(state, currency) {
			state.currency = currency;
		},
		setExchangeRateSC(state, rates) {
			state.exchangeRateSC = rates;
		},
		setExchangeRateSF(state, rates) {
			state.exchangeRateSF = rates;
		}
	},
	actions: {
		setCurrency({ commit }, currency) {
			localStorage.setItem('displayCurrency', currency);
			commit('setCurrency', currency);
		},
		setExchangeRateSC({ commit }, rates) {
			commit('setExchangeRateSC', rates);
		},
		setExchangeRateSF({ commit }, rates) {
			commit('setExchangeRateSF', rates);
		}
	}
});

export default store;