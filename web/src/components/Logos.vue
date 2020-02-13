<template>
	<div class="logos">
		<a href="https://siacentral.com" _target="blank" class="logo"><sia-central /></a>
		<a href="https://sia.tech" _target="blank" class="built-with"><built-with-sia /></a>
	</div>
</template>

<script>
import SiaCentral from '@/assets/siacentral.svg';
import BuiltWithSia from '@/assets/built-with-sia.svg';

export default {
	components: {
		SiaCentral,
		BuiltWithSia
	}
};
</script>

<style lang="stylus" scoped>
.logos {
	display: flex;
	align-content: center;
	justify-content: center;

	a {
		display: inline-block;
		width: 120px;
		padding: 15px;
	}
}
</style>