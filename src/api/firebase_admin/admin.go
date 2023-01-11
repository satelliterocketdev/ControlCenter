package firebase_admin

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"i69social/src/api/config"
	"i69social/src/api/models/fake_user"
	"i69social/src/api/models/user"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

var FirebaseApp *firebase.App
var ctx = context.Background()

func NewFirebaseApp() (err error) {
	rootDir, _ := os.Getwd()
	filePath := path.Join(rootDir, "src/api/config", "firebase.json")
	// check if file exist
	if _, errFileExist := os.Stat(filePath); os.IsNotExist(errFileExist) {
		err = fmt.Errorf("firebase.json not fount in config directory. Orig err: %+v", errFileExist)
		return
	}

	cfg := config.GetConfig()
	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{
		DatabaseURL: cfg.FirebaseDatabaseURL,
	}
	FirebaseApp, err = firebase.NewApp(ctx, config, opt)
	return
}

func CreateUser(email string) (u *auth.UserRecord, err error) {
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}
	pass := config.GetConfig().NewUserPassword
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(pass).
		Disabled(false)
	u, err = client.CreateUser(ctx, params)
	return
}

func GetUserDetails(uid string) (details map[string]interface{}, err error) {
	details = make(map[string]interface{})
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("users/"+uid).Get(ctx, &details)
	if err == nil && len(details) > 0 {
		return
	}
	return FillUser(uid)
}

type UsersReports map[string]int

func GetUsersReports() (details UsersReports, err error) {
	details = make(UsersReports)
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	var reports map[string]interface{}
	err = client.NewRef("reports/").Get(ctx, &reports)
	if err != nil {
		return
	}
	if len(reports) == 0 {
		return
	}

	details = make(UsersReports, len(reports))
	for k, v := range reports {
		whoReported, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		details[k] = len(whoReported)
	}

	return
}

func GetPresets() (details map[string]interface{}, err error) {
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("defaultPickers/").Get(ctx, &details)
	if err != nil {
		return
	}
	if len(details) == 0 {
		return
	}

	return
}

func RemoveUserReports(uid string) (err error) {

	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("reports/" + uid).Delete(ctx)
	if err != nil {
		return
	}

	return
}

func FillUser(uid string) (details map[string]interface{}, err error) {
	details = make(map[string]interface{})
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}
	details["name"] = ""
	err = client.NewRef("users/"+uid).Set(ctx, &details)
	if err != nil {
		return
	}

	return
}

var UpdateRequestInvalidParamsErr = errors.New("Update request params invalid")
var CreateRequestInvalidParamsErr = errors.New("Create request params invalid")

func UpdateUser(uid string, u user.User, details []byte) (err error) {
	var userDetails map[string]interface{}
	err = json.Unmarshal(details, &userDetails)
	if err != nil {
		err = UpdateRequestInvalidParamsErr
		return
	}
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}

	params := &auth.UserToUpdate{}
	if len(u.DisplayName) != 0 {
		params = params.DisplayName(u.DisplayName)
	}

	if len(u.Email) != 0 {
		params = params.Email(u.Email)
	}

	if len(u.PhoneNumber) != 0 {
		params = params.PhoneNumber(u.PhoneNumber)
	}

	_, err = client.UpdateUser(ctx, uid, params)
	if auth.IsUserNotFound(err) {
		err = UpdateRequestInvalidParamsErr
		return
	}
	if err != nil {
		return
	}

	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}
	err = database.NewRef("/users/"+uid).Update(ctx, userDetails)
	if err != nil {
		return
	}
	return
}

func BlockUser(uid string) (err error) {
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}
	params := (&auth.UserToUpdate{}).Disabled(true)
	_, err = client.UpdateUser(ctx, uid, params)
	if err != nil {
		return
	}
	return
}

func UnBlockUser(uid string) (err error) {
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}
	params := (&auth.UserToUpdate{}).Disabled(false)
	_, err = client.UpdateUser(ctx, uid, params)
	if err != nil {
		return
	}
	return
}

func RemoveUser(uid string) (err error) {
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}
	err = client.DeleteUser(ctx, uid)
	if err != nil {
		return
	}
	return
}

func RemoveUserDetails(uid string) (err error) {
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}
	err = client.NewRef("/users/" + uid).Delete(ctx)
	if err != nil {
		return
	}
	return
}

func GetUser(uid string) (u *auth.UserRecord, err error) {
	client, err := FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}
	u, err = client.GetUser(ctx, uid)
	if err != nil {
		return
	}
	return
}

