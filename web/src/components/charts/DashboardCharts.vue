<template>
	<div class="charts">
		<chart-display title="Revenue"
			:nodes="revenueData.data"
			:labels="revenueData.labels"
			:text="revenueData.values"
			:colors="revenueColors"
			:fills="revenueFills"
			@select="onSelectRevenue">

		</chart-display>
		<chart-display title="Contracts"
			:nodes="contractData.data"
			:labels="contractData.labels"
			:text="revenueData.values"
			:colors="contractColors"
			:fills="contractFills">

		</chart-display>
	</div>
</template>

<script>
import BigNumber from 'bignumber.js';

import ChartDisplay from '@/components/charts/ChartDisplay';
// import { formatPriceString } from '@/utils/format';

export default {
	components: {
		ChartDisplay
	},
	props: {
		snapshots: Array
	},
	computed: {
		revenueData() {
			const data = this.snapshots.reduce((d, s, i) => {
				const timestamp = new Date(s.timestamp),
					last = d.data[0].length - 1,
					prev = i - 1;

				timestamp.setMonth(timestamp.getMonth(), 1);

				if (i === 0) {
					d.data[0].push(new BigNumber(s.earned_revenue));
					d.data[1].push(new BigNumber(s.potential_revenue));
					d.labels.push(new Date(timestamp).toLocaleString([], {
						month: 'short',
						year: 'numeric'
					}));
					return d;
				}

				const prevTimestamp = new Date(this.snapshots[prev].timestamp);

				prevTimestamp.setMonth(prevTimestamp.getMonth(), 1);

				if (prevTimestamp.getTime() !== timestamp.getTime()) {
					d.data[0].push(new BigNumber(s.earned_revenue));
					d.data[1].push(new BigNumber(s.potential_revenue));
					d.labels.push(new Date(timestamp).toLocaleString([], {
						month: 'short',
						year: 'numeric'
					}));
				} else {
					d.data[0][last] = d.data[0][last].plus(new BigNumber(s.earned_revenue));
					d.data[1][last] = d.data[1][last].plus(s.potential_revenue);
				}

				return d;
			}, {
				data: [[], []],
				labels: []
			});

			return data;
		},
		contractData() {
			const data = this.snapshots.reduce((d, s, i) => {
				const timestamp = new Date(s.timestamp),
					last = d.data[0].length - 1,
					prev = i - 1;

				timestamp.setMonth(timestamp.getMonth(), 1);

				if (i === 0) {
					d.data[0].push(new BigNumber(s.successful_contracts));
					d.data[1].push(new BigNumber(s.failed_contracts));
					d.data[2].push(new BigNumber(s.expired_contracts));
					d.labels.push(new Date(timestamp).toLocaleString([], {
						month: 'short',
						year: 'numeric'
					}));
					return d;
				}

				const prevTimestamp = new Date(this.snapshots[prev].timestamp);

				prevTimestamp.setMonth(prevTimestamp.getMonth(), 1);

				if (prevTimestamp.getTime() !== timestamp.getTime()) {
					d.data[0].push(new BigNumber(s.successful_contracts));
					d.data[1].push(new BigNumber(s.failed_contracts));
					d.data[2].push(new BigNumber(s.expired_contracts));
					d.labels.push(new Date(timestamp).toLocaleString([], {
						month: 'short',
						year: 'numeric'
					}));
				} else {
					d.data[0][last] = d.data[0][last].plus(new BigNumber(s.successful_contracts));
					d.data[1][last] = d.data[1][last].plus(new BigNumber(s.failed_contracts));
					d.data[2][last] = d.data[2][last].plus(new BigNumber(s.expired_contracts));
				}

				return d;
			}, {
				data: [[], [], []],
				labels: []
			});

			return data;
		},
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
		contractColors() {
			return [
				'#19bdcf',
				'#da5454',
				'#19cf86'
			];
		},
		contractFills() {
			return [
				'#225e70',
				'#843b3b',
				'#227051'
			];
		}
	},
	methods: {
		onSelectRevenue() {

		}
	}
};
</script>

<style lang="stylus" scoped>
.charts {
	display: grid;
	grid-gap: 15px;
	margin-bottom: 15px;
}

.charts {
	@media screen and (min-width: 767px) {
		grid-template-columns: repeat(2, minmax(0, 1fr));
	}
}
</style>