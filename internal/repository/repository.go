package repository

import (
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
	User User
	Gym  Gym
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
		User: NewUserRepo(db),
		Gym:  NewGymRepo(db),
		//Files:          NewFilesRepo(db),
		//SurveyResults:  NewSurveyResultsRepo(db),
	}
}
