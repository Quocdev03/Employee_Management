<template>
	<div class="user-view">
		<!-- ===== Tiêu đề trang ===== -->
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý tài khoản</h1>
				<p class="page-subtitle">
					Kiểm soát quyền truy cập và tài khoản hệ thống (Chỉ dành cho Admin)
				</p>
			</div>
			<button class="btn btn--primary" @click="openAccountForm(null)">
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Tạo tài khoản
			</button>
		</header>

		<!-- ===== Nội dung chính (Bảng dữ liệu) ===== -->
		<main class="content-card">
			<!-- Thanh công cụ (Tìm kiếm & Chuyển tab) -->
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchKeyword"
						class="search-box__input form-control"
						:placeholder="
							viewMode === 'users'
								? 'Tìm kiếm email...'
								: 'Tìm kiếm tên nhân viên...'
						"
					/>
				</div>
				<div class="tab-switcher">
					<button
						class="tab-switcher__btn"
						:class="{ 'tab-switcher__btn--active': viewMode === 'users' }"
						@click="viewMode = 'users'"
					>
						Tài khoản ({{ userStore.users.length }})
					</button>
					<button
						class="tab-switcher__btn"
						:class="{ 'tab-switcher__btn--active': viewMode === 'missing' }"
						@click="viewMode = 'missing'"
					>
						Chưa có tài khoản ({{ employeesWithoutAccounts.length }})
					</button>
				</div>
			</div>

			<!-- Hiệu ứng tải dữ liệu (Skeleton) -->
			<div v-if="userStore.loading || employeeStore.loading" class="skeleton-list">
				<div v-for="i in 3" :key="i" class="skeleton-item"></div>
			</div>

			<!-- Bảng dữ liệu -->
			<div v-else class="table-wrapper">
				<table class="data-table">
					<thead>
						<tr v-if="viewMode === 'users'">
							<th>Email</th>
							<th>Quyền hạn</th>
							<th>Ngày tạo</th>
							<th class="text-right">Thao tác</th>
						</tr>
						<tr v-else>
							<th>Nhân viên</th>
							<th>Số điện thoại</th>
							<th>Phòng ban</th>
							<th class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<template v-if="viewMode === 'users'">
							<tr v-for="u in filteredUsers" :key="u.id">
								<td>
									<div class="user-info">
										<div class="user-info__avatar">
											{{ u.email ? u.email[0].toUpperCase() : "" }}
										</div>
										<div class="user-info__details">
											<div class="user-info__email">{{ u.email }}</div>
											<div class="user-info__name" v-if="u.employee">
												👤 {{ u.employee.name }}
											</div>
										</div>
									</div>
								</td>
								<td>
									<span :class="['badge', `badge--${u.role?.name || 'user'}`]">
										{{ u.role?.name === "admin" ? "Quản trị viên" : "Nhân viên" }}
									</span>
								</td>
								<td class="text-muted">{{ formatVietnameseDate(u.created_at) }}</td>
								<td class="text-right">
									<div class="action-group" v-if="canModify(u)">
										<button
											class="btn-icon btn-icon--edit"
											title="Sửa"
											@click="openAccountForm(u)"
										>
											<img :src="editIcon" alt="edit" />
										</button>
										<button
											class="btn-icon btn-icon--delete"
											title="Xoá"
											@click="prepareDelete(u)"
										>
											<img :src="deleteIcon" alt="delete" />
										</button>
									</div>
									<div v-else class="action-locked">
										<span title="Không có quyền chỉnh sửa">🔒</span>
									</div>
								</td>
							</tr>
						</template>
						<template v-else>
							<tr v-for="emp in filteredMissingAccounts" :key="emp.id">
								<td>
									<div class="user-info">
										<div class="user-info__avatar user-info__avatar--muted">
											{{ emp.name[0].toUpperCase() }}
										</div>
										<div class="user-info__details">
											<div class="user-info__name--primary">{{ emp.name }}</div>
										</div>
									</div>
								</td>
								<td>{{ emp.phone }}</td>
								<td>{{ emp.department?.name || "---" }}</td>
								<td class="text-right">
									<button
										class="btn btn--outline-primary btn--sm"
										@click="openFormForEmployee(emp)"
									>
										Cấp tài khoản
									</button>
								</td>
							</tr>
						</template>
					</tbody>
				</table>
			</div>
		</main>

		<!-- ===== Modal: Tạo/Sửa tài khoản ===== -->
		<div class="modal" v-if="isFormVisible" @click.self="isFormVisible = false">
			<div class="modal__dialog">
				<header class="modal__header">
					<div>
						<h2 class="modal__title">
							{{ selectedUser ? "Chỉnh sửa tài khoản" : "Tạo tài khoản mới" }}
						</h2>
						<p class="modal__subtitle">Thiết lập quyền truy cập cho nhân viên</p>
					</div>
					<button class="btn-close" @click="isFormVisible = false">
						<img :src="closeIcon" alt="close" />
					</button>
				</header>

				<form class="modal__body" @submit.prevent="handleFormSubmit">
					<div class="form-group">
						<label class="form-label">
							Email tài khoản <span class="required">*</span>
						</label>
						<div class="input-group">
							<img :src="mailIcon" class="input-group__icon" alt="mail" />
							<input
								v-model="accountForm.email"
								type="email"
								class="form-control input-group__input"
								placeholder="admin@company.dev"
							/>
						</div>
					</div>

					<div class="form-group">
						<label class="form-label">
							Mật khẩu <span class="required" v-if="!selectedUser">*</span>
						</label>
						<div class="input-group">
							<img :src="lockIcon" class="input-group__icon" alt="lock" />
							<input
								v-model="accountForm.password"
								type="password"
								class="form-control input-group__input"
								placeholder="••••••••"
							/>
						</div>
					</div>

					<div class="form-group">
						<label class="form-label">
							Nhập lại mật khẩu <span class="required" v-if="!selectedUser">*</span>
						</label>
						<div class="input-group">
							<img :src="lockIcon" class="input-group__icon" alt="lock" />
							<input
								v-model="accountForm.confirmPassword"
								type="password"
								class="form-control input-group__input"
								placeholder="••••••••"
							/>
						</div>
					</div>

					<div class="form-group">
						<label class="form-label">Phân quyền</label>
						<div class="input-group">
							<img :src="userIcon" class="input-group__icon" alt="role" />
							<select
								v-model="accountForm.role_id"
								class="form-control form-select input-group__input"
								:disabled="selectedUser && selectedUser.id === authStore.user?.id"
							>
								<option :value="2">Nhân viên (User)</option>
								<option :value="1">Quản trị viên (Admin)</option>
							</select>
						</div>
					</div>

					<div class="form-group">
						<label class="form-label">Nhân viên liên kết</label>
						<div class="input-group">
							<img :src="deptIcon" class="input-group__icon" alt="employee" />
							<select
								v-model="accountForm.employee_id"
								class="form-control form-select input-group__input"
							>
								<option :value="null" disabled>Chọn nhân viên</option>
								<option
									v-for="emp in availableEmployees"
									:key="emp.id"
									:value="emp.id"
								>
									{{ emp.id }} - {{ emp.name }} - {{ emp.email || emp.phone }}
								</option>
							</select>
						</div>
					</div>

					<footer class="modal__footer">
						<button type="button" class="btn btn--secondary" @click="isFormVisible = false">
							Huỷ bỏ
						</button>
						<button
							type="submit"
							class="btn btn--primary"
							:disabled="isSubmitting"
						>
							{{ isSubmitting ? "Đang lưu..." : (selectedUser ? "Cập nhật" : "Tạo tài khoản") }}
						</button>
					</footer>
				</form>
			</div>
		</div>

		<!-- ===== Modal: Xác nhận xoá ===== -->
		<div class="modal" v-if="isDeleteModalVisible" @click.self="isDeleteModalVisible = false">
			<div class="modal__dialog modal__dialog--confirm">
				<div class="confirm-icon">⚠️</div>
				<h3 class="modal__title text-center">Xác nhận xoá</h3>
				<p class="modal__subtitle text-center">
					Bạn có chắc chắn muốn xoá tài khoản <strong>{{ userToDelete?.email }}</strong>? Hành động này không thể hoàn tác.
				</p>
				<div class="confirm-actions">
					<button class="btn btn--secondary flex-1" @click="isDeleteModalVisible = false">
						Hủy bỏ
					</button>
					<button class="btn btn--danger flex-1" @click="executeDelete">
						Đồng ý xoá
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useToast } from "vue-toastification";
import { useUserStore } from "@/stores/user";
import { useEmployeeStore } from "@/stores/employee";
import { useAuthStore } from "@/stores/auth";

