<template>
	<div class="employee-view">
		<!-- ===== Tiêu đề trang ===== -->
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Danh sách nhân viên</h1>
				<p class="page-subtitle">
					Hệ thống có tổng cộng {{ employeeStore.total }} nhân viên
				</p>
			</div>
			<button
				v-if="isAdmin"
				class="btn btn--primary"
				@click="openEmployeeForm(null)"
			>
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm nhân viên
			</button>
		</header>

		<!-- ===== Nội dung chính ===== -->
		<main class="content-card">
			<!-- Thanh công cụ (Tìm kiếm & Bộ lọc) -->
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchKeyword"
						class="form-control search-box__input"
						placeholder="Tìm tên hoặc số điện thoại..."
						@input="handleSearchInput"
					/>
				</div>

				<div class="filter-group">
					<select
						v-model="selectedDeptId"
						@change="refreshToFirstPage"
						class="form-control form-select"
					>
						<option value="">Tất cả phòng ban</option>
						<option
							v-for="d in deptStore.departments"
							:key="d.id"
							:value="d.id"
						>
							{{ d.name }}
						</option>
					</select>
				</div>
			</div>

			<!-- Hiệu ứng tải dữ liệu (Skeleton) -->
			<div v-if="employeeStore.loading" class="skeleton-list">
				<div v-for="i in 5" :key="i" class="skeleton-item"></div>
			</div>

			<!-- Trạng thái trống -->
			<div v-else-if="isEmployeeListEmpty" class="empty-state">
				<div class="empty-state__icon">👥</div>
				<p class="empty-state__text">Không tìm thấy nhân viên nào phù hợp.</p>
			</div>

			<!-- Bảng dữ liệu -->
			<div v-else class="table-wrapper">
				<table class="data-table">
					<thead>
						<tr>
							<th>Nhân viên</th>
							<th>Ngày sinh</th>
							<th>Giới tính</th>
							<th>Liên hệ</th>
							<th>Phòng ban / Chức vụ</th>
							<th v-if="isAdmin">Lương</th>
							<th>Trạng thái</th>
							<th v-if="isAdmin" class="text-right">Thao tác</th>
						</tr>
					</thead>
					<tbody>
						<tr
							v-for="emp in employeeStore.employees"
							:key="emp.id"
						>
							<td>
								<div class="user-info">
									<div class="user-info__avatar-wrap">
										<img
											v-if="emp.avatar_url"
											:src="emp.avatar_url"
											class="user-info__avatar-img"
											alt="avatar"
										/>
										<div
											v-else
											:class="['user-info__initials', `user-info__initials--${emp.gender}`]"
										>
											{{ getInitials(emp.name) }}
										</div>
									</div>
									<div class="user-info__details">
										<div class="user-info__name">{{ emp.name }}</div>
										<div class="user-info__email">
											{{ emp.user?.email || "Chưa cấp tài khoản" }}
										</div>
									</div>
								</div>
							</td>
							<td class="text-muted">{{ formatVietnameseDate(emp.date_of_birth) }}</td>
							<td>
								<span :class="['tag', `tag--${emp.gender}`]">
									{{ getGenderLabel(emp.gender) }}
								</span>
							</td>
							<td>
								<div class="text-main fw-500">{{ emp.phone || "—" }}</div>
							</td>
							<td>
								<div class="job-info">
									<div class="job-info__dept">{{ emp.department?.name || "Vãng lai" }}</div>
									<div class="job-info__pos">{{ emp.position || "Nhân viên" }}</div>
								</div>
							</td>
							<td v-if="isAdmin">
								<span class="salary-amount">{{ formatCurrency(emp.salary) }}</span>
							</td>
							<td>
								<span :class="['status-badge', emp.status == 1 ? 'status-badge--active' : 'status-badge--inactive']">
									{{ emp.status == 1 ? "Đang làm" : "Đã nghỉ" }}
								</span>
							</td>
							<td v-if="isAdmin" class="text-right">
								<div class="action-group">
									<button
										class="btn-icon btn-icon--edit"
										title="Chỉnh sửa"
										@click="openEmployeeForm(emp)"
									>
										<img :src="editIcon" alt="edit" />
									</button>
									<button
										class="btn-icon btn-icon--delete"
										title="Xoá"
										@click="prepareDelete(emp.id, emp.name)"
									>
										<img :src="deleteIcon" alt="delete" />
									</button>
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>

			<!-- Phân trang -->
			<div class="pagination" v-if="shouldShowPagination">
				<button
					class="pagination__btn"
					:disabled="currentPage <= 1"
					@click="goToPrevPage"
				>
					<img :src="prevIcon" alt="prev" />
				</button>
				<div class="pagination__info">
					Trang <span>{{ currentPage }}</span> / {{ totalPages }}
				</div>
				<button
					class="pagination__btn"
					:disabled="currentPage >= totalPages"
					@click="goToNextPage"
				>
					<img :src="nextIcon" alt="next" />
				</button>
			</div>
		</main>

		<!-- ===== Form nhập liệu ===== -->
		<EmployeeForm
			v-if="isFormVisible"
			:employee="editingEmployee"
			@close="isFormVisible = false"
			@saved="handleSaveSuccess"
		/>

		<!-- ===== Modal Xác nhận xoá ===== -->
		<div
			v-if="isDeleteModalVisible"
			class="modal"
			@click.self="isDeleteModalVisible = false"
		>
			<div class="modal__dialog modal__dialog--confirm">
				<div class="confirm-icon">⚠️</div>
				<h3 class="modal__title text-center">Xác nhận xoá</h3>
				<p class="modal__subtitle text-center">
					Bạn có chắc chắn muốn xoá nhân viên <strong>{{ pendingDeleteEmployee?.name }}</strong>? Hành động này không thể hoàn tác.
				</p>
				<div class="confirm-actions">
					<button
						class="btn btn--secondary flex-1"
						@click="isDeleteModalVisible = false"
					>
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
import { ref, computed, onMounted } from "vue";
import { useEmployeeStore } from "@/stores/employee";
import { useDepartmentStore } from "@/stores/department";
import { useToast } from "vue-toastification";
import { useAuthStore } from "@/stores/auth";
import EmployeeForm from "@/components/EmployeeForm.vue";

