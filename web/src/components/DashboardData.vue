<template>
<div class="data-points">
	<data-panel
		:title="dayStr"
		:earnedRevenue="day.earned_revenue"
		:potentialRevenue="day.potential_revenue"
		:burntCollateral="day.burnt_collateral"
		:contracts="day.new_contracts"
		contractLabel="New Contracts"
		:successfulContracts="day.successful_contracts"
		:failedContracts="day.failed_contracts" />
	<data-panel
		:title="monthStr"
		:earnedRevenue="month.earned_revenue"
		:potentialRevenue="month.potential_revenue"
		:burntCollateral="month.burnt_collateral"
		:contracts="month.new_contracts"
		contractLabel="New Contracts"
		:successfulContracts="month.successful_contracts"
		:failedContracts="month.failed_contracts" />
	<data-panel
		:title="yearStr"
		:earnedRevenue="year.earned_revenue"
		:potentialRevenue="year.potential_revenue"
		:burntCollateral="year.burnt_collateral"
		:contracts="year.new_contracts"
		contractLabel="New Contracts"
		:successfulContracts="year.successful_contracts"
		:failedContracts="year.failed_contracts" />
	<data-panel
		title="Total"
		:earnedRevenue="total.earned_revenue"
		:potentialRevenue="total.potential_revenue"
		:burntCollateral="total.burnt_collateral"
		:contracts="total.active_contracts"
		contractLabel="Active Contracts"
		:successfulContracts="total.successful_contracts"
		:failedContracts="total.failed_contracts" />
</div>
</template>

<script>
import { formatDate, formatMonth } from '@/utils/format';

import DataPanel from '@/components/DataPanel';

export default {
	components: {
		DataPanel
	},
	props: {
		month: Object,
		year: Object,
		day: Object,
		total: Object
	},
	computed: {
		yearStr() {
			if (!this.year)
				return new Date().getFullYear().toString();

			return new Date(this.year.timestamp).getFullYear().toString();
		},
		monthStr() {
			if (!this.month)
				return formatMonth(new Date());

			return formatMonth(new Date(this.month.timestamp));
		},
		dayStr() {
			if (!this.day)
				return formatDate(new Date());

			return formatDate(new Date(this.day.timestamp));
		}
	}
};
</script>

<style lang="stylus" scoped>
.data-points {
	display: grid;
	grid-gap: 15px;
	margin-bottom: 15px;

	@media screen and (min-width: 600px) {
		grid-template-columns: repeat(2, minmax(0, 1fr));
	}

	@media screen and (min-width: 850px) {
		grid-template-columns: repeat(4, minmax(0, 1fr));
	}
}
</style>