// --- Tài nguyên và Biểu tượng ---
import plusIcon from "@/assets/icons/plus.svg";
import searchIcon from "@/assets/icons/search.svg";
import editIcon from "@/assets/icons/edit.svg";
import deleteIcon from "@/assets/icons/delete.svg";
import mailIcon from "@/assets/icons/mail.svg";
import lockIcon from "@/assets/icons/lock.svg";
import closeIcon from "@/assets/icons/close.svg";
import userIcon from "@/assets/icons/user.svg";
import deptIcon from "@/assets/icons/department.svg";

// --- Cấu hình và Hằng số ---
const ROLE_ADMIN_ID = 1;
const ROLE_USER_ID = 2;

// --- Store và Trạng thái ---
const userStore = useUserStore();
const employeeStore = useEmployeeStore();
const authStore = useAuthStore();
const toast = useToast();

const viewMode = ref("users"); // 'users' | 'missing'
const searchKeyword = ref("");
const isFormVisible = ref(false);
const isDeleteModalVisible = ref(false);
const isSubmitting = ref(false);

const selectedUser = ref(null);
const userToDelete = ref(null);

const accountForm = reactive({
	email: "",
	password: "",
	confirmPassword: "",
	role_id: ROLE_USER_ID,
	employee_id: null,
});