// --- Tài nguyên và Biểu tượng ---
import plusIcon from "@/assets/icons/plus.svg";
import searchIcon from "@/assets/icons/search.svg";
import editIcon from "@/assets/icons/edit.svg";
import deleteIcon from "@/assets/icons/delete.svg";
import prevIcon from "@/assets/icons/chevron-left.svg";
import nextIcon from "@/assets/icons/chevron-right.svg";

// --- Cấu hình và Hằng số ---
const ITEMS_PER_PAGE = 7;

// --- Store và Service ---
const employeeStore = useEmployeeStore();
const deptStore = useDepartmentStore();
const authStore = useAuthStore();
const toast = useToast();

// --- Trạng thái Reactive ---
const searchKeyword = ref("");
const selectedDeptId = ref("");
const currentPage = ref(1);
const isFormVisible = ref(false);
const isDeleteModalVisible = ref(false);
const editingEmployee = ref(null);
const pendingDeleteEmployee = ref(null);
let debounceTimer = null;

// --- Thuộc tính tính toán (Computed) ---

/**
 * Kiểm tra quyền quản trị của người dùng hiện tại
 */
const isAdmin = computed(() => authStore.user?.role === "admin");

/**
 * Tính toán tổng số trang dựa trên dữ liệu từ store
 */
const totalPages = computed(
	() => Math.ceil(employeeStore.total / ITEMS_PER_PAGE) || 1,
);

/**
 * Kiểm tra danh sách nhân viên có trống hay không
 */
const isEmployeeListEmpty = computed(() => {
	const employees = employeeStore.employees;
	return !employees || employees.length === 0;
});

/**
 * Kiểm tra xem có cần hiển thị phân trang hay không
 */
const shouldShowPagination = computed(() => totalPages.value > 1);

