<template>
	<BaseModal
		:visible="visible"
		:title="title"
		dialog-class="modal__dialog--confirm"
		@close="$emit('cancel')"
	>
		<template #default>
			<div class="confirm-icon">⚠️</div>
			<h3 class="modal__title text-center">{{ title }}</h3>
			<p class="modal__subtitle text-center">
				<slot></slot>
			</p>
			<div class="confirm-actions">
				<button
					class="btn btn--secondary flex-1"
					@click="$emit('cancel')"
				>
					{{ cancelText }}
				</button>
				<button
					class="btn btn--danger flex-1"
					@click="$emit('confirm')"
				>
					{{ confirmText }}
				</button>
			</div>
		</template>
	</BaseModal>
</template>

<script setup>
import BaseModal from "./BaseModal.vue";

defineProps({
	visible: {
		type: Boolean,
		required: true,
	},
	title: {
		type: String,
		default: "Xác nhận",
	},
	cancelText: {
		type: String,
		default: "Hủy bỏ",
	},
	confirmText: {
		type: String,
		default: "Đồng ý",
	},
});

defineEmits(["cancel", "confirm"]);
</script>

<style scoped>
:deep(.modal__header) {
	display: none;
}

:deep(.modal__body) {
	padding: 2.5rem 2rem;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 0.5rem;
}

:deep(.modal__dialog--confirm) {
	max-width: 400px;
}


.confirm-icon {
	font-size: 3.5rem;
	margin-bottom: 1rem;
	text-align: center;
}

.modal__title {
	font-size: 1.25rem;
	font-weight: 700;
	margin: 0 0 0.25rem 0;
	text-align: center;
}

.modal__subtitle {
	font-size: 0.85rem;
	color: var(--text-muted, #64748b);
	margin: 0;
	text-align: center;
}

.confirm-actions {
	display: flex;
	gap: 12px;
	margin-top: 1.5rem;
	width: 100%;
}
</style>
