<template>
	<!-- ===== Modal Overlay ===== -->
	<div class="modal-overlay" @click.self="$emit('close')">
		<div class="modal-container">
			<!-- ===== Header ===== -->
			<div class="modal-header">
				<div class="header-info">
					<h2>{{ isEdit ? "Chỉnh sửa nhân viên" : "Thêm nhân viên mới" }}</h2>
					<p class="header-desc">Điền đầy đủ thông tin nhân viên vào hệ thống</p>
				</div>
				<button class="btn-icon" @click="$emit('close')">
					<img :src="closeIcon" alt="close" />
				</button>
			</div>

			<!-- ===== Form Content ===== -->
			<form @submit.prevent="handleSubmit" class="modal-form">
				<div class="form-body">
					<!-- Section: Thông tin cơ bản -->
					<div class="form-section">
						<h3 class="section-title">Thông tin cơ bản</h3>
						<div class="form-grid">
							<!-- Họ và tên -->
							<div class="form-group col-span-2">
								<label>Họ và tên <span class="text-danger">*</span></label>
								<div class="input-group">
									<img :src="userIcon" class="input-icon" />
									<input
										v-model="form.name"
										placeholder="Nguyễn Văn A"
										class="form-control"
									/>
								</div>
								<span class="error-text" v-if="errors.name">{{ errors.name }}</span>
							</div>

							<!-- Giới tính -->
							<div class="form-group">
								<label>Giới tính</label>
								<div class="radio-group">
									<label class="radio-item">
										<input type="radio" v-model="form.gender" value="male" />
										<span class="radio-label">Nam</span>
									</label>
									<label class="radio-item">
										<input type="radio" v-model="form.gender" value="female" />
										<span class="radio-label">Nữ</span>
									</label>
								</div>
							</div>

							<!-- Ngày sinh -->
							<div class="form-group">
								<label>Ngày sinh</label>
								<div class="input-group">
									<img :src="calendarIcon" class="input-icon" />
									<input v-model="form.date_of_birth" type="date" class="form-control" />
								</div>
							</div>

							<!-- Số điện thoại -->
							<div class="form-group">
								<label>Số điện thoại <span class="text-danger">*</span></label>
								<div class="input-group">
									<img :src="phoneIcon" class="input-icon" />
									<input v-model="form.phone" placeholder="0901234567" class="form-control" />
								</div>
								<span class="error-text" v-if="errors.phone">{{ errors.phone }}</span>
							</div>
						</div>
					</div>

					<!-- Section: Công việc & Lương -->
					<div class="form-section">
						<h3 class="section-title">Công việc & Lương</h3>
						<div class="form-grid">
							<!-- Phòng ban -->
							<div class="form-group">
								<label>Phòng ban</label>
								<div class="input-group">
									<img :src="deptIcon" class="input-icon" />
									<select v-model="form.department_id" class="form-control form-select">
										<option :value="null">Chọn phòng ban</option>
										<option v-for="d in deptStore.departments" :key="d.id" :value="d.id">
											{{ d.name }}
										</option>
									</select>
								</div>
							</div>

							<!-- Chức vụ -->
							<div class="form-group">
								<label>Chức vụ</label>
								<div class="input-group">
									<img :src="positionIcon" class="input-icon" />
									<input v-model="form.position" placeholder="Ví dụ: Kế toán trưởng" class="form-control" />
								</div>
							</div>

							<!-- Mức lương -->
							<div class="form-group">
								<label>Mức lương (VNĐ) <span class="text-danger">*</span></label>
								<div class="input-group">
									<img :src="salaryIcon" class="input-icon" />
									<input v-model.number="form.salary" type="number" step="100000" placeholder="15000000" class="form-control" />
								</div>
								<span class="error-text" v-if="errors.salary">{{ errors.salary }}</span>
							</div>

							<!-- Ngày vào làm -->
							<div class="form-group">
								<label>Ngày vào làm</label>
								<div class="input-group">
									<img :src="calendarIcon" class="input-icon" />
									<input v-model="form.hire_date" type="date" class="form-control" />
								</div>
							</div>

							<!-- Trạng thái (chỉ hiện khi edit) -->
							<div class="form-group col-span-2" v-if="isEdit">
								<label>Trạng thái nhân sự</label>
								<div class="input-group">
									<select v-model="form.status" class="form-control form-select">
										<option :value="1">Đang làm việc</option>
										<option :value="0">Đã nghỉ việc</option>
									</select>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- ===== Footer ===== -->
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary" @click="$emit('close')">
						Huỷ bỏ
					</button>
					<button type="submit" class="btn btn-primary" :disabled="submitting">
						<span v-if="!submitting">{{ isEdit ? "Lưu thay đổi" : "Tạo nhân viên" }}</span>
						<span v-else class="loader"></span>
					</button>
				</div>
			</form>
		</div>
	</div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from "vue";
import { useEmployeeStore } from "@/stores/employee";
import { useDepartmentStore } from "@/stores/department";
import { useToast } from "vue-toastification";

