package dto

type DeptStat struct {
	DepartmentName string `json:"name"`
	Count          int64  `json:"count"`
}

type DashboardResponse struct {
	TotalEmployees    int64      `json:"totalEmployees"`
	ActiveEmployees   int64      `json:"activeEmployees"`
	InactiveEmployees int64      `json:"inactiveEmployees"`
	TotalDepartments  int64      `json:"totalDepartments"`
	TotalUsers        int64      `json:"totalUsers"`
	EmployeesByDept   []DeptStat `json:"employeesByDepartment"`
	TotalAdminRole    int64      `json:"totalAdminRole"`
}