// --- Thuộc tính tính toán (Computed) ---

/**
 * Danh sách tài khoản người dùng được lọc theo từ khóa tìm kiếm
 */
const filteredUsers = computed(() => {
	const query = searchKeyword.value.toLowerCase();
	return userStore.users.filter((u) => u.email?.toLowerCase().includes(query));
});

/**
 * Danh sách nhân viên chưa được cấp tài khoản hệ thống
 */
const employeesWithoutAccounts = computed(() => {
	return employeeStore.employees.filter((emp) => !emp.user?.id);
});

/**
 * Danh sách nhân viên chưa có tài khoản, được lọc theo tìm kiếm
 */
const filteredMissingAccounts = computed(() => {
	const query = searchKeyword.value.toLowerCase();
	return employeesWithoutAccounts.value.filter((emp) => 
		emp.name?.toLowerCase().includes(query)
	);
});

/**
 * Danh sách nhân viên có thể gán cho tài khoản (chưa có tk hoặc đang được gán cho chính user đang sửa)
 */
const availableEmployees = computed(() => {
	return employeeStore.employees.filter((emp) => {
		const isNotLinked = !emp.user?.id;
		const isLinkedToCurrent = selectedUser.value && emp.id === selectedUser.value.employee_id;
		return isNotLinked || isLinkedToCurrent;
	});
});

/**
 * Kiểm tra xem người dùng hiện tại có quyền chỉnh sửa tài khoản mục tiêu không
 * - Không thể sửa chính mình (phải qua Profile)
 * - Admin không thể sửa/xoá Admin khác
 */
const canModify = (targetUser) => {
	if (!authStore.user) return false;
	
	// Cho phép tự sửa chính mình
	if (targetUser.id === authStore.user.id) return true;
	
	// Quản trị viên không thể sửa Quản trị viên khác
	if (targetUser.role_id === ROLE_ADMIN_ID) return false;
	
	return true;
};

// --- Xử lý logic nghiệp vụ ---

/**
 * Tải dữ liệu từ máy chủ
 */
