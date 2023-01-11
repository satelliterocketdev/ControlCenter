import {getToken} from "./auth";

export class Api {

    static addAuthorizationHeader(headers) {
        var newHeaders = Object.assign({}, headers);
        newHeaders['Authorization'] = `Bearer ${getToken()}`;
        return newHeaders;
    }

    static login(credentials) {
        var request = new Request('/auth/signIn', {
            method: 'POST',
            headers: new Headers(Api.defaultHeaders),
            body: JSON.stringify(credentials)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static signUp(params) {
        var request = new Request('/auth/signUp', {
            method: 'POST',
            headers: new Headers(Api.defaultHeaders),
            body: JSON.stringify(params)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static applyInvitationKey(key) {
        var request = new Request('/internal/invitation/' + key, {
            method: 'GET',
            headers: new Headers(Api.defaultHeaders),
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static createUser(user) {
        var request = new Request('/internal/users/', {
            method: 'POST',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(user)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static blockUser(uid) {
        var request = new Request('/internal/users/'+uid+'/block/', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static unblockUser(uid) {
        var request = new Request('/internal/users/'+uid+'/unblock/', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static updateUser(uid, user) {
        var request = new Request('/internal/users/'+uid, {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(user)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static removeUser(uid) {
        var request = new Request('/internal/users/'+uid, {
            method: 'DELETE',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static removeUserReports(uid) {
        var request = new Request('/internal/users/'+uid+'/reports/', {
            method: 'DELETE',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getUser(uid) {
        var request = new Request('/internal/users/' + uid, {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getUserDetails(uid) {
        var request = new Request('/internal/users/' + uid+'/details', {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getUsers(limit, offset, ordered, orderedField) {

        let path = '/internal/users';
        // limit=25offset=0&ordered=1&field=sign_up
        let fields = [];
        if (typeof limit !== 'undefined' && limit !== 0) {
            fields.push({key: "limit", value: limit});
        }
        if (typeof offset !== 'undefined' && offset !== 0) {
            fields.push({key: "offset", value: offset});
        }
        if (typeof ordered !== 'undefined' && ordered !== 0) {
            fields.push({key:"field", value: orderedField});
            fields.push({key:"ordered", value: ordered});
        }
        if (fields.length !== 0) {
            var query = fields.map((i) => i.key + "=" + i.value).join('&');
            path += "?"+query;
        }
        const request = new Request(path, {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static searchUsers(value) {

        var path = '/internal/users';
        if (value && value.length > 0) {
            path += "?filter="+value
        } else {
            return new Error('filter not defined');
        }
        var request = new Request(path, {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static generateInvitation(params) {
        var request = new Request('/internal/invitation', {
            method: 'POST',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(params)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static createFakeUser(user) {
        var request = new Request('/internal/fakeUsers', {
            method: 'POST',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(user)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static updateFakeUser(uid, user) {
        var request = new Request('/internal/fakeUsers/'+uid+"/details", {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(user)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getFakeUser(uid) {
        var request = new Request('/internal/fakeUsers/'+uid+'/details', {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getFakeUsersPresets() {
        var request = new Request('/internal/presets/fakeUsers', {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getWorkerFakeUsers(id) {
        var request = new Request('/internal/workers/'+id+'/fakeUsers', {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static changeFakeUserOnlineStatus(uid, online) {
        var request = new Request('/internal/fakeUsers/'+uid+'/online', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify({online: online})
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getChatMessages(onPutFn, onPatchFn, onErrorFn) {
        const messages = new EventSource("https://i69social.firebaseio.com/messages.json");
        messages.onerror = e => {
            console.error(e);
            onErrorFn(e);
        };

        messages.addEventListener('put', event => {
            try {
                const data = JSON.parse(event.data);
                onPutFn(data);
            } catch (e) {
                console.error(e);
                onErrorFn(e);
            }
        });

        messages.addEventListener('patch', event => {
            try {
                const data = JSON.parse(event.data);
                onPatchFn(data);
            } catch (e) {
                console.error(e);
                onErrorFn(e);
            }
        });
    }

    static getChatUsersDetails(onDataFn, onErrorFn) {
        const users = new EventSource("https://i69social.firebaseio.com/users.json");
        users.onerror = e => {
            console.error(e);
            onErrorFn(e);
        };

        users.addEventListener('put', event => {
            try {
                const data = JSON.parse(event.data);
                onDataFn(data);
            } catch (e) {
                console.error(e);
                onErrorFn(e);
            }
        });

    }

    static getFakeUsersDetails() {
        var request = new Request("https://i69social.firebaseio.com/fakeUsers.json", {
            method: 'GET'
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static sendChatMessage(chatId, fromName, fromUserId, recipientToken, type, message) {
        const chatMessage = {
            from_name: fromName,
            from_user_id: fromUserId,
            recipient_token: recipientToken,
            type: type,
            message: message,
        };
        const request = new Request('/internal/chats/'+chatId+'/send', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(chatMessage)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static saveChatNotes(chatId, notes) {
        const message = {
            notes: notes
        };
        const request = new Request('/internal/chats/'+chatId+'/notes', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(message)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static startNewChat(fakeUserId, realUserId) {
        const message = {
            fake_user_id: fakeUserId,
            real_user_id: realUserId
        };
        const request = new Request('/internal/chats/new', {
            method: 'POST',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(message)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static readChatMessages(uid, chatId, messages) {
        const chatMessage = {
            user_id: uid,
            messages: messages,
        };
        const request = new Request('/internal/chats/'+chatId+'/read', {
            method: 'PUT',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders),
            body: JSON.stringify(chatMessage)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static getWorkers() {
        var request = new Request('/internal/workers', {
            method: 'GET',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }

    static deleteWorker(id) {
        var request = new Request('/internal/workers/'+id, {
            method: 'DELETE',
            headers: Api.addAuthorizationHeader(Api.defaultHeaders)
        });

        return fetch(request).then(response => {
            return response;
        }).catch(error => {
            return error;
        });
    }
}

Api.defaultHeaders = {
    'Content-Type': 'application/json'
};
