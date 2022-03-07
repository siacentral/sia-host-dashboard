<template>
	<chart-display title="Revenue"
		:nodes="revenueData.data"
		:labels="revenueData.labels"
		:colors="revenueColors"
		:fills="revenueFills"
		@selected="onSelectRevenue">
		<div class="active-labels labels-revenue">
			<div class="chart-label line-primary">
				<div class="label-title">Earned</div>
				<div v-html="revenueEarnedLabel" />
			</div>
			<div class="chart-label line-secondary">
				<div class="label-title">Potential</div>
				<div v-html="revenuePotLabel" />
			</div>
		</div>
	</chart-display>
</template>

<script>
import { mapState } from 'vuex';
import BigNumber from 'bignumber.js';

import ChartDisplay from '@/components/charts/ChartDisplay';
import { formatPriceString } from '@/utils/format';

export default {
	components: {
		ChartDisplay
	},
	props: {
		snapshots: Array
	},
	data() {
		return {
			active: -1
		};
	},
	created() {
		console.log('okay');
	},
	computed: {
		...mapState(['currency', 'exchangeRateSC']),
		revenueColors() {
			return [
				'#19bdcf',
				'#19cf86'
			];
		},
		revenueFills() {
			return [
				'#225e70',
				'#227051'
			];
		},
		revenueData() {
			console.log('revenueData');
			let data = this.snapshots.reduce((d, s, i) => {
				const timestamp = new Date(s.timestamp);
				timestamp.setMonth(timestamp.getMonth(), 1);
				timestamp.setHours(0, 0, 0, 0);
				const id = timestamp.getTime();

				if (!d[id]) {
					d[id] = {
						timestamp: timestamp,
						earned: new BigNumber(0),
						potential: new BigNumber(0)
					};
				}

				d[id].earned = d[id].earned.plus(s.earned_revenue);
				d[id].potential = d[id].potential.plus(s.potential_revenue);
				return d;
			}, {});
			const keys = Object.keys(data);
			data = keys.map(k => data[k]);
			data.sort((a, b) => a.timestamp - b.timestamp);

			const labels = data.map(d => d.timestamp.toLocaleString([], {
					month: 'short',
					year: 'numeric'
				})),
				earned = data.map(d => d.earned),
				potential = data.map(d => d.potential);

			return {
				data: [earned, potential],
				labels
			};
		},
		revenueEarnedLabel() {
			let i = this.active;

			if (i === -1 || i >= this.revenueData.data[0].length)
				i = this.revenueData.data[0].length - 4;

			const v = this.revenueData.data[0][i],
				sc = formatPriceString(v),
				curr = formatPriceString(v, 2, this.currency, this.exchangeRateSC[this.currency]);

			return [
				`<div class="data-label">${sc.value}<span class="currency-display">${sc.label}</span></div>`,
				`<div class="data-label-secondary">${curr.value}<span class="currency-display">${curr.label}</span></div>`
			].join('');
		},
		revenuePotLabel() {
			let i = this.active;

			if (i === -1 || i >= this.revenueData.data[1].length)
				i = this.revenueData.data[1].length - 4;

			const v = this.revenueData.data[1][i],
				sc = formatPriceString(v),
				curr = formatPriceString(v, 2, this.currency, this.exchangeRateSC[this.currency]);

			return [
				`<div class="data-label">${sc.value}<span class="currency-display">${sc.label}</span></div>`,
				`<div class="data-label-secondary">${curr.value}<span class="currency-display">${curr.label}</span></div>`
			].join('');
		}
	},
	methods: {
		onSelectRevenue(i) {
			try {
				this.active = i;
			} catch (ex) {
				console.error(ex);
			}
		}
	}
};
</script>

<style lang="stylus" scoped>
.active-labels {
	display: grid;
	grid-gap: 10px;
	justify-items: center;
	text-align: right;
	font-size: 1rem;
	grid-template-columns: repeat(2, minmax(0, 1fr));

	.chart-label {
		display: grid;
		grid-gap: 15px;
		grid-template-columns: repeat(2, auto);
		align-items: center;
	}

	.label-title {
		text-align: center;
		font-size: 0.8rem;
	}

	.line-primary {
		color: #19cf86;
	}

	.line-secondary {
		color: #19bdcf;
	}

	.line-tertiary {
		color: #da5454;
	}
}
</style>