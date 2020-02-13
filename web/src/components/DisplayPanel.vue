<template>
	<div>
		<div :class="{ 'panel': true,  'icon-panel': (icon && icon.length > 0) }">
			<div v-if="icon" :class="iconClasses">
				<icon :icon="icon" />
			</div>
			<div class="panel-content">
				<slot />
			</div>
		</div>
		<div class="panel-extras">
			<slot name="extras" />
		</div>
	</div>
</template>

<script>
export default {
	props: {
		icon: String,
		severity: String,
		extras: Boolean
	},
	computed: {
		iconClasses() {
			const classes = { 'panel-icon': true };

			if (this.severity)
				classes[`icon-${this.severity}`] = true;

			return classes;
		}
	}
};
</script>

<style lang="stylus" scoped>
.panel {
	height: 100%;
	z-index: 2;
}

.panel-extras {
	position: relative;
	top: -10px;
	padding: 25px 15px 15px;
	z-index: 1;
	background: lighten(bg-accent, 5%);
    border-bottom-left-radius: 8px;
    border-bottom-right-radius: 8px;
	box-shadow: 0 2px 3px 1px rgba(0, 0, 0, 0.1);
	font-size: 0.9rem;
	color: rgba(255, 255, 255, 0.75);

	&:empty {
		padding: 0;
	}
}

.extra-value {
	width: 100%;
	overflow-x: auto;
}
</style>
