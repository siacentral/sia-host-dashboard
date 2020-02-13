<template>
	<div class="data-chart">
		<div class="chart-title">{{ title }}</div>
		<div class="chart-selected">{{ selectedTimestamp }}</div>
		<div class="chart-data-controls">
			<button :class="{ 'chart-selected': selected === i }" @click="selected = i" v-for="(data, i) in charts" :key="i">{{ data.name }}</button>
		</div>
		<div class="chart-sub">{{ selectedLabel }}</div>
		<svg class="chart" :viewBox="viewBox">
			<defs>
				<clipPath :id="clipID">
					<path :d="fillPath" fill="none" stroke="white" stroke-width="6"></path>
				</clipPath>
			</defs>
			<path class="line-fill" :d="fillPath" stroke-width="3"></path>
			<!--<path class="line" :d="linePath" fill="none" stroke-width="3"></path>-->
			<rect :class="{ 'line-hover': true, 'active': isActive(i) }" :clip-path="rectClipPath" v-for="(point, i) in coordinates" :key="i" stroke-width="6" :x="dist * (i - 1)" :y="0" :width="dist" :height="height"></rect>
			<circle :key="`line-circle-${i}`" class="point" v-for="(point, i) in coordinates" :cx="point.coord[0]" :cy="point.coord[1]" r="4" stroke-width="2"></circle>
			<rect :key="`line-rect-${i}`" v-for="(point, i) in coordinates" fill="transparent" :x="dist * (i- 1)" :y="0" :width="dist" :height="height" @mouseover="setActive(i)" @mouseout="setActive(-1)"></rect>
		</svg>
		<div class="chart-label-left">{{ lines[0].timestamp }}</div>
		<div class="chart-label-right">{{ lines[lines.length - 1].timestamp }}</div>
	</div>
</template>

<script>
import { BigNumber } from 'bignumber.js';
import { line, curveCardinal } from 'd3-shape';

export default {
	name: 'line-chart',
	props: {
		title: String,
		charts: Array
	},
	data() {
		return {
			width: 500,
			height: 200,
			active: -1,
			selected: 0
		};
	},
	methods: {
		setActive(n) {
			this.active = n;
		},
		isActive(n) {
			return this.active === n;
		}
	},
	computed: {
		lines() {
			return this.charts[this.selected].data;
		},
		viewBox() {
			return `0 0 ${this.width} ${this.height}`;
		},
		dist() {
			return this.width / (this.lines.length || 1);
		},
		clipID() {
			return `clip-path-${this._uid}`;
		},
		rectClipPath() {
			return `url(#${this.clipID})`;
		},
		maxValue() {
			let max;

			for (let i = 0; i < this.lines.length; i++) {
				if (!max || max.isLessThan(this.lines[i].value))
					max = this.lines[i].value;
			}

			return max || new BigNumber(1);
		},
		minValue() {
			let min;

			for (let i = 0; i < this.lines.length; i++) {
				if (!min || min.isGreaterThan(this.lines[i].value))
					min = this.lines[i].value;
			}

			return min || new BigNumber(1);
		},
		coordinates() {
			const min = new BigNumber(this.minValue.toFixed(0, 3)),
				max = new BigNumber(this.maxValue.toFixed(0, 2));

			let yr = new BigNumber(max).minus(min).dividedBy(this.height * 0.5).toNumber(),
				xr = this.width / (this.lines.length || 1);

			if (yr === 0)
				yr = 1;

			const values = this.lines.map((v, i) => {
				const x = (xr * i) + (xr / 2),
					y = new BigNumber(this.height * 0.75).minus(v.value.minus(min).dividedBy(yr)).toNumber();

				v.coord = [x, y];

				return v;
			});

			if (values.length === 0)
				return values;

			const first = values[0].coord[1],
				last = values[values.length - 1].coord[1];

			values.unshift({ value: new BigNumber(0), label: '', coord: [(xr / 2) * -1, first] });
			values.push({ value: new BigNumber(0), label: '', coord: [this.width + (xr / 2), last] });

			return values;
		},
		linePath() {
			const generator = line()
				.x(d => d.coord[0])
				.y(d => d.coord[1])
				.curve(curveCardinal);

			return generator(this.coordinates);
		},
		fillPath() {
			return `${this.linePath} L ${this.coordinates[this.coordinates.length - 1].coord[0]}, ${this.height} L ${this.coordinates[0].coord[0]}, ${this.height} z`;
		},
		startPoint() {
			if (this.coordinates.length === 0)
				return { label: '' };

			return this.coordinates[0][1];
		},
		endPoint() {
			if (this.coordinates.length === 0)
				return { label: '' };

			const coords = this.coordinates[0];

			return coords[coords.length - 2];
		},
		selectedTimestamp() {
			if (this.active === -1 || !this.lines[this.active - 1])
				return this.lines[this.lines.length - 1].timestamp;

			return this.lines[this.active - 1].timestamp;
		},
		selectedLabel() {
			if (this.active === -1 || !this.lines[this.active - 1])
				return this.lines[this.lines.length - 1].label;

			if (this.active === 1 && this.lines[0])
				return this.lines[0].label;

			return `${this.lines[this.active - 1].label} (${this.lines[this.active - 1].delta})`;
		}
	}
};
</script>

<style lang="stylus" scoped>
svg {
	display: block;
}

.data-chart {
	position: relative;
	display: grid;
	grid-template-columns: repeat(2, min-content);
	justify-content: space-between;
	background: bg-dark-accent;
	border-radius: 8px;
	overflow: hidden;
	box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
	grid-gap: 10px;
	padding: 10px 10px 0;

	.chart-title, .chart-selected {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.54);
		white-space: nowrap;
	}

	.chart-sub {
		font-size: 1.2rem;
		color: primary;
	}

	.chart {
		margin: 0 -10px;
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

.line-hover {
	fill: primary;
	opacity: 0;
	transition: opacity 0.3s linear;
}

&.active {
	opacity: 1;
}

.line {
	stroke: primary;
}

.point {
	fill: primary;
	stroke: bg-dark-accent;
}

.line-fill {
	fill: primary;
	opacity: 0.4;
}
</style>