func GetUsers(limit int, offset int, ordered bool, direction bool, orderField string) (users []user.User, size int, err error) {
	client, err := FirebaseApp.Auth(context.Background())
	if err != nil {
		return
	}
	pager := iterator.NewPager(client.Users(context.Background(), ""), 1000, "")

	users = make([]user.User, 0, 1000)
	reports, err := GetUsersReports()
	if err != nil {
		return
	}
	var nextPageToken = ""
	for {
		var usersBatch []*auth.ExportedUserRecord
		nextPageToken, err = pager.NextPage(&usersBatch)
		if err != nil {
			return
		}

		for _, u := range usersBatch {
			preparedUser := user.User{
				Email:       u.Email,
				UID:         u.UID,
				DisplayName: u.DisplayName,
				Disabled:    u.Disabled,
				Reports:     reports[u.UID],
				PhoneNumber: u.PhoneNumber,
				PhotoURL:    u.PhotoURL,
			}

			if u.UserMetadata != nil {
				preparedUser.SignUp = u.UserMetadata.CreationTimestamp
			}

			users = append(users, preparedUser)
		}
		if nextPageToken == "" {
			break
		}
	}

	size = len(users)
	if !ordered {
		users = returnLimitedList(users, limit, offset)
		return
	}

	if orderField == "" {
		err = fmt.Errorf("ordered field cannot be empty")
		return
	}

	var sortFn func(i, j user.User) bool
	switch orderField {
	case "sign_up":
		sortFn = user.BySignUp
	case "disabled":
		sortFn = user.ByDisabled
	case "reports":
		sortFn = user.ByReports
	}

	var userComparator sort.Interface = user.NewUserComparator(users, sortFn)
	if !direction {
		userComparator = sort.Reverse(userComparator)
	}

	sort.Stable(userComparator)
	users = returnLimitedList(users, limit, offset)
	return
}

func returnLimitedList(users []user.User, limit, offset int) (limitedUsers []user.User) {
	if offset > len(users) {
		return
	}

	if len(users) <= limit && offset == 0 {
		return users
	}

	if len(users) <= offset+limit {
		return users[offset:]
	}

	return users[offset : offset+limit]
}

func GetUsersByFilter(filter string) (users []user.User, err error) {
	client, err := FirebaseApp.Auth(context.Background())
	if err != nil {
		return
	}
	users = make([]user.User, 0, 0)
	pager := iterator.NewPager(client.Users(context.Background(), ""), 1000, "")

	var nextPageToken = ""
	for {
		var usersBatch []*auth.ExportedUserRecord
		nextPageToken, err = pager.NextPage(&usersBatch)
		if err != nil {
			return
		}

		reports, _ := GetUsersReports()
		for _, u := range usersBatch {
			if !strings.Contains(strings.ToLower(u.DisplayName), strings.ToLower(filter)) {
				continue
			}

			preparedUser := user.User{
				Email:       u.Email,
				UID:         u.UID,
				DisplayName: u.DisplayName,
				Disabled:    u.Disabled,
				Reports:     reports[u.UID],
				PhoneNumber: u.PhoneNumber,
				PhotoURL:    u.PhotoURL,
				SignUp:      u.UserMetadata.CreationTimestamp,
			}
			if u.UserMetadata != nil {
				preparedUser.SignUp = u.UserMetadata.CreationTimestamp
			}
			users = append(users, preparedUser)
		}
		if nextPageToken == "" {
			return
		}
	}

}

func ExportUsersEmails(w *csv.Writer) (err error) {
	client, err := FirebaseApp.Auth(context.Background())
	if err != nil {
		return
	}

	pager := iterator.NewPager(client.Users(context.Background(), ""), 1000, "")

	var nextPageToken = ""
	for {
		var usersBatch []*auth.ExportedUserRecord
		nextPageToken, err = pager.NextPage(&usersBatch)
		if err != nil {
			return
		}

		for _, u := range usersBatch {
			if len(u.Email) == 0 {
				continue
			}
			err = w.Write([]string{u.Email})
			if err != nil {
				return
			}
		}
		if nextPageToken == "" {
			return
		}
	}

}

func CreateFakeUser(details []byte) (userDetails map[string]interface{}, err error) {
	err = json.Unmarshal(details, &userDetails)
	if err != nil {
		err = CreateRequestInvalidParamsErr
		return
	}

	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}
	uid := fake_user.GenerateNewId()
	userDetails["id"] = uid
	err = database.NewRef("/fakeUsers/"+uid).Set(ctx, userDetails)
	if err != nil {
		return
	}
	return
}