async function fetchAllData() {
	try {
		await Promise.all([
			userStore.fetchUsers(),
			employeeStore.fetchEmployees({ limit: 1000 }),
		]);
	} catch (error) {
		toast.error("Không thể tải dữ liệu hệ thống");
		console.error("[UserManagement] Fetch failed:", error);
	}
}

/**
 * Làm mới dữ liệu form về trạng thái ban đầu
 */
function resetAccountForm() {
	accountForm.email = "";
	accountForm.password = "";
	accountForm.confirmPassword = "";
	accountForm.role_id = ROLE_USER_ID;
	accountForm.employee_id = null;
}

/**
 * Mở modal form để tạo mới hoặc cập nhật tài khoản
 */
function openAccountForm(user = null) {
	selectedUser.value = user;
	resetAccountForm();

	if (user) {
		accountForm.email = user.email || "";
		accountForm.role_id = user.role_id;
		accountForm.employee_id = user.employee_id || null;
	}
	
	isFormVisible.value = true;
}

/**
 * Mở modal tạo tài khoản nhanh từ thông tin một nhân viên
 */
function openFormForEmployee(employee) {
	openAccountForm(null);
	accountForm.employee_id = employee.id;
}

/**
 * Kiểm tra tính hợp lệ của dữ liệu form trước khi gửi
 */
function validateAccountData() {
	const { email, password, confirmPassword, role_id, employee_id } = accountForm;
	
	// Kiểm tra các trường bắt buộc khi tạo mới
	if (!selectedUser.value) {
		if (!email) return "Vui lòng nhập email tài khoản";
		if (!password) return "Vui lòng nhập mật khẩu";
	}
	
	// Kiểm tra mật khẩu khớp nhau (nếu có nhập mật khẩu)
	if (password && password !== confirmPassword) {
		return "Mật khẩu xác nhận không khớp";
	}
	
	// Kiểm tra liên kết nhân viên cho quyền User
	if (role_id === ROLE_USER_ID && !employee_id) {
		return "Vui lòng chọn nhân viên liên kết";
	}
	
	return null;
}

/**
 * Xây dựng payload cập nhật chỉ chứa các trường thay đổi
 */
function buildUpdatePayload() {
	const payload = {};
	const current = selectedUser.value;

	if (accountForm.email !== current.email) payload.email = accountForm.email;
	if (accountForm.password) payload.password = accountForm.password;
	if (accountForm.role_id !== current.role_id) payload.role_id = accountForm.role_id;

	const isEmployeeChanged = accountForm.employee_id !== current.employee_id;
	if (isEmployeeChanged) {
		payload.employee_id = accountForm.employee_id;
	}

	return payload;
}

/**
 * Xử lý tạo mới tài khoản
 */
async function createNewAccount() {
	const payload = {
		email: accountForm.email,
		password: accountForm.password,
		role_id: accountForm.role_id,
		employee_id: accountForm.employee_id,
	};
	
	await userStore.createUser(payload);
	toast.success("Tạo tài khoản thành công");
}

/**
 * Xử lý cập nhật tài khoản hiện có
 */
async function updateExistingAccount() {
	const payload = buildUpdatePayload();
	const hasChanges = Object.keys(payload).length > 0;

	if (!hasChanges) {
		toast.info("Không có thay đổi nào để cập nhật");
		return false;
	}

	await userStore.updateUser(selectedUser.value.id, payload);
	toast.success("Cập nhật thông tin thành công");
	return true;
}

/**
 * Xử lý sự kiện Submit của form
 */
async function handleFormSubmit() {
	const error = validateAccountData();
	if (error) {
		toast.warning(error);
		return;
	}

	isSubmitting.value = true;
	try {
		let success = false;
		if (selectedUser.value) {
			success = await updateExistingAccount();
		} else {
			await createNewAccount();
			success = true;
		}

		if (success) {
			await fetchAllData();
			isFormVisible.value = false;
		}
	} catch (err) {
		toast.error(err.message || "Thao tác thất bại");
	} finally {
		isSubmitting.value = false;
	}
}