// --- Xử lý Logic nghiệp vụ ---

/**
 * Tải danh sách nhân viên với các tham số lọc và phân trang hiện tại
 */
async function fetchEmployeeData() {
	try {
		await employeeStore.fetchEmployees({
			page: currentPage.value,
			limit: ITEMS_PER_PAGE,
			search: searchKeyword.value,
			department_id: selectedDeptId.value || undefined,
		});
	} catch (error) {
		toast.error("Không thể tải danh sách nhân viên");
		console.error("[EmployeeList] Fetch error:", error);
	}
}

/**
 * Đưa trang về 1 và làm mới danh sách (dùng khi tìm kiếm hoặc lọc)
 */
function refreshToFirstPage() {
	currentPage.value = 1;
	fetchEmployeeData();
}

/**
 * Xử lý sự kiện nhập ô tìm kiếm với cơ chế debounce 300ms
 */
function handleSearchInput() {
	if (debounceTimer) clearTimeout(debounceTimer);
	debounceTimer = setTimeout(refreshToFirstPage, 300);
}

/**
 * Chuyển đến trang trước đó
 */
function goToPrevPage() {
	if (currentPage.value <= 1) return;
	currentPage.value--;
	fetchEmployeeData();
}

/**
 * Chuyển đến trang tiếp theo
 */
function goToNextPage() {
	if (currentPage.value >= totalPages.value) return;
	currentPage.value++;
	fetchEmployeeData();
}

/**
 * Mở form thêm mới hoặc cập nhật nhân viên
 */
function openEmployeeForm(employee = null) {
	editingEmployee.value = employee;
	isFormVisible.value = true;
}

/**
 * Xử lý sau khi lưu dữ liệu nhân viên thành công
 */
function handleSaveSuccess() {
	isFormVisible.value = false;
	const message = editingEmployee.value
		? "Cập nhật nhân viên thành công!"
		: "Thêm nhân viên mới thành công!";

	toast.success(message);
	fetchEmployeeData();
}

/**
 * Chuẩn bị thông tin và hiển thị modal xác nhận xoá
 */
function prepareDelete(id, name) {
	pendingDeleteEmployee.value = { id, name };
	isDeleteModalVisible.value = true;
}

/**
 * Thực hiện hành động xoá nhân viên sau khi xác nhận
 */
async function executeDelete() {
	const employee = pendingDeleteEmployee.value;
	if (!employee) return;

	try {
		await employeeStore.deleteEmployee(employee.id);
		toast.success(`Đã xoá nhân viên "${employee.name}"`);
		isDeleteModalVisible.value = false;
		fetchEmployeeData();
	} catch (error) {
		toast.error("Có lỗi xảy ra khi thực hiện xoá");
		console.error("[EmployeeList] Delete error:", error);
	}
}

// --- Xử lý Format dữ liệu ---

/**
 * Lấy chữ cái đầu của tên để làm avatar mặc định
 */
function getInitials(name) {
	if (!name) return "??";

	return name
		.trim()
		.split(" ")
		.map(word => word[0])
		.slice(-2)
		.join("")
		.toUpperCase();
}

/**
 * Định dạng giá trị lương sang chuẩn tiền tệ VNĐ
 */
function formatCurrency(amount) {
	if (!amount) return "—";

	return new Intl.NumberFormat("vi-VN", {
		style: "currency",
		currency: "VND",
	}).format(amount);
}

/**
 * Chuyển đổi mã giới tính sang nhãn hiển thị tiếng Việt
 */
function getGenderLabel(genderCode) {
	const labels = { male: "Nam", female: "Nữ", other: "Khác" };
	return labels[genderCode] || "Chưa rõ";
}

/**
 * Định dạng chuỗi ngày tháng sang dd/mm/yyyy
 */
function formatVietnameseDate(dateStr) {
	if (!dateStr) return "—";

	return new Date(dateStr).toLocaleDateString("vi-VN", {
		day: "2-digit",
		month: "2-digit",
		year: "numeric",
	});
}

// --- Lifecycle Hooks ---

