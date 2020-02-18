<template>
	<div class="host-stats">
		<div class="split">
			<div class="panel">
				<div class="title">Uptime</div>
				<div class="uptime-counter">
					<div v-html="uptimeStr" />
				</div>
			</div>
			<div class="panel">
				<div class="title">Storage Usage</div>
				<div class="data-points">
					<div class="data-panel">
						<div class="data-value" v-html="usedByteStr" />
						<div class="data-title">Used</div>
					</div>
					<div class="data-panel">
						<div class="data-value" v-html="totalByteStr" />
						<div class="data-title">Total</div>
					</div>
				</div>
			</div>
			<div class="panel">
				<div class="title">Bandwidth Usage (Last 30 Days)</div>
				<div class="data-points">
					<div class="data-panel">
						<div class="data-value" v-html="uploadBytesStr" />
						<div class="data-title">Sent</div>
					</div>
					<div class="data-panel">
						<div class="data-value" v-html="downloadBytesStr" />
						<div class="data-title">Received</div>
					</div>
				</div>
			</div>
		</div>
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
	</div>
</template>

<script>
import { mapState } from 'vuex';
import { formatPriceString, formatDataPriceString, formatMonthlyPriceString, formatByteString, formatDuration } from '@/utils/format';
import BigNumber from 'bignumber.js';

export default {
	props: {
		settings: Object,
		status: Object
	},
	computed: {
		...mapState(['exchangeRateSC', 'currency']),
		usedStorageBytes() {
			let val = new BigNumber(0);

			if (this.status && this.status.used_storage)
				val = new BigNumber(this.status.used_storage);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		totalStorageBytes() {
			let val = new BigNumber(0);

			if (this.status && this.status.total_storage)
				val = new BigNumber(this.status.total_storage);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		uploadBytes() {
			let val = new BigNumber(0);

			if (this.status && this.status.upload_bandwidth)
				val = new BigNumber(this.status.upload_bandwidth);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
		downloadBytes() {
			let val = new BigNumber(0);

			if (this.status && this.status.download_bandwidth)
				val = new BigNumber(this.status.download_bandwidth);

			if (!val.isFinite() || val.isNaN())
				val = new BigNumber(0);

			return val;
		},
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
		uptimeStr() {
			let v = 0;

			if (this.status && typeof this.status.start_time === 'string')
				v = (Date.now() - new Date(this.status.start_time).getTime()) / 1000;

			const format = formatDuration(v, true);

			return format.map(f => `${f.value} <span class="currency-display">${f.label}</span>`).join(' ');
		},
		usedByteStr() {
			const format = formatByteString(this.usedStorageBytes, 'decimal', 2);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		totalByteStr() {
			const format = formatByteString(this.totalStorageBytes, 'decimal', 2);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		uploadBytesStr() {
			const format = formatByteString(this.uploadBytes, 'decimal', 2);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
		},
		downloadBytesStr() {
			const format = formatByteString(this.downloadBytes, 'decimal', 2);

			return `${format.value} <span class="currency-display">${format.label}</span>`;
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
.host-stats {
	margin-bottom: 15px;
}

.panel {
	padding: 15px;
}

.data-points {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	grid-gap: 15px;
	justify-items: center;

	@media screen and (min-width: 600px) {
		grid-gap: 30px;
	}

	@media screen and (min-width: 850px) {
		grid-template-columns: repeat(4, minmax(0, 1fr));
	}
}

.split {
	display: grid;
	grid-gap: 15px;
	margin-bottom: 15px;

	.data-value {
		font-size: 1.5rem;
	}

	@media screen and (min-width: 757px) {
		grid-template-columns: auto repeat(2, minmax(0, 1fr));
	}

	@media screen and (min-width: 850px) {
		.data-points {
			grid-template-columns: repeat(2, minmax(0, 1fr));
			justify-items: center;
			grid-gap: 30px;
		}
	}
}

.uptime-counter {
	display: grid;
	font-size: 1.5rem;
	color: primary;
	align-content: center;
	justify-content: center;
}

.title {
	font-size: 1rem;
	margin-bottom: 15px;
    color: rgba(255,255,255,0.84);
	text-align: center;
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