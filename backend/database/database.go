package database

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/models"
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

func MigrateModels() {
	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.Department{},
		&models.Employee{},
		&models.User{},
	)

	if err != nil {
		panic("Failed to migrate models: " + err.Error())
	}

	fmt.Println("Auto migrate models successfully")
}

func Seed() {
	SeedRole()
	SeedDepartment()
	SeedEmployee()
	SeedAdmin()
	SeedUser()
}

func SeedRole() {
	roles := []string{"admin", "user"}

	for _, name := range roles {
		config.DB.FirstOrCreate(&models.Role{}, models.Role{Name: name})
	}
	fmt.Println("Seed roles successfully")
}

func SeedDepartment() {
	names := []string{"IT", "HR", "Finance", "Marketing", "Sales", "Operations"}

	for _, name := range names {
		var dept models.Department
		config.DB.Where("name = ?", name).FirstOrCreate(&dept, models.Department{Name: name})
	}

	fmt.Println("Seed departments successfully")
}

func SeedEmployee() {
	var count int64
	config.DB.Model(&models.Employee{}).Count(&count)
	if count > 0 {
		fmt.Println("Employees already seeded")
		return
	}

	departmentMap := map[string]uint{}
	var departments []models.Department
	config.DB.Find(&departments)

	for _, d := range departments {
		departmentMap[d.Name] = d.ID
	}

	getDept := func(name string) *uint {
		id, ok := departmentMap[name]
		if !ok {
			fmt.Println("Warning: department not found:", name)
			return nil
		}
		return &id
	}

	ptrTime := func(y, m, d int) *time.Time {
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		return &t
	}

	employees := []models.Employee{
		{Name: "Nguyễn Văn An", Gender: "male", Phone: "0901111111", DepartmentID: getDept("IT"), Position: "Senior Developer", Salary: 25000000, HireDate: ptrTime(2021, 3, 1), DateOfBirth: ptrTime(1992, 5, 15), AvatarURL: "https://ui-avatars.com/api/?name=Nguyen+Van+An&background=4F46E5&color=fff"},
		{Name: "Trần Thị Bình", Gender: "female", Phone: "0902222222", DepartmentID: getDept("HR"), Position: "HR Manager", Salary: 20000000, HireDate: ptrTime(2020, 6, 15), DateOfBirth: ptrTime(1988, 10, 20), AvatarURL: "https://ui-avatars.com/api/?name=Tran+Thi+Binh&background=EC4899&color=fff"},
		{Name: "Lê Hoàng Cường", Gender: "male", Phone: "0903333333", DepartmentID: getDept("Finance"), Position: "Accountant", Salary: 18000000, HireDate: ptrTime(2022, 1, 10), DateOfBirth: ptrTime(1995, 3, 12), AvatarURL: "https://ui-avatars.com/api/?name=Le+Hoang+Cuong&background=0EA5E9&color=fff"},
		{Name: "Phạm Thị Dung", Gender: "female", Phone: "0904444444", DepartmentID: getDept("Marketing"), Position: "Marketing Lead", Salary: 22000000, HireDate: ptrTime(2021, 9, 5), DateOfBirth: ptrTime(1993, 12, 1), AvatarURL: "https://ui-avatars.com/api/?name=Pham+Thi+Dung&background=F59E0B&color=fff"},
		{Name: "Hoàng Minh Đức", Gender: "male", Phone: "0905555555", DepartmentID: getDept("IT"), Position: "DevOps Engineer", Salary: 28000000, HireDate: ptrTime(2020, 12, 1), DateOfBirth: ptrTime(1990, 8, 25), AvatarURL: "https://ui-avatars.com/api/?name=Hoang+Minh+Duc&background=10B981&color=fff"},
		{Name: "Vũ Thị Hoa", Gender: "female", Phone: "0906666666", DepartmentID: getDept("Sales"), Position: "Sales Executive", Salary: 16000000, HireDate: ptrTime(2023, 2, 20), DateOfBirth: ptrTime(1998, 4, 10), AvatarURL: "https://ui-avatars.com/api/?name=Vu+Thi+Hoa&background=EF4444&color=fff"},
		{Name: "Đặng Văn Khoa", Gender: "male", Phone: "0907777777", DepartmentID: getDept("IT"), Position: "Frontend Developer", Salary: 20000000, HireDate: ptrTime(2022, 7, 1), DateOfBirth: ptrTime(1996, 11, 30), AvatarURL: "https://ui-avatars.com/api/?name=Dang+Van+Khoa&background=8B5CF6&color=fff"},
		{Name: "Bùi Thị Lan", Gender: "female", Phone: "0908888888", DepartmentID: getDept("Operations"), Position: "Operations Manager", Salary: 21000000, HireDate: ptrTime(2019, 5, 15), DateOfBirth: ptrTime(1985, 2, 14), AvatarURL: "https://ui-avatars.com/api/?name=Bui+Thi+Lan&background=06B6D4&color=fff"},
		{Name: "Ngô Văn Long", Gender: "male", Phone: "0911111111", DepartmentID: getDept("IT"), Position: "Backend Developer", Salary: 23000000, HireDate: ptrTime(2021, 4, 12), DateOfBirth: ptrTime(1994, 7, 22), AvatarURL: "https://ui-avatars.com/api/?name=Ngo+Van+Long&background=2563EB&color=fff"},
		{Name: "Đỗ Thị Mai", Gender: "female", Phone: "0912222222", DepartmentID: getDept("HR"), Position: "Recruitment Specialist", Salary: 17000000, HireDate: ptrTime(2022, 8, 3), DateOfBirth: ptrTime(1997, 1, 15), AvatarURL: "https://ui-avatars.com/api/?name=Do+Thi+Mai&background=DB2777&color=fff"},
		{Name: "Phan Quốc Bảo", Gender: "male", Phone: "0913333333", DepartmentID: getDept("Finance"), Position: "Financial Analyst", Salary: 24000000, HireDate: ptrTime(2020, 11, 18), DateOfBirth: ptrTime(1991, 6, 18), AvatarURL: "https://ui-avatars.com/api/?name=Phan+Quoc+Bao&background=0891B2&color=fff"},
		{Name: "Lý Thị Hạnh", Gender: "female", Phone: "0914444444", DepartmentID: getDept("Marketing"), Position: "Content Creator", Salary: 15000000, HireDate: ptrTime(2023, 1, 9), DateOfBirth: ptrTime(1999, 9, 9), AvatarURL: "https://ui-avatars.com/api/?name=Ly+Thi+Hanh&background=F97316&color=fff"},
		{Name: "Trương Minh Hải", Gender: "male", Phone: "0915555555", DepartmentID: getDept("IT"), Position: "QA Engineer", Salary: 19000000, HireDate: ptrTime(2021, 10, 20), DateOfBirth: ptrTime(1993, 4, 5), AvatarURL: "https://ui-avatars.com/api/?name=Truong+Minh+Hai&background=16A34A&color=fff"},
		{Name: "Nguyễn Thị Kiều", Gender: "female", Phone: "0916666666", DepartmentID: getDept("Sales"), Position: "Sales Manager", Salary: 26000000, HireDate: ptrTime(2019, 7, 14), DateOfBirth: ptrTime(1987, 8, 28), AvatarURL: "https://ui-avatars.com/api/?name=Nguyen+Thi+Kieu&background=DC2626&color=fff"},
		{Name: "Mai Văn Nam", Gender: "male", Phone: "0917777777", DepartmentID: getDept("Operations"), Position: "Logistics Coordinator", Salary: 18000000, HireDate: ptrTime(2022, 5, 30), DateOfBirth: ptrTime(1996, 12, 12), AvatarURL: "https://ui-avatars.com/api/?name=Mai+Van+Nam&background=7C3AED&color=fff"},
		{Name: "Tạ Thị Oanh", Gender: "female", Phone: "0918888888", DepartmentID: getDept("Customer Service"), Position: "Support Executive", Salary: 14000000, HireDate: ptrTime(2023, 3, 12), DateOfBirth: ptrTime(2000, 5, 20), AvatarURL: "https://ui-avatars.com/api/?name=Ta+Thi+Oanh&background=0F766E&color=fff"},
		{Name: "Dương Quốc Phong", Gender: "male", Phone: "0921111111", DepartmentID: getDept("IT"), Position: "Mobile Developer", Salary: 27000000, HireDate: ptrTime(2020, 9, 25), DateOfBirth: ptrTime(1992, 11, 11), AvatarURL: "https://ui-avatars.com/api/?name=Duong+Quoc+Phong&background=1D4ED8&color=fff"},
		{Name: "Cao Thị Quyên", Gender: "female", Phone: "0922222222", DepartmentID: getDept("Finance"), Position: "Payroll Specialist", Salary: 17500000, HireDate: ptrTime(2021, 6, 8), DateOfBirth: ptrTime(1994, 3, 3), AvatarURL: "https://ui-avatars.com/api/?name=Cao+Thi+Quyen&background=BE185D&color=fff"},
		{Name: "Huỳnh Văn Sơn", Gender: "male", Phone: "0923333333", DepartmentID: getDept("Marketing"), Position: "SEO Specialist", Salary: 18500000, HireDate: ptrTime(2022, 4, 17), DateOfBirth: ptrTime(1995, 10, 10), AvatarURL: "https://ui-avatars.com/api/?name=Huynh+Van+Son&background=0284C7&color=fff"},
		{Name: "Võ Thị Thảo", Gender: "female", Phone: "0924444444", DepartmentID: getDept("Operations"), Position: "Project Coordinator", Salary: 20000000, HireDate: ptrTime(2020, 2, 11), DateOfBirth: ptrTime(1989, 5, 5), AvatarURL: "https://ui-avatars.com/api/?name=Vo+Thi+Thao&background=EA580C&color=fff"},
		{Name: "Đinh Minh Tuấn", Gender: "male", Phone: "0925555555", DepartmentID: getDept("IT"), Position: "System Administrator", Salary: 25000000, HireDate: ptrTime(2019, 11, 1), DateOfBirth: ptrTime(1986, 1, 1), AvatarURL: "https://ui-avatars.com/api/?name=Dinh+Minh+Tuan&background=059669&color=fff"},
		{Name: "Trịnh Thị Uyên", Gender: "female", Phone: "0926666666", DepartmentID: getDept("HR"), Position: "Training Specialist", Salary: 16500000, HireDate: ptrTime(2023, 4, 22), DateOfBirth: ptrTime(1998, 7, 7), AvatarURL: "https://ui-avatars.com/api/?name=Trinh+Thi+Uyen&background=DC2626&color=fff"},
		{Name: "La Văn Việt", Gender: "male", Phone: "0927777777", DepartmentID: getDept("Sales"), Position: "Business Development Executive", Salary: 21000000, HireDate: ptrTime(2021, 8, 19), DateOfBirth: ptrTime(1993, 2, 2), AvatarURL: "https://ui-avatars.com/api/?name=La+Van+Viet&background=7E22CE&color=fff"},
		{Name: "Nguyễn Ngọc Yến", Gender: "female", Phone: "0928888888", DepartmentID: getDept("Customer Service"), Position: "Customer Success Manager", Salary: 22000000, HireDate: ptrTime(2020, 1, 27), DateOfBirth: ptrTime(1991, 4, 4), AvatarURL: "https://ui-avatars.com/api/?name=Nguyen+Ngoc+Yen&background=0D9488&color=fff"},
		{Name: "Phùng Thanh Bình", Gender: "male", Phone: "0931111111", DepartmentID: getDept("Finance"), Position: "Auditor", Salary: 23000000, HireDate: ptrTime(2022, 10, 6), DateOfBirth: ptrTime(1994, 8, 8), AvatarURL: "https://ui-avatars.com/api/?name=Phung+Thanh+Binh&background=2563EB&color=fff"},
		{Name: "Lưu Thị Châu", Gender: "female", Phone: "0932222222", DepartmentID: getDept("Marketing"), Position: "Brand Manager", Salary: 24000000, HireDate: ptrTime(2019, 9, 13), DateOfBirth: ptrTime(1985, 12, 12), AvatarURL: "https://ui-avatars.com/api/?name=Luu+Thi+Chau&background=C026D3&color=fff"},
		{Name: "Tôn Văn Đạt", Gender: "male", Phone: "0933333333", DepartmentID: getDept("Operations"), Position: "Warehouse Supervisor", Salary: 19000000, HireDate: ptrTime(2021, 12, 2), DateOfBirth: ptrTime(1996, 6, 6), AvatarURL: "https://ui-avatars.com/api/?name=Ton+Van+Dat&background=0369A1&color=fff"},
		{Name: "Kiều Thị Giang", Gender: "female", Phone: "0934444444", DepartmentID: getDept("HR"), Position: "HR Executive", Salary: 16000000, HireDate: ptrTime(2022, 6, 28), DateOfBirth: ptrTime(1997, 10, 10), AvatarURL: "https://ui-avatars.com/api/?name=Kieu+Thi+Giang&background=EA580C&color=fff"},
		{Name: "Châu Minh Hưng", Gender: "male", Phone: "0935555555", DepartmentID: getDept("IT"), Position: "AI Engineer", Salary: 32000000, HireDate: ptrTime(2023, 5, 10), DateOfBirth: ptrTime(1995, 1, 1), AvatarURL: "https://ui-avatars.com/api/?name=Chau+Minh+Hung&background=15803D&color=fff"},
		{Name: "Đoàn Thị Linh", Gender: "female", Phone: "0936666666", DepartmentID: getDept("Sales"), Position: "Key Account Manager", Salary: 27000000, HireDate: ptrTime(2020, 3, 16), DateOfBirth: ptrTime(1988, 3, 3), AvatarURL: "https://ui-avatars.com/api/?name=Doan+Thi+Linh&background=B91C1C&color=fff"},
		{Name: "Quách Văn Minh", Gender: "male", Phone: "0937777777", DepartmentID: getDept("Finance"), Position: "Tax Consultant", Salary: 21000000, HireDate: ptrTime(2021, 1, 5), DateOfBirth: ptrTime(1992, 2, 2), AvatarURL: "https://ui-avatars.com/api/?name=Quach+Van+Minh&background=6D28D9&color=fff"},
		{Name: "Hồ Thị Nga", Gender: "female", Phone: "0938888888", DepartmentID: getDept("Customer Service"), Position: "Call Center Lead", Salary: 17000000, HireDate: ptrTime(2022, 9, 9), DateOfBirth: ptrTime(1999, 11, 11), AvatarURL: "https://ui-avatars.com/api/?name=Ho+Thi+Nga&background=0F766E&color=fff"},
		{Name: "Tăng Quốc Nhật", Gender: "male", Phone: "0941111111", DepartmentID: getDept("IT"), Position: "Security Engineer", Salary: 29000000, HireDate: ptrTime(2020, 7, 7), DateOfBirth: ptrTime(1990, 7, 7), AvatarURL: "https://ui-avatars.com/api/?name=Tang+Quoc+Nhat&background=1E40AF&color=fff"},
		{Name: "Mạc Thị Phương", Gender: "female", Phone: "0942222222", DepartmentID: getDept("Marketing"), Position: "Digital Marketing Executive", Salary: 18000000, HireDate: ptrTime(2023, 2, 14), DateOfBirth: ptrTime(1997, 2, 14), AvatarURL: "https://ui-avatars.com/api/?name=Mac+Thi+Phuong&background=BE123C&color=fff"},
		{Name: "Lâm Văn Quý", Gender: "male", Phone: "0943333333", DepartmentID: getDept("Operations"), Position: "Supply Chain Specialist", Salary: 22000000, HireDate: ptrTime(2021, 11, 23), DateOfBirth: ptrTime(1994, 11, 23), AvatarURL: "https://ui-avatars.com/api/?name=Lam+Van+Quy&background=4338CA&color=fff"},
		{Name: "Đặng Thị Ruby", Gender: "female", Phone: "0944444444", DepartmentID: getDept("HR"), Position: "Compensation Analyst", Salary: 19500000, HireDate: ptrTime(2020, 5, 4), DateOfBirth: ptrTime(1991, 5, 4), AvatarURL: "https://ui-avatars.com/api/?name=Dang+Thi+Ruby&background=0E7490&color=fff"},
		{Name: "Bạch Minh Sang", Gender: "male", Phone: "0945555555", DepartmentID: getDept("Sales"), Position: "Regional Sales Manager", Salary: 30000000, HireDate: ptrTime(2019, 8, 18), DateOfBirth: ptrTime(1986, 8, 18), AvatarURL: "https://ui-avatars.com/api/?name=Bach+Minh+Sang&background=991B1B&color=fff"},
		{Name: "Từ Thị Vy", Gender: "female", Phone: "0946666666", DepartmentID: getDept("Finance"), Position: "Budget Analyst", Salary: 18500000, HireDate: ptrTime(2022, 12, 21), DateOfBirth: ptrTime(1995, 12, 21), AvatarURL: "https://ui-avatars.com/api/?name=Tu+Thi+Vy&background=7C2D12&color=fff"},
	}

	config.DB.Create(&employees)
	fmt.Println("Seed employees successfully")
}

