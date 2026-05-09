<template>
	<div class="app-layout">
		<AppSidebar />
		<div :class="['app-body', { collapsed: ui.isCollapsed }]">
			<AppHeader />
			<main class="app-main">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<script setup>
import AppSidebar from "@/components/AppSidebar.vue";
import AppHeader from "@/components/AppHeader.vue";
import { useUIStore } from "@/stores/ui";

const ui = useUIStore();
</script>
<style scoped>
.app-layout {
	min-height: 100vh;
}

.app-body {
	margin-left: 240px;
	display: flex;
	flex-direction: column;
	min-height: 100vh;
	transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.app-body.collapsed {
	margin-left: 80px;
}

.app-main {
	flex: 1;
	overflow-y: auto;
	padding: 1rem;
	background: #f9fafb;
}

@media (max-width: 1023px) {
	.app-main {
		padding: 0.6rem;
	}
	.app-body,
	.app-body.collapsed {
		margin-left: 0;
	}
}
</style>
