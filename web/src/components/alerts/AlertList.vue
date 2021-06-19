<template>
	<div class="alert-list">
		<alert v-for="(alert, i) in alerts"
			:key="i"
			:text="alert.text"
			:severity="alert.severity"
			:icon="alertIcon(alert.type)" />
	</div>
</template>

<script>
import Alert from '@/components/alerts/Alert';

export default {
	components: {
		Alert
	},
	props: {
		alerts: Array
	},
	methods: {
		alertIcon(type) {
			if (typeof type !== 'string')
				return 'exclamation';

			switch (type.toLowerCase()) {
			case 'storage':
				return 'hdd';
			default:
				return 'exclamation';
			}
		}
	}
};
</script>

<style lang="stylus" scoped>
.alert-list {
	margin-bottom: 15px;
}
</style>