/**
 * Chuẩn bị thông tin cho việc xoá tài khoản
 */
function prepareDelete(user) {
	userToDelete.value = user;
	isDeleteModalVisible.value = true;
}

/**
 * Thực hiện xoá tài khoản sau khi xác nhận
 */
async function executeDelete() {
	const target = userToDelete.value;
	if (!target) return;
	
	try {
		await userStore.deleteUser(target.id);
		await userStore.fetchUsers();
		toast.info("Đã xoá tài khoản thành công");
		isDeleteModalVisible.value = false;
	} catch (err) {
		toast.error("Lỗi khi thực hiện xoá tài khoản");
		console.error("[UserManagement] Delete failed:", err);
	} finally {
		userToDelete.value = null;
	}
}

// --- Xử lý Format dữ liệu ---

/**
 * Định dạng ngày tháng sang dd/mm/yyyy
 */
function formatVietnameseDate(dateStr) {
	if (!dateStr) return "—";
	return new Date(dateStr).toLocaleDateString("vi-VN");
}

// --- Lifecycle Hooks ---
onMounted(fetchAllData);
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Biến thiết kế (Design Tokens) ===== */
.user-view {
	--primary-color: #3b82f6;
	--primary-hover: #2563eb;
	--danger-color: #ef4444;
	--danger-hover: #dc2626;
	--text-main: #1e293b;
	--text-muted: #64748b;
	--bg-main: #f8fafc;
	--bg-card: #ffffff;
	--border-color: #e2e8f0;
	--border-hover: #cbd5e1;
	
	--radius-sm: 8px;
	--radius-md: 12px;
	--radius-lg: 20px;
	--radius-xl: 24px;
	
	--shadow-sm: 0 2px 4px rgba(0, 0, 0, 0.05);
	--shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
	--shadow-lg: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
	--shadow-modal: 0 25px 50px -12px rgba(0, 0, 0, 0.25);

	font-family: "Inter", sans-serif;
	color: var(--text-main);
}

/* ===== Tiện ích (Utilities) ===== */
.text-center { text-align: center; }
.text-right { text-align: right; }
.text-muted { color: var(--text-muted); font-size: 0.875rem; }
.flex-1 { flex: 1; }
.required { color: var(--danger-color); }

/* ===== Nút bấm (Buttons) ===== */
.btn {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 8px;
	padding: 0.75rem 1.5rem;
	border: none;
	border-radius: var(--radius-md);
	font-weight: 600;
	cursor: pointer;
	transition: all 0.2s ease;
}

.btn:disabled {
	opacity: 0.7;
	cursor: not-allowed;
}

.btn--sm {
	padding: 0.5rem 1rem;
	font-size: 0.8rem;
}

.btn--primary {
	background: var(--primary-color);
	color: white;
	box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
}

.btn--primary:hover:not(:disabled) {
	background: var(--primary-hover);
	transform: translateY(-1px);
	box-shadow: 0 6px 16px rgba(59, 130, 246, 0.3);
}

.btn--secondary {
	background: var(--bg-main);
	color: var(--text-muted);
	border: 1px solid var(--border-color);
}

.btn--secondary:hover:not(:disabled) {
	background: var(--border-color);
	color: var(--text-main);
}

.btn--danger {
	background: var(--danger-color);
	color: white;
	box-shadow: 0 4px 12px rgba(239, 68, 68, 0.25);
}

.btn--danger:hover:not(:disabled) {
	background: var(--danger-hover);
	transform: translateY(-1px);
	box-shadow: 0 6px 16px rgba(239, 68, 68, 0.35);
}

.btn--outline-primary {
	background: #eff6ff;
	color: var(--primary-color);
	border: 1px solid #dbeafe;
}

.btn--outline-primary:hover:not(:disabled) {
	background: var(--primary-color);
	color: white;
	transform: translateY(-1px);
	box-shadow: var(--shadow-sm);
}

.btn__icon {
	width: 20px;
	height: 20px;
	filter: brightness(0) invert(1);
}

