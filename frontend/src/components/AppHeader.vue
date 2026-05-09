<template>
	<header class="app-header">
		<!-- ===== Hamburger Mobile ===== -->
		<button class="hamburger-btn" @click="ui.toggleMobileSidebar()">
			<img :src="menuIcon" alt="menu" />
		</button>

		<!-- ===== Tiêu đề trang ===== -->
		<div class="header-title">{{ pageTitle }}</div>

		<!-- ===== Thông tin User & Đăng xuất ===== -->
		<div class="header-right">
			<div class="user-profile">
				<div class="avatar">
					<img
						v-if="auth.user?.avatar"
						:src="auth.user?.avatar"
						alt="avatar"
						class="avatar-circle"
					/>
				</div>
				<div class="user-info">
					<span class="user-email">{{ auth.user?.email }}</span>
					<span :class="['role-badge', auth.user?.role]">{{
						auth.user?.role
					}}</span>
				</div>
			</div>
			<div class="divider"></div>
			<button class="btn-logout" @click="handleLogout" title="Đăng xuất">
				<img :src="logoutIcon" alt="logout" />
				<span class="logout-text">Đăng xuất</span>
			</button>
		</div>
	</header>
</template>

<script setup>
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useUIStore } from "@/stores/ui";

// Import Icons
import menuIcon from "@/assets/icons/menu.svg";
import logoutIcon from "@/assets/icons/log-out.svg";

// --- Khởi tạo Store và Routing ---
const auth = useAuthStore();
const ui = useUIStore();
const route = useRoute();
const router = useRouter();

// --- Cấu hình Tiêu đề Trang ---

// Bản đồ mapping giữa đường dẫn và tiêu đề trang tương ứng
const PAGE_TITLES = {
	"/dashboard": "Tổng quan",
	"/employees": "Nhân viên",
	"/profile": "Hồ sơ",
	"/users": "Người dùng",
	"/departments": "Phòng Ban",
};

/**
 * Tự động lấy tiêu đề trang dựa trên đường dẫn hiện tại
 */
const pageTitle = computed(() => {
	return PAGE_TITLES[route.path] || "HR System";
});

// --- Xử lý sự kiện ---

/**
 * Thực hiện đăng xuất và điều hướng người dùng về trang đăng nhập
 */
function handleLogout() {
	auth.logout();
	router.push("/login");
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Layout chính ===== */
.app-header {
	font-family: "Inter", sans-serif;
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 1.2rem 1.5rem;
	background: rgba(255, 255, 255, 0.98);
	border-bottom: 1px solid #e2e8f0;
	box-shadow:
		0 4px 6px -1px rgba(0, 0, 0, 0.02),
		0 2px 4px -1px rgba(0, 0, 0, 0.01);
	position: sticky;
	top: 0;
	z-index: 10;
}

/* ===== Tiêu đề ===== */
.header-title {
	font-weight: 700;
	font-size: 1.25rem;
	color: #1e293b;
	letter-spacing: -0.01em;
}

/* ===== Nhóm bên phải ===== */
.header-right {
	display: flex;
	align-items: center;
	flex-wrap: wrap;
	gap: 1.5rem;
}

/* ===== User Profile ===== */
.user-profile {
	display: flex;
	align-items: center;
	gap: 0.875rem;
}

.avatar-circle {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 42px;
	height: 42px;
	background-color: #eff6ff;
	color: #3b82f6;
	border-radius: 50%;
}

.user-info {
	display: flex;
	flex-direction: column;
	justify-content: center;
	gap: 0.2rem;
}

.user-email {
	font-size: 0.875rem;
	font-weight: 600;
	color: #334155;
}

.role-badge {
	display: inline-flex;
	padding: 0.15rem 0.5rem;
	border-radius: 9999px;
	font-size: 0.7rem;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	width: fit-content;
}

.role-badge.admin {
	background: #f3e8ff;
	color: #7e22ce;
}

.role-badge.user,
.role-badge.employee {
	background: #e0f2fe;
	color: #0369a1;
}

/* ===== Nút Đăng xuất ===== */
.divider {
	width: 1px;
	height: 32px;
	background-color: #e2e8f0;
}

.btn-logout {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	padding: 0.6rem 1rem;
	border: 1px solid #e2e8f0;
	border-radius: 10px;
	background: #ffffff;
	cursor: pointer;
	font-size: 0.875rem;
	font-weight: 500;
	color: #64748b;
	transition: all 0.2s ease;
}

.btn-logout:hover {
	background: #fef2f2;
	color: #ef4444;
	border-color: #fecaca;
}

.btn-logout img {
	width: 18px;
	height: 18px;
}

/* ===== Mobile Button ===== */
.hamburger-btn {
	display: none;
}
.hamburger-btn img {
	width: 22px;
	height: 22px;
}

/* ===== Responsive ===== */
@media (max-width: 1023px) {
	.user-email,
	.divider,
	.logout-text {
		display: none;
	}

	.btn-logout {
		padding: 0.6rem;
	}

	.hamburger-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 38px;
		height: 38px;
		margin-right: 1rem;
		background: #f1f5f9;
		border: 1px solid #e2e8f0;
		border-radius: 10px;
		cursor: pointer;
		color: #64748b;
		flex-shrink: 0;
	}

	.hamburger-btn:hover {
		background: #e2e8f0;
		color: #1e293b;
	}
}

@media (max-width: 900px) {
	.header-title {
		display: none;
	}
	.header-right {
		width: 100%;
		justify-content: space-between;
	}
}
</style>
