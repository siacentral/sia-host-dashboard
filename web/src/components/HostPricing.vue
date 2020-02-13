<template>
	<div class="panel">
		<div class="title">Pricing</div>
		<div class="data-points">
			<div class="data-panel">
				<div class="data-title">Contract Price</div>
				<div class="data-value" v-html="contractPriceSC" />
				<div class="data-secondary" v-html="contractPriceCurrency" />
			</div>
			<div class="data-panel">
				<div class="data-title">Storage Price</div>
				<div class="data-value" v-html="storagePriceSC" />
				<div class="data-secondary" v-html="storagePriceCurrency" />
			</div>
			<div class="data-panel">
				<div class="data-title">Download Price</div>
				<div class="data-value" v-html="downloadPriceSC" />
				<div class="data-secondary" v-html="downloadPriceCurrency" />
			</div>
			<div class="data-panel">
				<div class="data-title">Upload Price</div>
				<div class="data-value" v-html="uploadPriceSC" />
				<div class="data-secondary" v-html="uploadPriceCurrency" />
			</div>
		</div>
	</div>
</template>

<script>
import { mapState } from 'vuex';
import { formatPriceString, formatDataPriceString, formatMonthlyPriceString } from '@/utils/format';
import BigNumber from 'bignumber.js';

export default {
	props: {
		settings: Object
	},
	computed: {
		...mapState(['exchangeRateSC', 'currency']),
		contractPrice() {
			let val = new BigNumber(0);

			if (this.settings && this.settings.contract_price)
				val = new BigNumber(this.settings.contract_price);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		storagePrice() {
			let val = new BigNumber(0);

			if (this.settings && this.settings.storage_price)
				val = new BigNumber(this.settings.storage_price);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		uploadPrice() {
			let val = new BigNumber(0);

			if (this.settings && this.settings.upload_price)
				val = new BigNumber(this.settings.upload_price);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		downloadPrice() {
			let val = new BigNumber(0);

			if (this.settings && this.settings.download_price)
				val = new BigNumber(this.settings.download_price);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		contractPriceSC() {
			const format = formatPriceString(this.contractPrice);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		contractPriceCurrency() {
			const format = formatPriceString(this.contractPrice, 2, this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		storagePriceSC() {
			const format = formatMonthlyPriceString(this.storagePrice);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		storagePriceCurrency() {
			const format = formatMonthlyPriceString(this.storagePrice, 2, 'decimal', this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		uploadPriceSC() {
			const format = formatDataPriceString(this.uploadPrice);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		uploadPriceCurrency() {
			const format = formatDataPriceString(this.uploadPrice, 2, 'decimal', this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		downloadPriceSC() {
			const format = formatDataPriceString(this.downloadPrice);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		downloadPriceCurrency() {
			const format = formatDataPriceString(this.downloadPrice, 2, 'decimal', this.currency, this.exchangeRateSC[this.currency]);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		}
	}
};
</script>

<style lang="stylus" scoped>
.panel {
	padding: 15px;
}

.data-points {
	display: grid;
	grid-gap: 15px;

	@media screen and (min-width: 600px) {
		grid-template-columns: repeat(2, auto);
		justify-content: space-between;
	}

	@media screen and (min-width: 850px) {
		grid-template-columns: repeat(4, auto);
	}
}

.title {
	font-size: 1rem;
	margin-bottom: 15px;
    color: rgba(255,255,255,0.84);
}

.data-title {
	margin-bottom: 5px;
	color: rgba(255, 255, 255, 0.54);
	font-size: 0.9rem;
	text-align: right;
}

.data-value {
	color: primary;
	font-size: 1.1rem;
	text-align: right;
}

.data-secondary {
	color: rgba(255, 255, 255, 0.54);
	font-size: 1rem;
	margin-top: 2px;
	text-align: right;
}
</style>