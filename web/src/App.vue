<template>
	<div id="app">
		<div id="dashboard-wrapper">
			<div id="dashboard">
				<host-pricing :settings="settings" />
				<div class="date-range">
					<button class="date-range-next" @click="onSetDate(-1)"><icon icon="chevron-left" /></button>
					<div>{{ dateStr }}</div>
					<button class="date-range-next" @click="onSetDate(1)"><icon icon="chevron-right" /></button>
				</div>
				<dashboard-charts :snapshots="snapshots" />
				<dashboard-data
					:month="totals.month || {}"
					:year="totals.year || {}"
					:day="totals.day || {}"
					:total="totals.total || {}" />
			</div>
		</div>
		<div class="extra-links">
			<a href="https://github.com/siacentral/host-dashboard"><icon :icon="['fab', 'github']" /></a>
			<a href="https://siacentral.com"><sia-central /></a>
			<a href="https://sia.tech"><built-with-sia /></a>
		</div>
	</div>
</template>

<script>
import { mapActions } from 'vuex';

import { getStatus, getSnapshots, getTotals, getCoinPrice, getAverageSettings } from '@/utils/api';
import { formatDate } from '@/utils/format';

import DashboardCharts from '@/components/charts/DashboardCharts';
import DashboardData from '@/components/DashboardData';
import HostPricing from '@/components/HostPricing';
import SiaCentral from '@/assets/siacentral.svg';
import BuiltWithSia from '@/assets/built-with-sia.svg';

export default {
	components: {
		BuiltWithSia,
		DashboardCharts,
		DashboardData,
		HostPricing,
		SiaCentral
	},
	data() {
		return {
			loaded: false,
			currentDate: new Date(),
			totals: {},
			settings: {},
			snapshots: [],
			averageSettings: {}
		};
	},
	computed: {
		dateStr() {
			return formatDate(this.currentDate);
		}
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
		...mapActions(['setExchangeRateSC', 'setExchangeRateSF']),
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
						this.settings = status.status.host_settings;
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
				const d = new Date(this.currentDate);

				d.setDate(d.getDate() + n);

				this.currentDate = d;

				this.loadHostData();
			} catch (ex) {
				console.error('App.onChangeDate', ex);
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
	flex: 1;
	margin: auto;
	align-content: safe center;
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