func GetFakeUserDetails(uid string) (details map[string]interface{}, err error) {
	details = make(map[string]interface{})
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("fakeUsers/"+uid).Get(ctx, &details)
	return
}

func UpdateFakeUser(uid string, details []byte) (userDetails map[string]interface{}, err error) {
	err = json.Unmarshal(details, &userDetails)
	if err != nil {
		return
	}

	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}
	err = database.NewRef("/fakeUsers/"+uid).Update(ctx, userDetails)
	if err != nil {
		return
	}

	return
}

func ChangeFakeUserOnlineStatus(uid string, online bool) (err error) {
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("fakeUsers/"+uid).Update(ctx, map[string]interface{}{"online": online})
	return
}

func AddChatMessage(uid, chatId, messageType, message string) (err error) {
	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	messages, err := GetChatMessages(chatId)
	if err != nil {
		return
	}

	chatParts := strings.Split(chatId, "___")
	if len(chatParts) != 2 {
		err = fmt.Errorf("Chat id should contain two parts")
		return
	}

	var userId = chatParts[0]
	if userId == uid {
		userId = chatParts[1]
	}

	if len(messages) == 0 {
		err = SetFirstMessageAsLike(uid, userId, message)
		if err != nil {
			return
		}
	}
	timestamp := time.Now().UnixNano() / 1e6
	chatMessageDetails := map[string]interface{}{
		"ir":  false, // ir - isRead
		"uid": uid,
		"ts":  timestamp,
	}
	if messageType == "text" {
		chatMessageDetails["txt"] = message
	} else {
		chatMessageDetails["url"] = message
	}
	newKey := fmt.Sprintf("%d", timestamp)
	err = database.NewRef("/messages/"+chatId+"/history/"+newKey).Set(ctx, chatMessageDetails)
	if err != nil {
		return
	}
	return
}

func ReadChatMessages(chatId string, messages []string) (err error) {
	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	for _, m := range messages {
		err = database.NewRef("/messages/"+chatId+"/history/"+m+"/ir").Set(ctx, true)
		if err != nil {
			return
		}
	}

	return
}

func SendChatMessageToFcm(token, name, messageType, message string) (err error) {
	c, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		return
	}

	body := message
	if messageType == "image" {
		body = "New image"
	}
	_, err = c.Send(ctx, &messaging.Message{
		Token:        token,
		Notification: &messaging.Notification{Title: name, Body: body}})
	return
}

func SaveChatNotes(chatId string, notes string) (err error) {
	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = database.NewRef("/messages/"+chatId+"/notes/").Set(ctx, notes)
	if err != nil {
		return
	}

	return
}

func SetMatch(fakeUserId string, newUserId string) (err error) {
	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	chatMessageDetails := map[string]interface{}{
		fakeUserId: map[string]interface{}{
			"option": 0,
			"text":   "_",
			"userId": fakeUserId,
		},
	}
	err = database.NewRef("/users/"+newUserId+"/matches").Update(ctx, chatMessageDetails)
	if err != nil {
		return
	}

	return
}

func SetFirstMessageAsLike(fakeUserId string, newUserId string, message string) (err error) {
	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	chatMessageDetails := map[string]interface{}{
		fakeUserId: message,
	}
	err = database.NewRef("/users/"+newUserId+"/likes").Update(ctx, chatMessageDetails)
	if err != nil {
		return
	}

	return
}

func StartNewChat(fakeUserId string, newUserId string) (err error) {
	var chat string
	if fakeUserId < newUserId {
		chat = fmt.Sprintf("%s___%s", fakeUserId, newUserId)
	} else {
		chat = fmt.Sprintf("%s___%s", newUserId, fakeUserId)
	}

	database, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	chatMessageDetails := map[string]interface{}{
		chat: "history",
	}
	err = database.NewRef("/messages/").Update(ctx, chatMessageDetails)
	if err != nil {
		return
	}

	err = SetMatch(fakeUserId, newUserId)
	log.Printf("Start new chat %s", chat)
	if err != nil {
		log.Printf("Start new chat err %s; chat %s", err, chat)
		return
	}

	return
}

func GetChatMessages(chatId string) (messages map[string]interface{}, err error) {
	messages = make(map[string]interface{})
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	err = client.NewRef("messages/"+chatId+"/history").Get(ctx, &messages)

	return
}