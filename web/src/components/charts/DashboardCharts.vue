<template>
	<div class="charts">
		<revenue-chart :snapshots="snapshots" />
		<contract-chart :snapshots="snapshots" />
	</div>
</template>

<script>
import ContractChart from '@/components/charts/ContractChart';
import RevenueChart from '@/components/charts/RevenueChart';

export default {
	components: {
		ContractChart,
		RevenueChart
	},
	props: {
		snapshots: Array
	}
};
</script>

<style lang="stylus" scoped>
.charts {
	display: grid;
	grid-gap: 15px;
}

.charts {
	@media screen and (min-width: 767px) {
		grid-template-columns: repeat(2, minmax(0, 1fr));
	}
}
</style>