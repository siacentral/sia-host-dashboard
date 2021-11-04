<template>
	<div id="app">
		<div id="dashboard-wrapper">
			<div class="control control-inline">
				<select v-model="displayCurrency" @change="onChangeCurrency">
					<optgroup label="Fiat">
						<option value="usd">USD</option>
						<option value="jpy">JPY</option>
						<option value="eur">EUR</option>
						<option value="gbp">GBP</option>
						<option value="aus">AUS</option>
						<option value="cad">CAD</option>
						<option value="rub">RUB</option>
						<option value="cny">CNY</option>
					</optgroup>
					<optgroup label="Crypto">
						<option value="btc">BTC</option>
						<option value="bch">BCH</option>
						<option value="eth">ETH</option>
						<option value="xrp">XRP</option>
						<option value="ltc">LTC</option>
					</optgroup>
				</select>
			</div>
			<div id="dashboard">
				<alert-list :alerts="alerts" />
				<host-stats :settings="settings" :status="status" />
				<dashboard-charts :snapshots="snapshots" />
				<div class="date-range">
					<button class="date-range-next" @click="onSetDate(-1)"><icon icon="chevron-left" /></button>
					<div>{{ dateStr }}</div>
					<button class="date-range-next" @click="onSetDate(1)"><icon icon="chevron-right" /></button>
				</div>
				<dashboard-data
					:month="totals.month || {}"
					:year="totals.year || {}"
					:day="totals.day || {}"
					:total="totals.total || {}" />
			</div>
		</div>
		<div class="extra-links">
			<a target="_blank" href="https://github.com/siacentral/sia-host-dashboard"><icon :icon="['fab', 'github']" /></a>
			<a target="_blank" href="https://siacentral.com"><sia-central /></a>
			<a target="_blank" href="https://sia.tech"><built-with-sia /></a>
		</div>
	</div>
</template>

<script>
import { mapActions, mapState } from 'vuex';

import { getStatus, getSnapshots, getTotals, getCoinPrice, getAverageSettings } from '@/utils/api';
import { formatDate } from '@/utils/format';

import AlertList from '@/components/alerts/AlertList';
import DashboardCharts from '@/components/charts/DashboardCharts';
import DashboardData from '@/components/DashboardData';
import HostStats from '@/components/HostStats';
import SiaCentral from '@/assets/siacentral.svg';
import BuiltWithSia from '@/assets/built-with-sia.svg';

export default {
	components: {
		AlertList,
		BuiltWithSia,
		DashboardCharts,
		DashboardData,
		HostStats,
		SiaCentral
	},
	data() {
		return {
			loaded: false,
			currentDate: new Date(),
			totals: {},
			status: {},
			alerts: [],
			snapshots: [],
			averageSettings: {},
			displayCurrency: 'usd',
			debounceTimeout: null
		};
	},
	computed: {
		...mapState(['currency']),
		dateStr() {
			return formatDate(this.currentDate);
		},
		settings() {
			if (!this.status || typeof this.status.host_settings !== 'object')
				return {};

			return this.status.host_settings;
		}
	},
	beforeMount() {
		const d = new Date();

		d.setHours(23, 0, 0, 0);
		this.currentDate = d;
		this.displayCurrency = this.currency;
	},
	async mounted() {
		try {
			await Promise.all([
				this.loadChainData(),
				this.loadHostData()
			]);

			this.loaded = true;
		} catch (ex) {
			console.error('AppMounted', ex.message);
		}
	},
	methods: {
		...mapActions(['setCurrency', 'setExchangeRateSC', 'setExchangeRateSF']),
		loadChainData() {
			return Promise.all([
				getCoinPrice()
					.then(currencies => {
						this.setExchangeRateSC(currencies.siacoin);
						this.setExchangeRateSF(currencies.siafund);
					}),
				getAverageSettings()
					.then(settings => {
						this.averageSettings = settings;
					})
			]);
		},
		loadHostData() {
			return Promise.all([
				getStatus()
					.then(status => {
						this.status = status.status;

						if (Array.isArray(status.alerts))
							this.alerts = status.alerts;
					}),
				getSnapshots(this.currentDate)
					.then(snapshots => {
						this.snapshots = snapshots;
					}),
				getTotals(this.currentDate)
					.then(totals => {
						this.totals = totals;
					})
			]);
		},
		onSetDate(n) {
			try {
				clearTimeout(this.debounceTimeout);

				const d = new Date(this.currentDate);

				d.setDate(d.getDate() + n);

				this.currentDate = d;

				this.debounceTimeout = setTimeout(this.loadHostData, 300);
			} catch (ex) {
				console.error('App.onChangeDate', ex);
			}
		},
		onChangeCurrency() {
			try {
				this.setCurrency(this.displayCurrency);
			} catch (ex) {
				console.error('App.onChangeCurrency', ex);
			}
		}
	}
};
</script>

<style lang="stylus">
@require "/styles/global.styl";

#app {
	display: flex;
	width: 100%;
	height: 100%;
	margin: auto;
	flex-direction: column;
	overflow-x: hidden;
	overflow-y: auto;
}

#dashboard-wrapper {
	display: grid;
	padding: 15px;
	width: 100%;
	max-width: 1200px;
	grid-template-columns: minmax(0, 1fr) auto;
	flex: 1;
	margin: auto;
	align-content: safe center;

	.control {
		margin-bottom: 10px;
		grid-column: 1 / -1;
		width: 100%;

		@media screen and (min-width: 767px) {
			min-width: 150px;
			width: auto;
			grid-column: 2;
		}
	}

	#dashboard {
		grid-column: 1 / -1;
	}
}

.extra-links {
	display: grid;
	grid-template-columns: repeat(3, auto);
	grid-gap: 30px;
	padding: 15px;
	align-items: center;
	justify-content: center;
	text-align: center;

	a {
		display: inline-block;
		color: rgba(255, 255, 255, 0.54);

		svg, svg.svg-inline--fa.svg-inline--fa {
			width: 28px;
			height: auto;

			g path {
				stroke: rgba(255, 255, 255, 0.54) !important;
			}
		}
	}
}

.date-range {
	display: grid;
	padding: 15px;
	grid-template-columns: repeat(3, auto);
	grid-gap: 15px;
	align-items: center;
	justify-content: center;

	.date-range-next, .date-range-prev {
		display: inline-block;
		color: rgba(255, 255, 255, 0.54);
		border: none;
		background: none;
		outline: none;
		transition: all 0.3s linear;

		&:hover, &:active, &:focus {
			cursor: pointer;
			color: primary;
		}
	}
}
</style>
