package repository

import (
	"context"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	//Schools        Schools
	//Students       Students
	//StudentLessons StudentLessons
	//Courses        Courses
	//Modules        Modules
	//Packages       Packages
	//LessonContent  LessonContent
	//Offers         Offers
	//PromoCodes     PromoCodes
	//Orders         Orders
	//Admins         Admins
	Users Users
	//Files         Files
	//SurveyResults SurveyResults
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		//Schools:        NewSchoolsRepo(db),
		//Students:       NewStudentsRepo(db),
		//StudentLessons: NewStudentLessonsRepo(db),
		//Courses:        NewCoursesRepo(db),
		//Modules:        NewModulesRepo(db),
		//LessonContent:  NewLessonContentRepo(db),
		//Offers:         NewOffersRepo(db),
		//PromoCodes:     NewPromocodeRepo(db),
		//Orders:         NewOrdersRepo(db),
		//Admins:         NewAdminsRepo(db),
		//Packages:       NewPackagesRepo(db),
		Users: NewUsersRepo(db),
		//Files:          NewFilesRepo(db),
		//SurveyResults:  NewSurveyResultsRepo(db),
	}
}

// Users interface
type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID primitive.ObjectID, code string) error
	//SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	//AttachSchool(ctx context.Context, userID, schoolID primitive.ObjectID) error
}
