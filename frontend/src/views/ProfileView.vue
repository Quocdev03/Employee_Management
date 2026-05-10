<template>
	<div class="profile-view">
		<!-- ===== Tiêu đề trang ===== -->
		<div class="page-header">
			<div class="header-content">
				<h1>Hồ sơ cá nhân</h1>
				<p class="subtitle">
					Quản lý và xem thông tin chi tiết của bạn
				</p>
			</div>
		</div>

		<!-- ===== Bố cục hồ sơ ===== -->
		<div class="profile-grid">
			<!-- ===== Cột trái: Thông tin tổng quan ===== -->
			<div class="profile-card sidebar-card">
				<div class="avatar-section">
					<div class="avatar-container">
						<img
							:src="userProfile.avatar_url"
							class="profile-avatar"
						/>
					</div>
					<h2 class="user-name">{{ userProfile.name }}</h2>
					<p class="user-role-tag">
						{{ getRoleLabel(userProfile.role) }}
					</p>
					<div class="status-indicator">
						<span
							:class="[
								'status-dot',
								userProfile.isActive ? 'active' : 'inactive',
							]"
						></span>
						<span>{{
							userProfile.isActive
								? "Đang làm việc"
								: "Đã nghỉ việc"
						}}</span>
					</div>
				</div>

				<div class="quick-stats">
					<div class="stat-item">
						<span class="stat-label">Ngày tham gia:</span>
						<span class="stat-value">{{
							formatVietnameseDate(userProfile.hire_date)
						}}</span>
					</div>
					<div class="stat-item">
						<span class="stat-label">Phòng ban:</span>
						<span class="stat-value">{{
							userProfile.department
						}}</span>
					</div>
				</div>
			</div>

			<!-- ===== Cột phải: Thông tin chi tiết ===== -->
			<div class="profile-main">
				<!-- ===== Block: Thông tin cá nhân ===== -->
				<div class="profile-card main-card">
					<div class="card-header">
						<img :src="userBlueIcon" alt="user" class="card-icon" />
						<h3>Thông tin cá nhân</h3>
					</div>
					<div class="info-grid">
						<div class="info-group">
							<label>Họ và tên</label>
							<div class="value">{{ userProfile.name }}</div>
						</div>
						<div class="info-group">
							<label>Giới tính</label>
							<div class="value">
								{{ getGenderLabel(userProfile.gender) }}
							</div>
						</div>
						<div class="info-group">
							<label>Ngày sinh</label>
							<div class="value">
								{{
									formatVietnameseDate(
										userProfile.date_of_birth,
									)
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Số điện thoại</label>
							<div class="value">{{ userProfile.phone }}</div>
						</div>
						<div class="info-group">
							<label>Email cá nhân</label>
							<div class="value">{{ userProfile.email }}</div>
						</div>
					</div>
				</div>

				<!-- ===== Block: Thông tin công việc ===== -->
				<div class="profile-card main-card mt-6">
					<div class="card-header">
						<img
							:src="briefcaseGreenIcon"
							alt="work"
							class="card-icon"
						/>
						<h3>Thông tin công việc</h3>
					</div>
					<div class="info-grid">
						<div class="info-group">
							<label>Chức vụ</label>
							<div class="value">{{ userProfile.position }}</div>
						</div>
						<div class="info-group">
							<label>Phòng ban</label>
							<div class="value">
								{{ userProfile.department }}
							</div>
						</div>
						<div class="info-group">
							<label>Ngày vào làm</label>
							<div class="value">
								{{
									formatVietnameseDate(userProfile.hire_date)
								}}
							</div>
						</div>
						<div class="info-group">
							<label>Mức lương hiện tại</label>
							<div class="value salary">
								{{ formatCurrencyVND(userProfile.salary) }}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";

// --- Tài nguyên và Biểu tượng ---
import userBlueIcon from "@/assets/icons/user-blue.svg";
import briefcaseGreenIcon from "@/assets/icons/briefcase-green.svg";

// --- Store và Dữ liệu nguồn ---
const authStore = useAuthStore();

// --- Thuộc tính tính toán (Computed) ---

/**
 * Tổng hợp và chuẩn hóa dữ liệu hồ sơ người dùng
 */
const userProfile = computed(() => {
	const rawUser = authStore.user;
	if (!rawUser) return {};

	const employee = rawUser.employee || {};

	return {
		email: rawUser.email,
		name: rawUser.name,
		gender: rawUser.gender,
		role: rawUser.role,
		isActive: employee.status === 1,
		isAdmin: rawUser.role === "admin",
		...employee,
		department: employee.department?.name || "N/A",
		position: employee.position?.name || "N/A",
	};
});

// --- Xử lý Format dữ liệu ---

/**
 * Định dạng ngày tháng sang kiểu Việt Nam (dd MMMM, yyyy)
 */
function formatVietnameseDate(dateInput) {
	if (!dateInput) return "Chưa cập nhật";

	const config = {
		year: "numeric",
		month: "long",
		day: "numeric",
	};

	return new Date(dateInput).toLocaleDateString("vi-VN", config);
}

/**
 * Định dạng số tiền sang chuẩn tiền tệ Việt Nam (VND)
 */
function formatCurrencyVND(amount) {
	if (!amount) return "—";

	return new Intl.NumberFormat("vi-VN", {
		style: "currency",
		currency: "VND",
	}).format(amount);
}

/**
 * Chuyển đổi mã giới tính sang nhãn tiếng Việt
 */
function getGenderLabel(genderCode) {
	const genderLabels = {
		male: "Nam",
		female: "Nữ",
		other: "Khác",
	};
	return genderLabels[genderCode] || "Chưa rõ";
}

/**
 * Chuyển đổi mã vai trò sang nhãn hiển thị tiếng Việt
 */
function getRoleLabel(roleCode) {
	return roleCode === "admin" ? "Quản trị viên" : "Nhân viên";
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Khai báo biến CSS ===== */
.profile-view {
	font-family: "Inter", sans-serif;
	width: 100%;
	margin: 0 auto;
}

/* ===== Tiêu đề trang ===== */
.page-header {
	margin-bottom: 2rem;
}

.page-header h1 {
	font-size: 1.75rem;
	font-weight: 700;
	color: var(--text-main);
	letter-spacing: -0.02em;
}

.subtitle {
	color: var(--text-muted);
	font-size: 0.95rem;
	margin-top: 4px;
}

/* ===== Bố cục chính ===== */
.profile-grid {
	display: grid;
	grid-template-columns: 340px 1fr;
	gap: 1.5rem;
}

/* ===== Card dùng chung ===== */
.profile-card {
	background: var(--bg-card);
	border-radius: var(--border-radius-card);
	border: 1px solid var(--border-color);
	box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.05);
	overflow: hidden;
}

/* ===== Cột trái: Sidebar ===== */
.sidebar-card {
	padding: 2.5rem 1.5rem;
	display: flex;
	flex-direction: column;
	align-items: center;
	height: fit-content;
}

.avatar-section {
	text-align: center;
	margin-bottom: 2rem;
}

.avatar-container {
	width: 140px;
	height: 140px;
	margin: 1.5rem auto;
}

.profile-avatar {
	width: 100%;
	height: 100%;
	border-radius: 40px;
	object-fit: cover;
	box-shadow: 0 10px 25px -5px rgba(59, 130, 246, 0.2);
}

.user-name {
	font-size: 1.5rem;
	font-weight: 700;
	color: var(--text-main);
	margin-bottom: 0.6rem;
}

.user-role-tag {
	display: inline-block;
	padding: 0.35rem 1rem;
	background: var(--bg-light);
	color: #475569;
	border-radius: 999px;
	font-size: 0.8125rem;
	font-weight: 600;
	margin-bottom: 0.75rem;
}

.status-indicator {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	font-size: 0.875rem;
	color: var(--text-muted);
}

.status-dot {
	width: 8px;
	height: 8px;
	border-radius: 50%;
}
.status-dot.active {
	background: var(--success-color);
	box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.1);
}
.status-dot.inactive {
	background: var(--text-muted);
}

.quick-stats {
	width: 100%;
	border-top: 1px solid var(--bg-light);
	padding-top: 1.5rem;
	margin-top: 0.5rem;
}

.stat-item {
	display: flex;
	justify-content: space-between;
	padding: 0.75rem 0;
}

.stat-label {
	font-size: 0.875rem;
	color: var(--text-muted);
}

.stat-value {
	font-size: 0.875rem;
	font-weight: 600;
	color: var(--text-main);
}

/* ===== Cột phải: Thông tin chi tiết ===== */
.main-card {
	padding: 1.75rem 2rem;
}

.card-header {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 1.5rem;
}

.card-header h3 {
	font-size: 1.125rem;
	font-weight: 700;
	color: var(--text-main);
}

.card-icon {
	width: 20px;
	height: 20px;
}

.info-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
	gap: 1.5rem;
}

.info-group label {
	display: block;
	font-size: 0.75rem;
	font-weight: 700;
	color: var(--text-light);
	text-transform: uppercase;
	letter-spacing: 0.05em;
	margin-bottom: 0.5rem;
}

.info-group .value {
	font-size: 1rem;
	font-weight: 500;
	color: #334155;
	padding: 0.75rem 1rem;
	background: var(--bg-lighter);
	border-radius: 12px;
	border: 1px solid var(--bg-light);
}

.info-group .value.salary {
	color: #059669;
	font-weight: 700;
}

/* ===== Utilities ===== */
.mt-6 {
	margin-top: 1.5rem;
}

/* ===== Responsive ===== */
@media (max-width: 1023px) {
	.profile-grid {
		grid-template-columns: 1fr;
	}

	.stat-item {
		display: flex;
		justify-content: flex-start;
		padding: 0.75rem 0;
		gap: 1rem;
	}
}

@media (max-width: 640px) {
	.avatar-section {
		flex-direction: column;
		text-align: center;
		gap: 1rem;
	}
	.info-grid {
		grid-template-columns: 1fr;
	}
	.quick-stats {
		grid-template-columns: 1fr;
		gap: 0;
	}
}
</style>