onMounted(() => {
	deptStore.fetchDepartments();
	fetchEmployeeData();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Biến thiết kế (Design Tokens) ===== */
.employee-view {
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
	--shadow-md: 0 4px 12px rgba(59, 130, 246, 0.25);
	--shadow-modal: 0 25px 50px -12px rgba(0, 0, 0, 0.25);

	font-family: "Inter", sans-serif;
	color: var(--text-main);
}

/* ===== Tiện ích (Utilities) ===== */
.text-center { text-align: center; }
.text-right { text-align: right; }
.text-muted { color: var(--text-muted); font-size: 0.875rem; }
.text-main { color: var(--text-main); }
.fw-500 { font-weight: 500; }
.flex-1 { flex: 1; }

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

.btn--primary {
	background: var(--primary-color);
	color: white;
	box-shadow: var(--shadow-md);
}

.btn--primary:hover {
	background: var(--primary-hover);
	transform: translateY(-1px);
	box-shadow: 0 6px 16px rgba(59, 130, 246, 0.3);
}

.btn--secondary {
	background: var(--bg-main);
	color: var(--text-muted);
	border: 1px solid var(--border-color);
}

.btn--secondary:hover {
	background: var(--border-hover);
	color: var(--text-main);
}

.btn--danger {
	background: var(--danger-color);
	color: white;
	box-shadow: 0 4px 12px rgba(239, 68, 68, 0.25);
}

.btn--danger:hover {
	background: var(--danger-hover);
	transform: translateY(-1px);
	box-shadow: 0 6px 16px rgba(239, 68, 68, 0.35);
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
	color: var(--text-muted);
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

/* ===== Form Controls ===== */
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
	cursor: pointer;
	min-width: 200px;
}

/* ===== Layout: Header ===== */
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

/* ===== Layout: Card & Toolbar ===== */
.content-card {
	background: var(--bg-card);
	border-radius: var(--radius-lg);
	border: 1px solid var(--border-color);
	box-shadow: var(--shadow-sm);
	overflow: hidden;
}

.toolbar {
	display: flex;
	justify-content: space-between;
	gap: 1.5rem;
	padding: 1.5rem;
	border-bottom: 1px solid var(--bg-main);
	flex-wrap: wrap;
}

.search-box {
	position: relative;
	flex: 1;
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
	z-index: 1;
}

.search-box__input {
	padding-left: 2.75rem;
}

/* ===== Layout: Data Table ===== */
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
	padding: 0.6rem 1.5rem;
	border-bottom: 1px solid var(--bg-main);
}

.data-table tr:last-child td {
	border-bottom: none;
}

.user-info {
	display: flex;
	align-items: center;
	gap: 12px;
}

.user-info__avatar-wrap {
	flex-shrink: 0;
}

.user-info__avatar-img {
	width: 44px;
	height: 44px;
	border-radius: var(--radius-md);
	object-fit: cover;
}

.user-info__initials {
	width: 44px;
	height: 44px;
	border-radius: var(--radius-md);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 700;
	font-size: 0.875rem;
}

.user-info__initials--male { background: #eff6ff; color: var(--primary-color); }
.user-info__initials--female { background: #fdf2f8; color: #db2777; }
.user-info__initials--other { background: #f5f3ff; color: #7c3aed; }

.user-info__details {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.user-info__name {
	font-weight: 600;
	color: var(--text-main);
}

.user-info__email {
	font-size: 0.8125rem;
	color: var(--text-muted);
}

/* Tags and Badges */
.tag {
	padding: 0.25rem 0.75rem;
	border-radius: 999px;
	font-size: 0.75rem;
	font-weight: 600;
	display: inline-block;
}

.tag--male { background: #eff6ff; color: var(--primary-color); }
.tag--female { background: #fdf2f8; color: #db2777; }
.tag--other { background: var(--bg-main); color: var(--text-muted); }

.status-badge {
	padding: 0.25rem 0.75rem;
	border-radius: var(--radius-sm);
	font-size: 0.75rem;
	font-weight: 700;
	display: inline-block;
}

.status-badge--active { background: #dcfce7; color: #15803d; }
.status-badge--inactive { background: #fee2e2; color: #b91c1c; }

.job-info {
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.job-info__dept {
	font-weight: 600;
	font-size: 0.9rem;
}

.job-info__pos {
	font-size: 0.8125rem;
	color: var(--text-muted);
}

.salary-amount {
	font-family: "JetBrains Mono", monospace;
	font-weight: 600;
	color: #059669;
}

.action-group {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
}

/* ===== Pagination ===== */
.pagination {
	padding: 1.5rem;
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 1.5rem;
	border-top: 1px solid var(--bg-main);
}

.pagination__btn {
	width: 40px;
	height: 40px;
	border-radius: var(--radius-md);
	border: 1px solid var(--border-color);
	background: white;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: all 0.2s;
	color: var(--text-muted);
}

.pagination__btn img {
	width: 18px;
	height: 18px;
	filter: invert(47%) sepia(13%) saturate(545%) hue-rotate(183deg) brightness(95%) contrast(87%);
}

.pagination__btn:hover:not(:disabled) {
	background: var(--bg-main);
	color: var(--primary-color);
	border-color: var(--primary-color);
}

.pagination__btn:hover:not(:disabled) img {
	filter: invert(44%) sepia(91%) saturate(1185%) hue-rotate(200deg) brightness(101%) contrast(92%);
}

.pagination__btn:disabled {
	opacity: 0.4;
	cursor: not-allowed;
}

.pagination__info {
	font-size: 0.9rem;
	color: var(--text-muted);
}

.pagination__info span {
	font-weight: 700;
	color: var(--text-main);
}

/* ===== Skeletons & Empty State ===== */
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

@keyframes pulse {
	0%, 100% { opacity: 1; }
	50% { opacity: 0.5; }
}

.empty-state {
	padding: 4rem;
	text-align: center;
	color: #94a3b8;
}

.empty-state__icon {
	font-size: 3rem;
	margin-bottom: 1rem;
}

.empty-state__text {
	margin: 0;
}

/* ===== Modal Styles ===== */
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
	border-radius: var(--radius-lg);
	width: 100%;
	max-width: 500px;
	max-height: 90vh;
	box-shadow: var(--shadow-modal);
	display: flex;
	flex-direction: column;
	overflow: hidden;
	animation: modalIn 0.3s ease-out;
}

.modal__dialog--confirm {
	max-width: 400px;
	padding: 2rem;
}

.modal__title {
	font-size: 1.5rem;
	font-weight: 700;
	margin: 0 0 0.75rem 0;
}

.modal__subtitle {
	color: var(--text-muted);
	line-height: 1.6;
	margin: 0 0 2rem 0;
}

@keyframes modalIn {
	from { opacity: 0; transform: translateY(20px); }
	to { opacity: 1; transform: translateY(0); }
}

.confirm-icon {
	font-size: 3.5rem;
	margin-bottom: 1rem;
	text-align: center;
}

.confirm-actions {
	display: flex;
	gap: 12px;
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
	.page-header {
		flex-direction: column;
		align-items: flex-start;
		gap: 1.25rem;
	}

	.btn--primary {
		width: 100%;
		justify-content: center;
	}

	.toolbar {
		flex-direction: column;
		gap: 1rem;
		padding: 1.25rem;
	}

	.search-box {
		max-width: none;
		width: 100%;
	}

	.form-select {
		width: 100%;
		min-width: 0;
	}

	.user-info__name {
		font-size: 0.875rem;
	}

	.data-table td {
		padding: 1rem;
	}

	.action-group {
		gap: 4px;
	}

	.btn-icon {
		width: 32px;
		height: 32px;
	}
}

@media (max-width: 480px) {
	.page-title {
		font-size: 1.5rem;
	}

	.pagination {
		flex-direction: column;
		gap: 1rem;
		padding: 1.25rem;
	}

	.content-card {
		border-radius: 0;
		border-left: none;
		border-right: none;
	}
}
</style>
