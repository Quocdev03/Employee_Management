<template>
	<div class="modal" v-if="visible" @click.self="$emit('close')">
		<component 
			:is="isForm ? 'form' : 'div'" 
			class="modal__dialog" 
			:class="dialogClass" 
			@submit.prevent="isForm ? $emit('submit', $event) : undefined"
		>
			<!-- Header -->
			<header class="modal__header" v-if="$slots.header || title">
				<slot name="header">
					<div>
						<h2 class="modal__title">{{ title }}</h2>
						<p class="modal__subtitle" v-if="subtitle">{{ subtitle }}</p>
					</div>
					<button type="button" class="btn-close" @click="$emit('close')">
						<img :src="closeIcon" alt="close" />
					</button>
				</slot>
			</header>

			<!-- Body -->
			<div class="modal__body" :class="bodyClass">
				<slot></slot>
			</div>

			<!-- Footer -->
			<footer class="modal__footer" v-if="$slots.footer">
				<slot name="footer"></slot>
			</footer>
		</component>
	</div>
</template>

<script setup>
import closeIcon from "@/assets/icons/close.svg";

defineProps({
	visible: {
		type: Boolean,
		default: false,
	},
	isForm: {
		type: Boolean,
		default: false,
	},
	title: {
		type: String,
		default: "",
	},
	subtitle: {
		type: String,
		default: "",
	},
	dialogClass: {
		type: String,
		default: "",
	},
	bodyClass: {
		type: String,
		default: "",
	},
});

defineEmits(["close", "submit"]);
</script>

<style scoped>
/* Scoped styles removed as they are now in common.css */
</style>
