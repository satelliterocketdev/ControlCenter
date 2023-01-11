package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"i69social/src/api/config"
	"i69social/src/api/models/fake_user"
	"i69social/src/api/models/profile"
	"i69social/src/api/models/tools"
	"time"
)

func GetProfileService(db *gorm.DB, redisDb *redis.Client) *Profile {
	return &Profile{repo: profile.GetProfileRepository(db, redisDb), fakeUsersRepo:fake_user.GetRepository(db)}
}

type Profile struct {
	repo *profile.ProfileRepository
	fakeUsersRepo *fake_user.Repository
}

func (s *Profile) Create(request profile.SignUpRequest) (p profile.Profile, err error) {
	err = tools.Validator.Struct(request)
	if err != nil {
		err = fmt.Errorf("Sign up request is not valid %w", tools.NewValidationErr(err))
		return
	}
	err = request.Validate()
	if err != nil {
		err = fmt.Errorf("Sign up request is not valid %w", tools.NewValidationErr(err))
		return
	}

	invitationParams, err := s.repo.GetInvitationParams(request.InvitationKey)
	if err == profile.KeyNotExists {
		err = tools.NewValidationErr(fmt.Errorf("Given invitation key `%s` not exists or expired", request.InvitationKey))
		return
	}

	_, err = s.repo.GetByEmail(invitationParams.Email)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if err == nil {
		err = tools.NewValidationErr(fmt.Errorf("User with %s email was already registered", invitationParams.Email))
		return
	}

	profileParams := profile.Profile{
		Email: invitationParams.Email,
		Password: request.Password,
		FirstName: request.FirstName,
		LastName: request.LastName,
		IsActive: true,
	}

	rolesBytes, err := json.Marshal(invitationParams.ToRoles())
	if err != nil {
		return
	}

	profileParams.Roles = profile.CustomJSONMessage{RawMessage: json.RawMessage(rolesBytes)}

	p, err = s.repo.Create(profileParams)
	if err != nil {
		return
	}

	s.repo.RemoveInvitationKey(request.InvitationKey)

	return
}

func (s *Profile) FindUser(request profile.SignInRequest) (p profile.Profile, err error) {

	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}
	profileParams := profile.Profile{Password: request.Password, Email: request.Email}
	p, err = s.repo.GetByUserParams(profileParams)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if err != nil {
		err = tools.NewValidationErr(err)
	}
	return
}

func (s *Profile) DeleteWorker(id uint) (err error) {
	err = s.repo.DeleteById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}

	return
}

func (s *Profile) GenerateInvitation(request profile.GenerateInvitationRequest) (link string, err error) {
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	_, err = s.repo.GetByEmail(request.Email)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if err == nil {
		err = tools.NewValidationErr(fmt.Errorf("User with %s email was already registered", request.Email))
		return
	}

	invitationKey := tools.GenerateRandomString(profile.InvitationKeySize)
	err = s.repo.StoreInvitationParams(invitationKey, request)
	if err != nil {
		return
	}

	link = s.generateInvitationLink(invitationKey)
	return
}

const invitationLinkPattern = "%s/#/signUp?key=%s"
func (s *Profile) generateInvitationLink(key string) (link string) {
	return fmt.Sprintf(invitationLinkPattern, config.GetConfig().Domain, key)
}

func (s *Profile) GetInvitationParams(key string) (params profile.GenerateInvitationRequest, err error) {
	params, err = s.repo.GetInvitationParams(key)
	if err != nil {
		return
	}
	return
}

type Workers struct {
	Workers []profile.Profile `json:"workers"`
	LastContact map[string]string `json:"last_contact"`
	LastMessage map[string]string `json:"last_message"`
	Stats []profile.WorkerStatExported `json:"stats"`
	ActiveChats []profile.WorkerActiveStat `json:"active_chats"`
	Labels []string `json:"labels"`
}

func (s *Profile) GetProfileById(id uint) (err error) {
	_, err = s.repo.GetById(id)
	return
}

func (s *Profile) GetWorkers() (workers Workers, err error) {
	allWorkers, err := s.repo.GetAll()
	if err != nil {
		return
	}

	var workersList = make([]profile.Profile, 0, len(allWorkers))
	for _, worker := range allWorkers {
		if worker.HasRole(profile.ChatterUserPrivilege) {
			worker.Password = ""
			workersList = append(workersList, worker)
		}
	}

	workers.Workers = workersList

	workers.LastContact, err = s.repo.GetWorkersLastItem(profile.UsersLastContact)
	if err != nil {
		return
	}

	workers.LastMessage, err = s.repo.GetWorkersLastItem(profile.UsersLastMessage)
	if err != nil {
		return
	}

	ws, err := s.repo.GetWorkersStats()
	if err != nil {
		return
	}


	for _, item := range ws {
		workers.Stats = append(workers.Stats, item.ToExported())
	}



	workers.ActiveChats, err = s.repo.GetActiveChats()
	if err != nil {
		return
	}

	workers.Labels = GetLabels()



	return
}

func GetLabels() []string {
	n := time.Now().UTC()
	l := make([]string, 0, 7)

	for i := 6; i >= 0; i-- {
		l = append(l, n.Add(-time.Duration(i)*time.Hour*24).Format("2006-01-02"))
	}

	return l
}

func (s *Profile) StoreWorkerLastContact(currentUserId uint) (err error) {
	err = s.repo.StoreWorkerLastContact(profile.UsersLastContact, fmt.Sprintf("%d", currentUserId))
	return
}

func (s *Profile) StoreWorkerLastMessage(currentUserId uint) (err error) {
	err = s.repo.StoreWorkerLastMessage(profile.UsersLastMessage, fmt.Sprintf("%d", currentUserId))
	return
}

func (s *Profile) GetWorkerFakeUsers(workerId uint) (f []fake_user.User, err error) {
	f, err = s.fakeUsersRepo.GetByOwnerId(workerId)
	if err != nil {
		return
	}
	return
}
