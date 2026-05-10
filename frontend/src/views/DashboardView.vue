<template>
	<div class="dashboard-container">
		<!-- ===== Thống kê tổng quan ===== -->
		<div class="stats-grid">
			<div class="stat-card stat-card--blue">
				<div class="stat-icon-wrapper">
					<img :src="usersIcon" class="stat-icon" alt="users" />
				</div>
				<div class="stat-info">
					<div class="stat-value">
						{{ dashboardStats.totalEmployees }}
					</div>
					<div class="stat-label">Tổng nhân viên</div>
				</div>
			</div>

			<div class="stat-card stat-card--green">
				<div class="stat-icon-wrapper">
					<img :src="checkIcon" class="stat-icon" alt="active" />
				</div>
				<div class="stat-info">
					<div class="stat-value">
						{{ dashboardStats.activeEmployees }}
					</div>
					<div class="stat-label">Đang làm việc</div>
				</div>
			</div>

			<div class="stat-card stat-card--amber">
				<div class="stat-icon-wrapper">
					<img :src="buildingIcon" class="stat-icon" alt="dept" />
				</div>
				<div class="stat-info">
					<div class="stat-value">
						{{ dashboardStats.totalDepartments }}
					</div>
					<div class="stat-label">Phòng ban</div>
				</div>
			</div>

			<div class="stat-card stat-card--purple" v-if="isAdmin">
				<div class="stat-icon-wrapper">
					<img :src="keyIcon" class="stat-icon" alt="users" />
				</div>
				<div class="stat-info">
					<div class="stat-value">
						{{ dashboardStats.totalAdminRole }}
					</div>
					<div class="stat-label">Tổng quản trị viên</div>
				</div>
			</div>
		</div>

		<!-- ===== Phân bổ nhân sự ===== -->
		<section class="dept-section">
			<h2 class="section-title">Nhân viên theo phòng ban</h2>
			<div class="dept-grid">
				<div
					v-for="d in departmentDistribution"
					:key="d.name"
					class="dept-card"
				>
					<div class="dept-header">
						<span class="dept-name">{{ d.name }}</span>
						<span class="dept-count">{{ d.count }} người</span>
					</div>
					<div class="progress-container">
						<div
							class="progress-bar"
							:style="{ width: d.percentage + '%' }"
						></div>
					</div>
				</div>
			</div>
		</section>
	</div>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useDashboardStore } from "@/stores/dashboard";

// --- Tài nguyên và Biểu tượng ---
import usersIcon from "@/assets/icons/users.svg";
import checkIcon from "@/assets/icons/check-circle.svg";
import buildingIcon from "@/assets/icons/building.svg";
import keyIcon from "@/assets/icons/key.svg";

// --- Store và Trạng thái ---
const authStore = useAuthStore();
const dashboardStore = useDashboardStore();

// --- Thuộc tính tính toán (Computed) ---

/**
 * Kiểm tra quyền quản trị của người dùng hiện tại
 */
const isAdmin = computed(() => authStore.user?.role === "admin");

/**
 * Trích xuất dữ liệu thống kê từ store với các giá trị mặc định
 */
const dashboardStats = computed(() => {
	const defaultStats = {
		totalEmployees: 0,
		activeEmployees: 0,
		inactiveEmployees: 0,
		totalDepartments: 0,
		totalUsers: 0,
		totalAdminRole: 0,
		employeesByDepartment: [],
	};

	return dashboardStore.dashboardData || defaultStats;
});

/**
 * Tính toán phân bổ nhân viên theo phòng ban và tỷ lệ phần trăm hiển thị
 */
const departmentDistribution = computed(() => {
	const departments = dashboardStats.value.employeesByDepartment || [];
	if (departments.length === 0) return [];

	const maxCount = Math.max(...departments.map(d => d.count), 1);

	return departments.map(d => ({
		name: d.name,
		count: d.count,
		percentage: Math.round((d.count / maxCount) * 100),
	}));
});

// --- Lifecycle Hooks ---

/**
 * Khởi tạo dữ liệu khi component được gắn vào DOM
 */