func SeedAdmin() {
	var adminRole models.Role
	config.DB.Where("name = ?", "admin").First(&adminRole)

	ptrTime := func(y, m, d int) *time.Time {
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		return &t
	}

	getDept := func(name string) *uint {
		var dept models.Department
		if err := config.DB.Where("name = ?", name).First(&dept).Error; err != nil {
			return nil
		}
		return &dept.ID
	}

	admins := []struct {
		Employee models.Employee
		Email    string
		Password string
	}{
		{
			Employee: models.Employee{
				Name:         "Cao Chí Quốc",
				Gender:       "male",
				Phone:        "0378286743",
				DepartmentID: getDept("IT"),
				Position:     "Intern",
				Salary:       7000000,
				HireDate:     ptrTime(2026, 5, 7),
				DateOfBirth:  ptrTime(2003, 6, 4),
				AvatarURL:    "https://scontent.fvca5-1.fna.fbcdn.net/v/t39.30808-6/637131180_2129830147835668_1383023946345921909_n.jpg?_nc_cat=109&ccb=1-7&_nc_sid=1d70fc&_nc_ohc=CNaPu6h71MAQ7kNvwF1ii0D&_nc_oc=Adq2kLZ1ePWQE7hjLhpoReHZjTJaehoF0pArLIqEvdlJIsfRFI25AklzfN-BBBaFH60&_nc_zt=23&_nc_ht=scontent.fvca5-1.fna&_nc_gid=4V-Zz1HXckMuUXtNkOY-Nw&_nc_ss=7b2a8&oh=00_Af6q8vuZdftfGfDaGdEH7XOuWJwLXpNdY3xo0yc1QVNh6Q&oe=6A026843",
			},
			Email:    "chiquoc64@admin.company.dev",
			Password: "admin123",
		},
	}

	for _, a := range admins {
		// Tạo Employee nếu chưa có
		var emp models.Employee
		if err := config.DB.Where("phone = ?", a.Employee.Phone).First(&emp).Error; err != nil {
			config.DB.Create(&a.Employee)
			emp = a.Employee
		}

		// Tạo User admin nếu chưa có
		var user models.User
		if err := config.DB.Where("email = ?", a.Email).First(&user).Error; err != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), 10)
			if err != nil {
				continue
			}

			config.DB.Create(&models.User{
				Email:        a.Email,
				PasswordHash: string(hash),
				RoleID:       adminRole.ID,
				EmployeeID:   &emp.ID,
			})
		}
	}

	fmt.Println("Seed admins successfully")
}