// Import Icons
import closeIcon from "@/assets/icons/close.svg";
import userIcon from "@/assets/icons/user.svg";
import phoneIcon from "@/assets/icons/phone.svg";
import deptIcon from "@/assets/icons/department.svg";
import positionIcon from "@/assets/icons/position.svg";
import salaryIcon from "@/assets/icons/salary.svg";
import calendarIcon from "@/assets/icons/calendar.svg";

const props = defineProps({ employee: { type: Object, default: null } });
const emit = defineEmits(["close", "saved"]);

// --- Quản lý state và khởi tạo ---
const store = useEmployeeStore();
const deptStore = useDepartmentStore();
const toast = useToast();

const submitting = ref(false);
const isEdit = computed(() => !!props.employee);

/**
 * Trả về dữ liệu mặc định của form
 */
function getDefaultForm() {
	return {
		name: "",
		gender: "male",
		phone: "",
		department_id: null,
		position: "",
		salary: null,
		hire_date: new Date().toISOString().substring(0, 10),
		date_of_birth: "",
		status: 1,
	};
}

const form = reactive(getDefaultForm());
const errors = reactive({});

// --- Xử lý format dữ liệu ---

/**
 * Lấy định dạng ngày YYYY-MM-DD từ chuỗi ISO
 */
function formatDateString(dateStr) {
	return dateStr ? dateStr.substring(0, 10) : "";
}

/**
 * Cập nhật form với dữ liệu nhân viên (nếu có) hoặc reset về mặc định
 */
function populateForm(employeeData) {
	clearErrors();

	if (!employeeData) {
		Object.assign(form, getDefaultForm());
		return;
	}

	Object.assign(form, {
		name: employeeData.name || "",
		gender: employeeData.gender || "male",
		phone: employeeData.phone || "",
		department_id: employeeData.department_id || null,
		position: employeeData.position || "",
		salary: employeeData.salary || 0,
		hire_date: formatDateString(employeeData.hire_date),
		date_of_birth: formatDateString(employeeData.date_of_birth),
		status: employeeData.status,
	});
}

// --- Validation và Xử lý lỗi ---

/**
 * Xóa tất cả các lỗi đang hiển thị
 */
function clearErrors() {
	Object.keys(errors).forEach(key => delete errors[key]);
}

/**
 * Kiểm tra tính hợp lệ của dữ liệu form trước khi gửi
 */
function validateForm() {
	clearErrors();
	
	const isNameEmpty = !form.name.trim();
	const isPhoneEmpty = !form.phone.trim();
	const isSalaryInvalid = !form.salary || form.salary < 0;

	if (isNameEmpty) errors.name = "Họ tên không được để trống";
	if (isPhoneEmpty) errors.phone = "Số điện thoại không được để trống";
	if (isSalaryInvalid) errors.salary = "Mức lương không hợp lệ";
	
	const hasErrors = Object.keys(errors).length > 0;
	return !hasErrors;
}

/**
 * Xử lý thông báo lỗi từ API
 */
function handleApiError(error) {
	console.error(error);
	const errorMessage = error.response?.data?.error || "Có lỗi xảy ra, vui lòng thử lại";
	
	const isPhoneError = errorMessage.includes("số điện thoại");
	if (isPhoneError) {
		errors.phone = errorMessage;
	} else {
		toast.error(errorMessage);
	}
}

// --- Chuẩn bị payload API ---

/**
 * Tạo payload chứa các trường bị thay đổi để cập nhật
 */
function buildUpdatePayload() {
	const payload = {};
	const emp = props.employee;
	
	const originalHireDate = formatDateString(emp.hire_date);
	const originalDob = formatDateString(emp.date_of_birth);

	if (form.name !== emp.name) payload.name = form.name;
	if (form.gender !== emp.gender) payload.gender = form.gender;
	if (form.phone !== emp.phone) payload.phone = form.phone;
	if (form.department_id !== emp.department_id) payload.department_id = form.department_id;
	if (form.position !== emp.position) payload.position = form.position;
	if (form.salary !== emp.salary) payload.salary = form.salary;
	if (form.hire_date !== originalHireDate) payload.hire_date = form.hire_date;
	if (form.date_of_birth !== originalDob) payload.date_of_birth = form.date_of_birth;
	if (form.status !== emp.status) payload.status = String(form.status);

	return payload;
}

/**
 * Tạo payload đầy đủ để thêm nhân viên mới
 */
function buildCreatePayload() {
	return {
		...form,
		status: String(form.status)
	};
}

// --- Xử lý sự kiện chính ---

/**
 * Gửi dữ liệu lên server (Thêm mới hoặc Cập nhật)
 */
async function handleSubmit() {
	const isValid = validateForm();
	if (!isValid) return;
	
	submitting.value = true;
	
	try {
		if (isEdit.value) {
			const payload = buildUpdatePayload();
			const hasChanges = Object.keys(payload).length > 0;
			
			if (!hasChanges) {
				emit("close");
				return;
			}
			
			await store.updateEmployee(props.employee.id, payload);
		} else {
			const payload = buildCreatePayload();
			await store.createEmployee(payload);
		}
		
		emit("saved");
	} catch (error) {
		handleApiError(error);
	} finally {
		submitting.value = false;
	}
}

