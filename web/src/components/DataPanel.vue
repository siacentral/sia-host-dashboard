<template>
	<div class="panel data-panel">
		<div class="data-title">{{ title }}</div>
		<div class="data-point">
			<div class="point-title">Earned Revenue</div>
			<div class="point-value" v-html="earnedRevenueSC" />
			<div class="point-value point-secondary" v-html="earnedRevenueCurrency" />
		</div>
		<div class="data-point">
			<div class="point-title">Potential Revenue</div>
			<div class="point-value" v-html="potentialRevenueSC" />
			<div class="point-value point-secondary" v-html="potentialRevenueCurrency" />
		</div>
		<div class="data-point">
			<div class="point-title">New Contracts</div>
			<div class="point-value" v-html="newContractsStr" />
		</div>
		<div class="data-point">
			<div class="point-title">Successful Contracts</div>
			<div class="point-value" v-html="successfulContractsStr" />
		</div>
	</div>
</template>

<script>
import { mapState } from 'vuex';
import { formatPriceString, formatNumber } from '@/utils/format';
import BigNumber from 'bignumber.js';

export default {
	props: {
		title: String,
		earnedRevenue: String,
		potentialRevenue: String,
		newContracts: Number,
		successfulContracts: Number
	},
	computed: {
		...mapState(['exchangeRateSC', 'currency']),
		earnedRevenueSC() {
			let val = new BigNumber(this.earnedRevenue);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			const format = formatPriceString(val);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		potentialRevenueSC() {
			let val = new BigNumber(this.potentialRevenue);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			const format = formatPriceString(val);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		earnedRevenueCurrency() {
			let val = new BigNumber(this.earnedRevenue);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			const format = formatPriceString(val, 2, this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		potentialRevenueCurrency() {
			let val = new BigNumber(this.potentialRevenue);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			const format = formatPriceString(val, 2, this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		newContractsStr() {
			let val = new BigNumber(this.newContracts);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return `${formatNumber(val)}`;
		},
		successfulContractsStr() {
			let val = new BigNumber(this.successfulContracts);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return `${formatNumber(val)}`;
		}
	}
};
</script>

<style lang="stylus" scoped>
.data-panel {
	padding: 15px;
}

.data-title, .point-title, .point-value {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.data-point {
	margin-bottom: 15px;

	&:last-of-type {
		margin-bottom: 0;
	}
}

.data-title {
	margin-bottom: 15px;
	text-align: center;
	color: rgba(255, 255, 255, 0.84);
}

.point-title {
	margin-bottom: 5px;
	color: rgba(255, 255, 255, 0.54);
	font-size: 0.9rem;
	text-align: right;
}

.point-value {
	color: primary;
	font-size: 1.1rem;
	text-align: right;
}

.point-secondary {
	color: rgba(255, 255, 255, 0.54);
	font-size: 1rem;
	margin-top: 2px;
}
</style>