.btn-icon {
	width: 36px;
	height: 36px;
	border-radius: 10px;
	border: 1px solid var(--border-color);
	background: white;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s;
}

.btn-icon img {
	width: 18px;
	height: 18px;
	filter: invert(47%) sepia(13%) saturate(545%) hue-rotate(183deg) brightness(95%) contrast(87%);
}

.btn-icon:hover {
	background: var(--bg-main);
	color: var(--text-main);
	border-color: var(--border-hover);
}

.btn-icon--edit:hover {
	border-color: var(--primary-color);
	background: #eff6ff;
}
.btn-icon--edit:hover img {
	filter: invert(44%) sepia(91%) saturate(1185%) hue-rotate(200deg) brightness(101%) contrast(92%);
}

.btn-icon--delete:hover {
	border-color: var(--danger-color);
	background: #fef2f2;
}
.btn-icon--delete:hover img {
	filter: invert(41%) sepia(82%) saturate(4529%) hue-rotate(341deg) brightness(98%) contrast(93%);
}

.btn-close {
	background: var(--bg-main);
	border: none;
	width: 36px;
	height: 36px;
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s;
}

.btn-close:hover {
	background: #fee2e2;
}

.btn-close:hover img {
	filter: invert(34%) sepia(85%) saturate(3485%) hue-rotate(342deg) brightness(98%) contrast(96%);
}

.btn-close img {
	width: 18px;
	height: 18px;
}

/* ===== Form Elements ===== */
.form-group {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.form-label {
	font-size: 0.875rem;
	font-weight: 600;
	color: #475569;
}

.form-control {
	width: 100%;
	padding: 0.75rem 1rem;
	border: 1px solid var(--border-color);
	border-radius: var(--radius-md);
	font-size: 0.95rem;
	background: var(--bg-main);
	transition: all 0.2s;
	outline: none;
}

.form-control:focus {
	background: white;
	border-color: var(--primary-color);
	box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
}

.form-select {
	appearance: none;
	background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%2394a3b8'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E");
	background-repeat: no-repeat;
	background-position: right 1rem center;
	background-size: 1rem;
}

.input-group {
	position: relative;
	display: flex;
	align-items: center;
}

.input-group__icon {
	position: absolute;
	left: 1rem;
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%) contrast(89%);
	z-index: 1;
}

.input-group__input {
	padding-left: 2.75rem;
}

/* ===== Layout: Trang chính ===== */
.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 2rem;
	flex-wrap: wrap;
	gap: 1rem;
}

.page-title {
	font-size: 1.75rem;
	font-weight: 700;
	letter-spacing: -0.02em;
	margin: 0 0 4px 0;
}

.page-subtitle {
	color: var(--text-muted);
	font-size: 0.95rem;
	margin: 0;
}

.content-card {
	background: var(--bg-card);
	border-radius: var(--radius-lg);
	border: 1px solid var(--border-color);
	box-shadow: var(--shadow-sm);
	overflow: hidden;
}

.toolbar {
	padding: 1.5rem;
	border-bottom: 1px solid var(--bg-main);
	display: flex;
	justify-content: space-between;
	align-items: center;
	flex-wrap: wrap;
	gap: 1rem;
}

.search-box {
	position: relative;
	width: 100%;
	max-width: 400px;
}

.search-box__icon {
	position: absolute;
	left: 1rem;
	top: 50%;
	transform: translateY(-50%);
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%) contrast(89%);
}

.search-box__input {
	padding-left: 2.75rem;
}

.tab-switcher {
	display: flex;
	background: var(--bg-main);
	padding: 4px;
	border-radius: var(--radius-md);
	gap: 4px;
}

.tab-switcher__btn {
	padding: 0.5rem 1rem;
	border: none;
	background: transparent;
	border-radius: var(--radius-sm);
	font-size: 0.875rem;
	font-weight: 600;
	color: var(--text-muted);
	cursor: pointer;
	transition: all 0.2s;
}

.tab-switcher__btn--active {
	background: white;
	color: var(--text-main);
	box-shadow: var(--shadow-sm);
}