// --- Lifecycle hooks ---

onMounted(async () => {
	await deptStore.fetchDepartments();
});

watch(
	() => props.employee,
	(employeeData) => populateForm(employeeData),
	{ immediate: true }
);
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Modal Overlay & Container ===== */
.modal-overlay {
	position: fixed;
	inset: 0;
	background: rgba(15, 23, 42, 0.4);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: 1.5rem;
}

.modal-container {
	background: #ffffff;
	border-radius: 24px;
	width: 100%;
	max-width: 680px;
	max-height: 90vh;
	box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.15);
	display: flex;
	flex-direction: column;
	overflow: hidden;
	animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
	from {
		transform: translateY(20px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}

/* ===== Header ===== */
.modal-header {
	padding: 1.75rem 2rem;
	border-bottom: 1px solid #f1f5f9;
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
}

.header-info h2 {
	font-size: 1.25rem;
	font-weight: 700;
	color: #1e293b;
	margin-bottom: 4px;
}

.header-desc {
	font-size: 0.875rem;
	color: #64748b;
}

/* ===== Reusable Utility Classes ===== */
.btn-icon {
	background: #f8fafc;
	border: none;
	width: 36px;
	height: 36px;
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	color: #64748b;
	transition: all 0.2s;
}

.btn-icon:hover {
	background: #fee2e2;
	color: #ef4444;
}

.text-danger {
	color: #ef4444;
}

/* ===== Form Layout ===== */
.modal-form {
	flex: 1;
	overflow-y: auto;
}

.form-body {
	padding: 2rem;
	display: flex;
	flex-direction: column;
	gap: 2rem;
}

.section-title {
	font-size: 0.875rem;
	font-weight: 700;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: #94a3b8;
	margin-bottom: 1.25rem;
}

.form-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 1.5rem;
}

.col-span-2 {
	grid-column: span 2;
}

/* ===== Form Controls ===== */
.form-group {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.form-group label {
	font-size: 0.875rem;
	font-weight: 600;
	color: #475569;
}

.input-group {
	position: relative;
	display: flex;
	align-items: center;
}

.input-icon {
	position: absolute;
	left: 1rem;
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%) contrast(89%);
}

.form-control {
	width: 100%;
	padding: 0.75rem 1rem 0.75rem 2.75rem;
	border: 1px solid #e2e8f0;
	border-radius: 12px;
	font-size: 0.95rem;
	background: #f8fafc;
	transition: all 0.2s;
	outline: none;
}

.form-control:focus {
	background: #ffffff;
	border-color: #3b82f6;
	box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
}

.form-select {
	appearance: none;
	background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%2394a3b8'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E");
	background-repeat: no-repeat;
	background-position: right 1rem center;
	background-size: 1rem;
}

/* ===== Radio Group ===== */
.radio-group {
	display: flex;
	gap: 0.75rem;
	padding-top: 0.25rem;
}

.radio-item {
	flex: 1;
	cursor: pointer;
}

.radio-item input {
	display: none;
}

.radio-label {
	display: block;
	text-align: center;
	padding: 0.625rem;
	background: #f1f5f9;
	border: 1px solid #e2e8f0;
	border-radius: 10px;
	font-size: 0.875rem;
	font-weight: 500;
	color: #64748b;
	transition: all 0.2s;
}

.radio-item input:checked + .radio-label {
	background: #eff6ff;
	border-color: #3b82f6;
	color: #3b82f6;
	font-weight: 700;
}

/* ===== Form Validation ===== */
.error-text {
	font-size: 0.75rem;
	color: #ef4444;
	font-weight: 500;
}

/* ===== Modal Footer ===== */
.modal-footer {
	padding: 1.5rem 2rem;
	border-top: 1px solid #f1f5f9;
	display: flex;
	justify-content: flex-end;
	gap: 1rem;
}

.btn {
	padding: 0.75rem 1.5rem;
	border-radius: 12px;
	font-weight: 600;
	cursor: pointer;
	transition: all 0.2s;
	display: inline-flex;
	align-items: center;
	justify-content: center;
	border: none;
}

.btn-secondary {
	background: #f8fafc;
	border: 1px solid #e2e8f0;
	color: #64748b;
}

.btn-secondary:hover {
	background: #f1f5f9;
	color: #1e293b;
}

.btn-primary {
	background: #3b82f6;
	color: #ffffff;
	min-width: 140px;
}

.btn-primary:hover:not(:disabled) {
	background: #2563eb;
	transform: translateY(-1px);
}

.btn-primary:disabled {
	opacity: 0.7;
	cursor: not-allowed;
}

/* ===== Loader ===== */
.loader {
	width: 20px;
	height: 20px;
	border: 3px solid rgba(255, 255, 255, 0.3);
	border-top-color: #ffffff;
	border-radius: 50%;
	display: inline-block;
	animation: spin 0.8s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

/* ===== Responsive ===== */
@media (max-width: 640px) {
	.form-grid {
		grid-template-columns: 1fr;
	}
	.col-span-2 {
		grid-column: span 1;
	}
	.modal-container {
		max-height: 100vh;
		border-radius: 0;
	}
}
</style>