// bỏ dấu tiếng Việt
func removeVietnameseTones(str string) string {
	t := norm.NFD.String(str)
	out := make([]rune, 0, len(t))

	for _, r := range t {
		// bỏ dấu
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		switch r {
		case 'Đ':
			r = 'D'
		case 'đ':
			r = 'd'
		}

		out = append(out, r)
	}

	return string(out)
}

// name -> email
func generateEmail(name string) string {
	name = strings.ToLower(name)
	name = removeVietnameseTones(name)

	// loại ký tự lạ
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			return r
		}
		return -1
	}, name)

	name = strings.Join(strings.Fields(name), ".")
	return name + "@company.dev"
}

func SeedUser() {
	var userRole models.Role
	config.DB.Where("name = ?", "user").First(&userRole)

	var employees []models.Employee
	config.DB.Find(&employees)

	defaultPassword := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), 10)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}
	for _, emp := range employees {

		var count int64
		config.DB.Model(&models.User{}).
			Where("employee_id = ?", emp.ID).
			Count(&count)

		if count > 0 {
			continue
		}

		baseEmail := generateEmail(emp.Name)
		email := baseEmail

		prefix := strings.Split(baseEmail, "@")[0]

		config.DB.Model(&models.User{}).
			Where("email LIKE ?", prefix+"%").
			Count(&count)

		if count > 0 {
			email = fmt.Sprintf("%s%d@company.dev", prefix, count+1)
		}

		empID := emp.ID

		config.DB.Create(&models.User{
			Email:        email,
			PasswordHash: string(hash),
			RoleID:       userRole.ID,
			EmployeeID:   &empID,
		})
	}

	fmt.Println("Seed users successfully")
	fmt.Println("Default user password: 123456")
}
