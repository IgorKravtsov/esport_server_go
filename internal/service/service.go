package service

import (
	"github.com/IgorKravtsov/esport_server_go/internal/repository"
	"github.com/IgorKravtsov/esport_server_go/internal/service/user"
	"github.com/IgorKravtsov/esport_server_go/pkg/auth"
	"github.com/IgorKravtsov/esport_server_go/pkg/hash"
	"time"
)

type Services struct {
	//Schools        Schools
	//Students       Students
	//StudentLessons StudentLessons
	//Courses        Courses
	//PromoCodes     PromoCodes
	//Offers         Offers
	//Packages       Packages
	//Modules        Modules
	//Lessons        Lessons
	//Payments       Payments
	//Orders         Orders
	//Admins         Admins
	//Files          Files
	Users user.Users
	//Surveys        Surveys
}

type Deps struct {
	Repos *repository.Repositories
	//Cache                  cache.Cache
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
	//EmailSender            email.Sender
	//EmailConfig            configs.EmailConfig
	//StorageProvider        storage.Provider
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	FondyCallbackURL string
	CacheTTL         int64
	//OtpGenerator           otp.Generator
	VerificationCodeLength int
	Environment            string
	Domain                 string
	//DNS                    dns.DomainManager
}

func NewServices(deps Deps) *Services {
	//schoolsService := NewSchoolsService(deps.Repos.Schools, deps.Cache, deps.CacheTTL)
	//emailsService := NewEmailsService(deps.EmailSender, deps.EmailConfig, *schoolsService, deps.Cache)
	//modulesService := NewModulesService(deps.Repos.Modules, deps.Repos.LessonContent)
	//coursesService := NewCoursesService(deps.Repos.Courses, modulesService)
	//packagesService := NewPackagesService(deps.Repos.Packages, deps.Repos.Modules)
	//offersService := NewOffersService(deps.Repos.Offers, modulesService, packagesService)
	//promoCodesService := NewPromoCodeService(deps.Repos.PromoCodes)
	//lessonsService := NewLessonsService(deps.Repos.Modules, deps.Repos.LessonContent)
	//studentLessonsService := NewStudentLessonsService(deps.Repos.StudentLessons)
	//studentsService := NewStudentsService(deps.Repos.Students, modulesService, offersService, lessonsService, deps.Hasher,
	//  deps.TokenManager, emailsService, studentLessonsService, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength)
	//ordersService := NewOrdersService(deps.Repos.Orders, offersService, promoCodesService, studentsService)
	usersService := user.NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeLength, deps.Domain)

	return &Services{
		//Schools:        schoolsService,
		//Students:       studentsService,
		//StudentLessons: studentLessonsService,
		//Courses:        coursesService,
		//PromoCodes:     promoCodesService,
		//Offers:         offersService,
		//Modules:        modulesService,
		//Payments: NewPaymentsService(ordersService, offersService, studentsService, emailsService, schoolsService,
		//  deps.FondyCallbackURL),
		//Orders: ordersService,
		//Admins: NewAdminsService(deps.Hasher, deps.TokenManager, deps.Repos.Admins, deps.Repos.Schools, deps.Repos.Students,
		//  deps.AccessTokenTTL, deps.RefreshTokenTTL),
		//Packages: packagesService,
		//Lessons:  lessonsService,
		//Files:    NewFilesService(deps.Repos.Files, deps.StorageProvider, deps.Environment),
		Users: usersService,
		//Surveys:  NewSurveysService(deps.Repos.Modules, deps.Repos.SurveyResults, deps.Repos.Students),
	}
}
