<template>
	<div class="data-chart">
		<div class="chart-title">{{ title }}</div>
		<div class="chart-selected">{{ selectedLabel }}</div>
		<div class="chart-sub">
			<div v-for="(v, i) in text" :key="i">

			</div>
		</div>
		<div class="chart">
			<stacked-chart :nodes="nodes" :colors="colors" :fills="fills" @selected="onSetSelected" />
		</div>
		<div class="chart-label-left">{{ labels[0] }}</div>
		<div class="chart-label-right">{{ labels[labels.length - 1] }}</div>
	</div>
</template>

<script>
import StackedChart from '@/components/StackedChart';

export default {
	components: {
		StackedChart
	},
	props: {
		title: String,
		nodes: Array,
		colors: Array,
		fills: Array,
		labels: Array,
		text: Array
	},
	data() {
		return {
			active: -1
		};
	},
	computed: {
		selectedLabel() {
			if (this.active === -1 || !this.labels[this.active - 1])
				return `${this.labels[0]} - ${this.labels[this.labels.length - 1]}`;

			return this.labels[this.active - 1];
		}
	},
	methods: {
		onSetSelected(i) {
			try {
				console.log(i);
				this.active = i;
			} catch (ex) {
				console.error('ChartDisplay.onSetSelected', ex);
			}
		}
	}
};
</script>

<style lang="stylus" scoped>
.data-chart {
	position: relative;
	display: grid;
	grid-template-columns: repeat(2, min-content);
	justify-content: space-between;
	background: bg-accent;
	border-radius: 8px;
	overflow: hidden;
	box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
	grid-gap: 15px;
	padding: 15px 15px 0;

	.chart-title, .chart-selected {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.84);
		white-space: nowrap;
	}

	.chart-sub {
		font-size: 1.2rem;
		color: primary;
	}

	.chart {
		margin: 0 -15px -5px;
		align-self: bottom;
	}

	.chart, .chart-sub, .chart-controls {
		grid-column: 1 / -1;
	}

	.chart-label-left, .chart-label-right {
		position: absolute;
		bottom: 10px;
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.84);
		z-index: 2;
	}

	.chart-label-left {
		left: 10px;
	}

	.chart-label-right {
		right: 10px;
	}

	.chart-data-controls {
		grid-column: 1 / -1;

		button {
			display: inline-block;
			padding: 4px 8px;
			background: none;
			color: rgba(255, 255, 255, 0.54);
			border: 1px solid #999;
			border-radius: 4px;
			outline: none;
			margin-right: 8px;
			font-size: 0.8rem;
			cursor: pointer;
			transition: all 0.3s linear;

			&.chart-selected {
				color: primary;
				border-color: primary;
			}
		}
	}
}
</style>