onMounted(async () => {
	try {
		await dashboardStore.fetchDashboardData();
	} catch (error) {
		console.error("[Dashboard] Fetch failed:", error);
	}
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Biến thiết kế ===== */
.dashboard-container {
	--bg-card: rgba(255, 255, 255, 0.98);
	--radius-lg: 16px;
	--radius-full: 9999px;
	--shadow-sm:
		0 4px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px -1px rgba(0, 0, 0, 0.01);
	--shadow-md:
		0 10px 25px -5px rgba(0, 0, 0, 0.03),
		0 8px 10px -6px rgba(0, 0, 0, 0.01);
	--shadow-hover:
		0 20px 25px -5px rgba(0, 0, 0, 0.05),
		0 10px 10px -5px rgba(0, 0, 0, 0.02);

	/* Colors for stat cards */
	--color-blue: #3b82f6;
	--color-green: #22c55e;
	--color-amber: #f59e0b;
	--color-purple: #a855f7;

	font-family: "Inter", sans-serif;
	color: var(--text-main);
	padding: 1rem 0;
}

/* ===== Thẻ thống kê (Stat Cards) ===== */
.stats-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
	gap: 1.5rem;
	margin-bottom: 3rem;
}

.stat-card {
	display: flex;
	align-items: center;
	gap: 1.25rem;
	padding: 1.5rem;
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	box-shadow: var(--shadow-md);
	transition:
		transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1),
		box-shadow 0.3s ease;
}

.stat-card:hover {
	transform: translateY(-4px);
	box-shadow: var(--shadow-hover);
}

.stat-icon-wrapper {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 56px;
	height: 56px;
	border-radius: 12px;
	flex-shrink: 0;
}

.stat-icon {
	width: 28px;
	height: 28px;
}

.stat-info {
	display: flex;
	flex-direction: column;
}

.stat-value {
	font-size: 1.875rem;
	font-weight: 700;
	line-height: 1.1;
	margin-bottom: 0.25rem;
}

.stat-label {
	font-size: 0.875rem;
	font-weight: 500;
	color: var(--text-muted);
}

/* Định nghĩa màu sắc cho từng loại thẻ */
.stat-card--blue .stat-icon {
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%);
}
.stat-card--green .stat-icon {
	filter: invert(61%) sepia(54%) saturate(518%) hue-rotate(97deg)
		brightness(92%) contrast(85%);
}
.stat-card--amber .stat-icon {
	filter: invert(68%) sepia(67%) saturate(3062%) hue-rotate(5deg)
		brightness(103%) contrast(93%);
}
.stat-card--purple .stat-icon {
	filter: invert(41%) sepia(87%) saturate(3268%) hue-rotate(256deg)
		brightness(101%) contrast(96%);
}

/* ===== Section: Nhân viên theo phòng ban ===== */
.section-title {
	font-size: 1.25rem;
	font-weight: 700;
	margin-bottom: 1.5rem;
	letter-spacing: -0.01em;
}

.dept-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
	gap: 1.25rem;
}

.dept-card {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: var(--radius-lg);
	padding: 1.5rem;
	box-shadow: var(--shadow-sm);
	transition: all 0.2s ease;
}

.dept-card:hover {
	transform: translateY(-2px);
	box-shadow: var(--shadow-md);
}

.dept-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1.25rem;
}

.dept-name {
	font-weight: 600;
	color: #334155;
}

.dept-count {
	font-size: 0.8125rem;
	font-weight: 600;
	color: var(--color-blue);
	background: #eff6ff;
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-full);
}

/* Progress bar styles */
.progress-container {
	height: 8px;
	background: #f1f5f9;
	border-radius: var(--radius-full);
	overflow: hidden;
}

.progress-bar {
	height: 100%;
	background: linear-gradient(90deg, var(--color-blue), #60a5fa);
	border-radius: var(--radius-full);
	transition: width 1.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
	.stats-grid,
	.dept-grid {
		grid-template-columns: 1fr;
	}

	.stat-card {
		padding: 1.25rem;
	}
}
</style>
