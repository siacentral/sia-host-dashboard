<template>
	<div :class="alertClasses">
		<div class="alert-icon"><icon :icon="icon" /></div>
		<div>{{ text }}</div>
	</div>
</template>

<script>
export default {
	props: {
		text: String,
		icon: String,
		severity: String
	},
	computed: {
		alertClasses() {
			const c = { 'panel': true, 'alert-panel': true };

			c[`alert-${this.severity || 'severe'}`] = true;

			return c;
		}
	}
};
</script>

<style lang="stylus" scoped>
.alert-panel {
	display: grid;
	grid-template-columns: 32px minmax(0, 1fr);
	grid-gap: 15px;
	align-content: center;
	padding: 15px;
	margin-bottom: 5px;
	border: 2px solid dark-gray;
	color: rgba(255, 255, 255, 0.84);

	&.alert-severe {
		background: negative-dark;
		border-color: negative-accent;
	}

	&.alert-warning {
		background: warning-dark;
		border-color: warning-accent;
	}

	.alert-icon {
		text-align: center;
	}
}
</style>