/* ===== Layout: Bảng dữ liệu ===== */
.table-wrapper {
	overflow-x: auto;
}

.data-table {
	width: 100%;
	border-collapse: collapse;
}

.data-table th {
	padding: 1rem 1.5rem;
	background: var(--bg-main);
	text-align: left;
	font-size: 0.75rem;
	font-weight: 700;
	color: var(--text-muted);
	text-transform: uppercase;
	letter-spacing: 0.05em;
}

.data-table td {
	padding: 1rem 1.5rem;
	border-bottom: 1px solid var(--bg-main);
}

.user-info {
	display: flex;
	align-items: center;
	gap: 12px;
}

.user-info__avatar {
	width: 32px;
	height: 32px;
	background: #eff6ff;
	color: var(--primary-color);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 700;
	font-size: 0.8rem;
}

.user-info__avatar--muted {
	background: var(--bg-main);
	color: var(--text-muted);
}

.user-info__details {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.user-info__email {
	font-weight: 500;
	color: #334155;
}

.user-info__name {
	font-size: 0.8rem;
	color: var(--text-muted);
}

.user-info__name--primary {
	font-weight: 500;
	color: #334155;
}

/* Badges */
.badge {
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-sm);
	font-size: 0.75rem;
	font-weight: 700;
	display: inline-block;
}

.badge--admin {
	background: #ede9fe;
	color: #5b21b6;
}

.badge--user {
	background: #e0f2fe;
	color: #0369a1;
}

/* Actions */
.action-group {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
}

.action-locked {
	color: #94a3b8;
	font-size: 1.1rem;
	padding-right: 0.5rem;
}

/* ===== Layout: Modals ===== */
.modal {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.6);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: 1.5rem;
}

.modal__dialog {
	background: var(--bg-card);
	border-radius: var(--radius-xl);
	width: 100%;
	max-width: 500px;
	max-height: 90vh;
	box-shadow: var(--shadow-modal);
	display: flex;
	flex-direction: column;
	overflow: hidden;
	animation: slideUp 0.3s ease-out;
}

.modal__dialog--confirm {
	max-width: 400px;
	padding: 2rem;
}

.modal__header {
	padding: 1.5rem;
	border-bottom: 1px solid var(--bg-main);
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.modal__title {
	font-size: 1.25rem;
	font-weight: 700;
	margin: 0 0 0.25rem 0;
}

.modal__subtitle {
	font-size: 0.85rem;
	color: var(--text-muted);
	margin: 0;
}

.modal__body {
	padding: 1.5rem;
	display: flex;
	flex-direction: column;
	gap: 1.25rem;
	overflow-y: auto;
}

.modal__footer {
	padding: 1.25rem 1.5rem;
	border-top: 1px solid var(--bg-main);
	display: flex;
	justify-content: flex-end;
	gap: 12px;
	background: var(--bg-card);
}

/* Confirm modal specific */
.confirm-icon {
	font-size: 3.5rem;
	margin-bottom: 1rem;
	text-align: center;
}

.confirm-actions {
	display: flex;
	gap: 12px;
	margin-top: 1.5rem;
}

/* ===== Animations ===== */
@keyframes slideUp {
	from { transform: translateY(20px); opacity: 0; }
	to { transform: translateY(0); opacity: 1; }
}

@keyframes pulse {
	0%, 100% { opacity: 1; }
	50% { opacity: 0.5; }
}

/* Skeleton Loaders */
.skeleton-list {
	padding: 1.5rem;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.skeleton-item {
	height: 64px;
	background: var(--bg-main);
	border-radius: var(--radius-md);
	animation: pulse 1.5s infinite;
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
	.page-header {
		flex-direction: column;
		align-items: flex-start;
	}
	
	.toolbar {
		flex-direction: column;
	}
	
	.search-box {
		max-width: 100%;
	}
	
	.tab-switcher {
		width: 100%;
		justify-content: space-between;
	}
	
	.tab-switcher__btn {
		flex: 1;
		text-align: center;
	}
}
</style>
