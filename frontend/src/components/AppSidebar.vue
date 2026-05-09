<template>
	<!-- ===== Overlay Mobile ===== -->
	<div
		class="sidebar-overlay"
		v-if="ui.mobileOpen"
		@click="ui.closeMobileSidebar()"
	></div>

	<!-- ===== Main Sidebar ===== -->
	<aside
		:class="[
			'sidebar',
			{ collapsed: ui.isCollapsed, 'mobile-open': ui.mobileOpen },
		]"
	>
		<!-- Logo Section -->
		<div class="sidebar-logo">
			<div class="logo-icon">
				<img :src="usersIcon" alt="logo" class="icon-filter-primary" />
			</div>
			<span class="logo-text" v-if="!ui.isCollapsed">HR System</span>
		</div>

		<!-- Navigation Links -->
		<nav class="sidebar-nav">
			<RouterLink
				v-for="item in menuItems"
				:key="item.path"
				:to="item.path"
				class="nav-item"
				:title="ui.isCollapsed ? item.label : ''"
				@click="ui.closeMobileSidebar()"
			>
				<img :src="item.icon" class="nav-icon" />
				<span class="nav-label" v-if="!ui.isCollapsed">{{
					item.label
				}}</span>
			</RouterLink>
		</nav>

		<!-- Toggle Collapse Button -->
		<button
			class="collapse-btn"
			@click="ui.toggleCollapse()"
			:title="ui.isCollapsed ? 'Mở rộng' : 'Thu gọn'"
		>
			<img :src="ui.isCollapsed ? nextIcon : prevIcon" alt="toggle" />
		</button>
	</aside>
</template>

<script setup>
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useUIStore } from "@/stores/ui";

// Import Icons
import gridIcon from "@/assets/icons/grid.svg";
import usersIcon from "@/assets/icons/users.svg";
import keyIcon from "@/assets/icons/key.svg";
import userIcon from "@/assets/icons/user.svg";
import prevIcon from "@/assets/icons/chevron-left.svg";
import nextIcon from "@/assets/icons/chevron-right.svg";
import deptIcon from "@/assets/icons/department.svg";

// --- Khởi tạo Store ---
const ui = useUIStore();
const auth = useAuthStore();

// --- Cấu hình Menu ---

/**
 * Kiểm tra xem người dùng hiện tại có quyền admin hay không
 */
const isAdmin = computed(() => {
	return auth.user?.role === "admin";
});

/**
 * Danh sách các menu item hiển thị trên sidebar
 * Sẽ hiển thị thêm menu "Tài khoản" nếu người dùng là admin
 */
const menuItems = computed(() => {
	const baseMenu = [
		{ path: "/dashboard", icon: gridIcon, label: "Dashboard" },
		{ path: "/employees", icon: usersIcon, label: "Nhân viên" },
		{ path: "/profile", icon: userIcon, label: "Hồ sơ" },
		{ path: "/departments", icon: deptIcon, label: "Phòng Ban" },
	];

	if (isAdmin.value) {
		baseMenu.push({ path: "/users", icon: keyIcon, label: "Tài khoản" });
	}

	return baseMenu;
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Constants/Utilities ===== */
.icon-filter-primary {
	width: 20px;
	height: 20px;
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%); /* #3b82f6 */
}

/* ===== Base Layout ===== */
.sidebar {
	position: fixed;
	width: 240px;
	height: 100vh;
	top: 0;
	left: 0;
	overflow: auto;
	font-family: "Inter", sans-serif;
	background: #ffffff;
	display: flex;
	flex-direction: column;
	transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed {
	width: 80px;
}

.sidebar.collapsed ~ :deep(.app-body) {
	margin-left: 80px;
}

/* ===== Logo Section ===== */
.sidebar-logo {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 12px;
	padding: 1.5rem;
	height: 80px;
	box-sizing: border-box;
}

.sidebar:not(.collapsed) .sidebar-logo {
	justify-content: flex-start;
}

.logo-icon {
	width: 40px;
	height: 40px;
	background: #eff6ff;
	color: #3b82f6;
	border-radius: 12px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}

.logo-text {
	font-weight: 700;
	font-size: 1.25rem;
	color: #1e293b;
	white-space: nowrap;
	letter-spacing: -0.01em;
}

/* ===== Navigation Menu ===== */
.sidebar-nav {
	flex: 1;
	padding: 1rem;
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.nav-item {
	display: flex;
	align-items: center;
	gap: 14px;
	padding: 0.875rem 1rem;
	text-decoration: none;
	color: #64748b;
	border-radius: 12px;
	transition: all 0.2s ease;
	font-weight: 500;
	white-space: nowrap;
	overflow: hidden;
}

.sidebar.collapsed .nav-item {
	padding: 0.875rem;
	justify-content: center;
}

.nav-item:hover {
	background: #f8fafc;
	color: #1e293b;
}

.nav-item.router-link-active {
	background: #eff6ff;
	color: #3b82f6;
	font-weight: 600;
}

/* Icon Styling */
.nav-icon {
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	width: 24px;
	height: 24px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg)
		brightness(88%) contrast(89%); /* #94a3b8 */
	transition: all 0.2s ease;
}

.nav-item:hover .nav-icon {
	filter: invert(47%) sepia(13%) saturate(545%) hue-rotate(183deg)
		brightness(95%) contrast(87%); /* #64748b */
}

.nav-item.router-link-active .nav-icon {
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%); /* #3b82f6 */
}

.nav-label {
	font-size: 0.95rem;
}

/* ===== Toggle Button ===== */
.collapse-btn {
	position: absolute;
	bottom: 2rem;
	left: 50%;
	transform: translateX(-50%);
	width: calc(100% - 2rem);
	padding: 0.875rem;
	background: #ffffff;
	border: 1px solid #e2e8f0;
	border-radius: 12px;
	color: #64748b;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s ease;
	box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
}

.sidebar.collapsed .collapse-btn {
	width: 48px;
	padding: 0.875rem 0;
}

.collapse-btn:hover {
	background: #f8fafc;
	color: #1e293b;
	border-color: #cbd5e1;
}

.collapse-btn img {
	width: 20px;
	height: 20px;
	filter: invert(47%) sepia(13%) saturate(545%) hue-rotate(183deg)
		brightness(95%) contrast(87%);
}

/* ===== Responsive Mobile ===== */
@media (max-width: 1023px) {
	.sidebar {
		position: fixed;
		top: 0;
		left: 0;
		height: 100vh;
		z-index: 20;
		width: 260px !important;
		transform: translateX(-100%);
		transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}

	.sidebar.mobile-open {
		transform: translateX(0);
	}

	.sidebar-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.25);
		z-index: 15;
		backdrop-filter: blur(1px);
	}

	.collapse-btn {
		display: none;
	}
}
</style>
