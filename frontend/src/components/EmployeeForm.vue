<template>
	<div class="form-body">
		<!-- Section: Thông tin cơ bản -->
		<div class="form-section">
			<h3 class="section-title">Thông tin cơ bản</h3>
			<div class="form-grid">
				<!-- Họ và tên -->
				<div class="form-group col-span-2">
					<label>Họ và tên <span class="required">*</span></label>
					<div class="input-group">
						<img :src="userIcon" class="input-group__icon" />
						<input
							v-model="form.name"
							placeholder="Nguyễn Văn A"
							class="form-control input-group__input"
						/>
					</div>
					<span class="error-text" v-if="errors.name">{{
						errors.name
					}}</span>
				</div>

				<!-- Giới tính -->
				<div class="form-group">
					<label>Giới tính</label>
					<div class="radio-group">
						<label class="radio-item">
							<input
								type="radio"
								v-model="form.gender"
								value="male"
							/>
							<span class="radio-label">Nam</span>
						</label>
						<label class="radio-item">
							<input
								type="radio"
								v-model="form.gender"
								value="female"
							/>
							<span class="radio-label">Nữ</span>
						</label>
					</div>
				</div>

				<!-- Ngày sinh -->
				<div class="form-group">
					<label>Ngày sinh</label>
					<div class="input-group">
						<img :src="calendarIcon" class="input-group__icon" />
						<input
							v-model="form.date_of_birth"
							type="date"
							class="form-control input-group__input"
						/>
					</div>
				</div>

				<!-- Số điện thoại -->
				<div class="form-group">
					<label>Số điện thoại <span class="required">*</span></label>
					<div class="input-group">
						<img :src="phoneIcon" class="input-group__icon" />
						<input
							v-model="form.phone"
							placeholder="0901234567"
							class="form-control input-group__input"
						/>
					</div>
					<span class="error-text" v-if="errors.phone">{{
						errors.phone
					}}</span>
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
						<img :src="deptIcon" class="input-group__icon" />
						<select
							v-model="form.department_id"
							class="form-control form-select input-group__input"
						>
							<option :value="null">Chọn phòng ban</option>
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

				<!-- Chức vụ -->
				<div class="form-group">
					<label>Chức vụ</label>
					<div class="input-group">
						<img :src="positionIcon" class="input-group__icon" />
						<select
							v-model="form.position_id"
							class="form-control form-select input-group__input"
							:disabled="!form.department_id"
						>
							<option :value="null">
								{{ form.department_id ? "Chọn chức vụ" : "Vui lòng chọn phòng ban trước" }}
							</option>
							<option
								v-for="p in availablePositions"
								:key="p.id"
								:value="p.id"
							>
								{{ p.name }}
							</option>
						</select>
					</div>
				</div>

				<!-- Mức lương -->
				<div class="form-group">
					<label
						>Mức lương cơ bản <span class="required">*</span></label
					>
					<div class="input-group">
						<img :src="salaryIcon" class="input-group__icon" />
						<input
							v-model.number="form.salary"
							type="number"
							step="100000"
							placeholder="15000000"
							class="form-control input-group__input"
						/>
					</div>
					<span class="error-text" v-if="errors.salary">{{
						errors.salary
					}}</span>
				</div>

				<!-- Ngày vào làm -->
				<div class="form-group">
					<label>Ngày vào làm</label>
					<div class="input-group">
						<img :src="calendarIcon" class="input-group__icon" />
						<input
							v-model="form.hire_date"
							type="date"
							class="form-control input-group__input"
						/>
					</div>
				</div>

				<!-- Trạng thái (chỉ hiện khi edit) -->
				<div class="form-group col-span-2" v-if="isEdit">
					<label>Trạng thái nhân sự</label>
					<div class="input-group">
						<select
							v-model="form.status"
							class="form-control form-select input-group__input"
						>
							<option :value="1">Đang làm việc</option>
							<option :value="0">Đã nghỉ việc</option>
						</select>
					</div>
				</div>

				<!-- Ảnh đại diện (URL) -->
				<div class="form-group col-span-2">
					<label>URL Ảnh đại diện</label>
					<div class="input-group">
						<img :src="userIcon" class="input-group__icon" />
						<input
							v-model="form.avatar_url"
							placeholder="https://example.com/avatar.jpg"
							class="form-control input-group__input"
						/>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from "vue";
import { useEmployeeStore } from "@/stores/employee";
import { useDepartmentStore } from "@/stores/department";
import { useToast } from "vue-toastification";
import BaseModal from "./BaseModal.vue";

// Import Icons
import userIcon from "@/assets/icons/user.svg";
import phoneIcon from "@/assets/icons/phone.svg";
import deptIcon from "@/assets/icons/department.svg";
import positionIcon from "@/assets/icons/position.svg";
import salaryIcon from "@/assets/icons/salary.svg";
import calendarIcon from "@/assets/icons/calendar.svg";

const props = defineProps({
	employee: { type: Object, default: null },
});
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
		position_id: null,
		salary: null,
		hire_date: new Date().toISOString().substring(0, 10),
		date_of_birth: "",
		status: 1,
		avatar_url: "",
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
		position_id: employeeData.position_id || null,
		salary: employeeData.salary || 0,
		hire_date: formatDateString(employeeData.hire_date),
		date_of_birth: formatDateString(employeeData.date_of_birth),
		status: employeeData.status,
		avatar_url: employeeData.avatar_url || "",
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
	const isSalaryInvalid = form.salary === null || form.salary < 0;

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
	const errorMessage =
		error.response?.data?.error || "Có lỗi xảy ra, vui lòng thử lại";

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
	if (form.department_id !== emp.department_id)
		payload.department_id = form.department_id;
	if (form.position_id !== emp.position_id)
		payload.position_id = form.position_id;
	if (form.salary !== emp.salary) payload.salary = form.salary;
	if (form.hire_date !== originalHireDate) payload.hire_date = form.hire_date;
	if (form.date_of_birth !== originalDob)
		payload.date_of_birth = form.date_of_birth;
	if (form.status !== emp.status) payload.status = String(form.status);
	if (form.avatar_url !== emp.avatar_url)
		payload.avatar_url = form.avatar_url;

	return payload;
}

/**
 * Tạo payload đầy đủ để thêm nhân viên mới
 */
function buildCreatePayload() {
	return {
		...form,
		status: String(form.status),
	};
}

// --- Computed cho Chức vụ ---
const availablePositions = computed(() => {
	if (!form.department_id) return [];
	const dept = deptStore.departments.find((d) => d.id === form.department_id);
	return dept ? dept.positions || [] : [];
});

// Reset chức vụ khi đổi phòng ban
watch(
	() => form.department_id,
	(newDeptId, oldDeptId) => {
		// Chỉ reset nếu thay đổi thực sự và không phải lần đầu load data
		if (oldDeptId !== undefined && newDeptId !== oldDeptId) {
			form.position_id = null;
		}
	}
);

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

// Reset form khi employee thay đổi
watch(
	() => props.employee,
	employeeData => populateForm(employeeData),
	{ immediate: true },
);

defineExpose({
	handleSubmit,
	submitting,
	isEdit,
});
</script>

<style scoped>
/* Scoped styles kept only for grid and custom radio layout */
.form-body {
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
}
</style>
