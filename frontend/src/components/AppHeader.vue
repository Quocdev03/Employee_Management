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
				<div class="avatar-wrap">
					<!-- ✅ Bỏ console.log(), dùng me.user?.avatar -->
					<img
						v-if="me.user?.avatar"
						:src="me.user?.avatar"
						alt="avatar"
						class="avatar-img"
					/>
					<div v-else class="avatar-initials">
						{{ getInitials(me.user?.name) }}
					</div>
				</div>
				<div class="user-info">
					<span class="user-name">{{
						me.user?.name || "Người dùng"
					}}</span>
					<div class="user-meta">
						<span class="user-email">{{ me.user?.email }}</span>
						<span :class="['role-badge', me.user?.role]">
							{{
								me.user?.role === "admin"
									? "Quản trị"
									: "Nhân viên"
							}}
						</span>
					</div>
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
import { computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useMeStore } from "@/stores/me";
import { useUIStore } from "@/stores/ui";

// Import Icons
import menuIcon from "@/assets/icons/menu.svg";
import logoutIcon from "@/assets/icons/log-out.svg";

// --- Khởi tạo Store và Routing ---
const auth = useAuthStore();
const me = useMeStore();
const ui = useUIStore();
const route = useRoute();
const router = useRouter();

// --- Cấu hình Tiêu đề Trang ---
const PAGE_TITLES = {
	"/dashboard": "Tổng quan",
	"/employees": "Nhân viên",
	"/profile": "Hồ sơ",
	"/users": "Người dùng",
	"/departments": "Phòng Ban",
};

const pageTitle = computed(() => {
	return PAGE_TITLES[route.path] || "HR System";
});

// --- Lifecycle ---
onMounted(() => {
	// Đảm bảo thông tin Me luôn sẵn sàng để hiển thị avatar/name
	if (!me.profile) {
		me.fetchProfile().catch(() => {});
	}
});

// --- Xử lý sự kiện ---

/**
 * Đăng xuất người dùng
 */
function handleLogout() {
	auth.logout();
	router.push("/login");
}

/**
 * Lấy chữ cái đầu (Initials) để làm avatar fallback
 */
function getInitials(name) {
	if (!name || typeof name !== "string") return "?";

	const parts = name.trim().split(" ");
	if (parts.length === 0) return "?";

	return parts.pop().charAt(0).toUpperCase();
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

.avatar-wrap {
	width: 40px;
	height: 40px;
	border-radius: 10px;
	overflow: hidden;
	background: #eff6ff;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 1px solid #dbeafe;
}

.avatar-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.avatar-initials {
	font-weight: 700;
	color: #3b82f6;
	font-size: 1rem;
}

.user-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.user-name {
	font-size: 0.95rem;
	font-weight: 700;
	color: #1e293b;
	line-height: 1;
}

.user-meta {
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.user-email {
	font-size: 0.75rem;
	font-weight: 500;
	color: #64748b;
	line-height: 1;
}

.role-badge {
	display: inline-flex;
	align-items: center;
	padding: 0 0.5rem;
	height: 18px;
	border-radius: 9999px;
	font-size: 0.65rem;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	width: fit-content;
}

.role-badge.admin {
	background: #f3e8ff;
	color: #7e22ce;
}

.role-badge.user {
	background: #f0fdf4;
	color: #15803d;
}

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
